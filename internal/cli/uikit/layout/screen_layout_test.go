package layout_test

import (
	"regexp"
	"strings"

	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/layout"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// stripAnsi removes ANSI escape codes from a string for reliable test comparisons.
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripAnsi(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// MockLogo implements LogoRenderer for testing.
type MockLogo struct {
	width   int
	content string
}

func NewMockLogo(content string) *MockLogo {
	return &MockLogo{content: content}
}

func (m *MockLogo) ViewStatic() string {
	return m.content
}

func (m *MockLogo) SetWidth(width int) {
	m.width = width
}

var _ = Describe("ScreenLayout Pinned Layout", func() {
	var (
		termInfo *terminal.Info
		theme    themes.Theme
	)

	BeforeEach(func() {
		termInfo = &terminal.Info{
			Width:   80,
			Height:  24,
			IsValid: true,
		}
		theme = themes.NewDefaultTheme()
	})

	Describe("Logo at Top (Top Pinned)", func() {
		It("should render logo at line 0 when LogoSpacing is 0", func() {
			logo := NewMockLogo("LOGO")
			view := layout.NewScreenLayout(termInfo).
				WithLogo(logo, 0).
				WithTheme(theme).
				WithContent("Content").
				WithHelp("Help text")

			rendered := view.Render()
			lines := strings.Split(stripAnsi(rendered), "\n")

			Expect(strings.TrimSpace(lines[0])).To(ContainSubstring("LOGO"))
		})

		It("should respect LogoSpacing parameter for blank lines before logo", func() {
			logo := NewMockLogo("LOGO")
			view := layout.NewScreenLayout(termInfo).
				WithLogo(logo, 5).
				WithTheme(theme).
				WithContent("Content")

			rendered := view.Render()
			lines := strings.Split(stripAnsi(rendered), "\n")

			for i := range 5 {
				Expect(strings.TrimSpace(lines[i])).To(BeEmpty(), "line %d should be blank (spacing=5)", i)
			}
			Expect(strings.TrimSpace(lines[5])).To(ContainSubstring("LOGO"))
		})
	})

	Describe("Footer at Bottom (Bottom Pinned)", func() {
		It("should render footer on the last lines of the terminal", func() {
			view := layout.NewScreenLayout(termInfo).
				WithTheme(theme).
				WithContent("Content").
				WithHelp("Press q to quit")

			rendered := view.Render()
			lines := strings.Split(stripAnsi(rendered), "\n")

			// Terminal height is 24, so footer should be near line 23 (0-indexed)
			// Find the last non-empty line
			var lastNonEmptyLine string
			for i := len(lines) - 1; i >= 0; i-- {
				if strings.TrimSpace(lines[i]) != "" {
					lastNonEmptyLine = lines[i]
					break
				}
			}

			Expect(lastNonEmptyLine).To(ContainSubstring("Press q to quit"))
		})

		It("should render footer at bottom even with logo", func() {
			logo := NewMockLogo("LOGO\nLINE 2")
			view := layout.NewScreenLayout(termInfo).
				WithLogo(logo, 0).
				WithTheme(theme).
				WithContent("Content").
				WithHelp("Help text")

			rendered := view.Render()
			lines := strings.Split(stripAnsi(rendered), "\n")

			// Find last non-empty line
			var lastNonEmptyLine string
			for i := len(lines) - 1; i >= 0; i-- {
				if strings.TrimSpace(lines[i]) != "" {
					lastNonEmptyLine = lines[i]
					break
				}
			}

			Expect(lastNonEmptyLine).To(ContainSubstring("Help text"))
		})
	})

	Describe("Content Placement", func() {
		It("should render content immediately after header section", func() {
			logo := NewMockLogo("LOGO")
			view := layout.NewScreenLayout(termInfo).
				WithLogo(logo, 0).
				WithBreadcrumbs("Home", "Settings").
				WithTheme(theme).
				WithContent("Main Content").
				WithHelp("Help")

			rendered := view.Render()
			lines := strings.Split(stripAnsi(rendered), "\n")

			// Find logo, breadcrumbs, then content should follow
			logoFound := false
			breadcrumbsFound := false
			contentIndex := -1

			for i, line := range lines {
				stripped := strings.TrimSpace(line)
				if !logoFound && strings.Contains(stripped, "LOGO") {
					logoFound = true
				} else if logoFound && !breadcrumbsFound && strings.Contains(stripped, "Settings") {
					breadcrumbsFound = true
				} else if breadcrumbsFound && strings.Contains(stripped, "Main Content") {
					contentIndex = i
					break
				}
			}

			Expect(logoFound).To(BeTrue(), "Logo should be found")
			Expect(breadcrumbsFound).To(BeTrue(), "Breadcrumbs should be found")
			Expect(contentIndex).To(BeNumerically(">", 0), "Content should be found after header")
		})
	})

	Describe("Spacer Fills Gap", func() {
		It("should fill vertical space between content and footer", func() {
			view := layout.NewScreenLayout(termInfo).
				WithTheme(theme).
				WithContent("Content").
				WithHelp("Help")

			rendered := view.Render()
			lines := strings.Split(rendered, "\n")

			// Total lines should equal terminal height
			Expect(lines).To(HaveLen(termInfo.Height))
		})

		It("should calculate spacer correctly with logo and footer", func() {
			logo := NewMockLogo("LOGO\nLINE2\nLINE3") // 3 lines
			view := layout.NewScreenLayout(termInfo).
				WithLogo(logo, 0).
				WithTheme(theme).
				WithContent("Content"). // 1 line
				WithHelp("Help text")   // ~2 lines (separator + help)

			rendered := view.Render()
			lines := strings.Split(rendered, "\n")

			// Should have exactly termInfo.Height lines
			Expect(lines).To(HaveLen(termInfo.Height))
		})
	})

	Describe("Overflow Handling (Graceful Degradation)", func() {
		It("should handle content taller than terminal without negative spacer", func() {
			// Create content that exceeds terminal height
			tallContent := strings.Repeat("Line\n", 30) // 30 lines in a 24-line terminal

			view := layout.NewScreenLayout(termInfo).
				WithTheme(theme).
				WithContent(tallContent).
				WithHelp("Help")

			// Should not panic
			Expect(func() {
				_ = view.Render()
			}).NotTo(Panic())
		})

		It("should constrain overflowing content to available height via viewport", func() {
			tallContent := strings.Repeat("Line\n", 30) // 30 lines in a 24-line terminal

			view := layout.NewScreenLayout(termInfo).
				WithTheme(theme).
				WithContent(tallContent).
				WithHelp("Help")

			rendered := view.Render()
			lines := strings.Split(rendered, "\n")

			// Total output must still be exactly terminal height
			Expect(lines).To(HaveLen(termInfo.Height),
				"output should be exactly terminal height even with overflowing content")

			// Footer must still be at the bottom
			strippedLines := strings.Split(stripAnsi(rendered), "\n")
			var lastNonEmptyLine string
			for i := len(strippedLines) - 1; i >= 0; i-- {
				if strings.TrimSpace(strippedLines[i]) != "" {
					lastNonEmptyLine = strippedLines[i]
					break
				}
			}
			Expect(lastNonEmptyLine).To(ContainSubstring("Help"),
				"footer should remain pinned at bottom even when content overflows")
		})
	})

	Describe("No Logo Mode", func() {
		It("should start with breadcrumbs at line 0 when no logo", func() {
			view := layout.NewScreenLayout(termInfo).
				WithBreadcrumbs("Home", "Settings").
				WithTheme(theme).
				WithContent("Content").
				WithHelp("Help")

			rendered := view.Render()
			lines := strings.Split(stripAnsi(rendered), "\n")

			// First non-empty line should contain breadcrumbs
			var firstNonEmpty string
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					firstNonEmpty = line
					break
				}
			}

			Expect(firstNonEmpty).To(ContainSubstring("Settings"))
		})
	})

	Describe("No Footer Mode", func() {
		It("should render without footer when help text not provided", func() {
			view := layout.NewScreenLayout(termInfo).
				WithTheme(theme).
				WithContent("Content")
			// No WithHelp() call

			rendered := view.Render()

			// Should still render successfully
			Expect(rendered).NotTo(BeEmpty())
		})
	})

	Describe("Vertical Alignment", func() {
		It("should use Top alignment instead of Center", func() {
			logo := NewMockLogo("LOGO")
			view := layout.NewScreenLayout(termInfo).
				WithLogo(logo, 0).
				WithTheme(theme).
				WithContent("Content").
				WithHelp("Help")

			rendered := view.Render()
			lines := strings.Split(stripAnsi(rendered), "\n")

			// Logo should be in the top portion (first 5 lines)
			logoInTop := false
			for i := 0; i < 5 && i < len(lines); i++ {
				if strings.Contains(lines[i], "LOGO") {
					logoInTop = true
					break
				}
			}

			Expect(logoInTop).To(BeTrue(), "Logo should be in top portion of screen")

			// Help should be in the bottom portion (last 5 lines)
			helpInBottom := false
			startIdx := len(lines) - 5
			if startIdx < 0 {
				startIdx = 0
			}
			for i := startIdx; i < len(lines); i++ {
				if strings.Contains(lines[i], "Help") {
					helpInBottom = true
					break
				}
			}

			Expect(helpInBottom).To(BeTrue(), "Help should be in bottom portion of screen")
		})
	})

	Describe("Horizontal Centering Preserved", func() {
		It("should still center content horizontally", func() {
			view := layout.NewScreenLayout(termInfo).
				WithTheme(theme).
				WithContent("X"). // Single character
				WithHelp("Help")

			rendered := view.Render()
			lines := strings.Split(stripAnsi(rendered), "\n")

			// Find the content line
			for _, line := range lines {
				if strings.Contains(line, "X") {
					// Content should not be at the very start of the line
					trimmed := strings.TrimLeft(line, " ")
					leadingSpaces := len(line) - len(trimmed)
					Expect(leadingSpaces).To(BeNumerically(">", 0), "Content should have leading spaces (centered)")
					break
				}
			}
		})
	})

	Describe("Total Height Matches Terminal", func() {
		It("should render exactly terminal height lines", func() {
			view := layout.NewScreenLayout(termInfo).
				WithTheme(theme).
				WithContent("Content").
				WithHelp("Help")

			rendered := view.Render()
			lines := strings.Split(rendered, "\n")

			Expect(lines).To(HaveLen(termInfo.Height))
		})

		It("should maintain terminal height with different sizes", func() {
			smallTerm := &terminal.Info{Width: 40, Height: 10, IsValid: true}
			view := layout.NewScreenLayout(smallTerm).
				WithTheme(theme).
				WithContent("Content").
				WithHelp("Help")

			rendered := view.Render()
			lines := strings.Split(rendered, "\n")

			Expect(lines).To(HaveLen(10))
		})
	})

	Describe("Available Content Height Calculation", func() {
		It("should calculate available content height correctly", func() {
			logo := NewMockLogo("LOGO\nLINE2\nLINE3") // 3 lines
			view := layout.NewScreenLayout(termInfo).
				WithLogo(logo, 0).
				WithBreadcrumbs("Home", "Settings").
				WithTheme(theme).
				WithHelp("Help text")

			// Expected calculation:
			// Terminal height: 24
			// Header: 2 (blank) + 3 (logo) + 1 (blank) + 1 (breadcrumbs) + 1 (blank) = 8 lines
			// Footer: 1 (blank) + 1 (help) = 2 lines
			// Available content: 24 - 8 - 2 = 14 lines
			availableHeight := view.GetAvailableContentHeight()

			Expect(availableHeight).To(BeNumerically(">=", 10), "Should have at least 10 lines for content")
			Expect(availableHeight).To(BeNumerically("<=", 18), "Should not exceed reasonable content height")
		})

		It("should return full height when no header or footer", func() {
			view := layout.NewScreenLayout(termInfo).
				WithTheme(theme)
			// No logo, no help = minimal header/footer

			availableHeight := view.GetAvailableContentHeight()

			// Should be close to terminal height (24) minus minimal margins
			Expect(availableHeight).To(BeNumerically(">", 20), "Should use most of terminal for content")
		})

		It("should handle small terminals gracefully", func() {
			smallTerm := &terminal.Info{Width: 40, Height: 10, IsValid: true}
			logo := NewMockLogo("LOGO")
			view := layout.NewScreenLayout(smallTerm).
				WithLogo(logo, 0).
				WithTheme(theme).
				WithHelp("Help")

			availableHeight := view.GetAvailableContentHeight()

			// Even in small terminal, should return positive height
			Expect(availableHeight).To(BeNumerically(">", 0), "Should always return positive content height")
			Expect(availableHeight).To(BeNumerically("<", 10), "Should be less than terminal height")
		})
	})
})

var _ = Describe("Header", func() {
	var (
		header *layout.Header
		theme  themes.Theme
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
		header = layout.NewHeader("Test Title", 80)
	})

	Describe("NewHeader", func() {
		It("should create a header with the given title", func() {
			Expect(header.GetTitle()).To(Equal("Test Title"))
		})

		It("should have an empty subtitle by default", func() {
			Expect(header.GetSubtitle()).To(BeEmpty())
		})

		It("should have no breadcrumbs by default", func() {
			Expect(header.GetBreadcrumbs()).To(BeNil())
		})
	})

	Describe("WithTheme", func() {
		It("should return the header for chaining", func() {
			result := header.WithTheme(theme)
			Expect(result).To(BeIdenticalTo(header))
		})
	})

	Describe("WithSubtitle", func() {
		It("should set the subtitle", func() {
			header.WithSubtitle("My Subtitle")
			Expect(header.GetSubtitle()).To(Equal("My Subtitle"))
		})

		It("should return the header for chaining", func() {
			result := header.WithSubtitle("Sub")
			Expect(result).To(BeIdenticalTo(header))
		})
	})

	Describe("WithBreadcrumbs", func() {
		It("should set the breadcrumbs", func() {
			header.WithBreadcrumbs([]string{"Home", "Settings"})
			Expect(header.GetBreadcrumbs()).To(Equal([]string{"Home", "Settings"}))
		})

		It("should return the header for chaining", func() {
			result := header.WithBreadcrumbs([]string{"Home"})
			Expect(result).To(BeIdenticalTo(header))
		})
	})

	Describe("WithBorder", func() {
		It("should return the header for chaining", func() {
			result := header.WithBorder()
			Expect(result).To(BeIdenticalTo(header))
		})
	})

	Describe("SetWidth", func() {
		It("should update the width used for rendering", func() {
			header.SetWidth(120)
			header.WithTheme(theme)
			view := stripAnsi(header.View())
			Expect(view).To(ContainSubstring("Test Title"))
		})
	})

	Describe("SetHeight", func() {
		It("should set the height without error", func() {
			header.SetHeight(3)
			header.WithTheme(theme)
			view := stripAnsi(header.View())
			Expect(view).To(ContainSubstring("Test Title"))
		})
	})

	Describe("GetTitle", func() {
		It("should return the title set at construction", func() {
			Expect(header.GetTitle()).To(Equal("Test Title"))
		})
	})

	Describe("GetSubtitle", func() {
		It("should return the subtitle when set", func() {
			header.WithSubtitle("Description")
			Expect(header.GetSubtitle()).To(Equal("Description"))
		})

		It("should return empty string when no subtitle is set", func() {
			Expect(header.GetSubtitle()).To(BeEmpty())
		})
	})

	Describe("GetBreadcrumbs", func() {
		It("should return the breadcrumbs when set", func() {
			header.WithBreadcrumbs([]string{"A", "B", "C"})
			Expect(header.GetBreadcrumbs()).To(Equal([]string{"A", "B", "C"}))
		})
	})

	Describe("AddBreadcrumb", func() {
		It("should append a breadcrumb to the existing list", func() {
			header.WithBreadcrumbs([]string{"Home"})
			header.AddBreadcrumb("Settings")
			Expect(header.GetBreadcrumbs()).To(Equal([]string{"Home", "Settings"}))
		})

		It("should add a breadcrumb when list is initially nil", func() {
			header.AddBreadcrumb("Home")
			Expect(header.GetBreadcrumbs()).To(ContainElement("Home"))
		})
	})

	Describe("ClearBreadcrumbs", func() {
		It("should remove all breadcrumbs", func() {
			header.WithBreadcrumbs([]string{"Home", "Settings"})
			header.ClearBreadcrumbs()
			Expect(header.GetBreadcrumbs()).To(BeEmpty())
		})
	})

	Describe("View", func() {
		Context("when width is valid", func() {
			It("should contain the title", func() {
				header.WithTheme(theme)
				view := stripAnsi(header.View())
				Expect(view).To(ContainSubstring("Test Title"))
			})

			It("should render subtitle when set", func() {
				header.WithTheme(theme).WithSubtitle("My Subtitle")
				view := stripAnsi(header.View())
				Expect(view).To(ContainSubstring("My Subtitle"))
			})

			It("should render breadcrumbs when set", func() {
				header.WithTheme(theme).WithBreadcrumbs([]string{"Home", "Settings"})
				view := stripAnsi(header.View())
				Expect(view).To(ContainSubstring("Home"))
				Expect(view).To(ContainSubstring("Settings"))
			})

			It("should render with border when enabled", func() {
				header.WithTheme(theme).WithBorder()
				view := header.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should use default theme when no theme is set", func() {
				h := layout.NewHeader("No Theme", 80)
				view := stripAnsi(h.View())
				Expect(view).To(ContainSubstring("No Theme"))
			})
		})

		Context("when width is zero or negative", func() {
			It("should return empty string for zero width", func() {
				h := layout.NewHeader("Title", 0)
				Expect(h.View()).To(BeEmpty())
			})

			It("should return empty string for negative width", func() {
				h := layout.NewHeader("Title", -1)
				Expect(h.View()).To(BeEmpty())
			})
		})

		Context("when title exceeds available width", func() {
			It("should truncate the title with ellipsis", func() {
				longTitle := strings.Repeat("A", 200)
				h := layout.NewHeader(longTitle, 50).WithTheme(theme)
				view := stripAnsi(h.View())
				Expect(view).To(ContainSubstring("..."))
			})
		})

		Context("when subtitle exceeds available width", func() {
			It("should truncate the subtitle with ellipsis", func() {
				h := layout.NewHeader("Title", 50).
					WithTheme(theme).
					WithSubtitle(strings.Repeat("B", 200))
				view := stripAnsi(h.View())
				Expect(view).To(ContainSubstring("..."))
			})
		})
	})

	Describe("GetClickedBreadcrumbIndex", func() {
		It("should return -1 when no breadcrumbs exist", func() {
			Expect(header.GetClickedBreadcrumbIndex(5, 0)).To(Equal(-1))
		})

		Context("with breadcrumbs set", func() {
			BeforeEach(func() {
				header.WithBreadcrumbs([]string{"Home", "Settings", "Profile"})
			})

			It("should return 0 when clicking on the first breadcrumb", func() {
				Expect(header.GetClickedBreadcrumbIndex(0, 0)).To(Equal(0))
			})

			It("should return 1 when clicking on the second breadcrumb", func() {
				Expect(header.GetClickedBreadcrumbIndex(7, 0)).To(Equal(1))
			})

			It("should return 2 when clicking on the third breadcrumb", func() {
				Expect(header.GetClickedBreadcrumbIndex(19, 0)).To(Equal(2))
			})

			It("should return -1 when clicking outside all breadcrumb bounds", func() {
				Expect(header.GetClickedBreadcrumbIndex(100, 0)).To(Equal(-1))
			})
		})
	})
})

var _ = Describe("Footer", func() {
	var (
		footer *layout.Footer
		theme  themes.Theme
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
		footer = layout.NewFooter(80)
	})

	Describe("NewFooter", func() {
		It("should create a footer with the given width", func() {
			Expect(footer).NotTo(BeNil())
		})

		It("should have an empty status message by default", func() {
			Expect(footer.GetStatusMessage()).To(BeEmpty())
		})

		It("should have an empty mode context by default", func() {
			Expect(footer.GetModeContext()).To(BeEmpty())
		})

		It("should render empty view when no content is set", func() {
			view := footer.View()
			Expect(stripAnsi(view)).To(BeEmpty())
		})
	})

	Describe("WithTheme", func() {
		It("should return the footer for chaining", func() {
			result := footer.WithTheme(theme)
			Expect(result).To(BeIdenticalTo(footer))
		})
	})

	Describe("WithStatus", func() {
		It("should set the status message", func() {
			footer.WithStatus("3/10 events")
			Expect(footer.GetStatusMessage()).To(Equal("3/10 events"))
		})

		It("should return the footer for chaining", func() {
			result := footer.WithStatus("status")
			Expect(result).To(BeIdenticalTo(footer))
		})
	})

	Describe("WithMode", func() {
		It("should set the mode context", func() {
			footer.WithMode("Capture Mode: Timeline")
			Expect(footer.GetModeContext()).To(Equal("Capture Mode: Timeline"))
		})

		It("should return the footer for chaining", func() {
			result := footer.WithMode("mode")
			Expect(result).To(BeIdenticalTo(footer))
		})
	})

	Describe("WithHelp", func() {
		It("should set the help text", func() {
			footer.WithTheme(theme).WithHelp("Press ? for help")
			view := stripAnsi(footer.View())
			Expect(view).To(ContainSubstring("Press ? for help"))
		})

		It("should return the footer for chaining", func() {
			result := footer.WithHelp("help")
			Expect(result).To(BeIdenticalTo(footer))
		})
	})

	Describe("SetWidth", func() {
		It("should return the footer for chaining", func() {
			result := footer.SetWidth(120)
			Expect(result).To(BeIdenticalTo(footer))
		})

		It("should affect rendering width", func() {
			footer.WithTheme(theme).WithStatus("status")
			footer.SetWidth(120)
			view := stripAnsi(footer.View())
			Expect(view).To(ContainSubstring("status"))
		})
	})

	Describe("SetHeight", func() {
		It("should set the height without error", func() {
			footer.SetHeight(3)
			footer.WithTheme(theme).WithStatus("status")
			view := stripAnsi(footer.View())
			Expect(view).To(ContainSubstring("status"))
		})
	})

	Describe("SetShowStatus", func() {
		It("should hide status when set to false", func() {
			footer.WithTheme(theme).WithStatus("visible status")
			footer.SetShowStatus(false)
			view := stripAnsi(footer.View())
			Expect(view).NotTo(ContainSubstring("visible status"))
		})

		It("should show status when set to true", func() {
			footer.WithTheme(theme).WithStatus("visible status")
			footer.SetShowStatus(true)
			view := stripAnsi(footer.View())
			Expect(view).To(ContainSubstring("visible status"))
		})
	})

	Describe("SetShowMode", func() {
		It("should hide mode when set to false", func() {
			footer.WithTheme(theme).WithMode("Edit Mode")
			footer.SetShowMode(false)
			view := stripAnsi(footer.View())
			Expect(view).NotTo(ContainSubstring("Edit Mode"))
		})

		It("should show mode when set to true", func() {
			footer.WithTheme(theme).WithMode("Edit Mode")
			footer.SetShowMode(true)
			view := stripAnsi(footer.View())
			Expect(view).To(ContainSubstring("Edit Mode"))
		})
	})

	Describe("SetShowHelp", func() {
		It("should hide help when set to false", func() {
			footer.WithTheme(theme).WithHelp("Press ? for help")
			footer.SetShowHelp(false)
			view := stripAnsi(footer.View())
			Expect(view).NotTo(ContainSubstring("Press ? for help"))
		})

		It("should show help when set to true", func() {
			footer.WithTheme(theme).WithHelp("Press ? for help")
			footer.SetShowHelp(true)
			view := stripAnsi(footer.View())
			Expect(view).To(ContainSubstring("Press ? for help"))
		})
	})

	Describe("GetStatusMessage", func() {
		It("should return the status message when set", func() {
			footer.WithStatus("5/20 items")
			Expect(footer.GetStatusMessage()).To(Equal("5/20 items"))
		})

		It("should return empty string when no status is set", func() {
			Expect(footer.GetStatusMessage()).To(BeEmpty())
		})
	})

	Describe("GetModeContext", func() {
		It("should return the mode context when set", func() {
			footer.WithMode("Timeline View")
			Expect(footer.GetModeContext()).To(Equal("Timeline View"))
		})

		It("should return empty string when no mode is set", func() {
			Expect(footer.GetModeContext()).To(BeEmpty())
		})
	})

	Describe("View", func() {
		Context("when width is zero or negative", func() {
			It("should return empty string for zero width", func() {
				f := layout.NewFooter(0).WithTheme(theme).WithStatus("status")
				Expect(f.View()).To(BeEmpty())
			})

			It("should return empty string for negative width", func() {
				f := layout.NewFooter(-1).WithTheme(theme).WithStatus("status")
				Expect(f.View()).To(BeEmpty())
			})
		})

		Context("when only status is shown", func() {
			It("should render the status message", func() {
				footer.WithTheme(theme).WithStatus("3/10 events")
				footer.SetShowMode(false)
				view := stripAnsi(footer.View())
				Expect(view).To(ContainSubstring("3/10 events"))
			})
		})

		Context("when only mode is shown", func() {
			It("should render the mode context", func() {
				footer.WithTheme(theme).WithMode("Capture Mode: Timeline")
				footer.SetShowStatus(false)
				view := stripAnsi(footer.View())
				Expect(view).To(ContainSubstring("Capture Mode: Timeline"))
			})
		})

		Context("when both status and mode are shown", func() {
			It("should render both on the same line with separator", func() {
				footer.WithTheme(theme).
					WithStatus("3/10 events").
					WithMode("Timeline")
				view := stripAnsi(footer.View())
				Expect(view).To(ContainSubstring("3/10 events"))
				Expect(view).To(ContainSubstring("Timeline"))
				Expect(view).To(ContainSubstring("|"))
			})

			It("should fall back to status only when combined line exceeds width", func() {
				narrowFooter := layout.NewFooter(20).WithTheme(theme).
					WithStatus("status text").
					WithMode("a very long mode context that exceeds width")
				view := stripAnsi(narrowFooter.View())
				Expect(view).To(ContainSubstring("status text"))
			})
		})

		Context("when help text is set", func() {
			It("should render help text", func() {
				footer.WithTheme(theme).WithHelp("Press q to quit")
				view := stripAnsi(footer.View())
				Expect(view).To(ContainSubstring("Press q to quit"))
			})

			It("should render help alongside status", func() {
				footer.WithTheme(theme).
					WithStatus("3/10 events").
					WithHelp("Press q to quit")
				view := stripAnsi(footer.View())
				Expect(view).To(ContainSubstring("3/10 events"))
				Expect(view).To(ContainSubstring("Press q to quit"))
			})
		})

		Context("when status text exceeds width", func() {
			It("should truncate the status with ellipsis", func() {
				longStatus := strings.Repeat("A", 200)
				f := layout.NewFooter(50).WithTheme(theme).WithStatus(longStatus)
				f.SetShowMode(false)
				view := stripAnsi(f.View())
				Expect(view).To(ContainSubstring("..."))
			})
		})

		Context("when mode text exceeds width", func() {
			It("should truncate the mode with ellipsis", func() {
				longMode := strings.Repeat("B", 200)
				f := layout.NewFooter(50).WithTheme(theme).WithMode(longMode)
				f.SetShowStatus(false)
				view := stripAnsi(f.View())
				Expect(view).To(ContainSubstring("..."))
			})
		})

		Context("when using default theme", func() {
			It("should render without explicit theme set", func() {
				f := layout.NewFooter(80).WithStatus("status")
				view := stripAnsi(f.View())
				Expect(view).To(ContainSubstring("status"))
			})
		})
	})
})
