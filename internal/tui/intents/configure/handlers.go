package configure

import (
	"github.com/baphled/kariya/internal/tui/intents"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
)

// HandleNavigate processes a navigation result from view interactions.
//
// Expected: result must be a valid ViewResult.
// Returns: A tea.Cmd value.
//
// Side effects: None.
func (i *Intent) HandleNavigate(result widgets.ViewResult) tea.Cmd {
	nav, ok := result.(*widgets.NavigateViewResult)
	if !ok {
		return nil
	}

	domain, ok := nav.ResultData.(ConfigurationDomain)
	if !ok {
		return nil
	}

	i.selectedDomain = domain
	return nil
}

// HandleCancel processes a cancellation from view interactions.
//
// Expected: result must be a valid ViewResult.
// Returns: A tea.Cmd value.
//
// Side effects: None.
func (i *Intent) HandleCancel(_ widgets.ViewResult) tea.Cmd {
	i.setCancelled()
	return nil
}

// HandleSubmit processes submission results from views (no-op for configure workflow).
//
// Expected: result must be a valid ViewResult.
// Returns: A tea.Cmd value.
//
// Side effects: None.
func (i *Intent) HandleSubmit(_ widgets.ViewResult) tea.Cmd {
	return nil
}

// HandleError processes error results from views (no-op for configure workflow).
//
// Expected: result must be a valid ViewResult.
// Returns: A tea.Cmd value.
//
// Side effects: None.
func (i *Intent) HandleError(_ widgets.ViewResult) tea.Cmd {
	return nil
}

// updateSettingsModal handles updates to the unified settings modal.
func (i *Intent) updateSettingsModal(msg tea.Msg) tea.Cmd {
	cmd := i.settingsModal.Update(msg)

	if i.settingsModal.IsCompleted() {
		changes := i.settingsModal.GetChanges()
		i.pendingChanges = changes.Changes
		i.settingsModal = nil
		i.state = ConfigStateSaving
		return i.startSaving()
	}
	if i.settingsModal.IsCancelled() {
		i.settingsModal = nil
		i.setCancelled()
		return nil
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
				i.openSettingsModal()
				return i.settingsModal.Init()
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
func (i *Intent) handleAsyncCompletion(_ tea.Msg) tea.Cmd {
	return nil
}
