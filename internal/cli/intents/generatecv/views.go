package generatecv

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/charmbracelet/lipgloss"
)

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
