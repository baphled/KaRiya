package themes

import (
	"errors"
	"sort"
	"sync"
)

// ThemeManager handles theme registration and switching.
// It maintains a registry of available themes and manages the active theme.
type ThemeManager struct {
	mu          sync.RWMutex
	themes      map[string]Theme
	active      Theme
	defaultName string
	onChange    map[int]func(Theme)
	nextID      int
}

// NewThemeManager creates a new ThemeManager with the default theme registered.
func NewThemeManager() *ThemeManager {
	tm := &ThemeManager{
		themes:      make(map[string]Theme),
		defaultName: "default",
		onChange:    make(map[int]func(Theme)),
		nextID:      0,
	}

	// Register and activate the default theme
	defaultTheme := NewDefaultTheme()
	tm.themes[defaultTheme.Name()] = defaultTheme
	tm.active = defaultTheme

	return tm
}

// NewEmptyThemeManager creates a ThemeManager without any themes registered.
// This is primarily for testing purposes.
func NewEmptyThemeManager() *ThemeManager {
	return &ThemeManager{
		themes:      make(map[string]Theme),
		defaultName: "",
		onChange:    make(map[int]func(Theme)),
		nextID:      0,
	}
}

// Register adds a new theme to the manager.
// Returns an error if the theme is nil, has an empty name, or already exists.
func (tm *ThemeManager) Register(theme Theme) error {
	if theme == nil {
		return errors.New("cannot register nil theme")
	}

	if theme.Name() == "" {
		return errors.New("theme name cannot be empty")
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, exists := tm.themes[theme.Name()]; exists {
		return errors.New("theme '" + theme.Name() + "' is already registered")
	}

	tm.themes[theme.Name()] = theme
	return nil
}

// SetActive sets the active theme by name.
// Returns an error if the theme is not found.
func (tm *ThemeManager) SetActive(name string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	theme, exists := tm.themes[name]
	if !exists {
		return errors.New("theme '" + name + "' not found")
	}

	tm.active = theme

	// Notify all listeners
	for _, callback := range tm.onChange {
		callback(theme)
	}

	return nil
}

// Active returns the currently active theme.
// Returns nil if no theme is set.
func (tm *ThemeManager) Active() Theme {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.active
}

// Get returns a theme by name.
// Returns an error if the theme is not found.
func (tm *ThemeManager) Get(name string) (Theme, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	theme, exists := tm.themes[name]
	if !exists {
		return nil, errors.New("theme '" + name + "' not found")
	}

	return theme, nil
}

// List returns the names of all registered themes in alphabetical order.
func (tm *ThemeManager) List() []string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	names := make([]string, 0, len(tm.themes))
	for name := range tm.themes {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}

// OnChange registers a callback that is called when the active theme changes.
// Returns an ID that can be used to remove the callback.
func (tm *ThemeManager) OnChange(callback func(Theme)) int {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	id := tm.nextID
	tm.nextID++
	tm.onChange[id] = callback

	return id
}

// RemoveChangeCallback removes a previously registered callback.
func (tm *ThemeManager) RemoveChangeCallback(id int) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	delete(tm.onChange, id)
}

// Styles returns the StyleSet of the active theme.
// Returns nil if no theme is active.
func (tm *ThemeManager) Styles() *StyleSet {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	if tm.active == nil {
		return nil
	}
	return tm.active.Styles()
}

// Palette returns the ColorPalette of the active theme.
// Returns nil if no theme is active.
func (tm *ThemeManager) Palette() *ColorPalette {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	if tm.active == nil {
		return nil
	}
	return tm.active.Palette()
}
