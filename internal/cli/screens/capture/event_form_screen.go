package capture

import (
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// EventFormScreen wraps CaptureForm and adapts it to the Screen interface.
//
// This screen handles event capture forms with two strategies:
// - Quick: Minimal fields (event text only), date defaults to today
// - Manual: Full form with optional fields (date, company, project, tags)
//
// The screen delegates form rendering and input handling to CaptureForm,
// and translates form messages (SubmitMsg) to ScreenResults (SubmitResult, CancelResult).
//
// Keyboard Shortcuts:
// - Tab: Navigate between fields
// - Enter: Submit form (when on confirm button)
// - Ctrl+S: Submit form (from any field)
// - Esc: Cancel and return to previous screen
//
// Related:
// - internal/cli/models/capture_form.go (CaptureForm model)
// - internal/cli/forms/capture_event_form.go (Form configuration)
// - tasks/tasks-42-tui-architecture-refactor.md (Phase 1: CaptureEvent Migration)
type EventFormScreen struct {
	*base.BaseScreen

	// captureForm is the underlying form model
	captureForm *models.CaptureForm

	// breadcrumbs for the view header
	breadcrumbs []string

	// strategy is the capture strategy (quick or manual)
	strategy intents.CaptureStrategy
}

// NewEventFormScreen creates a new EventFormScreen with the specified strategy.
//
// Parameters:
//   - cliService: CLI service for event operations (used by form for submission)
//   - breadcrumbs: Breadcrumb trail for header (e.g., ["Main Menu", "Capture Event", "Form"])
//   - strategy: Capture strategy (StrategyQuick or StrategyManual)
//
// Returns a EventFormScreen with an initialized CaptureForm.
func NewEventFormScreen(
	cliService *service.CLIEventService,
	breadcrumbs []string,
	strategy intents.CaptureStrategy,
) *EventFormScreen {
	captureForm := models.NewCaptureForm(cliService)
	captureForm.SetStrategy(string(strategy))

	return &EventFormScreen{
		BaseScreen:  base.NewBaseScreen(),
		captureForm: captureForm,
		breadcrumbs: breadcrumbs,
		strategy:    strategy,
	}
}

// NewEventFormScreenWithEvent creates a new EventFormScreen for editing an existing event.
//
// This is used when editing an event from BrowseTimeline or other intents.
//
// Parameters:
//   - cliService: CLI service for event operations
//   - breadcrumbs: Breadcrumb trail for header
//   - strategy: Capture strategy (typically StrategyManual for editing)
//   - event: Existing event to edit (nil creates a new event)
//
// Returns a EventFormScreen with form pre-populated with event data.
func NewEventFormScreenWithEvent(
	cliService *service.CLIEventService,
	breadcrumbs []string,
	strategy intents.CaptureStrategy,
	event *career.CareerEvent,
) *EventFormScreen {
	screen := NewEventFormScreen(cliService, breadcrumbs, strategy)

	if event != nil {
		screen.captureForm.LoadEventForEditing(event)
	}

	return screen
}

// Init implements the Screen interface.
//
// Initializes the underlying CaptureForm.
func (s *EventFormScreen) Init() tea.Cmd {
	return s.captureForm.Init()
}

// Update implements the Screen interface.
//
// Handles:
// - Escape key → returns CancelResult
// - WindowSizeMsg → updates form dimensions
// - SubmitMsg → returns SubmitResult with event data
// - Other messages → delegates to CaptureForm
func (s *EventFormScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	// Handle window size via BaseScreen
	if cmd := s.BaseScreen.HandleWindowSizeMsg(msg); cmd != nil {
		// Also update CaptureForm's dimensions
		if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
			s.captureForm.Update(wsMsg)
		}
		return cmd, nil
	}

	// Handle key messages
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// User pressed Escape - cancel and return to previous screen
			return nil, &screens.CancelResult{}
		}

	case models.SubmitMsg:
		// Form submission completed
		if msg.Err != nil {
			// Form submission failed - return error result
			return nil, &screens.ErrorResult{
				Err:     msg.Err,
				Message: "Form submission failed",
			}
		}

		// Form submission succeeded - return submit result with event data
		return nil, &screens.SubmitResult{
			FormData: msg.Event,
		}
	}

	// Delegate to CaptureForm
	model, cmd := s.captureForm.Update(msg)
	if captureForm, ok := model.(*models.CaptureForm); ok {
		s.captureForm = captureForm
	}

	return cmd, nil
}

// View implements the Screen interface.
//
// Renders the form using StandardView layout.
func (s *EventFormScreen) View() string {
	// Get form content
	content := s.captureForm.View()

	// Build footer with shortcuts
	footer := s.renderFooter()

	// Create view with StandardView
	return s.CreateView(s.breadcrumbs, content, footer)
}

// renderFooter renders footer with form shortcuts.
func (s *EventFormScreen) renderFooter() string {
	return "Tab: Next field  Ctrl+S: Submit  Esc: Cancel  q: Quit"
}
