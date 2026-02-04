//go:generate mockgen -destination=../../testutil/mocks/repository/burst_repository_mock.go -package=mockrepo github.com/baphled/kariya/internal/repository/career BurstRepository

package career

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

var (
	// ErrBurstNotFound is returned when a requested burst cannot be found.
	ErrBurstNotFound = errors.New("burst not found")

	// ErrDuplicateBurst is returned when attempting to create a burst that already exists.
	ErrDuplicateBurst = errors.New("burst already exists")
)

// BurstRepository defines the interface for burst persistence.
type BurstRepository interface {
	// Create adds a new burst to the repository.
	Create(ctx context.Context, burst *career.Burst) error

	// GetByID retrieves a specific burst by its unique identifier.
	GetByID(ctx context.Context, id string) (*career.Burst, error)

	// Update modifies an existing burst.
	Update(ctx context.Context, burst *career.Burst) error

	// Delete removes a burst from the repository.
	Delete(ctx context.Context, id string) error

	// List retrieves bursts with optional filtering.
	List(ctx context.Context, filters BurstListFilters) ([]*career.Burst, error)

	// Count returns the total number of bursts matching the given filters.
	Count(ctx context.Context, filters BurstListFilters) (int, error)
}

// BurstListFilters provides flexible filtering options for bursts.
type BurstListFilters struct {
	// Date range filters.
	StartDate *time.Time
	EndDate   *time.Time

	// Pagination.
	Offset int
	Limit  int

	// Sorting.
	SortBy    string // "name", "created_at", "event_count"
	SortOrder string // "asc", "desc"
}
