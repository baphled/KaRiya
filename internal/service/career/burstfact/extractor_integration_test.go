package burstfact_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Extractor Integration Tests", func() {
	var (
		ctx        context.Context
		extractor  *burstfact.Extractor
		classifier *burstfact.Classifier
	)

	BeforeEach(func() {
		ctx = context.Background()
		classifier = burstfact.NewClassifier()
		extractor = burstfact.NewExtractor(classifier)
	})

	Describe("Complete Event-to-Fact Extraction Workflow", func() {
		It("extracts enriched facts from detailed career event", func() {
			event := fixtures.EventWith("1", "Architected and delivered enterprise-scale microservices platform leading cross-functional team through complete migration from monolith", "TechCorp Inc.", "Platform Migration Initiative")
			event.Tags = []string{"technical", "leadership", "achievement"}

			facts := extractor.ExtractFromEvent(ctx, event)

			// Should have at least one fact
			Expect(facts).NotTo(BeEmpty())

			// First fact should reference the event
			Expect(facts[0].SourceEventID).To(Equal(event.ID))

			// Should have multiple competencies (technical + leadership from text)
			Expect(len(facts[0].CompetencyCategories)).To(BeNumerically(">=", 1))
			Expect(facts[0].CompetencyCategories).To(ContainElement("technical"))

			// Should preserve timestamps
			Expect(facts[0].CreatedAt).To(Equal(event.CreatedAt))
			Expect(facts[0].UpdatedAt).To(Equal(event.UpdatedAt))

			// Should have audience relevance
			Expect(facts[0].AudienceRelevance).NotTo(BeEmpty())
			Expect(facts[0].AudienceRelevance).To(ContainElement("peer"))
		})

		It("extracts facts with proper validation constraints", func() {
			event := fixtures.EventWith("1", "Led global engineering organization through digital transformation program", "Global Tech", "Digital Transformation")

			facts := extractor.ExtractFromEvent(ctx, event)

			// All extracted facts should be valid (no aspirational language)
			for _, fact := range facts {
				err := fact.Validate()
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})

	Describe("Complete Burst-to-Facts Extraction Workflow", func() {
		It("extracts comprehensive facts from burst with multiple related events", func() {
			event1 := fixtures.EventWith("1", "Led team through critical system refactoring", "TechCorp", "Platform Modernization")
			event2 := fixtures.EventWith("2", "Mentored engineering team on new architecture patterns", "TechCorp", "Platform Modernization")
			event3 := fixtures.EventWith("3", "Delivered optimized architecture reducing latency by 60%", "TechCorp", "Platform Modernization")

			burst := fixtures.Burst("burst-1", event1.ID, event2.ID, event3.ID)
			burst.Name = "Platform Modernization Initiative"

			facts := extractor.ExtractFromBurst(ctx, burst, []*career.Event{event1, event2, event3})

			// Should have multiple facts (burst-level + individual event facts)
			Expect(len(facts)).To(BeNumerically(">=", 4))

			// First fact should be burst-level
			Expect(facts[0].SourceBurstID).To(Equal(burst.ID))
			Expect(facts[0].Text).To(Equal("Platform Modernization Initiative"))

			// Should have competencies from multiple events combined
			Expect(facts[0].CompetencyCategories).To(ContainElement("technical"))
			Expect(facts[0].CompetencyCategories).To(ContainElement("leadership"))

			// All facts should be valid
			for _, fact := range facts {
				err := fact.Validate()
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("extracts facts demonstrating competency accumulation", func() {
			events := []*career.Event{
				fixtures.EventWith("1", "Led architectural review", "TechCorp", ""),
				fixtures.EventWith("2", "Mentored junior developers", "TechCorp", ""),
				fixtures.EventWith("3", "Designed product feature", "TechCorp", ""),
				fixtures.EventWith("4", "Consulted on optimization strategy", "TechCorp", ""),
			}

			burst := fixtures.Burst("burst-1", events[0].ID, events[1].ID, events[2].ID, events[3].ID)
			burst.Name = "Multi-faceted Growth Phase"

			facts := extractor.ExtractFromBurst(ctx, burst, events)

			burstFact := facts[0]
			Expect(burstFact.CompetencyCategories).To(ContainElement("architecture"))
			Expect(burstFact.CompetencyCategories).To(ContainElement("leadership"))
			Expect(burstFact.CompetencyCategories).To(ContainElement("mentoring"))
			Expect(burstFact.CompetencyCategories).To(ContainElement("product"))
			Expect(burstFact.CompetencyCategories).To(ContainElement("consulting"))

			// Should have multiple audience types
			Expect(burstFact.AudienceRelevance).To(ContainElement("peer"))
		})
	})

	Describe("Fact Extraction Quality Assurance", func() {
		It("ensures all extracted facts pass domain validation", func() {
			events := []*career.Event{
				fixtures.EventWith("1", "Implemented critical performance optimization", "TechCorp", ""),
				fixtures.EventWith("2", "Led migration to cloud infrastructure", "TechCorp", ""),
			}

			burst := fixtures.Burst("burst-1", events[0].ID, events[1].ID)
			burst.Name = "Infrastructure Improvement"

			facts := extractor.ExtractFromBurst(ctx, burst, events)

			// Every extracted fact must validate
			for i, fact := range facts {
				By("fact " + fact.ID)
				err := fact.Validate()
				Expect(err).NotTo(HaveOccurred(), "Fact %d should validate: %v", i, err)

				// Verify required fields are populated
				Expect(fact.ID).NotTo(BeEmpty())
				Expect(fact.Text).NotTo(BeEmpty())
				Expect(fact.CompetencyCategories).NotTo(BeEmpty())
				Expect(fact.AudienceRelevance).NotTo(BeEmpty())

				// Verify source traceability
				Expect(fact.SourceEventID != "" || fact.SourceBurstID != "").To(BeTrue())

				// Verify timestamps
				Expect(fact.CreatedAt).NotTo(BeZero())
				Expect(fact.UpdatedAt).NotTo(BeZero())
			}
		})

		It("maintains temporal integrity in extracted facts", func() {
			oldEvent := fixtures.EventWith("1", "Early career achievement", "StartupCorp", "")
			recentEvent := fixtures.EventWith("2", "Recent achievement", "TechCorp", "")

			burst := fixtures.Burst("burst-1", oldEvent.ID, recentEvent.ID)
			burst.Name = "Career Journey"
			burst.CreatedAt = oldEvent.CreatedAt
			burst.UpdatedAt = recentEvent.UpdatedAt

			facts := extractor.ExtractFromBurst(ctx, burst, []*career.Event{oldEvent, recentEvent})

			// Burst fact should preserve timeline from oldest to newest
			Expect(facts[0].CreatedAt).To(Equal(oldEvent.CreatedAt))
			Expect(facts[0].UpdatedAt).To(Equal(recentEvent.UpdatedAt))

			// Individual event facts should preserve their specific timestamps
			for _, fact := range facts[1:] {
				if fact.SourceEventID == oldEvent.ID {
					Expect(fact.CreatedAt).To(Equal(oldEvent.CreatedAt))
				}
				if fact.SourceEventID == recentEvent.ID {
					Expect(fact.CreatedAt).To(Equal(recentEvent.CreatedAt))
				}
			}
		})
	})
})
