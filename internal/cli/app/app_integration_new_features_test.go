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

var _ = Describe("CLI App - Event Action Menu Integration", func() {
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

	Describe("Action Menu Navigation", func() {
		var testEvent *career.CareerEvent

		BeforeEach(func() {
			testEvent = &career.CareerEvent{
				ID:      "action-test-1",
				Text:    "Test event for action menu",
				Date:    time.Now().Add(-5 * 24 * time.Hour),
				Company: "TestCorp",
				Project: "TestProject",
				Tags:    []string{"technical"},
			}
			err := svc.CaptureEvent(ctx, testEvent, careerservice.ManualEntry)
			Expect(err).To(BeNil())
		})

		It("should navigate to action menu when EventActionMenuMsg is received", func() {
			actionMenuMsg := models.EventActionMenuMsg{Event: testEvent}
			newModel, _ := model.Update(actionMenuMsg)
			updatedModel := newModel.(*Model)

			Expect(updatedModel.currentScreen).To(Equal(ActionMenuScreen))
			Expect(updatedModel.actionMenuModel).ToNot(BeNil())
		})

		It("should navigate to view screen when View action is selected", func() {
			model.Update(models.EventActionMenuMsg{Event: testEvent})
			viewAction := models.EventActionSelectedMsg{
				Event:  testEvent,
				Action: models.EventActionView,
			}
			newModel, _ := model.Update(viewAction)
			updatedModel := newModel.(*Model)

			Expect(updatedModel.currentScreen).To(Equal(ViewScreen))
			Expect(updatedModel.detailsModel).ToNot(BeNil())
		})

		It("should navigate to capture screen when Edit action is selected", func() {
			model.Update(models.EventActionMenuMsg{Event: testEvent})
			editAction := models.EventActionSelectedMsg{
				Event:  testEvent,
				Action: models.EventActionEdit,
			}
			newModel, _ := model.Update(editAction)
			updatedModel := newModel.(*Model)

			Expect(updatedModel.currentScreen).To(Equal(CaptureScreen))
			Expect(updatedModel.formModel.IsEditMode()).To(BeTrue())
		})

		It("should show confirmation dialog when Delete action is selected", func() {
			model.Update(models.EventActionMenuMsg{Event: testEvent})
			deleteAction := models.EventActionSelectedMsg{
				Event:  testEvent,
				Action: models.EventActionDelete,
			}
			newModel, _ := model.Update(deleteAction)
			updatedModel := newModel.(*Model)

			Expect(updatedModel.currentScreen).To(Equal(ConfirmationScreen))
			Expect(updatedModel.confirmationDialog).ToNot(BeNil())
			Expect(updatedModel.deleteEventID).To(Equal(testEvent.ID))
		})
	})

	Describe("Event Deletion Flow", func() {
		var testEvent *career.CareerEvent

		BeforeEach(func() {
			testEvent = &career.CareerEvent{
				ID:      "delete-test-1",
				Text:    "Event to be deleted",
				Date:    time.Now().Add(-3 * 24 * time.Hour),
				Company: "DeleteCorp",
				Tags:    []string{"technical"},
			}
			err := svc.CaptureEvent(ctx, testEvent, careerservice.ManualEntry)
			Expect(err).To(BeNil())
		})

		It("should show confirmation dialog before deleting", func() {
			deleteAction := models.EventActionSelectedMsg{
				Event:  testEvent,
				Action: models.EventActionDelete,
			}
			newModel, _ := model.Update(deleteAction)
			updatedModel := newModel.(*Model)

			Expect(updatedModel.currentScreen).To(Equal(ConfirmationScreen))
			Expect(updatedModel.confirmationDialog).ToNot(BeNil())
			Expect(updatedModel.deleteEventID).To(Equal(testEvent.ID))
		})

		It("should cancel deletion and keep event in repository", func() {
			// Verify event exists
			events, err := svc.ListEvents(ctx, careerrepo.ListFilters{})
			Expect(err).To(BeNil())
			Expect(len(events)).To(Equal(1))

			// Trigger delete
			model.Update(models.EventActionSelectedMsg{
				Event:  testEvent,
				Action: models.EventActionDelete,
			})

			// Cancel deletion
			keyMsg := tea.KeyMsg{Type: tea.KeyEsc}
			newModel, _ := model.Update(keyMsg)
			updatedModel := newModel.(*Model)

			// Verify event still exists
			events, err = svc.ListEvents(ctx, careerrepo.ListFilters{})
			Expect(err).To(BeNil())
			Expect(len(events)).To(Equal(1))

			// Verify we're on ActionMenuScreen
			Expect(updatedModel.currentScreen).To(Equal(ActionMenuScreen))
		})
	})

	Describe("Navigation State Management", func() {
		var testEvent *career.CareerEvent

		BeforeEach(func() {
			testEvent = &career.CareerEvent{
				ID:      "nav-test-1",
				Text:    "Navigation test event",
				Date:    time.Now().Add(-2 * 24 * time.Hour),
				Company: "NavCorp",
				Tags:    []string{"leadership"},
			}
			err := svc.CaptureEvent(ctx, testEvent, careerservice.ManualEntry)
			Expect(err).To(BeNil())
		})

		It("should preserve screen context when navigating from ListScreen through ActionMenu to ViewScreen", func() {
			model.currentScreen = ListScreen
			model.previousScreen = HomeScreen

			model.Update(models.EventActionMenuMsg{Event: testEvent})
			model.Update(models.EventActionSelectedMsg{
				Event:  testEvent,
				Action: models.EventActionView,
			})

			Expect(model.screenBeforeActionMenu).To(Equal(ListScreen))

			backMsg := models.BackMsg{}
			newModel, _ := model.Update(backMsg)
			updatedModel := newModel.(*Model)

			Expect(updatedModel.currentScreen).To(Equal(ListScreen))
		})

		It("should handle back navigation from ImportReviewScreen", func() {
			model.currentScreen = ImportReviewScreen
			model.previousScreen = HomeScreen

			backMsg := models.BackMsg{}
			newModel, _ := model.Update(backMsg)
			updatedModel := newModel.(*Model)

			Expect(updatedModel.currentScreen).To(Equal(HomeScreen))
			Expect(updatedModel.importReviewModel).To(BeNil())
			Expect(updatedModel.importFilePath).To(Equal(""))
		})
	})

	Describe("Global Navigation Shortcuts from HomeScreen", func() {
		It("should navigate to capture with 'c' key from home", func() {
			model.currentScreen = HomeScreen

			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
			newModel, _ := model.Update(keyMsg)
			updatedModel := newModel.(*Model)

			Expect(updatedModel.currentScreen).To(Equal(CaptureScreen))
		})

		It("should navigate to list with 'l' key from home", func() {
			model.currentScreen = HomeScreen

			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
			newModel, _ := model.Update(keyMsg)
			updatedModel := newModel.(*Model)

			Expect(updatedModel.currentScreen).To(Equal(ListScreen))
		})

		It("should update dimensions on window resize", func() {
			windowMsg := tea.WindowSizeMsg{Width: 150, Height: 50}
			newModel, _ := model.Update(windowMsg)
			updatedModel := newModel.(*Model)

			Expect(updatedModel.width).To(Equal(150))
			Expect(updatedModel.height).To(Equal(50))
		})
	})

	Describe("SetInitialScreen Method", func() {
		It("should set initial screen to CaptureScreen", func() {
			model.SetInitialScreen(CaptureScreen)

			Expect(model.currentScreen).To(Equal(CaptureScreen))
			Expect(model.previousScreen).To(Equal(CaptureScreen))
		})

		It("should set initial screen to ListScreen", func() {
			model.SetInitialScreen(ListScreen)

			Expect(model.currentScreen).To(Equal(ListScreen))
			Expect(model.previousScreen).To(Equal(ListScreen))
		})
	})

	Describe("SetInitialCaptureMode Method", func() {
		It("should set initial capture mode without panicking", func() {
			model.formModel = models.NewFormModel(cliService)

			Expect(func() {
				model.SetInitialCaptureMode("timeline")
			}).NotTo(Panic())
		})

		It("should handle nil form model gracefully", func() {
			model.formModel = nil

			Expect(func() {
				model.SetInitialCaptureMode("timeline")
			}).NotTo(Panic())
		})
	})

	Describe("Model Initialization", func() {
		It("should initialize with correct default state", func() {
			newModel := NewModel(cliService, svc)

			Expect(newModel.currentScreen).To(Equal(HomeScreen))
			Expect(newModel.previousScreen).To(Equal(HomeScreen))
			Expect(newModel.width).To(Equal(80))
			Expect(newModel.height).To(Equal(24))
			Expect(newModel.formModel).ToNot(BeNil())
			Expect(newModel.listModel).ToNot(BeNil())
			Expect(newModel.importService).ToNot(BeNil())
		})

		It("should maintain state across window resize updates", func() {
			model.currentScreen = CaptureScreen
			model.previousScreen = HomeScreen

			windowMsg := tea.WindowSizeMsg{Width: 100, Height: 30}
			newModel, _ := model.Update(windowMsg)
			updatedModel := newModel.(*Model)

			Expect(updatedModel.currentScreen).To(Equal(CaptureScreen))
			Expect(updatedModel.previousScreen).To(Equal(HomeScreen))
		})
	})

	Describe("SuccessModel Integration", func() {
		It("should transition to success screen after form submission", func() {
			testEvent := &career.CareerEvent{
				ID:      "success-event-1",
				Text:    "Test event",
				Date:    time.Now().Add(-1 * 24 * time.Hour),
				Company: "TestCorp",
				Tags:    []string{"technical"},
			}

			model.successModel = models.NewSuccessModel(testEvent)
			model.currentScreen = SuccessScreen

			view := model.View()
			Expect(view).To(ContainSubstring(testEvent.Text))
		})

		It("should handle CaptureAnotherMsg to reset form and return to capture", func() {
			model.currentScreen = SuccessScreen
			model.successModel = models.NewSuccessModel(&career.CareerEvent{
				Text: "Test",
				Date: time.Now(),
			})

			captureMsg := models.CaptureAnotherMsg{}
			newModel, _ := model.Update(captureMsg)
			updatedModel := newModel.(*Model)

			Expect(updatedModel.currentScreen).To(Equal(CaptureScreen))
			Expect(updatedModel.successModel).To(BeNil())
			Expect(updatedModel.formModel).ToNot(BeNil())
		})

		It("should handle ViewRecentMsg to navigate to list", func() {
			model.currentScreen = SuccessScreen
			model.successModel = models.NewSuccessModel(&career.CareerEvent{
				Text: "Test",
				Date: time.Now(),
			})

			viewMsg := models.ViewRecentMsg{}
			newModel, _ := model.Update(viewMsg)
			updatedModel := newModel.(*Model)

			Expect(updatedModel.currentScreen).To(Equal(ListScreen))
			Expect(updatedModel.successModel).To(BeNil())
			Expect(updatedModel.listModel).ToNot(BeNil())
		})
	})

	Describe("View Rendering", func() {
		It("should render home screen correctly", func() {
			model.currentScreen = HomeScreen

			view := model.View()
			Expect(view).To(ContainSubstring("KaRiya"))
			Expect(view).To(ContainSubstring("Career Journal"))
		})

		It("should render error message when form model is nil", func() {
			model.currentScreen = CaptureScreen
			model.formModel = nil

			view := model.View()
			Expect(view).To(ContainSubstring("Error"))
		})

		It("should render error message when list model is nil", func() {
			model.currentScreen = ListScreen
			model.listModel = nil

			view := model.View()
			Expect(view).To(ContainSubstring("Error"))
		})

		It("should render error message when action menu model is nil", func() {
			model.currentScreen = ActionMenuScreen
			model.actionMenuModel = nil

			view := model.View()
			Expect(view).To(ContainSubstring("Error"))
		})
	})
})
