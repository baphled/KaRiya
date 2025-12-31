package navigation

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Common errors
var (
	ErrScreenNotFound       = errors.New("screen not found in registry")
	ErrEmptyNavigationStack = errors.New("navigation stack is empty")
	ErrInvalidScreenID      = errors.New("invalid screen ID")
	ErrCannotNavigateBack   = errors.New("cannot navigate back from this screen")
	ErrContextKeyNotFound   = errors.New("context key not found")
	ErrNoNavigationHistory  = errors.New("no navigation history available")
	ErrInvalidHistoryIndex  = errors.New("invalid history index")
	ErrCircularNavigation   = errors.New("circular navigation detected")
)

// Screen represents a screen identifier in the application
type Screen string

// ScreenDefinition defines the properties and behavior of a screen
type ScreenDefinition struct {
	// ID is the unique identifier for this screen
	ID Screen

	// Label is the human-readable name shown in breadcrumbs
	Label string

	// Parent is the screen to return to when going back (if nil, uses stack)
	Parent *Screen

	// CanGoBack indicates if back navigation is allowed from this screen
	CanGoBack bool

	// Shortcuts are the keyboard shortcuts available on this screen
	Shortcuts []NavigationKey

	// HelpContext is the context string for contextual help
	HelpContext string

	// IsPersistent indicates if this screen's state should be preserved
	IsPersistent bool

	// IsModal indicates if this screen is a modal/overlay
	IsModal bool

	// Priority affects navigation behavior (higher = more important)
	Priority int

	// Metadata stores additional screen-specific data
	Metadata map[string]interface{}
}

// NavigationState captures the complete state at a point in navigation
type NavigationState struct {
	// Screen is the screen identifier
	Screen Screen

	// Model is the BubbleTea model for this screen (optional)
	Model tea.Model

	// Context is screen-specific context data
	Context map[string]interface{}

	// Timestamp records when this state was created
	Timestamp time.Time

	// Breadcrumbs are the breadcrumb trail at this point
	Breadcrumbs []BreadcrumbItem

	// Metadata stores additional state information
	Metadata map[string]interface{}
}

// BreadcrumbItem represents a single item in the breadcrumb trail
type BreadcrumbItem struct {
	// Label is the display text
	Label string

	// Screen is the screen this breadcrumb links to
	Screen Screen

	// Index is the position in the breadcrumb trail
	Index int

	// IsActive indicates if this is the current screen
	IsActive bool

	// Metadata stores additional breadcrumb data
	Metadata map[string]interface{}
}

// NavigationRegistry manages the navigation state and screen definitions
type NavigationRegistry struct {
	mu sync.RWMutex

	// screens contains all registered screen definitions
	screens map[Screen]*ScreenDefinition

	// history is the complete navigation history stack (for undo/back)
	history []NavigationState

	// future is the forward navigation stack (for redo/forward)
	future []NavigationState

	// current is the current navigation state
	current *NavigationState

	// contextStore is a centralized key-value store for context data
	contextStore map[string]interface{}

	// breadcrumbBuilder generates breadcrumbs from navigation state
	breadcrumbBuilder BreadcrumbBuilder

	// maxHistorySize limits the history stack size
	maxHistorySize int

	// maxFutureSize limits the future stack size
	maxFutureSize int

	// defaultContext is the default context for all screens
	defaultContext context.Context
}

// BreadcrumbBuilder is a function type for generating breadcrumbs
type BreadcrumbBuilder func(current Screen, history []NavigationState, screens map[Screen]*ScreenDefinition) []BreadcrumbItem

// NavigationOptions configures the navigation registry
type NavigationOptions struct {
	MaxHistorySize    int
	BreadcrumbBuilder BreadcrumbBuilder
	DefaultContext    context.Context
}

// NewNavigationRegistry creates a new navigation registry
func NewNavigationRegistry(opts *NavigationOptions) *NavigationRegistry {
	if opts == nil {
		opts = &NavigationOptions{}
	}

	if opts.MaxHistorySize <= 0 {
		opts.MaxHistorySize = 100 // Default max history
	}

	if opts.BreadcrumbBuilder == nil {
		opts.BreadcrumbBuilder = DefaultBreadcrumbBuilder
	}

	if opts.DefaultContext == nil {
		opts.DefaultContext = context.Background()
	}

	return &NavigationRegistry{
		screens:           make(map[Screen]*ScreenDefinition),
		history:           make([]NavigationState, 0),
		future:            make([]NavigationState, 0),
		contextStore:      make(map[string]interface{}),
		breadcrumbBuilder: opts.BreadcrumbBuilder,
		maxHistorySize:    opts.MaxHistorySize,
		maxFutureSize:     opts.MaxHistorySize, // Use same limit for future
		defaultContext:    opts.DefaultContext,
	}
}

// RegisterScreen registers a screen definition
func (r *NavigationRegistry) RegisterScreen(def *ScreenDefinition) error {
	if def == nil {
		return ErrInvalidScreenID
	}

	if def.ID == "" {
		return ErrInvalidScreenID
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for circular parent references
	if def.Parent != nil {
		if err := r.checkCircularReference(def.ID, *def.Parent); err != nil {
			return fmt.Errorf("%w: %v", ErrCircularNavigation, err)
		}
	}

	r.screens[def.ID] = def
	return nil
}

// RegisterScreens registers multiple screen definitions
func (r *NavigationRegistry) RegisterScreens(defs []*ScreenDefinition) error {
	for _, def := range defs {
		if err := r.RegisterScreen(def); err != nil {
			return err
		}
	}
	return nil
}

// GetScreen retrieves a screen definition
func (r *NavigationRegistry) GetScreen(screen Screen) (*ScreenDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	def, ok := r.screens[screen]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrScreenNotFound, screen)
	}

	return def, nil
}

// GetAllScreens returns all registered screen definitions
func (r *NavigationRegistry) GetAllScreens() []*ScreenDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	screens := make([]*ScreenDefinition, 0, len(r.screens))
	for _, def := range r.screens {
		screens = append(screens, def)
	}

	return screens
}

// Navigate pushes a new navigation state onto the stack
func (r *NavigationRegistry) Navigate(screen Screen, model tea.Model, ctx map[string]interface{}) error {
	def, err := r.GetScreen(screen)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Save current state to history if it exists
	if r.current != nil {
		r.history = append(r.history, *r.current)

		// Trim history if needed
		if len(r.history) > r.maxHistorySize {
			r.history = r.history[1:]
		}
	}

	// Clear future stack on new navigation (cannot redo after new action)
	r.future = make([]NavigationState, 0)

	// Build breadcrumbs
	breadcrumbs := r.breadcrumbBuilder(screen, r.history, r.screens)

	// Create new current state
	r.current = &NavigationState{
		Screen:      screen,
		Model:       model,
		Context:     ctx,
		Timestamp:   time.Now(),
		Breadcrumbs: breadcrumbs,
		Metadata:    make(map[string]interface{}),
	}

	// Store screen definition in metadata
	r.current.Metadata["screen_definition"] = def

	return nil
}

// Back navigates to the previous screen
func (r *NavigationRegistry) Back() (*NavigationState, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.current == nil {
		return nil, ErrEmptyNavigationStack
	}

	// Check if current screen allows back navigation
	if def, ok := r.current.Metadata["screen_definition"].(*ScreenDefinition); ok {
		if !def.CanGoBack {
			return nil, fmt.Errorf("%w: %s", ErrCannotNavigateBack, def.ID)
		}

		// If screen has explicit parent, navigate there
		if def.Parent != nil {
			// Find the parent in history or create new state
			return r.navigateToScreen(*def.Parent)
		}
	}

	// Pop from history
	if len(r.history) == 0 {
		return nil, ErrEmptyNavigationStack
	}

	// Get previous state
	previousState := r.history[len(r.history)-1]
	r.history = r.history[:len(r.history)-1]

	// Set as current
	r.current = &previousState

	return r.current, nil
}

// Forward navigates forward in history (redo)
func (r *NavigationRegistry) Forward() (*NavigationState, error) {
	// Note: This requires a separate "future" stack for redo functionality
	// Will be implemented in task 2.5
	return nil, errors.New("forward navigation not yet implemented")
}

// GetCurrent returns the current navigation state
func (r *NavigationRegistry) GetCurrent() *NavigationState {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.current == nil {
		return nil
	}

	// Return a copy to prevent external modification
	currentCopy := *r.current
	return &currentCopy
}

// GetHistory returns the navigation history
func (r *NavigationRegistry) GetHistory() []NavigationState {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy
	historyCopy := make([]NavigationState, len(r.history))
	copy(historyCopy, r.history)

	return historyCopy
}

// GetHistorySize returns the number of items in history
func (r *NavigationRegistry) GetHistorySize() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.history)
}

// NavigateToHistoryIndex navigates to a specific point in history
func (r *NavigationRegistry) NavigateToHistoryIndex(index int) (*NavigationState, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if index < 0 || index >= len(r.history) {
		return nil, fmt.Errorf("%w: %d", ErrInvalidHistoryIndex, index)
	}

	// Get target state
	targetState := r.history[index]

	// Save current state if needed
	if r.current != nil {
		// Move current and subsequent history items to a separate stack if implementing redo
		// For now, we just replace current
	}

	// Trim history up to index
	r.history = r.history[:index]

	// Set as current
	r.current = &targetState

	return r.current, nil
}

// GetBreadcrumbs returns the current breadcrumb trail
func (r *NavigationRegistry) GetBreadcrumbs() []BreadcrumbItem {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.current == nil {
		return []BreadcrumbItem{}
	}

	// Return a copy
	breadcrumbs := make([]BreadcrumbItem, len(r.current.Breadcrumbs))
	copy(breadcrumbs, r.current.Breadcrumbs)

	return breadcrumbs
}

// GetBreadcrumbPath returns the breadcrumb trail as a string
func (r *NavigationRegistry) GetBreadcrumbPath() string {
	breadcrumbs := r.GetBreadcrumbs()
	labels := make([]string, len(breadcrumbs))

	for i, crumb := range breadcrumbs {
		labels[i] = crumb.Label
	}

	return strings.Join(labels, " > ")
}

// NavigateToBreadcrumb navigates to a specific breadcrumb
func (r *NavigationRegistry) NavigateToBreadcrumb(index int) (*NavigationState, error) {
	breadcrumbs := r.GetBreadcrumbs()

	if index < 0 || index >= len(breadcrumbs) {
		return nil, fmt.Errorf("%w: %d", ErrInvalidHistoryIndex, index)
	}

	targetScreen := breadcrumbs[index].Screen

	// Find the screen in history
	return r.navigateToScreen(targetScreen)
}

// SetContext stores a value in the context store
func (r *NavigationRegistry) SetContext(key string, value interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.contextStore[key] = value
}

// GetContext retrieves a value from the context store
func (r *NavigationRegistry) GetContext(key string) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	value, ok := r.contextStore[key]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrContextKeyNotFound, key)
	}

	return value, nil
}

// ClearContext removes a value from the context store
func (r *NavigationRegistry) ClearContext(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.contextStore, key)
}

// ClearAllContext clears the entire context store
func (r *NavigationRegistry) ClearAllContext() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.contextStore = make(map[string]interface{})
}

// Reset clears all navigation state
func (r *NavigationRegistry) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.history = make([]NavigationState, 0)
	r.current = nil
	r.contextStore = make(map[string]interface{})
}

// CanGoBack returns true if back navigation is possible
func (r *NavigationRegistry) CanGoBack() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.current == nil {
		return false
	}

	// Check if current screen allows back navigation
	if def, ok := r.current.Metadata["screen_definition"].(*ScreenDefinition); ok {
		if !def.CanGoBack {
			return false
		}
	}

	// Check if there's history or a parent screen
	if len(r.history) > 0 {
		return true
	}

	if def, ok := r.current.Metadata["screen_definition"].(*ScreenDefinition); ok {
		return def.Parent != nil
	}

	return false
}

// GetBackTarget returns the screen that would be navigated to on back
func (r *NavigationRegistry) GetBackTarget() (*Screen, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.CanGoBack() {
		return nil, ErrCannotNavigateBack
	}

	// Check for explicit parent
	if def, ok := r.current.Metadata["screen_definition"].(*ScreenDefinition); ok {
		if def.Parent != nil {
			return def.Parent, nil
		}
	}

	// Use history
	if len(r.history) > 0 {
		screen := r.history[len(r.history)-1].Screen
		return &screen, nil
	}

	return nil, ErrCannotNavigateBack
}

// Helper methods

// checkCircularReference checks for circular parent references
func (r *NavigationRegistry) checkCircularReference(screen Screen, parent Screen) error {
	visited := make(map[Screen]bool)
	current := parent

	for {
		if current == screen {
			return fmt.Errorf("circular reference: %s -> ... -> %s", screen, parent)
		}

		if visited[current] {
			// Already checked this path
			break
		}

		visited[current] = true

		def, ok := r.screens[current]
		if !ok || def.Parent == nil {
			// Reached end of parent chain
			break
		}

		current = *def.Parent
	}

	return nil
}

// navigateToScreen finds a screen in history or creates new state
func (r *NavigationRegistry) navigateToScreen(screen Screen) (*NavigationState, error) {
	// Search backwards through history for the screen
	for i := len(r.history) - 1; i >= 0; i-- {
		if r.history[i].Screen == screen {
			// Found it - navigate to this point
			targetState := r.history[i]
			r.history = r.history[:i]
			r.current = &targetState
			return r.current, nil
		}
	}

	// Not in history - check if screen exists
	def, err := r.GetScreen(screen)
	if err != nil {
		return nil, err
	}

	// Current is not saved to history in this case - we're jumping to a new screen
	// The history remains as-is

	// Create new state for this screen
	// Note: Caller should provide model, but for parent navigation we create minimal state
	breadcrumbs := r.breadcrumbBuilder(screen, r.history, r.screens)

	newState := &NavigationState{
		Screen:      screen,
		Model:       nil,
		Context:     make(map[string]interface{}),
		Timestamp:   time.Now(),
		Breadcrumbs: breadcrumbs,
		Metadata: map[string]interface{}{
			"screen_definition": def,
		},
	}

	// Don't modify history when jumping to parent
	r.current = newState
	return r.current, nil
}

// DefaultBreadcrumbBuilder is the default breadcrumb generation function
func DefaultBreadcrumbBuilder(current Screen, history []NavigationState, screens map[Screen]*ScreenDefinition) []BreadcrumbItem {
	breadcrumbs := make([]BreadcrumbItem, 0)

	// Build breadcrumb trail from history and current
	seen := make(map[Screen]bool)

	// Add from history (skip duplicates and modals)
	for _, state := range history {
		if seen[state.Screen] {
			continue
		}

		def, ok := screens[state.Screen]
		if !ok {
			continue
		}

		// Skip modals in breadcrumbs
		if def.IsModal {
			continue
		}

		breadcrumbs = append(breadcrumbs, BreadcrumbItem{
			Label:    def.Label,
			Screen:   state.Screen,
			Index:    len(breadcrumbs),
			IsActive: false,
		})

		seen[state.Screen] = true
	}

	// Add current screen
	if def, ok := screens[current]; ok && !def.IsModal {
		breadcrumbs = append(breadcrumbs, BreadcrumbItem{
			Label:    def.Label,
			Screen:   current,
			Index:    len(breadcrumbs),
			IsActive: true,
		})
	}

	return breadcrumbs
}

// HierarchicalBreadcrumbBuilder builds breadcrumbs from parent hierarchy
func HierarchicalBreadcrumbBuilder(current Screen, history []NavigationState, screens map[Screen]*ScreenDefinition) []BreadcrumbItem {
	breadcrumbs := make([]BreadcrumbItem, 0)

	// Build parent chain from current screen
	chain := make([]Screen, 0)
	currentScreen := current

	visited := make(map[Screen]bool)

	for {
		if visited[currentScreen] {
			// Circular reference - break
			break
		}

		visited[currentScreen] = true
		chain = append([]Screen{currentScreen}, chain...) // Prepend

		def, ok := screens[currentScreen]
		if !ok || def.Parent == nil {
			break
		}

		currentScreen = *def.Parent
	}

	// Convert chain to breadcrumbs
	for i, screen := range chain {
		def, ok := screens[screen]
		if !ok {
			continue
		}

		breadcrumbs = append(breadcrumbs, BreadcrumbItem{
			Label:    def.Label,
			Screen:   screen,
			Index:    i,
			IsActive: screen == current,
		})
	}

	return breadcrumbs
}
