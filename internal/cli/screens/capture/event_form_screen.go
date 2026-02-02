// Package capture provides screen components for the capture event workflow.
package capture

import (
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/types"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
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
// - tasks/tasks-42-tui-architecture-refactor.md (Phase 1: CaptureEvent Migration).
type EventFormScreen struct {
	*base.Screen

	// captureForm is the underlying form model
	captureForm *models.CaptureForm

	// breadcrumbs for the view header
	breadcrumbs []string

	// strategy is the capture strategy (quick or manual)
	strategy types.CaptureStrategy
}

// NewEventFormScreen creates a new EventFormScreen with the specified strategy.
//
// Parameters:
//   - cliService: CLI service for event operations (used by form for submission)
//   - breadcrumbs: Breadcrumb trail for header (e.g., ["Main Menu", "Capture Event", "Form"])
//   - strategy: Capture strategy (StrategyQuick or StrategyManual)
//
// Expected:
//   - cliservice must be a valid *service.CLIEventService.
//   - breadcrumbs must be a valid slice of strings.
//   - strategy must be a valid CaptureStrategy.
//
// Returns:
//   - A fully initialized EventFormScreen ready for use.
//
// Side effects:
//   - Initializes CaptureForm.
func NewEventFormScreen(
	cliService *service.CLIEventService,
	breadcrumbs []string,
	strategy types.CaptureStrategy,
) *EventFormScreen {
	captureForm := models.NewCaptureForm(cliService)
	captureForm.SetStrategy(string(strategy))

	return &EventFormScreen{
		Screen:      base.NewBaseScreen(),
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
// Expected:
//   - cliservice must be a valid *service.CLIEventService.
//   - breadcrumbs must be a valid slice of strings.
//   - strategy must be a valid CaptureStrategy.
//   - event must be a valid *career.Event (can be nil).
//
// Returns:
//   - A fully initialized EventFormScreen ready for use with form pre-populated with event data.
//
// Side effects:
//   - Pre-populates form with event data if event is provided.
func NewEventFormScreenWithEvent(
	cliService *service.CLIEventService,
	breadcrumbs []string,
	strategy types.CaptureStrategy,
	event *career.Event,
) *EventFormScreen {
	screen := NewEventFormScreen(cliService, breadcrumbs, strategy)

	if event != nil {
		screen.captureForm.LoadEventForEditing(event)
	}

	return screen
}

// Init implements the Screen interface.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (s *EventFormScreen) Init() tea.Cmd {
	return s.captureForm.Init()
}

// Update implements the Screen interface.
//
// Handles:
// - Escape key → returns CancelResult
// - WindowSizeMsg → updates form dimensions
// - SubmitMsg → returns SubmitResult with event data
// - Other messages → delegates to CaptureForm.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - A tea.Cmd value.
//   - A screens.ScreenResult value.
//
// Side effects:
//   - May return CancelResult on escape.
//   - May return SubmitResult on form completion.
//   - May return ErrorResult on submission error.
func (s *EventFormScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	// Handle window size via Screen
	if cmd := s.HandleWindowSizeMsg(msg); cmd != nil {
		// Also update CaptureForm's dimensions
		if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
			s.captureForm.Update(wsMsg)
		}
		return cmd, nil
	}

	// Handle key messages.
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
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
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *EventFormScreen) View() string {
	// Get form content
	content := s.captureForm.View()

	// Build footer with shortcuts
	footer := s.renderFooter()

	// Create view with StandardView
	return s.CreateView(s.breadcrumbs, content, footer)
}

// renderFooter renders footer with form shortcuts using UIKit badge primitives.
func (s *EventFormScreen) renderFooter() string {
	th := s.resolveThemesTheme()

	return primitives.RenderHelpFooter(th,
		primitives.NextFieldBadge(th),
		primitives.HelpKeyBadge("Ctrl+S", "Submit", th),
		primitives.CancelBadge(th),
		primitives.QuitBadge(th),
	)
}

// resolveThemesTheme returns the screen's theme as themes.Theme for badge rendering.
func (s *EventFormScreen) resolveThemesTheme() themes.Theme {
	if screenTheme := s.Theme(); screenTheme != nil {
		if th, ok := screenTheme.(themes.Theme); ok {
			return th
		}
	}
	return themes.NewDefaultTheme()
}
