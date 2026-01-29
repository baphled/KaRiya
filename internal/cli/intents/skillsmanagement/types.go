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

	// Screen orchestration (new architecture).
	activeScreen screens.Screen
	listScreen   *skills.SkillsListScreen

	// Modal registry for unified modal handling.
	modalRegistry *intents.ModalRegistry
}

// SetContext sets the intent context.
func (i *Intent) SetContext(ctx *IntentContext) {
	i.context = ctx
}

// GetContext returns the intent context.
func (i *Intent) GetContext() *IntentContext {
	return i.context
}

// SetState sets the current state.
func (i *Intent) SetState(state State) {
	i.state = state
}

// GetState returns the current state.
func (i *Intent) GetState() State {
	return i.state
}

// SetActive sets whether the intent is active.
func (i *Intent) SetActive(active bool) {
	i.active = active
}

// IsActive returns whether the intent is active.
func (i *Intent) IsActive() bool {
	return i.active
}

// GetSkills returns the current skills list.
func (i *Intent) GetSkills() []*domain.Skill {
	return i.skills
}

// GetSelectedSkill returns the currently selected skill.
func (i *Intent) GetSelectedSkill() *domain.Skill {
	return i.selectedSkill
}

// GetSelectedIndex returns the current selection index.
func (i *Intent) GetSelectedIndex() int {
	return i.selectedIndex
}
