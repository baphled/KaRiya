// Package sql provides SQL-backed repository implementations for career domain entities.
package sql

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careermodel "github.com/baphled/kariya/internal/model/career"
	career_repo "github.com/baphled/kariya/internal/repository/career"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Compile-time interface check.
var _ career_repo.FactRepository = (*FactRepository)(nil)

// FactRepository implements career.FactRepository using GORM.
type FactRepository struct {
	db *gorm.DB
}

// NewFactRepository creates a new SQL fact repository.
//
// Expected:
//   - db must be valid.
//
// Returns:
//   - A fully initialized FactRepository ready for use.
//
// Side effects:
//   - None.
func NewFactRepository(db *gorm.DB) *FactRepository {
	return &FactRepository{db: db}
}

// Create adds a new fact to the database.
//
// Expected:
//   - fact must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *FactRepository) Create(ctx context.Context, fact *career.Fact) error {
	if fact.ID == "" {
		fact.ID = uuid.New().String()
	}

	now := time.Now()
	fact.CreatedAt = now
	fact.UpdatedAt = now

	return r.db.WithContext(ctx).Create(careermodel.FactFromDomain(fact)).Error
}

// GetByID retrieves a fact by its ID.
func (r *FactRepository) GetByID(ctx context.Context, id string) (*career.Fact, error) {
	var model careermodel.Fact
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, career_repo.ErrFactNotFound
	}
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
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
func (r *FactRepository) Update(ctx context.Context, fact *career.Fact) error {
	// Check if record exists first.
	var count int64
	if err := r.db.WithContext(ctx).Model(&careermodel.Fact{}).Where("id = ?", fact.ID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return career_repo.ErrFactNotFound
	}

	fact.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(careermodel.FactFromDomain(fact)).Error
}

// Delete removes a fact from the database.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *FactRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&careermodel.Fact{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return career_repo.ErrFactNotFound
	}
	return result.Error
}

// List retrieves facts with optional filtering.
func (r *FactRepository) List(ctx context.Context, filters career_repo.FactListFilters) ([]*career.Fact, error) {
	query := r.db.WithContext(ctx).Model(&careermodel.Fact{})
	query = r.applyFilters(query, filters)
	query = r.applySorting(query, filters)
	query = r.applyPagination(query, filters)

	var results []careermodel.Fact
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	facts := make([]*career.Fact, len(results))
	for i := range results {
		facts[i] = results[i].ToDomain()
	}
	return facts, nil
}

// Count returns the number of facts matching the given filters.
func (r *FactRepository) Count(ctx context.Context, filters career_repo.FactListFilters) (int, error) {
	query := r.db.WithContext(ctx).Model(&careermodel.Fact{})
	query = r.applyFilters(query, filters)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// GetBySourceEventID retrieves all facts for a specific event.
func (r *FactRepository) GetBySourceEventID(ctx context.Context, eventID string) ([]*career.Fact, error) {
	var results []careermodel.Fact
	err := r.db.WithContext(ctx).Where("source_event_id = ?", eventID).Find(&results).Error
	if err != nil {
		return nil, err
	}

	facts := make([]*career.Fact, len(results))
	for i := range results {
		facts[i] = results[i].ToDomain()
	}
	return facts, nil
}

// GetBySourceBurstID retrieves all facts for a specific burst.
func (r *FactRepository) GetBySourceBurstID(ctx context.Context, burstID string) ([]*career.Fact, error) {
	var results []careermodel.Fact
	err := r.db.WithContext(ctx).Where("source_burst_id = ?", burstID).Find(&results).Error
	if err != nil {
		return nil, err
	}

	facts := make([]*career.Fact, len(results))
	for i := range results {
		facts[i] = results[i].ToDomain()
	}
	return facts, nil
}

func (r *FactRepository) applyFilters(query *gorm.DB, filters career_repo.FactListFilters) *gorm.DB {
	if filters.CompetencyCategory != "" {
		query = query.Where("competencies LIKE ?", "%"+filters.CompetencyCategory+"%")
	}
	if filters.RoleFit != "" {
		query = query.Where("role_fit = ?", filters.RoleFit)
	}
	if filters.AudienceRelevance != "" {
		query = query.Where("audience_relevance LIKE ?", "%"+filters.AudienceRelevance+"%")
	}
	if filters.StartDate != nil {
		query = query.Where("created_at >= ?", filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("created_at <= ?", filters.EndDate)
	}
	return query
}

func (r *FactRepository) applySorting(query *gorm.DB, filters career_repo.FactListFilters) *gorm.DB {
	sortBy := filters.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	order := "DESC"
	if strings.EqualFold(filters.SortOrder, "asc") {
		order = "ASC"
	}

	switch sortBy {
	case "text":
		return query.Order("text " + order)
	default:
		return query.Order("created_at " + order)
	}
}

func (r *FactRepository) applyPagination(query *gorm.DB, filters career_repo.FactListFilters) *gorm.DB {
	limit := filters.Limit
	if limit == 0 {
		limit = defaultPaginationLimit
	}
	query = query.Limit(limit)
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}
	return query
}
