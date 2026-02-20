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
//
// Returns:
//   - A fully initialized FactRepository ready for use.
//
// Side effects:
//   - None.
func NewFactRepository() *FactRepository {
	return &FactRepository{
		facts: make(map[string]*career.Fact),
	}
}

// Create adds a new fact to the in-memory store.
//
// Expected:
//   - fact must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *FactRepository) Create(_ context.Context, fact *career.Fact) error {
	if fact.ID == "" {
		fact.ID = uuid.New().String()
	}

	if err := fact.Validate(); err != nil {
		return err
	}

	now := time.Now()
	fact.CreatedAt = now
	fact.UpdatedAt = now

	return storeNew(&r.mu, r.facts, fact.ID, fact, career_repo.ErrDuplicateFact)
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
//
// Expected:
//   - fact must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
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
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
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
	// Iterate in deterministic order by sorting keys first
	keys := make([]string, 0, len(r.facts))
	for k := range r.facts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		facts = append(facts, r.facts[k])
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
	return filterByDateRange(facts, func(f *career.Fact) time.Time { return f.CreatedAt }, filters.StartDate, filters.EndDate)
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
	return paginate(facts, filters.Offset, filters.Limit)
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
