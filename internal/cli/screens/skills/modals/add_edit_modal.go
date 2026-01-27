package modals

import (
	"strconv"
	"strings"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// AddEditModal provides a way to add or edit a skill.
// It shows a form with fields for name, category, level, and years of experience.
//
// Usage:
//
//	// For adding a new skill:
//	modal := components.NewAddEditModal(nil, width, height)
//
//	// For editing an existing skill:
//	modal := components.NewAddEditModal(existingSkill, width, height)
//
//	cmd := modal.Init()
//	// In Update:
//	cmd, completed, skillData := modal.Update(msg)
//	if completed && skillData != nil {
//	    // User completed form - save skill
//	    if modal.IsEditMode() {
//	        updatedSkill := skillData.ToSkill(existingSkill.ID)
//	    } else {
//	        newSkill := skillData.ToSkill("")
//	    }
//	} else if !modal.IsVisible() {
//	    // User cancelled (Esc)
//	}
type AddEditModal struct {
	form          forms.Form
	formData      *forms.SkillFormData
	originalSkill *career.Skill // nil for add mode
	visible       bool
	width         int
	height        int
}

// NewAddEditModal creates a new skill add/edit modal with the given
// terminal dimensions. If skill is nil, creates a form for adding a new skill.
// If skill is provided, creates a form for editing with pre-populated fields.
func NewAddEditModal(skill *career.Skill, width, height int) *AddEditModal {
	// Initialize form data from existing skill or empty
	formData := &forms.SkillFormData{}
	if skill != nil {
		formData = forms.GetSkillFormData(skill)
	}

	modal := &AddEditModal{
		formData:      formData,
		originalSkill: skill,
		visible:       true,
		width:         width,
		height:        height,
	}

	modal.buildForm()
	return modal
}

// buildForm creates the huh form with proper dimensions.
func (m *AddEditModal) buildForm() {
	// Calculate form width
	modalWidth := m.width - 10
	if modalWidth > 90 {
		modalWidth = 90
	}
	if modalWidth < 50 {
		modalWidth = 50
	}

	// Let Huh use natural height
	formHeight := 0

	// Create form with dimensions
	m.form = forms.NewSkillFormWithDataAndDimensions(m.formData, modalWidth, formHeight)
}

// Init initializes the modal and its form.
func (m *AddEditModal) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the skill add/edit modal.
//
// Returns:
//   - tea.Cmd: command to execute
//   - bool: true if form completed successfully
//   - *SkillEditData: skill data if completed, nil otherwise
func (m *AddEditModal) Update(msg tea.Msg) (tea.Cmd, bool, *SkillEditData) {
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
			skillData := &SkillEditData{
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
// for overlay compositing.
func (m *AddEditModal) View() string {
	if !m.visible {
		return ""
	}

	// Wrap the form in a styled box with solid background using UIKit
	theme := themes.NewDefaultTheme()
	return containers.NewBox(theme).
		Content(m.form.View()).
		Padding(2).
		Background(theme.BackgroundColor()).
		Render()
}

// IsVisible returns whether the modal is currently visible.
func (m *AddEditModal) IsVisible() bool {
	return m.visible
}

// IsEditMode returns true if editing an existing skill, false if adding new.
func (m *AddEditModal) IsEditMode() bool {
	return m.originalSkill != nil
}

// Show makes the modal visible.
func (m *AddEditModal) Show() {
	m.visible = true
}

// Hide hides the modal.
func (m *AddEditModal) Hide() {
	m.visible = false
}

// GetOriginalSkill returns the original skill being edited (nil for add mode).
func (m *AddEditModal) GetOriginalSkill() *career.Skill {
	return m.originalSkill
}

// SkillEditData holds the form data from adding or editing a skill.
type SkillEditData struct {
	Name      string
	Category  string
	Level     string
	YearsUsed string
}

// ToSkill converts the form data to a Skill domain object.
// skillID: the ID of the skill being updated (empty string for new skills)
// Returns a Skill ready to be saved.
func (d *SkillEditData) ToSkill(skillID string) *career.Skill {
	skill := &career.Skill{
		ID:       skillID,
		Name:     strings.TrimSpace(d.Name),
		Category: strings.TrimSpace(d.Category),
		Level:    strings.TrimSpace(d.Level),
	}

	// Parse years used
	yearsStr := strings.TrimSpace(d.YearsUsed)
	if yearsStr != "" {
		if years, err := strconv.Atoi(yearsStr); err == nil && years >= 0 && years <= 50 {
			skill.YearsUsed = &years
		}
	}

	return skill
}
