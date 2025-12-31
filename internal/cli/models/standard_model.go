package models

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// BreadcrumbItem represents a single item in the navigation breadcrumb trail
type BreadcrumbItem struct {
	Label    string
	ID       string
	Metadata map[string]interface{}
}

// ContextMetadata contains contextual information about the current model state
type ContextMetadata struct {
	ScreenID      string
	PreviousState interface{}
	Data          map[string]interface{}
}

// StandardModel defines a comprehensive interface for models in the KaRiya CLI
// It provides a standardized approach to model initialization, rendering,
// state management, and interaction.
type StandardModel interface {
	// Core Model Lifecycle Methods
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string

	// Context and Navigation Management
	GetContext() context.Context
	SetContext(ctx context.Context)
	GetContextMetadata() *ContextMetadata
	SetContextMetadata(metadata *ContextMetadata)

	// Enhanced Breadcrumb Management
	GetBreadcrumbs() []BreadcrumbItem
	AddBreadcrumb(item BreadcrumbItem)
	PopBreadcrumb() BreadcrumbItem
	PeekBreadcrumb() *BreadcrumbItem
	GetBreadcrumbPath() string
	NavigateToBreadcrumb(index int) error
	ClearBreadcrumbs()

	// Navigation History
	PushNavigationHistory(label string, state interface{})
	PopNavigationHistory() (string, interface{}, error)
	GetNavigationHistory() []string

	// Keyboard Shortcut Handling
	RegisterShortcuts(shortcuts map[string]key.Binding)
	GetShortcuts() map[string]key.Binding
	HandleShortcut(key key.Binding) (tea.Model, tea.Cmd)

	// Error Handling
	GetLastError() error
	SetError(err error)
	ClearError()

	// State Management
	GetState() interface{}
	SetState(state interface{})

	// Validation and Interaction
	Validate() error
	Reset()
}

// NavigationHistoryItem represents an entry in the navigation history
type NavigationHistoryItem struct {
	Label string
	State interface{}
}

// BaseStandardModel provides a default implementation of StandardModel
// to be embedded or extended by specific model implementations
type BaseStandardModel struct {
	ctx                context.Context
	contextMetadata    *ContextMetadata
	breadcrumbs        []BreadcrumbItem
	navigationHistory  []NavigationHistoryItem
	shortcuts          map[string]key.Binding
	lastError          error
	state              interface{}
}

// NewBaseStandardModel creates a new instance of BaseStandardModel with initialized fields
func NewBaseStandardModel() *BaseStandardModel {
	return &BaseStandardModel{
		ctx:               context.Background(),
		contextMetadata:   &ContextMetadata{Data: make(map[string]interface{})},
		breadcrumbs:       []BreadcrumbItem{},
		navigationHistory: []NavigationHistoryItem{},
		shortcuts:         make(map[string]key.Binding),
	}
}

// Implement tea.Model interface methods

func (b *BaseStandardModel) Init() tea.Cmd {
	return nil
}

func (b *BaseStandardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b *BaseStandardModel) View() string {
	return ""
}

// Implement context management methods

func (b *BaseStandardModel) GetContext() context.Context {
	if b.ctx == nil {
		b.ctx = context.Background()
	}
	return b.ctx
}

func (b *BaseStandardModel) SetContext(ctx context.Context) {
	b.ctx = ctx
}

func (b *BaseStandardModel) GetContextMetadata() *ContextMetadata {
	if b.contextMetadata == nil {
		b.contextMetadata = &ContextMetadata{Data: make(map[string]interface{})}
	}
	return b.contextMetadata
}

func (b *BaseStandardModel) SetContextMetadata(metadata *ContextMetadata) {
	b.contextMetadata = metadata
}

// Implement enhanced breadcrumb management methods

func (b *BaseStandardModel) GetBreadcrumbs() []BreadcrumbItem {
	return b.breadcrumbs
}

func (b *BaseStandardModel) AddBreadcrumb(item BreadcrumbItem) {
	if item.Metadata == nil {
		item.Metadata = make(map[string]interface{})
	}
	b.breadcrumbs = append(b.breadcrumbs, item)
}

func (b *BaseStandardModel) PopBreadcrumb() BreadcrumbItem {
	if len(b.breadcrumbs) == 0 {
		return BreadcrumbItem{}
	}
	lastIndex := len(b.breadcrumbs) - 1
	item := b.breadcrumbs[lastIndex]
	b.breadcrumbs = b.breadcrumbs[:lastIndex]
	return item
}

func (b *BaseStandardModel) PeekBreadcrumb() *BreadcrumbItem {
	if len(b.breadcrumbs) == 0 {
		return nil
	}
	return &b.breadcrumbs[len(b.breadcrumbs)-1]
}

func (b *BaseStandardModel) GetBreadcrumbPath() string {
	if len(b.breadcrumbs) == 0 {
		return "/"
	}
	path := "/"
	for _, crumb := range b.breadcrumbs {
		path += crumb.Label + "/"
	}
	return path
}

func (b *BaseStandardModel) NavigateToBreadcrumb(index int) error {
	if index < 0 || index >= len(b.breadcrumbs) {
		return ErrInvalidBreadcrumbIndex
	}
	// Remove all breadcrumbs after the target index
	b.breadcrumbs = b.breadcrumbs[:index+1]
	return nil
}

func (b *BaseStandardModel) ClearBreadcrumbs() {
	b.breadcrumbs = []BreadcrumbItem{}
}

// Implement navigation history methods

func (b *BaseStandardModel) PushNavigationHistory(label string, state interface{}) {
	b.navigationHistory = append(b.navigationHistory, NavigationHistoryItem{
		Label: label,
		State: state,
	})
}

func (b *BaseStandardModel) PopNavigationHistory() (string, interface{}, error) {
	if len(b.navigationHistory) == 0 {
		return "", nil, ErrEmptyNavigationHistory
	}
	lastIndex := len(b.navigationHistory) - 1
	item := b.navigationHistory[lastIndex]
	b.navigationHistory = b.navigationHistory[:lastIndex]
	return item.Label, item.State, nil
}

func (b *BaseStandardModel) GetNavigationHistory() []string {
	history := make([]string, len(b.navigationHistory))
	for i, item := range b.navigationHistory {
		history[i] = item.Label
	}
	return history
}

// Implement keyboard shortcut handling methods

func (b *BaseStandardModel) RegisterShortcuts(shortcuts map[string]key.Binding) {
	b.shortcuts = shortcuts
}

func (b *BaseStandardModel) GetShortcuts() map[string]key.Binding {
	return b.shortcuts
}

func (b *BaseStandardModel) HandleShortcut(key key.Binding) (tea.Model, tea.Cmd) {
	// Default implementation - can be overridden by specific models
	return b, nil
}

// Implement error handling methods

func (b *BaseStandardModel) GetLastError() error {
	return b.lastError
}

func (b *BaseStandardModel) SetError(err error) {
	b.lastError = err
}

func (b *BaseStandardModel) ClearError() {
	b.lastError = nil
}

// Implement state management methods

func (b *BaseStandardModel) GetState() interface{} {
	return b.state
}

func (b *BaseStandardModel) SetState(state interface{}) {
	b.state = state
}

// Validation and Interaction - default no-op implementations

func (b *BaseStandardModel) Validate() error {
	return nil
}

func (b *BaseStandardModel) Reset() {
	b.breadcrumbs = []BreadcrumbItem{}
	b.navigationHistory = []NavigationHistoryItem{}
	b.shortcuts = make(map[string]key.Binding)
	b.lastError = nil
	b.state = nil
	b.contextMetadata = &ContextMetadata{Data: make(map[string]interface{})}
}

