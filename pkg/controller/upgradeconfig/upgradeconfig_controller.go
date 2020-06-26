package upgradeconfig

import (
	"context"
	upgradev1alpha1 "github.com/openshift/managed-upgrade-operator/pkg/apis/upgrade/v1alpha1"
	"github.com/openshift/managed-upgrade-operator/pkg/cluster_upgrader"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
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

	// Once we have a history, skip above. We are Upgrading now.
	// Upgrade history
	// Manage History status

	cv, err := cluster_upgrader.GetClusterVersion(r.client);
	if err != nil {
		return reconcile.Result{}, err
	}

	desiredVersion := instance.Spec.Desired.Version
	if cv.Spec.DesiredUpdate.Version != desiredVersion && (instance.Status.History == nil && !instance.Status.History.HasHistory(desiredVersion)) {
		reqLogger.Info("No Upgrade required.")
		return reconcile.Result{}, nil
	}

	//var history upgradev1alpha1.UpgradeHistory
	//found := false
	//for _, h := range instance.Status.History {
	//	if h.Version == instance.Spec.Desired.Version {
	//		history = h
	//		found = true
	//	}
	//}
	//
	//if !found {
	//	history = upgradev1alpha1.UpgradeHistory{Version: instance.Spec.Desired.Version}
	//	history.Conditions = upgradev1alpha1.NewConditions()
	//	instance.Status.History = append([]upgradev1alpha1.UpgradeHistory{history}, instance.Status.History...)
	//
	//	reqLogger.Info("Checking whether it's time to do an upgrade")
	//	ready := cluster_upgrader.IsTimeToUpgrade(instance)
	//	if !ready {
	//		history.Phase = upgradev1alpha1.UpgradePhasePending
	//		err := r.client.Status().Update(context.TODO(), instance)
	//		if err != nil {
	//			reqLogger.Info("Failed to set pending status for upgrade!")
	//			return reconcile.Result{}, err
	//		}
	//		return reconcile.Result{}, err
	//	} else {
	//		history.Phase = upgradev1alpha1.UpgradePhaseUpgrading
	//	}
	//}
	//
	//reqLogger.Info("Current cluster status", "status", history.Phase)
	//upgrader, err := r.clusterUpgraderBuilder.NewClient(r.client)
	//if err != nil {
	//	return reconcile.Result{}, err
	//}
	//
	//reqLogger.Info("Reconciling upgrade ", "time", time.Now())
	//err = upgrader.UpgradeCluster(instance, reqLogger)
	//if err != nil {
	//	reqLogger.Error(err, "Failed to upgrade cluster")
	//	return reconcile.Result{}, err
	//}

	return reconcile.Result{}, nil
}

func upgradeRequested() {

}