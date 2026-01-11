package career

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

// SQLiteSkillRepository provides a SQLite implementation of SkillRepository
type SQLiteSkillRepository struct {
	db *sql.DB
}

// NewSQLiteSkillRepository creates a new SQLite skill repository.
// Deprecated: Use NewSQLiteSkillRepositoryWithDB after running migrations via career.RunMigrations().
// This constructor is kept for backward compatibility with existing tests.
func NewSQLiteSkillRepository(db *sql.DB) (*SQLiteSkillRepository, error) {
	// Run migrations to ensure schema is up to date
	if err := RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &SQLiteSkillRepository{db: db}, nil
}

// NewSQLiteSkillRepositoryWithDB creates a skill repository with an existing database connection.
// This is useful when migrations have already been run on the database connection.
// The caller is responsible for managing the database connection lifecycle.
func NewSQLiteSkillRepositoryWithDB(db *sql.DB) *SQLiteSkillRepository {
	return &SQLiteSkillRepository{db: db}
}

// Create adds a new skill to the database
func (r *SQLiteSkillRepository) Create(ctx context.Context, skill *career.Skill) error {
	// Generate a unique ID if not provided
	if skill.ID == "" {
		skill.ID = uuid.New().String()
	}

	// Validate the skill
	if err := skill.Validate(); err != nil {
		return fmt.Errorf("skill validation failed for ID %s: %w", skill.ID, err)
	}

	// Check for duplicate name (case-sensitive)
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM skills WHERE name = ?", skill.Name).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking for duplicate skill name: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("duplicate skill with name %s: %w", skill.Name, ErrDuplicateSkill)
	}

	// Set timestamps
	now := time.Now()
	skill.CreatedAt = now
	skill.UpdatedAt = now

	// Insert the skill
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO skills
		(id, name, category, level, years_used, last_used, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, skill.ID, skill.Name, skill.Category, skill.Level, skill.YearsUsed, skill.LastUsed, skill.CreatedAt, skill.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create skill with ID %s: %w", skill.ID, err)
	}

	return nil
}

// GetByID retrieves a skill by its ID
func (r *SQLiteSkillRepository) GetByID(ctx context.Context, id string) (*career.Skill, error) {
	var skill career.Skill
	var yearsUsed sql.NullInt64
	var lastUsed sql.NullTime

	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, category, level, years_used, last_used, created_at, updated_at
		FROM skills
		WHERE id = ?
	`, id)

	err := row.Scan(
		&skill.ID,
		&skill.Name,
		&skill.Category,
		&skill.Level,
		&yearsUsed,
		&lastUsed,
		&skill.CreatedAt,
		&skill.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("skill with ID %s not found: %w", id, ErrSkillNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve skill with ID %s: %w", id, err)
	}

	// Handle nullable fields
	if yearsUsed.Valid {
		years := int(yearsUsed.Int64)
		skill.YearsUsed = &years
	}
	if lastUsed.Valid {
		skill.LastUsed = &lastUsed.Time
	}

	return &skill, nil
}

// GetByName retrieves a skill by its name (case-sensitive)
func (r *SQLiteSkillRepository) GetByName(ctx context.Context, name string) (*career.Skill, error) {
	var skill career.Skill
	var yearsUsed sql.NullInt64
	var lastUsed sql.NullTime

	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, category, level, years_used, last_used, created_at, updated_at
		FROM skills
		WHERE name = ?
	`, name)

	err := row.Scan(
		&skill.ID,
		&skill.Name,
		&skill.Category,
		&skill.Level,
		&yearsUsed,
		&lastUsed,
		&skill.CreatedAt,
		&skill.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("skill with name %s not found: %w", name, ErrSkillNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve skill with name %s: %w", name, err)
	}

	// Handle nullable fields
	if yearsUsed.Valid {
		years := int(yearsUsed.Int64)
		skill.YearsUsed = &years
	}
	if lastUsed.Valid {
		skill.LastUsed = &lastUsed.Time
	}

	return &skill, nil
}

// List retrieves skills with optional filtering
func (r *SQLiteSkillRepository) List(ctx context.Context, filters *SkillFilters) ([]*career.Skill, error) {
	query := `SELECT id, name, category, level, years_used, last_used, created_at, updated_at FROM skills WHERE 1=1`
	args := []interface{}{}

	// Apply filters
	if filters != nil {
		if filters.Category != "" {
			query += " AND category = ?"
			args = append(args, filters.Category)
		}
		if filters.Level != "" {
			query += " AND level = ?"
			args = append(args, filters.Level)
		}
	}

	// Add ordering
	query += " ORDER BY name ASC"

	// Apply pagination
	if filters != nil {
		// SQLite requires LIMIT before OFFSET
		if filters.Limit > 0 || filters.Offset > 0 {
			if filters.Limit > 0 {
				query += " LIMIT ?"
				args = append(args, filters.Limit)
			} else {
				// If only offset is specified, use a large limit
				query += " LIMIT -1"
			}

			if filters.Offset > 0 {
				query += " OFFSET ?"
				args = append(args, filters.Offset)
			}
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list skills: %w", err)
	}
	defer rows.Close()

	var skills []*career.Skill
	for rows.Next() {
		var skill career.Skill
		var yearsUsed sql.NullInt64
		var lastUsed sql.NullTime

		err := rows.Scan(
			&skill.ID,
			&skill.Name,
			&skill.Category,
			&skill.Level,
			&yearsUsed,
			&lastUsed,
			&skill.CreatedAt,
			&skill.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan skill: %w", err)
		}

		// Handle nullable fields
		if yearsUsed.Valid {
			years := int(yearsUsed.Int64)
			skill.YearsUsed = &years
		}
		if lastUsed.Valid {
			skill.LastUsed = &lastUsed.Time
		}

		skills = append(skills, &skill)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating skills: %w", err)
	}

	return skills, nil
}

// Update modifies an existing skill
func (r *SQLiteSkillRepository) Update(ctx context.Context, skill *career.Skill) error {
	// Validate the skill
	if err := skill.Validate(); err != nil {
		return fmt.Errorf("skill validation failed for ID %s: %w", skill.ID, err)
	}

	// Check if skill exists
	var exists int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM skills WHERE id = ?", skill.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if skill exists: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("skill with ID %s: %w", skill.ID, ErrSkillNotFound)
	}

	// Update timestamp
	skill.UpdatedAt = time.Now()

	// Update the skill
	_, err = r.db.ExecContext(ctx, `
		UPDATE skills
		SET name = ?, category = ?, level = ?, years_used = ?, last_used = ?, updated_at = ?
		WHERE id = ?
	`, skill.Name, skill.Category, skill.Level, skill.YearsUsed, skill.LastUsed, skill.UpdatedAt, skill.ID)

	if err != nil {
		return fmt.Errorf("failed to update skill with ID %s: %w", skill.ID, err)
	}

	return nil
}

// Delete removes a skill from the database
func (r *SQLiteSkillRepository) Delete(ctx context.Context, id string) error {
	// Check if skill exists
	var exists int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM skills WHERE id = ?", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if skill exists: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("skill with ID %s: %w", id, ErrSkillNotFound)
	}

	// Delete the skill
	_, err = r.db.ExecContext(ctx, "DELETE FROM skills WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete skill with ID %s: %w", id, err)
	}

	return nil
}

// GetByCategory retrieves all skills in a specific category
func (r *SQLiteSkillRepository) GetByCategory(ctx context.Context, category string) ([]*career.Skill, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, category, level, years_used, last_used, created_at, updated_at
		FROM skills
		WHERE category = ?
		ORDER BY name ASC
	`, category)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills by category: %w", err)
	}
	defer rows.Close()

	var skills []*career.Skill
	for rows.Next() {
		var skill career.Skill
		var yearsUsed sql.NullInt64
		var lastUsed sql.NullTime

		err := rows.Scan(
			&skill.ID,
			&skill.Name,
			&skill.Category,
			&skill.Level,
			&yearsUsed,
			&lastUsed,
			&skill.CreatedAt,
			&skill.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan skill: %w", err)
		}

		// Handle nullable fields
		if yearsUsed.Valid {
			years := int(yearsUsed.Int64)
			skill.YearsUsed = &years
		}
		if lastUsed.Valid {
			skill.LastUsed = &lastUsed.Time
		}

		skills = append(skills, &skill)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating skills: %w", err)
	}

	return skills, nil
}

// GetSkillsForEvent retrieves all skills associated with an event
func (r *SQLiteSkillRepository) GetSkillsForEvent(ctx context.Context, eventID string) ([]*career.Skill, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.id, s.name, s.category, s.level, s.years_used, s.last_used, s.created_at, s.updated_at
		FROM skills s
		INNER JOIN event_skills es ON s.id = es.skill_id
		WHERE es.event_id = ?
		ORDER BY s.name ASC
	`, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills for event: %w", err)
	}
	defer rows.Close()

	var skills []*career.Skill
	for rows.Next() {
		var skill career.Skill
		var yearsUsed sql.NullInt64
		var lastUsed sql.NullTime

		err := rows.Scan(
			&skill.ID,
			&skill.Name,
			&skill.Category,
			&skill.Level,
			&yearsUsed,
			&lastUsed,
			&skill.CreatedAt,
			&skill.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan skill: %w", err)
		}

		// Handle nullable fields
		if yearsUsed.Valid {
			years := int(yearsUsed.Int64)
			skill.YearsUsed = &years
		}
		if lastUsed.Valid {
			skill.LastUsed = &lastUsed.Time
		}

		skills = append(skills, &skill)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating skills: %w", err)
	}

	return skills, nil
}

// GetEventCountsForSkills returns a map of skill IDs to event counts
func (r *SQLiteSkillRepository) GetEventCountsForSkills(ctx context.Context) (map[string]int, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT skill_id, COUNT(event_id) as event_count
		FROM event_skills
		GROUP BY skill_id
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get event counts for skills: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var skillID string
		var count int

		err := rows.Scan(&skillID, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event count: %w", err)
		}

		counts[skillID] = count
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating event counts: %w", err)
	}

	// Return empty map if no associations
	if len(counts) == 0 {
		return make(map[string]int), nil
	}

	return counts, nil
}

// Unused import check - remove if sql.NullString not needed
var _ = strings.Builder{}
