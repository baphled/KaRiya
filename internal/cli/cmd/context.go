package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careersql "github.com/baphled/kariya/internal/repository/career/sql"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// CLIContext holds shared state for CLI commands.
//
// It manages database configuration and service initialization for all CLI operations.
// Supports both in-memory and SQLite persistence modes.
type CLIContext struct {
	DBPath   string
	InMemory bool
	Service  *careerservice.Service
}

// NewCLIContext creates a new CLIContext with the specified configuration.
//
// Parameters:
//   - dbPath: Path to SQLite database file (ignored if inMemory is true)
//   - inMemory: If true, uses in-memory repositories; otherwise uses SQLite
//
// Returns:
//   - A new CLIContext with uninitialized service (call InitService to initialize)
//
// Side effects:
//   - None; service initialization is deferred to InitService
func NewCLIContext(dbPath string, inMemory bool) *CLIContext {
	return &CLIContext{
		DBPath:   dbPath,
		InMemory: inMemory,
	}
}

// InitService initializes the career service with appropriate repositories.
//
// Initializes either in-memory or SQLite repositories based on configuration.
// For SQLite mode, creates the database directory if needed and runs migrations.
//
// Parameters:
//   - errOut: io.Writer for error messages (allows testable error output)
//
// Returns:
//   - nil on success
//   - error if database initialization, migration, or repository setup fails
//
// Side effects:
//   - Sets ctx.Service to initialized service
//   - For SQLite mode: creates ~/.kariya directory if it doesn't exist
//   - For SQLite mode: creates/opens database file
//   - Writes error details to errOut on failure
func (ctx *CLIContext) InitService(errOut io.Writer) error {
	if ctx.InMemory {
		ctx.initInMemory()
		return nil
	}
	return ctx.initSQL(errOut)
}

func (ctx *CLIContext) initInMemory() {
	eventRepo := careermemory.NewEventRepository()
	skillRepo := careermemory.NewSkillRepository()
	skillRepo.SetEventRepository(eventRepo)
	eventRepo.SetSkillRepository(skillRepo)

	svc := careerservice.NewService(eventRepo)
	svc.SetFactRepository(careermemory.NewFactRepository())
	svc.SetBurstRepository(careermemory.NewBurstRepository())
	svc.SetSkillRepository(skillRepo)
	ctx.Service = svc
}

func (ctx *CLIContext) initSQL(errOut io.Writer) error {
	dbPath, err := ctx.getDBPath()
	if err != nil {
		fmt.Fprintf(errOut, "Error getting database path: %v\n", err)
		return err
	}

	db, err := careersql.OpenDB(dbPath)
	if err != nil {
		fmt.Fprintf(errOut, "Error opening database at '%s': %v\n", dbPath, err)
		return fmt.Errorf("opening database at '%s': %w", dbPath, err)
	}

	if err := career.RunMigrations(db); err != nil {
		fmt.Fprintf(errOut, "Error running migrations: %v\n", err)
		return fmt.Errorf("running migrations: %w", err)
	}

	repos, err := careersql.NewRepositories(db)
	if err != nil {
		fmt.Fprintf(errOut, "Error initializing repositories: %v\n", err)
		return fmt.Errorf("initializing GORM repositories: %w", err)
	}

	svc := careerservice.NewService(repos.Event)
	svc.SetFactRepository(repos.Fact)
	svc.SetBurstRepository(repos.Burst)
	svc.SetSkillRepository(repos.Skill)
	ctx.Service = svc
	return nil
}

func (ctx *CLIContext) getDBPath() (string, error) {
	dbPath := ctx.DBPath
	if dbPath != "" {
		return dbPath, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	kariyaDir := filepath.Join(homeDir, ".kariya")
	dbPath = filepath.Join(kariyaDir, "events.db")

	if err := os.MkdirAll(kariyaDir, 0o750); err != nil {
		return "", fmt.Errorf("creating kariya directory: %w", err)
	}
	return dbPath, nil
}
