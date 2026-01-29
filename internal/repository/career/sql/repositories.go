package sql

import (
	"database/sql"

	career "github.com/baphled/kariya/internal/repository/career"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// Register the SQLite driver for database/sql.Open("sqlite", ...) calls.
	_ "modernc.org/sqlite"
)

// NewGormDB creates a GORM database connection from an existing sql.DB.
func NewGormDB(sqlDB *sql.DB) (*gorm.DB, error) {
	return gorm.Open(sqlite.New(sqlite.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

// NewRepositoriesFromDB creates all SQL repositories from a GORM database connection.
func NewRepositoriesFromDB(db *gorm.DB) *career.Repositories {
	return &career.Repositories{
		Event: NewEventRepository(db),
		Skill: NewSkillRepository(db),
		Fact:  NewFactRepository(db),
		Burst: NewBurstRepository(db),
	}
}

// NewRepositories creates all SQL repositories from an existing sql.DB.
func NewRepositories(sqlDB *sql.DB) (*career.Repositories, error) {
	gormDB, err := NewGormDB(sqlDB)
	if err != nil {
		return nil, err
	}
	repos := NewRepositoriesFromDB(gormDB)
	return repos, nil
}

// NewRepositoriesFromPath creates all SQL repositories from a database file path.
// This opens a new database connection and runs migrations.
// The caller should call Close() when done.
func NewRepositoriesFromPath(dbPath string) (*career.Repositories, error) {
	sqlDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := career.RunMigrations(sqlDB); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	gormDB, err := NewGormDB(sqlDB)
	if err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	repos := NewRepositoriesFromDB(gormDB)
	repos.SetCloser(sqlDB)
	return repos, nil
}
