package browse_timeline

import (
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/timeline/modals"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
)

// Intent implements the Intent interface for browsing career events.
// It owns the complete lifecycle of timeline browsing, including:
// - Displaying a filtered and sorted timeline of events
// - Selecting and viewing event details
// - Returning the selected event or cancelling
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

	// filteredEvents are the events after applying current filters.
	filteredEvents []*career.CareerEvent

	// selectedIndex is the index of the currently selected event.
	selectedIndex int

	// filters is the current filter and sort state.
	filters *Filters

	// filterStack tracks active filters in FIFO order for progressive clearing.
	filterStack *intents.FilterStack

	// selectedEvent is the event currently being viewed.
	selectedEvent *career.CareerEvent

	// selectedFacts are facts selected from the event.
	selectedFacts []*career.Fact

	// viewedEvents tracks events viewed during the session.
	viewedEvents []*career.CareerEvent

	// deleteError stores any error from delete operation.
	deleteError error

	// --- Screen Orchestration ---

	// activeScreen holds the current screen being displayed.
	activeScreen screens.Screen

	// filterModal holds the filter modal (shown over the list).
	filterModal *modals.FilterModal

	// searchModal holds the search modal for text search.
	searchModal *modals.SearchModal

	// sortModal holds the sort modal for sorting events.
	sortModal *modals.SortModal

	// deleteModal holds the delete confirmation modal (shown over the list).
	deleteModal *feedback.ConfirmModal

	// quickAddModal holds the quick add event modal (shown over the list).
	quickAddModal *modals.QuickAddModal

	// editModal holds the edit event modal (shown over the list).
	editModal *modals.EditModal

	// viewDetailModal holds the event detail viewer modal (shown over the list).
	viewDetailModal *modals.EventDetailModal

	// viewSkillsModal holds the skills viewer modal (shown over event detail).
	viewSkillsModal *modals.SkillsDetailModal

	// errorModal holds the error modal (shown when operations fail).
	errorModal *feedback.Modal
}
