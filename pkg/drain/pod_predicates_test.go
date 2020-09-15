package drain

import (
	corev1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pod Predicates", func() {

	var (
		podList *corev1.PodList
	)

	Context("PDB Pods", func() {
		var (
			pdbPodName  = "test-pdb-pod"
			pdbAppKey   = "app"
			pdbAppValue = "app1"
			pdbList     *policyv1beta1.PodDisruptionBudgetList
		)
		BeforeEach(func() {
			pdbList = &policyv1beta1.PodDisruptionBudgetList{
				Items: []policyv1beta1.PodDisruptionBudget{
					{
						Spec: policyv1beta1.PodDisruptionBudgetSpec{
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									pdbAppKey: pdbAppValue,
								},
							},
						},
					},
					{
						Spec: policyv1beta1.PodDisruptionBudgetSpec{
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"non-existent-pod-selector": "",
								},
							},
						},
					},
				},
			}
			podList = &corev1.PodList{
				Items: []corev1.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: pdbPodName,
							Labels: map[string]string{
								pdbAppKey:     pdbAppValue,
								"other-label": "label1",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app":         "app2",
								"other-label": "label2",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app":         "app3",
								"other-label": "label3",
							},
						},
					},
				},
			}
		})
		It("should return pods that have an associated PodDisruptionBudget", func() {
			filteredPods := Filter(podList, isPdbPod(pdbList))
			Expect(len(filteredPods.Items)).To(Equal(1))
			Expect(filteredPods.Items[0].Name).To(Equal(pdbPodName))
		})
		It("should return pods that do not have an associated PodDisruptionBudget", func() {
			filteredPods := Filter(podList, isNotPdbPod(pdbList))
			Expect(len(filteredPods.Items)).To(Equal(2))
			Expect(filteredPods.Items[0].Name).To(Not(Equal(pdbPodName)))
			Expect(filteredPods.Items[1].Name).To(Not(Equal(pdbPodName)))
		})
	})

	Context("Pods on a Node", func() {
		var (
			podOnNode       = "test-pdb-pod"
			nodeWhichHasPod = "test-node"
			nodePodIsOn     = &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: nodeWhichHasPod,
				},
			}
			nodePodIsNotOn = &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dummy node",
				},
			}
		)
		BeforeEach(func() {
			podList = &corev1.PodList{
				Items: []corev1.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: podOnNode,
						},
						Spec: corev1.PodSpec{
							NodeName: nodeWhichHasPod,
						},
					},
					{
						Spec: corev1.PodSpec{
							NodeName: podOnNode + "no",
						},
					},
					{
						Spec: corev1.PodSpec{
							NodeName: podOnNode + "also no",
						},
					},
				},
			}
		})
		It("should return pods that are on a specific node", func() {
			filteredPods := Filter(podList, isOnNode(nodePodIsOn))
			Expect(len(filteredPods.Items)).To(Equal(1))
			Expect(filteredPods.Items[0].Name).To(Equal(podOnNode))
		})
		It("should not return pods that are on a specific node", func() {
			filteredPods := Filter(podList, isOnNode(nodePodIsNotOn))
			Expect(len(filteredPods.Items)).To(Equal(0))
		})
	})

	Context("DaemonSet Pods", func() {
		var (
			daemonsetPodName = "test-pdb-pod"
		)
		BeforeEach(func() {
			podList = &corev1.PodList{
				Items: []corev1.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: daemonsetPodName,
							OwnerReferences: []metav1.OwnerReference{
								{
									Kind: "DaemonSet",
								},
							},
						},
					},
					{},
					{},
				},
			}
		})
		It("should return pods that are part of a DaemonSet", func() {
			filteredPods := Filter(podList, isDaemonSet)
			Expect(len(filteredPods.Items)).To(Equal(1))
			Expect(filteredPods.Items[0].Name).To(Equal(daemonsetPodName))
		})
		It("should return pods that are not part of a DaemonSet", func() {
			filteredPods := Filter(podList, isNotDaemonSet)
			Expect(len(filteredPods.Items)).To(Equal(2))
			Expect(filteredPods.Items[0].Name).To(Not(Equal(daemonsetPodName)))
			Expect(filteredPods.Items[1].Name).To(Not(Equal(daemonsetPodName)))
		})
	})
})
