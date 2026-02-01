package burstfact

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Soft Skills Competency Inference", func() {
	var classifier *Classifier

	BeforeEach(func() {
		classifier = NewClassifier()
	})

	Describe("Communication Inference", func() {
		DescribeTable("should infer communication from text",
			func(text string) {
				competencies := classifier.InferCompetencies(text, []string{})
				Expect(competencies).To(ContainElement("communication"))
			},
			Entry("present keyword", "Presented quarterly results to the board"),
			Entry("document keyword", "Documented the API endpoints for the team"),
			Entry("explain keyword", "Explained complex concepts to non-technical audience"),
			Entry("stakeholder keyword", "Gathered requirements from stakeholders"),
			Entry("articulate keyword", "Articulated the vision for the product"),
			Entry("communicate keyword", "Communicated status updates regularly"),
			Entry("clarify keyword", "Helped clarify the acceptance criteria"),
			Entry("brief keyword", "Delivered a brief on security best practices"),
		)

		It("should infer communication from tag", func() {
			competencies := classifier.InferCompetencies("generic text", []string{"communication"})
			Expect(competencies).To(ContainElement("communication"))
		})

		It("should not infer communication from unrelated text", func() {
			competencies := classifier.InferCompetencies("deployed the application", []string{})
			Expect(competencies).NotTo(ContainElement("communication"))
		})
	})

	Describe("Collaboration Inference", func() {
		DescribeTable("should infer collaboration from text",
			func(text string) {
				competencies := classifier.InferCompetencies(text, []string{})
				Expect(competencies).To(ContainElement("collaboration"))
			},
			Entry("collaborate keyword", "Collaborated with the design team"),
			Entry("cross-functional keyword", "Worked in a cross-functional squad"),
			Entry("partner keyword", "Partnered with the QA team on test strategy"),
			Entry("facilitate keyword", "Helped facilitate the workshop"),
			Entry("align keyword", "Worked to align engineering and product"),
			Entry("together keyword", "Built the feature together with backend"),
			Entry("joint keyword", "Ran a joint planning session"),
			Entry("cooperate keyword", "Teams cooperated on the release"),
		)

		It("should infer collaboration from tag", func() {
			competencies := classifier.InferCompetencies("generic text", []string{"collaboration"})
			Expect(competencies).To(ContainElement("collaboration"))
		})

		It("should not infer collaboration from unrelated text", func() {
			competencies := classifier.InferCompetencies("wrote unit tests", []string{})
			Expect(competencies).NotTo(ContainElement("collaboration"))
		})
	})

	Describe("Problem-Solving Inference", func() {
		DescribeTable("should infer problem-solving from text",
			func(text string) {
				competencies := classifier.InferCompetencies(text, []string{})
				Expect(competencies).To(ContainElement("problem-solving"))
			},
			Entry("debug keyword", "Debugged a memory leak in production"),
			Entry("troubleshoot keyword", "Had to troubleshoot the authentication failures"),
			Entry("diagnose keyword", "Diagnosed the root cause of the outage"),
			Entry("root cause keyword", "Performed root cause analysis on latency"),
			Entry("resolve keyword", "Resolved the critical data loss incident"),
		)

		It("should infer problem-solving from tag", func() {
			competencies := classifier.InferCompetencies("generic text", []string{"problem-solving"})
			Expect(competencies).To(ContainElement("problem-solving"))
		})

		It("should not infer problem-solving from unrelated text", func() {
			competencies := classifier.InferCompetencies("shipped new landing page", []string{})
			Expect(competencies).NotTo(ContainElement("problem-solving"))
		})
	})

	Describe("Project Management Inference", func() {
		DescribeTable("should infer project-management from text",
			func(text string) {
				competencies := classifier.InferCompetencies(text, []string{})
				Expect(competencies).To(ContainElement("project-management"))
			},
			Entry("plan keyword", "Created the implementation plan for Q4"),
			Entry("estimate keyword", "Provided effort estimates for the epic"),
			Entry("schedule keyword", "Adjusted the release schedule for the holiday"),
			Entry("milestone keyword", "Tracked milestone completion across teams"),
			Entry("sprint keyword", "Facilitated sprint planning for the squad"),
			Entry("roadmap keyword", "Contributed to the product roadmap"),
			Entry("prioritize keyword", "Helped prioritize the backlog items"),
			Entry("deadline keyword", "Managed the deadline for the compliance audit"),
		)

		It("should infer project-management from tag", func() {
			competencies := classifier.InferCompetencies("generic text", []string{"project-management"})
			Expect(competencies).To(ContainElement("project-management"))
		})

		It("should not infer project-management from unrelated text", func() {
			competencies := classifier.InferCompetencies("refactored the service", []string{})
			Expect(competencies).NotTo(ContainElement("project-management"))
		})
	})

	Describe("Architecture Inference", func() {
		DescribeTable("should infer architecture from text",
			func(text string) {
				competencies := classifier.InferCompetencies(text, []string{})
				Expect(competencies).To(ContainElement("architecture"))
			},
			Entry("architect keyword", "Architected the event-driven pipeline"),
			Entry("distributed keyword", "Built a distributed caching layer"),
			Entry("microservices keyword", "Migrated the monolith to microservices"),
			Entry("scalable keyword", "Designed a scalable ingestion pipeline"),
			Entry("platform keyword", "Developed the internal developer platform"),
			Entry("modular keyword", "Refactored into a modular architecture"),
			Entry("decoupled keyword", "Created a decoupled notification system"),
		)

		It("should infer architecture from tag", func() {
			competencies := classifier.InferCompetencies("generic text", []string{"architecture"})
			Expect(competencies).To(ContainElement("architecture"))
		})

		It("should not infer architecture from unrelated text", func() {
			competencies := classifier.InferCompetencies("fixed the typo in the readme", []string{})
			Expect(competencies).NotTo(ContainElement("architecture"))
		})
	})

	Describe("Multiple Soft Skills Detection", func() {
		It("should detect communication and collaboration together", func() {
			text := "Presented the proposal and collaborated with cross-functional partners"
			competencies := classifier.InferCompetencies(text, []string{})
			Expect(competencies).To(ContainElement("communication"))
			Expect(competencies).To(ContainElement("collaboration"))
		})

		It("should detect problem-solving and architecture together", func() {
			text := "Diagnosed the scalability bottleneck and architected a distributed solution"
			competencies := classifier.InferCompetencies(text, []string{})
			Expect(competencies).To(ContainElement("problem-solving"))
			Expect(competencies).To(ContainElement("architecture"))
		})

		It("should detect project-management and communication together", func() {
			text := "Planned the sprint milestones and presented status updates to stakeholders"
			competencies := classifier.InferCompetencies(text, []string{})
			Expect(competencies).To(ContainElement("project-management"))
			Expect(competencies).To(ContainElement("communication"))
		})

		It("should detect all five soft skills from rich text", func() {
			text := "Presented the plan to stakeholders and collaborated with partners to diagnose the scalable platform milestone"
			competencies := classifier.InferCompetencies(text, []string{})
			Expect(competencies).To(ContainElement("communication"))
			Expect(competencies).To(ContainElement("collaboration"))
			Expect(competencies).To(ContainElement("problem-solving"))
			Expect(competencies).To(ContainElement("project-management"))
			Expect(competencies).To(ContainElement("architecture"))
		})
	})

	Describe("Edge Cases", func() {
		It("should default to technical when no soft skills match", func() {
			competencies := classifier.InferCompetencies("did some work", []string{})
			Expect(competencies).To(ContainElement("technical"))
			Expect(competencies).NotTo(ContainElement("communication"))
			Expect(competencies).NotTo(ContainElement("collaboration"))
			Expect(competencies).NotTo(ContainElement("problem-solving"))
			Expect(competencies).NotTo(ContainElement("project-management"))
			Expect(competencies).NotTo(ContainElement("architecture"))
		})

		It("should handle empty text with soft skill tags", func() {
			competencies := classifier.InferCompetencies("", []string{"communication", "problem-solving"})
			Expect(competencies).To(ContainElement("communication"))
			Expect(competencies).To(ContainElement("problem-solving"))
		})

		It("should combine text and tag inference", func() {
			competencies := classifier.InferCompetencies("debugged the issue", []string{"collaboration"})
			Expect(competencies).To(ContainElement("problem-solving"))
			Expect(competencies).To(ContainElement("collaboration"))
		})

		It("should handle case-insensitive text matching", func() {
			competencies := classifier.InferCompetencies("PRESENTED to STAKEHOLDERS", []string{})
			Expect(competencies).To(ContainElement("communication"))
		})
	})
})
