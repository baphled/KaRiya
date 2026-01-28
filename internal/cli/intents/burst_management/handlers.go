package burst_management

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/screens"
	burstscreens "github.com/baphled/kariya/internal/cli/screens/burst_management"
	"github.com/baphled/kariya/internal/cli/screens/facts"
	"github.com/baphled/kariya/internal/cli/screens/timeline"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// Ensure Intent implements ScreenResultHandler interface.
var _ behaviors.ScreenResultHandler = (*Intent)(nil)

// handleScreenResult processes a screen result and determines next action.
func (i *Intent) handleScreenResult(result interface{}) tea.Cmd {
	if result == nil {
		return nil
	}

	screenResult, ok := result.(screens.ScreenResult)
	if !ok {
		return nil
	}

	return behaviors.NewScreenResultDispatcher(i).Dispatch(screenResult)
}

// HandleCancel handles screen cancellation (back/escape).
//
//nolint:revive // result parameter required by ScreenResultHandler interface
func (i *Intent) HandleCancel(result *screens.CancelResult) tea.Cmd {
	_ = result // Satisfy interface - result data not needed for cancel handling
	switch i.state {
	case StateList:
		// Check if any modal is visible - if so, Esc should close the modal, not cancel the intent.
		// Use hasVisibleModal() to properly check visibility rather than just nil checks.
		// This prevents race conditions where async operations might leave stale modal references.
		if i.hasVisibleModal() {
			// Modal is active or loading - ignore cancel from screen.
			return nil
		}
		// No modals active - cancel from list returns to main menu.
		i.SetCancelled()
		return nil

	case StateDetail, StateDetailEvents, StateDetailFacts:
		// Return to list from detail views.
		i.state = StateList
		i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
		return nil

	case StateDeleteConfirm:
		// Cancel delete returns to list.
		i.state = StateList
		i.deleteModal = nil
		i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
		return nil

	default:
		// Default: return to list.
		i.state = StateList
		i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
		return nil
	}
}

// HandleNavigate handles screen navigation results.
func (i *Intent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
	// Handle action data (add, edit, delete, suggest).
	if actionData, ok := result.ResultData.(map[string]interface{}); ok {
		return i.handleActionData(actionData)
	}

	// Handle burst selection (view details).
	if burst, ok := result.ResultData.(*career.Burst); ok {
		i.selectedBurst = burst
		i.AddViewedBurst(burst)
		// Show detail modal instead of transitioning to detail screen.
		return i.showBurstDetailModal(burst)
	}

	return nil
}

// HandleSubmit handles form submission results.
//
//nolint:revive // result parameter required by ScreenResultHandler interface
func (i *Intent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
	_ = result // Not yet implemented - form submissions handled via modals
	return nil
}

// HandleError handles error results from screens.
func (i *Intent) HandleError(result *screens.ErrorResult) tea.Cmd {
	// Store error and show error modal.
	i.deleteError = result.Err
	i.errorModal = feedback.NewErrorModal("Operation Failed", result.Err.Error())
	return nil
}

// handleActionData processes action data from navigation results.
//
//nolint:funlen // Action routing requires handling multiple action types in one function
func (i *Intent) handleActionData(actionData map[string]interface{}) tea.Cmd {
	//nolint:errcheck // Type assertion ok value intentionally ignored - empty string is acceptable default
	action, _ := actionData["action"].(string)
	switch action {
	case "edit":
		// Extract burst from action data.
		if burst, ok := actionData["burst"].(*career.Burst); ok {
			i.selectedBurst = burst
			i.state = StateEdit
			// Open edit modal using helper.
			return i.openEditModal(burst)
		}
		return nil

	case "delete":
		// Extract burst from action data.
		if burst, ok := actionData["burst"].(*career.Burst); ok {
			i.selectedBurst = burst
			// Create delete confirmation modal.
			burstName := burst.Name
			if len(burstName) > 50 {
				burstName = burstName[:47] + "..."
			}
			i.deleteModal = feedback.NewConfirmModal(
				"Delete Burst",
				fmt.Sprintf("Are you sure you want to delete '%s'?", burstName),
			).WithVariant(feedback.ConfirmDestructive)
			i.state = StateDeleteConfirm
			return i.deleteModal.Init()
		}
		return nil

	case "suggest":
		// Trigger burst suggestion AI detection.
		i.state = StateSuggesting
		// Load all events and trigger burst detection.
		return i.startBurstDetection()

	case "view_events":
		// View events in burst.
		if burst, ok := actionData["burst"].(*career.Burst); ok {
			i.selectedBurst = burst
			i.state = StateDetailEvents
			// Load events for this burst.
			i.burstEvents = i.loadBurstEvents(burst)
			// Use timeline.TimelineEventListScreen to display burst events.
			i.transitionToScreen(timeline.NewTimelineEventListScreen(i.burstEvents))
		}
		return nil

	case "view_facts":
		// View facts extracted from burst.
		if burst, ok := actionData["burst"].(*career.Burst); ok {
			i.selectedBurst = burst
			i.state = StateDetailFacts
			// Load facts for this burst.
			i.burstFacts = i.loadBurstFacts(burst)
			// Use facts.FactListScreen to display burst facts.
			i.transitionToScreen(facts.NewFactListScreen(i.burstFacts))
		}
		return nil

	case "confirm":
		// Confirm burst (mark as confirmed).
		if burst, ok := actionData["burst"].(*career.Burst); ok {
			i.selectedBurst = burst
			// Show confirmation modal.
			burstName := burst.Name
			if len(burstName) > 50 {
				burstName = burstName[:47] + "..."
			}
			i.confirmModal = feedback.NewConfirmModal(
				"Confirm Burst",
				fmt.Sprintf("Mark '%s' as confirmed? This indicates the burst is validated and complete.", burstName),
			)
			i.state = StateConfirm
			return i.confirmModal.Init()
		}
		return nil

	default:
		return nil
	}
}
