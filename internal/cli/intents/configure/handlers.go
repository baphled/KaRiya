package configure

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	tea "github.com/charmbracelet/bubbletea"
)

// Compile-time check that Intent implements ScreenResultHandler.
var _ behaviors.ScreenResultHandler = (*Intent)(nil)

// HandleNavigate processes a navigation result from the domain selection screen.
//
// Expected: navigateresult must be valid.
// Returns: A tea.Cmd value.
//
// Side effects: None.
func (i *Intent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
	domain, ok := result.ResultData.(ConfigurationDomain)
	if !ok {
		return nil
	}

	i.selectedDomain = domain
	i.state = ConfigStateEditSettings
	i.openEditModal()
	return i.editModal.Init()
}

// HandleCancel processes a cancellation from the domain selection screen.
//
// Expected: cancelresult must be valid.
// Returns: A tea.Cmd value.
//
// Side effects: None.
func (i *Intent) HandleCancel(_ *screens.CancelResult) tea.Cmd {
	i.setCancelled()
	return nil
}

// HandleSubmit processes submission results from screens (no-op for configure workflow).
//
// Expected: submitresult must be valid.
// Returns: A tea.Cmd value.
//
// Side effects: None.
func (i *Intent) HandleSubmit(_ *screens.SubmitResult) tea.Cmd {
	return nil
}

// HandleError processes error results from screens (no-op for configure workflow).
//
// Expected: errorresult must be valid.
// Returns: A tea.Cmd value.
//
// Side effects: None.
func (i *Intent) HandleError(_ *screens.ErrorResult) tea.Cmd {
	return nil
}

// handleScreenResult dispatches a screen result to the appropriate handler method.
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

// updateDomainScreen handles updates to the domain selection screen.
func (i *Intent) updateDomainScreen(msg tea.Msg) tea.Cmd {
	cmd, result := i.domainScreen.Update(msg)

	if result != nil {
		if handlerCmd := i.handleScreenResult(result); handlerCmd != nil {
			return handlerCmd
		}
	}

	return cmd
}

// updateEditModal handles updates to the edit settings modal.
func (i *Intent) updateEditModal(msg tea.Msg) tea.Cmd {
	cmd := i.editModal.Update(msg)

	if !i.editModal.IsVisible() {
		if i.editModal.IsCompleted() {
			i.pendingChanges = i.editModal.GetChanges()
			i.editModal = nil
			i.state = ConfigStateReviewChanges
			i.openReviewModal()
			return nil
		}
		if i.editModal.IsCancelled() {
			i.editModal = nil
			i.state = ConfigStateSelectDomain
			return nil
		}
	}

	return cmd
}

// updateReviewModal handles updates to the review changes modal.
func (i *Intent) updateReviewModal(msg tea.Msg) tea.Cmd {
	cmd := i.reviewModal.Update(msg)

	if !i.reviewModal.IsVisible() {
		if i.reviewModal.IsConfirmed() {
			i.reviewModal = nil
			i.state = ConfigStateConfirm
			i.openConfirmModal()
			return nil
		}
		if i.reviewModal.IsCancelled() {
			i.reviewModal = nil
			i.state = ConfigStateEditSettings
			i.openEditModal()
			return i.editModal.Init()
		}
	}

	return cmd
}

// updateConfirmModal handles updates to the confirmation modal.
func (i *Intent) updateConfirmModal(msg tea.Msg) tea.Cmd {
	cmd := i.confirmModal.Update(msg)

	if !i.confirmModal.IsVisible() {
		if i.confirmModal.IsConfirmed() {
			i.confirmModal = nil
			i.state = ConfigStateSaving
			return i.startSaving()
		}
		if i.confirmModal.IsCancelled() {
			i.confirmModal = nil
			i.state = ConfigStateReviewChanges
			i.openReviewModal()
			return nil
		}
	}

	return cmd
}

// updateSavingModal handles updates during the saving phase.
func (i *Intent) updateSavingModal(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case ConfigCompleteMsg:
		i.savingModal = nil
		i.configResult = msg.Result
		i.state = ConfigStateComplete
		i.settings = settingsFromConfig(i.cfg)
		i.resultModal = feedback.NewSuccessModal("Configuration saved!")
		i.result = &intents.IntentResult[interface{}]{
			Status: intents.Completed,
			Data:   msg.Result,
		}
		return i.resultModal.Init()

	case ConfigErrorMsg:
		i.savingModal = nil
		i.state = ConfigStateFailed
		i.resultModal = feedback.NewErrorModal("Save Failed", msg.Error.Message)
		return nil

	case feedback.ModalSpinnerTickMsg:
		return i.savingModal.Update(msg)

	default:
		return nil
	}
}

// updateResultModal handles updates to the result modal (success/error).
func (i *Intent) updateResultModal(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "esc", "q", " ":
			i.resultModal = nil
			if i.configResult == nil || !i.configResult.Success {
				i.state = ConfigStateEditSettings
				i.openEditModal()
				return i.editModal.Init()
			}
			i.active = false
		}
	case feedback.ModalCountdownTickMsg:
		if i.resultModal != nil && i.resultModal.Type == feedback.ModalSuccess {
			return i.resultModal.Update(msg)
		}
	case feedback.ModalAutoDismissMsg:
		i.resultModal = nil
		i.active = false
	}

	return nil
}

// handleAsyncCompletion processes async save completion messages at intent level.
func (i *Intent) handleAsyncCompletion(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case ConfigCompleteMsg:
		i.clearAllModals()
		i.settings = settingsFromConfig(i.cfg)
		i.configResult = msg.Result
		i.state = ConfigStateComplete
		i.resultModal = feedback.NewSuccessModal("Configuration saved!")
		i.result = &intents.IntentResult[interface{}]{
			Status: intents.Completed,
			Data:   msg.Result,
		}
		return nil
	case ConfigErrorMsg:
		i.clearAllModals()
		i.state = ConfigStateFailed
		i.resultModal = feedback.NewErrorModal("Save Failed", msg.Error.Message)
		return nil
	}
	return nil
}
