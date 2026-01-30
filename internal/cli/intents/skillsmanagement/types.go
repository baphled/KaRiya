package skillsmanagement

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/skills"
	"github.com/baphled/kariya/internal/cli/screens/skills/modals"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	domain "github.com/baphled/kariya/internal/domain/career"
)

// Intent orchestrates the skill management workflow.
type Intent struct {
	*intents.BaseIntent

	// Context holds input parameters and business logic.
	context *IntentContext

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

	// Modals (new architecture with bubbletea-overlay).
	filterModal *modals.FilterModal
	sortModal   *modals.SortModal
	searchModal *modals.SearchModal

	// View/edit modals (modal overlay architecture).
	viewDetailModal  *modals.DetailModal
	addEditModal     *modals.AddEditModal
	deleteModal      *feedback.ConfirmModal
	skillEventsModal *modals.EventsModal
	eventDetailModal *components.ViewEventDetailModal

	// Skill inference modals.
	loadingModal *feedback.Modal
	errorModal   *feedback.Modal

	// Screen orchestration (new architecture).
	activeScreen screens.Screen
	listScreen   *skills.SkillsListScreen

	// Modal registry for unified modal handling.
	modalRegistry *intents.ModalRegistry
}

// SetContext replaces the intent's context, allowing reconfiguration of repositories and filters.
//
// Expected: ctx should be a valid, fully populated IntentContext.
//
// Side effects: replaces the current context reference on the intent.
func (i *Intent) SetContext(ctx *IntentContext) {
	i.context = ctx
}

// GetContext provides access to the intent's configuration and repository dependencies.
//
// Returns: the current IntentContext, or nil if not yet set.
//
// Side effects: None.
func (i *Intent) GetContext() *IntentContext {
	return i.context
}

// SetState transitions the intent to a new state in the workflow state machine.
//
// Expected: state must be a valid State constant defined in constants.go.
//
// Side effects: updates the intent's current state.
func (i *Intent) SetState(state State) {
	i.state = state
}

// GetState exposes the current workflow state for testing and screen orchestration decisions.
//
// Returns: the current State value of the intent.
//
// Side effects: None.
func (i *Intent) GetState() State {
	return i.state
}

// SetActive controls whether the intent processes messages and renders views.
//
// Expected: active is true to enable processing, false to deactivate.
//
// Side effects: updates the intent's active flag.
func (i *Intent) SetActive(active bool) {
	i.active = active
}

// IsActive indicates whether the intent is currently processing messages and rendering.
//
// Returns: true if the intent is active and accepting updates.
//
// Side effects: None.
func (i *Intent) IsActive() bool {
	return i.active
}

// HasActiveModal returns true if a loading or error modal is currently active.
func (i *Intent) HasActiveModal() bool {
	return i.loadingModal != nil || i.errorModal != nil
}

// GetSkills provides access to the loaded skills for testing and screen rendering.
//
// Returns: the current slice of skills held by the intent.
//
// Side effects: None.
func (i *Intent) GetSkills() []*domain.Skill {
	return i.skills
}

// GetSelectedSkill provides access to the skill currently highlighted in the table for detail views and actions.
//
// Returns: the selected skill, or nil if no skill is selected.
//
// Side effects: None.
func (i *Intent) GetSelectedSkill() *domain.Skill {
	return i.selectedSkill
}

// GetSelectedIndex provides the zero-based position of the highlighted skill in the table for navigation state.
//
// Returns: the index of the currently selected skill in the skills slice.
//
// Side effects: None.
func (i *Intent) GetSelectedIndex() int {
	return i.selectedIndex
}
