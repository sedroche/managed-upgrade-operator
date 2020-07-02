package maintenance

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mockMaintenance "github.com/openshift/managed-upgrade-operator/pkg/maintenance/mocks"
)

var _ = Describe("Alert Manager Silencer Client", func() {
	var (
		mockCrtl        *gomock.Controller
		mockMaintClient *mockMaintenance.MockAlertManagerSilencer
	)

	BeforeEach(func() {
		mockCrtl = gomock.NewController(GinkgoT())
		mockMaintClient = mockMaintenance.NewMockAlertManagerSilencer(mockCrtl)
	})

	Context("Creating a Control Plane silence", func() {
		It("Returns success", func() {
			start := strfmt.DateTime(time.Now().UTC())
			end := strfmt.DateTime(start.UTC().Add(90 * time.Minute))
			err := mockMaintClient.Create(createDefaultMatchers(), start, end, "Creator", "Create this maintenance")
			Expect(err).Should(HaveOccurred())

		})
	})
})
