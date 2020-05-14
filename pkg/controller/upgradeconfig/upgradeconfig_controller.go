package upgradeconfig

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"time"

	configv1 "github.com/openshift/api/config/v1"
	upgradev1alpha1 "github.com/openshift/managed-upgrade-operator/pkg/apis/upgrade/v1alpha1"
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
	return &ReconcileUpgradeConfig{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

func NewReconcileUpgradeConfig(client client.Client, scheme *runtime.Scheme) (reconcile.Reconciler, error) {
	if scheme == nil {
		return nil, fmt.Errorf("scheme cannot be nil")
	}

	return &ReconcileUpgradeConfig{client, scheme}, nil
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("upgradeconfig-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource UpgradeConfig
	err = c.Watch(&source.Kind{Type: &upgradev1alpha1.UpgradeConfig{}}, &handler.EnqueueRequestForObject{})
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
	client client.Client
	scheme *runtime.Scheme
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

	cv := &configv1.ClusterVersion{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "version"}, cv)
	if err != nil {
		return reconcile.Result{}, err
	}

	for _, update := range cv.Status.AvailableUpdates {
		if update.Force == false && update.Version == instance.Spec.Desired.Version {
			cv.Spec.Overrides = nil
			cv.Spec.Channel = instance.Spec.Desired.Channel
			cv.Spec.DesiredUpdate = &configv1.Update{
				Version: instance.Spec.Desired.Version,
			}
			err := r.client.Update(context.TODO(), cv)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
	}
	return reconcile.Result{
		RequeueAfter: time.Hour,
	}, nil
}
