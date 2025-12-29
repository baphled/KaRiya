package career

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/google/uuid"
)

var (
	// ErrEventNotFound is returned when a requested event cannot be found
	ErrEventNotFound = errors.New("career event not found")

	// ErrDuplicateEvent is returned when attempting to create an event that already exists
	ErrDuplicateEvent = errors.New("career event already exists")
)

// Repository defines the interface for career event persistence
type Repository interface {
	// Create adds a new career event to the repository
	Create(ctx context.Context, event *career.CareerEvent) error

	// GetByID retrieves a specific career event by its unique identifier
	GetByID(ctx context.Context, id string) (*career.CareerEvent, error)

	// Update modifies an existing career event
	Update(ctx context.Context, event *career.CareerEvent) error

	// Delete removes a career event from the repository
	Delete(ctx context.Context, id string) error

	// List retrieves career events with optional filtering
	List(ctx context.Context, filters ListFilters) ([]*career.CareerEvent, error)

	// Count returns the total number of career events matching the given filters
	Count(ctx context.Context, filters ListFilters) (int, error)
}

// ListFilters provides flexible filtering options for career events
type ListFilters struct {
	// Tags to filter events by
	Tags []string

	// Date range filters
	StartDate *time.Time
	EndDate   *time.Time

	// Pagination
	Offset int
	Limit  int

	// Sorting
	SortBy    string
	SortOrder string
}

// MemoryRepository provides an in-memory implementation of the Repository interface
// This is primarily useful for testing and development
type MemoryRepository struct {
	events map[string]*career.CareerEvent
	mu     sync.RWMutex
}

// NewMemoryRepository creates a new in-memory repository
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		events: make(map[string]*career.CareerEvent),
	}
}

// Create adds a new career event to the in-memory store
func (r *MemoryRepository) Create(ctx context.Context, event *career.CareerEvent) error {
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
		return ErrDuplicateEvent
	}

	// Set timestamps
	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now

	// Store the event
	r.events[event.ID] = event
	return nil
}

// GetByID retrieves a career event by its ID
func (r *MemoryRepository) GetByID(ctx context.Context, id string) (*career.CareerEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	event, exists := r.events[id]
	if !exists {
		return nil, ErrEventNotFound
	}

	return event, nil
}

// Update modifies an existing career event
func (r *MemoryRepository) Update(ctx context.Context, event *career.CareerEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validate the event
	if err := event.Validate(); err != nil {
		return err
	}

	// Check if event exists
	if _, exists := r.events[event.ID]; !exists {
		return ErrEventNotFound
	}

	// Update timestamp
	event.UpdatedAt = time.Now()

	// Update the event
	r.events[event.ID] = event
	return nil
}

// Delete removes a career event from the repository
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.events[id]; !exists {
		return ErrEventNotFound
	}

	delete(r.events, id)
	return nil
}

// List retrieves career events with optional filtering
func (r *MemoryRepository) List(ctx context.Context, filters ListFilters) ([]*career.CareerEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*career.CareerEvent

	for _, event := range r.events {
		// Apply tag filter
		if len(filters.Tags) > 0 {
			if !containsAnyTag(event.Tags, filters.Tags) {
				continue
			}
		}

		// Apply date range filter
		if filters.StartDate != nil && event.Date.Before(*filters.StartDate) {
			continue
		}
		if filters.EndDate != nil && event.Date.After(*filters.EndDate) {
			continue
		}

		filtered = append(filtered, event)
	}

	// Apply sorting
	sort.Slice(filtered, func(i, j int) bool {
		switch filters.SortBy {
		case "date":
			if filters.SortOrder == "desc" {
				return filtered[i].Date.After(filtered[j].Date)
			}
			return filtered[i].Date.Before(filtered[j].Date)
		default:
			return filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
		}
	})

	// Apply pagination
	start := filters.Offset
	end := start + filters.Limit

	// If Limit is 0, return all results (no pagination limit)
	if filters.Limit == 0 {
		end = len(filtered)
	}

	if start > len(filtered) {
		return []*career.CareerEvent{}, nil
	}

	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], nil
}

// Count returns the number of events matching the filters
func (r *MemoryRepository) Count(ctx context.Context, filters ListFilters) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var count int
	for _, event := range r.events {
		// Apply tag filter
		if len(filters.Tags) > 0 {
			if !containsAnyTag(event.Tags, filters.Tags) {
				continue
			}
		}

		// Apply date range filter
		if filters.StartDate != nil && event.Date.Before(*filters.StartDate) {
			continue
		}
		if filters.EndDate != nil && event.Date.After(*filters.EndDate) {
			continue
		}

		count++
	}

	return count, nil
}

// Helper function to check if any tags match
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
