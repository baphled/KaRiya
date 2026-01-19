package components

import (
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// ExportOptionsModal provides a simple 2-field form for export configuration.
// It allows users to select export format (Text/Markdown/YAML) and
// location (File/Clipboard).
//
// Usage:
//
//	modal := components.NewExportOptionsModal(120, 40)
//	cmd := modal.Init()
//	// In Update:
//	cmd := modal.Update(msg)
//	if modal.IsCompleted() {
//	    data := modal.GetExportData()
//	    // Use data.Format and data.Location
//	}
type ExportOptionsModal struct {
	form      *huh.Form
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

// NewExportOptionsModal creates a new export options modal.
// width, height: Terminal dimensions for responsive sizing
func NewExportOptionsModal(width, height int) *ExportOptionsModal {
	modal := &ExportOptionsModal{
		data: &ExportData{
			// huh.Select will auto-select first options:
			// Format: "text" (first option)
			// Location: "file" (first option)
		},
		visible: true,
		width:   width,
		height:  height,
	}

	modal.buildForm()
	return modal
}

// buildForm creates the huh form with 2 fields (Format + Location).
func (m *ExportOptionsModal) buildForm() {
	// Calculate modal width
	modalWidth := m.width - 20
	if modalWidth > 60 {
		modalWidth = 60
	}
	if modalWidth < 40 {
		modalWidth = 40
	}

	// Create fields
	formatField := huh.NewSelect[string]().
		Key("format").
		Title("Export Format").
		Description("Choose output format for the CV").
		Options(
			huh.NewOption("Plain Text", "text"),
			huh.NewOption("Markdown", "markdown"),
			huh.NewOption("YAML", "yaml"),
		).
		Value(&m.data.Format)

	locationField := huh.NewSelect[string]().
		Key("location").
		Title("Save To").
		Description("Where to save the exported CV").
		Options(
			huh.NewOption("File", "file"),
			huh.NewOption("Clipboard", "clipboard"),
		).
		Value(&m.data.Location)

	// Create form with single group
	group := huh.NewGroup(formatField, locationField)

	m.form = huh.NewForm(group).
		WithWidth(modalWidth).
		WithShowHelp(true).
		WithShowErrors(true)
}

// Init initializes the export options modal and its form.
func (m *ExportOptionsModal) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the export options modal.
func (m *ExportOptionsModal) Update(msg tea.Msg) tea.Cmd {
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
		// Rebuild form with new dimensions
		m.buildForm()
		return m.form.Init()
	}

	// Update form
	if m.form != nil {
		var cmd tea.Cmd
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}

		// Check if form completed
		if m.form.State == huh.StateCompleted {
			m.completed = true
			m.visible = false
		}

		return cmd
	}

	return nil
}

// View renders the export options modal.
func (m *ExportOptionsModal) View() string {
	if !m.visible {
		return ""
	}

	if m.form == nil {
		return ""
	}

	// Calculate modal dimensions
	modalWidth := m.width - 20
	if modalWidth > 60 {
		modalWidth = 60
	}
	if modalWidth < 40 {
		modalWidth = 40
	}

	modalHeight := m.height - 10
	if modalHeight < 15 {
		modalHeight = 15
	}

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(styles.ColorAccentTeal).
		Bold(true).
		Align(lipgloss.Center).
		Width(modalWidth - 4)

	title := titleStyle.Render("Export Options")

	// Render form
	formView := m.form.View()

	// Footer with keyboard shortcuts
	footer := m.buildFooter()

	// Combine all parts
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		formView,
		"",
		footer,
	)

	// Wrap in styled container with solid background
	styledContent := lipgloss.NewStyle().
		Width(modalWidth).
		MaxHeight(modalHeight).
		Background(styles.ColorBackground). // CRITICAL: Solid background
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorAccentTeal).
		Padding(1).
		Render(content)

	return styledContent
}

// buildFooter creates the keyboard shortcuts footer using UIKit primitives.
func (m *ExportOptionsModal) buildFooter() string {
	th := theme.Default()

	return primitives.RenderHelpFooter(th,
		primitives.KeyBadge("Enter", "Export", th),
		primitives.CancelBadge(th),
	)
}

// Show makes the modal visible.
func (m *ExportOptionsModal) Show() {
	m.visible = true
}

// Hide makes the modal invisible.
func (m *ExportOptionsModal) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is currently visible.
func (m *ExportOptionsModal) IsVisible() bool {
	return m.visible
}

// IsCompleted returns whether the form has been completed.
func (m *ExportOptionsModal) IsCompleted() bool {
	return m.completed
}

// GetExportData returns the export configuration data.
func (m *ExportOptionsModal) GetExportData() *ExportData {
	return m.data
}

// SetFormat sets the export format.
func (m *ExportOptionsModal) SetFormat(format string) {
	m.data.Format = format
}

// SetLocation sets the save location.
func (m *ExportOptionsModal) SetLocation(location string) {
	m.data.Location = location
}

// Complete marks the modal as completed and hides it.
func (m *ExportOptionsModal) Complete() {
	m.completed = true
	m.visible = false
}

// GetDimensions returns the current modal dimensions.
func (m *ExportOptionsModal) GetDimensions() (width, height int) {
	return m.width, m.height
}
