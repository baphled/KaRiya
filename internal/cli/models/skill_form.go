package models

import (
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// SkillFormCompleteMsg is sent when the skill form is completed
type SkillFormCompleteMsg struct {
	Data      *forms.SkillFormData
	Cancelled bool
}

// SkillForm wraps a huh form for managing skills.
type SkillForm struct {
	*BaseStandardModel
	formData *forms.SkillFormData
	form     *huh.Form
	width    int
	height   int
}

// NewSkillForm creates a new skill form.
func NewSkillForm() *SkillForm {
	m := &SkillForm{
		BaseStandardModel: NewBaseStandardModel(),
		formData:          &forms.SkillFormData{},
		width:             80,
		height:            24,
	}

	m.rebuildForm()
	return m
}

// NewSkillFormWithData creates a new skill form with existing data.
func NewSkillFormWithData(skill *career.Skill) *SkillForm {
	m := &SkillForm{
		BaseStandardModel: NewBaseStandardModel(),
		formData:          forms.GetSkillFormData(skill),
		width:             80,
		height:            24,
	}

	m.rebuildForm()
	return m
}

// rebuildForm creates a new form with current settings.
func (m *SkillForm) rebuildForm() {
	m.form = forms.NewSkillFormWithDataAndDimensions(
		m.formData,
		m.width-4,
		forms.DefaultFormHeight(m.height),
	)
}

// Init initializes the form.
func (m *SkillForm) Init() tea.Cmd {
	return m.form.Init()
}

// Update handles messages.
func (m *SkillForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.form = m.form.WithHeight(forms.DefaultFormHeight(m.height)).WithWidth(m.width - 4)
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+s" {
			m.formData.SubmitConfirmed = true
			return m, m.submitForm()
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted && m.formData.SubmitConfirmed {
		return m, m.submitForm()
	}

	if m.form.State == huh.StateAborted {
		return m, func() tea.Msg {
			return SkillFormCompleteMsg{Data: nil, Cancelled: true}
		}
	}

	return m, cmd
}

// View renders the form.
func (m *SkillForm) View() string {
	return m.form.View()
}

// submitForm creates a SkillFormCompleteMsg.
func (m *SkillForm) submitForm() tea.Cmd {
	return func() tea.Msg {
		return SkillFormCompleteMsg{
			Data:      m.formData,
			Cancelled: false,
		}
	}
}

// GetFormData returns the current form data.
func (m *SkillForm) GetFormData() *forms.SkillFormData {
	return m.formData
}

// SetFormData sets the form data and rebuilds the form.
func (m *SkillForm) SetFormData(data *forms.SkillFormData) {
	m.formData = data
	m.rebuildForm()
}
