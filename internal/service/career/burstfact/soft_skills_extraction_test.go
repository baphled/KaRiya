package burstfact_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Soft Skills Fact Extraction Workflow", func() {
	var (
		ctx        context.Context
		extractor  *burstfact.Extractor
		classifier *burstfact.Classifier
	)

	BeforeEach(func() {
		ctx = context.Background() //nolint:fatcontext // test setup
		classifier = burstfact.NewClassifier()
		extractor = burstfact.NewExtractor(classifier)
	})

	Describe("Single Event Extraction with Soft Skills", func() {
		DescribeTable("should assign correct soft skill competency to extracted facts",
			func(text string, expectedCompetency string) {
				event := fixtures.EventWith("1", text, "TestCorp", "TestProject")
				facts := extractor.ExtractFromEvent(ctx, event)
				Expect(facts).NotTo(BeEmpty())
				Expect(facts[0].CompetencyCategories).To(ContainElement(expectedCompetency))
			},
			Entry("communication from presenting", "Presented the architecture decision records and documented findings for stakeholders", "communication"),
			Entry("collaboration from partnering", "Collaborated with cross-functional partners to facilitate alignment", "collaboration"),
			Entry("problem-solving from debugging", "Debugged and diagnosed the root cause of the production outage", "problem-solving"),
			Entry("project-management from planning", "Planned sprint milestones and estimated the delivery schedule", "project-management"),
			Entry("architecture from designing systems", "Architected a distributed microservices platform with scalable infrastructure", "architecture"),
		)

		It("should extract valid facts with communication competency", func() {
			event := fixtures.EventWith("1", "Communicated project status and documented decisions for stakeholders", "TestCorp", "StatusReporting")
			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts).NotTo(BeEmpty())

			fact := facts[0]
			Expect(fact.SourceEventID).To(Equal(event.ID))
			Expect(fact.CompetencyCategories).To(ContainElement("communication"))
			Expect(fact.AudienceRelevance).NotTo(BeEmpty())
			err := fact.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should extract valid facts with problem-solving competency", func() {
			event := fixtures.EventWith("1", "Resolved critical incident by diagnosing the root cause of data corruption", "TestCorp", "IncidentResponse")
			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts).NotTo(BeEmpty())

			fact := facts[0]
			Expect(fact.CompetencyCategories).To(ContainElement("problem-solving"))
			err := fact.Validate()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Burst Extraction with Soft Skills", func() {
		It("should accumulate soft skill competencies from multiple events", func() {
			events := []*career.Event{
				fixtures.EventWith("1", "Presented quarterly results and documented findings for stakeholders", "TestCorp", ""),
				fixtures.EventWith("2", "Collaborated with cross-functional partners to facilitate delivery", "TestCorp", ""),
				fixtures.EventWith("3", "Debugged and diagnosed the root cause of production failures", "TestCorp", ""),
			}

			burst := fixtures.Burst("burst-1", events[0].ID, events[1].ID, events[2].ID)
			burst.Name = "Cross-Team Quality Initiative"

			facts := extractor.ExtractFromBurst(ctx, burst, events)

			burstFact := facts[0]
			Expect(burstFact.SourceBurstID).To(Equal(burst.ID))
			Expect(burstFact.CompetencyCategories).To(ContainElement("communication"))
			Expect(burstFact.CompetencyCategories).To(ContainElement("collaboration"))
			Expect(burstFact.CompetencyCategories).To(ContainElement("problem-solving"))
		})

		It("should mix soft skills and traditional competencies in burst facts", func() {
			events := []*career.Event{
				fixtures.EventWith("1", "Led the engineering team through a critical milestone", "TestCorp", ""),
				fixtures.EventWith("2", "Presented the project roadmap and documented decisions for stakeholders", "TestCorp", ""),
				fixtures.EventWith("3", "Architected a scalable distributed platform", "TestCorp", ""),
			}

			burst := fixtures.Burst("burst-1", events[0].ID, events[1].ID, events[2].ID)
			burst.Name = "Platform Leadership Phase"

			facts := extractor.ExtractFromBurst(ctx, burst, events)

			burstFact := facts[0]
			Expect(burstFact.CompetencyCategories).To(ContainElement("leadership"))
			Expect(burstFact.CompetencyCategories).To(ContainElement("communication"))
			Expect(burstFact.CompetencyCategories).To(ContainElement("architecture"))
		})

		It("should produce valid facts for all burst events with soft skills", func() {
			events := []*career.Event{
				fixtures.EventWith("1", "Planned sprint milestones and estimated delivery deadlines", "TestCorp", "SprintPlanning"),
				fixtures.EventWith("2", "Collaborated with partners to facilitate the joint release", "TestCorp", "JointRelease"),
			}

			burst := fixtures.Burst("burst-1", events[0].ID, events[1].ID)
			burst.Name = "Coordinated Delivery"

			facts := extractor.ExtractFromBurst(ctx, burst, events)

			for _, fact := range facts {
				err := fact.Validate()
				Expect(err).NotTo(HaveOccurred())
				Expect(fact.CompetencyCategories).NotTo(BeEmpty())
				Expect(fact.AudienceRelevance).NotTo(BeEmpty())
			}
		})
	})

	Describe("Tag-Driven Soft Skill Extraction", func() {
		It("should assign soft skill competency from explicit tags", func() {
			event := fixtures.EventWith("1", "Worked on the project", "TestCorp", "GenericProject")
			event.Tags = []string{"communication", "collaboration"}

			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts).NotTo(BeEmpty())
			Expect(facts[0].CompetencyCategories).To(ContainElement("communication"))
			Expect(facts[0].CompetencyCategories).To(ContainElement("collaboration"))
		})

		It("should combine tag and text soft skill signals", func() {
			event := fixtures.EventWith("1", "Debugged the production incident", "TestCorp", "Incident")
			event.Tags = []string{"architecture"}

			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts).NotTo(BeEmpty())
			Expect(facts[0].CompetencyCategories).To(ContainElement("problem-solving"))
			Expect(facts[0].CompetencyCategories).To(ContainElement("architecture"))
		})
	})
})
