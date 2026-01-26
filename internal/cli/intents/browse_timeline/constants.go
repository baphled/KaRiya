// Package browse_timeline implements the BrowseTimeline intent for browsing career events.
package browse_timeline

// State represents the current state of the BrowseTimeline intent.
type State string

// State constants for the BrowseTimeline intent.
const (
	// StateTimeline shows the timeline list view.
	StateTimeline State = "timeline"

	// StateDeleteConfirm shows the delete confirmation modal.
	StateDeleteConfirm State = "delete_confirm"
)
