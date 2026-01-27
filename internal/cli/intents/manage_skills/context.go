// Package manage_skills implements the ManageSkills intent for managing user-defined skills.
package manage_skills

import (
	"context"
	"errors"

	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/repository/career"
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

// HasActiveFilters returns true if any non-default filters are active.
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

// Clear resets filters in FIFO order (most recent filter first).
// Search is most recent (most specific), sort is least recent (most general).
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

	// Filters holds the active filter and sort state.
	Filters *Filters
}

// NewIntentContext creates a new context with default values.
func NewIntentContext(ctx context.Context, skillRepo career.SkillRepository) *IntentContext {
	return &IntentContext{
		Ctx:             ctx,
		SkillRepository: skillRepo,
		Filters:         &Filters{},
	}
}

// Validate ensures the context is complete.
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

// LoadSkills loads skills from the repository with current filters.
func (c *IntentContext) LoadSkills() ([]*domain.Skill, error) {
	var repoFilters *career.SkillFilters
	if c.Filters != nil {
		repoFilters = &career.SkillFilters{
			Category:  c.Filters.Category,
			Level:     c.Filters.Level,
			MinEvents: c.Filters.MinEvents,
			SortBy:    c.Filters.SortBy,
			SortOrder: c.Filters.SortOrder,
		}
	}

	return c.SkillRepository.List(c.Ctx, repoFilters)
}

// GetEventCounts returns the event counts for all skills.
func (c *IntentContext) GetEventCounts() (map[string]int, error) {
	return c.SkillRepository.GetEventCountsForSkills(c.Ctx)
}

// GetEventsForSkill returns the events associated with a specific skill.
func (c *IntentContext) GetEventsForSkill(skillID string) ([]*domain.CareerEvent, error) {
	return c.SkillRepository.GetEventsUsingSkill(c.Ctx, skillID)
}
