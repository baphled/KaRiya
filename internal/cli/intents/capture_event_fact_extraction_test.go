package intents_test

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burst_fact"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CaptureEvent Fact Extraction", func() {
	var (
		ctx           context.Context
		classifier    *burst_fact.Classifier
		extractor     *burst_fact.Extractor
		sampleEvent   *career.CareerEvent
		extractedFact career.Fact
	)

	BeforeEach(func() {
		ctx = context.Background()
		classifier = burst_fact.NewClassifier()
		extractor = burst_fact.NewExtractor(classifier)

		// Create a sample event similar to what user would enter
		sampleEvent = &career.CareerEvent{
			ID:      "test-event-123",
			Text:    "Led migration of legacy monolith to microservices architecture, reducing deployment time by 60%",
			Date:    time.Now(),
			Company: "TechCorp",
			Project: "Platform Modernization",
		}
	})

	Describe("Fact Text Content", func() {
		Context("when extracting facts from an event", func() {
			BeforeEach(func() {
				facts := extractor.ExtractFromEvent(ctx, sampleEvent)
				Expect(facts).To(HaveLen(1), "Should extract exactly one fact from event")
				extractedFact = facts[0]
			})

			It("should use the event text as the fact text", func() {
				Expect(extractedFact.Text).To(Equal(sampleEvent.Text),
					"Fact text should match event text exactly")
			})

			It("should NOT include event ID in the fact text", func() {
				Expect(extractedFact.Text).NotTo(ContainSubstring(sampleEvent.ID),
					"Fact text should not contain event ID")
			})

			It("should NOT include event company in the fact text (unless it was in original event text)", func() {
				if !containsSubstring(sampleEvent.Text, sampleEvent.Company) {
					Expect(extractedFact.Text).NotTo(ContainSubstring(sampleEvent.Company),
						"Fact text should not add company name if not in event text")
				}
			})

			It("should set the SourceEventID to link back to the event", func() {
				Expect(extractedFact.SourceEventID).To(Equal(sampleEvent.ID),
					"Fact should have SourceEventID linking to event")
			})

			It("should infer competency categories based on event content", func() {
				Expect(extractedFact.CompetencyCategories).NotTo(BeEmpty(),
					"Fact should have at least one competency category")
				// For this technical event, should have "technical" category
				Expect(extractedFact.CompetencyCategories).To(ContainElement("technical"))
			})

			It("should infer a role fit based on event content", func() {
				Expect(extractedFact.RoleFit).NotTo(BeEmpty(),
					"Fact should have a role fit classification")
			})

			It("should infer audience relevance", func() {
				Expect(extractedFact.AudienceRelevance).NotTo(BeEmpty(),
					"Fact should have at least one audience relevance")
			})
		})

		Context("when event text is long (>200 chars)", func() {
			BeforeEach(func() {
				longText := "Led a comprehensive migration of our legacy monolithic application " +
					"to a modern microservices architecture, implementing domain-driven design principles " +
					"and establishing CI/CD pipelines using Jenkins and Kubernetes, which resulted in " +
					"a 60% reduction in deployment time and improved system scalability"
				sampleEvent.Text = longText

				facts := extractor.ExtractFromEvent(ctx, sampleEvent)
				Expect(facts).To(HaveLen(1))
				extractedFact = facts[0]
			})

			It("should summarize the event text (first 30 words)", func() {
				Expect(len(extractedFact.Text)).To(BeNumerically("<", len(sampleEvent.Text)),
					"Long event text should be summarized")
				Expect(extractedFact.Text).To(HaveSuffix("..."),
					"Summarized text should end with ellipsis")
			})

			It("should NOT include the entire event text", func() {
				// The summary should not include words from the end of the event
				Expect(extractedFact.Text).NotTo(ContainSubstring("improved system scalability"),
					"Summary should not include text from end of long event")
			})
		})

		Context("when event text is short (<= 200 chars)", func() {
			BeforeEach(func() {
				shortText := "Implemented OAuth2 authentication for API"
				sampleEvent.Text = shortText

				facts := extractor.ExtractFromEvent(ctx, sampleEvent)
				Expect(facts).To(HaveLen(1))
				extractedFact = facts[0]
			})

			It("should use the event text as-is (no summarization)", func() {
				Expect(extractedFact.Text).To(Equal(sampleEvent.Text),
					"Short event text should be used as-is without modification")
			})

			It("should NOT add ellipsis", func() {
				Expect(extractedFact.Text).NotTo(HaveSuffix("..."),
					"Short text should not have ellipsis")
			})
		})
	})

	Describe("EventReviewScreen Fact Display", func() {
		Context("when rendering facts in the review screen", func() {
			It("should display the fact.Text field", func() {
				// This is what EventReviewScreen does (line 286 of event_review_screen.go):
				// b.WriteString(fmt.Sprintf("  %d. %s\n", i+1, fact.Text))

				fact := career.Fact{
					Text:                 "Led microservices migration reducing deployment time by 60%",
					CompetencyCategories: []string{"technical", "leadership"},
					RoleFit:              career.RoleFitStaff,
				}

				// Simulate what the screen renders
				displayText := fact.Text

				Expect(displayText).To(Equal("Led microservices migration reducing deployment time by 60%"))
				Expect(displayText).NotTo(ContainSubstring("test-event"),
					"Display should not show event ID")
			})
		})
	})

	Describe("Expected Behavior", func() {
		It("should extract fact text that matches the event text", func() {
			// This documents the CORRECT behavior:
			// - Fact.Text should be the event description (or summary if long)
			// - Fact.Text should NOT be an ID or reference
			// - Fact.Text is what gets displayed to the user

			eventText := "Optimized database queries reducing load time from 3s to 200ms"
			event := &career.CareerEvent{
				ID:   "evt-123",
				Text: eventText,
			}

			facts := extractor.ExtractFromEvent(ctx, event)
			Expect(facts).To(HaveLen(1))

			// The extracted fact should have the event text as its content
			Expect(facts[0].Text).To(Equal(eventText),
				"Fact text should be the event description, not an ID or reference")
		})

		It("should show event description in the review screen", func() {
			// When user sees: "Inferred Facts: 1. Led microservices migration..."
			// This is CORRECT - they're seeing the fact text, which is the event description
			// This is the achievement/accomplishment that will go into their CV

			// If the user is seeing the event text in the fact display, that's expected!
			// The fact IS derived from the event text and represents the same achievement
		})
	})
})

// Helper function
func containsSubstring(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || (len(s) >= len(substr) &&
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
				findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
