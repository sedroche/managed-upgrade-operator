package upgradeconfig

import (
	"context"
	configv1 "github.com/openshift/api/config/v1"
	rm "github.com/openshift/cluster-version-operator/lib/resourcemerge"
	upgradev1alpha1 "github.com/openshift/managed-upgrade-operator/pkg/apis/upgrade/v1alpha1"
	"github.com/openshift/managed-upgrade-operator/pkg/cluster_upgrader"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"
)

var log = logf.Log.WithName("controller_upgradeconfig")

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
		clusterUpgraderBuilder: cluster_upgrader.NewBuilder(),
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

	// TODO: Do I need to watch the ClusterVersion? To speed up status updates/is ready checks?
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
	clusterUpgraderBuilder cluster_upgrader.ClusterUpgraderBuilder
}

// Reconcile reads that state of the cluster for a UpgradeConfig object and makes changes based on the state read
// and what is in the UpgradeConfig.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileUpgradeConfig) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling UpgradeConfig")

	// Fetch the UpgradeConfig instance
	instance := &upgradev1alpha1.UpgradeConfig{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	//Repeat based on whether we have a history for this version
	// Is there a new upgrade version specified
	// Is it valid
	// Is it time to Upgrade

	// Once we have a history with not pending, skip above. We are Upgrading now.
	// Upgrade history
	// Manage History status

	cv, err := getClusterVersion(r.client)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Advantages:
	// Upgrade is handled by the reconcile controller not the clusterupgrade logic. clusterupgrade logic should be independent of our upgradeconfig logic of when and if valid
	desiredVersion := instance.Spec.Desired.Version
	if !isNewDesiredVersion(cv, desiredVersion) {
		reqLogger.Info("No Upgrade required.")
		return reconcile.Result{}, nil
	}

	reqLogger.Info("Validating UpgradeConfig")
	err = cluster_upgrader.IsValidUpgradeConfig(instance, cv)
	if err != nil {
		// TODO: We need to Alert here
		_, ok := err.(*cluster_upgrader.ValidationError)
		if ok {
			// TODO: Is this double logging?
			reqLogger.Info(err.Error())
		}
		return reconcile.Result{}, err
	}

	var history upgradev1alpha1.UpgradeHistory
	if instance.Status.History == nil {
		history = upgradev1alpha1.UpgradeHistory{Version: instance.Spec.Desired.Version, Phase: upgradev1alpha1.UpgradePhasePending}
	} else {
		history = *instance.Status.History.GetHistory(desiredVersion)
	}

	if history.Phase != upgradev1alpha1.UpgradePhasePending {

	}

	isTime := cluster_upgrader.IsTimeToUpgrade(instance)
	if history.StartTime == nil && !isTime {
		history.Conditions = upgradev1alpha1.NewConditions()
		instance.Status.History = append([]upgradev1alpha1.UpgradeHistory{history}, instance.Status.History...)
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	// We should now start upgrading. Set phase and record start time.
	if history.StartTime == nil {
		history.Phase = upgradev1alpha1.UpgradePhaseUpgrading
		history.StartTime = &metav1.Time{Time: time.Now()}

		instance.Status.History.SetHistory(history)
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	status := history.Phase
	reqLogger.Info("Current cluster status", "status", status)

	upgrader, err := r.clusterUpgraderBuilder.NewClient(r.client)
	if err != nil {
		return reconcile.Result{}, err
	}
	reqLogger.Info("Reconciling UpgradeConfig")
	upgradeconfig, err := upgrader.UpgradeCluster(instance, reqLogger)
	if err != nil {
		reqLogger.Error(err, "Failed to upgrade cluster")
	}

	h := upgradeconfig.Status.History.GetHistory(desiredVersion)
	if  h.Phase == upgradev1alpha1.UpgradePhaseUpgraded {
		history.CompleteTime = &metav1.Time{Time: time.Now()}
	}

	upgradeconfig.Status.History.SetHistory(*h)
	err = r.client.Status().Update(context.TODO(), upgradeconfig)
	if err != nil {
		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil
}

func getClusterVersion(c client.Client) (*configv1.ClusterVersion, error) {
	cv := &configv1.ClusterVersion{}
	err := c.Get(context.TODO(), types.NamespacedName{Name: "version"}, cv)
	if err != nil {
		return nil, err
	}

	return cv, nil
}

func isNewDesiredVersion(cv *configv1.ClusterVersion, desiredVersion string) bool {
	return cv.Spec.DesiredUpdate.Version != desiredVersion && !rm.IsOperatorStatusConditionTrue(cv.Status.Conditions, configv1.OperatorProgressing)
}
