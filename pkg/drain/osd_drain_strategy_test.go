package drain

import (
	"github.com/openshift/managed-upgrade-operator/pkg/machinery"
	"time"

	"github.com/golang/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mockMachinery "github.com/openshift/managed-upgrade-operator/pkg/machinery/mocks"
	"github.com/openshift/managed-upgrade-operator/pkg/pod"
	"github.com/openshift/managed-upgrade-operator/util/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OSD Drain Strategy", func() {

	var (
		mockCtrl            *gomock.Controller
		mockKubeClient      *mocks.MockClient
		mockMachineryClient *mockMachinery.MockMachinery
		osdDrain            NodeDrainStrategy
		mockTimedDrainOne   *MockTimedDrainStrategy
		mockStrategyOne     *MockDrainStrategy
		mockTimedDrainTwo   *MockTimedDrainStrategy
		mockStrategyTwo     *MockDrainStrategy
		nodeDrainConfig     *NodeDrain
	)

	Context("Node drain Time Based Strategy execution", func() {
		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			mockKubeClient = mocks.NewMockClient(mockCtrl)
			mockMachineryClient = mockMachinery.NewMockMachinery(mockCtrl)
			mockTimedDrainOne = NewMockTimedDrainStrategy(mockCtrl)
			mockStrategyOne = NewMockDrainStrategy(mockCtrl)
			mockTimedDrainTwo = NewMockTimedDrainStrategy(mockCtrl)
			mockStrategyTwo = NewMockDrainStrategy(mockCtrl)
		})
		It("should not error if there are no Strategies", func() {
			osdDrain = &osdDrainStrategy{
				mockKubeClient,
				mockMachineryClient,
				&NodeDrain{},
				[]TimedDrainStrategy{},
			}
			fiveMinsAgo := &metav1.Time{Time: time.Now().Add(-5 * time.Minute)}
			gomock.InOrder(
				mockMachineryClient.EXPECT().IsNodeCordoned(gomock.Any()).Return(&machinery.IsCordonedResult{IsCordoned: true, AddedAt: fiveMinsAgo}),
			)
			result, err := osdDrain.Execute(&corev1.Node{})
			Expect(result).To(Not(BeNil()))
			Expect(err).To(BeNil())
			Expect(len(result)).To(Equal(0))
		})
		It("should execute a Time Based Drain Strategy after the assigned wait duration", func() {
			osdDrain = &osdDrainStrategy{
				mockKubeClient,
				mockMachineryClient,
				&NodeDrain{},
				[]TimedDrainStrategy{mockTimedDrainOne},
			}
			fortyFiveMinsAgo := &metav1.Time{Time: time.Now().Add(-45 * time.Minute)}
			gomock.InOrder(
				mockMachineryClient.EXPECT().IsNodeCordoned(gomock.Any()).Return(&machinery.IsCordonedResult{IsCordoned: true, AddedAt: fortyFiveMinsAgo}),
				mockTimedDrainOne.EXPECT().GetWaitDuration().Return(time.Minute*30),
				mockTimedDrainOne.EXPECT().GetStrategy().Return(mockStrategyOne),
				mockStrategyOne.EXPECT().Execute(gomock.Any()).Times(1).Return(&DrainStrategyResult{Message: "", HasExecuted: true}, nil),
				mockTimedDrainOne.EXPECT().GetDescription().Times(1).Return("Drain one"),
			)
			result, err := osdDrain.Execute(&corev1.Node{})
			Expect(result).To(Not(BeNil()))
			Expect(err).To(BeNil())
			Expect(len(result)).To(Equal(1))
		})
		It("should not execute a Time Based Drain Strategy before the assigned duration", func() {
			osdDrain = &osdDrainStrategy{
				mockKubeClient,
				mockMachineryClient,
				&NodeDrain{},
				[]TimedDrainStrategy{mockTimedDrainOne},
			}
			fortyFiveMinsAgo := &metav1.Time{Time: time.Now().Add(-45 * time.Minute)}
			gomock.InOrder(
				mockMachineryClient.EXPECT().IsNodeCordoned(gomock.Any()).Return(&machinery.IsCordonedResult{IsCordoned: true, AddedAt: fortyFiveMinsAgo}),
				mockTimedDrainOne.EXPECT().GetWaitDuration().Return(time.Minute*60),
				mockTimedDrainOne.EXPECT().GetStrategy().Return(mockStrategyOne),
				mockStrategyOne.EXPECT().Execute(gomock.Any()).Times(0),
				mockTimedDrainOne.EXPECT().GetDescription().Times(0).Return("Drain one"),
			)
			result, err := osdDrain.Execute(&corev1.Node{})
			Expect(result).To(Not(BeNil()))
			Expect(err).To(BeNil())
			Expect(len(result)).To(Equal(0))
		})
		It("should only execute Time Based Drain Strategy at the correct time if multiple strategies exist", func() {
			osdDrain = &osdDrainStrategy{
				mockKubeClient,
				mockMachineryClient,
				&NodeDrain{},
				[]TimedDrainStrategy{mockTimedDrainOne, mockTimedDrainTwo},
			}
			fortyFiveMinsAgo := &metav1.Time{Time: time.Now().Add(-45 * time.Minute)}
			gomock.InOrder(
				mockMachineryClient.EXPECT().IsNodeCordoned(gomock.Any()).Return(&machinery.IsCordonedResult{IsCordoned: true, AddedAt: fortyFiveMinsAgo}),
				mockTimedDrainOne.EXPECT().GetWaitDuration().Return(time.Minute*30),
				mockTimedDrainOne.EXPECT().GetStrategy().Return(mockStrategyOne),
				mockStrategyOne.EXPECT().Execute(gomock.Any()).Times(1).Return(&DrainStrategyResult{Message: "", HasExecuted: true}, nil),
				mockTimedDrainOne.EXPECT().GetDescription().Times(1).Return("Drain one"),
				mockTimedDrainTwo.EXPECT().GetWaitDuration().Return(time.Minute*60),
				mockTimedDrainTwo.EXPECT().GetStrategy().Return(mockStrategyTwo),
				mockStrategyTwo.EXPECT().Execute(gomock.Any()).Times(0),
			)
			result, err := osdDrain.Execute(&corev1.Node{})
			Expect(result).To(Not(BeNil()))
			Expect(err).To(BeNil())
			Expect(len(result)).To(Equal(1))
		})
	})

	Context("Node Drain Strategies failures", func() {
		Context("When there are no strategies", func() {
			BeforeEach(func() {
				mockCtrl = gomock.NewController(GinkgoT())
				mockKubeClient = mocks.NewMockClient(mockCtrl)
				nodeDrainConfig = &NodeDrain{
					Timeout: 45,
				}
				osdDrain = &osdDrainStrategy{
					mockKubeClient,
					mockMachineryClient,
					nodeDrainConfig,
					[]TimedDrainStrategy{},
				}
			})
			It("should not fail before default timeout wait has elapsed", func() {
				notLongEnough := &metav1.Time{Time: time.Now().Add(nodeDrainConfig.GetTimeOutDuration() / 2)}
				gomock.InOrder(
					mockMachineryClient.EXPECT().IsNodeCordoned(gomock.Any()).Return(&machinery.IsCordonedResult{IsCordoned: true, AddedAt: notLongEnough}),
				)
				result, err := osdDrain.HasFailed(&corev1.Node{})
				Expect(result).To(BeFalse())
				Expect(err).To(BeNil())
			})
			It("should fail after default timeout wait has elapsed", func() {
				tooLongAgo := &metav1.Time{Time: time.Now().Add(-2 * nodeDrainConfig.GetTimeOutDuration())}
				gomock.InOrder(
					mockMachineryClient.EXPECT().IsNodeCordoned(gomock.Any()).Return(&machinery.IsCordonedResult{IsCordoned: true, AddedAt: tooLongAgo}),
				)
				result, err := osdDrain.HasFailed(&corev1.Node{})
				Expect(result).To(BeTrue())
				Expect(err).To(BeNil())
			})
		})

		Context("Node drain Time Based Strategy failure", func() {
			BeforeEach(func() {
				mockCtrl = gomock.NewController(GinkgoT())
				mockKubeClient = mocks.NewMockClient(mockCtrl)
				mockTimedDrainOne = NewMockTimedDrainStrategy(mockCtrl)
				mockTimedDrainTwo = NewMockTimedDrainStrategy(mockCtrl)
				nodeDrainConfig = &NodeDrain{
					ExpectedNodeDrainTime: 8,
					Timeout:               15,
				}
				osdDrain = &osdDrainStrategy{
					mockKubeClient,
					mockMachineryClient,
					nodeDrainConfig,
					[]TimedDrainStrategy{mockTimedDrainTwo, mockTimedDrainOne},
				}
			})
			It("should fail after the last strategy has failed + allowed time for drain to occur", func() {
				drainStartedSixtyNineMinsAgo := &metav1.Time{Time: time.Now().Add(-69 * time.Minute)}
				mockOneDuration := time.Minute * 30
				mockTwoDuration := time.Minute * 60
				gomock.InOrder(
					mockMachineryClient.EXPECT().IsNodeCordoned(gomock.Any()).Return(&machinery.IsCordonedResult{IsCordoned: true, AddedAt: drainStartedSixtyNineMinsAgo}),
					// Need to use 'Any' as the sort function calls these functions many times
					mockTimedDrainOne.EXPECT().GetWaitDuration().Return(mockOneDuration).AnyTimes(),
					mockTimedDrainTwo.EXPECT().GetWaitDuration().Return(mockTwoDuration).AnyTimes(),
					mockTimedDrainOne.EXPECT().GetWaitDuration().Return(mockOneDuration),
					mockTimedDrainTwo.EXPECT().GetWaitDuration().Return(mockTwoDuration),
					mockTimedDrainTwo.EXPECT().GetWaitDuration().Return(mockTwoDuration),
					mockTimedDrainTwo.EXPECT().GetWaitDuration().Return(mockTwoDuration),
				)
				result, err := osdDrain.HasFailed(&corev1.Node{})
				Expect(result).To(BeTrue())
				Expect(err).To(BeNil())
			})
			It("should fail after default timeout wait has elapsed + allowed time for drain to occur if max strategy wait duration is less", func() {
				mockOneDuration := time.Minute * 5
				mockTwoDuration := time.Minute * 10
				thirtyOneMinsAgo := &metav1.Time{Time: time.Now().Add(-16*time.Minute - nodeDrainConfig.GetTimeOutDuration())}
				gomock.InOrder(
					mockMachineryClient.EXPECT().IsNodeCordoned(gomock.Any()).Return(&machinery.IsCordonedResult{IsCordoned: true, AddedAt: thirtyOneMinsAgo}),
					// Need to use 'Any' as the sort function calls these functions many times
					mockTimedDrainOne.EXPECT().GetWaitDuration().Return(mockOneDuration).AnyTimes(),
					mockTimedDrainTwo.EXPECT().GetWaitDuration().Return(mockTwoDuration).AnyTimes(),
					mockTimedDrainOne.EXPECT().GetWaitDuration().Return(mockOneDuration),
					mockTimedDrainTwo.EXPECT().GetWaitDuration().Return(mockTwoDuration),
					mockTimedDrainTwo.EXPECT().GetWaitDuration().Return(mockTwoDuration),
					mockTimedDrainTwo.EXPECT().GetWaitDuration().Return(mockTwoDuration),
				)

				result, _ := osdDrain.HasFailed(&corev1.Node{})
				Expect(result).To(BeTrue())
			})
			It("should not fail if there are pending strategies", func() {
				mockOneDuration := time.Minute * 10
				mockTwoDuration := time.Minute * 30
				twentyMinsAgo := &metav1.Time{Time: time.Now().Add(-20 * time.Minute)}
				gomock.InOrder(
					mockMachineryClient.EXPECT().IsNodeCordoned(gomock.Any()).Return(&machinery.IsCordonedResult{IsCordoned: true, AddedAt: twentyMinsAgo}),
					// Need to use 'Any' as the sort function calls these functions many times
					mockTimedDrainOne.EXPECT().GetWaitDuration().Return(mockOneDuration).AnyTimes(),
					mockTimedDrainTwo.EXPECT().GetWaitDuration().Return(mockTwoDuration).AnyTimes(),
					mockTimedDrainOne.EXPECT().GetWaitDuration().Return(mockOneDuration),
					mockTimedDrainTwo.EXPECT().GetWaitDuration().Return(mockTwoDuration),
					mockTimedDrainTwo.EXPECT().GetStrategy().Return(mockStrategyOne),
					mockStrategyOne.EXPECT().IsValid(gomock.Any()).Return(true, nil),
					mockTimedDrainOne.EXPECT().GetWaitDuration().Return(mockOneDuration),
					mockTimedDrainOne.EXPECT().GetWaitDuration().Return(mockOneDuration),
				)

				result, _ := osdDrain.HasFailed(&corev1.Node{})
				Expect(result).To(BeFalse())
			})
		})
	})

	Context("Pod Predicates", func() {
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
				filteredPods := pod.FilterPods(podList, isPdbPod(pdbList))
				Expect(len(filteredPods.Items)).To(Equal(1))
				Expect(filteredPods.Items[0].Name).To(Equal(pdbPodName))
			})
			It("should return pods that do not have an associated PodDisruptionBudget", func() {
				filteredPods := pod.FilterPods(podList, isNotPdbPod(pdbList))
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
				filteredPods := pod.FilterPods(podList, isOnNode(nodePodIsOn))
				Expect(len(filteredPods.Items)).To(Equal(1))
				Expect(filteredPods.Items[0].Name).To(Equal(podOnNode))
			})
			It("should not return pods that are on a specific node", func() {
				filteredPods := pod.FilterPods(podList, isOnNode(nodePodIsNotOn))
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
				filteredPods := pod.FilterPods(podList, isDaemonSet)
				Expect(len(filteredPods.Items)).To(Equal(1))
				Expect(filteredPods.Items[0].Name).To(Equal(daemonsetPodName))
			})
			It("should return pods that are not part of a DaemonSet", func() {
				filteredPods := pod.FilterPods(podList, isNotDaemonSet)
				Expect(len(filteredPods.Items)).To(Equal(2))
				Expect(filteredPods.Items[0].Name).To(Not(Equal(daemonsetPodName)))
				Expect(filteredPods.Items[1].Name).To(Not(Equal(daemonsetPodName)))
			})
		})
		Context("Pod Finalizers", func() {
			It("should return pods that have a finalizer", func() {
				podList = &corev1.PodList{
					Items: []corev1.Pod{
						{
							ObjectMeta: metav1.ObjectMeta{
								Finalizers: []string{"a.finalizer.com"},
							},
						},
					},
				}
				filteredPods := pod.FilterPods(podList, hasFinalizers)
				Expect(len(filteredPods.Items)).To(Equal(1))
			})
			It("should not return pods that have no finalizers", func() {
				podList = &corev1.PodList{
					Items: []corev1.Pod{{}},
				}
				filteredPods := pod.FilterPods(podList, hasFinalizers)
				Expect(len(filteredPods.Items)).To(Equal(0))
			})
		})
	})
})
