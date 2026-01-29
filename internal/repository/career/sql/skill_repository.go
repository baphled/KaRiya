// Package sql provides SQL-backed repository implementations for career domain entities.
package sql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	career_repo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/repository/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Compile-time interface check.
var _ career_repo.SkillRepository = (*SkillRepository)(nil)

// SkillRepository implements career.SkillRepository using GORM.
type SkillRepository struct {
	db *gorm.DB
}

// NewSkillRepository creates a new SQL skill repository.
func NewSkillRepository(db *gorm.DB) *SkillRepository {
	return &SkillRepository{db: db}
}

// Create persists a new skill to the database.
//
// The skill parameter is modified in place: if its ID field is empty a UUID
// is generated, and CreatedAt and UpdatedAt are both set to the current time
// regardless of their incoming values. The ctx parameter carries
// request-scoped deadlines and cancellation signals through to the
// underlying GORM query.
//
// Returns nil on success. Returns a GORM error if the write fails, for
// example when a skill with the same name already exists (unique constraint
// violation).
func (r *SkillRepository) Create(ctx context.Context, skill *career.Skill) error {
	if skill.ID == "" {
		skill.ID = uuid.New().String()
	}

	now := time.Now()
	skill.CreatedAt = now
	skill.UpdatedAt = now

	return r.db.WithContext(ctx).Create(models.SkillFromDomain(skill)).Error
}

// GetByID retrieves a single skill by its unique identifier.
//
// The id parameter is matched against the primary key column. Returns the
// fully populated domain Skill and nil error on success. Returns nil and
// ErrSkillNotFound if no row matches. Returns nil and a GORM error for
// any other database failure.
func (r *SkillRepository) GetByID(ctx context.Context, id string) (*career.Skill, error) {
	var model models.Skill
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, career_repo.ErrSkillNotFound
	}
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

// GetByName retrieves a single skill by its exact name.
//
// The name parameter is compared case-sensitively against the unique name
// column. Returns the domain Skill and nil on a match. Returns nil and
// ErrSkillNotFound when no row matches. Returns nil and a GORM error for
// any other database failure.
func (r *SkillRepository) GetByName(ctx context.Context, name string) (*career.Skill, error) {
	var model models.Skill
	err := r.db.WithContext(ctx).First(&model, "name = ?", name).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, career_repo.ErrSkillNotFound
	}
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

// Update replaces an existing skill record with the values from skill.
//
// The skill parameter must have a non-empty ID that matches an existing row.
// UpdatedAt is overwritten with the current time before saving. All other
// fields on skill are persisted as-is, fully replacing the previous values.
//
// Returns nil on success. Returns ErrSkillNotFound if no row matches the
// ID. Returns a GORM error for any other database failure.
func (r *SkillRepository) Update(ctx context.Context, skill *career.Skill) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Skill{}).Where("id = ?", skill.ID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return career_repo.ErrSkillNotFound
	}

	skill.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(models.SkillFromDomain(skill)).Error
}

// Delete removes a skill record identified by id.
//
// Returns nil on success. Returns ErrSkillNotFound when no row matches the
// id. Returns a GORM error for any other database failure. Note that rows
// in the event_skills junction table that reference this skill are not
// cascade-deleted; callers should unlink events first if referential
// integrity is required.
func (r *SkillRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Skill{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return career_repo.ErrSkillNotFound
	}
	return result.Error
}

// List retrieves all skills that satisfy the given filters.
//
// The filters parameter controls which skills are returned and in what
// order. Supported filter fields are Category (exact match), Level (exact
// match), and MinEvents (skills linked to at least N events via the
// event_skills junction table). When filters is nil every skill is
// returned. Sorting defaults to name ascending; set SortBy to "name",
// "category", "events", or "last_used" with an optional SortOrder of
// "desc" to override. Pagination is controlled by Limit and Offset.
//
// Returns the matching domain skills and nil on success. Returns nil and
// a GORM error on database failure.
func (r *SkillRepository) List(ctx context.Context, filters *career_repo.SkillListFilters) ([]*career.Skill, error) {
	query := r.db.WithContext(ctx).Model(&models.Skill{})
	query = r.applyFilters(query, filters)
	query = r.applySorting(query, filters)
	query = r.applyPagination(query, filters)

	var results []models.Skill
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	skills := make([]*career.Skill, len(results))
	for i := range results {
		skills[i] = results[i].ToDomain()
	}
	return skills, nil
}

// Count returns the total number of skills that satisfy the given filters.
//
// The filters parameter accepts the same Category, Level, and MinEvents
// constraints as List, but Limit, Offset, SortBy, and SortOrder are
// ignored. When filters is nil the total count of all skills is returned.
//
// Returns the count and nil on success. Returns 0 and a GORM error on
// database failure.
func (r *SkillRepository) Count(ctx context.Context, filters *career_repo.SkillListFilters) (int, error) {
	query := r.db.WithContext(ctx).Model(&models.Skill{})
	query = r.applyFilters(query, filters)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *SkillRepository) applyFilters(query *gorm.DB, filters *career_repo.SkillListFilters) *gorm.DB {
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

func (r *SkillRepository) applySorting(query *gorm.DB, filters *career_repo.SkillListFilters) *gorm.DB {
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
		subquery := "(SELECT MAX(date) FROM career_events ce JOIN event_skills es ON ce.id = es.event_id WHERE es.skill_id = skills.id)"
		return query.Order(fmt.Sprintf(
			"CASE WHEN %s IS NULL THEN 1 ELSE 0 END, %s %s, name ASC",
			subquery, subquery, order,
		))
	default:
		return query.Order("name ASC")
	}
}

func (r *SkillRepository) applyPagination(query *gorm.DB, filters *career_repo.SkillListFilters) *gorm.DB {
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

// LinkToEvent creates an association between a skill and a career event
// in the event_skills junction table.
//
// The skillID and eventID parameters identify the two rows to link. If the
// pair already exists the operation is silently ignored via INSERT OR
// IGNORE, making this method idempotent.
//
// Returns nil on success or a database error on failure.
func (r *SkillRepository) LinkToEvent(ctx context.Context, skillID, eventID string) error {
	return r.db.WithContext(ctx).Exec(
		"INSERT OR IGNORE INTO event_skills (skill_id, event_id) VALUES (?, ?)",
		skillID, eventID,
	).Error
}

// UnlinkFromEvent removes the association between a skill and a career
// event from the event_skills junction table.
//
// The skillID and eventID parameters identify the link to remove. If the
// pair does not exist the operation succeeds silently.
//
// Returns nil on success or a database error on failure.
func (r *SkillRepository) UnlinkFromEvent(ctx context.Context, skillID, eventID string) error {
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM event_skills WHERE skill_id = ? AND event_id = ?",
		skillID, eventID,
	).Error
}

// GetEventIDs returns the IDs of all career events linked to the given skill.
//
// The skillID parameter identifies the skill whose event associations are
// queried from the event_skills junction table.
//
// Returns the event ID slice and nil on success. Returns nil and a
// database error on failure. An empty slice is returned when the skill
// has no linked events.
func (r *SkillRepository) GetEventIDs(ctx context.Context, skillID string) ([]string, error) {
	var ids []string
	err := r.db.WithContext(ctx).
		Table("event_skills").
		Where("skill_id = ?", skillID).
		Pluck("event_id", &ids).Error
	return ids, err
}

// GetByCategory retrieves all skills in a specific category.
func (r *SkillRepository) GetByCategory(ctx context.Context, category string) ([]*career.Skill, error) {
	return r.List(ctx, &career_repo.SkillListFilters{Category: category})
}

// GetSkillsForEvent retrieves all skills associated with an event.
func (r *SkillRepository) GetSkillsForEvent(ctx context.Context, eventID string) ([]*career.Skill, error) {
	var results []models.Skill
	err := r.db.WithContext(ctx).
		Joins("JOIN event_skills ON event_skills.skill_id = skills.id").
		Where("event_skills.event_id = ?", eventID).
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	skills := make([]*career.Skill, len(results))
	for i := range results {
		skills[i] = results[i].ToDomain()
	}
	return skills, nil
}

// GetEventCountsForSkills returns a map of skill IDs to event counts.
func (r *SkillRepository) GetEventCountsForSkills(ctx context.Context) (map[string]int, error) {
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
func (r *SkillRepository) GetLastUsedForSkills(ctx context.Context) (map[string]time.Time, error) {
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
func (r *SkillRepository) GetEventsUsingSkill(ctx context.Context, skillID string) ([]*career.Event, error) {
	var results []models.Event
	err := r.db.WithContext(ctx).
		Joins("JOIN event_skills ON event_skills.event_id = career_events.id").
		Where("event_skills.skill_id = ?", skillID).
		Order("date DESC").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	events := make([]*career.Event, len(results))
	for i := range results {
		events[i] = results[i].ToDomain()
	}
	return events, nil
}
