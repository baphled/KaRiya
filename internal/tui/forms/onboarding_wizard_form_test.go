package forms_test

import (
	"github.com/baphled/kariya/internal/tui/forms"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OnboardingWizardForm", func() {
	Describe("OnboardingFormData", func() {
		It("should initialize with zero values", func() {
			data := &forms.OnboardingFormData{}

			Expect(data.Name).To(BeEmpty())
			Expect(data.Email).To(BeEmpty())
			Expect(data.Location).To(BeEmpty())
			Expect(data.Title).To(BeEmpty())
			Expect(data.GitHub).To(BeEmpty())
			Expect(data.Portfolio).To(BeEmpty())
		})

		It("should store field values", func() {
			data := &forms.OnboardingFormData{
				Name:      "Jane Doe",
				Email:     "jane@example.com",
				Location:  "London, UK",
				Title:     "Senior Engineer",
				GitHub:    "janedoe",
				Portfolio: "https://janedoe.dev",
			}

			Expect(data.Name).To(Equal("Jane Doe"))
			Expect(data.Email).To(Equal("jane@example.com"))
			Expect(data.Location).To(Equal("London, UK"))
			Expect(data.Title).To(Equal("Senior Engineer"))
			Expect(data.GitHub).To(Equal("janedoe"))
			Expect(data.Portfolio).To(Equal("https://janedoe.dev"))
		})
	})

	Describe("NewOnboardingWizardForm", func() {
		var data *forms.OnboardingFormData

		BeforeEach(func() {
			data = &forms.OnboardingFormData{}
		})

		It("should create a form", func() {
			form := forms.NewOnboardingWizardForm(data, 80, 40)

			Expect(form).NotTo(BeNil())
		})

		It("should create a form with the given width", func() {
			form := forms.NewOnboardingWizardForm(data, 60, 40)

			Expect(form).NotTo(BeNil())
		})

		It("should create a form with pre-filled data", func() {
			data.Name = "Existing User"
			data.Email = "existing@example.com"

			form := forms.NewOnboardingWizardForm(data, 80, 40)

			Expect(form).NotTo(BeNil())
		})

		It("should bind to data fields via pointers", func() {
			data.Name = "Initial"

			_ = forms.NewOnboardingWizardForm(data, 80, 40)

			Expect(data.Name).To(Equal("Initial"))
		})
	})

	Describe("NewOnboardingWizardForm with contact step", func() {
		It("should include contact step in form", func() {
			data := &forms.OnboardingFormData{
				Email:    "test@example.com",
				Location: "London, UK",
			}
			form := forms.NewOnboardingWizardForm(data, 80, 40)

			Expect(form).NotTo(BeNil())
			Expect(data.Email).To(Equal("test@example.com"))
			Expect(data.Location).To(Equal("London, UK"))
		})

		It("should allow email field to be updated", func() {
			data := &forms.OnboardingFormData{}
			form := forms.NewOnboardingWizardForm(data, 80, 40)

			Expect(form).NotTo(BeNil())
			data.Email = "newemail@example.com"
			Expect(data.Email).To(Equal("newemail@example.com"))
		})

		It("should allow location field to be updated", func() {
			data := &forms.OnboardingFormData{}
			form := forms.NewOnboardingWizardForm(data, 80, 40)

			Expect(form).NotTo(BeNil())
			data.Location = "New York, USA"
			Expect(data.Location).To(Equal("New York, USA"))
		})
	})
})
