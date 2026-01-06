package intents

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/cli/terminal"
)

// Intent defines the contract for all intent implementations.
// Each intent MUST:
// - Own local navigation state
// - Call domain services
// - Emit artifacts
// - Return a type-safe IntentResult
//
// Intents MUST NOT:
// - Mutate global UI state
// - Navigate into another intent
// - Assume prior context unless explicitly passed
// - Dispatch Bubble Tea commands affecting other intents
// - Use runtime type assertions
// - Pass arbitrary data between intents
type Intent interface {
	// Init is called when the intent is first activated.
	// It may emit commands (e.g., fetch data, start async operations).
	Init() tea.Cmd

	// Update processes a message and returns a command.
	// This is Bubble Tea's standard Update, but constrained to intent-local state.
	Update(msg tea.Msg) tea.Cmd

	// View renders the intent's current state.
	// This is Bubble Tea's standard View.
	View() string

	// Result returns the intent's result if it has completed, or nil if still active.
	// The result communicates the intent's outcome back to the router and parent context.
	Result() *IntentResult[interface{}]
}

// IntentRouter manages the lifecycle of intents and enforces the boundary contract.
// The root Bubble Tea Model delegates to the router, which:
// - Activates/deactivates intents
// - Validates intent results
// - Propagates results back to the root model
// - Prevents illegal state transitions
type IntentRouter interface {
	// ActivateIntent activates an intent by name.
	// Returns an error if the intent doesn't exist or the transition is illegal.
	ActivateIntent(name string, context map[string]interface{}) (tea.Cmd, error)

	// GetActiveIntent returns the currently active intent, or nil if none.
	GetActiveIntent() Intent

	// HandleMessage processes a message in the active intent.
	// Returns a command and any intent result if the intent completed.
	HandleMessage(msg tea.Msg) (tea.Cmd, interface{})

	// View renders the active intent.
	View() string

	// Back navigates back to the previous intent in history.
	// Returns an error if already at the root.
	Back() (tea.Cmd, error)

	// GetHistory returns the navigation history.
	GetHistory() []Intent

	// GetHistoryDepth returns the current depth in the navigation history.
	GetHistoryDepth() int
}

// IntentContext encapsulates minimal, explicit context passed to an intent.
// This is the ONLY way intents receive external data.
// Type safety is enforced by the intent's specific context type.
type IntentContext interface {
	// Validate ensures the context is complete and correct.
	// Returns an error if validation fails.
	Validate() error
}

// ModalEditResult[T] is the result of a modal sub-flow (inline editing).
// It captures both the original and modified values, allowing the parent intent
// to decide whether to accept or discard changes.
type ModalEditResult[T any] struct {
	// Original is the value before editing.
	Original T

	// Modified is the value after editing.
	Modified T

	// Accepted indicates whether the user confirmed the changes.
	Accepted bool

	// Changes is a map of field names to their new values.
	// Only includes fields that were actually changed.
	Changes map[string]interface{}
}

// HasChanges returns true if any fields were modified.
func (m *ModalEditResult[T]) HasChanges() bool {
	return len(m.Changes) > 0
}

// WasAccepted returns true if the user confirmed the changes.
func (m *ModalEditResult[T]) WasAccepted() bool {
	return m.Accepted
}

// GetChange returns the new value for a specific field, or nil if not changed.
func (m *ModalEditResult[T]) GetChange(fieldName string) interface{} {
	if m.Changes == nil {
		return nil
	}
	return m.Changes[fieldName]
}

// NewModalEditResult creates a ModalEditResult with the given original and modified values.
func NewModalEditResult[T any](original, modified T, accepted bool, changes map[string]interface{}) *ModalEditResult[T] {
	if changes == nil {
		changes = make(map[string]interface{})
	}
	return &ModalEditResult[T]{
		Original: original,
		Modified: modified,
		Accepted: accepted,
		Changes:  changes,
	}
}

// NewCancelledModalEditResult creates a ModalEditResult indicating the user cancelled the edit.
func NewCancelledModalEditResult[T any](original T) *ModalEditResult[T] {
	return &ModalEditResult[T]{
		Original: original,
		Modified: original,
		Accepted: false,
		Changes:  make(map[string]interface{}),
	}
}

// TerminalAwareIntent extends the Intent interface with terminal size awareness.
// Intents that implement this interface will receive terminal dimension updates
// and can adapt their rendering accordingly.
type TerminalAwareIntent interface {
	Intent

	// UpdateTerminalInfo updates the intent with current terminal dimensions
	UpdateTerminalInfo(info *terminal.Info)

	// GetMinimumSize returns the minimum terminal size required for this intent
	// Returns width and height in columns and rows
	GetMinimumSize() (width, height int)
}

// BaseIntent provides common functionality that all intents can embed
// It handles terminal size tracking and provides default implementations
type BaseIntent struct {
	terminalInfo   *terminal.Info
	terminalConfig terminal.Config
}

// NewBaseIntent creates a new BaseIntent with default terminal configuration
func NewBaseIntent() *BaseIntent {
	return &BaseIntent{
		terminalInfo:   terminal.NewInfo(),
		terminalConfig: terminal.DefaultConfig,
	}
}

// UpdateTerminalInfo updates the terminal information
func (b *BaseIntent) UpdateTerminalInfo(info *terminal.Info) {
	b.terminalInfo = info
}

// GetTerminalInfo returns the current terminal information
func (b *BaseIntent) GetTerminalInfo() *terminal.Info {
	return b.terminalInfo
}

// GetMinimumSize returns the default minimum terminal size
func (b *BaseIntent) GetMinimumSize() (width, height int) {
	return b.terminalConfig.MinWidth, b.terminalConfig.MinHeight
}
