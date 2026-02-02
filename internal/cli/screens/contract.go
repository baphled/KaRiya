package screens

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Screen represents a single UI screen in the TUI application.
//
// A Screen is responsible for:
// - Rendering a specific view (UI state)
// - Handling user input for that view
// - Returning results to indicate navigation or state changes
//
// Screens are orchestrated by Intents, which manage workflow state transitions.
// Unlike the monolithic intent pattern, Screens are:
// - Focused on a single UI concern (one "page")
// - Reusable across multiple intents
// - Composable (can be combined or nested)
//
// Example workflow:
//
//	Intent receives msg → passes to active Screen → Screen returns ScreenResult →
//	Intent handles result (navigate to new Screen, cancel, submit data, etc.)
//
// Related:
// - tasks/tasks-42-tui-architecture-refactor.md
// - docs/TUI_DEVELOPER_GUIDE.md
//
//nolint:interfacebloat // Screen interface provides complete screen lifecycle and rendering contract
type Screen interface {
	// Update handles a BubbleTea message and returns:
	// - A BubbleTea command (for async operations, window updates, etc.)
	// - An optional ScreenResult indicating a state change (navigate, cancel, submit, error)
	//
	// The Screen should return nil for ScreenResult when the message only updates
	// internal state (e.g., cursor movement, typing).
	//
	// ScreenResult is returned when the Screen's purpose is complete:
	// - User selects an item → NavigateResult with selected data
	// - User presses Escape → CancelResult
	// - User submits a form → SubmitResult with form data
	// - An error occurs → ErrorResult with error info
	Update(msg tea.Msg) (tea.Cmd, ScreenResult)

	// View renders the Screen's current state as a string.
	// This should use StandardView for consistency across all screens.
	//
	// For intents that need more control over the layout, use
	// RenderContent() to get just the content without StandardView wrapper.
	View() string

	// RenderContent returns the screen's content without StandardView wrapper.
	// This allows intents to wrap the content in their own StandardView
	// with custom breadcrumbs, help text, and styling.
	// This is the preferred pattern for screen-based intents.
	RenderContent() string

	// SetTerminalInfo updates the Screen's knowledge of terminal dimensions.
	// Screens should store this and pass it to StandardView for proper layout.
	SetTerminalInfo(width, height int)

	// SetTheme updates the Screen's theme for styling.
	// Screens should store this and pass it to StandardView for consistent theming.
	SetTheme(theme interface{})

	// SetLogo updates the Screen's logo and spacing.
	// Screens should store this and pass it to StandardView for consistent branding.
	SetLogo(logo interface{}, spacing int)
}

// ScreenResultType indicates the type of result a Screen is returning.
type ScreenResultType string

const (
	// ResultNavigate indicates the user selected something or wants to move forward.
	// The Intent should transition to the next screen based on the data in the result.
	ResultNavigate ScreenResultType = "navigate"

	// ResultCancel indicates the user pressed Escape or wants to go back.
	// The Intent should transition to the previous screen or cancel the workflow.
	ResultCancel ScreenResultType = "cancel"

	// ResultSubmit indicates the user submitted a form or completed an action.
	// The Intent should process the submitted data and continue the workflow.
	ResultSubmit ScreenResultType = "submit"

	// ResultError indicates an error occurred in the Screen.
	// The Intent should display the error or transition to an error state.
	ResultError ScreenResultType = "error"
)

// ScreenResult represents the outcome of a Screen's Update operation.
//
// Screens return ScreenResults to communicate with their parent Intent:
// - "I'm done, here's the data" (Navigate, Submit)
// - "User cancelled" (Cancel)
// - "Something went wrong" (Error)
//
// The Intent examines the ScreenResult type and data to determine the next action.
type ScreenResult interface {
	// Type returns the type of result (Navigate, Cancel, Submit, Error).
	Type() ScreenResultType

	// Data returns the result data (selection, form data, error, etc.).
	// The type of data depends on the result type:
	// - Navigate: selected item, next state identifier, etc.
	// - Submit: form data, user input, etc.
	// - Cancel: typically nil
	// - Error: error message, error code, etc.
	Data() interface{}

	// Metadata returns additional context for the result.
	// This is useful for preserving state when navigating back:
	// - Scroll position
	// - Filter settings
	// - Previously selected items
	// - Error context
	Metadata() map[string]interface{}

	// WithMetadata adds metadata to the result (fluent API).
	// Returns the same result with metadata attached.
	WithMetadata(key string, value interface{}) ScreenResult
}

// NavigateResult is returned when the user selects an item or wants to move forward.
//
// Example: User selects "Profile A" from a list
//
//	return &NavigateResult{
//	    ResultData: profileA,
//	}
//
// The Intent receives this result and transitions to the next screen,
// passing the selected profile as context.
type NavigateResult struct {
	ResultData interface{}
	Meta       map[string]interface{}
}

// Type returns ResultNavigate, identifying this result as a forward
//
// Returns:
//   - A ScreenResultType value.
//
// Side effects:
//   - None.
func (r *NavigateResult) Type() ScreenResultType {
	return ResultNavigate
}

// Data returns the navigation payload stored in ResultData. The caller
//
// Returns:
//   - A interface{} value.
//
// Side effects:
//   - None.
func (r *NavigateResult) Data() interface{} {
	return r.ResultData
}

// Metadata returns the key-value metadata map attached to this result.
//
// Returns:
//   - A map[string]interface{} value.
//
// Side effects:
//   - None.
func (r *NavigateResult) Metadata() map[string]interface{} {
	if r.Meta == nil {
		r.Meta = make(map[string]interface{})
	}
	return r.Meta
}

// WithMetadata stores a key-value pair in the result metadata and returns
//
// Expected:
//   - Must be a valid string.
//   - interface{} must be valid.
//
// Returns:
//   - A ScreenResult value.
//
// Side effects:
//   - None.
func (r *NavigateResult) WithMetadata(key string, value interface{}) ScreenResult {
	if r.Meta == nil {
		r.Meta = make(map[string]interface{})
	}
	r.Meta[key] = value
	return r
}

// CancelResult is returned when the user presses Escape or cancels the screen.
//
// Example: User presses Escape on a selection screen
//
//	return &CancelResult{}
//
// The Intent receives this result and transitions back to the previous screen,
// optionally restoring context from metadata.
type CancelResult struct {
	Meta map[string]interface{}
}

// Type returns ResultCancel, identifying this result as a backward
//
// Returns:
//   - A ScreenResultType value.
//
// Side effects:
//   - None.
func (r *CancelResult) Type() ScreenResultType {
	return ResultCancel
}

// Data always returns nil because cancellation carries no payload.
//
// Returns:
//   - A interface{} value.
//
// Side effects:
//   - None.
func (r *CancelResult) Data() interface{} {
	return nil
}

// Metadata returns the key-value metadata map attached to this result.
//
// Returns:
//   - A map[string]interface{} value.
//
// Side effects:
//   - None.
func (r *CancelResult) Metadata() map[string]interface{} {
	if r.Meta == nil {
		r.Meta = make(map[string]interface{})
	}
	return r.Meta
}

// WithMetadata stores a key-value pair in the result metadata and returns
//
// Expected:
//   - Must be a valid string.
//   - interface{} must be valid.
//
// Returns:
//   - A ScreenResult value.
//
// Side effects:
//   - None.
func (r *CancelResult) WithMetadata(key string, value interface{}) ScreenResult {
	if r.Meta == nil {
		r.Meta = make(map[string]interface{})
	}
	r.Meta[key] = value
	return r
}

// SubmitResult is returned when the user submits a form or completes an action.
//
// Example: User submits a form with valid data
//
//	return &SubmitResult{
//	    FormData: &MyFormData{...},
//	}
//
// The Intent receives this result, processes the form data, and transitions
// to the next screen or completes the workflow.
type SubmitResult struct {
	FormData interface{}
	Meta     map[string]interface{}
}

// Type returns ResultSubmit, identifying this result as a form
//
// Returns:
//   - A ScreenResultType value.
//
// Side effects:
//   - None.
func (r *SubmitResult) Type() ScreenResultType {
	return ResultSubmit
}

// Data returns FormData, the submitted form values. The caller should
//
// Returns:
//   - A interface{} value.
//
// Side effects:
//   - None.
func (r *SubmitResult) Data() interface{} {
	return r.FormData
}

// Metadata returns the key-value metadata map attached to this result.
//
// Returns:
//   - A map[string]interface{} value.
//
// Side effects:
//   - None.
func (r *SubmitResult) Metadata() map[string]interface{} {
	if r.Meta == nil {
		r.Meta = make(map[string]interface{})
	}
	return r.Meta
}

// WithMetadata stores a key-value pair in the result metadata and returns
//
// Expected:
//   - Must be a valid string.
//   - interface{} must be valid.
//
// Returns:
//   - A ScreenResult value.
//
// Side effects:
//   - None.
func (r *SubmitResult) WithMetadata(key string, value interface{}) ScreenResult {
	if r.Meta == nil {
		r.Meta = make(map[string]interface{})
	}
	r.Meta[key] = value
	return r
}

// ErrorResult is returned when an error occurs in the Screen.
//
// Example: Database query fails during screen initialization
//
//	return &ErrorResult{
//	    Err: err,
//	    Message: "Failed to load profiles",
//	}
//
// The Intent receives this result and can:
// - Display the error in a modal
// - Transition to an error state
// - Log the error and retry.
type ErrorResult struct {
	Err     error
	Message string
	Meta    map[string]interface{}
}

// Type returns ResultError, identifying this result as a failure
//
// Returns:
//   - A ScreenResultType value.
//
// Side effects:
//   - None.
func (r *ErrorResult) Type() ScreenResultType {
	return ResultError
}

// Data returns a map with "error" and "message" keys. The "error" value
//
// Returns:
//   - A interface{} value.
//
// Side effects:
//   - None.
func (r *ErrorResult) Data() interface{} {
	return map[string]interface{}{
		"error":   r.Err,
		"message": r.Message,
	}
}

// Metadata returns the key-value metadata map attached to this result.
//
// Returns:
//   - A map[string]interface{} value.
//
// Side effects:
//   - None.
func (r *ErrorResult) Metadata() map[string]interface{} {
	if r.Meta == nil {
		r.Meta = make(map[string]interface{})
	}
	return r.Meta
}

// WithMetadata stores a key-value pair in the result metadata and returns
//
// Expected:
//   - Must be a valid string.
//   - interface{} must be valid.
//
// Returns:
//   - A ScreenResult value.
//
// Side effects:
//   - None.
func (r *ErrorResult) WithMetadata(key string, value interface{}) ScreenResult {
	if r.Meta == nil {
		r.Meta = make(map[string]interface{})
	}
	r.Meta[key] = value
	return r
}
