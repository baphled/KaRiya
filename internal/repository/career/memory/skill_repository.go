package memory

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	career_repo "github.com/baphled/kariya/internal/repository/career"
	"github.com/google/uuid"
)

// Compile-time interface check.
var _ career_repo.SkillRepository = (*SkillRepository)(nil)

// SkillRepository provides an in-memory implementation of the SkillRepository interface.
type SkillRepository struct {
	skills      map[string]*career.Skill
	eventSkills map[string][]string         // eventID -> []skillID
	skillEvents map[string][]string         // skillID -> []eventID
	eventRepo   career_repo.EventRepository // Reference to event repository for GetEventsUsingSkill
	mu          sync.RWMutex
}

// NewSkillRepository creates a new in-memory skill repository
func NewSkillRepository() *SkillRepository {
	return &SkillRepository{
		skills:      make(map[string]*career.Skill),
		eventSkills: make(map[string][]string),
		skillEvents: make(map[string][]string),
	}
}

// SetEventRepository sets the event repository for cross-repository queries
func (r *SkillRepository) SetEventRepository(repo career_repo.EventRepository) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.eventRepo = repo
}

// Create adds a new skill to the in-memory store.
func (r *SkillRepository) Create(_ context.Context, skill *career.Skill) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Generate a unique ID if not provided
	if skill.ID == "" {
		skill.ID = uuid.New().String()
	}

	// Validate the skill
	if err := skill.Validate(); err != nil {
		return err
	}

	// Check for duplicate name
	for _, s := range r.skills {
		if s.Name == skill.Name {
			return career_repo.ErrDuplicateSkill
		}
	}

	// Set timestamps
	now := time.Now()
	skill.CreatedAt = now
	skill.UpdatedAt = now

	// Store the skill
	r.skills[skill.ID] = skill

	return nil
}

// GetByID retrieves a specific skill by its unique identifier.
func (r *SkillRepository) GetByID(_ context.Context, id string) (*career.Skill, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	skill, exists := r.skills[id]
	if !exists {
		return nil, career_repo.ErrSkillNotFound
	}

	return skill, nil
}

// GetByName retrieves a skill by its name (case-sensitive)
func (r *SkillRepository) GetByName(_ context.Context, name string) (*career.Skill, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, skill := range r.skills {
		if skill.Name == name {
			return skill, nil
		}
	}

	return nil, career_repo.ErrSkillNotFound
}

// List retrieves skills with optional filtering.
func (r *SkillRepository) List(ctx context.Context, filters *career_repo.SkillListFilters) ([]*career.Skill, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	eventCounts := r.buildEventCounts()

	var skills []*career.Skill
	for _, skill := range r.skills {
		skills = append(skills, skill)
	}

	skills = r.applyFilters(skills, filters, eventCounts)
	r.applySorting(ctx, skills, filters, eventCounts)
	return r.applyPagination(skills, filters), nil
}

func (r *SkillRepository) buildEventCounts() map[string]int {
	counts := make(map[string]int, len(r.skillEvents))
	for skillID, eventIDs := range r.skillEvents {
		counts[skillID] = len(eventIDs)
	}
	return counts
}

func (r *SkillRepository) applyFilters(
	skills []*career.Skill, filters *career_repo.SkillListFilters, eventCounts map[string]int,
) []*career.Skill {
	if filters == nil {
		return skills
	}

	var filtered []*career.Skill
	for _, skill := range skills {
		if filters.Category != "" && skill.Category != filters.Category {
			continue
		}
		if filters.Level != "" && skill.Level != filters.Level {
			continue
		}
		if filters.MinEvents > 0 && eventCounts[skill.ID] < filters.MinEvents {
			continue
		}
		filtered = append(filtered, skill)
	}
	return filtered
}

func (r *SkillRepository) applySorting(
	ctx context.Context, skills []*career.Skill, filters *career_repo.SkillListFilters, eventCounts map[string]int,
) {
	sortBy := "name"
	desc := false
	if filters != nil {
		if filters.SortBy != "" {
			sortBy = filters.SortBy
		}
		desc = filters.SortOrder == "desc"
	}

	switch sortBy {
	case "name":
		r.sortByName(skills, desc)
	case "events":
		r.sortByEventCount(skills, desc, eventCounts)
	case "last_used":
		r.sortByLastUsed(ctx, skills, desc)
	case "category":
		r.sortByCategory(skills, desc)
	default:
		r.sortByName(skills, false)
	}
}

func (r *SkillRepository) sortByName(skills []*career.Skill, desc bool) {
	sort.Slice(skills, func(i, j int) bool {
		less := strings.ToLower(skills[i].Name) < strings.ToLower(skills[j].Name)
		if desc {
			return !less
		}
		return less
	})
}

func (r *SkillRepository) sortByEventCount(skills []*career.Skill, desc bool, eventCounts map[string]int) {
	sort.Slice(skills, func(i, j int) bool {
		countI := eventCounts[skills[i].ID]
		countJ := eventCounts[skills[j].ID]
		if countI == countJ {
			return strings.ToLower(skills[i].Name) < strings.ToLower(skills[j].Name)
		}
		if desc {
			return countI > countJ
		}
		return countI < countJ
	})
}

func (r *SkillRepository) sortByLastUsed(ctx context.Context, skills []*career.Skill, desc bool) {
	lastUsedDates := r.buildLastUsedDates(ctx)

	sort.Slice(skills, func(i, j int) bool {
		dateI, hasI := lastUsedDates[skills[i].ID]
		dateJ, hasJ := lastUsedDates[skills[j].ID]

		// Skills without dates sort last for DESC, first for ASC.
		if !hasI && !hasJ {
			return strings.ToLower(skills[i].Name) < strings.ToLower(skills[j].Name)
		}
		if !hasI {
			return !desc
		}
		if !hasJ {
			return desc
		}

		if dateI.Equal(dateJ) {
			return strings.ToLower(skills[i].Name) < strings.ToLower(skills[j].Name)
		}
		if desc {
			return dateI.After(dateJ)
		}
		return dateI.Before(dateJ)
	})
}

func (r *SkillRepository) buildLastUsedDates(ctx context.Context) map[string]time.Time {
	lastUsedDates := make(map[string]time.Time)
	if r.eventRepo == nil {
		return lastUsedDates
	}

	for skillID, eventIDs := range r.skillEvents {
		var maxDate time.Time
		for _, eventID := range eventIDs {
			if event, err := r.eventRepo.GetByID(ctx, eventID); err == nil {
				if event.Date.After(maxDate) {
					maxDate = event.Date
				}
			}
		}
		if !maxDate.IsZero() {
			lastUsedDates[skillID] = maxDate
		}
	}
	return lastUsedDates
}

func (r *SkillRepository) sortByCategory(skills []*career.Skill, desc bool) {
	sort.Slice(skills, func(i, j int) bool {
		if skills[i].Category == skills[j].Category {
			return strings.ToLower(skills[i].Name) < strings.ToLower(skills[j].Name)
		}
		less := skills[i].Category < skills[j].Category
		if desc {
			return !less
		}
		return less
	})
}

func (r *SkillRepository) applyPagination(skills []*career.Skill, filters *career_repo.SkillListFilters) []*career.Skill {
	if filters == nil {
		return skills
	}

	if filters.Offset > 0 {
		if filters.Offset >= len(skills) {
			return []*career.Skill{}
		}
		skills = skills[filters.Offset:]
	}

	if filters.Limit > 0 && filters.Limit < len(skills) {
		skills = skills[:filters.Limit]
	}
	return skills
}

// Update modifies an existing skill
func (r *SkillRepository) Update(_ context.Context, skill *career.Skill) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if skill exists
	if _, exists := r.skills[skill.ID]; !exists {
		return career_repo.ErrSkillNotFound
	}

	// Validate the skill
	if err := skill.Validate(); err != nil {
		return err
	}

	// Check for duplicate name (excluding current skill)
	for id, s := range r.skills {
		if id != skill.ID && s.Name == skill.Name {
			return career_repo.ErrDuplicateSkill
		}
	}

	// Update timestamp
	skill.UpdatedAt = time.Now()

	// Store the updated skill
	r.skills[skill.ID] = skill

	return nil
}

// Delete removes a skill from the repository
func (r *SkillRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if skill exists
	if _, exists := r.skills[id]; !exists {
		return career_repo.ErrSkillNotFound
	}

	// Remove the skill
	delete(r.skills, id)

	// Clean up event-skill associations
	delete(r.skillEvents, id)
	for eventID := range r.eventSkills {
		r.eventSkills[eventID] = removeString(r.eventSkills[eventID], id)
	}

	return nil
}

// GetByCategory retrieves all skills in a specific category
func (r *SkillRepository) GetByCategory(_ context.Context, category string) ([]*career.Skill, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*career.Skill

	for _, skill := range r.skills {
		if skill.Category == category {
			result = append(result, skill)
		}
	}

	// Sort by name
	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
	})

	return result, nil
}

// GetSkillsForEvent retrieves all skills associated with an event
func (r *SkillRepository) GetSkillsForEvent(_ context.Context, eventID string) ([]*career.Skill, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	skillIDs, exists := r.eventSkills[eventID]
	if !exists {
		return []*career.Skill{}, nil
	}

	var result []*career.Skill
	for _, skillID := range skillIDs {
		if skill, exists := r.skills[skillID]; exists {
			result = append(result, skill)
		}
	}

	return result, nil
}

// GetEventCountsForSkills returns a map of skill IDs to event counts
func (r *SkillRepository) GetEventCountsForSkills(_ context.Context) (map[string]int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	counts := make(map[string]int)

	for skillID, eventIDs := range r.skillEvents {
		counts[skillID] = len(eventIDs)
	}

	return counts, nil
}

// GetLastUsedForSkills returns a map of skill IDs to their last used dates (from events)
func (r *SkillRepository) GetLastUsedForSkills(_ context.Context) (map[string]time.Time, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	lastUsed := make(map[string]time.Time)

	// Note: This requires access to events, which we don't have in this simple implementation
	// For now, return empty map. In real usage with SetEventRepository, this would work.
	// The SQL/GORM implementation handles this properly with JOINs.

	return lastUsed, nil
}

// GetEventsUsingSkill returns all events that use a specific skill, ordered by date DESC
func (r *SkillRepository) GetEventsUsingSkill(ctx context.Context, skillID string) ([]*career.CareerEvent, error) {
	r.mu.RLock()
	eventIDs := r.skillEvents[skillID]
	eventRepo := r.eventRepo
	r.mu.RUnlock()

	if len(eventIDs) == 0 {
		return []*career.CareerEvent{}, nil
	}

	// If we don't have an event repository, return empty list
	if eventRepo == nil {
		return []*career.CareerEvent{}, nil
	}

	// Fetch events
	var events []*career.CareerEvent
	for _, eventID := range eventIDs {
		event, err := eventRepo.GetByID(ctx, eventID)
		if err == nil {
			events = append(events, event)
		}
	}

	// Sort by date descending
	sort.Slice(events, func(i, j int) bool {
		return events[i].Date.After(events[j].Date)
	})

	return events, nil
}

// AssociateSkillWithEvent associates a skill with an event (for test setup)
func (r *SkillRepository) AssociateSkillWithEvent(skillID, eventID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Add to eventSkills map
	if r.eventSkills[eventID] == nil {
		r.eventSkills[eventID] = []string{}
	}
	if !containsString(r.eventSkills[eventID], skillID) {
		r.eventSkills[eventID] = append(r.eventSkills[eventID], skillID)
	}

	// Add to skillEvents map
	if r.skillEvents[skillID] == nil {
		r.skillEvents[skillID] = []string{}
	}
	if !containsString(r.skillEvents[skillID], eventID) {
		r.skillEvents[skillID] = append(r.skillEvents[skillID], eventID)
	}
}

// Helper functions
func removeString(slice []string, s string) []string {
	var result []string
	for _, item := range slice {
		if item != s {
			result = append(result, item)
		}
	}
	return result
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
