package modals

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// SearchFormData holds the form field values for event searching.
type SearchFormData struct {
	SearchText string
}

// SearchModal manages the search modal form for event searching.
type SearchModal struct {
	form     *huh.Form
	formData *SearchFormData
	visible  bool
	width    int
	height   int
	theme    themes.Theme
}

// NewSearchModal creates a new event search modal.
//
// Expected:
//   - Must be a valid string.
//   - int must be valid.
//
// Returns:
//   - A fully initialized SearchModal ready for use.
//
// Side effects:
//   - None.
func NewSearchModal(currentSearch string, width, height int) *SearchModal {
	formData := &SearchFormData{
		SearchText: currentSearch,
	}

	m := &SearchModal{
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
func (m *SearchModal) rebuildForm() {
	// Create form fields.
	fields := []huh.Field{
		huh.NewInput().
			Title("Search Events").
			Description("Search by text, company, category, or tags").
			Placeholder("Enter search text...").
			Value(&m.formData.SearchText),
	}

	group := huh.NewGroup(fields...)

	// Calculate modal width (60% of screen, max 60 chars for simple search form).
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

// Init initializes the modal.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *SearchModal) Init() tea.Cmd {
	m.visible = true
	return m.form.Init()
}

// Update handles messages for the search modal.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Cmd: command to execute.
//   - bool: true if form completed successfully.
//   - *SearchFormData: search data if completed, nil otherwise.
//
// Side effects:
//   - May hide modal on completion or cancellation.
func (m *SearchModal) Update(msg tea.Msg) (tea.Cmd, bool, *SearchFormData) {
	// Handle WindowSizeMsg for responsive sizing.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.rebuildForm()
		return m.form.Init(), false, nil
	}

	// Handle KeyMsg.
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "esc" {
			// User cancelled - close modal without applying.
			m.visible = false
			return nil, false, nil
		}
	}

	// Forward ALL messages to form (not just KeyMsg).
	// This is CRITICAL for Tab/Enter to work in huh forms.
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	// Check if form just completed.
	if m.form.State == huh.StateCompleted {
		m.visible = false
		return cmd, true, m.formData
	}

	return cmd, false, nil
}

// View renders the search modal.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *SearchModal) View() string {
	if !m.visible {
		return ""
	}

	// Build footer with primitives showing keyboard shortcuts.
	footer := primitives.RenderHelpFooter(m.theme,
		primitives.NextFieldBadge(m.theme),
		primitives.SubmitBadge(m.theme),
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
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *SearchModal) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *SearchModal) Show() {
	m.visible = true
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *SearchModal) Hide() {
	m.visible = false
}

// RenderOverlay renders the search modal as an overlay on top of the base view.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *SearchModal) RenderOverlay(baseView string) string {
	if !m.visible {
		return baseView
	}

	return RenderOverlayModal(m.View(), baseView)
}
