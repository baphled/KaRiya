package skill

import (
	"strings"

	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/containers"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
)

// AddEdit provides a way to add or edit a skill.
// It shows a form with fields for name, category, level, and years of experience.
//
// Usage:
//
//	// For adding a new skill:
//	modal := components.NewAddEdit(nil, width, height)
//
//	// For editing an existing skill:
//	modal := components.NewAddEdit(existingSkill, width, height)
//
//	cmd := modal.Init()
//	// In Update:
//	cmd, completed, skillData := modal.Update(msg)
//	if completed && skillData != nil {
//	    // User completed form - use EditData fields directly
//	    // Pass to domain function: skills.NewSkillFromInput(input, skillID)
//	} else if !modal.IsVisible() {
//	    // User cancelled (Esc)
//	}
type AddEdit struct {
	form            forms.Form
	formData        *forms.SkillFormData
	originalSkillID string
	visible         bool
	width           int
	height          int
}

// NewAddEdit creates a new skill add/edit modal with the given
//
// Expected:
//   - skill must be valid.
//   - int must be valid.
//
// Returns:
//   - A fully initialized AddEdit ready for use.
//
// Side effects:
//   - None.
func NewAddEdit(formData *forms.SkillFormData, width, height int, originalSkillID ...string) *AddEdit {
	if formData == nil {
		formData = &forms.SkillFormData{}
	}

	skillID := ""
	if len(originalSkillID) > 0 {
		skillID = originalSkillID[0]
	}

	modal := &AddEdit{
		formData:        formData,
		originalSkillID: skillID,
		visible:         true,
		width:           width,
		height:          height,
	}

	modal.buildForm()
	return modal
}

// buildForm creates the huh form with proper dimensions.
func (m *AddEdit) buildForm() {
	// Calculate form width
	modalWidth := m.width - 10
	if modalWidth > 90 {
		modalWidth = 90
	}
	if modalWidth < 50 {
		modalWidth = 50
	}

	formWidth := forms.ModalFormWidth(modalWidth)
	formHeight := forms.ModalFormHeight(m.height)

	m.form = forms.NewSkillForm(m.formData, formWidth, formHeight)
}

// Init initializes the modal and its form.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *AddEdit) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the skill add/edit modal.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Cmd: command to execute.
//   - bool: true if form completed successfully.
//   - *EditData: skill data if completed, nil otherwise.
//
// Side effects:
//   - May hide modal on completion or cancellation.
//   - May rebuild form on window resize.
func (m *AddEdit) Update(msg tea.Msg) (tea.Cmd, bool, *EditData) {
	if !m.visible {
		return nil, false, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Rebuild form with new dimensions
		m.buildForm()
		return m.form.Init(), false, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			// Close modal without saving
			m.visible = false
			return nil, false, nil
		}
	}

	// Update form using forms package helper.
	var cmd tea.Cmd
	m.form, cmd = forms.Update(m.form, msg)

	// Check if form is complete AND user confirmed submission.
	if forms.IsCompleted(m.form) {
		m.visible = false
		// Only return data if user confirmed (pressed Submit, not Cancel)
		if m.formData.SubmitConfirmed {
			skillData := &EditData{
				Name:      m.formData.Name,
				Category:  m.formData.Category,
				Level:     m.formData.Level,
				YearsUsed: m.formData.YearsUsed,
			}
			return cmd, true, skillData
		}
		// User cancelled - close without returning data
		return cmd, false, nil
	}

	return cmd, false, nil
}

// View renders the skill add/edit modal with proper chrome (border, background)
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *AddEdit) View() string {
	if !m.visible {
		return ""
	}

	theme := themes.NewDefaultTheme()

	// Build footer with primitives showing keyboard shortcuts.
	footer := primitives.RenderHelpFooter(theme,
		primitives.NextFieldBadge(theme),
		primitives.SubmitBadge(theme),
		primitives.CancelBadge(theme),
	)

	// Build modal content with form and footer.
	var content strings.Builder
	content.WriteString(m.form.View())
	content.WriteString("\n\n")
	content.WriteString(footer)

	// Wrap the form in a styled box with solid background using UIKit.
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
func (m *AddEdit) IsVisible() bool {
	return m.visible
}

// IsEditMode returns true if editing an existing skill, false if adding new.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *AddEdit) IsEditMode() bool {
	return m.originalSkillID != ""
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *AddEdit) Show() {
	m.visible = true
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *AddEdit) Hide() {
	m.visible = false
}

// GetOriginalSkillID returns the original skill ID being edited.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *AddEdit) GetOriginalSkillID() string {
	return m.originalSkillID
}

// EditData holds the form data from adding or editing a skill.
type EditData struct {
	Name      string
	Category  string
	Level     string
	YearsUsed string
}
