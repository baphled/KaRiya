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
const CompleteState = "complete"

// CompleteScreen displays a success message after configuration is saved.
//
// Embeds BaseScreen for common functionality (terminal dimensions, theme, CreateView).
//
// Related:
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/UIKIT_GUIDE.md (Text primitives)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts)
type CompleteScreen struct {
	*base.BaseScreen

	domain      configtypes.ConfigurationDomain
	changeCount int

	// Theme for rendering (nil-safe via getTheme())
	theme themes.Theme
}

// NewCompleteScreen creates a new complete screen.
//
// Parameters:
//   - domain: The configuration domain that was modified
//   - changeCount: Number of settings that were changed
func NewCompleteScreen(domain configtypes.ConfigurationDomain, changeCount int) *CompleteScreen {
	return &CompleteScreen{
		BaseScreen:  base.NewBaseScreen(),
		domain:      domain,
		changeCount: changeCount,
		theme:       themes.NewDefaultTheme(),
	}
}

// Init initializes the screen.
func (s *CompleteScreen) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (s *CompleteScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	// Handle window resize via BaseScreen
	if cmd := s.BaseScreen.HandleWindowSizeMsg(msg); cmd != nil {
		return cmd, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "esc":
			// Both Enter and Esc dismiss the success screen
			return nil, &screens.SubmitResult{FormData: true}

		case "m":
			// Main menu
			result := &screens.SubmitResult{FormData: true}
			result.WithMetadata("main_menu", true)
			return nil, result

		case "q":
			// Quick exit to menu (same as success dismiss)
			return nil, &screens.SubmitResult{FormData: true}
		}
	}

	return nil, nil
}

// View renders the screen using UIKit layout and BaseScreen.CreateView().
func (s *CompleteScreen) View() string {
	theme := s.getTheme()
	domainLabel := formatDomainLabel(s.domain)

	// Build content using UIKit primitives
	var content strings.Builder

	// Success icon and title
	content.WriteString(primitives.SuccessText("Configuration Updated!", theme).Bold().Render())
	content.WriteString("\n\n")

	// Domain info
	domainText := fmt.Sprintf("Domain: %s", domainLabel)
	content.WriteString(primitives.Body(domainText, theme).Render())
	content.WriteString("\n")

	// Change count
	if s.changeCount > 0 {
		changesText := fmt.Sprintf("%d setting(s) saved successfully.", s.changeCount)
		content.WriteString(primitives.Body(changesText, theme).Render())
	} else {
		content.WriteString(primitives.Muted("No changes were made.", theme).Render())
	}
	content.WriteString("\n")

	// Help footer with UIKit badges (ALWAYS use predefined badge constructors)
	helpFooter := primitives.RenderHelpFooter(theme,
		primitives.ContinueBadge(theme),
		primitives.MenuBadge(theme),
	)

	// Use BaseScreen.CreateView() for StandardView integration
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
func (s *CompleteScreen) SetTheme(theme interface{}) {
	if t, ok := theme.(themes.Theme); ok {
		s.theme = t
	}
}
