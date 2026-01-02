package intents

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// IntentTestHarness provides utilities for testing intents.
// It ensures:
// - Intents are tested in isolation
// - Intent results are strongly typed
// - Back navigation works correctly
// - State transitions are valid
// - No global state is mutated
type IntentTestHarness struct {
	intent Intent
	result interface{}
	t      *testing.T
}

// NewIntentTestHarness creates a new test harness for an intent.
func NewIntentTestHarness(t *testing.T, intent Intent) *IntentTestHarness {
	return &IntentTestHarness{
		intent: intent,
		t:      t,
	}
}

// Init initializes the intent and captures any startup commands.
func (h *IntentTestHarness) Init() tea.Cmd {
	return h.intent.Init()
}

// SendMessage sends a message to the intent and returns the command.
func (h *IntentTestHarness) SendMessage(msg tea.Msg) tea.Cmd {
	return h.intent.Update(msg)
}

// GetView returns the current view of the intent.
func (h *IntentTestHarness) GetView() string {
	return h.intent.View()
}

// GetResult returns the current result of the intent.
func (h *IntentTestHarness) GetResult() *IntentResult[interface{}] {
	return h.intent.Result()
}

// AssertResultCompleted asserts that the intent has completed.
func (h *IntentTestHarness) AssertResultCompleted() {
	result := h.intent.Result()
	if result == nil || result.Status != Completed {
		h.t.Errorf("expected intent result to be Completed, got %v", result)
	}
}

// AssertResultCancelled asserts that the intent was cancelled.
func (h *IntentTestHarness) AssertResultCancelled() {
	result := h.intent.Result()
	if result == nil || result.Status != Cancelled {
		h.t.Errorf("expected intent result to be Cancelled, got %v", result)
	}
}

// AssertResultFailed asserts that the intent failed.
func (h *IntentTestHarness) AssertResultFailed() {
	result := h.intent.Result()
	if result == nil || result.Status != Failed {
		h.t.Errorf("expected intent result to be Failed, got %v", result)
	}
}

// AssertViewContains asserts that the view contains a substring.
func (h *IntentTestHarness) AssertViewContains(substring string) {
	view := h.intent.View()
	if !contains(view, substring) {
		h.t.Errorf("expected view to contain %q, got: %s", substring, view)
	}
}

// AssertViewNotContains asserts that the view does not contain a substring.
func (h *IntentTestHarness) AssertViewNotContains(substring string) {
	view := h.intent.View()
	if contains(view, substring) {
		h.t.Errorf("expected view to not contain %q, got: %s", substring, view)
	}
}

// contains is a helper function to check if a string contains a substring.
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

// NewIntentRouterTestHelper creates a new test helper for the intent router.
func NewIntentRouterTestHelper(t *testing.T, router IntentRouter) *IntentRouterTestHelper {
	return &IntentRouterTestHelper{
		router: router,
		t:      t,
	}
}

// ActivateIntent activates an intent by name.
func (h *IntentRouterTestHelper) ActivateIntent(name string, context map[string]interface{}) tea.Cmd {
	cmd, err := h.router.ActivateIntent(name, context)
	if err != nil {
		h.t.Fatalf("failed to activate intent: %v", err)
	}
	return cmd
}

// GetActiveIntent returns the currently active intent.
func (h *IntentRouterTestHelper) GetActiveIntent() Intent {
	return h.router.GetActiveIntent()
}

// AssertIntentActive asserts that an intent is active.
func (h *IntentRouterTestHelper) AssertIntentActive() {
	if h.router.GetActiveIntent() == nil {
		h.t.Errorf("expected an active intent, got nil")
	}
}

// AssertIntentInactive asserts that no intent is active.
func (h *IntentRouterTestHelper) AssertIntentInactive() {
	if h.router.GetActiveIntent() != nil {
		h.t.Errorf("expected no active intent")
	}
}

// GoBack navigates back to the previous intent.
func (h *IntentRouterTestHelper) GoBack() {
	_, err := h.router.Back()
	if err != nil {
		h.t.Fatalf("failed to go back: %v", err)
	}
}

// AssertGoBackFails asserts that going back fails.
func (h *IntentRouterTestHelper) AssertGoBackFails() {
	_, err := h.router.Back()
	if err == nil {
		h.t.Errorf("expected Back to fail, but it succeeded")
	}
}

// GetHistory returns the intent history.
func (h *IntentRouterTestHelper) GetHistory() []Intent {
	return h.router.GetHistory()
}

// AssertHistoryLength asserts that the history has the expected length.
func (h *IntentRouterTestHelper) AssertHistoryLength(expected int) {
	history := h.router.GetHistory()
	if len(history) != expected {
		h.t.Errorf("expected history length %d, got %d", expected, len(history))
	}
}
