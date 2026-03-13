package configure

import (
	"errors"
	"fmt"

	"github.com/baphled/kariya/internal/config"
)

// IntentValidator holds input parameters and business logic for the ConfigureSystem intent.
type IntentValidator struct {
	Cfg      *config.Config
	Settings map[ConfigurationDomain][]*ConfigurationSetting
}

// Validate ensures all required dependencies are present before the intent can be constructed.
//
// Returns: A error value.
// Side effects: None.
func (c *IntentValidator) Validate() error {
	if c.Cfg == nil {
		return errors.New("config is required")
	}
	return nil
}

// settingsFromConfig extracts settings from config for display/editing.
func settingsFromConfig(cfg *config.Config) map[ConfigurationDomain][]*ConfigurationSetting {
	return map[ConfigurationDomain][]*ConfigurationSetting{
		DomainSystem:  systemSettings(cfg),
		DomainProfile: profileSettings(cfg),
		DomainExport:  exportSettings(cfg),
		DomainUI:      displaySettings(cfg),
	}
}

func systemSettings(cfg *config.Config) []*ConfigurationSetting {
	return []*ConfigurationSetting{
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
	}
}

func profileSettings(cfg *config.Config) []*ConfigurationSetting {
	settings := profileIdentitySettings(cfg)
	settings = append(settings, profileSkillSettings(cfg)...)
	settings = append(settings, profileDefaultsSettings(cfg)...)
	return settings
}

func profileIdentitySettings(cfg *config.Config) []*ConfigurationSetting {
	return []*ConfigurationSetting{
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
	}
}

func profileSkillSettings(cfg *config.Config) []*ConfigurationSetting {
	return []*ConfigurationSetting{
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
	}
}

func profileDefaultsSettings(cfg *config.Config) []*ConfigurationSetting {
	return []*ConfigurationSetting{
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
	}
}

func exportSettings(cfg *config.Config) []*ConfigurationSetting {
	return []*ConfigurationSetting{
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
	}
}

func displaySettings(cfg *config.Config) []*ConfigurationSetting {
	return []*ConfigurationSetting{
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
	apply, ok := profileChangeHandlers[key]
	if !ok {
		return fmt.Errorf("unknown profile setting: %s", key)
	}
	return apply(prof, value)
}

var profileChangeHandlers = map[string]func(*config.ProfileConfig, interface{}) error{
	"name":             applyProfileString(func(cfg *config.ProfileConfig) *string { return &cfg.Name }),
	"email":            applyProfileString(func(cfg *config.ProfileConfig) *string { return &cfg.Email }),
	"title":            applyProfileString(func(cfg *config.ProfileConfig) *string { return &cfg.Title }),
	"location":         applyProfileString(func(cfg *config.ProfileConfig) *string { return &cfg.Location }),
	"github":           applyProfileString(func(cfg *config.ProfileConfig) *string { return &cfg.GitHub }),
	"portfolio":        applyProfileString(func(cfg *config.ProfileConfig) *string { return &cfg.Portfolio }),
	"languages":        applyProfileStringSlice(func(cfg *config.ProfileConfig) *[]string { return &cfg.Languages }),
	"frontend":         applyProfileStringSlice(func(cfg *config.ProfileConfig) *[]string { return &cfg.Frontend }),
	"systems":          applyProfileStringSlice(func(cfg *config.ProfileConfig) *[]string { return &cfg.Systems }),
	"core_strengths":   applyProfileStringSlice(func(cfg *config.ProfileConfig) *[]string { return &cfg.CoreStrengths }),
	"what_i_bring":     applyProfileStringSlice(func(cfg *config.ProfileConfig) *[]string { return &cfg.WhatIBring }),
	"default_role":     applyProfileString(func(cfg *config.ProfileConfig) *string { return &cfg.DefaultRole }),
	"default_audience": applyProfileString(func(cfg *config.ProfileConfig) *string { return &cfg.DefaultAudience }),
}

func applyProfileString(field func(*config.ProfileConfig) *string) func(*config.ProfileConfig, interface{}) error {
	return func(prof *config.ProfileConfig, value interface{}) error {
		if v, ok := value.(string); ok {
			*field(prof) = v
		}
		return nil
	}
}

func applyProfileStringSlice(field func(*config.ProfileConfig) *[]string) func(*config.ProfileConfig, interface{}) error {
	return func(prof *config.ProfileConfig, value interface{}) error {
		if v, ok := value.([]string); ok {
			*field(prof) = v
		}
		return nil
	}
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

// ApplyChanges applies a flat map of changes to the config, routing each key
// to the correct domain using the provided settings map.
//
// Expected: cfg, settings, and changes must be valid.
// Returns: An error value if any change fails to apply.
// Side effects: Mutates cfg in place.
func ApplyChanges(cfg *config.Config, settings map[ConfigurationDomain][]*ConfigurationSetting, changes map[string]interface{}) error {
	for domain, domainSettings := range settings {
		for _, s := range domainSettings {
			if val, ok := changes[s.Key]; ok {
				if err := applyConfigChange(cfg, domain, s.Key, val); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
