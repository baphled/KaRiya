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

// View returns the pre-rendered string content stored in this model. It satisfies the
// overlay.Viewable interface so that static markup (such as a modal or background view)
// can be composed via bubbletea-overlay without requiring a full Bubble Tea model.
func (m staticViewModel) View() string { return m.content }

// RenderOverlayModal renders a modal view over a background using bubbletea-overlay.
// This is a helper for modals that need overlay rendering.
func RenderOverlayModal(modalView, backgroundView string) string {
	modalContent := staticViewModel{content: modalView}
	bgModel := staticViewModel{content: backgroundView}

	overlayModel := overlay.New(
		modalContent,
		bgModel,
		overlay.Center,
		overlay.Center,
		0,
		-2,
	)

	return overlayModel.View()
}
