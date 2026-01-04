package intents

import (
	tea "github.com/charmbracelet/bubbletea"
)

// ExportArtifactType represents the type of artifact to export
type ExportArtifactType string

const (
	ExportTypeCV      ExportArtifactType = "cv"
	ExportTypeEvents  ExportArtifactType = "events"
	ExportTypeFacts   ExportArtifactType = "facts"
	ExportTypeBursts  ExportArtifactType = "bursts"
	ExportTypeProfile ExportArtifactType = "profile"
)

// ExportFormat represents the export file format
type ExportFormat string

const (
	ExportFormatPDF  ExportFormat = "pdf"
	ExportFormatJSON ExportFormat = "json"
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatTXT  ExportFormat = "txt"
	ExportFormatMD   ExportFormat = "markdown"
)

// ExportDestination represents where to export the artifact
type ExportDestination string

const (
	ExportDestinationFile      ExportDestination = "file"
	ExportDestinationClipboard ExportDestination = "clipboard"
	ExportDestinationEmail     ExportDestination = "email"
)

// ExportArtifactContext contains context for the ExportArtifact intent
type ExportArtifactContext struct {
	ArtifactTypes    []ExportArtifactType
	SupportedFormats map[ExportArtifactType][]ExportFormat
	DefaultFormat    map[ExportArtifactType]ExportFormat
	Destinations     []ExportDestination
}

// ExportConfiguration holds the current export configuration
type ExportConfiguration struct {
	ArtifactType ExportArtifactType
	Format       ExportFormat
	Destination  ExportDestination
	FilePath     string // For file destination
	Email        string // For email destination
}

// ExportArtifactResult contains the result of artifact export
type ExportArtifactResult struct {
	Success      bool
	ArtifactType ExportArtifactType
	Format       ExportFormat
	Destination  ExportDestination
	FilePath     string // Path where artifact was saved
	Size         int64  // Size of exported artifact in bytes
	Error        *IntentError
}

// ExportState represents the current state of the export intent
type ExportState string

const (
	ExportStateSelectType   ExportState = "select_type"
	ExportStateSelectFormat ExportState = "select_format"
	ExportStateSelectDest   ExportState = "select_destination"
	ExportStateConfigure    ExportState = "configure"
	ExportStatePreview      ExportState = "preview"
	ExportStateConfirm      ExportState = "confirm"
	ExportStateInProgress   ExportState = "in_progress"
	ExportStateComplete     ExportState = "complete"
	ExportStateFailed       ExportState = "failed"
)

// ExportProgressMsg represents progress of an export operation
type ExportProgressMsg struct {
	Percentage int
	Message    string
}

// ExportCompleteMsg represents completion of an export operation
type ExportCompleteMsg struct {
	Result *ExportArtifactResult
	Error  *IntentError
}

// ExportCancelledMsg represents cancellation of an export operation
type ExportCancelledMsg struct{}

// ExportErrorMsg represents an error during export
type ExportErrorMsg struct {
	Error *IntentError
}

// NewExportArtifactContext creates a new ExportArtifactContext with default values
func NewExportArtifactContext() *ExportArtifactContext {
	return &ExportArtifactContext{
		ArtifactTypes: []ExportArtifactType{
			ExportTypeCV,
			ExportTypeEvents,
			ExportTypeFacts,
			ExportTypeBursts,
			ExportTypeProfile,
		},
		SupportedFormats: map[ExportArtifactType][]ExportFormat{
			ExportTypeCV:      {ExportFormatPDF, ExportFormatJSON, ExportFormatMD},
			ExportTypeEvents:  {ExportFormatJSON, ExportFormatCSV, ExportFormatTXT},
			ExportTypeFacts:   {ExportFormatJSON, ExportFormatCSV, ExportFormatTXT},
			ExportTypeBursts:  {ExportFormatJSON, ExportFormatCSV, ExportFormatTXT},
			ExportTypeProfile: {ExportFormatJSON, ExportFormatPDF},
		},
		DefaultFormat: map[ExportArtifactType]ExportFormat{
			ExportTypeCV:      ExportFormatPDF,
			ExportTypeEvents:  ExportFormatJSON,
			ExportTypeFacts:   ExportFormatJSON,
			ExportTypeBursts:  ExportFormatJSON,
			ExportTypeProfile: ExportFormatJSON,
		},
		Destinations: []ExportDestination{
			ExportDestinationFile,
			ExportDestinationClipboard,
			ExportDestinationEmail,
		},
	}
}

// NewExportConfiguration creates a new export configuration with defaults
func NewExportConfiguration(artifactType ExportArtifactType, ctx *ExportArtifactContext) *ExportConfiguration {
	format := ctx.DefaultFormat[artifactType]
	return &ExportConfiguration{
		ArtifactType: artifactType,
		Format:       format,
		Destination:  ExportDestinationFile,
	}
}

// NewExportArtifactResult creates a new export result
func NewExportArtifactResult(success bool, artifactType ExportArtifactType, format ExportFormat, destination ExportDestination, filePath string, size int64) *ExportArtifactResult {
	return &ExportArtifactResult{
		Success:      success,
		ArtifactType: artifactType,
		Format:       format,
		Destination:  destination,
		FilePath:     filePath,
		Size:         size,
	}
}

// NewExportArtifactResultWithError creates a new export result with an error
func NewExportArtifactResultWithError(err *IntentError) *ExportArtifactResult {
	return &ExportArtifactResult{
		Success: false,
		Error:   err,
	}
}

// ExportArtifactModel represents the state of the ExportArtifact intent
type ExportArtifactModel struct {
	state         ExportState
	context       *ExportArtifactContext
	config        *ExportConfiguration
	preview       string // Preview of artifact to export
	selectedIndex int    // Currently selected index in list
	result        *ExportArtifactResult
	error         *IntentError
	active        bool
}

// NewExportArtifactModel creates a new ExportArtifact intent model
func NewExportArtifactModel(ctx *ExportArtifactContext) *ExportArtifactModel {
	return &ExportArtifactModel{
		state:         ExportStateSelectType,
		context:       ctx,
		selectedIndex: 0,
		active:        false,
	}
}

// Init initializes the intent
func (m *ExportArtifactModel) Init() tea.Cmd {
	m.active = true
	return nil
}

// Update handles messages
func (m *ExportArtifactModel) Update(msg tea.Msg) tea.Cmd {
	if !m.active {
		return nil
	}

	switch m.state {
	case ExportStateSelectType:
		return m.updateSelectType(msg)
	case ExportStateSelectFormat:
		return m.updateSelectFormat(msg)
	case ExportStateSelectDest:
		return m.updateSelectDest(msg)
	case ExportStateConfigure:
		return m.updateConfigure(msg)
	case ExportStatePreview:
		return m.updatePreview(msg)
	case ExportStateConfirm:
		return m.updateConfirm(msg)
	case ExportStateInProgress:
		return m.updateInProgress(msg)
	case ExportStateComplete:
		return m.updateComplete(msg)
	case ExportStateFailed:
		return m.updateFailed(msg)
	default:
		return nil
	}
}

// View renders the current state
func (m *ExportArtifactModel) View() string {
	if !m.active {
		return ""
	}

	switch m.state {
	case ExportStateSelectType:
		return m.viewSelectType()
	case ExportStateSelectFormat:
		return m.viewSelectFormat()
	case ExportStateSelectDest:
		return m.viewSelectDest()
	case ExportStateConfigure:
		return m.viewConfigure()
	case ExportStatePreview:
		return m.viewPreview()
	case ExportStateConfirm:
		return m.viewConfirm()
	case ExportStateInProgress:
		return m.viewInProgress()
	case ExportStateComplete:
		return m.viewComplete()
	case ExportStateFailed:
		return m.viewFailed()
	default:
		return "Unknown state"
	}
}

// Result returns the intent result
func (m *ExportArtifactModel) Result() *IntentResult[interface{}] {
	if m.result == nil {
		return nil
	}

	if m.result.Error != nil && m.result.Error.Code == "export_cancelled" {
		return &IntentResult[interface{}]{
			Status: Cancelled,
		}
	}

	if !m.result.Success {
		return &IntentResult[interface{}]{
			Status: Failed,
			Data:   m.result,
			Error:  m.result.Error,
		}
	}

	return &IntentResult[interface{}]{
		Status: Completed,
		Data:   m.result,
	}
}

// State-specific update handlers

func (m *ExportArtifactModel) updateSelectType(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "down":
			if m.selectedIndex < len(m.context.ArtifactTypes)-1 {
				m.selectedIndex++
			}
		case "enter":
			artifactType := m.context.ArtifactTypes[m.selectedIndex]
			m.config = NewExportConfiguration(artifactType, m.context)
			m.state = ExportStateSelectFormat
		case "esc":
			m.setResult(NewExportArtifactResultWithError(&IntentError{
				Code:    "export_cancelled",
				Message: "Export cancelled by user",
			}))
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateSelectFormat(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		formats := m.context.SupportedFormats[m.config.ArtifactType]
		switch msg.String() {
		case "up":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "down":
			if m.selectedIndex < len(formats)-1 {
				m.selectedIndex++
			}
		case "enter":
			formats := m.context.SupportedFormats[m.config.ArtifactType]
			m.config.Format = formats[m.selectedIndex]
			m.selectedIndex = 0
			m.state = ExportStateSelectDest
		case "esc":
			m.selectedIndex = 0
			m.state = ExportStateSelectType
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateSelectDest(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "down":
			if m.selectedIndex < len(m.context.Destinations)-1 {
				m.selectedIndex++
			}
		case "enter":
			m.config.Destination = m.context.Destinations[m.selectedIndex]
			m.selectedIndex = 0
			m.state = ExportStateConfigure
		case "esc":
			m.selectedIndex = 0
			m.state = ExportStateSelectFormat
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateConfigure(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.state = ExportStatePreview
		case "esc":
			m.selectedIndex = 0
			m.state = ExportStateSelectDest
		}
	}
	return nil
}

func (m *ExportArtifactModel) updatePreview(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.state = ExportStateConfirm
		case "esc":
			m.state = ExportStateConfigure
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateConfirm(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "enter":
			m.state = ExportStateInProgress
			return m.startExport()
		case "n", "esc":
			m.state = ExportStatePreview
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateInProgress(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case ExportProgressMsg:
		// Update progress (could be displayed in view)
		return nil
	case ExportCompleteMsg:
		m.result = msg.Result
		m.state = ExportStateComplete
	case ExportErrorMsg:
		m.error = msg.Error
		m.state = ExportStateFailed
	case tea.KeyMsg:
		if msg.String() == "esc" {
			// Allow cancellation during export
			return nil
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateComplete(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" || msg.String() == "esc" {
			m.active = false
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateFailed(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			// Retry
			m.state = ExportStateConfirm
		case "esc":
			m.active = false
		}
	}
	return nil
}

// startExport initiates the async export operation
func (m *ExportArtifactModel) startExport() tea.Cmd {
	return func() tea.Msg {
		// Simulate export operation (would be replaced with actual service call)
		result := NewExportArtifactResult(
			true,
			m.config.ArtifactType,
			m.config.Format,
			m.config.Destination,
			"/tmp/export."+string(m.config.Format),
			1024,
		)
		return ExportCompleteMsg{Result: result}
	}
}

// setResult sets the intent result and marks intent as complete
func (m *ExportArtifactModel) setResult(result *ExportArtifactResult) {
	m.result = result
	m.active = false
}

// View rendering methods

func (m *ExportArtifactModel) viewSelectType() string {
	s := "Select Artifact Type:\n\n"
	for i, t := range m.context.ArtifactTypes {
		prefix := "  "
		if i == m.selectedIndex {
			prefix = "> "
		}
		s += prefix + string(t) + "\n"
	}
	s += "\n[↑/↓] Navigate | [Enter] Select | [Esc] Cancel"
	return s
}

func (m *ExportArtifactModel) viewSelectFormat() string {
	formats := m.context.SupportedFormats[m.config.ArtifactType]
	s := "Select Export Format for " + string(m.config.ArtifactType) + ":\n\n"
	for i, f := range formats {
		prefix := "  "
		if i == m.selectedIndex {
			prefix = "> "
		}
		s += prefix + string(f) + "\n"
	}
	s += "\n[↑/↓] Navigate | [Enter] Select | [Esc] Back"
	return s
}

func (m *ExportArtifactModel) viewSelectDest() string {
	s := "Select Export Destination:\n\n"
	for i, d := range m.context.Destinations {
		prefix := "  "
		if i == m.selectedIndex {
			prefix = "> "
		}
		s += prefix + string(d) + "\n"
	}
	s += "\n[↑/↓] Navigate | [Enter] Select | [Esc] Back"
	return s
}

func (m *ExportArtifactModel) viewConfigure() string {
	s := "Configure Export:\n\n"
	s += "Artifact Type: " + string(m.config.ArtifactType) + "\n"
	s += "Format: " + string(m.config.Format) + "\n"
	s += "Destination: " + string(m.config.Destination) + "\n"
	if m.config.Destination == ExportDestinationFile {
		s += "File Path: " + m.config.FilePath + "\n"
	}
	s += "\n[Enter] Continue | [Esc] Back"
	return s
}

func (m *ExportArtifactModel) viewPreview() string {
	s := "Preview Export:\n\n"
	s += m.preview + "\n\n"
	s += "[Enter] Proceed | [Esc] Back"
	return s
}

func (m *ExportArtifactModel) viewConfirm() string {
	s := "Confirm Export?\n\n"
	s += "Artifact: " + string(m.config.ArtifactType) + "\n"
	s += "Format: " + string(m.config.Format) + "\n"
	s += "Destination: " + string(m.config.Destination) + "\n\n"
	s += "[Y/Enter] Confirm | [N/Esc] Cancel"
	return s
}

func (m *ExportArtifactModel) viewInProgress() string {
	return "Exporting artifact...\n\nPlease wait..."
}

func (m *ExportArtifactModel) viewComplete() string {
	s := "Export Complete!\n\n"
	s += "File: " + m.result.FilePath + "\n"
	s += "Size: " + formatBytes(m.result.Size) + "\n\n"
	s += "[Enter] Done"
	return s
}

func (m *ExportArtifactModel) viewFailed() string {
	s := "Export Failed\n\n"
	if m.error != nil {
		s += "Error: " + m.error.Message + "\n"
	}
	s += "\n[R] Retry | [Esc] Cancel"
	return s
}

// Helper functions

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return string(rune(bytes)) + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"", "K", "M", "G", "T"}
	return string(rune(bytes/div)) + " " + units[exp] + "B"
}
