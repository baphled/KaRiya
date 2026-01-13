package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

// staticViewModel implements overlay.Viewable for rendering static content.
type staticViewModel struct {
	content string
}

func (m staticViewModel) View() string { return m.content }

// SkillSearchFormData holds the form field values for skill searching.
type SkillSearchFormData struct {
	SearchText string
}

// SkillSearchModal manages the search modal form for skill searching.
type SkillSearchModal struct {
	form     *huh.Form
	formData *SkillSearchFormData
	visible  bool
	width    int
	height   int
}

// NewSkillSearchModal creates a new skill search modal.
func NewSkillSearchModal(currentSearch string, width, height int) *SkillSearchModal {
	formData := &SkillSearchFormData{
		SearchText: currentSearch,
	}

	m := &SkillSearchModal{
		formData: formData,
		visible:  false,
		width:    width,
		height:   height,
	}

	m.rebuildForm()
	return m
}

// rebuildForm creates a new form instance with current dimensions.
func (m *SkillSearchModal) rebuildForm() {
	// Create form fields
	fields := []huh.Field{
		huh.NewInput().
			Title("Search Skills").
			Description("Search by name, category, or description").
			Placeholder("Enter search text...").
			Value(&m.formData.SearchText),
	}

	// Create form
	m.form = huh.NewForm(
		huh.NewGroup(fields...),
	).WithWidth(m.width - 4)
}

// Init initializes the modal.
func (m *SkillSearchModal) Init() tea.Cmd {
	m.visible = true
	return m.form.Init()
}

// Update handles messages for the search modal.
// Returns (cmd, applied, searchData).
// - cmd: Command to execute
// - applied: true if user confirmed search (Enter), false if cancelled (Esc)
// - searchData: The search text if applied=true, nil otherwise
func (m *SkillSearchModal) Update(msg tea.KeyMsg) (tea.Cmd, bool, *SkillSearchFormData) {
	switch msg.String() {
	case "esc":
		// User cancelled - close modal without applying
		m.visible = false
		return nil, false, nil

	case "enter":
		// Check if form is complete
		if m.form.State == huh.StateCompleted {
			// User confirmed - close modal and apply search
			m.visible = false
			return nil, true, m.formData
		}
		// Form not complete yet, let it handle Enter
		form, cmd := m.form.Update(msg)
		m.form = form.(*huh.Form)

		// Check again if form just completed
		if m.form.State == huh.StateCompleted {
			m.visible = false
			return cmd, true, m.formData
		}
		return cmd, false, nil

	default:
		// Forward all other keys to form
		form, cmd := m.form.Update(msg)
		m.form = form.(*huh.Form)
		return cmd, false, nil
	}
}

// View renders the search modal.
func (m *SkillSearchModal) View() string {
	if !m.visible {
		return ""
	}

	return m.form.View()
}

// IsVisible returns whether the modal is currently visible.
func (m *SkillSearchModal) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible.
func (m *SkillSearchModal) Show() {
	m.visible = true
}

// Hide hides the modal.
func (m *SkillSearchModal) Hide() {
	m.visible = false
}

// SetSize updates the modal's dimensions and rebuilds the form.
func (m *SkillSearchModal) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.rebuildForm()
}

// GetSearchText returns the current search text.
func (m *SkillSearchModal) GetSearchText() string {
	return m.formData.SearchText
}

// RenderOverlay renders the search modal as an overlay on top of the base view.
// This follows the bubbletea-overlay pattern used in BrowseTimeline modals.
func (m *SkillSearchModal) RenderOverlay(baseView string) string {
	if !m.visible {
		return baseView
	}

	// Create static view model for modal content
	modalContent := staticViewModel{content: m.View()}
	bgModel := staticViewModel{content: baseView}

	// Use bubbletea-overlay to composite the modal on top of base view
	overlayModel := overlay.New(
		modalContent,   // Foreground: the search form
		bgModel,        // Background: the rendered view
		overlay.Center, // X position
		overlay.Center, // Y position
		0,              // X offset
		0,              // Y offset
	)

	return overlayModel.View()
}
