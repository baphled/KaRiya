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
	// ErrFactNotFound is returned when a requested fact cannot be found
	ErrFactNotFound = errors.New("fact not found")

	// ErrDuplicateFact is returned when attempting to create a fact that already exists
	ErrDuplicateFact = errors.New("fact already exists")
)

// FactRepository defines the interface for fact persistence
//
//nolint:interfacebloat // Repository interfaces require CRUD + query methods
type FactRepository interface {
	// Create adds a new fact to the repository
	Create(ctx context.Context, fact *career.Fact) error

	// GetByID retrieves a specific fact by its unique identifier
	GetByID(ctx context.Context, id string) (*career.Fact, error)

	// Update modifies an existing fact
	Update(ctx context.Context, fact *career.Fact) error

	// Delete removes a fact from the repository
	Delete(ctx context.Context, id string) error

	// List retrieves facts with optional filtering
	List(ctx context.Context, filters FactListFilters) ([]*career.Fact, error)

	// Count returns the total number of facts matching the given filters
	Count(ctx context.Context, filters FactListFilters) (int, error)

	// GetBySourceEventID retrieves all facts for a specific event
	GetBySourceEventID(ctx context.Context, eventID string) ([]*career.Fact, error)

	// GetBySourceBurstID retrieves all facts for a specific burst
	GetBySourceBurstID(ctx context.Context, burstID string) ([]*career.Fact, error)
}

// FactListFilters provides flexible filtering options for facts
type FactListFilters struct {
	// CompetencyCategory to filter facts by
	CompetencyCategory string

	// RoleFit to filter facts by
	RoleFit string

	// AudienceRelevance to filter facts by
	AudienceRelevance string

	// Date range filters
	StartDate *time.Time
	EndDate   *time.Time

	// Pagination
	Offset int
	Limit  int

	// Sorting
	SortBy    string // "created_at", "text"
	SortOrder string // "asc", "desc"
}

// MemoryFactRepository provides an in-memory implementation of the FactRepository interface
type MemoryFactRepository struct {
	facts map[string]*career.Fact
	mu    sync.RWMutex
}

// NewMemoryFactRepository creates a new in-memory fact repository
func NewMemoryFactRepository() *MemoryFactRepository {
	return &MemoryFactRepository{
		facts: make(map[string]*career.Fact),
	}
}

// Create adds a new fact to the in-memory store
func (r *MemoryFactRepository) Create(ctx context.Context, fact *career.Fact) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Generate a unique ID if not provided
	if fact.ID == "" {
		fact.ID = uuid.New().String()
	}

	// Validate the fact
	if err := fact.Validate(); err != nil {
		return err
	}

	// Check for duplicates
	if _, exists := r.facts[fact.ID]; exists {
		return ErrDuplicateFact
	}

	// Set timestamps
	now := time.Now()
	fact.CreatedAt = now
	fact.UpdatedAt = now

	// Store the fact
	r.facts[fact.ID] = fact

	return nil
}

// GetByID retrieves a fact by its ID.
func (r *MemoryFactRepository) GetByID(_ context.Context, id string) (*career.Fact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	fact, exists := r.facts[id]
	if !exists {
		return nil, ErrFactNotFound
	}

	return fact, nil
}

// Update modifies an existing fact
func (r *MemoryFactRepository) Update(ctx context.Context, fact *career.Fact) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validate the fact
	if err := fact.Validate(); err != nil {
		return err
	}

	// Check if fact exists
	if _, exists := r.facts[fact.ID]; !exists {
		return ErrFactNotFound
	}

	// Update the timestamp
	fact.UpdatedAt = time.Now()

	// Store the updated fact
	r.facts[fact.ID] = fact

	return nil
}

// Delete removes a fact from the repository
func (r *MemoryFactRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if fact exists
	if _, exists := r.facts[id]; !exists {
		return ErrFactNotFound
	}

	// Delete the fact
	delete(r.facts, id)

	return nil
}

// List retrieves facts with optional filtering
func (r *MemoryFactRepository) List(ctx context.Context, filters FactListFilters) ([]*career.Fact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Collect all facts
	var facts []*career.Fact
	for _, fact := range r.facts {
		facts = append(facts, fact)
	}

	// Apply competency category filter
	if filters.CompetencyCategory != "" {
		var filtered []*career.Fact
		for _, fact := range facts {
			for _, cat := range fact.CompetencyCategories {
				if cat == filters.CompetencyCategory {
					filtered = append(filtered, fact)
					break
				}
			}
		}
		facts = filtered
	}

	// Apply role fit filter
	if filters.RoleFit != "" {
		var filtered []*career.Fact
		for _, fact := range facts {
			if string(fact.RoleFit) == filters.RoleFit {
				filtered = append(filtered, fact)
			}
		}
		facts = filtered
	}

	// Apply audience relevance filter
	if filters.AudienceRelevance != "" {
		var filtered []*career.Fact
		for _, fact := range facts {
			for _, aud := range fact.AudienceRelevance {
				if aud == filters.AudienceRelevance {
					filtered = append(filtered, fact)
					break
				}
			}
		}
		facts = filtered
	}

	// Apply date range filters
	if filters.StartDate != nil || filters.EndDate != nil {
		var filtered []*career.Fact
		for _, fact := range facts {
			if filters.StartDate != nil && fact.CreatedAt.Before(*filters.StartDate) {
				continue
			}
			if filters.EndDate != nil && fact.CreatedAt.After(*filters.EndDate) {
				continue
			}
			filtered = append(filtered, fact)
		}
		facts = filtered
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

	sort.Slice(facts, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "text":
			less = facts[i].Text < facts[j].Text
		case "created_at":
			fallthrough
		default:
			less = facts[i].CreatedAt.Before(facts[j].CreatedAt)
		}

		if sortOrder == "asc" {
			return less
		}
		return !less
	})

	// Apply pagination
	offset := filters.Offset
	limit := filters.Limit

	if offset > len(facts) {
		return []*career.Fact{}, nil
	}

	end := offset + limit
	if end > len(facts) {
		end = len(facts)
	}

	if limit == 0 {
		return facts[offset:], nil
	}

	return facts[offset:end], nil
}

// Count returns the total number of facts matching the given filters
func (r *MemoryFactRepository) Count(ctx context.Context, filters FactListFilters) (int, error) {
	facts, err := r.List(ctx, filters)
	if err != nil {
		return 0, err
	}
	return len(facts), nil
}

// GetBySourceEventID retrieves all facts for a specific event
func (r *MemoryFactRepository) GetBySourceEventID(ctx context.Context, eventID string) ([]*career.Fact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var facts []*career.Fact
	for _, fact := range r.facts {
		if fact.SourceEventID == eventID {
			facts = append(facts, fact)
		}
	}

	return facts, nil
}

// GetBySourceBurstID retrieves all facts for a specific burst
func (r *MemoryFactRepository) GetBySourceBurstID(ctx context.Context, burstID string) ([]*career.Fact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var facts []*career.Fact
	for _, fact := range r.facts {
		if fact.SourceBurstID == burstID {
			facts = append(facts, fact)
		}
	}

	return facts, nil
}
