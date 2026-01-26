package intents

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/configtypes"
	"github.com/baphled/kariya/internal/config"
)

// Re-export types from configtypes for backward compatibility.
type ConfigurationDomain = configtypes.ConfigurationDomain
type ConfigurationSetting = configtypes.ConfigurationSetting
type ConfigurationState = configtypes.ConfigurationState
type ConfigurationChanges = configtypes.ConfigurationChanges

// Re-export constants.
const (
	DomainSystem  = configtypes.DomainSystem
	DomainProfile = configtypes.DomainProfile
	DomainExport  = configtypes.DomainExport
	DomainUI      = configtypes.DomainUI
)

const (
	ConfigStateSelectDomain  = configtypes.ConfigStateSelectDomain
	ConfigStateEditSettings  = configtypes.ConfigStateEditSettings
	ConfigStateReviewChanges = configtypes.ConfigStateReviewChanges
	ConfigStateConfirm       = configtypes.ConfigStateConfirm
	ConfigStateSaving        = configtypes.ConfigStateSaving
	ConfigStateComplete      = configtypes.ConfigStateComplete
	ConfigStateFailed        = configtypes.ConfigStateFailed
)

// ConfigureSystemResult contains the result of configuration.
type ConfigureSystemResult struct {
	Success bool
	Domain  ConfigurationDomain
	Changes *ConfigurationChanges
	Error   *IntentError
}

// ConfigCompleteMsg represents successful configuration save.
type ConfigCompleteMsg struct {
	Result *ConfigureSystemResult
}

// ConfigErrorMsg represents an error during configuration save.
type ConfigErrorMsg struct {
	Error *IntentError
}

// settingsFromConfig extracts settings from config for display/editing.
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
				Key:          "title",
				Label:        "Professional Title",
				Value:        cfg.Profile.Title,
				DefaultValue: "",
				Type:         "string",
				Description:  "Your professional title (e.g., Senior Software Engineer)",
			},
			{
				Key:          "location",
				Label:        "Location",
				Value:        cfg.Profile.Location,
				DefaultValue: "",
				Type:         "string",
				Description:  "Your location (e.g., London, UK)",
			},
			{
				Key:          "github",
				Label:        "GitHub URL",
				Value:        cfg.Profile.GitHub,
				DefaultValue: "",
				Type:         "string",
				Description:  "Your GitHub profile URL",
			},
			{
				Key:          "portfolio",
				Label:        "Portfolio URL",
				Value:        cfg.Profile.Portfolio,
				DefaultValue: "",
				Type:         "string",
				Description:  "Your portfolio or personal website URL",
			},
			{
				Key:          "languages",
				Label:        "Programming Languages",
				Value:        cfg.Profile.Languages,
				DefaultValue: []string{},
				Type:         "list",
				Description:  "Your programming languages (comma-separated)",
			},
			{
				Key:          "frontend",
				Label:        "Frontend Technologies",
				Value:        cfg.Profile.Frontend,
				DefaultValue: []string{},
				Type:         "list",
				Description:  "Your frontend technologies (comma-separated)",
			},
			{
				Key:          "systems",
				Label:        "Systems/Infrastructure",
				Value:        cfg.Profile.Systems,
				DefaultValue: []string{},
				Type:         "list",
				Description:  "Your systems/infrastructure expertise (comma-separated)",
			},
			{
				Key:          "core_strengths",
				Label:        "Core Strengths",
				Value:        cfg.Profile.CoreStrengths,
				DefaultValue: []string{},
				Type:         "list",
				Description:  "Your core professional strengths (comma-separated)",
			},
			{
				Key:          "what_i_bring",
				Label:        "What I Bring",
				Value:        cfg.Profile.WhatIBring,
				DefaultValue: []string{},
				Type:         "list",
				Description:  "Your unique value propositions (comma-separated)",
			},
			{
				Key:          "default_role",
				Label:        "Default Role",
				Value:        cfg.Profile.DefaultRole,
				DefaultValue: "senior_ic",
				Type:         "select",
				Options: []string{
					"junior_ic", "mid_ic", "senior_ic", "staff_ic",
					"principal_ic", "manager", "senior_manager", "director",
				},
				Description: "Default role for CV generation",
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

// applyConfigChange applies a configuration change to the config object.
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

// applySystemChange applies a change to system configuration.
func applySystemChange(sys *config.SystemConfig, key string, value interface{}) error {
	switch key {
	case "data_dir":
		if v, ok := value.(string); ok {
			sys.DataDir = v
		}
	case "log_level":
		if v, ok := value.(string); ok {
			sys.LogLevel = v
		}
	case "auto_backup":
		if v, ok := value.(bool); ok {
			sys.AutoBackup = v
		}
	case "backup_count":
		if v, ok := value.(int); ok {
			sys.BackupCount = v
		}
	default:
		return fmt.Errorf("unknown system setting: %s", key)
	}
	return nil
}

// applyProfileChange applies a change to profile configuration.
func applyProfileChange(prof *config.ProfileConfig, key string, value interface{}) error {
	switch key {
	case "name":
		if v, ok := value.(string); ok {
			prof.Name = v
		}
	case "email":
		if v, ok := value.(string); ok {
			prof.Email = v
		}
	case "title":
		if v, ok := value.(string); ok {
			prof.Title = v
		}
	case "location":
		if v, ok := value.(string); ok {
			prof.Location = v
		}
	case "github":
		if v, ok := value.(string); ok {
			prof.GitHub = v
		}
	case "portfolio":
		if v, ok := value.(string); ok {
			prof.Portfolio = v
		}
	case "languages":
		if v, ok := value.([]string); ok {
			prof.Languages = v
		}
	case "frontend":
		if v, ok := value.([]string); ok {
			prof.Frontend = v
		}
	case "systems":
		if v, ok := value.([]string); ok {
			prof.Systems = v
		}
	case "core_strengths":
		if v, ok := value.([]string); ok {
			prof.CoreStrengths = v
		}
	case "what_i_bring":
		if v, ok := value.([]string); ok {
			prof.WhatIBring = v
		}
	case "default_role":
		if v, ok := value.(string); ok {
			prof.DefaultRole = v
		}
	case "default_audience":
		if v, ok := value.(string); ok {
			prof.DefaultAudience = v
		}
	default:
		return fmt.Errorf("unknown profile setting: %s", key)
	}
	return nil
}

// applyExportChange applies a change to export configuration.
func applyExportChange(exp *config.ExportConfig, key string, value interface{}) error {
	switch key {
	case "default_destination":
		if v, ok := value.(string); ok {
			exp.DefaultDestination = v
		}
	case "auto_open":
		if v, ok := value.(bool); ok {
			exp.AutoOpen = v
		}
	default:
		return fmt.Errorf("unknown export setting: %s", key)
	}
	return nil
}

// applyDisplayChange applies a change to display configuration.
func applyDisplayChange(disp *config.DisplayConfig, key string, value interface{}) error {
	switch key {
	case "theme":
		if v, ok := value.(string); ok {
			disp.Theme = v
		}
	case "animations":
		if v, ok := value.(bool); ok {
			disp.Animations = v
		}
	default:
		return fmt.Errorf("unknown display setting: %s", key)
	}
	return nil
}
