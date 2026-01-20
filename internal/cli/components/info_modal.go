package components

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InfoModalVariant defines the visual style variant of the modal.
type InfoModalVariant int

const (
	// InfoModalInfo uses informational styling (teal/blue border).
	InfoModalInfo InfoModalVariant = iota
	// InfoModalWarning uses warning styling (amber/yellow border).
	InfoModalWarning
	// InfoModalSuccess uses success styling (green border).
	InfoModalSuccess
)

// InfoModal is a reusable informational modal for displaying messages to the user.
// Unlike DeleteConfirmModal, it has a single action: dismiss. It does not ask for
// confirmation - it simply displays information and allows the user to acknowledge it.
//
// Usage:
//
//	modal := components.NewInfoModal("Title", "Message explaining something...")
//	modal := components.NewWarningInfoModal("Warning Title", "Warning message...")
//
//	// In Update:
//	if modal.Update(msg) {
//	    // User dismissed the modal
//	    modal = nil
//	}
//
//	// In View:
//	if modal != nil && modal.IsVisible() {
//	    return modal.View()
//	}
//
// Dismisses on: Enter, Space, Esc
type InfoModal struct {
	title   string
	message string
	visible bool
	variant InfoModalVariant
	theme   themes.Theme
	width   int
	height  int
}

// NewInfoModal creates a new info modal with informational styling (teal/blue border).
//
// Parameters:
//   - title: Modal title (displayed prominently at top)
//   - message: Message body (can include newlines for formatting)
//
// Returns a new InfoModal ready to use.
func NewInfoModal(title, message string) *InfoModal {
	return &InfoModal{
		title:   title,
		message: message,
		visible: true,
		variant: InfoModalInfo,
		theme:   themes.NewDefaultTheme(),
		width:   100,
		height:  24,
	}
}

// NewWarningInfoModal creates a new info modal with warning styling (amber/yellow border).
// Use this when the message is a warning or caution to the user.
//
// Parameters:
//   - title: Modal title (displayed prominently at top)
//   - message: Warning message (can include newlines for formatting)
//
// Returns a new InfoModal with warning variant.
func NewWarningInfoModal(title, message string) *InfoModal {
	return &InfoModal{
		title:   title,
		message: message,
		visible: true,
		variant: InfoModalWarning,
		theme:   themes.NewDefaultTheme(),
		width:   100,
		height:  24,
	}
}

// NewSuccessInfoModal creates a new info modal with success styling (green border).
// Use this when displaying a success or completion message.
//
// Parameters:
//   - title: Modal title (displayed prominently at top)
//   - message: Success message (can include newlines for formatting)
//
// Returns a new InfoModal with success variant.
func NewSuccessInfoModal(title, message string) *InfoModal {
	return &InfoModal{
		title:   title,
		message: message,
		visible: true,
		variant: InfoModalSuccess,
		theme:   themes.NewDefaultTheme(),
		width:   100,
		height:  24,
	}
}

// Init initializes the modal (required by BubbleTea lifecycle).
func (m *InfoModal) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input for the info modal.
//
// Returns true if the modal was dismissed, false otherwise.
//
// The modal dismisses on:
//   - Enter: acknowledges the message
//   - Space: acknowledges the message
//   - Esc: closes the modal
func (m *InfoModal) Update(msg tea.Msg) bool {
	if !m.visible {
		return false
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return false

	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ", "esc":
			// User dismissed the modal
			m.visible = false
			return true
		}
	}

	return false
}

// View renders the info modal as a centered box.
//
// The modal displays:
//   - Title (styled based on variant)
//   - Message body
//   - Dismissal instructions
//
// Returns empty string if modal is not visible.
func (m *InfoModal) View() string {
	if !m.visible {
		return ""
	}

	// Get border color based on variant
	borderColor := m.getBorderColor()

	// Build footer with KeyBadge components
	footer := RenderHelpFooter(m.theme,
		NewKeyBadge("Enter/Esc", "Close"),
	)

	// Build modal content
	var content strings.Builder

	// Title (styled based on variant)
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(borderColor).
		MarginBottom(1)
	content.WriteString(titleStyle.Render(m.title))
	content.WriteString("\n\n")

	// Message (wrapped to modal width)
	messageStyle := lipgloss.NewStyle().
		Width(50).
		MarginBottom(1)
	content.WriteString(messageStyle.Render(m.message))
	content.WriteString("\n\n")

	// Footer
	content.WriteString(footer)

	// Create modal box (centered, with border)
	modalWidth := 60
	if modalWidth > m.width-4 {
		modalWidth = m.width - 4
	}
	if modalWidth < 40 {
		modalWidth = 40
	}

	modalStyle := lipgloss.NewStyle().
		Width(modalWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Background(m.theme.BackgroundColor()). // Solid background for overlay
		Padding(1, 2).
		Align(lipgloss.Center)

	modalBox := modalStyle.Render(content.String())

	return modalBox
}

// getBorderColor returns the appropriate border color based on variant.
func (m *InfoModal) getBorderColor() lipgloss.Color {
	switch m.variant {
	case InfoModalWarning:
		return m.theme.WarningColor()
	case InfoModalSuccess:
		return m.theme.SuccessColor()
	case InfoModalInfo:
		fallthrough
	default:
		return m.theme.InfoColor()
	}
}

// IsVisible returns whether the modal is currently visible.
func (m *InfoModal) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible.
func (m *InfoModal) Show() {
	m.visible = true
}

// Hide hides the modal.
func (m *InfoModal) Hide() {
	m.visible = false
}

// SetTheme sets the theme for the modal (useful for testing or custom themes).
func (m *InfoModal) SetTheme(theme themes.Theme) {
	m.theme = theme
}

// SetDimensions sets the terminal dimensions for responsive sizing.
func (m *InfoModal) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

// GetTitle returns the modal title.
func (m *InfoModal) GetTitle() string {
	return m.title
}

// GetMessage returns the modal message.
func (m *InfoModal) GetMessage() string {
	return m.message
}

// GetVariant returns the modal variant.
func (m *InfoModal) GetVariant() InfoModalVariant {
	return m.variant
}
