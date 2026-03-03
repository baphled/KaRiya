package support

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/cli/bootstrap"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careersql "github.com/baphled/kariya/internal/repository/career/sql"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/harness"
	tea "github.com/charmbracelet/bubbletea"
	"gorm.io/gorm"
)

// onboardingKey is the context key for storing OnboardingTestModel.
type onboardingKey struct{}

// OnboardingEnv wraps the onboarding test model for BDD testing.
// It provides a simplified interface for interacting with the onboarding wizard.
type OnboardingEnv struct {
	model *bootstrap.OnboardingTestModel
}

// NewOnboardingEnv creates a new OnboardingEnv with an initialized model.
//
// Expected: None.
// Returns: A fully initialized OnboardingEnv ready for use.
// Side effects: None.
func NewOnboardingEnv() *OnboardingEnv {
	model := bootstrap.NewOnboardingTestModel(nil)
	env := &OnboardingEnv{
		model: model,
	}
	initCmd := model.Init()
	env.sendWindowSize(80, 24)
	env.processFormCmds(initCmd, 10)
	return env
}

// View returns the current rendered view.
//
// Expected: None.
// Returns: The current view as a string.
// Side effects: None.
func (e *OnboardingEnv) View() string {
	return e.model.View()
}

// TypeText types text character by character.
//
// Expected: text is the string to type.
// Returns: None.
// Side effects: Updates the model state.
func (e *OnboardingEnv) TypeText(text string) {
	for _, r := range text {
		_, cmd := e.model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		e.processFormCmds(cmd, 5)
	}
}

// PressKey sends a key press and processes resulting commands.
//
// Expected: key is one of "enter", "tab", "shift+tab", "up", "down", or a rune string.
// Returns: None.
// Side effects: Updates the model state.
func (e *OnboardingEnv) PressKey(key string) {
	var msg tea.KeyMsg
	switch key {
	case "enter":
		msg = tea.KeyMsg{Type: tea.KeyEnter}
	case "tab":
		msg = tea.KeyMsg{Type: tea.KeyTab}
	case "shift+tab":
		msg = tea.KeyMsg{Type: tea.KeyShiftTab}
	case "up":
		msg = tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		msg = tea.KeyMsg{Type: tea.KeyDown}
	default:
		msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	}
	_, cmd := e.model.Update(msg)
	e.processFormCmds(cmd, 10)
}

// PressEnter presses the enter key.
//
// Expected: None.
// Returns: None.
// Side effects: Updates the model state.
func (e *OnboardingEnv) PressEnter() {
	e.PressKey("enter")
}

// PressTab presses the tab key.
//
// Expected: None.
// Returns: None.
// Side effects: Updates the model state.
func (e *OnboardingEnv) PressTab() {
	e.PressKey("tab")
}

// IsCompleted returns whether the wizard has completed.
//
// Expected: None.
// Returns: true if the wizard is complete, false otherwise.
// Side effects: None.
func (e *OnboardingEnv) IsCompleted() bool {
	return e.model.IsCompleted()
}

// Result returns the profile config result after completion.
//
// Expected: Wizard should be completed.
// Returns: The profile configuration or nil if not completed.
// Side effects: None.
func (e *OnboardingEnv) Result() *config.ProfileConfig {
	return e.model.Result()
}

// sendWindowSize sends a window size message to the model.
func (e *OnboardingEnv) sendWindowSize(width, height int) {
	e.model.Update(tea.WindowSizeMsg{Width: width, Height: height})
}

// processFormCmds processes commands from form interactions, including batch messages.
func (e *OnboardingEnv) processFormCmds(cmd tea.Cmd, maxDepth int) {
	if cmd == nil || maxDepth <= 0 {
		return
	}

	done := make(chan tea.Msg, 1)
	go func() {
		msg := cmd()
		done <- msg
	}()

	var msg tea.Msg
	select {
	case msg = <-done:
	case <-time.After(50 * time.Millisecond):
		return
	}

	if msg == nil {
		return
	}

	switch m := msg.(type) {
	case tea.BatchMsg:
		for _, batchCmd := range m {
			e.processFormCmds(batchCmd, maxDepth-1)
		}
	default:
		_, nextCmd := e.model.Update(msg)
		e.processFormCmds(nextCmd, maxDepth-1)
	}
}

// GetOnboardingEnv retrieves the OnboardingEnv from context.
//
// Expected: ctx contains an OnboardingEnv stored with WithOnboardingEnv.
// Returns: The OnboardingEnv or nil if not found.
// Side effects: None.
func GetOnboardingEnv(ctx context.Context) *OnboardingEnv {
	env, ok := ctx.Value(onboardingKey{}).(*OnboardingEnv)
	if !ok {
		return nil
	}
	return env
}

// WithOnboardingEnv stores an OnboardingEnv in the context.
//
// Expected: ctx is a valid context, env is a valid OnboardingEnv.
// Returns: A new context with the OnboardingEnv stored.
// Side effects: None.
func WithOnboardingEnv(ctx context.Context, env *OnboardingEnv) context.Context {
	return context.WithValue(ctx, onboardingKey{}, env)
}

// appEnvKey is the context key for storing the full app TestEnv.
type appEnvKey struct{}

// eventDataKey is the context key for storing event data being built.
type eventDataKey struct{}

// EventData holds event data being built during a scenario.
type EventData struct {
	Description string
	Date        string
	Company     string
	Project     string
	Tags        []string
	Categories  []string
	Skills      []string
}

// BDDTestingT adapts context for TestEnv's TestingT interface.
type BDDTestingT struct {
	t *testing.T
}

// NewBDDTestingT creates a new BDDTestingT adapter.
//
// Expected:
//   - t must be a valid *testing.T.
//
// Returns:
//   - A new BDDTestingT instance.
//
// Side effects:
//   - None.
//
//nolint:thelper // Constructor for BDDTestingT adapter - not a test helper that needs t.Helper().
func NewBDDTestingT(t *testing.T) *BDDTestingT {
	return &BDDTestingT{t: t}
}

// Helper marks this as a test helper.
//
// Returns:
//   - A {} value.
//
// Side effects:
//   - None.
func (b *BDDTestingT) Helper() {}

// TempDir returns a temporary directory.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (b *BDDTestingT) TempDir() string {
	return b.t.TempDir()
}

// Fatalf logs a fatal error.
//
// Expected:
//   - Must be a valid string.
//   - interface{} must be valid.
//
// Side effects:
//   - None.
func (b *BDDTestingT) Fatalf(format string, args ...interface{}) {
	b.t.Fatalf(format, args...)
}

// Errorf logs an error.
//
// Expected:
//   - Must be a valid string.
//   - interface{} must be valid.
//
// Side effects:
//   - None.
func (b *BDDTestingT) Errorf(format string, args ...interface{}) {
	b.t.Errorf(format, args...)
}

// Fatal logs a fatal error.
//
// Expected:
//   - interface{} must be valid.
//
// Side effects:
//   - None.
func (b *BDDTestingT) Fatal(args ...interface{}) {
	b.t.Fatal(args...)
}

// Error logs an error.
//
// Expected:
//   - interface{} must be valid.
//
// Side effects:
//   - None.
func (b *BDDTestingT) Error(args ...interface{}) {
	b.t.Error(args...)
}

// NewAppEnv creates a new full application TestEnv for BDD testing.
//
// Expected:
//   - t must be a valid *testing.T.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - Creates a temporary database for testing.
//
// this function is not a test assertion helper; it constructs the environment and any failure
// should point to this function, not the caller.
//
//nolint:thelper // Factory constructor that returns a new TestEnv — t.Helper() is inappropriate here because
func NewAppEnv(t *testing.T) *harness.TestEnv {
	return harness.Setup(&BDDTestingT{t: t})
}

// GetAppEnv retrieves the TestEnv from context.
//
// Expected:
//   - ctx is a valid context.Context.
//
// Returns:
//   - TestEnv pointer or nil if not found.
//
// Side effects:
//   - None.
func GetAppEnv(ctx context.Context) *harness.TestEnv {
	env, ok := ctx.Value(appEnvKey{}).(*harness.TestEnv)
	if !ok {
		return nil
	}
	return env
}

// ErrEnvNotInitialised is returned when the test environment has not been set up.
// A nil env means scenario setup failed — this is a test harness bug, not a pending step.
var ErrEnvNotInitialised = errors.New("test environment not initialised: check BeforeScenario hook")

// RequireEnv retrieves the test environment from context.
// Returns a hard error if env is nil — setup failure should be red, not pending.
//
// Expected:
//   - ctx is a valid context.Context with a TestEnv set via WithAppEnv.
//
// Returns:
//   - TestEnv pointer and nil error on success.
//   - nil and ErrEnvNotInitialised if env was not set.
//
// Side effects:
//   - None.
func RequireEnv(ctx context.Context) (*harness.TestEnv, error) {
	env := GetAppEnv(ctx)
	if env == nil {
		return nil, ErrEnvNotInitialised
	}
	return env, nil
}

// WithAppEnv stores a TestEnv in the context.
//
// Expected:
//   - testenv must be valid.
//
// Returns:
//   - A context.Context value.
//
// Side effects:
//   - None.
func WithAppEnv(ctx context.Context, env *harness.TestEnv) context.Context {
	return context.WithValue(ctx, appEnvKey{}, env)
}

// GetEventData retrieves event data being built from context.
//
// Expected:
//   - ctx is a valid context.Context.
//
// Returns:
//   - EventData pointer or empty EventData if not found.
//
// Side effects:
//   - None.
func GetEventData(ctx context.Context) *EventData {
	data, ok := ctx.Value(eventDataKey{}).(*EventData)
	if !ok {
		return &EventData{}
	}
	return data
}

// WithEventData stores event data in the context.
//
// Expected:
//   - ctx is a valid context.Context.
//   - data is a valid EventData pointer.
//
// Returns:
//   - New context with EventData stored.
//
// Side effects:
//   - None.
func WithEventData(ctx context.Context, data *EventData) context.Context {
	return context.WithValue(ctx, eventDataKey{}, data)
}

// BuildEvent creates a career.Event from EventData.
//
// Expected:
//   - EventData fields are populated.
//
// Returns:
//   - career.Event pointer if successful.
//   - error if date parsing fails.
//
// Side effects:
//   - None.
func (d *EventData) BuildEvent() (*career.Event, error) {
	event := &career.Event{
		Text:       d.Description,
		Company:    d.Company,
		Project:    d.Project,
		Tags:       d.Tags,
		Categories: d.Categories,
		Skills:     d.Skills,
	}

	if d.Date != "" {
		parsed, err := forms.ParseDateString(d.Date)
		if err != nil {
			return nil, err
		}
		event.Date = parsed
	}

	return event, nil
}

// NavigateToTableItem navigates to a specific item in a table by name.
// It searches the current view for the item text and navigates down until found.
//
// Expected:
//   - env must be a valid *harness.TestEnv.
//   - itemName is the text to search for in the table.
//   - maxAttempts is the maximum number of down presses (default 20).
//
// Returns:
//   - error if item not found after maxAttempts.
//
// Side effects:
//   - Presses 'Home' to go to table start.
//   - Presses 'Down' repeatedly until item found.
func NavigateToTableItem(env *harness.TestEnv, itemName string, maxAttempts int) error {
	if maxAttempts == 0 {
		maxAttempts = 20
	}

	// Go to start of table
	env.PressKey(tea.KeyHome)

	// Search for item by navigating down
	for range maxAttempts {
		view := env.GetView()
		if containsSubstring(view, itemName) && isItemSelected(view, itemName) {
			return nil
		}
		env.NavigateDown()
	}

	return fmt.Errorf("item %q not found in table after %d attempts", itemName, maxAttempts)
}

// WaitForViewContains polls the view until it contains the expected string.
//
// Expected:
//   - env must be a valid *harness.TestEnv.
//   - expected is the substring to wait for.
//   - maxAttempts is the maximum number of polls (default 10).
//   - delayMs is the delay between polls in milliseconds (default 50ms).
//
// Returns:
//   - error if string not found after maxAttempts.
//
// Side effects:
//   - Polls the view multiple times with delay.
func WaitForViewContains(env *harness.TestEnv, expected string, maxAttempts int, delayMs int) error {
	if maxAttempts == 0 {
		maxAttempts = 10
	}
	if delayMs == 0 {
		delayMs = 50
	}

	for range maxAttempts {
		view := env.GetView()
		if containsSubstring(view, expected) {
			return nil
		}
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}

	return fmt.Errorf("view does not contain %q after %d attempts", expected, maxAttempts)
}

// containsSubstring checks if a string contains a substring.
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

// findSubstring performs a simple substring search.
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// isItemSelected checks if the item appears to be selected (has "▶" marker).
func isItemSelected(view, itemName string) bool {
	lines := splitLines(view)
	for _, line := range lines {
		if containsSubstring(line, "▶") && containsSubstring(line, itemName) {
			return true
		}
	}
	return false
}

// splitLines splits a string into lines.
func splitLines(s string) []string {
	var lines []string
	var current string
	for _, r := range s {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

// =============================================================================
// CLI Environment for testing command-line execution
// =============================================================================

// cliEnvKey is the context key for storing CLIEnv.
type cliEnvKey struct{}

// CLIEnv wraps CLI command execution for BDD testing.
type CLIEnv struct {
	binaryPath string
	stdout     *bytes.Buffer
	stderr     *bytes.Buffer
	exitCode   int
}

// NewCLIEnv creates a new CLIEnv with the given binary path.
//
// Expected: binaryPath is the path to the kariya binary.
// Returns: A new CLIEnv instance.
// Side effects: None.
func NewCLIEnv(binaryPath string) *CLIEnv {
	return &CLIEnv{
		binaryPath: binaryPath,
		stdout:     &bytes.Buffer{},
		stderr:     &bytes.Buffer{},
		exitCode:   0,
	}
}

// Run executes the CLI with the given arguments.
//
// Expected: args are command-line arguments.
// Returns: error if execution fails.
// Side effects: Captures stdout, stderr, and exit code.
func (e *CLIEnv) Run(args ...string) error {
	e.stdout.Reset()
	e.stderr.Reset()
	e.exitCode = 0

	// #nosec G204 - binaryPath is controlled by test code, args come from test scenarios
	cmd := exec.Command(e.binaryPath, args...)
	cmd.Stdout = e.stdout
	cmd.Stderr = e.stderr

	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			return err
		}
		e.exitCode = exitErr.ExitCode()
	}

	return nil
}

// GetOutput returns the combined stdout output.
//
// Expected: Run has been called.
// Returns: stdout as string.
// Side effects: None.
func (e *CLIEnv) GetOutput() string {
	return e.stdout.String()
}

// GetErrorOutput returns the stderr output.
//
// Expected: Run has been called.
// Returns: stderr as string.
// Side effects: None.
func (e *CLIEnv) GetErrorOutput() string {
	return e.stderr.String()
}

// GetExitCode returns the exit code from the last command.
//
// Expected: Run has been called.
// Returns: exit code (0 for success).
// Side effects: None.
func (e *CLIEnv) GetExitCode() int {
	return e.exitCode
}

// GetCLIEnv retrieves the CLIEnv from context.
//
// Expected: ctx contains a CLIEnv stored with WithCLIEnv.
// Returns: The CLIEnv or nil if not found.
// Side effects: None.
func GetCLIEnv(ctx context.Context) *CLIEnv {
	env, ok := ctx.Value(cliEnvKey{}).(*CLIEnv)
	if !ok {
		return nil
	}
	return env
}

// WithCLIEnv stores a CLIEnv in the context.
//
// Expected: ctx is a valid context, env is a valid CLIEnv.
// Returns: A new context with the CLIEnv stored.
// Side effects: None.
func WithCLIEnv(ctx context.Context, env *CLIEnv) context.Context {
	return context.WithValue(ctx, cliEnvKey{}, env)
}

// NewAppEnvFromGormDB creates a TestEnv from a GORM database connection (typically a transaction).
// This avoids per-scenario DB creation + migration overhead by reusing a shared connection.
//
// Expected:
//   - t must be a valid *testing.T.
//   - gormDB must be a valid *gorm.DB (can be a transaction).
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - Creates repositories, services, and application model from the GORM connection.
//
// this function is not a test assertion helper; it constructs the environment and any failure
// should point to this function, not the caller.
//
//nolint:thelper // Factory constructor that returns a new TestEnv — t.Helper() is inappropriate here because
func NewAppEnvFromGormDB(t *testing.T, gormDB *gorm.DB) *harness.TestEnv {
	ctx := context.Background()

	repos := careersql.NewRepositoriesFromDB(gormDB)

	svc := careerservice.NewService(repos.Event)
	svc.SetBurstRepository(repos.Burst)
	svc.SetFactRepository(repos.Fact)
	svc.SetSkillRepository(repos.Skill)

	cliService := service.NewCLIEventService(svc)

	log := logger.DefaultLogger()
	bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

	model := app.NewModel(cliService, svc, bootstrapResult)
	model.Update(tea.WindowSizeMsg{Width: harness.TerminalWidth, Height: harness.TerminalHeightLarge})

	return &harness.TestEnv{
		T:          &BDDTestingT{t: t},
		Model:      model,
		EventRepo:  repos.Event,
		BurstRepo:  repos.Burst,
		FactRepo:   repos.Fact,
		SkillRepo:  repos.Skill,
		Service:    svc,
		CLIService: cliService,
		Ctx:        ctx,
	}
}
