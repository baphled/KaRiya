package base

import (
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/uikit/layout"
	tea "github.com/charmbracelet/bubbletea"
)

// LogoModel defines the interface for logo components.
// The uikit/display.Logo package provides the standard implementation.
type LogoModel interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
	ViewStatic() string
	SetWidth(width int)
}

// Screen provides common functionality for all Screen implementations.
//
// Screens should embed Screen to get:
// - Terminal dimension management (width, height)
// - Theme management
// - StandardView creation helpers
// - WindowSizeMsg handling
//
// Example usage:
//
//	type MyScreen struct {
//	    *base.Screen
//	    // ... screen-specific fields
//	}
//
//	func NewMyScreen() *MyScreen {
//	    return &MyScreen{
//	        Screen: base.NewBaseScreen(),
//	    }
//	}
//
//	func (s *MyScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
//	    // Handle WindowSizeMsg automatically
//	    if cmd := s.Screen.HandleWindowSizeMsg(msg); cmd != nil {
//	        return cmd, nil
//	    }
//
//	    // ... screen-specific update logic
//	}
//
//	func (s *MyScreen) View() string {
//	    return s.CreateView(
//	        "My Screen",           // breadcrumb
//	        "Screen content here", // content
//	        "↑/↓: Navigate",       // footer
//	    )
//	}
//
// Related:
// - internal/cli/screens/contract.go (Screen interface)
// - internal/cli/components/standard_view.go (StandardView).
type Screen struct {
	// terminalWidth is the current terminal width in characters
	terminalWidth int

	// terminalHeight is the current terminal height in characters
	terminalHeight int

	// theme holds the current theme for styling.
	// Using interface{} until theme system type is finalized.
	theme interface{}

	// logo holds the logo to display (shared from intent)
	logo LogoModel

	// logoSpacing is the vertical spacing before the logo
	logoSpacing int
}

// NewBaseScreen creates a new Screen with default dimensions.
//
// Returns:
//   - A fully initialized Screen ready for use.
//
// Side effects:
//   - None.
func NewBaseScreen() *Screen {
	return &Screen{
		terminalWidth:  120,
		terminalHeight: 40,
	}
}

// SetTerminalInfo updates the terminal dimensions.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (b *Screen) SetTerminalInfo(width, height int) {
	b.terminalWidth = width
	b.terminalHeight = height
}

// SetTheme updates the theme used for styling.
//
// Expected:
//   - interface{} must be valid.
//
// Side effects:
//   - None.
func (b *Screen) SetTheme(theme interface{}) {
	b.theme = theme
}

// SetLogo sets the logo to be displayed in views.
//
// Expected:
//   - interface{} must be valid.
//   - int must be valid.
//
// Side effects:
//   - None.
func (b *Screen) SetLogo(logo interface{}, spacing int) {
	// Type assert to LogoModel interface
	if logoModel, ok := logo.(LogoModel); ok {
		b.logo = logoModel
		b.logoSpacing = spacing
	}
}

// GetLogo returns the currently set logo.
//
// Returns:
//   - A LogoModel value.
//
// Side effects:
//   - None.
func (b *Screen) GetLogo() LogoModel {
	return b.logo
}

// GetLogoSpacing returns the logo spacing.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (b *Screen) GetLogoSpacing() int {
	return b.logoSpacing
}

// Width returns the current terminal width.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (b *Screen) Width() int {
	return b.terminalWidth
}

// Height returns the current terminal height.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (b *Screen) Height() int {
	return b.terminalHeight
}

// Theme returns the current theme.
//
// Returns:
//   - A interface{} value.
//
// Side effects:
//   - None.
func (b *Screen) Theme() interface{} {
	return b.theme
}

// CreateView is a helper method to create a StandardView with current dimensions and theme.
//
// Expected:
//   - Must be a valid string.
//   - Must be a valid string.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (b *Screen) CreateView(breadcrumbs []string, content, footer string) string {
	// Create terminal info from current dimensions
	termInfo := &terminal.Info{
		Width:  b.terminalWidth,
		Height: b.terminalHeight,
	}

	// Build ScreenLayout using builder pattern
	view := layout.NewScreenLayout(termInfo).
		WithBreadcrumbs(breadcrumbs...).
		WithContent(content).
		WithHelp(footer).
		WithFooterSeparator(true)

	// Add logo if available
	if b.logo != nil {
		view = view.WithLogo(b.logo, b.logoSpacing)
	}

	return view.Render()
}

// HandleWindowSizeMsg is a helper to handle WindowSizeMsg uniformly across all screens.
//
// Expected:
//   - msg must be valid.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (b *Screen) HandleWindowSizeMsg(msg tea.Msg) tea.Cmd {
	if wsm, ok := msg.(tea.WindowSizeMsg); ok {
		b.SetTerminalInfo(wsm.Width, wsm.Height)
	}
	return nil
}

// RenderContent is a default implementation that returns empty string.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (b *Screen) RenderContent() string {
	return ""
}
