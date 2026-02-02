package skillsmanagement

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	burstmodals "github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
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
	loadingModal          *feedback.Modal
	feedbackModal         *feedback.Modal
	skillSuggestionModal  *burstmodals.SuggestionReviewModal
	suggestionEventsModal *modals.EventsModal

	// Screen orchestration (new architecture).
	activeScreen screens.Screen
	listScreen   *skills.SkillsListScreen

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
func (i *Intent) SetContext(ctx *IntentContext) {
	i.context = ctx
}

// GetContext provides access to the intent's configuration and repository dependencies.
//
// Returns:
//   - A fully initialized IntentContext ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetContext() *IntentContext {
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
