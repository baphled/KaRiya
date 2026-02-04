//go:generate mockgen -destination=../../testutil/mocks/repository/skill_repository_mock.go -package=mockrepo github.com/baphled/kariya/internal/repository/career SkillRepository

package career

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

var (
	// ErrSkillNotFound is returned when a requested skill cannot be found.
	ErrSkillNotFound = errors.New("skill not found")

	// ErrDuplicateSkill is returned when attempting to create a skill that already exists.
	ErrDuplicateSkill = errors.New("skill already exists")
)

// SkillRepository defines the persistence contract for career skills.
//
// It covers full CRUD operations, category-based lookups, event-skill
// relationship queries, and aggregate statistics such as event counts and
// last-used dates. Implementations back the skill management intent and the
// timeline detail view where associated skills are displayed. The interface
// follows the repository pattern to decouple domain logic from storage
// details.
//
// Every method accepts a ctx parameter that carries request-scoped deadlines
// and cancellation signals.
//
// Methods:
//
//   - Create persists a new skill. The skill parameter is the domain object
//     to store. Returns nil on success or an error if the write fails.
//
//   - GetByID looks up a single skill by its UUID. The id parameter is the
//     unique identifier string. Returns the matching *career.Skill and nil
//     on success, or nil and ErrSkillNotFound when no record exists.
//
//   - GetByName looks up a skill by name using case-insensitive matching.
//     The name parameter is the display name to match. Returns the matching
//     *career.Skill and nil on success, or nil and ErrSkillNotFound when no
//     record exists.
//
//   - List returns all skills that match the optional filters. The filters
//     parameter may be nil for unfiltered results. Returns the matching
//     []*career.Skill slice and nil, or nil and an error on failure.
//
//   - Update replaces the stored skill with the supplied version. The skill
//     parameter is the updated domain object whose ID must already exist.
//     Returns nil on success or an error if the skill does not exist.
//
//   - Delete removes a skill by its unique identifier. The id parameter is
//     the skill UUID. Returns nil on success or an error if the skill does
//     not exist.
//
//   - GetByCategory returns every skill whose Category field equals the
//     category parameter string. Returns the matching []*career.Skill slice
//     and nil, or nil and an error on failure.
//
//   - GetSkillsForEvent returns all skills linked to the given event. The
//     eventID parameter is the event UUID. Returns the matching
//     []*career.Skill slice and nil, or nil and an error on failure.
//
//   - GetEventCountsForSkills returns a map[string]int keyed by skill ID
//     whose values are the number of events associated with each skill.
//     The ctx parameter is the only input. Returns the map and nil, or nil
//     and an error on failure.
//
//   - GetLastUsedForSkills returns a map[string]time.Time keyed by skill ID
//     whose values are the most recent event date for each skill. The ctx
//     parameter is the only input. Returns the map and nil, or nil and an
//     error on failure.
//
//   - GetEventsUsingSkill returns all events linked to the given skill,
//     ordered by date descending. The skillID parameter is the skill UUID.
//     Returns the matching []*career.Event slice and nil, or nil and an
//     error on failure.
type SkillRepository interface {
	Create(ctx context.Context, skill *career.Skill) error
	GetByID(ctx context.Context, id string) (*career.Skill, error)
	GetByName(ctx context.Context, name string) (*career.Skill, error)
	List(ctx context.Context, filters *SkillListFilters) ([]*career.Skill, error)
	Update(ctx context.Context, skill *career.Skill) error
	Delete(ctx context.Context, id string) error
	GetByCategory(ctx context.Context, category string) ([]*career.Skill, error)
	GetSkillsForEvent(ctx context.Context, eventID string) ([]*career.Skill, error)
	GetEventCountsForSkills(ctx context.Context) (map[string]int, error)
	GetLastUsedForSkills(ctx context.Context) (map[string]time.Time, error)
	GetEventsUsingSkill(ctx context.Context, skillID string) ([]*career.Event, error)
}

// SkillListFilters provides flexible filtering options for skill queries.
//
// Category restricts results to a single category string. Level restricts
// results to a proficiency tier (beginner, intermediate, advanced, expert).
// MinEvents filters to skills with at least that many event associations;
// set to 1 to show only skills that have been used. SortBy selects the
// ordering field ("name", "events", "last_used", or "category"); the default
// is "name". SortOrder is "asc" or "desc"; the default is "asc". Offset and
// Limit control pagination.
type SkillListFilters struct {
	Category  string
	Level     string
	MinEvents int
	SortBy    string
	SortOrder string
	Offset    int
	Limit     int
}
