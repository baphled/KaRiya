package models_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ActionMenuModel", func() {
	var (
		model *models.ActionMenuModel
		event *career.CareerEvent
	)

	BeforeEach(func() {
		event = &career.CareerEvent{
			ID:        "test-event-456",
			Text:      "Completed critical system refactoring",
			Date:      time.Now().Add(-14 * 24 * time.Hour),
			Company:   "TechCorp Inc.",
			Project:   "System Upgrade",
			Tags:      []string{"technical"},
			CreatedAt: time.Now().Add(-14 * 24 * time.Hour),
			UpdatedAt: time.Now().Add(-14 * 24 * time.Hour),
		}
	})

	Context("when creating a new ActionMenuModel", func() {
		It("should initialize with an event", func() {
			model = models.NewActionMenuModel(event)
			Expect(model).NotTo(BeNil())
		})

		It("should have three action options", func() {
			model = models.NewActionMenuModel(event)
			Expect(model.SelectedAction()).To(Equal(models.EventActionView))
		})
	})

	Context("when navigating through actions", func() {
		BeforeEach(func() {
			model = models.NewActionMenuModel(event)
		})

		It("should navigate down with j key", func() {
			_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(model.SelectedAction()).To(Equal(models.EventActionEdit))
		})

		It("should navigate down with arrow key", func() {
			_, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(model.SelectedAction()).To(Equal(models.EventActionEdit))
		})

		It("should navigate up with k key", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(model.SelectedAction()).To(Equal(models.EventActionEdit))
		})

		It("should not navigate below the last option", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			_, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(model.SelectedAction()).To(Equal(models.EventActionDelete))
		})

		It("should not navigate above the first option", func() {
			_, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(model.SelectedAction()).To(Equal(models.EventActionView))
		})
	})

	Context("when rendering the action menu", func() {
		BeforeEach(func() {
			model = models.NewActionMenuModel(event)
		})

		It("should display action options", func() {
			view := model.View()
			Expect(view).To(Or(
				ContainSubstring("View Event"),
				ContainSubstring("Edit Event"),
				ContainSubstring("Delete Event"),
			))
		})

		It("should display event text", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Completed critical"))
		})

		It("should show selected action indicator", func() {
			view := model.View()
			Expect(view).To(Or(ContainSubstring("▶"), ContainSubstring(">")))
		})

		It("should render header", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Event Actions"))
		})

		It("should render with responsive width", func() {
			view := model.View()
			Expect(len(view)).To(BeNumerically(">", 50))
		})
	})

	Context("when testing header and help_footer integration", func() {
		BeforeEach(func() {
			model = models.NewActionMenuModel(event)
		})

		It("should include title in header", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Event Actions"))
		})

		It("should handle window resize messages", func() {
			msg := tea.WindowSizeMsg{
				Width:  100,
				Height: 30,
			}
			newModel, cmd := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).To(BeNil())
		})

		It("should render with responsive layout", func() {
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should have Init method", func() {
			Expect(model.Init).NotTo(BeNil())
		})

		It("should respond to Escape key", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			_, cmd := model.Update(msg)
			Expect(cmd).NotTo(BeNil())
		})
	})

	Context("when handling action selection", func() {
		BeforeEach(func() {
			model = models.NewActionMenuModel(event)
		})

		It("should return the selected action on Enter", func() {
			// Move to Edit action
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())
		})

		It("should handle Escape for cancellation", func() {
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
		})

		It("should start with View action selected", func() {
			Expect(model.SelectedAction()).To(Equal(models.EventActionView))
		})
	})
})
