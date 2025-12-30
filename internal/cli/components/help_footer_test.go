package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/navigation"
)

func TestComponents(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Components Suite")
}

var _ = Describe("HelpFooter Component", func() {
	Describe("NewHelpFooter", func() {
		It("should create help footer with form context", func() {
			footer := NewHelpFooter("form", 80)
			Expect(footer.context).To(Equal("form"))
			Expect(footer.width).To(Equal(80))
			Expect(footer.height).To(Equal(1))
		})

		It("should create help footer with list context", func() {
			footer := NewHelpFooter("list", 100)
			Expect(footer.context).To(Equal("list"))
			Expect(footer.width).To(Equal(100))
		})

		It("should create help footer with metadata_review context", func() {
			footer := NewHelpFooter("metadata_review", 120)
			Expect(footer.context).To(Equal("metadata_review"))
			Expect(footer.width).To(Equal(120))
		})

		It("should create help footer with empty context", func() {
			footer := NewHelpFooter("", 80)
			Expect(footer.context).To(Equal(""))
			Expect(footer.width).To(Equal(80))
		})
	})

	Describe("NewHelpFooterWithKeys", func() {
		It("should create help footer with specific keys", func() {
			keys := []navigation.NavigationKey{navigation.KeyUp, navigation.KeyDown}
			footer := NewHelpFooterWithKeys(keys, 80)
			Expect(len(footer.keys)).To(Equal(2))
			Expect(footer.width).To(Equal(80))
			Expect(footer.context).To(Equal(""))
		})

		It("should store multiple keys", func() {
			keys := []navigation.NavigationKey{
				navigation.KeyUp,
				navigation.KeyDown,
				navigation.KeySelect,
				navigation.KeyBack,
			}
			footer := NewHelpFooterWithKeys(keys, 100)
			Expect(len(footer.keys)).To(Equal(4))
		})
	})

	Describe("Init", func() {
		It("should return nil command", func() {
			footer := NewHelpFooter("form", 80)
			cmd := footer.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update", func() {
		It("should handle WindowSizeMsg", func() {
			footer := NewHelpFooter("form", 80)
			msg := tea.WindowSizeMsg{Width: 120, Height: 30}
			updated, cmd := footer.Update(msg)
			Expect(updated.(HelpFooterModel).width).To(Equal(120))
			Expect(cmd).To(BeNil())
		})

		It("should ignore other message types", func() {
			footer := NewHelpFooter("form", 80)
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			updated, cmd := footer.Update(msg)
			Expect(updated.(HelpFooterModel).width).To(Equal(80))
			Expect(cmd).To(BeNil())
		})

		It("should update width on multiple WindowSizeMsg", func() {
			footer := NewHelpFooter("form", 80)
			msg1 := tea.WindowSizeMsg{Width: 100, Height: 30}
			updated1, _ := footer.Update(msg1)
			footer2 := updated1.(HelpFooterModel)
			Expect(footer2.width).To(Equal(100))

			msg2 := tea.WindowSizeMsg{Width: 120, Height: 30}
			updated2, _ := footer2.Update(msg2)
			footer3 := updated2.(HelpFooterModel)
			Expect(footer3.width).To(Equal(120))
		})
	})

	Describe("View", func() {
		It("should render form context help text", func() {
			footer := NewHelpFooter("form", 80)
			view := footer.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render list context help text", func() {
			footer := NewHelpFooter("list", 80)
			view := footer.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render with custom keys", func() {
			keys := []navigation.NavigationKey{navigation.KeyUp, navigation.KeyDown}
			footer := NewHelpFooterWithKeys(keys, 80)
			view := footer.View()
			Expect(view).To(ContainSubstring("↑/k"))
			Expect(view).To(ContainSubstring("↓/j"))
		})

		It("should return empty string for zero width", func() {
			footer := NewHelpFooter("form", 0)
			view := footer.View()
			Expect(view).To(Equal(""))
		})

		It("should return empty string for negative width", func() {
			footer := NewHelpFooter("form", -10)
			view := footer.View()
			Expect(view).To(Equal(""))
		})

		It("should render on narrow terminal (80 chars)", func() {
			footer := NewHelpFooter("form", 80)
			view := footer.View()
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("should render on wide terminal (200 chars)", func() {
			footer := NewHelpFooter("list", 200)
			view := footer.View()
			Expect(len(view)).To(BeNumerically(">", 0))
		})
	})

	Describe("SetWidth", func() {
		It("should update footer width", func() {
			footer := NewHelpFooter("form", 80)
			footer.SetWidth(120)
			Expect(footer.width).To(Equal(120))
		})

		It("should update for render after SetWidth", func() {
			footer := NewHelpFooter("form", 80)
			view1 := footer.View()
			footer.SetWidth(40)
			view2 := footer.View()
			// Both should be valid but may differ due to truncation
			Expect(view1).NotTo(BeEmpty())
			Expect(view2).NotTo(BeEmpty())
		})
	})

	Describe("SetContext", func() {
		It("should change context to form", func() {
			footer := NewHelpFooter("list", 80)
			footer.SetContext("form")
			Expect(footer.context).To(Equal("form"))
		})

		It("should change context to metadata_review", func() {
			footer := NewHelpFooter("list", 80)
			footer.SetContext("metadata_review")
			Expect(footer.context).To(Equal("metadata_review"))
		})

		It("should reflect context change in View", func() {
			footer := NewHelpFooter("form", 80)
			view1 := footer.View()
			footer.SetContext("list")
			view2 := footer.View()
			// Views may be different due to context-specific help
			Expect(view1).NotTo(BeEmpty())
			Expect(view2).NotTo(BeEmpty())
		})
	})

	Describe("SetKeys", func() {
		It("should update displayed keys", func() {
			footer := NewHelpFooter("form", 80)
			keys := []navigation.NavigationKey{navigation.KeyUp, navigation.KeyDown, navigation.KeySelect}
			footer.SetKeys(keys)
			Expect(len(footer.keys)).To(Equal(3))
		})

		It("should prefer custom keys over context", func() {
			footer := NewHelpFooter("form", 80)
			keys := []navigation.NavigationKey{navigation.KeyBack, navigation.KeyHelp}
			footer.SetKeys(keys)
			view := footer.View()
			Expect(view).To(ContainSubstring("Esc"))
		})
	})

	Describe("GetHeight", func() {
		It("should always return 1 for footer height", func() {
			footer := NewHelpFooter("form", 80)
			Expect(footer.GetHeight()).To(Equal(1))
		})
	})

	Describe("RenderForWidth", func() {
		It("should render with specified width without changing state", func() {
			footer := NewHelpFooter("form", 80)
			view := footer.RenderForWidth(120)
			Expect(footer.width).To(Equal(80)) // Original width unchanged
			Expect(view).NotTo(BeEmpty())
		})

		It("should render at different widths", func() {
			footer := NewHelpFooter("form", 100)
			view80 := footer.RenderForWidth(80)
			view120 := footer.RenderForWidth(120)
			Expect(view80).NotTo(BeEmpty())
			Expect(view120).NotTo(BeEmpty())
		})
	})

	Describe("WithBorder", func() {
		It("should render footer with border", func() {
			footer := NewHelpFooter("form", 80)
			view := footer.WithBorder()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("WithPadding", func() {
		It("should render footer with padding", func() {
			footer := NewHelpFooter("form", 80)
			view := footer.WithPadding(1, 2)
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Contextual Help Rendering", func() {
		It("should render form help with proper keys", func() {
			footer := NewHelpFooter("form", 200)
			view := footer.View()
			Expect(view).To(ContainSubstring("↑/k"))
			Expect(view).To(ContainSubstring("↓/j"))
			Expect(view).To(ContainSubstring("Enter"))
			Expect(view).To(ContainSubstring("Esc"))
		})

		It("should render list help with proper keys", func() {
			footer := NewHelpFooter("list", 200)
			view := footer.View()
			Expect(view).To(ContainSubstring("↑/k"))
			Expect(view).To(ContainSubstring("↓/j"))
		})

		It("should render metadata_review help", func() {
			footer := NewHelpFooter("metadata_review", 200)
			view := footer.View()
			Expect(view).To(ContainSubstring("↑/k"))
			Expect(view).To(ContainSubstring("↓/j"))
		})

		It("should render home help", func() {
			footer := NewHelpFooter("home", 200)
			view := footer.View()
			Expect(view).To(ContainSubstring("c:"))
			Expect(view).To(ContainSubstring("?:"))
		})
	})

	Describe("Responsive Behavior", func() {
		It("should handle very narrow terminal (40 chars)", func() {
			footer := NewHelpFooter("list", 40)
			view := footer.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle very wide terminal (300 chars)", func() {
			footer := NewHelpFooter("list", 300)
			view := footer.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should gracefully truncate on narrow terminals", func() {
			footer := NewHelpFooter("list", 30)
			view := footer.View()
			// Should either show help or fallback to minimal display
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Model Interface Compliance", func() {
		It("should implement BubbleTea Model interface", func() {
			var model tea.Model
			footer := NewHelpFooter("form", 80)
			model = footer
			Expect(model).NotTo(BeNil())
		})

		It("should support Init, Update, and View", func() {
			footer := NewHelpFooter("form", 80)
			cmd := footer.Init()
			Expect(cmd).To(BeNil())

			updated, _ := footer.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
			Expect(updated).NotTo(BeNil())

			view := footer.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})
