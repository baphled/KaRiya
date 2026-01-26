// Package configtypes provides shared types for configuration functionality.
// This package exists to avoid import cycles between intents and screens.
package configtypes

// ConfigurationDomain represents a configuration domain (e.g., "system", "profile", "export").
type ConfigurationDomain string

const (
	DomainSystem  ConfigurationDomain = "system"
	DomainProfile ConfigurationDomain = "profile"
	DomainExport  ConfigurationDomain = "export"
	DomainUI      ConfigurationDomain = "ui"
)

// ConfigurationSetting represents a single configuration setting.
type ConfigurationSetting struct {
	Key          string      // e.g., "theme"
	Label        string      // e.g., "Theme"
	Value        interface{} // Current value
	DefaultValue interface{} // Default value
	Type         string      // "string", "bool", "int", "select"
	Options      []string    // For "select" type
	Description  string      // Help text
}

// ConfigurationState represents the current state of the configuration.
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

// ConfigurationChanges tracks all changes made during editing.
type ConfigurationChanges struct {
	Domain   ConfigurationDomain
	Original map[string]interface{} // Original values
	Modified map[string]interface{} // Modified values
}
