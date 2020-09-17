package drain

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/managed-upgrade-operator/pkg/pod"
)

type forceDeletePodStrategy struct {
	client  client.Client
	filters []pod.PodPredicate
}

func (fdps *forceDeletePodStrategy) Execute() (*DrainStrategyResult, error) {
	allPods := &corev1.PodList{}
	err := fdps.client.List(context.TODO(), allPods)
	if err != nil {
		return &DrainStrategyResult{Message: "Error listing pods for deletion"}, err
	}

	podsToDelete := pod.FilterPods(allPods, fdps.filters...)
	res, err := pod.DeletePods(fdps.client, podsToDelete)
	if err != nil {
		return &DrainStrategyResult{Message: res.Message}, err
	}

	return &DrainStrategyResult{Message: res.Message}, nil
}

func (fdps *forceDeletePodStrategy) HasFailed() (bool, error) {
	allPods := &corev1.PodList{}
	err := fdps.client.List(context.TODO(), allPods)
	if err != nil {
		return false, err
	}

	filterPods := pod.FilterPods(allPods, fdps.filters...)
	if len(filterPods.Items) == 0 {
		return false, nil
	}

	return true, nil
}
