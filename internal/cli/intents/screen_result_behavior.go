package intents

import (
	"github.com/baphled/kariya/internal/cli/screens"
	tea "github.com/charmbracelet/bubbletea"
)

// ScreenResultHandler defines how an intent handles screen results.
//
// This interface extracts the common pattern of handling screen results
// across all screen-based intents (GenerateCV, ManageSkills, BrowseTimeline).
//
// Benefits:
//   - Eliminates repetitive type switching (20+ lines → 1 line)
//   - Forces implementation of all result handlers (compile-time safety)
//   - Clear contract for screen result handling
//   - Easy to test (mock interface)
//
// Usage Pattern:
//
//	// 1. Implement the interface on your intent
//	var _ ScreenResultHandler = (*YourIntent)(nil)
//
//	func (i *YourIntent) HandleNavigate(result *screens.NavigateResult) tea.Cmd { ... }
//	func (i *YourIntent) HandleCancel(result *screens.CancelResult) tea.Cmd { ... }
//	func (i *YourIntent) HandleSubmit(result *screens.SubmitResult) tea.Cmd { ... }
//	func (i *YourIntent) HandleError(result *screens.ErrorResult) tea.Cmd { ... }
//
//	// 2. Create dispatcher in Update()
//	dispatcher := NewScreenResultDispatcher(i)
//
//	// 3. Replace handleScreenResult() with single line
//	return dispatcher.Dispatch(result)
type ScreenResultHandler interface {
	// HandleNavigate is called when the screen returns a NavigateResult.
	// This typically means the user selected something and wants to proceed.
	//
	// Example: User selects a profile from a list.
	// The intent should extract the data and transition to the next state.
	HandleNavigate(result *screens.NavigateResult) tea.Cmd

	// HandleCancel is called when the screen returns a CancelResult.
	// This typically means the user pressed Escape.
	//
	// Example: User presses Escape on a form.
	// The intent should transition back to the previous state.
	HandleCancel(result *screens.CancelResult) tea.Cmd

	// HandleSubmit is called when the screen returns a SubmitResult.
	// This typically means the user submitted a form or completed an action.
	//
	// Example: User submits a form with valid data.
	// The intent should process the data and transition to the next state.
	HandleSubmit(result *screens.SubmitResult) tea.Cmd

	// HandleError is called when the screen returns an ErrorResult.
	// This typically means an error occurred in the screen.
	//
	// Example: Database query fails during screen initialization.
	// The intent should display the error or transition to an error state.
	HandleError(result *screens.ErrorResult) tea.Cmd
}

// ScreenResultDispatcher routes screen results to the appropriate handler methods.
//
// This eliminates the need for repetitive type switching in every intent:
//
// Before (20+ lines per intent):
//
//	func (i *YourIntent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
//	    switch r := result.(type) {
//	    case *screens.NavigateResult:
//	        return i.handleNavigateResult(r)
//	    case *screens.CancelResult:
//	        return i.handleCancelResult(r)
//	    case *screens.SubmitResult:
//	        return i.handleSubmitResult(r)
//	    case *screens.ErrorResult:
//	        return i.handleErrorResult(r)
//	    default:
//	        return nil
//	    }
//	}
//
// After (1 line):
//
//	func (i *YourIntent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
//	    return NewScreenResultDispatcher(i).Dispatch(result)
//	}
//
// Or even simpler, replace the entire handleScreenResult method with inline calls:
//
//	case ScreenResultMsg:
//	    return NewScreenResultDispatcher(i).Dispatch(msg.Result)
type ScreenResultDispatcher struct {
	handler ScreenResultHandler
}

// NewScreenResultDispatcher creates a new dispatcher for the given handler.
//
// Usage:
//
//	dispatcher := NewScreenResultDispatcher(i) // where i implements ScreenResultHandler
//	return dispatcher.Dispatch(result)
func NewScreenResultDispatcher(handler ScreenResultHandler) *ScreenResultDispatcher {
	return &ScreenResultDispatcher{
		handler: handler,
	}
}

// Dispatch routes the screen result to the appropriate handler method.
//
// Returns:
//   - The tea.Cmd returned by the handler method
//   - nil if result is nil or unknown type
//
// Type Safety:
//   - NavigateResult → HandleNavigate
//   - CancelResult → HandleCancel
//   - SubmitResult → HandleSubmit
//   - ErrorResult → HandleError
//   - Unknown types → nil (safe fallback)
func (d *ScreenResultDispatcher) Dispatch(result screens.ScreenResult) tea.Cmd {
	if result == nil {
		return nil
	}

	// Type switch to route to appropriate handler
	switch r := result.(type) {
	case *screens.NavigateResult:
		return d.handler.HandleNavigate(r)

	case *screens.CancelResult:
		return d.handler.HandleCancel(r)

	case *screens.SubmitResult:
		return d.handler.HandleSubmit(r)

	case *screens.ErrorResult:
		return d.handler.HandleError(r)

	default:
		// Unknown result type - safe fallback
		return nil
	}
}
