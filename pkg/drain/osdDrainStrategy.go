package drain

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/hashicorp/go-multierror"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	upgradev1alpha1 "github.com/openshift/managed-upgrade-operator/pkg/apis/upgrade/v1alpha1"
)

var (
	pdbStrategyName        = "PDB"
	defaultPodStrategyName = "Default"
)

func NewOSDDrainStrategy(c client.Client, uc *upgradev1alpha1.UpgradeConfig, node *corev1.Node, cfg *NodeDrain, ts []TimeBasedDrainStrategy) (DrainStrategy, error) {
	return &osdDrainStrategy{
		c,
		node,
		cfg.GetExpectedDrainDuration(),
		ts,
	}, nil
}

type osdDrainStrategy struct {
	client                client.Client
	node                  *corev1.Node
	expectedDrainDuration time.Duration
	timedDrainStrategies  []TimeBasedDrainStrategy
}

func (ds *osdDrainStrategy) Execute(startTime *metav1.Time) ([]*DrainStrategyResult, error) {
	me := &multierror.Error{}
	res := []*DrainStrategyResult{}
	for _, ds := range ds.timedDrainStrategies {
		if isAfter(startTime, ds.GetWaitDuration()) {
			r, err := ds.Execute()
			me = multierror.Append(err, me)
			res = append(res, &DrainStrategyResult{Message: fmt.Sprintf("Drain strategy %s has been executed. %s", ds.GetDescription(), r.Message)})
		}
	}

	return res, me.ErrorOrNil()
}

func (ds *osdDrainStrategy) HasFailed(startTime *metav1.Time) (bool, error) {
	if len(ds.timedDrainStrategies) == 0 {
		return isAfter(startTime, ds.expectedDrainDuration), nil
	}

	// TODO: check what happens when it errors in here during multiple strategies
	timedFails := []bool{}
	for _, ds := range ds.timedDrainStrategies {
		failed, err := ds.HasFailed(startTime)
		if err != nil {
			return false, err
		}
		if failed {
			timedFails = append(timedFails, failed)
		}
	}

	if len(timedFails) < len(ds.timedDrainStrategies) {
		return false, nil
	}

	return isAfter(startTime, maxWaitDuration(ds.timedDrainStrategies)+ds.expectedDrainDuration), nil
}

type timedPodDeleteStrategy struct {
	client       client.Client
	name         string
	description  string
	waitDuration time.Duration
	filters      []podPredicate
}

func (tpds *timedPodDeleteStrategy) GetWaitDuration() time.Duration {
	return tpds.waitDuration
}

func (tpds *timedPodDeleteStrategy) Execute() (*DrainStrategyResult, error) {
	allPods := &corev1.PodList{}
	err := tpds.client.List(context.TODO(), allPods)
	if err != nil {
		return &DrainStrategyResult{Message: "Error listing pods for deletion"}, err
	}

	podsToDelete := Filter(allPods, tpds.filters...)
	return Delete(tpds.client, podsToDelete)
}

func (tpds *timedPodDeleteStrategy) GetName() string {
	return tpds.name
}

func (tpds *timedPodDeleteStrategy) GetDescription() string {
	return tpds.description
}

func (tpds *timedPodDeleteStrategy) HasFailed(startTime *metav1.Time) (bool, error) {
	if !isAfter(startTime, tpds.GetWaitDuration()) {
		return false, nil
	}

	allPods := &corev1.PodList{}
	err := tpds.client.List(context.TODO(), allPods)
	if err != nil {
		return false, err
	}

	filterPods := Filter(allPods, tpds.filters...)
	if len(filterPods.Items) == 0 {
		return false, nil
	}

	return true, nil
}

func isAfter(t *metav1.Time, d time.Duration) bool {
	return t != nil && t.Add(d).Before(metav1.Now().Time)
}

func maxWaitDuration(ts []TimeBasedDrainStrategy) time.Duration {
	sort.Slice(ts, func(i, j int) bool {
		return ts[i].GetWaitDuration() > ts[j].GetWaitDuration()
	})
	return ts[len(ts)-1].GetWaitDuration()
}
