package support

import (
	"context"
	"testing"
	"time"

	"github.com/baphled/kariya/internal/cli/bootstrap"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
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
//nolint:thelper // Factory function, not a test helper.
func NewBDDTestingT(t *testing.T) *BDDTestingT {
	return &BDDTestingT{t: t}
}

// Helper marks this as a test helper.
//
// Side effects:
//   - None.
func (b *BDDTestingT) Helper() {}

// TempDir returns a temporary directory.
//
// Returns:
//   - A string path to a temporary directory.
//
// Side effects:
//   - Creates a temporary directory on the filesystem.
func (b *BDDTestingT) TempDir() string {
	return b.t.TempDir()
}

// Fatalf logs a fatal error.
//
// Expected:
//   - format is a printf-style format string.
//
// Side effects:
//   - Terminates the test with a fatal error.
func (b *BDDTestingT) Fatalf(format string, args ...interface{}) {
	b.t.Fatalf(format, args...)
}

// Errorf logs an error.
//
// Expected:
//   - format is a printf-style format string.
//
// Side effects:
//   - Logs an error to the test output.
func (b *BDDTestingT) Errorf(format string, args ...interface{}) {
	b.t.Errorf(format, args...)
}

// Fatal logs a fatal error.
//
// Expected:
//   - args are the values to log.
//
// Side effects:
//   - Terminates the test with a fatal error.
func (b *BDDTestingT) Fatal(args ...interface{}) {
	b.t.Fatal(args...)
}

// Error logs an error.
//
// Expected:
//   - args are the values to log.
//
// Side effects:
//   - Logs an error to the test output.
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
//nolint:thelper // Factory function, not a test helper.
func NewAppEnv(t *testing.T) *e2e.TestEnv {
	return e2e.Setup(&BDDTestingT{t: t})
}

// GetAppEnv retrieves the TestEnv from context.
func GetAppEnv(ctx context.Context) *e2e.TestEnv {
	env, ok := ctx.Value(appEnvKey{}).(*e2e.TestEnv)
	if !ok {
		return nil
	}
	return env
}

// WithAppEnv stores a TestEnv in the context.
func WithAppEnv(ctx context.Context, env *e2e.TestEnv) context.Context {
	return context.WithValue(ctx, appEnvKey{}, env)
}

// GetEventData retrieves event data being built from context.
func GetEventData(ctx context.Context) *EventData {
	data, ok := ctx.Value(eventDataKey{}).(*EventData)
	if !ok {
		return &EventData{}
	}
	return data
}

// WithEventData stores event data in the context.
func WithEventData(ctx context.Context, data *EventData) context.Context {
	return context.WithValue(ctx, eventDataKey{}, data)
}

// BuildEvent creates a career.Event from EventData.
// Returns an error if date parsing fails.
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
