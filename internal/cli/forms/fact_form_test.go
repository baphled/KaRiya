package forms_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("FactForm", func() {
	var testFact *career.Fact

	BeforeEach(func() {
		testFact = &career.Fact{
			ID:                   "fact-123",
			Text:                 "Led migration of monolith to microservices",
			CompetencyCategories: []string{"technical", "leadership"},
			RoleFit:              career.RoleFitStaff,
			AudienceRelevance:    []string{"startup", "technical"},
			StrengthSignal:       "0.85",
		}
	})

	Describe("NewFactEditorForm", func() {
		It("should create a form with fact data", func() {
			form := forms.NewFactEditorForm(testFact)

			Expect(form).NotTo(BeNil())
		})
	})

	Describe("FactFormData", func() {
		It("should extract data from fact", func() {
			data := forms.GetFactFormData(testFact)

			Expect(data.Text).To(Equal("Led migration of monolith to microservices"))
			Expect(data.CompetencyCategories).To(Equal("technical, leadership"))
			Expect(data.RoleFit).To(Equal(string(career.RoleFitStaff)))
			Expect(data.AudienceRelevance).To(Equal("startup, technical"))
			Expect(data.StrengthSignal).To(Equal("0.85"))
		})

		It("should apply data to fact", func() {
			newFact := &career.Fact{}
			data := &forms.FactFormData{
				Text:                 "Improved system performance by 50%",
				CompetencyCategories: "performance, optimization",
				RoleFit:              string(career.RoleFitSeniorIC),
				AudienceRelevance:    "enterprise, technical",
				StrengthSignal:       "0.90",
			}

			err := forms.ApplyFactFormData(newFact, data)

			Expect(err).NotTo(HaveOccurred())
			Expect(newFact.Text).To(Equal("Improved system performance by 50%"))
			Expect(newFact.CompetencyCategories).To(Equal([]string{"performance", "optimization"}))
			Expect(newFact.RoleFit).To(Equal(career.RoleFitSeniorIC))
			Expect(newFact.AudienceRelevance).To(Equal([]string{"enterprise", "technical"}))
			Expect(newFact.StrengthSignal).To(Equal("0.90"))
		})

		It("should handle empty competency categories", func() {
			newFact := &career.Fact{}
			data := &forms.FactFormData{
				Text:                 "Sample fact",
				CompetencyCategories: "",
				RoleFit:              string(career.RoleFitStaff),
				AudienceRelevance:    "",
				StrengthSignal:       "0.5",
			}

			err := forms.ApplyFactFormData(newFact, data)

			Expect(err).NotTo(HaveOccurred())
			Expect(newFact.CompetencyCategories).To(Equal([]string{}))
			Expect(newFact.AudienceRelevance).To(Equal([]string{}))
		})
	})

	Describe("RoleFitOptions", func() {
		It("should return all role fit options", func() {
			options := forms.RoleFitOptions()

			Expect(options).To(HaveLen(4))
			Expect(options[0].Key).To(Equal(string(career.RoleFitPrincipal)))
			Expect(options[1].Key).To(Equal(string(career.RoleFitEM)))
			Expect(options[2].Key).To(Equal(string(career.RoleFitStaff)))
			Expect(options[3].Key).To(Equal(string(career.RoleFitSeniorIC)))
		})
	})

	Describe("NewFactEditorFormWithData", func() {
		It("should create a form with initial data", func() {
			data := &forms.FactFormData{
				Text:                 "Initial fact",
				CompetencyCategories: "cat1, cat2",
				RoleFit:              string(career.RoleFitPrincipal),
				AudienceRelevance:    "aud1, aud2",
				StrengthSignal:       "0.75",
			}

			form := forms.NewFactEditorFormWithData(data)

			Expect(form).NotTo(BeNil())
		})
	})
})
