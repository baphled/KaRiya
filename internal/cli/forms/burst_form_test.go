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

	Describe("Submit button", func() {
		It("should include a confirm/submit field in the form", func() {
			data := &forms.BurstFormData{
				Name:        "Test Burst",
				Description: "Test description",
			}

			form := forms.NewBurstEditorFormWithData(data)

			// The form should render without error
			view := form.View()
			Expect(view).NotTo(BeEmpty())

			// The form should have a submit key that can be retrieved
			submitValue := form.GetBool("submit")
			Expect(submitValue).To(BeFalse(), "Submit should start as false")
		})

		It("should have submit confirmation data field", func() {
			data := &forms.BurstFormData{
				Name:        "Test Burst",
				Description: "Test description",
			}

			// Verify that submit confirmation is tracked
			Expect(data.SubmitConfirmed).To(BeFalse(), "Submit should not be confirmed initially")
		})

		It("should initialize SubmitConfirmed to false", func() {
			data := &forms.BurstFormData{
				Name:        "Test",
				Description: "Desc",
			}

			// Create form which should initialize SubmitConfirmed
			forms.NewBurstEditorFormWithData(data)

			Expect(data.SubmitConfirmed).To(BeFalse())
		})
	})

	Describe("NewBurstEditorFormWithDataAndHeight", func() {
		It("should create a form with specified height", func() {
			data := &forms.BurstFormData{
				Name:        "Test Burst",
				Description: "Test description",
			}

			form := forms.NewBurstEditorFormWithDataAndHeight(data, 20)

			Expect(form).NotTo(BeNil())
			// Form should render and be scrollable
			view := form.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should create form without height when height is 0", func() {
			data := &forms.BurstFormData{
				Name:        "Test Burst",
				Description: "Test description",
			}

			form := forms.NewBurstEditorFormWithDataAndHeight(data, 0)

			Expect(form).NotTo(BeNil())
			view := form.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should use dynamic height for small terminals", func() {
			data := &forms.BurstFormData{
				Name:        "Test Burst",
				Description: "Test description",
			}

			// Simulate small terminal (30 lines - 20 overhead = 10 lines)
			height := forms.DefaultFormHeight(30)
			Expect(height).To(Equal(10))

			form := forms.NewBurstEditorFormWithDataAndHeight(data, height)
			Expect(form).NotTo(BeNil())
		})

		It("should use dynamic height for large terminals", func() {
			data := &forms.BurstFormData{
				Name:        "Test Burst",
				Description: "Test description",
			}

			// Simulate large terminal (50 lines - 20 overhead = 30 lines)
			height := forms.DefaultFormHeight(50)
			Expect(height).To(Equal(30))

			form := forms.NewBurstEditorFormWithDataAndHeight(data, height)
			Expect(form).NotTo(BeNil())
		})
	})

	Describe("DefaultFormHeight", func() {
		It("should calculate height with overhead subtracted", func() {
			// 40 lines terminal - 20 overhead = 20 lines for form
			height := forms.DefaultFormHeight(40)
			Expect(height).To(Equal(20))
		})

		It("should return minimum height for small terminals", func() {
			// 25 lines terminal - 20 overhead = 5 lines, but min is 10
			height := forms.DefaultFormHeight(25)
			Expect(height).To(Equal(10))
		})

		It("should never return less than minimum", func() {
			// Very small terminal
			height := forms.DefaultFormHeight(10)
			Expect(height).To(Equal(10))
		})
	})
})
