package components

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
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
	theme    themes.Theme
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
		theme:    themes.NewDefaultTheme(),
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

	group := huh.NewGroup(fields...)

	// Calculate modal width (60% of screen, max 60 chars for simple search form)
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
//
// CRITICAL: Takes tea.Msg (not tea.KeyMsg) to allow huh forms to process
// Tab and Enter keys correctly. Huh forms require full tea.Msg interface.
func (m *SkillSearchModal) Update(msg tea.Msg) (tea.Cmd, bool, *SkillSearchFormData) {
	// Handle WindowSizeMsg for responsive sizing
	if wsm, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = wsm.Width
		m.height = wsm.Height
		m.rebuildForm()
		return m.form.Init(), false, nil
	}

	// Handle KeyMsg
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "esc" {
			// User cancelled - close modal without applying
			m.visible = false
			return nil, false, nil
		}
	}

	// Forward ALL messages to form (not just KeyMsg).
	// This is CRITICAL for Tab/Enter to work in huh forms.
	form, cmd := m.form.Update(msg)
	//nolint:errcheck // Type assertion is safe - form.Update always returns *huh.Form.
	m.form = form.(*huh.Form)

	// Check if form just completed.
	if m.form.State == huh.StateCompleted {
		m.visible = false
		return cmd, true, m.formData
	}

	return cmd, false, nil
}

// View renders the search modal.
func (m *SkillSearchModal) View() string {
	if !m.visible {
		return ""
	}

	// Build footer with primitives showing keyboard shortcuts
	footer := primitives.RenderHelpFooter(m.theme,
		primitives.NextFieldBadge(m.theme),
		primitives.SubmitBadge(m.theme),
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
