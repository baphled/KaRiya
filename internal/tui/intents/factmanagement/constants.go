// Package factmanagement implements the FactManagement intent for managing career facts.
package factmanagement

// State represents the current state of the FactManagement intent.
type State string

// State constants for the FactManagement intent.
const (
	// StateList shows the facts list view.
	StateList State = "list"

	// StateView shows a single fact detail.
	StateView State = "view"

	// StateEditor shows the fact editor form.
	StateEditor State = "editor"

	// StateDeleteConfirm shows the delete confirmation modal.
	StateDeleteConfirm State = "delete_confirm"

	// StateResults shows operation results.
	StateResults State = "results"

	// StateCompleted indicates the intent has finished.
	StateCompleted State = "completed"
)
