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
	events map[string]*career.CareerEvent
	mu     sync.RWMutex
}

// NewEventRepository creates a new in-memory event repository.
func NewEventRepository() *EventRepository {
	return &EventRepository{
		events: make(map[string]*career.CareerEvent),
	}
}

// Create adds a new career event to the in-memory store.
func (r *EventRepository) Create(_ context.Context, event *career.CareerEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validate the event before storing
	if err := event.Validate(); err != nil {
		return err
	}

	// Generate a unique ID if not provided
	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	// Check for duplicates
	if _, exists := r.events[event.ID]; exists {
		return career_repo.ErrDuplicateEvent
	}

	// Set timestamps
	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now

	// Store the event
	r.events[event.ID] = event
	return nil
}

// GetByID retrieves a career event by its ID.
func (r *EventRepository) GetByID(_ context.Context, id string) (*career.CareerEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	event, exists := r.events[id]
	if !exists {
		return nil, career_repo.ErrEventNotFound
	}

	return event, nil
}

// Update modifies an existing career event.
func (r *EventRepository) Update(_ context.Context, event *career.CareerEvent) error {
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
func (r *EventRepository) List(_ context.Context, filters career_repo.EventListFilters) ([]*career.CareerEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var events []*career.CareerEvent
	for _, event := range r.events {
		events = append(events, event)
	}

	events = r.applyFilters(events, filters)
	r.applySorting(events, filters)
	return r.applyPagination(events, filters), nil
}

func (r *EventRepository) applyFilters(events []*career.CareerEvent, filters career_repo.EventListFilters) []*career.CareerEvent {
	var filtered []*career.CareerEvent
	for _, event := range events {
		if len(filters.Tags) > 0 && !containsAnyTag(event.Tags, filters.Tags) {
			continue
		}
		if filters.StartDate != nil && event.Date.Before(*filters.StartDate) {
			continue
		}
		if filters.EndDate != nil && event.Date.After(*filters.EndDate) {
			continue
		}
		filtered = append(filtered, event)
	}
	return filtered
}

func (r *EventRepository) applySorting(events []*career.CareerEvent, filters career_repo.EventListFilters) {
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

func (r *EventRepository) applyPagination(events []*career.CareerEvent, filters career_repo.EventListFilters) []*career.CareerEvent {
	start := filters.Offset
	if start > len(events) {
		return []*career.CareerEvent{}
	}

	if filters.Limit == 0 {
		return events[start:]
	}

	end := start + filters.Limit
	if end > len(events) {
		end = len(events)
	}
	return events[start:end]
}

// Count returns the number of events matching the filters.
func (r *EventRepository) Count(_ context.Context, filters career_repo.EventListFilters) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var events []*career.CareerEvent
	for _, event := range r.events {
		events = append(events, event)
	}

	return len(r.applyFilters(events, filters)), nil
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
