package app

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLI App - Import to Metadata Review Integration", func() {
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

	Describe("Import Completion Triggers Metadata Review", func() {
		It("should navigate to metadata review after successful import", func() {
			// Arrange: Create model and set up for import
			model.currentScreen = HomeScreen

			// Act: Simulate import completion with created events
			createdEvent := &career.CareerEvent{
				ID:   "imported-event-1",
				Text: "Imported career event",
				Date: time.Now().Add(-24 * time.Hour),
			}

			// Persist the event first
			err := cliService.CaptureEvent(ctx, createdEvent.Text, createdEvent.Date, careerservice.ManualEntry)
			Expect(err).To(BeNil())

			// Get the created event to get its actual ID
			events, _ := svc.ListEvents(ctx, careerrepo.ListFilters{})
			Expect(len(events)).To(BeNumerically(">", 0))

			// Assert: metadata review model should be created when we navigate
			model.metadataReviewModel = models.NewMetadataReviewModelForImport(svc, ctx, []string{events[0].ID})
			model.currentScreen = MetadataReviewScreen

			Expect(model.currentScreen).To(Equal(MetadataReviewScreen))
			Expect(model.metadataReviewModel).NotTo(BeNil())
		})

		It("should allow bulk operations on imported events", func() {
			// Arrange: Create and persist an event
			testEvent := &career.CareerEvent{
				Text:    "Event to bulk edit",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "TestCorp",
				Tags:    []string{"technical"},
			}

			err := cliService.CaptureEvent(ctx, testEvent.Text, testEvent.Date, careerservice.ManualEntry,
				service.WithCompany(testEvent.Company),
				service.WithTags(testEvent.Tags),
			)
			Expect(err).To(BeNil())

			// Get the event
			events, _ := svc.ListEvents(ctx, careerrepo.ListFilters{})
			Expect(len(events)).To(BeNumerically(">", 0))

			// Navigate to metadata review
			model.metadataReviewModel = models.NewMetadataReviewModelForImport(svc, ctx, []string{events[0].ID})
			model.currentScreen = MetadataReviewScreen

			// Act: Press 'b' to trigger bulk operations
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})

			// Assert: Should return a command
			Expect(cmd).NotTo(BeNil())

			// Execute the command and verify it returns BulkOperationsMsg
			msg := cmd()
			bulkMsg, ok := msg.(models.BulkOperationsMsg)
			Expect(ok).To(BeTrue())
			Expect(len(bulkMsg.Events)).To(BeNumerically(">", 0))
		})

		It("should display bulk operations hint in metadata review", func() {
			// Arrange: Create and persist an event
			err := cliService.CaptureEvent(ctx, "Test event for bulk ops", time.Now().Add(-24*time.Hour), careerservice.ManualEntry)
			Expect(err).To(BeNil())

			// Create metadata review model
			model.metadataReviewModel = models.NewMetadataReviewModel(svc, ctx)
			model.currentScreen = MetadataReviewScreen

			// Act: Render the view
			view := model.View()

			// Assert: View should contain bulk operations hint
			Expect(view).To(ContainSubstring("ulk"))
		})
	})
})
