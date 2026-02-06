package generatecv

import (
	"time"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/cv"
	cvmodals "github.com/baphled/kariya/internal/cli/screens/cv/modals"
	"github.com/baphled/kariya/internal/domain/career"
	cvsvc "github.com/baphled/kariya/internal/service/career/cv"
	tea "github.com/charmbracelet/bubbletea"
)

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

// HandleCancel processes screen cancellation by navigating back to the
// previous state or marking the intent as cancelled when at root.
//
// Expected:
//   - result must be non-nil.
//
// Returns:
//   - A tea.Cmd or nil.
//
// Side effects:
//   - May transition state or mark intent as cancelled.
func (i *Intent) HandleCancel(result *screens.CancelResult) tea.Cmd {
	switch i.state {
	case StateReview:
		return i.returnToWizard()

	case StatePreview:
		i.state = StateReview
		i.transitionToReviewScreen()
		return nil

	case StateConfiguring:
		i.setCancelled()
		return nil

	default:
		i.setCancelled()
		return nil
	}
}

// HandleNavigate dispatches navigation results based on the data type.
//
// Expected:
//   - result must be non-nil with a populated ResultData field.
//
// Returns:
//   - A tea.Cmd from the dispatched handler, or nil.
//
// Side effects:
//   - May open modals or transition screens.
func (i *Intent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
	switch data := result.ResultData.(type) {
	case string:
		return i.handleNavigateString(data)
	case *career.CVView:
		// CV confirmation from preview screen
		i.setCompleted()
		return nil
	default:
		return nil
	}
}

// handleNavigateString handles string-based navigation commands.
func (i *Intent) handleNavigateString(action string) tea.Cmd {
	switch action {
	case "preview":
		i.state = StatePreview
		i.transitionToPreviewScreen()
		return nil

	case "export":
		return i.showExportModal()

	case "edit":
		return i.returnToWizard()

	default:
		return nil
	}
}

// HandleSubmit processes form submissions - confirms the CV.
//
// Expected:
//   - result must be non-nil.
//
// Returns:
//   - Always nil.
//
// Side effects:
//   - Marks intent as completed.
func (i *Intent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
	i.setCompleted()
	return nil
}

// HandleError captures errors surfaced by screens.
//
// Expected:
//   - result must be non-nil.
//
// Returns:
//   - Always nil.
//
// Side effects:
//   - Stores the error in exportError field.
func (i *Intent) HandleError(result *screens.ErrorResult) tea.Cmd {
	i.exportError = result.Err
	return nil
}

// handleWizardComplete processes wizard completion by starting tech extraction.
func (i *Intent) handleWizardComplete(msg WizardCompleteMsg) tea.Cmd {
	i.wizardModal.Complete()

	for _, p := range i.context.AvailableProfiles {
		if p.ID == msg.ProfileID {
			i.selectedProfile = p
			break
		}
	}
	i.selectedAudience = msg.Audience

	i.selectedTechnologyFocus = cvsvc.TechnologyFocus(msg.TechFocus)
	i.selectedTechnologies = msg.Technologies
	i.selectedFocusArea = cvsvc.FocusArea(msg.FocusArea)

	i.selectedSkillsFormat = msg.SkillsFormat
	i.selectedSkillsLimit = msg.SkillsLimit
	i.selectedCVLength = msg.CVLength

	termInfo := i.GetTerminalInfo()
	i.progressModal = cvmodals.NewExtractingTechsProgress(termInfo.Width, termInfo.Height)
	i.progressModal.Show()

	i.state = StateExtracting
	return i.extractTechnologiesAsync()
}

// handleTechExtracted processes tech extraction completion by starting CV generation.
func (i *Intent) handleTechExtracted(msg TechnologiesExtractedMsg) tea.Cmd {
	if i.progressModal != nil {
		i.progressModal.Hide()
	}

	i.extractedTechnologies = msg.Technologies
	i.focusAreaSuggestion = msg.Suggestion

	termInfo := i.GetTerminalInfo()
	profileName := "Default Profile"
	if i.selectedProfile != nil {
		profileName = i.selectedProfile.Name
	}
	i.progressModal = cvmodals.NewGeneratingCVProgress(profileName, i.selectedAudience, termInfo.Width, termInfo.Height)
	i.progressModal.Show()

	i.state = StateGenerating
	return i.generateCVAsync()
}

// handleCVGenerated processes CV generation completion by showing review screen.
func (i *Intent) handleCVGenerated(msg CVGenerationCompleteMsg) tea.Cmd {
	if i.progressModal != nil {
		i.progressModal.Hide()
	}

	i.generatedCV = msg.CV
	i.state = StateReview
	i.transitionToReviewScreen()
	return nil
}

// transitionToReviewScreen creates and activates the review screen.
func (i *Intent) transitionToReviewScreen() {
	i.reviewScreen = cv.NewCVReviewScreenWithProfile(i.generatedCV, i.context.ProfileConfig)
	i.activeScreen = i.reviewScreen
}

// transitionToPreviewScreen creates and activates the preview screen.
func (i *Intent) transitionToPreviewScreen() {
	i.previewScreen = cv.NewCVPreviewScreenWithProfile(i.generatedCV, i.context.ProfileConfig)
	i.activeScreen = i.previewScreen
}

// handleExportComplete processes export modal completion by starting async export.
func (i *Intent) handleExportComplete(exportData *cvmodals.ExportData) tea.Cmd {
	i.exportModal.Hide()

	termInfo := i.GetTerminalInfo()
	i.progressModal = cvmodals.NewExportingProgress(exportData.Format, termInfo.Width, termInfo.Height)
	i.progressModal.Show()

	switch exportData.Format {
	case "text":
		i.selectedExportFormat = ExportFormatText
	case "markdown":
		i.selectedExportFormat = ExportFormatMarkdown
	case "yaml":
		i.selectedExportFormat = ExportFormatYAML
	}

	switch exportData.Location {
	case "file":
		i.selectedExportOption = ExportOptionSaveToFile
	case "clipboard":
		i.selectedExportOption = ExportOptionClipboard
	}

	i.state = StateExporting
	return i.exportCVAsync()
}

// handleExportCompleteMsg processes the async export completion message.
func (i *Intent) handleExportCompleteMsg(msg ExportCompleteMsg) {
	i.isExporting = false
	if i.progressModal != nil {
		i.progressModal.Hide()
	}
	if i.exportModal != nil {
		i.exportModal.Hide()
	}
	if msg.Error != nil {
		i.exportError = msg.Error
		i.state = StateExportComplete
		return
	}
	i.exportedPath = msg.Path
	i.state = StateExportComplete
}

// handleExportCompleteKeypress handles key events in export complete state.
func (i *Intent) handleExportCompleteKeypress(keyStr string) tea.Cmd {
	switch keyStr {
	case "enter":
		now := time.Now()
		i.result = &intents.IntentResult[*Result]{
			Status: intents.Completed,
			Data: &Result{
				GeneratedCV:     i.generatedCV,
				SelectedProfile: i.selectedProfile,
				AcceptedFields:  make(map[string]bool),
				ExportPath:      i.exportedPath,
				CVExportFormat:  string(i.selectedExportFormat),
				ExportedAt:      &now,
			},
			Metadata: map[string]interface{}{
				"profile":         i.selectedProfile.ID,
				"audience":        i.selectedAudience,
				"export_format":   string(i.selectedExportFormat),
				"export_location": i.exportedPath,
			},
		}
		i.active = false
		return nil
	case "esc":
		i.state = StateExportSelectLocation
		i.exportError = nil
		return nil
	}
	return nil
}
