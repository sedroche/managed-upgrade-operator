package machineconfigpool

import (
	"context"
	"github.com/openshift/managed-upgrade-operator/pkg/metrics"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"time"

	machineconfigapi "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	upgradev1alpha1 "github.com/openshift/managed-upgrade-operator/pkg/apis/upgrade/v1alpha1"
)

var log = logf.Log.WithName("controller_machineconfigpool")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new MachineConfigPool Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMachineConfigPool{
		client:               mgr.GetClient(),
		scheme:               mgr.GetScheme(),
		metricsClientBuilder: metrics.NewBuilder(),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("machineconfigpool-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MachineConfigPool
	err = c.Watch(&source.Kind{Type: &machineconfigapi.MachineConfigPool{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileMachineConfigPool implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileMachineConfigPool{}

// ReconcileMachineConfigPool reconciles a MachineConfigPool object
type ReconcileMachineConfigPool struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client               client.Client
	scheme               *runtime.Scheme
	metricsClientBuilder metrics.MetricsBuilder
}

func (r *ReconcileMachineConfigPool) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling MachineConfigPool")

	// Fetch the MachineConfigPool instance
	instance := &machineconfigapi.MachineConfigPool{}
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

	if instance.Name != "worker" {
		return reconcile.Result{}, nil
	}

	uc, err := getUpgradeConfig(r.client)
	if err != nil {
		return reconcile.Result{}, err
	}

	history := uc.Status.History.GetHistory(uc.Spec.Desired.Version)
	if history != nil && history.Phase == upgradev1alpha1.UpgradePhaseUpgrading && instance.Status.MachineCount == instance.Status.UpdatedMachineCount {
		metricsClient, err := r.metricsClientBuilder.NewClient(r.client)
		if err != nil {
			return reconcile.Result{}, err
		}
		isSet, err := metricsClient.IsMetricNodeUpgradeEndTimeSet(uc.Name, uc.Spec.Desired.Version)
		if err != nil {
			return reconcile.Result{}, err
		}
		if !isSet {
			metricsClient.UpdateMetricNodeUpgradeEndTime(time.Now(), uc.Name, uc.Spec.Desired.Version)
		}
	}

	return reconcile.Result{}, nil
}

func getUpgradeConfig(c client.Client) (*upgradev1alpha1.UpgradeConfig, error) {
	uCList := &upgradev1alpha1.UpgradeConfigList{}

	err := c.List(context.TODO(), uCList)
	if err != nil {
		return nil, err
	}

	for _, uC := range uCList.Items {
		return &uC, nil
	}

	return nil, errors.NewNotFound(schema.GroupResource{Group: upgradev1alpha1.SchemeGroupVersion.Group, Resource: "UpgradeConfig"}, "UpgradeConfig")
}
