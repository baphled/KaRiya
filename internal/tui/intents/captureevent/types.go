package captureevent

import (
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/tui/intents"
	burstviews "github.com/baphled/kariya/internal/tui/views/burst"
	"github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

// Intent orchestrates the complete lifecycle of event capture.
//
// Views: Strategy selection and event review are now View-based (presentational).
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
	context *IntentValidator

	// active is true while the intent is running.
	active bool

	// result is the typed intent result returned to the router on completion.
	result *intents.IntentResult[*Result]

	// eventService is the CLI event service used for persistence.
	eventService *service.CLIEventService

	// activeView is the currently displayed view component (presentational-only).
	activeView widgets.View

	// currentState tracks which workflow step is active.
	currentState State

	// reviewState holds the review sub-flow's data (event, bursts, facts, modals).
	reviewState *ReviewInferredEventState

	// strategy is the capture strategy chosen by the user.
	strategy CaptureStrategy

	// submitModal is the loading/success/error modal shown during async submission.
	submitModal *feedback.Modal
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

	// InferredBurstSuggestions are the raw burst suggestions from the inference
	// service, preserving fields (e.g. ConfidenceScore) that career.Burst does not carry.
	InferredBurstSuggestions []burstfact.BurstSuggestion

	// InferredFacts are the career facts suggested by enrichment.
	InferredFacts []*career.Fact

	// InferredSkills are the skills suggested by inference.
	InferredSkills []skillinference.SkillSuggestion

	// InferredDisplayFacts caches display facts presented in the fact suggestion modal.
	InferredDisplayFacts []display.Fact

	// EditingMode indicates which editing sub-flow is active (metadata, bursts, or facts).
	EditingMode EditingMode

	// EditingIndex is the index of the item being edited within bursts or facts.
	EditingIndex int

	// AcceptedBursts holds the bursts the user has accepted during review.
	AcceptedBursts []*career.Burst

	// AcceptedFacts holds the facts the user has accepted during review.
	AcceptedFacts []*career.Fact

	// AcceptedSkills holds the skills the user has accepted during review.
	AcceptedSkills []*career.Skill

	// RejectedItems maps item IDs to rejection reasons for items the user rejected.
	RejectedItems map[string]string

	// metadataModal is the form model for editing event metadata fields.
	metadataModal *event.ReviewEnrichment

	// burstModal is the modal for reviewing burst suggestions.
	burstModal *burstviews.SuggestionReview

	// factSuggestionModal is the modal for reviewing fact suggestions.
	factSuggestionModal *burstviews.SuggestionReview

	// skillModal is the modal for editing skill suggestions.
	skillModal *burstviews.SuggestionReview

	// SelectedItemType tracks which item kind ("burst" or "fact") is highlighted.
	SelectedItemType string

	// SelectedIndex is the zero-based position of the highlighted item within its type.
	SelectedIndex int
}

// ReviewEnrichmentDimensions holds terminal dimensions for modal sizing.
// Pass these from the intent so the form sizes correctly inside the overlay.
type ReviewEnrichmentDimensions struct {
	TerminalWidth  int
	TerminalHeight int
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

// GetFactSuggestionModal returns the fact suggestion review modal for test assertions.
//
// Returns:
//   - A fully initialized burstviews.SuggestionReview ready for use.
//
// Side effects:
//   - None.
func (r *ReviewInferredEventState) GetFactSuggestionModal() *burstviews.SuggestionReview {
	if r == nil {
		return nil
	}
	return r.factSuggestionModal
}
