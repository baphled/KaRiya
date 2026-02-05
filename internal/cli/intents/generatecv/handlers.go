package generatecv

import (
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/service/career/cv"
	tea "github.com/charmbracelet/bubbletea"
)

// handleWizardComplete processes wizard completion by starting tech extraction.
func (i *Intent) handleWizardComplete(msg WizardCompleteMsg) tea.Cmd {
	i.wizardModal.Complete()

	for _, p := range i.context.AvailableProfiles {
		if p.ID == msg.ProfileID {
			i.state.selectedProfile = p
			break
		}
	}
	i.state.selectedAudience = msg.Audience

	i.state.selectedTechnologyFocus = cv.TechnologyFocus(msg.TechFocus)
	i.state.selectedTechnologies = msg.Technologies
	i.state.selectedFocusArea = cv.FocusArea(msg.FocusArea)

	i.state.selectedSkillsFormat = msg.SkillsFormat
	i.state.selectedSkillsLimit = msg.SkillsLimit
	i.state.selectedCVLength = msg.CVLength

	termInfo := i.GetTerminalInfo()
	i.progressModal = components.NewExtractingTechsProgress(termInfo.Width, termInfo.Height)
	i.progressModal.Show()

	i.state.currentState = StateExtracting
	return i.extractTechnologiesAsync()
}

// handleTechExtracted processes tech extraction completion by starting CV generation.
func (i *Intent) handleTechExtracted(msg TechnologiesExtractedMsg) tea.Cmd {
	if i.progressModal != nil {
		i.progressModal.Hide()
	}

	i.state.extractedTechnologies = msg.Technologies
	i.state.focusAreaSuggestion = msg.Suggestion

	termInfo := i.GetTerminalInfo()
	profileName := "Default Profile"
	if i.state.selectedProfile != nil {
		profileName = i.state.selectedProfile.Name
	}
	i.progressModal = components.NewGeneratingCVProgress(profileName, i.state.selectedAudience, termInfo.Width, termInfo.Height)
	i.progressModal.Show()

	i.state.currentState = StateGenerating
	return i.generateCVAsync()
}

// handleCVGenerated processes CV generation completion by showing review screen.
func (i *Intent) handleCVGenerated(msg CVGenerationCompleteMsg) tea.Cmd {
	if i.progressModal != nil {
		i.progressModal.Hide()
	}

	i.state.generatedCV = msg.CV

	if i.context.ReviewScreenFactory != nil {
		i.wizardReviewScreen = i.context.ReviewScreenFactory(msg.CV)
	}
	i.state.currentState = StateReview
	return nil
}

// handleReviewScreenResult processes review screen results.
func (i *Intent) handleReviewScreenResult(result screens.ScreenResult) tea.Cmd {
	switch result.Type() {
	case screens.ResultNavigate:
		switch result.Data() {
		case "preview":
			if i.context.PreviewScreenFactory != nil {
				i.wizardPreviewScreen = i.context.PreviewScreenFactory(i.state.generatedCV)
			}
			i.state.currentState = StatePreview
			return nil

		case "export":
			return i.showExportModal()

		case "edit":
			return i.returnToWizard()
		}

	case screens.ResultCancel:
		return i.returnToWizard()
	}
	return nil
}

// handlePreviewScreenResult processes preview screen results.
func (i *Intent) handlePreviewScreenResult(result screens.ScreenResult) tea.Cmd {
	switch result.Type() {
	case screens.ResultSubmit:
		i.setCompleted()
		return nil

	case screens.ResultNavigate:
		switch result.Data() {
		case "export":
			return i.showExportModal()

		case "edit":
			return i.returnToWizard()
		}

	case screens.ResultCancel:
		i.state.currentState = StateReview
		return nil
	}
	return nil
}

// handleExportComplete processes export modal completion by starting async export.
func (i *Intent) handleExportComplete(exportData *components.ExportData) tea.Cmd {
	i.exportModal.Hide()

	termInfo := i.GetTerminalInfo()
	i.progressModal = components.NewExportingProgress(exportData.Format, termInfo.Width, termInfo.Height)
	i.progressModal.Show()

	switch exportData.Format {
	case "text":
		i.state.selectedExportFormat = ExportFormatText
	case "markdown":
		i.state.selectedExportFormat = ExportFormatMarkdown
	case "yaml":
		i.state.selectedExportFormat = ExportFormatYAML
	}

	switch exportData.Location {
	case "file":
		i.state.selectedExportOption = ExportOptionSaveToFile
	case "clipboard":
		i.state.selectedExportOption = ExportOptionClipboard
	}

	i.state.currentState = StateExporting
	return i.exportCVAsync()
}

// handleExportCompleteMsg processes the async export completion message.
func (i *Intent) handleExportCompleteMsg(msg ExportCompleteMsg) {
	i.state.isExporting = false
	if i.progressModal != nil {
		i.progressModal.Hide()
	}
	if i.exportModal != nil {
		i.exportModal.Hide()
	}
	if msg.Error != nil {
		i.state.exportError = msg.Error
		i.state.currentState = StateExportComplete
		return
	}
	i.state.exportedPath = msg.Path
	i.state.currentState = StateExportComplete
}

// handleExportCompleteKeypress handles key events in export complete state.
func (i *Intent) handleExportCompleteKeypress(keyStr string) tea.Cmd {
	switch keyStr {
	case "enter":
		now := time.Now()
		i.result = &intents.IntentResult[*Result]{
			Status: intents.Completed,
			Data: &Result{
				GeneratedCV:     i.state.generatedCV,
				SelectedProfile: i.state.selectedProfile,
				AcceptedFields:  make(map[string]bool),
				ExportPath:      i.state.exportedPath,
				CVExportFormat:  string(i.state.selectedExportFormat),
				ExportedAt:      &now,
			},
			Metadata: map[string]interface{}{
				"profile":         i.state.selectedProfile.ID,
				"audience":        i.state.selectedAudience,
				"export_format":   string(i.state.selectedExportFormat),
				"export_location": i.state.exportedPath,
			},
		}
		i.active = false
		return nil
	case "esc":
		i.state.currentState = StateExportSelectLocation
		i.state.exportError = nil
		return nil
	}
	return nil
}
