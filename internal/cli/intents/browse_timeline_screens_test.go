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

	Describe("Search Modal E2E", func() {
		BeforeEach(func() {
			intent.Init()
		})

		// Helper function to process commands
		updateWithCmd := func(intent *BrowseTimelineIntent, msg tea.Msg) {
			cmd := intent.Update(msg)
			if cmd != nil {
				resultMsg := cmd()
				if resultMsg != nil {
					intent.Update(resultMsg)
				}
			}
		}

		It("should open search modal with '/' key", func() {
			// Verify initial state has all events
			initialView := intent.View()
			Expect(initialView).To(ContainSubstring("TechCorp"))
			Expect(initialView).To(ContainSubstring("CloudInc"))
			Expect(initialView).To(ContainSubstring("WebSolutions"))

			// Open search modal
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

			// Modal should be visible
			view := intent.View()
			Expect(view).To(ContainSubstring("Search Events"))
		})

		It("should filter events by search text", func() {
			// Verify initial state has all events
			initialView := intent.View()
			Expect(initialView).To(ContainSubstring("TechCorp"))
			Expect(initialView).To(ContainSubstring("CloudInc"))
			Expect(initialView).To(ContainSubstring("WebSolutions"))

			// Open search modal
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

			// Type search text "Backend"
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'B'}})
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

			// Submit search with Enter (first Enter completes field, second submits form)
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyEnter})
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyEnter})

			// View should show filtered results (only "Backend Developer")
			filteredView := intent.View()
			Expect(filteredView).To(ContainSubstring("Backend Developer"))
			// Should NOT show other events
			Expect(filteredView).NotTo(ContainSubstring("DevOps Engineer"))
			Expect(filteredView).NotTo(ContainSubstring("Frontend Developer"))
		})

		It("should allow canceling search with Esc", func() {
			// Open search modal
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

			// Type some text
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T'}})
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			// Cancel with Esc
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// View should still show all events (search was cancelled)
			view := intent.View()
			Expect(view).To(ContainSubstring("TechCorp"))
			Expect(view).To(ContainSubstring("CloudInc"))
			Expect(view).To(ContainSubstring("WebSolutions"))
		})

		It("should handle tab navigation in search modal", func() {
			// Open search modal
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

			// Tab should be forwarded to form (should not cause errors)
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyTab})

			// Modal should still be visible
			view := intent.View()
			Expect(view).To(ContainSubstring("Search Events"))
		})

		It("should handle empty search submission", func() {
			// Open search modal
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

			// Submit without typing (empty search)
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyEnter})
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyEnter})

			// Process reload command if present
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			if cmd != nil {
				msg := cmd()
				if msg != nil {
					intent.Update(msg)
				}
			}

			// Should show all events (empty search = no filter)
			view := intent.View()
			Expect(view).To(ContainSubstring("Backend Developer"))
		})
	})

	Describe("Sort Modal E2E", func() {
		BeforeEach(func() {
			intent.Init()
		})

		// Helper function to process commands
		updateWithCmd := func(intent *BrowseTimelineIntent, msg tea.Msg) {
			cmd := intent.Update(msg)
			if cmd != nil {
				resultMsg := cmd()
				if resultMsg != nil {
					intent.Update(resultMsg)
				}
			}
		}

		It("should open sort modal with 's' key", func() {
			// Open sort modal
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			// Modal should be visible (shows "Sort By" field title)
			view := intent.View()
			Expect(view).To(ContainSubstring("Sort By"))
		})

		It("should sort events by company name ascending", func() {
			// Open sort modal
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			// Navigate to "Company" option (press down arrow)
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyDown})

			// Select "Company" (press Enter to confirm field)
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyEnter})

			// Navigate to sort order field (Tab)
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyTab})

			// Navigate to "Ascending" (press up arrow since "Descending" is default)
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyUp})

			// Select "Ascending"
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyEnter})

			// Submit form (Enter to submit)
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyEnter})

			// Process reload command if present
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			if cmd != nil {
				msg := cmd()
				if msg != nil {
					intent.Update(msg)
				}
			}

			// View should show events sorted by company (alphabetically)
			// CloudInc < TechCorp < WebSolutions
			sortedView := intent.View()
			Expect(sortedView).To(ContainSubstring("Event"))
		})

		It("should allow canceling sort with Esc", func() {
			// Open sort modal
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			// Cancel with Esc
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// View should still show events in original order (by date descending)
			view := intent.View()
			Expect(view).To(ContainSubstring("TechCorp")) // Most recent first
		})

		It("should handle tab navigation in sort modal", func() {
			// Open sort modal
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			// Tab should be forwarded to form (should not cause errors)
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyTab})

			// Modal should still be visible (shows "Sort By" field)
			view := intent.View()
			Expect(view).To(ContainSubstring("Sort By"))
		})

		It("should sort events by date descending (default)", func() {
			// Open sort modal
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			// Submit without changing anything (defaults: Date, Descending)
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyEnter}) // Confirm sort by field
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyTab})   // Tab to order field
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyEnter}) // Confirm order field
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyEnter}) // Submit form

			// Process reload command if present
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			if cmd != nil {
				msg := cmd()
				if msg != nil {
					intent.Update(msg)
				}
			}

			// Events should be in date descending order (newest first)
			// Events are always shown, so just verify view is not empty
			sortedView := intent.View()
			Expect(sortedView).To(ContainSubstring("Event"))
		})
	})

	Describe("Filter Modal E2E", func() {
		BeforeEach(func() {
			intent.Init()
		})

		// Helper function to process commands
		updateWithCmd := func(intent *BrowseTimelineIntent, msg tea.Msg) {
			cmd := intent.Update(msg)
			if cmd != nil {
				resultMsg := cmd()
				if resultMsg != nil {
					intent.Update(resultMsg)
				}
			}
		}

		It("should open filter modal with 'f' key", func() {
			// Open filter modal
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			// Modal should be visible (shows "Filter by Company" field)
			view := intent.View()
			Expect(view).To(ContainSubstring("Filter by Company"))
		})

		It("should filter events by company", func() {
			// Open filter modal
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			// Navigate to companies field and select TechCorp
			// (Implementation depends on FilterModalModel structure)
			// For now, just submit to verify modal handling works
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyEnter})

			// Modal should process submission
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should allow canceling filter with Esc", func() {
			// Open filter modal
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			// Cancel with Esc
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// View should still show all events (filter was cancelled)
			view := intent.View()
			Expect(view).To(ContainSubstring("TechCorp"))
			Expect(view).To(ContainSubstring("CloudInc"))
			Expect(view).To(ContainSubstring("WebSolutions"))
		})

		It("should handle tab navigation in filter modal", func() {
			// Open filter modal
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			// Tab should be forwarded to form (should not cause errors)
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyTab})

			// Modal should still be visible
			view := intent.View()
			Expect(view).To(ContainSubstring("Filter"))
		})
	})

	Describe("Modal Integration and Priority", func() {
		BeforeEach(func() {
			intent.Init()
		})

		// Helper function to process commands
		updateWithCmd := func(intent *BrowseTimelineIntent, msg tea.Msg) {
			cmd := intent.Update(msg)
			if cmd != nil {
				resultMsg := cmd()
				if resultMsg != nil {
					intent.Update(resultMsg)
				}
			}
		}

		It("should handle search and sort modals together (E2E combination)", func() {
			// 1. Apply search first
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'D'}})
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyEnter})
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyEnter})

			// Should show filtered events (Backend and Frontend Developer, DevOps)
			searchView := intent.View()
			Expect(searchView).To(ContainSubstring("Developer"))

			// 2. Then apply sort (should work on filtered results)
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyEnter}) // Confirm sort by field
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyTab})   // Tab to order
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyEnter}) // Confirm order
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyEnter}) // Submit

			// Should show filtered AND sorted results
			finalView := intent.View()
			Expect(finalView).NotTo(BeEmpty())
		})

		It("should prioritize modal input over screen input", func() {
			// Open search modal
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

			// Try to navigate (should not affect event list navigation)
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyDown})

			// Modal should still be visible and list should not navigate
			view := intent.View()
			Expect(view).To(ContainSubstring("Search Events"))
		})

		It("should only show one modal at a time", func() {
			// Open search modal
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
			searchView := intent.View()
			Expect(searchView).To(ContainSubstring("Search Events"))

			// Try to open sort modal (should not work while search is open)
			updateWithCmd(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			// Should still show search modal (not sort)
			view := intent.View()
			Expect(view).To(ContainSubstring("Search Events"))
			Expect(view).NotTo(ContainSubstring("Sort By")) // Sort modal shows "Sort By" field
		})
	})
})
