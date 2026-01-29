package capture

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// SubmitCompleteMsg indicates submission completed successfully.
type SubmitCompleteMsg struct {
	Event  *career.Event
	Bursts []*career.Burst
	Facts  []*career.Fact
}

// SubmitErrorMsg indicates submission failed.
type SubmitErrorMsg struct {
	Err error
}

// EventSubmitScreen handles async submission of captured event to database.
//
// NOTE(TECHNICAL DEBT): This should be a modal overlay, not a full screen.
// Currently this renders as a full screen replacement, but the correct
// architecture is to:
// 1. Keep the form screen visible in the background
// 2. Show a loading modal overlay during submission
// 3. Show success/error modal when complete
// 4. Return to previous screen (or complete intent) on modal dismiss
//
// See BrowseTimeline intent for the correct modal overlay pattern.
// This refactor requires:
// - Moving submission logic to intent level
// - Using components.NewLoadingModal() during submission
// - Using components.NewSuccessModal() / NewErrorModal() for results
// - Removing this full screen entirely
//
// Priority: MEDIUM (works but not ideal UX)
// Effort: ~2 hours
//
// This screen:
// - Shows progress indicator while submitting
// - Executes async submission command on Init
// - Returns SubmitResult on success
// - Returns ErrorResult on failure
//
// The actual submission is handled by the intent (via service),
// this screen just coordinates the async flow and displays progress.
//
// Keyboard Shortcuts:
// - Esc: (No effect during submission - prevents accidental cancel)
//
// Related:
// - internal/cli/intents/capture_event.go (performSubmit)
// - tasks/tasks-42-tui-architecture-refactor.md (Phase 1: CaptureEvent Migration).
type EventSubmitScreen struct {
	*base.Screen

	// event being submitted
	event *career.Event

	// bursts to submit
	bursts []*career.Burst

	// facts to submit
	facts []*career.Fact

	// breadcrumbs for the view header
	breadcrumbs []string

	// submitting indicates submission is in progress
	submitting bool

	// completed indicates submission finished (success or error)
	completed bool

	// error stores submission error (if any)
	error error

	// simulateError forces an error for testing
	simulateError error
}

// NewEventSubmitScreen creates a new EventSubmitScreen.
//
// Parameters:
//   - breadcrumbs: Breadcrumb trail for header (e.g., ["Main Menu", "Capture Event", "Submit"])
//   - event: Event to submit
//   - bursts: Bursts to submit
//   - facts: Facts to submit
//
// Returns a EventSubmitScreen that will submit on Init().
func NewEventSubmitScreen(
	breadcrumbs []string,
	event *career.Event,
	bursts []*career.Burst,
	facts []*career.Fact,
) *EventSubmitScreen {
	return &EventSubmitScreen{
		Screen:      base.NewBaseScreen(),
		event:       event,
		bursts:      bursts,
		facts:       facts,
		breadcrumbs: breadcrumbs,
		submitting:  false,
		completed:   false,
	}
}

// NewEventSubmitScreenWithError creates a screen that simulates submission error.
//
// This is used for testing error handling.
func NewEventSubmitScreenWithError(
	breadcrumbs []string,
	event *career.Event,
	bursts []*career.Burst,
	facts []*career.Fact,
	err error,
) *EventSubmitScreen {
	screen := NewEventSubmitScreen(breadcrumbs, event, bursts, facts)
	screen.simulateError = err
	return screen
}

// Init implements the Screen interface.
//
// Returns a command that triggers async submission.
func (s *EventSubmitScreen) Init() tea.Cmd {
	s.submitting = true
	return s.performSubmit()
}

// performSubmit simulates async submission.
//
// In the real implementation, this would call the service layer.
// For screens architecture, the intent will provide the actual submission command.
func (s *EventSubmitScreen) performSubmit() tea.Cmd {
	return func() tea.Msg {
		// Simulate submission error if configured
		if s.simulateError != nil {
			return SubmitErrorMsg{Err: s.simulateError}
		}

		// Simulate successful submission
		return SubmitCompleteMsg{
			Event:  s.event,
			Bursts: s.bursts,
			Facts:  s.facts,
		}
	}
}

// Update implements the Screen interface.
//
// Handles:
// - SubmitCompleteMsg → returns SubmitResult
// - SubmitErrorMsg → returns ErrorResult
// - WindowSizeMsg → updates dimensions
// - Esc → ignored during submission (prevents accidental cancel).
func (s *EventSubmitScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	// Handle window size via Screen
	if cmd := s.Screen.HandleWindowSizeMsg(msg); cmd != nil {
		return cmd, nil
	}

	// Handle submission messages
	switch msg := msg.(type) {
	case SubmitCompleteMsg:
		// Submission succeeded
		s.submitting = false
		s.completed = true

		return nil, &screens.SubmitResult{
			FormData: map[string]interface{}{
				"event":  msg.Event,
				"bursts": msg.Bursts,
				"facts":  msg.Facts,
			},
		}

	case SubmitErrorMsg:
		// Submission failed
		s.submitting = false
		s.completed = true
		s.error = msg.Err

		return nil, &screens.ErrorResult{
			Err:     msg.Err,
			Message: "Event submission failed",
		}

	case tea.KeyMsg:
		// Ignore Esc during submission to prevent accidental cancel
		// User must wait for submission to complete (or fail)
		if msg.String() == "esc" {
			return nil, nil
		}
	}

	return nil, nil
}

// View implements the Screen interface.
//
// Renders the submission progress screen.
func (s *EventSubmitScreen) View() string {
	content := s.renderContent()
	footer := s.renderFooter()

	return s.CreateView(s.breadcrumbs, content, footer)
}

// renderContent renders the submission status.
func (s *EventSubmitScreen) renderContent() string {
	var b strings.Builder

	b.WriteString("\n")

	if s.error != nil {
		// Show error state
		b.WriteString("❌ Submission Failed\n")
		b.WriteString("═══════════════════\n\n")
		b.WriteString(fmt.Sprintf("Error: %s\n", s.error.Error()))
		b.WriteString("\nPress Esc to go back and try again.\n")
	} else if s.completed {
		// Show success state
		b.WriteString("✅ Event Submitted Successfully\n")
		b.WriteString("══════════════════════════════\n\n")
		if s.event != nil {
			b.WriteString(fmt.Sprintf("Event: %s\n", s.event.Text))
		}
		if len(s.bursts) > 0 {
			b.WriteString(fmt.Sprintf("Bursts: %d\n", len(s.bursts)))
		}
		if len(s.facts) > 0 {
			b.WriteString(fmt.Sprintf("Facts: %d\n", len(s.facts)))
		}
	} else if s.submitting {
		// Show progress state
		b.WriteString("⏳ Submitting Event...\n")
		b.WriteString("═══════════════════\n\n")
		if s.event != nil {
			b.WriteString(fmt.Sprintf("Event: %s\n", s.event.Text))
		}
		b.WriteString("\nPlease wait...\n")
	}

	return b.String()
}

// renderFooter renders footer with status message.
func (s *EventSubmitScreen) renderFooter() string {
	if s.error != nil {
		return "Esc: Back  q: Quit"
	} else if s.completed {
		return "Submission complete"
	}
	return "Submitting... Please wait"
}
