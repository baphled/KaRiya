package forms_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("FactForm", func() {
	var testFact *career.Fact

	BeforeEach(func() {
		testFact = fixtures.FactWithCategories("fact-123", "Led migration of monolith to microservices", "", []string{"technical", "leadership"}, []string{"hiring_manager", "peer"})
		testFact.StrengthSignal = "leadership"
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
			newFact := fixtures.FactWith("", "")
			newFact.StrengthSignal = "existing_signal"
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
			newFact := fixtures.FactWith("", "")
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
			factWithNilSlices := fixtures.FactWith("fact-456", "Some fact")
			factWithNilSlices.RoleFit = career.RoleFitStaff
			factWithNilSlices.StrengthSignal = "technical expertise"

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

	Describe("NewFactEditorFormWithDataAndHeight", func() {
		It("should create a form with specified height", func() {
			data := &forms.FactFormData{
				Text:                 "Test fact for height testing",
				CompetencyCategories: []string{"technical"},
				RoleFit:              string(career.RoleFitStaff),
				AudienceRelevance:    []string{"hiring_manager"},
			}

			form := forms.NewFactEditorFormWithDataAndHeight(data, 20)

			Expect(form).NotTo(BeNil())
			// Form should render and be scrollable
			view := form.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should create form without height when height is 0", func() {
			data := &forms.FactFormData{
				Text:                 "Test fact for height testing",
				CompetencyCategories: []string{"technical"},
				RoleFit:              string(career.RoleFitStaff),
				AudienceRelevance:    []string{"hiring_manager"},
			}

			form := forms.NewFactEditorFormWithDataAndHeight(data, 0)

			Expect(form).NotTo(BeNil())
			view := form.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should use dynamic height for scrolling", func() {
			data := &forms.FactFormData{
				Text:                 "Test fact for height testing",
				CompetencyCategories: []string{"technical"},
				RoleFit:              string(career.RoleFitStaff),
				AudienceRelevance:    []string{"hiring_manager"},
			}

			// Simulate terminal with 40 lines (40 - 20 overhead = 20 lines)
			height := forms.DefaultFormHeight(40)
			form := forms.NewFactEditorFormWithDataAndHeight(data, height)

			Expect(form).NotTo(BeNil())
			Expect(height).To(Equal(20))
		})
	})
})
