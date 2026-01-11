package intents

import (
	"fmt"
	"sync"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
)

// IntentFactory is a function that creates an Intent without context.
type IntentFactory = func() Intent

// IntentFactoryWithContext is a function that creates an Intent with activation context.
// The context allows passing data (like an event to edit) during intent activation.
type IntentFactoryWithContext = func(ctx map[string]interface{}) Intent

// DefaultIntentRouter implements the IntentRouter interface.
// It manages the lifecycle of intents and enforces the boundary contract.
type DefaultIntentRouter struct {
	// intents maps intent names to their factory functions (context-less).
	intents map[string]IntentFactory

	// intentsWithContext maps intent names to context-aware factory functions.
	intentsWithContext map[string]IntentFactoryWithContext

	// activeIntent is the currently active intent.
	activeIntent Intent

	// intentHistory tracks the history of activated intents for back navigation.
	intentHistory []Intent

	// mu protects concurrent access to state.
	mu sync.RWMutex

	// resultHandlers maps intent names to result handlers.
	// These are called when an intent completes.
	resultHandlers map[string]func(result *IntentResult[interface{}]) tea.Cmd

	// terminalInfo holds the current terminal dimensions
	terminalInfo *terminal.Info

	// logo is the shared logo instance passed to all intents
	logo *components.ASCIILogo

	// themeManager manages the application's theme system
	themeManager *themes.ThemeManager
}

// NewDefaultIntentRouter creates a new intent router.
func NewDefaultIntentRouter() *DefaultIntentRouter {
	return &DefaultIntentRouter{
		intents:            make(map[string]IntentFactory),
		intentsWithContext: make(map[string]IntentFactoryWithContext),
		intentHistory:      make([]Intent, 0),
		resultHandlers:     make(map[string]func(result *IntentResult[interface{}]) tea.Cmd),
		terminalInfo:       terminal.NewInfo(),
		themeManager:       themes.NewThemeManager(),
	}
}

// RegisterIntent registers an intent factory with the router.
// The factory function is called each time the intent is activated.
// For intents that need activation context (like edit mode), use RegisterIntentWithContext.
func (r *DefaultIntentRouter) RegisterIntent(name string, factory IntentFactory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.intents[name]; exists {
		return fmt.Errorf("intent %q already registered", name)
	}
	if _, exists := r.intentsWithContext[name]; exists {
		return fmt.Errorf("intent %q already registered (with context)", name)
	}

	r.intents[name] = factory
	return nil
}

// RegisterIntentWithContext registers a context-aware intent factory with the router.
// The factory receives activation context (e.g., event to edit) when the intent is activated.
// This enables patterns like edit mode where the same intent handles both create and edit.
func (r *DefaultIntentRouter) RegisterIntentWithContext(name string, factory IntentFactoryWithContext) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.intentsWithContext[name]; exists {
		return fmt.Errorf("intent %q already registered", name)
	}
	if _, exists := r.intents[name]; exists {
		return fmt.Errorf("intent %q already registered (without context)", name)
	}

	r.intentsWithContext[name] = factory
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
// The context parameter is passed to context-aware factories (registered via RegisterIntentWithContext).
// For context-less factories, the context is ignored.
func (r *DefaultIntentRouter) ActivateIntent(name string, context map[string]interface{}) (tea.Cmd, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for context-aware factory first
	var intent Intent
	if factoryWithCtx, exists := r.intentsWithContext[name]; exists {
		intent = factoryWithCtx(context)
	} else if factory, exists := r.intents[name]; exists {
		intent = factory()
	} else {
		return nil, fmt.Errorf("intent %q not found", name)
	}

	if intent == nil {
		return nil, fmt.Errorf("failed to create intent %q: factory returned nil", name)
	}

	// Push the current intent to history (if any).
	if r.activeIntent != nil {
		r.intentHistory = append(r.intentHistory, r.activeIntent)
	}

	// Set the new active intent
	r.activeIntent = intent

	// Propagate logo to the new intent if it has BaseIntent
	if r.logo != nil {
		if setter, ok := intent.(interface{ SetLogo(*components.ASCIILogo) }); ok {
			setter.SetLogo(r.logo)
		}
	}

	// Propagate theme manager to the new intent if it's theme-aware
	if r.themeManager != nil {
		if setter, ok := intent.(interface{ SetThemeManager(*themes.ThemeManager) }); ok {
			setter.SetThemeManager(r.themeManager)
		}
	}

	// Propagate terminal info to the new intent if it's terminal-aware
	// This must happen BEFORE Init() is called so the intent has terminal info during initialization
	if termAware, ok := intent.(TerminalAwareIntent); ok && r.terminalInfo.IsValid {
		termAware.UpdateTerminalInfo(r.terminalInfo)
	}

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
	// Handle WindowSizeMsg to update terminal info
	if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
		r.mu.Lock()
		r.terminalInfo.Update(wsMsg)
		// Propagate to active intent if terminal-aware
		if r.activeIntent != nil {
			if termAware, ok := r.activeIntent.(TerminalAwareIntent); ok {
				termAware.UpdateTerminalInfo(r.terminalInfo)
			}
		}
		r.mu.Unlock()
	}

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

// UpdateTerminalInfo updates the router's terminal information
// This should be called by the root app when it receives WindowSizeMsg
func (r *DefaultIntentRouter) UpdateTerminalInfo(info *terminal.Info) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.terminalInfo = info

	// Propagate to active intent if terminal-aware
	if r.activeIntent != nil {
		if termAware, ok := r.activeIntent.(TerminalAwareIntent); ok {
			termAware.UpdateTerminalInfo(info)
		}
	}
}

// GetTerminalInfo returns the current terminal information
func (r *DefaultIntentRouter) GetTerminalInfo() *terminal.Info {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.terminalInfo
}

// SetLogo sets the shared logo instance for all intents
func (r *DefaultIntentRouter) SetLogo(logo *components.ASCIILogo) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logo = logo

	// Propagate to active intent if it has BaseIntent
	if r.activeIntent != nil {
		if setter, ok := r.activeIntent.(interface{ SetLogo(*components.ASCIILogo) }); ok {
			setter.SetLogo(logo)
		}
	}
}

// GetLogo returns the shared logo instance
func (r *DefaultIntentRouter) GetLogo() *components.ASCIILogo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.logo
}

// SetThemeManager sets the theme manager for the router
func (r *DefaultIntentRouter) SetThemeManager(tm *themes.ThemeManager) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.themeManager = tm

	// Propagate to active intent if it's theme-aware
	if r.activeIntent != nil {
		if setter, ok := r.activeIntent.(interface{ SetThemeManager(*themes.ThemeManager) }); ok {
			setter.SetThemeManager(tm)
		}
	}
}

// GetThemeManager returns the theme manager
func (r *DefaultIntentRouter) GetThemeManager() *themes.ThemeManager {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.themeManager
}

// Theme returns the currently active theme for convenience
func (r *DefaultIntentRouter) Theme() themes.Theme {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.themeManager == nil {
		return nil
	}
	return r.themeManager.Active()
}
