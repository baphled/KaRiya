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

// SQLiteBurstRepository provides a SQLite implementation of BurstRepository
type SQLiteBurstRepository struct {
	db *sql.DB
}

// NewSQLiteBurstRepository creates a new SQLite burst repository
func NewSQLiteBurstRepository(db *sql.DB) (*SQLiteBurstRepository, error) {
	repo := &SQLiteBurstRepository{db: db}

	// Initialize the schema
	if err := repo.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize burst schema: %w", err)
	}

	return repo, nil
}

// initSchema creates the bursts table if it doesn't exist
func (r *SQLiteBurstRepository) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS bursts (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		event_ids TEXT NOT NULL,
		competency_focus TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	`

	_, err := r.db.Exec(schema)
	return err
}

// Create adds a new burst to the database
func (r *SQLiteBurstRepository) Create(ctx context.Context, burst *career.Burst) error {
	// Generate a unique ID if not provided
	if burst.ID == "" {
		burst.ID = uuid.New().String()
	}

	// Validate the burst
	if err := burst.Validate(); err != nil {
		return fmt.Errorf("burst validation failed for ID %s: %w", burst.ID, err)
	}

	// Check for duplicate
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM bursts WHERE id = ?", burst.ID).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking for duplicate burst: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("duplicate burst with ID %s: %w", burst.ID, ErrDuplicateBurst)
	}

	// Set timestamps
	now := time.Now()
	burst.CreatedAt = now
	burst.UpdatedAt = now

	// Convert event IDs to string
	eventIDsString := strings.Join(burst.EventIDs, ",")

	// Insert the burst
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO bursts
		(id, name, description, event_ids, competency_focus, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, burst.ID, burst.Name, burst.Description, eventIDsString, burst.CompetencyFocus, burst.CreatedAt, burst.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create burst with ID %s: %w", burst.ID, err)
	}

	return nil
}

// GetByID retrieves a burst by its ID
func (r *SQLiteBurstRepository) GetByID(ctx context.Context, id string) (*career.Burst, error) {
	var burst career.Burst
	var eventIDsString string

	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, description, event_ids, competency_focus, created_at, updated_at
		FROM bursts
		WHERE id = ?
	`, id)

	err := row.Scan(
		&burst.ID,
		&burst.Name,
		&burst.Description,
		&eventIDsString,
		&burst.CompetencyFocus,
		&burst.CreatedAt,
		&burst.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("burst with ID %s not found: %w", id, ErrBurstNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve burst with ID %s: %w", id, err)
	}

	// Parse event IDs
	if eventIDsString != "" {
		burst.EventIDs = strings.Split(eventIDsString, ",")
	}

	return &burst, nil
}

// Update modifies an existing burst
func (r *SQLiteBurstRepository) Update(ctx context.Context, burst *career.Burst) error {
	// Validate the burst
	if err := burst.Validate(); err != nil {
		return fmt.Errorf("burst validation failed for ID %s: %w", burst.ID, err)
	}

	// Check if burst exists
	var exists int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM bursts WHERE id = ?", burst.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if burst exists: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("burst with ID %s: %w", burst.ID, ErrBurstNotFound)
	}

	// Update timestamp
	burst.UpdatedAt = time.Now()

	// Convert event IDs to string
	eventIDsString := strings.Join(burst.EventIDs, ",")

	// Update the burst
	_, err = r.db.ExecContext(ctx, `
		UPDATE bursts
		SET name = ?, description = ?, event_ids = ?, competency_focus = ?, updated_at = ?
		WHERE id = ?
	`, burst.Name, burst.Description, eventIDsString, burst.CompetencyFocus, burst.UpdatedAt, burst.ID)

	if err != nil {
		return fmt.Errorf("failed to update burst with ID %s: %w", burst.ID, err)
	}

	return nil
}

// Delete removes a burst from the database
func (r *SQLiteBurstRepository) Delete(ctx context.Context, id string) error {
	// Check if burst exists
	var exists int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM bursts WHERE id = ?", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if burst exists: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("burst with ID %s: %w", id, ErrBurstNotFound)
	}

	// Delete the burst
	_, err = r.db.ExecContext(ctx, "DELETE FROM bursts WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete burst with ID %s: %w", id, err)
	}

	return nil
}

// List retrieves bursts with optional filtering
func (r *SQLiteBurstRepository) List(ctx context.Context, filters BurstListFilters) ([]*career.Burst, error) {
	query := "SELECT id, name, description, event_ids, competency_focus, created_at, updated_at FROM bursts WHERE 1=1"
	var args []interface{}

	// Apply competency focus filter
	if filters.CompetencyFocus != "" {
		query += " AND competency_focus = ?"
		args = append(args, filters.CompetencyFocus)
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
	case "name":
		query += fmt.Sprintf(" ORDER BY name %s", sortOrder)
	case "event_count":
		query += fmt.Sprintf(" ORDER BY LENGTH(event_ids) - LENGTH(REPLACE(event_ids, ',', '')) + 1 %s", sortOrder)
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
		return nil, fmt.Errorf("failed to list bursts: %w", err)
	}
	defer rows.Close()

	var bursts []*career.Burst
	for rows.Next() {
		var burst career.Burst
		var eventIDsString string

		err := rows.Scan(
			&burst.ID,
			&burst.Name,
			&burst.Description,
			&eventIDsString,
			&burst.CompetencyFocus,
			&burst.CreatedAt,
			&burst.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan burst: %w", err)
		}

		// Parse event IDs
		if eventIDsString != "" {
			burst.EventIDs = strings.Split(eventIDsString, ",")
		}

		bursts = append(bursts, &burst)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating burst rows: %w", err)
	}

	return bursts, nil
}

// Count returns the total number of bursts matching the given filters
func (r *SQLiteBurstRepository) Count(ctx context.Context, filters BurstListFilters) (int, error) {
	query := "SELECT COUNT(*) FROM bursts WHERE 1=1"
	var args []interface{}

	// Apply competency focus filter
	if filters.CompetencyFocus != "" {
		query += " AND competency_focus = ?"
		args = append(args, filters.CompetencyFocus)
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
		return 0, fmt.Errorf("failed to count bursts: %w", err)
	}

	return count, nil
}
