package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careersql "github.com/baphled/kariya/internal/repository/career/sql"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// CLIContext holds shared configuration and service instances for Cobra commands.
//
// It manages both in-memory and SQLite database initialization, providing a
// centralized way to access the career service across all CLI commands.
type CLIContext struct {
	dbPath   string
	inMemory bool
	svc      *careerservice.Service
	db       *sql.DB
}

// NewCLIContext creates a new CLIContext with the specified database configuration.
//
// The service is not initialized in the constructor; call InitService() to
// initialize it after creating the context.
//
// Returns: A new CLIContext with lazy-initialized service.
//
// Side effects: None.
func NewCLIContext(dbPath string, inMemory bool) *CLIContext {
	return &CLIContext{
		dbPath:   dbPath,
		inMemory: inMemory,
		svc:      nil,
	}
}

// InitService initializes the career service with either in-memory or SQLite repositories.
//
// For in-memory mode, creates memory-backed repositories with bidirectional wiring.
// For SQLite mode, creates or opens the database at the specified path (or default
// ~/.kariya/events.db), runs migrations, and initializes GORM repositories.
//
// Returns: nil on success, or an error if initialization fails.
//
// Side effects: Creates database file and directories if using SQLite mode;
// writes error messages to errOut on failure.
func (ctx *CLIContext) InitService(errOut io.Writer) error {
	if ctx.inMemory {
		return ctx.initInMemoryService()
	}
	return ctx.initSQLiteService(errOut)
}

// Service returns the initialized career service.
//
// Returns nil if InitService has not been called yet.
//
// Returns: The initialized *careerservice.Service, or nil if not yet initialized.
//
// Side effects: None.
func (ctx *CLIContext) Service() *careerservice.Service {
	return ctx.svc
}

// Close closes the database connection if using SQLite.
// For in-memory databases, this is a no-op.
func (ctx *CLIContext) Close() error {
	if ctx.db != nil {
		if _, err := ctx.db.ExecContext(context.Background(), "PRAGMA wal_checkpoint(TRUNCATE)"); err != nil {
			return fmt.Errorf("WAL checkpoint failed: %w", err)
		}
		return ctx.db.Close()
	}
	return nil
}

func (ctx *CLIContext) initInMemoryService() error {
	eventRepo := careermemory.NewEventRepository()
	skillRepo := careermemory.NewSkillRepository()
	skillRepo.SetEventRepository(eventRepo)
	eventRepo.SetSkillRepository(skillRepo)

	ctx.svc = careerservice.NewService(eventRepo)
	ctx.svc.SetFactRepository(careermemory.NewFactRepository())
	ctx.svc.SetBurstRepository(careermemory.NewBurstRepository())
	ctx.svc.SetSkillRepository(skillRepo)

	return nil
}

func (ctx *CLIContext) initSQLiteService(errOut io.Writer) error {
	dbPath := ctx.dbPath
	if dbPath == "" {
		var err error
		dbPath, err = ctx.getDefaultDBPath(errOut)
		if err != nil {
			return err
		}
	}

	db, err := careersql.OpenDB(dbPath)
	if err != nil {
		fmt.Fprintf(errOut, "Error opening database at '%s': %v\n", dbPath, err)
		return err
	}
	ctx.db = db

	if err := career.RunMigrations(db); err != nil {
		fmt.Fprintf(errOut, "Error running migrations: %v\n", err)
		return err
	}

	repos, err := careersql.NewRepositories(db)
	if err != nil {
		fmt.Fprintf(errOut, "Error initializing GORM repositories: %v\n", err)
		return err
	}

	ctx.svc = careerservice.NewService(repos.Event)
	ctx.svc.SetFactRepository(repos.Fact)
	ctx.svc.SetBurstRepository(repos.Burst)
	ctx.svc.SetSkillRepository(repos.Skill)

	return nil
}

func (ctx *CLIContext) getDefaultDBPath(errOut io.Writer) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(errOut, "Error getting home directory: %v\n", err)
		return "", err
	}
	kariyaDir := filepath.Join(homeDir, ".kariya")
	dbPath := filepath.Join(kariyaDir, "events.db")

	if err := os.MkdirAll(kariyaDir, 0o750); err != nil {
		fmt.Fprintf(errOut, "Error creating kariya directory: %v\n", err)
		return "", err
	}

	return dbPath, nil
}
