package components

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

// SkillSortConfig represents the current sort configuration for skills.
type SkillSortConfig struct {
	SortBy    string
	SortOrder string
}

// SkillSortFormData holds the form field values for skill sorting.
type SkillSortFormData struct {
	SortBy    string
	SortOrder string
}

// SkillSortModal manages the sort modal form for skill sorting.
type SkillSortModal struct {
	form     *huh.Form
	formData *SkillSortFormData
	visible  bool
	width    int
	height   int
	theme    themes.Theme
}

// NewSkillSortModal creates a new skill sort modal.
func NewSkillSortModal(skills []*career.Skill, current *SkillSortConfig, width, height int) *SkillSortModal {
	formData := &SkillSortFormData{
		SortBy:    "name",
		SortOrder: "asc",
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

	modal := &SkillSortModal{
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
func (m *SkillSortModal) buildForm() {
	// Create form fields
	fields := []huh.Field{
		huh.NewSelect[string]().
			Title("Sort By").
			Options(
				huh.NewOption("Name", "name"),
				huh.NewOption("Category", "category"),
				huh.NewOption("Level", "level"),
				huh.NewOption("Years of Experience", "years"),
				huh.NewOption("Events Count", "events"),
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
func (m *SkillSortModal) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the sort modal
func (m *SkillSortModal) Update(msg tea.Msg) (tea.Cmd, bool, *SkillSortFormData) {
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

	// Update form.
	form, cmd := m.form.Update(msg)
	//nolint:errcheck // Type assertion is safe - form.Update always returns *huh.Form.
	m.form = form.(*huh.Form)

	// Check if form is complete.
	if m.form.State == huh.StateCompleted {
		m.visible = false
		return cmd, true, m.formData
	}

	return cmd, false, nil
}

// View renders the sort modal with proper chrome (border, background)
func (m *SkillSortModal) View() string {
	if !m.visible {
		return ""
	}

	// Build footer with primitives showing keyboard shortcuts
	footer := primitives.RenderHelpFooter(m.theme,
		primitives.NextFieldBadge(m.theme),
		primitives.ApplyBadge(m.theme),
		primitives.CancelBadge(m.theme),
	)

	// Build modal content with form and footer
	var content strings.Builder
	content.WriteString(m.form.View())
	content.WriteString("\n\n")
	content.WriteString(footer)

	// Wrap the form in a styled box with solid background using UIKit
	theme := m.theme
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}
	return containers.NewBox(theme).
		Content(content.String()).
		Padding(2).
		Background(theme.BackgroundColor()).
		Render()
}

// IsVisible returns whether the modal is currently visible.
func (m *SkillSortModal) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible.
func (m *SkillSortModal) Show() {
	m.visible = true
}

// Hide hides the modal.
func (m *SkillSortModal) Hide() {
	m.visible = false
}

// ToSkillSortConfig converts form data to SkillSortConfig
func (m *SkillSortModal) ToSkillSortConfig() *SkillSortConfig {
	return &SkillSortConfig{
		SortBy:    m.formData.SortBy,
		SortOrder: m.formData.SortOrder,
	}
}

// RenderOverlay renders the sort modal as an overlay on top of the base view.
// This follows the bubbletea-overlay pattern used in BrowseTimeline modals.
func (m *SkillSortModal) RenderOverlay(baseView string) string {
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
