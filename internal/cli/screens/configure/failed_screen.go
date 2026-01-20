package configure

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/configtypes"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
)

// State constant for state matrix tracking (REQUIRED)
const FailedState = "failed"

// FailedScreen displays an error message when configuration save fails.
//
// Embeds BaseScreen for common functionality (terminal dimensions, theme, CreateView).
// Supports retry action via Enter or 'r' key.
//
// Related:
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/UIKIT_GUIDE.md (Text primitives)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts)
type FailedScreen struct {
	*base.BaseScreen

	domain       configtypes.ConfigurationDomain
	errorMessage string

	// Theme for rendering (nil-safe via getTheme())
	theme themes.Theme
}

// NewFailedScreen creates a new failed screen.
//
// Parameters:
//   - domain: The configuration domain that failed to save
//   - errorMessage: The error message to display (defaults to generic message if empty)
func NewFailedScreen(domain configtypes.ConfigurationDomain, errorMessage string) *FailedScreen {
	if errorMessage == "" {
		errorMessage = "An unknown error occurred while saving configuration."
	}
	return &FailedScreen{
		BaseScreen:   base.NewBaseScreen(),
		domain:       domain,
		errorMessage: errorMessage,
		theme:        themes.NewDefaultTheme(),
	}
}

// Init initializes the screen.
func (s *FailedScreen) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (s *FailedScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	// Handle window resize via BaseScreen
	if cmd := s.BaseScreen.HandleWindowSizeMsg(msg); cmd != nil {
		return cmd, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Cancel - go back
			return nil, &screens.CancelResult{}

		case "m":
			// Main menu - cancel intent entirely
			result := &screens.CancelResult{}
			result.WithMetadata("main_menu", true)
			return nil, result

		case "enter", "r":
			// Retry
			return nil, &screens.NavigateResult{ResultData: "retry"}

		case "q":
			// Quick exit to menu (same as cancel)
			result := &screens.CancelResult{}
			result.WithMetadata("main_menu", true)
			return nil, result
		}
	}

	return nil, nil
}

// View renders the screen using UIKit layout and BaseScreen.CreateView().
func (s *FailedScreen) View() string {
	theme := s.getTheme()
	domainLabel := formatDomainLabel(s.domain)

	// Build content using UIKit primitives
	var content strings.Builder

	// Error icon and title
	content.WriteString(primitives.ErrorText("Configuration Failed", theme).Bold().Render())
	content.WriteString("\n\n")

	// Error message
	content.WriteString(primitives.Body(s.errorMessage, theme).Render())
	content.WriteString("\n\n")

	// Helpful hint
	content.WriteString(primitives.Muted("Press Enter or 'r' to retry, Esc to go back.", theme).Render())
	content.WriteString("\n")

	// Help footer with UIKit badges (ALWAYS use predefined badge constructors)
	helpFooter := primitives.RenderHelpFooter(theme,
		primitives.RetryEnterBadge(theme),
		primitives.BackBadge(theme),
		primitives.MenuBadge(theme),
	)

	// Use BaseScreen.CreateView() for StandardView integration
	breadcrumbs := []string{"Main Menu", "Configure System", domainLabel, "Error"}
	return s.CreateView(breadcrumbs, content.String(), helpFooter)
}

// getTheme returns the theme with nil-safe fallback.
func (s *FailedScreen) getTheme() themes.Theme {
	if s.theme == nil {
		return themes.NewDefaultTheme()
	}
	return s.theme
}

// SetTheme updates the theme.
func (s *FailedScreen) SetTheme(theme interface{}) {
	if t, ok := theme.(themes.Theme); ok {
		s.theme = t
	}
}

// GetErrorMessage returns the error message.
func (s *FailedScreen) GetErrorMessage() string {
	return s.errorMessage
}
