package career

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

// RunMigrations executes all pending database migrations.
// For existing databases (pre-goose), it auto-detects and marks baseline migrations as applied.
// This ensures smooth migration from the old schema management approach to goose-based migrations.
func RunMigrations(db *sql.DB) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Check if this is a pre-goose database (has tables but no goose tracking)
	if err := handleBaselineMigration(db); err != nil {
		return fmt.Errorf("failed to handle baseline: %w", err)
	}

	// Run any pending migrations
	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// handleBaselineMigration detects existing databases and marks baseline migrations as applied.
// This prevents goose from attempting to re-create tables that already exist.
//
// Detection logic:
// - If goose_db_version table exists, goose is already tracking migrations -> do nothing
// - If career_events table doesn't exist, this is a fresh database -> do nothing
// - If career_events exists but goose isn't tracking, detect which migrations to mark as applied:
//   - Version 1: career_events table exists
//   - Version 2: categories column exists
//   - Version 3: bursts table exists
//   - Version 4: facts table exists
func handleBaselineMigration(db *sql.DB) error {
	// Check if goose version table exists
	var gooseTableExists int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM sqlite_master 
		WHERE type='table' AND name='goose_db_version'
	`).Scan(&gooseTableExists)
	if err != nil {
		return fmt.Errorf("failed to check goose table: %w", err)
	}

	if gooseTableExists > 0 {
		// Goose is already initialized, nothing to do
		return nil
	}

	// Check if career_events table exists (indicates pre-goose database)
	var careerEventsExists int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM sqlite_master 
		WHERE type='table' AND name='career_events'
	`).Scan(&careerEventsExists)
	if err != nil {
		return fmt.Errorf("failed to check career_events table: %w", err)
	}

	if careerEventsExists == 0 {
		// Fresh database, no baseline needed
		return nil
	}

	// This is a pre-goose database - determine what migrations to mark as applied
	baselineVersion := int64(1) // At minimum, career_events exists

	// Check for categories column (migration 002)
	if hasColumn(db, "career_events", "categories") {
		baselineVersion = 2
	}

	// Check for bursts table (migration 003)
	if hasTable(db, "bursts") {
		baselineVersion = 3
	}

	// Check for facts table (migration 004)
	if hasTable(db, "facts") {
		baselineVersion = 4
	}

	// Create the goose version table manually
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS goose_db_version (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version_id INTEGER NOT NULL,
			is_applied INTEGER NOT NULL,
			tstamp TIMESTAMP DEFAULT (datetime('now'))
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create goose_db_version table: %w", err)
	}

	// Insert version records for each baseline migration
	// These migrations are already applied (tables exist), so we just record them
	for v := int64(1); v <= baselineVersion; v++ {
		_, err := db.Exec(`
			INSERT INTO goose_db_version (version_id, is_applied)
			VALUES (?, 1)
		`, v)
		if err != nil {
			return fmt.Errorf("failed to mark version %d as applied: %w", v, err)
		}
	}

	return nil
}

// hasTable checks if a table exists in the database
func hasTable(db *sql.DB, tableName string) bool {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM sqlite_master 
		WHERE type='table' AND name=?
	`, tableName).Scan(&count)
	return err == nil && count > 0
}

// hasColumn checks if a column exists in a table using PRAGMA table_info
func hasColumn(db *sql.DB, tableName, columnName string) bool {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, typ string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &typ, &notnull, &dflt, &pk); err != nil {
			continue
		}
		if name == columnName {
			return true
		}
	}
	return false
}

// MigrationStatus returns the current migration version for debugging and status reporting.
// Returns 0 if no migrations have been applied yet.
func MigrationStatus(db *sql.DB) (int64, error) {
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return 0, err
	}
	return goose.GetDBVersion(db)
}
