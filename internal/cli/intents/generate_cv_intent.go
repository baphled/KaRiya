package intents

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/service/career/technology"
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

	// Screen orchestration (Phase 2.2 TUI Architecture Refactoring)
	// When non-nil, this Screen handles Update/View for the current state.
	// Allows gradual migration from monolithic intent to screen-based architecture.
	activeScreen screens.Screen

	// useScreens enables the new screen-based architecture (opt-in for now)
	// Set to false to use legacy code and pass existing tests
	// TODO: Remove this flag once all states are migrated and tests updated
	useScreens bool
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
		useScreens: false, // Disabled by default to maintain backward compatibility
	}, nil
}

// Init is called when the intent is activated.
func (i *GenerateCVIntent) Init() tea.Cmd {
	if i.context.DefaultProfile != nil {
		i.state.selectedProfile = i.context.DefaultProfile
	}

	// Initialize active screen based on current state (Phase 2.2)
	// Only if screen-based architecture is enabled
	if i.useScreens {
		i.transitionToScreen(i.state.currentState)
	}

	return nil
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
// This avoids import cycle by creating BaseSelectScreen directly here.
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

	// Delegate to active screen if present (Phase 2.2 screen orchestration)
	// Only if useScreens is enabled
	if i.useScreens && i.activeScreen != nil {
		cmd, result := i.activeScreen.Update(msg)
		if result != nil {
			return i.handleScreenResult(result)
		}
		return cmd
	}

	// Legacy state machine (default for backward compatibility)
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
			i.activeScreen = nil // Clear screen to use legacy code for now
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
func (i *GenerateCVIntent) handleCancelResult(result screens.ScreenResult) tea.Cmd {
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
func (i *GenerateCVIntent) handleSubmitResult(result screens.ScreenResult) tea.Cmd {
	// Most screens use Navigate instead of Submit for now
	// This will be used more when we add form-based screens
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
		case " ": // Space to toggle
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

			// TODO: Implement length format selection UI
			// For now, default to Standard and proceed to generation
			i.state.selectedLengthFormat = cv.LengthStandard
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
	case GenerateCVStateExtractingTechnologies:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("...", "Please wait"),
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
				components.NewKeyBadge("Space", "Toggle"),
				components.NewKeyBadge("Enter", "Confirm"),
			),
			ThemedNavigationFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case GenerateCVStateSelectFocusArea:
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

	// Delegate to active screen if present (Phase 2.2 screen orchestration)
	// Only if useScreens is enabled
	if i.useScreens && i.activeScreen != nil {
		return i.activeScreen.View()
	}

	// Legacy view rendering (default for backward compatibility)
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
		if !i.state.technologiesAvailable && (option.focus == cv.TechnologyFocusGeneralist || option.focus == cv.TechnologyFocusSpecialist) {
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
		evidenceKeys []string // Categories that map to this area
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

		// Suggested indicator
		suggested := ""
		if i.state.focusAreaSuggestion != nil && option.area == cv.FocusArea(i.state.focusAreaSuggestion.Area) {
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
