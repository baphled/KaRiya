//go:generate mockgen -destination=../../../testutil/mocks/intent/skills_skill_inference_service_mock.go -package=mockintent -mock_names=SkillInferenceService=MockSkillsSkillInferenceService github.com/baphled/kariya/internal/cli/intents/skillsmanagement SkillInferenceService

package skillsmanagement

import (
	"context"
	"errors"

	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
)

var (
	// ErrNoSkillSelected is returned when no skill is selected.
	ErrNoSkillSelected = errors.New("no skill selected")

	// ErrInvalidSkillData is returned when skill data is invalid.
	ErrInvalidSkillData = errors.New("invalid skill data")

	// ErrRepositoryNotAvailable is returned when the repository is not available.
	ErrRepositoryNotAvailable = errors.New("repository not available")

	// ErrContextNotAvailable is returned when context is not available.
	ErrContextNotAvailable = errors.New("context not available")
)

// Filters holds the active filter and sort state for skills.
type Filters struct {
	Category   string
	Level      string
	MinEvents  int
	SearchText string
	SortBy     string
	SortOrder  string
}

// HasActiveFilters checks whether any filter or sort criteria differ from the default empty state.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (f *Filters) HasActiveFilters() bool {
	if f == nil {
		return false
	}
	return f.Category != "" ||
		f.Level != "" ||
		f.MinEvents > 0 ||
		f.SearchText != "" ||
		f.SortBy != ""
}

// Clear progressively removes the most recently applied filter layer in FIFO order.
//
// Side effects:
//   - None.
func (f *Filters) Clear() {
	if f == nil {
		return
	}

	// Clear in FIFO order: search -> filter -> sort.
	if f.SearchText != "" {
		f.SearchText = ""
		return
	}

	if f.Category != "" || f.Level != "" || f.MinEvents > 0 {
		f.Category = ""
		f.Level = ""
		f.MinEvents = 0
		return
	}

	// Clear sort (least specific).
	f.SortBy = ""
	f.SortOrder = ""
}

// IntentContext holds input parameters and business logic for the ManageSkills intent.
type IntentContext struct {
	// Ctx is the context for repository operations.
	Ctx context.Context

	// SkillRepository is the repository for skill CRUD operations.
	SkillRepository career.SkillRepository

	// EventRepository is the repository for event operations (needed for skill inference).
	EventRepository career.EventRepository

	// SkillInferenceService provides skill inference from events.
	SkillInferenceService SkillInferenceService

	// Filters holds the active filter and sort state.
	Filters *Filters
}

// SkillInferenceService defines the interface for skill inference operations.
type SkillInferenceService interface {
	InferSkillsFromEvents(ctx context.Context, events []*domain.Event) (*skillinference.InferenceResult, error)
	CreateSkillsFromSuggestions(ctx context.Context, suggestions []skillinference.SkillSuggestion) ([]*domain.Skill, error)
}

// NewIntentContext constructs an IntentContext wired to the given repository with default empty filters.
//
// Expected: ctx must be a non-nil context.Context; skillRepo must be a non-nil SkillRepository.
//
// Returns: a fully initialized IntentContext ready for Validate and use by NewIntent.
//
// Side effects: None.
func NewIntentContext(ctx context.Context, skillRepo career.SkillRepository) *IntentContext {
	return &IntentContext{
		Ctx:             ctx,
		SkillRepository: skillRepo,
		Filters:         &Filters{},
	}
}

// Validate ensures all required dependencies are present before the intent can be constructed.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (c *IntentContext) Validate() error {
	if c.Ctx == nil {
		return ErrContextNotAvailable
	}
	if c.SkillRepository == nil {
		return ErrRepositoryNotAvailable
	}
	if c.Filters == nil {
		c.Filters = &Filters{}
	}
	return nil
}

// LoadSkills fetches skills from the repository, translating the current Filters into repository query parameters.
//
// Returns: the filtered list of skills, or an error if the repository call fails.
//
// Side effects: performs a repository read operation against the underlying data store.
func (c *IntentContext) LoadSkills() ([]*domain.Skill, error) {
	var repoFilters *career.SkillListFilters
	if c.Filters != nil {
		repoFilters = &career.SkillListFilters{
			Category:  c.Filters.Category,
			Level:     c.Filters.Level,
			MinEvents: c.Filters.MinEvents,
			SortBy:    c.Filters.SortBy,
			SortOrder: c.Filters.SortOrder,
		}
	}

	return c.SkillRepository.List(c.Ctx, repoFilters)
}

// GetEventCounts retrieves the number of associated events per skill for display in the list view.
//
// Returns: a map of skill ID to event count, or an error if the repository call fails.
//
// Side effects: performs a repository read operation against the underlying data store.
func (c *IntentContext) GetEventCounts() (map[string]int, error) {
	return c.SkillRepository.GetEventCountsForSkills(c.Ctx)
}

// GetEventsForSkill fetches the career events linked to a specific skill for the detail and events modals.
//
// Expected: skillID must be a non-empty string identifying an existing skill.
//
// Returns: the list of events associated with the skill, or an error if the repository call fails.
//
// Side effects: performs a repository read operation against the underlying data store.
func (c *IntentContext) GetEventsForSkill(skillID string) ([]*domain.Event, error) {
	return c.SkillRepository.GetEventsUsingSkill(c.Ctx, skillID)
}
