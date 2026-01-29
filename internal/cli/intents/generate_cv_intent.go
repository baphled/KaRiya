package intents

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/service/career/technology"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

	// Screen orchestration (Phase 2.2 TUI Architecture Refactoring)
	// When non-nil, this Screen handles Update/View for the current state.
	// Allows gradual migration from monolithic intent to screen-based architecture.
	activeScreen screens.Screen

	// useScreens enables the new screen-based architecture (opt-in for now).
	// Set to false to use legacy code and pass existing tests.
	// NOTE: Feature flag to be removed once all states are migrated.
	useScreens bool

	// Wizard-based workflow (Phase 5 - Task 43)
	// Feature flag to enable new modal-based workflow
	useWizardFlow bool

	// Modals for wizard workflow
	wizardModal   *components.CVConfigWizardModal
	progressModal *components.CVProgressModal
	exportModal   *components.ExportOptionsModal

	// Screens for wizard workflow (using base Screen interface to avoid import cycle)
	// Will be *cvscreens.ReviewScreen and *cvscreens.CVPreviewScreen at runtime
	wizardReviewScreen  screens.Screen
	wizardPreviewScreen screens.Screen
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
			context:         context,
			currentState:    GenerateCVStateSelectProfile,
			selectedProfile: selectedProfile,
			selectedIndex:   0,
		},
		active:     true,
		logger:     nil,
		useScreens: false,
		// Disabled by default for test compatibility.
		// PRODUCTION: Enable via EnableWizardFlow() (see app.go line 616).
		useWizardFlow: false,
		// IMPORTANT: The wizard flow is the RECOMMENDED approach and is enabled by default in production.
		// Legacy 17-state flow is DEPRECATED and maintained only for backward compatibility with existing tests.
		// See deprecation comments at lines 669-2264 (update methods) and 1672-2475 (view methods).
	}, nil
}

// Init is called when the intent is activated.
func (i *GenerateCVIntent) Init() tea.Cmd {
	if i.context.DefaultProfile != nil {
		i.state.selectedProfile = i.context.DefaultProfile
	}

	// Wizard-based workflow (Phase 5 - Task 43)
	if i.useWizardFlow {
		return i.initWizardFlow()
	}

	// Initialize active screen based on current state (Phase 2.2)
	// Only if screen-based architecture is enabled
	if i.useScreens {
		i.transitionToScreen(i.state.currentState)
	}

	return nil
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

// transitionToScreen creates and activates a screen for the given state.
// This is part of the incremental migration to screen-based architecture.
func (i *GenerateCVIntent) transitionToScreen(state GenerateCVState) {
	termInfo := i.GetTerminalInfo()
	width, height := termInfo.Width, termInfo.Height

	switch state {
	case GenerateCVStateSelectProfile:
		screen := NewCVProfileSelectScreenFromIntent(i.context.AvailableProfiles)
		if screen != nil {
			screen.SetTerminalInfo(width, height)
			screen.SetTheme(i.Theme())
			// Set logo if available
			if logo := i.GetLogo(); logo != nil {
				screen.SetLogo(logo, i.GetLogoSpacing())
			}
		}
		i.activeScreen = screen

	// Other states will be added incrementally
	// case GenerateCVStateSelectAudience:
	// case GenerateCVStateGenerating:
	// case GenerateCVStatePreview:

	default:
		// State not yet migrated to screens - use legacy code
		i.activeScreen = nil
	}
}

// NewCVProfileSelectScreenFromIntent creates a CV profile select screen.
// This avoids import cycle by creating SelectScreen directly here.
func NewCVProfileSelectScreenFromIntent(profiles []*CVProfile) screens.Screen {
	// Create item renderer for CV profiles
	renderer := func(item *CVProfile) string {
		// Build profile display
		lines := []string{
			item.Name,
			fmt.Sprintf("  Role: %s | Audience: %s", item.TargetRole, item.TargetAudience),
		}

		if item.Description != "" {
			lines = append(lines, fmt.Sprintf("  %s", item.Description))
		}

		return strings.Join(lines, "\n")
	}

	breadcrumbs := []string{"Main Menu", "Generate CV", "Select Profile"}
	title := "Select CV Profile"

	baseScreen := base.NewBaseSelectScreen(
		profiles,
		renderer,
		breadcrumbs,
		title,
	)

	return baseScreen
}

// Theme helper methods for consistent themed styling.

// getCardStyle returns a themed card style, with fallback to default styling.
// getTheme returns the theme or a default.
func (i *GenerateCVIntent) getTheme() themes.Theme {
	if theme := i.Theme(); theme != nil {
		return theme
	}
	return themes.NewDefaultTheme()
}

func (i *GenerateCVIntent) getCardStyle() lipgloss.Style {
	return i.getTheme().Styles().CardBase
}

// getPrimaryColor returns the primary text color from theme.
func (i *GenerateCVIntent) getPrimaryColor() lipgloss.Color {
	return i.getTheme().ForegroundColor()
}

// getAccentColor returns the accent color from theme.
func (i *GenerateCVIntent) getAccentColor() lipgloss.Color {
	return i.getTheme().PrimaryColor()
}

// getErrorColor returns the error color from theme.
func (i *GenerateCVIntent) getErrorColor() lipgloss.Color {
	return i.getTheme().ErrorColor()
}

// getBorderColor returns the border color from theme.
func (i *GenerateCVIntent) getBorderColor() lipgloss.Color {
	return i.getTheme().BorderColor()
}

// Update processes a message in the intent.
func (i *GenerateCVIntent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	// Wizard-based workflow (Phase 5 - Task 43)
	// 3-tier priority: Modal → Global → Screen
	if i.useWizardFlow {
		return i.updateWizardFlow(msg)
	}

	// Delegate to active screen if present (Phase 2.2 screen orchestration)
	// Only if useScreens is enabled
	if i.useScreens && i.activeScreen != nil {
		cmd, result := i.activeScreen.Update(msg)
		if result != nil {
			return i.handleScreenResult(result)
		}
		return cmd
	}

	// DEPRECATED: Legacy state machine (maintained for backward compatibility with tests only)
	// This code path is only reached when useWizardFlow=false. The wizard workflow is now
	// the default (see updateWizardFlow). This code will be removed in a future release.
	switch i.state.currentState {
	case GenerateCVStateSelectProfile:
		return i.updateSelectProfile(msg)
	case GenerateCVStateSelectAudience:
		return i.updateSelectAudience(msg)
	case GenerateCVStateExtractingTechnologies:
		return i.updateExtractingTechnologies(msg)
	case GenerateCVStateSelectTechnologyFocus:
		return i.updateSelectTechnologyFocus(msg)
	case GenerateCVStateSelectTechnologies:
		return i.updateSelectTechnologies(msg)
	case GenerateCVStateSelectFocusArea:
		return i.updateSelectFocusArea(msg)
	case GenerateCVStateSelectSkillsConfig:
		return i.updateSelectSkillsConfig(msg)
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

// updateWizardFlow handles wizard-based workflow updates with 3-tier priority.
// Tier 1 (Highest): Modal updates
// Tier 2: Global keys (quit, help, main menu)
// Tier 3: Screen updates
func (i *GenerateCVIntent) updateWizardFlow(msg tea.Msg) tea.Cmd {
	// Handle window size messages
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		// Update BaseIntent terminal info
		termInfo := i.BaseIntent.GetTerminalInfo()
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

	// Handle export complete (wizard flow uses same handler as legacy)
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
			i.state.currentState = GenerateCVStateExportComplete
			return nil
		}
		i.state.exportedPath = msg.Path
		i.state.currentState = GenerateCVStateExportComplete
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
		if i.state.currentState == GenerateCVStateExportComplete {
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
				i.state.currentState = GenerateCVStateExportSelectLocation
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
	termInfo := i.BaseIntent.GetTerminalInfo()
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
	termInfo := i.BaseIntent.GetTerminalInfo()
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
			termInfo := i.BaseIntent.GetTerminalInfo()
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
			termInfo := i.BaseIntent.GetTerminalInfo()
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
	termInfo := i.BaseIntent.GetTerminalInfo()
	width, height := termInfo.Width, termInfo.Height
	i.progressModal = components.NewExportingProgress(exportData.Format, width, height)
	i.progressModal.Show()

	// Start export async by transitioning to legacy export states
	// Map modal selection to legacy state values
	switch exportData.Format {
	case "text":
		i.state.selectedExportFormat = CVExportFormatText
	case "markdown":
		i.state.selectedExportFormat = CVExportFormatMarkdown
	case "yaml":
		i.state.selectedExportFormat = CVExportFormatYAML
	}

	// Map location to export option
	switch exportData.Location {
	case "file":
		i.state.selectedExportOption = CVExportOptionSaveToFile
	case "clipboard":
		i.state.selectedExportOption = CVExportOptionClipboard
	}

	// Trigger export process (will use legacy export workflow)
	i.state.currentState = GenerateCVStateExporting
	return i.exportCVAsync()
}

// handleScreenResult processes a ScreenResult from the active screen.
// This is the bridge between screen-based UI and intent-based workflow orchestration.
func (i *GenerateCVIntent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
	switch result.Type() {
	case screens.ResultNavigate:
		// User selected something and wants to proceed
		// Extract the data and transition to the next state
		return i.handleNavigateResult(result)

	case screens.ResultCancel:
		// User pressed Escape - go back to previous state
		return i.handleCancelResult(result)

	case screens.ResultSubmit:
		// User submitted a form or completed an action
		return i.handleSubmitResult(result)

	case screens.ResultError:
		// An error occurred in the screen
		return i.handleErrorResult(result)
	}

	return nil
}

// handleNavigateResult processes a NavigateResult from a screen.
func (i *GenerateCVIntent) handleNavigateResult(result screens.ScreenResult) tea.Cmd {
	data := result.Data()

	switch i.state.currentState {
	case GenerateCVStateSelectProfile:
		// User selected a profile
		if profile, ok := data.(*CVProfile); ok {
			i.state.selectedProfile = profile
			i.state.currentState = GenerateCVStateSelectAudience
			i.activeScreen = nil
		}
		return nil

	case GenerateCVStateSelectAudience:
		// User selected an audience
		if audience, ok := data.(string); ok {
			i.state.selectedAudience = audience
			i.state.currentState = GenerateCVStateGenerating
			i.state.isGenerating = true
			i.activeScreen = nil
			return i.generateCVAsync()
		}
		return nil

	case GenerateCVStateGenerating:
		// CV generation complete, move to preview
		if cv, ok := data.(*career.CVView); ok {
			i.state.generatedCV = cv
			i.state.currentState = GenerateCVStatePreview
			i.activeScreen = nil
		}
		return nil

	case GenerateCVStatePreview:
		// User finished previewing, move to review
		i.state.currentState = GenerateCVStateReview
		i.activeScreen = nil
		return nil
	}

	return nil
}

// handleCancelResult processes a CancelResult from a screen.
func (i *GenerateCVIntent) handleCancelResult(_ screens.ScreenResult) tea.Cmd {
	switch i.state.currentState {
	case GenerateCVStateSelectProfile:
		// Root state - cancel the intent
		i.setCancelled()
		return nil

	case GenerateCVStateSelectAudience:
		// Go back to profile selection
		i.state.currentState = GenerateCVStateSelectProfile
		i.activeScreen = nil
		return nil

	case GenerateCVStateGenerating:
		// Go back to audience selection
		i.state.currentState = GenerateCVStateSelectAudience
		i.activeScreen = nil
		return nil

	case GenerateCVStatePreview:
		// Go back to audience selection (regenerate)
		i.state.currentState = GenerateCVStateSelectAudience
		i.activeScreen = nil
		return nil
	}

	return nil
}

// handleSubmitResult processes a SubmitResult from a screen.
func (i *GenerateCVIntent) handleSubmitResult(_ screens.ScreenResult) tea.Cmd {
	// Most screens use Navigate instead of Submit for now.
	// This will be used more when we add form-based screens.
	return nil
}

// handleErrorResult processes an ErrorResult from a screen.
func (i *GenerateCVIntent) handleErrorResult(result screens.ScreenResult) tea.Cmd {
	data := result.Data()
	if errData, ok := data.(map[string]interface{}); ok {
		if err, ok := errData["error"].(error); ok {
			i.state.generationError = err
		}
	}
	return nil
}

// ============================================================================
// LEGACY UPDATE METHODS (DEPRECATED)
// ============================================================================
// The following update methods implement the legacy 17-state workflow.
//
// DEPRECATED: This code is maintained for backward compatibility with existing
// tests only. The wizard-based workflow (useWizardFlow=true) is now the default
// and recommended approach. See updateWizardFlow() for the current implementation.
//
// These methods will be removed in a future release after all tests are migrated
// to the wizard workflow. Do NOT use these methods in new code.
//
// Legacy update methods: lines 669-2264 (~1595 lines)
// - updateSelectProfile
// - updateSelectAudience
// - updateExtractingTechnologies
// - updateSelectTechnologyFocus
// - updateSelectTechnologies
// - updateSelectFocusArea
// - updateSelectSkillsConfig
// - updateGenerating
// - updatePreview
// - updateReview
// - updateConfirm
// - updateExportSelectFormat
// - updateExportSelectLocation
// - updateExporting
// - updateExportComplete
// ============================================================================

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
			// Transition to technology extraction
			i.state.currentState = GenerateCVStateExtractingTechnologies
			return i.extractTechnologiesAsync()
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

			// Technology selections (Phase 8 - Task 43 wizard integration)
			TechnologyFocus:      string(i.state.selectedTechnologyFocus),
			SelectedTechnologies: i.state.selectedTechnologies,
			FocusArea:            string(i.state.selectedFocusArea),
			LengthFormat:         string(cv.MapUILengthToFormat(i.state.selectedCVLength)),

			// Skills section configuration
			SkillsFormat: i.state.selectedSkillsFormat,
			SkillsLimit:  i.state.selectedSkillsLimit,
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

// updateExtractingTechnologies handles messages while extracting technologies.
func (i *GenerateCVIntent) updateExtractingTechnologies(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case TechnologiesExtractedMsg:
		if msg.Error != nil {
			// Error extracting - go back to audience selection
			i.state.currentState = GenerateCVStateSelectAudience
			return nil
		}

		// Store extracted technologies and suggestion
		i.state.extractedTechnologies = msg.Technologies
		i.state.focusAreaSuggestion = msg.Suggestion
		i.state.technologiesAvailable = len(msg.Technologies) >= 3

		// Transition to technology focus selection
		i.state.currentState = GenerateCVStateSelectTechnologyFocus
		i.state.technologyFocusIndex = 0
		return nil

	case tea.KeyMsg:
		// Handle global keys (allow cancellation during extraction)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyBack:
			// Let extraction complete in background, navigate back
			i.state.currentState = GenerateCVStateSelectAudience
			return nil
		}
	}
	return nil
}

// updateSelectTechnologyFocus handles messages while selecting technology focus.
func (i *GenerateCVIntent) updateSelectTechnologyFocus(msg tea.Msg) tea.Cmd {
	// Available technology focus options
	options := []cv.TechnologyFocus{
		cv.TechnologyFocusLanguageAgnostic,
		cv.TechnologyFocusGeneralist,
		cv.TechnologyFocusSpecialist,
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if i.state.technologyFocusIndex > 0 {
				i.state.technologyFocusIndex--
			}
			return nil
		case "down", "j":
			if i.state.technologyFocusIndex < len(options)-1 {
				i.state.technologyFocusIndex++
			}
			return nil
		case "enter":
			// Set selected technology focus
			i.state.selectedTechnologyFocus = options[i.state.technologyFocusIndex]

			// Transition based on selection
			switch i.state.selectedTechnologyFocus {
			case cv.TechnologyFocusLanguageAgnostic:
				// Skip technology selection, go to focus area
				i.state.currentState = GenerateCVStateSelectFocusArea
				i.state.focusAreaCursor = 0
			case cv.TechnologyFocusGeneralist, cv.TechnologyFocusSpecialist:
				// Go to technology selection (multi or single select)
				i.state.currentState = GenerateCVStateSelectTechnologies
				i.state.technologyCursor = 0
				i.state.technologySelected = make(map[int]bool)
				i.state.selectedTechnologies = []string{}
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
			// Go back to extracting technologies (or previous state)
			i.state.currentState = GenerateCVStateExtractingTechnologies
			return nil
		}
	}
	return nil
}

// updateSelectTechnologies handles messages while selecting technologies (multi or single select).
func (i *GenerateCVIntent) updateSelectTechnologies(msg tea.Msg) tea.Cmd {
	if len(i.state.extractedTechnologies) == 0 {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if i.state.technologyCursor > 0 {
				i.state.technologyCursor--
			}
			return nil
		case "down", "j":
			if i.state.technologyCursor < len(i.state.extractedTechnologies)-1 {
				i.state.technologyCursor++
			}
			return nil
		case " ":
			// Determine if multi-select or single-select
			if i.state.selectedTechnologyFocus == cv.TechnologyFocusSpecialist {
				// Single-select: clear all others, select current
				i.state.technologySelected = make(map[int]bool)
				i.state.technologySelected[i.state.technologyCursor] = true
			} else {
				// Multi-select: toggle current
				i.state.technologySelected[i.state.technologyCursor] = !i.state.technologySelected[i.state.technologyCursor]
			}
			return nil
		case "enter":
			// Collect selected technology IDs
			var selectedIDs []string
			for idx, selected := range i.state.technologySelected {
				if selected && idx < len(i.state.extractedTechnologies) {
					selectedIDs = append(selectedIDs, i.state.extractedTechnologies[idx].ID)
				}
			}

			// Validate selection count based on focus type
			var valid bool
			if i.state.selectedTechnologyFocus == cv.TechnologyFocusSpecialist {
				valid = len(selectedIDs) == 1
			} else if i.state.selectedTechnologyFocus == cv.TechnologyFocusGeneralist {
				valid = len(selectedIDs) >= 2 && len(selectedIDs) <= 5
			}

			if !valid {
				// Don't transition - invalid selection
				return nil
			}

			// Store selected technologies
			i.state.selectedTechnologies = selectedIDs

			// Transition to focus area selection
			i.state.currentState = GenerateCVStateSelectFocusArea
			i.state.focusAreaCursor = 0
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
			// Go back to technology focus selection
			i.state.currentState = GenerateCVStateSelectTechnologyFocus
			return nil
		}
	}
	return nil
}

// updateSelectFocusArea handles the focus area selection state.
func (i *GenerateCVIntent) updateSelectFocusArea(msg tea.Msg) tea.Cmd {
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
			// Determine where to go back based on where we came from
			if i.state.selectedTechnologyFocus == cv.TechnologyFocusLanguageAgnostic {
				// If Language Agnostic, go back to technology focus selection
				i.state.currentState = GenerateCVStateSelectTechnologyFocus
			} else {
				// If Generalist/Specialist, go back to technology selection
				i.state.currentState = GenerateCVStateSelectTechnologies
			}
			return nil
		}

		// Handle navigation and selection
		switch msg.String() {
		case "up", "k":
			if i.state.focusAreaCursor > 0 {
				i.state.focusAreaCursor--
			}
		case "down", "j":
			// 4 focus areas: Backend, Frontend, Fullstack, DevOps
			if i.state.focusAreaCursor < 3 {
				i.state.focusAreaCursor++
			}
		case "enter":
			// Map cursor position to focus area
			focusAreas := []cv.FocusArea{
				cv.FocusAreaBackend,
				cv.FocusAreaFrontend,
				cv.FocusAreaFullstack,
				cv.FocusAreaDevOps,
			}
			i.state.selectedFocusArea = focusAreas[i.state.focusAreaCursor]

			// Transition to skills configuration
			i.state.currentState = GenerateCVStateSelectSkillsConfig
			i.state.skillsConfigCursor = 0

			// Set defaults for skills config if not already set
			if i.state.selectedSkillsFormat == "" {
				i.state.selectedSkillsFormat = "flat"
			}
			if i.state.selectedSkillsLimit == 0 {
				i.state.selectedSkillsLimit = 0
			}

			return nil
		}
	}
	return nil
}

// updateSelectSkillsConfig handles the skills configuration selection state.
func (i *GenerateCVIntent) updateSelectSkillsConfig(msg tea.Msg) tea.Cmd {
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
			// Go back to focus area selection
			i.state.currentState = GenerateCVStateSelectFocusArea
			return nil
		}

		// Handle navigation and selection
		switch msg.String() {
		case "up", "k":
			if i.state.skillsConfigCursor > 0 {
				i.state.skillsConfigCursor--
			}
		case "down", "j":
			// 2 options: format (0) and limit (1)
			if i.state.skillsConfigCursor < 1 {
				i.state.skillsConfigCursor++
			}
		case " ":
			// Toggle format when cursor is on format row (0)
			if i.state.skillsConfigCursor == 0 {
				if i.state.selectedSkillsFormat == "flat" {
					i.state.selectedSkillsFormat = "grouped"
				} else {
					i.state.selectedSkillsFormat = "flat"
				}
			}
		case "left", "h":
			// Decrease limit when cursor is on limit row (1)
			if i.state.skillsConfigCursor == 1 {
				if i.state.selectedSkillsLimit > 0 {
					i.state.selectedSkillsLimit -= 5
					if i.state.selectedSkillsLimit < 0 {
						i.state.selectedSkillsLimit = 0
					}
				}
			}
		case "right", "l":
			// Increase limit when cursor is on limit row (1)
			if i.state.skillsConfigCursor == 1 {
				if i.state.selectedSkillsLimit < 50 {
					i.state.selectedSkillsLimit += 5
					if i.state.selectedSkillsLimit > 50 {
						i.state.selectedSkillsLimit = 50
					}
				}
			}
		case "enter":
			// Proceed to CV generation
			i.state.currentState = GenerateCVStateGenerating
			i.state.isGenerating = true
			return i.generateCVAsync()
		}
	}
	return nil
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
// Deprecated: This method routes to legacy view methods and is only used when
// useWizardFlow=false (for backward compatibility with tests). The wizard workflow
// uses wizardView() instead. This method will be removed in a future release.
func (i *GenerateCVIntent) getStateContent() string {
	switch i.state.currentState {
	case GenerateCVStateSelectProfile:
		return i.viewSelectProfile()
	case GenerateCVStateSelectAudience:
		return i.viewSelectAudience()
	case GenerateCVStateSelectTechnologyFocus:
		return i.viewSelectTechnologyFocus()
	case GenerateCVStateSelectTechnologies:
		return i.viewSelectTechnologies()
	case GenerateCVStateSelectFocusArea:
		return i.viewSelectFocusArea()
	case GenerateCVStateSelectSkillsConfig:
		return i.viewSelectSkillsConfig()
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
// Deprecated: This method provides help for legacy 17-state workflow and is only used when
// useWizardFlow=false (for backward compatibility with tests). The wizard workflow
// uses getWizardContextHelp() instead. This method will be removed in a future release.
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
	case GenerateCVStateExtractingTechnologies:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("...", "Please wait", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateSelectTechnologyFocus:
		return CombineThemedFooters(
			ThemedNavigationFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateSelectTechnologies:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Space", "Toggle", theme),
				primitives.HelpKeyBadge("Enter", "Confirm", theme),
			),
			ThemedNavigationFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateSelectFocusArea:
		return CombineThemedFooters(
			ThemedNavigationFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateSelectSkillsConfig:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Space", "Toggle Format", theme),
				primitives.HelpKeyBadge("←→", "Adjust Limit", theme),
				primitives.HelpKeyBadge("Enter", "Continue", theme),
			),
			ThemedNavigationFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateGenerating:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("...", "Please wait", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStatePreview:
		return CombineThemedFooters(
			ThemedDetailViewFooter(theme),
			ThemedCustomFooter(theme,
				primitives.EditBadge(theme),
				primitives.HelpKeyBadge("c", "Continue", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateReview:
		return CombineThemedFooters(
			ThemedDetailViewFooter(theme),
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Enter", "Continue", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateConfirm:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("y/Enter", "Confirm", theme),
				primitives.HelpKeyBadge("e/x", "Export", theme),
				primitives.HelpKeyBadge("n/Esc", "Back", theme),
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
				primitives.HelpKeyBadge("...", "Please wait", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateExportComplete:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Enter", "Continue", theme),
				primitives.BackBadge(theme),
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

	// Use wizard flow if enabled (Task 43 refactored architecture)
	if i.useWizardFlow {
		return i.wizardView()
	}

	// Delegate to active screen if present (Phase 2.2 screen orchestration)
	// Only if useScreens is enabled
	if i.useScreens && i.activeScreen != nil {
		return i.activeScreen.View()
	}

	// DEPRECATED: Legacy view rendering (maintained for backward compatibility with tests only)
	// This code path is only reached when useWizardFlow=false. The wizard workflow uses
	// wizardView() instead. This code will be removed in a future release.
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

// wizardView renders the wizard-based workflow with modal overlays.
func (i *GenerateCVIntent) wizardView() string {
	// Get terminal dimensions
	termInfo := i.BaseIntent.GetTerminalInfo()
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
		// For preview state, let the screen render itself fully
		if i.wizardPreviewScreen != nil {
			// Screen renders its own StandardView, so return it directly
			// without wrapping in another StandardView
			return i.renderPreviewScreenWithModalOverlay(width, height)
		}
		content = "Loading preview..."
	case CVStateExporting:
		// Note: CVStateExporting == GenerateCVStateExporting == "exporting"
		content = "Exporting CV..."
	case GenerateCVStateExportComplete:
		// Use the same export complete view as legacy flow
		content = i.viewExportComplete()
	default:
		// For other states, show minimal content (modals will overlay)
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

// getBreadcrumbs returns breadcrumbs for the current state.
func (i *GenerateCVIntent) getBreadcrumbs() []string {
	crumbs := []string{"Main Menu", "Generate CV"}

	switch i.state.currentState {
	case GenerateCVStateSelectProfile:
		crumbs = append(crumbs, "Select Profile")
	case GenerateCVStateSelectAudience:
		crumbs = append(crumbs, "Select Audience")
	case GenerateCVStateExtractingTechnologies:
		crumbs = append(crumbs, "Extracting Technologies")
	case GenerateCVStateSelectTechnologyFocus:
		crumbs = append(crumbs, "Select Technology Focus")
	case GenerateCVStateSelectTechnologies:
		crumbs = append(crumbs, "Select Technologies")
	case GenerateCVStateSelectFocusArea:
		crumbs = append(crumbs, "Select Focus Area")
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

// ============================================================================
// LEGACY VIEW METHODS (DEPRECATED)
// ============================================================================
// The following view methods implement rendering for the legacy 17-state workflow.
//
// DEPRECATED: This code is maintained for backward compatibility with existing
// tests only. The wizard-based workflow (useWizardFlow=true) is now the default
// and uses wizardView() for rendering. See wizardView() for the current implementation.
//
// These methods will be removed in a future release after all tests are migrated
// to the wizard workflow. Do NOT use these methods in new code.
//
// Legacy view methods: lines 1672-2475 (~803 lines)
// - viewSelectProfile
// - viewSelectAudience
// - viewSelectTechnologyFocus
// - viewSelectTechnologies
// - viewSelectFocusArea
// - viewSelectSkillsConfig
// - viewGenerating
// - viewPreview
// - viewReview
// - viewConfirm
// - viewExportSelectFormat
// - viewExportSelectLocation
// - viewExporting
// - viewExportComplete
// ============================================================================

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

// viewSelectTechnologyFocus renders the technology focus selection view.
func (i *GenerateCVIntent) viewSelectTechnologyFocus() string {
	var content strings.Builder

	// Header
	content.WriteString("\n🎯 Select Technology Focus\n\n")

	// Show technology count
	techCount := len(i.state.extractedTechnologies)
	if techCount > 0 {
		content.WriteString(fmt.Sprintf("Found %d technologies across your career events\n\n", techCount))
	}

	// Technology focus options
	options := []struct {
		focus       cv.TechnologyFocus
		name        string
		description string
	}{
		{cv.TechnologyFocusLanguageAgnostic, "Language Agnostic", "Technology-agnostic narrative (emphasizes adaptability)"},
		{cv.TechnologyFocusGeneralist, "Generalist (2-5 technologies)", "Highlight 2-5 key technologies"},
		{cv.TechnologyFocusSpecialist, "Specialist (1 technology)", "Focus deeply on a single technology"},
	}

	for idx, option := range options {
		prefix := "  "
		if idx == i.state.technologyFocusIndex {
			prefix = "▶ "
		}

		// Check if option is available
		disabled := ""
		isAdvancedFocus := option.focus == cv.TechnologyFocusGeneralist ||
			option.focus == cv.TechnologyFocusSpecialist
		if !i.state.technologiesAvailable && isAdvancedFocus {
			disabled = " (unavailable - requires 3+ technologies)"
		}

		content.WriteString(fmt.Sprintf("%s%s%s\n", prefix, option.name, disabled))
		content.WriteString(fmt.Sprintf("   %s\n\n", option.description))
	}

	// Show warning if < 3 technologies
	if !i.state.technologiesAvailable {
		content.WriteString("\n⚠️  Only Language Agnostic is available (requires 3+ technologies for other options)\n")
	}

	return i.getCardStyle().Render(content.String())
}

// viewSelectTechnologies renders the technology selection view (multi or single select).
func (i *GenerateCVIntent) viewSelectTechnologies() string {
	var content strings.Builder

	// Header based on mode
	if i.state.selectedTechnologyFocus == cv.TechnologyFocusSpecialist {
		content.WriteString("\n🎯 Select 1 Technology\n\n")
		content.WriteString("Choose the technology you want to specialize in:\n\n")
	} else {
		content.WriteString("\n🎯 Select 2-5 Technologies\n\n")
		content.WriteString("Choose technologies to highlight (select 2-5):\n\n")
	}

	// Technology list with checkboxes
	for idx, tech := range i.state.extractedTechnologies {
		// Cursor indicator
		cursor := "  "
		if idx == i.state.technologyCursor {
			cursor = "▶ "
		}

		// Selection checkbox
		checkbox := "☐"
		if i.state.technologySelected[idx] {
			checkbox = "☑"
		}

		// Event count
		eventInfo := fmt.Sprintf("(%d events)", tech.EventCount)

		content.WriteString(fmt.Sprintf("%s%s %s %s\n", cursor, checkbox, tech.Name, eventInfo))
	}

	// Footer with count
	selectedCount := 0
	for _, selected := range i.state.technologySelected {
		if selected {
			selectedCount++
		}
	}

	content.WriteString(fmt.Sprintf("\nSelected: %d", selectedCount))

	if i.state.selectedTechnologyFocus == cv.TechnologyFocusGeneralist {
		content.WriteString(" (need 2-5)")
	} else {
		content.WriteString(" (need 1)")
	}

	content.WriteString("\n\nSpace to toggle, Enter to confirm\n")

	return i.getCardStyle().Render(content.String())
}

// viewSelectFocusArea renders the focus area selection view.
func (i *GenerateCVIntent) viewSelectFocusArea() string {
	var content strings.Builder
	content.WriteString("\n🎯 Select Focus Area\n\n")
	content.WriteString("Choose the primary focus area for your CV:\n\n")

	// Focus areas with descriptions and evidence keys
	focusAreas := []struct {
		area         cv.FocusArea
		name         string
		description  string
		evidenceKeys []string
	}{
		{cv.FocusAreaBackend, "Backend", "Server-side, databases, APIs, infrastructure", []string{"backend", "database"}},
		{cv.FocusAreaFrontend, "Frontend", "UI/UX, web apps, client-side frameworks", []string{"frontend", "ui"}},
		{cv.FocusAreaFullstack, "Fullstack", "Both frontend and backend development", []string{"fullstack"}},
		{cv.FocusAreaDevOps, "DevOps", "CI/CD, deployment, monitoring, cloud", []string{"devops", "cloud", "infrastructure"}},
	}

	for idx, option := range focusAreas {
		// Cursor indicator
		cursor := "  "
		if idx == i.state.focusAreaCursor {
			cursor = "▶ "
		}

		// Suggested indicator.
		suggested := ""
		if i.state.focusAreaSuggestion != nil && option.area == i.state.focusAreaSuggestion.Area {
			suggested = " ⭐ (suggested)"
		}

		content.WriteString(fmt.Sprintf("%s%s%s\n", cursor, option.name, suggested))
		content.WriteString(fmt.Sprintf("   %s\n", option.description))

		// Show individual skill category counts if we have evidence
		if i.state.focusAreaSuggestion != nil {
			var evidenceParts []string
			for _, key := range option.evidenceKeys {
				if count, ok := i.state.focusAreaSuggestion.Evidence[key]; ok && count > 0 {
					evidenceParts = append(evidenceParts, fmt.Sprintf("%s: %d", key, count))
				}
			}
			if len(evidenceParts) > 0 {
				content.WriteString(fmt.Sprintf("   (%s)\n", strings.Join(evidenceParts, ", ")))
			}
		}

		content.WriteString("\n")
	}

	return i.getCardStyle().Render(content.String())
}

// viewSelectSkillsConfig renders the skills configuration selection view.
func (i *GenerateCVIntent) viewSelectSkillsConfig() string {
	var content strings.Builder
	content.WriteString("\n⚙️  Configure Skills Section\n\n")
	content.WriteString("Customize how skills appear in your CV:\n\n")

	// Option 1: Format selection
	cursor1 := "  "
	if i.state.skillsConfigCursor == 0 {
		cursor1 = "▶ "
	}

	formatCheckmark := ""
	formatDesc := ""
	if i.state.selectedSkillsFormat == "flat" {
		formatCheckmark = " ✓"
		formatDesc = " (one skill per line)"
	} else {
		formatCheckmark = " ✓"
		formatDesc = " (skills grouped by category)"
	}

	caser := cases.Title(language.English)
	content.WriteString(fmt.Sprintf("%sFormat: %s%s%s\n", cursor1,
		caser.String(i.state.selectedSkillsFormat), formatCheckmark, formatDesc))
	content.WriteString("   Press Space to toggle between Flat / Grouped\n\n")

	// Option 2: Limit selection
	cursor2 := "  "
	if i.state.skillsConfigCursor == 1 {
		cursor2 = "▶ "
	}

	limitDesc := ""
	if i.state.selectedSkillsLimit == 0 {
		limitDesc = " (no limit - show all skills)"
	} else {
		limitDesc = fmt.Sprintf(" (%d max per section/group)", i.state.selectedSkillsLimit)
	}

	content.WriteString(fmt.Sprintf("%sLimit: %d%s\n", cursor2,
		i.state.selectedSkillsLimit, limitDesc))
	content.WriteString("   Use ← → to adjust (0 = no limit, max 50)\n\n")

	content.WriteString("\nPress Enter to continue\n")

	return i.getCardStyle().Render(content.String())
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
func (i *GenerateCVIntent) viewPreview() string {
	if i.state.generatedCV == nil {
		return "No CV generated yet"
	}

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
		sectionTitle := primitives.NewText(strings.ToUpper(section.Title), i.Theme()).
			Bold().
			Foreground(i.getAccentColor()).
			MarginTop(1)
		content.WriteString(sectionTitle.Render() + "\n")
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

	// Set viewport content
	i.state.previewViewport.SetContent(content.String())

	// Render viewport
	viewportContent := i.state.previewViewport.View()

	return viewportContent
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

// EnableScreens enables the new screen-based architecture for this intent.
// This is opt-in during Phase 2.2 migration to maintain backward compatibility.
// Once all states are migrated and tests updated, this will become the default.
//
// Usage:
//
//	intent.EnableScreens()
//	intent.Init() // Initializes screens based on current state
func (i *GenerateCVIntent) EnableScreens() {
	i.useScreens = true
}

// EnableWizardFlow enables the wizard-based workflow for this intent.
// This is opt-in during Phase 5 (Task 43) migration to maintain backward compatibility.
// Once fully tested and validated, this will become the default.
//
// The wizard flow simplifies the CV generation workflow from 17 states to 5:
//   - CVStateConfiguring: Configuration wizard modal (replaces 6 states)
//   - CVStateExtracting: Technology extraction progress (replaces 1 state)
//   - CVStateGenerating: CV generation progress (replaces 1 state)
//   - CVStatePreview: Preview screen with actions (replaces 3 states)
//   - CVStateExporting: Export modal (replaces 6 states)
//
// Usage:
//
//	intent.EnableWizardFlow()
//	intent.Init() // Initializes wizard modal
func (i *GenerateCVIntent) EnableWizardFlow() {
	i.useWizardFlow = true
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
			i.state.currentState = GenerateCVStateExportComplete
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
