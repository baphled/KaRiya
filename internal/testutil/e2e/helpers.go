// Package e2e provides E2E test utilities for integration testing of the KaRiya TUI.
package e2e

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/cli/bootstrap"
	"github.com/baphled/kariya/internal/cli/intents"
	burstmanagement "github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/cli/intents/captureevent"
	"github.com/baphled/kariya/internal/cli/intents/factmanagement"
	"github.com/baphled/kariya/internal/cli/intents/generatecv"
	"github.com/baphled/kariya/internal/cli/intents/skillsmanagement"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careersql "github.com/baphled/kariya/internal/repository/career/sql"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	_ "modernc.org/sqlite"
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
// This allows the e2e package to work with both standard Go tests and Ginkgo.
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

// ============================================================================
// Navigation Helpers
// ============================================================================

// SelectIntent navigates to and selects a menu item by its index (0-based).
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SelectIntent(index int) *TestEnv {
	e.T.Helper()

	// Ensure we're in menu state before attempting navigation
	if !e.IsInMenuState() {
		e.T.Fatalf("Cannot select intent: not in menu state. Current view:\n%s", e.GetView())
	}

	// Navigate to the menu item from position 0
	// Press 'g' to ensure we're at the top of the menu first
	e.PressKeyRune('g')

	// Navigate down to the desired index
	for range index {
		e.PressKeyRune('j')
	}

	// Select the intent
	e.PressKey(tea.KeyEnter)

	return e
}

// SelectIntentByName navigates to and selects a menu item by its intent name.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SelectIntentByName(name string) *TestEnv {
	e.T.Helper()

	intentOrder := map[string]int{
		"capture_event":    0,
		"browse_timeline":  1,
		"manage_skills":    2,
		"generate_cv":      3,
		"configure_system": 4,
		"burst_management": 5,
		"fact_management":  6,
	}

	index, ok := intentOrder[name]
	if !ok {
		e.T.Fatalf("unknown intent name: %s", name)
	}

	return e.SelectIntent(index)
}

// PressKey sends a key message to the model.
//
// Expected:
//   - keytype must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) PressKey(key tea.KeyType) *TestEnv {
	e.T.Helper()

	modelInterface, cmd := e.Model.Update(tea.KeyMsg{Type: key})
	model, ok := modelInterface.(*app.Model)
	if !ok {
		e.T.Fatal("model type assertion failed: expected *app.Model")
	}
	e.Model = model

	// Execute any returned command
	e.executeCmd(cmd)

	return e
}

// PressKeyRune sends a rune key message to the model.
//
// Expected:
//   - rune must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) PressKeyRune(r rune) *TestEnv {
	e.T.Helper()

	modelInterface, cmd := e.Model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	model, ok := modelInterface.(*app.Model)
	if !ok {
		e.T.Fatal("model type assertion failed: expected *app.Model")
	}
	e.Model = model

	// Execute any returned command
	e.executeCmd(cmd)

	return e
}

// PressKeys sends multiple keys in sequence.
//
// Expected:
//   - interface{} must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) PressKeys(keys ...interface{}) *TestEnv {
	e.T.Helper()

	for _, key := range keys {
		switch k := key.(type) {
		case tea.KeyType:
			e.PressKey(k)
		case rune:
			e.PressKeyRune(k)
		case string:
			for _, r := range k {
				e.PressKeyRune(r)
			}
		default:
			e.T.Fatalf("unsupported key type: %T", key)
		}
	}

	return e
}

// TypeText types a string character by character.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) TypeText(text string) *TestEnv {
	e.T.Helper()

	for _, r := range text {
		e.PressKeyRune(r)
	}

	return e
}

// NavigateDown moves down in a list (j or down arrow).
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) NavigateDown() *TestEnv {
	return e.PressKeyRune('j')
}

// NavigateUp moves up in a list (k or up arrow).
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) NavigateUp() *TestEnv {
	return e.PressKeyRune('k')
}

// Confirm presses Enter to confirm an action.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) Confirm() *TestEnv {
	return e.PressKey(tea.KeyEnter)
}

// Cancel presses Escape to cancel/go back.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) Cancel() *TestEnv {
	return e.PressKey(tea.KeyEscape)
}

// GoBack presses Escape to go back.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) GoBack() *TestEnv {
	return e.Cancel()
}

// Quit presses 'q' to quit.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) Quit() *TestEnv {
	return e.PressKeyRune('q')
}

// Tab presses Tab to move to next field.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) Tab() *TestEnv {
	return e.PressKey(tea.KeyTab)
}

// SubmitHuhForm submits a huh form by pressing Enter.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SubmitHuhForm() *TestEnv {
	return e.Confirm()
}

// ClearTextField clears a text field by moving to end and pressing backspace.
// This is useful for clearing pre-populated huh form fields.
//
// Expected:
//   - maxChars should be a reasonable upper bound for the text length.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) ClearTextField(maxChars int) *TestEnv {
	e.T.Helper()

	e.PressKey(tea.KeyCtrlE)

	for range maxChars {
		e.PressKey(tea.KeyBackspace)
	}

	return e
}

// NextFormField moves to the next field in a huh form.
// This sends the huh.NextField message directly to properly navigate forms.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) NextFormField() *TestEnv {
	e.T.Helper()

	msg := huh.NextField()
	e.updateModelAndExecute(msg)

	return e
}

// executeCmd executes commands returned by Update, but only for specific message types
// that are essential for state transitions (like form submission).
//
// Most Bubble Tea commands (cursor blink, window resize) are ignored because they
// cause infinite loops or stuck goroutines in tests. We only care about messages
// that actually change application state.
//
// Commands that take longer than 500ms to execute (tick commands with delays) are skipped
// to avoid slow tests from cursor blink animations (530ms each). The 500ms timeout
// allows database operations to complete while still filtering out cursor blinks.
func (e *TestEnv) executeCmd(cmd tea.Cmd) {
	if cmd == nil {
		return
	}

	// Execute command with timeout to skip slow tick commands
	// Cursor blink ticks take 530ms, database operations typically complete in <100ms
	type result struct {
		msg tea.Msg
	}
	done := make(chan result, 1)
	go func() {
		done <- result{msg: cmd()}
	}()

	select {
	case r := <-done:
		if r.msg == nil {
			return
		}
		e.processCmdResult(r.msg)
	case <-time.After(500 * time.Millisecond):
		// Command is a slow tick (cursor blink, etc.) - skip it
		return
	}
}

// updateModelAndExecute updates the model with a message and executes any returned command.
// This is a helper to reduce cognitive complexity in processCmdResult.
//
// Expected:
//   - msg must be a valid tea.Msg.
//
// Returns:
//   - None.
//
// Side effects:
//   - Updates the Model field of the TestEnv.
//   - Recursively executes any returned command.
func (e *TestEnv) updateModelAndExecute(msg tea.Msg) {
	modelInterface, nextCmd := e.Model.Update(msg)
	model, ok := modelInterface.(*app.Model)
	if !ok {
		e.T.Fatal("model type assertion failed: expected *app.Model")
	}
	e.Model = model
	e.executeCmd(nextCmd)
}

// processCmdResult processes a message returned from a command.
// Only essential state transition messages are processed.
//
// Expected:
//   - msg must be a valid tea.Msg.
//
// Returns:
//   - None.
//
// Side effects:
//   - May update the Model field of the TestEnv.
//   - May recursively execute commands.
//
//nolint:gocyclo // Type switch handler inherently requires many cases for different message types.
func (e *TestEnv) processCmdResult(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.QuitMsg:
		e.QuitRequested = true
	case tea.BatchMsg:
		e.processBatchMsg(msg)
	case captureevent.SubmitMsg:
		e.updateModelAndExecute(msg)
	case captureevent.SubmitCompleteMsg:
		e.processSubmitCompleteMsg(msg)
	case captureevent.InferenceCompleteMsg:
		e.updateModelAndExecute(msg)
	case captureevent.PostSavePersistenceCompleteMsg:
		e.updateModelAndExecute(msg)
	case captureevent.SubmitErrorMsg:
		e.updateModelAndExecute(msg)
	case captureevent.DismissModalMsg:
		e.updateModelAndExecute(msg)
	case intents.ConfigCompleteMsg:
		e.updateModelAndExecute(msg)
	case generatecv.TechnologiesExtractedMsg:
		e.updateModelAndExecute(msg)
	case generatecv.CVGenerationCompleteMsg:
		e.updateModelAndExecute(msg)
	case generatecv.ExportCompleteMsg:
		e.updateModelAndExecute(msg)
	case skillsmanagement.SkillsLoadedMsg:
		e.updateModelAndExecute(msg)
	case skillsmanagement.SkillCreatedMsg:
		e.updateModelAndExecute(msg)
	case skillsmanagement.SkillUpdatedMsg:
		e.updateModelAndExecute(msg)
	case skillsmanagement.SkillDeletedMsg:
		e.updateModelAndExecute(msg)
	case skillsmanagement.SkillFormCompleteMsg:
		e.updateModelAndExecute(msg)
	case skillsmanagement.SkillEventsForModalLoadedMsg:
		e.updateModelAndExecute(msg)
	case burstmanagement.EditBurstMsg:
		e.updateModelAndExecute(msg)
	case burstmanagement.BurstDeletedMsg:
		e.updateModelAndExecute(msg)
	case burstmanagement.BurstConfirmedMsg:
		e.updateModelAndExecute(msg)
	case burstmanagement.BurstEventsLoadedMsg:
		e.updateModelAndExecute(msg)
	case burstmanagement.BurstFactsLoadedMsg:
		e.updateModelAndExecute(msg)
	case burstmanagement.BurstSkillsLoadedMsg:
		e.updateModelAndExecute(msg)
	case factmanagement.FactsLoadedMsg:
		e.updateModelAndExecute(msg)
	case factmanagement.FactSavedMsg:
		e.updateModelAndExecute(msg)
	case factmanagement.FactDeletedMsg:
		e.updateModelAndExecute(msg)
	case feedback.ModalCountdownTickMsg:
		e.updateModelAndExecute(msg)
	case feedback.ModalAutoDismissMsg:
		e.updateModelAndExecute(msg)
	}
}

// processBatchMsg processes a batch of commands.
//
// Expected:
//   - msg must be a valid tea.BatchMsg.
//
// Returns:
//   - None.
//
// Side effects:
//   - Executes each non-nil command in the batch.
func (e *TestEnv) processBatchMsg(msg tea.BatchMsg) {
	for _, cmd := range msg {
		if cmd != nil {
			e.executeCmd(cmd)
		}
	}
}

// processSubmitCompleteMsg handles submit completion and dismisses the modal.
//
// Expected:
//   - msg must be a valid captureevent.SubmitCompleteMsg.
//
// Returns:
//   - None.
//
// Side effects:
//   - Updates the Model with the submit completion message.
//   - Automatically dismisses the success modal to skip the timer.
func (e *TestEnv) processSubmitCompleteMsg(msg captureevent.SubmitCompleteMsg) {
	e.updateModelAndExecute(msg)
	e.updateModelAndExecute(captureevent.DismissModalMsg{})
}

// SendMessage sends a message directly to the model.
//
// Expected:
//   - msg must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SendMessage(msg tea.Msg) *TestEnv {
	e.T.Helper()

	modelInterface, cmd := e.Model.Update(msg)
	model, ok := modelInterface.(*app.Model)
	if !ok {
		e.T.Fatal("model type assertion failed: expected *app.Model")
	}
	e.Model = model
	e.executeCmd(cmd)

	return e
}

// SubmitEvent sends a SubmitMsg directly to the model with the given event.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SubmitEvent(event *career.Event) *TestEnv {
	e.T.Helper()

	return e.SendMessage(captureevent.SubmitMsg{Event: event, Err: nil})
}

// SubmitEventWithError sends a SubmitMsg with an error to trigger validation error handling.
//
// Expected:
//   - err should describe the validation failure.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SubmitEventWithError(event *career.Event, err error) *TestEnv {
	e.T.Helper()

	return e.SendMessage(captureevent.SubmitMsg{Event: event, Err: err})
}

// SubmitSkill creates a skill in the repository and sends a SkillCreatedMsg.
// This bypasses the UI form submission path, similar to SubmitEvent.
//
// Expected:
//   - skill must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - Creates the skill in the repository.
func (e *TestEnv) SubmitSkill(skill *career.Skill) *TestEnv {
	e.T.Helper()

	skillRepo := e.Service.GetSkillRepository()
	if skillRepo == nil {
		e.T.Fatal("skill repository not set")
	}

	err := skillRepo.Create(e.Ctx, skill)

	return e.SendMessage(skillsmanagement.SkillCreatedMsg{Skill: skill, Error: err})
}

// SubmitSkillWithError sends a SkillCreatedMsg with an error to trigger validation error handling.
//
// Expected:
//   - err should describe the validation failure.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SubmitSkillWithError(skill *career.Skill, err error) *TestEnv {
	e.T.Helper()

	return e.SendMessage(skillsmanagement.SkillCreatedMsg{Skill: skill, Error: err})
}

// SubmitSkillUpdate updates a skill in the repository and sends a SkillUpdatedMsg.
// This bypasses the UI form submission path for skill editing.
//
// Expected:
//   - skill must be valid with existing ID.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - Updates the skill in the repository.
func (e *TestEnv) SubmitSkillUpdate(skill *career.Skill) *TestEnv {
	e.T.Helper()

	skillRepo := e.Service.GetSkillRepository()
	if skillRepo == nil {
		e.T.Fatal("skill repository not set")
	}

	err := skillRepo.Update(e.Ctx, skill)

	return e.SendMessage(skillsmanagement.SkillUpdatedMsg{Skill: skill, Error: err})
}

// SubmitSkillUpdateWithError sends a SkillUpdatedMsg with an error.
//
// Expected:
//   - err should describe the update failure.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SubmitSkillUpdateWithError(skill *career.Skill, err error) *TestEnv {
	e.T.Helper()

	return e.SendMessage(skillsmanagement.SkillUpdatedMsg{Skill: skill, Error: err})
}

// SubmitFact creates a fact in the repository and sends a FactSavedMsg.
// This bypasses the UI form submission path for fact creation.
//
// Expected:
//   - fact must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - Creates the fact in the repository.
//
//nolint:dupl // Acceptable duplication in test helpers - similar to SubmitSkill pattern
func (e *TestEnv) SubmitFact(fact *career.Fact) *TestEnv {
	e.T.Helper()

	factRepo := e.Service.GetFactRepository()
	if factRepo == nil {
		e.T.Fatal("fact repository not set")
	}

	err := factRepo.Create(e.Ctx, fact)
	if err != nil {
		return e.SendMessage(factmanagement.FactSavedMsg{Fact: fact, IsNew: true, Message: err.Error()})
	}

	return e.SendMessage(factmanagement.FactSavedMsg{Fact: fact, IsNew: true, Message: "Fact saved successfully"})
}

// SubmitFactUpdate updates a fact in the repository and sends a FactSavedMsg.
// This bypasses the UI form submission path for fact editing.
//
// Expected:
//   - fact must be valid with existing ID.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - Updates the fact in the repository.
//
//nolint:dupl // Acceptable duplication in test helpers - similar to SubmitSkill pattern
func (e *TestEnv) SubmitFactUpdate(fact *career.Fact) *TestEnv {
	e.T.Helper()

	factRepo := e.Service.GetFactRepository()
	if factRepo == nil {
		e.T.Fatal("fact repository not set")
	}

	err := factRepo.Update(e.Ctx, fact)
	if err != nil {
		return e.SendMessage(factmanagement.FactSavedMsg{Fact: fact, IsNew: false, Message: err.Error()})
	}

	return e.SendMessage(factmanagement.FactSavedMsg{Fact: fact, IsNew: false, Message: "Fact updated successfully"})
}

// SubmitBurstUpdate updates a burst in the repository and sends a BurstEditCompleteMsg.
// This bypasses the UI form submission path for burst editing.
//
// Expected:
//   - burst must be valid with existing ID.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - Updates the burst in the repository.
func (e *TestEnv) SubmitBurstUpdate(burst *career.Burst) *TestEnv {
	e.T.Helper()

	burstRepo := e.Service.GetBurstRepository()
	if burstRepo == nil {
		e.T.Fatal("burst repository not set")
	}

	err := burstRepo.Update(e.Ctx, burst)

	return e.SendMessage(burstmanagement.BurstEditCompleteMsg{Burst: burst, Cancelled: false, Error: err})
}

// ConfirmBurst confirms a burst and shows the loading modal for fact extraction.
// This bypasses the UI confirmation modal but triggers the loading state that tests expect.
//
// Expected:
//   - burst must be valid with existing ID.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - Updates the burst as confirmed in the repository.
//   - Shows loading modal (fact extraction runs but modal stays visible for test assertions).
func (e *TestEnv) ConfirmBurst(burst *career.Burst) *TestEnv {
	e.T.Helper()

	burstRepo := e.Service.GetBurstRepository()
	if burstRepo == nil {
		e.T.Fatal("burst repository not set")
	}

	// Bypass service ConfirmBurst (which triggers async fact extraction)
	// Just update the burst directly in the repository
	burst.Confirmed = true
	now := time.Now()
	burst.ConfirmedAt = &now
	burst.UpdatedAt = now

	err := burstRepo.Update(e.Ctx, burst)
	if err != nil {
		e.T.Fatalf("failed to update burst: %v", err)
	}

	// Don't send any message - tests will check repository state
	return e
}

// DismissSuccessModal bypasses the auto-dismiss countdown and immediately
// dismisses the success modal. Use this to speed up tests that don't need
// to verify countdown behavior. To test the actual countdown, send
// ModalCountdownTickMsg messages explicitly instead.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - Sends a captureevent.DismissModalMsg through the test environment.
//   - Advances the Bubble Tea update loop and updates e.Model to the
//     post-dismissal state of the success modal.
func (e *TestEnv) DismissSuccessModal() *TestEnv {
	e.T.Helper()

	return e.SendMessage(captureevent.DismissModalMsg{})
}

// ============================================================================
// View Assertion Helpers
// ============================================================================

// GetView returns the current view output.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (e *TestEnv) GetView() string {
	return e.Model.View()
}

// AssertViewContains checks that the view contains the given substring.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AssertViewContains(substr string) *TestEnv {
	e.T.Helper()

	view := e.GetView()
	if !strings.Contains(view, substr) {
		e.T.Errorf("expected view to contain %q, but it doesn't.\nView:\n%s", substr, view)
	}

	return e
}

// AssertViewNotContains checks that the view does NOT contain the given substring.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AssertViewNotContains(substr string) *TestEnv {
	e.T.Helper()

	view := e.GetView()
	if strings.Contains(view, substr) {
		e.T.Errorf("expected view NOT to contain %q, but it does.\nView:\n%s", substr, view)
	}

	return e
}

// AssertViewContainsAny checks that the view contains at least one of the given substrings.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AssertViewContainsAny(substrs ...string) *TestEnv {
	e.T.Helper()

	view := e.GetView()
	for _, substr := range substrs {
		if strings.Contains(view, substr) {
			return e
		}
	}

	e.T.Errorf("expected view to contain one of %v, but none found.\nView:\n%s", substrs, view)
	return e
}

// ============================================================================
// Data Verification Helpers
// ============================================================================

// AssertEventCount verifies the number of events in the database.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AssertEventCount(expected int) *TestEnv {
	e.T.Helper()

	events, err := e.Service.ListEvents(e.Ctx, careerrepo.EventListFilters{})
	if err != nil {
		e.T.Fatalf("failed to get events: %v", err)
	}

	if len(events) != expected {
		e.T.Errorf("expected %d events, got %d", expected, len(events))
	}

	return e
}

// AssertBurstCount verifies the number of bursts in the database.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AssertBurstCount(expected int) *TestEnv {
	e.T.Helper()

	burstRepo := e.Service.GetBurstRepository()
	if burstRepo == nil {
		e.T.Fatal("burst repository not set")
	}

	bursts, err := burstRepo.List(e.Ctx, careerrepo.BurstListFilters{})
	if err != nil {
		e.T.Fatalf("failed to get bursts: %v", err)
	}

	if len(bursts) != expected {
		e.T.Errorf("expected %d bursts, got %d", expected, len(bursts))
	}

	return e
}

// AssertFactCount verifies the number of facts in the database.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AssertFactCount(expected int) *TestEnv {
	e.T.Helper()

	factRepo := e.Service.GetFactRepository()
	if factRepo == nil {
		e.T.Fatal("fact repository not set")
	}

	facts, err := factRepo.List(e.Ctx, careerrepo.FactListFilters{})
	if err != nil {
		e.T.Fatalf("failed to get facts: %v", err)
	}

	if len(facts) != expected {
		e.T.Errorf("expected %d facts, got %d", expected, len(facts))
	}

	return e
}

// GetEvents returns all events from the database.
//
// Returns:
//   - A []*career.Event value.
//
// Side effects:
//   - None.
func (e *TestEnv) GetEvents() []*career.Event {
	e.T.Helper()

	events, err := e.Service.ListEvents(e.Ctx, careerrepo.EventListFilters{})
	if err != nil {
		e.T.Fatalf("failed to get events: %v", err)
	}

	return events
}

// GetBursts returns all bursts from the database.
//
// Returns:
//   - A []*career.Burst value.
//
// Side effects:
//   - None.
func (e *TestEnv) GetBursts() []*career.Burst {
	e.T.Helper()

	burstRepo := e.Service.GetBurstRepository()
	if burstRepo == nil {
		e.T.Fatal("burst repository not set")
	}

	bursts, err := burstRepo.List(e.Ctx, careerrepo.BurstListFilters{})
	if err != nil {
		e.T.Fatalf("failed to get bursts: %v", err)
	}

	return bursts
}

// GetFacts returns all facts from the database.
//
// Returns:
//   - A []*career.Fact value.
//
// Side effects:
//   - None.
func (e *TestEnv) GetFacts() []*career.Fact {
	e.T.Helper()

	factRepo := e.Service.GetFactRepository()
	if factRepo == nil {
		e.T.Fatal("fact repository not set")
	}

	facts, err := factRepo.List(e.Ctx, careerrepo.FactListFilters{})
	if err != nil {
		e.T.Fatalf("failed to get facts: %v", err)
	}

	return facts
}

// GetSkills returns all skills from the database.
//
// Returns:
//   - A []*career.Skill value.
//
// Side effects:
//   - None.
func (e *TestEnv) GetSkills() []*career.Skill {
	e.T.Helper()

	skillRepo := e.Service.GetSkillRepository()
	if skillRepo == nil {
		e.T.Fatal("skill repository not set")
	}

	skills, err := skillRepo.List(e.Ctx, nil)
	if err != nil {
		e.T.Fatalf("failed to get skills: %v", err)
	}

	return skills
}

// AssertSkillCount verifies the number of skills in the database.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AssertSkillCount(expected int) *TestEnv {
	e.T.Helper()

	skills := e.GetSkills()

	if len(skills) != expected {
		e.T.Errorf("expected %d skills, got %d", expected, len(skills))
	}

	return e
}

// AddSkill creates a skill in the database.
//
// Expected:
//   - skill must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AddSkill(skill *career.Skill) *TestEnv {
	e.T.Helper()

	skillRepo := e.Service.GetSkillRepository()
	if skillRepo == nil {
		e.T.Fatal("skill repository not set")
	}

	err := skillRepo.Create(e.Ctx, skill)
	if err != nil {
		e.T.Fatalf("failed to create skill: %v", err)
	}

	return e
}

// ============================================================================
// Session Simulation
// ============================================================================

// SimulateRestart recreates the application model while preserving the database.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SimulateRestart() *TestEnv {
	e.T.Helper()

	if e.DB == nil {
		e.T.Fatal("SimulateRestart requires SQLite persistence (use Setup, not SetupWithMemory)")
	}

	// Create new repositories pointing to the same database using ORM
	repos, err := careersql.NewRepositories(e.DB)
	if err != nil {
		e.T.Fatalf("failed to create repositories: %v", err)
	}
	eventRepo := repos.Event
	burstRepo := repos.Burst
	factRepo := repos.Fact
	skillRepo := repos.Skill

	// Create new service
	svc := careerservice.NewService(eventRepo)
	svc.SetBurstRepository(burstRepo)
	svc.SetFactRepository(factRepo)
	svc.SetSkillRepository(skillRepo)

	// Create new CLI service
	cliService := service.NewCLIEventService(svc)

	// Create bootstrap result (skipping onboarding for tests)
	log := logger.DefaultLogger()
	bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

	// Create new application model
	model := app.NewModel(cliService, svc, bootstrapResult)

	// Update environment
	e.EventRepo = eventRepo
	e.BurstRepo = burstRepo
	e.FactRepo = factRepo
	e.SkillRepo = skillRepo
	e.Service = svc
	e.CLIService = cliService
	e.Model = model

	return e
}

// ============================================================================
// Data Population Helpers
// ============================================================================

// AddEvent creates an event in the database.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AddEvent(event *career.Event) *TestEnv {
	e.T.Helper()

	eventRepo := e.Service.GetEventRepository()
	if eventRepo == nil {
		e.T.Fatal("event repository not set")
	}

	err := eventRepo.Create(e.Ctx, event)
	if err != nil {
		e.T.Fatalf("failed to create event: %v", err)
	}

	return e
}

// AddBurst creates a burst in the database.
//
// Expected:
//   - burst must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AddBurst(burst *career.Burst) *TestEnv {
	e.T.Helper()

	burstRepo := e.Service.GetBurstRepository()
	if burstRepo == nil {
		e.T.Fatal("burst repository not set")
	}

	err := burstRepo.Create(e.Ctx, burst)
	if err != nil {
		e.T.Fatalf("failed to create burst: %v", err)
	}

	return e
}

// AddFact creates a fact in the database.
//
// Expected:
//   - fact must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AddFact(fact *career.Fact) *TestEnv {
	e.T.Helper()

	factRepo := e.Service.GetFactRepository()
	if factRepo == nil {
		e.T.Fatal("fact repository not set")
	}

	err := factRepo.Create(e.Ctx, fact)
	if err != nil {
		e.T.Fatalf("failed to create fact: %v", err)
	}

	return e
}

// ============================================================================
// Menu Item Helpers
// ============================================================================

// GetMenuItems returns the list of menu items from the model.
//
// Returns:
//   - A []app.MenuItem value.
//
// Side effects:
//   - None.
func (e *TestEnv) GetMenuItems() []app.MenuItem {
	return e.Model.GetMenuItems()
}

// IsInMenuState checks if the application is currently showing the main menu.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (e *TestEnv) IsInMenuState() bool {
	view := e.GetView()
	// Menu shows the tagline and menu items
	return strings.Contains(view, "Career Event Management System") &&
		strings.Contains(view, "Capture Event")
}

// ============================================================================
// Onboarding Helpers
// ============================================================================

// IsInOnboardingState checks if the application is currently showing the onboarding wizard.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (e *TestEnv) IsInOnboardingState() bool {
	// Onboarding is now a separate program that runs before the main app.
	// The main app is never in an "onboarding state".
	return false
}

// SkipOnboarding is deprecated - onboarding is now skipped by default.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SkipOnboarding() *TestEnv {
	// No-op - onboarding is handled by bootstrap before app creation.
	return e
}

// InitModel initializes the model by calling Init() and sending a WindowSizeMsg.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) InitModel() *TestEnv {
	e.T.Helper()

	_ = e.Model.Init()

	// Send a WindowSizeMsg to trigger form layout
	e.SendMessage(tea.WindowSizeMsg{Width: TerminalWidth, Height: TerminalHeightShared})

	// Type and delete a character to force the form to render its fields
	// This workaround activates huh's internal rendering state
	e.PressKeyRune('x')
	e.PressKey(tea.KeyBackspace)

	return e
}

// PressEnterWithFormProcessing presses Enter and processes any internal form messages.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) PressEnterWithFormProcessing() *TestEnv {
	e.T.Helper()

	// Send Enter key
	modelInterface, cmd := e.Model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model, ok := modelInterface.(*app.Model)
	if !ok {
		e.T.Fatal("model type assertion failed: expected *app.Model")
	}
	e.Model = model

	// Process all resulting messages (including internal form messages)
	// This allows huh's group transitions to complete
	e.processFormCmds(cmd, 10)

	return e
}

// SendMessageWithFormProcessing sends a message and processes all resulting
// internal form messages (including huh init, focus, and group transitions).
//
// Expected:
//   - msg must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SendMessageWithFormProcessing(msg tea.Msg) *TestEnv {
	e.T.Helper()

	modelInterface, cmd := e.Model.Update(msg)
	model, ok := modelInterface.(*app.Model)
	if ok {
		e.Model = model
	}
	e.processFormCmds(cmd, 10)

	return e
}

// PressKeyRuneWithFormProcessing sends a rune key and processes all resulting
// internal form messages. Use this when a key press opens or initializes a
// huh form (e.g., pressing 'e' to open the metadata editor).
//
// Expected:
//   - r must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) PressKeyRuneWithFormProcessing(r rune) *TestEnv {
	e.T.Helper()

	return e.SendMessageWithFormProcessing(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
}

// TabWithFormProcessing sends a Tab key and processes all resulting internal
// form messages. Use this when navigating between fields in a huh form.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) TabWithFormProcessing() *TestEnv {
	e.T.Helper()

	return e.SendMessageWithFormProcessing(tea.KeyMsg{Type: tea.KeyTab})
}

// processFormCmds processes commands from form interactions.
// Unlike executeCmd, this processes ALL messages (including huh internals)
// but has a depth limit to prevent infinite loops. Commands that take longer
// than 600ms (cursor blink ticks) are skipped to avoid blocking.
func (e *TestEnv) processFormCmds(cmd tea.Cmd, maxDepth int) {
	if cmd == nil || maxDepth <= 0 {
		return
	}

	type result struct {
		msg tea.Msg
	}
	done := make(chan result, 1)
	go func() {
		done <- result{msg: cmd()}
	}()

	var msg tea.Msg
	select {
	case r := <-done:
		msg = r.msg
	case <-time.After(600 * time.Millisecond):
		return
	}

	if msg == nil {
		return
	}

	switch typedMsg := msg.(type) {
	case tea.BatchMsg:
		for _, batchCmd := range typedMsg {
			e.processFormCmds(batchCmd, maxDepth-1)
		}
	case nil:
		return
	default:
		modelInterface, nextCmd := e.Model.Update(msg)
		model, ok := modelInterface.(*app.Model)
		if !ok {
			e.T.Fatal("model type assertion failed: expected *app.Model")
		}
		e.Model = model
		e.processFormCmds(nextCmd, maxDepth-1)
	}
}

// CompleteOnboarding simulates completing the onboarding wizard.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) CompleteOnboarding(name, email string) *TestEnv {
	e.T.Helper()

	if !e.IsInOnboardingState() {
		return e
	}

	// Initialize the model first (required for huh forms to work)
	e.InitModel()

	// Step 1: Welcome + Name
	// Type the name
	e.TypeText(name)
	// Press Enter to advance to next step (uses form processing to handle internal huh messages)
	e.PressEnterWithFormProcessing()

	// Step 2: Contact - Email + Location
	// Type the email
	e.TypeText(email)
	// Press Enter to accept email and move to Location field
	e.PressEnterWithFormProcessing()
	// Press Enter again to accept empty Location and advance to Step 3
	e.PressEnterWithFormProcessing()

	// Step 3: Professional Details - 3 optional fields (Title, GitHub, Portfolio)
	// Press Enter 3 times to accept all empty fields and complete
	e.PressEnterWithFormProcessing()
	e.PressEnterWithFormProcessing()
	e.PressEnterWithFormProcessing()

	return e
}
