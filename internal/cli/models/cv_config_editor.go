package models

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
)

// CVConfigEditorModel handles editing of CV configurations.
type CVConfigEditorModel struct {
	*BaseStandardModel
	configManager cv.ConfigManager
	config        *career.CVConfig
	isNew         bool
	formFields    []*CVConfigField
	focusIndex    int
	saving        bool
	formErrors    map[string]string
}

// CVConfigField represents a single form field in the CV config editor.
type CVConfigField struct {
	Name     string
	Label    string
	Value    string
	FieldErr string
	Required bool
}

// NewCVConfigEditorModel creates a new CV Config Editor model.
func NewCVConfigEditorModel(
	baseModel *BaseStandardModel,
	configManager cv.ConfigManager,
	config *career.CVConfig,
) *CVConfigEditorModel {
	model := &CVConfigEditorModel{
		BaseStandardModel: baseModel,
		configManager:     configManager,
		config:            config,
		focusIndex:        0,
		saving:            false,
		formErrors:        make(map[string]string),
	}

	if config == nil {
		model.config = &career.CVConfig{
			EventFilters:   make(map[string]interface{}),
			TargetAudience: make([]string, 0),
		}
		model.isNew = true
	} else {
		model.isNew = false
	}

	model.initializeFormFields()
	return model
}

// initializeFormFields initializes the form fields based on the current config.
func (m *CVConfigEditorModel) initializeFormFields() {
	m.formFields = []*CVConfigField{
		{
			Name:     "cv_name",
			Label:    "CV Configuration Name",
			Value:    m.config.Name,
			Required: true,
		},
		{
			Name:     "target_role",
			Label:    "Target Role (senior_ic, staff, em, principal)",
			Value:    m.config.TargetRole,
			Required: true,
		},
		{
			Name:     "target_audience",
			Label:    "Target Audience (hiring_manager, recruiter, peer)",
			Value:    audienceListToString(m.config.TargetAudience),
			Required: true,
		},
	}
}

// audienceListToString converts audience list to comma-separated string.
func audienceListToString(audiences []string) string {
	if len(audiences) == 0 {
		return ""
	}
	result := ""
	for i, a := range audiences {
		if i > 0 {
			result += ", "
		}
		result += a
	}
	return result
}

// audienceStringToList converts comma-separated audience string to list.
func audienceStringToList(s string) []string {
	if s == "" {
		return []string{}
	}
	// Simple split by comma
	audiences := make([]string, 0)
	for _, a := range stringToList(s) {
		audiences = append(audiences, a)
	}
	return audiences
}

// stringToList is a helper to parse comma-separated string.
func stringToList(s string) []string {
	if s == "" {
		return []string{}
	}
	// Simplified - in production would use proper parsing
	var result []string
	var current string
	for _, ch := range s {
		if ch == ',' {
			if current != "" {
				result = append(result, current)
			}
			current = ""
		} else if ch != ' ' {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// Init initializes the editor.
func (m *CVConfigEditorModel) Init() tea.Cmd {
	return m.BaseStandardModel.Init()
}

// Update handles messages and updates the model state.
func (m *CVConfigEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.saving {
			return m, nil
		}

		switch msg.String() {
		case "tab":
			m.focusIndex = (m.focusIndex + 1) % len(m.formFields)
			return m, nil

		case "shift+tab":
			m.focusIndex = (m.focusIndex - 1 + len(m.formFields)) % len(m.formFields)
			return m, nil

		case "up":
			if m.focusIndex > 0 {
				m.focusIndex--
			}
			return m, nil

		case "down":
			if m.focusIndex < len(m.formFields)-1 {
				m.focusIndex++
			}
			return m, nil

		case "enter":
			if m.focusIndex == len(m.formFields)-1 {
				// Save button
				return m, m.saveConfig()
			}
			return m, nil

		case "esc":
			return m, func() tea.Msg {
				return BackToCVConfigManagerMsg{}
			}
		}

		// Handle text input for the focused field
		if m.focusIndex < len(m.formFields) {
			field := m.formFields[m.focusIndex]
			switch msg.Type {
			case tea.KeyRunes:
				field.Value += string(msg.Runes)
				m.formErrors[field.Name] = "" // Clear error on edit
				return m, nil
			case tea.KeyBackspace:
				if len(field.Value) > 0 {
					field.Value = field.Value[:len(field.Value)-1]
				}
				return m, nil
			}
		}
	}

	// Return model unchanged for unhandled messages
	return m, nil
}

// saveConfig validates and saves the configuration.
func (m *CVConfigEditorModel) saveConfig() tea.Cmd {
	return func() tea.Msg {
		// Validate fields
		m.formErrors = make(map[string]string)

		for _, field := range m.formFields {
			if field.Required && field.Value == "" {
				m.formErrors[field.Name] = fmt.Sprintf("%s is required", field.Label)
			}
		}

		if len(m.formErrors) > 0 {
			return ConfigValidationErrorMsg{errors: m.formErrors}
		}

		// Update config from form fields
		if len(m.formFields) > 0 {
			m.config.Name = m.formFields[0].Value
		}
		if len(m.formFields) > 1 {
			m.config.TargetRole = m.formFields[1].Value
		}
		if len(m.formFields) > 2 {
			m.config.TargetAudience = audienceStringToList(m.formFields[2].Value)
		}

		// Validate the entire config
		if err := m.config.Validate(); err != nil {
			return ConfigSaveErrorMsg{err: err}
		}

		// Save the config
		err := m.configManager.SaveConfig(context.Background(), m.config)
		if err != nil {
			return ConfigSaveErrorMsg{err: err}
		}

		return ConfigSavedMsg{config: m.config}
	}
}

// View renders the CV Configuration Editor screen.
func (m *CVConfigEditorModel) View() string {
	title := "New CV Configuration"
	if !m.isNew {
		title = "Edit CV Configuration"
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("6")).
		Render(title)

	var content string
	for i, field := range m.formFields {
		focused := i == m.focusIndex
		content += m.renderFormField(field, focused)
	}

	// Save button
	saveBtn := lipgloss.NewStyle().
		Bold(true).
		Render("[ Save ]")
	if len(m.formFields) == m.focusIndex {
		saveBtn = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("4")).
			Foreground(lipgloss.Color("15")).
			Render("[ Save ]")
	}
	content += "\n\n" + saveBtn

	// Error messages
	var errors string
	if len(m.formErrors) > 0 {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("1"))
		errors = "\n\nErrors:\n"
		for field, err := range m.formErrors {
			errors += errorStyle.Render(fmt.Sprintf("  %s: %s\n", field, err))
		}
	}

	footer := "\n\nTab: Next field | Shift+Tab: Previous field | Esc: Cancel"

	return fmt.Sprintf("%s\n\n%s%s%s", header, content, errors, footer)
}

// renderFormField renders a single form field.
func (m *CVConfigEditorModel) renderFormField(field *CVConfigField, focused bool) string {
	label := field.Label
	if field.Required {
		label += " *"
	}

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("6"))

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(0, 1)

	if focused {
		inputStyle = inputStyle.
			BorderForeground(lipgloss.Color("4")).
			Bold(true)
	}

	return fmt.Sprintf("%s\n%s\n\n",
		labelStyle.Render(label),
		inputStyle.Render(field.Value))
}

// Messages for CV Config Editor

// ConfigSavedMsg is sent when a configuration has been saved successfully.
type ConfigSavedMsg struct {
	config *career.CVConfig
}

// ConfigSaveErrorMsg is sent when there's an error saving a configuration.
type ConfigSaveErrorMsg struct {
	err error
}

// ConfigValidationErrorMsg is sent when there are validation errors.
type ConfigValidationErrorMsg struct {
	errors map[string]string
}

// BackToCVConfigManagerMsg navigates back to the config manager.
type BackToCVConfigManagerMsg struct{}

