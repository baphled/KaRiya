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

// FailedState identifies the configuration failure screen in the state
// matrix. On this screen the user sees a bold error title, the error
// message explaining why the save failed, and a muted hint about
// available actions. Pressing Enter or r retries the save, Escape goes
// back to the previous step, and q exits to the main menu.
const FailedState = "failed"

// FailedScreen displays an error message when configuration save fails.
//
// Embeds Screen for common functionality (terminal dimensions, theme, CreateView).
// Supports retry action via Enter or 'r' key.
//
// Related:
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/UIKIT_GUIDE.md (Text primitives)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts).
type FailedScreen struct {
	*base.Screen

	domain       configtypes.ConfigurationDomain
	errorMessage string

	theme themes.Theme
}

// NewFailedScreen creates a new failed screen.
//
// Expected:
//   - config must be a valid configuration object.
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized FailedScreen ready for use.
//
// Side effects:
//   - None.
func NewFailedScreen(domain configtypes.ConfigurationDomain, errorMessage string) *FailedScreen {
	if errorMessage == "" {
		errorMessage = "An unknown error occurred while saving configuration."
	}
	return &FailedScreen{
		Screen:       base.NewBaseScreen(),
		domain:       domain,
		errorMessage: errorMessage,
		theme:        themes.NewDefaultTheme(),
	}
}

// Init initializes the screen.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (s *FailedScreen) Init() tea.Cmd {
	return nil
}

// Update handles messages.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Cmd: command to execute.
//   - screens.ScreenResult: result indicating user action.
//
// Side effects:
//   - May return CancelResult or NavigateResult.
func (s *FailedScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	if cmd := s.HandleWindowSizeMsg(msg); cmd != nil {
		return cmd, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "esc":
			return nil, &screens.CancelResult{}

		case "enter", "r":
			return nil, &screens.NavigateResult{ResultData: "retry"}

		case "q":
			result := &screens.CancelResult{}
			result.WithMetadata("main_menu", true)
			return nil, result
		}
	}

	return nil, nil
}

// View renders the screen using UIKit layout and Screen.CreateView().
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *FailedScreen) View() string {
	theme := s.getTheme()
	domainLabel := formatDomainLabel(s.domain)

	var content strings.Builder

	content.WriteString(primitives.ErrorText("Configuration Failed", theme).Bold().Render())
	content.WriteString("\n\n")

	content.WriteString(primitives.Body(s.errorMessage, theme).Render())
	content.WriteString("\n\n")

	content.WriteString(primitives.Muted("Press Enter or 'r' to retry, Esc to go back.", theme).Render())
	content.WriteString("\n")

	helpFooter := primitives.RenderHelpFooter(theme,
		primitives.RetryEnterBadge(theme),
		primitives.BackBadge(theme),
	)

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
//
// Expected:
//   - interface{} must be valid.
//
// Side effects:
//   - None.
func (s *FailedScreen) SetTheme(theme interface{}) {
	if t, ok := theme.(themes.Theme); ok {
		s.theme = t
	}
}

// GetErrorMessage returns the error message.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *FailedScreen) GetErrorMessage() string {
	return s.errorMessage
}
