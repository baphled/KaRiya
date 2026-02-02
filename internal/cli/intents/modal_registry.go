// Package intents provides intent implementations for the KaRiya TUI.
package intents

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	tea "github.com/charmbracelet/bubbletea"
)

// ModalUpdateResult holds the normalized outcome of a modal update.
// This provides a unified return type for all modal updates regardless of
// the underlying modal's specific signature.
type ModalUpdateResult struct {
	// Cmd is the tea.Cmd to return to the runtime.
	Cmd tea.Cmd

	// Closed indicates the modal was closed (via Esc, completion, etc.).
	Closed bool

	// Applied indicates the user submitted/confirmed (vs. cancelled).
	// For form modals: true if form was submitted.
	// For confirm modals: true if user confirmed.
	Applied bool

	// Data holds type-specific result data.
	// For form modals: the form data struct.
	// For confirm modals: the bool confirmation result.
	// For view modals: nil.
	Data interface{}
}

// ManagedModal is the unified interface for all modals in the registry.
// It normalizes the different modal APIs (form, confirm, view) into a
// consistent interface for the intent to interact with.
type ManagedModal interface {
	// IsVisible returns true if the modal is currently visible.
	IsVisible() bool

	// View returns the modal's rendered content.
	View() string

	// HandleUpdate processes a message and returns a normalized result.
	HandleUpdate(msg tea.Msg) ModalUpdateResult
}

// ModalRegistry manages a prioritized collection of modals.
// It provides unified handling for Update and View operations,
// following the "one modal at a time" principle from the docs.
//
// Modals are registered in priority order (first registered = highest priority).
// Only the highest-priority visible modal receives updates and renders.
type ModalRegistry struct {
	modals []ManagedModal
}

// NewModalRegistry creates a new empty modal registry with pre-allocated
//
// Returns:
//   - A fully initialized ModalRegistry ready for use.
//
// Side effects:
//   - None.
func NewModalRegistry() *ModalRegistry {
	return &ModalRegistry{
		modals: make([]ManagedModal, 0, 10), // Pre-allocate for typical modal count
	}
}

// Register adds a modal to the registry in priority order, where the first
//
// Expected:
//   - managedmodal must be valid.
//
// Side effects:
//   - None.
func (r *ModalRegistry) Register(modal ManagedModal) {
	if modal != nil {
		r.modals = append(r.modals, modal)
	}
}

// Clear removes all modals from the registry while retaining the
//
// Side effects:
//   - None.
func (r *ModalRegistry) Clear() {
	r.modals = r.modals[:0]
}

// GetFirstVisible locates the highest-priority modal that is currently
//
// Returns:
//   - A ManagedModal value.
//
// Side effects:
//   - None.
func (r *ModalRegistry) GetFirstVisible() ManagedModal {
	for _, m := range r.modals {
		if m != nil && m.IsVisible() {
			return m
		}
	}
	return nil
}

// HasVisibleModal checks whether any registered modal is currently displayed,
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (r *ModalRegistry) HasVisibleModal() bool {
	return r.GetFirstVisible() != nil
}

// HandleUpdate delegates a Bubble Tea message to the highest-priority visible
// modal, enforcing the single-active-modal policy.
//
// Expected:
//   - msg is the Bubble Tea message to forward to the active modal.
//
// Returns:
//   - A pointer to the ModalUpdateResult from the active modal, or nil if no modal is visible.
//
// Side effects:
//   - Invokes HandleUpdate on the highest-priority visible modal, which may mutate its state.
func (r *ModalRegistry) HandleUpdate(msg tea.Msg) *ModalUpdateResult {
	if m := r.GetFirstVisible(); m != nil {
		result := m.HandleUpdate(msg)
		return &result
	}
	return nil
}

// RenderOverlay composites the highest-priority visible modal on top of the
// provided base view using the standard overlay compositing behavior.
//
// Expected:
//   - baseView is the rendered screen content to use as the background layer.
//
// Returns:
//   - The composited view with the modal overlay, or baseView unchanged if no modal is visible.
//
// Side effects:
//   - None.
func (r *ModalRegistry) RenderOverlay(baseView string) string {
	if m := r.GetFirstVisible(); m != nil {
		return behaviors.RenderModalOverlay(m, baseView)
	}
	return baseView
}
