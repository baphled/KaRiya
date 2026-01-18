// Package testutil provides reusable test utilities for KaRiya tests.
package testutil

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/baphled/kariya/internal/repository/career"
	_ "modernc.org/sqlite"
)

// SetupTestDB creates a test database with all migrations applied.
// Returns the database connection and a cleanup function.
// The cleanup function closes the database connection and removes the temporary directory.
//
// Usage:
//
//	db, cleanup := testutil.SetupTestDB(t)
//	defer cleanup()
//	// use db for testing
func SetupTestDB(t testing.TB) (*sql.DB, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// Use RunMigrationsForTests which skips baseline detection for fresh test databases
	if err := career.RunMigrationsForTests(db); err != nil {
		closeErr := db.Close()
		if closeErr != nil {
			t.Fatalf("failed to run test migrations: %v (also failed to close db: %v)", err, closeErr)
		}
		t.Fatalf("failed to run test migrations: %v", err)
	}

	cleanup := func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close test db: %v", err)
		}
	}

	return db, cleanup
}

// SetupTestDBWithPath creates a test database at a specific path and returns
// both the path and the database connection. This is useful when constructors
// need the database path in addition to the connection.
//
// Returns the database path, connection, and cleanup function.
//
// Usage:
//
//	dbPath, db, cleanup := testutil.SetupTestDBWithPath(t)
//	defer cleanup()
//	// use dbPath and db for testing
func SetupTestDBWithPath(t testing.TB) (string, *sql.DB, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open test db at %s: %v", dbPath, err)
	}

	// Use RunMigrationsForTests which skips baseline detection for fresh test databases
	if err := career.RunMigrationsForTests(db); err != nil {
		closeErr := db.Close()
		if closeErr != nil {
			t.Fatalf("failed to run test migrations: %v (also failed to close db: %v)", err, closeErr)
		}
		t.Fatalf("failed to run test migrations: %v", err)
	}

	cleanup := func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close test db: %v", err)
		}
	}

	return dbPath, db, cleanup
}
