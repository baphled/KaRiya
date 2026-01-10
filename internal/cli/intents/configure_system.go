package intents

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/baphled/kariya/internal/cli/navigation"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/config"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// titleCase converts a string to title case using the proper Go 1.18+ API
func titleCase(s string) string {
	return cases.Title(language.English).String(s)
}

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
	Config   *config.Config // Configuration loaded from file
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

	// Editing state
	settingsInputs map[string]textinput.Model // One input per setting
	focusedSetting int                        // Which setting is focused
	editingValue   bool                       // Currently editing a value

	// Theme support - set by parent intent
	theme themes.Theme

	// Navigation handler for domain selection
	domainNavHandler *navigation.ListNavigationHandler
}

// SetTheme sets the theme for the model (called by parent intent)
func (m *ConfigureSystemModel) SetTheme(theme themes.Theme) {
	m.theme = theme
}

// Theme helper methods for consistent themed styling.

// getCardStyle returns a themed card style, with fallback to default styling.
func (m *ConfigureSystemModel) getCardStyle() lipgloss.Style {
	if m.theme != nil {
		return m.theme.Styles().CardBase
	}
	// Fallback to default styling
	return lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)
}

// getCardStyleWithBorder returns a card style with custom border color.
func (m *ConfigureSystemModel) getCardStyleWithBorder(borderColor lipgloss.Color) lipgloss.Style {
	if m.theme != nil {
		return m.theme.Styles().CardBase.BorderForeground(borderColor)
	}
	return lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)
}

// getPrimaryColor returns the primary text color from theme or fallback.
func (m *ConfigureSystemModel) getPrimaryColor() lipgloss.Color {
	if m.theme != nil {
		return m.theme.ForegroundColor()
	}
	return styles.ColorTextPrimary
}

// getAccentColor returns the accent color from theme or fallback.
func (m *ConfigureSystemModel) getAccentColor() lipgloss.Color {
	if m.theme != nil {
		return m.theme.PrimaryColor()
	}
	return styles.ColorAccentTeal
}

// getMutedColor returns the muted text color from theme or fallback.
func (m *ConfigureSystemModel) getMutedColor() lipgloss.Color {
	if m.theme != nil {
		return m.theme.MutedColor()
	}
	return styles.ColorTextMuted
}

// getSuccessColor returns the success color from theme or fallback.
func (m *ConfigureSystemModel) getSuccessColor() lipgloss.Color {
	if m.theme != nil {
		return m.theme.SuccessColor()
	}
	return styles.ColorSuccess
}

// getErrorColor returns the error color from theme or fallback.
func (m *ConfigureSystemModel) getErrorColor() lipgloss.Color {
	if m.theme != nil {
		return m.theme.ErrorColor()
	}
	return styles.ColorError
}

// getInfoColor returns the info color from theme or fallback.
func (m *ConfigureSystemModel) getInfoColor() lipgloss.Color {
	if m.theme != nil {
		return m.theme.InfoColor()
	}
	return styles.ColorInfo
}

// NewConfigureSystemContext creates a new ConfigureSystemContext with values loaded from config file
func NewConfigureSystemContext() *ConfigureSystemContext {
	// Load config from file (or use defaults if file doesn't exist)
	cfg, err := config.LoadConfig()
	if err != nil {
		// Log error but continue with defaults
		cfg = config.DefaultConfig()
	}

	// Populate settings from loaded config
	settings := settingsFromConfig(cfg)

	return &ConfigureSystemContext{
		Domains: []ConfigurationDomain{
			DomainSystem,
			DomainProfile,
			DomainExport,
			DomainUI,
		},
		Settings: settings,
		Config:   cfg,
	}
}

// settingsFromConfig populates settings from config values
func settingsFromConfig(cfg *config.Config) map[ConfigurationDomain][]*ConfigurationSetting {
	return map[ConfigurationDomain][]*ConfigurationSetting{
		DomainSystem: {
			{
				Key:          "log_level",
				Label:        "Log Level",
				Value:        cfg.System.LogLevel,
				DefaultValue: "info",
				Type:         "select",
				Options:      []string{"debug", "info", "warn", "error"},
				Description:  "Logging verbosity level",
			},
			{
				Key:          "data_dir",
				Label:        "Data Directory",
				Value:        cfg.System.DataDir,
				DefaultValue: cfg.System.DataDir,
				Type:         "string",
				Description:  "Directory for storing application data",
			},
			{
				Key:          "auto_backup",
				Label:        "Auto Backup",
				Value:        cfg.System.AutoBackup,
				DefaultValue: true,
				Type:         "bool",
				Description:  "Automatically backup data",
			},
			{
				Key:          "backup_count",
				Label:        "Backup Count",
				Value:        cfg.System.BackupCount,
				DefaultValue: 5,
				Type:         "int",
				Description:  "Number of backups to keep",
			},
		},
		DomainProfile: {
			{
				Key:          "name",
				Label:        "Full Name",
				Value:        cfg.Profile.Name,
				DefaultValue: "",
				Type:         "string",
				Description:  "Your full name",
			},
			{
				Key:          "email",
				Label:        "Email",
				Value:        cfg.Profile.Email,
				DefaultValue: "",
				Type:         "string",
				Description:  "Your email address",
			},
			{
				Key:          "default_role",
				Label:        "Default Role",
				Value:        cfg.Profile.DefaultRole,
				DefaultValue: "senior_ic",
				Type:         "select",
				Options:      []string{"junior_ic", "mid_ic", "senior_ic", "staff_ic", "principal_ic", "manager", "senior_manager", "director"},
				Description:  "Default role for CV generation",
			},
			{
				Key:          "default_audience",
				Label:        "Default Audience",
				Value:        cfg.Profile.DefaultAudience,
				DefaultValue: "technical",
				Type:         "select",
				Options:      []string{"technical", "executive", "general"},
				Description:  "Default audience for CV generation",
			},
		},
		DomainExport: {
			{
				Key:          "default_destination",
				Label:        "Default Destination",
				Value:        cfg.Export.DefaultDestination,
				DefaultValue: "file",
				Type:         "select",
				Options:      []string{"file", "clipboard"},
				Description:  "Default export destination",
			},
			{
				Key:          "auto_open",
				Label:        "Auto-Open Exported Files",
				Value:        cfg.Export.AutoOpen,
				DefaultValue: false,
				Type:         "bool",
				Description:  "Automatically open exported files",
			},
		},
		DomainUI: {
			{
				Key:          "theme",
				Label:        "Theme",
				Value:        cfg.Display.Theme,
				DefaultValue: "dark",
				Type:         "select",
				Options:      []string{"light", "dark"},
				Description:  "UI color theme",
			},
			{
				Key:          "animations",
				Label:        "Animations",
				Value:        cfg.Display.Animations,
				DefaultValue: true,
				Type:         "bool",
				Description:  "Enable UI animations",
			},
		},
	}
}

// NewConfigureSystemModel creates a new ConfigureSystem intent model
func NewConfigureSystemModel(ctx *ConfigureSystemContext) *ConfigureSystemModel {
	model := &ConfigureSystemModel{
		state:         ConfigStateSelectDomain,
		context:       ctx,
		selectedIndex: 0,
		active:        false,
	}
	// Initialize domain navigation handler
	model.domainNavHandler = navigation.NewListNavigationHandler(model)
	return model
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
		return nil
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

// Helper methods for input management

// initializeInputs creates text input components for all settings in the current domain
func (m *ConfigureSystemModel) initializeInputs() {
	m.settingsInputs = make(map[string]textinput.Model)
	settings := m.context.Settings[m.domain]

	for _, setting := range settings {
		input := textinput.New()
		input.Placeholder = setting.Description
		input.SetValue(fmt.Sprintf("%v", setting.Value))

		// Configure based on type
		switch setting.Type {
		case "int":
			input.CharLimit = 10
			input.Validate = validateInt
		case "string":
			input.CharLimit = 200
		case "bool":
			// For bool, we'll use "true"/"false" strings
			input.CharLimit = 5
		case "select":
			// For select, show current value (user will need to type exact match)
			input.CharLimit = 50
		}

		m.settingsInputs[setting.Key] = input
	}

	// Focus first input
	if len(settings) > 0 {
		first := m.settingsInputs[settings[0].Key]
		first.Focus()
		m.settingsInputs[settings[0].Key] = first
	}
}

// validateInt validates that input is a valid integer
func validateInt(s string) error {
	if s == "" {
		return nil
	}
	_, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("must be a number")
	}
	return nil
}

// updateInputFocus moves focus to the currently selected setting
func (m *ConfigureSystemModel) updateInputFocus() {
	settings := m.context.Settings[m.domain]

	// Blur all inputs
	for key, input := range m.settingsInputs {
		input.Blur()
		m.settingsInputs[key] = input
	}

	// Focus current
	if m.focusedSetting < len(settings) {
		key := settings[m.focusedSetting].Key
		input := m.settingsInputs[key]
		input.Focus()
		m.settingsInputs[key] = input
	}
}

// updateFocusedInput updates the currently focused input with the message
func (m *ConfigureSystemModel) updateFocusedInput(msg tea.Msg) tea.Cmd {
	settings := m.context.Settings[m.domain]
	if m.focusedSetting >= len(settings) {
		return nil
	}

	key := settings[m.focusedSetting].Key
	input := m.settingsInputs[key]

	var cmd tea.Cmd
	input, cmd = input.Update(msg)
	m.settingsInputs[key] = input

	return cmd
}

// saveCurrentValue saves the current input value to the changes map
func (m *ConfigureSystemModel) saveCurrentValue() {
	settings := m.context.Settings[m.domain]
	if m.focusedSetting >= len(settings) {
		return
	}

	setting := settings[m.focusedSetting]
	input := m.settingsInputs[setting.Key]

	// Convert value based on type
	var newValue interface{}
	switch setting.Type {
	case "string", "select":
		newValue = input.Value()
	case "int":
		val, err := strconv.Atoi(input.Value())
		if err != nil {
			// Invalid int, revert
			m.revertCurrentValue()
			return
		}
		newValue = val
	case "bool":
		v := input.Value()
		newValue = (v == "true" || v == "True" || v == "TRUE")
	}

	// Store in changes
	if m.changes == nil {
		m.changes = &ConfigurationChanges{
			Domain:   m.domain,
			Original: make(map[string]interface{}),
			Modified: make(map[string]interface{}),
		}
	}

	if _, exists := m.changes.Original[setting.Key]; !exists {
		m.changes.Original[setting.Key] = setting.Value
	}
	m.changes.Modified[setting.Key] = newValue

	// Update setting value for display
	setting.Value = newValue
}

// revertCurrentValue reverts the current input to its original value
func (m *ConfigureSystemModel) revertCurrentValue() {
	settings := m.context.Settings[m.domain]
	if m.focusedSetting >= len(settings) {
		return
	}

	setting := settings[m.focusedSetting]
	input := m.settingsInputs[setting.Key]

	// Reset to original value
	input.SetValue(fmt.Sprintf("%v", setting.Value))
	m.settingsInputs[setting.Key] = input
}

// State-specific update handlers

func (m *ConfigureSystemModel) updateSelectDomain(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			// TODO: Toggle help modal when integrated into BaseIntent
			return nil
		case KeyBack:
			// At root state, back means cancel
			m.setResult(&ConfigureSystemResult{
				Success: false,
				Error: &IntentError{
					Code:    "config_cancelled",
					Message: "Configuration cancelled by user",
				},
			})
			return nil
		}

		// Try list navigation handler (handles up/down/j/k/pgup/pgdn/home/end/g/G)
		if m.domainNavHandler.HandleKey(msg.String()) {
			return nil
		}

		switch msg.String() {
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
			m.focusedSetting = 0
			m.state = ConfigStateEditSettings
			// Initialize inputs for editing
			m.initializeInputs()
		}
	}
	return nil
}

func (m *ConfigureSystemModel) updateEditSettings(msg tea.Msg) tea.Cmd {
	settings := m.context.Settings[m.domain]

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if !m.editingValue && m.focusedSetting > 0 {
				m.focusedSetting--
				m.updateInputFocus()
			}

		case "down", "j":
			if !m.editingValue && m.focusedSetting < len(settings)-1 {
				m.focusedSetting++
				m.updateInputFocus()
			}

		case "enter":
			if m.editingValue {
				// Save current value
				m.saveCurrentValue()
				m.editingValue = false
			} else {
				// Start editing
				m.editingValue = true
			}

		case "esc":
			if m.editingValue {
				// Cancel editing
				m.revertCurrentValue()
				m.editingValue = false
			} else {
				// Go back to domain selection
				m.selectedIndex = 0
				m.state = ConfigStateSelectDomain
			}

		case "ctrl+s":
			// Save all changes and move to review
			m.state = ConfigStateReviewChanges

		default:
			// Pass to focused input if editing
			if m.editingValue {
				return m.updateFocusedInput(msg)
			}
		}
	}
	return nil
}

func (m *ConfigureSystemModel) updateReviewChanges(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			// TODO: Toggle help modal when integrated into BaseIntent
			return nil
		case KeyBack:
			m.state = ConfigStateEditSettings
			return nil
		}

		switch msg.String() {
		case "enter":
			m.state = ConfigStateConfirm
		}
	}
	return nil
}

func (m *ConfigureSystemModel) updateConfirm(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			// TODO: Toggle help modal when integrated into BaseIntent
			return nil
		case KeyBack:
			m.state = ConfigStateReviewChanges
			return nil
		}

		switch msg.String() {
		case "y", "enter":
			m.state = ConfigStateSaving
			return m.startSave()
		case "n":
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
	case tea.KeyMsg:
		// Handle global keys (q=quit, ?=help)
		// Note: esc allows going back even during save
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			m.setResult(&ConfigureSystemResult{
				Success: false,
				Error: &IntentError{
					Code:    "config_cancelled",
					Message: "User quit during save",
				},
			})
			return tea.Quit
		case KeyHelp:
			// TODO: Toggle help modal when integrated into BaseIntent
			return nil
		case KeyBack:
			// Note: Save operation continues in background per user decision
			// Navigate back to review changes
			m.state = ConfigStateReviewChanges
			return nil
		}
	}
	return nil
}

func (m *ConfigureSystemModel) updateComplete(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyBack:
			m.active = false
			return nil
		}

		switch msg.String() {
		case "enter":
			m.active = false
		}
	}
	return nil
}

func (m *ConfigureSystemModel) updateFailed(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			// TODO: Toggle help modal when integrated into BaseIntent
			return nil
		case KeyBack:
			m.active = false
			return nil
		}

		switch msg.String() {
		case "r":
			// Retry - clear error and go back to Confirm
			m.error = nil
			m.state = ConfigStateConfirm
		}
	}
	return nil
}

// startSave initiates the async save operation
func (m *ConfigureSystemModel) startSave() tea.Cmd {
	return func() tea.Msg {
		// Apply changes to config
		cfg := m.context.Config

		for key, value := range m.changes.Modified {
			if err := applyConfigChange(cfg, m.domain, key, value); err != nil {
				return ConfigCompleteMsg{
					Result: &ConfigureSystemResult{
						Success: false,
						Error: &IntentError{
							Code:    "save_failed",
							Message: fmt.Sprintf("Failed to apply change: %s", err),
						},
					},
				}
			}
		}

		// Save to file
		if err := config.SaveConfig(cfg); err != nil {
			return ConfigCompleteMsg{
				Result: &ConfigureSystemResult{
					Success: false,
					Error: &IntentError{
						Code:    "save_failed",
						Message: fmt.Sprintf("Failed to save config: %s", err),
					},
				},
			}
		}

		// Success
		return ConfigCompleteMsg{
			Result: &ConfigureSystemResult{
				Success: true,
				Domain:  m.domain,
				Changes: m.changes,
			},
		}
	}
}

// applyConfigChange applies a configuration change to the appropriate domain
func applyConfigChange(cfg *config.Config, domain ConfigurationDomain, key string, value interface{}) error {
	switch domain {
	case DomainSystem:
		return applySystemChange(&cfg.System, key, value)
	case DomainProfile:
		return applyProfileChange(&cfg.Profile, key, value)
	case DomainExport:
		return applyExportChange(&cfg.Export, key, value)
	case DomainUI:
		return applyDisplayChange(&cfg.Display, key, value)
	}
	return fmt.Errorf("unknown domain: %s", domain)
}

// applySystemChange applies a change to system configuration
func applySystemChange(sys *config.SystemConfig, key string, value interface{}) error {
	switch key {
	case "data_dir":
		sys.DataDir = value.(string)
	case "log_level":
		sys.LogLevel = value.(string)
	case "auto_backup":
		sys.AutoBackup = value.(bool)
	case "backup_count":
		sys.BackupCount = value.(int)
	default:
		return fmt.Errorf("unknown system setting: %s", key)
	}
	return nil
}

// applyProfileChange applies a change to profile configuration
func applyProfileChange(prof *config.ProfileConfig, key string, value interface{}) error {
	switch key {
	case "name":
		prof.Name = value.(string)
	case "email":
		prof.Email = value.(string)
	case "default_role":
		prof.DefaultRole = value.(string)
	case "default_audience":
		prof.DefaultAudience = value.(string)
	default:
		return fmt.Errorf("unknown profile setting: %s", key)
	}
	return nil
}

// applyExportChange applies a change to export configuration
func applyExportChange(exp *config.ExportConfig, key string, value interface{}) error {
	switch key {
	case "default_destination":
		exp.DefaultDestination = value.(string)
	case "auto_open":
		exp.AutoOpen = value.(bool)
	default:
		return fmt.Errorf("unknown export setting: %s", key)
	}
	return nil
}

// applyDisplayChange applies a change to display configuration
func applyDisplayChange(disp *config.DisplayConfig, key string, value interface{}) error {
	switch key {
	case "theme":
		disp.Theme = value.(string)
	case "animations":
		disp.Animations = value.(bool)
	default:
		return fmt.Errorf("unknown display setting: %s", key)
	}
	return nil
}

// setResult sets the intent result and marks intent as complete
func (m *ConfigureSystemModel) setResult(result *ConfigureSystemResult) {
	m.result = result
	m.active = false
}

// View rendering methods

func (m *ConfigureSystemModel) viewSelectDomain() string {
	var content strings.Builder
	content.WriteString("\n⚙️  Select Configuration Domain\n\n")

	// Domain descriptions
	domainDescs := map[ConfigurationDomain]string{
		DomainSystem:  "Log level, data directory, backups",
		DomainProfile: "Name, email, default roles",
		DomainExport:  "Default formats and destinations",
		DomainUI:      "Theme and animations",
	}

	for i, d := range m.context.Domains {
		prefix := "  "
		itemStyle := lipgloss.NewStyle().Foreground(m.getPrimaryColor())

		if i == m.selectedIndex {
			prefix = "▶ "
			itemStyle = itemStyle.Foreground(m.getAccentColor()).Bold(true)
		}

		// Domain name with capitalization
		domainName := titleCase(string(d))
		line := fmt.Sprintf("%s%s", prefix, domainName)
		content.WriteString(itemStyle.Render(line))

		// Add description in muted color
		if desc, ok := domainDescs[d]; ok {
			descStyle := lipgloss.NewStyle().Foreground(m.getMutedColor())
			content.WriteString(" " + descStyle.Render("- "+desc))
		}
		content.WriteString("\n")
	}

	return m.getCardStyle().Render(content.String())
}

func (m *ConfigureSystemModel) viewEditSettings() string {
	var content strings.Builder
	content.WriteString(fmt.Sprintf("\n📝 Edit %s Settings\n\n", titleCase(string(m.domain))))

	settings := m.context.Settings[m.domain]

	for i, setting := range settings {
		prefix := "  "
		labelStyle := lipgloss.NewStyle().Foreground(m.getPrimaryColor())

		if i == m.focusedSetting {
			prefix = "▶ "
			labelStyle = labelStyle.Foreground(m.getAccentColor()).Bold(true)
		}

		input, hasInput := m.settingsInputs[setting.Key]

		// Label
		content.WriteString(labelStyle.Render(fmt.Sprintf("%s%s: ", prefix, setting.Label)))

		// Value
		if hasInput {
			if i == m.focusedSetting && m.editingValue {
				content.WriteString(input.View() + " ")
				editIndicator := lipgloss.NewStyle().Foreground(m.getInfoColor()).Render("(editing)")
				content.WriteString(editIndicator)
			} else {
				content.WriteString(input.Value())
			}
		} else {
			content.WriteString(fmt.Sprintf("%v", setting.Value))
		}
		content.WriteString("\n")

		// Description (muted)
		descStyle := lipgloss.NewStyle().Foreground(m.getMutedColor()).PaddingLeft(4)
		content.WriteString(descStyle.Render(setting.Description) + "\n\n")
	}

	return m.getCardStyle().Render(content.String())
}

func (m *ConfigureSystemModel) viewReviewChanges() string {
	var content strings.Builder
	content.WriteString(fmt.Sprintf("\n📋 Review Changes - %s\n\n", titleCase(string(m.domain))))

	if m.changes == nil {
		noChanges := lipgloss.NewStyle().Foreground(m.getMutedColor()).Render("No changes to review")
		content.WriteString(noChanges)
	} else {
		hasChanges := false
		for key, modified := range m.changes.Modified {
			original := m.changes.Original[key]
			if original != modified {
				hasChanges = true
				keyStyle := lipgloss.NewStyle().Foreground(m.getPrimaryColor()).Bold(true)
				oldStyle := lipgloss.NewStyle().Foreground(m.getErrorColor()).Strikethrough(true)
				newStyle := lipgloss.NewStyle().Foreground(m.getSuccessColor())

				content.WriteString(keyStyle.Render(key) + ": ")
				content.WriteString(oldStyle.Render(fmt.Sprintf("%v", original)))
				content.WriteString(" → ")
				content.WriteString(newStyle.Render(fmt.Sprintf("%v", modified)) + "\n")
			}
		}

		if !hasChanges {
			noChanges := lipgloss.NewStyle().Foreground(m.getMutedColor()).Render("No changes made")
			content.WriteString(noChanges)
		}
	}

	return m.getCardStyle().Render(content.String())
}

func (m *ConfigureSystemModel) viewConfirm() string {
	var content strings.Builder
	content.WriteString("\n❓ Confirm Configuration Changes?\n\n")

	content.WriteString(fmt.Sprintf("Domain: %s\n", titleCase(string(m.domain))))

	if m.changes != nil {
		changeCount := 0
		for key, modified := range m.changes.Modified {
			original := m.changes.Original[key]
			if original != modified {
				changeCount++
			}
		}
		countStyle := lipgloss.NewStyle().Foreground(m.getAccentColor()).Bold(true)
		content.WriteString(fmt.Sprintf("Changes: %s\n", countStyle.Render(fmt.Sprintf("%d", changeCount))))
	}

	return m.getCardStyle().Render(content.String())
}

func (m *ConfigureSystemModel) viewSaving() string {
	var content strings.Builder
	content.WriteString("\n⏳ Saving Configuration...\n\n")

	content.WriteString("Please wait while changes are being saved.\n\n")
	content.WriteString("• Validating configuration\n")
	content.WriteString("• Writing to disk\n")
	content.WriteString("• Applying changes\n")

	return m.getCardStyleWithBorder(m.getInfoColor()).Render(content.String())
}

func (m *ConfigureSystemModel) viewComplete() string {
	var content strings.Builder
	content.WriteString("\n✅ Configuration Updated!\n\n")

	content.WriteString(fmt.Sprintf("Domain: %s\n", titleCase(string(m.domain))))
	content.WriteString("Changes saved successfully.\n")

	return m.getCardStyleWithBorder(m.getSuccessColor()).Render(content.String())
}

func (m *ConfigureSystemModel) viewFailed() string {
	var content strings.Builder
	content.WriteString("\n❌ Configuration Failed\n\n")

	if m.error != nil {
		errorStyle := lipgloss.NewStyle().Foreground(m.getErrorColor())
		content.WriteString(errorStyle.Render("Error: "+m.error.Message) + "\n")
	}

	return m.getCardStyleWithBorder(m.getErrorColor()).Render(content.String())
}

// ConfigCompleteMsg represents completion of a configuration save
type ConfigCompleteMsg struct {
	Result *ConfigureSystemResult
}

// ConfigErrorMsg represents an error during configuration save
type ConfigErrorMsg struct {
	Error *IntentError
}

// ListNavigator interface implementation for domain selection

// GetTotalItems returns the total number of domains.
func (m *ConfigureSystemModel) GetTotalItems() int {
	return len(m.context.Domains)
}

// GetSelectedIndex returns the current selection index.
func (m *ConfigureSystemModel) GetSelectedIndex() int {
	return m.selectedIndex
}

// SetSelectedIndex sets the selection index.
func (m *ConfigureSystemModel) SetSelectedIndex(idx int) {
	// Clamp index to valid range
	if idx < 0 {
		idx = 0
	}
	if idx >= len(m.context.Domains) {
		idx = len(m.context.Domains) - 1
	}
	m.selectedIndex = idx
}

// GetPageSize returns the page size for pagination.
func (m *ConfigureSystemModel) GetPageSize() int {
	return 10
}
