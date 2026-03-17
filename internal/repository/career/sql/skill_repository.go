// Package sql provides SQL-backed repository implementations for career domain entities.
package sql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	models "github.com/baphled/kariya/internal/model/career"
	career_repo "github.com/baphled/kariya/internal/repository/career"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Compile-time interface check.
var _ career_repo.SkillRepository = (*SkillRepository)(nil)

// sqliteDateFormats lists the timestamp formats that SQLite may return for
// date columns, ordered from most specific to least specific. The pure-Go
// SQLite driver (modernc.org/sqlite) returns aggregated date values as raw
// strings, so we must parse them manually.
var sqliteDateFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

// SkillRepository implements career.SkillRepository using GORM.
type SkillRepository struct {
	db *gorm.DB
}

// NewSkillRepository creates a new SQL skill repository.
//
// Expected:
//   - db must be valid.
//
// Returns:
//   - A fully initialized SkillRepository ready for use.
//
// Side effects:
//   - None.
func NewSkillRepository(db *gorm.DB) *SkillRepository {
	return &SkillRepository{db: db}
}

// Create persists a new skill to the database.
//
// Expected:
//   - skill must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
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
// Expected:
//   - id must be a valid skill identifier.
//
// Returns:
//   - The skill if found, or nil with ErrSkillNotFound if not found.
//
// Side effects:
//   - None.
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

// GetByName retrieves a single skill by name using case-insensitive matching.
//
// Expected:
//   - name must be a non-empty string.
//
// Returns:
//   - The skill if found, or nil with ErrSkillNotFound if not found.
//
// Side effects:
//   - None.
func (r *SkillRepository) GetByName(ctx context.Context, name string) (*career.Skill, error) {
	var model models.Skill
	err := r.db.WithContext(ctx).First(&model, "LOWER(name) = LOWER(?)", name).Error
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
// Expected:
//   - skill must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
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
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *SkillRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Skill{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return career_repo.ErrSkillNotFound
	}
	return result.Error
}

// List retrieves all skills that satisfy the given filters.
//
// Expected:
//   - filters specifies the filtering, sorting, and pagination options.
//
// Returns:
//   - A slice of skills matching the filters, or an error.
//
// Side effects:
//   - None.
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
// Expected:
//   - filters specifies the filtering options.
//
// Returns:
//   - The count of matching skills, or an error.
//
// Side effects:
//   - None.
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
		return query.Order("name " + order)
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
		return query.Limit(defaultPaginationLimit)
	}
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

// LinkToEvent creates an association between a skill and a career event
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *SkillRepository) LinkToEvent(ctx context.Context, skillID, eventID string) error {
	return r.db.WithContext(ctx).Exec(
		"INSERT OR IGNORE INTO event_skills (skill_id, event_id) VALUES (?, ?)",
		skillID, eventID,
	).Error
}

// UnlinkFromEvent removes the association between a skill and a career
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *SkillRepository) UnlinkFromEvent(ctx context.Context, skillID, eventID string) error {
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM event_skills WHERE skill_id = ? AND event_id = ?",
		skillID, eventID,
	).Error
}

// GetEventIDs returns the IDs of all career events linked to the given skill.
//
// Expected:
//   - skillID must be a valid skill identifier.
//
// Returns:
//   - A slice of event IDs, or an error.
//
// Side effects:
//   - None.
func (r *SkillRepository) GetEventIDs(ctx context.Context, skillID string) ([]string, error) {
	var ids []string
	err := r.db.WithContext(ctx).
		Table("event_skills").
		Where("skill_id = ?", skillID).
		Pluck("event_id", &ids).Error
	return ids, err
}

// GetByCategory retrieves all skills in a specific category.
//
// Expected:
//   - category must be a valid category string.
//
// Returns:
//   - A slice of skills in the category, or an error.
//
// Side effects:
//   - None.
func (r *SkillRepository) GetByCategory(ctx context.Context, category string) ([]*career.Skill, error) {
	return r.List(ctx, &career_repo.SkillListFilters{Category: category})
}

// GetSkillsForEvent retrieves all skills associated with an event.
//
// Expected:
//   - eventID must be a valid event identifier.
//
// Returns:
//   - A slice of skills for the event, or an error.
//
// Side effects:
//   - None.
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

// GetSkillsForEvents retrieves all unique skills linked to any of the given events.
//
// Expected:
//   - eventIDs must be a slice of valid event identifiers.
//
// Returns:
//   - A slice of unique skills across all events, or an error.
//
// Side effects:
//   - None.
func (r *SkillRepository) GetSkillsForEvents(ctx context.Context, eventIDs []string) ([]*career.Skill, error) {
	if len(eventIDs) == 0 {
		return []*career.Skill{}, nil
	}

	var results []models.Skill
	err := r.db.WithContext(ctx).
		Joins("JOIN event_skills ON event_skills.skill_id = skills.id").
		Where("event_skills.event_id IN ?", eventIDs).
		Distinct().
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
//
// Expected:
//   - None.
//
// Returns:
//   - A map of skill IDs to their event counts, or an error.
//
// Side effects:
//   - None.
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
//
// Expected:
//   - None.
//
// Returns:
//   - A map of skill IDs to last used times, or an error.
//
// Side effects:
//   - None.
func (r *SkillRepository) GetLastUsedForSkills(ctx context.Context) (map[string]time.Time, error) {
	type result struct {
		SkillID  string
		LastUsed string
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
	for _, row := range results {
		if row.LastUsed == "" {
			continue
		}
		var parsed time.Time
		var parseErr error
		for _, layout := range sqliteDateFormats {
			parsed, parseErr = time.ParseInLocation(layout, row.LastUsed, time.UTC)
			if parseErr == nil {
				break
			}
		}
		if parseErr != nil {
			return nil, fmt.Errorf("parsing last_used date for skill %s (value %q): %w", row.SkillID, row.LastUsed, parseErr)
		}
		lastUsed[row.SkillID] = parsed
	}
	return lastUsed, nil
}

// GetEventsUsingSkill returns all events that use a specific skill, ordered by date DESC.
//
// Expected:
//   - skillID must be a valid skill identifier.
//
// Returns:
//   - A slice of events using the skill, or an error.
//
// Side effects:
//   - None.
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
