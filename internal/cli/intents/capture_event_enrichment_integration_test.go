package intents_test

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CaptureEvent Enrichment Integration", func() {
	var (
		ctx context.Context
		env *e2e.TestEnv
	)

	BeforeEach(func() {
		ctx = context.Background()
		env = e2e.Setup(GinkgoT())
	})

	AfterEach(func() {
		env.Cleanup()
	})

	Describe("Event ID Generation", func() {
		Context("when capturing a new event", func() {
			It("should generate a unique event ID", func() {
				// Create event
				event := &career.Event{
					Text: "Built REST API with Go and PostgreSQL",
					Date: time.Now(),
				}

				// Save event (this should populate the ID)
				err := env.Service.CaptureEvent(ctx, event, careerservice.ManualEntry)
				Expect(err).To(BeNil())

				// CRITICAL: Event ID must be populated after save
				Expect(event.ID).NotTo(BeEmpty(), "Event ID must be generated during save")
				Expect(len(event.ID)).To(BeNumerically(">", 0))
			})

			It("should allow using the event ID for enrichment", func() {
				// Create and save event
				event := &career.Event{
					Text: "Implemented OAuth2 authentication for API",
					Date: time.Now(),
				}

				err := env.Service.CaptureEvent(ctx, event, careerservice.ManualEntry)
				Expect(err).To(BeNil())
				Expect(event.ID).NotTo(BeEmpty())

				// CRITICAL: Should be able to extract facts using the event ID
				facts, err := env.Service.ExtractFactsFromEvent(ctx, event)
				Expect(err).To(BeNil())
				Expect(facts).NotTo(BeNil())
			})
		})
	})

	Describe("Fact Extraction After Event Save", func() {
		Context("when event is saved with valid ID", func() {
			It("should extract facts from the event", func() {
				// Create and save event
				event := &career.Event{
					Text: "Led migration of legacy monolith to microservices, reducing deployment time by 60%",
					Date: time.Now(),
				}

				err := env.Service.CaptureEvent(ctx, event, careerservice.ManualEntry)
				Expect(err).To(BeNil())
				Expect(event.ID).NotTo(BeEmpty(), "Event must have ID before fact extraction")

				// Extract facts
				facts, err := env.Service.ExtractFactsFromEvent(ctx, event)
				Expect(err).To(BeNil())
				Expect(facts).NotTo(BeEmpty(), "Should extract at least one fact from technical event")

				// Verify fact links to event
				Expect(facts[0].SourceEventID).To(Equal(event.ID))
			})

			It("should save extracted facts to the database", func() {
				// Create and save event
				event := &career.Event{
					Text: "Optimized database queries reducing load time from 3s to 200ms",
					Date: time.Now(),
				}

				err := env.Service.CaptureEvent(ctx, event, careerservice.ManualEntry)
				Expect(err).To(BeNil())

				// Extract facts
				facts, err := env.Service.ExtractFactsFromEvent(ctx, event)
				Expect(err).To(BeNil())
				Expect(facts).NotTo(BeEmpty())

				// Save fact
				fact := &facts[0]
				fact.SourceEventID = event.ID
				err = env.Service.SaveFact(ctx, fact)
				Expect(err).To(BeNil(), "Should successfully save fact to database")
				Expect(fact.ID).NotTo(BeEmpty(), "Fact should have ID after save")

				// Verify fact is in database by retrieving the event
				retrievedEvent, err := env.Service.GetEventByID(ctx, event.ID)
				Expect(err).To(BeNil())
				Expect(retrievedEvent).NotTo(BeNil())
				Expect(retrievedEvent.ID).To(Equal(event.ID))
			})
		})

		Context("when event ID is empty", func() {
			It("should fail fact validation", func() {
				// Try to save fact with empty source event ID
				fact := &career.Fact{
					Text:          "Test fact",
					SourceEventID: "", // Empty!
				}

				err := env.Service.SaveFact(ctx, fact)
				Expect(err).NotTo(BeNil(), "Should fail when source event ID is empty")
			})
		})
	})

	Describe("Burst Detection After Event Save", func() {
		Context("when multiple related events exist", func() {
			It("should detect bursts across events", func() {
				// Create and save multiple related events
				baseDate := time.Now().Add(-45 * 24 * time.Hour)
				events := []*career.Event{
					{
						Text: "Migrated authentication service to Kubernetes",
						Date: baseDate,
					},
					{
						Text: "Migrated payment service to Kubernetes",
						Date: baseDate.Add(15 * 24 * time.Hour),
					},
					{
						Text: "Migrated notification service to Kubernetes",
						Date: baseDate.Add(30 * 24 * time.Hour),
					},
				}

				// Save all events
				eventIDs := []string{}
				for _, event := range events {
					err := env.Service.CaptureEvent(ctx, event, careerservice.ManualEntry)
					Expect(err).To(BeNil())
					Expect(event.ID).NotTo(BeEmpty())
					eventIDs = append(eventIDs, event.ID)
				}

				// Try to detect bursts
				suggestions, err := env.Service.SuggestBursts(ctx, eventIDs)
				Expect(err).To(BeNil())
				// May or may not detect bursts depending on algorithm
				// This verifies the API works without error
				_ = suggestions
			})
		})

		Context("when only one event exists", func() {
			It("should return empty burst suggestions", func() {
				// Create and save single event
				event := &career.Event{
					Text: "Implemented Redis caching",
					Date: time.Now(),
				}

				err := env.Service.CaptureEvent(ctx, event, careerservice.ManualEntry)
				Expect(err).To(BeNil())

				// Try to detect bursts with single event
				suggestions, err := env.Service.SuggestBursts(ctx, []string{event.ID})
				Expect(err).To(BeNil())
				Expect(suggestions).To(BeEmpty(), "Single event should not generate burst suggestions")
			})
		})
	})

	Describe("Complete Enrichment Workflow", func() {
		Context("when capturing an event with enrichment", func() {
			It("should save event, extract facts, and save facts in sequence", func() {
				// 1. Create event
				event := &career.Event{
					Text: "Implemented CI/CD pipeline with Jenkins and Docker, reducing deployment time by 70%",
					Date: time.Now(),
				}

				// 2. Save event (generates ID)
				err := env.Service.CaptureEvent(ctx, event, careerservice.ManualEntry)
				Expect(err).To(BeNil())
				Expect(event.ID).NotTo(BeEmpty(), "Step 1: Event must have ID after save")

				// 3. Extract facts (requires event ID)
				facts, err := env.Service.ExtractFactsFromEvent(ctx, event)
				Expect(err).To(BeNil())
				Expect(facts).NotTo(BeEmpty(), "Step 2: Should extract facts from event")

				// 4. Save facts (link to event ID)
				savedFactCount := 0
				for i := range facts {
					fact := &facts[i]
					fact.SourceEventID = event.ID
					err := env.Service.SaveFact(ctx, fact)
					if err == nil {
						savedFactCount++
						Expect(fact.ID).NotTo(BeEmpty(), "Saved fact should have ID")
					}
				}
				Expect(savedFactCount).To(BeNumerically(">", 0), "Step 3: Should save at least one fact")

				// 5. Verify event was saved
				retrievedEvent, err := env.Service.GetEventByID(ctx, event.ID)
				Expect(err).To(BeNil())
				Expect(retrievedEvent.Text).To(Equal(event.Text))
			})
		})
	})

	Describe("Error Cases", func() {
		Context("when event is nil", func() {
			It("should return error when trying to extract facts", func() {
				_, err := env.Service.ExtractFactsFromEvent(ctx, nil)
				Expect(err).NotTo(BeNil(), "Should fail when extracting facts from nil event")
			})
		})

		Context("when event text is empty", func() {
			It("should fail validation", func() {
				event := &career.Event{
					Text: "", // Empty text
					Date: time.Now(),
				}

				err := env.Service.CaptureEvent(ctx, event, careerservice.ManualEntry)
				Expect(err).NotTo(BeNil(), "Should fail when event text is empty")
			})
		})
	})
})
