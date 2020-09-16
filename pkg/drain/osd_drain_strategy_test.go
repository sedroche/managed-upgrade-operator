package drain

import (
	"time"

	"github.com/golang/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	upgradev1alpha1 "github.com/openshift/managed-upgrade-operator/pkg/apis/upgrade/v1alpha1"
	"github.com/openshift/managed-upgrade-operator/util/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OSD Drain Strategy", func() {

	var (
		mockCtrl          *gomock.Controller
		mockKubeClient    *mocks.MockClient
		osdDrain          DrainStrategy
		mockTimedDrainOne *MockTimeBasedDrainStrategy
		mockTimedDrainTwo *MockTimeBasedDrainStrategy
		nodeDrainConfig   *NodeDrain
	)

	Context("OSD Time Based Drain Strategy execution", func() {
		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			mockKubeClient = mocks.NewMockClient(mockCtrl)
			mockTimedDrainOne = NewMockTimeBasedDrainStrategy(mockCtrl)
			mockTimedDrainTwo = NewMockTimeBasedDrainStrategy(mockCtrl)
		})
		It("should not error if there are no Strategies", func() {
			osdDrain = &osdDrainStrategy{
				mockKubeClient,
				&corev1.Node{},
				&NodeDrain{},
				[]TimeBasedDrainStrategy{},
			}
			drainStartedFiveMinsAgo := &metav1.Time{Time: time.Now().Add(-5 * time.Minute)}
			result, err := osdDrain.Execute(drainStartedFiveMinsAgo)
			Expect(result).To(Not(BeNil()))
			Expect(err).To(BeNil())
			Expect(len(result)).To(Equal(0))
		})
		It("should execute a Time Based Drain Strategy after the assigned wait duration", func() {
			osdDrain = &osdDrainStrategy{
				mockKubeClient,
				&corev1.Node{},
				&NodeDrain{},
				[]TimeBasedDrainStrategy{mockTimedDrainOne},
			}
			gomock.InOrder(
				mockTimedDrainOne.EXPECT().GetWaitDuration().Return(time.Minute*30),
				mockTimedDrainOne.EXPECT().Execute().Times(1).Return(&DrainStrategyResult{Message: ""}, nil),
				mockTimedDrainOne.EXPECT().GetDescription().Times(1).Return("Drain one"),
			)
			drainStartedFortyFiveMinsAgo := &metav1.Time{Time: time.Now().Add(-45 * time.Minute)}
			result, err := osdDrain.Execute(drainStartedFortyFiveMinsAgo)
			Expect(result).To(Not(BeNil()))
			Expect(err).To(BeNil())
			Expect(len(result)).To(Equal(1))
		})
		It("should not execute a Time Based Drain Strategy before the assigned duration", func() {
			osdDrain = &osdDrainStrategy{
				mockKubeClient,
				&corev1.Node{},
				&NodeDrain{},
				[]TimeBasedDrainStrategy{mockTimedDrainOne},
			}
			gomock.InOrder(
				mockTimedDrainOne.EXPECT().GetWaitDuration().Return(time.Minute*60),
				mockTimedDrainOne.EXPECT().Execute().Times(0),
				mockTimedDrainOne.EXPECT().GetDescription().Times(0).Return("Drain one"),
			)
			drainStartedFortyFiveMinsAgo := &metav1.Time{Time: time.Now().Add(-45 * time.Minute)}
			result, err := osdDrain.Execute(drainStartedFortyFiveMinsAgo)
			Expect(result).To(Not(BeNil()))
			Expect(err).To(BeNil())
			Expect(len(result)).To(Equal(0))
		})
		It("should only execute Time Based Drain Strategy at the correct time if multiple strategies exist", func() {
			osdDrain = &osdDrainStrategy{
				mockKubeClient,
				&corev1.Node{},
				&NodeDrain{},
				[]TimeBasedDrainStrategy{mockTimedDrainOne, mockTimedDrainTwo},
			}
			gomock.InOrder(
				mockTimedDrainOne.EXPECT().GetWaitDuration().Return(time.Minute*30),
				mockTimedDrainOne.EXPECT().Execute().Times(1).Return(&DrainStrategyResult{Message: ""}, nil),
				mockTimedDrainOne.EXPECT().GetDescription().Times(1).Return("Drain one"),
				mockTimedDrainTwo.EXPECT().GetWaitDuration().Return(time.Minute*60),
				mockTimedDrainTwo.EXPECT().Execute().Times(0),
			)
			drainStartedFortyFiveMinsAgo := &metav1.Time{Time: time.Now().Add(-45 * time.Minute)}
			result, err := osdDrain.Execute(drainStartedFortyFiveMinsAgo)
			Expect(result).To(Not(BeNil()))
			Expect(err).To(BeNil())
			Expect(len(result)).To(Equal(1))
		})
	})

	Context("OSD Drain Strategies failures", func() {
		Context("When there are no strategies", func() {
			BeforeEach(func() {
				mockCtrl = gomock.NewController(GinkgoT())
				mockKubeClient = mocks.NewMockClient(mockCtrl)
				nodeDrainConfig = &NodeDrain{
					Timeout: 45,
				}
				osdDrain = &osdDrainStrategy{
					mockKubeClient,
					&corev1.Node{},
					nodeDrainConfig,
					[]TimeBasedDrainStrategy{},
				}
			})
			It("should not fail before default timeout has elapsed", func() {
				notLongEnough := &metav1.Time{Time: time.Now().Add(nodeDrainConfig.GetTimeOutDuration() / 2)}
				result, err := osdDrain.HasFailed(notLongEnough)
				Expect(result).To(BeFalse())
				Expect(err).To(BeNil())
			})
			It("should fail after default timeout has elapsed", func() {
				tooLongAgo := &metav1.Time{Time: time.Now().Add(-2 * nodeDrainConfig.GetTimeOutDuration())}
				result, err := osdDrain.HasFailed(tooLongAgo)
				Expect(result).To(BeTrue())
				Expect(err).To(BeNil())
			})
		})

		Context("OSD Time Based Drain Strategy failure", func() {
			BeforeEach(func() {
				mockCtrl = gomock.NewController(GinkgoT())
				mockKubeClient = mocks.NewMockClient(mockCtrl)
				mockTimedDrainOne = NewMockTimeBasedDrainStrategy(mockCtrl)
				mockTimedDrainTwo = NewMockTimeBasedDrainStrategy(mockCtrl)
				nodeDrainConfig = &NodeDrain{
					WorkerNodeTime: 8,
				}
				osdDrain = &osdDrainStrategy{
					mockKubeClient,
					&corev1.Node{},
					nodeDrainConfig,
					[]TimeBasedDrainStrategy{mockTimedDrainTwo, mockTimedDrainOne},
				}
			})
			It("should fail after the last strategy has failed + allowed time for drain to occur", func() {
				drainStartedSixtyNineMinsAgo := &metav1.Time{Time: time.Now().Add(-69 * time.Minute)}
				mockOneDuration := time.Minute * 30
				mockTwoDuration := time.Minute * 60
				gomock.InOrder(
					// Need to use 'Any' as the sort function calls these functions many times
					mockTimedDrainOne.EXPECT().GetWaitDuration().Return(mockOneDuration).AnyTimes(),
					mockTimedDrainTwo.EXPECT().GetWaitDuration().Return(mockTwoDuration).AnyTimes(),
					mockTimedDrainTwo.EXPECT().HasFailed(drainStartedSixtyNineMinsAgo).Return(true, nil),
					mockTimedDrainTwo.EXPECT().GetWaitDuration().Return(mockTwoDuration),
				)
				result, err := osdDrain.HasFailed(drainStartedSixtyNineMinsAgo)
				Expect(result).To(BeTrue())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("OSD PDB drain strategy", func() {
		var (
			pdbPodName  = "test-pdb-pod"
			pdbAppKey   = "app"
			pdbAppValue = "app1"
			pdbList     policyv1beta1.PodDisruptionBudgetList
			podList     corev1.PodList
			testNode    *corev1.Node
		)

		BeforeEach(func() {
			testNode = &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "test-node"}}
			pdbList = policyv1beta1.PodDisruptionBudgetList{
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
				},
			}
			podList = corev1.PodList{
				Items: []corev1.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: pdbPodName,
							Labels: map[string]string{
								pdbAppKey:     pdbAppValue,
								"other-label": "label1",
							},
						},
						Spec: corev1.PodSpec{
							NodeName: testNode.Name,
						},
					},
				},
			}
		})
		It("should have a PDB drain strategy if the node has a PDB Pod", func() {
			gomock.InOrder(
				mockKubeClient.EXPECT().List(gomock.Any(), gomock.Any()).SetArg(1, pdbList).Return(nil),
				mockKubeClient.EXPECT().List(gomock.Any(), gomock.Any()).SetArg(1, podList).Return(nil),
			)
			ts, err := getOsdTimedStrategies(mockKubeClient, &upgradev1alpha1.UpgradeConfig{}, testNode, &NodeDrain{})
			hasPdbStrategy := false
			for _, ts := range ts {
				if ts.GetName() == pdbStrategyName {
					hasPdbStrategy = true
				}
			}
			Expect(hasPdbStrategy).To(BeTrue())
			Expect(err).To(BeNil())
		})
	})
})
