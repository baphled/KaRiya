package configure

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/configtypes"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ReviewChangesModal displays pending configuration changes for review.
type ReviewChangesModal struct {
	domain  configtypes.ConfigurationDomain
	changes map[string]interface{}

	// Dimensions
	width  int
	height int

	// Theme
	theme themes.Theme

	// Modal state
	visible   bool
	confirmed bool
	cancelled bool
}

// NewReviewChangesModal creates a new review changes modal.
//
// Expected:
//   - config must be a valid configuration object.
//   - Must be a valid string.
//   - int must be valid.
//
// Returns:
//   - A fully initialized ReviewChangesModal ready for use.
//
// Side effects:
//   - None.
func NewReviewChangesModal(domain configtypes.ConfigurationDomain, changes map[string]interface{}, width, height int) *ReviewChangesModal {
	return &ReviewChangesModal{
		domain:  domain,
		changes: changes,
		width:   width,
		height:  height,
		theme:   themes.NewDefaultTheme(),
		visible: true,
	}
}

// Init initializes the modal.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *ReviewChangesModal) Init() tea.Cmd {
	return nil
}

// Update handles messages.
//
// Expected:
//   - msg must be valid.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *ReviewChangesModal) Update(msg tea.Msg) tea.Cmd {
	if !m.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.cancelled = true
			m.visible = false
			return nil

		case "enter", "y":
			m.confirmed = true
			m.visible = false
			return nil

		case "n":
			m.cancelled = true
			m.visible = false
			return nil
		}
	}

	return nil
}

// View renders the modal content.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *ReviewChangesModal) View() string {
	if !m.visible {
		return ""
	}

	domainLabel := formatDomainLabel(m.domain)
	title := fmt.Sprintf("Review %s Changes", domainLabel)

	// Build changes list
	var content string
	if len(m.changes) == 0 {
		content = primitives.Muted("No changes to review.", m.theme).Render()
	} else {
		// Style for change rows
		keyStyle := lipgloss.NewStyle().
			Foreground(m.theme.PrimaryColor()).
			Bold(true)

		valueStyle := lipgloss.NewStyle().
			Foreground(m.theme.SuccessColor())

		var rows []string
		for key, value := range m.changes {
			row := fmt.Sprintf("  %s: %s",
				keyStyle.Render(key),
				valueStyle.Render(fmt.Sprintf("%v", value)),
			)
			rows = append(rows, row)
		}

		header := primitives.Subtitle(fmt.Sprintf("Changes to apply (%d):", len(m.changes)), m.theme).Render()
		content = lipgloss.JoinVertical(lipgloss.Left,
			header,
			"",
			lipgloss.JoinVertical(lipgloss.Left, rows...),
		)
	}

	// Build help footer
	helpFooter := primitives.RenderHelpFooter(m.theme,
		primitives.HelpKeyBadge("Enter/y", "Confirm", m.theme),
		primitives.HelpKeyBadge("n/Esc", "Cancel", m.theme),
	)

	// Combine content
	fullContent := lipgloss.JoinVertical(lipgloss.Left,
		content,
		"",
		helpFooter,
	)

	// Render in a box
	boxWidth := m.width - 20
	if boxWidth < 40 {
		boxWidth = 40
	}

	box := containers.NewBox(m.theme).
		Title(title).
		Content(fullContent).
		Width(boxWidth).
		Background(m.theme.BackgroundColor())

	return box.Render()
}

// Render renders the modal at the specified dimensions.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *ReviewChangesModal) Render(width, height int) string {
	m.width = width
	m.height = height
	return m.View()
}

// IsVisible returns whether the modal is visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ReviewChangesModal) IsVisible() bool {
	return m.visible
}

// IsConfirmed returns whether the changes were confirmed.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ReviewChangesModal) IsConfirmed() bool {
	return m.confirmed
}

// IsCancelled returns whether the modal was cancelled.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ReviewChangesModal) IsCancelled() bool {
	return m.cancelled
}

// SetTheme updates the theme.
//
// Expected:
//   - th must be a valid theme instance (can be nil).
//
// Side effects:
//   - None.
func (m *ReviewChangesModal) SetTheme(theme themes.Theme) {
	m.theme = theme
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *ReviewChangesModal) Show() {
	m.visible = true
	m.confirmed = false
	m.cancelled = false
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *ReviewChangesModal) Hide() {
	m.visible = false
}

// GetChanges returns the changes being reviewed.
//
// Returns:
//   - A map[string]interface{} value.
//
// Side effects:
//   - None.
func (m *ReviewChangesModal) GetChanges() map[string]interface{} {
	return m.changes
}
