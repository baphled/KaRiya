package components

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

// EventSortConfig represents the current sort configuration for events.
type EventSortConfig struct {
	SortBy    string
	SortOrder string
}

// EventSortFormData holds the form field values for event sorting.
type EventSortFormData struct {
	SortBy    string
	SortOrder string
}

// EventSortModal manages the sort modal form for event sorting.
type EventSortModal struct {
	form     *huh.Form
	formData *EventSortFormData
	visible  bool
	width    int
	height   int
	theme    themes.Theme
}

// NewEventSortModal creates a new event sort modal.
func NewEventSortModal(events []*career.CareerEvent, current *EventSortConfig, width, height int) *EventSortModal {
	formData := &EventSortFormData{
		SortBy:    "date",
		SortOrder: "desc",
	}

	// Pre-populate from current config
	if current != nil {
		if current.SortBy != "" {
			formData.SortBy = current.SortBy
		}
		if current.SortOrder != "" {
			formData.SortOrder = current.SortOrder
		}
	}

	modal := &EventSortModal{
		formData: formData,
		visible:  true,
		width:    width,
		height:   height,
		theme:    themes.NewDefaultTheme(),
	}

	modal.buildForm()
	return modal
}

// buildForm creates the huh form with sort options
func (m *EventSortModal) buildForm() {
	// Create form fields
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

	// Calculate modal width (60% of screen, max 60 chars for simpler form)
	modalWidth := m.width * 60 / 100
	if modalWidth > 60 {
		modalWidth = 60
	}
	if modalWidth < 40 {
		modalWidth = 40
	}

	// Let Huh use natural height - bubbletea-overlay will handle positioning
	m.form = huh.NewForm(group).
		WithWidth(modalWidth)
}

// Init initializes the sort modal and its form.
func (m *EventSortModal) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the sort modal
func (m *EventSortModal) Update(msg tea.Msg) (tea.Cmd, bool, *EventSortFormData) {
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
			// Close modal without applying
			m.visible = false
			return nil, false, nil
		}
	}

	// Update form
	form, cmd := m.form.Update(msg)
	m.form = form.(*huh.Form)

	// Check if form is complete
	if m.form.State == huh.StateCompleted {
		m.visible = false
		return cmd, true, m.formData
	}

	return cmd, false, nil
}

// View renders the sort modal with proper chrome (border, background)
func (m *EventSortModal) View() string {
	if !m.visible {
		return ""
	}

	// Build footer with KeyBadge components showing keyboard shortcuts
	footer := RenderHelpFooter(m.theme,
		NewKeyBadge("Tab", "Next field"),
		NewKeyBadge("Enter", "Apply"),
		NewKeyBadge("Esc", "Cancel"),
	)

	// Build modal content with form and footer
	var content strings.Builder
	content.WriteString(m.form.View())
	content.WriteString("\n\n")
	content.WriteString(footer)

	// Wrap the form in a styled box with solid background, border, and padding
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackground).
		Padding(1, 2).
		Render(content.String())
}

// IsVisible returns whether the modal is currently visible.
func (m *EventSortModal) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible.
func (m *EventSortModal) Show() {
	m.visible = true
}

// Hide hides the modal.
func (m *EventSortModal) Hide() {
	m.visible = false
}

// ToEventSortConfig converts form data to EventSortConfig
func (m *EventSortModal) ToEventSortConfig() *EventSortConfig {
	return &EventSortConfig{
		SortBy:    m.formData.SortBy,
		SortOrder: m.formData.SortOrder,
	}
}

// RenderOverlay renders the sort modal as an overlay on top of the base view.
// This follows the bubbletea-overlay pattern used in BrowseTimeline modals.
func (m *EventSortModal) RenderOverlay(baseView string) string {
	if !m.visible {
		return baseView
	}

	// Create static view model for modal content
	modalContent := staticViewModel{content: m.View()}
	bgModel := staticViewModel{content: baseView}

	// Use bubbletea-overlay to composite the modal on top of base view
	overlayModel := overlay.New(
		modalContent,   // Foreground: the sort form
		bgModel,        // Background: the rendered view
		overlay.Center, // X position
		overlay.Center, // Y position
		0,              // X offset
		-2,             // Y offset (avoid footer overlap)
	)

	return overlayModel.View()
}
