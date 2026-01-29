package modals_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens/skills/modals"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("EventsModal", func() {
	var (
		modal     *modals.EventsModal
		theme     themes.Theme
		testEvent *career.Event
		events    []*career.Event
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
		testEvent = &career.Event{
			ID:      "event-1",
			Text:    "Implemented new feature",
			Date:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Company: "Test Corp",
			Project: "Project Alpha",
		}
		events = []*career.Event{testEvent}
		modal = modals.NewEventsModal("skill-1", "Go Programming", events, theme)
	})

	Describe("NewEventsModal", func() {
		It("creates a new modal with correct initial state", func() {
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.GetSkillID()).To(Equal("skill-1"))
			Expect(modal.GetSkillName()).To(Equal("Go Programming"))
			Expect(modal.HasSelection()).To(BeFalse())
		})

		It("sets default dimensions", func() {
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Show and Hide", func() {
		It("shows the modal", func() {
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("hides the modal", func() {
			modal.Show()
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("clears selection when showing", func() {
			modal.Show()
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.HasSelection()).To(BeTrue())

			modal.Show() // Re-show should clear selection
			Expect(modal.HasSelection()).To(BeFalse())
		})

		It("resets selection index when showing", func() {
			// Move selection down then re-show
			modal.Show()
			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			modal.Show()
			Expect(modal.GetSelectedIndex()).To(Equal(0))
		})
	})

	Describe("View", func() {
		It("returns empty string when not visible", func() {
			view := modal.View()
			Expect(view).To(BeEmpty())
		})

		It("shows skill name in title", func() {
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("Go Programming"))
		})

		It("shows event count in title", func() {
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("(1)"))
		})

		It("shows event text", func() {
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("Implemented new feature"))
		})

		It("shows event date formatted", func() {
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("2024-01-15"))
		})

		It("shows event company", func() {
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("Test Corp"))
		})

		It("shows table header columns", func() {
			modal.Show()
			view := modal.View()
			// Table shows Date, Event, Company columns
			Expect(view).To(ContainSubstring("Date"))
			Expect(view).To(ContainSubstring("Event"))
			Expect(view).To(ContainSubstring("Company"))
		})

		It("shows footer with keyboard hints", func() {
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("Enter: View Details"))
			Expect(view).To(ContainSubstring("Navigate"))
			Expect(view).To(ContainSubstring("Esc: Close"))
		})

		It("shows empty state message when no events", func() {
			emptyModal := modals.NewEventsModal("skill-1", "Go", nil, theme)
			emptyModal.Show()
			view := emptyModal.View()
			Expect(view).To(ContainSubstring("No events use this skill"))
		})

		It("shows current position in footer", func() {
			modal.Show()
			view := modal.View()
			// Table shows position indicator like [1/1]
			Expect(view).To(ContainSubstring("[1/1]"))
		})

		It("truncates long event text", func() {
			longEvent := &career.Event{
				ID:   "event-long",
				Text: "This is a very long event description that should be truncated when displayed in the modal to prevent layout issues",
				Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			}
			longModal := modals.NewEventsModal("skill-1", "Go", []*career.Event{longEvent}, theme)
			longModal.SetDimensions(60, 24) // Small width to force truncation
			longModal.Show()
			view := longModal.View()
			Expect(view).To(ContainSubstring("..."))
		})
	})

	Describe("Navigation", func() {
		var multiEventModal *modals.EventsModal

		BeforeEach(func() {
			events := []*career.Event{
				{ID: "event-1", Text: "First event", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
				{ID: "event-2", Text: "Second event", Date: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
				{ID: "event-3", Text: "Third event", Date: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
			}
			multiEventModal = modals.NewEventsModal("skill-1", "Go", events, theme)
			multiEventModal.Show()
		})

		It("navigates down with down arrow", func() {
			multiEventModal.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(multiEventModal.GetSelectedIndex()).To(Equal(1))
		})

		It("navigates down with j key", func() {
			multiEventModal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(multiEventModal.GetSelectedIndex()).To(Equal(1))
		})

		It("navigates up with up arrow", func() {
			multiEventModal.Update(tea.KeyMsg{Type: tea.KeyDown})
			multiEventModal.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(multiEventModal.GetSelectedIndex()).To(Equal(0))
		})

		It("navigates up with k key", func() {
			multiEventModal.Update(tea.KeyMsg{Type: tea.KeyDown})
			multiEventModal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(multiEventModal.GetSelectedIndex()).To(Equal(0))
		})

		It("does not go below zero", func() {
			multiEventModal.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(multiEventModal.GetSelectedIndex()).To(Equal(0))
		})

		It("does not go beyond last item", func() {
			multiEventModal.Update(tea.KeyMsg{Type: tea.KeyDown})
			multiEventModal.Update(tea.KeyMsg{Type: tea.KeyDown})
			multiEventModal.Update(tea.KeyMsg{Type: tea.KeyDown}) // Should stay at index 2
			Expect(multiEventModal.GetSelectedIndex()).To(Equal(2))
		})

		It("goes to first item with home key", func() {
			multiEventModal.Update(tea.KeyMsg{Type: tea.KeyDown})
			multiEventModal.Update(tea.KeyMsg{Type: tea.KeyDown})
			multiEventModal.Update(tea.KeyMsg{Type: tea.KeyHome})
			Expect(multiEventModal.GetSelectedIndex()).To(Equal(0))
		})

		It("goes to first item with g key", func() {
			multiEventModal.Update(tea.KeyMsg{Type: tea.KeyDown})
			multiEventModal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
			Expect(multiEventModal.GetSelectedIndex()).To(Equal(0))
		})

		It("goes to last item with end key", func() {
			multiEventModal.Update(tea.KeyMsg{Type: tea.KeyEnd})
			Expect(multiEventModal.GetSelectedIndex()).To(Equal(2))
		})

		It("goes to last item with G key", func() {
			multiEventModal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
			Expect(multiEventModal.GetSelectedIndex()).To(Equal(2))
		})
	})

	Describe("Selection", func() {
		It("selects event on Enter", func() {
			modal.Show()
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(modal.HasSelection()).To(BeTrue())
			Expect(modal.GetSelectedEvent()).To(Equal(testEvent))
		})

		It("hides modal after selection", func() {
			modal.Show()
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("does not select when no events", func() {
			emptyModal := modals.NewEventsModal("skill-1", "Go", nil, theme)
			emptyModal.Show()
			emptyModal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(emptyModal.HasSelection()).To(BeFalse())
		})

		It("clears selection with ClearSelection", func() {
			modal.Show()
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.HasSelection()).To(BeTrue())

			modal.ClearSelection()
			Expect(modal.HasSelection()).To(BeFalse())
			Expect(modal.GetSelectedEvent()).To(BeNil())
		})
	})

	Describe("Close without selection", func() {
		BeforeEach(func() {
			modal.Show()
		})

		It("closes with Escape", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEscape})
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.HasSelection()).To(BeFalse())
		})

		It("closes with backspace", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyBackspace})
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.HasSelection()).To(BeFalse())
		})

		It("closes with q", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.HasSelection()).To(BeFalse())
		})
	})

	Describe("SetDimensions", func() {
		It("updates modal dimensions", func() {
			modal.SetDimensions(100, 50)
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("resets viewport on dimension change", func() {
			modal.Show()
			modal.SetDimensions(100, 50)
			// View should rebuild with new dimensions
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("SetEvents", func() {
		It("updates events list", func() {
			newEvents := []*career.Event{
				{ID: "new-1", Text: "New event", Date: time.Now()},
			}
			modal.SetEvents(newEvents)
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("New event"))
		})

		It("resets selection index", func() {
			modal.Show()
			modal.Update(tea.KeyMsg{Type: tea.KeyDown})

			modal.SetEvents([]*career.Event{
				{ID: "new-1", Text: "New event"},
			})
			Expect(modal.GetSelectedIndex()).To(Equal(0))
		})

		It("clears previous selection", func() {
			modal.Show()
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.HasSelection()).To(BeTrue())

			modal.SetEvents([]*career.Event{})
			Expect(modal.HasSelection()).To(BeFalse())
		})
	})

	Describe("WindowSizeMsg handling", func() {
		It("updates dimensions on window resize", func() {
			modal.Show()
			modal.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Init", func() {
		It("returns nil command", func() {
			cmd := modal.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update when not visible", func() {
		It("ignores key events", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.HasSelection()).To(BeFalse())
		})
	})

	Describe("Handles nil events in list", func() {
		It("skips nil events when rendering", func() {
			eventsWithNil := []*career.Event{
				{ID: "event-1", Text: "First event", Date: time.Now()},
				nil,
				{ID: "event-3", Text: "Third event", Date: time.Now()},
			}
			nilModal := modals.NewEventsModal("skill-1", "Go", eventsWithNil, theme)
			nilModal.Show()
			view := nilModal.View()
			Expect(view).To(ContainSubstring("First event"))
			Expect(view).To(ContainSubstring("Third event"))
		})
	})

	Describe("Handles events with missing optional fields", func() {
		It("renders event without company", func() {
			event := &career.Event{
				ID:   "event-1",
				Text: "Event without company",
				Date: time.Now(),
			}
			simpleModal := modals.NewEventsModal("skill-1", "Go", []*career.Event{event}, theme)
			simpleModal.Show()
			view := simpleModal.View()
			Expect(view).To(ContainSubstring("Event without company"))
		})

		It("renders event without project", func() {
			event := &career.Event{
				ID:      "event-1",
				Text:    "Event without project",
				Date:    time.Now(),
				Company: "Company",
			}
			simpleModal := modals.NewEventsModal("skill-1", "Go", []*career.Event{event}, theme)
			simpleModal.Show()
			view := simpleModal.View()
			Expect(view).To(ContainSubstring("Event without project"))
			Expect(view).To(ContainSubstring("Company"))
		})

		It("renders event with zero date", func() {
			event := &career.Event{
				ID:   "event-1",
				Text: "Event without date",
			}
			simpleModal := modals.NewEventsModal("skill-1", "Go", []*career.Event{event}, theme)
			simpleModal.Show()
			view := simpleModal.View()
			Expect(view).To(ContainSubstring("Event without date"))
		})

		It("shows placeholder for empty text", func() {
			event := &career.Event{
				ID:   "event-1",
				Text: "",
				Date: time.Now(),
			}
			simpleModal := modals.NewEventsModal("skill-1", "Go", []*career.Event{event}, theme)
			simpleModal.Show()
			view := simpleModal.View()
			Expect(view).To(ContainSubstring("No description"))
		})
	})
})
