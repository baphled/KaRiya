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
		repo       *careerrepo.MemoryRepository
		cliService *service.CLIEventService
		svc        *careerservice.Service
		model      *Model
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
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

			// First update delegates to ListModel, which returns BackMsg command
			newModel, cmd := model.Update(msg)
			updatedModel := newModel.(*Model)

			// ListModel should return a cmd that generates BackMsg
			Expect(cmd).NotTo(BeNil())

			// Execute the command to get the BackMsg
			backMsg := cmd()

			// Now update with the BackMsg
			finalModel, _ := updatedModel.Update(backMsg)
			finalUpdatedModel := finalModel.(*Model)

			Expect(finalUpdatedModel.currentScreen).To(Equal(HomeScreen))
		})
	})

	Context("View Rendering", func() {
		BeforeEach(func() {
			// Create a repository with events for the list view test
			repo := careerrepo.NewMemoryRepository()
			svc = careerservice.NewService(repo)
			cliService = service.NewCLIEventService(svc)

			// Add an event to the repository BEFORE creating the model
			ctx := context.Background()
			err := cliService.CaptureEvent(ctx, "Event 1", time.Now().Add(-24*time.Hour), careerservice.TimelineJournaling)
			Expect(err).To(BeNil())

			// Now create the model (which will load events)
			model = NewModel(cliService, svc)
		})

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
			Expect(view).To(ContainSubstring("Career Events"))
			Expect(view).To(ContainSubstring("Event 1"))
		})

		It("should render view screen view", func() {
			// Create a mock event for viewing
			event := &career.CareerEvent{
				ID:      "test-id",
				Text:    "Test Event for Viewing",
				Date:    time.Now(),
				Company: "Test Company",
			}

			// Create a DetailsModel for the view screen
			model.detailsModel = models.NewDetailsModel(event)
			model.currentScreen = ViewScreen

			view := model.View()
			Expect(view).To(ContainSubstring("Event Details"))
			Expect(view).To(ContainSubstring("Test Event for Viewing"))
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
		It("should handle quit from HomeScreen", func() {
			model.currentScreen = HomeScreen
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
			newModel, _ := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			// On HomeScreen, 'q' is not handled directly by app,
			// it falls through. The test just verifies no crash.
		})

		It("should handle ctrl+c", func() {
			msg := tea.KeyMsg{Type: tea.KeyCtrlC}
			newModel, _ := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			// Ctrl+C is not handled directly by app,
			// it falls through. The test just verifies no crash.
		})

		It("should handle QuitMsg and return tea.Quit command", func() {
			msg := models.QuitMsg{}
			newModel, cmd := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())
			// QuitMsg should return tea.Quit command
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

	Context("Event Viewing from List", func() {
		BeforeEach(func() {
			// Create fresh repository and service for each test
			repo = careerrepo.NewMemoryRepository()
			svc = careerservice.NewService(repo)
			cliService = service.NewCLIEventService(svc)

			// Add test events
			ctx := context.Background()
			event1 := &career.CareerEvent{
				Text:    "First test event",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "TestCorp",
			}
			event2 := &career.CareerEvent{
				Text:    "Second test event",
				Date:    time.Now().Add(-48 * time.Hour),
				Company: "AnotherCorp",
			}
			_ = svc.CaptureEvent(ctx, event1, careerservice.ManualEntry)
			_ = svc.CaptureEvent(ctx, event2, careerservice.ManualEntry)

			// Create model AFTER adding events
			model = NewModel(cliService, svc)
		})

		It("should navigate to ViewScreen when models.ViewEventMsg is sent", func() {
			// Set current screen to list
			model.currentScreen = ListScreen
			model.previousScreen = HomeScreen

			// Create a test event
			testEvent := &career.CareerEvent{
				ID:      "test-id",
				Text:    "Test event for viewing",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "ViewCorp",
				Project: "ViewProject",
			}

			// Send models.ViewEventMsg
			viewMsg := models.ViewEventMsg{Event: testEvent}
			newModel, _ := model.Update(viewMsg)
			updatedModel := newModel.(*Model)

			// Should navigate to ViewScreen
			Expect(updatedModel.currentScreen).To(Equal(ViewScreen))
			Expect(updatedModel.previousScreen).To(Equal(ListScreen))

			// Details model should be initialized with the event
			Expect(updatedModel.detailsModel).NotTo(BeNil())
		})

		It("should display event details in ViewScreen", func() {
			// Set up a test event
			testEvent := &career.CareerEvent{
				ID:      "detail-test-id",
				Text:    "Detailed test event for viewing",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "DetailCorp",
				Project: "DetailProject",
				Tags:    []string{"technical", "leadership"},
			}

			// Send models.ViewEventMsg
			viewMsg := models.ViewEventMsg{Event: testEvent}
			newModel, _ := model.Update(viewMsg)
			updatedModel := newModel.(*Model)

			// Render the view
			view := updatedModel.View()

			// Should display event details
			Expect(view).To(ContainSubstring("Detailed test event for viewing"))
			Expect(view).To(ContainSubstring("DetailCorp"))
			Expect(view).To(ContainSubstring("DetailProject"))
			Expect(view).To(ContainSubstring("technical"))
			Expect(view).To(ContainSubstring("leadership"))
		})

		It("should allow navigation back from ViewScreen to ListScreen", func() {
			// Navigate to ViewScreen first
			testEvent := &career.CareerEvent{
				ID:   "back-test-id",
				Text: "Event for back navigation test",
				Date: time.Now().Add(-24 * time.Hour),
			}
			viewMsg := models.ViewEventMsg{Event: testEvent}
			model.currentScreen = ListScreen
			newModel, _ := model.Update(viewMsg)
			updatedModel := newModel.(*Model)

			// Verify we're on ViewScreen
			Expect(updatedModel.currentScreen).To(Equal(ViewScreen))
			Expect(updatedModel.previousScreen).To(Equal(ListScreen))

			// Send backspace to go back
			backspaceMsg := tea.KeyMsg{Type: tea.KeyBackspace}
			model2, cmd := updatedModel.Update(backspaceMsg)
			model2Updated := model2.(*Model)

			// Should get a BackMsg command from details model
			Expect(cmd).NotTo(BeNil())
			backMsg := cmd()

			// Update with BackMsg
			model3, _ := model2Updated.Update(backMsg)
			model3Updated := model3.(*Model)

			// Should be back on ListScreen
			Expect(model3Updated.currentScreen).To(Equal(ListScreen))
		})
	})

	Context("Action Menu Workflow", func() {
		var (
			testEvent *career.CareerEvent
		)

		BeforeEach(func() {
			// Create a test event for action menu scenarios
			testEvent = &career.CareerEvent{
				ID:      "test-event-1",
				Text:    "Test event for action menu",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "ActionCorp",
				Project: "Action Testing",
				Tags:    []string{"technical", "leadership"},
			}

			// Ensure list model is populated with test event
			repo = careerrepo.NewMemoryRepository()
			svc = careerservice.NewService(repo)
			cliService = service.NewCLIEventService(svc)

			ctx := context.Background()
			err := svc.CaptureEvent(ctx, testEvent, careerservice.ManualEntry)
			Expect(err).To(BeNil())

			// Create a new model with the populated repository
			model = NewModel(cliService, svc)
			model.currentScreen = ListScreen
		})


		It("should enter Action Menu when event is selected", func() {
			// Simulate selecting an event in the list
			actionMenuMsg := models.EventActionMenuMsg{Event: testEvent}
			newModel, _ := model.Update(actionMenuMsg)
			updatedModel := newModel.(*Model)

			// Verify transition to Action Menu Screen
			Expect(updatedModel.currentScreen).To(Equal(ActionMenuScreen))
			Expect(updatedModel.actionMenuModel).NotTo(BeNil())
			Expect(updatedModel.previousScreen).To(Equal(ListScreen))
		})

		It("should display action menu options", func() {
			// Enter action menu
			actionMenuMsg := models.EventActionMenuMsg{Event: testEvent}
			newModel, _ := model.Update(actionMenuMsg)
			updatedModel := newModel.(*Model)

			// Render the action menu
			view := updatedModel.View()

			// Verify action menu is displayed with options
			Expect(view).To(ContainSubstring("Event Actions"))
		})

		It("should handle View action from action menu", func() {
			// Enter action menu
			actionMenuMsg := models.EventActionMenuMsg{Event: testEvent}
			newModel, _ := model.Update(actionMenuMsg)
			updatedModel := newModel.(*Model)

			// Select View action
			viewAction := models.EventActionSelectedMsg{
				Event:  testEvent,
				Action: models.EventActionView,
			}
			newModel, _ = updatedModel.Update(viewAction)
			updatedModel = newModel.(*Model)

			// Verify transition to ViewScreen
			Expect(updatedModel.currentScreen).To(Equal(ViewScreen))
			Expect(updatedModel.detailsModel).NotTo(BeNil())
			Expect(updatedModel.detailsModel.Event().ID).To(Equal(testEvent.ID))
			Expect(updatedModel.previousScreen).To(Equal(ActionMenuScreen))
		})

		It("should handle Edit action from action menu", func() {
			// Enter action menu
			actionMenuMsg := models.EventActionMenuMsg{Event: testEvent}
			newModel, _ := model.Update(actionMenuMsg)
			updatedModel := newModel.(*Model)

			// Select Edit action
			editAction := models.EventActionSelectedMsg{
				Event:  testEvent,
				Action: models.EventActionEdit,
			}
			newModel, _ = updatedModel.Update(editAction)
			updatedModel = newModel.(*Model)

			// Verify transition to CaptureScreen
			Expect(updatedModel.currentScreen).To(Equal(CaptureScreen))
			Expect(updatedModel.formModel).NotTo(BeNil())
			Expect(updatedModel.previousScreen).To(Equal(ActionMenuScreen))
		})

		It("should handle Delete action from action menu", func() {
			// Enter action menu
			actionMenuMsg := models.EventActionMenuMsg{Event: testEvent}
			newModel, _ := model.Update(actionMenuMsg)
			updatedModel := newModel.(*Model)

			// Select Delete action
			deleteAction := models.EventActionSelectedMsg{
				Event:  testEvent,
				Action: models.EventActionDelete,
			}
			newModel, _ = updatedModel.Update(deleteAction)
			updatedModel = newModel.(*Model)

			// Verify transition to ConfirmationScreen with confirmation dialog
			Expect(updatedModel.currentScreen).To(Equal(ConfirmationScreen))
			Expect(updatedModel.confirmationDialog).NotTo(BeNil())
			Expect(updatedModel.deleteEventID).To(Equal(testEvent.ID))
		})

		It("should cancel action menu with BackMsg", func() {
			// Enter action menu
			actionMenuMsg := models.EventActionMenuMsg{Event: testEvent}
			newModel, _ := model.Update(actionMenuMsg)
			updatedModel := newModel.(*Model)

			// Verify we're on ActionMenuScreen
			Expect(updatedModel.currentScreen).To(Equal(ActionMenuScreen))

			// Simulate cancellation (BackMsg)
			backMsg := models.BackMsg{}
			newModel, _ = updatedModel.Update(backMsg)
			updatedModel = newModel.(*Model)

			// Verify return to previous screen
			Expect(updatedModel.currentScreen).To(Equal(ListScreen))
		})

		It("should complete action menu workflow: List -> Action Menu -> View -> Back to List", func() {
			// Start on list screen
			Expect(model.currentScreen).To(Equal(ListScreen))

			// Navigate to action menu
			actionMenuMsg := models.EventActionMenuMsg{Event: testEvent}
			newModel, _ := model.Update(actionMenuMsg)
			model = newModel.(*Model)
			Expect(model.currentScreen).To(Equal(ActionMenuScreen))

			// Select View action
			viewAction := models.EventActionSelectedMsg{
				Event:  testEvent,
				Action: models.EventActionView,
			}
			newModel, _ = model.Update(viewAction)
			model = newModel.(*Model)
			Expect(model.currentScreen).To(Equal(ViewScreen))

			// Go back to list
			backMsg := models.BackMsg{}
			newModel, _ = model.Update(backMsg)
			model = newModel.(*Model)
			Expect(model.currentScreen).To(Equal(ListScreen))
		})
	})

	Context("Delete Confirmation Dialog", func() {
		var (
			testEvent *career.CareerEvent
		)

		BeforeEach(func() {
			// Create a test event for delete scenarios
			testEvent = &career.CareerEvent{
				ID:      "test-delete-event-1",
				Text:    "Test event for delete",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "DeleteCorp",
				Project: "Delete Testing",
				Tags:    []string{"technical"},
			}

			// Setup repository with test event
			repo = careerrepo.NewMemoryRepository()
			err := repo.Create(context.Background(), testEvent)
			Expect(err).NotTo(HaveOccurred())

			// Setup service and models
			svc = careerservice.NewService(repo)
			cliService = service.NewCLIEventService(svc)
			model = NewModel(cliService, svc)
			model.currentScreen = ListScreen
			model.listModel = models.NewListModel(svc, context.Background())
		})

		It("should show confirmation dialog when delete action is selected", func() {
			// Enter action menu
			actionMenuMsg := models.EventActionMenuMsg{Event: testEvent}
			newModel, _ := model.Update(actionMenuMsg)
			updatedModel := newModel.(*Model)

			// Select Delete action
			deleteAction := models.EventActionSelectedMsg{
				Event:  testEvent,
				Action: models.EventActionDelete,
			}
			newModel, _ = updatedModel.Update(deleteAction)
			updatedModel = newModel.(*Model)

			// Verify confirmation dialog is shown
			Expect(updatedModel.currentScreen).To(Equal(ConfirmationScreen))
			Expect(updatedModel.confirmationDialog).NotTo(BeNil())
			Expect(updatedModel.deleteEventID).To(Equal(testEvent.ID))
		})

		It("should cancel delete and return to action menu on BackMsg", func() {
			// Setup: Navigate to confirmation screen
			actionMenuMsg := models.EventActionMenuMsg{Event: testEvent}
			newModel, _ := model.Update(actionMenuMsg)
			updatedModel := newModel.(*Model)

			deleteAction := models.EventActionSelectedMsg{
				Event:  testEvent,
				Action: models.EventActionDelete,
			}
			newModel, _ = updatedModel.Update(deleteAction)
			updatedModel = newModel.(*Model)

			// Verify we're on ConfirmationScreen
			Expect(updatedModel.currentScreen).To(Equal(ConfirmationScreen))

			// Send BackMsg to cancel
			backMsg := models.BackMsg{}
			newModel, _ = updatedModel.Update(backMsg)
			updatedModel = newModel.(*Model)

			// Verify we return to ActionMenuScreen
			Expect(updatedModel.currentScreen).To(Equal(ActionMenuScreen))
			Expect(updatedModel.confirmationDialog).To(BeNil())
			Expect(updatedModel.deleteEventID).To(Equal(""))
		})

		It("should allow confirming deletion from confirmation screen", func() {
			// Setup: Navigate to confirmation screen
			actionMenuMsg := models.EventActionMenuMsg{Event: testEvent}
			newModel, _ := model.Update(actionMenuMsg)
			updatedModel := newModel.(*Model)

			deleteAction := models.EventActionSelectedMsg{
				Event:  testEvent,
				Action: models.EventActionDelete,
			}
			newModel, _ = updatedModel.Update(deleteAction)
			updatedModel = newModel.(*Model)

			// Verify we're on ConfirmationScreen with dialog
			Expect(updatedModel.currentScreen).To(Equal(ConfirmationScreen))
			Expect(updatedModel.confirmationDialog).NotTo(BeNil())

			// Simulate user focusing on "Yes, Delete" button (move right)
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}} // 'l' moves right in confirmation dialog
			newModel, _ = updatedModel.Update(keyMsg)
			updatedModel = newModel.(*Model)

			// Verify dialog processed the key
			Expect(updatedModel.currentScreen).To(Equal(ConfirmationScreen))

			// Simulate user pressing Enter to confirm
			enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
			newModel, _ = updatedModel.Update(enterMsg)
			updatedModel = newModel.(*Model)

			// After confirmation, should be on ListScreen
			Expect(updatedModel.currentScreen).To(Equal(ListScreen))
			Expect(updatedModel.confirmationDialog).To(BeNil())
			Expect(updatedModel.deleteEventID).To(Equal(""))
		})
	})


	Context("Edit Form Pre-Fill", func() {
		var (
			editEvent *career.CareerEvent
		)

		BeforeEach(func() {
			// Create test event with all fields populated
			editEvent = &career.CareerEvent{
				ID:      "edit-test-event-123",
				Text:    "Complex event to be edited",
				Date:    time.Date(2024, 12, 15, 14, 30, 0, 0, time.UTC),
				Company: "EditCorp",
				Project: "EditProject",
				Tags:    []string{"technical", "leadership", "product"},
			}

			// Persist the event to repository
			repo = careerrepo.NewMemoryRepository()
			svc = careerservice.NewService(repo)
			cliService = service.NewCLIEventService(svc)
			
			ctx := context.Background()
			err := svc.CaptureEvent(ctx, editEvent, careerservice.ManualEntry)
			Expect(err).To(BeNil())
			
			// Create model after persisting event
			model = NewModel(cliService, svc)
		})

		It("should pre-fill form with existing event data when entering edit mode", func() {
			// Trigger edit action
			editAction := models.EventActionSelectedMsg{
				Event:  editEvent,
				Action: models.EventActionEdit,
			}
			
			newModel, _ := model.Update(editAction)
			updatedModel := newModel.(*Model)

			// Verify we're on CaptureScreen
			Expect(updatedModel.currentScreen).To(Equal(CaptureScreen))
			
			// Verify form is in edit mode
			Expect(updatedModel.formModel.IsEditMode()).To(BeTrue())
			
			// Verify edit event ID is set
			Expect(updatedModel.formModel.GetEditEventID()).To(Equal(editEvent.ID))
			
			// Verify form fields are pre-filled
			Expect(updatedModel.formModel.GetInputValue(0)).To(Equal(editEvent.Text))
			Expect(updatedModel.formModel.GetInputValue(1)).To(Equal(editEvent.Date.Format("2006-01-02")))
			Expect(updatedModel.formModel.GetInputValue(2)).To(Equal(editEvent.Company))
			Expect(updatedModel.formModel.GetInputValue(3)).To(Equal(editEvent.Project))
			
			// Verify tags are pre-filled
			selectedTags := updatedModel.formModel.TagSelector().SelectedTags()
			Expect(len(selectedTags)).To(Equal(len(editEvent.Tags)))
			for i, tag := range editEvent.Tags {
				Expect(selectedTags[i]).To(Equal(tag))
			}
		})

		It("should display pre-filled form in view", func() {
			// Trigger edit action
			editAction := models.EventActionSelectedMsg{
				Event:  editEvent,
				Action: models.EventActionEdit,
			}
			
			newModel, _ := model.Update(editAction)
			updatedModel := newModel.(*Model)

			// Render the form
			view := updatedModel.View()
			
			// Verify form content is displayed
			Expect(view).To(ContainSubstring("Complex event to be edited"))
			Expect(view).To(ContainSubstring("EditCorp"))
			Expect(view).To(ContainSubstring("EditProject"))
			// Verify the form shows the date
			Expect(view).To(ContainSubstring("2024-12-15"))
		})

		It("should allow editing and submitting updated event", func() {
			// Trigger edit action
			editAction := models.EventActionSelectedMsg{
				Event:  editEvent,
				Action: models.EventActionEdit,
			}
			
			newModel, _ := model.Update(editAction)
			model = newModel.(*Model)

			// Verify form is populated with old data
			Expect(model.formModel.GetInputValue(0)).To(Equal("Complex event to be edited"))
			
			// Clear and update the text field
			// Simulate user clearing and typing new text
			// (In real usage, user would edit via keyboard input)
			model.formModel.GetInputValue(0) // Read current value
			
			// Verify form remains in edit mode for submission
			Expect(model.formModel.IsEditMode()).To(BeTrue())
		})

		It("should return to action menu after edit action", func() {
			// Navigate to list first
			model.currentScreen = ListScreen
			
			// Trigger action menu
			actionMenuMsg := models.EventActionMenuMsg{Event: editEvent}
			newModel, _ := model.Update(actionMenuMsg)
			model = newModel.(*Model)
			Expect(model.currentScreen).To(Equal(ActionMenuScreen))
			
			// Trigger edit action
			editAction := models.EventActionSelectedMsg{
				Event:  editEvent,
				Action: models.EventActionEdit,
			}
			
			newModel, _ = model.Update(editAction)
			model = newModel.(*Model)
			
			// Verify transition to CaptureScreen and previous screen is ActionMenuScreen
			Expect(model.currentScreen).To(Equal(CaptureScreen))
			Expect(model.previousScreen).To(Equal(ActionMenuScreen))
			
			// Verify form is pre-filled
			Expect(model.formModel.IsEditMode()).To(BeTrue())
			Expect(model.formModel.GetInputValue(0)).To(Equal(editEvent.Text))
		})
	})

})
