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

// State constant for state matrix tracking (REQUIRED)
const ConfirmState = "confirm"

// ConfirmScreen displays a confirmation dialog for saving configuration changes.
//
// Uses UIKit ButtonGroup for yes/no selection with keyboard navigation.
// Embeds BaseScreen for common functionality (terminal dimensions, theme, CreateView).
//
// Related:
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/UIKIT_GUIDE.md (ButtonGroup usage)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts)
type ConfirmScreen struct {
	*base.BaseScreen

	domain      configtypes.ConfigurationDomain
	changeCount int

	// Button group for yes/no selection
	buttonGroup *primitives.ButtonGroup

	// Theme for rendering (nil-safe via getTheme())
	theme themes.Theme
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

	// Create button group with Cancel as default (index 1)
	buttonGroup := primitives.NewButtonGroup(theme).
		AddPrimary("Save Changes").
		AddSecondary("Cancel")
	buttonGroup.FocusLast() // Default to Cancel (safer)

	return &ConfirmScreen{
		BaseScreen:  base.NewBaseScreen(),
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
	// Handle window resize via BaseScreen
	if cmd := s.BaseScreen.HandleWindowSizeMsg(msg); cmd != nil {
		return cmd, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "esc":
			// Cancel - go back
			return nil, &screens.CancelResult{}

		case "m":
			// Main menu - cancel intent entirely
			result := &screens.CancelResult{}
			result.WithMetadata("main_menu", true)
			return nil, result

		case "enter":
			// Confirm current selection
			// Index 0 = Save Changes (true), Index 1 = Cancel (false)
			confirmed := s.buttonGroup.FocusIndex() == 0
			return nil, &screens.NavigateResult{ResultData: confirmed}

		case "y", "Y":
			// Direct Yes
			return nil, &screens.NavigateResult{ResultData: true}

		case "n", "N":
			// Direct No
			return nil, &screens.NavigateResult{ResultData: false}

		case "left", "right", "h", "l", "tab":
			// Delegate navigation to ButtonGroup
			s.buttonGroup.Update(keyMsg)
			return nil, nil
		}
	}

	return nil, nil
}

// View renders the screen using UIKit layout and BaseScreen.CreateView().
func (s *ConfirmScreen) View() string {
	theme := s.getTheme()
	domainLabel := formatDomainLabel(s.domain)

	// Build content using UIKit primitives
	var content strings.Builder

	// Title
	content.WriteString(primitives.Title("Confirm Configuration Changes?", theme).Render())
	content.WriteString("\n\n")

	// Domain info
	domainText := fmt.Sprintf("Domain: %s", domainLabel)
	content.WriteString(primitives.Body(domainText, theme).Render())
	content.WriteString("\n")

	// Change count
	changesText := fmt.Sprintf("Changes: %d setting(s) modified", s.changeCount)
	content.WriteString(primitives.Body(changesText, theme).Render())
	content.WriteString("\n\n")

	// Buttons via ButtonGroup
	content.WriteString(s.buttonGroup.Render())
	content.WriteString("\n")

	// Help footer with UIKit badges (ALWAYS use predefined badge constructors)
	helpFooter := primitives.RenderHelpFooter(theme,
		primitives.YesBadge(theme),
		primitives.NoBadge(theme),
		primitives.ConfirmBadge(theme),
		primitives.ToggleBadge(theme),
		primitives.BackBadge(theme),
		primitives.MenuBadge(theme),
	)

	// Use BaseScreen.CreateView() for StandardView integration
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
		// Update button group theme
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
