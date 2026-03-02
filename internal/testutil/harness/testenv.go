package harness

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/cli/bootstrap"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/config"

	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careersql "github.com/baphled/kariya/internal/repository/career/sql"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
)

// Terminal dimensions for test environments.
// Large dimensions accommodate forms with many fields (e.g. Profile with 11+ fields).
// Compact dimensions are used for shared environments where full form height is not needed.
const (
	TerminalWidth        = 120
	TerminalHeightLarge  = 100
	TerminalHeightShared = 40
)

// sharedEnv holds the shared test environment for BeforeSuite/AfterSuite pattern.
// This avoids recreating the database for every test.
var sharedEnv *TestEnv

// sharedTmpDir holds the temp directory for the shared environment.
var sharedTmpDir string

// TestingT is an interface that matches both *testing.T and GinkgoT()
// This allows the harness package to work with both standard Go tests and Ginkgo.
//
//nolint:interfacebloat // Matches standard testing.T interface which has many methods
type TestingT interface {
	Helper()
	TempDir() string
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Error(args ...interface{})
}

// TestEnv holds all test dependencies for E2E testing.
// It provides a complete test environment with SQLite persistence,
// repositories, services, and the application model.
type TestEnv struct {
	// T is the testing context (supports both *testing.T and GinkgoT())
	T TestingT

	// Model is the root application model
	Model *app.Model

	// DB is the SQLite database connection
	DB *sql.DB

	// DBPath is the path to the SQLite database file
	DBPath string

	// Repositories (interface types for flexibility)
	EventRepo careerrepo.EventRepository
	BurstRepo careerrepo.BurstRepository
	FactRepo  careerrepo.FactRepository
	SkillRepo careerrepo.SkillRepository

	// Memory repositories (for fast tests)
	MemEventRepo *careermemory.EventRepository
	MemBurstRepo *careermemory.BurstRepository
	MemFactRepo  *careermemory.FactRepository

	// Services
	Service    *careerservice.Service
	CLIService *service.CLIEventService

	// Context for async operations
	Ctx context.Context

	// Cleanup function to call when done
	cleanup func()

	// QuitRequested is true if tea.Quit was returned by the model
	QuitRequested bool
}

// Setup creates a complete E2E test environment with SQLite persistence.
//
// Expected:
//   - testingt must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func Setup(t TestingT) *TestEnv {
	t.Helper()

	ctx := context.Background()
	// Use os.MkdirTemp instead of t.TempDir() to control cleanup timing
	// t.TempDir() registers auto-cleanup that runs after AfterEach, causing
	// Windows file lock errors when the DB file is still being released
	tmpDir, err := os.MkdirTemp("", "e2e_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	dbPath := filepath.Join(tmpDir, "e2e_test.db")

	// Issue-007 fix: Isolate config file writes to temp directory
	// Use SwapConfigPathForTesting to preserve BeforeSuite's path for restoration
	configPath := filepath.Join(tmpDir, "config.yaml")
	prevConfigPath := config.SwapConfigPathForTesting(configPath)

	// Open database connection with WAL mode and busy timeout.
	db, err := careersql.OpenDB(dbPath)
	if err != nil {
		config.SetConfigPathForTesting(prevConfigPath)
		t.Fatalf("failed to open test db: %v", err)
	}

	// Run migrations
	if err := careerrepo.RunMigrationsForTests(db); err != nil {
		config.SetConfigPathForTesting(prevConfigPath)
		_ = db.Close()
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Create GORM repositories
	repos, err := careersql.NewRepositories(db)
	if err != nil {
		config.SetConfigPathForTesting(prevConfigPath)
		_ = db.Close()
		t.Fatalf("failed to create GORM repositories: %v", err)
	}

	// Create service
	svc := careerservice.NewService(repos.Event)
	svc.SetBurstRepository(repos.Burst)
	svc.SetFactRepository(repos.Fact)
	svc.SetSkillRepository(repos.Skill)

	// Create CLI service
	cliService := service.NewCLIEventService(svc)

	// Create bootstrap result (skipping onboarding for tests)
	log := logger.DefaultLogger()
	bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

	// Create application model
	model := app.NewModel(cliService, svc, bootstrapResult)

	// Set terminal dimensions to ensure forms render correctly.
	// Without this, viewport calculations may use default 80x24, truncating forms.
	model.Update(tea.WindowSizeMsg{Width: TerminalWidth, Height: TerminalHeightLarge})

	cleanup := func() {
		// Issue-007 fix: Restore previous config path (from BeforeSuite) instead of clearing
		// This allows nested isolation without breaking suite-level isolation
		config.SetConfigPathForTesting(prevConfigPath)
		// Close database connection
		if err := db.Close(); err != nil {
			// Log but don't fail - this is cleanup
			t.Errorf("warning: failed to close db: %v", err)
		}
		// Give Windows time to release file handles before temp dir cleanup
		// This prevents "file in use" errors on Windows CI
		// 100ms is needed for reliable cleanup on Windows CI runners
		time.Sleep(100 * time.Millisecond)
		// Manually remove temp dir since we used os.MkdirTemp() instead of t.TempDir()
		// This gives us control over cleanup timing (after DB close + sleep)
		_ = os.RemoveAll(tmpDir)
	}

	return &TestEnv{
		T:          t,
		Model:      model,
		DB:         db,
		DBPath:     dbPath,
		EventRepo:  repos.Event,
		BurstRepo:  repos.Burst,
		FactRepo:   repos.Fact,
		SkillRepo:  repos.Skill,
		Service:    svc,
		CLIService: cliService,
		Ctx:        ctx,
		cleanup:    cleanup,
	}
}

// SetupShared creates a shared E2E test environment for use with BeforeSuite.
//
// Side effects:
//   - None.
func SetupShared() {
	var err error
	sharedTmpDir, err = os.MkdirTemp("", "e2e_test_*")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}

	// Issue-007 fix: Isolate config file writes to temp directory
	configPath := filepath.Join(sharedTmpDir, "config.yaml")
	config.SetConfigPathForTesting(configPath)

	dbPath := filepath.Join(sharedTmpDir, "e2e_shared.db")

	// Open database connection with WAL mode and busy timeout.
	db, err := careersql.OpenDB(dbPath)
	if err != nil {
		panic("failed to open shared test db: " + err.Error())
	}

	if err := careerrepo.RunMigrationsForTests(db); err != nil {
		_ = db.Close()
		panic("failed to run migrations: " + err.Error())
	}

	// Create ORM repositories
	repos, err := careersql.NewRepositories(db)
	if err != nil {
		_ = db.Close()
		panic("failed to create repositories: " + err.Error())
	}

	// Create service
	svc := careerservice.NewService(repos.Event)
	svc.SetBurstRepository(repos.Burst)
	svc.SetFactRepository(repos.Fact)
	svc.SetSkillRepository(repos.Skill)

	// Create CLI service
	cliService := service.NewCLIEventService(svc)

	// Create bootstrap result (skipping onboarding for tests)
	log := logger.DefaultLogger()
	bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

	// Create application model
	model := app.NewModel(cliService, svc, bootstrapResult)

	// Set terminal dimensions to ensure modals render correctly.
	// Without this, viewport calculations may use 0 height, showing only last lines.
	model.Update(tea.WindowSizeMsg{Width: TerminalWidth, Height: TerminalHeightShared})

	sharedEnv = &TestEnv{
		T:          nil,
		Model:      model,
		DB:         db,
		DBPath:     dbPath,
		EventRepo:  repos.Event,
		BurstRepo:  repos.Burst,
		FactRepo:   repos.Fact,
		SkillRepo:  repos.Skill,
		Service:    svc,
		CLIService: cliService,
		Ctx:        context.Background(),
		cleanup:    nil,
	}
}

// CleanupShared releases all shared test resources.
//
// Side effects:
//   - None.
func CleanupShared() {
	if sharedEnv != nil && sharedEnv.DB != nil {
		_ = sharedEnv.DB.Close()
	}
	// Issue-007 fix: Reset config path override
	config.ResetConfigPath()
	if sharedTmpDir != "" {
		_ = os.RemoveAll(sharedTmpDir)
	}
	sharedEnv = nil
	sharedTmpDir = ""
}

// GetSharedEnv returns the shared test environment for use in BeforeEach.
//
// Expected:
//   - testingt must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func GetSharedEnv(t TestingT) *TestEnv {
	if sharedEnv == nil {
		panic("shared env not initialized - call SetupShared() in BeforeSuite")
	}

	// Reset database state (truncate all tables)
	sharedEnv.resetDatabase()

	// Create fresh bootstrap result (skipping onboarding for tests)
	log := logger.DefaultLogger()
	bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), sharedEnv.Service, log)

	// Create fresh application model (only thing with UI state)
	// Repositories and services are stateless, so we reuse them
	model := app.NewModel(sharedEnv.CLIService, sharedEnv.Service, bootstrapResult)

	// Set terminal dimensions so modals and overlays render correctly.
	// Without this, viewport calculations use 0x0, causing empty views.
	model.Update(tea.WindowSizeMsg{Width: TerminalWidth, Height: TerminalHeightShared})

	// Update only what changes per-test
	sharedEnv.T = t
	sharedEnv.Model = model
	sharedEnv.Ctx = context.Background()

	return sharedEnv
}

// GetSharedEnvWithOnboarding returns the shared test environment for onboarding tests.
//
// Expected:
//   - testingt must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func GetSharedEnvWithOnboarding(t TestingT) *TestEnv {
	if sharedEnv == nil {
		panic("shared env not initialized - call SetupShared() in BeforeSuite")
	}

	// Reset database state (truncate all tables)
	sharedEnv.resetDatabase()

	// Create fresh bootstrap result
	log := logger.DefaultLogger()
	bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), sharedEnv.Service, log)

	// Create fresh application model
	model := app.NewModel(sharedEnv.CLIService, sharedEnv.Service, bootstrapResult)

	// Update only what changes per-test
	sharedEnv.T = t
	sharedEnv.Model = model
	sharedEnv.Ctx = context.Background()

	return sharedEnv
}

// GetOnboardingTestModel returns a model for testing the onboarding wizard UI.
//
// Expected:
//   - config must be a valid configuration object.
//
// Returns:
//   - A fully initialized bootstrap.OnboardingTestModel ready for use.
//
// Side effects:
//   - None.
func GetOnboardingTestModel(existingProfile *config.ProfileConfig) *bootstrap.OnboardingTestModel {
	return bootstrap.NewOnboardingTestModel(existingProfile)
}

// resetDatabase truncates all tables to reset state between tests.
func (e *TestEnv) resetDatabase() {
	if e.DB == nil {
		return
	}

	// Truncate tables in order (respecting foreign key constraints)
	// Using explicit statements to avoid SQL string concatenation warnings
	// Errors during cleanup are logged but not fatal for test teardown
	e.DB.Exec("DELETE FROM event_skills")  // #nosec G104 -- Test cleanup errors are not critical
	e.DB.Exec("DELETE FROM facts")         // #nosec G104 -- Test cleanup errors are not critical
	e.DB.Exec("DELETE FROM bursts")        // #nosec G104 -- Test cleanup errors are not critical
	e.DB.Exec("DELETE FROM skills")        // #nosec G104 -- Test cleanup errors are not critical
	e.DB.Exec("DELETE FROM career_events") // #nosec G104 -- Test cleanup errors are not critical
}

// SetupWithOnboarding creates an E2E test environment with the onboarding wizard active.
//
// Expected:
//   - testingt must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func SetupWithOnboarding(t TestingT) *TestEnv {
	t.Helper()

	ctx := context.Background()
	// Use os.MkdirTemp instead of t.TempDir() to control cleanup timing
	// t.TempDir() registers auto-cleanup that runs after AfterEach, causing
	// Windows file lock errors when the DB file is still being released
	tmpDir, err := os.MkdirTemp("", "e2e_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	dbPath := filepath.Join(tmpDir, "e2e_test.db")

	// Issue-007 fix: Isolate config file writes to temp directory
	// Use SwapConfigPathForTesting to preserve BeforeSuite's path for restoration
	configPath := filepath.Join(tmpDir, "config.yaml")
	prevConfigPath := config.SwapConfigPathForTesting(configPath)

	// Open database connection with WAL mode and busy timeout.
	db, err := careersql.OpenDB(dbPath)
	if err != nil {
		config.SetConfigPathForTesting(prevConfigPath)
		_ = os.RemoveAll(tmpDir)
		t.Fatalf("failed to open test db: %v", err)
	}

	// Run migrations
	if err := careerrepo.RunMigrationsForTests(db); err != nil {
		config.SetConfigPathForTesting(prevConfigPath)
		_ = db.Close()
		_ = os.RemoveAll(tmpDir)
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Create ORM repositories
	repos, err := careersql.NewRepositories(db)
	if err != nil {
		config.SetConfigPathForTesting(prevConfigPath)
		_ = db.Close()
		_ = os.RemoveAll(tmpDir)
		t.Fatalf("failed to create repositories: %v", err)
	}

	// Create service
	svc := careerservice.NewService(repos.Event)
	svc.SetBurstRepository(repos.Burst)
	svc.SetFactRepository(repos.Fact)
	svc.SetSkillRepository(repos.Skill)

	// Create CLI service
	cliService := service.NewCLIEventService(svc)

	// Create bootstrap result (skipping onboarding for main app)
	log := logger.DefaultLogger()
	bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

	// Create application model
	model := app.NewModel(cliService, svc, bootstrapResult)

	cleanup := func() {
		// Issue-007 fix: Restore previous config path (from BeforeSuite) instead of clearing
		config.SetConfigPathForTesting(prevConfigPath)
		// Close database connection
		if err := db.Close(); err != nil {
			// Log but don't fail - this is cleanup
			t.Errorf("warning: failed to close db: %v", err)
		}
		// Give Windows time to release file handles before temp dir cleanup
		// This prevents "file in use" errors on Windows CI
		// 100ms is needed for reliable cleanup on Windows CI runners
		time.Sleep(100 * time.Millisecond)
		// Manually remove temp dir since we used os.MkdirTemp() instead of t.TempDir()
		// This gives us control over cleanup timing (after DB close + sleep)
		_ = os.RemoveAll(tmpDir)
	}

	return &TestEnv{
		T:          t,
		Model:      model,
		DB:         db,
		DBPath:     dbPath,
		EventRepo:  repos.Event,
		BurstRepo:  repos.Burst,
		FactRepo:   repos.Fact,
		SkillRepo:  repos.Skill,
		Service:    svc,
		CLIService: cliService,
		Ctx:        ctx,
		cleanup:    cleanup,
	}
}

// SetupWithMemory creates an E2E test environment using in-memory repositories.
//
// Expected:
//   - testingt must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func SetupWithMemory(t TestingT) *TestEnv {
	t.Helper()

	ctx := context.Background()
	// Use os.MkdirTemp instead of t.TempDir() to control cleanup timing
	// This maintains consistency with Setup() and SetupWithOnboarding()
	tmpDir, err := os.MkdirTemp("", "e2e_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Issue-007 fix: Isolate config file writes to temp directory
	// Use SwapConfigPathForTesting to preserve BeforeSuite's path for restoration
	configPath := filepath.Join(tmpDir, "config.yaml")
	prevConfigPath := config.SwapConfigPathForTesting(configPath)

	// Create in-memory repositories
	eventRepo := careermemory.NewEventRepository()
	burstRepo := careermemory.NewBurstRepository()
	factRepo := careermemory.NewFactRepository()
	skillRepo := careermemory.NewSkillRepository()

	// Cross-link event and skill repos for LinkSkill/GetSkillsForEvent sync.
	skillRepo.SetEventRepository(eventRepo)
	eventRepo.SetSkillRepository(skillRepo)

	// Create service
	svc := careerservice.NewService(eventRepo)
	svc.SetBurstRepository(burstRepo)
	svc.SetFactRepository(factRepo)
	svc.SetSkillRepository(skillRepo)

	// Create CLI service
	cliService := service.NewCLIEventService(svc)

	// Create bootstrap result (skipping onboarding for tests)
	log := logger.DefaultLogger()
	bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

	// Create application model
	model := app.NewModel(cliService, svc, bootstrapResult)

	return &TestEnv{
		T:            t,
		Model:        model,
		DB:           nil,
		DBPath:       "",
		MemEventRepo: eventRepo,
		MemBurstRepo: burstRepo,
		MemFactRepo:  factRepo,
		Service:      svc,
		CLIService:   cliService,
		Ctx:          ctx,
		cleanup: func() {
			// Issue-007 fix: Restore previous config path (from BeforeSuite) instead of clearing
			config.SetConfigPathForTesting(prevConfigPath)
			// Manually remove temp dir since we used os.MkdirTemp() instead of t.TempDir()
			_ = os.RemoveAll(tmpDir)
		},
	}
}

// Cleanup releases all test resources.
//
// Side effects:
//   - None.
func (e *TestEnv) Cleanup() {
	// Skip cleanup for shared environment (cleanup is nil)
	if e.cleanup != nil {
		e.cleanup()
	}
}
