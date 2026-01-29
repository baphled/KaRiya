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
func NewBurstRepository() *BurstRepository {
	return &BurstRepository{
		bursts: make(map[string]*career.Burst),
	}
}

// Create adds a new burst to the in-memory store.
func (r *BurstRepository) Create(_ context.Context, burst *career.Burst) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Generate a unique ID if not provided
	if burst.ID == "" {
		burst.ID = uuid.New().String()
	}

	// Validate the burst before storing
	if err := burst.Validate(); err != nil {
		return err
	}

	// Check for duplicates
	if _, exists := r.bursts[burst.ID]; exists {
		return career_repo.ErrDuplicateBurst
	}

	// Set timestamps
	now := time.Now()
	burst.CreatedAt = now
	burst.UpdatedAt = now

	// Store the burst
	r.bursts[burst.ID] = burst

	return nil
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

	// Collect all bursts
	var bursts []*career.Burst
	for _, burst := range r.bursts {
		bursts = append(bursts, burst)
	}

	// Apply filters
	if filters.StartDate != nil || filters.EndDate != nil {
		var filtered []*career.Burst
		for _, burst := range bursts {
			if filters.StartDate != nil && burst.CreatedAt.Before(*filters.StartDate) {
				continue
			}
			if filters.EndDate != nil && burst.CreatedAt.After(*filters.EndDate) {
				continue
			}
			filtered = append(filtered, burst)
		}
		bursts = filtered
	}

	// Apply sorting
	sortBy := filters.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	sortOrder := filters.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}

	sort.Slice(bursts, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "name":
			less = bursts[i].Name < bursts[j].Name
		case "event_count":
			less = len(bursts[i].EventIDs) < len(bursts[j].EventIDs)
		case "created_at":
			fallthrough
		default:
			less = bursts[i].CreatedAt.Before(bursts[j].CreatedAt)
		}

		if sortOrder == "asc" {
			return less
		}
		return !less
	})

	// Apply pagination
	offset := filters.Offset
	limit := filters.Limit

	if offset > len(bursts) {
		return []*career.Burst{}, nil
	}

	end := offset + limit
	if end > len(bursts) {
		end = len(bursts)
	}

	if limit == 0 {
		return bursts[offset:], nil
	}

	return bursts[offset:end], nil
}

// Count returns the total number of bursts matching the given filters.
func (r *BurstRepository) Count(ctx context.Context, filters career_repo.BurstListFilters) (int, error) {
	bursts, err := r.List(ctx, filters)
	if err != nil {
		return 0, err
	}
	return len(bursts), nil
}
