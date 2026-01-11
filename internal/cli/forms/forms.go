// Package forms provides common form utilities and configurations for the KaRiya TUI.
// It uses Charm's huh library for consistent, accessible form handling.
package forms

import (
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Theme returns the Catppuccin theme configured for KaRiya forms.
// Deprecated: Use ThemedForm or themes.GenerateHuhTheme for theme-aware forms.
func Theme() *huh.Theme {
	return huh.ThemeCatppuccin()
}

// ThemedForm returns a huh.Theme that matches the given KaRiya theme.
// If theme is nil, falls back to Catppuccin theme.
func ThemedForm(theme themes.Theme) *huh.Theme {
	return themes.GenerateHuhTheme(theme)
}

// NewForm creates a new form with KaRiya's default theme and configuration.
// Note: For theme-aware forms, use NewThemedForm instead.
func NewForm(groups ...*huh.Group) *huh.Form {
	return huh.NewForm(groups...).WithTheme(Theme())
}

// NewThemedForm creates a new form with the given KaRiya theme.
// This ensures forms match the rest of the TUI styling.
func NewThemedForm(theme themes.Theme, groups ...*huh.Group) *huh.Form {
	return huh.NewForm(groups...).WithTheme(ThemedForm(theme))
}

// NewFormWithAccessible creates a form optimized for accessibility.
func NewFormWithAccessible(groups ...*huh.Group) *huh.Form {
	return huh.NewForm(groups...).
		WithTheme(Theme()).
		WithAccessible(true)
}

// NewThemedFormWithAccessible creates an accessible form with the given KaRiya theme.
func NewThemedFormWithAccessible(theme themes.Theme, groups ...*huh.Group) *huh.Form {
	return huh.NewForm(groups...).
		WithTheme(ThemedForm(theme)).
		WithAccessible(true)
}

// FormColors defines the color scheme for form elements.
var FormColors = struct {
	Title       lipgloss.Color
	Description lipgloss.Color
	Error       lipgloss.Color
	Success     lipgloss.Color
	Placeholder lipgloss.Color
}{
	Title:       lipgloss.Color("#89B4FA"), // Catppuccin Blue
	Description: lipgloss.Color("#94E2D5"), // Catppuccin Teal
	Error:       lipgloss.Color("#F38BA8"), // Catppuccin Red
	Success:     lipgloss.Color("#A6E3A1"), // Catppuccin Green
	Placeholder: lipgloss.Color("#6C7086"), // Catppuccin Overlay0
}

// Common form helper functions

// IsCompleted checks if the form has been completed by the user.
func IsCompleted(form *huh.Form) bool {
	return form.State == huh.StateCompleted
}

// IsAborted checks if the form was cancelled/aborted by the user.
func IsAborted(form *huh.Form) bool {
	return form.State == huh.StateAborted
}

// GetString safely retrieves a string value from the form.
func GetString(form *huh.Form, key string) string {
	val := form.GetString(key)
	return val
}

// GetBool safely retrieves a boolean value from the form.
func GetBool(form *huh.Form, key string) bool {
	return form.GetBool(key)
}

// GetInt safely retrieves an int value from the form.
func GetInt(form *huh.Form, key string) int {
	return form.GetInt(key)
}

// GetStrings safely retrieves a slice of strings from the form (for MultiSelect).
func GetStrings(form *huh.Form, key string) []string {
	// huh stores MultiSelect values as interface{} containing []string
	val := form.Get(key)
	if val == nil {
		return []string{}
	}

	switch v := val.(type) {
	case []string:
		return v
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	default:
		return []string{}
	}
}

// FieldConfig represents common field configuration options.
type FieldConfig struct {
	Key         string
	Title       string
	Description string
	Placeholder string
	CharLimit   int
	Required    bool
	Validate    func(string) error
}

// NewInput creates a pre-configured input field.
// Note: Prompt("> ") is set explicitly to fix a huh library display issue
// where empty fields show only the first character of the placeholder.
func NewInput(config FieldConfig) *huh.Input {
	input := huh.NewInput().
		Key(config.Key).
		Title(config.Title).
		Prompt("> ") // Explicit prompt fixes placeholder display issue

	if config.Description != "" {
		input = input.Description(config.Description)
	}

	if config.Placeholder != "" {
		input = input.Placeholder(config.Placeholder)
	}

	if config.CharLimit > 0 {
		input = input.CharLimit(config.CharLimit)
	}

	if config.Validate != nil {
		input = input.Validate(config.Validate)
	}

	return input
}

// NewText creates a pre-configured text area field.
func NewText(config FieldConfig) *huh.Text {
	text := huh.NewText().
		Key(config.Key).
		Title(config.Title)

	if config.Description != "" {
		text = text.Description(config.Description)
	}

	if config.Placeholder != "" {
		text = text.Placeholder(config.Placeholder)
	}

	if config.CharLimit > 0 {
		text = text.CharLimit(config.CharLimit)
	}

	if config.Validate != nil {
		text = text.Validate(config.Validate)
	}

	return text
}

// SelectOption represents a selectable option in a Select or MultiSelect field.
type SelectOption struct {
	Key   string
	Value string
}

// NewSelect creates a pre-configured select field.
func NewSelect(key, title, description string, options []SelectOption) *huh.Select[string] {
	huhOptions := make([]huh.Option[string], len(options))
	for i, opt := range options {
		huhOptions[i] = huh.NewOption(opt.Value, opt.Key)
	}

	sel := huh.NewSelect[string]().
		Key(key).
		Title(title).
		Options(huhOptions...)

	if description != "" {
		sel = sel.Description(description)
	}

	return sel
}

// NewMultiSelect creates a pre-configured multi-select field.
func NewMultiSelect(key, title, description string, options []SelectOption, limit int) *huh.MultiSelect[string] {
	huhOptions := make([]huh.Option[string], len(options))
	for i, opt := range options {
		huhOptions[i] = huh.NewOption(opt.Value, opt.Key)
	}

	multi := huh.NewMultiSelect[string]().
		Key(key).
		Title(title).
		Options(huhOptions...)

	if description != "" {
		multi = multi.Description(description)
	}

	if limit > 0 {
		multi = multi.Limit(limit)
	}

	return multi
}

// NewConfirm creates a pre-configured confirmation field.
func NewConfirm(key, title, description, affirmative, negative string) *huh.Confirm {
	confirm := huh.NewConfirm().
		Key(key).
		Title(title)

	if description != "" {
		confirm = confirm.Description(description)
	}

	if affirmative != "" && negative != "" {
		confirm = confirm.Affirmative(affirmative).Negative(negative)
	}

	return confirm
}
