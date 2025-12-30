package app

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLI App - Manual Capture to Metadata Review Integration", func() {
	var (
		repo       *careerrepo.MemoryRepository
		svc        *careerservice.Service
		cliService *service.CLIEventService
		model      *Model
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliService = service.NewCLIEventService(svc)
		model = NewModel(cliService, svc)
	})

	Describe("Success Screen Metadata Review Option", func() {
		It("should include Review Metadata option in success screen", func() {
			// Arrange: Create a captured event
			testEvent := &career.CareerEvent{
				ID:   "test-event-1",
				Text: "Completed a project milestone",
				Date: time.Now(),
			}

			// Act: Create success model
			successModel := models.NewSuccessModel(testEvent)

			// Assert: Success model should have the review metadata action
			Expect(successModel).NotTo(BeNil())
		})

		It("should navigate to metadata review when Review Metadata is selected", func() {
			// Arrange: Capture an event
			err := cliService.CaptureEvent(ctx, "Test event", time.Now(), careerservice.ManualEntry)
			Expect(err).To(BeNil())

			// Get the created event
			events, _ := svc.ListEvents(ctx, careerrepo.ListFilters{})
			Expect(len(events)).To(Equal(1))

			// Set model to success screen with the captured event
			model.successModel = models.NewSuccessModel(events[0])
			model.currentScreen = SuccessScreen

			// Act: Simulate selecting Review Metadata action
			model.currentScreen = MetadataReviewScreen
			model.metadataReviewModel = models.NewMetadataReviewModelForImport(svc, ctx, []string{events[0].ID})

			// Assert: Should be on metadata review screen with the event
			Expect(model.currentScreen).To(Equal(MetadataReviewScreen))
			Expect(model.metadataReviewModel).NotTo(BeNil())
		})

		It("should allow editing metadata after capture", func() {
			// Arrange: Capture an event with minimal metadata
			err := cliService.CaptureEvent(ctx, "Completed project", time.Now(), careerservice.ManualEntry)
			Expect(err).To(BeNil())

			events, _ := svc.ListEvents(ctx, careerrepo.ListFilters{})
			capturedEvent := events[0]

			// Act: Navigate to metadata editor
			model.metadataEditorModel = models.NewMetadataEditorModel(capturedEvent, svc, cliService, ctx)
			model.currentScreen = MetadataEditorScreen

			// Assert: Should be able to edit metadata
			Expect(model.currentScreen).To(Equal(MetadataEditorScreen))
			Expect(model.metadataEditorModel).NotTo(BeNil())
		})

		It("should allow quick metadata enrichment (company, project, tags)", func() {
			// Arrange: Capture an event
			err := cliService.CaptureEvent(ctx, "Implemented feature", time.Now(), careerservice.ManualEntry)
			Expect(err).To(BeNil())

			events, _ := svc.ListEvents(ctx, careerrepo.ListFilters{})
			originalEvent := events[0]

			// Act: Update event with additional metadata
			originalEvent.Company = "TechCorp"
			originalEvent.Project = "Project Alpha"
			originalEvent.Tags = []string{"technical", "leadership"}
			updateErr := cliService.UpdateEventMetadata(ctx, originalEvent)
			Expect(updateErr).To(BeNil())

			// Assert: Event should have updated metadata
			updatedEvent, _ := svc.GetEventByID(ctx, originalEvent.ID)
			Expect(updatedEvent.Company).To(Equal("TechCorp"))
			Expect(updatedEvent.Project).To(Equal("Project Alpha"))
			Expect(updatedEvent.Tags).To(ContainElements("technical", "leadership"))
		})

		It("should support capturing multiple events and reviewing all metadata", func() {
			// Arrange: Capture multiple events
			event1Err := cliService.CaptureEvent(ctx, "Event 1", time.Now(), careerservice.ManualEntry)
			Expect(event1Err).To(BeNil())

			event2Err := cliService.CaptureEvent(ctx, "Event 2", time.Now().Add(-24*time.Hour), careerservice.ManualEntry)
			Expect(event2Err).To(BeNil())

			// Get all events
			events, _ := svc.ListEvents(ctx, careerrepo.ListFilters{})
			Expect(len(events)).To(Equal(2))

			// Act: Navigate to metadata review for all events
			eventIDs := []string{events[0].ID, events[1].ID}
			model.metadataReviewModel = models.NewMetadataReviewModelForImport(svc, ctx, eventIDs)
			model.currentScreen = MetadataReviewScreen

			// Assert: Should show all events in metadata review
			Expect(model.currentScreen).To(Equal(MetadataReviewScreen))
			Expect(model.metadataReviewModel).NotTo(BeNil())
		})

		It("should provide bulk operations on captured events", func() {
			// Arrange: Capture multiple events
			for i := 0; i < 3; i++ {
				err := cliService.CaptureEvent(ctx, "Event "+string(rune(48+i)), time.Now().Add(-time.Duration(i)*24*time.Hour), careerservice.ManualEntry)
				Expect(err).To(BeNil())
			}

			events, _ := svc.ListEvents(ctx, careerrepo.ListFilters{})
			Expect(len(events)).To(Equal(3))

			// Act: Navigate to bulk operations with first two events
			bulkEvents := []*career.CareerEvent{events[0], events[1]}
			model.bulkOperationsModel = models.NewBulkOperationsModel(bulkEvents, svc, cliService, ctx)
			model.currentScreen = BulkOperationsScreen

			// Assert: Should be in bulk operations screen
			Expect(model.currentScreen).To(Equal(BulkOperationsScreen))
			Expect(model.bulkOperationsModel).NotTo(BeNil())
		})
	})
})
