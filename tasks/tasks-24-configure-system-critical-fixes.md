# Task 24: ConfigureSystem Critical Fixes

**Created**: 2026-01-08
**Status**: Ready for Implementation
**Priority**: CRITICAL
**Estimated Time**: 4-5 hours
**Related**: Codebase Audit (2026-01-08)

---

## Overview

The ConfigureSystem intent's save operation is **completely stubbed** - it always returns success without persisting any changes. The EditSettings state has no actual input components for editing values.

**Current State**:
- `startSave()` is a stub that always returns success (line 448-459)
- `updateEditSettings()` has no input handling - just Enter/Esc (line 309-330)
- No text input components for editing string/int values
- No toggle mechanism for boolean values
- Configuration values are hardcoded placeholders (line 86-162)
- No configuration file persistence

**Impact**: Users cannot actually configure the system. Settings changes are lost immediately. The feature is non-functional.

---

## Files to Modify

- [ ] `internal/cli/intents/configure_system.go` (main implementation)
- [ ] `internal/cli/intents/configure_system_intent.go` (add config file handling)
- [ ] `internal/cli/intents/configure_system_test.go` (update tests)
- [ ] Create: `internal/config/config.go` (configuration persistence)

---

## Implementation Plan

### Phase 1: Configuration Persistence (1-2 hours)

#### 1.1 Create Configuration File Handler
**New File**: `internal/config/config.go`

```go
package config

import (
    "fmt"
    "os"
    "path/filepath"
    "gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
    System  SystemConfig  `yaml:"system"`
    Profile ProfileConfig `yaml:"profile"`
    CV      CVConfig      `yaml:"cv"`
    Export  ExportConfig  `yaml:"export"`
    Display DisplayConfig `yaml:"display"`
}

type SystemConfig struct {
    DataDir     string `yaml:"data_dir"`
    LogLevel    string `yaml:"log_level"`
    AutoBackup  bool   `yaml:"auto_backup"`
    BackupCount int    `yaml:"backup_count"`
}

type ProfileConfig struct {
    Name          string `yaml:"name"`
    Email         string `yaml:"email"`
    DefaultRole   string `yaml:"default_role"`
    DefaultAudience string `yaml:"default_audience"`
}

// ... other config types

// LoadConfig loads configuration from file
func LoadConfig() (*Config, error) {
    path, err := GetConfigPath()
    if err != nil {
        return nil, err
    }
    
    // Check if config file exists
    if _, err := os.Stat(path); os.IsNotExist(err) {
        // Return default config
        return DefaultConfig(), nil
    }
    
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }
    
    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }
    
    return &cfg, nil
}

// SaveConfig saves configuration to file
func SaveConfig(cfg *Config) error {
    path, err := GetConfigPath()
    if err != nil {
        return err
    }
    
    // Create directory if it doesn't exist
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create config directory: %w", err)
    }
    
    data, err := yaml.Marshal(cfg)
    if err != nil {
        return fmt.Errorf("failed to marshal config: %w", err)
    }
    
    if err := os.WriteFile(path, data, 0644); err != nil {
        return fmt.Errorf("failed to write config: %w", err)
    }
    
    return nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", fmt.Errorf("failed to get home directory: %w", err)
    }
    
    return filepath.Join(homeDir, ".kariya", "config.yaml"), nil
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
    homeDir, _ := os.UserHomeDir()
    dataDir := filepath.Join(homeDir, ".kariya")
    
    return &Config{
        System: SystemConfig{
            DataDir:     dataDir,
            LogLevel:    "info",
            AutoBackup:  true,
            BackupCount: 5,
        },
        Profile: ProfileConfig{
            Name:            "",  // User should set
            Email:           "",  // User should set
            DefaultRole:     "senior_ic",
            DefaultAudience: "technical",
        },
        // ... other defaults
    }
}
```

**Tasks**:
- [ ] Create `internal/config/config.go`
- [ ] Define Config struct with all setting domains
- [ ] Implement LoadConfig() function
- [ ] Implement SaveConfig() function
- [ ] Implement GetConfigPath() function
- [ ] Implement DefaultConfig() function
- [ ] Add tests for config loading/saving

#### 1.2 Update Intent Context
**Location**: `internal/cli/intents/configure_system_intent.go`

**Add**:
```go
type ConfigureSystemContext struct {
    Domains  []ConfigurationDomain
    Settings map[ConfigurationDomain][]*ConfigurationSetting
    Config   *config.Config  // ADD THIS
}

// Update initialization
func NewConfigureSystemIntent(context *ConfigureSystemContext) *ConfigureSystemIntent {
    // Load config from file
    cfg, err := config.LoadConfig()
    if err != nil {
        // Use defaults, log error
        cfg = config.DefaultConfig()
    }
    context.Config = cfg
    
    // Populate settings from loaded config
    context.Settings = settingsFromConfig(cfg)
    
    // ... rest of initialization
}
```

**Tasks**:
- [ ] Add Config field to context
- [ ] Load config in intent initialization
- [ ] Populate settings from loaded config
- [ ] Handle config load errors gracefully

---

### Phase 2: Settings Editing UI (2-3 hours)

#### 2.1 Add Text Input Components
**Location**: `internal/cli/intents/configure_system.go`

**Add to ConfigureSystemModel**:
```go
type ConfigureSystemModel struct {
    // ... existing fields
    
    // Editing state
    settingsInputs   map[string]textinput.Model  // One input per setting
    focusedSetting   int                         // Which setting is focused
    editingValue     bool                        // Currently editing a value
}
```

**Initialize inputs when entering EditSettings state**:
```go
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
            input.CharLimit = 100
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
```

**Tasks**:
- [ ] Add settingsInputs map to model
- [ ] Add focusedSetting and editingValue fields
- [ ] Implement initializeInputs()
- [ ] Add validation functions (validateInt, validateString, etc.)
- [ ] Call initializeInputs() when entering EditSettings state

#### 2.2 Implement Settings Navigation and Editing
**Location**: `internal/cli/intents/configure_system.go:309-330`

**Replace stub updateEditSettings**:
```go
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
            
        case "m":
            // Return to main menu
            m.setResult(&ConfigureSystemResult{
                Success: false,
                Error: &IntentError{
                    Code:    "config_cancelled",
                    Message: "User returned to main menu",
                },
            })
            
        default:
            // Pass to focused input if editing
            if m.editingValue {
                return m.updateFocusedInput(msg)
            }
        }
    }
    return nil
}

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
    case "string":
        newValue = input.Value()
    case "int":
        newValue, _ = strconv.Atoi(input.Value())
    case "bool":
        newValue = input.Value() == "true"
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
```

**Tasks**:
- [ ] Implement full updateEditSettings() with navigation
- [ ] Add j/k vim-style navigation
- [ ] Implement value editing mode
- [ ] Implement saveCurrentValue()
- [ ] Implement revertCurrentValue()
- [ ] Implement updateInputFocus()
- [ ] Test navigation and editing

#### 2.3 Update View to Show Inputs
**Location**: `internal/cli/intents/configure_system.go`

**Add to viewEditSettings()**:
```go
func (m *ConfigureSystemModel) viewEditSettings() string {
    settings := m.context.Settings[m.domain]
    
    s := "Edit " + string(m.domain) + " Settings:\n\n"
    
    for i, setting := range settings {
        prefix := "  "
        if i == m.focusedSetting {
            prefix = "> "
        }
        
        // Show input if available
        input, hasInput := m.settingsInputs[setting.Key]
        
        s += prefix + setting.Key + ": "
        
        if hasInput {
            if i == m.focusedSetting && m.editingValue {
                s += input.View() + " (editing)"
            } else {
                s += input.Value()
            }
        } else {
            s += fmt.Sprintf("%v", setting.Value)
        }
        
        s += "\n"
        s += "    " + setting.Description + "\n\n"
    }
    
    footer := "↑/↓ or j/k: Navigate | Enter: Edit | Esc: Back | Ctrl+S: Save All"
    if m.editingValue {
        footer = "Type to edit | Enter: Confirm | Esc: Cancel"
    }
    
    return s + "\n" + footer
}
```

**Tasks**:
- [ ] Update viewEditSettings() to show text inputs
- [ ] Show editing indicator when active
- [ ] Update footer based on editing state
- [ ] Test view rendering

---

### Phase 3: Real Save Operation (1 hour)

#### 3.1 Implement Actual Persistence
**Location**: `internal/cli/intents/configure_system.go:448-459`

**Replace stub startSave()**:
```go
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

func applyConfigChange(cfg *config.Config, domain ConfigurationDomain, key string, value interface{}) error {
    switch domain {
    case DomainSystem:
        return applySystemChange(&cfg.System, key, value)
    case DomainProfile:
        return applyProfileChange(&cfg.Profile, key, value)
    // ... other domains
    }
    return fmt.Errorf("unknown domain: %s", domain)
}

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

// Similar for other domains...
```

**Tasks**:
- [ ] Replace stub startSave() with real implementation
- [ ] Implement applyConfigChange()
- [ ] Implement domain-specific apply functions
- [ ] Add error handling for file operations
- [ ] Test save operation writes to file
- [ ] Test save operation handles errors

#### 3.2 Load Config on Startup
**Location**: `internal/cli/app/app.go`

**Update GlobalContext initialization**:
```go
// Load configuration
cfg, err := config.LoadConfig()
if err != nil {
    // Log error, use defaults
    cfg = config.DefaultConfig()
}

// Initialize ConfigureSystem context with loaded config
configContext := &intents.ConfigureSystemContext{
    Domains:  configDomains,
    Settings: settingsFromConfig(cfg),
    Config:   cfg,
}
```

**Tasks**:
- [ ] Load config on app startup
- [ ] Pass config to ConfigureSystem context
- [ ] Handle load errors gracefully
- [ ] Test config loading

---

### Phase 4: Additional Fixes (30 min)

#### 4.1 Add Vim Navigation to Domain Selection
**Location**: `internal/cli/intents/configure_system.go:261-307`

**Add to updateSelectDomain()**:
```go
case "up", "k":  // ADD "k"
    if m.selectedIndex > 0 {
        m.selectedIndex--
    }
case "down", "j":  // ADD "j"
    if m.selectedIndex < len(m.context.Domains)-1 {
        m.selectedIndex++
    }
```

**Tasks**:
- [ ] Add "k" key handler (up)
- [ ] Add "j" key handler (down)
- [ ] Test vim navigation

#### 4.2 Use LoadingRotator
**Similar to Task 23** - use rotating messages in viewSaving()

**Tasks**:
- [ ] Use loadingRotator.GetMessage() in viewSaving()
- [ ] Add tick command for rotation
- [ ] Test loading animation

#### 4.3 Clear Error on Retry
**Location**: `internal/cli/intents/configure_system.go:431-446`

**Fix updateFailed()**:
```go
case "r":
    // Retry - clear error and go back to Confirm
    m.error = nil  // ADD THIS
    m.state = ConfigStateConfirm
```

**Tasks**:
- [ ] Clear m.error when retrying
- [ ] Test error clears on retry

---

## Acceptance Criteria

### Must Have
- [ ] Can navigate between settings with j/k or arrow keys
- [ ] Can edit string values with text input
- [ ] Can edit integer values with validation
- [ ] Can edit boolean values (toggle or true/false input)
- [ ] Changes are saved to `~/.kariya/config.yaml`
- [ ] Config is loaded on application startup
- [ ] Config file is created if it doesn't exist
- [ ] Default values are used if config file doesn't exist
- [ ] Error handling for file system operations
- [ ] All tests passing (maintain 2,078/2,078)
- [ ] Zero race conditions
- [ ] Build successful

### Should Have
- [ ] Vim j/k navigation in domain selection
- [ ] LoadingRotator shows during save
- [ ] Error state clears on retry
- [ ] Validation errors shown for invalid input
- [ ] Changes can be cancelled (Esc from editing)

### Nice to Have
- [ ] Select dropdown for enum values
- [ ] Color picker for color values
- [ ] Path browser for directory values
- [ ] Confirmation before overwriting changes

---

## Testing Strategy

### TDD Approach

**Red**: Write failing tests
```go
It("should save changes to config file", func() {
    // Setup
    model := setupConfigIntent()
    model.domain = DomainSystem
    model.changes = &ConfigurationChanges{
        Modified: map[string]interface{}{
            "log_level": "debug",
        },
    }
    
    // Execute save
    msg := model.startSave()()
    
    // Verify
    result := msg.(ConfigCompleteMsg).Result
    Expect(result.Success).To(BeTrue())
    
    // Verify file written
    cfg, err := config.LoadConfig()
    Expect(err).NotTo(HaveOccurred())
    Expect(cfg.System.LogLevel).To(Equal("debug"))
})

It("should allow editing setting value", func() {
    model := setupConfigIntent()
    model.state = ConfigStateEditSettings
    model.focusedSetting = 0
    
    // Start editing
    model.Update(tea.KeyMsg{Type: tea.KeyEnter})
    Expect(model.editingValue).To(BeTrue())
    
    // Type value
    model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("new value")})
    
    // Confirm
    model.Update(tea.KeyMsg{Type: tea.KeyEnter})
    Expect(model.editingValue).To(BeFalse())
    Expect(model.changes.Modified).To(HaveKey("setting_key"))
})
```

**Green**: Implement
**Refactor**: Clean up

### Verification Commands
```bash
# Unit tests
go test ./internal/cli/intents -v -run TestConfigureSystem
go test ./internal/config -v

# Integration test (manual)
go build -o kariya ./cmd/cli
./kariya
# Navigate to: Configure > System > Edit settings
# Change a value, save, restart app
# Verify value persisted

# Check config file
cat ~/.kariya/config.yaml

# Full suite
go test ./... -v
go test -race ./...

# Compliance
make check-compliance
```

---

## Risk Assessment

### High Risk
1. **Config file corruption**
   - Mitigation: Validate before saving, keep backup

2. **Breaking existing config files**
   - Mitigation: Version config format, migration path

### Medium Risk
3. **File permissions issues**
   - Mitigation: Check permissions, show clear errors

4. **Invalid values breaking app**
   - Mitigation: Validate all inputs, use defaults on error

### Low Risk
5. **Input validation edge cases**
   - Mitigation: Comprehensive validation tests

---

## References

### Related Documents
- `docs/rules/master-task-prompt.md` - Task workflow
- `docs/TUI_STANDARDS.md` - TUI design standards

### Related Files
- `internal/cli/intents/configure_system.go` (561 lines)
- `internal/cli/intents/configure_system_intent.go` (174 lines)
- `internal/cli/models/form.go` - Example of text input usage

### Example Config File
```yaml
system:
  data_dir: /home/user/.kariya
  log_level: info
  auto_backup: true
  backup_count: 5

profile:
  name: "John Doe"
  email: "john@example.com"
  default_role: senior_ic
  default_audience: technical

cv:
  default_format: markdown
  max_bullets: 50

export:
  default_destination: file
  auto_open: false

display:
  theme: dark
  animations: true
```

---

**Last Updated**: 2026-01-08
**Author**: AI Assistant (via OpenCode)
**Status**: Ready for implementation
