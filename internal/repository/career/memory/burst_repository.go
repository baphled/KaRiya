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
var _ career_repo.BurstRepository = (*BurstRepository)(nil)

// BurstRepository provides an in-memory implementation of the BurstRepository interface.
type BurstRepository struct {
	bursts map[string]*career.Burst
	mu     sync.RWMutex
}

// NewBurstRepository creates a new in-memory burst repository.
//
// Returns:
//   - A fully initialized BurstRepository ready for use.
//
// Side effects:
//   - None.
func NewBurstRepository() *BurstRepository {
	return &BurstRepository{
		bursts: make(map[string]*career.Burst),
	}
}

// Create adds a new burst to the in-memory store.
//
// Expected:
//   - burst must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *BurstRepository) Create(_ context.Context, burst *career.Burst) error {
	if burst.ID == "" {
		burst.ID = uuid.New().String()
	}

	if err := burst.Validate(); err != nil {
		return err
	}

	now := time.Now()
	burst.CreatedAt = now
	burst.UpdatedAt = now

	return storeNew(&r.mu, r.bursts, burst.ID, burst, career_repo.ErrDuplicateBurst)
}

// GetByID retrieves a burst by its ID.
func (r *BurstRepository) GetByID(_ context.Context, id string) (*career.Burst, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	burst, exists := r.bursts[id]
	if !exists {
		return nil, career_repo.ErrBurstNotFound
	}

	return burst, nil
}

// Update modifies an existing burst.
//
// Expected:
//   - burst must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *BurstRepository) Update(_ context.Context, burst *career.Burst) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validate the burst before updating
	if err := burst.Validate(); err != nil {
		return err
	}

	// Check if burst exists
	if _, exists := r.bursts[burst.ID]; !exists {
		return career_repo.ErrBurstNotFound
	}

	// Update the timestamp
	burst.UpdatedAt = time.Now()

	// Store the updated burst
	r.bursts[burst.ID] = burst

	return nil
}

// Delete removes a burst from the repository.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *BurstRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if burst exists
	if _, exists := r.bursts[id]; !exists {
		return career_repo.ErrBurstNotFound
	}

	// Delete the burst
	delete(r.bursts, id)

	return nil
}

// List retrieves bursts with optional filtering.
func (r *BurstRepository) List(_ context.Context, filters career_repo.BurstListFilters) ([]*career.Burst, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var bursts []*career.Burst
	// Iterate in deterministic order by sorting keys first
	keys := make([]string, 0, len(r.bursts))
	for k := range r.bursts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		bursts = append(bursts, r.bursts[k])
	}

	bursts = r.applyFilters(bursts, filters)
	r.applySorting(bursts, filters)
	return r.applyPagination(bursts, filters), nil
}

func (r *BurstRepository) applyFilters(bursts []*career.Burst, filters career_repo.BurstListFilters) []*career.Burst {
	return filterByDateRange(bursts, func(b *career.Burst) time.Time { return b.CreatedAt }, filters.StartDate, filters.EndDate)
}

const (
	sortFieldCreatedAt  = "created_at"
	sortFieldName       = "name"
	sortFieldEventCount = "event_count"
)

func (r *BurstRepository) applySorting(bursts []*career.Burst, filters career_repo.BurstListFilters) {
	sortBy := filters.SortBy
	if sortBy == "" {
		sortBy = sortFieldCreatedAt
	}

	asc := filters.SortOrder == "asc"

	sort.Slice(bursts, func(i, j int) bool {
		var less bool
		switch sortBy {
		case sortFieldName:
			less = bursts[i].Name < bursts[j].Name
		case sortFieldEventCount:
			less = len(bursts[i].EventIDs) < len(bursts[j].EventIDs)
		default:
			less = bursts[i].CreatedAt.Before(bursts[j].CreatedAt)
		}

		if asc {
			return less
		}
		return !less
	})
}

func (r *BurstRepository) applyPagination(bursts []*career.Burst, filters career_repo.BurstListFilters) []*career.Burst {
	return paginate(bursts, filters.Offset, filters.Limit)
}

// Count returns the total number of bursts matching the given filters.
func (r *BurstRepository) Count(ctx context.Context, filters career_repo.BurstListFilters) (int, error) {
	bursts, err := r.List(ctx, filters)
	if err != nil {
		return 0, err
	}
	return len(bursts), nil
}
