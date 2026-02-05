package intents

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/service/career/technology"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// GenerateCVIntent orchestrates the multi-step CV generation workflow, managing
// state transitions between profile selection, technology extraction, generation,
// preview, and export phases.
//
// Side effects:
//   - None at construction; state mutations occur through Update and Init.
type GenerateCVIntent struct {
	// Embed BaseIntent for terminal awareness, logo, and state management
	*BaseIntent

	context *GenerateCVContext
	state   *GenerateCVModel
	active  bool
	result  *IntentResult[*GenerateCVResult]
	logger  *logger.Logger

	// Modals for wizard workflow
	wizardModal   *components.CVConfigWizardModal
	progressModal *components.CVProgressModal
	exportModal   *components.ExportOptionsModal

	// Screens for wizard workflow (using base Screen interface to avoid import cycle)
	// Will be *cvscreens.ReviewScreen and *cvscreens.CVPreviewScreen at runtime
	wizardReviewScreen  screens.Screen
	wizardPreviewScreen screens.Screen
}

// NewGenerateCVIntent constructs a GenerateCV intent initialised with the given
// context, selecting a default profile when one is available.
//
// Expected:
//   - context must pass Validate (non-nil, at least one profile and one event).
//
// Returns:
//   - A ready-to-activate intent and nil error on success.
//   - Nil intent and a validation error when context is invalid.
//
// Side effects:
//   - None.
func NewGenerateCVIntent(ctx *GenerateCVContext) (*GenerateCVIntent, error) {
	if err := ctx.Validate(); err != nil {
		return nil, err
	}

	selectedProfile := ctx.DefaultProfile
	if selectedProfile == nil && len(ctx.AvailableProfiles) > 0 {
		selectedProfile = ctx.AvailableProfiles[0]
	}

	// Create BaseIntent for terminal awareness and state management
	baseIntent := NewBaseIntent()

	return &GenerateCVIntent{
		BaseIntent: baseIntent,
		context:    ctx,
		state: &GenerateCVModel{
			context:         ctx,
			currentState:    CVStateConfiguring,
			selectedProfile: selectedProfile,
		},
		active: true,
		logger: nil,
	}, nil
}

// Init prepares the intent for its first render cycle, optionally bootstrapping
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (i *GenerateCVIntent) Init() tea.Cmd {
	if i.context.DefaultProfile != nil {
		i.state.selectedProfile = i.context.DefaultProfile
	}

	return i.initWizardFlow()
}

// initWizardFlow initializes the wizard-based workflow.
func (i *GenerateCVIntent) initWizardFlow() tea.Cmd {
	// Get terminal dimensions
	termInfo := i.GetTerminalInfo()

	// Convert CVProfiles to ProfileOptions
	profileOptions := make([]components.ProfileOption, len(i.context.AvailableProfiles))
	for idx, profile := range i.context.AvailableProfiles {
		profileOptions[idx] = components.ProfileOption{
			ID:   profile.ID,
			Name: profile.Name,
		}
	}

	// Create configuration wizard modal
	i.wizardModal = components.NewCVConfigWizardModalWithProfiles(
		termInfo.Width,
		termInfo.Height,
		profileOptions,
	)

	// Set default profile if available
	if i.context.DefaultProfile != nil {
		i.wizardModal.SetProfileID(i.context.DefaultProfile.ID)
		i.wizardModal.SetAudience(i.context.DefaultProfile.TargetAudience)
	}

	// Show the modal
	i.wizardModal.Show()

	// Set initial state
	i.state.currentState = CVStateConfiguring

	return i.wizardModal.Init()
}

// Theme helper methods for consistent themed styling.

// getCardStyle returns a themed card style, with fallback to default styling.
// getTheme returns the theme or a default.
func (i *GenerateCVIntent) getTheme() themes.Theme {
	if themeVal := i.Theme(); themeVal != nil {
		return themeVal
	}
	return themes.NewDefaultTheme()
}

func (i *GenerateCVIntent) getCardStyle() lipgloss.Style {
	return i.getTheme().Styles().CardBase
}

// Update advances the intent state machine by processing a single Bubble Tea
// message, delegating to the wizard flow or active screen handlers.
//
// Expected:
//   - msg must be a valid tea.Msg (key press, window resize, or async result).
//
// Returns:
//   - A tea.Cmd for follow-up work (async generation, quit, etc.), or nil.
//
// Side effects:
//   - Mutates internal state, active screens, and modal visibility.
func (i *GenerateCVIntent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	return i.updateWizardFlow(msg)
}

// updateWizardFlow handles wizard-based workflow updates with 3-tier priority.
// Tier 1 (Highest): Modal updates
// Tier 2: Global keys (quit, help, main menu)
// Tier 3: Screen updates
func (i *GenerateCVIntent) updateWizardFlow(msg tea.Msg) tea.Cmd {
	// Handle window size messages
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		// Update BaseIntent terminal info
		termInfo := i.GetTerminalInfo()
		termInfo.Update(msg)

		// Modals handle their own window size internally via their Update methods

		// Update screen terminal info
		if i.wizardReviewScreen != nil {
			i.wizardReviewScreen.SetTerminalInfo(msg.Width, msg.Height)
		}
		if i.wizardPreviewScreen != nil {
			i.wizardPreviewScreen.SetTerminalInfo(msg.Width, msg.Height)
		}
		return nil
	}

	// Handle wizard complete message
	if msg, ok := msg.(WizardCompleteMsg); ok {
		return i.handleWizardComplete(msg)
	}

	// Handle tech extraction complete
	if msg, ok := msg.(TechnologiesExtractedMsg); ok {
		return i.handleTechExtracted(msg)
	}

	// Handle CV generation complete
	if msg, ok := msg.(CVGenerationCompleteMsg); ok {
		return i.handleCVGenerated(msg)
	}

	if msg, ok := msg.(CVExportCompleteMsg); ok {
		i.state.isExporting = false
		// Hide any visible modals
		if i.progressModal != nil {
			i.progressModal.Hide()
		}
		if i.exportModal != nil {
			i.exportModal.Hide()
		}
		if msg.Error != nil {
			i.state.exportError = msg.Error
			i.state.currentState = CVStateExportComplete
			return nil
		}
		i.state.exportedPath = msg.Path
		i.state.currentState = CVStateExportComplete
		return nil
	}

	// Tier 1: Global keys (HIGHEST PRIORITY - must work everywhere)
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+c":
			// Quit immediately
			i.setCancelled()
			return tea.Quit
		case "q":
			// Quit with cancellation
			i.setCancelled()
			return tea.Quit
		case "?", "h":
			i.ToggleHelp()
			return nil
		}
	}

	// Tier 2: Modal updates (MUST be outside keyMsg check to receive ALL messages)
	// Wizard modal
	if i.wizardModal != nil && i.wizardModal.IsVisible() {
		cmd := i.wizardModal.Update(msg)
		// Check if wizard was completed
		if i.wizardModal.IsCompleted() {
			config := i.wizardModal.GetConfigData()
			return i.handleWizardComplete(WizardCompleteMsg{
				// Step 1: WHO
				ProfileID: config.ProfileID,
				Audience:  config.Audience,
				// Step 2: TECH
				TechFocus:    config.TechFocus,
				Technologies: config.Technologies,
				FocusArea:    config.FocusArea,
				// Step 3: FORMAT
				SkillsFormat: config.SkillsFormat,
				SkillsLimit:  config.SkillsLimit,
				CVLength:     config.CVLength,
			})
		}
		// Check if wizard was cancelled (hidden without completing)
		if !i.wizardModal.IsVisible() && !i.wizardModal.IsCompleted() {
			i.setCancelled()
			return nil
		}
		return cmd
	}

	// Export modal
	if i.exportModal != nil && i.exportModal.IsVisible() {
		cmd := i.exportModal.Update(msg)
		// Check if export was completed
		if i.exportModal.IsCompleted() {
			exportData := i.exportModal.GetExportData()
			return i.handleExportComplete(exportData)
		}
		// Check if export was cancelled (hidden without completing)
		if !i.exportModal.IsVisible() && !i.exportModal.IsCompleted() {
			i.exportModal.Hide()
			i.state.currentState = CVStatePreview
			return nil
		}
		return cmd
	}

	// Progress modal - handle Esc to cancel async operations
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if i.progressModal != nil && i.progressModal.IsVisible() {
			if keyMsg.String() == "esc" && i.progressModal.IsCancellable() {
				// Cancel current async operation and return to wizard
				i.progressModal.Hide()
				i.wizardModal.Show()
				i.state.currentState = CVStateConfiguring
				return nil
			}
			// Progress modal doesn't handle other keys
		}

		// Tier 3: Screen delegation - Review screen (metadata/stats view)
		if i.wizardReviewScreen != nil && i.state.currentState == CVStateReview {
			cmd, result := i.wizardReviewScreen.Update(msg)
			if result != nil {
				return i.handleReviewScreenResult(result)
			}
			return cmd
		}

		// Tier 4: Screen delegation - Preview screen (full scrollable content)
		if i.wizardPreviewScreen != nil && i.state.currentState == CVStatePreview {
			cmd, result := i.wizardPreviewScreen.Update(msg)
			if result != nil {
				return i.handlePreviewScreenResult(result)
			}
			return cmd
		}

		// Tier 4: Export complete state handling
		if i.state.currentState == CVStateExportComplete {
			switch keyMsg.String() {
			case "enter":
				// Complete workflow with export info
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
					},
				}
				i.active = false
				return nil
			case "esc":
				// Go back to export location selection to retry
				i.state.currentState = CVStateExportSelectLocation
				i.state.exportError = nil
				return nil
			}
		}
	}

	return nil
}

// handleWizardComplete processes wizard completion by starting tech extraction.
func (i *GenerateCVIntent) handleWizardComplete(msg WizardCompleteMsg) tea.Cmd {
	// Mark wizard as completed and hide it
	i.wizardModal.Complete()

	// Store selected profile (Step 1: WHO)
	for _, p := range i.context.AvailableProfiles {
		if p.ID == msg.ProfileID {
			i.state.selectedProfile = p
			break
		}
	}
	i.state.selectedAudience = msg.Audience

	// Store technology selections (Step 2: TECH)
	i.state.selectedTechnologyFocus = cv.TechnologyFocus(msg.TechFocus)
	i.state.selectedTechnologies = msg.Technologies
	i.state.selectedFocusArea = cv.FocusArea(msg.FocusArea)

	// Store format selections (Step 3: FORMAT)
	i.state.selectedSkillsFormat = msg.SkillsFormat
	i.state.selectedSkillsLimit = msg.SkillsLimit
	i.state.selectedCVLength = msg.CVLength

	// Show progress modal for tech extraction
	termInfo := i.GetTerminalInfo()
	width, height := termInfo.Width, termInfo.Height
	i.progressModal = components.NewExtractingTechsProgress(width, height)
	i.progressModal.Show()

	// Set state and start extraction
	i.state.currentState = CVStateExtracting
	return i.extractTechnologiesAsync()
}

// handleTechExtracted processes tech extraction completion by starting CV generation.
func (i *GenerateCVIntent) handleTechExtracted(msg TechnologiesExtractedMsg) tea.Cmd {
	// Hide progress modal
	if i.progressModal != nil {
		i.progressModal.Hide()
	}

	// Store extracted technologies
	i.state.extractedTechnologies = msg.Technologies
	i.state.focusAreaSuggestion = msg.Suggestion

	// Show progress modal for CV generation
	termInfo := i.GetTerminalInfo()
	width, height := termInfo.Width, termInfo.Height
	profileName := "Default Profile"
	if i.state.selectedProfile != nil {
		profileName = i.state.selectedProfile.Name
	}
	i.progressModal = components.NewGeneratingCVProgress(profileName, i.state.selectedAudience, width, height)
	i.progressModal.Show()

	// Set state and start generation
	i.state.currentState = CVStateGenerating
	return i.generateCVAsync()
}

// handleCVGenerated processes CV generation completion by showing review screen.
func (i *GenerateCVIntent) handleCVGenerated(msg CVGenerationCompleteMsg) tea.Cmd {
	// Hide progress modal
	if i.progressModal != nil {
		i.progressModal.Hide()
	}

	// Store generated CV
	i.state.generatedCV = msg.CV

	// Always go to review first, then preview
	if i.context.ReviewScreenFactory != nil {
		i.wizardReviewScreen = i.context.ReviewScreenFactory(msg.CV)
	}
	i.state.currentState = CVStateReview
	return nil
}

// handleReviewScreenResult processes review screen results.
func (i *GenerateCVIntent) handleReviewScreenResult(result screens.ScreenResult) tea.Cmd {
	switch result.Type() {
	case screens.ResultNavigate:
		switch result.Data() {
		case "preview":
			// User wants to see full preview
			// Create preview screen using the PreviewScreenFactory if available
			if i.context.PreviewScreenFactory != nil {
				i.wizardPreviewScreen = i.context.PreviewScreenFactory(i.state.generatedCV)
			}
			i.state.currentState = CVStatePreview
			return nil

		case "export":
			// User wants to export directly from review
			termInfo := i.GetTerminalInfo()
			width, height := termInfo.Width, termInfo.Height
			i.exportModal = components.NewExportOptionsModal(width, height)
			i.exportModal.Show()
			i.state.currentState = CVStateExporting
			return i.exportModal.Init()

		case "edit":
			// User wants to edit - go back to wizard
			i.state.currentState = CVStateConfiguring
			i.wizardModal.Reset()
			return i.wizardModal.Init()
		}

	case screens.ResultCancel:
		// User cancelled - go back to wizard with preserved data
		i.state.currentState = CVStateConfiguring
		i.wizardModal.Reset()
		return i.wizardModal.Init()
	}
	return nil
}

// handlePreviewScreenResult processes preview screen results.
func (i *GenerateCVIntent) handlePreviewScreenResult(result screens.ScreenResult) tea.Cmd {
	switch result.Type() {
	case screens.ResultSubmit:
		// User completed CV - mark intent as complete
		i.setCompleted()
		return nil

	case screens.ResultNavigate:
		switch result.Data() {
		case "export":
			// User wants to export
			termInfo := i.GetTerminalInfo()
			width, height := termInfo.Width, termInfo.Height
			i.exportModal = components.NewExportOptionsModal(width, height)
			i.exportModal.Show()
			i.state.currentState = CVStateExporting
			return i.exportModal.Init()

		case "edit":
			// User wants to edit - go back to wizard
			i.state.currentState = CVStateConfiguring
			i.wizardModal.Reset()
			return i.wizardModal.Init()
		}

	case screens.ResultCancel:
		// User cancelled from preview - go back to review
		i.state.currentState = CVStateReview
		return nil
	}
	return nil
}

// handleExportComplete processes export completion.
func (i *GenerateCVIntent) handleExportComplete(exportData *components.ExportData) tea.Cmd {
	// Hide export modal
	i.exportModal.Hide()

	// Show progress modal for export
	termInfo := i.GetTerminalInfo()
	width, height := termInfo.Width, termInfo.Height
	i.progressModal = components.NewExportingProgress(exportData.Format, width, height)
	i.progressModal.Show()

	switch exportData.Format {
	case "text":
		i.state.selectedExportFormat = CVExportFormatText
	case "markdown":
		i.state.selectedExportFormat = CVExportFormatMarkdown
	case "yaml":
		i.state.selectedExportFormat = CVExportFormatYAML
	}

	switch exportData.Location {
	case "file":
		i.state.selectedExportOption = CVExportOptionSaveToFile
	case "clipboard":
		i.state.selectedExportOption = CVExportOptionClipboard
	}

	i.state.currentState = CVStateExporting
	return i.exportCVAsync()
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

			TechnologyFocus:      string(i.state.selectedTechnologyFocus),
			SelectedTechnologies: i.state.selectedTechnologies,
			FocusArea:            string(i.state.selectedFocusArea),
			LengthFormat:         string(cv.MapUILengthToFormat(i.state.selectedCVLength)),
			SkillsFormat:         i.state.selectedSkillsFormat,
			SkillsLimit:          i.state.selectedSkillsLimit,
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

// extractTechnologiesAsync extracts technologies from user skills asynchronously.
func (i *GenerateCVIntent) extractTechnologiesAsync() tea.Cmd {
	return func() tea.Msg {
		// Check if repositories are available
		if i.context.SkillRepository == nil || i.context.EventRepository == nil {
			// Return empty results if repositories not configured (for testing)
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
			ctx = context.Background()
		}

		// Create extractor and extract technologies
		extractor := technology.NewExtractor(i.context.SkillRepository, i.context.EventRepository)
		techs, err := extractor.ExtractFromUser(ctx)
		if err != nil {
			return TechnologiesExtractedMsg{
				Technologies: nil,
				Suggestion:   nil,
				Error:        err,
			}
		}

		// Filter to skills with 3+ events
		filtered := extractor.FilterByThreshold(techs, 3)

		// Analyze skills to suggest focus area
		analyzer := &technology.Analyzer{}
		suggestion := analyzer.AnalyzeSkills(filtered)

		return TechnologiesExtractedMsg{
			Technologies: filtered,
			Suggestion:   suggestion,
			Error:        nil,
		}
	}
}

// View produces the terminal UI string for the intent's current state,
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (i *GenerateCVIntent) View() string {
	if !i.active {
		return "GenerateCV intent is not active"
	}

	return i.wizardView()
}

// wizardView renders the wizard-based workflow with modal overlays.
func (i *GenerateCVIntent) wizardView() string {
	// Get terminal dimensions
	termInfo := i.GetTerminalInfo()
	width, height := termInfo.Width, termInfo.Height

	// Create StandardView with breadcrumbs
	breadcrumbs := i.getWizardBreadcrumbs()
	view := CreateStandardViewWithBreadcrumbs(i.BaseIntent, breadcrumbs...)

	// Render content based on current state
	var content string
	switch i.state.currentState {
	case CVStateReview:
		// For review state, let the review screen render itself
		if i.wizardReviewScreen != nil {
			return i.renderReviewScreenWithModalOverlay(width, height)
		}
		content = "Loading review..."
	case CVStatePreview:
		if i.wizardPreviewScreen != nil {
			return i.renderPreviewScreenWithModalOverlay(width, height)
		}
		content = "Loading preview..."
	case CVStateExporting:
		content = "Exporting CV..."
	case CVStateExportComplete:
		content = i.viewExportComplete()
	default:
		content = ""
	}

	// Configure view
	view.WithContent(content)
	view.WithHelp(i.getWizardContextHelp())
	view.WithFooterSeparator(true)

	// Render base view
	baseView := view.Render()

	// Overlay visible modal (LAST STEP - highest priority renders last)
	// Order matters: wizard > progress > export
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

// getWizardBreadcrumbs returns breadcrumbs for wizard workflow states.
func (i *GenerateCVIntent) getWizardBreadcrumbs() []string {
	crumbs := []string{"Main Menu", "Generate CV"}

	switch i.state.currentState {
	case CVStateConfiguring:
		crumbs = append(crumbs, "Configure")
	case CVStateExtracting:
		crumbs = append(crumbs, "Extracting Technologies")
	case CVStateGenerating:
		crumbs = append(crumbs, "Generating CV")
	case CVStateReview:
		crumbs = append(crumbs, "Review")
	case CVStatePreview:
		crumbs = append(crumbs, "Preview")
	case CVStateExporting:
		crumbs = append(crumbs, "Exporting")
	}

	return crumbs
}

// getWizardContextHelp returns context-aware help for wizard workflow.
func (i *GenerateCVIntent) getWizardContextHelp() string {
	// If modal is visible, it provides its own help
	if i.wizardModal != nil && i.wizardModal.IsVisible() {
		return ""
	}
	if i.exportModal != nil && i.exportModal.IsVisible() {
		return ""
	}
	if i.progressModal != nil && i.progressModal.IsVisible() {
		return "⏳ Please wait   q Quit   m Main Menu"
	}

	// State-specific help
	switch i.state.currentState {
	case CVStateReview:
		return "Enter/p Preview   x Export   e Edit   Esc Back   q Quit"
	case CVStatePreview:
		return "↑↓ Scroll   Enter/y Confirm   x Export   Esc Back   q Quit"
	default:
		return "q Quit   m Main Menu"
	}
}

// renderWizardModalOverlay renders the wizard modal over the base view.
func (i *GenerateCVIntent) renderWizardModalOverlay(baseView string, width, height int) string {
	if i.wizardModal == nil {
		return baseView
	}

	// Use UIKit containers.Overlay for proper modal compositing
	modalView := i.wizardModal.View()
	return containers.NewOverlay(width, height).Content(modalView).Dimmed().Render()
}

// renderProgressModalOverlay renders the progress modal over the base view.
func (i *GenerateCVIntent) renderProgressModalOverlay(baseView string, width, height int) string {
	if i.progressModal == nil {
		return baseView
	}

	// Use UIKit containers.Overlay for proper modal compositing
	modalView := i.progressModal.View()
	return containers.NewOverlay(width, height).Content(modalView).Dimmed().Render()
}

// renderExportModalOverlay renders the export modal over the base view.
func (i *GenerateCVIntent) renderExportModalOverlay(baseView string, width, height int) string {
	if i.exportModal == nil {
		return baseView
	}

	// Use UIKit containers.Overlay for proper modal compositing
	modalView := i.exportModal.View()
	return containers.NewOverlay(width, height).Content(modalView).Dimmed().Render()
}

// renderReviewScreenWithModalOverlay renders the review screen and overlays export modal if visible.
func (i *GenerateCVIntent) renderReviewScreenWithModalOverlay(width, height int) string {
	// Create StandardView with breadcrumbs
	breadcrumbs := i.getWizardBreadcrumbs()
	view := CreateStandardViewWithBreadcrumbs(i.BaseIntent, breadcrumbs...)

	// Get content from review screen
	content := i.wizardReviewScreen.View()
	view.WithContent(content)

	// Get context-aware help
	help := "Enter/p Preview   x Export   e Edit   Esc Back   q Quit"
	view.WithHelp(help).WithFooterSeparator(true)

	// Render base view with StandardView
	baseView := view.Render()

	// Overlay export modal if visible
	if i.exportModal != nil && i.exportModal.IsVisible() {
		return i.renderExportModalOverlay(baseView, width, height)
	}

	return baseView
}

// renderPreviewScreenWithModalOverlay renders the preview screen and overlays export modal if visible.
func (i *GenerateCVIntent) renderPreviewScreenWithModalOverlay(width, height int) string {
	// Create StandardView with breadcrumbs
	breadcrumbs := i.getWizardBreadcrumbs()
	view := CreateStandardViewWithBreadcrumbs(i.BaseIntent, breadcrumbs...)

	// Get content from preview screen
	content := i.wizardPreviewScreen.View()
	view.WithContent(content)

	// Get context-aware help
	help := "↑↓/jk Scroll   g/G Top/Bottom   Enter/y Confirm   x Export   Esc Back   q Quit"
	view.WithHelp(help).WithFooterSeparator(true)

	// Render base view with StandardView
	baseView := view.Render()

	// Overlay export modal if visible
	if i.exportModal != nil && i.exportModal.IsVisible() {
		return i.renderExportModalOverlay(baseView, width, height)
	}

	return baseView
}

// Result retrieves the outcome of the intent after it becomes inactive,
//
// Returns:
//   - A fully initialized IntentResult[interface{}] ready for use.
//
// Side effects:
//   - None.
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
			"audience":    i.state.selectedAudience,
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

// exportCVAsync exports the CV asynchronously
func (i *GenerateCVIntent) exportCVAsync() tea.Cmd {
	return func() tea.Msg {
		// Check if export service is available
		if i.context.ExportService == nil {
			return CVExportCompleteMsg{Path: "", Error: errors.New("export service not available")}
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
			return CVExportCompleteMsg{Path: "", Error: errors.New("unknown export format")}
		}

		if err != nil {
			return CVExportCompleteMsg{Path: "", Error: fmt.Errorf("failed to export: %w", err)}
		}

		// Save based on option
		switch i.state.selectedExportOption {
		case CVExportOptionSaveToFile:
			path, err := i.context.ExportService.SaveToFile(ctx, i.state.generatedCV.Name, exportFormat, content)
			if err != nil {
				return CVExportCompleteMsg{Path: "", Error: fmt.Errorf("failed to save file: %w", err)}
			}
			return CVExportCompleteMsg{Path: path, Error: nil}

		case CVExportOptionClipboard:
			err := i.context.ExportService.CopyToClipboard(ctx, content)
			if err != nil {
				return CVExportCompleteMsg{Path: "", Error: fmt.Errorf("failed to copy to clipboard: %w", err)}
			}
			return CVExportCompleteMsg{Path: "clipboard", Error: nil}

		default:
			return CVExportCompleteMsg{Path: "", Error: errors.New("unknown export option")}
		}
	}
}

// viewExportComplete renders the export completion view using UIKit primitives.
func (i *GenerateCVIntent) viewExportComplete() string {
	th := theme.Default()
	var content strings.Builder

	if i.state.exportError != nil {
		// Use UIKit ErrorText for error header
		errorHeader := primitives.ErrorText("Export Failed", th).Bold().Render()
		content.WriteString("\n" + errorHeader + "\n\n")
		content.WriteString(fmt.Sprintf("Error: %v\n\n", i.state.exportError))
		content.WriteString("Try a different location or format.\n")
	} else {
		// Use UIKit SuccessText for success header
		successHeader := primitives.SuccessText("Export Complete!", th).Bold().Render()
		content.WriteString("\n" + successHeader + "\n\n")

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
