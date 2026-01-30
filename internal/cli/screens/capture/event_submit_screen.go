package capture

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/cli/uikit/widgets"
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
// NOTE(TECHNICAL DEBT): This screen is no longer used by the capture_event
// intent, which now handles submission via feedback.Modal overlays directly.
// This screen is retained for backward compatibility and testing but may be
// removed in a future cleanup pass.
//
// This screen:
// - Shows progress indicator while submitting
// - Executes async submission command on Init
// - Returns SubmitResult on success
// - Returns ErrorResult on failure
//
// Keyboard Shortcuts:
// - Esc: (No effect during submission - prevents accidental cancel)
//
// Related:
// - internal/cli/intents/capture_event/intent.go (submitModal handling)
// - tasks/tasks-42-tui-architecture-refactor.md (Phase 1: CaptureEvent Migration).
type EventSubmitScreen struct {
	*base.Screen

	// event being submitted.
	event *career.Event

	// bursts to submit.
	bursts []*career.Burst

	// facts to submit.
	facts []*career.Fact

	// breadcrumbs for the view header.
	breadcrumbs []string

	// submitting indicates submission is in progress.
	submitting bool

	// completed indicates submission finished (success or error).
	completed bool

	// error stores submission error (if any).
	error error

	// simulateError forces an error for testing.
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
		if s.simulateError != nil {
			return SubmitErrorMsg{Err: s.simulateError}
		}

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
// - SubmitCompleteMsg -> returns SubmitResult
// - SubmitErrorMsg -> returns ErrorResult
// - WindowSizeMsg -> updates dimensions
// - Esc -> ignored during submission (prevents accidental cancel).
func (s *EventSubmitScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	if cmd := s.Screen.HandleWindowSizeMsg(msg); cmd != nil {
		return cmd, nil
	}

	switch msg := msg.(type) {
	case SubmitCompleteMsg:
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
		s.submitting = false
		s.completed = true
		s.error = msg.Err

		return nil, &screens.ErrorResult{
			Err:     msg.Err,
			Message: "Event submission failed",
		}

	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			return nil, nil
		}
	}

	return nil, nil
}

// View implements the Screen interface.
//
// Renders the submission progress screen using UIKit components.
func (s *EventSubmitScreen) View() string {
	content := s.renderContent()
	footer := s.renderFooter()

	return s.CreateView(s.breadcrumbs, content, footer)
}

// renderContent renders the submission status using UIKit primitives.
func (s *EventSubmitScreen) renderContent() string {
	th := s.resolveTheme()

	if s.error != nil {
		return s.renderErrorState(th)
	}
	if s.completed {
		return s.renderSuccessState(th)
	}
	return s.renderSubmittingState(th)
}

// renderErrorState renders the error view using UIKit primitives.
func (s *EventSubmitScreen) renderErrorState(th theme.Theme) string {
	var parts []string

	parts = append(parts, primitives.ErrorText("Submission Failed", th).
		Bold().
		MarginBottom(1).
		Render())

	parts = append(parts, primitives.Body(
		fmt.Sprintf("Error: %s", s.error.Error()), th).
		MarginBottom(1).
		Render())

	parts = append(parts, primitives.Muted(
		"Press Esc to go back and try again.", th).
		Render())

	return primitives.JoinVertical(primitives.AlignLeft, parts...)
}

// renderSuccessState renders the success view using UIKit DetailView.
func (s *EventSubmitScreen) renderSuccessState(th theme.Theme) string {
	var parts []string

	parts = append(parts, primitives.SuccessText("Event Submitted Successfully", th).
		Bold().
		MarginBottom(1).
		Render())

	dv := widgets.NewDetailView(th)
	if s.event != nil {
		dv.Field("Event", s.event.Text)
	}
	if len(s.bursts) > 0 {
		dv.Field("Bursts", fmt.Sprintf("%d", len(s.bursts)))
	}
	if len(s.facts) > 0 {
		dv.Field("Facts", fmt.Sprintf("%d", len(s.facts)))
	}
	parts = append(parts, dv.Render())

	return primitives.JoinVertical(primitives.AlignLeft, parts...)
}

// renderSubmittingState renders the progress view using UIKit primitives.
func (s *EventSubmitScreen) renderSubmittingState(th theme.Theme) string {
	var parts []string

	parts = append(parts, primitives.InfoText("Submitting Event...", th).
		Bold().
		MarginBottom(1).
		Render())

	if s.event != nil {
		dv := widgets.NewDetailView(th).
			Field("Event", s.event.Text)
		parts = append(parts, dv.Render())
	}

	parts = append(parts, primitives.Muted("Please wait...", th).
		MarginTop(1).
		Render())

	return primitives.JoinVertical(primitives.AlignLeft, parts...)
}

// renderFooter renders footer with status message using UIKit badge primitives.
func (s *EventSubmitScreen) renderFooter() string {
	th := s.resolveThemesTheme()

	if s.error != nil {
		return primitives.RenderHelpFooter(th,
			primitives.BackBadge(th),
			primitives.QuitBadge(th),
		)
	}
	if s.completed {
		return primitives.RenderHelpFooter(th,
			primitives.ContinueBadge(th),
		)
	}

	return primitives.Muted("Submitting... Please wait", th).Render()
}

// resolveTheme returns the screen's theme or a default if none is set.
func (s *EventSubmitScreen) resolveTheme() theme.Theme {
	if screenTheme := s.Theme(); screenTheme != nil {
		if th, ok := screenTheme.(theme.Theme); ok {
			return th
		}
	}
	return theme.Default()
}

// resolveThemesTheme returns the screen's theme as themes.Theme for badge rendering.
func (s *EventSubmitScreen) resolveThemesTheme() themes.Theme {
	if screenTheme := s.Theme(); screenTheme != nil {
		if th, ok := screenTheme.(themes.Theme); ok {
			return th
		}
	}
	return themes.NewDefaultTheme()
}
