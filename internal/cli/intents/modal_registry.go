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

// NewModalRegistry creates a new empty modal registry.
func NewModalRegistry() *ModalRegistry {
	return &ModalRegistry{
		modals: make([]ManagedModal, 0, 10), // Pre-allocate for typical modal count
	}
}

// Register adds a modal to the registry.
// Modals are stored in priority order - first registered has highest priority.
// Nil modals are ignored.
func (r *ModalRegistry) Register(modal ManagedModal) {
	if modal != nil {
		r.modals = append(r.modals, modal)
	}
}

// Clear removes all modals from the registry.
// Call this before rebuilding the registry with current modal instances.
func (r *ModalRegistry) Clear() {
	r.modals = r.modals[:0]
}

// GetFirstVisible returns the highest-priority visible modal.
// Returns nil if no modal is visible.
func (r *ModalRegistry) GetFirstVisible() ManagedModal {
	for _, m := range r.modals {
		if m != nil && m.IsVisible() {
			return m
		}
	}
	return nil
}

// HasVisibleModal returns true if any modal in the registry is visible.
func (r *ModalRegistry) HasVisibleModal() bool {
	return r.GetFirstVisible() != nil
}

// HandleUpdate updates the highest-priority visible modal.
// Returns nil if no modal is visible.
// Following the docs: only one modal is active at a time.
func (r *ModalRegistry) HandleUpdate(msg tea.Msg) *ModalUpdateResult {
	if m := r.GetFirstVisible(); m != nil {
		result := m.HandleUpdate(msg)
		return &result
	}
	return nil
}

// RenderOverlay renders the first visible modal over the base view.
// Returns the base view unchanged if no modal is visible.
// Uses behaviors.RenderModalOverlay for consistent overlay compositing.
func (r *ModalRegistry) RenderOverlay(baseView string) string {
	if m := r.GetFirstVisible(); m != nil {
		return behaviors.RenderModalOverlay(m, baseView)
	}
	return baseView
}
