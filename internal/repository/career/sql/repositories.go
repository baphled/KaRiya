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

// defaultPaginationLimit caps the number of rows returned when no explicit limit is set.
const defaultPaginationLimit = 100

// NewGormDB creates a GORM database connection from an existing sql.DB.
//
// Expected:
//   - sqlDB must be a valid sql.DB connection.
//
// Returns:
//   - A GORM database connection, or an error if setup fails.
//
// Side effects:
//   - None.
func NewGormDB(sqlDB *sql.DB) (*gorm.DB, error) {
	return gorm.Open(sqlite.New(sqlite.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

// NewRepositoriesFromDB creates all SQL repositories from a GORM database connection.
//
// Expected:
//   - db must be valid.
//
// Returns:
//   - A fully initialized career.Repositories ready for use.
//
// Side effects:
//   - None.
func NewRepositoriesFromDB(db *gorm.DB) *career.Repositories {
	return &career.Repositories{
		Event: NewEventRepository(db),
		Skill: NewSkillRepository(db),
		Fact:  NewFactRepository(db),
		Burst: NewBurstRepository(db),
	}
}

// NewRepositories creates all SQL repositories from an existing sql.DB.
//
// Expected:
//   - sqlDB must be a valid sql.DB connection.
//
// Returns:
//   - A fully initialized career.Repositories, or an error.
//
// Side effects:
//   - None.
func NewRepositories(sqlDB *sql.DB) (*career.Repositories, error) {
	gormDB, err := NewGormDB(sqlDB)
	if err != nil {
		return nil, err
	}
	repos := NewRepositoriesFromDB(gormDB)
	return repos, nil
}

// OpenDB opens a SQLite database with WAL journal mode, a 5-second busy
// timeout, and a single-connection pool. These settings prevent
// SQLITE_BUSY errors on Windows where file locking is stricter.
//
// Expected:
//   - dbPath must be a valid file path.
//
// Returns:
//   - An open sql.DB connection, or an error if open fails.
//
// Side effects:
//   - Creates the database file if it does not exist.
func OpenDB(dbPath string) (*sql.DB, error) {
	dsn := dbPath + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	// SQLite does not support concurrent writers; serialise all access.
	db.SetMaxOpenConns(1)

	return db, nil
}

// NewRepositoriesFromPath creates all SQL repositories from a database file path.
// This opens a new database connection and runs migrations.
// The caller should call Close() when done.
//
// Expected:
//   - dbPath must be a valid file path.
//
// Returns:
//   - A fully initialized career.Repositories, or an error.
//
// Side effects:
//   - Creates the database file and runs migrations if needed.
func NewRepositoriesFromPath(dbPath string) (*career.Repositories, error) {
	sqlDB, err := OpenDB(dbPath)
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
