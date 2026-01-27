package skills_management

import domain "github.com/baphled/kariya/internal/domain/career"

// Custom message types for ManageSkills state transitions.
// ALL *Msg structs MUST be in this file.

// SkillsLoadedMsg is sent when skills are loaded from the repository.
type SkillsLoadedMsg struct {
	Skills []*domain.Skill
	Error  error
}

// SkillFormCompleteMsg is sent when the skill form is completed or cancelled.
type SkillFormCompleteMsg struct {
	Skill     *domain.Skill
	Cancelled bool
	Error     error
}

// SkillCreatedMsg is sent when a skill is created.
type SkillCreatedMsg struct {
	Skill *domain.Skill
	Error error
}

// SkillUpdatedMsg is sent when a skill is updated.
type SkillUpdatedMsg struct {
	Skill *domain.Skill
	Error error
}

// SkillDeletedMsg is sent when a skill is deleted.
type SkillDeletedMsg struct {
	SkillID string
	Error   error
}

// SkillEventsLoadedMsg is sent when events for a skill are loaded (state-based flow).
type SkillEventsLoadedMsg struct {
	Events []*domain.CareerEvent
	Error  error
}

// SkillEventsForModalLoadedMsg is sent when events for a skill are loaded (modal flow).
type SkillEventsForModalLoadedMsg struct {
	Events []*domain.CareerEvent
	Error  error
}

// RequestBrowseEventMsg requests that the app route to BrowseTimeline intent.
// This is sent to the app router which will activate BrowseTimeline with the selected event.
type RequestBrowseEventMsg struct {
	Event     *domain.CareerEvent
	AllEvents []*domain.CareerEvent
	SkillName string // For context in breadcrumbs.
}
