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
var _ career_repo.FactRepository = (*FactRepository)(nil)

// FactRepository provides an in-memory implementation of the FactRepository interface.
type FactRepository struct {
	facts map[string]*career.Fact
	mu    sync.RWMutex
}

// NewFactRepository creates a new in-memory fact repository.
func NewFactRepository() *FactRepository {
	return &FactRepository{
		facts: make(map[string]*career.Fact),
	}
}

// Create adds a new fact to the in-memory store.
func (r *FactRepository) Create(_ context.Context, fact *career.Fact) error {
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
		return career_repo.ErrDuplicateFact
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
func (r *FactRepository) GetByID(_ context.Context, id string) (*career.Fact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	fact, exists := r.facts[id]
	if !exists {
		return nil, career_repo.ErrFactNotFound
	}

	return fact, nil
}

// Update modifies an existing fact.
func (r *FactRepository) Update(_ context.Context, fact *career.Fact) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validate the fact
	if err := fact.Validate(); err != nil {
		return err
	}

	// Check if fact exists
	if _, exists := r.facts[fact.ID]; !exists {
		return career_repo.ErrFactNotFound
	}

	// Update the timestamp
	fact.UpdatedAt = time.Now()

	// Store the updated fact
	r.facts[fact.ID] = fact

	return nil
}

// Delete removes a fact from the repository.
func (r *FactRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if fact exists
	if _, exists := r.facts[id]; !exists {
		return career_repo.ErrFactNotFound
	}

	// Delete the fact
	delete(r.facts, id)

	return nil
}

// List retrieves facts with optional filtering.
func (r *FactRepository) List(_ context.Context, filters career_repo.FactListFilters) ([]*career.Fact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var facts []*career.Fact
	for _, fact := range r.facts {
		facts = append(facts, fact)
	}

	facts = r.applyFilters(facts, filters)
	r.applySorting(facts, filters)
	return r.applyPagination(facts, filters), nil
}

func (r *FactRepository) applyFilters(facts []*career.Fact, filters career_repo.FactListFilters) []*career.Fact {
	if filters.CompetencyCategory != "" {
		facts = filterBySliceContains(facts, func(f *career.Fact) []string {
			return f.CompetencyCategories
		}, filters.CompetencyCategory)
	}

	if filters.RoleFit != "" {
		var filtered []*career.Fact
		for _, fact := range facts {
			if string(fact.RoleFit) == filters.RoleFit {
				filtered = append(filtered, fact)
			}
		}
		facts = filtered
	}

	if filters.AudienceRelevance != "" {
		facts = filterBySliceContains(facts, func(f *career.Fact) []string {
			return f.AudienceRelevance
		}, filters.AudienceRelevance)
	}

	return r.applyDateFilter(facts, filters)
}

func (r *FactRepository) applyDateFilter(facts []*career.Fact, filters career_repo.FactListFilters) []*career.Fact {
	if filters.StartDate == nil && filters.EndDate == nil {
		return facts
	}

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
	return filtered
}

// filterBySliceContains filters facts where a slice field contains the target value.
func filterBySliceContains(facts []*career.Fact, getSlice func(*career.Fact) []string, target string) []*career.Fact {
	var filtered []*career.Fact
	for _, fact := range facts {
		for _, val := range getSlice(fact) {
			if val == target {
				filtered = append(filtered, fact)
				break
			}
		}
	}
	return filtered
}

func (r *FactRepository) applySorting(facts []*career.Fact, filters career_repo.FactListFilters) {
	sortBy := filters.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	asc := filters.SortOrder == "asc"

	sort.Slice(facts, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "text":
			less = facts[i].Text < facts[j].Text
		default:
			less = facts[i].CreatedAt.Before(facts[j].CreatedAt)
		}

		if asc {
			return less
		}
		return !less
	})
}

func (r *FactRepository) applyPagination(facts []*career.Fact, filters career_repo.FactListFilters) []*career.Fact {
	offset := filters.Offset
	if offset > len(facts) {
		return []*career.Fact{}
	}

	if filters.Limit == 0 {
		return facts[offset:]
	}

	end := offset + filters.Limit
	if end > len(facts) {
		end = len(facts)
	}
	return facts[offset:end]
}

// Count returns the total number of facts matching the given filters.
func (r *FactRepository) Count(ctx context.Context, filters career_repo.FactListFilters) (int, error) {
	facts, err := r.List(ctx, filters)
	if err != nil {
		return 0, err
	}
	return len(facts), nil
}

// GetBySourceEventID retrieves all facts for a specific event.
func (r *FactRepository) GetBySourceEventID(_ context.Context, eventID string) ([]*career.Fact, error) {
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

// GetBySourceBurstID retrieves all facts for a specific burst.
func (r *FactRepository) GetBySourceBurstID(_ context.Context, burstID string) ([]*career.Fact, error) {
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
