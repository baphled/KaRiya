package intents

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ConfigurationDomain represents a configuration domain (e.g., "system", "profile", "export")
type ConfigurationDomain string

const (
	DomainSystem  ConfigurationDomain = "system"
	DomainProfile ConfigurationDomain = "profile"
	DomainExport  ConfigurationDomain = "export"
	DomainUI      ConfigurationDomain = "ui"
)

// ConfigurationSetting represents a single configuration setting
type ConfigurationSetting struct {
	Key          string      // e.g., "theme"
	Label        string      // e.g., "Theme"
	Value        interface{} // Current value
	DefaultValue interface{} // Default value
	Type         string      // "string", "bool", "int", "select"
	Options      []string    // For "select" type
	Description  string      // Help text
}

// ConfigurationState represents the current state of the configuration
type ConfigurationState string

const (
	ConfigStateSelectDomain  ConfigurationState = "select_domain"
	ConfigStateEditSettings  ConfigurationState = "edit_settings"
	ConfigStateReviewChanges ConfigurationState = "review_changes"
	ConfigStateConfirm       ConfigurationState = "confirm"
	ConfigStateSaving        ConfigurationState = "saving"
	ConfigStateComplete      ConfigurationState = "complete"
	ConfigStateFailed        ConfigurationState = "failed"
)

// ConfigurationChanges tracks all changes made during editing
type ConfigurationChanges struct {
	Domain   ConfigurationDomain
	Original map[string]interface{} // Original values
	Modified map[string]interface{} // Modified values
}

// ConfigureSystemContext contains context for the ConfigureSystem intent
type ConfigureSystemContext struct {
	Domains  []ConfigurationDomain
	Settings map[ConfigurationDomain][]*ConfigurationSetting
}

// ConfigureSystemResult contains the result of configuration
type ConfigureSystemResult struct {
	Success bool
	Domain  ConfigurationDomain
	Changes *ConfigurationChanges
	Error   *IntentError
}

// ConfigureSystemModel represents the state of the ConfigureSystem intent
type ConfigureSystemModel struct {
	state         ConfigurationState
	context       *ConfigureSystemContext
	selectedIndex int
	domain        ConfigurationDomain
	changes       *ConfigurationChanges
	result        *ConfigureSystemResult
	error         *IntentError
	active        bool
}

// NewConfigureSystemContext creates a new ConfigureSystemContext with default values
func NewConfigureSystemContext() *ConfigureSystemContext {
	return &ConfigureSystemContext{
		Domains: []ConfigurationDomain{
			DomainSystem,
			DomainProfile,
			DomainExport,
			DomainUI,
		},
		Settings: map[ConfigurationDomain][]*ConfigurationSetting{
			DomainSystem: {
				{
					Key:          "log_level",
					Label:        "Log Level",
					Value:        "info",
					DefaultValue: "info",
					Type:         "select",
					Options:      []string{"debug", "info", "warn", "error"},
					Description:  "Logging verbosity level",
				},
				{
					Key:          "data_dir",
					Label:        "Data Directory",
					Value:        "/home/user/.kariya",
					DefaultValue: "/home/user/.kariya",
					Type:         "string",
					Description:  "Directory for storing application data",
				},
			},
			DomainProfile: {
				{
					Key:          "name",
					Label:        "Full Name",
					Value:        "John Doe",
					DefaultValue: "",
					Type:         "string",
					Description:  "Your full name",
				},
				{
					Key:          "email",
					Label:        "Email",
					Value:        "john@example.com",
					DefaultValue: "",
					Type:         "string",
					Description:  "Your email address",
				},
			},
			DomainExport: {
				{
					Key:          "default_format",
					Label:        "Default Export Format",
					Value:        "pdf",
					DefaultValue: "pdf",
					Type:         "select",
					Options:      []string{"pdf", "json", "csv", "txt"},
					Description:  "Default format for exports",
				},
				{
					Key:          "auto_open",
					Label:        "Auto-Open Exported Files",
					Value:        true,
					DefaultValue: false,
					Type:         "bool",
					Description:  "Automatically open exported files",
				},
			},
			DomainUI: {
				{
					Key:          "theme",
					Label:        "Theme",
					Value:        "dark",
					DefaultValue: "dark",
					Type:         "select",
					Options:      []string{"light", "dark"},
					Description:  "UI color theme",
				},
				{
					Key:          "animations",
					Label:        "Animations",
					Value:        true,
					DefaultValue: true,
					Type:         "bool",
					Description:  "Enable UI animations",
				},
			},
		},
	}
}

// NewConfigureSystemModel creates a new ConfigureSystem intent model
func NewConfigureSystemModel(ctx *ConfigureSystemContext) *ConfigureSystemModel {
	return &ConfigureSystemModel{
		state:         ConfigStateSelectDomain,
		context:       ctx,
		selectedIndex: 0,
		active:        false,
	}
}

// Init initializes the intent
func (m *ConfigureSystemModel) Init() tea.Cmd {
	m.active = true
	return nil
}

// Update handles messages
func (m *ConfigureSystemModel) Update(msg tea.Msg) tea.Cmd {
	if !m.active {
		return nil
	}

	switch m.state {
	case ConfigStateSelectDomain:
		return m.updateSelectDomain(msg)
	case ConfigStateEditSettings:
		return m.updateEditSettings(msg)
	case ConfigStateReviewChanges:
		return m.updateReviewChanges(msg)
	case ConfigStateConfirm:
		return m.updateConfirm(msg)
	case ConfigStateSaving:
		return m.updateSaving(msg)
	case ConfigStateComplete:
		return m.updateComplete(msg)
	case ConfigStateFailed:
		return m.updateFailed(msg)
	default:
		return nil
	}
}

// View renders the current state
func (m *ConfigureSystemModel) View() string {
	if !m.active {
		return ""
	}

	switch m.state {
	case ConfigStateSelectDomain:
		return m.viewSelectDomain()
	case ConfigStateEditSettings:
		return m.viewEditSettings()
	case ConfigStateReviewChanges:
		return m.viewReviewChanges()
	case ConfigStateConfirm:
		return m.viewConfirm()
	case ConfigStateSaving:
		return m.viewSaving()
	case ConfigStateComplete:
		return m.viewComplete()
	case ConfigStateFailed:
		return m.viewFailed()
	default:
		return "Unknown state"
	}
}

// Result returns the intent result
func (m *ConfigureSystemModel) Result() *IntentResult[interface{}] {
	if m.result == nil {
		return &IntentResult[interface{}]{
			Status: Cancelled,
		}
	}

	if m.result.Error != nil && m.result.Error.Code == "config_cancelled" {
		return &IntentResult[interface{}]{
			Status: Cancelled,
		}
	}

	if !m.result.Success {
		return &IntentResult[interface{}]{
			Status: Failed,
			Data:   m.result,
			Error:  m.result.Error,
		}
	}

	return &IntentResult[interface{}]{
		Status: Completed,
		Data:   m.result,
	}
}

// State-specific update handlers

func (m *ConfigureSystemModel) updateSelectDomain(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "down":
			if m.selectedIndex < len(m.context.Domains)-1 {
				m.selectedIndex++
			}
		case "enter":
			m.domain = m.context.Domains[m.selectedIndex]
			m.changes = &ConfigurationChanges{
				Domain:   m.domain,
				Original: make(map[string]interface{}),
				Modified: make(map[string]interface{}),
			}
			// Copy current values as original
			for _, setting := range m.context.Settings[m.domain] {
				m.changes.Original[setting.Key] = setting.Value
				m.changes.Modified[setting.Key] = setting.Value
			}
			m.selectedIndex = 0
			m.state = ConfigStateEditSettings
		case "esc":
			m.setResult(&ConfigureSystemResult{
				Success: false,
				Error: &IntentError{
					Code:    "config_cancelled",
					Message: "Configuration cancelled by user",
				},
			})
		}
	}
	return nil
}

func (m *ConfigureSystemModel) updateEditSettings(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.state = ConfigStateReviewChanges
		case "esc":
			m.selectedIndex = 0
			m.state = ConfigStateSelectDomain
		}
	}
	return nil
}

func (m *ConfigureSystemModel) updateReviewChanges(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.state = ConfigStateConfirm
		case "esc":
			m.state = ConfigStateEditSettings
		}
	}
	return nil
}

func (m *ConfigureSystemModel) updateConfirm(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "enter":
			m.state = ConfigStateSaving
			return m.startSave()
		case "n", "esc":
			m.state = ConfigStateReviewChanges
		}
	}
	return nil
}

func (m *ConfigureSystemModel) updateSaving(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case ConfigCompleteMsg:
		m.state = ConfigStateComplete
		m.result = msg.Result
	case ConfigErrorMsg:
		m.state = ConfigStateFailed
		m.error = msg.Error
	}
	return nil
}

func (m *ConfigureSystemModel) updateComplete(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" || msg.String() == "esc" {
			m.active = false
		}
	}
	return nil
}

func (m *ConfigureSystemModel) updateFailed(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			// Retry
			m.state = ConfigStateConfirm
		case "esc":
			m.active = false
		}
	}
	return nil
}

// startSave initiates the async save operation
func (m *ConfigureSystemModel) startSave() tea.Cmd {
	return func() tea.Msg {
		// Simulate save operation (would be replaced with actual service call)
		result := &ConfigureSystemResult{
			Success: true,
			Domain:  m.domain,
			Changes: m.changes,
		}
		return ConfigCompleteMsg{Result: result}
	}
}

// setResult sets the intent result and marks intent as complete
func (m *ConfigureSystemModel) setResult(result *ConfigureSystemResult) {
	m.result = result
	m.active = false
}

// View rendering methods

func (m *ConfigureSystemModel) viewSelectDomain() string {
	s := "Select Configuration Domain:\n\n"
	for i, d := range m.context.Domains {
		prefix := "  "
		if i == m.selectedIndex {
			prefix = "> "
		}
		s += prefix + string(d) + "\n"
	}
	s += "\n[↑/↓] Navigate | [Enter] Select | [Esc] Cancel"
	return s
}

func (m *ConfigureSystemModel) viewEditSettings() string {
	s := "Edit Settings for " + string(m.domain) + ":\n\n"
	settings := m.context.Settings[m.domain]
	for i, setting := range settings {
		prefix := "  "
		if i == m.selectedIndex {
			prefix = "> "
		}
		s += prefix + setting.Label + ": " + fmt.Sprintf("%v", setting.Value) + "\n"
	}
	s += "\n[↑/↓] Navigate | [Enter] Edit | [Esc] Back"
	return s
}

func (m *ConfigureSystemModel) viewReviewChanges() string {
	s := "Review Changes for " + string(m.domain) + ":\n\n"
	if m.changes == nil {
		return s + "No changes"
	}
	for key, modified := range m.changes.Modified {
		original := m.changes.Original[key]
		if original != modified {
			s += key + ": " + fmt.Sprintf("%v", original) + " → " + fmt.Sprintf("%v", modified) + "\n"
		}
	}
	s += "\n[Enter] Confirm | [Esc] Back"
	return s
}

func (m *ConfigureSystemModel) viewConfirm() string {
	s := "Confirm Configuration Changes?\n\n"
	s += "Domain: " + string(m.domain) + "\n"
	if m.changes != nil {
		changeCount := 0
		for key, modified := range m.changes.Modified {
			original := m.changes.Original[key]
			if original != modified {
				changeCount++
			}
		}
		s += fmt.Sprintf("Changes: %d\n", changeCount)
	}
	s += "\n[Y/Enter] Confirm | [N/Esc] Cancel"
	return s
}

func (m *ConfigureSystemModel) viewSaving() string {
	return "Saving configuration...\n\nPlease wait..."
}

func (m *ConfigureSystemModel) viewComplete() string {
	s := "Configuration Updated!\n\n"
	s += "Domain: " + string(m.domain) + "\n"
	s += "Changes saved successfully.\n\n"
	s += "[Enter] Done"
	return s
}

func (m *ConfigureSystemModel) viewFailed() string {
	s := "Configuration Failed\n\n"
	if m.error != nil {
		s += "Error: " + m.error.Message + "\n"
	}
	s += "\n[R] Retry | [Esc] Cancel"
	return s
}

// ConfigCompleteMsg represents completion of a configuration save
type ConfigCompleteMsg struct {
	Result *ConfigureSystemResult
}

// ConfigErrorMsg represents an error during configuration save
type ConfigErrorMsg struct {
	Error *IntentError
}
