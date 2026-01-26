package modals

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// SortConfig represents the current sort configuration for events.
type SortConfig struct {
	SortBy    string
	SortOrder string
}

// SortFormData holds the form field values for event sorting.
type SortFormData struct {
	SortBy    string
	SortOrder string
}

// SortModal manages the sort modal form for event sorting.
type SortModal struct {
	form     *huh.Form
	formData *SortFormData
	visible  bool
	width    int
	height   int
	theme    themes.Theme
}

// NewSortModal creates a new event sort modal.
// The events parameter is reserved for future use (e.g., dynamic sort options based on data).
func NewSortModal(_ []*career.CareerEvent, current *SortConfig, width, height int) *SortModal {
	formData := &SortFormData{
		SortBy:    "date",
		SortOrder: "desc",
	}

	// Pre-populate from current config.
	if current != nil {
		if current.SortBy != "" {
			formData.SortBy = current.SortBy
		}
		if current.SortOrder != "" {
			formData.SortOrder = current.SortOrder
		}
	}

	modal := &SortModal{
		formData: formData,
		visible:  true,
		width:    width,
		height:   height,
		theme:    themes.NewDefaultTheme(),
	}

	modal.buildForm()
	return modal
}

// buildForm creates the huh form with sort options.
func (m *SortModal) buildForm() {
	// Create form fields.
	fields := []huh.Field{
		huh.NewSelect[string]().
			Title("Sort By").
			Options(
				huh.NewOption("Date", "date"),
				huh.NewOption("Company", "company"),
				huh.NewOption("Category", "category"),
			).
			Value(&m.formData.SortBy),

		huh.NewSelect[string]().
			Title("Sort Order").
			Options(
				huh.NewOption("Ascending", "asc"),
				huh.NewOption("Descending", "desc"),
			).
			Value(&m.formData.SortOrder),
	}

	group := huh.NewGroup(fields...)

	// Calculate modal width (60% of screen, max 60 chars for simpler form).
	modalWidth := m.width * 60 / 100
	if modalWidth > 60 {
		modalWidth = 60
	}
	if modalWidth < 40 {
		modalWidth = 40
	}

	m.form = huh.NewForm(group).
		WithWidth(modalWidth)
}

// Init initializes the sort modal and its form.
func (m *SortModal) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the sort modal.
func (m *SortModal) Update(msg tea.Msg) (tea.Cmd, bool, *SortFormData) {
	if !m.visible {
		return nil, false, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return nil, false, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Close modal without applying.
			m.visible = false
			return nil, false, nil
		}
	}

	// Update form.
	form, cmd := m.form.Update(msg)
	m.form = form.(*huh.Form)

	// Check if form is complete.
	if m.form.State == huh.StateCompleted {
		m.visible = false
		return cmd, true, m.formData
	}

	return cmd, false, nil
}

// View renders the sort modal with proper chrome (border, background).
func (m *SortModal) View() string {
	if !m.visible {
		return ""
	}

	// Build footer with primitives showing keyboard shortcuts.
	footer := primitives.RenderHelpFooter(m.theme,
		primitives.NextFieldBadge(m.theme),
		primitives.ApplyBadge(m.theme),
		primitives.CancelBadge(m.theme),
	)

	// Build modal content with form and footer.
	var content strings.Builder
	content.WriteString(m.form.View())
	content.WriteString("\n\n")
	content.WriteString(footer)

	// Wrap the form in a styled box with solid background using UIKit.
	return containers.NewBox(m.theme).
		Content(content.String()).
		Padding(2).
		Background(m.theme.BackgroundColor()).
		Render()
}

// IsVisible returns whether the modal is currently visible.
func (m *SortModal) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible.
func (m *SortModal) Show() {
	m.visible = true
}

// Hide hides the modal.
func (m *SortModal) Hide() {
	m.visible = false
}

// ToSortConfig converts form data to SortConfig.
func (m *SortModal) ToSortConfig() *SortConfig {
	return &SortConfig{
		SortBy:    m.formData.SortBy,
		SortOrder: m.formData.SortOrder,
	}
}

// RenderOverlay renders the sort modal as an overlay on top of the base view.
// This follows the bubbletea-overlay pattern used in BrowseTimeline modals.
func (m *SortModal) RenderOverlay(baseView string) string {
	if !m.visible {
		return baseView
	}

	return RenderOverlayModal(m.View(), baseView)
}
