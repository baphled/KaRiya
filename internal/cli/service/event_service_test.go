package service

import (
	"context"
	"time"

	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLI Event Service", func() {
	var (
		cliEventService *CLIEventService
		careerSvc       *careerservice.Service
		repo            careerrepo.EventRepository
	)

	BeforeEach(func() {
		// Set up in-memory repository for testing
		repo = careermemory.NewEventRepository()
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
			event := fixtures.Event("test-list-1")
			err := repo.Create(ctx, event)
			Expect(err).To(BeNil())

			// List events using CLI service
			filters := fixtures.EventListFilters()
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
			event := fixtures.EventWith("", "Completed important milestone", "", "")
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
		repo   *careermemory.EventRepository
		svc    *careerservice.Service
		cliSvc *CLIEventService
		ctx    context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = careermemory.NewEventRepository()
		svc = careerservice.NewService(repo)
		cliSvc = NewCLIEventService(svc)
	})

	It("should update event metadata without changing text or date", func() {
		// Create an event
		originalEvent := fixtures.EventWith("test-event-1", "Original text", "", "")
		originalEvent.Date = time.Now().Add(-24 * time.Hour)
		err := repo.Create(ctx, originalEvent)
		Expect(err).To(BeNil())

		// Update metadata
		updatedEvent := fixtures.EventWith("test-event-1", "This should be ignored", "NewCompany", "NewProject")
		updatedEvent.Tags = []string{"technical", "leadership"}
		updatedEvent.Categories = []string{"Technical"}

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
		event := fixtures.EventWith("", "Test", "Company", "")
		err := cliSvc.UpdateEventMetadata(ctx, event)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("event ID cannot be empty"))
	})

	It("should return error when event does not exist", func() {
		event := fixtures.EventWith("non-existent", "", "Company", "")
		err := cliSvc.UpdateEventMetadata(ctx, event)
		Expect(err).NotTo(BeNil())
	})
})
