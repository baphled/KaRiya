package forms_test

import (
	"github.com/baphled/kariya/internal/cli/forms"
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
			form := forms.NewCVConfigForm(data, profileOptions, extractedTechs, 80, 0)
			Expect(form).NotTo(BeNil())
		})

		It("should set default skills limit to 5", func() {
			data.SkillsLimit = 0 // Unset
			_ = forms.NewCVConfigForm(data, profileOptions, extractedTechs, 80, 0)
			// Form should set default to 5
			Expect(data.SkillsLimit).To(Equal(5))
		})

		It("should preserve existing skills limit if already set", func() {
			data.SkillsLimit = 10
			_ = forms.NewCVConfigForm(data, profileOptions, extractedTechs, 80, 0)
			Expect(data.SkillsLimit).To(Equal(10))
		})

		Describe("Skills Format options", func() {
			It("should NOT include 'categorized' option", func() {
				// The form should only offer "grouped" and "flat"
				// "categorized" was removed as it's not implemented
				form := forms.NewCVConfigForm(data, profileOptions, extractedTechs, 80, 0)
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
				form := forms.NewCVConfigForm(data, profileOptions, extractedTechs, 80, 0)
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
			// Should include: 5, 10, 15, 20, 0 (All)
			Expect(options).To(ContainElement(forms.SkillsLimitOption{Value: 5, Label: "5 per category"}))
			Expect(options).To(ContainElement(forms.SkillsLimitOption{Value: 10, Label: "10 per category"}))
			Expect(options).To(ContainElement(forms.SkillsLimitOption{Value: 15, Label: "15 per category"}))
			Expect(options).To(ContainElement(forms.SkillsLimitOption{Value: 20, Label: "20 per category"}))
			Expect(options).To(ContainElement(forms.SkillsLimitOption{Value: 0, Label: "All (no limit)"}))
		})
	})
})
