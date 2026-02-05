package generatecv

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	cvmodals "github.com/baphled/kariya/internal/cli/screens/cv/modals"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/service/career/technology"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func (i *Intent) wizardView() string {
	termInfo := i.GetTerminalInfo()
	width, height := termInfo.Width, termInfo.Height

	breadcrumbs := i.getWizardBreadcrumbs()
	view := intents.CreateStandardViewWithBreadcrumbs(i.BaseIntent, breadcrumbs...)

	var content string
	switch i.state {
	case StateReview:
		if i.activeScreen != nil {
			return i.renderReviewScreenWithModalOverlay(width, height)
		}
		content = "Loading review..."
	case StatePreview:
		if i.activeScreen != nil {
			return i.renderPreviewScreenWithModalOverlay(width, height)
		}
		content = "Loading preview..."
	case StateExporting:
		content = "Exporting CV..."
	case StateExportComplete:
		content = i.viewExportComplete()
	default:
		content = ""
	}

	view.WithContent(content)
	view.WithHelp(i.getWizardContextHelp())
	view.WithFooterSeparator(true)

	baseView := view.Render()

	if i.wizardModal != nil && i.wizardModal.IsVisible() {
		return i.renderWizardModalOverlay(baseView, width, height)
	}
	if i.progressModal != nil && i.progressModal.IsVisible() {
		return i.renderProgressModalOverlay(baseView, width, height)
	}
	if i.exportModal != nil && i.exportModal.IsVisible() {
		return i.renderExportModalOverlay(baseView, width, height)
	}

	return baseView
}

func (i *Intent) getTheme() themes.Theme {
	if themeVal := i.Theme(); themeVal != nil {
		return themeVal
	}
	return themes.NewDefaultTheme()
}

func (i *Intent) getCardStyle() lipgloss.Style {
	return i.getTheme().Styles().CardBase
}

func (i *Intent) getWizardBreadcrumbs() []string {
	crumbs := []string{"Main Menu", "Generate CV"}

	switch i.state {
	case StateConfiguring:
		crumbs = append(crumbs, "Configure")
	case StateExtracting:
		crumbs = append(crumbs, "Extracting Technologies")
	case StateGenerating:
		crumbs = append(crumbs, "Generating CV")
	case StateReview:
		crumbs = append(crumbs, "Review")
	case StatePreview:
		crumbs = append(crumbs, "Preview")
	case StateExporting:
		crumbs = append(crumbs, "Exporting")
	}

	return crumbs
}

func (i *Intent) getWizardContextHelp() string {
	if i.wizardModal != nil && i.wizardModal.IsVisible() {
		return ""
	}
	if i.exportModal != nil && i.exportModal.IsVisible() {
		return ""
	}
	if i.progressModal != nil && i.progressModal.IsVisible() {
		return "Please wait   q Quit   m Main Menu"
	}

	switch i.state {
	case StateReview:
		return "Enter/p Preview   x Export   e Edit   Esc Back   q Quit"
	case StatePreview:
		return "Scroll   Enter/y Confirm   x Export   Esc Back   q Quit"
	default:
		return "q Quit   m Main Menu"
	}
}

func (i *Intent) renderWizardModalOverlay(baseView string, width, height int) string {
	if i.wizardModal == nil {
		return baseView
	}
	modalView := i.wizardModal.View()
	return containers.NewOverlay(width, height).Content(modalView).Dimmed().Render()
}

func (i *Intent) renderProgressModalOverlay(baseView string, width, height int) string {
	if i.progressModal == nil {
		return baseView
	}
	modalView := i.progressModal.View()
	return containers.NewOverlay(width, height).Content(modalView).Dimmed().Render()
}

func (i *Intent) renderExportModalOverlay(baseView string, width, height int) string {
	if i.exportModal == nil {
		return baseView
	}
	modalView := i.exportModal.View()
	return containers.NewOverlay(width, height).Content(modalView).Dimmed().Render()
}

func (i *Intent) renderReviewScreenWithModalOverlay(width, height int) string {
	breadcrumbs := i.getWizardBreadcrumbs()
	view := intents.CreateStandardViewWithBreadcrumbs(i.BaseIntent, breadcrumbs...)

	content := i.activeScreen.View()
	view.WithContent(content)

	help := "Enter/p Preview   x Export   e Edit   Esc Back   q Quit"
	view.WithHelp(help).WithFooterSeparator(true)

	baseView := view.Render()

	if i.exportModal != nil && i.exportModal.IsVisible() {
		return i.renderExportModalOverlay(baseView, width, height)
	}

	return baseView
}

func (i *Intent) renderPreviewScreenWithModalOverlay(width, height int) string {
	breadcrumbs := i.getWizardBreadcrumbs()
	view := intents.CreateStandardViewWithBreadcrumbs(i.BaseIntent, breadcrumbs...)

	content := i.activeScreen.View()
	view.WithContent(content)

	help := "jk Scroll   g/G Top/Bottom   Enter/y Confirm   x Export   Esc Back   q Quit"
	view.WithHelp(help).WithFooterSeparator(true)

	baseView := view.Render()

	if i.exportModal != nil && i.exportModal.IsVisible() {
		return i.renderExportModalOverlay(baseView, width, height)
	}

	return baseView
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
			content, err = i.context.ExportService.ExportToYAML(ctx, i.generatedCV, sections, bulletsMap)
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

func (i *Intent) viewExportComplete() string {
	th := theme.Default()
	var content strings.Builder

	if i.exportError != nil {
		errorHeader := primitives.ErrorText("Export Failed", th).Bold().Render()
		content.WriteString("\n" + errorHeader + "\n\n")
		content.WriteString(fmt.Sprintf("Error: %v\n\n", i.exportError))
		content.WriteString("Try a different location or format.\n")
	} else {
		successHeader := primitives.SuccessText("Export Complete!", th).Bold().Render()
		content.WriteString("\n" + successHeader + "\n\n")

		formatName := exportFormatDisplayName(i.selectedExportFormat)
		content.WriteString(fmt.Sprintf("Format: %s\n", formatName))

		if i.selectedExportOption == ExportOptionSaveToFile {
			content.WriteString(fmt.Sprintf("Location: %s\n\n", i.exportedPath))
			content.WriteString("You can now share this file!\n")
		} else {
			content.WriteString("Location: Clipboard\n\n")
			content.WriteString("You can now paste the CV anywhere!\n")
		}
	}

	return i.getCardStyle().Render(content.String())
}

func exportFormatDisplayName(format ExportFormat) string {
	switch format {
	case ExportFormatMarkdown:
		return "Markdown"
	case ExportFormatYAML:
		return "YAML"
	default:
		return "Text"
	}
}
