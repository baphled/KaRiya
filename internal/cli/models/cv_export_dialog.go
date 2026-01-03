package models

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
)

// CVExportDialogModel displays export format options for a generated CV
type CVExportDialogModel struct {
	*BaseStandardModel
	cvView         *career.CVView
	sections       []*career.CVSection
	exportService  *cv.ExportService
	selectedFormat cv.ExportFormat
	formats        []ExportFormatOption
	cursor         int
	exporting      bool
	header         components.HeaderModel
	helpFooter     components.HelpFooterModel
	width          int
	height         int
}

// ExportFormatOption represents an export format choice
type ExportFormatOption struct {
	Format      cv.ExportFormat
	Name        string
	Description string
}

// NewCVExportDialogModel creates a new CV Export Dialog model
func NewCVExportDialogModel(
	baseModel *BaseStandardModel,
	cvView *career.CVView,
	sections []*career.CVSection,
	exportService *cv.ExportService,
) *CVExportDialogModel {
	formats := []ExportFormatOption{
		{
			Format:      cv.ExportFormatText,
			Name:        "Plain Text",
			Description: "Export as plain text file",
		},
		{
			Format:      cv.ExportFormatMarkdown,
			Name:        "Markdown",
			Description: "Export as Markdown file (best for GitHub/web)",
		},
		{
			Format:      cv.ExportFormatYAML,
			Name:        "YAML",
			Description: "Export as YAML file (for data processing)",
		},
	}

	m := &CVExportDialogModel{
		BaseStandardModel: baseModel,
		cvView:            cvView,
		sections:          sections,
		exportService:     exportService,
		selectedFormat:    cv.ExportFormatText,
		formats:           formats,
		cursor:            0,
		exporting:         false,
		header:            components.NewHeader("📤 Export CV", 80),
		helpFooter:        components.NewHelpFooter("cv_export_dialog", 80),
		width:             80,
		height:            20,
	}
	m.header.SetBreadcrumbs([]string{"Home", "CV Management", "Preview", "Export"})
	return m
}

// Init initializes the model
func (m *CVExportDialogModel) Init() tea.Cmd {
	return m.BaseStandardModel.Init()
}

// Update handles messages and updates the model state
func (m *CVExportDialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.header.SetWidth(msg.Width)
		m.helpFooter.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		if m.exporting {
			// Don't handle input while exporting
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			m.selectedFormat = m.formats[m.cursor].Format
			return m, nil

		case "down", "j":
			if m.cursor < len(m.formats)-1 {
				m.cursor++
			}
			m.selectedFormat = m.formats[m.cursor].Format
			return m, nil

		case "enter":
			// Perform export immediately (it's very fast)
			return m, m.performExport()

		case "esc", "q":
			// Cancel export
			return m, func() tea.Msg {
				return BackMsg{}
			}
		}
	}

	return m, nil
}

// performExport performs the actual export operation
func (m *CVExportDialogModel) performExport() tea.Cmd {
	return func() tea.Msg {
		if m.exportService == nil {
			return CVExportErrorMsg{Err: fmt.Errorf("export service not initialized")}
		}

		ctx := context.Background()

		// Export to the selected format
		var content string
		var err error

		switch m.selectedFormat {
		case cv.ExportFormatText:
			content, err = m.exportService.ExportToText(ctx, m.cvView, m.sections, nil)
		case cv.ExportFormatMarkdown:
			content, err = m.exportService.ExportToMarkdown(ctx, m.cvView, m.sections, nil)
		case cv.ExportFormatYAML:
			content, err = m.exportService.ExportToYAML(ctx, m.cvView, m.sections, nil)
		default:
			err = fmt.Errorf("unsupported export Format: %v", m.selectedFormat)
		}

		if err != nil {
			return CVExportErrorMsg{Err: err}
		}

		// Save to file
		filePath, err := m.exportService.SaveToFile(ctx, m.cvView.Name, m.selectedFormat, content)
		if err != nil {
			return CVExportErrorMsg{Err: err}
		}

		return CVExportedMsg{
			CVView:   m.cvView,
			Format:   m.formats[m.cursor].Name,
			FilePath: filePath,
		}
	}
}

// View renders the export dialog
func (m *CVExportDialogModel) View() string {
	headerView := m.header.View()

	if m.exporting {
		return fmt.Sprintf("%s\n\nExporting CV in %s format...", headerView, m.formats[m.cursor].Name)
	}

	// Render format options
	var optionsView string
	for i, format := range m.formats {
		indicator := "  "
		if i == m.cursor {
			indicator = "▶ "
		}

		line := fmt.Sprintf("%s%s\n", indicator, format.Name)
		if i == m.cursor {
			line = styles.InputHint.Render(line)
		}
		optionsView += line
	}

	content := fmt.Sprintf(`
Select export Format:

%s
Press 'j'/'k' to navigate, 'enter' to export, 'esc' to cancel
`, optionsView)

	return fmt.Sprintf("%s\n%s", headerView, content)
}

// Messages for CV Export Dialog

// StartCVExportMsg initiates the CV export process
type StartCVExportMsg struct {
	CVView       *career.CVView
	Format       string
	ExportFormat cv.ExportFormat
}
