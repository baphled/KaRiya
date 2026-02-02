package modals

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
)

// SearchFormData holds the form field values for skill searching.
type SearchFormData struct {
	SearchText string
}

// SearchModal manages the search modal form for skill searching.
type SearchModal struct {
	form     forms.Form
	formData *SearchFormData
	visible  bool
	width    int
	height   int
	theme    themes.Theme
}

// NewSearchModal creates a new skill search modal.
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
	// Create form fields using forms package.
	fields := []forms.Field{
		forms.NewInput(forms.FieldConfig{
			Title:       "Search Skills",
			Description: "Search by name, category, or description",
			Placeholder: "Enter search text...",
		}).Value(&m.formData.SearchText),
	}

	group := forms.NewGroup(fields...)

	// Calculate modal width (60% of screen, max 60 chars for simple search form).
	modalWidth := m.width * 60 / 100
	if modalWidth > 60 {
		modalWidth = 60
	}
	if modalWidth < 40 {
		modalWidth = 40
	}

	// Create form with dimensions.
	m.form = forms.NewFormWithDimensions(modalWidth, 0, group)
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
	if wsm, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = wsm.Width
		m.height = wsm.Height
		m.rebuildForm()
		return m.form.Init(), false, nil
	}

	// Handle KeyMsg.
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyEsc {
			// User cancelled - close modal without applying.
			m.visible = false
			return nil, false, nil
		}
	}

	// Forward ALL messages to form (not just KeyMsg).
	// This is CRITICAL for Tab/Enter to work in huh forms.
	var cmd tea.Cmd
	m.form, cmd = forms.Update(m.form, msg)

	// Check if form just completed.
	if forms.IsCompleted(m.form) {
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

// SetSize updates the modal's dimensions and rebuilds the form.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (m *SearchModal) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.rebuildForm()
}

// GetSearchText returns the current search text.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *SearchModal) GetSearchText() string {
	return m.formData.SearchText
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
