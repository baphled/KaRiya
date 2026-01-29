package base

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/screens"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmScreen provides a reusable confirmation dialog screen.
//
// This screen handles:
// - Yes/No selection with keyboard navigation
// - Left/Right arrows and h/l (vim) to toggle selection
// - Enter to confirm current selection
// - y/n keys for direct submission
// - Escape to cancel
// - Customizable button text
// - StandardView integration
//
// Example usage:
//
//	screen := base.NewBaseConfirmScreen(
//	    []string{"Main Menu", "Delete Item"},
//	    "Delete Confirmation",
//	    "Are you sure you want to delete this item? This action cannot be undone.",
//	)
//	screen.SetYesText("Delete")
//	screen.SetNoText("Keep")
//
//	// In intent:
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultNavigate {
//	    confirmed := result.Data().(bool)
//	    if confirmed {
//	        // User confirmed - proceed with action
//	    } else {
//	        // User declined - go back
//	    }
//	}
//
// Related:
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts).
type ConfirmScreen struct {
	*Screen

	// breadcrumbs for navigation context
	breadcrumbs []string

	// title is the confirmation dialog title
	title string

	// message is the confirmation message/question
	message string

	// selectedYes indicates whether Yes is selected (true) or No is selected (false)
	selectedYes bool

	// yesText is the text for the Yes button (default: "Yes")
	yesText string

	// noText is the text for the No button (default: "No")
	noText string

	// footer is the help text shown at the bottom
	footer string
}

// NewBaseConfirmScreen creates a new confirmation screen.
//
// Parameters:
//   - breadcrumbs: Navigation breadcrumb trail (e.g., []string{"Main Menu", "Delete Item"})
//   - title: The confirmation dialog title
//   - message: The confirmation message/question
//
// Default behavior:
//   - Starts with "No" selected (safer default)
//   - Left/Right arrows and h/l toggle selection
//   - Enter confirms current selection
//   - y/n keys submit directly
//   - Escape cancels
func NewBaseConfirmScreen(
	breadcrumbs []string,
	title, message string,
) *ConfirmScreen {
	return &ConfirmScreen{
		Screen:      NewBaseScreen(),
		breadcrumbs: breadcrumbs,
		title:       title,
		message:     message,
		selectedYes: false,
		yesText:     "Yes",
		noText:      "No",
		footer:      "y/n: Choose  Enter: Confirm  ←→/hl: Toggle  Esc: Cancel",
	}
}

// Update handles messages and returns result when user makes a choice.
func (s *ConfirmScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle window resize
		s.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Cancel confirmation
			return nil, &screens.CancelResult{}

		case "enter":
			// Confirm current selection
			return nil, &screens.NavigateResult{
				ResultData: s.selectedYes,
			}

		case "left", "right", "h", "l":
			// Toggle selection
			s.selectedYes = !s.selectedYes
			return nil, nil

		case "y", "Y":
			// Direct Yes submission
			return nil, &screens.NavigateResult{
				ResultData: true,
			}

		case "n", "N":
			// Direct No submission
			return nil, &screens.NavigateResult{
				ResultData: false,
			}
		}
	}

	return nil, nil
}

// RenderContent returns the confirmation content without StandardView wrapper.
// This allows intents to wrap in their own StandardView with custom breadcrumbs/help.
func (s *ConfirmScreen) RenderContent() string {
	var b strings.Builder

	// Title (bold and centered)
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	b.WriteString(titleStyle.Render(s.title))
	b.WriteString("\n\n")

	// Message
	b.WriteString(s.message)
	b.WriteString("\n\n")

	// Buttons (side by side)
	yesButtonStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("10"))

	noButtonStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("9"))

	// Highlight selected button
	if s.selectedYes {
		yesButtonStyle = yesButtonStyle.Bold(true).Background(lipgloss.Color("10")).Foreground(lipgloss.Color("0"))
	} else {
		noButtonStyle = noButtonStyle.Bold(true).Background(lipgloss.Color("9")).Foreground(lipgloss.Color("0"))
	}

	yesButton := yesButtonStyle.Render(s.yesText)
	noButton := noButtonStyle.Render(s.noText)

	// Render buttons side by side with spacing
	buttons := lipgloss.JoinHorizontal(lipgloss.Left, yesButton, "  ", noButton)
	b.WriteString(buttons)
	b.WriteString("\n")

	return b.String()
}

// View renders the confirmation screen using StandardView.
func (s *ConfirmScreen) View() string {
	// Use Screen's CreateView helper for StandardView integration
	return s.CreateView(s.breadcrumbs, s.RenderContent(), s.footer)
}

// SetFooter updates the footer help text.
func (s *ConfirmScreen) SetFooter(footer string) {
	s.footer = footer
}

// GetTitle returns the confirmation title.
func (s *ConfirmScreen) GetTitle() string {
	return s.title
}

// GetMessage returns the confirmation message.
func (s *ConfirmScreen) GetMessage() string {
	return s.message
}

// GetSelection returns the current selection (true = Yes, false = No).
func (s *ConfirmScreen) GetSelection() bool {
	return s.selectedYes
}

// SetSelection sets the current selection (true = Yes, false = No).
func (s *ConfirmScreen) SetSelection(yes bool) {
	s.selectedYes = yes
}

// SetYesText customizes the Yes button text.
//
// Example:
//
//	screen.SetYesText("Delete")   // "Delete" instead of "Yes"
//	screen.SetYesText("Proceed")  // "Proceed" instead of "Yes"
func (s *ConfirmScreen) SetYesText(text string) {
	s.yesText = text
}

// SetNoText customizes the No button text.
//
// Example:
//
//	screen.SetNoText("Cancel")  // "Cancel" instead of "No"
//	screen.SetNoText("Keep")    // "Keep" instead of "No"
func (s *ConfirmScreen) SetNoText(text string) {
	s.noText = text
}
