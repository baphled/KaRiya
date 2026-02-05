// Package capture provides screen components for the capture event workflow.
package capture

import (
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/types"
	"github.com/baphled/kariya/internal/domain/career"
)

// EventFormState identifies the event capture form view in the state
// matrix. On this screen the user sees an interactive form with fields for
// event details. The form has two strategies: quick (minimal fields) and
// manual (all fields). Tab advances between fields, Enter submits the
// completed form, and Escape cancels without saving.
const EventFormState = "event_form"

// EventFormScreen provides a form for capturing career events.
//
// This screen wraps FormScreen[*forms.CaptureEventFormData] with capture-specific context:
// - Supports two strategies: quick (text + date) and manual (all fields)
// - Pre-populates form with existing event data when editing
// - Uses CaptureEventFormData with validation from forms package
// - Handles terminal resize by rebuilding form
// - Returns SubmitResult with form data on submission
// - Returns CancelResult on escape
//
// Usage (New Event - Quick Strategy):
//
//	screen := capture.NewEventFormScreen(nil, []string{"Main", "Capture"}, types.StrategyQuick)
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultSubmit {
//	    submitResult := result.(*screens.SubmitResult)
//	    data := submitResult.FormData.(*forms.CaptureEventFormData)
//	    if data.SubmitConfirmed {
//	        // Convert form data to event and save
//	    }
//	}
//
// Usage (Edit Event - Manual Strategy):
//
//	screen := capture.NewEventFormScreen(existingEvent, []string{"Main", "Edit"}, types.StrategyManual)
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultSubmit {
//	    submitResult := result.(*screens.SubmitResult)
//	    data := submitResult.FormData.(*forms.CaptureEventFormData)
//	    if data.SubmitConfirmed {
//	        // Convert form data to event and save
//	    }
//	}
//
// Related:
// - FormScreen provides the form UI
// - forms.CaptureEventFormData defines form structure
// - forms.NewCaptureEventForm creates the huh form
// - docs/FORMS_GUIDE.md (Form patterns and best practices).
type EventFormScreen struct {
	*base.FormScreen[*forms.CaptureEventFormData]

	// strategy is the capture strategy (quick or manual)
	strategy types.CaptureStrategy
}

// NewEventFormScreen creates a new event form screen with the specified strategy.
//
// Parameters:
//   - event: Existing event to edit (nil creates a new event)
//   - breadcrumbs: Breadcrumb trail for header (e.g., ["Main Menu", "Capture Event"])
//   - strategy: Capture strategy (StrategyQuick or StrategyManual)
//
// Expected:
//   - event can be nil (creates empty form for new capture).
//   - breadcrumbs must be a valid slice of strings.
//   - strategy must be a valid CaptureStrategy.
//
// Returns:
//   - A fully initialized EventFormScreen ready for use.
//
// Side effects:
//   - Creates form with appropriate fields based on strategy.
func NewEventFormScreen(
	event *career.Event,
	breadcrumbs []string,
	strategy types.CaptureStrategy,
) *EventFormScreen {
	var formData *forms.CaptureEventFormData
	if event == nil {
		formData = forms.NewCaptureEventFormData()
	} else {
		formData = forms.GetCaptureEventFormData(event)
	}

	strategyStr := string(strategy)

	builder := func(data *forms.CaptureEventFormData, w, h int) forms.Form {
		return forms.NewCaptureEventForm(data, strategyStr, w, h)
	}

	baseScreen := base.NewBaseFormScreen(breadcrumbs, builder, formData)

	return &EventFormScreen{
		FormScreen: baseScreen,
		strategy:   strategy,
	}
}

// GetStrategy returns the current capture strategy.
//
// Returns:
//   - A types.CaptureStrategy value.
//
// Side effects:
//   - None.
func (s *EventFormScreen) GetStrategy() types.CaptureStrategy {
	return s.strategy
}
