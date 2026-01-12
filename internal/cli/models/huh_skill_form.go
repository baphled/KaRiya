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

// HuhSkillForm wraps a huh form for managing skills.
type HuhSkillForm struct {
	*BaseStandardModel
	formData *forms.SkillFormData
	form     *huh.Form
	width    int
	height   int
}

// NewHuhSkillForm creates a new huh-based skill form.
func NewHuhSkillForm() *HuhSkillForm {
	m := &HuhSkillForm{
		BaseStandardModel: NewBaseStandardModel(),
		formData:          &forms.SkillFormData{},
		width:             80,
		height:            24,
	}

	m.rebuildForm()
	return m
}

// NewHuhSkillFormWithData creates a new huh-based skill form with existing data.
func NewHuhSkillFormWithData(skill *career.Skill) *HuhSkillForm {
	m := &HuhSkillForm{
		BaseStandardModel: NewBaseStandardModel(),
		formData:          forms.GetSkillFormData(skill),
		width:             80,
		height:            24,
	}

	m.rebuildForm()
	return m
}

// rebuildForm creates a new form with current settings.
func (m *HuhSkillForm) rebuildForm() {
	m.form = forms.NewSkillFormWithDataAndDimensions(
		m.formData,
		m.width-4,
		forms.DefaultFormHeight(m.height),
	)
}

// Init initializes the form.
func (m *HuhSkillForm) Init() tea.Cmd {
	return m.form.Init()
}

// Update handles messages.
func (m *HuhSkillForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (m *HuhSkillForm) View() string {
	return m.form.View()
}

// submitForm creates a SkillFormCompleteMsg.
func (m *HuhSkillForm) submitForm() tea.Cmd {
	return func() tea.Msg {
		return SkillFormCompleteMsg{
			Data:      m.formData,
			Cancelled: false,
		}
	}
}

// GetFormData returns the current form data.
func (m *HuhSkillForm) GetFormData() *forms.SkillFormData {
	return m.formData
}

// SetFormData sets the form data and rebuilds the form.
func (m *HuhSkillForm) SetFormData(data *forms.SkillFormData) {
	m.formData = data
	m.rebuildForm()
}
