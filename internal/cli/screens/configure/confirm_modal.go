package configure

import (
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmModal is a simple yes/no confirmation modal.
type ConfirmModal struct {
	title   string
	message string

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

// NewConfirmModal creates a new confirmation modal.
//
// Expected:
//   - Must be a valid string.
//   - int must be valid.
//
// Returns:
//   - A fully initialized ConfirmModal ready for use.
//
// Side effects:
//   - None.
func NewConfirmModal(title, message string, width, height int) *ConfirmModal {
	return &ConfirmModal{
		title:   title,
		message: message,
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
func (m *ConfirmModal) Init() tea.Cmd {
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
func (m *ConfirmModal) Update(msg tea.Msg) tea.Cmd {
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
		case "esc", "n", "q":
			m.cancelled = true
			m.visible = false
			return nil

		case "enter", "y":
			m.confirmed = true
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
func (m *ConfirmModal) View() string {
	if !m.visible {
		return ""
	}

	// Message
	messageText := primitives.Body(m.message, m.theme).Render()

	// Build help footer
	helpFooter := primitives.RenderHelpFooter(m.theme,
		primitives.YesBadge(m.theme),
		primitives.NoBadge(m.theme),
	)

	// Combine content
	content := lipgloss.JoinVertical(lipgloss.Center,
		messageText,
		"",
		helpFooter,
	)

	// Render in a box - smaller for confirmation
	boxWidth := 50
	if m.width < 60 {
		boxWidth = m.width - 10
	}

	box := containers.NewBox(m.theme).
		Title(m.title).
		Content(content).
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
func (m *ConfirmModal) Render(width, height int) string {
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
func (m *ConfirmModal) IsVisible() bool {
	return m.visible
}

// IsConfirmed returns whether the user confirmed.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ConfirmModal) IsConfirmed() bool {
	return m.confirmed
}

// IsCancelled returns whether the modal was cancelled.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ConfirmModal) IsCancelled() bool {
	return m.cancelled
}

// SetTheme updates the theme.
//
// Expected:
//   - th must be a valid theme instance (can be nil).
//
// Side effects:
//   - None.
func (m *ConfirmModal) SetTheme(theme themes.Theme) {
	m.theme = theme
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *ConfirmModal) Show() {
	m.visible = true
	m.confirmed = false
	m.cancelled = false
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *ConfirmModal) Hide() {
	m.visible = false
}
