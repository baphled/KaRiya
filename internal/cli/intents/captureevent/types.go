package captureevent

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
)

// Ensure Intent implements ScreenResultHandler interface.
var _ behaviors.ScreenResultHandler = (*Intent)(nil)

// Intent orchestrates the complete lifecycle of event capture.
//
// Responsibilities:
//   - Choosing capture strategy (quick or manual)
//   - Presenting the event capture form
//   - Managing the review phase for inferred bursts and facts
//   - Coordinating event persistence and enrichment
//
// This struct is the broker only; rendering is delegated to screens and
// business logic is delegated to the career service.
type Intent struct {
	*intents.BaseIntent

	// context holds the input parameters passed to the intent at creation.
	context *IntentContext

	// active is true while the intent is running.
	active bool

	// result is the typed intent result returned to the router on completion.
	result *intents.IntentResult[*Result]

	// eventService is the CLI event service used for persistence.
	eventService *service.CLIEventService

	// activeScreen is the currently displayed screen component.
	activeScreen screens.Screen

	// currentState tracks which workflow step is active.
	currentState State

	// captureFormScreen is the screen-level form for capturing event details.
	captureFormScreen *capture.EventFormScreen

	// reviewState holds the review sub-flow's data (event, bursts, facts, modals).
	reviewState *ReviewInferredEventState

	// strategy is the capture strategy chosen by the user.
	strategy CaptureStrategy

	// submitModal is the loading/success/error modal shown during async submission.
	submitModal *feedback.Modal

	// postSaveReview is true when reviewing enriched data after a successful save,
	// as opposed to the pre-save review.
	postSaveReview bool
}

// ReviewInferredEventState holds all data for the review sub-flow where the user
// inspects and edits inferred bursts and facts.
//
// This struct owns the editing modals and selection cursor used during review.
type ReviewInferredEventState struct {
	// Event is the career event being reviewed for enrichment.
	Event *career.Event

	// InferredBursts are the activity bursts suggested by enrichment.
	InferredBursts []*career.Burst

	// InferredFacts are the career facts suggested by enrichment.
	InferredFacts []*career.Fact

	// InferredSkills are the skills suggested by inference.
	InferredSkills []skillinference.SkillSuggestion

	// EditingMode indicates which editing sub-flow is active (metadata, bursts, or facts).
	EditingMode EditingMode

	// EditingIndex is the index of the item being edited within bursts or facts.
	EditingIndex int

	// AcceptedBursts holds the bursts the user has accepted during review.
	AcceptedBursts []*career.Burst

	// AcceptedFacts holds the facts the user has accepted during review.
	AcceptedFacts []*career.Fact

	// AcceptedSkills holds the skills the user has accepted during review.
	AcceptedSkills []skillinference.SkillSuggestion

	// RejectedItems maps item IDs to rejection reasons for items the user rejected.
	RejectedItems map[string]string

	// metadataModal is the form model for editing event metadata fields.
	metadataModal *MetadataEditorModelNew

	// burstModal is the form model for editing burst suggestions.
	burstModal *BurstSuggestionModelNew

	// factModal is the form model for editing fact suggestions.
	factModal *FactEditorModelNew

	// SelectedItemType tracks which item kind ("burst" or "fact") is highlighted.
	SelectedItemType string

	// SelectedIndex is the zero-based position of the highlighted item within its type.
	SelectedIndex int
}

// IsActive reports whether this intent is currently running.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (i *Intent) IsActive() bool {
	return i.active
}

// GetResult returns the typed intent result, or nil if not yet complete.
//
// Returns:
//   - A fully initialized intents.IntentResult[*Result] ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetResult() *intents.IntentResult[*Result] {
	return i.result
}

// GetState returns the current workflow step as a string.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (i *Intent) GetState() string {
	return string(i.currentState)
}

// GetFormScreen returns the capture form screen instance.
//
// Returns:
//   - A fully initialized capture.EventFormScreen ready for use, or nil.
//
// Side effects:
//   - None.
func (i *Intent) GetFormScreen() *capture.EventFormScreen {
	if i == nil {
		return nil
	}
	return i.captureFormScreen
}

// SetStateForTesting sets the current workflow state for cross-package tests.
//
// Expected:
//   - state must be valid.
//
// Side effects:
//   - None.
func (i *Intent) SetStateForTesting(state State) {
	i.currentState = state
}

// GetReviewState returns the review sub-flow state for test assertions.
//
// Returns:
//   - A fully initialized ReviewInferredEventState ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetReviewState() *ReviewInferredEventState {
	return i.reviewState
}
