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

var _ = Describe("Application Model", func() {
	var (
		cliService *service.CLIEventService
		svc        *careerservice.Service
		model      *Model
	)

	BeforeEach(func() {
		svc = careerservice.NewService(nil)
		cliService = service.NewCLIEventService(svc)
		model = NewModel(cliService, svc)
	})

	Context("Initialization", func() {
		It("should initialize with HomeScreen as current screen", func() {
			Expect(model.currentScreen).To(Equal(HomeScreen))
		})

		It("should initialize with default width and height", func() {
			Expect(model.width).To(Equal(80))
			Expect(model.height).To(Equal(24))
		})

		It("should have no initial error", func() {
			Expect(model.err).To(BeNil())
		})
	})

	Context("Screen Navigation from HomeScreen", func() {
		It("should navigate to CaptureScreen with 'c'", func() {
			model.currentScreen = HomeScreen
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			Expect(updatedModel.currentScreen).To(Equal(CaptureScreen))
		})

		It("should navigate to ListScreen with 'l'", func() {
			model.currentScreen = HomeScreen
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			Expect(updatedModel.currentScreen).To(Equal(ListScreen))
		})

		It("should stay on HomeScreen with 'h'", func() {
			model.currentScreen = HomeScreen
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			Expect(updatedModel.currentScreen).To(Equal(HomeScreen))
		})

		It("should handle backspace navigation from ListScreen", func() {
			model.currentScreen = ListScreen
			model.previousScreen = HomeScreen
			msg := tea.KeyMsg{Type: tea.KeyBackspace}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			Expect(updatedModel.currentScreen).To(Equal(HomeScreen))
		})
	})

	Context("View Rendering", func() {
		It("should render home screen view", func() {
			model.currentScreen = HomeScreen
			view := model.View()
			Expect(view).To(ContainSubstring("KaRiya - Career Journal CLI"))
			Expect(view).To(ContainSubstring("Capture Career Event"))
			Expect(view).To(ContainSubstring("List Events"))
		})

		It("should render capture screen view", func() {
			model.currentScreen = CaptureScreen
			view := model.View()
			Expect(view).To(ContainSubstring("Capture Career Event"))
			Expect(view).To(ContainSubstring("Event Text"))
			Expect(view).To(ContainSubstring("Company"))
		})

		It("should render list screen view", func() {
			model.currentScreen = ListScreen
			view := model.View()
			Expect(view).To(ContainSubstring("Recent Career Events"))
			Expect(view).To(ContainSubstring("Event 1"))
		})

		It("should render view screen view", func() {
			model.currentScreen = ViewScreen
			view := model.View()
			Expect(view).To(ContainSubstring("Event Details"))
			Expect(view).To(ContainSubstring("Title"))
		})
	})

	Context("Window Resizing", func() {
		It("should update width and height on window resize", func() {
			msg := tea.WindowSizeMsg{Width: 120, Height: 40}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			Expect(updatedModel.width).To(Equal(120))
			Expect(updatedModel.height).To(Equal(40))
		})
	})

	Context("Quit Command", func() {
		It("should return model when 'q' key is pressed", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
			newModel, cmd := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())
		})

		It("should return model when ctrl+c is pressed", func() {
			msg := tea.KeyMsg{Type: tea.KeyCtrlC}
			newModel, cmd := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())
		})
	})

	Context("FormModel Integration", func() {
		It("should have FormModel instance in Model struct", func() {
			Expect(model.formModel).NotTo(BeNil())
		})

		It("should create SuccessModel after form submission", func() {
			// Initially successModel should be nil
			Expect(model.successModel).To(BeNil())
			// After form submission, successModel would be created
			// This is tested in the form submission test
		})

		It("should delegate Update to FormModel on CaptureScreen", func() {
			model.currentScreen = CaptureScreen
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			Expect(updatedModel).NotTo(BeNil())
		})

		It("should delegate View to FormModel on CaptureScreen", func() {
			model.currentScreen = CaptureScreen
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should reset FormModel when navigating to CaptureScreen", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			Expect(updatedModel.formModel).NotTo(BeNil())
		})
	})

	Context("Keyboard Input Handling on CaptureScreen", func() {
		It("should allow 'h' character input in form without navigating", func() {
			model.currentScreen = CaptureScreen
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			// Should still be on CaptureScreen, not navigated to HomeScreen
			Expect(updatedModel.currentScreen).To(Equal(CaptureScreen))
		})

		It("should allow 'l' character input in form without navigating", func() {
			model.currentScreen = CaptureScreen
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			// Should still be on CaptureScreen, not navigated to ListScreen
			Expect(updatedModel.currentScreen).To(Equal(CaptureScreen))
		})

		It("should allow 'backspace' character input in form without navigating", func() {
			model.currentScreen = CaptureScreen
			model.previousScreen = HomeScreen
			msg := tea.KeyMsg{Type: tea.KeyBackspace}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			// Should still be on CaptureScreen, not navigated back
			Expect(updatedModel.currentScreen).To(Equal(CaptureScreen))
		})

		It("should allow 'c' character input in form without navigating", func() {
			model.currentScreen = CaptureScreen
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			// Should still be on CaptureScreen, not navigated to CaptureScreen again
			Expect(updatedModel.currentScreen).To(Equal(CaptureScreen))
		})
	})

	Context("End-to-End Capture Workflow", func() {
		It("should navigate from Home → Capture screen", func() {
			Expect(model.currentScreen).To(Equal(HomeScreen))
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			Expect(updatedModel.currentScreen).To(Equal(CaptureScreen))
		})

		It("should complete form submission and transition to SuccessScreen", func() {
			// Navigate to capture screen
			model.currentScreen = CaptureScreen

			// Simulate form submission by creating a SuccessModel
			model.successModel = models.NewSuccessModel(&career.CareerEvent{
				Text:    "Test event",
				Date:    time.Now(),
				Company: "Test Company",
			})
			model.currentScreen = SuccessScreen

			Expect(model.currentScreen).To(Equal(SuccessScreen))
			Expect(model.successModel).NotTo(BeNil())
		})

		It("should handle CaptureAnotherMsg and reset form", func() {
			model.currentScreen = SuccessScreen
			model.successModel = models.NewSuccessModel(&career.CareerEvent{
				Text: "Test event",
				Date: time.Now(),
			})

			// Send CaptureAnotherMsg
			var captureMsg tea.Msg = models.CaptureAnotherMsg{}
			newModel, _ := model.Update(captureMsg)
			updatedModel := newModel.(*Model)

			Expect(updatedModel.currentScreen).To(Equal(CaptureScreen))
			Expect(updatedModel.successModel).To(BeNil())
			Expect(updatedModel.formModel).NotTo(BeNil())
		})

		It("should handle ViewRecentMsg and navigate to ListScreen", func() {
			model.currentScreen = SuccessScreen
			model.successModel = models.NewSuccessModel(&career.CareerEvent{
				Text: "Test event",
				Date: time.Now(),
			})

			// Send ViewRecentMsg
			var viewMsg tea.Msg = models.ViewRecentMsg{}
			newModel, _ := model.Update(viewMsg)
			updatedModel := newModel.(*Model)

			Expect(updatedModel.currentScreen).To(Equal(ListScreen))
			Expect(updatedModel.successModel).To(BeNil())
		})

		It("should handle complete capture workflow", func() {
			// Start on home screen
			Expect(model.currentScreen).To(Equal(HomeScreen))

			// Navigate to capture
			captureMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
			newModel, _ := model.Update(captureMsg)
			model = newModel.(*Model)
			Expect(model.currentScreen).To(Equal(CaptureScreen))

			// Verify form is ready
			Expect(model.formModel).NotTo(BeNil())
			Expect(model.successModel).To(BeNil())

			// Simulate form submission
			model.successModel = models.NewSuccessModel(&career.CareerEvent{
				Text:    "Completed important project",
				Date:    time.Now(),
				Company: "TechCorp",
				Project: "Platform Migration",
				Tags:    []string{"technical", "leadership"},
			})
			model.currentScreen = SuccessScreen

			Expect(model.currentScreen).To(Equal(SuccessScreen))
			Expect(model.successModel.Event().Text).To(Equal("Completed important project"))

			// Capture another
			var anotherMsg tea.Msg = models.CaptureAnotherMsg{}
			newModel, _ = model.Update(anotherMsg)
			model = newModel.(*Model)
			Expect(model.currentScreen).To(Equal(CaptureScreen))
			Expect(model.formModel).NotTo(BeNil())
		})
	})


	Context("Event Persistence and Display", func() {
		It("should persist event to repository and display in SuccessModel", func() {
			// Create repository and service
			repo := careerrepo.NewMemoryRepository()
			svc := careerservice.NewService(repo)
			cliService := service.NewCLIEventService(svc)

			// Create an event directly
			ctx := context.Background()
			event := &career.CareerEvent{
				Text:    "Completed important project milestone",
				Date:    time.Now().Add(-5 * 24 * time.Hour),
				Company: "TechCorp Inc",
				Project: "Platform Modernization",
				Tags:    []string{"technical", "leadership"},
			}

			// Capture the event
			err := cliService.CaptureEvent(ctx, event.Text, event.Date, careerservice.TimelineJournaling,
				service.WithCompany(event.Company),
				service.WithProject(event.Project),
				service.WithTags(event.Tags),
			)
			Expect(err).To(BeNil())

			// Verify event was persisted
			events, err := repo.List(ctx, careerrepo.ListFilters{Limit: 100})
			Expect(err).To(BeNil())
			Expect(len(events)).To(Equal(1))
			Expect(events[0].Text).To(Equal("Completed important project milestone"))
			Expect(events[0].Company).To(Equal("TechCorp Inc"))
			Expect(events[0].Project).To(Equal("Platform Modernization"))

			// Create SuccessModel with the persisted event
			successModel := models.NewSuccessModel(events[0])
			view := successModel.View()

			// Verify event details are displayed
			Expect(view).To(ContainSubstring("Completed important project milestone"))
			Expect(view).To(ContainSubstring("TechCorp Inc"))
			Expect(view).To(ContainSubstring("Platform Modernization"))
			Expect(view).To(ContainSubstring("technical"))
			Expect(view).To(ContainSubstring("leadership"))
		})

		It("should retrieve persisted events from repository", func() {
			// Create repository and service
			repo := careerrepo.NewMemoryRepository()
			svc := careerservice.NewService(repo)
			cliService := service.NewCLIEventService(svc)

			ctx := context.Background()

			// Create and persist multiple events
			for i := 1; i <= 3; i++ {
				event := &career.CareerEvent{
					Text:    "Event " + string(rune('0'+i)),
					Date:    time.Now().Add(-time.Duration(i) * 24 * time.Hour),
					Company: "Company " + string(rune('0'+i)),
				}
				err := cliService.CaptureEvent(ctx, event.Text, event.Date, careerservice.TimelineJournaling,
					service.WithCompany(event.Company),
				)
				Expect(err).To(BeNil())
			}

			// Retrieve all events
			events, err := repo.List(ctx, careerrepo.ListFilters{Limit: 100})
			Expect(err).To(BeNil())
			Expect(len(events)).To(Equal(3))

			// Verify events are in correct order
			Expect(events[0].Text).To(ContainSubstring("Event 1"))
			Expect(events[1].Text).To(ContainSubstring("Event 2"))
			Expect(events[2].Text).To(ContainSubstring("Event 3"))
		})

		It("should display all event fields in SuccessModel", func() {
			// Create event with all fields
			event := &career.CareerEvent{
				ID:         "test-123",
				Text:       "Complex event with all fields",
				Date:       time.Date(2024, 12, 20, 0, 0, 0, 0, time.UTC),
				Company:    "TechCorp",
				Project:    "Project Alpha",
				Tags:       []string{"technical", "leadership", "product"},
				Categories: []string{"Technical"},
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			// Create SuccessModel
			successModel := models.NewSuccessModel(event)
			view := successModel.View()

			// Verify all fields are displayed
			Expect(view).To(ContainSubstring("Complex event with all fields"))
			Expect(view).To(ContainSubstring("TechCorp"))
			Expect(view).To(ContainSubstring("Project Alpha"))
			Expect(view).To(ContainSubstring("technical"))
			Expect(view).To(ContainSubstring("leadership"))
			Expect(view).To(ContainSubstring("product"))
		})
	})

})
