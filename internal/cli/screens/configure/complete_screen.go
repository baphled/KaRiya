package configure

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/configtypes"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
)

// CompleteState identifies the configuration completion screen in the state
// matrix. On this screen the user sees a bold success title, the name of
// the configuration domain that was modified, and the count of settings
// saved. If no changes were made the screen shows a muted notice instead.
// Pressing Enter, Escape, or q dismisses the screen and returns to the
// configuration menu.
const CompleteState = "complete"

// CompleteScreen displays a success message after configuration is saved.
//
// Embeds Screen for common functionality (terminal dimensions, theme, CreateView).
//
// Related:
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/UIKIT_GUIDE.md (Text primitives)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts).
type CompleteScreen struct {
	*base.Screen

	domain      configtypes.ConfigurationDomain
	changeCount int

	theme themes.Theme
}

// NewCompleteScreen creates a new complete screen.
//
// Expected:
//   - config must be a valid configuration object.
//   - int must be valid.
//
// Returns:
//   - A fully initialized CompleteScreen ready for use.
//
// Side effects:
//   - None.
func NewCompleteScreen(domain configtypes.ConfigurationDomain, changeCount int) *CompleteScreen {
	return &CompleteScreen{
		Screen:      base.NewBaseScreen(),
		domain:      domain,
		changeCount: changeCount,
		theme:       themes.NewDefaultTheme(),
	}
}

// Init initializes the screen.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (s *CompleteScreen) Init() tea.Cmd {
	return nil
}

// Update handles messages.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Cmd: command to execute.
//   - screens.ScreenResult: result indicating dismissal.
//
// Side effects:
//   - May return SubmitResult on enter/esc/q.
func (s *CompleteScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	if cmd := s.Screen.HandleWindowSizeMsg(msg); cmd != nil {
		return cmd, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "enter", "esc":
			return nil, &screens.SubmitResult{FormData: true}

		case "q":
			return nil, &screens.SubmitResult{FormData: true}
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
func (s *CompleteScreen) View() string {
	theme := s.getTheme()
	domainLabel := formatDomainLabel(s.domain)

	var content strings.Builder

	content.WriteString(primitives.SuccessText("Configuration Updated!", theme).Bold().Render())
	content.WriteString("\n\n")

	domainText := fmt.Sprintf("Domain: %s", domainLabel)
	content.WriteString(primitives.Body(domainText, theme).Render())
	content.WriteString("\n")

	if s.changeCount > 0 {
		changesText := fmt.Sprintf("%d setting(s) saved successfully.", s.changeCount)
		content.WriteString(primitives.Body(changesText, theme).Render())
	} else {
		content.WriteString(primitives.Muted("No changes were made.", theme).Render())
	}
	content.WriteString("\n")

	helpFooter := primitives.RenderHelpFooter(theme,
		primitives.ContinueBadge(theme),
	)

	breadcrumbs := []string{"Main Menu", "Configure System", domainLabel, "Complete"}
	return s.CreateView(breadcrumbs, content.String(), helpFooter)
}

// getTheme returns the theme with nil-safe fallback.
func (s *CompleteScreen) getTheme() themes.Theme {
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
func (s *CompleteScreen) SetTheme(theme interface{}) {
	if t, ok := theme.(themes.Theme); ok {
		s.theme = t
	}
}
