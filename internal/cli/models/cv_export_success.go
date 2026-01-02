package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
)

// CVExportSuccessModel displays successful CV export information
type CVExportSuccessModel struct {
	*BaseStandardModel
	cvView       *career.CVView
	format       string
	filePath     string
	header       components.HeaderModel
	helpFooter   components.HelpFooterModel
	width        int
	height       int
	breadcrumbs  []string
	selectedIdx  int
	options      []string
}

// NewCVExportSuccessModel creates a new CV Export Success model
func NewCVExportSuccessModel(
	baseModel *BaseStandardModel,
	cvView *career.CVView,
	format string,
	filePath string,
) *CVExportSuccessModel {
	m := &CVExportSuccessModel{
		BaseStandardModel: baseModel,
		cvView:            cvView,
		format:            format,
		filePath:          filePath,
		header:            components.NewHeader("✅ CV Exported Successfully", 80),
		helpFooter:        components.NewHelpFooter("cv_export_success", 80),
		width:             80,
		height:            20,
		breadcrumbs:       []string{"Home", "CV Management", "Export", "Success"},
		selectedIdx:       0,
		options:           []string{"View Export Details", "Back to Preview", "Back to Home"},
	}
	m.header.SetBreadcrumbs(m.breadcrumbs)
	return m
}

// Init initializes the model
func (m *CVExportSuccessModel) Init() tea.Cmd {
	return m.BaseStandardModel.Init()
}

// Update handles messages and updates the model state
func (m *CVExportSuccessModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.header.SetWidth(msg.Width)
		m.helpFooter.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
			return m, nil

		case "down", "j":
			if m.selectedIdx < len(m.options)-1 {
				m.selectedIdx++
			}
			return m, nil

		case "enter":
			switch m.selectedIdx {
			case 0:
				// Show export details
				return m, func() tea.Msg {
					return ShowExportDetailsMsg{
						CVView:   m.cvView,
						Format:   m.format,
						FilePath: m.filePath,
					}
				}
			case 1:
				// Back to preview
				return m, func() tea.Msg {
					return BackMsg{}
				}
			case 2:
				// Back to home
				return m, func() tea.Msg {
					return NavigateToHomeMsg{}
				}
			}
			return m, nil

		case "esc", "q":
			return m, func() tea.Msg {
				return BackMsg{}
			}
		}
	}

	return m, nil
}

// View renders the export success screen
func (m *CVExportSuccessModel) View() string {
	headerView := m.header.View()

	// Create the success message with export details
	successMessage := fmt.Sprintf(`
📄 CV Export Successful!

Configuration: %s
Format: %s
Location: %s

File Details:
  • Name: %s
  • Format: %s
  • Ready for use

What would you like to do next?

`, m.cvView.Name, m.format, m.filePath, m.cvView.Name, m.format)

	// Render options
	optionsView := "\n"
	for i, option := range m.options {
		indicator := "  "
		if i == m.selectedIdx {
			indicator = "▶ "
			optionsView += styles.InputHint.Render(indicator + option) + "\n"
		} else {
			optionsView += indicator + option + "\n"
		}
	}

	footer := "\nPress 'j'/'k' to navigate, 'enter' to select, 'esc' to go back"

	content := successMessage + optionsView + footer

	return fmt.Sprintf("%s\n%s", headerView, content)
}

// Messages for CV Export Success

// ShowExportDetailsMsg shows detailed export information
type ShowExportDetailsMsg struct {
	CVView   *career.CVView
	Format   string
	FilePath string
}

// NavigateToHomeMsg navigates to the home screen
type NavigateToHomeMsg struct{}
