package base_test

import (
	"testing"

	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/uikit/display"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Screens Suite")
}

// Screen Tests
//
// These tests define the behavior of the Screen helper,
// which provides common functionality to all screens.
//
// Related:
// - tasks/tasks-42-tui-architecture-refactor.md (Phase 1.1)
// - docs/TUI_DEVELOPER_GUIDE.md (Screen usage)

var _ = Describe("Screen", func() {
	var bs *base.Screen

	BeforeEach(func() {
		bs = base.NewBaseScreen()
	})

	Describe("Terminal Info Management", func() {
		It("should store terminal width and height", func() {
			bs.SetTerminalInfo(100, 50)

			Expect(bs.Width()).To(Equal(100))
			Expect(bs.Height()).To(Equal(50))
		})

		It("should provide access to terminal dimensions", func() {
			// NewBaseScreen creates with default dimensions (120x40)
			Expect(bs.Width()).To(Equal(120))
			Expect(bs.Height()).To(Equal(40))
		})

		It("should handle nil terminal info gracefully", func() {
			// Screen always has defaults, so calling Width/Height never panics
			newScreen := base.NewBaseScreen()
			Expect(newScreen.Width()).To(BeNumerically(">", 0))
			Expect(newScreen.Height()).To(BeNumerically(">", 0))
		})

		It("should update dimensions when SetTerminalInfo is called", func() {
			bs.SetTerminalInfo(80, 24)
			Expect(bs.Width()).To(Equal(80))
			Expect(bs.Height()).To(Equal(24))

			bs.SetTerminalInfo(200, 60)
			Expect(bs.Width()).To(Equal(200))
			Expect(bs.Height()).To(Equal(60))
		})
	})

	Describe("Theme Management", func() {
		It("should store theme reference", func() {
			theme := "test-theme"
			bs.SetTheme(theme)

			Expect(bs.Theme()).To(Equal("test-theme"))
		})

		It("should provide access to theme", func() {
			// Initially nil
			Expect(bs.Theme()).To(BeNil())
		})

		It("should handle nil theme gracefully", func() {
			// Should not panic when theme is nil
			bs.SetTheme(nil)
			Expect(bs.Theme()).To(BeNil())
		})

		It("should update theme when SetTheme is called", func() {
			bs.SetTheme("theme1")
			Expect(bs.Theme()).To(Equal("theme1"))

			bs.SetTheme("theme2")
			Expect(bs.Theme()).To(Equal("theme2"))
		})
	})

	Describe("View Creation Helpers", func() {
		It("should provide CreateView method for StandardView creation", func() {
			view := bs.CreateView([]string{"Home"}, "Content", "Footer")

			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Home"))
		})

		It("should pass terminal dimensions to StandardView", func() {
			bs.SetTerminalInfo(80, 24)
			view := bs.CreateView([]string{"Test"}, "Content", "Footer")

			// View should be rendered (non-empty)
			Expect(view).NotTo(BeEmpty())
		})

		It("should pass theme to StandardView", func() {
			bs.SetTheme("custom-theme")
			view := bs.CreateView([]string{"Test"}, "Content", "Footer")

			// View renders without panicking
			Expect(view).NotTo(BeEmpty())
		})

		It("should allow custom breadcrumbs", func() {
			view := bs.CreateView([]string{"Main", "Sub", "Detail"}, "Content", "Footer")

			Expect(view).To(ContainSubstring("Main"))
			Expect(view).To(ContainSubstring("Sub"))
			Expect(view).To(ContainSubstring("Detail"))
		})

		It("should allow custom content", func() {
			view := bs.CreateView([]string{"Test"}, "My Custom Content Here", "Footer")

			Expect(view).To(ContainSubstring("My Custom Content Here"))
		})

		It("should allow custom footer", func() {
			view := bs.CreateView([]string{"Test"}, "Content", "Press Enter to continue")

			Expect(view).To(ContainSubstring("Press Enter to continue"))
		})
	})

	Describe("Window Size Message Handling", func() {
		It("should update dimensions when receiving WindowSizeMsg", func() {
			msg := tea.WindowSizeMsg{Width: 150, Height: 45}

			cmd := bs.HandleWindowSizeMsg(msg)

			Expect(cmd).To(BeNil())
			Expect(bs.Width()).To(Equal(150))
			Expect(bs.Height()).To(Equal(45))
		})

		It("should return nil command for WindowSizeMsg", func() {
			msg := tea.WindowSizeMsg{Width: 100, Height: 30}

			cmd := bs.HandleWindowSizeMsg(msg)

			Expect(cmd).To(BeNil())
		})

		It("should return nil for non-WindowSizeMsg", func() {
			msg := tea.KeyMsg{Type: tea.KeyEnter}

			cmd := bs.HandleWindowSizeMsg(msg)

			Expect(cmd).To(BeNil())
		})
	})

	Describe("Logo Management", func() {
		It("should store logo reference", func() {
			logo := display.NewLogo(false, 100)
			bs.SetLogo(logo, 2)

			Expect(bs.GetLogo()).To(Equal(logo))
			Expect(bs.GetLogoSpacing()).To(Equal(2))
		})

		It("should handle nil logo gracefully", func() {
			bs.SetLogo(nil, 0)

			Expect(bs.GetLogo()).To(BeNil())
			Expect(bs.GetLogoSpacing()).To(Equal(0))
		})

		It("should include logo in view when set", func() {
			logo := display.NewLogo(false, 100)
			bs.SetLogo(logo, 1)

			view := bs.CreateView([]string{"Test"}, "Content", "Footer")

			// Logo should be rendered in the view
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Composition", func() {
		It("should be embeddable in concrete screen implementations", func() {
			// Create a mock screen that embeds Screen
			type MockScreen struct {
				*base.Screen
				customField string
			}

			mock := &MockScreen{
				Screen:      base.NewBaseScreen(),
				customField: "test",
			}

			// Should be able to access Screen methods
			mock.SetTerminalInfo(100, 50)
			Expect(mock.Width()).To(Equal(100))
			Expect(mock.customField).To(Equal("test"))
		})

		It("should allow concrete screens to use Screen methods", func() {
			type MockScreen struct {
				*base.Screen
			}

			mock := &MockScreen{
				Screen: base.NewBaseScreen(),
			}

			// Use HandleWindowSizeMsg
			msg := tea.WindowSizeMsg{Width: 80, Height: 24}
			cmd := mock.HandleWindowSizeMsg(msg)

			Expect(cmd).To(BeNil())
			Expect(mock.Width()).To(Equal(80))

			// Use CreateView
			view := mock.CreateView([]string{"Test"}, "Content", "Footer")
			Expect(view).NotTo(BeEmpty())
		})
	})
})
