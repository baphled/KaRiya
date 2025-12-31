package models

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/key"
)

// ShortcutInfo contains metadata about a keyboard shortcut
type ShortcutInfo struct {
	ID          string
	Binding     key.Binding
	Description string
	Contexts    []string
	Category    string
}

// ShortcutHelpSystem manages help information for keyboard shortcuts
type ShortcutHelpSystem struct {
	mu         sync.RWMutex
	shortcuts  map[string]*ShortcutInfo
	categories map[string]string
	contexts   map[string]bool
}

// NewShortcutHelpSystem creates a new shortcut help system
func NewShortcutHelpSystem() *ShortcutHelpSystem {
	return &ShortcutHelpSystem{
		shortcuts:  make(map[string]*ShortcutInfo),
		categories: make(map[string]string),
		contexts:   make(map[string]bool),
	}
}

// RegisterShortcut registers a shortcut with help information
func (shs *ShortcutHelpSystem) RegisterShortcut(id string, binding key.Binding, description string, contexts []string) {
	shs.mu.Lock()
	defer shs.mu.Unlock()

	info := &ShortcutInfo{
		ID:          id,
		Binding:     binding,
		Description: description,
		Contexts:    make([]string, len(contexts)),
	}
	copy(info.Contexts, contexts)

	shs.shortcuts[id] = info

	for _, ctx := range contexts {
		shs.contexts[ctx] = true
	}
}

// GetShortcutInfo retrieves information about a shortcut
func (shs *ShortcutHelpSystem) GetShortcutInfo(id string) (*ShortcutInfo, bool) {
	shs.mu.RLock()
	defer shs.mu.RUnlock()

	info, exists := shs.shortcuts[id]
	if !exists {
		return nil, false
	}

	infoCopy := *info
	infoCopy.Contexts = make([]string, len(info.Contexts))
	copy(infoCopy.Contexts, info.Contexts)

	return &infoCopy, true
}

// GetAllShortcuts returns all registered shortcuts
func (shs *ShortcutHelpSystem) GetAllShortcuts() map[string]*ShortcutInfo {
	shs.mu.RLock()
	defer shs.mu.RUnlock()

	result := make(map[string]*ShortcutInfo)
	for id, info := range shs.shortcuts {
		infoCopy := *info
		infoCopy.Contexts = make([]string, len(info.Contexts))
		copy(infoCopy.Contexts, info.Contexts)
		result[id] = &infoCopy
	}

	return result
}

// GetShortcutsForContext returns all shortcuts available in a context
func (shs *ShortcutHelpSystem) GetShortcutsForContext(context string) map[string]*ShortcutInfo {
	shs.mu.RLock()
	defer shs.mu.RUnlock()

	result := make(map[string]*ShortcutInfo)

	for id, info := range shs.shortcuts {
		for _, ctx := range info.Contexts {
			if ctx == context {
				infoCopy := *info
				infoCopy.Contexts = make([]string, len(info.Contexts))
				copy(infoCopy.Contexts, info.Contexts)
				result[id] = &infoCopy
				break
			}
		}
	}

	return result
}

// CategorizeShortcut assigns a shortcut to a category
func (shs *ShortcutHelpSystem) CategorizeShortcut(id, category string) {
	shs.mu.Lock()
	defer shs.mu.Unlock()

	shs.categories[id] = category
}

// GetShortcutCategory returns the category of a shortcut
func (shs *ShortcutHelpSystem) GetShortcutCategory(id string) (string, bool) {
	shs.mu.RLock()
	defer shs.mu.RUnlock()

	category, exists := shs.categories[id]
	return category, exists
}

// GetShortcutsByCategory returns all shortcuts in a category
func (shs *ShortcutHelpSystem) GetShortcutsByCategory(category string) map[string]*ShortcutInfo {
	shs.mu.RLock()
	defer shs.mu.RUnlock()

	result := make(map[string]*ShortcutInfo)

	for id, info := range shs.shortcuts {
		if cat, exists := shs.categories[id]; exists && cat == category {
			infoCopy := *info
			infoCopy.Contexts = make([]string, len(info.Contexts))
			copy(infoCopy.Contexts, info.Contexts)
			result[id] = &infoCopy
		}
	}

	return result
}

// GenerateHelpText creates formatted help text for a shortcut
func (shs *ShortcutHelpSystem) GenerateHelpText(id string) string {
	shs.mu.RLock()
	info, exists := shs.shortcuts[id]
	shs.mu.RUnlock()

	if !exists {
		return ""
	}

	help := info.Binding.Help()
	keyStr := help.Key
	if keyStr == "" {
		keyStr = id
	}
	return fmt.Sprintf("%s - %s", keyStr, info.Description)
}

// GenerateContextHelp creates formatted help for all shortcuts in a context
func (shs *ShortcutHelpSystem) GenerateContextHelp(context string) string {
	shortcuts := shs.GetShortcutsForContext(context)

	if len(shortcuts) == 0 {
		return ""
	}

	var lines []string
	for _, info := range shortcuts {
		help := info.Binding.Help()
		keyStr := help.Key
		if keyStr == "" {
			keyStr = info.ID
		}
		lines = append(lines, fmt.Sprintf("  %s - %s", keyStr, info.Description))
	}

	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

// SearchShortcuts finds shortcuts matching a query
func (shs *ShortcutHelpSystem) SearchShortcuts(query string) []*ShortcutInfo {
	shs.mu.RLock()
	defer shs.mu.RUnlock()

	query = strings.ToLower(query)
	var results []*ShortcutInfo

	for _, info := range shs.shortcuts {
		help := info.Binding.Help()
		if strings.Contains(strings.ToLower(help.Key), query) ||
			strings.Contains(strings.ToLower(info.Description), query) ||
			strings.Contains(strings.ToLower(info.ID), query) {
			infoCopy := *info
			infoCopy.Contexts = make([]string, len(info.Contexts))
			copy(infoCopy.Contexts, info.Contexts)
			results = append(results, &infoCopy)
		}
	}

	return results
}

// GetAllContexts returns all registered contexts
func (shs *ShortcutHelpSystem) GetAllContexts() []string {
	shs.mu.RLock()
	defer shs.mu.RUnlock()

	var contexts []string
	for ctx := range shs.contexts {
		contexts = append(contexts, ctx)
	}

	sort.Strings(contexts)
	return contexts
}

// GetAllCategories returns all registered categories
func (shs *ShortcutHelpSystem) GetAllCategories() []string {
	shs.mu.RLock()
	defer shs.mu.RUnlock()

	categoryMap := make(map[string]bool)
	for _, category := range shs.categories {
		categoryMap[category] = true
	}

	var categories []string
	for category := range categoryMap {
		categories = append(categories, category)
	}

	sort.Strings(categories)
	return categories
}

// UnregisterShortcut removes a shortcut from help
func (shs *ShortcutHelpSystem) UnregisterShortcut(id string) {
	shs.mu.Lock()
	defer shs.mu.Unlock()

	delete(shs.shortcuts, id)
	delete(shs.categories, id)
}

// Clear removes all registered shortcuts and categories
func (shs *ShortcutHelpSystem) Clear() {
	shs.mu.Lock()
	defer shs.mu.Unlock()

	shs.shortcuts = make(map[string]*ShortcutInfo)
	shs.categories = make(map[string]string)
	shs.contexts = make(map[string]bool)
}
