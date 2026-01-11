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
			AudienceRelevance:    []string{"hiring_manager", "peer"},
			StrengthSignal:       "leadership",
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
			Expect(data.CompetencyCategories).To(Equal([]string{"technical", "leadership"}))
			Expect(data.RoleFit).To(Equal(string(career.RoleFitStaff)))
			Expect(data.AudienceRelevance).To(Equal([]string{"hiring_manager", "peer"}))
		})

		It("should apply data to fact", func() {
			newFact := &career.Fact{
				StrengthSignal: "existing_signal", // Pre-existing signal
			}
			data := &forms.FactFormData{
				Text:                 "Improved system performance by 50%",
				CompetencyCategories: []string{"technical", "product"},
				RoleFit:              string(career.RoleFitSeniorIC),
				AudienceRelevance:    []string{"hiring_manager", "recruiter"},
			}

			err := forms.ApplyFactFormData(newFact, data)

			Expect(err).NotTo(HaveOccurred())
			Expect(newFact.Text).To(Equal("Improved system performance by 50%"))
			Expect(newFact.CompetencyCategories).To(Equal([]string{"technical", "product"}))
			Expect(newFact.RoleFit).To(Equal(career.RoleFitSeniorIC))
			Expect(newFact.AudienceRelevance).To(Equal([]string{"hiring_manager", "recruiter"}))
			// StrengthSignal is NOT modified by ApplyFactFormData
			Expect(newFact.StrengthSignal).To(Equal("existing_signal"))
		})

		It("should handle empty competency categories", func() {
			newFact := &career.Fact{}
			data := &forms.FactFormData{
				Text:                 "Sample fact text here",
				CompetencyCategories: []string{},
				RoleFit:              string(career.RoleFitStaff),
				AudienceRelevance:    []string{},
			}

			err := forms.ApplyFactFormData(newFact, data)

			Expect(err).NotTo(HaveOccurred())
			Expect(newFact.CompetencyCategories).To(Equal([]string{}))
			Expect(newFact.AudienceRelevance).To(Equal([]string{}))
		})

		It("should handle nil slices from fact", func() {
			factWithNilSlices := &career.Fact{
				ID:                   "fact-456",
				Text:                 "Some fact",
				CompetencyCategories: nil,
				RoleFit:              career.RoleFitStaff,
				AudienceRelevance:    nil,
				StrengthSignal:       "technical expertise",
			}

			data := forms.GetFactFormData(factWithNilSlices)

			Expect(data.CompetencyCategories).NotTo(BeNil())
			Expect(data.CompetencyCategories).To(Equal([]string{}))
			Expect(data.AudienceRelevance).NotTo(BeNil())
			Expect(data.AudienceRelevance).To(Equal([]string{}))
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

	Describe("CompetencyCategoryOptions", func() {
		It("should return all 6 competency category options", func() {
			options := forms.CompetencyCategoryOptions()

			Expect(options).To(HaveLen(6))
		})
	})

	Describe("AudienceRelevanceOptions", func() {
		It("should return all 3 audience relevance options", func() {
			options := forms.AudienceRelevanceOptions()

			Expect(options).To(HaveLen(3))
		})
	})

	Describe("NewFactEditorFormWithData", func() {
		It("should create a form with initial data", func() {
			data := &forms.FactFormData{
				Text:                 "Initial fact text here",
				CompetencyCategories: []string{"technical", "leadership"},
				RoleFit:              string(career.RoleFitPrincipal),
				AudienceRelevance:    []string{"hiring_manager"},
			}

			form := forms.NewFactEditorFormWithData(data)

			Expect(form).NotTo(BeNil())
		})
	})
})
