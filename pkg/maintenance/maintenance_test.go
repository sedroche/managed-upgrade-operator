package maintenance

import (
	"fmt"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Alert Manager Maintenance Client", func() {
	var (
		mockCrtl      *gomock.Controller
		silenceClient *MockAlertManagerSilence
	)

	BeforeEach(func() {
		mockCrtl = gomock.NewController(GinkgoT())
		silenceClient = NewMockAlertManagerSilence(mockCrtl)
	})

	Context("Creating a Control Plane silence", func() {
		It("should not error on successful maintenance start", func() {
			silenceClient.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			end := time.Now().Add(90 * time.Minute)
			amm := alertManagerMaintenance{client: silenceClient}
			err := amm.StartControlPlane(end)
			Expect(err).Should(Not(HaveOccurred()))
		})
		It("should error on failing to start maintenance", func() {
			silenceClient.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("fake error"))
			end := time.Now().Add(90 * time.Minute)
			amm := alertManagerMaintenance{client: silenceClient}
			err := amm.StartControlPlane(end)
			Expect(err).Should(HaveOccurred())
		})
	})
})
