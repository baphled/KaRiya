// Package configtypes provides shared types for the configuration wizard. It
// exists as a separate package to avoid import cycles between intents and
// screens that both need access to configuration domain enumerations and
// setting structures.
package configtypes

// ConfigurationDomain partitions application settings into logical groups that
// the configuration wizard displays as top-level menu items. Each domain maps
// to a distinct set of ConfigurationSetting values loaded from the config store.
type ConfigurationDomain string

const (
	// DomainSystem covers application-wide settings such as database path, log level,
	// and data storage preferences.
	DomainSystem ConfigurationDomain = "system"
	// DomainProfile covers personal profile settings including name, summary, contact
	// information, and professional headline used in CV generation.
	DomainProfile ConfigurationDomain = "profile"
	// DomainExport covers default export preferences such as output directory, preferred
	// format, and default destination.
	DomainExport ConfigurationDomain = "export"
	// DomainUI covers visual and interaction preferences including theme selection,
	// color scheme, and table page sizes.
	DomainUI ConfigurationDomain = "ui"
)

// ConfigurationSetting holds the metadata and current value for a single
// editable setting within a configuration domain. The wizard uses these
// fields to render the appropriate input control and to detect changes.
type ConfigurationSetting struct {
	Key          string
	Label        string
	Value        interface{}
	DefaultValue interface{}
	Type         string
	Options      []string
	Description  string
}

// ConfigurationState models the linear wizard flow of the configuration intent.
// The user advances through these states in order; failure or cancellation can
// short-circuit the progression back to the menu.
type ConfigurationState string

const (
	// ConfigStateSelectDomain is the initial state where the user picks which configuration
	// domain (system, profile, export, UI) to edit.
	ConfigStateSelectDomain ConfigurationState = "select_domain"
	// ConfigStateEditSettings presents the editable fields for the selected domain, allowing
	// the user to modify individual key-value settings.
	ConfigStateEditSettings ConfigurationState = "edit_settings"
	// ConfigStateReviewChanges shows a side-by-side diff of original versus modified values
	// so the user can verify changes before committing them.
	ConfigStateReviewChanges ConfigurationState = "review_changes"
	// ConfigStateConfirm prompts the user with a yes/no confirmation before persisting the
	// reviewed changes.
	ConfigStateConfirm ConfigurationState = "confirm"
	// ConfigStateSaving is a transient state during which the confirmed changes are being
	// written to the config store. A spinner is typically displayed.
	ConfigStateSaving ConfigurationState = "saving"
	// ConfigStateComplete indicates that all changes have been successfully persisted. The
	// intent shows a success message and prepares to return to the menu.
	ConfigStateComplete ConfigurationState = "complete"
	// ConfigStateFailed indicates that the save operation encountered an error. The intent
	// displays the error details and allows the user to retry or cancel.
	ConfigStateFailed ConfigurationState = "failed"
)

// ConfigurationChanges captures the before-and-after state of a single editing
// session within one configuration domain. The review screen diffs Original
// against Modified to highlight what the user changed before confirming.
type ConfigurationChanges struct {
	Domain   ConfigurationDomain
	Original map[string]interface{}
	Modified map[string]interface{}
}
