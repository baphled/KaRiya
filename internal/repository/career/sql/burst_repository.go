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
var _ career_repo.BurstRepository = (*BurstRepository)(nil)

// BurstRepository implements career.BurstRepository using GORM.
type BurstRepository struct {
	db *gorm.DB
}

// NewBurstRepository creates a new SQL burst repository.
//
// Expected:
//   - db must be valid.
//
// Returns:
//   - A fully initialized BurstRepository ready for use.
//
// Side effects:
//   - None.
func NewBurstRepository(db *gorm.DB) *BurstRepository {
	return &BurstRepository{db: db}
}

// Create adds a new burst to the database.
//
// Expected:
//   - burst must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *BurstRepository) Create(ctx context.Context, burst *career.Burst) error {
	if burst.ID == "" {
		burst.ID = uuid.New().String()
	}

	now := time.Now()
	burst.CreatedAt = now
	burst.UpdatedAt = now

	return r.db.WithContext(ctx).Create(careermodel.BurstFromDomain(burst)).Error
}

// GetByID retrieves a burst by its ID.
func (r *BurstRepository) GetByID(ctx context.Context, id string) (*career.Burst, error) {
	var model careermodel.Burst
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, career_repo.ErrBurstNotFound
	}
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
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
func (r *BurstRepository) Update(ctx context.Context, burst *career.Burst) error {
	// Check if record exists first.
	var count int64
	if err := r.db.WithContext(ctx).Model(&careermodel.Burst{}).Where("id = ?", burst.ID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return career_repo.ErrBurstNotFound
	}

	burst.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(careermodel.BurstFromDomain(burst)).Error
}

// Delete removes a burst from the database.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *BurstRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&careermodel.Burst{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return career_repo.ErrBurstNotFound
	}
	return result.Error
}

// List retrieves bursts with optional filtering.
func (r *BurstRepository) List(ctx context.Context, filters career_repo.BurstListFilters) ([]*career.Burst, error) {
	query := r.db.WithContext(ctx).Model(&careermodel.Burst{})
	query = r.applyFilters(query, filters)
	query = r.applySorting(query, filters)
	query = r.applyPagination(query, filters)

	var results []careermodel.Burst
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	bursts := make([]*career.Burst, len(results))
	for i := range results {
		bursts[i] = results[i].ToDomain()
	}
	return bursts, nil
}

// Count returns the number of bursts matching the given filters.
func (r *BurstRepository) Count(ctx context.Context, filters career_repo.BurstListFilters) (int, error) {
	query := r.db.WithContext(ctx).Model(&careermodel.Burst{})
	query = r.applyFilters(query, filters)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *BurstRepository) applyFilters(query *gorm.DB, filters career_repo.BurstListFilters) *gorm.DB {
	if filters.StartDate != nil {
		query = query.Where("created_at >= ?", filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("created_at <= ?", filters.EndDate)
	}
	return query
}

func (r *BurstRepository) applySorting(query *gorm.DB, filters career_repo.BurstListFilters) *gorm.DB {
	sortBy := filters.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	order := "DESC"
	if strings.EqualFold(filters.SortOrder, "asc") {
		order = "ASC"
	}

	switch sortBy {
	case "name":
		return query.Order("name " + order)
	case "event_count":
		// Count commas in event_ids to estimate event count, handling empty strings correctly.
		return query.Order("CASE WHEN event_ids = '' THEN 0 ELSE LENGTH(event_ids) - LENGTH(REPLACE(event_ids, ',', '')) + 1 END " + order)
	default:
		return query.Order("created_at " + order)
	}
}

func (r *BurstRepository) applyPagination(query *gorm.DB, filters career_repo.BurstListFilters) *gorm.DB {
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
