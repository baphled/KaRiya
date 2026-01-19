package burst_fact

import (
	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Classifier", func() {
	var classifier *Classifier

	BeforeEach(func() {
		classifier = NewClassifier()
	})

	Context("when classifying role fit", func() {
		It("should classify principal role from keywords", func() {
			text := "Defined technical vision and strategy for the entire organization"
			roleFit := classifier.ClassifyRoleFit(text)
			Expect(roleFit).To(Equal(career.RoleFitPrincipal))
		})

		It("should classify EM role from keywords", func() {
			text := "Managed team of 10 engineers and hiring new talent"
			roleFit := classifier.ClassifyRoleFit(text)
			Expect(roleFit).To(Equal(career.RoleFitEM))
		})

		It("should classify staff role from keywords", func() {
			text := "Staff engineer providing deep technical expertise on complex systems"
			roleFit := classifier.ClassifyRoleFit(text)
			Expect(roleFit).To(Equal(career.RoleFitStaff))
		})

		It("should default to senior IC role", func() {
			text := "Developed new feature for the platform"
			roleFit := classifier.ClassifyRoleFit(text)
			Expect(roleFit).To(Equal(career.RoleFitSeniorIC))
		})
	})

	Context("when classifying role fit with categories", func() {
		It("should use technical category for senior_ic", func() {
			text := "Architected the backend infrastructure" // Would be principal by keywords alone
			categories := []string{"technical"}
			roleFit := classifier.ClassifyRoleFitWithCategories(text, categories)
			Expect(roleFit).To(Equal(career.RoleFitSeniorIC))
		})

		It("should use leadership category for EM", func() {
			text := "Delivered new feature" // Would be senior_ic by keywords alone
			categories := []string{"leadership"}
			roleFit := classifier.ClassifyRoleFitWithCategories(text, categories)
			Expect(roleFit).To(Equal(career.RoleFitEM))
		})

		It("should use leadership+technical categories for staff", func() {
			text := "Delivered new feature"
			categories := []string{"leadership", "technical"}
			roleFit := classifier.ClassifyRoleFitWithCategories(text, categories)
			Expect(roleFit).To(Equal(career.RoleFitStaff))
		})

		It("should elevate to principal with strong principal keywords", func() {
			text := "Set company-wide technical direction as founding engineer"
			categories := []string{"technical"}
			roleFit := classifier.ClassifyRoleFitWithCategories(text, categories)
			Expect(roleFit).To(Equal(career.RoleFitPrincipal))
		})

		It("should fall back to keyword matching when no categories", func() {
			text := "Managed team of engineers"
			roleFit := classifier.ClassifyRoleFitWithCategories(text, nil)
			Expect(roleFit).To(Equal(career.RoleFitEM))
		})

		It("should use architecture category same as technical for senior_ic", func() {
			text := "Designed the system"
			categories := []string{"architecture"}
			roleFit := classifier.ClassifyRoleFitWithCategories(text, categories)
			Expect(roleFit).To(Equal(career.RoleFitSeniorIC))
		})
	})

	Context("when classifying audience relevance", func() {
		It("should include peer for all facts", func() {
			text := "Delivered project"
			audiences := classifier.ClassifyAudienceRelevance(text, career.RoleFitSeniorIC)
			Expect(audiences).To(ContainElement("peer"))
		})

		It("should include hiring_manager and recruiter for leadership facts", func() {
			text := "Led team to deliver critical project"
			audiences := classifier.ClassifyAudienceRelevance(text, career.RoleFitEM)
			Expect(audiences).To(ContainElement("hiring_manager"))
			Expect(audiences).To(ContainElement("recruiter"))
		})

		It("should include hiring_manager for technical facts", func() {
			text := "Architected microservices system"
			audiences := classifier.ClassifyAudienceRelevance(text, career.RoleFitSeniorIC)
			Expect(audiences).To(ContainElement("hiring_manager"))
		})
	})

	Context("when extracting strength signal", func() {
		It("should extract delivery capability from 'delivered'", func() {
			text := "Delivered critical project on time"
			signal := classifier.ExtractStrengthSignal(text)
			Expect(signal).To(Equal("delivery capability"))
		})

		It("should extract leadership from 'led'", func() {
			text := "Led cross-functional team"
			signal := classifier.ExtractStrengthSignal(text)
			Expect(signal).To(Equal("leadership"))
		})

		It("should extract scalability expertise from 'scaled'", func() {
			text := "Scaled system to handle 10x traffic"
			signal := classifier.ExtractStrengthSignal(text)
			Expect(signal).To(Equal("scalability expertise"))
		})

		It("should default to professional accomplishment", func() {
			text := "Worked on interesting problem"
			signal := classifier.ExtractStrengthSignal(text)
			Expect(signal).To(Equal("professional accomplishment"))
		})
	})

	Context("when inferring competencies", func() {
		It("should infer leadership from text", func() {
			text := "Led team of engineers"
			competencies := classifier.InferCompetencies(text, []string{})
			Expect(competencies).To(ContainElement("leadership"))
		})

		It("should infer technical from text", func() {
			text := "Architected backend system"
			competencies := classifier.InferCompetencies(text, []string{})
			Expect(competencies).To(ContainElement("technical"))
		})

		It("should infer product from text", func() {
			text := "Designed new product feature for users"
			competencies := classifier.InferCompetencies(text, []string{})
			Expect(competencies).To(ContainElement("product"))
		})

		It("should infer mentoring from text", func() {
			text := "Mentored junior engineers to develop their skills"
			competencies := classifier.InferCompetencies(text, []string{})
			Expect(competencies).To(ContainElement("mentoring"))
		})

		It("should use explicit tags", func() {
			text := "Worked on project"
			tags := []string{"leadership"}
			competencies := classifier.InferCompetencies(text, tags)
			Expect(competencies).To(ContainElement("leadership"))
		})

		It("should default to technical if no competencies found", func() {
			text := "Did some work"
			competencies := classifier.InferCompetencies(text, []string{})
			Expect(competencies).To(ContainElement("technical"))
		})

		It("should infer multiple competencies", func() {
			text := "Led team to develop and ship new product feature"
			competencies := classifier.InferCompetencies(text, []string{})
			Expect(competencies).To(ContainElement("leadership"))
			Expect(competencies).To(ContainElement("mentoring"))
			Expect(competencies).To(ContainElement("product"))
		})
	})
})
