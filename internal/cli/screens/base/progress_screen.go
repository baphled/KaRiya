package base

import (
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/screens"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TickMsg is sent periodically to animate the spinner.
type TickMsg struct{}

// CompleteMsg is sent when the async operation completes successfully.
type CompleteMsg struct {
	Data interface{}
}

// ErrorMsg is sent when the async operation fails.
type ErrorMsg struct {
	Error string
}

// TickCmd returns a command that sends TickMsg after a delay.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func TickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
		return TickMsg{}
	})
}

// ProgressScreen provides a reusable progress/loading screen.
//
// This screen handles:
// - Animated spinner for visual feedback
// - Optional cancellation (escape key)
// - Async operation completion (CompleteMsg)
// - Error handling (ErrorMsg)
// - Progress message updates
// - StandardView integration
//
// Example usage:
//
//	// Create a progress view.
//	screen := base.NewBaseProgressScreen(
//	    []string{"Main Menu", "Generate CV"},
//	    "Generating CV",
//	    "Analyzing career events and extracting insights...",
//	)
//	screen.SetAllowCancel(false) // Disable cancellation during critical operation
//
//	// In intent, start async operation:
//	cmd := func() tea.Msg {
//	    result, err := performLongOperation()
//	    if err != nil {
//	        return base.ErrorMsg{Error: err.Error()}
//	    }
//	    return base.CompleteMsg{Data: result}
//	}
//
//	// Update progress message mid-operation (optional):
//	screen.SetMessage("Processing facts...")
//
//	// Handle completion:
//	if result.Type() == screens.ResultNavigate {
//	    data := result.Data() // Access completed data
//	}
//
// Related:
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts)
// - internal/cli/screens/cv/generating.go (Example progress screen).
type ProgressScreen struct {
	*Screen

	// breadcrumbs for navigation context
	breadcrumbs []string

	// title is the operation title
	title string

	// message is the progress message
	message string

	// spinnerFrame tracks the current spinner animation frame
	spinnerFrame int

	// allowCancel determines whether escape key cancels the operation
	allowCancel bool

	// footer is the help text shown at the bottom
	footer string

	// spinnerChars are the characters used for the spinner animation
	spinnerChars []string
}

// NewBaseProgressScreen creates a new progress screen.
//
// Parameters:
//   - breadcrumbs: Navigation breadcrumb trail
//   - title: The operation title (e.g., "Generating CV")
//   - message: The progress message (e.g., "Analyzing career events...")
//
// Expected:
//   - breadcrumbs must be a valid slice of strings.
//   - title must be a valid string.
//   - message must be a valid string.
//
// Returns:
//   - A fully initialized ProgressScreen ready for use.
//
// Side effects:
//   - None.
func NewBaseProgressScreen(
	breadcrumbs []string,
	title, message string,
) *ProgressScreen {
	return &ProgressScreen{
		Screen:       NewBaseScreen(),
		breadcrumbs:  breadcrumbs,
		title:        title,
		message:      message,
		spinnerFrame: 0,
		allowCancel:  true,
		footer:       "Esc: Cancel (operation will complete in background)",
		spinnerChars: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
}

// Update handles messages and returns result when operation completes or is cancelled.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Cmd: command to execute.
//   - screens.ScreenResult: result indicating completion state.
//
// Side effects:
//   - May advance spinner frame.
//   - May return NavigateResult on success.
//   - May return ErrorResult on failure.
//   - May return CancelResult on cancellation.
func (s *ProgressScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle window resize
		s.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil

	case TickMsg:
		// Advance spinner animation
		s.spinnerFrame++
		return TickCmd(), nil

	case CompleteMsg:
		// Operation completed successfully
		return nil, &screens.NavigateResult{
			ResultData: msg.Data,
		}

	case ErrorMsg:
		// Operation failed
		return nil, &screens.ErrorResult{
			Message: msg.Error,
		}

	case tea.KeyMsg:
		if s.allowCancel && msg.String() == "esc" {
			// Cancel operation
			return nil, &screens.CancelResult{}
		}
		// Ignore other keys during progress
	}

	return nil, nil
}

// Init initializes the screen and starts the spinner animation.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (s *ProgressScreen) Init() tea.Cmd {
	return TickCmd()
}

// RenderContent returns the progress content without StandardView wrapper.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *ProgressScreen) RenderContent() string {
	var b strings.Builder

	// Title (bold and colored)
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	b.WriteString(titleStyle.Render(s.title))
	b.WriteString("\n\n")

	// Spinner + Message
	spinner := s.spinnerChars[s.spinnerFrame%len(s.spinnerChars)]
	spinnerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	b.WriteString(spinnerStyle.Render(spinner + " "))
	b.WriteString(s.message)
	b.WriteString("\n\n")

	// Additional info
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	b.WriteString(infoStyle.Render("This may take a few moments..."))
	b.WriteString("\n")

	return b.String()
}

// View renders the progress screen using StandardView.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *ProgressScreen) View() string {
	// Update footer based on cancellation setting
	footer := s.footer
	if !s.allowCancel {
		footer = "Please wait... (cancellation disabled)"
	}

	// Use Screen's CreateView helper for StandardView integration
	return s.CreateView(s.breadcrumbs, s.RenderContent(), footer)
}

// SetFooter updates the footer help text.
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (s *ProgressScreen) SetFooter(footer string) {
	s.footer = footer
}

// GetTitle returns the operation title.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *ProgressScreen) GetTitle() string {
	return s.title
}

// SetTitle updates the operation title.
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (s *ProgressScreen) SetTitle(title string) {
	s.title = title
}

// GetMessage returns the progress message.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *ProgressScreen) GetMessage() string {
	return s.message
}

// SetMessage updates the progress message.
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (s *ProgressScreen) SetMessage(message string) {
	s.message = message
}

// GetSpinnerFrame returns the current spinner frame.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (s *ProgressScreen) GetSpinnerFrame() int {
	return s.spinnerFrame
}

// SetSpinnerFrame sets the spinner frame.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (s *ProgressScreen) SetSpinnerFrame(frame int) {
	s.spinnerFrame = frame
}

// SetAllowCancel sets whether the escape key cancels the operation.
//
// Expected:
//   - bool must be valid.
//
// Side effects:
//   - None.
func (s *ProgressScreen) SetAllowCancel(allow bool) {
	s.allowCancel = allow
}
