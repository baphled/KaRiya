// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"context"
	"errors"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/tui/intents"
	burstviews "github.com/baphled/kariya/internal/tui/views/burst"
	skillviews "github.com/baphled/kariya/internal/tui/views/skill"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

var (
	// ErrInvalidContext is returned when the context is invalid.
	ErrInvalidContext = errors.New("invalid context")
)

// Intent implements the Intent interface for managing career bursts.
// It owns the complete lifecycle of burst management, including:
// - Displaying a list of bursts
// - Viewing burst details, events, and facts
// - Creating, editing, and deleting bursts
// - Confirming bursts and extracting facts
// - Suggesting bursts from events using AI detection.
type Intent struct {
	// Embed BaseIntent for terminal awareness, logo, and state management.
	*intents.BaseIntent

	// context is the input context passed to the intent.
	context *IntentValidator

	// state tracks the current state of the intent (typed enum).
	state State

	// active indicates whether this intent is currently active.
	active bool

	// result is the final result of the intent (set when complete).
	result *intents.IntentResult[*Result]

	// --- Flattened state fields ---

	// filteredBursts are the bursts after applying current filters.
	filteredBursts []*career.Burst

	// selectedIndex is the index of the currently selected burst.
	selectedIndex int

	// selectedBurst is the burst currently being viewed.
	selectedBurst *career.Burst

	// viewedBursts tracks bursts viewed during the session.
	viewedBursts []*career.Burst

	// burstEvents are the events for the current burst.
	burstEvents []*career.Event

	// burstFacts are the facts for the current burst.
	burstFacts []*career.Fact

	// burstSuggestions cache the current domain burst suggestions for reverse lookup.
	burstSuggestions []burstfact.BurstSuggestion

	// skillSuggestions cache the current domain skill suggestions for reverse lookup.
	skillSuggestions []skillinference.SkillSuggestion

	// --- Loading states ---

	// loadingEvents indicates if events are being loaded.
	loadingEvents bool

	// loadingFacts indicates if facts are being loaded.
	loadingFacts bool

	// extractingFacts indicates if facts are being extracted.
	extractingFacts bool

	// suggestionsLoading indicates if burst suggestions are being loaded.
	suggestionsLoading bool

	// inferringSkills indicates if skills are being inferred.
	inferringSkills bool

	// --- Error states ---

	// deleteError stores any error from delete operation.
	deleteError error

	// editError stores any error from edit operation.
	editError error

	// confirmError stores any error from confirm operation.
	confirmError error

	// suggestionsError stores any error from suggestion detection.
	suggestionsError error

	// skillInferenceError stores any error from skill inference.
	skillInferenceError error

	// --- Progress tracking ---

	// extractedFactsCount tracks the number of facts extracted in current operation.
	extractedFactsCount int

	// --- View Orchestration ---

	// activeView holds the current view being displayed (burst List only).
	activeView widgets.View

	// detailModal holds the burst detail modal (replaces detail screen).
	detailModal *burstviews.Detail

	// eventsModal holds the burst events modal (replaces events screen).
	eventsModal *burstviews.Events

	// factsModal holds the burst facts modal (replaces facts screen).
	factsModal *burstviews.Facts

	// skillsModal holds the burst skills modal for viewing skills associated with a burst.
	skillsModal *burstviews.Skills

	// editModal holds the edit burst modal.
	editModal *burstviews.Edit

	// deleteModal holds the delete confirmation modal.
	deleteModal *feedback.ConfirmModal

	// confirmModal holds the burst confirmation modal.
	confirmModal *feedback.ConfirmModal

	// suggestionModal holds the suggestion review modal.
	suggestionModal *burstviews.SuggestionReview

	// skillSuggestionModal holds the skill suggestion review modal.
	skillSuggestionModal *burstviews.SuggestionReview

	// suggestionEventsModal holds the events sub-modal for drill-down from skill suggestions.
	suggestionEventsModal *skillviews.Events

	// inferredFromDetail tracks whether skill inference was triggered from burst detail modal.
	inferredFromDetail bool

	// feedbackModal holds the feedback modal (shown for errors, warnings, and success messages).
	feedbackModal *feedback.Modal

	// loadingModal holds the loading modal (shown during async operations).
	loadingModal *feedback.Modal

	// modalRegistry manages all modals with unified Update/View handling.
	modalRegistry *intents.ModalRegistry

	// cancelFunc is the cancel function for the current async operation.
	// Call this to cancel loading operations when user presses Esc.
	cancelFunc context.CancelFunc
}

// GetState reports the intent's current lifecycle phase.
//
// Returns:
//   - A State value.
//
// Side effects:
//   - None.
func (i *Intent) GetState() State {
	return i.state
}

// SetState transitions the intent to a new lifecycle phase.
//
// Expected:
//   - state must be valid.
//
// Side effects:
//   - None.
func (i *Intent) SetState(state State) {
	i.state = state
}

// IsActive indicates whether the intent is still processing user interactions.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (i *Intent) IsActive() bool {
	return i.active
}

// Deactivate marks the intent as inactive.
//
// Side effects:
//   - None.
func (i *Intent) Deactivate() {
	i.active = false
}

// GetFilteredBursts provides the burst list after all active filters have been applied.
//
// Returns:
//   - A []*career.Burst value.
//
// Side effects:
//   - None.
func (i *Intent) GetFilteredBursts() []*career.Burst {
	return i.filteredBursts
}

// GetSelectedBurst provides the burst at the current cursor position.
//
// Returns:
//   - A fully initialized career.Burst ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedBurst() *career.Burst {
	return i.selectedBurst
}

// SetSelectedBurst updates the cursor to point to a specific burst.
//
// Expected:
//   - burst must be valid.
//
// Side effects:
//   - None.
func (i *Intent) SetSelectedBurst(burst *career.Burst) {
	i.selectedBurst = burst
}

// GetSelectedIndex reports the cursor position in the filtered burst list.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedIndex() int {
	return i.selectedIndex
}

// SetSelectedIndex moves the cursor to a specific position in the filtered list.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (i *Intent) SetSelectedIndex(index int) {
	i.selectedIndex = index
}

// GetViewedBursts tracks which bursts the user has already inspected during this session.
//
// Returns:
//   - A []*career.Burst value.
//
// Side effects:
//   - None.
func (i *Intent) GetViewedBursts() []*career.Burst {
	return i.viewedBursts
}

// AddViewedBurst marks a burst as viewed so the UI can distinguish visited items.
//
// Expected:
//   - burst must be valid.
//
// Side effects:
//   - None.
func (i *Intent) AddViewedBurst(burst *career.Burst) {
	i.viewedBursts = append(i.viewedBursts, burst)
}

// Result provides the intent's outcome for the router to process.
//
// Returns:
//   - A fully initialized intents.IntentResult[interface{}] ready for use.
//
// Side effects:
//   - None.
func (i *Intent) Result() *intents.IntentResult[interface{}] {
	if i.result == nil {
		return nil
	}

	return &intents.IntentResult[interface{}]{
		Status:   i.result.Status,
		Data:     i.result.Data,
		Error:    i.result.Error,
		Metadata: i.result.Metadata,
	}
}

// SetCompleted marks the intent as completed with the selected burst.
//
// Expected:
//   - burst must be valid.
//
// Side effects:
//   - None.
func (i *Intent) SetCompleted(burst *career.Burst) {
	i.result = &intents.IntentResult[*Result]{
		Status: intents.Completed,
		Data: &Result{
			Action:        "selected",
			Burst:         burst,
			Bursts:        i.filteredBursts,
			ViewedBursts:  i.viewedBursts,
			SelectedIndex: i.selectedIndex,
		},
	}
	i.active = false
}

// SetCancelled marks the intent as cancelled.
//
// Side effects:
//   - None.
func (i *Intent) SetCancelled() {
	i.result = &intents.IntentResult[*Result]{
		Status: intents.Cancelled,
	}
	i.active = false
}

// GetModalRegistry provides access to the modal manager for registration and lookup.
//
// Returns:
//   - A fully initialized intents.ModalRegistry ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetModalRegistry() *intents.ModalRegistry {
	return i.modalRegistry
}

// HasActiveModal checks if any modal overlay is currently visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (i *Intent) HasActiveModal() bool {
	// Rebuild registry to ensure it's current with modal state.
	i.rebuildModalRegistry()

	// Use registry to check for visible modals.
	return i.modalRegistry != nil && i.modalRegistry.HasVisibleModal()
}

// GetDetailModal provides the burst detail overlay for rendering and testing.
//
// Returns:
//   - A fully initialized burstviews.Detail ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetDetailModal() *burstviews.Detail {
	return i.detailModal
}

// GetSuggestionModal provides the AI suggestion overlay for rendering and testing.
//
// Returns:
//   - A fully initialized burstviews.SuggestionReview ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetSuggestionModal() *burstviews.SuggestionReview {
	return i.suggestionModal
}

// IsExtractingFacts checks if an AI fact extraction operation is currently running.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (i *Intent) IsExtractingFacts() bool {
	return i.extractingFacts
}

// GetExtractedFactsCount reports how many facts the AI has identified so far.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (i *Intent) GetExtractedFactsCount() int {
	return i.extractedFactsCount
}

// SetLoadingEventsForTesting overrides the loading events flag for tests.
//
// Expected:
//   - bool must be valid.
//
// Side effects:
//   - None.
func (i *Intent) SetLoadingEventsForTesting(loading bool) {
	i.loadingEvents = loading
}

// SetLoadingFactsForTesting overrides the loading facts flag for tests.
//
// Expected:
//   - bool must be valid.
//
// Side effects:
//   - None.
func (i *Intent) SetLoadingFactsForTesting(loading bool) {
	i.loadingFacts = loading
}

// NewIntent creates a new BurstManagement intent.
//
// Expected:
//   - ctx must be non-nil and pass validation.
//
// Returns:
//   - A fully initialized Intent ready for use, or an error if ctx is invalid.
//
// Side effects:
//   - None.
func NewIntent(ctx *IntentValidator) (*Intent, error) {
	if ctx == nil {
		return nil, ErrInvalidContext
	}

	if err := ctx.Validate(); err != nil {
		return nil, err
	}

	intent := &Intent{
		BaseIntent:       intents.NewBaseIntent(),
		context:          ctx,
		state:            StateList,
		active:           true,
		filteredBursts:   ctx.Bursts,
		selectedIndex:    0,
		viewedBursts:     make([]*career.Burst, 0),
		burstEvents:      make([]*career.Event, 0),
		burstFacts:       make([]*career.Fact, 0),
		burstSuggestions: make([]burstfact.BurstSuggestion, 0),
		skillSuggestions: make([]skillinference.SkillSuggestion, 0),
		modalRegistry:    intents.NewModalRegistry(),
	}

	return intent, nil
}
