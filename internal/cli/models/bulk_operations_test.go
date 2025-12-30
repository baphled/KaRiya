package models_test

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	cliservice "github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BulkOperationsModel", func() {
	var (
		repo      *careerrepo.MemoryRepository
		svc       *careerservice.Service
		cliSvc    *cliservice.CLIEventService
		ctx       context.Context
		events    []*career.CareerEvent
		model     *models.BulkOperationsModel
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliSvc = cliservice.NewCLIEventService(svc)
		ctx = context.Background()

		// Create test events
		events = []*career.CareerEvent{
			{
				ID:        "event-1",
				Text:      "First event",
				Date:      time.Now().Add(-24 * time.Hour),
				Company:   "Company A",
				Project:   "Project 1",
				Tags:      []string{"technical"},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "event-2",
				Text:      "Second event",
				Date:      time.Now().Add(-48 * time.Hour),
				Company:   "Company B",
				Project:   "Project 2",
				Tags:      []string{"leadership"},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "event-3",
				Text:      "Third event",
				Date:      time.Now().Add(-72 * time.Hour),
				Company:   "",
				Project:   "",
				Tags:      []string{},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		// Add events to repository
		for _, event := range events {
			err := svc.CaptureEvent(ctx, event, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	Context("when creating BulkOperationsModel", func() {
		It("should create model with events and selection state", func() {
			model = models.NewBulkOperationsModel(events, svc, cliSvc, ctx)
			Expect(model).NotTo(BeNil())
			Expect(len(model.GetEvents())).To(Equal(3))
			Expect(model.GetSelectedCount()).To(Equal(0))
		})

		It("should initialize with zero selected count", func() {
			model = models.NewBulkOperationsModel(events, svc, cliSvc, ctx)
			Expect(model.GetSelectedCount()).To(Equal(0))
		})
	})

	Context("when initializing model", func() {
		It("should implement Init() method", func() {
			model = models.NewBulkOperationsModel(events, svc, cliSvc, ctx)
			cmd := model.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Context("when rendering View", func() {
		It("should render events with checkboxes", func() {
			model = models.NewBulkOperationsModel(events, svc, cliSvc, ctx)
			model.SetSize(80, 24)
			view := model.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("First event"))
		})

		It("should show selection count", func() {
			model = models.NewBulkOperationsModel(events, svc, cliSvc, ctx)
			model.SetSize(80, 24)
			model.ToggleSelection(0)
			view := model.View()
			Expect(view).To(ContainSubstring("1"))
		})

		It("should display checkboxes for each event", func() {
			model = models.NewBulkOperationsModel(events, svc, cliSvc, ctx)
			model.SetSize(80, 24)
			view := model.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("[ ]"),
				ContainSubstring("[x]"),
			))
		})

		It("should handle window resize", func() {
			model = models.NewBulkOperationsModel(events, svc, cliSvc, ctx)
			msg := tea.WindowSizeMsg{Width: 120, Height: 30}
			updatedModel, _ := model.Update(msg)
			bulkModel := updatedModel.(*models.BulkOperationsModel)
			Expect(bulkModel).NotTo(BeNil())
		})
	})

	Context("when handling keyboard input", func() {
		BeforeEach(func() {
			model = models.NewBulkOperationsModel(events, svc, cliSvc, ctx)
			model.SetSize(80, 24)
		})

		It("should toggle selection with Space key", func() {
			initialCount := model.GetSelectedCount()
			msg := tea.KeyMsg{Type: tea.KeySpace}
			_, _ = model.Update(msg)
			Expect(model.GetSelectedCount()).To(Equal(initialCount + 1))

			_, _ = model.Update(msg)
			Expect(model.GetSelectedCount()).To(Equal(initialCount))
		})

		It("should select all events with 'a' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
			_, _ = model.Update(msg)
			Expect(model.GetSelectedCount()).To(Equal(3))
		})

		It("should deselect all events with 'd' key", func() {
			model.SelectAll()
			Expect(model.GetSelectedCount()).To(Equal(3))

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
			_, _ = model.Update(msg)
			Expect(model.GetSelectedCount()).To(Equal(0))
		})

		It("should navigate between events with up/down arrows", func() {
			msgDown := tea.KeyMsg{Type: tea.KeyDown}
			_, _ = model.Update(msgDown)
			Expect(model.GetCurrentIndex()).To(Equal(1))

			msgUp := tea.KeyMsg{Type: tea.KeyUp}
			_, _ = model.Update(msgUp)
			Expect(model.GetCurrentIndex()).To(Equal(0))
		})

		It("should wrap around when navigating past last event", func() {
			model.SetCurrentIndex(2)
			msgDown := tea.KeyMsg{Type: tea.KeyDown}
			_, _ = model.Update(msgDown)
			Expect(model.GetCurrentIndex()).To(Equal(0))
		})

		It("should wrap around when navigating before first event", func() {
			model.SetCurrentIndex(0)
			msgUp := tea.KeyMsg{Type: tea.KeyUp}
			_, _ = model.Update(msgUp)
			Expect(model.GetCurrentIndex()).To(Equal(2))
		})
	})

	Context("when handling bulk field editing", func() {
		BeforeEach(func() {
			model = models.NewBulkOperationsModel(events, svc, cliSvc, ctx)
			model.SetSize(80, 24)
			model.ToggleSelection(0)
			model.ToggleSelection(1)
		})

		It("should enter bulk edit mode with 'e' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			_, _ = model.Update(msg)
			Expect(model.IsInEditMode()).To(BeTrue())
		})

		It("should display bulk edit form", func() {
			model.EnterEditMode()
			model.SetBulkCompany("New Company")
			view := model.View()
			Expect(view).To(ContainSubstring("New Company"))
		})

		It("should allow editing company field in bulk mode", func() {
			model.EnterEditMode()
			model.SetBulkCompany("Test Company")
			Expect(model.GetBulkCompany()).To(Equal("Test Company"))
		})

		It("should allow editing project field in bulk mode", func() {
			model.EnterEditMode()
			model.SetBulkProject("Test Project")
			Expect(model.GetBulkProject()).To(Equal("Test Project"))
		})
	})

	Context("when previewing bulk changes", func() {
		BeforeEach(func() {
			model = models.NewBulkOperationsModel(events, svc, cliSvc, ctx)
			model.SetSize(80, 24)
			model.ToggleSelection(0)
			model.ToggleSelection(1)
			model.EnterEditMode()
			model.SetBulkCompany("New Company")
		})

		It("should show preview of changes", func() {
			model.ShowPreview()
			Expect(model.IsShowingPreview()).To(BeTrue())
		})

		It("should display affected events in preview", func() {
			model.ShowPreview()
			view := model.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("First event"),
				ContainSubstring("Second event"),
			))
		})

		It("should display changes to be applied", func() {
			model.ShowPreview()
			view := model.View()
			Expect(view).To(ContainSubstring("New Company"))
		})
	})

	Context("when confirming bulk changes", func() {
		BeforeEach(func() {
			model = models.NewBulkOperationsModel(events, svc, cliSvc, ctx)
			model.SetSize(80, 24)
			model.ToggleSelection(0)
			model.ToggleSelection(1)
		})

		It("should show confirmation dialog", func() {
			model.ShowConfirmation()
			Expect(model.IsShowingConfirmation()).To(BeTrue())
		})

		It("should accept confirmation with 'y' key", func() {
			model.ShowConfirmation()
			model.SetBulkCompany("New Company")
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
			_, _ = model.Update(msg)
			Expect(model.IsShowingConfirmation()).To(BeFalse())
		})

		It("should cancel confirmation with 'n' key", func() {
			model.ShowConfirmation()
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
			_, _ = model.Update(msg)
			Expect(model.IsShowingConfirmation()).To(BeFalse())
		})
	})

	Context("when undoing bulk operations", func() {
		BeforeEach(func() {
			model = models.NewBulkOperationsModel(events, svc, cliSvc, ctx)
			model.SetSize(80, 24)
			model.ToggleSelection(0)
		})

		It("should revert to previous state with 'u' key", func() {
			originalCompany := model.GetEvents()[0].Company
			model.SetBulkCompany("New Company")
			model.ApplyChanges()
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
			_, _ = model.Update(msg)
			Expect(model.GetEvents()[0].Company).To(Equal(originalCompany))
		})
	})

	Context("when getting selection state", func() {
		BeforeEach(func() {
			model = models.NewBulkOperationsModel(events, svc, cliSvc, ctx)
		})

		It("should return correct selected count", func() {
			model.ToggleSelection(0)
			model.ToggleSelection(2)
			Expect(model.GetSelectedCount()).To(Equal(2))
		})

		It("should return selected event IDs", func() {
			model.ToggleSelection(0)
			model.ToggleSelection(1)
			ids := model.GetSelectedEventIDs()
			Expect(ids).To(ContainElement("event-1"))
			Expect(ids).To(ContainElement("event-2"))
			Expect(ids).NotTo(ContainElement("event-3"))
		})

		It("should return empty list when no events selected", func() {
			ids := model.GetSelectedEventIDs()
			Expect(ids).To(BeEmpty())
		})

		It("should indicate if changes were applied", func() {
			Expect(model.ChangesApplied()).To(BeFalse())
			model.MarkChangesApplied()
			Expect(model.ChangesApplied()).To(BeTrue())
		})

		It("should indicate if operation was cancelled", func() {
			Expect(model.WasCancelled()).To(BeFalse())
			model.Cancel()
			Expect(model.WasCancelled()).To(BeTrue())
		})
	})

	Context("when handling empty event list", func() {
		It("should handle zero events gracefully", func() {
			model = models.NewBulkOperationsModel([]*career.CareerEvent{}, svc, cliSvc, ctx)
			Expect(len(model.GetEvents())).To(Equal(0))
			view := model.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("No events"),
				ContainSubstring("empty"),
			))
		})
	})

	Context("when handling single event", func() {
		It("should handle single event selection", func() {
			singleEvent := []*career.CareerEvent{events[0]}
			model = models.NewBulkOperationsModel(singleEvent, svc, cliSvc, ctx)
			model.ToggleSelection(0)
			Expect(model.GetSelectedCount()).To(Equal(1))
		})
	})

	Context("when handling large event list", func() {
		It("should handle 100+ events without error", func() {
			largeEventList := make([]*career.CareerEvent, 100)
			for i := 0; i < 100; i++ {
				largeEventList[i] = &career.CareerEvent{
					ID:        "event-" + string(rune(i)),
					Text:      "Event " + string(rune(i)),
					Date:      time.Now().Add(-time.Duration(i) * time.Hour),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
			}
			model = models.NewBulkOperationsModel(largeEventList, svc, cliSvc, ctx)
			Expect(len(model.GetEvents())).To(Equal(100))
		})
	})
})

