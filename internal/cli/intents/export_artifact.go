package intents

import (
	"context"
	"fmt"
	"strings"

	careerdomain "github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/service/career/cv"
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
	ExportService       *cv.ExportService
	CVGenerationService cv.CVGenerationService
	CareerService       *career.Service

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
	ExportStateSelectCV     ExportState = "select_cv"
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

// CVsLoadedMsg represents successful loading of available CVs
type CVsLoadedMsg struct {
	CVs []*careerdomain.CVView
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

	// CV selection (for ExportTypeCV only)
	availableCVs []*careerdomain.CVView
	selectedCV   *careerdomain.CVView
}

// NewExportArtifactModel creates a new ExportArtifact intent model
func NewExportArtifactModel(ctx *ExportArtifactContext) *ExportArtifactModel {
	return &ExportArtifactModel{
		state:         ExportStateSelectType,
		context:       ctx,
		selectedIndex: 0,
		active:        false,
		scrollOffset:  0,
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
	case ExportStateSelectCV:
		return m.updateSelectCV(msg)
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
	case ExportStateSelectCV:
		return m.viewSelectCV()
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
		case "up", "k":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "down", "j":
			if m.selectedIndex < len(m.context.ArtifactTypes)-1 {
				m.selectedIndex++
			}
		case "enter":
			artifactType := m.context.ArtifactTypes[m.selectedIndex]
			m.config = NewExportConfiguration(artifactType, m.context)
			// If CV type, go to CV selection first
			if artifactType == ExportTypeCV {
				m.state = ExportStateSelectCV
				m.selectedIndex = 0
				return m.loadAvailableCVs()
			}
			// For other types, go directly to format selection
			m.state = ExportStateSelectFormat
		case "esc":
			m.setResult(NewExportArtifactResultWithError(&IntentError{
				Code:    "export_cancelled",
				Message: "Export cancelled by user",
			}))
		case "m":
			m.setResult(NewExportArtifactResultWithError(&IntentError{
				Code:    "export_cancelled",
				Message: "Export cancelled by user",
			}))
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateSelectCV(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case CVsLoadedMsg:
		m.availableCVs = msg.CVs
		return nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			} else if len(m.availableCVs) > 0 {
				// Wrap to bottom
				m.selectedIndex = len(m.availableCVs) - 1
			}
		case "down", "j":
			if m.selectedIndex < len(m.availableCVs)-1 {
				m.selectedIndex++
			} else if len(m.availableCVs) > 0 {
				// Wrap to top
				m.selectedIndex = 0
			}
		case "enter":
			if len(m.availableCVs) > 0 {
				m.selectedCV = m.availableCVs[m.selectedIndex]
				m.selectedIndex = 0
				m.state = ExportStateSelectFormat
			}
		case "esc":
			m.selectedIndex = 0
			m.state = ExportStateSelectType
		case "m":
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
		case "up", "k":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "down", "j":
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
			// Go back to CV selection if artifact type is CV, otherwise SelectType
			if m.config.ArtifactType == ExportTypeCV {
				m.state = ExportStateSelectCV
			} else {
				m.state = ExportStateSelectType
			}
		case "m":
			m.setResult(NewExportArtifactResultWithError(&IntentError{
				Code:    "export_cancelled",
				Message: "Export cancelled by user",
			}))
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateSelectDest(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "down", "j":
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
		case "m":
			m.setResult(NewExportArtifactResultWithError(&IntentError{
				Code:    "export_cancelled",
				Message: "Export cancelled by user",
			}))
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateConfigure(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.generatePreview()
			m.scrollOffset = 0
			m.state = ExportStatePreview
		case "esc":
			m.selectedIndex = 0
			m.state = ExportStateSelectDest
		case "m":
			m.setResult(NewExportArtifactResultWithError(&IntentError{
				Code:    "export_cancelled",
				Message: "Export cancelled by user",
			}))
		}
	}
	return nil
}

func (m *ExportArtifactModel) updatePreview(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
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
		case "esc":
			m.state = ExportStateConfigure
		case "m":
			m.setResult(NewExportArtifactResultWithError(&IntentError{
				Code:    "export_cancelled",
				Message: "Export cancelled by user",
			}))
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
		case "m":
			m.setResult(NewExportArtifactResultWithError(&IntentError{
				Code:    "export_cancelled",
				Message: "Export cancelled by user",
			}))
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
		switch msg.String() {
		case "esc":
			// Let export complete in background
			return nil
		case "m":
			// Return to main menu immediately (cancel)
			m.setResult(NewExportArtifactResultWithError(&IntentError{
				Code:    "export_cancelled",
				Message: "Export cancelled by user",
			}))
			return nil
		}
	}
	return nil
}

func (m *ExportArtifactModel) updateComplete(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "esc":
			m.active = false
		case "m":
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
			m.state = ExportStateConfirm
		case "esc":
			m.active = false
		case "m":
			m.active = false
		}
	}
	return nil
}

// loadAvailableCVs fetches available CVs from the CV generation service
func (m *ExportArtifactModel) loadAvailableCVs() tea.Cmd {
	return func() tea.Msg {
		// TODO: In real implementation, fetch from CVGenerationService
		// For now, return empty list (tests will handle mocking)
		return CVsLoadedMsg{CVs: []*careerdomain.CVView{}}
	}
}

func (m *ExportArtifactModel) startExport() tea.Cmd {
	return func() tea.Msg {
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

func (m *ExportArtifactModel) viewSelectCV() string {
	s := "Select CV to Export:\n\n"

	if len(m.availableCVs) == 0 {
		s += "No CVs available. Please generate a CV first.\n"
		s += "\nEsc: Back | m: Main menu"
		return s
	}

	for i, cv := range m.availableCVs {
		prefix := "  "
		if i == m.selectedIndex {
			prefix = "▶ "
		}
		s += fmt.Sprintf("%s%s - %s\n", prefix, cv.Name, cv.TargetRole)
	}
	s += "\n↑/↓ or j/k: Navigate | Enter: Select | Esc: Back | m: Main menu"
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
	case ExportTypeCV:
		preview = m.generateCVPreview()
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

func (m *ExportArtifactModel) generateCVPreview() string {
	switch m.config.Format {
	case ExportFormatJSON:
		return `{
  "cv": {
    "profile": {
      "name": "John Doe",
      "title": "Senior Software Engineer",
      "summary": "Experienced software engineer with 8+ years in full-stack development"
    },
    "experience": [
      {
        "company": "Tech Corp",
        "title": "Senior Engineer",
        "duration": "2022 - Present",
        "description": "Led team of 5 engineers on microservices architecture"
      }
    ],
    "skills": ["Go", "Rust", "Python", "TypeScript"],
    "education": [
      {
        "institution": "University",
        "degree": "BS Computer Science",
        "year": "2016"
      }
    ]
  }
}`
	case ExportFormatMD:
		return `# John Doe - Senior Software Engineer

## Summary
Experienced software engineer with 8+ years in full-stack development

## Experience

### Tech Corp | Senior Engineer (2022 - Present)
- Led team of 5 engineers on microservices architecture
- Improved system performance by 40%
- Mentored 3 junior developers

### Previous Company | Engineer (2020 - 2022)
- Developed REST APIs serving 1M+ requests daily

## Skills
Go, Rust, Python, TypeScript, Kubernetes, Docker

## Education
**BS Computer Science** - University (2016)`
	case ExportFormatPDF:
		return "[PDF Preview - Binary format]\n\nWhen exported, this will contain a formatted PDF version of your CV with professional styling."
	default:
		return "Preview not available for this format"
	}
}

func (m *ExportArtifactModel) generateEventsPreview() string {
	switch m.config.Format {
	case ExportFormatJSON:
		return `{
  "events": [
    {
      "id": "evt-001",
      "title": "Led team standup meeting",
      "date": "2025-12-15",
      "company": "Acme Corp",
      "project": "Platform Redesign",
      "description": "Facilitated daily standup with 8 team members",
      "impact": "high",
      "tags": ["leadership", "communication"]
    },
    {
      "id": "evt-002",
      "title": "Completed API integration",
      "date": "2025-12-14",
      "company": "Acme Corp",
      "project": "Platform Redesign",
      "description": "Integrated payment gateway API",
      "impact": "medium",
      "tags": ["technical", "backend"]
    }
  ],
  "total_events": 24,
  "date_range": "2025-12-01 to 2025-12-31"
}`
	case ExportFormatCSV:
		return `date,title,company,project,description,impact,tags
2025-12-15,Led team standup meeting,Acme Corp,Platform Redesign,Facilitated daily standup with 8 team members,high,leadership;communication
2025-12-14,Completed API integration,Acme Corp,Platform Redesign,Integrated payment gateway API,medium,technical;backend
2025-12-13,Code review completed,Acme Corp,Platform Redesign,Reviewed 12 pull requests from team,medium,code-review
...
Total: 24 events`
	case ExportFormatTXT:
		return `CAREER EVENTS EXPORT
====================

Event 1: Led team standup meeting
Date: 2025-12-15
Company: Acme Corp
Project: Platform Redesign
Description: Facilitated daily standup with 8 team members
Impact: High
Tags: leadership, communication

Event 2: Completed API integration
Date: 2025-12-14
Company: Acme Corp
Project: Platform Redesign
Description: Integrated payment gateway API
Impact: Medium
Tags: technical, backend

...

Total Events: 24
Date Range: 2025-12-01 to 2025-12-31`
	default:
		return "Preview not available for this format"
	}
}

func (m *ExportArtifactModel) generateFactsPreview() string {
	switch m.config.Format {
	case ExportFormatJSON:
		return `{
  "facts": [
    {
      "id": "fact-001",
      "title": "Led 8-person team",
      "category": "leadership",
      "verified": true,
      "source_events": ["evt-001", "evt-002"],
      "confidence": 0.95,
      "extracted_date": "2025-12-15"
    },
    {
      "id": "fact-002",
      "title": "Proficient in Go and Rust",
      "category": "technical",
      "verified": true,
      "source_events": ["evt-003", "evt-004"],
      "confidence": 0.98,
      "extracted_date": "2025-12-14"
    }
  ],
  "total_facts": 12
}`
	case ExportFormatCSV:
		return `id,title,category,verified,confidence,source_events,extracted_date
fact-001,Led 8-person team,leadership,true,0.95,"evt-001,evt-002",2025-12-15
fact-002,Proficient in Go and Rust,technical,true,0.98,"evt-003,evt-004",2025-12-14
fact-003,5+ years backend development,experience,true,0.92,"evt-005,evt-006",2025-12-13
...
Total: 12 facts`
	case ExportFormatTXT:
		return `EXTRACTED FACTS
================

Fact 1: Led 8-person team
Category: Leadership
Verified: Yes
Confidence: 95%
Source Events: 2
Extracted: 2025-12-15

Fact 2: Proficient in Go and Rust
Category: Technical
Verified: Yes
Confidence: 98%
Source Events: 2
Extracted: 2025-12-14

...

Total Facts: 12`
	default:
		return "Preview not available for this format"
	}
}

func (m *ExportArtifactModel) generateBurstsPreview() string {
	switch m.config.Format {
	case ExportFormatJSON:
		return `{
  "bursts": [
    {
      "id": "burst-001",
      "title": "Platform Redesign Phase 1",
      "start_date": "2025-11-01",
      "end_date": "2025-12-15",
      "description": "Complete redesign of user-facing platform",
      "impact": "high",
      "team_size": 8,
      "key_achievements": [
        "30% improvement in page load time",
        "Reduced API calls by 40%",
        "Improved user satisfaction from 3.2 to 4.5 stars"
      ]
    }
  ],
  "total_bursts": 3
}`
	case ExportFormatCSV:
		return `id,title,start_date,end_date,duration_days,team_size,impact,key_achievements
burst-001,Platform Redesign Phase 1,2025-11-01,2025-12-15,45,8,high,"30% improvement in page load time; Reduced API calls by 40%"
burst-002,Backend Optimization,2025-10-15,2025-11-30,46,5,medium,"Reduced database queries by 50%"
...
Total: 3 bursts`
	case ExportFormatTXT:
		return `CAREER BURSTS
==============

Burst 1: Platform Redesign Phase 1
Duration: 2025-11-01 to 2025-12-15 (45 days)
Team Size: 8 people
Impact: High

Key Achievements:
- 30% improvement in page load time
- Reduced API calls by 40%
- Improved user satisfaction from 3.2 to 4.5 stars

Description: Complete redesign of user-facing platform

...

Total Bursts: 3`
	default:
		return "Preview not available for this format"
	}
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
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"K", "M", "G", "T", "P"}
	return fmt.Sprintf("%d %sB", bytes/div, units[exp])
}

// DefaultArtifactTypes returns the default set of exportable artifact types
func DefaultArtifactTypes() []ExportArtifactType {
	return []ExportArtifactType{
		ExportTypeCV,
		ExportTypeEvents,
		ExportTypeFacts,
		ExportTypeBursts,
	}
}

// DefaultSupportedFormats returns the default format support mapping
func DefaultSupportedFormats() map[ExportArtifactType][]ExportFormat {
	return map[ExportArtifactType][]ExportFormat{
		ExportTypeCV: {
			ExportFormatTXT,
			ExportFormatMD,
			ExportFormatYAML,
		},
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
		ExportTypeCV:     ExportFormatMD,
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
