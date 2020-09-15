package drain

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	upgradev1alpha1 "github.com/openshift/managed-upgrade-operator/pkg/apis/upgrade/v1alpha1"
)

//go:generate mockgen -destination=mocks/drainStrategyBuilder.go -package=mocks github.com/openshift/managed-upgrade-operator/pkg/drain DrainStrategyBuilder
type DrainStrategyBuilder interface {
	NewDrainStrategy(c client.Client, uc *upgradev1alpha1.UpgradeConfig, node *corev1.Node, cfg *NodeDrain) (DrainStrategy, error)
}

//go:generate mockgen -destination=mocks/drainStrategy.go -package=mocks github.com/openshift/managed-upgrade-operator/pkg/drain DrainStrategy
type DrainStrategy interface {
	Execute(*metav1.Time) ([]*DrainStrategyResult, error)
	HasFailed(*metav1.Time) (bool, error)
}

//go:generate mockgen -destination=./timeBasedDrainStrategy.go -package=drain github.com/openshift/managed-upgrade-operator/pkg/drain TimeBasedDrainStrategy
type TimeBasedDrainStrategy interface {
	GetWaitDuration() time.Duration
	Execute() (*DrainStrategyResult, error)
	GetName() string
	GetDescription() string
	HasFailed(*metav1.Time) (bool, error)
}

func NewBuilder() DrainStrategyBuilder {
	return &drainStrategyBuilder{}
}

type drainStrategyBuilder struct{}

func (dsb *drainStrategyBuilder) NewDrainStrategy(c client.Client, uc *upgradev1alpha1.UpgradeConfig, node *corev1.Node, cfg *NodeDrain) (DrainStrategy, error) {
	// TODO: move pdb budget to filters
	pdbList := &policyv1beta1.PodDisruptionBudgetList{}
	err := c.List(context.TODO(), pdbList)
	if err != nil {
		return nil, err
	}

	return NewOSDDrainStrategy(c, uc, node, cfg, []TimeBasedDrainStrategy{
		&timedPodDeleteStrategy{
			client:       c,
			name:         defaultPodStrategyName,
			description:  "Default pod deletion",
			waitDuration: cfg.GetTimeOutDuration(),
			filters:      []podPredicate{isOnNode(node), isNotDaemonSet, isNotPdbPod(pdbList)},
		},
		&timedPodDeleteStrategy{
			client:       c,
			name:         pdbStrategyName,
			description:  "PDB pod deletion",
			waitDuration: uc.GetPDBDrainTimeoutDuration(),
			filters:      []podPredicate{isOnNode(node), isNotDaemonSet, isPdbPod(pdbList)},
		},
	})
}

type DrainStrategyResult struct {
	Message string
}
