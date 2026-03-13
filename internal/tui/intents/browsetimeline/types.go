package browsetimeline

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/tui/intents"
	burstviews "github.com/baphled/kariya/internal/tui/views/burst"
	eventviews "github.com/baphled/kariya/internal/tui/views/event"
	skillviews "github.com/baphled/kariya/internal/tui/views/skill"
	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

// Intent implements the Intent interface for browsing career events.
// It owns the complete lifecycle of timeline browsing, including:
// - Displaying a filtered and sorted timeline of events
// - Selecting and viewing event details
// - Returning the selected event or cancelling.
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

	// filteredEvents are the events after applying current filters.
	filteredEvents []*career.Event

	// selectedIndex is the index of the currently selected event.
	selectedIndex int

	// filters is the current filter and sort state.
	filters *Filters

	// filterStack tracks active filters in FIFO order for progressive clearing.
	filterStack *behaviors.FilterStack

	// selectedEvent is the event currently being viewed.
	selectedEvent *career.Event

	// selectedFacts are facts selected from the event.
	selectedFacts []*career.Fact

	// viewedEvents tracks events viewed during the session.
	viewedEvents []*career.Event

	// deleteError stores any error from delete operation.
	deleteError error

	// --- Screen Orchestration ---

	// activeView holds the current view being displayed.
	activeView widgets.View

	// --- Modals (priority order: error > delete > edit > detail > skills > filter > search > sort > quickAdd) ---
	// The modalRegistry manages priority and ensures only the highest-priority visible modal
	// receives updates and renders. Priority is defined in initializeModalRegistry().

	// searchAdapter wraps the shared SearchView as a managed modal.
	searchAdapter *intents.FormViewAdapter

	// sortAdapter wraps the shared SortView as a managed modal.
	sortAdapter *intents.FormViewAdapter

	// filterAdapter wraps the event Filter view as a managed modal.
	filterAdapter *intents.FormViewAdapter

	// deleteModal holds the delete confirmation modal (shown over the list).
	deleteModal *feedback.ConfirmModal

	// quickAdd holds the quick add event modal (shown over the list).
	quickAdd *eventviews.QuickAdd

	// edit holds the edit event modal (shown over the list).
	edit *eventviews.Edit

	// viewDetail holds the event detail viewer modal (shown over the list).
	viewDetail *eventviews.Detail

	// viewSkills holds the skills viewer modal (shown over event detail).
	viewSkills *eventviews.SkillsDetail

	// skillPicker holds the skill picker modal for linking existing skills.
	skillPicker *eventviews.SkillPicker

	// availableSkills caches the current set of domain skills available for linking.
	availableSkills []*career.Skill

	// selectedEventSkills caches the current domain skills linked to the selected event.
	selectedEventSkills []*career.Skill

	// skillAddModal holds the add/edit modal for creating new skills.
	skillAddModal *skillviews.AddEdit

	// skillSuggestionModal holds the skill suggestion review modal.
	skillSuggestionModal *burstviews.SuggestionReview

	// skillSuggestions caches the domain skill suggestions shown in the review modal.
	skillSuggestions []skillinference.SkillSuggestion

	// skillService for creating new skills.
	skillService SkillCreator

	// errorModal holds the error modal (shown when operations fail).
	errorModal *feedback.Modal

	// modalRegistry manages all overlays with unified Update/View handling.
	modalRegistry *intents.ModalRegistry
}
