package components

import (
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// SpinnerModel wraps the bubbles spinner with our custom styling
type SpinnerModel struct {
	spinner spinner.Model
}

// NewSpinner creates a new spinner with default styling
func NewSpinner() *SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.SpinnerStyle
	return &SpinnerModel{
		spinner: s,
	}
}

// Init initializes the spinner model
func (m *SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles spinner tick messages
func (m *SpinnerModel) Update(msg tea.Msg) (*SpinnerModel, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// View renders the spinner
func (m *SpinnerModel) View() string {
	return m.spinner.View()
}
