package base_test

import (
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestDetailData is a simple data structure for testing detail views.
type TestDetailData struct {
	Title       string
	Description string
	Tags        []string
	Metadata    map[string]string
}

var _ = Describe("DetailScreen", func() {
	var (
		screen     *base.DetailScreen[*TestDetailData]
		detailData *TestDetailData
	)

	BeforeEach(func() {
		detailData = &TestDetailData{
			Title:       "Test Item",
			Description: "This is a test description",
			Tags:        []string{"tag1", "tag2", "tag3"},
			Metadata: map[string]string{
				"Author": "Test Author",
				"Date":   "2026-01-13",
			},
		}
	})

	Describe("Construction", func() {
		It("should create a detail screen with content renderer", func() {
			renderer := func(data *TestDetailData, _, _ int) string {
				return "Title: " + data.Title + "\n" + data.Description
			}

			screen = base.NewBaseDetailScreen(
				[]string{"Test", "Detail"},
				renderer,
				detailData,
			)

			Expect(screen).NotTo(BeNil())
			Expect(screen.GetData()).To(Equal(detailData))
		})

		It("should use default dimensions if not set", func() {
			renderer := func(_ *TestDetailData, _, _ int) string {
				return "content"
			}

			screen = base.NewBaseDetailScreen([]string{"Test"}, renderer, detailData)
			Expect(screen.Width()).To(Equal(120)) // Default width
			Expect(screen.Height()).To(Equal(40)) // Default height
		})

		It("should use custom footer if provided", func() {
			renderer := func(_ *TestDetailData, _, _ int) string {
				return "content"
			}

			screen = base.NewBaseDetailScreen([]string{"Test"}, renderer, detailData)
			screen.SetFooter("Custom footer")

			view := screen.View()
			Expect(view).To(ContainSubstring("Custom footer"))
		})
	})

	Describe("Terminal Info", func() {
		BeforeEach(func() {
			renderer := func(_ *TestDetailData, _, _ int) string {
				return "content"
			}
			screen = base.NewBaseDetailScreen([]string{"Test"}, renderer, detailData)
		})

		It("should update dimensions on SetTerminalInfo", func() {
			screen.SetTerminalInfo(100, 30)
			Expect(screen.Width()).To(Equal(100))
			Expect(screen.Height()).To(Equal(30))
		})

		It("should handle WindowSizeMsg", func() {
			msg := tea.WindowSizeMsg{Width: 80, Height: 24}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
			Expect(screen.Width()).To(Equal(80))
			Expect(screen.Height()).To(Equal(24))
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			renderer := func(_ *TestDetailData, _, _ int) string {
				return "content"
			}
			screen = base.NewBaseDetailScreen([]string{"Test"}, renderer, detailData)
		})

		It("should return CancelResult on escape key", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})

		It("should return NavigateResult with action on enter key", func() {
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			Expect(result.Data()).To(Equal("confirm"))
		})

		It("should support custom action keys", func() {
			renderer := func(_ *TestDetailData, _, _ int) string {
				return "content"
			}
			screen = base.NewBaseDetailScreen([]string{"Test"}, renderer, detailData)

			// Register custom action
			screen.AddAction("e", "edit")

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			Expect(result.Data()).To(Equal("edit"))
		})

		It("should support multiple custom actions", func() {
			renderer := func(_ *TestDetailData, _, _ int) string {
				return "content"
			}
			screen = base.NewBaseDetailScreen([]string{"Test"}, renderer, detailData)

			screen.AddAction("d", "delete")
			screen.AddAction("c", "copy")

			// Test delete action
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
			_, result := screen.Update(msg)
			Expect(result.Data()).To(Equal("delete"))

			// Test copy action
			msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
			_, result = screen.Update(msg)
			Expect(result.Data()).To(Equal("copy"))
		})
	})

	Describe("Scrolling", func() {
		BeforeEach(func() {
			renderer := func(_ *TestDetailData, _, _ int) string {
				// Generate multi-line content for scrolling
				content := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\n"
				content += "Line 6\nLine 7\nLine 8\nLine 9\nLine 10\n"
				return content
			}
			screen = base.NewBaseDetailScreen([]string{"Test"}, renderer, detailData)
			screen.SetTerminalInfo(80, 10) // Small height to enable scrolling
		})

		It("should scroll down on down arrow", func() {
			initialOffset := screen.GetScrollOffset()

			msg := tea.KeyMsg{Type: tea.KeyDown}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
			Expect(screen.GetScrollOffset()).To(Equal(initialOffset + 1))
		})

		It("should scroll down on j key (vim-style)", func() {
			initialOffset := screen.GetScrollOffset()

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
			Expect(screen.GetScrollOffset()).To(Equal(initialOffset + 1))
		})

		It("should scroll up on up arrow", func() {
			// First scroll down
			screen.SetScrollOffset(5)

			msg := tea.KeyMsg{Type: tea.KeyUp}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
			Expect(screen.GetScrollOffset()).To(Equal(4))
		})

		It("should scroll up on k key (vim-style)", func() {
			screen.SetScrollOffset(5)

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
			Expect(screen.GetScrollOffset()).To(Equal(4))
		})

		It("should not scroll up beyond zero", func() {
			screen.SetScrollOffset(0)

			msg := tea.KeyMsg{Type: tea.KeyUp}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
			Expect(screen.GetScrollOffset()).To(Equal(0))
		})

		It("should jump to top on g key", func() {
			screen.SetScrollOffset(10)

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
			_, _ = screen.Update(msg)

			Expect(screen.GetScrollOffset()).To(Equal(0))
		})

		It("should jump to bottom on G key (shift+g)", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}}
			_, _ = screen.Update(msg)

			// Should scroll to a large number (implementation detail)
			Expect(screen.GetScrollOffset()).To(BeNumerically(">", 0))
		})

		It("should preserve scroll position in metadata on cancel", func() {
			screen.SetScrollOffset(5)

			msg := tea.KeyMsg{Type: tea.KeyEsc}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Metadata()["scroll_offset"]).To(Equal(5))
		})

		It("should restore scroll position from metadata", func() {
			screen.RestoreFromMetadata(map[string]interface{}{
				"scroll_offset": 10,
			})

			Expect(screen.GetScrollOffset()).To(Equal(10))
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			renderer := func(data *TestDetailData, _, _ int) string {
				return "Title: " + data.Title + "\n" + data.Description
			}
			screen = base.NewBaseDetailScreen([]string{"Test", "Detail"}, renderer, detailData)
			screen.SetTerminalInfo(120, 40)
		})

		It("should render a view with StandardView structure", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should include rendered content in view", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Title: Test Item"))
			Expect(view).To(ContainSubstring("This is a test description"))
		})

		It("should pass dimensions to renderer", func() {
			var capturedWidth, capturedHeight int
			renderer := func(_ *TestDetailData, width, height int) string {
				capturedWidth = width
				capturedHeight = height
				return "content"
			}

			screen = base.NewBaseDetailScreen([]string{"Test"}, renderer, detailData)
			screen.SetTerminalInfo(100, 50)
			_ = screen.View()

			Expect(capturedWidth).To(Equal(100))
			Expect(capturedHeight).To(Equal(50))
		})

		It("should re-render when data changes", func() {
			view1 := screen.View()

			detailData.Title = "Updated Title"
			view2 := screen.View()

			Expect(view1).NotTo(Equal(view2))
			Expect(view2).To(ContainSubstring("Updated Title"))
		})
	})

	Describe("Data Access", func() {
		BeforeEach(func() {
			renderer := func(_ *TestDetailData, _, _ int) string {
				return "content"
			}
			screen = base.NewBaseDetailScreen([]string{"Test"}, renderer, detailData)
		})

		It("should provide access to data", func() {
			data := screen.GetData()
			Expect(data).To(Equal(detailData))
			Expect(data).To(BeIdenticalTo(detailData)) // Same pointer
		})

		It("should allow modifying data through pointer", func() {
			data := screen.GetData()
			data.Title = "Modified"

			Expect(detailData.Title).To(Equal("Modified"))
			Expect(screen.GetData().Title).To(Equal("Modified"))
		})
	})

	Describe("Screen Interface", func() {
		BeforeEach(func() {
			renderer := func(_ *TestDetailData, _, _ int) string {
				return "content"
			}
			screen = base.NewBaseDetailScreen([]string{"Test"}, renderer, detailData)
		})

		It("should implement Screen interface", func() {
			var _ screens.Screen = screen
		})

		It("should support SetTheme", func() {
			theme := "test-theme"
			screen.SetTheme(theme)
			Expect(screen.Theme()).To(Equal(theme))
		})
	})

	Describe("Edge Cases", func() {
		It("should handle nil data gracefully", func() {
			renderer := func(data *TestDetailData, _, _ int) string {
				if data == nil {
					return "No data"
				}
				return data.Title
			}

			screen := base.NewBaseDetailScreen([]string{"Test"}, renderer, (*TestDetailData)(nil))
			view := screen.View()

			Expect(view).To(ContainSubstring("No data"))
		})

		It("should handle empty breadcrumbs", func() {
			renderer := func(_ *TestDetailData, _, _ int) string {
				return "content"
			}
			screen := base.NewBaseDetailScreen([]string{}, renderer, detailData)

			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle very small terminal dimensions", func() {
			renderer := func(_ *TestDetailData, _, _ int) string {
				return "content"
			}
			screen := base.NewBaseDetailScreen([]string{"Test"}, renderer, detailData)
			screen.SetTerminalInfo(20, 10)

			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle renderer returning empty string", func() {
			renderer := func(_ *TestDetailData, _, _ int) string {
				return ""
			}
			screen := base.NewBaseDetailScreen([]string{"Test"}, renderer, detailData)

			view := screen.View()
			// Should still render StandardView structure
			Expect(view).NotTo(BeEmpty())
		})
	})
})
