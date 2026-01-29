// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"context"
	"errors"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	burstmodals "github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
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
	context *IntentContext

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

	// --- Loading states ---

	// loadingEvents indicates if events are being loaded.
	loadingEvents bool

	// loadingFacts indicates if facts are being loaded.
	loadingFacts bool

	// extractingFacts indicates if facts are being extracted.
	extractingFacts bool

	// suggestionsLoading indicates if burst suggestions are being loaded.
	suggestionsLoading bool

	// --- Error states ---

	// deleteError stores any error from delete operation.
	deleteError error

	// editError stores any error from edit operation.
	editError error

	// confirmError stores any error from confirm operation.
	confirmError error

	// suggestionsError stores any error from suggestion detection.
	suggestionsError error

	// --- Progress tracking ---

	// extractedFactsCount tracks the number of facts extracted in current operation.
	extractedFactsCount int

	// --- Screen Orchestration ---

	// activeScreen holds the current screen being displayed (BurstListScreen only).
	activeScreen screens.Screen

	// detailModal holds the burst detail modal (replaces detail screen).
	detailModal *burstmodals.BurstDetailModal

	// eventsModal holds the burst events modal (replaces events screen).
	eventsModal *burstmodals.BurstEventsModal

	// factsModal holds the burst facts modal (replaces facts screen).
	factsModal *burstmodals.BurstFactsModal

	// editModal holds the edit burst modal.
	editModal *burstmodals.EditBurstModal

	// deleteModal holds the delete confirmation modal.
	deleteModal *feedback.ConfirmModal

	// confirmModal holds the burst confirmation modal.
	confirmModal *feedback.ConfirmModal

	// suggestionModal holds the suggestion review modal.
	suggestionModal *burstmodals.SuggestionReviewModal

	// errorModal holds the error modal (shown when operations fail).
	errorModal *feedback.Modal

	// loadingModal holds the loading modal (shown during async operations).
	loadingModal *feedback.Modal

	// modalRegistry manages all modals with unified Update/View handling.
	modalRegistry *intents.ModalRegistry

	// cancelFunc is the cancel function for the current async operation.
	// Call this to cancel loading operations when user presses Esc.
	cancelFunc context.CancelFunc
}

// GetState returns the current state of the intent.
func (i *Intent) GetState() State {
	return i.state
}

// SetState sets the current state of the intent.
func (i *Intent) SetState(state State) {
	i.state = state
}

// IsActive returns whether the intent is currently active.
func (i *Intent) IsActive() bool {
	return i.active
}

// Deactivate marks the intent as inactive.
func (i *Intent) Deactivate() {
	i.active = false
}

// GetFilteredBursts returns the filtered bursts.
func (i *Intent) GetFilteredBursts() []*career.Burst {
	return i.filteredBursts
}

// GetSelectedBurst returns the currently selected burst.
func (i *Intent) GetSelectedBurst() *career.Burst {
	return i.selectedBurst
}

// SetSelectedBurst sets the currently selected burst.
func (i *Intent) SetSelectedBurst(burst *career.Burst) {
	i.selectedBurst = burst
}

// GetSelectedIndex returns the selected index.
func (i *Intent) GetSelectedIndex() int {
	return i.selectedIndex
}

// SetSelectedIndex sets the selected index.
func (i *Intent) SetSelectedIndex(index int) {
	i.selectedIndex = index
}

// GetViewedBursts returns the bursts viewed during the session.
func (i *Intent) GetViewedBursts() []*career.Burst {
	return i.viewedBursts
}

// AddViewedBurst adds a burst to the viewed bursts list.
func (i *Intent) AddViewedBurst(burst *career.Burst) {
	i.viewedBursts = append(i.viewedBursts, burst)
}

// Result returns the final result of the intent.
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
func (i *Intent) SetCancelled() {
	i.result = &intents.IntentResult[*Result]{
		Status: intents.Cancelled,
	}
	i.active = false
}

// GetModalRegistry returns the modal registry.
func (i *Intent) GetModalRegistry() *intents.ModalRegistry {
	return i.modalRegistry
}

// HasActiveModal returns true if any modal is currently visible.
// Uses the modal registry to check for visible modals.
func (i *Intent) HasActiveModal() bool {
	// Rebuild registry to ensure it's current with modal state.
	i.rebuildModalRegistry()

	// Use registry to check for visible modals.
	return i.modalRegistry != nil && i.modalRegistry.HasVisibleModal()
}

// GetDetailModal returns the detail modal (for testing).
func (i *Intent) GetDetailModal() *burstmodals.BurstDetailModal {
	return i.detailModal
}

// GetSuggestionModal returns the suggestion modal (for testing).
func (i *Intent) GetSuggestionModal() *burstmodals.SuggestionReviewModal {
	return i.suggestionModal
}

// IsExtractingFacts returns true if fact extraction is in progress.
func (i *Intent) IsExtractingFacts() bool {
	return i.extractingFacts
}

// GetExtractedFactsCount returns the count of facts extracted in the last operation.
func (i *Intent) GetExtractedFactsCount() int {
	return i.extractedFactsCount
}

// SetLoadingEventsForTesting sets the loadingEvents flag for testing purposes.
// This allows tests to simulate the loading state without requiring a full service setup.
func (i *Intent) SetLoadingEventsForTesting(loading bool) {
	i.loadingEvents = loading
}

// SetLoadingFactsForTesting sets the loadingFacts flag for testing purposes.
func (i *Intent) SetLoadingFactsForTesting(loading bool) {
	i.loadingFacts = loading
}

// NewIntent creates a new BurstManagement intent.
func NewIntent(ctx *IntentContext) (*Intent, error) {
	if ctx == nil {
		return nil, ErrInvalidContext
	}

	if err := ctx.Validate(); err != nil {
		return nil, err
	}

	intent := &Intent{
		BaseIntent:     intents.NewBaseIntent(),
		context:        ctx,
		state:          StateList,
		active:         true,
		filteredBursts: ctx.Bursts,
		selectedIndex:  0,
		viewedBursts:   make([]*career.Burst, 0),
		burstEvents:    make([]*career.Event, 0),
		burstFacts:     make([]*career.Fact, 0),
		modalRegistry:  intents.NewModalRegistry(),
	}

	return intent, nil
}
