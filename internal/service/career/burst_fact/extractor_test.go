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

var _ = Describe("Extractor", func() {
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

	Describe("ExtractFromEvent", func() {
		It("extracts fact from simple event", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led development of critical microservice architecture",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				Project:   "Platform",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			facts := extractor.ExtractFromEvent(ctx, event)

			Expect(facts).To(HaveLen(1))
			Expect(facts[0].Text).To(ContainSubstring("Led development"))
			Expect(facts[0].SourceEventID).To(Equal(event.ID))
			Expect(facts[0].CompetencyCategories).NotTo(BeEmpty())
		})

		It("infers leadership competency from event text", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led cross-functional team to deliver critical project",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			facts := extractor.ExtractFromEvent(ctx, event)

			Expect(facts).To(HaveLen(1))
			Expect(facts[0].CompetencyCategories).To(ContainElement("leadership"))
		})

		It("returns empty list for nil event", func() {
			facts := extractor.ExtractFromEvent(ctx, nil)
			Expect(facts).To(BeEmpty())
		})

		It("defaults to technical competency when no specific competency detected", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "General work on various projects",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			facts := extractor.ExtractFromEvent(ctx, event)

			Expect(facts).To(HaveLen(1))
			Expect(facts[0].CompetencyCategories).To(ContainElement("technical"))
		})

		It("creates fact with unique ID", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Implementation work",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			facts := extractor.ExtractFromEvent(ctx, event)

			Expect(facts).To(HaveLen(1))
			Expect(facts[0].ID).NotTo(BeEmpty())
		})
	})

	Describe("ExtractFromBurst", func() {
		It("extracts facts from burst with multiple events", func() {
			event1 := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led development team",
				Date:      time.Now().Add(-60 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-60 * 24 * time.Hour),
			}
			event2 := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Architected microservices platform",
				Date:      time.Now().Add(-50 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-50 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-50 * 24 * time.Hour),
			}

			burst := &career.Burst{
				ID:        uuid.New().String(),
				Name:      "Platform Modernization",
				EventIDs:  []string{event1.ID, event2.ID},
				CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-50 * 24 * time.Hour),
			}

			facts := extractor.ExtractFromBurst(ctx, burst, []*career.CareerEvent{event1, event2})

			// Should have burst-level fact + individual event facts
			Expect(len(facts)).To(BeNumerically(">=", 3))

			// First fact should be burst-level
			Expect(facts[0].SourceBurstID).To(Equal(burst.ID))
		})

		It("generates burst fact with name", func() {
			event1 := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led development",
				Date:      time.Now().Add(-60 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-60 * 24 * time.Hour),
			}

			burst := &career.Burst{
				ID:        uuid.New().String(),
				Name:      "Platform Initiative",
				EventIDs:  []string{event1.ID},
				CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt: time.Now(),
			}

			facts := extractor.ExtractFromBurst(ctx, burst, []*career.CareerEvent{event1})

			Expect(facts).NotTo(BeEmpty())
			Expect(facts[0].Text).To(Equal("Platform Initiative"))
		})

		It("returns empty list for nil burst", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Some work",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			facts := extractor.ExtractFromBurst(ctx, nil, []*career.CareerEvent{event})
			Expect(facts).To(BeEmpty())
		})

		It("returns empty list for burst with no events", func() {
			burst := &career.Burst{
				ID:        uuid.New().String(),
				Name:      "Empty Burst",
				EventIDs:  []string{},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			facts := extractor.ExtractFromBurst(ctx, burst, []*career.CareerEvent{})
			Expect(facts).To(BeEmpty())
		})

		It("includes individual event facts in burst extraction", func() {
			event1 := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led development",
				Date:      time.Now().Add(-60 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-60 * 24 * time.Hour),
			}

			burst := &career.Burst{
				ID:        uuid.New().String(),
				Name:      "Initiative",
				EventIDs:  []string{event1.ID},
				CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt: time.Now(),
			}

			facts := extractor.ExtractFromBurst(ctx, burst, []*career.CareerEvent{event1})

			// Should have at least burst fact + event fact
			Expect(len(facts)).To(BeNumerically(">=", 2))
		})
	})

	Describe("RoleFitInference", func() {
		It("correctly identifies principal-level work", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Established company-wide technical vision and architecture roadmap",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts[0].RoleFit).To(Equal(career.RoleFitPrincipal))
		})

		It("correctly identifies EM-level work", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Managed engineering team and handled hiring decisions",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts[0].RoleFit).To(Equal(career.RoleFitEM))
		})

		It("defaults to senior ic for generic work", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Worked on various features and improvements",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts[0].RoleFit).To(Equal(career.RoleFitSeniorIC))
		})
	})

	Describe("AudienceRelevanceInference", func() {
		It("includes hiring manager for leadership work", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led engineering team through major transformation",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts[0].AudienceRelevance).To(ContainElement("hiring_manager"))
		})

		It("always includes peer as audience", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Fixed a bug in the system",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts[0].AudienceRelevance).To(ContainElement("peer"))
		})
	})

	Describe("StrengthSignalExtraction", func() {
		It("extracts achievement signal from event text", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Delivered critical feature that improved system performance by 40%",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts[0].StrengthSignal).NotTo(BeEmpty())
		})
	})

	Describe("FactTextGeneration", func() {
		It("uses event text directly for concise descriptions", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led critical infrastructure project",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts[0].Text).To(Equal(event.Text))
		})

		It("preserves burst name as fact text when available", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Work on platform",
				Date:      time.Now().Add(-60 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-60 * 24 * time.Hour),
			}

			burst := &career.Burst{
				ID:        uuid.New().String(),
				Name:      "Enterprise Platform Initiative",
				EventIDs:  []string{event.ID},
				CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt: time.Now(),
			}

			facts := extractor.ExtractFromBurst(ctx, burst, []*career.CareerEvent{event})
			Expect(facts[0].Text).To(Equal("Enterprise Platform Initiative"))
		})
	})
})
