package burst_fact_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burst_fact"
	"github.com/google/uuid"
)

var _ = Describe("Extractor Integration Tests", func() {
	var (
		ctx        context.Context
		extractor  *burst_fact.Extractor
		classifier *burst_fact.Classifier
	)

	BeforeEach(func() {
		ctx = context.Background()
		classifier = burst_fact.NewClassifier()
		extractor = burst_fact.NewExtractor(classifier)
	})

	Describe("Complete Event-to-Fact Extraction Workflow", func() {
		It("extracts enriched facts from detailed career event", func() {
			event := &career.CareerEvent{
				ID:   uuid.New().String(),
				Text: "Architected and delivered enterprise-scale microservices platform leading cross-functional team through complete migration from monolith",
				Date: time.Now().Add(-90 * 24 * time.Hour),
				Company: "TechCorp Inc.",
				Project: "Platform Migration Initiative",
				Tags: []string{"technical", "leadership", "achievement"},
				CreatedAt: time.Now().Add(-90 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-90 * 24 * time.Hour),
			}

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
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led global engineering organization through digital transformation program",
				Date:      time.Now().Add(-180 * 24 * time.Hour),
				Company:   "Global Tech",
				Project:   "Digital Transformation",
				CreatedAt: time.Now().Add(-180 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-180 * 24 * time.Hour),
			}

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
			event1 := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led team through critical system refactoring",
				Date:      time.Now().Add(-120 * 24 * time.Hour),
				Company:   "TechCorp",
				Project:   "Platform Modernization",
				CreatedAt: time.Now().Add(-120 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-120 * 24 * time.Hour),
			}

			event2 := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Mentored engineering team on new architecture patterns",
				Date:      time.Now().Add(-110 * 24 * time.Hour),
				Company:   "TechCorp",
				Project:   "Platform Modernization",
				CreatedAt: time.Now().Add(-110 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-110 * 24 * time.Hour),
			}

			event3 := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Delivered optimized architecture reducing latency by 60%",
				Date:      time.Now().Add(-100 * 24 * time.Hour),
				Company:   "TechCorp",
				Project:   "Platform Modernization",
				CreatedAt: time.Now().Add(-100 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-100 * 24 * time.Hour),
			}

			burst := &career.Burst{
				ID:        uuid.New().String(),
				Name:      "Platform Modernization Initiative",
				EventIDs:  []string{event1.ID, event2.ID, event3.ID},
				CreatedAt: time.Now().Add(-120 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-100 * 24 * time.Hour),
			}

			facts := extractor.ExtractFromBurst(ctx, burst, []*career.CareerEvent{event1, event2, event3})

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
			events := []*career.CareerEvent{
				{
					ID:        uuid.New().String(),
					Text:      "Led architectural review",
					Date:      time.Now().Add(-90 * 24 * time.Hour),
					Company:   "TechCorp",
					CreatedAt: time.Now().Add(-90 * 24 * time.Hour),
					UpdatedAt: time.Now().Add(-90 * 24 * time.Hour),
				},
				{
					ID:        uuid.New().String(),
					Text:      "Mentored junior developers",
					Date:      time.Now().Add(-80 * 24 * time.Hour),
					Company:   "TechCorp",
					CreatedAt: time.Now().Add(-80 * 24 * time.Hour),
					UpdatedAt: time.Now().Add(-80 * 24 * time.Hour),
				},
				{
					ID:        uuid.New().String(),
					Text:      "Designed product feature",
					Date:      time.Now().Add(-70 * 24 * time.Hour),
					Company:   "TechCorp",
					CreatedAt: time.Now().Add(-70 * 24 * time.Hour),
					UpdatedAt: time.Now().Add(-70 * 24 * time.Hour),
				},
				{
					ID:        uuid.New().String(),
					Text:      "Consulted on optimization strategy",
					Date:      time.Now().Add(-60 * 24 * time.Hour),
					Company:   "TechCorp",
					CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
					UpdatedAt: time.Now().Add(-60 * 24 * time.Hour),
				},
			}

			burst := &career.Burst{
				ID:        uuid.New().String(),
				Name:      "Multi-faceted Growth Phase",
				EventIDs:  []string{events[0].ID, events[1].ID, events[2].ID, events[3].ID},
				CreatedAt: time.Now().Add(-90 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-60 * 24 * time.Hour),
			}

			facts := extractor.ExtractFromBurst(ctx, burst, events)

			// Verify burst-level fact includes diverse competencies
			burstFact := facts[0]
			Expect(burstFact.CompetencyCategories).To(ContainElement("technical"))
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
			events := []*career.CareerEvent{
				{
					ID:        uuid.New().String(),
					Text:      "Implemented critical performance optimization",
					Date:      time.Now().Add(-60 * 24 * time.Hour),
					Company:   "TechCorp",
					CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
					UpdatedAt: time.Now().Add(-60 * 24 * time.Hour),
				},
				{
					ID:        uuid.New().String(),
					Text:      "Led migration to cloud infrastructure",
					Date:      time.Now().Add(-50 * 24 * time.Hour),
					Company:   "TechCorp",
					CreatedAt: time.Now().Add(-50 * 24 * time.Hour),
					UpdatedAt: time.Now().Add(-50 * 24 * time.Hour),
				},
			}

			burst := &career.Burst{
				ID:        uuid.New().String(),
				Name:      "Infrastructure Improvement",
				EventIDs:  []string{events[0].ID, events[1].ID},
				CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-50 * 24 * time.Hour),
			}

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
			oldEvent := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Early career achievement",
				Date:      time.Now().Add(-365 * 24 * time.Hour),
				Company:   "StartupCorp",
				CreatedAt: time.Now().Add(-365 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-365 * 24 * time.Hour),
			}

			recentEvent := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Recent achievement",
				Date:      time.Now().Add(-10 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-10 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-10 * 24 * time.Hour),
			}

			burst := &career.Burst{
				ID:        uuid.New().String(),
				Name:      "Career Journey",
				EventIDs:  []string{oldEvent.ID, recentEvent.ID},
				CreatedAt: oldEvent.CreatedAt,
				UpdatedAt: recentEvent.UpdatedAt,
			}

			facts := extractor.ExtractFromBurst(ctx, burst, []*career.CareerEvent{oldEvent, recentEvent})

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

