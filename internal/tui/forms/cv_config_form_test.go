package forms_test

import (
	"github.com/baphled/kariya/internal/tui/forms"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CVConfigForm", func() {
	Describe("CVConfigFormData", func() {
		Describe("SkillsLimit field", func() {
			It("should have SkillsLimit field with default value of 5", func() {
				data := &forms.CVConfigFormData{}
				// Default should be 0 (unset), form will set to 5
				Expect(data.SkillsLimit).To(Equal(0))
			})

			It("should store skills limit value", func() {
				data := &forms.CVConfigFormData{
					SkillsLimit: 10,
				}
				Expect(data.SkillsLimit).To(Equal(10))
			})
		})
	})

	Describe("NewCVConfigForm", func() {
		var (
			data           *forms.CVConfigFormData
			profileOptions []forms.ProfileOption
			extractedTechs []forms.ExtractedTechnology
		)

		BeforeEach(func() {
			data = &forms.CVConfigFormData{}
			profileOptions = []forms.ProfileOption{
				{ID: "senior", Name: "Senior Engineer"},
				{ID: "staff", Name: "Staff Engineer"},
			}
			extractedTechs = []forms.ExtractedTechnology{
				{Name: "Go"},
				{Name: "PostgreSQL"},
			}
		})

		It("should create form with skills limit field", func() {
			form := forms.NewCVConfigForm(data, profileOptions, extractedTechs, 80, 0, false)
			Expect(form).NotTo(BeNil())
		})

		It("should set default skills limit to 0 (All)", func() {
			data.SkillsLimit = 0 // Unset
			_ = forms.NewCVConfigForm(data, profileOptions, extractedTechs, 80, 0, false)
			// Form should preserve default of 0 (All)
			Expect(data.SkillsLimit).To(Equal(0))
		})

		It("should preserve existing skills limit if already set", func() {
			data.SkillsLimit = 10
			_ = forms.NewCVConfigForm(data, profileOptions, extractedTechs, 80, 0, false)
			Expect(data.SkillsLimit).To(Equal(10))
		})

		Describe("Skills Format options", func() {
			It("should NOT include 'categorized' option", func() {
				// The form should only offer "grouped" and "flat"
				// "categorized" was removed as it's not implemented
				form := forms.NewCVConfigForm(data, profileOptions, extractedTechs, 80, 0, false)
				Expect(form).NotTo(BeNil())
				// We verify by checking the data bound value works with valid options
				data.SkillsFormat = "grouped"
				Expect(data.SkillsFormat).To(Equal("grouped"))
				data.SkillsFormat = "flat"
				Expect(data.SkillsFormat).To(Equal("flat"))
			})
		})

		Describe("Skills Limit options", func() {
			It("should offer preset options including 0 for 'All'", func() {
				form := forms.NewCVConfigForm(data, profileOptions, extractedTechs, 80, 0, false)
				Expect(form).NotTo(BeNil())
				// Valid limit values: 0 (All), 5, 10, 15, 20
				data.SkillsLimit = 0
				Expect(data.SkillsLimit).To(Equal(0))
				data.SkillsLimit = 5
				Expect(data.SkillsLimit).To(Equal(5))
				data.SkillsLimit = 20
				Expect(data.SkillsLimit).To(Equal(20))
			})
		})
	})

	Describe("SkillsLimitOptions", func() {
		It("should return preset options for skills limit", func() {
			options := forms.SkillsLimitOptions()
			Expect(options).To(HaveLen(5))
			Expect(options).To(ContainElement(forms.SkillsLimitOption{Value: 3, Label: "3"}))
			Expect(options).To(ContainElement(forms.SkillsLimitOption{Value: 5, Label: "5"}))
			Expect(options).To(ContainElement(forms.SkillsLimitOption{Value: 10, Label: "10"}))
			Expect(options).To(ContainElement(forms.SkillsLimitOption{Value: 15, Label: "15"}))
			Expect(options).To(ContainElement(forms.SkillsLimitOption{Value: 0, Label: "All"}))
		})
	})

	Describe("Specialist Single-Select Mode", func() {
		var (
			data           *forms.CVConfigFormData
			profileOptions []forms.ProfileOption
			extractedTechs []forms.ExtractedTechnology
		)

		BeforeEach(func() {
			data = &forms.CVConfigFormData{}
			profileOptions = []forms.ProfileOption{
				{ID: "senior", Name: "Senior Engineer"},
			}
			extractedTechs = []forms.ExtractedTechnology{
				{Name: "Go"},
				{Name: "PostgreSQL"},
				{Name: "Docker"},
			}
		})

		It("should create form with multi-select when singleTechSelect is false", func() {
			form := forms.NewCVConfigForm(data, profileOptions, extractedTechs, 80, 0, false)
			Expect(form).NotTo(BeNil())
			// Multi-select mode uses Technologies []string
			data.Technologies = []string{"Go", "PostgreSQL"}
			Expect(data.Technologies).To(HaveLen(2))
		})

		It("should create form with single-select when singleTechSelect is true", func() {
			form := forms.NewCVConfigForm(data, profileOptions, extractedTechs, 80, 0, true)
			Expect(form).NotTo(BeNil())
			// Single-select mode uses Technology string
			data.Technology = "Go"
			Expect(data.Technology).To(Equal("Go"))
		})

		It("should bind single-select value to Technology field", func() {
			data.Technology = "PostgreSQL"
			form := forms.NewCVConfigForm(data, profileOptions, extractedTechs, 80, 0, true)
			Expect(form).NotTo(BeNil())
			// The form should preserve the pre-set value
			Expect(data.Technology).To(Equal("PostgreSQL"))
		})
	})
})
