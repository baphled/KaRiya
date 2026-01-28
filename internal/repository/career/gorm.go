// Package career provides repository implementations for career domain entities.
package career

import (
	"database/sql"

	// SQLite driver for database/sql (pure Go, no CGO required).
	_ "modernc.org/sqlite"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewGormDB creates a GORM database connection from an existing sql.DB.
// This allows sharing the same database connection between GORM and raw SQL.
func NewGormDB(sqlDB *sql.DB) (*gorm.DB, error) {
	return gorm.Open(sqlite.New(sqlite.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

// Repositories holds all ORM-based repositories.
type Repositories struct {
	Event *Event
	Skill *Skill
	Fact  *Fact
	Burst *Burst

	// db holds the underlying database connection for closing.
	db     *gorm.DB
	sqlDB  *sql.DB
	ownSQL bool // true if we created the sql.DB and should close it
}

// Close closes the underlying database connection.
func (r *Repositories) Close() error {
	if r.ownSQL && r.sqlDB != nil {
		return r.sqlDB.Close()
	}
	return nil
}

// NewRepositories creates all ORM repositories from a database connection.
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Event:  NewEvent(db),
		Skill:  NewSkill(db),
		Fact:   NewFact(db),
		Burst:  NewBurst(db),
		db:     db,
		ownSQL: false,
	}
}

// NewRepositoriesFromSQL creates all ORM repositories from an existing sql.DB.
// This is a convenience function that wraps NewGormDB and NewRepositories.
func NewRepositoriesFromSQL(sqlDB *sql.DB) (*Repositories, error) {
	gormDB, err := NewGormDB(sqlDB)
	if err != nil {
		return nil, err
	}
	repos := NewRepositories(gormDB)
	repos.sqlDB = sqlDB
	repos.ownSQL = false // Caller owns the sql.DB
	return repos, nil
}

// NewRepositoriesFromPath creates all ORM repositories from a database file path.
// This opens a new database connection and runs migrations.
// The caller should call Close() when done.
func NewRepositoriesFromPath(dbPath string) (*Repositories, error) {
	// Open database connection
	sqlDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := RunMigrations(sqlDB); err != nil {
		_ = sqlDB.Close() // Best effort close on error path
		return nil, err
	}

	// Create ORM connection
	gormDB, err := NewGormDB(sqlDB)
	if err != nil {
		_ = sqlDB.Close() // Best effort close on error path
		return nil, err
	}

	repos := NewRepositories(gormDB)
	repos.sqlDB = sqlDB
	repos.ownSQL = true // We own the sql.DB
	return repos, nil
}
