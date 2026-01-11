package career

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

var (
	// ErrSkillNotFound is returned when a requested skill cannot be found
	ErrSkillNotFound = errors.New("skill not found")

	// ErrDuplicateSkill is returned when attempting to create a skill that already exists
	ErrDuplicateSkill = errors.New("skill already exists")
)

// SkillRepository defines the interface for skill persistence
type SkillRepository interface {
	// Create adds a new skill to the repository
	Create(ctx context.Context, skill *career.Skill) error

	// GetByID retrieves a specific skill by its unique identifier
	GetByID(ctx context.Context, id string) (*career.Skill, error)

	// GetByName retrieves a skill by its name (case-sensitive)
	GetByName(ctx context.Context, name string) (*career.Skill, error)

	// List retrieves skills with optional filtering
	List(ctx context.Context, filters *SkillFilters) ([]*career.Skill, error)

	// Update modifies an existing skill
	Update(ctx context.Context, skill *career.Skill) error

	// Delete removes a skill from the repository
	Delete(ctx context.Context, id string) error

	// GetByCategory retrieves all skills in a specific category
	GetByCategory(ctx context.Context, category string) ([]*career.Skill, error)

	// GetSkillsForEvent retrieves all skills associated with an event
	GetSkillsForEvent(ctx context.Context, eventID string) ([]*career.Skill, error)

	// GetEventCountsForSkills returns a map of skill IDs to event counts
	GetEventCountsForSkills(ctx context.Context) (map[string]int, error)

	// GetLastUsedForSkills returns a map of skill IDs to their last used dates (from events)
	GetLastUsedForSkills(ctx context.Context) (map[string]time.Time, error)

	// GetEventsUsingSkill returns all events that use a specific skill, ordered by date DESC
	GetEventsUsingSkill(ctx context.Context, skillID string) ([]*career.CareerEvent, error)
}

// SkillFilters provides flexible filtering options for skills
type SkillFilters struct {
	// Category to filter by
	Category string

	// Level to filter by (beginner, intermediate, advanced, expert)
	Level string

	// Pagination
	Offset int
	Limit  int
}
