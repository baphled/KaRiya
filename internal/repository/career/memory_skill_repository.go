package career

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/google/uuid"
)

// MemorySkillRepository provides an in-memory implementation of the SkillRepository interface
type MemorySkillRepository struct {
	skills      map[string]*career.Skill
	eventSkills map[string][]string // eventID -> []skillID
	skillEvents map[string][]string // skillID -> []eventID
	eventRepo   Repository          // Reference to event repository for GetEventsUsingSkill
	mu          sync.RWMutex
}

// NewMemorySkillRepository creates a new in-memory skill repository
func NewMemorySkillRepository() *MemorySkillRepository {
	return &MemorySkillRepository{
		skills:      make(map[string]*career.Skill),
		eventSkills: make(map[string][]string),
		skillEvents: make(map[string][]string),
	}
}

// SetEventRepository sets the event repository for cross-repository queries
func (r *MemorySkillRepository) SetEventRepository(repo Repository) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.eventRepo = repo
}

// Create adds a new skill to the in-memory store
func (r *MemorySkillRepository) Create(ctx context.Context, skill *career.Skill) error {
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
			return ErrDuplicateSkill
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

// GetByID retrieves a specific skill by its unique identifier
func (r *MemorySkillRepository) GetByID(ctx context.Context, id string) (*career.Skill, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	skill, exists := r.skills[id]
	if !exists {
		return nil, ErrSkillNotFound
	}

	return skill, nil
}

// GetByName retrieves a skill by its name (case-sensitive)
func (r *MemorySkillRepository) GetByName(ctx context.Context, name string) (*career.Skill, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, skill := range r.skills {
		if skill.Name == name {
			return skill, nil
		}
	}

	return nil, ErrSkillNotFound
}

// List retrieves skills with optional filtering
func (r *MemorySkillRepository) List(ctx context.Context, filters *SkillFilters) ([]*career.Skill, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Get event counts for MinEvents filtering
	eventCounts := make(map[string]int)
	for skillID, eventIDs := range r.skillEvents {
		eventCounts[skillID] = len(eventIDs)
	}

	var result []*career.Skill

	// Collect all skills that match the filters
	for _, skill := range r.skills {
		if filters != nil {
			// Filter by category
			if filters.Category != "" && skill.Category != filters.Category {
				continue
			}

			// Filter by level
			if filters.Level != "" && skill.Level != filters.Level {
				continue
			}

			// Filter by minimum events
			if filters.MinEvents > 0 && eventCounts[skill.ID] < filters.MinEvents {
				continue
			}
		}

		result = append(result, skill)
	}

	// Determine sort order and field
	sortOrder := "asc"
	sortBy := "name"
	if filters != nil {
		if filters.SortOrder == "desc" {
			sortOrder = "desc"
		}
		if filters.SortBy != "" {
			sortBy = filters.SortBy
		}
	}

	// Sort based on sortBy and sortOrder
	switch sortBy {
	case "name":
		sort.Slice(result, func(i, j int) bool {
			cmp := strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
			if sortOrder == "desc" {
				return !cmp
			}
			return cmp
		})
	case "events":
		sort.Slice(result, func(i, j int) bool {
			countI := eventCounts[result[i].ID]
			countJ := eventCounts[result[j].ID]
			if countI == countJ {
				// Secondary sort by name ascending
				return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
			}
			if sortOrder == "desc" {
				return countI > countJ
			}
			return countI < countJ
		})
	case "last_used":
		// Get last used dates from event repository if available
		lastUsedDates := make(map[string]time.Time)
		if r.eventRepo != nil {
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
		}

		sort.Slice(result, func(i, j int) bool {
			dateI, hasI := lastUsedDates[result[i].ID]
			dateJ, hasJ := lastUsedDates[result[j].ID]

			// Skills without dates should sort last for DESC, first for ASC
			if !hasI && !hasJ {
				// Both have no date, sort by name
				return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
			}
			if !hasI {
				return sortOrder == "asc"
			}
			if !hasJ {
				return sortOrder == "desc"
			}

			// Both have dates
			if dateI.Equal(dateJ) {
				return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
			}
			if sortOrder == "desc" {
				return dateI.After(dateJ)
			}
			return dateI.Before(dateJ)
		})
	case "category":
		sort.Slice(result, func(i, j int) bool {
			cmp := result[i].Category < result[j].Category
			if result[i].Category == result[j].Category {
				// Secondary sort by name
				return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
			}
			if sortOrder == "desc" {
				return !cmp
			}
			return cmp
		})
	default:
		// Default to name sort
		sort.Slice(result, func(i, j int) bool {
			return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
		})
	}

	// Apply pagination
	if filters != nil {
		if filters.Offset > 0 {
			if filters.Offset >= len(result) {
				return []*career.Skill{}, nil
			}
			result = result[filters.Offset:]
		}

		if filters.Limit > 0 && filters.Limit < len(result) {
			result = result[:filters.Limit]
		}
	}

	return result, nil
}

// Update modifies an existing skill
func (r *MemorySkillRepository) Update(ctx context.Context, skill *career.Skill) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if skill exists
	if _, exists := r.skills[skill.ID]; !exists {
		return ErrSkillNotFound
	}

	// Validate the skill
	if err := skill.Validate(); err != nil {
		return err
	}

	// Check for duplicate name (excluding current skill)
	for id, s := range r.skills {
		if id != skill.ID && s.Name == skill.Name {
			return ErrDuplicateSkill
		}
	}

	// Update timestamp
	skill.UpdatedAt = time.Now()

	// Store the updated skill
	r.skills[skill.ID] = skill

	return nil
}

// Delete removes a skill from the repository
func (r *MemorySkillRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if skill exists
	if _, exists := r.skills[id]; !exists {
		return ErrSkillNotFound
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
func (r *MemorySkillRepository) GetByCategory(ctx context.Context, category string) ([]*career.Skill, error) {
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
func (r *MemorySkillRepository) GetSkillsForEvent(ctx context.Context, eventID string) ([]*career.Skill, error) {
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
func (r *MemorySkillRepository) GetEventCountsForSkills(ctx context.Context) (map[string]int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	counts := make(map[string]int)

	for skillID, eventIDs := range r.skillEvents {
		counts[skillID] = len(eventIDs)
	}

	return counts, nil
}

// GetLastUsedForSkills returns a map of skill IDs to their last used dates (from events)
func (r *MemorySkillRepository) GetLastUsedForSkills(ctx context.Context) (map[string]time.Time, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	lastUsed := make(map[string]time.Time)

	// Note: This requires access to events, which we don't have in this simple implementation
	// For now, return empty map. In real usage with SetEventRepository, this would work.
	// The SQLite implementation handles this properly with JOINs.

	return lastUsed, nil
}

// GetEventsUsingSkill returns all events that use a specific skill, ordered by date DESC
func (r *MemorySkillRepository) GetEventsUsingSkill(ctx context.Context, skillID string) ([]*career.CareerEvent, error) {
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
func (r *MemorySkillRepository) AssociateSkillWithEvent(skillID, eventID string) {
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
