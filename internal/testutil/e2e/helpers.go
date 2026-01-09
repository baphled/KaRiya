// Package e2e provides E2E test utilities for integration testing of the KaRiya TUI.
package e2e

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	_ "modernc.org/sqlite"
)

// TestEnv holds all test dependencies for E2E testing.
// It provides a complete test environment with SQLite persistence,
// repositories, services, and the application model.
type TestEnv struct {
	// T is the testing context
	T *testing.T

	// Model is the root application model
	Model *app.Model

	// DB is the SQLite database connection
	DB *sql.DB

	// DBPath is the path to the SQLite database file
	DBPath string

	// Repositories (SQLite versions)
	EventRepo *careerrepo.SQLiteRepository
	BurstRepo *careerrepo.SQLiteBurstRepository
	FactRepo  *careerrepo.SQLiteFactRepository

	// Memory repositories (for fast tests)
	MemEventRepo *careerrepo.MemoryRepository
	MemBurstRepo *careerrepo.MemoryBurstRepository
	MemFactRepo  *careerrepo.MemoryFactRepository

	// Services
	Service    *careerservice.Service
	CLIService *service.CLIEventService

	// Context for async operations
	Ctx context.Context

	// Cleanup function to call when done
	cleanup func()
}

// Setup creates a complete E2E test environment with SQLite persistence.
// This sets up all repositories, services, and the application model.
//
// Usage:
//
//	env := e2e.Setup(t)
//	defer env.Cleanup()
//	// use env for testing
func Setup(t *testing.T) *TestEnv {
	t.Helper()

	ctx := context.Background()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "e2e_test.db")

	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// Run migrations
	if err := careerrepo.RunMigrations(db); err != nil {
		db.Close()
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Create repositories with SQLite
	eventRepo := careerrepo.NewSQLiteRepositoryWithDB(db)
	burstRepo := careerrepo.NewSQLiteBurstRepositoryWithDB(db)
	factRepo := careerrepo.NewSQLiteFactRepositoryWithDB(db)

	// Create service
	svc := careerservice.NewService(eventRepo)
	svc.SetBurstRepository(burstRepo)
	svc.SetFactRepository(factRepo)

	// Create CLI service
	cliService := service.NewCLIEventService(svc)

	// Create application model
	model := app.NewModel(cliService, svc)

	cleanup := func() {
		db.Close()
	}

	return &TestEnv{
		T:          t,
		Model:      model,
		DB:         db,
		DBPath:     dbPath,
		EventRepo:  eventRepo,
		BurstRepo:  burstRepo,
		FactRepo:   factRepo,
		Service:    svc,
		CLIService: cliService,
		Ctx:        ctx,
		cleanup:    cleanup,
	}
}

// SetupWithMemory creates an E2E test environment using in-memory repositories.
// This is faster but doesn't test actual SQLite persistence.
func SetupWithMemory(t *testing.T) *TestEnv {
	t.Helper()

	ctx := context.Background()

	// Create in-memory repositories
	eventRepo := careerrepo.NewMemoryRepository()
	burstRepo := careerrepo.NewMemoryBurstRepository()
	factRepo := careerrepo.NewMemoryFactRepository()

	// Create service
	svc := careerservice.NewService(eventRepo)
	svc.SetBurstRepository(burstRepo)
	svc.SetFactRepository(factRepo)

	// Create CLI service
	cliService := service.NewCLIEventService(svc)

	// Create application model
	model := app.NewModel(cliService, svc)

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
		cleanup:      func() {},
	}
}

// Cleanup releases all test resources.
// Should be called with defer immediately after Setup.
func (e *TestEnv) Cleanup() {
	if e.cleanup != nil {
		e.cleanup()
	}
}

// ============================================================================
// Navigation Helpers
// ============================================================================

// SelectIntent navigates to and selects a menu item by its index (0-based).
// Returns the environment for method chaining.
func (e *TestEnv) SelectIntent(index int) *TestEnv {
	e.T.Helper()

	// Navigate to the menu item
	for i := 0; i < index; i++ {
		e.PressKeyRune('j')
	}

	// Select the intent
	e.PressKey(tea.KeyEnter)

	return e
}

// SelectIntentByName navigates to and selects a menu item by its intent name.
// Valid names: "capture_event", "browse_timeline", "generate_cv", "export_artifact",
// "configure_system", "burst_management", "fact_management", "import_wizard",
// "metadata_editor", "bulk_operations"
func (e *TestEnv) SelectIntentByName(name string) *TestEnv {
	e.T.Helper()

	intentOrder := map[string]int{
		"capture_event":    0,
		"browse_timeline":  1,
		"generate_cv":      2,
		"export_artifact":  3,
		"configure_system": 4,
		"burst_management": 5,
		"fact_management":  6,
		"import_wizard":    7,
		"metadata_editor":  8,
		"bulk_operations":  9,
	}

	index, ok := intentOrder[name]
	if !ok {
		e.T.Fatalf("unknown intent name: %s", name)
	}

	return e.SelectIntent(index)
}

// PressKey sends a key message to the model.
// Returns the environment for method chaining.
func (e *TestEnv) PressKey(key tea.KeyType) *TestEnv {
	e.T.Helper()

	modelInterface, cmd := e.Model.Update(tea.KeyMsg{Type: key})
	e.Model = modelInterface.(*app.Model)

	// Execute any returned command
	e.executeCmd(cmd)

	return e
}

// PressKeyRune sends a rune key message to the model.
// Returns the environment for method chaining.
func (e *TestEnv) PressKeyRune(r rune) *TestEnv {
	e.T.Helper()

	modelInterface, cmd := e.Model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	e.Model = modelInterface.(*app.Model)

	// Execute any returned command
	e.executeCmd(cmd)

	return e
}

// PressKeys sends multiple keys in sequence.
// Accepts tea.KeyType or rune values.
// Returns the environment for method chaining.
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
// Returns the environment for method chaining.
func (e *TestEnv) TypeText(text string) *TestEnv {
	e.T.Helper()

	for _, r := range text {
		e.PressKeyRune(r)
	}

	return e
}

// NavigateDown moves down in a list (j or down arrow).
func (e *TestEnv) NavigateDown() *TestEnv {
	return e.PressKeyRune('j')
}

// NavigateUp moves up in a list (k or up arrow).
func (e *TestEnv) NavigateUp() *TestEnv {
	return e.PressKeyRune('k')
}

// Confirm presses Enter to confirm an action.
func (e *TestEnv) Confirm() *TestEnv {
	return e.PressKey(tea.KeyEnter)
}

// Cancel presses Escape to cancel/go back.
func (e *TestEnv) Cancel() *TestEnv {
	return e.PressKey(tea.KeyEscape)
}

// GoBack presses Escape to go back.
func (e *TestEnv) GoBack() *TestEnv {
	return e.Cancel()
}

// Quit presses 'q' to quit.
func (e *TestEnv) Quit() *TestEnv {
	return e.PressKeyRune('q')
}

// Tab presses Tab to move to next field.
func (e *TestEnv) Tab() *TestEnv {
	return e.PressKey(tea.KeyTab)
}

// executeCmd executes a tea.Cmd and processes the result.
func (e *TestEnv) executeCmd(cmd tea.Cmd) {
	if cmd == nil {
		return
	}

	msg := cmd()
	if msg != nil {
		modelInterface, newCmd := e.Model.Update(msg)
		e.Model = modelInterface.(*app.Model)
		// Recursively execute any new commands (but limit depth to prevent infinite loops)
		if newCmd != nil {
			e.executeCmd(newCmd)
		}
	}
}

// ============================================================================
// View Assertion Helpers
// ============================================================================

// GetView returns the current view output.
func (e *TestEnv) GetView() string {
	return e.Model.View()
}

// AssertViewContains checks that the view contains the given substring.
// Returns the environment for method chaining.
func (e *TestEnv) AssertViewContains(substr string) *TestEnv {
	e.T.Helper()

	view := e.GetView()
	if !strings.Contains(view, substr) {
		e.T.Errorf("expected view to contain %q, but it doesn't.\nView:\n%s", substr, view)
	}

	return e
}

// AssertViewNotContains checks that the view does NOT contain the given substring.
// Returns the environment for method chaining.
func (e *TestEnv) AssertViewNotContains(substr string) *TestEnv {
	e.T.Helper()

	view := e.GetView()
	if strings.Contains(view, substr) {
		e.T.Errorf("expected view NOT to contain %q, but it does.\nView:\n%s", substr, view)
	}

	return e
}

// AssertViewContainsAny checks that the view contains at least one of the given substrings.
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
func (e *TestEnv) AssertEventCount(expected int) *TestEnv {
	e.T.Helper()

	events, err := e.Service.ListEvents(e.Ctx, careerrepo.ListFilters{})
	if err != nil {
		e.T.Fatalf("failed to get events: %v", err)
	}

	if len(events) != expected {
		e.T.Errorf("expected %d events, got %d", expected, len(events))
	}

	return e
}

// AssertBurstCount verifies the number of bursts in the database.
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
func (e *TestEnv) GetEvents() []*career.CareerEvent {
	e.T.Helper()

	events, err := e.Service.ListEvents(e.Ctx, careerrepo.ListFilters{})
	if err != nil {
		e.T.Fatalf("failed to get events: %v", err)
	}

	return events
}

// GetBursts returns all bursts from the database.
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

// ============================================================================
// Session Simulation
// ============================================================================

// SimulateRestart recreates the application model while preserving the database.
// This simulates an application restart.
func (e *TestEnv) SimulateRestart() *TestEnv {
	e.T.Helper()

	if e.DB == nil {
		e.T.Fatal("SimulateRestart requires SQLite persistence (use Setup, not SetupWithMemory)")
	}

	// Create new repositories pointing to the same database
	eventRepo := careerrepo.NewSQLiteRepositoryWithDB(e.DB)
	burstRepo := careerrepo.NewSQLiteBurstRepositoryWithDB(e.DB)
	factRepo := careerrepo.NewSQLiteFactRepositoryWithDB(e.DB)

	// Create new service
	svc := careerservice.NewService(eventRepo)
	svc.SetBurstRepository(burstRepo)
	svc.SetFactRepository(factRepo)

	// Create new CLI service
	cliService := service.NewCLIEventService(svc)

	// Create new application model
	model := app.NewModel(cliService, svc)

	// Update environment
	e.EventRepo = eventRepo
	e.BurstRepo = burstRepo
	e.FactRepo = factRepo
	e.Service = svc
	e.CLIService = cliService
	e.Model = model

	return e
}

// ============================================================================
// Data Population Helpers
// ============================================================================

// AddEvent creates an event in the database.
func (e *TestEnv) AddEvent(event *career.CareerEvent) *TestEnv {
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
func (e *TestEnv) GetMenuItems() []app.MenuItem {
	return e.Model.GetMenuItems()
}

// IsInMenuState checks if the application is currently showing the main menu.
func (e *TestEnv) IsInMenuState() bool {
	view := e.GetView()
	// Menu shows the tagline and menu items
	return strings.Contains(view, "Career Event Management System") &&
		strings.Contains(view, "Capture Event")
}
