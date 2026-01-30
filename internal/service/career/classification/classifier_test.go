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

	Context("Soft Skill Classification", func() {
		DescribeTable("Soft Skill Event Classification",
			func(eventText string, tags []string, expectedCategory CompetencyCategory) {
				event := &career.Event{
					Text: eventText,
					Tags: tags,
				}

				category := classifier.Classify(event)
				Expect(category).To(Equal(expectedCategory))
			},
			Entry("Communication event via explicit tag",
				"Presented quarterly results to stakeholders",
				[]string{"communication"},
				CommunicationCompetency,
			),
			Entry("Collaboration event via explicit tag",
				"Worked with cross-functional teams on delivery",
				[]string{"collaboration"},
				CollaborationCompetency,
			),
			Entry("Problem-solving event via explicit tag",
				"Debugged critical production issue",
				[]string{"problem-solving"},
				ProblemSolvingCompetency,
			),
			Entry("Project management event via explicit tag",
				"Planned sprint milestones for Q3 delivery",
				[]string{"project-management"},
				ProjectManagementCompetency,
			),
			Entry("Architecture event via explicit tag",
				"Designed scalable microservices platform",
				[]string{"architecture"},
				ArchitectureCompetency,
			),
		)

		DescribeTable("Soft Skill Keyword Detection",
			func(eventText string, expectedCategory CompetencyCategory) {
				event := &career.Event{
					Text: eventText,
					Tags: []string{},
				}

				category := classifier.Classify(event)
				Expect(category).To(Equal(expectedCategory))
			},
			Entry("Communication via presentation keywords",
				"Communicated updates and clarified the brief to explain the report",
				CommunicationCompetency,
			),
			Entry("Collaboration via teamwork keywords",
				"Collaborated with a joint partner to align and facilitate together",
				CollaborationCompetency,
			),
			Entry("Problem-solving via debugging keywords",
				"Debugged and diagnosed the root cause then resolved the identified issue",
				ProblemSolvingCompetency,
			),
			Entry("Project management via planning keywords",
				"Estimated the deadline and scheduled the milestone for the timeline",
				ProjectManagementCompetency,
			),
			Entry("Architecture via design keywords",
				"Architected a distributed and modular decoupled platform",
				ArchitectureCompetency,
			),
		)

		It("should classify communication keywords in categoryKeywords map", func() {
			keywords := classifier.GetCategoryKeywords(CommunicationCompetency)
			Expect(keywords).ToNot(BeEmpty())
			Expect(len(keywords)).To(BeNumerically(">=", 8))
		})

		It("should classify collaboration keywords in categoryKeywords map", func() {
			keywords := classifier.GetCategoryKeywords(CollaborationCompetency)
			Expect(keywords).ToNot(BeEmpty())
			Expect(len(keywords)).To(BeNumerically(">=", 8))
		})

		It("should classify problem-solving keywords in categoryKeywords map", func() {
			keywords := classifier.GetCategoryKeywords(ProblemSolvingCompetency)
			Expect(keywords).ToNot(BeEmpty())
			Expect(len(keywords)).To(BeNumerically(">=", 8))
		})

		It("should classify project-management keywords in categoryKeywords map", func() {
			keywords := classifier.GetCategoryKeywords(ProjectManagementCompetency)
			Expect(keywords).ToNot(BeEmpty())
			Expect(len(keywords)).To(BeNumerically(">=", 8))
		})

		It("should classify architecture keywords in categoryKeywords map", func() {
			keywords := classifier.GetCategoryKeywords(ArchitectureCompetency)
			Expect(keywords).ToNot(BeEmpty())
			Expect(len(keywords)).To(BeNumerically(">=", 8))
		})
	})
})
