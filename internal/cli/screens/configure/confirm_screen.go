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

// ConfirmState identifies the configuration save confirmation screen in the
// state matrix. On this screen the user sees a title asking to confirm
// changes, the configuration domain name, the number of modified settings,
// and a two-button group (Save Changes / Cancel) rendered via UIKit
// ButtonGroup. Left/right arrows or h/l toggle between buttons, Enter
// submits the focused choice, y/n act as direct shortcuts, and Escape
// cancels without saving.
const ConfirmState = "confirm"

// ConfirmScreen displays a confirmation dialog for saving configuration changes.
//
// Uses UIKit ButtonGroup for yes/no selection with keyboard navigation.
// Embeds Screen for common functionality (terminal dimensions, theme, CreateView).
//
// Related:
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/UIKIT_GUIDE.md (ButtonGroup usage)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts).
type ConfirmScreen struct {
	*base.Screen

	domain      configtypes.ConfigurationDomain
	changeCount int

	buttonGroup *primitives.ButtonGroup
	theme       themes.Theme
}

// NewConfirmScreen creates a new confirmation screen.
//
// Parameters:
//   - domain: The configuration domain being modified
//   - changeCount: Number of settings that will be changed
//
// Default behavior:
//   - Starts with "Cancel" selected (safer default)
//   - Left/Right arrows and h/l toggle selection
//   - Enter confirms current selection
//   - y/n keys submit directly
//   - Escape cancels
//   - 'm' returns to main menu
func NewConfirmScreen(domain configtypes.ConfigurationDomain, changeCount int) *ConfirmScreen {
	theme := themes.NewDefaultTheme()

	buttonGroup := primitives.NewButtonGroup(theme).
		AddPrimary("Save Changes").
		AddSecondary("Cancel")
	buttonGroup.FocusLast()

	return &ConfirmScreen{
		Screen:      base.NewBaseScreen(),
		domain:      domain,
		changeCount: changeCount,
		buttonGroup: buttonGroup,
		theme:       theme,
	}
}

// Init initializes the screen.
func (s *ConfirmScreen) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (s *ConfirmScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	if cmd := s.Screen.HandleWindowSizeMsg(msg); cmd != nil {
		return cmd, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "esc":
			return nil, &screens.CancelResult{}

		case "enter":
			confirmed := s.buttonGroup.FocusIndex() == 0
			return nil, &screens.NavigateResult{ResultData: confirmed}

		case "y", "Y":
			return nil, &screens.NavigateResult{ResultData: true}

		case "n", "N":
			return nil, &screens.NavigateResult{ResultData: false}

		case "left", "right", "h", "l", "tab":
			s.buttonGroup.Update(keyMsg)
			return nil, nil
		}
	}

	return nil, nil
}

// View renders the screen using UIKit layout and Screen.CreateView().
func (s *ConfirmScreen) View() string {
	theme := s.getTheme()
	domainLabel := formatDomainLabel(s.domain)

	var content strings.Builder

	content.WriteString(primitives.Title("Confirm Configuration Changes?", theme).Render())
	content.WriteString("\n\n")

	domainText := fmt.Sprintf("Domain: %s", domainLabel)
	content.WriteString(primitives.Body(domainText, theme).Render())
	content.WriteString("\n")

	changesText := fmt.Sprintf("Changes: %d setting(s) modified", s.changeCount)
	content.WriteString(primitives.Body(changesText, theme).Render())
	content.WriteString("\n\n")

	content.WriteString(s.buttonGroup.Render())
	content.WriteString("\n")

	helpFooter := primitives.RenderHelpFooter(theme,
		primitives.YesBadge(theme),
		primitives.NoBadge(theme),
		primitives.ConfirmBadge(theme),
		primitives.ToggleBadge(theme),
		primitives.BackBadge(theme),
	)

	breadcrumbs := []string{"Main Menu", "Configure System", domainLabel, "Confirm"}
	return s.CreateView(breadcrumbs, content.String(), helpFooter)
}

// getTheme returns the theme with nil-safe fallback.
func (s *ConfirmScreen) getTheme() themes.Theme {
	if s.theme == nil {
		return themes.NewDefaultTheme()
	}
	return s.theme
}

// SetTheme updates the theme.
func (s *ConfirmScreen) SetTheme(theme interface{}) {
	if t, ok := theme.(themes.Theme); ok {
		s.theme = t
		if s.buttonGroup != nil {
			s.buttonGroup.SetTheme(t)
		}
	}
}

// GetSelection returns the current selection (true = Yes/Save, false = No/Cancel).
func (s *ConfirmScreen) GetSelection() bool {
	return s.buttonGroup.FocusIndex() == 0
}

// SetSelection sets the current selection.
func (s *ConfirmScreen) SetSelection(yes bool) {
	if yes {
		s.buttonGroup.FocusFirst()
	} else {
		s.buttonGroup.FocusLast()
	}
}
