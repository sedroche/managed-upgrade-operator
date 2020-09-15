package drain

import (
	"time"

	"github.com/golang/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/managed-upgrade-operator/util/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OSD Drain Strategy", func() {

	var (
		mockCtrl              *gomock.Controller
		mockKubeClient        *mocks.MockClient
		osdDrain              DrainStrategy
		mockTimedDrainOne     *MockTimeBasedDrainStrategy
		mockTimedDrainTwo     *MockTimeBasedDrainStrategy
		expectedDrainDuration = time.Minute * 5
	)

	Context("Default timeout behaviour", func() {
		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			mockKubeClient = mocks.NewMockClient(mockCtrl)
			osdDrain = &osdDrainStrategy{
				mockKubeClient,
				&corev1.Node{},
				expectedDrainDuration,
				[]TimeBasedDrainStrategy{},
			}
		})
		It("should fail after expectedDrainDuration has elapsed if there are no strategies", func() {
			tooLongAgo := &metav1.Time{Time: time.Now().Add(-2 * expectedDrainDuration)}
			result, err := osdDrain.HasFailed(tooLongAgo)
			Expect(result).To(BeTrue())
			Expect(err).To(BeNil())
		})
	})

	Context("OSD Time Based Drain Strategies", func() {
		Context("OSD Time Based Drain Strategies", func() {
			BeforeEach(func() {
				mockCtrl = gomock.NewController(GinkgoT())
				mockKubeClient = mocks.NewMockClient(mockCtrl)
				mockTimedDrainOne = NewMockTimeBasedDrainStrategy(mockCtrl)
				mockTimedDrainTwo = NewMockTimeBasedDrainStrategy(mockCtrl)
			})
			It("should execute a Time Based Drain Strategy after the assigned wait duration", func() {
				osdDrain = &osdDrainStrategy{
					mockKubeClient,
					&corev1.Node{},
					expectedDrainDuration,
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
					expectedDrainDuration,
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
			It("should only execute a Time Based Drain Strategy after the assigned duration if multiple strategies", func() {
				osdDrain = &osdDrainStrategy{
					mockKubeClient,
					&corev1.Node{},
					expectedDrainDuration,
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
	})
})

// TODO: HasFailed tests