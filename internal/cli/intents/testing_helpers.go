package intents

import (
	"time"

	"github.com/baphled/kariya/internal/cli/themes"
	tea "github.com/charmbracelet/bubbletea"
)

// ErrorMsg is a message type for communicating errors from mocks.
type ErrorMsg struct {
	Err error
}

// MockIntent is a mock implementation of Intent for testing.
type MockIntent struct {
	InitCalled   bool
	UpdateCalled int
	ViewCalled   bool
	Result_      *IntentResult[interface{}]
	Messages     []tea.Msg
}

// NewMockIntent creates a new MockIntent for testing.
func NewMockIntent() *MockIntent {
	return &MockIntent{
		Messages: make([]tea.Msg, 0),
	}
}

// Init implements Intent.Init.
func (m *MockIntent) Init() tea.Cmd {
	m.InitCalled = true
	return nil
}

// Update implements Intent.Update.
func (m *MockIntent) Update(msg tea.Msg) tea.Cmd {
	m.UpdateCalled++
	m.Messages = append(m.Messages, msg)
	return nil
}

// View implements Intent.View.
func (m *MockIntent) View() string {
	m.ViewCalled = true
	return "Mock Intent View"
}

// Result implements Intent.Result.
func (m *MockIntent) Result() *IntentResult[interface{}] {
	return m.Result_
}

// SetResult sets the result for this mock intent.
func (m *MockIntent) SetResult(result *IntentResult[interface{}]) {
	m.Result_ = result
}

// ThemeAwareMockIntent is a mock intent that supports theme management.
type ThemeAwareMockIntent struct {
	*MockIntent
	ThemeManager_ *themes.ThemeManager
}

// NewThemeAwareMockIntent creates a new ThemeAwareMockIntent for testing.
func NewThemeAwareMockIntent() *ThemeAwareMockIntent {
	return &ThemeAwareMockIntent{
		MockIntent: NewMockIntent(),
	}
}

// SetThemeManager implements ThemeAware.SetThemeManager.
func (m *ThemeAwareMockIntent) SetThemeManager(tm *themes.ThemeManager) {
	m.ThemeManager_ = tm
}

// GetThemeManager implements ThemeAware.GetThemeManager.
func (m *ThemeAwareMockIntent) GetThemeManager() *themes.ThemeManager {
	return m.ThemeManager_
}

// MockIntentWithSelection extends MockIntent with selection state for testing.
type MockIntentWithSelection struct {
	*MockIntent
	SelectedIndex int
}

// NewMockIntentWithSelection creates a new MockIntentWithSelection for testing.
func NewMockIntentWithSelection() *MockIntentWithSelection {
	return &MockIntentWithSelection{
		MockIntent:    NewMockIntent(),
		SelectedIndex: 0,
	}
}

// GetSelectedIndex returns the selected index.
func (m *MockIntentWithSelection) GetSelectedIndex() int {
	return m.SelectedIndex
}

// SetSelectedIndex sets the selected index.
func (m *MockIntentWithSelection) SetSelectedIndex(idx int) {
	m.SelectedIndex = idx
}

// StateTransition records a state transition with metadata.
type StateTransition struct {
	From      string
	To        string
	Trigger   string
	Timestamp time.Time
}

// MockIntentWithStateMachine extends MockIntent with state machine tracking for testing.
type MockIntentWithStateMachine struct {
	*MockIntent
	CurrentState string
	StateHistory []string
	Transitions  []StateTransition
}

// NewMockIntentWithStateMachine creates a new MockIntentWithStateMachine for testing.
func NewMockIntentWithStateMachine(initialState string) *MockIntentWithStateMachine {
	return &MockIntentWithStateMachine{
		MockIntent:   NewMockIntent(),
		CurrentState: initialState,
		StateHistory: []string{initialState},
		Transitions:  make([]StateTransition, 0),
	}
}

// TransitionTo transitions to a new state and records it in history.
func (m *MockIntentWithStateMachine) TransitionTo(state string) {
	m.StateHistory = append(m.StateHistory, state)
	m.CurrentState = state
}

// TransitionToWithTrigger transitions to a new state with trigger metadata.
func (m *MockIntentWithStateMachine) TransitionToWithTrigger(state, trigger string) {
	transition := StateTransition{
		From:      m.CurrentState,
		To:        state,
		Trigger:   trigger,
		Timestamp: time.Now(),
	}
	m.Transitions = append(m.Transitions, transition)
	m.TransitionTo(state)
}

// GetStateHistory returns the complete state history.
func (m *MockIntentWithStateMachine) GetStateHistory() []string {
	return m.StateHistory
}

// GetTransitions returns all recorded transitions.
func (m *MockIntentWithStateMachine) GetTransitions() []StateTransition {
	return m.Transitions
}

// Reset clears state history and transitions, starting fresh.
func (m *MockIntentWithStateMachine) Reset(initialState string) {
	m.CurrentState = initialState
	m.StateHistory = []string{initialState}
	m.Transitions = make([]StateTransition, 0)
}

// MockIntentWithErrors extends MockIntent with error injection capabilities for testing.
type MockIntentWithErrors struct {
	*MockIntent
	ErrorOnInit      error
	ErrorOnUpdate    error
	ErrorOnNthUpdate int // Fail on Nth update call (0-indexed)
	ErrorCount       int
	updateCount      int
}

// NewMockIntentWithErrors creates a new MockIntentWithErrors for testing.
func NewMockIntentWithErrors() *MockIntentWithErrors {
	return &MockIntentWithErrors{
		MockIntent:       NewMockIntent(),
		ErrorOnNthUpdate: -1, // Default: no delayed error
	}
}

// Init implements Intent.Init with error injection.
func (m *MockIntentWithErrors) Init() tea.Cmd {
	m.MockIntent.Init()
	if m.ErrorOnInit != nil {
		m.ErrorCount++
		return func() tea.Msg {
			return ErrorMsg{Err: m.ErrorOnInit}
		}
	}
	return nil
}

// Update implements Intent.Update with error injection.
func (m *MockIntentWithErrors) Update(msg tea.Msg) tea.Cmd {
	m.MockIntent.Update(msg)

	// Check for immediate error
	if m.ErrorOnUpdate != nil && m.ErrorOnNthUpdate < 0 {
		m.ErrorCount++
		return func() tea.Msg {
			return ErrorMsg{Err: m.ErrorOnUpdate}
		}
	}

	// Check for delayed error (on Nth update)
	if m.ErrorOnUpdate != nil && m.ErrorOnNthUpdate >= 0 {
		if m.updateCount == m.ErrorOnNthUpdate {
			m.updateCount++
			m.ErrorCount++
			return func() tea.Msg {
				return ErrorMsg{Err: m.ErrorOnUpdate}
			}
		}
		m.updateCount++
	}

	return nil
}

// SetInitError configures the mock to return an error on Init.
func (m *MockIntentWithErrors) SetInitError(err error) {
	m.ErrorOnInit = err
}

// SetUpdateError configures the mock to return an error on Update.
// afterN specifies which update call should fail (0-indexed).
// Use -1 to fail immediately on every update.
func (m *MockIntentWithErrors) SetUpdateError(err error, afterN int) {
	m.ErrorOnUpdate = err
	m.ErrorOnNthUpdate = afterN
}

// GetErrorCount returns the number of errors that have been triggered.
func (m *MockIntentWithErrors) GetErrorCount() int {
	return m.ErrorCount
}

// ClearErrors resets all error configuration and counts.
func (m *MockIntentWithErrors) ClearErrors() {
	m.ErrorOnInit = nil
	m.ErrorOnUpdate = nil
	m.ErrorOnNthUpdate = -1
	m.ErrorCount = 0
	m.updateCount = 0
}

// MockIntentWithContext extends MockIntent with context awareness for testing.
// It captures the context passed during activation for verification.
type MockIntentWithContext struct {
	*MockIntent
	// Context is the activation context passed to the factory
	Context map[string]interface{}
}

// NewMockIntentWithContext creates a new MockIntentWithContext for testing.
// It accepts the activation context for verification in tests.
func NewMockIntentWithContext(ctx map[string]interface{}) *MockIntentWithContext {
	return &MockIntentWithContext{
		MockIntent: NewMockIntent(),
		Context:    ctx,
	}
}

// GetContext returns the context passed during activation.
func (m *MockIntentWithContext) GetContext() map[string]interface{} {
	return m.Context
}

// GetContextValue returns a specific value from the activation context.
func (m *MockIntentWithContext) GetContextValue(key string) interface{} {
	if m.Context == nil {
		return nil
	}
	return m.Context[key]
}
