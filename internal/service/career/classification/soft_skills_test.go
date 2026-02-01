package classification

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Soft Skills Keyword Detection", func() {
	var classifier *Classifier

	BeforeEach(func() {
		classifier = NewClassifier()
	})

	Describe("Communication Keywords", func() {
		DescribeTable("should detect communication competency",
			func(text string) {
				event := fixtures.EventWith("", text, "", "")
				category := classifier.Classify(event)
				Expect(category).To(Equal(CommunicationCompetency))
			},
			Entry("from present keyword", "Communicated the update and presented the brief to clarify the report"),
			Entry("from document keyword", "Documented the brief and communicated the update to clarify findings"),
			Entry("from explain keyword", "Explained the brief and communicated the report to clarify the update"),
			Entry("from articulate keyword", "Articulated the brief and communicated the report to update stakeholders"),
		)

		It("should have at least 8 communication keywords", func() {
			keywords := classifier.GetCategoryKeywords(CommunicationCompetency)
			Expect(len(keywords)).To(BeNumerically(">=", 8))
		})

		It("should include key communication terms", func() {
			keywords := classifier.GetCategoryKeywords(CommunicationCompetency)
			Expect(keywords).To(ContainElement("communicate"))
			Expect(keywords).To(ContainElement("present"))
			Expect(keywords).To(ContainElement("document"))
			Expect(keywords).To(ContainElement("stakeholder"))
		})
	})

	Describe("Collaboration Keywords", func() {
		DescribeTable("should detect collaboration competency",
			func(text string) {
				event := fixtures.EventWith("", text, "", "")
				category := classifier.Classify(event)
				Expect(category).To(Equal(CollaborationCompetency))
			},
			Entry("from collaborate keyword", "Collaborated with a joint partner to align and facilitate together"),
			Entry("from partner keyword", "Partnered with a joint group to align and cooperate together"),
			Entry("from facilitate keyword", "Facilitated alignment with a joint partner to cooperate together"),
			Entry("from cooperate keyword", "Cooperated with a joint partner to align and facilitate together"),
		)

		It("should have at least 8 collaboration keywords", func() {
			keywords := classifier.GetCategoryKeywords(CollaborationCompetency)
			Expect(len(keywords)).To(BeNumerically(">=", 8))
		})

		It("should include key collaboration terms", func() {
			keywords := classifier.GetCategoryKeywords(CollaborationCompetency)
			Expect(keywords).To(ContainElement("collaborate"))
			Expect(keywords).To(ContainElement("partner"))
			Expect(keywords).To(ContainElement("facilitate"))
			Expect(keywords).To(ContainElement("together"))
		})
	})

	Describe("Problem-Solving Keywords", func() {
		DescribeTable("should detect problem-solving competency",
			func(text string) {
				event := fixtures.EventWith("", text, "", "")
				category := classifier.Classify(event)
				Expect(category).To(Equal(ProblemSolvingCompetency))
			},
			Entry("from debug keyword", "Debugged and diagnosed the root cause then resolved the identified issue"),
			Entry("from troubleshoot keyword", "Troubleshot and diagnosed the root cause then resolved the identified issue"),
			Entry("from diagnose keyword", "Diagnosed the root cause and resolved the identified issue after debugging"),
			Entry("from resolve keyword", "Resolved the identified issue after diagnosing the root cause"),
		)

		It("should have at least 8 problem-solving keywords", func() {
			keywords := classifier.GetCategoryKeywords(ProblemSolvingCompetency)
			Expect(len(keywords)).To(BeNumerically(">=", 8))
		})

		It("should include key problem-solving terms", func() {
			keywords := classifier.GetCategoryKeywords(ProblemSolvingCompetency)
			Expect(keywords).To(ContainElement("debug"))
			Expect(keywords).To(ContainElement("troubleshoot"))
			Expect(keywords).To(ContainElement("diagnose"))
			Expect(keywords).To(ContainElement("resolve"))
		})
	})

	Describe("Project Management Keywords", func() {
		DescribeTable("should detect project-management competency",
			func(text string) {
				event := fixtures.EventWith("", text, "", "")
				category := classifier.Classify(event)
				Expect(category).To(Equal(ProjectManagementCompetency))
			},
			Entry("from estimate keyword", "Estimated the deadline and scheduled the milestone for the timeline"),
			Entry("from schedule keyword", "Scheduled the deadline and estimated the milestone for the timeline"),
			Entry("from milestone keyword", "Tracked milestones and deadlines for the estimated timeline"),
			Entry("from deadline keyword", "Met the deadline on schedule and estimated the timeline for milestones"),
		)

		It("should have at least 8 project-management keywords", func() {
			keywords := classifier.GetCategoryKeywords(ProjectManagementCompetency)
			Expect(len(keywords)).To(BeNumerically(">=", 8))
		})

		It("should include key project-management terms", func() {
			keywords := classifier.GetCategoryKeywords(ProjectManagementCompetency)
			Expect(keywords).To(ContainElement("estimate"))
			Expect(keywords).To(ContainElement("schedule"))
			Expect(keywords).To(ContainElement("milestone"))
			Expect(keywords).To(ContainElement("deadline"))
		})
	})

	Describe("Architecture Keywords", func() {
		DescribeTable("should detect architecture competency",
			func(text string) {
				event := fixtures.EventWith("", text, "", "")
				category := classifier.Classify(event)
				Expect(category).To(Equal(ArchitectureCompetency))
			},
			Entry("from architect keyword", "Architected a distributed and modular decoupled platform"),
			Entry("from distributed keyword", "Built a distributed and modular decoupled platform with scalable architecture"),
			Entry("from modular keyword", "Built a modular and decoupled platform with distributed scalable architecture"),
			Entry("from decoupled keyword", "Created a decoupled and modular platform with distributed scalable architecture"),
		)

		It("should have at least 8 architecture keywords", func() {
			keywords := classifier.GetCategoryKeywords(ArchitectureCompetency)
			Expect(len(keywords)).To(BeNumerically(">=", 8))
		})

		It("should include key architecture terms", func() {
			keywords := classifier.GetCategoryKeywords(ArchitectureCompetency)
			Expect(keywords).To(ContainElement("architect"))
			Expect(keywords).To(ContainElement("distributed"))
			Expect(keywords).To(ContainElement("modular"))
			Expect(keywords).To(ContainElement("decoupled"))
		})
	})

	Describe("Tag-Based Soft Skill Classification", func() {
		DescribeTable("should classify by tag for all soft skill categories",
			func(tag string, expected CompetencyCategory) {
				event := fixtures.EventWith("", "Generic event text", "", "")
				event.Tags = []string{tag}
				Expect(classifier.Classify(event)).To(Equal(expected))
			},
			Entry("communication tag", "communication", CommunicationCompetency),
			Entry("collaboration tag", "collaboration", CollaborationCompetency),
			Entry("problem-solving tag", "problem-solving", ProblemSolvingCompetency),
			Entry("project-management tag", "project-management", ProjectManagementCompetency),
			Entry("architecture tag", "architecture", ArchitectureCompetency),
		)

		DescribeTable("should return multiple categories from ClassifyMulti with tag",
			func(tag string, expected CompetencyCategory) {
				event := fixtures.EventWith("", "Generic event text", "", "")
				event.Tags = []string{tag}
				categories := classifier.ClassifyMulti(event)
				Expect(categories).To(ContainElement(expected))
			},
			Entry("communication tag multi", "communication", CommunicationCompetency),
			Entry("collaboration tag multi", "collaboration", CollaborationCompetency),
			Entry("problem-solving tag multi", "problem-solving", ProblemSolvingCompetency),
			Entry("project-management tag multi", "project-management", ProjectManagementCompetency),
			Entry("architecture tag multi", "architecture", ArchitectureCompetency),
		)
	})

	Describe("Priority Order", func() {
		It("should prioritize original categories over soft skills when both match", func() {
			event := fixtures.EventWith("", "Leading the strategy and guiding the team while communicating updates", "", "")
			category := classifier.Classify(event)
			Expect(category).To(Equal(LeadershipCompetency))
		})

		It("should detect soft skills when no original category matches", func() {
			event := fixtures.EventWith("", "Communicated the brief and documented the update to clarify the report", "", "")
			category := classifier.Classify(event)
			Expect(category).To(Equal(CommunicationCompetency))
		})
	})
})
