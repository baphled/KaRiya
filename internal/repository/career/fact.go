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
var _ FactRepository = (*Fact)(nil)

// Fact implements FactRepository using GORM.
type Fact struct {
	db *gorm.DB
}

// NewFact creates a new fact repository.
func NewFact(db *gorm.DB) *Fact {
	return &Fact{db: db}
}

// Create adds a new fact to the database.
func (r *Fact) Create(ctx context.Context, fact *career.Fact) error {
	if fact.ID == "" {
		fact.ID = uuid.New().String()
	}

	now := time.Now()
	fact.CreatedAt = now
	fact.UpdatedAt = now

	return r.db.WithContext(ctx).Create(models.FactFromDomain(fact)).Error
}

// GetByID retrieves a fact by its ID.
func (r *Fact) GetByID(ctx context.Context, id string) (*career.Fact, error) {
	var model models.Fact
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrFactNotFound
	}
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

// Update modifies an existing fact.
func (r *Fact) Update(ctx context.Context, fact *career.Fact) error {
	// Check if record exists first.
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Fact{}).Where("id = ?", fact.ID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrFactNotFound
	}

	fact.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(models.FactFromDomain(fact)).Error
}

// Delete removes a fact from the database.
func (r *Fact) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Fact{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return ErrFactNotFound
	}
	return result.Error
}

// List retrieves facts with optional filtering.
func (r *Fact) List(ctx context.Context, filters FactListFilters) ([]*career.Fact, error) {
	query := r.db.WithContext(ctx).Model(&models.Fact{})
	query = r.applyFilters(query, filters)
	query = r.applySorting(query, filters)
	query = r.applyPagination(query, filters)

	var results []models.Fact
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	facts := make([]*career.Fact, len(results))
	for i, m := range results {
		facts[i] = m.ToDomain()
	}
	return facts, nil
}

// Count returns the number of facts matching the given filters.
func (r *Fact) Count(ctx context.Context, filters FactListFilters) (int, error) {
	query := r.db.WithContext(ctx).Model(&models.Fact{})
	query = r.applyFilters(query, filters)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// GetBySourceEventID retrieves all facts for a specific event.
func (r *Fact) GetBySourceEventID(ctx context.Context, eventID string) ([]*career.Fact, error) {
	var results []models.Fact
	err := r.db.WithContext(ctx).Where("source_event_id = ?", eventID).Find(&results).Error
	if err != nil {
		return nil, err
	}

	facts := make([]*career.Fact, len(results))
	for i, m := range results {
		facts[i] = m.ToDomain()
	}
	return facts, nil
}

// GetBySourceBurstID retrieves all facts for a specific burst.
func (r *Fact) GetBySourceBurstID(ctx context.Context, burstID string) ([]*career.Fact, error) {
	var results []models.Fact
	err := r.db.WithContext(ctx).Where("source_burst_id = ?", burstID).Find(&results).Error
	if err != nil {
		return nil, err
	}

	facts := make([]*career.Fact, len(results))
	for i, m := range results {
		facts[i] = m.ToDomain()
	}
	return facts, nil
}

func (r *Fact) applyFilters(query *gorm.DB, filters FactListFilters) *gorm.DB {
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

func (r *Fact) applySorting(query *gorm.DB, filters FactListFilters) *gorm.DB {
	sortBy := filters.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	order := "DESC"
	if strings.ToLower(filters.SortOrder) == "asc" {
		order = "ASC"
	}

	switch sortBy {
	case "text":
		return query.Order("text " + order)
	default:
		return query.Order("created_at " + order)
	}
}

func (r *Fact) applyPagination(query *gorm.DB, filters FactListFilters) *gorm.DB {
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}
	return query
}
