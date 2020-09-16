package drain

import (
	"time"

	corev1 "k8s.io/api/core/v1"
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
	ts, err := getOsdTimedStrategies(c, uc, node, cfg)
	if err != nil {
		return nil, err
	}
	return NewOSDDrainStrategy(c, uc, node, cfg, ts)
}

type DrainStrategyResult struct {
	Message string
}
