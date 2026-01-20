// Package behaviors provides embeddable components for table-based UIs.
package behaviors

import (
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

// Viewable is the interface required for modal overlay rendering.
// Any type with a View() method can be used as a modal.
type Viewable interface {
	View() string
}

// StaticViewModel is a simple Viewable that returns static content.
// It's used as the background layer when rendering modal overlays.
type StaticViewModel struct {
	Content string
}

// View implements Viewable.
func (m *StaticViewModel) View() string {
	return m.Content
}

// RenderModalOverlay renders a modal on top of a background view using
// bubbletea-overlay. The modal is centered with a small upward offset
// to avoid the footer.
//
// Parameters:
//   - modal: The foreground modal (must implement Viewable - just needs View() string)
//   - background: The rendered background view string
//
// Returns the composited view string.
func RenderModalOverlay(modal Viewable, background string) string {
	bgModel := &StaticViewModel{Content: background}
	overlayModel := overlay.New(
		modal,          // Foreground: the modal
		bgModel,        // Background: the rendered view
		overlay.Center, // X position
		overlay.Center, // Y position
		0,              // X offset
		-2,             // Y offset (move up 2 lines to avoid footer)
	)
	return overlayModel.View()
}

// RenderModalOverlayWithOffset renders a modal with custom positioning.
//
// Parameters:
//   - modal: The foreground modal (must implement Viewable - just needs View() string)
//   - background: The rendered background view string
//   - xOffset: Horizontal offset from center (positive = right)
//   - yOffset: Vertical offset from center (positive = down)
//
// Returns the composited view string.
func RenderModalOverlayWithOffset(modal Viewable, background string, xOffset, yOffset int) string {
	bgModel := &StaticViewModel{Content: background}
	overlayModel := overlay.New(
		modal,          // Foreground: the modal
		bgModel,        // Background: the rendered view
		overlay.Center, // X position
		overlay.Center, // Y position
		xOffset,        // X offset
		yOffset,        // Y offset
	)
	return overlayModel.View()
}

// DefaultModalDimensions returns sensible default dimensions for modals
// when terminal info is not available.
func DefaultModalDimensions() (width, height int) {
	return 120, 40
}
