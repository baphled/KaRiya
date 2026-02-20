// Package sql provides SQL-backed repository implementations for career domain entities.
package sql

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	career_repo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/repository/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Compile-time interface check.
var _ career_repo.EventRepository = (*EventRepository)(nil)

// EventRepository implements career.EventRepository using GORM.
type EventRepository struct {
	db *gorm.DB
}

// NewEventRepository creates a new SQL event repository.
//
// Expected:
//   - db must be valid.
//
// Returns:
//   - A fully initialized EventRepository ready for use.
//
// Side effects:
//   - None.
func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

// Create adds a new career event to the database.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *EventRepository) Create(ctx context.Context, event *career.Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now

	model := models.EventFromDomain(event)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// Save skill associations.
	if err := r.saveSkillAssociations(ctx, event.ID, event.Skills); err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a career event by its ID.
func (r *EventRepository) GetByID(ctx context.Context, id string) (*career.Event, error) {
	var model models.Event
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, career_repo.ErrEventNotFound
	}
	if err != nil {
		return nil, err
	}

	event := model.ToDomain()

	// Load skill IDs.
	skillIDs, err := r.loadSkillIDs(ctx, id)
	if err != nil {
		return nil, err
	}
	event.Skills = skillIDs

	return event, nil
}

// Update modifies an existing career event.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *EventRepository) Update(ctx context.Context, event *career.Event) error {
	// Check if record exists first.
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Event{}).Where("id = ?", event.ID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return career_repo.ErrEventNotFound
	}

	event.UpdatedAt = time.Now()
	model := models.EventFromDomain(event)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}

	// Update skill associations.
	if err := r.saveSkillAssociations(ctx, event.ID, event.Skills); err != nil {
		return err
	}

	return nil
}

// Delete removes a career event from the database.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *EventRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Event{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return career_repo.ErrEventNotFound
	}
	return result.Error
}

// List retrieves career events with optional filtering.
func (r *EventRepository) List(ctx context.Context, filters career_repo.EventListFilters) ([]*career.Event, error) {
	query := r.db.WithContext(ctx).Model(&models.Event{})
	query = r.applyFilters(query, filters)
	query = r.applySorting(query, filters)
	query = r.applyPagination(query, filters)

	var results []models.Event
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	// Collect event IDs to batch-load associated skill IDs.
	eventIDs := make([]string, 0, len(results))
	for i := range results {
		eventIDs = append(eventIDs, results[i].ID)
	}

	// Batch load all skill associations in a single query.
	eventSkills, err := r.loadSkillIDsForEvents(ctx, eventIDs)
	if err != nil {
		return nil, err
	}

	events := make([]*career.Event, len(results))
	for i := range results {
		events[i] = results[i].ToDomain()
		// Assign preloaded skill IDs from the batched lookup.
		if skills, ok := eventSkills[results[i].ID]; ok {
			events[i].Skills = skills
		}
	}
	return events, nil
}

// Count returns the number of events matching the given filters.
func (r *EventRepository) Count(ctx context.Context, filters career_repo.EventListFilters) (int, error) {
	query := r.db.WithContext(ctx).Model(&models.Event{})
	query = r.applyFilters(query, filters)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *EventRepository) applyFilters(query *gorm.DB, filters career_repo.EventListFilters) *gorm.DB {
	// Tag filtering - use OR logic (match any tag).
	if len(filters.Tags) > 0 {
		tagConditions := make([]string, len(filters.Tags))
		tagArgs := make([]interface{}, len(filters.Tags))
		for i, tag := range filters.Tags {
			tagConditions[i] = "tags LIKE ?"
			tagArgs[i] = "%" + tag + "%"
		}
		// Join conditions with OR to match original behavior.
		query = query.Where("("+strings.Join(tagConditions, " OR ")+")", tagArgs...)
	}

	// Date range filtering.
	if filters.StartDate != nil {
		query = query.Where("date >= ?", filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("date <= ?", filters.EndDate)
	}

	return query
}

func (r *EventRepository) applySorting(query *gorm.DB, filters career_repo.EventListFilters) *gorm.DB {
	order := "ASC"
	if filters.SortOrder == "desc" {
		order = "DESC"
	}

	switch filters.SortBy {
	case "date":
		return query.Order("date " + order + ", created_at " + order)
	default:
		return query.Order("created_at " + order)
	}
}

func (r *EventRepository) applyPagination(query *gorm.DB, filters career_repo.EventListFilters) *gorm.DB {
	limit := filters.Limit
	if limit == 0 {
		limit = defaultPaginationLimit
	}
	return query.Limit(limit).Offset(filters.Offset)
}

func (r *EventRepository) saveSkillAssociations(ctx context.Context, eventID string, skillIDs []string) error {
	// Delete existing associations.
	if err := r.db.WithContext(ctx).Exec("DELETE FROM event_skills WHERE event_id = ?", eventID).Error; err != nil {
		return err
	}

	// Insert new associations.
	for _, skillID := range skillIDs {
		if err := r.db.WithContext(ctx).Exec(
			"INSERT INTO event_skills (event_id, skill_id) VALUES (?, ?)",
			eventID, skillID,
		).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *EventRepository) loadSkillIDs(ctx context.Context, eventID string) ([]string, error) {
	var ids []string
	err := r.db.WithContext(ctx).
		Table("event_skills").
		Where("event_id = ?", eventID).
		Order("skill_id").
		Pluck("skill_id", &ids).Error
	return ids, err
}

// LinkSkill creates an association between an event and a skill.
//
// Expected:
//   - Must be a valid string.
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *EventRepository) LinkSkill(ctx context.Context, eventID string, skillID string) error {
	// First check if event exists
	var exists bool
	err := r.db.WithContext(ctx).
		Table("career_events").
		Select("1").
		Where("id = ?", eventID).
		Limit(1).
		Scan(&exists).Error

	if err != nil {
		return err
	}

	if !exists {
		return career_repo.ErrEventNotFound
	}

	// Insert the skill association if it doesn't already exist (ON CONFLICT DO NOTHING)
	// This prevents duplicate entries in the junction table
	err = r.db.WithContext(ctx).Exec(
		"INSERT INTO event_skills (event_id, skill_id) VALUES (?, ?) ON CONFLICT (event_id, skill_id) DO NOTHING",
		eventID, skillID,
	).Error

	return err
}

// UnlinkSkill removes an association between an event and a skill.
//
// Expected:
//   - eventID must be a valid string identifier for an existing event.
//   - skillID must be a valid string identifier for an existing skill.
//
// Returns:
//   - An error value.
//
// Side effects:
//   - Deletes the event-skill association from the database.
func (r *EventRepository) UnlinkSkill(ctx context.Context, eventID string, skillID string) error {
	var exists bool
	err := r.db.WithContext(ctx).
		Table("career_events").
		Select("1").
		Where("id = ?", eventID).
		Limit(1).
		Scan(&exists).Error

	if err != nil {
		return err
	}

	if !exists {
		return career_repo.ErrEventNotFound
	}

	err = r.db.WithContext(ctx).Exec(
		"DELETE FROM event_skills WHERE event_id = ? AND skill_id = ?",
		eventID, skillID,
	).Error

	return err
}

// loadSkillIDsForEvents batch loads skill IDs for multiple events in a single query.
func (r *EventRepository) loadSkillIDsForEvents(ctx context.Context, eventIDs []string) (map[string][]string, error) {
	if len(eventIDs) == 0 {
		return make(map[string][]string), nil
	}

	// Struct to scan event-skill associations.
	type eventSkillRow struct {
		EventID string `gorm:"column:event_id"`
		SkillID string `gorm:"column:skill_id"`
	}

	var rows []eventSkillRow
	if err := r.db.WithContext(ctx).
		Table("event_skills").
		Select("event_id, skill_id").
		Where("event_id IN ?", eventIDs).
		Order("event_id, skill_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	// Build map from event ID to skill IDs.
	eventSkills := make(map[string][]string, len(eventIDs))
	for _, row := range rows {
		eventSkills[row.EventID] = append(eventSkills[row.EventID], row.SkillID)
	}

	return eventSkills, nil
}
