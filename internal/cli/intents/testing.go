package intents

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// IntentTestHarness provides utilities for testing intents.
type IntentTestHarness struct {
	intent Intent
	t      *testing.T
}

// NewIntentTestHarness creates a new test harness for an intent.
func NewIntentTestHarness(t *testing.T, intent Intent) *IntentTestHarness {
	return &IntentTestHarness{
		intent: intent,
		t:      t,
	}
}

// Init initializes the intent.
func (h *IntentTestHarness) Init() tea.Cmd {
	return h.intent.Init()
}

// SendMessage sends a message to the intent.
func (h *IntentTestHarness) SendMessage(msg tea.Msg) tea.Cmd {
	return h.intent.Update(msg)
}

// GetView returns the current view.
func (h *IntentTestHarness) GetView() string {
	return h.intent.View()
}

// GetResult returns the current result.
func (h *IntentTestHarness) GetResult() *IntentResult[interface{}] {
	return h.intent.Result()
}

// AssertResultCompleted asserts completion.
func (h *IntentTestHarness) AssertResultCompleted() {
	result := h.intent.Result()
	if result == nil || result.Status != Completed {
		h.t.Errorf("expected Completed, got %v", result)
	}
}

// AssertResultCancelled asserts cancellation.
func (h *IntentTestHarness) AssertResultCancelled() {
	result := h.intent.Result()
	if result == nil || result.Status != Cancelled {
		h.t.Errorf("expected Cancelled, got %v", result)
	}
}

// AssertResultFailed asserts failure.
func (h *IntentTestHarness) AssertResultFailed() {
	result := h.intent.Result()
	if result == nil || result.Status != Failed {
		h.t.Errorf("expected Failed, got %v", result)
	}
}

// AssertViewContains asserts view contains substring.
func (h *IntentTestHarness) AssertViewContains(substring string) {
	view := h.intent.View()
	if !contains(view, substring) {
		h.t.Errorf("expected view to contain %q", substring)
	}
}

// AssertViewNotContains asserts view doesn't contain substring.
func (h *IntentTestHarness) AssertViewNotContains(substring string) {
	view := h.intent.View()
	if contains(view, substring) {
		h.t.Errorf("expected view to not contain %q", substring)
	}
}

// contains helper function.
func contains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// IntentRouterTestHelper provides utilities for testing the intent router.
type IntentRouterTestHelper struct {
	router IntentRouter
	t      *testing.T
}

// NewIntentRouterTestHelper creates a new test helper.
func NewIntentRouterTestHelper(t *testing.T, router IntentRouter) *IntentRouterTestHelper {
	return &IntentRouterTestHelper{
		router: router,
		t:      t,
	}
}

// ActivateIntent activates an intent.
func (h *IntentRouterTestHelper) ActivateIntent(name string, context map[string]interface{}) tea.Cmd {
	cmd, err := h.router.ActivateIntent(name, context)
	if err != nil {
		h.t.Fatalf("failed to activate intent: %v", err)
	}
	return cmd
}

// GetActiveIntent returns the active intent.
func (h *IntentRouterTestHelper) GetActiveIntent() Intent {
	return h.router.GetActiveIntent()
}

// AssertIntentActive asserts intent is active.
func (h *IntentRouterTestHelper) AssertIntentActive() {
	if h.router.GetActiveIntent() == nil {
		h.t.Errorf("expected an active intent")
	}
}

// AssertIntentInactive asserts no intent is active.
func (h *IntentRouterTestHelper) AssertIntentInactive() {
	if h.router.GetActiveIntent() != nil {
		h.t.Errorf("expected no active intent")
	}
}

// GoBack navigates back.
func (h *IntentRouterTestHelper) GoBack() {
	_, err := h.router.Back()
	if err != nil {
		h.t.Fatalf("failed to go back: %v", err)
	}
}

// AssertGoBackFails asserts back navigation fails.
func (h *IntentRouterTestHelper) AssertGoBackFails() {
	_, err := h.router.Back()
	if err == nil {
		h.t.Errorf("expected Back to fail")
	}
}

// GetHistory returns navigation history.
func (h *IntentRouterTestHelper) GetHistory() []Intent {
	return h.router.GetHistory()
}

// AssertHistoryLength asserts history length.
func (h *IntentRouterTestHelper) AssertHistoryLength(expected int) {
	history := h.router.GetHistory()
	if len(history) != expected {
		h.t.Errorf("expected history length %d, got %d", expected, len(history))
	}
}

// TestIntentFactory creates test intents.
type TestIntentFactory struct {
	intents map[string]func() Intent
}

// NewTestIntentFactory creates a new factory.
func NewTestIntentFactory() *TestIntentFactory {
	return &TestIntentFactory{
		intents: make(map[string]func() Intent),
	}
}

// Register registers a factory function.
func (f *TestIntentFactory) Register(name string, factory func() Intent) {
	f.intents[name] = factory
}

// Create creates an intent.
func (f *TestIntentFactory) Create(name string) Intent {
	if factory, ok := f.intents[name]; ok {
		return factory()
	}
	return nil
}

// IntentWithState wraps an intent with state inspection.
type IntentWithState struct {
	intent Intent
	state  interface{}
}

// NewIntentWithState creates a wrapper.
func NewIntentWithState(intent Intent, state interface{}) *IntentWithState {
	return &IntentWithState{
		intent: intent,
		state:  state,
	}
}

// Init implements Intent.
func (i *IntentWithState) Init() tea.Cmd {
	return i.intent.Init()
}

// Update implements Intent.
func (i *IntentWithState) Update(msg tea.Msg) tea.Cmd {
	return i.intent.Update(msg)
}

// View implements Intent.
func (i *IntentWithState) View() string {
	return i.intent.View()
}

// Result implements Intent.
func (i *IntentWithState) Result() *IntentResult[interface{}] {
	return i.intent.Result()
}

// GetState returns the state.
func (i *IntentWithState) GetState() interface{} {
	return i.state
}

// SetState sets the state.
func (i *IntentWithState) SetState(state interface{}) {
	i.state = state
}
