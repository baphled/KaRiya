package intents

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/navigation"
	careerdomain "github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
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
	ExportFormatYAML ExportFormat = "yaml"
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
	// Configuration
	ArtifactTypes    []ExportArtifactType
	SupportedFormats map[ExportArtifactType][]ExportFormat
	DefaultFormat    map[ExportArtifactType]ExportFormat
	Destinations     []ExportDestination

	// Services
	ExportService *cv.ExportService
	CareerService *career.Service

	// Repositories
	EventRepository careerrepo.Repository
	FactRepository  careerrepo.FactRepository
	BurstRepository careerrepo.BurstRepository

	// Context
	AppContext context.Context
}

// ExportConfiguration holds the current export configuration
type ExportConfiguration struct {
	ArtifactType ExportArtifactType
	Format       ExportFormat
	Destination  ExportDestination
	FilePath     string
	Email        string
}

// ExportArtifactResult contains the result of artifact export
type ExportArtifactResult struct {
	Success      bool
	ArtifactType ExportArtifactType
	Format       ExportFormat
	Destination  ExportDestination
	FilePath     string
	Size         int64
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
	preview       string
	previewLines  []string
	scrollOffset  int
	selectedIndex int
	result        *ExportArtifactResult
	error         *IntentError
	active        bool
	navHandler    *navigation.ListNavigationHandler
}

// NewExportArtifactModel creates a new ExportArtifact intent model
func NewExportArtifactModel(ctx *ExportArtifactContext) *ExportArtifactModel {
	model := &ExportArtifactModel{
		state:         ExportStateSelectType,
		context:       ctx,
		selectedIndex: 0,
		active:        false,
		scrollOffset:  0,
	}
	// Initialize navigation handler for list navigation
	model.navHandler = navigation.NewListNavigationHandler(model)
	return model
}

// ListNavigator interface implementation for state-aware list navigation

// GetTotalItems returns the total number of items in the current list based on state
func (m *ExportArtifactModel) GetTotalItems() int {
	switch m.state {
	case ExportStateSelectType:
		return len(m.context.ArtifactTypes)
	case ExportStateSelectFormat:
		if m.config != nil {
			return len(m.context.SupportedFormats[m.config.ArtifactType])
		}
		return 0
	case ExportStateSelectDest:
		return len(m.context.Destinations)
	default:
		return 0
	}
}

// GetSelectedIndex returns the current selection index
func (m *ExportArtifactModel) GetSelectedIndex() int {
	return m.selectedIndex
}

// SetSelectedIndex sets the selection index with bounds checking
func (m *ExportArtifactModel) SetSelectedIndex(idx int) {
	if idx < 0 {
		idx = 0
	}
	total := m.GetTotalItems()
	if total > 0 && idx >= total {
		idx = total - 1
	}
	m.selectedIndex = idx
}

// GetPageSize returns items per page for pagination
func (m *ExportArtifactModel) GetPageSize() int {
	return 10
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
		// Handle global keys (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			// Help modal is handled at the intent wrapper level (ExportArtifactIntent)
			return nil
		case KeyBack:
			m.setResult(NewExportArtifactResultWithError(&IntentError{
				Code:    "export_cancelled",
				Message: "Export cancelled by user",
			}))
			return nil
		}

		// Handle navigation via centralized handler
		if m.navHandler.HandleKey(msg.String()) {
			return nil
		}

		switch msg.String() {
		case "enter":
			artifactType := m.context.ArtifactTypes[m.selectedIndex]
			m.config = NewExportConfiguration(artifactType, m.context)
			// Go directly to format selection
			m.selectedIndex = 0
			m.state = ExportStateSelectFormat
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateSelectFormat(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			// Help modal is handled at the intent wrapper level (ExportArtifactIntent)
			return nil
		case KeyBack:
			m.selectedIndex = 0
			m.state = ExportStateSelectType
			return nil
		}

		// Handle navigation via centralized handler
		if m.navHandler.HandleKey(msg.String()) {
			return nil
		}

		switch msg.String() {
		case "enter":
			formats := m.context.SupportedFormats[m.config.ArtifactType]
			m.config.Format = formats[m.selectedIndex]
			m.selectedIndex = 0
			m.state = ExportStateSelectDest
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateSelectDest(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			// Help modal is handled at the intent wrapper level (ExportArtifactIntent)
			return nil
		case KeyBack:
			m.selectedIndex = 0
			m.state = ExportStateSelectFormat
			return nil
		}

		// Handle navigation via centralized handler
		if m.navHandler.HandleKey(msg.String()) {
			return nil
		}

		switch msg.String() {
		case "enter":
			m.config.Destination = m.context.Destinations[m.selectedIndex]
			m.selectedIndex = 0
			m.state = ExportStateConfigure
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateConfigure(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			// Help modal is handled at the intent wrapper level (ExportArtifactIntent)
			return nil
		case KeyBack:
			m.selectedIndex = 0
			m.state = ExportStateSelectDest
			return nil
		}

		switch msg.String() {
		case "enter":
			m.generatePreview()
			m.scrollOffset = 0
			m.state = ExportStatePreview
		}
	}
	return nil
}

func (m *ExportArtifactModel) updatePreview(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			// Help modal is handled at the intent wrapper level (ExportArtifactIntent)
			return nil
		case KeyBack:
			m.state = ExportStateConfigure
			return nil
		}

		switch msg.String() {
		case "up", "k":
			if m.scrollOffset > 0 {
				m.scrollOffset--
			}
		case "down", "j":
			if m.scrollOffset < len(m.previewLines)-1 {
				m.scrollOffset++
			}
		case "pageup":
			m.scrollOffset -= 10
			if m.scrollOffset < 0 {
				m.scrollOffset = 0
			}
		case "pagedown":
			m.scrollOffset += 10
			if m.scrollOffset >= len(m.previewLines) {
				m.scrollOffset = len(m.previewLines) - 1
			}
		case "enter":
			m.state = ExportStateConfirm
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateConfirm(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			// Help modal is handled at the intent wrapper level (ExportArtifactIntent)
			return nil
		case KeyBack:
			m.state = ExportStatePreview
			return nil
		}

		switch msg.String() {
		case "y", "enter":
			m.state = ExportStateInProgress
			return m.startExport()
		case "n":
			m.state = ExportStatePreview
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateInProgress(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case ExportProgressMsg:
		return nil
	case ExportCompleteMsg:
		m.result = msg.Result
		m.state = ExportStateComplete
	case ExportErrorMsg:
		m.error = msg.Error
		m.state = ExportStateFailed
	case tea.KeyMsg:
		// Handle global keys (q=quit, ?=help)
		// Note: esc doesn't go back during export - let it complete
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			m.setResult(NewExportArtifactResultWithError(&IntentError{
				Code:    "export_cancelled",
				Message: "Export cancelled by user",
			}))
			return tea.Quit
		case KeyHelp:
			// Help modal is handled at the intent wrapper level (ExportArtifactIntent)
			return nil
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateComplete(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyBack:
			m.active = false
			return nil
		}

		switch msg.String() {
		case "enter":
			m.active = false
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateFailed(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			// Help modal is handled at the intent wrapper level (ExportArtifactIntent)
			return nil
		case KeyBack:
			m.active = false
			return nil
		}

		switch msg.String() {
		case "r":
			m.state = ExportStateConfirm
		}
	}
	return nil
}

// loadAvailableCVs fetches available CVs from the CV generation service
func (m *ExportArtifactModel) startExport() tea.Cmd {
	return func() tea.Msg {
		var result *ExportArtifactResult
		var err error

		// Call appropriate export function based on artifact type
		switch m.config.ArtifactType {
		case ExportTypeEvents:
			result, err = m.exportEvents()
		case ExportTypeFacts:
			result, err = m.exportFacts()
		case ExportTypeBursts:
			result, err = m.exportBursts()
		default:
			err = fmt.Errorf("unknown artifact type: %s", m.config.ArtifactType)
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

// exportEvents exports career events to the selected format and destination
func (m *ExportArtifactModel) exportEvents() (*ExportArtifactResult, error) {
	ctx := m.context.AppContext
	if ctx == nil {
		ctx = context.Background()
	}

	// Fetch all events from repository
	events, err := m.context.EventRepository.List(ctx, careerrepo.ListFilters{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events: %w", err)
	}

	// Marshal to selected format
	var content string
	switch m.config.Format {
	case ExportFormatJSON:
		content, err = marshalToJSON(events)
	case ExportFormatCSV:
		content, err = marshalEventsToCSV(events)
	case ExportFormatYAML:
		content, err = marshalToYAML(events)
	case ExportFormatTXT:
		content, err = marshalEventsToText(events)
	default:
		return nil, fmt.Errorf("unsupported events export format: %s", m.config.Format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to marshal events: %w", err)
	}

	// Save to destination
	filePath, err := m.saveToDestination(ctx, "career_events", content)
	if err != nil {
		return nil, err
	}

	size := int64(len(content))
	return NewExportArtifactResult(true, m.config.ArtifactType, m.config.Format, m.config.Destination, filePath, size), nil
}

// exportFacts exports facts to the selected format and destination
func (m *ExportArtifactModel) exportFacts() (*ExportArtifactResult, error) {
	ctx := m.context.AppContext
	if ctx == nil {
		ctx = context.Background()
	}

	// Fetch all facts from repository
	facts, err := m.context.FactRepository.List(ctx, careerrepo.FactListFilters{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch facts: %w", err)
	}

	// Marshal to selected format
	var content string
	switch m.config.Format {
	case ExportFormatJSON:
		content, err = marshalToJSON(facts)
	case ExportFormatCSV:
		content, err = marshalFactsToCSV(facts)
	case ExportFormatYAML:
		content, err = marshalToYAML(facts)
	case ExportFormatTXT:
		content, err = marshalFactsToText(facts)
	default:
		return nil, fmt.Errorf("unsupported facts export format: %s", m.config.Format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to marshal facts: %w", err)
	}

	// Save to destination
	filePath, err := m.saveToDestination(ctx, "facts", content)
	if err != nil {
		return nil, err
	}

	size := int64(len(content))
	return NewExportArtifactResult(true, m.config.ArtifactType, m.config.Format, m.config.Destination, filePath, size), nil
}

// exportBursts exports bursts to the selected format and destination
func (m *ExportArtifactModel) exportBursts() (*ExportArtifactResult, error) {
	ctx := m.context.AppContext
	if ctx == nil {
		ctx = context.Background()
	}

	// Fetch all bursts from repository
	bursts, err := m.context.BurstRepository.List(ctx, careerrepo.BurstListFilters{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bursts: %w", err)
	}

	// Marshal to selected format
	var content string
	switch m.config.Format {
	case ExportFormatJSON:
		content, err = marshalToJSON(bursts)
	case ExportFormatCSV:
		content, err = marshalBurstsToCSV(bursts)
	case ExportFormatYAML:
		content, err = marshalToYAML(bursts)
	case ExportFormatTXT:
		content, err = marshalBurstsToText(bursts)
	default:
		return nil, fmt.Errorf("unsupported bursts export format: %s", m.config.Format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to marshal bursts: %w", err)
	}

	// Save to destination
	filePath, err := m.saveToDestination(ctx, "bursts", content)
	if err != nil {
		return nil, err
	}

	size := int64(len(content))
	return NewExportArtifactResult(true, m.config.ArtifactType, m.config.Format, m.config.Destination, filePath, size), nil
}

// saveToDestination saves content to file or clipboard based on config
func (m *ExportArtifactModel) saveToDestination(ctx context.Context, name string, content string) (string, error) {
	var filePath string
	var err error

	switch m.config.Destination {
	case ExportDestinationFile:
		// Save to file
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}

		exportDir := filepath.Join(homeDir, ".kariya", "exports")
		if err := os.MkdirAll(exportDir, 0750); err != nil {
			return "", fmt.Errorf("failed to create export directory: %w", err)
		}

		timestamp := time.Now().Format("20060102_150405")
		extension := getFileExtensionForFormat(m.config.Format)
		filename := fmt.Sprintf("%s_%s%s", name, timestamp, extension)
		filePath = filepath.Join(exportDir, filename)

		if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
			return "", fmt.Errorf("failed to write file: %w", err)
		}

	case ExportDestinationClipboard:
		if m.context.ExportService == nil {
			return "", fmt.Errorf("export service not available")
		}
		if err = m.context.ExportService.CopyToClipboard(ctx, content); err != nil {
			return "", err
		}
		filePath = "clipboard"

	default:
		return "", fmt.Errorf("unsupported destination: %s", m.config.Destination)
	}

	return filePath, nil
}

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
	s += "\n↑/↓: Navigate | Enter: Select | Esc: Cancel | m: Main menu"
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
	s += "\n↑/↓: Navigate | Enter: Select | Esc: Back | m: Main menu"
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
	s += "\n↑/↓: Navigate | Enter: Select | Esc: Back | m: Main menu"
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
	s += "\nEnter: Continue | Esc: Back | m: Main menu"
	return s
}

func (m *ExportArtifactModel) viewPreview() string {
	if m.config == nil {
		return "Preview not available. Please configure the export first.\n\n[Enter] Continue | [Esc] Back"
	}
	s := "Preview Export (" + string(m.config.ArtifactType) + " as " + string(m.config.Format) + "):\n\n"
	s += "════════════════════════════════════════════════════════\n\n"

	viewHeight := 15
	endOffset := m.scrollOffset + viewHeight
	if endOffset > len(m.previewLines) {
		endOffset = len(m.previewLines)
	}

	for i := m.scrollOffset; i < endOffset; i++ {
		if i < len(m.previewLines) {
			s += m.previewLines[i] + "\n"
		}
	}

	s += "\n════════════════════════════════════════════════════════\n"
	s += "↑/↓: Scroll | PgUp/PgDn: Page | Enter: Confirm | Esc: Back | m: Main menu"

	if len(m.previewLines) > viewHeight {
		scrollPercent := (m.scrollOffset * 100) / len(m.previewLines)
		s += fmt.Sprintf("\n[%d%% scrolled]", scrollPercent)
	}

	return s
}

func (m *ExportArtifactModel) viewConfirm() string {
	s := "Confirm Export?\n\n"
	s += "Artifact: " + string(m.config.ArtifactType) + "\n"
	s += "Format: " + string(m.config.Format) + "\n"
	s += "Destination: " + string(m.config.Destination) + "\n\n"
	s += "Y/Enter: Confirm | N/Esc: Cancel | m: Main menu"
	return s
}

func (m *ExportArtifactModel) viewInProgress() string {
	return "Exporting artifact...\n\nPlease wait...\n\nEsc: Let export complete in background | m: Cancel and return to menu"
}

func (m *ExportArtifactModel) viewComplete() string {
	s := "Export Complete!\n\n"
	s += "File: " + m.result.FilePath + "\n"
	s += "Size: " + formatBytes(m.result.Size) + "\n\n"
	s += "Enter: Done | m: Main menu"
	return s
}

func (m *ExportArtifactModel) viewFailed() string {
	s := "Export Failed\n\n"
	if m.error != nil {
		s += "Error: " + m.error.Message + "\n"
	}
	s += "\nR: Retry | Esc: Cancel | m: Main menu"
	return s
}

// Preview generation functions

func (m *ExportArtifactModel) generatePreview() {
	var preview string

	switch m.config.ArtifactType {
	case ExportTypeEvents:
		preview = m.generateEventsPreview()
	case ExportTypeFacts:
		preview = m.generateFactsPreview()
	case ExportTypeBursts:
		preview = m.generateBurstsPreview()
	case ExportTypeProfile:
		preview = m.generateProfilePreview()
	default:
		preview = "No preview available for this artifact type"
	}

	m.preview = preview
	m.previewLines = strings.Split(preview, "\n")
}

func (m *ExportArtifactModel) generateEventsPreview() string {
	// Check if repository is available
	if m.context.EventRepository == nil {
		return "Preview not available - repository not initialized"
	}

	ctx := m.context.AppContext
	if ctx == nil {
		ctx = context.Background()
	}

	// Fetch events from repository
	events, err := m.context.EventRepository.List(ctx, careerrepo.ListFilters{
		Limit: 10, // Limit preview to first 10 events
	})
	if err != nil {
		return fmt.Sprintf("Error loading events: %v", err)
	}

	if len(events) == 0 {
		return "No events found.\n\nCreate some career events first to preview exports."
	}

	// Generate preview based on format
	var content string
	switch m.config.Format {
	case ExportFormatJSON:
		content, err = marshalToJSON(events)
	case ExportFormatCSV:
		content, err = marshalEventsToCSV(events)
	case ExportFormatYAML:
		content, err = marshalToYAML(events)
	case ExportFormatTXT:
		content, err = marshalEventsToText(events)
	default:
		return "Preview not available for format: " + string(m.config.Format)
	}

	if err != nil {
		return fmt.Sprintf("Error generating events preview: %v", err)
	}

	// Truncate if preview is too long
	if len(content) > 2000 {
		content = content[:2000] + "\n\n...(preview truncated, showing first 10 events)..."
	}

	return content
}

func (m *ExportArtifactModel) generateFactsPreview() string {
	// Check if repository is available
	if m.context.FactRepository == nil {
		return "Preview not available - repository not initialized"
	}

	ctx := m.context.AppContext
	if ctx == nil {
		ctx = context.Background()
	}

	// Fetch facts from repository
	facts, err := m.context.FactRepository.List(ctx, careerrepo.FactListFilters{
		Limit: 10, // Limit preview to first 10 facts
	})
	if err != nil {
		return fmt.Sprintf("Error loading facts: %v", err)
	}

	if len(facts) == 0 {
		return "No facts found.\n\nExtract some facts from your career events first to preview exports."
	}

	// Generate preview based on format
	var content string
	switch m.config.Format {
	case ExportFormatJSON:
		content, err = marshalToJSON(facts)
	case ExportFormatCSV:
		content, err = marshalFactsToCSV(facts)
	case ExportFormatYAML:
		content, err = marshalToYAML(facts)
	case ExportFormatTXT:
		content, err = marshalFactsToText(facts)
	default:
		return "Preview not available for format: " + string(m.config.Format)
	}

	if err != nil {
		return fmt.Sprintf("Error generating facts preview: %v", err)
	}

	// Truncate if preview is too long
	if len(content) > 2000 {
		content = content[:2000] + "\n\n...(preview truncated, showing first 10 facts)..."
	}

	return content
}

func (m *ExportArtifactModel) generateBurstsPreview() string {
	// Check if repository is available
	if m.context.BurstRepository == nil {
		return "Preview not available - repository not initialized"
	}

	ctx := m.context.AppContext
	if ctx == nil {
		ctx = context.Background()
	}

	// Fetch bursts from repository
	bursts, err := m.context.BurstRepository.List(ctx, careerrepo.BurstListFilters{
		Limit: 10, // Limit preview to first 10 bursts
	})
	if err != nil {
		return fmt.Sprintf("Error loading bursts: %v", err)
	}

	if len(bursts) == 0 {
		return "No bursts found.\n\nCreate some bursts from your career events first to preview exports."
	}

	// Generate preview based on format
	var content string
	switch m.config.Format {
	case ExportFormatJSON:
		content, err = marshalToJSON(bursts)
	case ExportFormatCSV:
		content, err = marshalBurstsToCSV(bursts)
	case ExportFormatYAML:
		content, err = marshalToYAML(bursts)
	case ExportFormatTXT:
		content, err = marshalBurstsToText(bursts)
	default:
		return "Preview not available for format: " + string(m.config.Format)
	}

	if err != nil {
		return fmt.Sprintf("Error generating bursts preview: %v", err)
	}

	// Truncate if preview is too long
	if len(content) > 2000 {
		content = content[:2000] + "\n\n...(preview truncated, showing first 10 bursts)..."
	}

	return content
}

func (m *ExportArtifactModel) generateProfilePreview() string {
	switch m.config.Format {
	case ExportFormatJSON:
		return `{
  "profile": {
    "personal": {
      "name": "John Doe",
      "email": "john@example.com",
      "location": "San Francisco, CA"
    },
    "professional": {
      "current_title": "Senior Software Engineer",
      "current_company": "Acme Corp",
      "years_experience": 8,
      "industries": ["Technology", "SaaS"]
    },
    "statistics": {
      "total_events": 24,
      "total_facts": 12,
      "total_bursts": 3,
      "avg_impact": "high"
    },
    "skills": {
      "technical": ["Go", "Rust", "Python", "TypeScript"],
      "soft": ["Leadership", "Communication", "Problem-solving"]
    }
  }
}`
	case ExportFormatPDF:
		return "[PDF Preview - Binary format]\n\nWhen exported, this will contain a formatted PDF version of your profile with professional styling."
	default:
		return "Preview not available for this format"
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	units := []string{"K", "M", "G", "T", "P"}
	for n := bytes / unit; n >= unit && exp < len(units)-1; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%d %sB", bytes/div, units[exp]) // #nosec G602 -- exp bounded by loop condition exp < len(units)-1
}

// DefaultArtifactTypes returns the default set of exportable artifact types
func DefaultArtifactTypes() []ExportArtifactType {
	return []ExportArtifactType{
		ExportTypeEvents,
		ExportTypeFacts,
		ExportTypeBursts,
	}
}

// DefaultSupportedFormats returns the default format support mapping
func DefaultSupportedFormats() map[ExportArtifactType][]ExportFormat {
	return map[ExportArtifactType][]ExportFormat{
		ExportTypeEvents: {
			ExportFormatJSON,
			ExportFormatCSV,
			ExportFormatTXT,
			ExportFormatYAML,
		},
		ExportTypeFacts: {
			ExportFormatJSON,
			ExportFormatCSV,
			ExportFormatTXT,
			ExportFormatYAML,
		},
		ExportTypeBursts: {
			ExportFormatJSON,
			ExportFormatCSV,
			ExportFormatTXT,
			ExportFormatYAML,
		},
	}
}

// DefaultFormats returns the default format for each artifact type
func DefaultFormats() map[ExportArtifactType]ExportFormat {
	return map[ExportArtifactType]ExportFormat{
		ExportTypeEvents: ExportFormatJSON,
		ExportTypeFacts:  ExportFormatJSON,
		ExportTypeBursts: ExportFormatJSON,
	}
}

// DefaultDestinations returns the default set of export destinations
func DefaultDestinations() []ExportDestination {
	return []ExportDestination{
		ExportDestinationFile,
		ExportDestinationClipboard,
	}
}

// mapToExportServiceFormat maps ExportFormat to cv.ExportFormat
// This is needed because ExportArtifact uses its own format constants
func mapToExportServiceFormat(format ExportFormat) cv.ExportFormat {
	switch format {
	case ExportFormatTXT:
		return cv.ExportFormatText
	case ExportFormatMD:
		return cv.ExportFormatMarkdown
	case ExportFormatYAML:
		return cv.ExportFormatYAML
	default:
		// Default to text for unknown formats (JSON, CSV, etc.)
		return cv.ExportFormatText
	}
}

// getFileExtensionForFormat returns the file extension for a given format
func getFileExtensionForFormat(format ExportFormat) string {
	switch format {
	case ExportFormatTXT:
		return ".txt"
	case ExportFormatMD:
		return ".md"
	case ExportFormatJSON:
		return ".json"
	case ExportFormatCSV:
		return ".csv"
	case ExportFormatYAML:
		return ".yaml"
	default:
		return ".txt"
	}
}

// marshalToJSON marshals any data to JSON format
func marshalToJSON(data interface{}) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return string(bytes), nil
}

// marshalToYAML marshals any data to YAML format
func marshalToYAML(data interface{}) (string, error) {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to YAML: %w", err)
	}
	return string(bytes), nil
}

// marshalEventsToCSV marshals career events to CSV format
func marshalEventsToCSV(events []*careerdomain.CareerEvent) (string, error) {
	var sb strings.Builder
	sb.WriteString("id,date,text,company,project,tags,categories\n")

	for _, event := range events {
		tags := strings.Join(event.Tags, ";")
		categories := strings.Join(event.Categories, ";")
		sb.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s\n",
			escapeCSV(event.ID),
			event.Date.Format("2006-01-02"),
			escapeCSV(event.Text),
			escapeCSV(event.Company),
			escapeCSV(event.Project),
			escapeCSV(tags),
			escapeCSV(categories),
		))
	}

	return sb.String(), nil
}

// marshalEventsToText marshals career events to plain text format
func marshalEventsToText(events []*careerdomain.CareerEvent) (string, error) {
	var sb strings.Builder
	sb.WriteString("Career Events\n")
	sb.WriteString(strings.Repeat("=", 80) + "\n\n")

	for i, event := range events {
		sb.WriteString(fmt.Sprintf("Event %d\n", i+1))
		sb.WriteString(fmt.Sprintf("Date: %s\n", event.Date.Format("2006-01-02")))
		if event.Company != "" {
			sb.WriteString(fmt.Sprintf("Company: %s\n", event.Company))
		}
		if event.Project != "" {
			sb.WriteString(fmt.Sprintf("Project: %s\n", event.Project))
		}
		sb.WriteString(fmt.Sprintf("Text: %s\n", event.Text))
		if len(event.Tags) > 0 {
			sb.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(event.Tags, ", ")))
		}
		if len(event.Categories) > 0 {
			sb.WriteString(fmt.Sprintf("Categories: %s\n", strings.Join(event.Categories, ", ")))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("\nTotal Events: %d\n", len(events)))
	return sb.String(), nil
}

// marshalFactsToCSV marshals facts to CSV format
func marshalFactsToCSV(facts []*careerdomain.Fact) (string, error) {
	var sb strings.Builder
	sb.WriteString("id,text,competency_categories,role_fit,audience_relevance,strength_signal,source_event_id,source_burst_id\n")

	for _, fact := range facts {
		categories := strings.Join(fact.CompetencyCategories, ";")
		audiences := strings.Join(fact.AudienceRelevance, ";")
		sb.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s\n",
			escapeCSV(fact.ID),
			escapeCSV(fact.Text),
			escapeCSV(categories),
			escapeCSV(string(fact.RoleFit)),
			escapeCSV(audiences),
			escapeCSV(fact.StrengthSignal),
			escapeCSV(fact.SourceEventID),
			escapeCSV(fact.SourceBurstID),
		))
	}

	return sb.String(), nil
}

// marshalFactsToText marshals facts to plain text format
func marshalFactsToText(facts []*careerdomain.Fact) (string, error) {
	var sb strings.Builder
	sb.WriteString("Facts\n")
	sb.WriteString(strings.Repeat("=", 80) + "\n\n")

	for i, fact := range facts {
		sb.WriteString(fmt.Sprintf("Fact %d: %s\n", i+1, fact.Text))
		sb.WriteString(fmt.Sprintf("Competency Categories: %s\n", strings.Join(fact.CompetencyCategories, ", ")))
		sb.WriteString(fmt.Sprintf("Role Fit: %s\n", fact.RoleFit))
		sb.WriteString(fmt.Sprintf("Audience Relevance: %s\n", strings.Join(fact.AudienceRelevance, ", ")))
		sb.WriteString(fmt.Sprintf("Strength Signal: %s\n", fact.StrengthSignal))
		if fact.SourceEventID != "" {
			sb.WriteString(fmt.Sprintf("Source Event: %s\n", fact.SourceEventID))
		}
		if fact.SourceBurstID != "" {
			sb.WriteString(fmt.Sprintf("Source Burst: %s\n", fact.SourceBurstID))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("\nTotal Facts: %d\n", len(facts)))
	return sb.String(), nil
}

// marshalBurstsToCSV marshals bursts to CSV format
func marshalBurstsToCSV(bursts []*careerdomain.Burst) (string, error) {
	var sb strings.Builder
	sb.WriteString("id,name,description,event_ids,confirmed\n")

	for _, burst := range bursts {
		eventIDs := strings.Join(burst.EventIDs, ";")
		sb.WriteString(fmt.Sprintf("%s,%s,%s,%s,%t\n",
			escapeCSV(burst.ID),
			escapeCSV(burst.Name),
			escapeCSV(burst.Description),
			escapeCSV(eventIDs),
			burst.Confirmed,
		))
	}

	return sb.String(), nil
}

// marshalBurstsToText marshals bursts to plain text format
func marshalBurstsToText(bursts []*careerdomain.Burst) (string, error) {
	var sb strings.Builder
	sb.WriteString("Bursts\n")
	sb.WriteString(strings.Repeat("=", 80) + "\n\n")

	for i, burst := range bursts {
		sb.WriteString(fmt.Sprintf("Burst %d: %s\n", i+1, burst.Name))
		if burst.Description != "" {
			sb.WriteString(fmt.Sprintf("Description: %s\n", burst.Description))
		}
		sb.WriteString(fmt.Sprintf("Event IDs: %s\n", strings.Join(burst.EventIDs, ", ")))
		sb.WriteString(fmt.Sprintf("Confirmed: %t\n", burst.Confirmed))
		if burst.ConfirmedAt != nil {
			sb.WriteString(fmt.Sprintf("Confirmed At: %s\n", burst.ConfirmedAt.Format("2006-01-02 15:04:05")))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("\nTotal Bursts: %d\n", len(bursts)))
	return sb.String(), nil
}

// escapeCSV escapes a string for CSV format
func escapeCSV(s string) string {
	// If string contains comma, quote, or newline, wrap in quotes and escape quotes
	if strings.ContainsAny(s, ",\"\n") {
		s = strings.ReplaceAll(s, "\"", "\"\"")
		return fmt.Sprintf("\"%s\"", s)
	}
	return s
}
