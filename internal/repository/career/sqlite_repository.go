package career

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"

	domain "github.com/baphled/kariya/internal/domain/career"
)

// SQLiteRepository provides a SQLite implementation of the Repository interface
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite-backed repository.
// Deprecated: Use NewSQLiteRepositoryWithDB after running migrations via career.RunMigrations().
// This constructor is kept for backward compatibility with existing tests.
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	// Open the SQLite database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Run migrations to ensure schema is up to date
	if err := RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &SQLiteRepository{db: db}, nil
}

// NewSQLiteRepositoryWithDB creates a repository with an existing database connection.
// This is useful when migrations have already been run on the database connection.
// The caller is responsible for managing the database connection lifecycle.
func NewSQLiteRepositoryWithDB(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

// GetDB returns the underlying database connection for sharing with other repositories
func (r *SQLiteRepository) GetDB() *sql.DB {
	return r.db
}

// Create adds a new career event to the SQLite database
func (r *SQLiteRepository) Create(ctx context.Context, event *domain.CareerEvent) error {
	// Validate the event
	if err := event.Validate(); err != nil {
		return fmt.Errorf("event validation failed: %w", err)
	}

	// Generate UUID if not provided
	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now

	// Convert tags to comma-separated strings
	tagString := ""
	if len(event.Tags) > 0 {
		tagString = event.Tags[0]
		for _, tag := range event.Tags[1:] {
			tagString += "," + tag
		}
	}

	// Convert categories to comma-separated strings
	categoriesString := ""
	if len(event.Categories) > 0 {
		categoriesString = event.Categories[0]
		for _, category := range event.Categories[1:] {
			categoriesString += "," + category
		}
	}

	// Check for duplicate event
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM career_events WHERE id = ?", event.ID).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking for duplicate event: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("duplicate event with ID %s: %w", event.ID, ErrDuplicateEvent)
	}

	// Insert the event
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO career_events
		(id, text, date, tags, categories, company, project, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, event.ID, event.Text, event.Date, tagString, categoriesString, event.Company, event.Project, event.CreatedAt, event.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create event with ID %s: %w", event.ID, err)
	}

	return nil
}

// GetByID retrieves a career event by its ID
func (r *SQLiteRepository) GetByID(ctx context.Context, id string) (*domain.CareerEvent, error) {
	var event domain.CareerEvent
	var tagString, categoriesString sql.NullString

	row := r.db.QueryRowContext(ctx, `
		SELECT id, text, date, tags, categories, company, project, created_at, updated_at
		FROM career_events
		WHERE id = ?
	`, id)

	err := row.Scan(
		&event.ID,
		&event.Text,
		&event.Date,
		&tagString,
		&categoriesString,
		&event.Company,
		&event.Project,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("event with ID %s not found: %w", id, ErrEventNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve event with ID %s: %w", id, err)
	}

	// Parse tags
	if tagString.Valid && tagString.String != "" {
		event.Tags = splitStringList(tagString.String)
	}

	// Parse categories
	if categoriesString.Valid && categoriesString.String != "" {
		event.Categories = parseCategories(categoriesString.String)
	}

	return &event, nil
}

// Update modifies an existing career event
func (r *SQLiteRepository) Update(ctx context.Context, event *domain.CareerEvent) error {
	// Validate the event
	if err := event.Validate(); err != nil {
		return fmt.Errorf("event validation failed for ID %s: %w", event.ID, err)
	}

	// Check if event exists
	var exists int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM career_events WHERE id = ?", event.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking event existence for ID %s: %w", event.ID, err)
	}
	if exists == 0 {
		return fmt.Errorf("event with ID %s not found: %w", event.ID, ErrEventNotFound)
	}

	// Convert tags to comma-separated strings
	tagString := ""
	if len(event.Tags) > 0 {
		tagString = event.Tags[0]
		for _, tag := range event.Tags[1:] {
			tagString += "," + tag
		}
	}

	// Convert categories to comma-separated strings
	categoriesString := ""
	if len(event.Categories) > 0 {
		categoriesString = event.Categories[0]
		for _, category := range event.Categories[1:] {
			categoriesString += "," + category
		}
	}

	// Update the event
	updatedTime := time.Now()
	event.UpdatedAt = updatedTime
	_, err = r.db.ExecContext(ctx, `
		UPDATE career_events
		SET text = ?, date = ?, tags = ?, categories = ?, company = ?, project = ?, updated_at = ?
		WHERE id = ?
	`, event.Text, event.Date, tagString, categoriesString, event.Company, event.Project, event.UpdatedAt, event.ID)

	if err != nil {
		return fmt.Errorf("failed to update event with ID %s: %w", event.ID, err)
	}

	return nil
}

// Delete removes a career event from the repository
func (r *SQLiteRepository) Delete(ctx context.Context, id string) error {
	// Check if event exists
	var exists int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM career_events WHERE id = ?", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking event existence for ID %s: %w", id, err)
	}
	if exists == 0 {
		return fmt.Errorf("event with ID %s not found: %w", id, ErrEventNotFound)
	}

	// Delete the event
	_, err = r.db.ExecContext(ctx, "DELETE FROM career_events WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete event with ID %s: %w", id, err)
	}

	return nil
}

// List retrieves career events with optional filtering
func (r *SQLiteRepository) List(ctx context.Context, filters ListFilters) ([]*domain.CareerEvent, error) {
	// Build query with dynamic filtering
	query := "SELECT id, text, date, tags, categories, company, project, created_at, updated_at FROM career_events WHERE 1=1"
	args := []interface{}{}

	// Tag filtering
	if len(filters.Tags) > 0 {
		// Use parameterized query to prevent SQL injection
		query += " AND ("
		for i, tag := range filters.Tags {
			if i > 0 {
				query += " OR "
			}
			query += "tags LIKE ?"
			args = append(args, "%"+tag+"%")
		}
		query += ")"
	}

	// Date range filtering
	if filters.StartDate != nil {
		query += " AND date >= ?"
		args = append(args, filters.StartDate)
	}
	if filters.EndDate != nil {
		query += " AND date <= ?"
		args = append(args, filters.EndDate)
	}

	// Sorting with secondary sort
	switch filters.SortBy {
	case "date":
		query += " ORDER BY date"
		if filters.SortOrder == "desc" {
			query += " DESC"
		}
		// Secondary sort by created_at
		query += ", created_at"
		if filters.SortOrder == "desc" {
			query += " DESC"
		}
	default:
		query += " ORDER BY created_at"
		if filters.SortOrder == "desc" {
			query += " DESC"
		}
	}

	// Pagination
	query += " LIMIT ? OFFSET ?"
	limit := filters.Limit
	if limit == 0 {
		limit = 100 // Default limit
	}
	args = append(args, limit, filters.Offset)

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}
	defer rows.Close()

	var events []*domain.CareerEvent
	for rows.Next() {
		var event domain.CareerEvent
		var tagString, categoriesString sql.NullString

		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&event.ID,
			&event.Text,
			&event.Date,
			&tagString,
			&categoriesString,
			&event.Company,
			&event.Project,
			&createdAt,
			&updatedAt,
		)
		event.CreatedAt = createdAt
		event.UpdatedAt = updatedAt
		if err != nil {
			return nil, fmt.Errorf("error scanning event: %w", err)
		}

		// Parse tags
		if tagString.Valid && tagString.String != "" {
			event.Tags = splitStringList(tagString.String)
		}

		// Parse categories
		if categoriesString.Valid && categoriesString.String != "" {
			event.Categories = parseCategories(categoriesString.String)
		}

		events = append(events, &event)
	}

	return events, nil
}

// Count returns the number of events matching the given filters
func (r *SQLiteRepository) Count(ctx context.Context, filters ListFilters) (int, error) {
	// Build query with dynamic filtering
	query := "SELECT COUNT(*) FROM career_events WHERE 1=1"
	args := []interface{}{}

	// Tag filtering
	if len(filters.Tags) > 0 {
		// Use parameterized query to prevent SQL injection
		query += " AND ("
		for i, tag := range filters.Tags {
			if i > 0 {
				query += " OR "
			}
			query += "tags LIKE ?"
			args = append(args, "%"+tag+"%")
		}
		query += ")"
	}

	// Date range filtering
	if filters.StartDate != nil {
		query += " AND date >= ?"
		args = append(args, filters.StartDate)
	}
	if filters.EndDate != nil {
		query += " AND date <= ?"
		args = append(args, filters.EndDate)
	}

	// Execute query
	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count events: %w", err)
	}

	return count, nil
}

// Close closes the database connection
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

// Helper function to split comma-separated string into slice
func splitStringList(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

// Helper function to parse categories from string
func parseCategories(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}
