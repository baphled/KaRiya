// Package career provides repository implementations for career domain entities.
package career

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/repository/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Compile-time interface check.
var _ BurstRepository = (*Burst)(nil)

// Burst implements BurstRepository using GORM.
type Burst struct {
	db *gorm.DB
}

// NewBurst creates a new burst repository.
func NewBurst(db *gorm.DB) *Burst {
	return &Burst{db: db}
}

// Create adds a new burst to the database.
func (r *Burst) Create(ctx context.Context, burst *career.Burst) error {
	if burst.ID == "" {
		burst.ID = uuid.New().String()
	}

	now := time.Now()
	burst.CreatedAt = now
	burst.UpdatedAt = now

	return r.db.WithContext(ctx).Create(models.BurstFromDomain(burst)).Error
}

// GetByID retrieves a burst by its ID.
func (r *Burst) GetByID(ctx context.Context, id string) (*career.Burst, error) {
	var model models.Burst
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrBurstNotFound
	}
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

// Update modifies an existing burst.
func (r *Burst) Update(ctx context.Context, burst *career.Burst) error {
	// Check if record exists first.
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Burst{}).Where("id = ?", burst.ID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrBurstNotFound
	}

	burst.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(models.BurstFromDomain(burst)).Error
}

// Delete removes a burst from the database.
func (r *Burst) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Burst{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return ErrBurstNotFound
	}
	return result.Error
}

// List retrieves bursts with optional filtering.
func (r *Burst) List(ctx context.Context, filters BurstListFilters) ([]*career.Burst, error) {
	query := r.db.WithContext(ctx).Model(&models.Burst{})
	query = r.applyFilters(query, filters)
	query = r.applySorting(query, filters)
	query = r.applyPagination(query, filters)

	var results []models.Burst
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	bursts := make([]*career.Burst, len(results))
	for i, m := range results {
		bursts[i] = m.ToDomain()
	}
	return bursts, nil
}

// Count returns the number of bursts matching the given filters.
func (r *Burst) Count(ctx context.Context, filters BurstListFilters) (int, error) {
	query := r.db.WithContext(ctx).Model(&models.Burst{})
	query = r.applyFilters(query, filters)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *Burst) applyFilters(query *gorm.DB, filters BurstListFilters) *gorm.DB {
	if filters.StartDate != nil {
		query = query.Where("created_at >= ?", filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("created_at <= ?", filters.EndDate)
	}
	return query
}

func (r *Burst) applySorting(query *gorm.DB, filters BurstListFilters) *gorm.DB {
	sortBy := filters.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	order := "DESC"
	if strings.ToLower(filters.SortOrder) == "asc" {
		order = "ASC"
	}

	switch sortBy {
	case "name":
		return query.Order("name " + order)
	case "event_count":
		// Count commas in event_ids to estimate event count.
		return query.Order("LENGTH(event_ids) - LENGTH(REPLACE(event_ids, ',', '')) + 1 " + order)
	default:
		return query.Order("created_at " + order)
	}
}

func (r *Burst) applyPagination(query *gorm.DB, filters BurstListFilters) *gorm.DB {
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}
	return query
}
