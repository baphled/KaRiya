package skillsmanagement

import (
	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	event "github.com/baphled/kariya/internal/tui/views/event"

	"github.com/baphled/kariya/internal/tui/intents"
	burstviews "github.com/baphled/kariya/internal/tui/views/burst"
	skillviews "github.com/baphled/kariya/internal/tui/views/skill"
	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

// Intent orchestrates the skill management workflow.
type Intent struct {
	*intents.BaseIntent

	// Context holds input parameters and business logic.
	context *IntentValidator

	// State machine fields.
	state  State
	active bool
	result *intents.IntentResult[*Result]

	// Skills data.
	skills        []*domain.Skill
	selectedIndex int
	selectedSkill *domain.Skill

	// TableBehavior provides type-safe table operations for skills.
	tableBehavior *behaviors.TableBehavior[*domain.Skill]

	// Detail view data.
	eventCounts  map[string]int
	skillEvents  []*domain.Event
	eventsLoaded bool

	// Form view adapters (search, filter, sort) — implement ManagedModal directly.
	filterAdapter *intents.FormViewAdapter
	sortAdapter   *intents.FormViewAdapter
	searchAdapter *intents.FormViewAdapter

	// View/edit modals (modal overlay architecture).
	viewDetailModal  *skillviews.Detail
	addEditModal     *skillviews.AddEdit
	deleteModal      *feedback.ConfirmModal
	skillEventsModal *skillviews.Events
	eventDetailModal *event.Detail

	// Skill inference modals.
	loadingModal          *feedback.Modal
	feedbackModal         *feedback.Modal
	skillSuggestionModal  *burstviews.SuggestionReview
	suggestionEventsModal *skillviews.Events
	skillSuggestions      []skillinference.SkillSuggestion

	// Screen orchestration (new architecture).
	activeView widgets.View

	// Modal registry for unified modal handling.
	modalRegistry *intents.ModalRegistry
}

// SetContext replaces the intent's context, allowing reconfiguration of repositories and filters.
//
// Expected:
//   - intentcontext must be valid.
//
// Side effects:
//   - None.
func (i *Intent) SetContext(ctx *IntentValidator) {
	i.context = ctx
}

// GetContext provides access to the intent's configuration and repository dependencies.
//
// Returns:
//   - A fully initialized IntentValidator ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetContext() *IntentValidator {
	return i.context
}

// SetState transitions the intent to a new state in the workflow state machine.
//
// Expected:
//   - state must be valid.
//
// Side effects:
//   - None.
func (i *Intent) SetState(state State) {
	i.state = state
}

// GetState exposes the current workflow state for testing and screen orchestration decisions.
//
// Returns:
//   - A State value.
//
// Side effects:
//   - None.
func (i *Intent) GetState() State {
	return i.state
}

// SetActive controls whether the intent processes messages and renders views.
//
// Expected:
//   - bool must be valid.
//
// Side effects:
//   - None.
func (i *Intent) SetActive(active bool) {
	i.active = active
}

// IsActive indicates whether the intent is currently processing messages and rendering.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (i *Intent) IsActive() bool {
	return i.active
}

// HasActiveModal returns true if a loading, feedback, suggestion, or events modal is currently active.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (i *Intent) HasActiveModal() bool {
	return i.loadingModal != nil || i.feedbackModal != nil ||
		(i.skillSuggestionModal != nil && i.skillSuggestionModal.IsVisible()) ||
		(i.suggestionEventsModal != nil && i.suggestionEventsModal.IsVisible())
}

// GetFeedbackModal returns the current feedback modal for testing.
//
// Returns:
//   - A fully initialized feedback.Modal ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetFeedbackModal() *feedback.Modal {
	return i.feedbackModal
}

// GetLoadingModal returns the current loading modal for testing.
//
// Returns:
//   - A fully initialized feedback.Modal ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetLoadingModal() *feedback.Modal {
	return i.loadingModal
}

// GetViewDetailModal returns the current view detail modal for testing.
//
// Returns:
//   - A fully initialized skillviews.Detail modal ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetViewDetailModal() *skillviews.Detail {
	return i.viewDetailModal
}

// GetSkills provides access to the loaded skills for testing and screen rendering.
//
// Returns:
//   - A []*domain.Skill value.
//
// Side effects:
//   - None.
func (i *Intent) GetSkills() []*domain.Skill {
	return i.skills
}

// GetSelectedSkill provides access to the skill currently highlighted in the table for detail views and actions.
//
// Returns:
//   - A fully initialized domain.Skill ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedSkill() *domain.Skill {
	return i.selectedSkill
}

// GetSelectedIndex provides the zero-based position of the highlighted skill in the table for navigation state.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedIndex() int {
	return i.selectedIndex
}
