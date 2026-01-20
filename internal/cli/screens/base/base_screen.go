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

// BaseScreen provides common functionality for all Screen implementations.
//
// Screens should embed BaseScreen to get:
// - Terminal dimension management (width, height)
// - Theme management
// - StandardView creation helpers
// - WindowSizeMsg handling
//
// Example usage:
//
//	type MyScreen struct {
//	    *base.BaseScreen
//	    // ... screen-specific fields
//	}
//
//	func NewMyScreen() *MyScreen {
//	    return &MyScreen{
//	        BaseScreen: base.NewBaseScreen(),
//	    }
//	}
//
//	func (s *MyScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
//	    // Handle WindowSizeMsg automatically
//	    if cmd := s.BaseScreen.HandleWindowSizeMsg(msg); cmd != nil {
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
// - internal/cli/components/standard_view.go (StandardView)
type BaseScreen struct {
	// terminalWidth is the current terminal width in characters
	terminalWidth int

	// terminalHeight is the current terminal height in characters
	terminalHeight int

	// theme holds the current theme for styling
	// TODO: Replace interface{} with actual Theme type when theme system is defined
	theme interface{}

	// logo holds the logo to display (shared from intent)
	logo LogoModel

	// logoSpacing is the vertical spacing before the logo
	logoSpacing int
}

// NewBaseScreen creates a new BaseScreen with default dimensions.
//
// Default dimensions (120x40) are used until SetTerminalInfo is called
// with actual terminal dimensions.
func NewBaseScreen() *BaseScreen {
	return &BaseScreen{
		terminalWidth:  120, // Default width
		terminalHeight: 40,  // Default height
	}
}

// SetTerminalInfo updates the terminal dimensions.
//
// This should be called when the screen receives a WindowSizeMsg,
// or when the screen is initialized with known dimensions.
func (b *BaseScreen) SetTerminalInfo(width, height int) {
	b.terminalWidth = width
	b.terminalHeight = height
}

// SetTheme updates the theme used for styling.
//
// This should be called when the intent sets up the screen,
// passing the global or intent-specific theme.
func (b *BaseScreen) SetTheme(theme interface{}) {
	b.theme = theme
}

// SetLogo sets the logo to be displayed in views.
//
// This should be called when the intent sets up the screen,
// passing the shared logo instance and optional spacing.
// Accepts any LogoModel implementation (typically display.Logo).
func (b *BaseScreen) SetLogo(logo interface{}, spacing int) {
	// Type assert to LogoModel interface
	if logoModel, ok := logo.(LogoModel); ok {
		b.logo = logoModel
		b.logoSpacing = spacing
	}
}

// GetLogo returns the currently set logo.
func (b *BaseScreen) GetLogo() LogoModel {
	return b.logo
}

// GetLogoSpacing returns the logo spacing.
func (b *BaseScreen) GetLogoSpacing() int {
	return b.logoSpacing
}

// Width returns the current terminal width.
func (b *BaseScreen) Width() int {
	return b.terminalWidth
}

// Height returns the current terminal height.
func (b *BaseScreen) Height() int {
	return b.terminalHeight
}

// Theme returns the current theme.
func (b *BaseScreen) Theme() interface{} {
	return b.theme
}

// CreateView is a helper method to create a StandardView with current dimensions and theme.
//
// Parameters:
//   - breadcrumbs: Variadic breadcrumb trail (e.g., "Main Menu", "Generate CV", "Select Profile")
//   - content: The main content of the screen
//   - footer: The footer text (shortcuts, help, etc.)
//
// Example:
//
//	view := s.CreateView(
//	    []string{"Main Menu", "My Feature"},
//	    s.renderContent(),
//	    "Enter: Select  Esc: Back  q: Quit",
//	)
//
// This automatically uses the current terminal dimensions and theme.
func (b *BaseScreen) CreateView(breadcrumbs []string, content, footer string) string {
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

	// TODO: Apply theme to view when theme system is fully integrated
	// if b.theme != nil {
	//     view = view.WithTheme(b.theme)
	// }

	return view.Render()
}

// HandleWindowSizeMsg is a helper to handle WindowSizeMsg uniformly across all screens.
//
// Call this at the start of your screen's Update method:
//
//	func (s *MyScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
//	    if cmd := s.BaseScreen.HandleWindowSizeMsg(msg); cmd != nil {
//	        return cmd, nil
//	    }
//	    // ... rest of update logic
//	}
//
// Returns nil if msg is not a WindowSizeMsg.
// Returns a command (usually nil) if msg is a WindowSizeMsg.
// Never returns a ScreenResult for WindowSizeMsg (window resize is not a user action).
func (b *BaseScreen) HandleWindowSizeMsg(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.SetTerminalInfo(msg.Width, msg.Height)
		return nil
	}
	return nil
}

// RenderContent is a default implementation that returns empty string.
// Screens should override this method to provide their content.
//
// This method allows intents to get just the content without StandardView wrapper,
// enabling them to apply their own StandardView with custom breadcrumbs and help.
func (b *BaseScreen) RenderContent() string {
	return ""
}
