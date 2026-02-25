package generatecv

import (
	"errors"
	"fmt"
	"time"

	cvmodals "github.com/baphled/kariya/internal/cli/screens/cv/modals"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/service/career/technology"
	tea "github.com/charmbracelet/bubbletea"
)

func (i *Intent) initWizardFlow() tea.Cmd {
	termInfo := i.GetTerminalInfo()

	profileOptions := make([]cvmodals.ProfileOption, len(i.context.AvailableProfiles))
	for idx, profile := range i.context.AvailableProfiles {
		profileOptions[idx] = cvmodals.ProfileOption{
			ID:   profile.ID,
			Name: profile.Name,
		}
	}

	i.wizardModal = cvmodals.NewConfigWizardModalWithProfiles(
		termInfo.Width,
		termInfo.Height,
		profileOptions,
	)

	if i.context.DefaultProfile != nil {
		i.wizardModal.SetProfileID(i.context.DefaultProfile.ID)
		i.wizardModal.SetAudience(i.context.DefaultProfile.TargetAudience)
	}

	i.wizardModal.Show()
	i.state = StateConfiguring

	return i.wizardModal.Init()
}

func (i *Intent) updateWizardFlow(msg tea.Msg) tea.Cmd {
	if windowMsg, ok := msg.(tea.WindowSizeMsg); ok {
		return i.handleWindowResize(windowMsg)
	}

	if msg, ok := msg.(WizardCompleteMsg); ok {
		return i.handleWizardComplete(msg)
	}

	if msg, ok := msg.(TechnologiesExtractedMsg); ok {
		return i.handleTechExtracted(msg)
	}

	if msg, ok := msg.(CVGenerationCompleteMsg); ok {
		return i.handleCVGenerated(msg)
	}

	if msg, ok := msg.(ExportCompleteMsg); ok {
		i.handleExportCompleteMsg(msg)
		return nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if cmd := i.handleGlobalKeys(keyMsg); cmd != nil {
			return cmd
		}
	}

	if cmd, handled := i.delegateToWizardModal(msg); handled {
		return cmd
	}

	if cmd, handled := i.delegateToExportModal(msg); handled {
		return cmd
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		return i.handleKeyDelegation(keyMsg)
	}

	return nil
}

func (i *Intent) handleWindowResize(msg tea.WindowSizeMsg) tea.Cmd {
	termInfo := i.GetTerminalInfo()
	termInfo.Update(msg)

	if i.reviewScreen != nil {
		i.reviewScreen.SetTerminalInfo(msg.Width, msg.Height)
	}
	if i.previewScreen != nil {
		i.previewScreen.SetTerminalInfo(msg.Width, msg.Height)
	}
	return nil
}

func (i *Intent) handleGlobalKeys(keyMsg tea.KeyMsg) tea.Cmd {
	switch keyMsg.String() {
	case "ctrl+c":
		i.setCancelled()
		return tea.Quit
	case "q":
		i.setCancelled()
		return tea.Quit
	case "?", "h":
		i.ToggleHelp()
		return nil
	}
	return nil
}

func (i *Intent) delegateToWizardModal(msg tea.Msg) (tea.Cmd, bool) {
	if i.wizardModal == nil || !i.wizardModal.IsVisible() {
		return nil, false
	}

	cmd := i.wizardModal.Update(msg)
	if i.wizardModal.IsCompleted() {
		config := i.wizardModal.GetConfigData()
		return i.handleWizardComplete(WizardCompleteMsg{
			ProfileID:    config.ProfileID,
			Audience:     config.Audience,
			TechFocus:    config.TechFocus,
			Technologies: config.Technologies,
			FocusArea:    config.FocusArea,
			SkillsFormat: config.SkillsFormat,
			SkillsLimit:  config.SkillsLimit,
			CVLength:     config.CVLength,
		}), true
	}
	if !i.wizardModal.IsVisible() && !i.wizardModal.IsCompleted() {
		i.setCancelled()
		return nil, true
	}
	return cmd, true
}

func (i *Intent) delegateToExportModal(msg tea.Msg) (tea.Cmd, bool) {
	if i.exportModal == nil || !i.exportModal.IsVisible() {
		return nil, false
	}

	cmd := i.exportModal.Update(msg)
	if i.exportModal.IsCompleted() {
		exportData := i.exportModal.GetExportData()
		return i.handleExportComplete(exportData), true
	}
	if !i.exportModal.IsVisible() && !i.exportModal.IsCompleted() {
		i.exportModal.Hide()
		i.state = StatePreview
		return nil, true
	}
	return cmd, true
}

func (i *Intent) handleKeyDelegation(keyMsg tea.KeyMsg) tea.Cmd {
	if i.progressModal != nil && i.progressModal.IsVisible() {
		if keyMsg.String() == "esc" && i.progressModal.IsCancellable() {
			i.progressModal.Hide()
			i.wizardModal.Show()
			i.state = StateConfiguring
			return nil
		}
	}

	if i.activeScreen != nil && (i.state == StateReview || i.state == StatePreview) {
		cmd, result := i.activeScreen.Update(keyMsg)
		if result != nil {
			return i.handleScreenResult(result)
		}
		return cmd
	}

	if i.state == StateExportComplete {
		return i.handleExportCompleteKeypress(keyMsg.String())
	}

	return nil
}

func (i *Intent) generateCVAsync() tea.Cmd {
	return func() tea.Msg {
		if i.context.CVGenerationService == nil || i.context.AppContext == nil {
			cvView := &career.CVView{
				ID:               fmt.Sprintf("cv_%d", time.Now().Unix()),
				Name:             i.selectedProfile.Name,
				TargetRole:       i.selectedProfile.TargetRole,
				TargetAudience:   i.selectedAudience,
				GeneratedAt:      time.Now(),
				SourceEventCount: len(i.context.Events),
				SourceFactCount:  len(i.context.Facts),
			}
			return CVGenerationCompleteMsg{CV: cvView, Error: nil}
		}

		ctx := i.context.AppContext
		config := &career.CVConfig{
			Name:           i.selectedProfile.Name,
			TargetRole:     i.selectedProfile.TargetRole,
			TargetAudience: i.selectedAudience,

			TechnologyFocus:      string(i.selectedTechnologyFocus),
			SelectedTechnologies: i.selectedTechnologies,
			FocusArea:            string(i.selectedFocusArea),
			LengthFormat:         string(cv.MapUILengthToFormat(i.selectedCVLength)),
			SkillsFormat:         i.selectedSkillsFormat,
			SkillsLimit:          i.selectedSkillsLimit,
		}

		cvView, err := i.context.CVGenerationService.GenerateCVFromConfig(ctx, config)
		if err != nil {
			if i.logger != nil {
				i.logger.Error("Failed to generate CV: %v", err)
			}
			return CVGenerationCompleteMsg{CV: nil, Error: err}
		}

		if i.logger != nil {
			i.logger.Info("Successfully generated CV: %s", cvView.ID)
		}
		return CVGenerationCompleteMsg{CV: cvView, Error: nil}
	}
}

func (i *Intent) extractTechnologiesAsync() tea.Cmd {
	return func() tea.Msg {
		if i.context.SkillRepository == nil || i.context.EventRepository == nil {
			return TechnologiesExtractedMsg{
				Technologies: []*ExtractedTechnology{},
				Suggestion: &FocusAreaSuggestion{
					Area:       technology.FocusAreaBackend,
					Confidence: 0.0,
					Evidence:   map[string]int{},
				},
				Error: nil,
			}
		}

		ctx := i.context.AppContext
		if ctx == nil {
			return TechnologiesExtractedMsg{
				Technologies: nil,
				Suggestion:   nil,
				Error:        errors.New("application context is required for technology extraction"),
			}
		}

		extractor := technology.NewExtractor(i.context.SkillRepository, i.context.EventRepository)
		techs, err := extractor.ExtractFromUser(ctx)
		if err != nil {
			return TechnologiesExtractedMsg{
				Technologies: nil,
				Suggestion:   nil,
				Error:        err,
			}
		}

		filtered := extractor.FilterByThreshold(techs, 3)

		analyzer := &technology.Analyzer{}
		suggestion := analyzer.AnalyzeSkills(filtered)

		return TechnologiesExtractedMsg{
			Technologies: filtered,
			Suggestion:   suggestion,
			Error:        nil,
		}
	}
}

func (i *Intent) exportCVAsync() tea.Cmd {
	return func() tea.Msg {
		if i.context.ExportService == nil {
			return ExportCompleteMsg{Path: "", Error: errors.New("export service not available")}
		}

		ctx := i.context.AppContext
		if ctx == nil {
			return ExportCompleteMsg{
				Path:  "",
				Error: errors.New("application context is required for export"),
			}
		}

		sections := i.generatedCV.Sections
		bulletsMap := make(map[string][]*career.CVBullet)

		var content string
		var err error
		var exportFormat cv.ExportFormat

		switch i.selectedExportFormat {
		case ExportFormatText:
			content, err = i.context.ExportService.ExportToText(ctx, i.generatedCV, sections, bulletsMap)
			exportFormat = cv.ExportFormatText
		case ExportFormatMarkdown:
			content, err = i.context.ExportService.ExportToMarkdown(ctx, i.generatedCV, sections, bulletsMap)
			exportFormat = cv.ExportFormatMarkdown
		case ExportFormatYAML:
			content, err = i.context.ExportService.ExportToYAML(ctx, i.generatedCV, sections, i.context.ProfileConfig)
			exportFormat = cv.ExportFormatYAML
		default:
			return ExportCompleteMsg{Path: "", Error: errors.New("unknown export format")}
		}

		if err != nil {
			return ExportCompleteMsg{Path: "", Error: fmt.Errorf("failed to export: %w", err)}
		}

		switch i.selectedExportOption {
		case ExportOptionSaveToFile:
			path, saveErr := i.context.ExportService.SaveToFile(ctx, i.generatedCV.Name, exportFormat, content)
			if saveErr != nil {
				return ExportCompleteMsg{Path: "", Error: fmt.Errorf("failed to save file: %w", saveErr)}
			}
			return ExportCompleteMsg{Path: path, Error: nil}

		case ExportOptionClipboard:
			clipErr := i.context.ExportService.CopyToClipboard(ctx, content)
			if clipErr != nil {
				return ExportCompleteMsg{Path: "", Error: fmt.Errorf("failed to copy to clipboard: %w", clipErr)}
			}
			return ExportCompleteMsg{Path: "clipboard", Error: nil}

		default:
			return ExportCompleteMsg{Path: "", Error: errors.New("unknown export option")}
		}
	}
}

func (i *Intent) showExportModal() tea.Cmd {
	termInfo := i.GetTerminalInfo()
	i.exportModal = cvmodals.NewExportModal(termInfo.Width, termInfo.Height)
	i.exportModal.Show()
	i.state = StateExporting
	return i.exportModal.Init()
}

func (i *Intent) returnToWizard() tea.Cmd {
	i.state = StateConfiguring
	i.wizardModal.Reset()
	return i.wizardModal.Init()
}
