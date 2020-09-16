package drain

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/hashicorp/go-multierror"
	corev1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
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
		cfg,
		ts,
	}, nil
}

type osdDrainStrategy struct {
	client               client.Client
	node                 *corev1.Node
	cfg                  *NodeDrain
	timedDrainStrategies []TimeBasedDrainStrategy
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
		return isAfter(startTime, ds.cfg.GetTimeOutDuration()), nil
	}

	maxWaitStrategy := maxWaitDuration(ds.timedDrainStrategies)
	failed, err := maxWaitStrategy.HasFailed(startTime)
	if err != nil {
		return false, err
	}

	return failed && isAfter(startTime, maxWaitStrategy.GetWaitDuration()+ds.cfg.GetExpectedDrainDuration()), nil
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
	res, err := Delete(tpds.client, podsToDelete)
	if err != nil {
		return &DrainStrategyResult{Message: res.Message}, err
	}

	return &DrainStrategyResult{Message: res.Message}, nil
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

func maxWaitDuration(ts []TimeBasedDrainStrategy) TimeBasedDrainStrategy {
	sort.Slice(ts, func(i, j int) bool {
		iWait := ts[i].GetWaitDuration()
		jWait := ts[j].GetWaitDuration()
		return iWait < jWait
	})
	return ts[len(ts)-1]
}

func getOsdTimedStrategies(c client.Client, uc *upgradev1alpha1.UpgradeConfig, node *corev1.Node, cfg *NodeDrain) ([]TimeBasedDrainStrategy, error) {
	pdbList := &policyv1beta1.PodDisruptionBudgetList{}
	err := c.List(context.TODO(), pdbList)
	if err != nil {
		return nil, err
	}

	allPods := &corev1.PodList{}
	err = c.List(context.TODO(), allPods)
	if err != nil {
		return nil, err
	}

	pdbPodsOnNode := Filter(allPods, isOnNode(node), isNotDaemonSet, isPdbPod(pdbList))
	hasPdbPod := len(pdbPodsOnNode.Items) > 0
	if hasPdbPod {
		return []TimeBasedDrainStrategy{
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
		}, nil
	} else {
		return []TimeBasedDrainStrategy{
			&timedPodDeleteStrategy{
				client:       c,
				name:         defaultPodStrategyName,
				description:  "Default pod deletion",
				waitDuration: cfg.GetTimeOutDuration(),
				filters:      []podPredicate{isOnNode(node), isNotDaemonSet, isNotPdbPod(pdbList)},
			},
		}, nil
	}
}
