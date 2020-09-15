package drain

import (
	corev1 "k8s.io/api/core/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pod Filter", func() {

	var (
		podList       *corev1.PodList
		passPredicate podPredicate = func(p corev1.Pod) bool {
			return true
		}
		failPredicate podPredicate = func(p corev1.Pod) bool {
			return false
		}
	)

	BeforeEach(func() {
		podList = &corev1.PodList{
			Items: []corev1.Pod{
				{}, {}, {},
			},
		}
	})

	Context("Filtering", func() {
		It("should return pods that match a predicate", func() {
			filteredPods := Filter(podList, passPredicate)
			Expect(len(filteredPods.Items)).To(Equal(len(podList.Items)))
		})
		It("should return pods that matches all predicates", func() {
			filteredPods := Filter(podList, passPredicate, passPredicate)
			Expect(len(filteredPods.Items)).To(Equal(len(podList.Items)))
		})
		It("should filter pods that do not match the predicate(s)", func() {
			filteredPods := Filter(podList, passPredicate, failPredicate, passPredicate)
			Expect(len(filteredPods.Items)).To(Equal(0))
		})
	})
})
