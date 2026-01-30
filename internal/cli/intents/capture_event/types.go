package capture_event

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
)

// Ensure Intent implements ScreenResultHandler interface.
var _ behaviors.ScreenResultHandler = (*Intent)(nil)

// Intent implements the Intent interface for capturing career events.
// It owns the complete lifecycle of event capture, including:
// - Choosing capture strategy (manual, quick, enriched)
// - Capturing event details via form
// - Reviewing and refining inferred bursts and facts
// - Submitting the final event.
type Intent struct {
	*intents.BaseIntent

	context      *IntentContext
	state        *Model
	active       bool
	result       *intents.IntentResult[*Result]
	eventService *service.CLIEventService

	activeScreen screens.Screen
	useScreens   bool
}

// Model represents the local state of the CaptureEvent intent.
// This is the ONLY mutable state owned by the intent.
type Model struct {
	// context is the input context passed to the intent.
	context *IntentContext

	// currentState tracks which sub-flow is active.
	currentState State

	// captureForm is the form for capturing event details.
	captureForm *models.CaptureForm

	// reviewState is the state of the ReviewInferredEvent sub-flow.
	reviewState *ReviewInferredEventState

	// result is the final result to be returned.
	result *Result

	// error tracks any errors during the intent.
	error *intents.IntentError

	// selectedStrategyIndex tracks the selected strategy in the choose strategy
	// screen (0=Quick, 1=Manual).
	selectedStrategyIndex int

	// strategy is the chosen capture strategy.
	strategy CaptureStrategy

	// submitModal is shown during async submission (modal overlay pattern).
	submitModal *feedback.Modal

	// postSaveReview indicates we're reviewing enriched data after save (not pre-save review).
	postSaveReview bool
}

// ReviewInferredEventState represents the state of the ReviewInferredEvent sub-flow.
type ReviewInferredEventState struct {
	// Event is the career event being reviewed for enrichment.
	Event *career.Event

	// InferredBursts are the activity bursts suggested by enrichment.
	InferredBursts []*career.Burst

	// InferredFacts are the career facts suggested by enrichment.
	InferredFacts []*career.Fact

	// EditingMode indicates which editing sub-flow is active (metadata, bursts, or facts).
	EditingMode EditingMode

	// EditingIndex is the index of the item being edited within bursts or facts.
	EditingIndex int

	// AcceptedBursts holds the bursts the user has accepted during review.
	AcceptedBursts []*career.Burst

	// AcceptedFacts holds the facts the user has accepted during review.
	AcceptedFacts []*career.Fact

	// RejectedItems maps item IDs to rejection reasons for items the user rejected.
	RejectedItems map[string]string

	// metadataModal is the form model for editing event metadata fields.
	metadataModal *models.MetadataEditorModelNew

	// burstModal is the form model for editing burst suggestions.
	burstModal *models.BurstSuggestionModelNew

	// factModal is the form model for editing fact suggestions.
	factModal *models.FactEditorModelNew

	// SelectedItemType tracks which item kind ("burst" or "fact") is currently highlighted.
	SelectedItemType string

	// SelectedIndex is the zero-based index of the currently highlighted item within its type.
	SelectedIndex int
}

// IsActive returns true if this intent is currently active.
func (i *Intent) IsActive() bool {
	return i.active
}

// GetResult returns the final result of the intent (if complete).
// This is called by the intent router to check if the intent has completed.
func (i *Intent) GetResult() *intents.IntentResult[*Result] {
	return i.result
}

// GetState returns the current workflow step of the capture event intent.
func (i *Intent) GetState() string {
	return string(i.state.currentState)
}

// GetForm returns the current form model instance (for test and debug).
func (i *Intent) GetForm() *models.CaptureForm {
	if i == nil || i.state == nil {
		return nil
	}
	return i.state.captureForm
}

// DisableScreens disables the screens architecture for testing legacy flows.
// This is primarily used in tests that need to test legacy message handling.
func (i *Intent) DisableScreens() {
	i.useScreens = false
	i.activeScreen = nil
}

// SetStateForTesting sets the current state for testing.
// This is an exported helper for cross-package tests that need to manipulate internal state.
func (i *Intent) SetStateForTesting(state State) {
	i.state.currentState = state
}

// GetReviewState returns the review state for testing.
func (i *Intent) GetReviewState() *ReviewInferredEventState {
	if i.state == nil {
		return nil
	}
	return i.state.reviewState
}
