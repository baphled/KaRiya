package intents

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/screens"
	exportscreens "github.com/baphled/kariya/internal/cli/screens/export"
	"github.com/baphled/kariya/internal/cli/types"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	tea "github.com/charmbracelet/bubbletea"
)

// ExportArtifactIntent implements the Intent interface for artifact export.
// It uses a wizard modal for configuration followed by preview and export modals.
//
// The workflow uses:
// - Wizard modal for: Configuration (type, format, destination selection in 2 steps)
// - Screen for: Preview (scrollable content)
// - Modal overlays for: Confirm, Progress, Success, Error (dialog overlays)
//
// This matches the fluid pattern used in GenerateCV with its wizard modal.
type ExportArtifactIntent struct {
	// Embed BaseIntent for terminal awareness, logo, and state management
	*BaseIntent

	// context contains the input context with services and repositories
	context *ExportArtifactContext

	// state tracks the current state
	currentState ExportState

	// config holds the current export configuration
	config *ExportConfiguration

	// preview holds the generated preview content
	preview string

	// previewStats holds statistics about the preview
	previewStats *exportscreens.PreviewStats

	// result holds the final export result
	exportResult *ExportArtifactResult

	// intentResult is the final result of the intent
	intentResult *IntentResult[*ExportArtifactResult]

	// error holds any error that occurred
	exportError *IntentError

	// activeScreen holds the current screen being displayed
	activeScreen screens.Screen

	// Wizard modal for configuration (replaces type/format/dest screens)
	wizardModal *components.ExportConfigWizardModal

	// Modal overlays for confirm, progress, success, and error states
	confirmModal  *components.ExportConfirmModal
	progressModal *components.ExportProgressModal
	successModal  *components.ExportSuccessModal
	errorModal    *components.ExportErrorModal

	// active indicates whether this intent is currently active
	active bool
}

// NewExportArtifactIntent creates a new ExportArtifact intent.
func NewExportArtifactIntent(ctx *ExportArtifactContext) (*ExportArtifactIntent, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	baseIntent := NewBaseIntent()

	return &ExportArtifactIntent{
		BaseIntent:   baseIntent,
		context:      ctx,
		currentState: ExportStateSelectType,
		active:       true,
	}, nil
}

// Init initializes the intent and shows the configuration wizard.
func (e *ExportArtifactIntent) Init() tea.Cmd {
	e.active = true
	e.currentState = ExportStateConfigure

	// Get terminal dimensions
	width, height := 100, 40
	termInfo := e.GetTerminalInfo()
	if termInfo != nil {
		width = termInfo.Width
		height = termInfo.Height
	}

	// Create and show the configuration wizard modal
	e.wizardModal = components.NewExportConfigWizardModal(width, height)
	return e.wizardModal.Init()
}

// Update handles messages and delegates to the wizard, screen, or modal.
func (e *ExportArtifactIntent) Update(msg tea.Msg) tea.Cmd {
	if !e.active {
		return nil
	}

	// Handle window resize for wizard
	if _, ok := msg.(tea.WindowSizeMsg); ok {
		if e.wizardModal != nil && e.wizardModal.IsVisible() {
			return e.wizardModal.Update(msg)
		}
	}

	// Handle global keys first
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch HandleGlobalKeys(keyMsg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			e.ToggleHelp()
			return nil
		}
	}

	// Handle wizard modal first (highest priority when visible)
	if e.wizardModal != nil && e.wizardModal.IsVisible() {
		cmd := e.wizardModal.Update(msg)

		// Check if wizard completed or cancelled
		if e.wizardModal.IsCompleted() {
			// Get configuration from wizard
			e.config = &ExportConfiguration{
				ArtifactType: ExportArtifactType(e.wizardModal.GetArtifactType()),
				Format:       ExportFormat(e.wizardModal.GetFormat()),
				Destination:  ExportDestination(e.wizardModal.GetDestination()),
			}
			e.wizardModal = nil

			// Generate preview and show preview screen
			e.preview = e.generatePreview()
			e.currentState = ExportStatePreview
			e.transitionToScreen(e.newPreviewScreen())
			return nil
		}

		if e.wizardModal.IsCancelled() {
			e.wizardModal = nil
			e.setCancelled()
			return nil
		}

		return cmd
	}

	// Handle modal interactions (modals take priority over screens)
	if e.confirmModal != nil && e.confirmModal.IsVisible() {
		cmd, confirmed := e.confirmModal.Update(msg)
		if confirmed {
			// User confirmed - start export with progress modal
			e.confirmModal = nil
			e.showProgressModal()
			return tea.Batch(e.progressModal.Init(), e.startExport())
		} else if !e.confirmModal.IsVisible() {
			// User cancelled - go back to preview
			e.confirmModal = nil
			e.currentState = ExportStatePreview
		}
		return cmd
	}

	if e.progressModal != nil && e.progressModal.IsVisible() {
		// Handle spinner tick
		if _, ok := msg.(components.SpinnerTickMsg); ok {
			return e.progressModal.Update(msg)
		}
		// Handle export completion messages
		switch msg := msg.(type) {
		case ExportCompleteMsg:
			e.progressModal.Complete()
			e.progressModal = nil
			e.exportResult = msg.Result
			e.currentState = ExportStateComplete
			e.showSuccessModal()
			return nil
		case ExportErrorMsg:
			e.progressModal.SetError(nil)
			e.progressModal = nil
			e.exportError = msg.Error
			e.currentState = ExportStateFailed
			e.showErrorModal()
			return nil
		}
		// Handle cancellation
		cmd := e.progressModal.Update(msg)
		if e.progressModal.IsCancelled() {
			e.progressModal = nil
			e.currentState = ExportStatePreview
		}
		return cmd
	}

	if e.successModal != nil && e.successModal.IsVisible() {
		cmd, done := e.successModal.Update(msg)
		if done {
			e.successModal = nil
			e.setCompleted()
		}
		return cmd
	}

	if e.errorModal != nil && e.errorModal.IsVisible() {
		cmd, result := e.errorModal.Update(msg)
		if result == components.ErrorResultRetry {
			e.errorModal = nil
			e.showProgressModal()
			return tea.Batch(e.progressModal.Init(), e.startExport())
		} else if result == components.ErrorResultCancel {
			e.errorModal = nil
			e.setCancelled()
		}
		return cmd
	}

	// Delegate to active screen
	if e.activeScreen != nil {
		cmd, result := e.activeScreen.Update(msg)
		if result != nil {
			screenCmd := e.handleScreenResult(result)
			if screenCmd != nil {
				return tea.Batch(cmd, screenCmd)
			}
		}
		return cmd
	}

	return nil
}

// handleScreenResult routes screen results to appropriate handlers.
func (e *ExportArtifactIntent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
	switch result.Type() {
	case screens.ResultNavigate:
		return e.handleNavigateResult(result)
	case screens.ResultCancel:
		return e.handleCancelResult()
	case screens.ResultSubmit:
		return e.handleSubmitResult(result)
	case screens.ResultError:
		return e.handleErrorResult(result)
	}
	return nil
}

// handleNavigateResult processes navigation results from screens.
func (e *ExportArtifactIntent) handleNavigateResult(_ screens.ScreenResult) tea.Cmd {
	switch e.currentState {
	// Note: ExportStateConfigure (type/format/dest selection) is now handled by the wizard modal
	// in Update(), not by screen navigation results

	case ExportStatePreview:
		// User confirmed preview, show confirmation modal
		e.currentState = ExportStateConfirm
		e.showConfirmModal()

		// Note: ExportStateConfirm, ExportStateComplete, ExportStateFailed are now handled by modals
		// in the Update() method, not by screen results
	}

	return nil
}

// handleCancelResult processes cancel results (Esc key).
func (e *ExportArtifactIntent) handleCancelResult() tea.Cmd {
	switch e.currentState {
	case ExportStateConfigure:
		// Cancel from wizard - handled by wizard itself
		e.setCancelled()

	case ExportStatePreview:
		// Go back to configuration wizard
		e.currentState = ExportStateConfigure

		// Get terminal dimensions
		width, height := 100, 40
		termInfo := e.GetTerminalInfo()
		if termInfo != nil {
			width = termInfo.Width
			height = termInfo.Height
		}

		// Re-create wizard with previous configuration
		e.wizardModal = components.NewExportConfigWizardModal(width, height)
		if e.config != nil {
			e.wizardModal.SetArtifactType(string(e.config.ArtifactType))
			e.wizardModal.SetFormat(string(e.config.Format))
			e.wizardModal.SetDestination(string(e.config.Destination))
		}
		e.activeScreen = nil

		// Note: ExportStateConfirm, ExportStateComplete, ExportStateFailed are now handled by modals
	}

	return nil
}

// handleSubmitResult processes submit results.
func (e *ExportArtifactIntent) handleSubmitResult(_ screens.ScreenResult) tea.Cmd {
	// Not used in export workflow (screens use NavigateResult)
	return nil
}

// handleErrorResult processes error results.
func (e *ExportArtifactIntent) handleErrorResult(result screens.ScreenResult) tea.Cmd {
	if data, ok := result.Data().(map[string]interface{}); ok {
		if msg, ok := data["message"].(string); ok {
			e.exportError = &IntentError{
				Code:    "screen_error",
				Message: msg,
			}
		}
	}
	e.currentState = ExportStateFailed
	e.showErrorModal()
	return nil
}

// transitionToScreen sets up and transitions to a new screen.
func (e *ExportArtifactIntent) transitionToScreen(screen screens.Screen) {
	// Set screen properties from intent
	termInfo := e.GetTerminalInfo()
	if termInfo != nil {
		screen.SetTerminalInfo(termInfo.Width, termInfo.Height)
	}
	screen.SetTheme(e.Theme())
	screen.SetLogo(e.GetLogo(), e.GetLogoSpacing())

	e.activeScreen = screen
}

// View renders the current state using StandardView with modal overlays.
func (e *ExportArtifactIntent) View() string {
	// If wizard is visible, render it as the main content
	if e.wizardModal != nil && e.wizardModal.IsVisible() {
		// Create standard view with breadcrumbs showing we're in configuration
		view := e.CreateViewWithBreadcrumbs("Main Menu", "Export Artifact", "Configure")
		view.WithContent(e.wizardModal.View())
		view.WithHelp(e.getContextHelp()).WithFooterSeparator(true)
		return view.Render()
	}

	if e.activeScreen == nil {
		return "No active screen"
	}

	// Create standard view with breadcrumbs
	view := e.CreateViewWithBreadcrumbs("Main Menu", "Export Artifact", e.getStateName())

	// Get content from screen (RenderContent returns just the content without StandardView wrapper)
	view.WithContent(e.activeScreen.RenderContent())

	// Get context-aware help (themed badges)
	help := e.getContextHelp()
	view.WithHelp(help).WithFooterSeparator(true)

	baseView := view.Render()

	// Render modal overlays on top of the base view using bubbletea-overlay
	// This composites the modal centered on the background content
	if e.confirmModal != nil && e.confirmModal.IsVisible() {
		return behaviors.RenderModalOverlay(e.confirmModal, baseView)
	}
	if e.progressModal != nil && e.progressModal.IsVisible() {
		return behaviors.RenderModalOverlay(e.progressModal, baseView)
	}
	if e.successModal != nil && e.successModal.IsVisible() {
		return behaviors.RenderModalOverlay(e.successModal, baseView)
	}
	if e.errorModal != nil && e.errorModal.IsVisible() {
		return behaviors.RenderModalOverlay(e.errorModal, baseView)
	}

	return baseView
}

// getStateName returns a human-readable name for the current state.
func (e *ExportArtifactIntent) getStateName() string {
	switch e.currentState {
	case ExportStateSelectType:
		return "Select Type"
	case ExportStateSelectFormat:
		return "Select Format"
	case ExportStateSelectDest:
		return "Select Destination"
	case ExportStateConfigure:
		return "Configure"
	case ExportStatePreview:
		return "Preview"
	case ExportStateConfirm:
		return "Confirm"
	case ExportStateInProgress:
		return "Exporting"
	case ExportStateComplete:
		return "Complete"
	case ExportStateFailed:
		return "Failed"
	default:
		return string(e.currentState)
	}
}

// getContextHelp returns context-aware help text for the current state.
func (e *ExportArtifactIntent) getContextHelp() string {
	theme := e.Theme()

	switch e.currentState {
	case ExportStateConfigure:
		// Wizard handles its own footer, so we just provide global badges
		return ThemedGlobalBadges(theme)
	case ExportStateSelectType, ExportStateSelectFormat, ExportStateSelectDest:
		// Legacy states - kept for backward compatibility but no longer used
		return CombineThemedFooters(
			ThemedNavigationFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case ExportStatePreview:
		return CombineThemedFooters(
			ThemedDetailViewFooter(theme),
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Enter", "Continue", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case ExportStateConfirm:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("y/Enter", "Confirm", theme),
				primitives.HelpKeyBadge("n/Esc", "Cancel", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case ExportStateInProgress:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("...", "Please wait", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case ExportStateComplete:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Enter", "Done", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case ExportStateFailed:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("r", "Retry", theme),
				primitives.CancelBadge(theme),
			),
			ThemedGlobalBadges(theme),
		)
	default:
		return ThemedGlobalBadges(theme)
	}
}

// Result returns the intent result.
func (e *ExportArtifactIntent) Result() *IntentResult[interface{}] {
	if e.intentResult == nil {
		return nil
	}

	return &IntentResult[interface{}]{
		Status:   e.intentResult.Status,
		Data:     e.intentResult.Data,
		Error:    e.intentResult.Error,
		Metadata: e.intentResult.Metadata,
	}
}

// setCompleted marks the intent as completed.
func (e *ExportArtifactIntent) setCompleted() {
	e.intentResult = &IntentResult[*ExportArtifactResult]{
		Status: Completed,
		Data:   e.exportResult,
	}
	e.active = false
}

// setCancelled marks the intent as cancelled.
func (e *ExportArtifactIntent) setCancelled() {
	e.intentResult = &IntentResult[*ExportArtifactResult]{
		Status: Cancelled,
	}
	e.active = false
}

// --- Modal Helper Methods ---

// showConfirmModal creates and displays the export confirmation modal.
func (e *ExportArtifactIntent) showConfirmModal() {
	artifactName := getArtifactTypeName(e.config.ArtifactType)
	formatName := getFormatName(e.config.Format)
	destName := getDestinationName(e.config.Destination)

	e.confirmModal = components.NewExportConfirmModal(artifactName, formatName, destName)
	e.confirmModal.SetTheme(e.Theme())

	termInfo := e.GetTerminalInfo()
	if termInfo != nil {
		e.confirmModal.SetDimensions(termInfo.Width, termInfo.Height)
	}
}

// showProgressModal creates and displays the export progress modal.
func (e *ExportArtifactIntent) showProgressModal() {
	artifactName := getArtifactTypeName(e.config.ArtifactType)
	formatName := getFormatName(e.config.Format)
	destName := getDestinationName(e.config.Destination)

	width, height := 100, 40
	termInfo := e.GetTerminalInfo()
	if termInfo != nil {
		width = termInfo.Width
		height = termInfo.Height
	}

	e.progressModal = components.NewExportProgressModal(artifactName, formatName, destName, width, height)
	e.progressModal.SetTheme(e.Theme())
	e.currentState = ExportStateInProgress
}

// showSuccessModal creates and displays the export success modal.
func (e *ExportArtifactIntent) showSuccessModal() {
	artifactName := getArtifactTypeName(e.exportResult.ArtifactType)
	formatName := getFormatName(e.exportResult.Format)
	destName := getDestinationName(e.exportResult.Destination)

	e.successModal = components.NewExportSuccessModal(
		artifactName,
		formatName,
		destName,
		e.exportResult.FilePath,
		e.exportResult.Size,
	)
	e.successModal.SetTheme(e.Theme())

	termInfo := e.GetTerminalInfo()
	if termInfo != nil {
		e.successModal.SetDimensions(termInfo.Width, termInfo.Height)
	}
}

// showErrorModal creates and displays the export error modal.
func (e *ExportArtifactIntent) showErrorModal() {
	e.errorModal = components.NewExportErrorModal(
		e.exportError.Code,
		e.exportError.Message,
	)
	e.errorModal.SetTheme(e.Theme())

	termInfo := e.GetTerminalInfo()
	if termInfo != nil {
		e.errorModal.SetDimensions(termInfo.Width, termInfo.Height)
	}
}

// Helper functions for display names
func getArtifactTypeName(t ExportArtifactType) string {
	switch t {
	case ExportTypeEvents:
		return "Career Events"
	case ExportTypeFacts:
		return "Facts"
	case ExportTypeBursts:
		return "Bursts"
	default:
		return string(t)
	}
}

func getFormatName(f ExportFormat) string {
	switch f {
	case ExportFormatJSON:
		return "JSON"
	case ExportFormatCSV:
		return "CSV"
	case ExportFormatYAML:
		return "YAML"
	case ExportFormatTXT:
		return "Text"
	default:
		return string(f)
	}
}

func getDestinationName(d ExportDestination) string {
	switch d {
	case ExportDestinationFile:
		return "File"
	case ExportDestinationClipboard:
		return "Clipboard"
	default:
		return string(d)
	}
}

// --- Screen Factory Methods ---

// newPreviewScreen creates the preview screen.
func (e *ExportArtifactIntent) newPreviewScreen() screens.Screen {
	return exportscreens.NewPreviewWithStats(
		e.preview,
		types.ExportArtifactType(e.config.ArtifactType),
		types.ExportFormat(e.config.Format),
		types.ExportDestination(e.config.Destination),
		[]string{"Main Menu", "Export Artifact", "Preview"},
		e.previewStats,
	)
}

// Note: Type, Format, and Destination selection are now handled by the wizard modal.
// Confirm, Progress, Complete, and Failed states use modal overlays.

// --- Export Business Logic ---

// startExport returns a command that performs the actual export.
func (e *ExportArtifactIntent) startExport() tea.Cmd {
	return func() tea.Msg {
		var result *ExportArtifactResult
		var err error

		switch e.config.ArtifactType {
		case ExportTypeEvents:
			result, err = e.exportEvents()
		case ExportTypeFacts:
			result, err = e.exportFacts()
		case ExportTypeBursts:
			result, err = e.exportBursts()
		default:
			err = fmt.Errorf("unknown artifact type: %s", e.config.ArtifactType)
		}

		if err != nil {
			return ExportErrorMsg{Error: &IntentError{
				Code:    "export_failed",
				Message: fmt.Sprintf("Export failed: %v", err),
			}}
		}
		return ExportCompleteMsg{Result: result}
	}
}

// exportEvents exports career events.
func (e *ExportArtifactIntent) exportEvents() (*ExportArtifactResult, error) {
	ctx := e.getContext()

	events, err := e.context.EventRepository.List(ctx, careerrepo.ListFilters{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events: %w", err)
	}

	var content string
	switch e.config.Format {
	case ExportFormatJSON:
		content, err = marshalToJSON(events)
	case ExportFormatCSV:
		content, err = marshalEventsToCSV(events)
	case ExportFormatYAML:
		content, err = marshalToYAML(events)
	case ExportFormatTXT:
		content, err = marshalEventsToText(events)
	default:
		return nil, fmt.Errorf("unsupported format: %s", e.config.Format)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to marshal events: %w", err)
	}

	filePath, err := e.saveToDestination(ctx, "career_events", content)
	if err != nil {
		return nil, err
	}

	return NewExportArtifactResult(true, e.config.ArtifactType, e.config.Format, e.config.Destination, filePath, int64(len(content))), nil
}

// exportFacts exports facts.
func (e *ExportArtifactIntent) exportFacts() (*ExportArtifactResult, error) {
	ctx := e.getContext()

	facts, err := e.context.FactRepository.List(ctx, careerrepo.FactListFilters{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch facts: %w", err)
	}

	var content string
	switch e.config.Format {
	case ExportFormatJSON:
		content, err = marshalToJSON(facts)
	case ExportFormatCSV:
		content, err = marshalFactsToCSV(facts)
	case ExportFormatYAML:
		content, err = marshalToYAML(facts)
	case ExportFormatTXT:
		content, err = marshalFactsToText(facts)
	default:
		return nil, fmt.Errorf("unsupported format: %s", e.config.Format)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to marshal facts: %w", err)
	}

	filePath, err := e.saveToDestination(ctx, "facts", content)
	if err != nil {
		return nil, err
	}

	return NewExportArtifactResult(true, e.config.ArtifactType, e.config.Format, e.config.Destination, filePath, int64(len(content))), nil
}

// exportBursts exports bursts.
func (e *ExportArtifactIntent) exportBursts() (*ExportArtifactResult, error) {
	ctx := e.getContext()

	bursts, err := e.context.BurstRepository.List(ctx, careerrepo.BurstListFilters{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bursts: %w", err)
	}

	var content string
	switch e.config.Format {
	case ExportFormatJSON:
		content, err = marshalToJSON(bursts)
	case ExportFormatCSV:
		content, err = marshalBurstsToCSV(bursts)
	case ExportFormatYAML:
		content, err = marshalToYAML(bursts)
	case ExportFormatTXT:
		content, err = marshalBurstsToText(bursts)
	default:
		return nil, fmt.Errorf("unsupported format: %s", e.config.Format)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bursts: %w", err)
	}

	filePath, err := e.saveToDestination(ctx, "bursts", content)
	if err != nil {
		return nil, err
	}

	return NewExportArtifactResult(true, e.config.ArtifactType, e.config.Format, e.config.Destination, filePath, int64(len(content))), nil
}

// saveToDestination saves content to file or clipboard.
func (e *ExportArtifactIntent) saveToDestination(ctx context.Context, name string, content string) (string, error) {
	switch e.config.Destination {
	case ExportDestinationFile:
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}

		exportDir := filepath.Join(homeDir, ".kariya", "exports")
		if err := os.MkdirAll(exportDir, 0750); err != nil {
			return "", fmt.Errorf("failed to create export directory: %w", err)
		}

		timestamp := time.Now().Format("20060102_150405")
		extension := getFileExtensionForFormat(e.config.Format)
		filename := fmt.Sprintf("%s_%s%s", name, timestamp, extension)
		filePath := filepath.Join(exportDir, filename)

		if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
			return "", fmt.Errorf("failed to write file: %w", err)
		}
		return filePath, nil

	case ExportDestinationClipboard:
		if e.context.ExportService == nil {
			return "", fmt.Errorf("export service not available")
		}
		if err := e.context.ExportService.CopyToClipboard(ctx, content); err != nil {
			return "", err
		}
		return "clipboard", nil

	default:
		return "", fmt.Errorf("unsupported destination: %s", e.config.Destination)
	}
}

// generatePreview generates a preview of the export content.
func (e *ExportArtifactIntent) generatePreview() string {
	ctx := e.getContext()

	switch e.config.ArtifactType {
	case ExportTypeEvents:
		return e.generateEventsPreview(ctx)
	case ExportTypeFacts:
		return e.generateFactsPreview(ctx)
	case ExportTypeBursts:
		return e.generateBurstsPreview(ctx)
	default:
		return "No preview available for this artifact type"
	}
}

func (e *ExportArtifactIntent) generateEventsPreview(ctx context.Context) string {
	if e.context.EventRepository == nil {
		return "Preview not available - repository not initialized"
	}

	// Get total count for stats
	allEvents, _ := e.context.EventRepository.List(ctx, careerrepo.ListFilters{})
	totalCount := len(allEvents)

	// Get preview items (limited)
	events, err := e.context.EventRepository.List(ctx, careerrepo.ListFilters{Limit: 10})
	if err != nil {
		return fmt.Sprintf("Error loading events: %v", err)
	}
	if len(events) == 0 {
		return "No events found."
	}

	var content string
	switch e.config.Format {
	case ExportFormatJSON:
		content, err = marshalToJSON(events)
	case ExportFormatCSV:
		content, err = marshalEventsToCSV(events)
	case ExportFormatYAML:
		content, err = marshalToYAML(events)
	case ExportFormatTXT:
		content, err = marshalEventsToText(events)
	default:
		return "Preview not available for format: " + string(e.config.Format)
	}
	if err != nil {
		return fmt.Sprintf("Error generating preview: %v", err)
	}

	// Calculate stats
	isTruncated := len(content) > 2000
	if isTruncated {
		content = content[:2000] + "\n\n...(preview truncated)..."
	}

	e.previewStats = &exportscreens.PreviewStats{
		ItemCount:     len(events),
		TotalCount:    totalCount,
		EstimatedSize: int64(len(content) * totalCount / max(len(events), 1)),
		ContentLines:  countLines(content),
		IsTruncated:   isTruncated,
	}

	return content
}

func (e *ExportArtifactIntent) generateFactsPreview(ctx context.Context) string {
	if e.context.FactRepository == nil {
		return "Preview not available - repository not initialized"
	}

	// Get total count for stats
	allFacts, _ := e.context.FactRepository.List(ctx, careerrepo.FactListFilters{})
	totalCount := len(allFacts)

	facts, err := e.context.FactRepository.List(ctx, careerrepo.FactListFilters{Limit: 10})
	if err != nil {
		return fmt.Sprintf("Error loading facts: %v", err)
	}
	if len(facts) == 0 {
		return "No facts found."
	}

	var content string
	switch e.config.Format {
	case ExportFormatJSON:
		content, err = marshalToJSON(facts)
	case ExportFormatCSV:
		content, err = marshalFactsToCSV(facts)
	case ExportFormatYAML:
		content, err = marshalToYAML(facts)
	case ExportFormatTXT:
		content, err = marshalFactsToText(facts)
	default:
		return "Preview not available for format: " + string(e.config.Format)
	}
	if err != nil {
		return fmt.Sprintf("Error generating preview: %v", err)
	}

	// Calculate stats
	isTruncated := len(content) > 2000
	if isTruncated {
		content = content[:2000] + "\n\n...(preview truncated)..."
	}

	e.previewStats = &exportscreens.PreviewStats{
		ItemCount:     len(facts),
		TotalCount:    totalCount,
		EstimatedSize: int64(len(content) * totalCount / max(len(facts), 1)),
		ContentLines:  countLines(content),
		IsTruncated:   isTruncated,
	}

	return content
}

func (e *ExportArtifactIntent) generateBurstsPreview(ctx context.Context) string {
	if e.context.BurstRepository == nil {
		return "Preview not available - repository not initialized"
	}

	// Get total count for stats
	allBursts, _ := e.context.BurstRepository.List(ctx, careerrepo.BurstListFilters{})
	totalCount := len(allBursts)

	bursts, err := e.context.BurstRepository.List(ctx, careerrepo.BurstListFilters{Limit: 10})
	if err != nil {
		return fmt.Sprintf("Error loading bursts: %v", err)
	}
	if len(bursts) == 0 {
		return "No bursts found."
	}

	var content string
	switch e.config.Format {
	case ExportFormatJSON:
		content, err = marshalToJSON(bursts)
	case ExportFormatCSV:
		content, err = marshalBurstsToCSV(bursts)
	case ExportFormatYAML:
		content, err = marshalToYAML(bursts)
	case ExportFormatTXT:
		content, err = marshalBurstsToText(bursts)
	default:
		return "Preview not available for format: " + string(e.config.Format)
	}
	if err != nil {
		return fmt.Sprintf("Error generating preview: %v", err)
	}

	// Calculate stats
	isTruncated := len(content) > 2000
	if isTruncated {
		content = content[:2000] + "\n\n...(preview truncated)..."
	}

	e.previewStats = &exportscreens.PreviewStats{
		ItemCount:     len(bursts),
		TotalCount:    totalCount,
		EstimatedSize: int64(len(content) * totalCount / max(len(bursts), 1)),
		ContentLines:  countLines(content),
		IsTruncated:   isTruncated,
	}

	return content
}

// countLines counts the number of lines in a string.
func countLines(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n") + 1
}

// getContext returns a context for service calls.
func (e *ExportArtifactIntent) getContext() context.Context {
	if e.context.AppContext != nil {
		return e.context.AppContext
	}
	return context.Background()
}

// --- Test Helpers ---

// GetState returns the current state (for testing).
func (e *ExportArtifactIntent) GetState() ExportState {
	return e.currentState
}

// GetConfig returns the current export configuration (for testing).
func (e *ExportArtifactIntent) GetConfig() *ExportConfiguration {
	return e.config
}

// GetResult returns the export result (for testing).
func (e *ExportArtifactIntent) GetResult() *ExportArtifactResult {
	return e.exportResult
}

// SetState sets the state (for testing).
func (e *ExportArtifactIntent) SetState(state ExportState) {
	e.currentState = state
}

// SetConfig sets the configuration (for testing).
func (e *ExportArtifactIntent) SetConfig(config *ExportConfiguration) {
	e.config = config
}

// IsActive returns whether the intent is active.
func (e *ExportArtifactIntent) IsActive() bool {
	return e.active
}
