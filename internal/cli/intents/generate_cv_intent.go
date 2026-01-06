package intents

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/service/career/cv"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// GenerateCVIntent implements the Intent interface for generating CVs.
type GenerateCVIntent struct {
	context   *GenerateCVContext
	state     *GenerateCVModel
	active    bool
	result    *IntentResult[*GenerateCVResult]
	logger    *logger.Logger
	cvPreview *models.CVPreviewModel
}

// NewGenerateCVIntent creates a new GenerateCV intent.
func NewGenerateCVIntent(context *GenerateCVContext) (*GenerateCVIntent, error) {
	if err := context.Validate(); err != nil {
		return nil, err
	}

	selectedProfile := context.DefaultProfile
	if selectedProfile == nil && len(context.AvailableProfiles) > 0 {
		selectedProfile = context.AvailableProfiles[0]
	}

	return &GenerateCVIntent{
		context: context,
		state: &GenerateCVModel{
			context:         context,
			currentState:    GenerateCVStateSelectProfile,
			selectedProfile: selectedProfile,
			selectedIndex:   0,
		},
		active: true,
		logger: nil,
	}, nil
}

// Init is called when the intent is activated.
func (i *GenerateCVIntent) Init() tea.Cmd {
	if i.context.DefaultProfile != nil {
		i.state.selectedProfile = i.context.DefaultProfile
	}
	return nil
}

// Update processes a message in the intent.
func (i *GenerateCVIntent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	switch i.state.currentState {
	case GenerateCVStateSelectProfile:
		return i.updateSelectProfile(msg)
	case GenerateCVStateSelectAudience:
		return i.updateSelectAudience(msg)
	case GenerateCVStateGenerating:
		return i.updateGenerating(msg)
	case GenerateCVStatePreview:
		return i.updatePreview(msg)
	case GenerateCVStateReview:
		return i.updateReview(msg)
	case GenerateCVStateConfirm:
		return i.updateConfirm(msg)
	case GenerateCVStateExportSelectFormat:
		return i.updateExportSelectFormat(msg)
	case GenerateCVStateExportSelectLocation:
		return i.updateExportSelectLocation(msg)
	case GenerateCVStateExporting:
		return i.updateExporting(msg)
	case GenerateCVStateExportComplete:
		return i.updateExportComplete(msg)
	}
	return nil
}

// updateSelectProfile handles messages while selecting a profile.
func (i *GenerateCVIntent) updateSelectProfile(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if i.state.selectedIndex > 0 {
				i.state.selectedIndex--
				i.state.selectedProfile = i.context.AvailableProfiles[i.state.selectedIndex]
			}
			return nil
		case "down", "j":
			if i.state.selectedIndex < len(i.context.AvailableProfiles)-1 {
				i.state.selectedIndex++
				i.state.selectedProfile = i.context.AvailableProfiles[i.state.selectedIndex]
			}
			return nil
		case "enter":
			if i.state.selectedIndex >= 0 && i.state.selectedIndex < len(i.context.AvailableProfiles) {
				i.state.selectedProfile = i.context.AvailableProfiles[i.state.selectedIndex]
				i.state.currentState = GenerateCVStateSelectAudience
				i.state.selectedAudiences = i.state.selectedProfile.TargetAudience
			}
			return nil
		case "esc":
			i.setCancelled()
			return nil
		case "q", "ctrl+c":
			i.setCancelled()
			return nil
		}
	case ProfileSelectedMsg:
		i.state.selectedProfile = msg.Profile
		i.state.selectedIndex = msg.Index
		i.state.currentState = GenerateCVStateSelectAudience
		return nil
	}
	return nil
}

// updateSelectAudience handles messages while selecting audience(s).
func (i *GenerateCVIntent) updateSelectAudience(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if len(i.state.selectedAudiences) == 0 {
				i.state.selectedAudiences = i.state.selectedProfile.TargetAudience
			}
			i.state.currentState = GenerateCVStateGenerating
			i.state.isGenerating = true
			return i.generateCVAsync()
		case "q", "ctrl+c":
			i.setCancelled()
			return nil
		case "esc":
			i.state.currentState = GenerateCVStateSelectProfile
			return nil
		}
	case AudienceSelectedMsg:
		i.state.selectedAudiences = msg.Audiences
		i.state.currentState = GenerateCVStateGenerating
		i.state.isGenerating = true
		return i.generateCVAsync()
	}
	return nil
}

// generateCVAsync generates the CV asynchronously.
func (i *GenerateCVIntent) generateCVAsync() tea.Cmd {
	return func() tea.Msg {
		// Check if services are available
		if i.context.CVGenerationService == nil || i.context.AppContext == nil {
			// Fallback: create a minimal CV for testing
			cvView := &career.CVView{
				ID:               fmt.Sprintf("cv_%d", time.Now().Unix()),
				Name:             i.state.selectedProfile.Name,
				TargetRole:       i.state.selectedProfile.TargetRole,
				TargetAudience:   i.state.selectedAudiences,
				GeneratedAt:      time.Now(),
				SourceEventCount: len(i.context.Events),
				SourceFactCount:  len(i.context.Facts),
			}
			return CVGenerationCompleteMsg{CV: cvView, Error: nil}
		}

		ctx := i.context.AppContext
		config := &career.CVConfig{
			Name:           i.state.selectedProfile.Name,
			TargetRole:     i.state.selectedProfile.TargetRole,
			TargetAudience: i.state.selectedAudiences,
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

// updateGenerating handles messages while CV is being generated.
func (i *GenerateCVIntent) updateGenerating(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case CVGenerationCompleteMsg:
		i.state.isGenerating = false
		if msg.Error != nil {
			i.state.generationError = msg.Error
			i.state.currentState = GenerateCVStateSelectAudience
			return nil
		}
		i.state.generatedCV = msg.CV
		i.state.currentState = GenerateCVStatePreview
		// i.cvPreview = models.NewCVPreviewModel(msg.CV) // TODO: implement CV preview
		return nil
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			i.setCancelled()
			return nil
		}
	}
	return nil
}

// updatePreview handles messages while previewing the CV.
func (i *GenerateCVIntent) updatePreview(msg tea.Msg) tea.Cmd {
	// if i.cvPreview != nil {
	// 	_, cmd := i.cvPreview.Update(msg)
	// 	if cmd != nil {
	// 		return cmd
	// 	}
	// }

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "e":
			i.state.currentState = GenerateCVStateReview
			return nil
		case "c":
			i.state.currentState = GenerateCVStateConfirm
			return nil
		case "q", "ctrl+c":
			i.setCancelled()
			return nil
		case "esc":
			i.state.currentState = GenerateCVStateSelectAudience
			return nil
		}
	}
	return nil
}

// updateReview handles messages while reviewing/editing the CV.
func (i *GenerateCVIntent) updateReview(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i.state.currentState = GenerateCVStateConfirm
			return nil
		case "q", "ctrl+c":
			i.setCancelled()
			return nil
		case "esc":
			i.state.currentState = GenerateCVStatePreview
			return nil
		}
	}
	return nil
}

// updateConfirm handles messages while confirming the CV.
func (i *GenerateCVIntent) updateConfirm(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "enter":
			i.setCompleted()
			return nil
		case "e", "x":
			// Transition to export format selection
			i.state.currentState = GenerateCVStateExportSelectFormat
			i.state.selectedIndex = 0
			return nil
		case "n", "esc":
			i.state.currentState = GenerateCVStateReview
			return nil
		case "q", "ctrl+c":
			i.setCancelled()
			return nil
		}
	}
	return nil
}

// View renders the intent's current state.
func (i *GenerateCVIntent) View() string {
	switch i.state.currentState {
	case GenerateCVStateSelectProfile:
		return i.viewSelectProfile()
	case GenerateCVStateSelectAudience:
		return i.viewSelectAudience()
	case GenerateCVStateGenerating:
		return i.viewGenerating()
	case GenerateCVStatePreview:
		return i.viewPreview()
	case GenerateCVStateReview:
		return i.viewReview()
	case GenerateCVStateConfirm:
		return i.viewConfirm()
	case GenerateCVStateExportSelectFormat:
		return i.viewExportSelectFormat()
	case GenerateCVStateExportSelectLocation:
		return i.viewExportSelectLocation()
	case GenerateCVStateExporting:
		return i.viewExporting()
	case GenerateCVStateExportComplete:
		return i.viewExportComplete()
	}
	return ""
}

// viewSelectProfile renders the profile selection view.
func (i *GenerateCVIntent) viewSelectProfile() string {
	var content strings.Builder
	content.WriteString("\n📋 Select CV Profile\n\n")

	if len(i.context.AvailableProfiles) == 0 {
		content.WriteString("No profiles available.\n")
	} else {
		for idx, profile := range i.context.AvailableProfiles {
			prefix := "  "
			if idx == i.state.selectedIndex {
				prefix = "▶ "
			}
			content.WriteString(fmt.Sprintf("%s%s (%s)\n", prefix, profile.Name, profile.TargetRole))
		}
	}

	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("↑/k up, ↓/j down, Enter to select, q to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewSelectAudience renders the audience selection view.
func (i *GenerateCVIntent) viewSelectAudience() string {
	var content strings.Builder
	content.WriteString("\n👥 Select Target Audience(s)\n\n")

	if i.state.selectedProfile != nil {
		content.WriteString(fmt.Sprintf("Profile: %s\n", i.state.selectedProfile.Name))
		content.WriteString(fmt.Sprintf("Role: %s\n\n", i.state.selectedProfile.TargetRole))

		if len(i.state.selectedProfile.TargetAudience) > 0 {
			content.WriteString("Default audiences:\n")
			for _, aud := range i.state.selectedProfile.TargetAudience {
				content.WriteString(fmt.Sprintf("  • %s\n", aud))
			}
		}
	}

	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Enter to generate CV, Esc to go back, q to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewGenerating renders the CV generation progress view.
func (i *GenerateCVIntent) viewGenerating() string {
	var content strings.Builder
	content.WriteString("\n⏳ Generating CV...\n\n")

	if i.state.selectedProfile != nil {
		content.WriteString(fmt.Sprintf("Profile: %s\n", i.state.selectedProfile.Name))
		content.WriteString(fmt.Sprintf("Role: %s\n", i.state.selectedProfile.TargetRole))
		content.WriteString(fmt.Sprintf("Audiences: %s\n\n", strings.Join(i.state.selectedAudiences, ", ")))
	}

	content.WriteString("Processing career events and facts...\n")
	content.WriteString("Extracting achievements and metrics...\n")
	content.WriteString("Generating professional bullets...\n")

	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Press q to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewPreview renders the CV preview view.
func (i *GenerateCVIntent) viewPreview() string {
	// if i.cvPreview != nil {
	// 	return i.cvPreview.View()
	// }

	var content strings.Builder
	content.WriteString("\n📄 CV Preview\n\n")

	if i.state.generatedCV != nil {
		content.WriteString(fmt.Sprintf("Name: %s\n", i.state.generatedCV.Name))
		content.WriteString(fmt.Sprintf("Target Role: %s\n", i.state.generatedCV.TargetRole))
		content.WriteString(fmt.Sprintf("Target Audiences: %s\n", strings.Join(i.state.generatedCV.TargetAudience, ", ")))
		content.WriteString(fmt.Sprintf("Source Events: %d\n", i.state.generatedCV.SourceEventCount))
		content.WriteString(fmt.Sprintf("Source Facts: %d\n", i.state.generatedCV.SourceFactCount))
		content.WriteString(fmt.Sprintf("Generated: %s\n", i.state.generatedCV.GeneratedAt.Format(time.RFC822)))
	}

	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("e to edit, c to confirm, Esc to go back, q to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewReview renders the CV review/edit view.
func (i *GenerateCVIntent) viewReview() string {
	var content strings.Builder
	content.WriteString("\n✏️  Review & Edit CV\n\n")

	if i.state.generatedCV != nil {
		content.WriteString(fmt.Sprintf("Name: %s\n", i.state.generatedCV.Name))
		content.WriteString(fmt.Sprintf("Target Role: %s\n", i.state.generatedCV.TargetRole))
		content.WriteString(fmt.Sprintf("Target Audiences: %s\n\n", strings.Join(i.state.generatedCV.TargetAudience, ", ")))
		content.WriteString("CV content can be edited here.\n")
		content.WriteString("(Full editing interface would be implemented here)\n")
	}

	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Enter to continue, Esc to go back, q to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewConfirm renders the CV confirmation view.
func (i *GenerateCVIntent) viewConfirm() string {
	var content strings.Builder
	content.WriteString("\n✅ Confirm CV Generation\n\n")

	if i.state.generatedCV != nil {
		content.WriteString(fmt.Sprintf("Name: %s\n", i.state.generatedCV.Name))
		content.WriteString(fmt.Sprintf("Target Role: %s\n", i.state.generatedCV.TargetRole))
		content.WriteString(fmt.Sprintf("Target Audiences: %s\n\n", strings.Join(i.state.generatedCV.TargetAudience, ", ")))
		content.WriteString("Are you sure you want to complete CV generation?\n")
	}

	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("y/Enter to confirm, e/x to export, n/Esc to go back, q to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// Result returns the final result of the intent.
func (i *GenerateCVIntent) Result() *IntentResult[interface{}] {
	if i.result == nil {
		return nil
	}

	return &IntentResult[interface{}]{
		Status:   i.result.Status,
		Data:     i.result.Data,
		Error:    i.result.Error,
		Metadata: i.result.Metadata,
	}
}

// setCompleted marks the intent as completed with success.
func (i *GenerateCVIntent) setCompleted() {
	i.result = &IntentResult[*GenerateCVResult]{
		Status: Completed,
		Data: &GenerateCVResult{
			GeneratedCV:     i.state.generatedCV,
			SelectedProfile: i.state.selectedProfile,
			AcceptedFields:  make(map[string]bool),
		},
		Metadata: map[string]interface{}{
			"profile":     i.state.selectedProfile.ID,
			"audiences":   strings.Join(i.state.selectedAudiences, ","),
			"timestamp":   time.Now(),
			"event_count": len(i.context.Events),
			"fact_count":  len(i.context.Facts),
		},
	}
	i.active = false
}

// setCancelled marks the intent as cancelled by the user.
func (i *GenerateCVIntent) setCancelled() {
	i.result = &IntentResult[*GenerateCVResult]{
		Status: Cancelled,
	}
	i.active = false
}

// Export state handlers

// updateExportSelectFormat handles format selection
func (i *GenerateCVIntent) updateExportSelectFormat(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if i.state.selectedIndex > 0 {
				i.state.selectedIndex--
			}
			return nil
		case "down", "j":
			if i.state.selectedIndex < 2 {
				i.state.selectedIndex++
			}
			return nil
		case "enter":
			formats := []CVExportFormat{CVExportFormatText, CVExportFormatMarkdown, CVExportFormatYAML}
			if i.state.selectedIndex < len(formats) {
				i.state.selectedExportFormat = formats[i.state.selectedIndex]
				i.state.currentState = GenerateCVStateExportSelectLocation
				i.state.selectedIndex = 0
			}
			return nil
		case "esc":
			i.state.currentState = GenerateCVStateConfirm
			return nil
		case "q", "ctrl+c":
			i.setCancelled()
			return nil
		}
	}
	return nil
}

// updateExportSelectLocation handles location selection
func (i *GenerateCVIntent) updateExportSelectLocation(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if i.state.selectedIndex > 0 {
				i.state.selectedIndex--
			}
			return nil
		case "down", "j":
			if i.state.selectedIndex < 2 {
				i.state.selectedIndex++
			}
			return nil
		case "enter":
			options := []CVExportOption{CVExportOptionSaveToFile, CVExportOptionClipboard, CVExportOptionCancel}
			if i.state.selectedIndex < len(options) {
				selectedOption := options[i.state.selectedIndex]
				if selectedOption == CVExportOptionCancel {
					i.state.currentState = GenerateCVStateConfirm
					return nil
				}
				i.state.selectedExportOption = selectedOption
				i.state.currentState = GenerateCVStateExporting
				i.state.isExporting = true
				return i.exportCVAsync()
			}
			return nil
		case "esc":
			i.state.currentState = GenerateCVStateExportSelectFormat
			return nil
		case "q", "ctrl+c":
			i.setCancelled()
			return nil
		}
	}
	return nil
}

// updateExporting handles export progress
func (i *GenerateCVIntent) updateExporting(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case CVExportCompleteMsg:
		i.state.isExporting = false
		if msg.Error != nil {
			i.state.exportError = msg.Error
			i.state.currentState = GenerateCVStateExportSelectLocation
			return nil
		}
		i.state.exportedPath = msg.Path
		i.state.currentState = GenerateCVStateExportComplete
		return nil
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			i.state.isExporting = false
			i.setCancelled()
			return nil
		}
	}
	return nil
}

// updateExportComplete handles export completion
func (i *GenerateCVIntent) updateExportComplete(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Complete with export info
			now := time.Now()
			i.result = &IntentResult[*GenerateCVResult]{
				Status: Completed,
				Data: &GenerateCVResult{
					GeneratedCV:     i.state.generatedCV,
					SelectedProfile: i.state.selectedProfile,
					AcceptedFields:  make(map[string]bool),
					ExportPath:      i.state.exportedPath,
					CVExportFormat:  string(i.state.selectedExportFormat),
					ExportedAt:      &now,
				},
				Metadata: map[string]interface{}{
					"profile":         i.state.selectedProfile.ID,
					"audiences":       strings.Join(i.state.selectedAudiences, ","),
					"export_format":   string(i.state.selectedExportFormat),
					"export_location": i.state.exportedPath,
					"timestamp":       time.Now(),
					"event_count":     len(i.context.Events),
					"fact_count":      len(i.context.Facts),
				},
			}
			i.active = false
			return nil
		case "esc":
			i.state.currentState = GenerateCVStateExportSelectLocation
			i.state.exportError = nil
			return nil
		case "q", "ctrl+c":
			i.setCancelled()
			return nil
		}
	}
	return nil
}

// exportCVAsync exports the CV asynchronously
func (i *GenerateCVIntent) exportCVAsync() tea.Cmd {
	return func() tea.Msg {
		// Check if export service is available
		if i.context.ExportService == nil {
			return CVExportCompleteMsg{Path: "", Error: fmt.Errorf("export service not available")}
		}

		ctx := i.context.AppContext
		if ctx == nil {
			ctx = context.Background()
		}

		// Get export content based on format
		var content string
		var err error
		var exportFormat cv.ExportFormat

		switch i.state.selectedExportFormat {
		case CVExportFormatText:
			content, err = i.context.ExportService.ExportToText(ctx, i.state.generatedCV, nil, nil)
			exportFormat = cv.ExportFormatText
		case CVExportFormatMarkdown:
			content, err = i.context.ExportService.ExportToMarkdown(ctx, i.state.generatedCV, nil, nil)
			exportFormat = cv.ExportFormatMarkdown
		case CVExportFormatYAML:
			content, err = i.context.ExportService.ExportToYAML(ctx, i.state.generatedCV, nil, nil)
			exportFormat = cv.ExportFormatYAML
		default:
			return CVExportCompleteMsg{Path: "", Error: fmt.Errorf("unknown export format")}
		}

		if err != nil {
			return CVExportCompleteMsg{Path: "", Error: fmt.Errorf("failed to export: %v", err)}
		}

		// Save based on option
		switch i.state.selectedExportOption {
		case CVExportOptionSaveToFile:
			path, err := i.context.ExportService.SaveToFile(ctx, i.state.generatedCV.Name, exportFormat, content)
			if err != nil {
				return CVExportCompleteMsg{Path: "", Error: fmt.Errorf("failed to save file: %v", err)}
			}
			return CVExportCompleteMsg{Path: path, Error: nil}

		case CVExportOptionClipboard:
			err := i.context.ExportService.CopyToClipboard(ctx, content)
			if err != nil {
				return CVExportCompleteMsg{Path: "", Error: fmt.Errorf("failed to copy to clipboard: %v", err)}
			}
			return CVExportCompleteMsg{Path: "clipboard", Error: nil}

		default:
			return CVExportCompleteMsg{Path: "", Error: fmt.Errorf("unknown export option")}
		}
	}
}

// Export view methods

// viewExportSelectFormat renders the export format selection view
func (i *GenerateCVIntent) viewExportSelectFormat() string {
	var content strings.Builder
	content.WriteString("\n📤 Export Format\n\n")

	formats := []struct {
		name        string
		description string
	}{
		{"Text", "Plain text for email and sharing"},
		{"Markdown", "Formatted for GitHub and docs"},
		{"YAML", "Structured for data interchange"},
	}

	for idx, format := range formats {
		prefix := "  "
		if idx == i.state.selectedIndex {
			prefix = "▶ "
		}
		content.WriteString(fmt.Sprintf("%s%s\n", prefix, format.name))
		content.WriteString(fmt.Sprintf("   %s\n\n", format.description))
	}

	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("↑/k up, ↓/j down, Enter to select, Esc to go back")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewExportSelectLocation renders the save location selection view
func (i *GenerateCVIntent) viewExportSelectLocation() string {
	var content strings.Builder
	content.WriteString("\n💾 Save Location\n\n")

	formatName := "Text"
	switch i.state.selectedExportFormat {
	case CVExportFormatMarkdown:
		formatName = "Markdown"
	case CVExportFormatYAML:
		formatName = "YAML"
	}

	content.WriteString(fmt.Sprintf("Format: %s\n", formatName))
	content.WriteString(fmt.Sprintf("CV Name: %s\n\n", i.state.generatedCV.Name))
	content.WriteString("Where to save?\n\n")

	options := []struct {
		name        string
		description string
	}{
		{"Save to file", "~/kariya-cvs/"},
		{"Copy to clipboard", "Paste anywhere"},
		{"Cancel", "Return to review"},
	}

	for idx, option := range options {
		prefix := "  "
		if idx == i.state.selectedIndex {
			prefix = "▶ "
		}
		content.WriteString(fmt.Sprintf("%s%s\n", prefix, option.name))
		content.WriteString(fmt.Sprintf("   %s\n\n", option.description))
	}

	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("↑/k up, ↓/j down, Enter to select, Esc to go back")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewExporting renders the export progress view
func (i *GenerateCVIntent) viewExporting() string {
	var content strings.Builder
	content.WriteString("\n⏳ Exporting CV...\n\n")

	formatName := "Text"
	switch i.state.selectedExportFormat {
	case CVExportFormatMarkdown:
		formatName = "Markdown"
	case CVExportFormatYAML:
		formatName = "YAML"
	}

	content.WriteString(fmt.Sprintf("Format: %s\n", formatName))

	if i.state.selectedExportOption == CVExportOptionSaveToFile {
		content.WriteString("Destination: ~/kariya-cvs/\n")
	} else {
		content.WriteString("Destination: Clipboard\n")
	}

	content.WriteString(fmt.Sprintf("CV Name: %s\n\n", i.state.generatedCV.Name))
	content.WriteString("Processing...\n")
	content.WriteString("• Formatting content\n")
	content.WriteString("• Generating filename\n")
	if i.state.selectedExportOption == CVExportOptionSaveToFile {
		content.WriteString("• Writing to disk\n")
	} else {
		content.WriteString("• Copying to clipboard\n")
	}

	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Press q to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewExportComplete renders the export completion view
func (i *GenerateCVIntent) viewExportComplete() string {
	var content strings.Builder
	content.WriteString("\n✅ Export Complete!\n\n")

	if i.state.exportError != nil {
		content.WriteString("❌ Error during export\n\n")
		content.WriteString(fmt.Sprintf("Error: %v\n\n", i.state.exportError))
		content.WriteString("Try a different location or format.\n")
	} else {
		formatName := "Text"
		switch i.state.selectedExportFormat {
		case CVExportFormatMarkdown:
			formatName = "Markdown"
		case CVExportFormatYAML:
			formatName = "YAML"
		}

		content.WriteString(fmt.Sprintf("Format: %s\n", formatName))

		if i.state.selectedExportOption == CVExportOptionSaveToFile {
			content.WriteString(fmt.Sprintf("Location: %s\n\n", i.state.exportedPath))
			content.WriteString("You can now share this file!\n")
		} else {
			content.WriteString("Location: Clipboard\n\n")
			content.WriteString("You can now paste the CV anywhere!\n")
		}
	}

	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Enter to finish, Esc to go back, q to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}
