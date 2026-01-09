package intents

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/service/career/cv"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// GenerateCVIntent implements the Intent interface for generating CVs.
type GenerateCVIntent struct {
	// Embed BaseIntent for terminal awareness, logo, and state management
	*BaseIntent

	context *GenerateCVContext
	state   *GenerateCVModel
	active  bool
	result  *IntentResult[*GenerateCVResult]
	logger  *logger.Logger
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

	// Create BaseIntent for terminal awareness and state management
	base := NewBaseIntent()

	return &GenerateCVIntent{
		BaseIntent: base,
		context:    context,
		state: &GenerateCVModel{
			context:             context,
			currentState:        GenerateCVStateSelectProfile,
			selectedProfile:     selectedProfile,
			selectedIndex:       0,
			selectedCVStructure: CVStructureStandard, // Default to standard
			structureIndex:      0,
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

// Theme helper methods for consistent themed styling.

// getCardStyle returns a themed card style, with fallback to default styling.
func (i *GenerateCVIntent) getCardStyle() lipgloss.Style {
	if theme := i.Theme(); theme != nil {
		return theme.Styles().CardBase
	}
	// Fallback to default styling
	return lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)
}

// getPrimaryColor returns the primary text color from theme or fallback.
func (i *GenerateCVIntent) getPrimaryColor() lipgloss.Color {
	if theme := i.Theme(); theme != nil {
		return theme.ForegroundColor()
	}
	return styles.ColorTextPrimary
}

// getAccentColor returns the accent color from theme or fallback.
func (i *GenerateCVIntent) getAccentColor() lipgloss.Color {
	if theme := i.Theme(); theme != nil {
		return theme.PrimaryColor()
	}
	return styles.ColorAccentTeal
}

// getErrorColor returns the error color from theme or fallback.
func (i *GenerateCVIntent) getErrorColor() lipgloss.Color {
	if theme := i.Theme(); theme != nil {
		return theme.ErrorColor()
	}
	return styles.ColorError
}

// getBorderColor returns the border color from theme or fallback.
func (i *GenerateCVIntent) getBorderColor() lipgloss.Color {
	if theme := i.Theme(); theme != nil {
		return theme.BorderColor()
	}
	return styles.ColorBorder
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
	case GenerateCVStateSelectStructure:
		return i.updateSelectStructure(msg)
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
				i.state.selectedAudience = i.state.selectedProfile.TargetAudience
				i.state.audienceIndex = 0
			}
			return nil
		}

		// Handle global keys (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			// At root state, back means cancel
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
	// Available audiences (must match viewSelectAudience)
	audiences := []string{"hiring_manager", "recruiter", "peer"}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if i.state.audienceIndex > 0 {
				i.state.audienceIndex--
			}
			return nil
		case "down", "j":
			if i.state.audienceIndex < len(audiences)-1 {
				i.state.audienceIndex++
			}
			return nil
		case "enter":
			// Set selected audience based on current index
			i.state.selectedAudience = audiences[i.state.audienceIndex]
			i.state.currentState = GenerateCVStateSelectStructure
			return nil
		}

		// Handle global keys (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			i.state.currentState = GenerateCVStateSelectProfile
			return nil
		}
	case AudienceSelectedMsg:
		i.state.selectedAudience = msg.Audience
		i.state.currentState = GenerateCVStateSelectStructure
	}
	return nil
}

// updateSelectStructure handles messages while selecting CV structure.
func (i *GenerateCVIntent) updateSelectStructure(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if i.state.structureIndex > 0 {
				i.state.structureIndex--
			}
			return nil
		case "down", "j":
			if i.state.structureIndex < 1 { // Only 2 options: Standard (0) and Narrative (1)
				i.state.structureIndex++
			}
			return nil
		case "enter":
			// Set selected structure based on current index
			structures := []CVStructure{CVStructureStandard, CVStructureNarrative}
			i.state.selectedCVStructure = structures[i.state.structureIndex]
			i.state.currentState = GenerateCVStateGenerating
			i.state.isGenerating = true
			return i.generateCVAsync()
		case "esc":
			i.state.currentState = GenerateCVStateSelectAudience
			return nil
		case "m":
			// Return to main menu
			i.setCancelled()
			return nil
		case "q", "ctrl+c":
			i.setCancelled()
			return nil
		}
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
				TargetAudience:   i.state.selectedAudience,
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
			TargetAudience: i.state.selectedAudience,
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

		// Initialize viewport for preview (80 cols x 20 rows)
		i.state.previewViewport = viewport.New(80, 20)
		i.state.previewViewport.Style = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(i.getBorderColor()).
			Padding(1, 2)

		i.state.currentState = GenerateCVStatePreview
		return nil
	case tea.KeyMsg:
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			// Let generation complete in background, navigate back
			i.state.currentState = GenerateCVStateSelectAudience
			return nil
		}
	}
	return nil
}

// updatePreview handles messages while previewing the CV.
func (i *GenerateCVIntent) updatePreview(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			i.state.currentState = GenerateCVStateSelectAudience
			return nil
		}

		switch msg.String() {
		case "up", "k", "down", "j", "pgup", "pgdown":
			// Handle viewport scrolling
			i.state.previewViewport, cmd = i.state.previewViewport.Update(msg)
			return cmd
		case "enter", "e":
			i.state.currentState = GenerateCVStateReview
			return nil
		case "c":
			i.state.currentState = GenerateCVStateConfirm
			return nil
		}
	}
	return nil
}

// updateReview handles messages while reviewing/editing the CV.
func (i *GenerateCVIntent) updateReview(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			i.state.currentState = GenerateCVStatePreview
			return nil
		}

		switch msg.String() {
		case "enter":
			i.state.currentState = GenerateCVStateConfirm
			return nil
		}
	}
	return nil
}

// updateConfirm handles messages while confirming the CV.
func (i *GenerateCVIntent) updateConfirm(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			i.state.currentState = GenerateCVStateReview
			return nil
		}

		switch msg.String() {
		case "y", "enter":
			i.setCompleted()
			return nil
		case "e", "x":
			// Transition to export format selection
			i.state.currentState = GenerateCVStateExportSelectFormat
			i.state.selectedIndex = 0
			return nil
		case "n":
			i.state.currentState = GenerateCVStateReview
			return nil
		}
	}
	return nil
}

// getStateContent returns the content for the current state.
func (i *GenerateCVIntent) getStateContent() string {
	switch i.state.currentState {
	case GenerateCVStateSelectProfile:
		return i.viewSelectProfile()
	case GenerateCVStateSelectAudience:
		return i.viewSelectAudience()
	case GenerateCVStateSelectStructure:
		return i.viewSelectStructure()
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
	default:
		return ""
	}
}

// getContextHelp returns context-aware help text for the current state.
func (i *GenerateCVIntent) getContextHelp() string {
	theme := i.Theme()

	switch i.state.currentState {
	case GenerateCVStateSelectProfile:
		return CombineThemedFooters(
			ThemedNavigationFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateSelectAudience:
		return CombineThemedFooters(
			ThemedNavigationFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateSelectStructure:
		return CombineThemedFooters(
			ThemedNavigationFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateGenerating:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("...", "Please wait"),
			),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStatePreview:
		return CombineThemedFooters(
			ThemedDetailViewFooter(theme),
			ThemedCustomFooter(theme,
				components.EditBadge(),
				components.NewKeyBadge("c", "Continue"),
			),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateReview:
		return CombineThemedFooters(
			ThemedDetailViewFooter(theme),
			ThemedCustomFooter(theme,
				components.NewKeyBadge("Enter", "Continue"),
			),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateConfirm:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("y/Enter", "Confirm"),
				components.NewKeyBadge("e/x", "Export"),
				components.NewKeyBadge("n/Esc", "Back"),
			),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateExportSelectFormat:
		return CombineThemedFooters(
			ThemedNavigationFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateExportSelectLocation:
		return CombineThemedFooters(
			ThemedNavigationFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateExporting:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("...", "Please wait"),
			),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateExportComplete:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("Enter", "Continue"),
				components.BackBadge(),
			),
			ThemedGlobalBadges(theme),
		)
	default:
		return ThemedGlobalBadges(theme)
	}
}

// View renders the intent's current state using StandardView.
func (i *GenerateCVIntent) View() string {
	if !i.active {
		return "GenerateCV intent is not active"
	}

	// Create standard view with breadcrumbs
	breadcrumbs := i.getBreadcrumbs()
	view := CreateStandardViewWithBreadcrumbs(i.BaseIntent, breadcrumbs...)

	// Get content for current state
	content := i.getStateContent()
	view.WithContent(content)

	// Get context-aware help
	help := i.getContextHelp()
	view.WithHelp(help).WithFooterSeparator(true)

	return view.Render()
}

// getBreadcrumbs returns breadcrumbs for the current state.
func (i *GenerateCVIntent) getBreadcrumbs() []string {
	crumbs := []string{"Main Menu", "Generate CV"}

	switch i.state.currentState {
	case GenerateCVStateSelectProfile:
		crumbs = append(crumbs, "Select Profile")
	case GenerateCVStateSelectAudience:
		crumbs = append(crumbs, "Select Audience")
	case GenerateCVStateSelectStructure:
		crumbs = append(crumbs, "Select Structure")
	case GenerateCVStateGenerating:
		crumbs = append(crumbs, "Generating")
	case GenerateCVStatePreview:
		crumbs = append(crumbs, "Preview")
	case GenerateCVStateReview:
		crumbs = append(crumbs, "Review")
	case GenerateCVStateConfirm:
		crumbs = append(crumbs, "Confirm")
	case GenerateCVStateExportSelectFormat:
		crumbs = append(crumbs, "Export", "Select Format")
	case GenerateCVStateExportSelectLocation:
		crumbs = append(crumbs, "Export", "Select Location")
	case GenerateCVStateExporting:
		crumbs = append(crumbs, "Export", "Exporting")
	case GenerateCVStateExportComplete:
		crumbs = append(crumbs, "Export", "Complete")
	}

	return crumbs
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

	return i.getCardStyle().Render(content.String())
}

// viewSelectAudience renders the audience selection view.
func (i *GenerateCVIntent) viewSelectAudience() string {
	var content strings.Builder

	// Display error message prominently if generation failed
	if i.state.generationError != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(i.getErrorColor()).
			Bold(true)
		content.WriteString(errorStyle.Render("❌ CV Generation Failed") + "\n\n")
		content.WriteString(fmt.Sprintf("Error: %s\n\n", i.state.generationError.Error()))
		content.WriteString("Please try selecting a different audience or check your data.\n\n")
		content.WriteString(strings.Repeat("─", 50) + "\n\n")
		// Clear error after displaying so it doesn't persist
		i.state.generationError = nil
	}

	content.WriteString("👥 Select Target Audience\n\n")

	if i.state.selectedProfile != nil {
		content.WriteString(fmt.Sprintf("Profile: %s\n", i.state.selectedProfile.Name))
		content.WriteString(fmt.Sprintf("Role: %s\n\n", i.state.selectedProfile.TargetRole))

		// Define available audiences with descriptions
		audiences := []struct {
			value       string
			label       string
			description string
		}{
			{"hiring_manager", "Hiring Manager", "Focus on outcomes, ownership, and business impact"},
			{"recruiter", "Recruiter", "Emphasize skills, competencies, and achievements"},
			{"peer", "Technical Peer", "Highlight technical depth, collaboration, and problem-solving"},
		}

		content.WriteString("Select target audience:\n\n")
		for idx, aud := range audiences {
			prefix := "  "
			if idx == i.state.audienceIndex {
				prefix = "▶ "
			}

			audienceStyle := lipgloss.NewStyle().Foreground(i.getPrimaryColor())
			if idx == i.state.audienceIndex {
				audienceStyle = audienceStyle.Foreground(i.getAccentColor()).Bold(true)
			}

			line := fmt.Sprintf("%s%s - %s", prefix, aud.label, aud.description)
			content.WriteString(audienceStyle.Render(line) + "\n")
		}
	}

	return i.getCardStyle().Render(content.String())
}

// viewSelectStructure renders the CV structure selection view.
func (i *GenerateCVIntent) viewSelectStructure() string {
	var content strings.Builder
	content.WriteString("\n📐 Select CV Structure\n\n")

	if i.state.selectedProfile != nil {
		content.WriteString(fmt.Sprintf("Profile: %s\n", i.state.selectedProfile.Name))
		content.WriteString(fmt.Sprintf("Role: %s\n", i.state.selectedProfile.TargetRole))
		content.WriteString(fmt.Sprintf("Audience: %s\n\n", i.state.selectedAudience))
	}

	// Define available structures with descriptions
	structures := []struct {
		value       string
		label       string
		description string
	}{
		{"standard", "Standard", "Traditional CV with Experience, Projects, Skills, Summary sections"},
		{"narrative", "Narrative", "Language-agnostic professional format with Core Strengths, Technologies, What I Bring sections"},
	}

	content.WriteString("Select CV structure:\n\n")
	for idx, structure := range structures {
		prefix := "  "
		if idx == i.state.structureIndex {
			prefix = "▶ "
		}

		structureStyle := lipgloss.NewStyle().Foreground(styles.ColorTextPrimary)
		if idx == i.state.structureIndex {
			structureStyle = structureStyle.Foreground(styles.ColorAccentTeal).Bold(true)
		}

		line := fmt.Sprintf("%s%s\n   %s", prefix, structure.label, structure.description)
		content.WriteString(structureStyle.Render(line) + "\n\n")
	}

	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	return card
}

// viewGenerating renders the CV generation progress view.
func (i *GenerateCVIntent) viewGenerating() string {
	var content strings.Builder
	content.WriteString("\n⏳ Generating CV...\n\n")

	if i.state.selectedProfile != nil {
		content.WriteString(fmt.Sprintf("Profile: %s\n", i.state.selectedProfile.Name))
		content.WriteString(fmt.Sprintf("Role: %s\n", i.state.selectedProfile.TargetRole))
		content.WriteString(fmt.Sprintf("Audience: %s\n\n", i.state.selectedAudience))
	}

	content.WriteString("Processing career events and facts...\n")
	content.WriteString("Extracting achievements and metrics...\n")
	content.WriteString("Generating professional bullets...\n")

	return i.getCardStyle().Render(content.String())
}

// viewPreview renders the CV preview view with scrollable content.
// Routes to structure-specific preview based on selectedCVStructure.
func (i *GenerateCVIntent) viewPreview() string {
	if i.state.generatedCV == nil {
		return "No CV generated yet"
	}

	var content string
	switch i.state.selectedCVStructure {
	case CVStructureNarrative:
		content = i.viewPreviewNarrative()
	default:
		content = i.viewPreviewStandard()
	}

	// Set viewport content
	i.state.previewViewport.SetContent(content)

	// Render viewport
	return i.state.previewViewport.View()
}

// viewPreviewStandard renders the standard CV preview (traditional format).
func (i *GenerateCVIntent) viewPreviewStandard() string {
	var content strings.Builder

	// Header with metadata
	content.WriteString(fmt.Sprintf("📄 CV Preview: %s\n", i.state.generatedCV.Name))
	content.WriteString(fmt.Sprintf("Target Role: %s | Audience: %s\n",
		i.state.generatedCV.TargetRole,
		i.state.generatedCV.TargetAudience))
	content.WriteString(fmt.Sprintf("Events: %d | Facts: %d | Generated: %s\n\n",
		i.state.generatedCV.SourceEventCount,
		i.state.generatedCV.SourceFactCount,
		i.state.generatedCV.GeneratedAt.Format(time.RFC822)))

	content.WriteString(strings.Repeat("─", 80) + "\n\n")

	// Render all sections with their content
	for idx, section := range i.state.generatedCV.Sections {
		if idx > 0 {
			content.WriteString("\n")
		}

		// Section title with themed color
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(i.getAccentColor()).
			MarginTop(1)
		content.WriteString(titleStyle.Render(strings.ToUpper(section.Title)) + "\n")
		content.WriteString(strings.Repeat("─", len(section.Title)) + "\n")

		// Handle summary section (prose)
		if section.SectionType == "summary" && section.Summary != "" {
			content.WriteString(section.Summary + "\n")
			continue
		}

		// Handle content groups (experience, projects, skills)
		for _, group := range section.Content {
			// Group header with dates
			if group.Header != "" {
				if group.StartDate != "" && group.EndDate != "" {
					if group.StartDate == group.EndDate {
						content.WriteString(fmt.Sprintf("  %s - %s\n", group.Header, group.StartDate))
					} else {
						content.WriteString(fmt.Sprintf("  %s - %s - %s\n", group.Header, group.StartDate, group.EndDate))
					}
				} else {
					content.WriteString(fmt.Sprintf("  %s\n", group.Header))
				}
			}

			// Bullets
			for _, bullet := range group.Bullets {
				content.WriteString(fmt.Sprintf("    • %s\n", bullet.Text))
			}
			content.WriteString("\n")
		}
	}

	return content.String()
}

// viewPreviewNarrative renders the narrative CV preview (professional format).
func (i *GenerateCVIntent) viewPreviewNarrative() string {
	var content strings.Builder
	profile := DefaultNarrativeProfile()

	// Profile header
	content.WriteString(fmt.Sprintf("# %s\n\n", profile.Name))
	content.WriteString(fmt.Sprintf("**%s**\n", profile.Role))
	content.WriteString(fmt.Sprintf("%s\n", profile.Location))
	content.WriteString(fmt.Sprintf("Email: %s\n", profile.Email))
	content.WriteString(fmt.Sprintf("GitHub: %s\n", profile.GitHub))
	content.WriteString(fmt.Sprintf("Portfolio: %s\n\n", profile.Portfolio))

	content.WriteString(strings.Repeat("─", 80) + "\n\n")

	// Summary section
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorAccentTeal)

	summary := getSummaryFromSections(i.state.generatedCV.Sections)
	content.WriteString(titleStyle.Render("SUMMARY") + "\n")
	content.WriteString(strings.Repeat("─", 7) + "\n")
	if summary != "" {
		content.WriteString(summary + "\n\n")
	} else {
		content.WriteString("Experienced software engineer with strong technical leadership skills.\n\n")
	}

	// Core Strengths section
	content.WriteString(titleStyle.Render("CORE STRENGTHS") + "\n")
	content.WriteString(strings.Repeat("─", 14) + "\n")
	strengths := extractStrengthsFromSections(i.state.generatedCV.Sections)
	for _, strength := range strengths {
		content.WriteString(fmt.Sprintf("  • %s\n", strength))
	}
	content.WriteString("\n")

	// Languages & Technologies section
	content.WriteString(titleStyle.Render("LANGUAGES & TECHNOLOGIES") + "\n")
	content.WriteString(strings.Repeat("─", 24) + "\n")
	languages, frontend, systems := extractTechnologiesFromSections(i.state.generatedCV.Sections)
	content.WriteString(fmt.Sprintf("**Languages:** %s\n", languages))
	content.WriteString(fmt.Sprintf("**Frontend:** %s\n", frontend))
	content.WriteString(fmt.Sprintf("**Systems:** %s\n\n", systems))

	// Selected Experience section (filtered by confidence)
	content.WriteString(titleStyle.Render("SELECTED EXPERIENCE") + "\n")
	content.WriteString(strings.Repeat("─", 19) + "\n")

	experienceSections := getExperienceSections(i.state.generatedCV.Sections)
	for _, section := range experienceSections {
		for _, group := range section.Content {
			// Filter bullets by confidence
			highConfidenceBullets := filterBulletsByConfidence(group.Bullets, MinConfidenceForNarrative)
			if len(highConfidenceBullets) == 0 {
				continue
			}

			// Group header with dates
			if group.Header != "" {
				if group.StartDate != "" && group.EndDate != "" {
					content.WriteString(fmt.Sprintf("\n### %s\n", group.Header))
					content.WriteString(fmt.Sprintf("*%s - %s*\n\n", group.StartDate, group.EndDate))
				} else {
					content.WriteString(fmt.Sprintf("\n### %s\n\n", group.Header))
				}
			}

			// High-confidence bullets only
			for _, bullet := range highConfidenceBullets {
				content.WriteString(fmt.Sprintf("  • %s\n", bullet.Text))
			}
		}
	}
	content.WriteString("\n")

	// What I Bring section
	content.WriteString(titleStyle.Render("WHAT I BRING") + "\n")
	content.WriteString(strings.Repeat("─", 12) + "\n")
	valueProps := extractValuePropositions(i.state.generatedCV.Sections)
	for _, prop := range valueProps {
		content.WriteString(fmt.Sprintf("  • %s\n", prop))
	}
	content.WriteString("\n")

	content.WriteString(strings.Repeat("─", 80) + "\n")
	content.WriteString("**References available on request.**\n")

	return content.String()
}

// viewReview renders the CV review/edit view.
func (i *GenerateCVIntent) viewReview() string {
	var content strings.Builder
	content.WriteString("\n✏️  Review & Edit CV\n\n")

	if i.state.generatedCV != nil {
		content.WriteString(fmt.Sprintf("Name: %s\n", i.state.generatedCV.Name))
		content.WriteString(fmt.Sprintf("Target Role: %s\n", i.state.generatedCV.TargetRole))
		content.WriteString(fmt.Sprintf("Target Audience: %s\n\n", i.state.generatedCV.TargetAudience))
		content.WriteString("CV content can be edited here.\n")
		content.WriteString("(Full editing interface would be implemented here)\n")
	}

	return i.getCardStyle().Render(content.String())
}

// viewConfirm renders the CV confirmation view.
func (i *GenerateCVIntent) viewConfirm() string {
	var content strings.Builder
	content.WriteString("\n✅ Confirm CV Generation\n\n")

	if i.state.generatedCV != nil {
		content.WriteString(fmt.Sprintf("Name: %s\n", i.state.generatedCV.Name))
		content.WriteString(fmt.Sprintf("Target Role: %s\n", i.state.generatedCV.TargetRole))
		content.WriteString(fmt.Sprintf("Target Audience: %s\n\n", i.state.generatedCV.TargetAudience))
		content.WriteString("Are you sure you want to complete CV generation?\n")
	}

	return i.getCardStyle().Render(content.String())
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
			GeneratedCV:       i.state.generatedCV,
			SelectedProfile:   i.state.selectedProfile,
			SelectedStructure: i.state.selectedCVStructure,
			AcceptedFields:    make(map[string]bool),
		},
		Metadata: map[string]interface{}{
			"profile":     i.state.selectedProfile.ID,
			"audience":    i.state.selectedAudience,
			"structure":   i.state.selectedCVStructure,
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
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			i.state.currentState = GenerateCVStateConfirm
			return nil
		}

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
		}
	}
	return nil
}

// updateExportSelectLocation handles location selection
func (i *GenerateCVIntent) updateExportSelectLocation(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			i.state.currentState = GenerateCVStateExportSelectFormat
			return nil
		}

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
			i.state.currentState = GenerateCVStateExportComplete // Show error in complete view, not selection view
			return nil
		}
		i.state.exportedPath = msg.Path
		i.state.currentState = GenerateCVStateExportComplete
		return nil
	case tea.KeyMsg:
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			i.state.isExporting = false
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			// Let export complete in background, navigate back
			i.state.currentState = GenerateCVStateExportSelectLocation
			return nil
		}
	}
	return nil
}

// updateExportComplete handles export completion
func (i *GenerateCVIntent) updateExportComplete(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			i.state.currentState = GenerateCVStateExportSelectLocation
			i.state.exportError = nil
			return nil
		}

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
					"audience":        i.state.selectedAudience,
					"export_format":   string(i.state.selectedExportFormat),
					"export_location": i.state.exportedPath,
					"timestamp":       time.Now(),
					"event_count":     len(i.context.Events),
					"fact_count":      len(i.context.Facts),
				},
			}
			i.active = false
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

		// Extract sections from generated CV
		sections := i.state.generatedCV.Sections
		// Build empty bullets map (kept for backward compatibility with export interface)
		bulletsMap := make(map[string][]*career.CVBullet)

		// Get export content based on format
		var content string
		var err error
		var exportFormat cv.ExportFormat

		switch i.state.selectedExportFormat {
		case CVExportFormatText:
			content, err = i.context.ExportService.ExportToText(ctx, i.state.generatedCV, sections, bulletsMap)
			exportFormat = cv.ExportFormatText
		case CVExportFormatMarkdown:
			content, err = i.context.ExportService.ExportToMarkdown(ctx, i.state.generatedCV, sections, bulletsMap)
			exportFormat = cv.ExportFormatMarkdown
		case CVExportFormatYAML:
			content, err = i.context.ExportService.ExportToYAML(ctx, i.state.generatedCV, sections, bulletsMap)
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

	return i.getCardStyle().Render(content.String())
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

	return i.getCardStyle().Render(content.String())
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

	return i.getCardStyle().Render(content.String())
}

// viewExportComplete renders the export completion view
func (i *GenerateCVIntent) viewExportComplete() string {
	var content strings.Builder

	if i.state.exportError != nil {
		content.WriteString("\n❌ Export Failed\n\n")
		content.WriteString(fmt.Sprintf("Error: %v\n\n", i.state.exportError))
		content.WriteString("Try a different location or format.\n")
	} else {
		content.WriteString("\n✅ Export Complete!\n\n")

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

	return i.getCardStyle().Render(content.String())
}
