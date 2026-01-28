package skills_management

import "github.com/baphled/kariya/internal/domain/career"

// Result is the output of the ManageSkills intent.
type Result struct {
	// Action describes the action performed (created, updated, deleted, cancelled).
	Action string

	// Skill is the affected skill (may be nil if no specific skill was affected).
	Skill *career.Skill
}
