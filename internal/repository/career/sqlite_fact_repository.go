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

// SQLiteFactRepository provides a SQLite implementation of FactRepository
type SQLiteFactRepository struct {
	db *sql.DB
}

// NewSQLiteFactRepository creates a new SQLite fact repository.
// Deprecated: Use NewSQLiteFactRepositoryWithDB after running migrations via career.RunMigrations().
// This constructor is kept for backward compatibility with existing tests.
func NewSQLiteFactRepository(db *sql.DB) (*SQLiteFactRepository, error) {
	// Run migrations to ensure schema is up to date
	if err := RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &SQLiteFactRepository{db: db}, nil
}

// NewSQLiteFactRepositoryWithDB creates a fact repository with an existing database connection.
// This is useful when migrations have already been run on the database connection.
// The caller is responsible for managing the database connection lifecycle.
func NewSQLiteFactRepositoryWithDB(db *sql.DB) *SQLiteFactRepository {
	return &SQLiteFactRepository{db: db}
}

// Create adds a new fact to the database
func (r *SQLiteFactRepository) Create(ctx context.Context, fact *career.Fact) error {
	// Generate a unique ID if not provided
	if fact.ID == "" {
		fact.ID = uuid.New().String()
	}

	// Validate the fact
	if err := fact.Validate(); err != nil {
		return fmt.Errorf("fact validation failed for ID %s: %w", fact.ID, err)
	}

	// Check for duplicate
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM facts WHERE id = ?", fact.ID).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking for duplicate fact: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("duplicate fact with ID %s: %w", fact.ID, ErrDuplicateFact)
	}

	// Set timestamps
	now := time.Now()
	fact.CreatedAt = now
	fact.UpdatedAt = now

	// Convert slices to strings
	competenciesString := strings.Join(fact.CompetencyCategories, ",")
	audienceString := strings.Join(fact.AudienceRelevance, ",")

	// Insert the fact
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO facts
		(id, text, competencies, role_fit, audience_relevance, strength_signal, source_event_id, source_burst_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, fact.ID, fact.Text, competenciesString, string(fact.RoleFit), audienceString, fact.StrengthSignal, fact.SourceEventID, fact.SourceBurstID, fact.CreatedAt, fact.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create fact with ID %s: %w", fact.ID, err)
	}

	return nil
}

// GetByID retrieves a fact by its ID
func (r *SQLiteFactRepository) GetByID(ctx context.Context, id string) (*career.Fact, error) {
	var fact career.Fact
	var competenciesString, audienceString string

	row := r.db.QueryRowContext(ctx, `
		SELECT id, text, competencies, role_fit, audience_relevance, strength_signal, source_event_id, source_burst_id, created_at, updated_at
		FROM facts
		WHERE id = ?
	`, id)

	err := row.Scan(
		&fact.ID,
		&fact.Text,
		&competenciesString,
		&fact.RoleFit,
		&audienceString,
		&fact.StrengthSignal,
		&fact.SourceEventID,
		&fact.SourceBurstID,
		&fact.CreatedAt,
		&fact.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("fact with ID %s not found: %w", id, ErrFactNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve fact with ID %s: %w", id, err)
	}

	// Parse competencies and audience
	if competenciesString != "" {
		fact.CompetencyCategories = strings.Split(competenciesString, ",")
	}
	if audienceString != "" {
		fact.AudienceRelevance = strings.Split(audienceString, ",")
	}

	return &fact, nil
}

// Update modifies an existing fact
func (r *SQLiteFactRepository) Update(ctx context.Context, fact *career.Fact) error {
	// Validate the fact
	if err := fact.Validate(); err != nil {
		return fmt.Errorf("fact validation failed for ID %s: %w", fact.ID, err)
	}

	// Check if fact exists
	var exists int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM facts WHERE id = ?", fact.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if fact exists: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("fact with ID %s: %w", fact.ID, ErrFactNotFound)
	}

	// Update timestamp
	fact.UpdatedAt = time.Now()

	// Convert slices to strings
	competenciesString := strings.Join(fact.CompetencyCategories, ",")
	audienceString := strings.Join(fact.AudienceRelevance, ",")

	// Update the fact
	_, err = r.db.ExecContext(ctx, `
		UPDATE facts
		SET text = ?, competencies = ?, role_fit = ?, audience_relevance = ?, strength_signal = ?, source_event_id = ?, source_burst_id = ?, updated_at = ?
		WHERE id = ?
	`, fact.Text, competenciesString, string(fact.RoleFit), audienceString, fact.StrengthSignal, fact.SourceEventID, fact.SourceBurstID, fact.UpdatedAt, fact.ID)

	if err != nil {
		return fmt.Errorf("failed to update fact with ID %s: %w", fact.ID, err)
	}

	return nil
}

// Delete removes a fact from the database
func (r *SQLiteFactRepository) Delete(ctx context.Context, id string) error {
	// Check if fact exists
	var exists int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM facts WHERE id = ?", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if fact exists: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("fact with ID %s: %w", id, ErrFactNotFound)
	}

	// Delete the fact
	_, err = r.db.ExecContext(ctx, "DELETE FROM facts WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete fact with ID %s: %w", id, err)
	}

	return nil
}

// List retrieves facts with optional filtering
func (r *SQLiteFactRepository) List(ctx context.Context, filters FactListFilters) ([]*career.Fact, error) {
	query := "SELECT id, text, competencies, role_fit, audience_relevance, strength_signal, source_event_id, source_burst_id, created_at, updated_at FROM facts WHERE 1=1"
	var args []interface{}

	// Apply competency category filter
	if filters.CompetencyCategory != "" {
		query += " AND competencies LIKE ?"
		args = append(args, "%"+filters.CompetencyCategory+"%")
	}

	// Apply role fit filter
	if filters.RoleFit != "" {
		query += " AND role_fit = ?"
		args = append(args, filters.RoleFit)
	}

	// Apply audience relevance filter
	if filters.AudienceRelevance != "" {
		query += " AND audience_relevance LIKE ?"
		args = append(args, "%"+filters.AudienceRelevance+"%")
	}

	// Apply date range filters
	if filters.StartDate != nil {
		query += " AND created_at >= ?"
		args = append(args, filters.StartDate)
	}
	if filters.EndDate != nil {
		query += " AND created_at <= ?"
		args = append(args, filters.EndDate)
	}

	// Apply sorting
	sortBy := filters.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	sortOrder := filters.SortOrder
	if sortOrder == "" {
		sortOrder = "DESC"
	} else if strings.ToUpper(sortOrder) == "ASC" {
		sortOrder = "ASC"
	} else {
		sortOrder = "DESC"
	}

	switch sortBy {
	case "text":
		query += fmt.Sprintf(" ORDER BY text %s", sortOrder)
	case "created_at":
		fallthrough
	default:
		query += fmt.Sprintf(" ORDER BY created_at %s", sortOrder)
	}

	// Apply pagination
	if filters.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filters.Limit)
	}
	if filters.Offset > 0 {
		query += " OFFSET ?"
		args = append(args, filters.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list facts: %w", err)
	}
	defer rows.Close()

	var facts []*career.Fact
	for rows.Next() {
		var fact career.Fact
		var competenciesString, audienceString string

		err := rows.Scan(
			&fact.ID,
			&fact.Text,
			&competenciesString,
			&fact.RoleFit,
			&audienceString,
			&fact.StrengthSignal,
			&fact.SourceEventID,
			&fact.SourceBurstID,
			&fact.CreatedAt,
			&fact.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fact: %w", err)
		}

		// Parse competencies and audience
		if competenciesString != "" {
			fact.CompetencyCategories = strings.Split(competenciesString, ",")
		}
		if audienceString != "" {
			fact.AudienceRelevance = strings.Split(audienceString, ",")
		}

		facts = append(facts, &fact)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating fact rows: %w", err)
	}

	return facts, nil
}

// Count returns the total number of facts matching the given filters
func (r *SQLiteFactRepository) Count(ctx context.Context, filters FactListFilters) (int, error) {
	query := "SELECT COUNT(*) FROM facts WHERE 1=1"
	var args []interface{}

	// Apply competency category filter
	if filters.CompetencyCategory != "" {
		query += " AND competencies LIKE ?"
		args = append(args, "%"+filters.CompetencyCategory+"%")
	}

	// Apply role fit filter
	if filters.RoleFit != "" {
		query += " AND role_fit = ?"
		args = append(args, filters.RoleFit)
	}

	// Apply audience relevance filter
	if filters.AudienceRelevance != "" {
		query += " AND audience_relevance LIKE ?"
		args = append(args, "%"+filters.AudienceRelevance+"%")
	}

	// Apply date range filters
	if filters.StartDate != nil {
		query += " AND created_at >= ?"
		args = append(args, filters.StartDate)
	}
	if filters.EndDate != nil {
		query += " AND created_at <= ?"
		args = append(args, filters.EndDate)
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count facts: %w", err)
	}

	return count, nil
}

// GetBySourceEventID retrieves all facts for a specific event
func (r *SQLiteFactRepository) GetBySourceEventID(ctx context.Context, eventID string) ([]*career.Fact, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, text, competencies, role_fit, audience_relevance, strength_signal, source_event_id, source_burst_id, created_at, updated_at
		FROM facts
		WHERE source_event_id = ?
	`, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to query facts by event: %w", err)
	}
	defer rows.Close()

	var facts []*career.Fact
	for rows.Next() {
		var fact career.Fact
		var competenciesString, audienceString string

		err := rows.Scan(
			&fact.ID,
			&fact.Text,
			&competenciesString,
			&fact.RoleFit,
			&audienceString,
			&fact.StrengthSignal,
			&fact.SourceEventID,
			&fact.SourceBurstID,
			&fact.CreatedAt,
			&fact.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fact: %w", err)
		}

		// Parse competencies and audience
		if competenciesString != "" {
			fact.CompetencyCategories = strings.Split(competenciesString, ",")
		}
		if audienceString != "" {
			fact.AudienceRelevance = strings.Split(audienceString, ",")
		}

		facts = append(facts, &fact)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating fact rows: %w", err)
	}

	return facts, nil
}

// GetBySourceBurstID retrieves all facts for a specific burst
func (r *SQLiteFactRepository) GetBySourceBurstID(ctx context.Context, burstID string) ([]*career.Fact, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, text, competencies, role_fit, audience_relevance, strength_signal, source_event_id, source_burst_id, created_at, updated_at
		FROM facts
		WHERE source_burst_id = ?
	`, burstID)
	if err != nil {
		return nil, fmt.Errorf("failed to query facts by burst: %w", err)
	}
	defer rows.Close()

	var facts []*career.Fact
	for rows.Next() {
		var fact career.Fact
		var competenciesString, audienceString string

		err := rows.Scan(
			&fact.ID,
			&fact.Text,
			&competenciesString,
			&fact.RoleFit,
			&audienceString,
			&fact.StrengthSignal,
			&fact.SourceEventID,
			&fact.SourceBurstID,
			&fact.CreatedAt,
			&fact.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fact: %w", err)
		}

		// Parse competencies and audience
		if competenciesString != "" {
			fact.CompetencyCategories = strings.Split(competenciesString, ",")
		}
		if audienceString != "" {
			fact.AudienceRelevance = strings.Split(audienceString, ",")
		}

		facts = append(facts, &fact)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating fact rows: %w", err)
	}

	return facts, nil
}
