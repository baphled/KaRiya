// Package modals provides skill-specific modal components.
// These modals handle filtering, searching, sorting, and skill CRUD operations
// for the ManageSkills intent.
package modals

import (
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

// staticViewModel implements overlay.Viewable for rendering static content.
// This is used with bubbletea-overlay for modal compositing.
type staticViewModel struct {
	content string
}

func (m staticViewModel) View() string { return m.content }

// RenderOverlayModal renders a modal view over a background using bubbletea-overlay.
// This is a helper for modals that need overlay rendering.
func RenderOverlayModal(modalView, backgroundView string) string {
	modalContent := staticViewModel{content: modalView}
	bgModel := staticViewModel{content: backgroundView}

	overlayModel := overlay.New(
		modalContent,   // Foreground: the modal
		bgModel,        // Background: the rendered view
		overlay.Center, // X position
		overlay.Center, // Y position
		0,              // X offset
		-2,             // Y offset (avoid footer overlap)
	)

	return overlayModel.View()
}
