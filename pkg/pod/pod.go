package pod

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PodPredicate func(corev1.Pod) bool

func FilterPods(podList *corev1.PodList, predicates ...PodPredicate) *corev1.PodList {
	filteredPods := &corev1.PodList{}
	for _, pod := range podList.Items {
		var match = true
		for _, p := range predicates {
			if !p(pod) {
				match = false
				break
			}
		}
		if match {
			filteredPods.Items = append(filteredPods.Items, pod)
		}
	}

	return filteredPods
}

type DeleteResult struct {
	Message string
}

func DeletePods(c client.Client, pl *corev1.PodList) (*DeleteResult, error) {
	me := &multierror.Error{}
	var podsMarkedForDeletion []string
	for _, p := range pl.Items {
		if len(p.ObjectMeta.GetFinalizers()) != 0 {
			err := removeFinalizers(c, &p)
			if err != nil {
				return &DeleteResult{Message: fmt.Sprintf("Error removing finalizer for pod %s", p.Name)}, err
			}
		}

		if p.DeletionTimestamp == nil {
			err := c.Delete(context.TODO(), &p)
			if err != nil {
				me = multierror.Append(err, me)
			} else {
				podsMarkedForDeletion = append(podsMarkedForDeletion, p.Name)
			}
		}
	}

	// TODO: Log removing finalizers. and log no pods deleted better or don't log
	// TODO: need to test both types on the one node
	return &DeleteResult{Message: fmt.Sprintf("Pod(s) %s have been marked for deletion", strings.Join(podsMarkedForDeletion, ","))}, nil
}

func removeFinalizers(c client.Client, p *corev1.Pod) error {
	emptyFinalizer := make([]string, 0)
	p.ObjectMeta.SetFinalizers(emptyFinalizer)

	err := c.Update(context.TODO(), p)
	if err != nil {
		return err
	}

	return nil
}
