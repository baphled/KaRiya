package burstfact_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Extractor", func() {
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

	Describe("ExtractFromEvent", func() {
		It("extracts fact from simple event", func() {
			event := fixtures.EventWith("1", "Led development of critical microservice architecture", "TechCorp", "Platform")

			facts := extractor.ExtractFromEvent(ctx, event)

			Expect(facts).To(HaveLen(1))
			Expect(facts[0].Text).To(ContainSubstring("Led development"))
			Expect(facts[0].SourceEventID).To(Equal(event.ID))
			Expect(facts[0].CompetencyCategories).NotTo(BeEmpty())
		})

		It("infers leadership competency from event text", func() {
			event := fixtures.EventWith("1", "Led cross-functional team to deliver critical project", "TechCorp", "")

			facts := extractor.ExtractFromEvent(ctx, event)

			Expect(facts).To(HaveLen(1))
			Expect(facts[0].CompetencyCategories).To(ContainElement("leadership"))
		})

		It("returns empty list for nil event", func() {
			facts := extractor.ExtractFromEvent(ctx, nil)
			Expect(facts).To(BeEmpty())
		})

		It("defaults to technical competency when no specific competency detected", func() {
			event := fixtures.EventWith("1", "General work on various projects", "TechCorp", "")

			facts := extractor.ExtractFromEvent(ctx, event)

			Expect(facts).To(HaveLen(1))
			Expect(facts[0].CompetencyCategories).To(ContainElement("technical"))
		})

		It("creates fact with unique ID", func() {
			event := fixtures.EventWith("1", "Implementation work", "TechCorp", "")

			facts := extractor.ExtractFromEvent(ctx, event)

			Expect(facts).To(HaveLen(1))
			Expect(facts[0].ID).NotTo(BeEmpty())
		})
	})

	Describe("ExtractFromBurst", func() {
		It("extracts facts from burst with multiple events", func() {
			event1 := fixtures.EventWith("1", "Led development team", "TechCorp", "")
			event2 := fixtures.EventWith("2", "Architected microservices platform", "TechCorp", "")

			burst := fixtures.Burst("burst-1", event1.ID, event2.ID)
			burst.Name = "Platform Modernization"

			facts := extractor.ExtractFromBurst(ctx, burst, []*career.Event{event1, event2})

			// Should have burst-level fact + individual event facts
			Expect(len(facts)).To(BeNumerically(">=", 3))

			// First fact should be burst-level
			Expect(facts[0].SourceBurstID).To(Equal(burst.ID))
		})

		It("generates burst fact with name", func() {
			event1 := fixtures.EventWith("1", "Led development", "TechCorp", "")

			burst := fixtures.Burst("burst-1", event1.ID, "event-2")
			burst.Name = "Platform Initiative"

			facts := extractor.ExtractFromBurst(ctx, burst, []*career.Event{event1})

			Expect(facts).NotTo(BeEmpty())
			Expect(facts[0].Text).To(Equal("Platform Initiative"))
		})

		It("returns empty list for nil burst", func() {
			event := fixtures.Event("1")

			facts := extractor.ExtractFromBurst(ctx, nil, []*career.Event{event})
			Expect(facts).To(BeEmpty())
		})

		It("returns empty list for burst with no events", func() {
			burst := fixtures.Burst("burst-1")
			burst.EventIDs = []string{}
			burst.Name = "Empty Burst"

			facts := extractor.ExtractFromBurst(ctx, burst, []*career.Event{})
			Expect(facts).To(BeEmpty())
		})

		It("includes individual event facts in burst extraction", func() {
			event1 := fixtures.EventWith("1", "Led development", "TechCorp", "")

			burst := fixtures.Burst("burst-1", event1.ID, "event-2")
			burst.Name = "Initiative"

			facts := extractor.ExtractFromBurst(ctx, burst, []*career.Event{event1})

			// Should have at least burst fact + event fact
			Expect(len(facts)).To(BeNumerically(">=", 2))
		})
	})

	Describe("RoleFitInference", func() {
		It("correctly identifies principal-level work", func() {
			event := fixtures.EventWith("1", "Established company-wide technical vision and architecture roadmap", "TechCorp", "")

			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts[0].RoleFit).To(Equal(career.RoleFitPrincipal))
		})

		It("correctly identifies EM-level work", func() {
			event := fixtures.EventWith("1", "Managed engineering team and handled hiring decisions", "TechCorp", "")

			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts[0].RoleFit).To(Equal(career.RoleFitEM))
		})

		It("defaults to senior ic for generic work", func() {
			event := fixtures.EventWith("1", "Worked on various features and improvements", "TechCorp", "")

			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts[0].RoleFit).To(Equal(career.RoleFitSeniorIC))
		})
	})

	Describe("AudienceRelevanceInference", func() {
		It("includes hiring manager for leadership work", func() {
			event := fixtures.EventWith("1", "Led engineering team through major transformation", "TechCorp", "")

			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts[0].AudienceRelevance).To(ContainElement("hiring_manager"))
		})

		It("always includes peer as audience", func() {
			event := fixtures.EventWith("1", "Fixed a bug in the system", "TechCorp", "")

			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts[0].AudienceRelevance).To(ContainElement("peer"))
		})
	})

	Describe("StrengthSignalExtraction", func() {
		It("extracts achievement signal from event text", func() {
			event := fixtures.EventWith("1", "Delivered critical feature that improved system performance by 40%", "TechCorp", "")

			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts[0].StrengthSignal).NotTo(BeEmpty())
		})
	})

	Describe("FactTextGeneration", func() {
		It("uses event text directly for concise descriptions", func() {
			event := fixtures.EventWith("1", "Led critical infrastructure project", "TechCorp", "")

			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts[0].Text).To(Equal(event.Text))
		})

		It("preserves burst name as fact text when available", func() {
			event := fixtures.EventWith("1", "Work on platform", "TechCorp", "")

			burst := fixtures.Burst("burst-1", event.ID, "event-2")
			burst.Name = "Enterprise Platform Initiative"

			facts := extractor.ExtractFromBurst(ctx, burst, []*career.Event{event})
			Expect(facts[0].Text).To(Equal("Enterprise Platform Initiative"))
		})
	})
})
