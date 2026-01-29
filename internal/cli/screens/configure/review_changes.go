package configure

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/configtypes"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/layout"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ReviewChangesState identifies the configuration review screen in the state
// matrix. On this screen the user sees a summary count of modified settings
// followed by a three-column table listing each setting name, its original
// value, and the proposed new value. Pressing Enter, c, or y advances to
// the confirmation step. Escape returns to the editing form. When no
// changes exist the table is replaced with a muted notice.
const ReviewChangesState = "review_changes"

// ReviewChangesScreen displays a summary of configuration changes for user review.
type ReviewChangesScreen struct {
	domain configtypes.ConfigurationDomain

	// changes maps setting key to proposed new value.
	changes map[string]interface{}

	// original maps setting key to its value before editing.
	original map[string]interface{}

	// labels maps setting key to a human-readable label.
	labels map[string]string

	termInfo *terminal.Info
	theme    themes.Theme
	logo     layout.LogoRenderer
}

// NewReviewChangesScreen creates a new review changes screen.
func NewReviewChangesScreen(
	domain configtypes.ConfigurationDomain,
	changes map[string]interface{},
	original map[string]interface{},
	labels map[string]string,
) *ReviewChangesScreen {
	return &ReviewChangesScreen{
		domain:   domain,
		changes:  changes,
		original: original,
		labels:   labels,
		termInfo: &terminal.Info{Width: 120, Height: 40},
		theme:    themes.NewDefaultTheme(),
	}
}

// Init initializes the screen.
func (s *ReviewChangesScreen) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (s *ReviewChangesScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.termInfo = &terminal.Info{Width: msg.Width, Height: msg.Height}
		return nil, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return nil, &screens.CancelResult{}

		case tea.KeyEnter:
			return nil, &screens.NavigateResult{ResultData: "confirm"}
		}

		switch msg.String() {
		case "c", "y":
			return nil, &screens.NavigateResult{ResultData: "confirm"}
		}
	}

	return nil, nil
}

// View renders the screen using UIKit layout.
func (s *ReviewChangesScreen) View() string {
	domainLabel := formatDomainLabel(s.domain)

	var content strings.Builder

	if len(s.changes) == 0 {
		noChanges := primitives.Muted("No changes to review.", s.theme).Render()
		content.WriteString(noChanges)
	} else {
		countText := fmt.Sprintf("%d setting(s) changed", len(s.changes))
		summary := primitives.Body(countText, s.theme).Render()
		content.WriteString(summary)
		content.WriteString("\n\n")

		content.WriteString(s.renderChangesTable())
	}

	var helpFooter string
	if len(s.changes) == 0 {
		helpFooter = primitives.RenderHelpFooter(s.theme,
			primitives.BackBadge(s.theme),
		)
	} else {
		helpFooter = primitives.RenderHelpFooter(s.theme,
			primitives.ConfirmBadge(s.theme),
			primitives.BackBadge(s.theme),
		)
	}

	screenLayout := layout.NewScreenLayout(s.termInfo).
		WithTheme(s.theme).
		WithBreadcrumbs("Main Menu", "Configure System", domainLabel, "Review Changes").
		WithContent(content.String()).
		WithHelp(helpFooter).
		WithFooterSeparator(true)

	if s.logo != nil {
		screenLayout = screenLayout.WithLogo(s.logo, 2)
	}

	return screenLayout.Render()
}

// renderChangesTable renders the changes as a formatted table.
func (s *ReviewChangesScreen) renderChangesTable() string {
	var rows []string

	headerStyle := lipgloss.NewStyle().
		Foreground(s.theme.MutedColor()).
		Bold(true)

	labelCol := headerStyle.Width(20).Render("Setting")
	fromCol := headerStyle.Width(25).Render("Original")
	toCol := headerStyle.Width(25).Render("New Value")

	header := labelCol + " │ " + fromCol + " │ " + toCol
	rows = append(rows, header)

	separator := strings.Repeat("─", 20) + "─┼─" + strings.Repeat("─", 25) + "─┼─" + strings.Repeat("─", 25)
	rows = append(rows, separator)

	labelStyle := lipgloss.NewStyle().
		Foreground(s.theme.ForegroundColor()).
		Width(20)

	fromStyle := lipgloss.NewStyle().
		Foreground(s.theme.MutedColor()).
		Width(25)

	toStyle := lipgloss.NewStyle().
		Foreground(s.theme.SuccessColor()).
		Width(25)

	for key, newValue := range s.changes {
		label := s.labels[key]
		if label == "" {
			label = key
		}

		original := s.original[key]
		originalStr := fmt.Sprintf("%v", original)
		newStr := fmt.Sprintf("%v", newValue)

		row := labelStyle.Render(truncate(label, 18)) + " │ " +
			fromStyle.Render(truncate(originalStr, 23)) + " │ " +
			toStyle.Render(truncate(newStr, 23))

		rows = append(rows, row)
	}

	return strings.Join(rows, "\n")
}

// truncate truncates a string to the given length, adding "..." if needed.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// SetTerminalInfo updates terminal dimensions.
func (s *ReviewChangesScreen) SetTerminalInfo(width, height int) {
	s.termInfo = &terminal.Info{Width: width, Height: height}
}

// SetTheme updates the theme.
func (s *ReviewChangesScreen) SetTheme(theme interface{}) {
	if t, ok := theme.(themes.Theme); ok {
		s.theme = t
	}
}

// SetLogo updates the logo.
func (s *ReviewChangesScreen) SetLogo(logo interface{}, _ int) {
	if l, ok := logo.(layout.LogoRenderer); ok {
		s.logo = l
	}
}

// GetChanges returns the changes map.
func (s *ReviewChangesScreen) GetChanges() map[string]interface{} {
	return s.changes
}

// HasChanges returns true if there are any changes.
func (s *ReviewChangesScreen) HasChanges() bool {
	return len(s.changes) > 0
}
