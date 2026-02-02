package intents

import (
	"errors"
	"fmt"
	"sync"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
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

	// terminalInfo holds the current terminal dimensions
	terminalInfo *terminal.Info

	// logo is the shared logo instance passed to all intents
	logo LogoModel

	// themeManager manages the application's theme system
	themeManager *themes.ThemeManager
}

// NewDefaultIntentRouter creates a new intent router with empty registries and default terminal info.
//
// Returns: a fully initialized router ready for intent registration.
//
// Side effects: allocates a new ThemeManager and terminal Info.
func NewDefaultIntentRouter() *DefaultIntentRouter {
	return &DefaultIntentRouter{
		intents:        make(map[string]func() Intent),
		intentHistory:  make([]Intent, 0),
		resultHandlers: make(map[string]func(result *IntentResult[interface{}]) tea.Cmd),
		terminalInfo:   terminal.NewInfo(),
		themeManager:   themes.NewThemeManager(),
	}
}

// RegisterIntent registers an intent factory with the router so it can be activated by name.
// The factory function is called each time the intent is activated, producing a fresh instance.
//
// Expected: name must be unique across all registered intents; factory must not be nil.
//
// Returns: an error if an intent with the same name is already registered.
//
// Side effects: stores the factory in the router's intent registry under a write lock.
func (r *DefaultIntentRouter) RegisterIntent(name string, factory func() Intent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.intents[name]; exists {
		return fmt.Errorf("intent %q already registered", name)
	}

	r.intents[name] = factory
	return nil
}

// RegisterResultHandler registers a callback invoked when the named intent completes,
// allowing the root model to react to intent results.
//
// Expected: intentName must correspond to a registered intent; handler must not be nil.
//
// Side effects: stores the handler in the result handlers map under a write lock;
// overwrites any previously registered handler for the same intent name.
func (r *DefaultIntentRouter) RegisterResultHandler(intentName string, handler func(result *IntentResult[interface{}]) tea.Cmd) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.resultHandlers[intentName] = handler
}

// ActivateIntent creates and activates an intent by name, enforcing strict transition rules.
// This is the sole entry point for intent activation, ensuring consistent lifecycle management.
//
// Expected: name must match a previously registered intent; the factory must produce a non-nil intent.
//
// Returns: the startup command from the intent's Init method, or an error if the intent
// is not found or the factory returns nil.
//
// Side effects: pushes the current active intent onto the history stack, creates a new intent
// via its factory, propagates logo/theme/terminal info, and calls Init on the new intent.
func (r *DefaultIntentRouter) ActivateIntent(name string, _ map[string]interface{}) (tea.Cmd, error) {
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

	// Propagate logo to the new intent if it has BaseIntent
	if r.logo != nil {
		if setter, ok := intent.(interface{ SetLogo(LogoModel) }); ok {
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

// GetActiveIntent provides access to the currently active intent for external inspection or testing.
//
// Returns: the active intent, or nil if no intent has been activated.
//
// Side effects: acquires a read lock for thread-safe access.
func (r *DefaultIntentRouter) GetActiveIntent() Intent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.activeIntent
}

// HandleMessage delegates a Bubble Tea message to the active intent and checks for completion.
// WindowSizeMsg messages are intercepted to keep terminal dimensions in sync across the router.
//
// Expected: msg must be a valid tea.Msg; an active intent should be set for meaningful processing.
//
// Returns: a command from the intent's Update, and a non-nil result if the intent has completed.
//
// Side effects: updates terminal info on WindowSizeMsg; propagates terminal info to the active
// intent if it implements TerminalAwareIntent; calls Update and Result on the active intent.
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

// View renders the active intent's UI for display in the terminal.
//
// Returns: the rendered string from the active intent, or a placeholder if no intent is active.
//
// Side effects: acquires a read lock for thread-safe access.
func (r *DefaultIntentRouter) View() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.activeIntent == nil {
		return "No active intent"
	}

	return r.activeIntent.View()
}

// Back navigates to the previous intent by popping the history stack,
// enabling breadcrumb-style back navigation between intents.
//
// Returns: the startup command from the restored intent's Init, or an error if the history is empty.
//
// Side effects: pops the most recent intent from the history stack, sets it as active,
// and calls Init to restore its state.
func (r *DefaultIntentRouter) Back() (tea.Cmd, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.intentHistory) == 0 {
		return nil, errors.New("no previous intent to go back to")
	}

	// Pop the previous intent from history.
	previousIntent := r.intentHistory[len(r.intentHistory)-1]
	r.intentHistory = r.intentHistory[:len(r.intentHistory)-1]

	// Activate the previous intent.
	r.activeIntent = previousIntent

	// Call the intent's Init method to restore state.
	return previousIntent.Init(), nil
}

// GetHistory provides a snapshot of the navigation history for inspection or testing.
//
// Returns: a defensive copy of the intent history slice, safe for external mutation.
//
// Side effects: acquires a read lock for thread-safe access.
func (r *DefaultIntentRouter) GetHistory() []Intent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	history := make([]Intent, len(r.intentHistory))
	copy(history, r.intentHistory)
	return history
}

// GetHistoryDepth reports how many intents deep the user is in the navigation stack,
// useful for displaying breadcrumb depth or determining if back navigation is available.
//
// Returns: the total number of intents in the stack (history plus active), or 0 if no intent is active.
//
// Side effects: acquires a read lock for thread-safe access.
func (r *DefaultIntentRouter) GetHistoryDepth() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.activeIntent == nil {
		return 0
	}
	return len(r.intentHistory) + 1
}

// UpdateTerminalInfo synchronizes the router's terminal dimensions with the latest window state.
// The root app should call this when it receives a WindowSizeMsg.
//
// Expected: info must be a valid, non-nil terminal.Info with current dimensions.
//
// Side effects: replaces the stored terminal info and propagates the update to the active
// intent if it implements TerminalAwareIntent.
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

// GetTerminalInfo provides access to the cached terminal dimensions for layout calculations.
//
// Returns: the current terminal info held by the router.
//
// Side effects: acquires a read lock for thread-safe access.
func (r *DefaultIntentRouter) GetTerminalInfo() *terminal.Info {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.terminalInfo
}

// SetLogo configures the shared logo model used by all intents for consistent branding.
//
// Expected: logo should be a valid LogoModel; nil disables logo propagation.
//
// Side effects: stores the logo and propagates it to the active intent if it supports SetLogo.
func (r *DefaultIntentRouter) SetLogo(logo LogoModel) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logo = logo

	// Propagate to active intent if it has BaseIntent
	if r.activeIntent != nil {
		if setter, ok := r.activeIntent.(interface{ SetLogo(LogoModel) }); ok {
			setter.SetLogo(logo)
		}
	}
}

// GetLogo provides access to the shared logo model for external rendering or inspection.
//
// Returns: the currently configured LogoModel, or nil if none has been set.
//
// Side effects: acquires a read lock for thread-safe access.
func (r *DefaultIntentRouter) GetLogo() LogoModel {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.logo
}

// SetThemeManager configures the theme system used by the router and all managed intents.
//
// Expected: tm should be a valid, initialized ThemeManager; nil disables theme propagation.
//
// Side effects: stores the theme manager and propagates it to the active intent if it
// supports SetThemeManager.
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

// GetThemeManager provides access to the router's theme manager for theme queries or configuration.
//
// Returns: the current ThemeManager, or nil if none has been set.
//
// Side effects: acquires a read lock for thread-safe access.
func (r *DefaultIntentRouter) GetThemeManager() *themes.ThemeManager {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.themeManager
}

// Theme provides a convenience accessor for the active theme without requiring direct
// ThemeManager interaction.
//
// Returns: the currently active theme, or nil if no ThemeManager has been configured.
//
// Side effects: acquires a read lock for thread-safe access.
func (r *DefaultIntentRouter) Theme() themes.Theme {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.themeManager == nil {
		return nil
	}
	return r.themeManager.Active()
}
