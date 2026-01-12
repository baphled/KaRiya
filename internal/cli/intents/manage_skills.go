package intents

import (
	"context"
	"errors"

	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

var (
	// ErrNoSkillSelected is returned when no skill is selected
	ErrNoSkillSelected = errors.New("no skill selected")

	// ErrInvalidSkillData is returned when skill data is invalid
	ErrInvalidSkillData = errors.New("invalid skill data")
)

// SkillsState represents the current state of the ManageSkills intent
type SkillsState string

const (
	SkillsStateList              SkillsState = "list"                // View all skills grouped by category
	SkillsStateDetail            SkillsState = "detail"              // View single skill details
	SkillsStateDetailEvents      SkillsState = "detail_events"       // View events using this skill
	SkillsStateDetailEventDetail SkillsState = "detail_event_detail" // View single event details from skill events
	SkillsStateAdd               SkillsState = "add"                 // Add new skill (huh form)
	SkillsStateEdit              SkillsState = "edit"                // Edit existing skill (huh form)
	SkillsStateDelete            SkillsState = "delete"              // Confirm deletion
	SkillsStateFilter            SkillsState = "filter"              // Filter menu
	SkillsStateSort              SkillsState = "sort"                // Sort menu
)

// ManageSkillsContext holds the context and dependencies for ManageSkills intent
type ManageSkillsContext struct {
	// Services and repositories
	Service         *careerservice.Service
	SkillRepository career.SkillRepository

	// Context
	Ctx context.Context
}

// ManageSkillsResult is returned when the intent completes
type ManageSkillsResult struct {
	// Action performed (if any)
	Action string // "created", "updated", "deleted", "cancelled"

	// Skill involved (if any)
	Skill *domain.Skill
}

// Custom message types for ManageSkills state transitions

// SkillsLoadedMsg is sent when skills are loaded from the repository
type SkillsLoadedMsg struct {
	Skills []*domain.Skill
	Error  error
}

// SkillFormCompleteMsg is sent when the skill form is completed or cancelled
type SkillFormCompleteMsg struct {
	Skill     *domain.Skill
	Cancelled bool
	Error     error
}

// SkillCreatedMsg is sent when a skill is created
type SkillCreatedMsg struct {
	Skill *domain.Skill
	Error error
}

// SkillUpdatedMsg is sent when a skill is updated
type SkillUpdatedMsg struct {
	Skill *domain.Skill
	Error error
}

// SkillDeletedMsg is sent when a skill is deleted
type SkillDeletedMsg struct {
	SkillID string
	Error   error
}

// SkillEventsLoadedMsg is sent when events for a skill are loaded
type SkillEventsLoadedMsg struct {
	Events []*domain.CareerEvent
	Error  error
}

// RequestBrowseEventMsg requests that the app route to BrowseTimeline intent.
// This is sent to the app router which will activate BrowseTimeline with the selected event.
type RequestBrowseEventMsg struct {
	Event     *domain.CareerEvent
	AllEvents []*domain.CareerEvent
	SkillName string // For context in breadcrumbs
}
