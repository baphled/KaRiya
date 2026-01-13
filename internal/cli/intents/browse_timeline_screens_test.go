package intents

import (
	"time"

	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BrowseTimelineIntent - Screen Architecture", func() {
	var (
		intent  *BrowseTimelineIntent
		context *BrowseTimelineContext
		events  []*career.CareerEvent
	)

	BeforeEach(func() {
		// Create test events
		events = []*career.CareerEvent{
			{
				ID:      "event-1",
				Date:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Text:    "Backend Developer at TechCorp - Built scalable APIs",
				Company: "TechCorp",
			},
			{
				ID:      "event-2",
				Date:    time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
				Text:    "DevOps Engineer at CloudInc - Managed Kubernetes",
				Company: "CloudInc",
			},
			{
				ID:      "event-3",
				Date:    time.Date(2022, 3, 10, 0, 0, 0, 0, time.UTC),
				Text:    "Frontend Developer at WebSolutions - React apps",
				Company: "WebSolutions",
			},
		}

		// Create context
		context = &BrowseTimelineContext{
			Events:          events,
			CLIEventService: &service.CLIEventService{},
		}

		// Create intent (active by default)
		var err error
		intent, err = NewBrowseTimelineIntent(context)
		Expect(err).NotTo(HaveOccurred())

		// Enable screen-based architecture
		intent.EnableScreens()
	})

	Describe("Initialization with Screens", func() {
		It("should initialize with timeline list screen", func() {
			intent.Init()

			// View should render timeline list screen with StandardView breadcrumbs
			view := intent.View()
			Expect(view).To(ContainSubstring("Timeline")) // Breadcrumb text
			Expect(view).To(ContainSubstring("Backend Developer"))
		})

		It("should show event count", func() {
			intent.Init()

			view := intent.View()
			Expect(view).To(ContainSubstring("Events: 3")) // StandardView footer format
		})

		It("should show all events in list", func() {
			intent.Init()

			view := intent.View()
			Expect(view).To(ContainSubstring("TechCorp"))
			Expect(view).To(ContainSubstring("CloudInc"))
			Expect(view).To(ContainSubstring("WebSolutions"))
		})
	})

	Describe("Navigation with Screens", func() {
		BeforeEach(func() {
			intent.Init()
		})

		It("should navigate down with arrow key", func() {
			// Initial view should show first event selected
			view := intent.View()
			Expect(view).To(ContainSubstring("▶"))

			// Navigate down
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})

			// View should update
			view = intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should navigate up with arrow key", func() {
			// Navigate down first
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})

			// Navigate up
			intent.Update(tea.KeyMsg{Type: tea.KeyUp})

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should support vim-style navigation (j/k)", func() {
			// Navigate with j (down)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

			view := intent.View()
			Expect(view).NotTo(BeEmpty())

			// Navigate with k (up)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

			view = intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Event Selection with Modal", func() {
		BeforeEach(func() {
			intent.Init()
		})

		It("should show ViewEventDetailModal on enter", func() {
			// Press enter to select event
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// View should show modal overlay with event details
			view := intent.View()
			// Modal should overlay the list (list content still present)
			Expect(view).To(ContainSubstring("Timeline"))
			// Modal should contain event data
			Expect(view).To(ContainSubstring("Backend Developer"))
			Expect(view).To(ContainSubstring("TechCorp"))
		})

		It("should show full event information in modal", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			view := intent.View()
			// Modal renders event fields
			Expect(view).To(ContainSubstring("2024-01-01"))
			Expect(view).To(ContainSubstring("TechCorp"))
			Expect(view).To(ContainSubstring("Backend Developer"))
		})

		It("should close modal and return to list with escape", func() {
			// Show modal
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			modalView := intent.View()
			Expect(modalView).To(ContainSubstring("Backend Developer"))

			// Press escape to close modal
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should be back at list view without modal
			listView := intent.View()
			Expect(listView).To(ContainSubstring("Timeline")) // Breadcrumb still present
			// List should show all events again
			Expect(listView).To(ContainSubstring("DevOps Engineer"))
		})

		It("should preserve list state after closing modal", func() {
			// Navigate to second event
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})

			// Show modal for second event
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Close modal
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// List should still show all events
			view := intent.View()
			Expect(view).To(ContainSubstring("Timeline"))
			Expect(view).To(ContainSubstring("Backend Developer"))
			Expect(view).To(ContainSubstring("DevOps Engineer"))
		})
	})

	Describe("Actions with Screens", func() {
		BeforeEach(func() {
			intent.Init()
		})

		It("should handle add action from list", func() {
			// Press 'a' for add
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			// Action is triggered (TODO: routing not implemented yet)
			// Intent should remain active, not completed
			// We just verify it doesn't crash
		})

		It("should handle edit action from list", func() {
			// Press 'e' for edit
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			// Action is triggered (TODO: routing not implemented yet)
			// We just verify it doesn't crash
		})

		It("should handle delete action from list", func() {
			// Press 'd' for delete
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

			// Action is triggered (TODO: routing not implemented yet)
			// We just verify it doesn't crash
		})

		It("should handle edit action from detail view", func() {
			// Go to detail view
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Press 'e' for edit
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			// Action is triggered (TODO: routing not implemented yet)
			// We just verify it doesn't crash
		})

		It("should handle delete action from detail view", func() {
			// Go to detail view
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Press 'd' for delete
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

			// Action is triggered (TODO: routing not implemented yet)
			// We just verify it doesn't crash
		})
	})

	Describe("Cancellation with Screens", func() {
		BeforeEach(func() {
			intent.Init()
		})

		It("should cancel intent from list with escape", func() {
			// Press escape from root state (list)
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Intent should be cancelled
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Cancelled))
		})

		It("should quit app from list with q", func() {
			// Press 'q' to quit app (global key handled by intent)
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

			// Should return tea.Quit command (quits entire app, not just intent)
			Expect(cmd).NotTo(BeNil())
			// tea.Quit is a function, we can't directly compare it
			// but we can verify it's not nil which means 'q' was handled
		})

		It("should not cancel intent from detail with escape", func() {
			// Go to detail view
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Press escape (should go back to list, not cancel)
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should be back at list with breadcrumbs, not cancelled
			view := intent.View()
			Expect(view).To(ContainSubstring("Timeline")) // Breadcrumb text
		})
	})

	Describe("Terminal Resize with Screens", func() {
		BeforeEach(func() {
			intent.Init()
		})

		It("should handle window size message", func() {
			msg := tea.WindowSizeMsg{Width: 80, Height: 24}
			intent.Update(msg)

			// Should still render correctly
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should propagate size to active screen", func() {
			// Resize
			intent.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

			// Screen should still render with breadcrumbs
			view := intent.View()
			Expect(view).To(ContainSubstring("Timeline")) // Breadcrumb text
		})
	})

	Describe("Empty Event List with Screens", func() {
		BeforeEach(func() {
			// Create intent with no events
			context = &BrowseTimelineContext{
				Events:          []*career.CareerEvent{},
				CLIEventService: &service.CLIEventService{},
			}

			var err error
			intent, err = NewBrowseTimelineIntent(context)
			Expect(err).NotTo(HaveOccurred())

			intent.EnableScreens()
			intent.Init()
		})

		It("should show empty state message", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("No events"))
		})

		It("should prompt to add event", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Add")) // KeyBadge format (capital A)
		})

		It("should not crash on navigation", func() {
			// Try to navigate (should do nothing gracefully)
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			intent.Update(tea.KeyMsg{Type: tea.KeyUp})
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Should still render empty state
			view := intent.View()
			Expect(view).To(ContainSubstring("No events"))
		})
	})

	Describe("Screen vs Legacy Mode", func() {
		It("should use legacy mode when screens not enabled", func() {
			// Create new intent without enabling screens
			legacyIntent, err := NewBrowseTimelineIntent(context)
			Expect(err).NotTo(HaveOccurred())
			legacyIntent.Init()

			// View should use legacy table-based rendering
			view := legacyIntent.View()
			Expect(view).NotTo(BeEmpty())
			// Legacy view won't have screen-specific markers like "Career Timeline"
		})

		It("should use screen mode when screens enabled", func() {
			// Already enabled in BeforeEach
			intent.Init()

			view := intent.View()
			Expect(view).To(ContainSubstring("Timeline")) // Breadcrumb text in screen mode
		})
	})
})
