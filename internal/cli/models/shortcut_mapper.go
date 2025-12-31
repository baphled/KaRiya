package models

import (
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// ShortcutMapper manages global and context-specific keyboard shortcuts
type ShortcutMapper struct {
	mu                sync.RWMutex
	globalShortcuts   map[string]key.Binding
	globalActions     map[string]ShortcutAction
	contextShortcuts  map[string]map[string]key.Binding
	contextActions    map[string]map[string]ShortcutAction
	conflictResolvers map[string]string // Maps key string to shortcut ID for conflict detection
}

// NewShortcutMapper creates a new global shortcut mapper instance
func NewShortcutMapper() *ShortcutMapper {
	return &ShortcutMapper{
		globalShortcuts:   make(map[string]key.Binding),
		globalActions:     make(map[string]ShortcutAction),
		contextShortcuts:  make(map[string]map[string]key.Binding),
		contextActions:    make(map[string]map[string]ShortcutAction),
		conflictResolvers: make(map[string]string),
	}
}

// RegisterGlobalShortcut registers a shortcut that works across all contexts
func (sm *ShortcutMapper) RegisterGlobalShortcut(id string, binding key.Binding, action ShortcutAction) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.globalShortcuts[id] = binding
	sm.globalActions[id] = action

	// Track for conflict detection
	keyStr := getKeyString(binding)
	sm.conflictResolvers[keyStr] = id
}

// UnregisterGlobalShortcut removes a global shortcut
func (sm *ShortcutMapper) UnregisterGlobalShortcut(id string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	binding, exists := sm.globalShortcuts[id]
	if exists {
		delete(sm.globalShortcuts, id)
		delete(sm.globalActions, id)

		keyStr := getKeyString(binding)
		delete(sm.conflictResolvers, keyStr)
	}
}

// GetGlobalShortcut retrieves a global shortcut by ID
func (sm *ShortcutMapper) GetGlobalShortcut(id string) (key.Binding, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	binding, exists := sm.globalShortcuts[id]
	return binding, exists
}

// GetAllGlobalShortcuts returns all registered global shortcuts
func (sm *ShortcutMapper) GetAllGlobalShortcuts() map[string]key.Binding {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	result := make(map[string]key.Binding)
	for id, binding := range sm.globalShortcuts {
		result[id] = binding
	}
	return result
}

// RegisterContextShortcut registers a shortcut that only works in a specific context
func (sm *ShortcutMapper) RegisterContextShortcut(context, id string, binding key.Binding, action ShortcutAction) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.contextShortcuts[context] == nil {
		sm.contextShortcuts[context] = make(map[string]key.Binding)
		sm.contextActions[context] = make(map[string]ShortcutAction)
	}

	sm.contextShortcuts[context][id] = binding
	sm.contextActions[context][id] = action
}

// UnregisterContextShortcut removes a context-specific shortcut
func (sm *ShortcutMapper) UnregisterContextShortcut(context, id string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if contextMap, exists := sm.contextShortcuts[context]; exists {
		delete(contextMap, id)
		if len(contextMap) == 0 {
			delete(sm.contextShortcuts, context)
			delete(sm.contextActions, context)
		}
	}
	if contextActions, exists := sm.contextActions[context]; exists {
		delete(contextActions, id)
	}
}

// GetContextShortcut retrieves a context-specific shortcut
func (sm *ShortcutMapper) GetContextShortcut(context, id string) (key.Binding, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if contextMap, exists := sm.contextShortcuts[context]; exists {
		binding, found := contextMap[id]
		return binding, found
	}
	return key.Binding{}, false
}

// GetContextShortcuts returns all shortcuts for a specific context
func (sm *ShortcutMapper) GetContextShortcuts(context string) map[string]key.Binding {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	result := make(map[string]key.Binding)
	if contextMap, exists := sm.contextShortcuts[context]; exists {
		for id, binding := range contextMap {
			result[id] = binding
		}
	}
	return result
}

// CheckConflicts detects conflicting key bindings with existing shortcuts
func (sm *ShortcutMapper) CheckConflicts(id string, binding key.Binding) []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var conflicts []string
	keyStr := getKeyString(binding)

	// Check against existing shortcuts
	for existingID, existingBinding := range sm.globalShortcuts {
		if existingID != id && getKeyString(existingBinding) == keyStr {
			conflicts = append(conflicts, existingID)
		}
	}

	// Check context shortcuts
	for context, contextMap := range sm.contextShortcuts {
		for contextID, contextBinding := range contextMap {
			if getKeyString(contextBinding) == keyStr {
				conflicts = append(conflicts, fmt.Sprintf("%s:%s", context, contextID))
			}
		}
	}

	return conflicts
}

// GetMergedShortcuts returns all shortcuts for a context, with global ones merged in
// Context shortcuts take priority over global shortcuts
func (sm *ShortcutMapper) GetMergedShortcuts(context string) map[string]key.Binding {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	result := make(map[string]key.Binding)

	// First add all global shortcuts
	for id, binding := range sm.globalShortcuts {
		result[id] = binding
	}

	// Then override with context-specific shortcuts
	if contextMap, exists := sm.contextShortcuts[context]; exists {
		for id, binding := range contextMap {
			result[id] = binding
		}
	}

	return result
}

// InvokeShortcut invokes the action associated with a global shortcut
func (sm *ShortcutMapper) InvokeShortcut(id string) (tea.Cmd, bool) {
	sm.mu.RLock()
	action, exists := sm.globalActions[id]
	sm.mu.RUnlock()

	if !exists {
		return nil, false
	}

	return action(), true
}

// InvokeContextShortcut invokes the action associated with a context-specific shortcut
func (sm *ShortcutMapper) InvokeContextShortcut(context, id string) (tea.Cmd, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if contextActions, exists := sm.contextActions[context]; exists {
		if action, found := contextActions[id]; found {
			return action(), true
		}
	}

	return nil, false
}

// ClearGlobalShortcuts removes all global shortcuts
func (sm *ShortcutMapper) ClearGlobalShortcuts() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.globalShortcuts = make(map[string]key.Binding)
	sm.globalActions = make(map[string]ShortcutAction)

	// Clear conflict resolvers for global shortcuts
	for keyStr := range sm.conflictResolvers {
		delete(sm.conflictResolvers, keyStr)
	}
}

// ClearContextShortcuts removes all shortcuts for a specific context
func (sm *ShortcutMapper) ClearContextShortcuts(context string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.contextShortcuts, context)
	delete(sm.contextActions, context)
}

// Clear removes all shortcuts (both global and context-specific)
func (sm *ShortcutMapper) Clear() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.globalShortcuts = make(map[string]key.Binding)
	sm.globalActions = make(map[string]ShortcutAction)
	sm.contextShortcuts = make(map[string]map[string]key.Binding)
	sm.contextActions = make(map[string]map[string]ShortcutAction)
	sm.conflictResolvers = make(map[string]string)
}

// getKeyString extracts a string representation of a key binding for comparison
func getKeyString(binding key.Binding) string {
	help := binding.Help()
	return strings.ToLower(help.Key)
}
