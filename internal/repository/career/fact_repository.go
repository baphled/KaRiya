//go:generate mockgen -destination=../../testutil/mocks/repository/fact_repository_mock.go -package=mockrepo github.com/baphled/kariya/internal/repository/career FactRepository

package career

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

var (
	// ErrFactNotFound is returned when a requested fact cannot be found.
	ErrFactNotFound = errors.New("fact not found")

	// ErrDuplicateFact is returned when attempting to create a fact that already exists.
	ErrDuplicateFact = errors.New("fact already exists")
)

// FactRepository defines the interface for fact persistence.
type FactRepository interface {
	// Create adds a new fact to the repository.
	Create(ctx context.Context, fact *career.Fact) error

	// GetByID retrieves a specific fact by its unique identifier.
	GetByID(ctx context.Context, id string) (*career.Fact, error)

	// Update modifies an existing fact.
	Update(ctx context.Context, fact *career.Fact) error

	// Delete removes a fact from the repository.
	Delete(ctx context.Context, id string) error

	// List retrieves facts with optional filtering.
	List(ctx context.Context, filters FactListFilters) ([]*career.Fact, error)

	// Count returns the total number of facts matching the given filters.
	Count(ctx context.Context, filters FactListFilters) (int, error)

	// GetBySourceEventID retrieves all facts for a specific event.
	GetBySourceEventID(ctx context.Context, eventID string) ([]*career.Fact, error)

	// GetBySourceBurstID retrieves all facts for a specific burst.
	GetBySourceBurstID(ctx context.Context, burstID string) ([]*career.Fact, error)
}

// FactListFilters provides flexible filtering options for facts.
type FactListFilters struct {
	// CompetencyCategory to filter facts by.
	CompetencyCategory string

	// RoleFit to filter facts by.
	RoleFit string

	// AudienceRelevance to filter facts by.
	AudienceRelevance string

	// Date range filters.
	StartDate *time.Time
	EndDate   *time.Time

	// Pagination.
	Offset int
	Limit  int

	// Sorting.
	SortBy    string // "created_at", "text"
	SortOrder string // "asc", "desc"
}
