package modals

import (
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ExportModal provides a simple 2-field form for export configuration.
// It allows users to select export format (Text/Markdown/YAML) and
// location (File/Clipboard).
//
// Usage:
//
//	modal := modals.NewExportModal(120, 40)
//	cmd := modal.Init()
//	// In Update:
//	cmd := modal.Update(msg)
//	if modal.IsCompleted() {
//	    data := modal.GetExportData()
//	    // Use data.Format and data.Location
//	}
type ExportModal struct {
	form      forms.Form
	data      *ExportData
	visible   bool
	completed bool
	width     int
	height    int
}

// ExportData holds the export configuration selections.
type ExportData struct {
	Format   string // "text" | "markdown" | "yaml"
	Location string // "file" | "clipboard"
}

// NewExportModal creates a new export options modal.
//
// Expected:
//   - width and height must be valid positive integers.
//
// Returns:
//   - A fully initialized ExportModal ready for use.
//
// Side effects:
//   - None.
func NewExportModal(width, height int) *ExportModal {
	modal := &ExportModal{
		data:    &ExportData{},
		visible: true,
		width:   width,
		height:  height,
	}

	modal.buildForm()
	return modal
}

// buildForm creates the form with 2 fields (Format + Location).
func (m *ExportModal) buildForm() {
	modalWidth := calcModalWidth(m.width)

	formatOptions := []forms.SelectOption{
		{Key: "text", Value: "Plain Text"},
		{Key: "markdown", Value: "Markdown"},
		{Key: "yaml", Value: "YAML"},
	}

	locationOptions := []forms.SelectOption{
		{Key: "file", Value: "File"},
		{Key: "clipboard", Value: "Clipboard"},
	}

	formatField := forms.NewSelect("format", "Export Format", "Choose output format for the CV", formatOptions).
		Value(&m.data.Format)

	locationField := forms.NewSelect("location", "Save To", "Where to save the exported CV", locationOptions).
		Value(&m.data.Location)

	group := forms.NewGroup(formatField, locationField)

	m.form = forms.NewFormWithDimensions(modalWidth, 0, group)
}

// Init initializes the export options modal and its form.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *ExportModal) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the export options modal.
//
// Expected:
//   - msg must be a valid tea.Msg.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - May update internal state based on key messages.
func (m *ExportModal) Update(msg tea.Msg) tea.Cmd {
	if !m.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			m.visible = false
			return nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.buildForm()
		return m.form.Init()
	}

	if m.form != nil {
		var cmd tea.Cmd
		m.form, cmd = forms.Update(m.form, msg)

		if forms.IsCompleted(m.form) {
			m.completed = true
			m.visible = false
		}

		return cmd
	}

	return nil
}

// View renders the export options modal.
//
// Returns:
//   - A string containing the rendered modal view.
//
// Side effects:
//   - None.
func (m *ExportModal) View() string {
	if !m.visible {
		return ""
	}

	if m.form == nil {
		return ""
	}

	modalWidth := calcModalWidth(m.width)

	modalHeight := m.height - 10
	if modalHeight < 15 {
		modalHeight = 15
	}

	th := theme.Default()

	title := primitives.Title("Export Options", th).
		Width(modalWidth - 4).
		Center().
		Render()

	formView := m.form.View()

	footer := m.buildFooter()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		formView,
		"",
		footer,
	)

	return containers.NewBox(th).
		Content(content).
		Width(modalWidth).
		MaxHeight(modalHeight).
		Padding(1).
		Background(th.BackgroundColor()).
		Variant(containers.BoxInfo).
		Render()
}

// buildFooter creates the keyboard shortcuts footer using UIKit primitives.
func (m *ExportModal) buildFooter() string {
	th := theme.Default()

	return primitives.RenderHelpFooter(th,
		primitives.KeyBadge("Enter", "Export", th),
		primitives.CancelBadge(th),
	)
}

// Show makes the modal visible.
//
// Side effects:
//   - Updates modal visibility state.
func (m *ExportModal) Show() {
	m.visible = true
}

// Hide makes the modal invisible.
//
// Side effects:
//   - Updates modal visibility state.
func (m *ExportModal) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is currently visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ExportModal) IsVisible() bool {
	return m.visible
}

// IsCompleted returns whether the form has been completed.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ExportModal) IsCompleted() bool {
	return m.completed
}

// GetExportData returns the export configuration data.
//
// Returns:
//   - A pointer to ExportData containing format and location selections.
//
// Side effects:
//   - None.
func (m *ExportModal) GetExportData() *ExportData {
	return m.data
}

// SetFormat sets the export format.
//
// Expected:
//   - format must be a valid string ("text", "markdown", or "yaml").
//
// Side effects:
//   - Updates the format in export data.
func (m *ExportModal) SetFormat(format string) {
	m.data.Format = format
}

// SetLocation sets the save location.
//
// Expected:
//   - location must be a valid string ("file" or "clipboard").
//
// Side effects:
//   - Updates the location in export data.
func (m *ExportModal) SetLocation(location string) {
	m.data.Location = location
}

// Complete marks the modal as completed and hides it.
//
// Side effects:
//   - Sets completed flag and hides the modal.
func (m *ExportModal) Complete() {
	m.completed = true
	m.visible = false
}

// GetDimensions returns the current modal dimensions.
//
// Returns:
//   - width and height as int values.
//
// Side effects:
//   - None.
func (m *ExportModal) GetDimensions() (width, height int) {
	return m.width, m.height
}
