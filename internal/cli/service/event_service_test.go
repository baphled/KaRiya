package service

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLI Event Service", func() {
	var (
		cliEventService *CLIEventService
		careerSvc       *careerservice.Service
		repo            careerrepo.Repository
	)

	BeforeEach(func() {
		// Set up in-memory repository for testing
		repo = careerrepo.NewMemoryRepository()
		careerSvc = careerservice.NewService(repo)
		cliEventService = NewCLIEventService(careerSvc)
	})

	Context("Capturing Events", func() {
		It("should capture an event with default options", func() {
			ctx := context.Background()
			eventText := "Completed a significant project"
			eventDate := time.Now()
			mode := careerservice.TimelineJournaling

			// Capture the event
			err := cliEventService.CaptureEvent(ctx, eventText, eventDate, mode)

			// Assertions
			Expect(err).To(BeNil())
		})

		It("should capture an event with optional configurations", func() {
			ctx := context.Background()
			eventText := "Led a cross-functional team"
			eventDate := time.Now()
			mode := careerservice.TimelineJournaling

			// Capture the event with optional configurations
			err := cliEventService.CaptureEvent(
				ctx,
				eventText,
				eventDate,
				mode,
				WithCompany("TechCorp"),
				WithProject("Platform Migration"),
				WithTags([]string{"leadership", "technical"}),
			)

			// Assertions
			Expect(err).To(BeNil())
		})
	})

	Context("Listing Events", func() {
		It("should support listing events", func() {
			ctx := context.Background()

			// Create a test event directly via repository
			event := &career.CareerEvent{
				Text: "Test event",
				Date: time.Now(),
			}
			err := repo.Create(ctx, event)
			Expect(err).To(BeNil())

			// List events using CLI service
			filters := &careerrepo.ListFilters{}
			events, err := cliEventService.ListEvents(ctx, filters)

			// Assertions
			Expect(err).To(BeNil())
			Expect(events).ToNot(BeNil())
		})
	})

	Context("Get Event By ID", func() {
		It("should retrieve a specific event", func() {
			ctx := context.Background()

			// Create a test event
			event := &career.CareerEvent{
				Text: "Completed important milestone",
				Date: time.Now(),
			}
			err := careerSvc.CaptureEvent(ctx, event, careerservice.ManualEntry)
			Expect(err).To(BeNil())
			eventID := event.ID

			// Get event by ID using CLI service
			retrievedEvent, err := cliEventService.GetEventByID(ctx, eventID)

			// Assertions
			Expect(err).To(BeNil())
			Expect(retrievedEvent).ToNot(BeNil())
			Expect(retrievedEvent.ID).To(Equal(eventID))
			Expect(retrievedEvent.Text).To(Equal("Completed important milestone"))
		})
	})
})

var _ = Describe("UpdateEventMetadata", func() {
	var (
		repo   *careerrepo.MemoryRepository
		svc    *careerservice.Service
		cliSvc *CLIEventService
		ctx    context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliSvc = NewCLIEventService(svc)
	})

	It("should update event metadata without changing text or date", func() {
		// Create an event
		originalEvent := &career.CareerEvent{
			ID:   "test-event-1",
			Text: "Original text",
			Date: time.Now().Add(-24 * time.Hour),
		}
		err := repo.Create(ctx, originalEvent)
		Expect(err).To(BeNil())

		// Update metadata
		updatedEvent := &career.CareerEvent{
			ID:         "test-event-1",
			Text:       "This should be ignored",
			Date:       time.Now(), // This should be ignored
			Company:    "NewCompany",
			Project:    "NewProject",
			Tags:       []string{"technical", "leadership"},
			Categories: []string{"Technical"},
		}

		err = cliSvc.UpdateEventMetadata(ctx, updatedEvent)
		Expect(err).To(BeNil())

		// Verify metadata was updated
		retrieved, err := svc.GetEventByID(ctx, "test-event-1")
		Expect(err).To(BeNil())
		Expect(retrieved.Company).To(Equal("NewCompany"))
		Expect(retrieved.Project).To(Equal("NewProject"))
		Expect(retrieved.Tags).To(ContainElements("technical", "leadership"))

		// Verify text and date were not changed
		Expect(retrieved.Text).To(Equal("Original text"))
		Expect(retrieved.Date).To(Equal(originalEvent.Date))
	})

	It("should return error when event is nil", func() {
		err := cliSvc.UpdateEventMetadata(ctx, nil)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("event cannot be nil"))
	})

	It("should return error when event ID is empty", func() {
		event := &career.CareerEvent{
			ID:      "",
			Text:    "Test",
			Company: "Company",
		}
		err := cliSvc.UpdateEventMetadata(ctx, event)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("event ID cannot be empty"))
	})

	It("should return error when event does not exist", func() {
		event := &career.CareerEvent{
			ID:      "non-existent",
			Company: "Company",
		}
		err := cliSvc.UpdateEventMetadata(ctx, event)
		Expect(err).NotTo(BeNil())
	})
})

var _ = Describe("BulkUpdateMetadata", func() {
	var (
		repo  *careerrepo.MemoryRepository
		svc   *careerservice.Service
		cliSvc *CLIEventService
		ctx   context.Context
		eventIDs []string
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliSvc = NewCLIEventService(svc)
		ctx = context.Background()

		// Create test events
		for i := 0; i < 3; i++ {
			text := "Event " + string(rune('1'+i))
			date := time.Now().Add(-time.Duration(i) * time.Hour)
			err := cliSvc.CaptureEvent(ctx, text, date, careerservice.ManualEntry,
				WithCompany("OldCompany"),
				WithProject("OldProject"),
			)
			Expect(err).To(BeNil())
			
			// Get the created event to retrieve its ID
			events, _ := svc.ListEvents(ctx, careerrepo.ListFilters{SortBy: "date", SortOrder: "desc", Limit: 1})
			if len(events) > 0 {
				eventIDs = append(eventIDs, events[0].ID)
			}
		}
	})

	It("should implement BulkUpdateMetadata method", func() {
		// Verify method exists and is callable
		Expect(cliSvc.BulkUpdateMetadata).NotTo(BeNil())
	})

	It("should update company for multiple events", func() {
		summary, err := cliSvc.BulkUpdateMetadata(ctx, eventIDs, "NewCompany", "", nil, nil)
		Expect(err).To(BeNil())
		Expect(summary).NotTo(BeNil())
		Expect(summary.UpdatedCount).To(Equal(3))

		// Verify all events were updated
		for _, id := range eventIDs {
			event, _ := svc.GetEventByID(ctx, id)
			Expect(event.Company).To(Equal("NewCompany"))
		}
	})

	It("should update project for multiple events", func() {
		summary, err := cliSvc.BulkUpdateMetadata(ctx, eventIDs, "", "NewProject", nil, nil)
		Expect(err).To(BeNil())
		Expect(summary.UpdatedCount).To(Equal(3))

		// Verify all events were updated
		for _, id := range eventIDs {
			event, _ := svc.GetEventByID(ctx, id)
			Expect(event.Project).To(Equal("NewProject"))
		}
	})

	It("should update tags for multiple events", func() {
		tags := []string{"technical", "leadership"}
		summary, err := cliSvc.BulkUpdateMetadata(ctx, eventIDs, "", "", tags, nil)
		Expect(err).To(BeNil())
		Expect(summary.UpdatedCount).To(Equal(3))

		// Verify all events were updated
		for _, id := range eventIDs {
			event, _ := svc.GetEventByID(ctx, id)
			Expect(event.Tags).To(ContainElements(tags))
		}
	})

	It("should return error when event IDs list is empty", func() {
		_, err := cliSvc.BulkUpdateMetadata(ctx, []string{}, "NewCompany", "", nil, nil)
		Expect(err).NotTo(BeNil())
	})

	It("should return error when non-existent event ID is provided", func() {
		invalidIDs := []string{"non-existent-1", "non-existent-2"}
		_, err := cliSvc.BulkUpdateMetadata(ctx, invalidIDs, "NewCompany", "", nil, nil)
		Expect(err).NotTo(BeNil())
	})

	It("should handle mixed valid and invalid event IDs", func() {
		mixedIDs := append([]string{eventIDs[0]}, "non-existent")
		_, err := cliSvc.BulkUpdateMetadata(ctx, mixedIDs, "NewCompany", "", nil, nil)
		Expect(err).NotTo(BeNil())
		// Should not have updated any events (transaction-like behavior)
		event, _ := svc.GetEventByID(ctx, eventIDs[0])
		Expect(event.Company).To(Equal("OldCompany"))
	})

	It("should validate all events before updating any", func() {
		// Create an event with invalid data
		invalidEvent := &career.CareerEvent{
			ID:   "invalid-event",
			Text: "", // Invalid: empty text
			Date: time.Now(),
		}
		_ = repo.Create(ctx, invalidEvent)

		mixedIDs := []string{eventIDs[0], "invalid-event"}
		_, err := cliSvc.BulkUpdateMetadata(ctx, mixedIDs, "NewCompany", "", nil, nil)
		Expect(err).NotTo(BeNil())
		// Should not have updated the valid event (transaction-like)
		event, _ := svc.GetEventByID(ctx, eventIDs[0])
		Expect(event.Company).To(Equal("OldCompany"))
	})

	It("should return summary of updated events", func() {
		summary, err := cliSvc.BulkUpdateMetadata(ctx, eventIDs, "NewCompany", "NewProject", nil, nil)
		Expect(err).To(BeNil())
		Expect(summary).NotTo(BeNil())
		Expect(summary.UpdatedCount).To(Equal(3))
		Expect(summary.FailedCount).To(Equal(0))
		Expect(summary.TotalCount).To(Equal(3))
	})

	It("should preserve original text and date during bulk update", func() {
		originalEvent, _ := svc.GetEventByID(ctx, eventIDs[0])
		originalText := originalEvent.Text
		originalDate := originalEvent.Date

		_, _ = cliSvc.BulkUpdateMetadata(ctx, []string{eventIDs[0]}, "NewCompany", "NewProject", nil, nil)

		updated, _ := svc.GetEventByID(ctx, eventIDs[0])
		Expect(updated.Text).To(Equal(originalText))
		Expect(updated.Date).To(Equal(originalDate))
	})
})
