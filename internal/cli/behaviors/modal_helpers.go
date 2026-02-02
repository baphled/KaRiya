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

// View renders the static background content for overlay composition.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *StaticViewModel) View() string {
	return m.Content
}

// RenderModalOverlay renders a modal on top of a background view using
// bubbletea-overlay. The modal is centered with a small upward offset
// to avoid the footer.
//
// Expected:
//   - modal must be non-nil and implement Viewable.
//   - background must be a pre-rendered view string.
//
// Returns:
//   - The composited view string with the modal overlaid on the background.
//
// Side effects:
//   - None.
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
// Expected:
//   - modal must be non-nil and implement Viewable.
//   - background must be a pre-rendered view string.
//   - xOffset is the horizontal offset from center (positive = right).
//   - yOffset is the vertical offset from center (positive = down).
//
// Returns:
//   - The composited view string with the modal overlaid at the specified offset.
//
// Side effects:
//   - None.
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

// DefaultModalDimensions provides fallback modal dimensions for use when
// terminal size information is unavailable.
//
// Returns:
//   - width: The default modal width (120).
//   - height: The default modal height (40).
//
// Side effects:
//   - None.
func DefaultModalDimensions() (width, height int) {
	return 120, 40
}
