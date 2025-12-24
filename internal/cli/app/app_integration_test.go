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

var _ = Describe("CLI App Integration", func() {
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

	Context("Form Submission Integration", func() {
		It("should display event on SuccessModel after form submission", func() {
			// Navigate to capture screen
			model.currentScreen = CaptureScreen

			// Create a test event
			testEvent := &career.CareerEvent{
				Text:    "Completed important project milestone",
				Date:    time.Now().Add(-5 * 24 * time.Hour),
				Company: "TechCorp Inc",
				Project: "Platform Modernization",
				Tags:    []string{"technical", "leadership"},
			}

			// Create success model with the event
			model.successModel = models.NewSuccessModel(testEvent)
			model.currentScreen = SuccessScreen

			// Verify success model displays the event
			Expect(model.currentScreen).To(Equal(SuccessScreen))
			Expect(model.successModel).NotTo(BeNil())
			Expect(model.successModel.Event().Text).To(Equal(testEvent.Text))
		})

		It("should persist event to repository during form submission", func() {
			// Create a test event
			testEvent := &career.CareerEvent{
				Text:    "Built feature that improved performance by 40%",
				Date:    time.Now().Add(-3 * 24 * time.Hour),
				Company: "TechCorp Inc",
				Project: "Performance Optimization",
				Tags:    []string{"technical"},
			}

			// Capture the event using CLI service
			err := cliService.CaptureEvent(ctx, testEvent.Text, testEvent.Date,
				careerservice.TimelineJournaling,
				service.WithCompany(testEvent.Company),
				service.WithProject(testEvent.Project),
				service.WithTags(testEvent.Tags),
			)
			Expect(err).To(BeNil())

			// Verify event was persisted to repository
			events, err := repo.List(ctx, careerrepo.ListFilters{Limit: 100})
			Expect(err).To(BeNil())
			Expect(len(events)).To(Equal(1))
			Expect(events[0].Text).To(Equal(testEvent.Text))
			Expect(events[0].Company).To(Equal(testEvent.Company))
			Expect(events[0].Project).To(Equal(testEvent.Project))
		})

		It("should display all persisted event fields in SuccessModel", func() {
			// Create and persist event
			testEvent := &career.CareerEvent{
				Text:    "Led team through critical migration",
				Date:    time.Now().Add(-2 * 24 * time.Hour),
				Company: "TechCorp",
				Project: "Migration Project",
				Tags:    []string{"leadership", "technical"},
			}

			err := cliService.CaptureEvent(ctx, testEvent.Text, testEvent.Date,
				careerservice.TimelineJournaling,
				service.WithCompany(testEvent.Company),
				service.WithProject(testEvent.Project),
				service.WithTags(testEvent.Tags),
			)
			Expect(err).To(BeNil())

			// Retrieve the persisted event
			events, err := repo.List(ctx, careerrepo.ListFilters{Limit: 100})
			Expect(err).To(BeNil())
			Expect(len(events)).To(Equal(1))

			// Create SuccessModel with persisted event
			successModel := models.NewSuccessModel(events[0])
			view := successModel.View()

			// Verify all fields are displayed
			Expect(view).To(ContainSubstring("Led team through critical migration"))
			Expect(view).To(ContainSubstring("TechCorp"))
			Expect(view).To(ContainSubstring("Migration Project"))
			Expect(view).To(ContainSubstring("leadership"))
			Expect(view).To(ContainSubstring("technical"))
		})
	})

	Context("SuccessScreen Navigation", func() {
		It("should reset form and return to capture screen via CaptureAnotherMsg", func() {
			// Simulate successful submission and transition to success screen
			model.currentScreen = SuccessScreen
			model.successModel = models.NewSuccessModel(&career.CareerEvent{
				Text: "Test event",
				Date: time.Now(),
			})

			// Send CaptureAnotherMsg
			var captureMsg tea.Msg = models.CaptureAnotherMsg{}
			newModel, _ := model.Update(captureMsg)
			updatedModel := newModel.(*Model)

			// Verify state changes
			Expect(updatedModel.currentScreen).To(Equal(CaptureScreen))
			Expect(updatedModel.formModel).NotTo(BeNil())
			Expect(updatedModel.successModel).To(BeNil())
		})

		It("should navigate to list screen via ViewRecentMsg", func() {
			// Simulate successful submission and transition to success screen
			model.currentScreen = SuccessScreen
			model.successModel = models.NewSuccessModel(&career.CareerEvent{
				Text: "Test event",
				Date: time.Now(),
			})

			// Send ViewRecentMsg
			var viewMsg tea.Msg = models.ViewRecentMsg{}
			newModel, _ := model.Update(viewMsg)
			updatedModel := newModel.(*Model)

			// Verify state changes
			Expect(updatedModel.currentScreen).To(Equal(ListScreen))
			Expect(updatedModel.successModel).To(BeNil())
		})

		It("should handle exit from success screen", func() {
			// Simulate successful submission
			model.currentScreen = SuccessScreen
			model.successModel = models.NewSuccessModel(&career.CareerEvent{
				Text: "Test event",
				Date: time.Now(),
			})

			// Send quit message
			var quitMsg tea.Msg = tea.KeyMsg{Type: tea.KeyCtrlC}
			newModel, cmd := model.Update(quitMsg)

			// Verify quit command is returned
			Expect(cmd).NotTo(BeNil())
			Expect(newModel).NotTo(BeNil())
		})
	})

	Context("Complete Capture Workflow", func() {
		It("should handle full event capture workflow from capture to success", func() {
			// Start on home screen
			Expect(model.currentScreen).To(Equal(HomeScreen))

			// Navigate to capture screen
			captureMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
			newModel, _ := model.Update(captureMsg)
			model = newModel.(*Model)
			Expect(model.currentScreen).To(Equal(CaptureScreen))

			// Verify form is ready
			Expect(model.formModel).NotTo(BeNil())
			Expect(model.successModel).To(BeNil())

			// Create and persist an event through CLI service
			testEvent := &career.CareerEvent{
				Text:    "Architected new microservices platform",
				Date:    time.Now().Add(-10 * 24 * time.Hour),
				Company: "TechCorp",
				Project: "Platform Redesign",
				Tags:    []string{"technical", "technical"},
			}

			err := cliService.CaptureEvent(ctx, testEvent.Text, testEvent.Date,
				careerservice.TimelineJournaling,
				service.WithCompany(testEvent.Company),
				service.WithProject(testEvent.Project),
				service.WithTags(testEvent.Tags),
			)
			Expect(err).To(BeNil())

			// Retrieve the event and simulate form submission
			events, err := repo.List(ctx, careerrepo.ListFilters{Limit: 100})
			Expect(err).To(BeNil())
			Expect(len(events)).To(Equal(1))

			// Create success screen with the event
			model.successModel = models.NewSuccessModel(events[0])
			model.currentScreen = SuccessScreen

			// Verify success screen displays the event
			Expect(model.currentScreen).To(Equal(SuccessScreen))
			Expect(model.successModel.Event().Text).To(Equal(testEvent.Text))

			view := model.View()
			Expect(view).To(ContainSubstring("Architected new microservices platform"))
			Expect(view).To(ContainSubstring("TechCorp"))

			// Test navigation to capture another
			anotherMsg := models.CaptureAnotherMsg{}
			newModel, _ = model.Update(anotherMsg)
			model = newModel.(*Model)

			Expect(model.currentScreen).To(Equal(CaptureScreen))
			Expect(model.formModel).NotTo(BeNil())
			Expect(model.successModel).To(BeNil())
		})

		It("should maintain event data through capture and display workflow", func() {
			// Create multiple events
			eventTexts := []string{
				"Led team through critical migration",
				"Architected microservices platform",
				"Mentored junior developers",
			}

			for _, text := range eventTexts {
				err := cliService.CaptureEvent(ctx, text, time.Now().Add(-24*time.Hour),
					careerservice.TimelineJournaling,
					service.WithCompany("TechCorp"),
				)
				Expect(err).To(BeNil())
			}

			// Retrieve all events
			events, err := repo.List(ctx, careerrepo.ListFilters{Limit: 100})
			Expect(err).To(BeNil())
			Expect(len(events)).To(Equal(3))

			// Verify each event can be displayed
			for i, event := range events {
				successModel := models.NewSuccessModel(event)
				view := successModel.View()
				Expect(view).To(ContainSubstring(eventTexts[i]))
			}
		})
	})

	Context("Event Persistence Across Sessions", func() {
		It("should retrieve events persisted from previous session", func() {
			// Simulate events from previous session
			oldEvents := []string{
				"Event from yesterday",
				"Event from last week",
			}

			for _, text := range oldEvents {
				err := cliService.CaptureEvent(ctx, text, time.Now().Add(-24*time.Hour),
					careerservice.ManualEntry,
				)
				Expect(err).To(BeNil())
			}

			// Create new model (simulating new session)
			newModel := NewModel(cliService, svc)

			// Verify can still retrieve events
			events, err := repo.List(ctx, careerrepo.ListFilters{Limit: 100})
			Expect(err).To(BeNil())
			Expect(len(events)).To(Equal(2))

			// Verify new model can reference the same events
			Expect(newModel).NotTo(BeNil())
		})
	})
})
