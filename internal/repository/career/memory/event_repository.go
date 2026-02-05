package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	career_repo "github.com/baphled/kariya/internal/repository/career"
	"github.com/google/uuid"
)

// Compile-time interface check.
var _ career_repo.EventRepository = (*EventRepository)(nil)

// EventRepository provides an in-memory implementation of the EventRepository interface.
// This is primarily useful for testing and development.
type EventRepository struct {
	events    map[string]*career.Event
	skillRepo *SkillRepository
	mu        sync.RWMutex
}

// NewEventRepository creates a new in-memory event repository.
//
// Returns:
//   - A fully initialized EventRepository ready for use.
//
// Side effects:
//   - None.
func NewEventRepository() *EventRepository {
	return &EventRepository{
		events: make(map[string]*career.Event),
	}
}

// SetSkillRepository sets the skill repository for cross-repository association sync.
//
// Expected:
//   - skillrepository must be valid.
//
// Side effects:
//   - None.
func (r *EventRepository) SetSkillRepository(repo *SkillRepository) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.skillRepo = repo
}

// Create adds a new career event to the in-memory store.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *EventRepository) Create(_ context.Context, event *career.Event) error {
	if err := event.Validate(); err != nil {
		return err
	}

	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now

	return storeNew(&r.mu, r.events, event.ID, event, career_repo.ErrDuplicateEvent)
}

// GetByID retrieves a career event by its ID.
func (r *EventRepository) GetByID(_ context.Context, id string) (*career.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	event, exists := r.events[id]
	if !exists {
		return nil, career_repo.ErrEventNotFound
	}

	return event, nil
}

// Update modifies an existing career event.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *EventRepository) Update(_ context.Context, event *career.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validate the event
	if err := event.Validate(); err != nil {
		return err
	}

	// Check if event exists
	if _, exists := r.events[event.ID]; !exists {
		return career_repo.ErrEventNotFound
	}

	// Update timestamp
	event.UpdatedAt = time.Now()

	// Update the event
	r.events[event.ID] = event
	return nil
}

// Delete removes a career event from the repository.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *EventRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.events[id]; !exists {
		return career_repo.ErrEventNotFound
	}

	delete(r.events, id)
	return nil
}

// List retrieves career events with optional filtering.
func (r *EventRepository) List(_ context.Context, filters career_repo.EventListFilters) ([]*career.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var events []*career.Event
	for _, event := range r.events {
		events = append(events, event)
	}

	events = r.applyFilters(events, filters)
	r.applySorting(events, filters)
	return r.applyPagination(events, filters), nil
}

func (r *EventRepository) applyFilters(events []*career.Event, filters career_repo.EventListFilters) []*career.Event {
	if len(filters.Tags) > 0 {
		var tagged []*career.Event
		for _, event := range events {
			if containsAnyTag(event.Tags, filters.Tags) {
				tagged = append(tagged, event)
			}
		}
		events = tagged
	}

	return filterByDateRange(events, func(e *career.Event) time.Time { return e.Date }, filters.StartDate, filters.EndDate)
}

func (r *EventRepository) applySorting(events []*career.Event, filters career_repo.EventListFilters) {
	desc := filters.SortOrder == "desc"

	sort.Slice(events, func(i, j int) bool {
		switch filters.SortBy {
		case "date":
			if events[i].Date.Equal(events[j].Date) {
				if desc {
					return events[i].CreatedAt.After(events[j].CreatedAt)
				}
				return events[i].CreatedAt.Before(events[j].CreatedAt)
			}
			if desc {
				return events[i].Date.After(events[j].Date)
			}
			return events[i].Date.Before(events[j].Date)
		default:
			if desc {
				return events[i].CreatedAt.After(events[j].CreatedAt)
			}
			return events[i].CreatedAt.Before(events[j].CreatedAt)
		}
	})
}

func (r *EventRepository) applyPagination(events []*career.Event, filters career_repo.EventListFilters) []*career.Event {
	return paginate(events, filters.Offset, filters.Limit)
}

// Count returns the number of events matching the filters.
func (r *EventRepository) Count(_ context.Context, filters career_repo.EventListFilters) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var events []*career.Event
	for _, event := range r.events {
		events = append(events, event)
	}

	return len(r.applyFilters(events, filters)), nil
}

// LinkSkill creates an association between an event and a skill.
//
// Expected:
//   - Must be a valid string.
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *EventRepository) LinkSkill(_ context.Context, eventID string, skillID string) error {
	r.mu.Lock()

	event, exists := r.events[eventID]
	if !exists {
		r.mu.Unlock()
		return career_repo.ErrEventNotFound
	}

	alreadyLinked := false
	for _, existingSkillID := range event.Skills {
		if existingSkillID == skillID {
			alreadyLinked = true
			break
		}
	}

	if !alreadyLinked {
		event.Skills = append(event.Skills, skillID)
		event.UpdatedAt = time.Now()
	}

	skillRepo := r.skillRepo
	r.mu.Unlock()

	if skillRepo != nil {
		skillRepo.AssociateSkillWithEvent(skillID, eventID)
	}

	return nil
}

// UnlinkSkill removes an association between an event and a skill.
//
// Expected:
//   - eventID must be a valid string identifier for an existing event.
//   - skillID must be a valid string identifier for an existing skill.
//
// Returns:
//   - An error value.
//
// Side effects:
//   - Removes the skill from the event's skill list.
func (r *EventRepository) UnlinkSkill(_ context.Context, eventID string, skillID string) error {
	r.mu.Lock()

	event, exists := r.events[eventID]
	if !exists {
		r.mu.Unlock()
		return career_repo.ErrEventNotFound
	}

	newSkills := make([]string, 0, len(event.Skills))
	found := false
	for _, existingSkillID := range event.Skills {
		if existingSkillID == skillID {
			found = true
			continue
		}
		newSkills = append(newSkills, existingSkillID)
	}

	if found {
		event.Skills = newSkills
		event.UpdatedAt = time.Now()
	}

	skillRepo := r.skillRepo
	r.mu.Unlock()

	if skillRepo != nil && found {
		skillRepo.DisassociateSkillFromEvent(skillID, eventID)
	}

	return nil
}

// containsAnyTag checks if any tags match.
func containsAnyTag(eventTags, filterTags []string) bool {
	tagSet := make(map[string]bool)
	for _, tag := range eventTags {
		tagSet[tag] = true
	}

	for _, filterTag := range filterTags {
		if tagSet[filterTag] {
			return true
		}
	}

	return false
}
