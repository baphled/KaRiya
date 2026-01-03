package context

import (
	"sync"
)

// GlobalContext holds application-wide state that is accessible to all intents.
// It provides read-only access to user preferences and configuration, plus
// transient state that can be shared between intents.
//
// Thread-safe: All methods are protected by a mutex for concurrent access.
type GlobalContext struct {
	mu sync.RWMutex

	// preferences holds user preferences (read-only after initialization)
	preferences map[string]interface{}

	// transientState holds temporary state that can be shared between intents
	// (e.g., selected items, filters, scroll positions)
	transientState map[string]interface{}

	// config holds application configuration (read-only after initialization)
	config map[string]interface{}
}

// NewGlobalContext creates a new GlobalContext with optional initial values.
func NewGlobalContext(
	preferences map[string]interface{},
	config map[string]interface{},
) *GlobalContext {
	if preferences == nil {
		preferences = make(map[string]interface{})
	}
	if config == nil {
		config = make(map[string]interface{})
	}

	return &GlobalContext{
		preferences:    preferences,
		transientState: make(map[string]interface{}),
		config:         config,
	}
}

// GetPreference retrieves a user preference by key.
// Returns the value and a boolean indicating whether the key exists.
// Preferences are read-only and shared across all intents.
func (gc *GlobalContext) GetPreference(key string) (interface{}, bool) {
	gc.mu.RLock()
	defer gc.mu.RUnlock()

	val, exists := gc.preferences[key]
	return val, exists
}

// SetPreference sets a user preference by key.
// WARNING: This modifies the shared preference state. Use with caution.
// Typically, preferences should be set only during initialization.
func (gc *GlobalContext) SetPreference(key string, value interface{}) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	gc.preferences[key] = value
}

// GetAllPreferences returns a copy of all preferences.
// Safe for iteration without holding the lock.
func (gc *GlobalContext) GetAllPreferences() map[string]interface{} {
	gc.mu.RLock()
	defer gc.mu.RUnlock()

	copy := make(map[string]interface{})
	for k, v := range gc.preferences {
		copy[k] = v
	}
	return copy
}

// SetTransientState sets transient state that can be shared between intents.
// This is useful for passing data between intents without modifying persistent state.
// Examples: selected items, filter state, scroll positions.
func (gc *GlobalContext) SetTransientState(key string, value interface{}) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	gc.transientState[key] = value
}

// GetTransientState retrieves transient state by key.
// Returns the value and a boolean indicating whether the key exists.
func (gc *GlobalContext) GetTransientState(key string) (interface{}, bool) {
	gc.mu.RLock()
	defer gc.mu.RUnlock()

	val, exists := gc.transientState[key]
	return val, exists
}

// ClearTransientState removes transient state by key.
// Useful for cleanup when exiting an intent.
func (gc *GlobalContext) ClearTransientState(key string) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	delete(gc.transientState, key)
}

// GetAllTransientState returns a copy of all transient state.
// Safe for iteration without holding the lock.
func (gc *GlobalContext) GetAllTransientState() map[string]interface{} {
	gc.mu.RLock()
	defer gc.mu.RUnlock()

	copy := make(map[string]interface{})
	for k, v := range gc.transientState {
		copy[k] = v
	}
	return copy
}

// ClearAllTransientState removes all transient state.
// Useful for cleanup when resetting the app state.
func (gc *GlobalContext) ClearAllTransientState() {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	gc.transientState = make(map[string]interface{})
}

// GetConfig retrieves application configuration by key.
// Returns the value and a boolean indicating whether the key exists.
// Configuration is read-only and shared across all intents.
func (gc *GlobalContext) GetConfig(key string) (interface{}, bool) {
	gc.mu.RLock()
	defer gc.mu.RUnlock()

	val, exists := gc.config[key]
	return val, exists
}

// SetConfig sets application configuration by key.
// WARNING: This modifies the shared configuration state. Use with caution.
// Typically, configuration should be set only during initialization.
func (gc *GlobalContext) SetConfig(key string, value interface{}) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	gc.config[key] = value
}

// GetAllConfig returns a copy of all configuration.
// Safe for iteration without holding the lock.
func (gc *GlobalContext) GetAllConfig() map[string]interface{} {
	gc.mu.RLock()
	defer gc.mu.RUnlock()

	copy := make(map[string]interface{})
	for k, v := range gc.config {
		copy[k] = v
	}
	return copy
}

// PreferenceKeys returns all preference keys.
// Useful for debugging and inspection.
func (gc *GlobalContext) PreferenceKeys() []string {
	gc.mu.RLock()
	defer gc.mu.RUnlock()

	keys := make([]string, 0, len(gc.preferences))
	for k := range gc.preferences {
		keys = append(keys, k)
	}
	return keys
}

// TransientStateKeys returns all transient state keys.
// Useful for debugging and inspection.
func (gc *GlobalContext) TransientStateKeys() []string {
	gc.mu.RLock()
	defer gc.mu.RUnlock()

	keys := make([]string, 0, len(gc.transientState))
	for k := range gc.transientState {
		keys = append(keys, k)
	}
	return keys
}

// ConfigKeys returns all configuration keys.
// Useful for debugging and inspection.
func (gc *GlobalContext) ConfigKeys() []string {
	gc.mu.RLock()
	defer gc.mu.RUnlock()

	keys := make([]string, 0, len(gc.config))
	for k := range gc.config {
		keys = append(keys, k)
	}
	return keys
}
