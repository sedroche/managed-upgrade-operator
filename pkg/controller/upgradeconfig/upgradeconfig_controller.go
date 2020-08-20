package upgradeconfig

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	upgradev1alpha1 "github.com/openshift/managed-upgrade-operator/pkg/apis/upgrade/v1alpha1"
	cub "github.com/openshift/managed-upgrade-operator/pkg/cluster_upgrader_builder"
	"github.com/openshift/managed-upgrade-operator/pkg/configmanager"
	"github.com/openshift/managed-upgrade-operator/pkg/metrics"
	"github.com/openshift/managed-upgrade-operator/pkg/scheduler"
	"github.com/openshift/managed-upgrade-operator/pkg/validation"

	"github.com/openshift-online/ocm-sdk-go"
	configv1 "github.com/openshift/api/config/v1"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

var (
	log = logf.Log.WithName("controller_upgradeconfig")
)

// Add creates a new UpgradeConfig Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileUpgradeConfig{
		client:                 mgr.GetClient(),
		scheme:                 mgr.GetScheme(),
		metricsClientBuilder:   metrics.NewBuilder(),
		clusterUpgraderBuilder: cub.NewBuilder(),
		validationBuilder:      validation.NewBuilder(),
		configManagerBuilder:   configmanager.NewBuilder(),
		scheduler:              scheduler.NewScheduler(),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("upgradeconfig-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource UpgradeConfig, status change will not trigger a reconcile
	err = c.Watch(&source.Kind{Type: &upgradev1alpha1.UpgradeConfig{}}, &handler.EnqueueRequestForObject{}, StatusChangedPredicate{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileUpgradeConfig implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileUpgradeConfig{}

// ReconcileUpgradeConfig reconciles a UpgradeConfig object
type ReconcileUpgradeConfig struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client                 client.Client
	scheme                 *runtime.Scheme
	metricsClientBuilder   metrics.MetricsBuilder
	clusterUpgraderBuilder cub.ClusterUpgraderBuilder
	validationBuilder      validation.ValidationBuilder
	configManagerBuilder   configmanager.ConfigManagerBuilder
	scheduler              scheduler.Scheduler
}

// Reconcile reads that state of the cluster for a UpgradeConfig object and makes changes based on the state read
// and what is in the UpgradeConfig.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileUpgradeConfig) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling UpgradeConfig")

	// Get the token
	token := os.Getenv("OCM_TOKEN")

	cv := &configv1.ClusterVersion{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: "version"} , cv)
	if err != nil {
		return reconcile.Result{}, nil
	}

	// Create the connection, and remember to close it:
	connection, err := sdk.NewConnectionBuilder().
		URL("https://api.stage.openshift.com").
		Tokens(token).
		Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't build connection: %v\n", err)
		os.Exit(1)
	}
	defer connection.Close()

	cluster, err := GetCluster(connection.ClustersMgmt().V1().Clusters(), string(cv.Spec.ClusterID))
	uc, err := getUpgradeConfig(cluster)

	err = r.client.Update(context.TODO(), uc)
	if err != nil {
		return reconcile.Result{}, err
	}

	fmt.Printf("%v", cluster)
	return reconcile.Result{}, nil
}

func getUpgradeConfig(*cmv1.Cluster) (*upgradev1alpha1.UpgradeConfig, error) {
	return &upgradev1alpha1.UpgradeConfig{
		Spec: upgradev1alpha1.UpgradeConfigSpec{
			Type: upgradev1alpha1.OSD,
			UpgradeAt: "2020-01-01T00:00:00Z",
			PDBForceDrainTimeout: 60,
			Desired: upgradev1alpha1.Update{
				Channel: "fast-4.4",
				Version: "4.3.19",
			},
		},
	}, nil
}

// TODO: SDK does not work on external ids
// connection.ClustersMgmt().V1().Clusters().Cluster(id)
func GetCluster(client *cmv1.ClustersClient, clusterKey string) (*cmv1.Cluster, error) {
	query := fmt.Sprintf(
		"(id = '%s' or name = '%s' or external_id = '%s')",
		clusterKey, clusterKey,
	)
	response, err := client.List().
		Search(query).
		Page(1).
		Size(1).
		Send()
	if err != nil {
		return nil, fmt.Errorf("Failed to locate cluster '%s': %v", clusterKey, err)
	}

	switch response.Total() {
	case 0:
		return nil, fmt.Errorf("There is no cluster with identifier or name '%s'", clusterKey)
	case 1:
		return response.Items().Slice()[0], nil
	default:
		return nil, fmt.Errorf("There are %d clusters with identifier or name '%s'", response.Total(), clusterKey)
	}
}
