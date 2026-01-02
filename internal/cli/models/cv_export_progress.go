package models

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/components"
)

// CVExportProgressModel shows export progress
type CVExportProgressModel struct {
	*BaseStandardModel
	cvName      string
	format      string
	startTime   time.Time
	header      components.HeaderModel
	width       int
	height      int
	breadcrumbs []string
	spinnerIdx  int
	done        bool
}

// NewCVExportProgressModel creates a new CV Export Progress model
func NewCVExportProgressModel(
	baseModel *BaseStandardModel,
	cvName string,
	format string,
) *CVExportProgressModel {
	m := &CVExportProgressModel{
		BaseStandardModel: baseModel,
		cvName:            cvName,
		format:            format,
		startTime:         time.Now(),
		header:            components.NewHeader("📤 Exporting CV", 80),
		width:             80,
		height:            20,
		breadcrumbs:       []string{"Home", "CV Management", "Export", "Processing"},
		spinnerIdx:        0,
		done:              false,
	}
	m.header.SetBreadcrumbs(m.breadcrumbs)
	return m
}

// Init initializes the model
func (m *CVExportProgressModel) Init() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg{t}
	})
}

// Update handles messages and updates the model state
func (m *CVExportProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.header.SetWidth(msg.Width)
		return m, nil

	case TickMsg:
		if m.done {
			return m, nil
		}
		m.spinnerIdx++
		// Continue ticking
		return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return TickMsg{t}
		})

	case CVExportCompleteMsg:
		// Export is complete, show success
		m.done = true
		return m, func() tea.Msg {
			return msg // Pass through to parent
		}

	case tea.KeyMsg:
		// Don't allow cancellation during export
		return m, nil
	}

	return m, nil
}

// View renders the export progress screen
func (m *CVExportProgressModel) View() string {
	headerView := m.header.View()

	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	elapsed := time.Since(m.startTime)
	spinnerFrame := spinner[(m.spinnerIdx)%len(spinner)]

	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("6")).
		Bold(true)

	progressView := fmt.Sprintf(`
%s Exporting CV...

Configuration: %s
Format: %s
Elapsed: %v

Processing your CV and preparing the export file...
`, spinnerFrame, m.cvName, m.format, elapsed.Round(100*time.Millisecond))

	return fmt.Sprintf("%s\n%s", headerView, progressStyle.Render(progressView))
}

// TickMsg is sent to update the spinner
type TickMsg struct {
	Time time.Time
}

// CVExportCompleteMsg signals export completion
type CVExportCompleteMsg struct {
	CVView   interface{}
	Format   string
	FilePath string
}
