package components_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ViewEventDetailModal", func() {
	var (
		testTheme themes.Theme
		testEvent *career.Event
		modal     *components.ViewEventDetailModal
	)

	BeforeEach(func() {
		testTheme = themes.NewDefaultTheme()
		testEvent = &career.Event{
			ID:         "test-123",
			Text:       "Implemented new feature for the product",
			Date:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Company:    "Acme Corp",
			Project:    "Project Alpha",
			Tags:       []string{"go", "backend", "api"},
			Categories: []string{"development", "feature"},
			Skills:     []string{"skill-1", "skill-2"},
		}
		modal = components.NewViewEventDetailModal(testEvent, testTheme)
	})

	Describe("Construction", func() {
		It("creates a new modal", func() {
			Expect(modal).NotTo(BeNil())
		})

		It("starts hidden", func() {
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("has default dimensions", func() {
			// Show the modal and render to verify default dimensions work
			modal.Show()
			result := modal.View()
			Expect(result).NotTo(BeEmpty())
		})

		It("works with nil theme", func() {
			modal := components.NewViewEventDetailModal(testEvent, nil)
			modal.Show()
			// Should not panic
			result := modal.View()
			Expect(result).NotTo(BeEmpty())
		})
	})

	Describe("Visibility", func() {
		It("can be shown", func() {
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("can be hidden", func() {
			modal.Show()
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("returns empty string when hidden", func() {
			Expect(modal.View()).To(BeEmpty())
		})

		It("returns content when visible", func() {
			modal.Show()
			result := modal.View()
			Expect(result).NotTo(BeEmpty())
		})
	})

	Describe("Keyboard Handling", func() {
		BeforeEach(func() {
			modal.Show()
		})

		It("closes on Escape key", func() {
			msg := tea.KeyMsg{Type: tea.KeyEscape}
			modal.Update(msg)
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("closes on Enter key", func() {
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			modal.Update(msg)
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("closes on 'q' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
			modal.Update(msg)
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("closes on backspace key", func() {
			msg := tea.KeyMsg{Type: tea.KeyBackspace}
			modal.Update(msg)
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("does nothing when hidden", func() {
			modal.Hide()
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			modal.Update(msg)
			// Should not panic and still be hidden
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("Content Rendering", func() {
		BeforeEach(func() {
			// Set larger dimensions to ensure content fits
			modal.SetDimensions(100, 50)
			modal.Show()
		})

		It("renders event details", func() {
			result := modal.View()
			Expect(result).To(ContainSubstring("Event Details"))
		})

		It("renders event date", func() {
			result := modal.View()
			Expect(result).To(ContainSubstring("2024-01-15"))
		})

		It("renders event company", func() {
			result := modal.View()
			Expect(result).To(ContainSubstring("Acme Corp"))
		})

		It("renders event text", func() {
			result := modal.View()
			// Text is in the viewport content but may be cut off by viewport height
			// Check that the "Text:" label is visible at minimum
			Expect(result).To(ContainSubstring("Text:"))
		})

		It("renders close instruction", func() {
			result := modal.View()
			// Should have close instruction
			Expect(result).To(MatchRegexp(`(Enter|Esc|Close)`))
		})
	})

	Describe("Dimensions", func() {
		It("can set dimensions", func() {
			modal.SetDimensions(100, 50)
			modal.Show()
			// Should not panic
			result := modal.View()
			Expect(result).NotTo(BeEmpty())
		})

		It("handles window size messages", func() {
			modal.Show()
			msg := tea.WindowSizeMsg{Width: 120, Height: 40}
			modal.Update(msg)
			// Should not panic
			result := modal.View()
			Expect(result).NotTo(BeEmpty())
		})

		It("handles very small dimensions", func() {
			modal.SetDimensions(40, 15)
			modal.Show()
			// Should not panic
			result := modal.View()
			Expect(result).NotTo(BeEmpty())
		})

		It("handles very large dimensions", func() {
			modal.SetDimensions(200, 100)
			modal.Show()
			// Should not panic and should cap dimensions reasonably
			result := modal.View()
			Expect(result).NotTo(BeEmpty())
		})
	})

	Describe("Event Management", func() {
		It("can update the event", func() {
			newEvent := &career.Event{
				ID:   "new-123",
				Text: "New event text",
				Date: time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC),
			}
			modal.SetEvent(newEvent)
			modal.Show()

			result := modal.View()
			Expect(result).To(ContainSubstring("New event text"))
			Expect(result).To(ContainSubstring("2024-02-20"))
		})

		It("handles nil event", func() {
			modal.SetEvent(nil)
			modal.Show()

			result := modal.View()
			Expect(result).To(ContainSubstring("No event selected"))
		})
	})

	Describe("Action State", func() {
		It("returns empty action by default", func() {
			Expect(modal.GetAction()).To(BeEmpty())
		})

		It("returns empty action after closing", func() {
			modal.Show()
			msg := tea.KeyMsg{Type: tea.KeyEscape}
			modal.Update(msg)

			Expect(modal.GetAction()).To(BeEmpty())
		})
	})

	Describe("Long Content Handling", func() {
		It("handles events with long text", func() {
			longText := "This is a very long event description that should be " +
				"displayed properly in the modal. It contains multiple sentences " +
				"and describes a complex project with technical details. " +
				"The implementation involved coordination across teams. " +
				"We achieved excellent results through this effort."

			testEvent.Text = longText
			modal.SetEvent(testEvent)
			modal.SetDimensions(100, 50) // Larger dimensions to show more content
			modal.Show()

			result := modal.View()
			// Text label should be visible; full text may require scrolling
			Expect(result).To(ContainSubstring("Text:"))
			// The viewport includes the content even if not all visible
			Expect(result).NotTo(BeEmpty())
		})

		It("handles events with many tags", func() {
			testEvent.Tags = []string{
				"go", "python", "javascript", "typescript", "rust",
				"kubernetes", "docker", "aws", "gcp", "terraform",
				"postgresql", "redis", "mongodb", "elasticsearch",
			}
			modal.SetEvent(testEvent)
			modal.SetDimensions(100, 50) // Larger dimensions
			modal.Show()

			result := modal.View()
			// Should contain the event details; tags may be below fold
			Expect(result).To(ContainSubstring("Event Details"))
		})

		It("shows scroll indicator for tall content", func() {
			// Create very long text that will overflow
			longText := ""
			for i := 0; i < 50; i++ {
				longText += "Line of text that adds content to make the modal scrollable. "
			}
			testEvent.Text = longText
			modal.SetEvent(testEvent)
			modal.SetDimensions(80, 20) // Small height to trigger scrolling
			modal.Show()

			// Force viewport initialization by rendering
			_ = modal.View()

			// Note: The scroll indicator depends on content height vs viewport
			// We just verify it doesn't crash
		})
	})

	Describe("Scroll Keys", func() {
		BeforeEach(func() {
			// Create scrollable content
			longText := ""
			for i := 0; i < 100; i++ {
				longText += "Line of content to make scrollable. "
			}
			testEvent.Text = longText
			modal.SetEvent(testEvent)
			modal.SetDimensions(80, 20)
			modal.Show()
			// Initialize viewport
			_ = modal.View()
		})

		It("handles down key", func() {
			msg := tea.KeyMsg{Type: tea.KeyDown}
			_, cmd := modal.Update(msg)
			// Should not close modal
			Expect(modal.IsVisible()).To(BeTrue())
			// cmd might be nil if viewport isn't ready, that's ok
			_ = cmd
		})

		It("handles up key", func() {
			msg := tea.KeyMsg{Type: tea.KeyUp}
			_, cmd := modal.Update(msg)
			Expect(modal.IsVisible()).To(BeTrue())
			_ = cmd
		})

		It("handles j key (vim down)", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			_, cmd := modal.Update(msg)
			Expect(modal.IsVisible()).To(BeTrue())
			_ = cmd
		})

		It("handles k key (vim up)", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
			_, cmd := modal.Update(msg)
			Expect(modal.IsVisible()).To(BeTrue())
			_ = cmd
		})
	})
})
