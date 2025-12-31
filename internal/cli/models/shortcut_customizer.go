package models

import (
	"errors"
	"sync"

	"github.com/charmbracelet/bubbles/key"
)

// ShortcutProfile represents a collection of custom shortcuts
type ShortcutProfile struct {
	Name      string
	Overrides map[string]key.Binding
	Disabled  map[string]bool
	mu        sync.RWMutex
}

// NewShortcutProfile creates a new shortcut profile
func NewShortcutProfile(name string) *ShortcutProfile {
	return &ShortcutProfile{
		Name:      name,
		Overrides: make(map[string]key.Binding),
		Disabled:  make(map[string]bool),
	}
}

// OverrideShortcut sets a custom binding for a shortcut in the profile
func (sp *ShortcutProfile) OverrideShortcut(id string, binding key.Binding) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.Overrides[id] = binding
}

// DisableShortcut marks a shortcut as disabled in the profile
func (sp *ShortcutProfile) DisableShortcut(id string) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.Disabled[id] = true
}

// EnableShortcut marks a shortcut as enabled in the profile
func (sp *ShortcutProfile) EnableShortcut(id string) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	delete(sp.Disabled, id)
}

// ShortcutCustomizer manages shortcut customization and profiles
type ShortcutCustomizer struct {
	mu               sync.RWMutex
	overrides        map[string]key.Binding
	disabled         map[string]bool
	profiles         map[string]*ShortcutProfile
	activeProfile    string
	defaultShortcuts map[string]key.Binding
}

// NewShortcutCustomizer creates a new shortcut customizer
func NewShortcutCustomizer() *ShortcutCustomizer {
	return &ShortcutCustomizer{
		overrides:        make(map[string]key.Binding),
		disabled:         make(map[string]bool),
		profiles:         make(map[string]*ShortcutProfile),
		defaultShortcuts: make(map[string]key.Binding),
	}
}

// OverrideShortcut sets a custom binding for a shortcut
func (sc *ShortcutCustomizer) OverrideShortcut(id string, binding key.Binding) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.overrides[id] = binding
}

// GetOverride retrieves a custom binding for a shortcut
func (sc *ShortcutCustomizer) GetOverride(id string) (key.Binding, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	binding, exists := sc.overrides[id]
	return binding, exists
}

// GetAllOverrides returns all overridden shortcuts
func (sc *ShortcutCustomizer) GetAllOverrides() map[string]key.Binding {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	result := make(map[string]key.Binding)
	for id, binding := range sc.overrides {
		result[id] = binding
	}
	return result
}

// DisableShortcut marks a shortcut as disabled
func (sc *ShortcutCustomizer) DisableShortcut(id string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.disabled[id] = true
}

// EnableShortcut marks a shortcut as enabled
func (sc *ShortcutCustomizer) EnableShortcut(id string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	delete(sc.disabled, id)
}

// IsShortcutDisabled checks if a shortcut is disabled
func (sc *ShortcutCustomizer) IsShortcutDisabled(id string) bool {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.disabled[id]
}

// GetDisabledShortcuts returns all disabled shortcuts
func (sc *ShortcutCustomizer) GetDisabledShortcuts() []string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	var disabled []string
	for id := range sc.disabled {
		disabled = append(disabled, id)
	}
	return disabled
}

// CreateProfile creates a new shortcut profile
func (sc *ShortcutCustomizer) CreateProfile(name string) *ShortcutProfile {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	profile := NewShortcutProfile(name)
	sc.profiles[name] = profile
	return profile
}

// GetProfile retrieves a profile by name
func (sc *ShortcutCustomizer) GetProfile(name string) *ShortcutProfile {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.profiles[name]
}

// ApplyProfile applies a profile's customizations
func (sc *ShortcutCustomizer) ApplyProfile(name string) error {
	sc.mu.RLock()
	profile, exists := sc.profiles[name]
	sc.mu.RUnlock()

	if !exists {
		return errors.New("profile not found")
	}

	profile.mu.RLock()
	defer profile.mu.RUnlock()

	// Apply overrides
	for id, binding := range profile.Overrides {
		sc.OverrideShortcut(id, binding)
	}

	// Apply disabled shortcuts
	for id := range profile.Disabled {
		sc.DisableShortcut(id)
	}

	sc.mu.Lock()
	sc.activeProfile = name
	sc.mu.Unlock()

	return nil
}

// GetActiveProfile returns the currently active profile name
func (sc *ShortcutCustomizer) GetActiveProfile() string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.activeProfile
}

// GetAllProfiles returns all available profile names
func (sc *ShortcutCustomizer) GetAllProfiles() []string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	var names []string
	for name := range sc.profiles {
		names = append(names, name)
	}
	return names
}

// DeleteProfile removes a profile
func (sc *ShortcutCustomizer) DeleteProfile(name string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	delete(sc.profiles, name)
}

// ResetAll removes all customizations
func (sc *ShortcutCustomizer) ResetAll() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.overrides = make(map[string]key.Binding)
	sc.disabled = make(map[string]bool)
	sc.activeProfile = ""
}

// ResetShortcuts removes all shortcut overrides
func (sc *ShortcutCustomizer) ResetShortcuts() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.overrides = make(map[string]key.Binding)
}

// ValidateOverride checks if an override is valid
func (sc *ShortcutCustomizer) ValidateOverride(id string, binding key.Binding) error {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	// Check if shortcut is disabled
	if sc.disabled[id] {
		return errors.New("cannot override a disabled shortcut")
	}

	return nil
}

// GetDefaults returns the default shortcuts
func (sc *ShortcutCustomizer) GetDefaults() map[string]key.Binding {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	defaults := make(map[string]key.Binding)
	for id, binding := range sc.defaultShortcuts {
		defaults[id] = binding
	}
	return defaults
}

// SetDefaults sets the default shortcuts
func (sc *ShortcutCustomizer) SetDefaults(defaults map[string]key.Binding) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.defaultShortcuts = make(map[string]key.Binding)
	for id, binding := range defaults {
		sc.defaultShortcuts[id] = binding
	}
}

// RestoreDefaults restores all shortcuts to defaults
func (sc *ShortcutCustomizer) RestoreDefaults() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.overrides = make(map[string]key.Binding)
	sc.disabled = make(map[string]bool)
}
