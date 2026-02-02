package base_test

import (
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// SelectScreen Tests
//
// SelectScreen[T] is a generic screen for selecting items from a list.
// It handles navigation (↑/↓/j/k/g/G), selection (Enter), and cancellation (Esc).
//
// Related:
// - tasks/tasks-42-tui-architecture-refactor.md (Phase 1.2)
// - docs/KEYBOARD_SHORTCUTS_GUIDE.md (navigation keys)

var _ = Describe("SelectScreen", func() {
	var (
		screen *base.SelectScreen[string]
		items  []string
	)

	BeforeEach(func() {
		items = []string{"Item A", "Item B", "Item C", "Item D", "Item E"}
		screen = base.NewBaseSelectScreen[string](
			items,
			func(s string) string { return s },
			[]string{"Main", "Select"},
			"Select an Item",
		)
	})

	Describe("Creation", func() {
		It("should create with items and default selection at index 0", func() {
			view := screen.View()

			Expect(view).To(ContainSubstring("Item A"))
			Expect(view).To(ContainSubstring("▶")) // Selection indicator
		})

		It("should accept generic type parameter for items", func() {
			// Test with string type (already done in BeforeEach)
			Expect(screen).NotTo(BeNil())

			// Test with custom struct type
			type CustomItem struct {
				Name string
				ID   int
			}

			customItems := []CustomItem{
				{Name: "First", ID: 1},
				{Name: "Second", ID: 2},
			}
			customScreen := base.NewBaseSelectScreen[CustomItem](
				customItems,
				func(c CustomItem) string { return c.Name },
				[]string{"Test"},
				"Custom Items",
			)

			view := customScreen.View()
			Expect(view).To(ContainSubstring("First"))
		})

		It("should accept item renderer function", func() {
			// Custom renderer that adds prefix
			rendererScreen := base.NewBaseSelectScreen[string](
				[]string{"One", "Two"},
				func(s string) string { return ">>> " + s },
				[]string{"Test"},
				"Custom Render",
			)

			view := rendererScreen.View()
			Expect(view).To(ContainSubstring(">>> One"))
			Expect(view).To(ContainSubstring(">>> Two"))
		})

		It("should initialize with provided breadcrumbs", func() {
			view := screen.View()

			Expect(view).To(ContainSubstring("Main"))
			Expect(view).To(ContainSubstring("Select"))
		})

		It("should initialize with provided title", func() {
			view := screen.View()

			Expect(view).To(ContainSubstring("Select an Item"))
		})
	})

	Describe("Navigation - Arrow Keys", func() {
		It("should navigate down with Down arrow key", func() {
			// Initial: index 0 (Item A selected)
			msg := tea.KeyMsg{Type: tea.KeyDown}

			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())

			// View should now have Item B selected
			view := screen.View()
			// Item B should be the selected item
			Expect(view).To(ContainSubstring("Item B"))
		})

		It("should navigate up with Up arrow key", func() {
			// First move down
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})

			// Then move back up
			msg := tea.KeyMsg{Type: tea.KeyUp}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should not go below first item (boundary check)", func() {
			// Press Up multiple times at index 0
			for i := range 5 {
				_ = i
				screen.Update(tea.KeyMsg{Type: tea.KeyUp})
			}

			// Should still be at first item
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Data()).To(Equal("Item A"))
		})

		It("should not go above last item (boundary check)", func() {
			// Press Down multiple times past the end
			for i := range 10 {
				_ = i
				screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			}

			// Should be at last item
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Data()).To(Equal("Item E"))
		})
	})

	Describe("Navigation - Vim Keys", func() {
		It("should navigate down with 'j' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())

			// Verify moved down by selecting
			_, selectResult := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(selectResult.Data()).To(Equal("Item B"))
		})

		It("should navigate up with 'k' key", func() {
			// First move down
			screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

			// Then move up with k
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())

			// Verify back at first item
			_, selectResult := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(selectResult.Data()).To(Equal("Item A"))
		})

		It("should jump to top with 'g' key", func() {
			// Move to somewhere in the middle
			for i := range 3 {
				_ = i
				screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			}

			// Jump to top
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
			screen.Update(msg)

			// Verify at first item
			_, selectResult := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(selectResult.Data()).To(Equal("Item A"))
		})

		It("should jump to bottom with 'G' key", func() {
			// Jump to bottom
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}}
			screen.Update(msg)

			// Verify at last item
			_, selectResult := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(selectResult.Data()).To(Equal("Item E"))
		})
	})

	Describe("Selection - Enter Key", func() {
		It("should return NavigateResult when Enter is pressed", func() {
			msg := tea.KeyMsg{Type: tea.KeyEnter}

			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
		})

		It("should include selected item in result data", func() {
			// Move to Item C (index 2)
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})

			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := screen.Update(msg)

			Expect(result.Data()).To(Equal("Item C"))
		})

		It("should include selected index in metadata", func() {
			// Move to Item C (index 2)
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})

			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := screen.Update(msg)

			Expect(result.Metadata()["selected_index"]).To(Equal(2))
		})

		It("should not return result for navigation keys", func() {
			msgs := []tea.Msg{
				tea.KeyMsg{Type: tea.KeyDown},
				tea.KeyMsg{Type: tea.KeyUp},
				tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
				tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}},
			}

			for _, msg := range msgs {
				_, result := screen.Update(msg)
				Expect(result).To(BeNil(), "Navigation key should not return result")
			}
		})
	})

	Describe("Cancellation - Escape Key", func() {
		It("should return CancelResult when Esc is pressed", func() {
			msg := tea.KeyMsg{Type: tea.KeyEscape}

			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})

		It("should include current index in cancel metadata", func() {
			for range 3 {
				screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			}

			msg := tea.KeyMsg{Type: tea.KeyEscape}
			_, result := screen.Update(msg)

			Expect(result.Metadata()["selected_index"]).To(Equal(3))
		})
	})

	Describe("View Rendering", func() {
		It("should render items with selection indicator", func() {
			view := screen.View()

			Expect(view).To(ContainSubstring("▶"))
			Expect(view).To(ContainSubstring("Item A"))
		})

		It("should use StandardView for consistent layout", func() {
			view := screen.View()

			// Should have breadcrumbs, content, and footer
			Expect(view).To(ContainSubstring("Main"))
			Expect(view).To(ContainSubstring("Select"))
		})

		It("should show breadcrumbs in view", func() {
			customScreen := base.NewBaseSelectScreen[string](
				[]string{"Item"},
				func(s string) string { return s },
				[]string{"Level1", "Level2", "Level3"},
				"Test",
			)

			view := customScreen.View()

			Expect(view).To(ContainSubstring("Level1"))
			Expect(view).To(ContainSubstring("Level2"))
			Expect(view).To(ContainSubstring("Level3"))
		})

		It("should show title in view", func() {
			view := screen.View()

			Expect(view).To(ContainSubstring("Select an Item"))
		})

		It("should show footer with navigation hints", func() {
			view := screen.View()

			Expect(view).To(ContainSubstring("Navigate"))
			Expect(view).To(ContainSubstring("Select"))
			Expect(view).To(ContainSubstring("Back"))
		})

		It("should call item renderer for each item", func() {
			renderCount := 0
			countingScreen := base.NewBaseSelectScreen[string](
				[]string{"A", "B", "C"},
				func(s string) string {
					renderCount++
					return "Rendered: " + s
				},
				[]string{"Test"},
				"Test",
			)

			view := countingScreen.View()

			// Should have rendered all visible items
			Expect(view).To(ContainSubstring("Rendered: A"))
			Expect(view).To(ContainSubstring("Rendered: B"))
			Expect(view).To(ContainSubstring("Rendered: C"))
		})
	})

	Describe("Window Size Handling", func() {
		It("should handle WindowSizeMsg via Screen", func() {
			msg := tea.WindowSizeMsg{Width: 100, Height: 50}

			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should update view dimensions when terminal resizes", func() {
			msg := tea.WindowSizeMsg{Width: 80, Height: 24}
			screen.Update(msg)

			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should not return ScreenResult for WindowSizeMsg", func() {
			msg := tea.WindowSizeMsg{Width: 120, Height: 40}

			_, result := screen.Update(msg)

			Expect(result).To(BeNil())
		})
	})

	Describe("Empty List Handling", func() {
		var emptyScreen *base.SelectScreen[string]

		BeforeEach(func() {
			emptyScreen = base.NewBaseSelectScreen[string](
				[]string{},
				func(s string) string { return s },
				[]string{"Test"},
				"Empty List",
			)
		})

		It("should handle empty item list gracefully", func() {
			view := emptyScreen.View()

			Expect(view).To(ContainSubstring("No items"))
		})

		It("should not panic when navigating empty list", func() {
			// Should not panic
			emptyScreen.Update(tea.KeyMsg{Type: tea.KeyDown})
			emptyScreen.Update(tea.KeyMsg{Type: tea.KeyUp})
			emptyScreen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			emptyScreen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			emptyScreen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
			emptyScreen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
		})

		It("should allow cancellation with empty list", func() {
			msg := tea.KeyMsg{Type: tea.KeyEscape}

			_, result := emptyScreen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})

		It("should not allow selection with empty list", func() {
			msg := tea.KeyMsg{Type: tea.KeyEnter}

			_, result := emptyScreen.Update(msg)

			Expect(result).To(BeNil())
		})
	})

	Describe("Single Item List", func() {
		var singleItemScreen *base.SelectScreen[string]

		BeforeEach(func() {
			singleItemScreen = base.NewBaseSelectScreen[string](
				[]string{"Only Item"},
				func(s string) string { return s },
				[]string{"Test"},
				"Single Item",
			)
		})

		It("should show single item without navigation", func() {
			view := singleItemScreen.View()

			Expect(view).To(ContainSubstring("Only Item"))
		})

		It("should select single item on Enter", func() {
			msg := tea.KeyMsg{Type: tea.KeyEnter}

			_, result := singleItemScreen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Data()).To(Equal("Only Item"))
		})

		It("should not change selection with navigation keys", func() {
			// Try navigating
			singleItemScreen.Update(tea.KeyMsg{Type: tea.KeyDown})
			singleItemScreen.Update(tea.KeyMsg{Type: tea.KeyUp})

			// Should still select the only item
			_, result := singleItemScreen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(result.Data()).To(Equal("Only Item"))
		})
	})

	Describe("Large List Handling", func() {
		var largeScreen *base.SelectScreen[string]

		BeforeEach(func() {
			largeItems := make([]string, 50)
			for i := range 50 {
				largeItems[i] = "Item " + string(rune('A'+i%26)) + string(rune('0'+i/26))
			}

			largeScreen = base.NewBaseSelectScreen[string](
				largeItems,
				func(s string) string { return s },
				[]string{"Test"},
				"Large List",
			)
		})

		It("should handle lists larger than terminal height", func() {
			view := largeScreen.View()

			// Should render without error
			Expect(view).NotTo(BeEmpty())
		})

		It("should show visible window of items", func() {
			view := largeScreen.View()

			// First item should be visible
			Expect(view).To(ContainSubstring("Item A0"))
		})

		It("should scroll as user navigates", func() {
			for range 30 {
				largeScreen.Update(tea.KeyMsg{Type: tea.KeyDown})
			}

			view := largeScreen.View()

			// Should show scroll indicator
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("ItemRenderer Function", func() {
		It("should call renderer for each item in view", func() {
			called := make(map[string]bool)
			rendererScreen := base.NewBaseSelectScreen[string](
				[]string{"A", "B", "C"},
				func(s string) string {
					called[s] = true
					return s
				},
				[]string{"Test"},
				"Test",
			)

			rendererScreen.View()

			Expect(called["A"]).To(BeTrue())
			Expect(called["B"]).To(BeTrue())
			Expect(called["C"]).To(BeTrue())
		})

		It("should handle multi-line item renderings", func() {
			multiLineScreen := base.NewBaseSelectScreen[string](
				[]string{"Item"},
				func(s string) string { return s + "\n  Subtitle" },
				[]string{"Test"},
				"Test",
			)

			view := multiLineScreen.View()

			Expect(view).To(ContainSubstring("Item"))
			Expect(view).To(ContainSubstring("Subtitle"))
		})

		It("should handle renderer returning empty string", func() {
			emptyRenderScreen := base.NewBaseSelectScreen[string](
				[]string{"A", "B"},
				func(_ string) string { return "" },
				[]string{"Test"},
				"Test",
			)

			view := emptyRenderScreen.View()

			// Should still render (with selection indicator)
			Expect(view).To(ContainSubstring("▶"))
		})
	})

	Describe("State Preservation", func() {
		It("should restore selection index from metadata", func() {
			// Use WithInitialSelection
			restoredScreen := base.NewBaseSelectScreen[string](
				items,
				func(s string) string { return s },
				[]string{"Test"},
				"Test",
			).WithInitialSelection(3)

			// Verify selection is at index 3
			_, result := restoredScreen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(result.Data()).To(Equal("Item D"))
		})

		It("should accept initial selection index", func() {
			initialScreen := base.NewBaseSelectScreen[string](
				items,
				func(s string) string { return s },
				[]string{"Test"},
				"Test",
			).WithInitialSelection(2)

			_, result := initialScreen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(result.Data()).To(Equal("Item C"))
		})

		It("should clamp out-of-bounds initial selection", func() {
			// Index too high
			highScreen := base.NewBaseSelectScreen[string](
				items,
				func(s string) string { return s },
				[]string{"Test"},
				"Test",
			).WithInitialSelection(100)

			_, result := highScreen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(result.Data()).To(Equal("Item E")) // Last item

			// Index negative
			lowScreen := base.NewBaseSelectScreen[string](
				items,
				func(s string) string { return s },
				[]string{"Test"},
				"Test",
			).WithInitialSelection(-5)

			_, result = lowScreen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(result.Data()).To(Equal("Item A")) // First item
		})
	})

	Describe("Integration with Screen", func() {
		It("should embed Screen for common functionality", func() {
			// Verify SetTerminalInfo works
			screen.SetTerminalInfo(80, 24)

			// Verify SetTheme works
			screen.SetTheme("test-theme")

			// Should render without error
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should use Screen.CreateView", func() {
			// View should have StandardView structure
			view := screen.View()

			// Has breadcrumbs
			Expect(view).To(ContainSubstring("Main"))

			// Has content
			Expect(view).To(ContainSubstring("Item A"))

			// Has footer
			Expect(view).To(ContainSubstring("Back"))
		})

		It("should use Screen.HandleWindowSizeMsg", func() {
			msg := tea.WindowSizeMsg{Width: 100, Height: 50}

			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})
	})
})
