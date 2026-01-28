// Package career provides repository implementations for career domain entities.
package career

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/repository/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Compile-time interface check.
var _ SkillRepository = (*Skill)(nil)

// Skill implements SkillRepository using GORM.
type Skill struct {
	db *gorm.DB
}

// NewSkill creates a new skill repository.
func NewSkill(db *gorm.DB) *Skill {
	return &Skill{db: db}
}

func (r *Skill) Create(ctx context.Context, skill *career.Skill) error {
	if skill.ID == "" {
		skill.ID = uuid.New().String()
	}

	now := time.Now()
	skill.CreatedAt = now
	skill.UpdatedAt = now

	return r.db.WithContext(ctx).Create(models.SkillFromDomain(skill)).Error
}

func (r *Skill) GetByID(ctx context.Context, id string) (*career.Skill, error) {
	var model models.Skill
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSkillNotFound
	}
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *Skill) GetByName(ctx context.Context, name string) (*career.Skill, error) {
	var model models.Skill
	err := r.db.WithContext(ctx).First(&model, "name = ?", name).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSkillNotFound
	}
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *Skill) Update(ctx context.Context, skill *career.Skill) error {
	// Check if record exists first.
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Skill{}).Where("id = ?", skill.ID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrSkillNotFound
	}

	skill.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(models.SkillFromDomain(skill)).Error
}

func (r *Skill) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Skill{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return ErrSkillNotFound
	}
	return result.Error
}

func (r *Skill) List(ctx context.Context, filters *SkillFilters) ([]*career.Skill, error) {
	query := r.db.WithContext(ctx).Model(&models.Skill{})
	query = r.applyFilters(query, filters)
	query = r.applySorting(query, filters)
	query = r.applyPagination(query, filters)

	var results []models.Skill
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	skills := make([]*career.Skill, len(results))
	for i, m := range results {
		skills[i] = m.ToDomain()
	}
	return skills, nil
}

func (r *Skill) Count(ctx context.Context, filters *SkillFilters) (int, error) {
	query := r.db.WithContext(ctx).Model(&models.Skill{})
	query = r.applyFilters(query, filters)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *Skill) applyFilters(query *gorm.DB, filters *SkillFilters) *gorm.DB {
	if filters == nil {
		return query
	}
	if filters.Category != "" {
		query = query.Where("category = ?", filters.Category)
	}
	if filters.Level != "" {
		query = query.Where("level = ?", filters.Level)
	}
	if filters.MinEvents > 0 {
		query = query.Where(
			"(SELECT COUNT(*) FROM event_skills WHERE skill_id = skills.id) >= ?",
			filters.MinEvents,
		)
	}
	return query
}

func (r *Skill) applySorting(query *gorm.DB, filters *SkillFilters) *gorm.DB {
	if filters == nil || filters.SortBy == "" {
		return query.Order("name ASC")
	}

	order := "ASC"
	if filters.SortOrder == "desc" {
		order = "DESC"
	}

	switch filters.SortBy {
	case "name":
		return query.Order(fmt.Sprintf("name %s", order))
	case "category":
		return query.Order(fmt.Sprintf("category %s, name ASC", order))
	case "events":
		return query.Order(fmt.Sprintf(
			"(SELECT COUNT(*) FROM event_skills WHERE skill_id = skills.id) %s, name ASC",
			order,
		))
	case "last_used":
		return query.Order(fmt.Sprintf(
			"(SELECT MAX(date) FROM career_events ce JOIN event_skills es ON ce.id = es.event_id WHERE es.skill_id = skills.id) %s NULLS LAST, name ASC",
			order,
		))
	default:
		return query.Order("name ASC")
	}
}

func (r *Skill) applyPagination(query *gorm.DB, filters *SkillFilters) *gorm.DB {
	if filters == nil {
		return query
	}
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}
	return query
}

func (r *Skill) LinkToEvent(ctx context.Context, skillID, eventID string) error {
	return r.db.WithContext(ctx).Exec(
		"INSERT OR IGNORE INTO event_skills (skill_id, event_id) VALUES (?, ?)",
		skillID, eventID,
	).Error
}

func (r *Skill) UnlinkFromEvent(ctx context.Context, skillID, eventID string) error {
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM event_skills WHERE skill_id = ? AND event_id = ?",
		skillID, eventID,
	).Error
}

func (r *Skill) GetEventIDs(ctx context.Context, skillID string) ([]string, error) {
	var ids []string
	err := r.db.WithContext(ctx).
		Table("event_skills").
		Where("skill_id = ?", skillID).
		Pluck("event_id", &ids).Error
	return ids, err
}

// GetByCategory retrieves all skills in a specific category.
func (r *Skill) GetByCategory(ctx context.Context, category string) ([]*career.Skill, error) {
	return r.List(ctx, &SkillFilters{Category: category})
}

// GetSkillsForEvent retrieves all skills associated with an event.
func (r *Skill) GetSkillsForEvent(ctx context.Context, eventID string) ([]*career.Skill, error) {
	var results []models.Skill
	err := r.db.WithContext(ctx).
		Joins("JOIN event_skills ON event_skills.skill_id = skills.id").
		Where("event_skills.event_id = ?", eventID).
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	skills := make([]*career.Skill, len(results))
	for i, m := range results {
		skills[i] = m.ToDomain()
	}
	return skills, nil
}

// GetEventCountsForSkills returns a map of skill IDs to event counts.
func (r *Skill) GetEventCountsForSkills(ctx context.Context) (map[string]int, error) {
	type result struct {
		SkillID string
		Count   int
	}
	var results []result
	err := r.db.WithContext(ctx).
		Table("event_skills").
		Select("skill_id, COUNT(*) as count").
		Group("skill_id").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int, len(results))
	for _, r := range results {
		counts[r.SkillID] = r.Count
	}
	return counts, nil
}

// GetLastUsedForSkills returns a map of skill IDs to their last used dates.
func (r *Skill) GetLastUsedForSkills(ctx context.Context) (map[string]time.Time, error) {
	type result struct {
		SkillID  string
		LastUsed time.Time
	}
	var results []result
	err := r.db.WithContext(ctx).
		Table("event_skills").
		Select("event_skills.skill_id, MAX(career_events.date) as last_used").
		Joins("JOIN career_events ON career_events.id = event_skills.event_id").
		Group("event_skills.skill_id").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	lastUsed := make(map[string]time.Time, len(results))
	for _, r := range results {
		lastUsed[r.SkillID] = r.LastUsed
	}
	return lastUsed, nil
}

// GetEventsUsingSkill returns all events that use a specific skill, ordered by date DESC.
func (r *Skill) GetEventsUsingSkill(ctx context.Context, skillID string) ([]*career.CareerEvent, error) {
	var results []models.Event
	err := r.db.WithContext(ctx).
		Joins("JOIN event_skills ON event_skills.event_id = career_events.id").
		Where("event_skills.skill_id = ?", skillID).
		Order("date DESC").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	events := make([]*career.CareerEvent, len(results))
	for i, m := range results {
		events[i] = m.ToDomain()
	}
	return events, nil
}
