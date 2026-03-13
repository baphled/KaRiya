package harness

import (
	"github.com/baphled/kariya/internal/tui/intents/browsetimeline"
	"github.com/baphled/kariya/internal/tui/intents/burst_management"
	"github.com/baphled/kariya/internal/tui/intents/captureevent"
	configure "github.com/baphled/kariya/internal/tui/intents/configure"
	"github.com/baphled/kariya/internal/tui/intents/factmanagement"
	"github.com/baphled/kariya/internal/tui/intents/generatecv"
	"github.com/baphled/kariya/internal/tui/intents/skillsmanagement"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	tea "github.com/charmbracelet/bubbletea"
)

// processCmdResult processes a message returned from a command.
// Only essential state transition messages are processed.
//
// Expected:
//   - msg must be a valid tea.Msg.
//
// Returns:
//   - None.
//
// Side effects:
//   - May update the Model field of the TestEnv.
//   - May recursively execute commands.
//
//nolint:gocyclo // Type switch handler inherently requires many cases for different message types.
func (e *TestEnv) processCmdResult(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.QuitMsg:
		e.QuitRequested = true
	case tea.BatchMsg:
		e.processBatchMsg(msg)
	case captureevent.SubmitMsg:
		e.updateModelAndExecute(msg)
	case captureevent.SubmitCompleteMsg:
		e.processSubmitCompleteMsg(msg)
	case captureevent.InferenceCompleteMsg:
		e.updateModelAndExecute(msg)
	case captureevent.PostSavePersistenceCompleteMsg:
		e.updateModelAndExecute(msg)
	case captureevent.SubmitErrorMsg:
		e.updateModelAndExecute(msg)
	case captureevent.DismissModalMsg:
		e.updateModelAndExecute(msg)
	case feedback.ModalAutoDismissMsg:
		e.updateModelAndExecute(msg)
	case configure.ConfigCompleteMsg:
		e.updateModelAndExecute(msg)
	case generatecv.TechnologiesExtractedMsg:
		e.updateModelAndExecute(msg)
	case generatecv.CVGenerationCompleteMsg:
		e.updateModelAndExecute(msg)
	case generatecv.ExportCompleteMsg:
		e.updateModelAndExecute(msg)
	case skillsmanagement.SkillsLoadedMsg:
		e.updateModelAndExecute(msg)
	case skillsmanagement.SkillCreatedMsg:
		e.updateModelAndExecute(msg)
	case skillsmanagement.SkillUpdatedMsg:
		e.updateModelAndExecute(msg)
	case skillsmanagement.SkillDeletedMsg:
		e.updateModelAndExecute(msg)
	case skillsmanagement.SkillFormCompleteMsg:
		e.updateModelAndExecute(msg)
	case skillsmanagement.SkillEventsForModalLoadedMsg:
		e.updateModelAndExecute(msg)
	case skillsmanagement.SkillSuggestionsLoadedMsg:
		e.updateModelAndExecute(msg)
	case skillsmanagement.SkillsCreatedMsg:
		e.updateModelAndExecute(msg)
	case browsetimeline.SkillsForModalLoadedMsg:
		e.updateModelAndExecute(msg)
	case browsetimeline.SkillPickerDataLoadedMsg:
		e.updateModelAndExecute(msg)
	case browsetimeline.SkillsRefreshedMsg:
		e.updateModelAndExecute(msg)
	case burst_management.EditBurstMsg:
		e.updateModelAndExecute(msg)
	case burst_management.BurstEditCompleteMsg:
		e.updateModelAndExecute(msg)
	case burst_management.BurstDeletedMsg:
		e.updateModelAndExecute(msg)
	case burst_management.BurstConfirmedMsg:
		e.updateModelAndExecute(msg)
	case burst_management.BurstEventsLoadedMsg:
		e.updateModelAndExecute(msg)
	case burst_management.BurstFactsLoadedMsg:
		e.updateModelAndExecute(msg)
	case burst_management.BurstSkillsLoadedMsg:
		e.updateModelAndExecute(msg)
	case burst_management.BurstSuggestionsLoadedMsg:
		e.updateModelAndExecute(msg)
	case burst_management.SkillSuggestionsLoadedMsg:
		e.updateModelAndExecute(msg)
	case burst_management.FactExtractionCompleteMsg:
		e.updateModelAndExecute(msg)
	case burst_management.BurstSavedMsg:
		e.updateModelAndExecute(msg)
	case burst_management.ConfirmBurstFactsLoadedMsg:
		e.updateModelAndExecute(msg)
	case factmanagement.FactsLoadedMsg:
		e.updateModelAndExecute(msg)
	case factmanagement.FactSavedMsg:
		e.updateModelAndExecute(msg)
	case factmanagement.FactDeletedMsg:
		e.updateModelAndExecute(msg)
	case feedback.ModalCountdownTickMsg:
		e.updateModelAndExecute(msg)
	}
}

// processBatchMsg processes a batch of commands.
//
// Expected:
//   - msg must be a valid tea.BatchMsg.
//
// Returns:
//   - None.
//
// Side effects:
//   - Executes each non-nil command in the batch.
func (e *TestEnv) processBatchMsg(msg tea.BatchMsg) {
	for _, cmd := range msg {
		if cmd != nil {
			e.executeCmd(cmd)
		}
	}
}

// processSubmitCompleteMsg handles submit completion.
// The modal is NOT auto-dismissed; test steps handle dismissal explicitly.
//
// Expected:
//   - msg must be a valid captureevent.SubmitCompleteMsg.
//
// Returns:
//   - None.
//
// Side effects:
//   - Updates the Model with the submit completion message.
func (e *TestEnv) processSubmitCompleteMsg(msg captureevent.SubmitCompleteMsg) {
	e.updateModelAndExecute(msg)
}
