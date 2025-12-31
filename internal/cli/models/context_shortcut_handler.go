package models

import (
	"sync"

	"github.com/charmbracelet/bubbles/key"
)

// ContextShortcutHandler manages context-aware shortcut handling with a stack-based context system
type ContextShortcutHandler struct {
	mu              sync.RWMutex
	mapper          *ShortcutMapper
	contextStack    []string
	currentContext  string
	contextMetadata map[string]map[string]interface{}
}

// NewContextShortcutHandler creates a new context-aware shortcut handler
func NewContextShortcutHandler() *ContextShortcutHandler {
	return &ContextShortcutHandler{
		mapper:          NewShortcutMapper(),
		contextStack:    make([]string, 0),
		currentContext:  "",
		contextMetadata: make(map[string]map[string]interface{}),
	}
}

// SetMapper sets the shortcut mapper for the handler
func (csh *ContextShortcutHandler) SetMapper(mapper *ShortcutMapper) {
	csh.mu.Lock()
	defer csh.mu.Unlock()
	csh.mapper = mapper
}

// SetCurrentContext sets the current context
func (csh *ContextShortcutHandler) SetCurrentContext(context string) {
	csh.mu.Lock()
	defer csh.mu.Unlock()
	csh.currentContext = context
}

// GetCurrentContext returns the current context
func (csh *ContextShortcutHandler) GetCurrentContext() string {
	csh.mu.RLock()
	defer csh.mu.RUnlock()
	return csh.currentContext
}

// PushContext pushes a context onto the stack and makes it current
func (csh *ContextShortcutHandler) PushContext(context string) {
	csh.mu.Lock()
	defer csh.mu.Unlock()
	csh.contextStack = append(csh.contextStack, csh.currentContext)
	csh.currentContext = context
}

// PopContext pops the current context and returns to the previous one
func (csh *ContextShortcutHandler) PopContext() string {
	csh.mu.Lock()
	defer csh.mu.Unlock()

	if len(csh.contextStack) == 0 {
		popped := csh.currentContext
		csh.currentContext = ""
		return popped
	}

	popped := csh.currentContext
	csh.currentContext = csh.contextStack[len(csh.contextStack)-1]
	csh.contextStack = csh.contextStack[:len(csh.contextStack)-1]
	return popped
}

// GetContextStack returns a copy of the context stack
func (csh *ContextShortcutHandler) GetContextStack() []string {
	csh.mu.RLock()
	defer csh.mu.RUnlock()

	stack := make([]string, len(csh.contextStack))
	copy(stack, csh.contextStack)
	return stack
}

// GetActiveShortcuts returns all active shortcuts for the current context
// This includes both global and context-specific shortcuts
func (csh *ContextShortcutHandler) GetActiveShortcuts() map[string]key.Binding {
	csh.mu.RLock()
	defer csh.mu.RUnlock()

	if csh.mapper == nil {
		return make(map[string]key.Binding)
	}

	return csh.mapper.GetMergedShortcuts(csh.currentContext)
}

// LookupShortcut finds a shortcut in the current context or falls back to global
func (csh *ContextShortcutHandler) LookupShortcut(id string) (key.Binding, bool) {
	csh.mu.RLock()
	defer csh.mu.RUnlock()

	if csh.mapper == nil {
		return key.Binding{}, false
	}

	// Try context-specific shortcut first
	if binding, exists := csh.mapper.GetContextShortcut(csh.currentContext, id); exists {
		return binding, true
	}

	// Fall back to global shortcut
	return csh.mapper.GetGlobalShortcut(id)
}

// SetContextMetadata sets metadata for a context
func (csh *ContextShortcutHandler) SetContextMetadata(context string, key string, value interface{}) {
	csh.mu.Lock()
	defer csh.mu.Unlock()

	if csh.contextMetadata[context] == nil {
		csh.contextMetadata[context] = make(map[string]interface{})
	}
	csh.contextMetadata[context][key] = value
}

// GetContextMetadata retrieves metadata for a context
func (csh *ContextShortcutHandler) GetContextMetadata(context, key string) (interface{}, bool) {
	csh.mu.RLock()
	defer csh.mu.RUnlock()

	if metadata, exists := csh.contextMetadata[context]; exists {
		value, found := metadata[key]
		return value, found
	}
	return nil, false
}

// ClearContextMetadata clears all metadata for a context
func (csh *ContextShortcutHandler) ClearContextMetadata(context string) {
	csh.mu.Lock()
	defer csh.mu.Unlock()

	delete(csh.contextMetadata, context)
}

// IsInContext checks if a specific context is active (in stack or current)
func (csh *ContextShortcutHandler) IsInContext(context string) bool {
	csh.mu.RLock()
	defer csh.mu.RUnlock()

	if csh.currentContext == context {
		return true
	}

	for _, ctx := range csh.contextStack {
		if ctx == context {
			return true
		}
	}

	return false
}

// GetContextDepth returns the depth of the context stack
func (csh *ContextShortcutHandler) GetContextDepth() int {
	csh.mu.RLock()
	defer csh.mu.RUnlock()

	// Include current context in depth
	if csh.currentContext == "" {
		return len(csh.contextStack)
	}
	return len(csh.contextStack) + 1
}

// ResetContextStack clears the context stack and resets to root
func (csh *ContextShortcutHandler) ResetContextStack() {
	csh.mu.Lock()
	defer csh.mu.Unlock()

	csh.contextStack = make([]string, 0)
	csh.currentContext = ""
}
