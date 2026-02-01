package classification

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Classifier", func() {
	var classifier *Classifier

	BeforeEach(func() {
		classifier = NewClassifier()
	})

	Context("Classification", func() {
		DescribeTable("Event Classification",
			func(eventText string, tags []string, expectedCategory CompetencyCategory, expectedMulti []CompetencyCategory) {
				event := fixtures.EventWith("", eventText, "", "")
				event.Tags = tags

				// Single classification
				category := classifier.Classify(event)
				Expect(category).To(Equal(expectedCategory), "Classification should match expected category")

				// Multi-classification
				multiCategories := classifier.ClassifyMulti(event)
				Expect(multiCategories).To(Equal(expectedMulti), "Multi-classification should match expected categories")
			},
			Entry("Technical Event with Explicit Tag",
				"Developed a scalable microservices architecture",
				[]string{"technical"},
				TechnicalCompetency,
				[]CompetencyCategory{TechnicalCompetency},
			),
			Entry("Leadership Event with Keyword",
				"Led a cross-functional team to deliver critical project milestones",
				[]string{},
				LeadershipCompetency,
				[]CompetencyCategory{LeadershipCompetency},
			),
			Entry("Mixed Competency Event",
				"Developed and mentored junior engineers, guiding their technical growth and career development",
				[]string{},
				TechnicalCompetency,
				[]CompetencyCategory{TechnicalCompetency, MentoringCompetency},
			),
			Entry("Consulting Event",
				"Provided strategic consulting to optimize client's cloud infrastructure",
				[]string{"consulting"},
				ConsultingCompetency,
				[]CompetencyCategory{ConsultingCompetency},
			),
			Entry("Research-Oriented Event",
				"Conducted extensive research on machine learning algorithms",
				[]string{},
				ResearchCompetency,
				[]CompetencyCategory{ResearchCompetency},
			),
			Entry("Product Management Event",
				"Designed and launched MVP for innovative SaaS product",
				[]string{},
				ProductCompetency,
				[]CompetencyCategory{ProductCompetency},
			),
			Entry("Mentoring Event",
				"Coached junior developers through complex technical challenges",
				[]string{},
				MentoringCompetency,
				[]CompetencyCategory{MentoringCompetency},
			),
		)

		It("Should have default technical classification", func() {
			defaultEvent := fixtures.EventWith("", "Some generic event description", "", "")

			// Single classification should default to technical
			defaultCategory := classifier.Classify(defaultEvent)
			Expect(defaultCategory).To(Equal(TechnicalCompetency), "Default classification should be technical")

			// Multi-classification should also return technical
			defaultMultiCategories := classifier.ClassifyMulti(defaultEvent)
			Expect(defaultMultiCategories).To(Equal([]CompetencyCategory{TechnicalCompetency}), "Default multi-classification should be technical")
		})
	})
})
