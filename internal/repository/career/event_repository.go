package career

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

var (
	// ErrEventNotFound is returned when a requested event cannot be found.
	ErrEventNotFound = errors.New("career event not found")

	// ErrDuplicateEvent is returned when attempting to create an event that already exists.
	ErrDuplicateEvent = errors.New("career event already exists")
)

// EventRepository defines the interface for career event persistence.
type EventRepository interface {
	// Create adds a new career event to the repository.
	Create(ctx context.Context, event *career.Event) error

	// GetByID retrieves a specific career event by its unique identifier.
	GetByID(ctx context.Context, id string) (*career.Event, error)

	// Update modifies an existing career event.
	Update(ctx context.Context, event *career.Event) error

	// Delete removes a career event from the repository.
	Delete(ctx context.Context, id string) error

	// List retrieves career events with optional filtering.
	List(ctx context.Context, filters EventListFilters) ([]*career.Event, error)

	// Count returns the total number of career events matching the given filters.
	Count(ctx context.Context, filters EventListFilters) (int, error)

	// LinkSkill creates an association between an event and a skill.
	// This is used by the skill inference service to record detected skills.
	LinkSkill(ctx context.Context, eventID string, skillID string) error
}

// EventListFilters provides flexible filtering options for career events.
type EventListFilters struct {
	// Tags to filter events by.
	Tags []string

	// Date range filters.
	StartDate *time.Time
	EndDate   *time.Time

	// Pagination.
	Offset int
	Limit  int

	// Sorting.
	SortBy    string
	SortOrder string
}
