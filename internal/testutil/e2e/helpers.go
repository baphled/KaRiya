// Package e2e provides E2E test utilities for integration testing of the KaRiya TUI.
package e2e

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"

	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	_ "modernc.org/sqlite"
)

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

	// Repositories (SQLite versions)
	EventRepo *careerrepo.SQLiteRepository
	BurstRepo *careerrepo.SQLiteBurstRepository
	FactRepo  *careerrepo.SQLiteFactRepository
	SkillRepo *careerrepo.SQLiteSkillRepository

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
//
// Works with both *testing.T and GinkgoT().
func Setup(t TestingT) *TestEnv {
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
		_ = db.Close() // Ignore error as we're already in failure path
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Create repositories with SQLite
	eventRepo := careerrepo.NewSQLiteRepositoryWithDB(db)
	burstRepo := careerrepo.NewSQLiteBurstRepositoryWithDB(db)
	factRepo := careerrepo.NewSQLiteFactRepositoryWithDB(db)
	skillRepo := careerrepo.NewSQLiteSkillRepositoryWithDB(db)

	// Create service
	svc := careerservice.NewService(eventRepo)
	svc.SetBurstRepository(burstRepo)
	svc.SetFactRepository(factRepo)
	svc.SetSkillRepository(skillRepo)

	// Create CLI service
	cliService := service.NewCLIEventService(svc)

	// Create application model
	model := app.NewModel(cliService, svc)

	// Skip onboarding by default for E2E tests
	// Tests that need to test onboarding should use SetupWithOnboarding
	model.SkipOnboarding()

	cleanup := func() {
		_ = db.Close()
		_ = db.Close() // Error ignored as this is test cleanup
	}

	return &TestEnv{
		T:          t,
		Model:      model,
		DB:         db,
		DBPath:     dbPath,
		EventRepo:  eventRepo,
		BurstRepo:  burstRepo,
		FactRepo:   factRepo,
		SkillRepo:  skillRepo,
		Service:    svc,
		CLIService: cliService,
		Ctx:        ctx,
		cleanup:    cleanup,
	}
}

// SetupWithOnboarding creates an E2E test environment with the onboarding wizard active.
// Use this to test the onboarding workflow specifically.
// This forces onboarding to appear regardless of the user's config file.
//
// IMPORTANT: This function isolates config file writes to a temporary directory
// to prevent tests from polluting the user's real config file (BUG-007 fix).
//
// Works with both *testing.T and GinkgoT().
func SetupWithOnboarding(t TestingT) *TestEnv {
	t.Helper()

	ctx := context.Background()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "e2e_test.db")

	// BUG-007 FIX: Isolate config file writes to temp directory
	// This prevents onboarding completion from polluting the user's real config
	configPath := filepath.Join(tmpDir, "config.yaml")
	config.SetConfigPathForTesting(configPath)

	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		config.ResetConfigPath() // Clean up on failure
		t.Fatalf("failed to open test db: %v", err)
	}

	// Run migrations
	if err := careerrepo.RunMigrations(db); err != nil {
		config.ResetConfigPath() // Clean up on failure
		_ = db.Close()           // Ignore error as we're already in failure path
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Create repositories with SQLite
	eventRepo := careerrepo.NewSQLiteRepositoryWithDB(db)
	burstRepo := careerrepo.NewSQLiteBurstRepositoryWithDB(db)
	factRepo := careerrepo.NewSQLiteFactRepositoryWithDB(db)
	skillRepo := careerrepo.NewSQLiteSkillRepositoryWithDB(db)

	// Create service
	svc := careerservice.NewService(eventRepo)
	svc.SetBurstRepository(burstRepo)
	svc.SetFactRepository(factRepo)
	svc.SetSkillRepository(skillRepo)

	// Create CLI service
	cliService := service.NewCLIEventService(svc)

	// Create application model and FORCE onboarding
	// This ensures onboarding appears regardless of user's config file
	model := app.NewModel(cliService, svc)
	model.ForceOnboarding()

	cleanup := func() {
		// BUG-007 FIX: Reset config path override to prevent pollution
		config.ResetConfigPath()
		_ = db.Close()
		_ = db.Close() // Error ignored as this is test cleanup
	}

	return &TestEnv{
		T:          t,
		Model:      model,
		DB:         db,
		DBPath:     dbPath,
		EventRepo:  eventRepo,
		BurstRepo:  burstRepo,
		FactRepo:   factRepo,
		SkillRepo:  skillRepo,
		Service:    svc,
		CLIService: cliService,
		Ctx:        ctx,
		cleanup:    cleanup,
	}
}

// SetupWithMemory creates an E2E test environment using in-memory repositories.
// This is faster but doesn't test actual SQLite persistence.
//
// Works with both *testing.T and GinkgoT().
func SetupWithMemory(t TestingT) *TestEnv {
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

	// Skip onboarding by default for E2E tests
	model.SkipOnboarding()

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
// Valid names: "capture_event", "browse_timeline", "manage_skills", "generate_cv",
// "configure_system", "burst_management", "fact_management"
//
// NOTE: This order must match the menu items defined in internal/cli/app/app.go
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

// SubmitHuhForm submits a huh form by pressing Enter.
// Huh forms are submitted with Enter when the form is complete.
// This is equivalent to Confirm() but with a more descriptive name for form contexts.
func (e *TestEnv) SubmitHuhForm() *TestEnv {
	return e.Confirm()
}

// executeCmd executes commands returned by Update, but only for specific message types
// that are essential for state transitions (like form submission).
//
// Most Bubble Tea commands (cursor blink, window resize) are ignored because they
// cause infinite loops or stuck goroutines in tests. We only care about messages
// that actually change application state.
func (e *TestEnv) executeCmd(cmd tea.Cmd) {
	if cmd == nil {
		return
	}

	msg := cmd()
	if msg == nil {
		return
	}

	// Only process messages that are essential for state transitions
	// Skip all other messages to avoid infinite loops from huh forms (cursor blink, etc.)
	switch msg.(type) {
	case models.SubmitMsg:
		// Form submission - essential for form → review state transition
		modelInterface, nextCmd := e.Model.Update(msg)
		e.Model = modelInterface.(*app.Model)
		// Recursively execute any returned command
		e.executeCmd(nextCmd)

	case intents.SubmitCompleteMsg, intents.SubmitErrorMsg:
		// Submit completion - essential for submit → complete state transition
		modelInterface, nextCmd := e.Model.Update(msg)
		e.Model = modelInterface.(*app.Model)
		// For SubmitCompleteMsg, immediately send DismissModalMsg to skip the 2s timer
		if _, ok := msg.(intents.SubmitCompleteMsg); ok {
			// Skip the tea.Tick timer by directly sending DismissModalMsg
			modelInterface, nextCmd = e.Model.Update(intents.DismissModalMsg{})
			e.Model = modelInterface.(*app.Model)
		}
		// Recursively execute any returned command
		e.executeCmd(nextCmd)

	case intents.DismissModalMsg:
		// Modal dismissal - essential for success modal → enrichment review transition
		modelInterface, nextCmd := e.Model.Update(msg)
		e.Model = modelInterface.(*app.Model)
		// Recursively execute any returned command
		e.executeCmd(nextCmd)

	case intents.ConfigCompleteMsg:
		// Configuration save completion - essential for saving → complete state transition
		modelInterface, nextCmd := e.Model.Update(msg)
		e.Model = modelInterface.(*app.Model)
		// Recursively execute any returned command
		e.executeCmd(nextCmd)
	default:
		// Ignore all other messages (cursor blink, window resize, etc.)
		return
	}
}

// SendMessage sends a message directly to the model.
// This is useful for testing state transitions without simulating keystrokes.
// Returns the environment for method chaining.
func (e *TestEnv) SendMessage(msg tea.Msg) *TestEnv {
	e.T.Helper()

	modelInterface, cmd := e.Model.Update(msg)
	e.Model = modelInterface.(*app.Model)
	e.executeCmd(cmd)

	return e
}

// SubmitEvent sends a SubmitMsg directly to the model with the given event.
// This bypasses huh form navigation issues in E2E tests.
// Use this when you need to test the workflow after form submission.
func (e *TestEnv) SubmitEvent(event *career.CareerEvent) *TestEnv {
	e.T.Helper()

	return e.SendMessage(models.SubmitMsg{Event: event, Err: nil})
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
	skillRepo := careerrepo.NewSQLiteSkillRepositoryWithDB(e.DB)

	// Create new service
	svc := careerservice.NewService(eventRepo)
	svc.SetBurstRepository(burstRepo)
	svc.SetFactRepository(factRepo)
	svc.SetSkillRepository(skillRepo)

	// Create new CLI service
	cliService := service.NewCLIEventService(svc)

	// Create new application model
	model := app.NewModel(cliService, svc)

	// Skip onboarding by default for E2E tests (same as Setup)
	// Without this, tests that use SimulateRestart would show the onboarding wizard
	// instead of the expected main menu
	model.SkipOnboarding()

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

// ============================================================================
// Onboarding Helpers
// ============================================================================

// IsInOnboardingState checks if the application is currently showing the onboarding wizard.
func (e *TestEnv) IsInOnboardingState() bool {
	return e.Model.GetState() == app.StateOnboarding
}

// SkipOnboarding skips the onboarding wizard.
// Use Setup() instead of SetupWithOnboarding() to automatically skip onboarding.
func (e *TestEnv) SkipOnboarding() *TestEnv {
	e.Model.SkipOnboarding()
	return e
}

// InitModel initializes the model by calling Init() and sending a WindowSizeMsg.
// This is required for huh forms to render their content properly.
// Returns the environment for method chaining.
func (e *TestEnv) InitModel() *TestEnv {
	e.T.Helper()

	// Call Init() to set up the model
	cmd := e.Model.Init()

	// Process the init commands - this triggers huh form setup
	e.processFormCmds(cmd, 20)

	// Send a WindowSizeMsg to trigger form layout
	e.SendMessage(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Type and delete a character to force the form to render its fields
	// This workaround activates huh's internal rendering state
	e.PressKeyRune('x')
	e.PressKey(tea.KeyBackspace)

	return e
}

// PressEnterWithFormProcessing presses Enter and processes any internal form messages.
// This is needed because huh forms use internal messages (nextGroupMsg) to transition
// between groups. These messages must be processed for the form to advance.
//
// Note: This method loops to process messages but has a safety limit to prevent infinite loops.
func (e *TestEnv) PressEnterWithFormProcessing() *TestEnv {
	e.T.Helper()

	// Send Enter key
	modelInterface, cmd := e.Model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	e.Model = modelInterface.(*app.Model)

	// Process all resulting messages (including internal form messages)
	// This allows huh's group transitions to complete
	e.processFormCmds(cmd, 10) // max 10 iterations for safety

	return e
}

// processFormCmds processes commands from form interactions.
// Unlike executeCmd, this processes ALL messages (including huh internals)
// but has a depth limit to prevent infinite loops.
func (e *TestEnv) processFormCmds(cmd tea.Cmd, maxDepth int) {
	if cmd == nil || maxDepth <= 0 {
		return
	}

	msg := cmd()
	if msg == nil {
		return
	}

	switch m := msg.(type) {
	case tea.BatchMsg:
		// BatchMsg contains multiple commands - process each one
		for _, batchCmd := range m {
			e.processFormCmds(batchCmd, maxDepth-1)
		}
	case nil:
		return
	default:
		// Process the message and any follow-up commands
		// This includes internal huh messages like nextGroupMsg
		modelInterface, nextCmd := e.Model.Update(msg)
		e.Model = modelInterface.(*app.Model)
		e.processFormCmds(nextCmd, maxDepth-1)
	}
}

// CompleteOnboarding simulates completing the onboarding wizard.
// This types name and email, then presses Enter to advance through steps.
// name: Required field (e.g., "Test User")
// email: Required field (e.g., "test@example.com")
//
// Returns the environment for method chaining.
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
	e.PressEnterWithFormProcessing() // Title (optional)
	e.PressEnterWithFormProcessing() // GitHub (optional)
	e.PressEnterWithFormProcessing() // Portfolio (optional) - completes form

	return e
}
