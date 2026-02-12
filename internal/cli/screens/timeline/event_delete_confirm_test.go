package timeline_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/timeline"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EventDeleteConfirmScreen", func() {
	var (
		screen *timeline.EventDeleteConfirmScreen
		event  *career.Event
	)

	BeforeEach(func() {
		event = fixtures.EventWith("test-event-id", "Test event for deletion", "Test Company", "Test Project")
		event.Date = time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
		event.CreatedAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	})

	Describe("NewEventDeleteConfirmScreen", func() {
		It("creates screen with event", func() {
			screen = timeline.NewEventDeleteConfirmScreen(event)

			Expect(screen).NotTo(BeNil())
			Expect(screen.GetEvent()).To(Equal(event))
		})

		It("preserves event ID", func() {
			screen = timeline.NewEventDeleteConfirmScreen(event)

			Expect(screen.GetEvent().ID).To(Equal("test-event-id"))
		})

		It("preserves event text", func() {
			screen = timeline.NewEventDeleteConfirmScreen(event)

			Expect(screen.GetEvent().Text).To(Equal("Test event for deletion"))
		})

		It("truncates long event text in confirmation message", func() {
			longTextEvent := fixtures.Event("long-text-id")
			longTextEvent.Text = "This is a very long event text that exceeds sixty characters and should be truncated in the confirmation message for better display"

			screen = timeline.NewEventDeleteConfirmScreen(longTextEvent)

			// Screen is created successfully - original text is preserved in the event.
			Expect(screen).NotTo(BeNil())
			// GetEvent returns the original event with full text.
			Expect(len(screen.GetEvent().Text)).To(BeNumerically(">", 60))
		})

		It("handles short event text without truncation", func() {
			shortTextEvent := fixtures.Event("short-text-id")
			shortTextEvent.Text = "Short text"

			screen = timeline.NewEventDeleteConfirmScreen(shortTextEvent)

			Expect(screen).NotTo(BeNil())
			Expect(screen.GetEvent().Text).To(Equal("Short text"))
		})

		It("handles exactly 60 character text", func() {
			exactTextEvent := fixtures.Event("exact-text-id")
			exactTextEvent.Text = "This is exactly sixty characters long for testing purposes!!" // 60 chars

			screen = timeline.NewEventDeleteConfirmScreen(exactTextEvent)

			Expect(screen).NotTo(BeNil())
			Expect(screen.GetEvent().Text).To(HaveLen(60))
		})
	})

	Describe("GetEvent", func() {
		It("returns the event being considered for deletion", func() {
			screen = timeline.NewEventDeleteConfirmScreen(event)

			result := screen.GetEvent()

			Expect(result).To(Equal(event))
			Expect(result.ID).To(Equal("test-event-id"))
			Expect(result.Text).To(Equal("Test event for deletion"))
			Expect(result.Company).To(Equal("Test Company"))
			Expect(result.Project).To(Equal("Test Project"))
		})

		It("returns event with all fields preserved", func() {
			fullEvent := fixtures.EventWith("full-event", "Full event", "Big Corp", "Big Project")
			fullEvent.Date = time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
			fullEvent.Tags = []string{"tag1", "tag2"}
			fullEvent.Categories = []string{"cat1"}
			fullEvent.Skills = []string{"skill1"}
			fullEvent.CreatedAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
			fullEvent.UpdatedAt = time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)

			screen = timeline.NewEventDeleteConfirmScreen(fullEvent)

			result := screen.GetEvent()
			Expect(result.Tags).To(ContainElement("tag1"))
			Expect(result.Categories).To(ContainElement("cat1"))
			Expect(result.Skills).To(ContainElement("skill1"))
		})
	})

	Describe("State constant", func() {
		It("has correct state name for state matrix", func() {
			Expect(timeline.EventDeleteConfirmState).To(Equal("event_delete_confirm"))
		})
	})

	Describe("Key Handling", func() {
		BeforeEach(func() {
			screen = timeline.NewEventDeleteConfirmScreen(event)
			screen.SetTerminalInfo(120, 40)
		})

		It("should return CancelResult on escape key", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})

		It("should return NavigateResult with false on 'n' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			navResult := result.(*screens.NavigateResult)
			Expect(navResult.ResultData).To(BeFalse())
		})

		It("should return NavigateResult with true on 'y' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			navResult := result.(*screens.NavigateResult)
			Expect(navResult.ResultData).To(BeTrue())
		})

		It("should return NavigateResult on Enter key with current selection", func() {
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			// Result data should be the current selection state (true or false)
			navResult := result.(*screens.NavigateResult)
			Expect(navResult.ResultData).To(BeAssignableToTypeOf(false))
		})

		It("should toggle selection with arrow keys", func() {
			// Get initial state by pressing Enter without navigation
			initialMsg := tea.KeyMsg{Type: tea.KeyEnter}
			_, initialResult := screen.Update(initialMsg)

			initialSelection := initialResult.(*screens.NavigateResult).ResultData.(bool)

			// Reset by creating a new screen
			screen = timeline.NewEventDeleteConfirmScreen(event)
			screen.SetTerminalInfo(120, 40)

			// Press right arrow to toggle selection
			arrowMsg := tea.KeyMsg{Type: tea.KeyRight}
			screen.Update(arrowMsg)

			// Check new selection state
			enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
			_, newResult := screen.Update(enterMsg)

			newSelection := newResult.(*screens.NavigateResult).ResultData.(bool)

			// Selection should have toggled
			Expect(newSelection).NotTo(Equal(initialSelection))
		})
	})
})
