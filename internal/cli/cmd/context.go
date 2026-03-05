package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careersql "github.com/baphled/kariya/internal/repository/career/sql"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// CLIContext holds shared state for CLI commands.
type CLIContext struct {
	DBPath   string
	InMemory bool
	Service  *careerservice.Service
}

// InitService initializes the career service with appropriate repositories.
func (ctx *CLIContext) InitService() error {
	if ctx.InMemory {
		ctx.initInMemory()
		return nil
	}
	return ctx.initSQL()
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

func (ctx *CLIContext) initSQL() error {
	dbPath, err := ctx.getDBPath()
	if err != nil {
		return err
	}

	db, err := careersql.OpenDB(dbPath)
	if err != nil {
		return fmt.Errorf("opening database at '%s': %w", dbPath, err)
	}

	if err := career.RunMigrations(db); err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}

	repos, err := careersql.NewRepositories(db)
	if err != nil {
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
