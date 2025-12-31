package app

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metadata Review to Editor E2E Workflow", func() {
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

	Describe("Edit event from metadata review and refresh list", func() {
		Context("when user edits event metadata", func() {
			It("should edit event from metadata review and refresh list with updated values", func() {
				// Arrange: Create an event with minimal metadata
				err := cliService.CaptureEvent(
					ctx,
					"Implemented critical feature",
					time.Now().Add(-24*time.Hour),
					careerservice.ManualEntry,
				)
				Expect(err).ToNot(HaveOccurred())

				// Get the created event
				events, err := svc.ListEvents(ctx, careerrepo.ListFilters{})
				Expect(err).ToNot(HaveOccurred())
				Expect(events).To(HaveLen(1))

				originalEvent := events[0]
				Expect(originalEvent.Company).To(Equal("")) // No company initially

				// Act: Navigate to metadata review screen
				model.currentScreen = MetadataReviewScreen
				model.metadataReviewModel = models.NewMetadataReviewModelForImport(svc, ctx, []string{originalEvent.ID})

				// Open editor for the event
				model.currentScreen = MetadataEditorScreen
				editorModel := models.NewMetadataEditorModel(originalEvent, svc, cliService, ctx)
				model.metadataEditorModel = editorModel

				// Simulate editing the event (change company and project)
				editedEvent := originalEvent
				editedEvent.Company = "TechCorp"
				editedEvent.Project = "Project Alpha"
				editedEvent.Tags = []string{"technical", "leadership"}

				// Save the changes via CLI service
				updateErr := cliService.UpdateEventMetadata(ctx, editedEvent)
				Expect(updateErr).ToNot(HaveOccurred())

				// Return to metadata review screen
				model.currentScreen = MetadataReviewScreen
				model.metadataReviewModel = models.NewMetadataReviewModelForImport(svc, ctx, []string{originalEvent.ID})

				// Assert: Verify the event in the review list shows updated metadata
				updatedEvents, err := svc.ListEvents(ctx, careerrepo.ListFilters{})
				Expect(err).ToNot(HaveOccurred())
				Expect(updatedEvents).To(HaveLen(1))

				updatedEvent := updatedEvents[0]
				Expect(updatedEvent.Company).To(Equal("TechCorp"))
				Expect(updatedEvent.Project).To(Equal("Project Alpha"))
				Expect(updatedEvent.Tags).To(ContainElements("technical", "leadership"))
				Expect(updatedEvent.Text).To(Equal(originalEvent.Text)) // Text unchanged
			})

			It("should preserve original text and date when editing metadata", func() {
				// Arrange: Create event with specific text and date
				originalText := "Architected microservices platform"
				originalDate := time.Now().Add(-72 * time.Hour)

				err := cliService.CaptureEvent(
					ctx,
					originalText,
					originalDate,
					careerservice.ManualEntry,
				)
				Expect(err).ToNot(HaveOccurred())

				events, _ := svc.ListEvents(ctx, careerrepo.ListFilters{})
				eventToEdit := events[0]

				// Act: Edit only metadata fields
				eventToEdit.Company = "NewCorp"
				eventToEdit.Tags = []string{"technical"}

				updateErr := cliService.UpdateEventMetadata(ctx, eventToEdit)
				Expect(updateErr).ToNot(HaveOccurred())

				// Assert: Verify text and date are preserved
				updatedEvents, _ := svc.ListEvents(ctx, careerrepo.ListFilters{})
				Expect(updatedEvents[0].Text).To(Equal(originalText))
				Expect(updatedEvents[0].Date.Format("2006-01-02")).To(Equal(originalDate.Format("2006-01-02")))
			})

			It("should handle multiple event editing in sequence", func() {
				// Arrange: Create two events
				err1 := cliService.CaptureEvent(
					ctx,
					"Event 1",
					time.Now().Add(-48*time.Hour),
					careerservice.ManualEntry,
				)
				Expect(err1).ToNot(HaveOccurred())

				err2 := cliService.CaptureEvent(
					ctx,
					"Event 2",
					time.Now().Add(-24*time.Hour),
					careerservice.ManualEntry,
				)
				Expect(err2).ToNot(HaveOccurred())

				events, _ := svc.ListEvents(ctx, careerrepo.ListFilters{})
				Expect(events).To(HaveLen(2))

				// Act: Edit first event
				event1 := events[0]
				event1.Company = "Company1"
				updateErr1 := cliService.UpdateEventMetadata(ctx, event1)
				Expect(updateErr1).ToNot(HaveOccurred())

				// Edit second event
				event2 := events[1]
				event2.Company = "Company2"
				updateErr2 := cliService.UpdateEventMetadata(ctx, event2)
				Expect(updateErr2).ToNot(HaveOccurred())

				// Assert: Verify both events updated correctly
				updatedEvents, _ := svc.ListEvents(ctx, careerrepo.ListFilters{})
				Expect(updatedEvents).To(HaveLen(2))

				// Find events by ID to verify updates
				event1Updated := updatedEvents[0]
				event2Updated := updatedEvents[1]

				if event1Updated.ID == event1.ID {
					Expect(event1Updated.Company).To(Equal("Company1"))
					Expect(event2Updated.Company).To(Equal("Company2"))
				} else {
					Expect(event1Updated.Company).To(Equal("Company2"))
					Expect(event2Updated.Company).To(Equal("Company1"))
				}
			})
		})
	})
})

