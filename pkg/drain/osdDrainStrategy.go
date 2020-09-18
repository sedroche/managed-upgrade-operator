package drain

import (
	"fmt"
	"sort"
	"time"

	"github.com/hashicorp/go-multierror"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	upgradev1alpha1 "github.com/openshift/managed-upgrade-operator/pkg/apis/upgrade/v1alpha1"
)

func NewOSDDrainStrategy(c client.Client, uc *upgradev1alpha1.UpgradeConfig, node *corev1.Node, cfg *NodeDrain, ts []TimedDrainStrategy) (NodeDrainStrategy, error) {
	return &osdDrainStrategy{
		c,
		node,
		cfg,
		ts,
	}, nil
}

type osdDrainStrategy struct {
	client               client.Client
	node                 *corev1.Node
	cfg                  *NodeDrain
	timedDrainStrategies []TimedDrainStrategy
}

func (ds *osdDrainStrategy) Execute(startTime *metav1.Time) ([]*DrainStrategyResult, error) {
	me := &multierror.Error{}
	res := []*DrainStrategyResult{}
	for _, ds := range ds.timedDrainStrategies {
		if isAfter(startTime, ds.GetWaitDuration()) {
			r, err := ds.GetStrategy().Execute()
			me = multierror.Append(err, me)
			if r.HasExecuted {
				res = append(res, &DrainStrategyResult{Message: fmt.Sprintf("Drain strategy %s has been executed. %s", ds.GetDescription(), r.Message)})
			}
		}
	}

	return res, me.ErrorOrNil()
}

func (ds *osdDrainStrategy) HasFailed(startTime *metav1.Time) (bool, error) {
	if len(ds.timedDrainStrategies) == 0 {
		return isAfter(startTime, ds.cfg.GetTimeOutDuration()), nil
	}

	maxWaitStrategy := maxWaitDuration(ds.timedDrainStrategies)
	if !isAfter(startTime, maxWaitStrategy.GetWaitDuration()) {
		return false, nil
	}
	failed, err := maxWaitStrategy.GetStrategy().HasFailed()
	if err != nil {
		return false, err
	}

	return failed && isAfter(startTime, maxWaitStrategy.GetWaitDuration()+ds.cfg.GetExpectedDrainDuration()), nil
}

type timedStrategy struct {
	name         string
	description  string
	waitDuration time.Duration
	strategy     DrainStrategy
}

func (ts *timedStrategy) GetWaitDuration() time.Duration {
	return ts.waitDuration
}

func (ts *timedStrategy) GetName() string {
	return ts.name
}

func (ts *timedStrategy) GetDescription() string {
	return ts.description
}

func (ts *timedStrategy) GetStrategy() DrainStrategy {
	return ts.strategy
}

func isAfter(t *metav1.Time, d time.Duration) bool {
	return t != nil && t.Add(d).Before(metav1.Now().Time)
}

func maxWaitDuration(ts []TimedDrainStrategy) TimedDrainStrategy {
	sort.Slice(ts, func(i, j int) bool {
		iWait := ts[i].GetWaitDuration()
		jWait := ts[j].GetWaitDuration()
		return iWait < jWait
	})
	return ts[len(ts)-1]
}

func getOsdTimedStrategies(c client.Client, uc *upgradev1alpha1.UpgradeConfig, node *corev1.Node, cfg *NodeDrain) ([]TimedDrainStrategy, error) {
	return []TimedDrainStrategy{}, nil
}