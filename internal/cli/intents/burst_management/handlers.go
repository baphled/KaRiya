package burst_management

import (
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	tea "github.com/charmbracelet/bubbletea"
)

// Ensure Intent implements ScreenResultHandler interface.
var _ intents.ScreenResultHandler = (*Intent)(nil)

// handleScreenResult processes a screen result and determines next action.
func (i *Intent) handleScreenResult(result interface{}) tea.Cmd {
	if result == nil {
		return nil
	}

	screenResult, ok := result.(screens.ScreenResult)
	if !ok {
		return nil
	}

	return intents.NewScreenResultDispatcher(i).Dispatch(screenResult)
}

// HandleCancel handles screen cancellation (back/escape).
func (i *Intent) HandleCancel(result *screens.CancelResult) tea.Cmd {
	switch i.state {
	case StateList:
		// Cancel from list returns to main menu
		i.SetCancelled()
		return nil

	case StateDetail, StateDetailEvents, StateDetailFacts:
		// Return to list from detail views
		i.state = StateList
		// TODO: Transition to list screen
		return nil

	case StateDeleteConfirm:
		// Cancel delete returns to list
		i.state = StateList
		// TODO: Transition to list screen
		return nil

	default:
		// Default: return to list
		i.state = StateList
		return nil
	}
}

// HandleNavigate handles screen navigation results.
func (i *Intent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
	// Handle action data (add, edit, delete)
	if actionData, ok := result.ResultData.(map[string]interface{}); ok {
		return i.handleActionData(actionData)
	}

	// Handle burst selection (view details)
	// TODO: Add burst selection handling when screens are implemented

	return nil
}

// HandleSubmit handles form submission results.
func (i *Intent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
	// TODO: Implement submit handling when forms are added
	return nil
}

// HandleError handles error results from screens.
func (i *Intent) HandleError(result *screens.ErrorResult) tea.Cmd {
	i.deleteError = result.Err
	// TODO: Show error modal
	return nil
}

// handleActionData processes action data from navigation results.
func (i *Intent) handleActionData(actionData map[string]interface{}) tea.Cmd {
	action, _ := actionData["action"].(string)
	switch action {
	case "add":
		// TODO: Open create burst modal
		return nil

	case "edit":
		// TODO: Open edit burst modal
		return nil

	case "delete":
		// TODO: Open delete confirmation modal
		return nil

	case "suggest":
		// TODO: Trigger burst suggestion AI detection
		return nil

	default:
		return nil
	}
}
