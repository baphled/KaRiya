package generatecv

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	cvsvc "github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/tui/intents"
	cvviews "github.com/baphled/kariya/internal/tui/views/cv"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

// handleViewResult processes a view result and determines next action.
func (i *Intent) handleViewResult(result widgets.ViewResult) tea.Cmd {
	if result == nil {
		return nil
	}

	switch result.(type) {
	case *widgets.CancelViewResult:
		return i.HandleCancel(result)
	case *widgets.NavigateViewResult:
		return i.HandleNavigate(result)
	case *widgets.SubmitViewResult:
		return i.HandleSubmit(result)
	case *widgets.ErrorViewResult:
		return i.HandleError(result)
	default:
		return nil
	}
}

// HandleCancel processes view cancellation by navigating back to the
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
func (i *Intent) HandleCancel(result widgets.ViewResult) tea.Cmd {
	switch i.state {
	case StateReview:
		return i.returnToWizard()

	case StatePreview:
		i.state = StateReview
		i.transitionToReviewView()
		return nil

	default:
		i.setCancelled()
		return nil
	}
}

// HandleNavigate dispatches navigation results based on the data type.
//
// Expected:
//   - result must be non-nil with a populated Data field.
//
// Returns:
//   - A tea.Cmd from the dispatched handler, or nil.
//
// Side effects:
//   - May open modals or transition views.
func (i *Intent) HandleNavigate(result widgets.ViewResult) tea.Cmd {
	data := result.Data()

	if nav, ok := data.(cvviews.Nav); ok {
		switch nav.Action {
		case cvviews.ActionPreview:
			i.state = StatePreview
			i.transitionToPreviewView()
			return nil
		case cvviews.ActionExport:
			return i.showExportModal()
		case cvviews.ActionEdit:
			return i.returnToWizard()
		case cvviews.ActionConfirm:
			i.setCompleted()
			return nil
		}
	}

	return nil
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
func (i *Intent) HandleSubmit(result widgets.ViewResult) tea.Cmd {
	i.setCompleted()
	return nil
}

// HandleError captures errors surfaced by views.
//
// Expected:
//   - result must be non-nil.
//
// Returns:
//   - Always nil.
//
// Side effects:
//   - Stores the error in exportError field.
func (i *Intent) HandleError(result widgets.ViewResult) tea.Cmd {
	if errResult, ok := result.(*widgets.ErrorViewResult); ok {
		i.exportError = errResult.Err
	}
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
	i.progressModal = cvviews.NewExtractingTechsProgress(termInfo.Width, termInfo.Height)
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
	i.progressModal = cvviews.NewGeneratingCVProgress(profileName, i.selectedAudience, termInfo.Width, termInfo.Height)
	i.progressModal.Show()

	i.state = StateGenerating
	return i.generateCVAsync()
}

// handleCVGenerated processes CV generation completion by showing review view.
func (i *Intent) handleCVGenerated(msg CVGenerationCompleteMsg) tea.Cmd {
	if i.progressModal != nil {
		i.progressModal.Hide()
	}

	i.generatedCV = msg.CV
	i.state = StateReview
	i.transitionToReviewView()
	return nil
}

// transitionToReviewView creates and activates the review view.
func (i *Intent) transitionToReviewView() {
	summary := i.buildGenerationSummary()
	i.reviewView = cvviews.NewReview(display.CVViewFromDomain(i.generatedCV), i.context.ProfileConfig, summary)
	i.activeView = i.reviewView
}

// buildGenerationSummary constructs a GenerationSummary from intent state.
func (i *Intent) buildGenerationSummary() *cvviews.GenerationSummary {
	summary := &cvviews.GenerationSummary{
		SelectedAudience: i.selectedAudience,
		TechnologyFocus:  string(i.selectedTechnologyFocus),
		Technologies:     i.selectedTechnologies,
		FocusArea:        string(i.selectedFocusArea),
		SkillsFormat:     i.selectedSkillsFormat,
		SkillsLimit:      i.selectedSkillsLimit,
		CVLength:         i.selectedCVLength,
	}
	if i.selectedProfile != nil {
		summary.SelectedProfile = i.selectedProfile
	}
	if i.generatedCV != nil {
		summary.SectionCount = len(i.generatedCV.Sections)
		summary.SourceEventCount = i.generatedCV.SourceEventCount
		summary.SourceFactCount = i.generatedCV.SourceFactCount
		summary.TotalBullets = countCVBullets(display.CVSectionsFromDomain(i.generatedCV.Sections))
	}
	return summary
}

// countCVBullets counts the total number of bullets across all sections.
func countCVBullets(sections []display.CVSection) int {
	count := 0
	for sectionIdx := range sections {
		section := sections[sectionIdx]
		for groupIdx := range section.Content {
			group := section.Content[groupIdx]
			count += len(group.Bullets)
		}
	}
	return count
}

// transitionToPreviewView creates and activates the preview view.
func (i *Intent) transitionToPreviewView() {
	i.previewView = cvviews.NewPreview(display.CVViewFromDomain(i.generatedCV), i.context.ProfileConfig)
	i.activeView = i.previewView
}

// handleExportComplete processes export modal completion by starting async export.
func (i *Intent) handleExportComplete(exportData *cvviews.ExportData) tea.Cmd {
	i.exportModal.Hide()

	termInfo := i.GetTerminalInfo()
	i.progressModal = cvviews.NewExportingProgress(exportData.Format, termInfo.Width, termInfo.Height)
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
