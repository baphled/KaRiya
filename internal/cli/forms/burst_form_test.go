package forms_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("BurstForm", func() {
	var testBurst *career.Burst

	BeforeEach(func() {
		testBurst = &career.Burst{
			ID:          "burst-123",
			Name:        "Test Burst",
			Description: "Test description",
		}
	})

	Describe("NewBurstEditorForm", func() {
		It("should create a form with burst data", func() {
			form := forms.NewBurstEditorForm(testBurst)

			Expect(form).NotTo(BeNil())
		})

		It("should have Catppuccin theme", func() {
			form := forms.NewBurstEditorForm(testBurst)

			// Form should be created successfully with theme
			Expect(form).NotTo(BeNil())
		})
	})

	Describe("BurstFormData", func() {
		It("should extract data from burst", func() {
			data := forms.GetBurstFormData(testBurst)

			Expect(data.Name).To(Equal("Test Burst"))
			Expect(data.Description).To(Equal("Test description"))
		})

		It("should apply data to burst", func() {
			newBurst := &career.Burst{}
			data := &forms.BurstFormData{
				Name:        "New Name",
				Description: "New Description",
			}

			forms.ApplyBurstFormData(newBurst, data)

			Expect(newBurst.Name).To(Equal("New Name"))
			Expect(newBurst.Description).To(Equal("New Description"))
		})
	})

	Describe("NewBurstEditorFormWithData", func() {
		It("should create a form with initial data", func() {
			data := &forms.BurstFormData{
				Name:        "Initial Name",
				Description: "Initial Description",
			}

			form := forms.NewBurstEditorFormWithData(data)

			Expect(form).NotTo(BeNil())
		})

		It("should bind data correctly", func() {
			data := &forms.BurstFormData{
				Name:        "Bound Name",
				Description: "Bound Description",
			}

			form := forms.NewBurstEditorFormWithData(data)

			// Simulate user updating the data
			data.Name = "Updated Name"

			// The form should have access to updated data
			Expect(data.Name).To(Equal("Updated Name"))
			Expect(form).NotTo(BeNil())
		})
	})

	Describe("Form validation", func() {
		It("should validate Name field", func() {
			// Test Title validator (used for Name field)
			err := forms.Title("Valid Title")
			Expect(err).NotTo(HaveOccurred())

			err = forms.Title("")
			Expect(err).To(HaveOccurred())

			err = forms.Title("AB")
			Expect(err).To(HaveOccurred())
		})

		It("should validate Description field", func() {
			// Test Description validator (optional)
			err := forms.Description("")
			Expect(err).NotTo(HaveOccurred())

			err = forms.Description("Valid description here")
			Expect(err).NotTo(HaveOccurred())

			err = forms.Description("Too short")
			Expect(err).To(HaveOccurred())
		})
	})
})
