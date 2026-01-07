package intents

import (
	"fmt"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

// DefaultIntentRouter implements the IntentRouter interface.
// It manages the lifecycle of intents and enforces the boundary contract.
type DefaultIntentRouter struct {
	// intents maps intent names to their factory functions.
	intents map[string]func() Intent

	// activeIntent is the currently active intent.
	activeIntent Intent

	// intentHistory tracks the history of activated intents for back navigation.
	intentHistory []Intent

	// mu protects concurrent access to state.
	mu sync.RWMutex

	// resultHandlers maps intent names to result handlers.
	// These are called when an intent completes.
	resultHandlers map[string]func(result *IntentResult[interface{}]) tea.Cmd
}

// NewDefaultIntentRouter creates a new intent router.
func NewDefaultIntentRouter() *DefaultIntentRouter {
	return &DefaultIntentRouter{
		intents:        make(map[string]func() Intent),
		intentHistory:  make([]Intent, 0),
		resultHandlers: make(map[string]func(result *IntentResult[interface{}]) tea.Cmd),
	}
}

// RegisterIntent registers an intent factory with the router.
// The factory function is called each time the intent is activated.
func (r *DefaultIntentRouter) RegisterIntent(name string, factory func() Intent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.intents[name]; exists {
		return fmt.Errorf("intent %q already registered", name)
	}

	r.intents[name] = factory
	return nil
}

// RegisterResultHandler registers a handler called when an intent completes.
// The handler returns a command to propagate the result to the root model.
func (r *DefaultIntentRouter) RegisterResultHandler(intentName string, handler func(result *IntentResult[interface{}]) tea.Cmd) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.resultHandlers[intentName] = handler
}

// ActivateIntent activates an intent by name.
// This is the ONLY way intents are activated, enforcing strict transition rules.
func (r *DefaultIntentRouter) ActivateIntent(name string, context map[string]interface{}) (tea.Cmd, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	factory, exists := r.intents[name]
	if !exists {
		return nil, fmt.Errorf("intent %q not found", name)
	}

	// Push the current intent to history (if any).
	if r.activeIntent != nil {
		r.intentHistory = append(r.intentHistory, r.activeIntent)
	}

	// Create and activate the new intent.
	intent := factory()
	if intent == nil {
		return nil, fmt.Errorf("failed to create intent %q: factory returned nil", name)
	}
	r.activeIntent = intent

	// Call the intent's Init method to get any startup commands.
	return intent.Init(), nil
}

// GetActiveIntent returns the currently active intent.
func (r *DefaultIntentRouter) GetActiveIntent() Intent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.activeIntent
}

// HandleMessage processes a message in the active intent.
// Returns a command and any intent result if the intent completed.
func (r *DefaultIntentRouter) HandleMessage(msg tea.Msg) (tea.Cmd, interface{}) {
	r.mu.RLock()
	intent := r.activeIntent
	r.mu.RUnlock()

	if intent == nil {
		return nil, nil
	}

	// Delegate to the active intent's Update method.
	cmd := intent.Update(msg)

	// Check if the intent has completed.
	result := intent.Result()
	if result != nil {
		return cmd, result
	}

	return cmd, nil
}

// View renders the active intent.
func (r *DefaultIntentRouter) View() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.activeIntent == nil {
		return "No active intent"
	}

	return r.activeIntent.View()
}

// Back navigates back to the previous intent in the history.
// Returns an error if there's no previous intent.
func (r *DefaultIntentRouter) Back() (tea.Cmd, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.intentHistory) == 0 {
		return nil, fmt.Errorf("no previous intent to go back to")
	}

	// Pop the previous intent from history.
	previousIntent := r.intentHistory[len(r.intentHistory)-1]
	r.intentHistory = r.intentHistory[:len(r.intentHistory)-1]

	// Activate the previous intent.
	r.activeIntent = previousIntent

	// Call the intent's Init method to restore state.
	return previousIntent.Init(), nil
}

// GetHistory returns a copy of the intent history.
func (r *DefaultIntentRouter) GetHistory() []Intent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	history := make([]Intent, len(r.intentHistory))
	copy(history, r.intentHistory)
	return history
}

// GetHistoryDepth returns the current depth in the navigation history.
func (r *DefaultIntentRouter) GetHistoryDepth() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.activeIntent == nil {
		return 0
	}
	return len(r.intentHistory) + 1
}
