package intents

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	tea "github.com/charmbracelet/bubbletea"
)

// ExportArtifactIntent implements the Intent interface for artifact export.
// It orchestrates screens and handles the export workflow.
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

	// result holds the final export result
	exportResult *ExportArtifactResult

	// intentResult is the final result of the intent
	intentResult *IntentResult[*ExportArtifactResult]

	// error holds any error that occurred
	exportError *IntentError

	// activeScreen holds the current screen being displayed
	activeScreen screens.Screen

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

// Init initializes the intent and shows the first screen.
func (e *ExportArtifactIntent) Init() tea.Cmd {
	e.active = true
	e.transitionToScreen(e.newTypeSelectScreen())
	return nil
}

// Update handles messages and delegates to the active screen.
func (e *ExportArtifactIntent) Update(msg tea.Msg) tea.Cmd {
	if !e.active {
		return nil
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

	// Handle progress screen messages
	if e.currentState == ExportStateInProgress {
		switch msg := msg.(type) {
		case base.TickMsg:
			// Forward tick to progress screen for spinner animation
			if e.activeScreen != nil {
				cmd, _ := e.activeScreen.Update(msg)
				return tea.Batch(cmd, base.TickCmd())
			}
		case ExportCompleteMsg:
			e.exportResult = msg.Result
			e.currentState = ExportStateComplete
			e.transitionToScreen(e.newCompleteScreen())
			return nil
		case ExportErrorMsg:
			e.exportError = msg.Error
			e.currentState = ExportStateFailed
			e.transitionToScreen(e.newFailedScreen())
			return nil
		}
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
func (e *ExportArtifactIntent) handleNavigateResult(result screens.ScreenResult) tea.Cmd {
	switch e.currentState {
	case ExportStateSelectType:
		// User selected artifact type
		if artifactType, ok := result.Data().(ExportArtifactType); ok {
			e.config = NewExportConfiguration(artifactType, e.context)
			e.currentState = ExportStateSelectFormat
			e.transitionToScreen(e.newFormatSelectScreen())
		}

	case ExportStateSelectFormat:
		// User selected format
		if format, ok := result.Data().(ExportFormat); ok {
			e.config.Format = format
			e.currentState = ExportStateSelectDest
			e.transitionToScreen(e.newDestSelectScreen())
		}

	case ExportStateSelectDest:
		// User selected destination
		if dest, ok := result.Data().(ExportDestination); ok {
			e.config.Destination = dest
			// Generate preview and show it
			e.preview = e.generatePreview()
			e.currentState = ExportStatePreview
			e.transitionToScreen(e.newPreviewScreen())
		}

	case ExportStatePreview:
		// User confirmed preview, show confirmation
		e.currentState = ExportStateConfirm
		e.transitionToScreen(e.newConfirmScreen())

	case ExportStateConfirm:
		// User made confirmation choice
		if confirmed, ok := result.Data().(bool); ok {
			if confirmed {
				// Start export
				e.currentState = ExportStateInProgress
				e.transitionToScreen(e.newProgressScreen())
				return tea.Batch(base.TickCmd(), e.startExport())
			}
			// User declined - go back to preview
			e.currentState = ExportStatePreview
			e.transitionToScreen(e.newPreviewScreen())
		}

	case ExportStateComplete:
		// User acknowledged completion
		e.setCompleted()

	case ExportStateFailed:
		// User chose action
		if action, ok := result.Data().(string); ok && action == "retry" {
			// Retry export
			e.currentState = ExportStateInProgress
			e.transitionToScreen(e.newProgressScreen())
			return tea.Batch(base.TickCmd(), e.startExport())
		}
		e.setCancelled()
	}

	return nil
}

// handleCancelResult processes cancel results (Esc key).
func (e *ExportArtifactIntent) handleCancelResult() tea.Cmd {
	switch e.currentState {
	case ExportStateSelectType:
		// Cancel from first screen - exit intent
		e.setCancelled()

	case ExportStateSelectFormat:
		// Go back to type selection
		e.currentState = ExportStateSelectType
		e.transitionToScreen(e.newTypeSelectScreen())

	case ExportStateSelectDest:
		// Go back to format selection
		e.currentState = ExportStateSelectFormat
		e.transitionToScreen(e.newFormatSelectScreen())

	case ExportStatePreview:
		// Go back to destination selection
		e.currentState = ExportStateSelectDest
		e.transitionToScreen(e.newDestSelectScreen())

	case ExportStateConfirm:
		// Go back to preview
		e.currentState = ExportStatePreview
		e.transitionToScreen(e.newPreviewScreen())

	case ExportStateComplete, ExportStateFailed:
		// Exit intent
		e.setCancelled()
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
	e.transitionToScreen(e.newFailedScreen())
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

// View renders the current state using StandardView.
func (e *ExportArtifactIntent) View() string {
	if e.activeScreen == nil {
		return "No active screen"
	}

	// Create standard view with breadcrumbs
	view := e.CreateViewWithBreadcrumbs("Main Menu", "Export Artifact", e.getStateName())

	// Get content from screen
	view.WithContent(e.activeScreen.View())

	// Get context-aware help
	help := e.getContextHelp()
	view.WithHelp(help).WithFooterSeparator(true)

	return view.Render()
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
	case ExportStateSelectType, ExportStateSelectFormat, ExportStateSelectDest:
		return CombineThemedFooters(
			ThemedNavigationFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case ExportStateConfigure:
		return CombineThemedFooters(
			ThemedFormFooter(theme),
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

// --- Screen Factory Methods ---

// newTypeSelectScreen creates the artifact type selection screen.
func (e *ExportArtifactIntent) newTypeSelectScreen() screens.Screen {
	renderer := func(item ExportArtifactType) string {
		descriptions := map[ExportArtifactType]string{
			ExportTypeEvents:  "Export career events",
			ExportTypeFacts:   "Export extracted facts",
			ExportTypeBursts:  "Export career bursts",
			ExportTypeProfile: "Export profile data",
			ExportTypeCV:      "Export generated CV",
		}
		desc := descriptions[item]
		if desc == "" {
			desc = "Export " + string(item)
		}
		return fmt.Sprintf("%s\n  %s", string(item), desc)
	}

	return base.NewBaseSelectScreen(
		e.context.ArtifactTypes,
		renderer,
		[]string{"Main Menu", "Export Artifact", "Select Type"},
		"Select Artifact Type",
	)
}

// newFormatSelectScreen creates the format selection screen.
func (e *ExportArtifactIntent) newFormatSelectScreen() screens.Screen {
	renderer := func(item ExportFormat) string {
		descriptions := map[ExportFormat]string{
			ExportFormatJSON: "JavaScript Object Notation",
			ExportFormatCSV:  "Comma Separated Values",
			ExportFormatYAML: "YAML Ain't Markup Language",
			ExportFormatTXT:  "Plain text format",
			ExportFormatMD:   "Markdown format",
			ExportFormatPDF:  "Portable Document Format",
		}
		desc := descriptions[item]
		if desc == "" {
			desc = string(item) + " format"
		}
		return fmt.Sprintf("%s\n  %s", string(item), desc)
	}

	formats := e.context.SupportedFormats[e.config.ArtifactType]
	return base.NewBaseSelectScreen(
		formats,
		renderer,
		[]string{"Main Menu", "Export Artifact", "Select Format"},
		fmt.Sprintf("Select Export Format for %s", e.config.ArtifactType),
	)
}

// newDestSelectScreen creates the destination selection screen.
func (e *ExportArtifactIntent) newDestSelectScreen() screens.Screen {
	renderer := func(item ExportDestination) string {
		descriptions := map[ExportDestination]string{
			ExportDestinationFile:      "Save to a file on disk",
			ExportDestinationClipboard: "Copy to system clipboard",
			ExportDestinationEmail:     "Send via email",
		}
		desc := descriptions[item]
		if desc == "" {
			desc = string(item)
		}
		return fmt.Sprintf("%s\n  %s", string(item), desc)
	}

	return base.NewBaseSelectScreen(
		e.context.Destinations,
		renderer,
		[]string{"Main Menu", "Export Artifact", "Select Destination"},
		"Select Export Destination",
	)
}

// newPreviewScreen creates the preview screen.
func (e *ExportArtifactIntent) newPreviewScreen() screens.Screen {
	renderer := func(data string, _, _ int) string {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("Preview Export (%s as %s):\n\n", e.config.ArtifactType, e.config.Format))
		b.WriteString("════════════════════════════════════════════════════════\n\n")
		b.WriteString(data)
		b.WriteString("\n\n════════════════════════════════════════════════════════")
		return b.String()
	}

	screen := base.NewBaseDetailScreen(
		[]string{"Main Menu", "Export Artifact", "Preview"},
		renderer,
		e.preview,
	)
	screen.SetFooter("↑/↓/j/k: Scroll  Enter: Continue  Esc: Back")

	return screen
}

// newConfirmScreen creates the confirmation screen.
func (e *ExportArtifactIntent) newConfirmScreen() screens.Screen {
	message := fmt.Sprintf(
		"Confirm export?\n\n"+
			"Artifact: %s\n"+
			"Format: %s\n"+
			"Destination: %s",
		e.config.ArtifactType,
		e.config.Format,
		e.config.Destination,
	)

	screen := base.NewBaseConfirmScreen(
		[]string{"Main Menu", "Export Artifact", "Confirm"},
		"Confirm Export",
		message,
	)
	screen.SetYesText("Export")
	screen.SetNoText("Cancel")

	return screen
}

// newProgressScreen creates the progress screen.
func (e *ExportArtifactIntent) newProgressScreen() screens.Screen {
	screen := base.NewBaseProgressScreen(
		[]string{"Main Menu", "Export Artifact", "Exporting"},
		"Exporting",
		"Exporting "+string(e.config.ArtifactType)+"...",
	)
	screen.SetAllowCancel(false)

	return screen
}

// newCompleteScreen creates the completion screen.
func (e *ExportArtifactIntent) newCompleteScreen() screens.Screen {
	renderer := func(data *ExportArtifactResult, _, _ int) string {
		var b strings.Builder
		b.WriteString("Export Complete!\n\n")
		b.WriteString(fmt.Sprintf("Artifact: %s\n", data.ArtifactType))
		b.WriteString(fmt.Sprintf("Format: %s\n", data.Format))
		b.WriteString(fmt.Sprintf("Destination: %s\n", data.Destination))
		b.WriteString(fmt.Sprintf("File: %s\n", data.FilePath))
		b.WriteString(fmt.Sprintf("Size: %s\n", formatBytes(data.Size)))
		return b.String()
	}

	screen := base.NewBaseDetailScreen(
		[]string{"Main Menu", "Export Artifact", "Complete"},
		renderer,
		e.exportResult,
	)
	screen.SetFooter("Enter: Done  Esc: Back")

	return screen
}

// newFailedScreen creates the failed screen.
func (e *ExportArtifactIntent) newFailedScreen() screens.Screen {
	renderer := func(data *IntentError, _, _ int) string {
		var b strings.Builder
		b.WriteString("Export Failed\n\n")
		b.WriteString(fmt.Sprintf("Error: %s\n", data.Message))
		if data.Code != "" {
			b.WriteString(fmt.Sprintf("Code: %s\n", data.Code))
		}
		return b.String()
	}

	screen := base.NewBaseDetailScreen(
		[]string{"Main Menu", "Export Artifact", "Failed"},
		renderer,
		e.exportError,
	)
	screen.SetFooter("r: Retry  Esc: Back")
	screen.AddAction("r", "retry")

	return screen
}

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

	if len(content) > 2000 {
		content = content[:2000] + "\n\n...(preview truncated)..."
	}
	return content
}

func (e *ExportArtifactIntent) generateFactsPreview(ctx context.Context) string {
	if e.context.FactRepository == nil {
		return "Preview not available - repository not initialized"
	}

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

	if len(content) > 2000 {
		content = content[:2000] + "\n\n...(preview truncated)..."
	}
	return content
}

func (e *ExportArtifactIntent) generateBurstsPreview(ctx context.Context) string {
	if e.context.BurstRepository == nil {
		return "Preview not available - repository not initialized"
	}

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

	if len(content) > 2000 {
		content = content[:2000] + "\n\n...(preview truncated)..."
	}
	return content
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
