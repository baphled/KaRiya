// Package forms provides common form utilities and configurations for the KaRiya TUI.
// It uses Charm's huh library for consistent, accessible form handling.
//
// # Architecture
//
// This package serves as the ONLY allowed import point for huh in the codebase.
// Screens, modals, and intents should import forms/ instead of huh directly.
//
// # Type Aliases
//
// The package provides type aliases for common huh types:
//   - forms.Form = *huh.Form
//   - forms.Group = *huh.Group
//   - forms.Field = huh.Field
//
// # Usage in Modals
//
//	modal := &MyModal{
//	    form: forms.NewFormWithDimensions(width, height, group),
//	}
//
//	// In Update:
//	updatedForm, cmd := forms.Update(m.form, msg)
//	m.form = updatedForm
//	if forms.IsCompleted(m.form) { ... }
package forms

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/themes"
)

// Form is the form type used throughout KaRiya.
// This is an alias for *huh.Form to avoid direct huh imports in screens/modals.
type Form = *huh.Form

// Group is a group of form fields.
// This is an alias for *huh.Group to avoid direct huh imports.
type Group = *huh.Group

// Field is a form field interface.
// This is an alias for huh.Field to avoid direct huh imports.
type Field = huh.Field

// Input is a text input field.
type Input = *huh.Input

// Text is a text area field.
type Text = *huh.Text

// Select is a single-select field.
type Select = *huh.Select[string]

// MultiSelect is a multi-select field.
type MultiSelect = *huh.MultiSelect[string]

// Confirm is a confirmation field.
type Confirm = *huh.Confirm

// NewGroup creates a new form group from fields.
// Use this instead of huh.NewGroup.
func NewGroup(fields ...Field) Group {
	return huh.NewGroup(fields...)
}

// Update handles form message updates.
// Returns the updated form and any command to execute.
// Use this instead of calling form.Update directly.
func Update(form Form, msg tea.Msg) (Form, tea.Cmd) {
	model, cmd := form.Update(msg)
	f, ok := model.(*huh.Form)
	if !ok {
		return form, cmd
	}
	return f, cmd
}

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
	return huh.NewForm(groups...).
		WithTheme(Theme()).
		WithShowHelp(false).
		WithShowErrors(false)
}

// NewFormWithHeight creates a new form with KaRiya's default theme and a fixed height.
// When height is set, the form becomes scrollable if content exceeds the height.
// Use this for forms displayed in modals or constrained containers.
func NewFormWithHeight(height int, groups ...*huh.Group) *huh.Form {
	return huh.NewForm(groups...).WithTheme(Theme()).WithHeight(height)
}

// NewFormWithDimensions creates a new form with KaRiya's default theme and fixed dimensions.
// When height is set, the form becomes scrollable if content exceeds the height.
// Width controls the form's rendering width.
func NewFormWithDimensions(width, height int, groups ...*huh.Group) *huh.Form {
	form := huh.NewForm(groups...).WithTheme(Theme())
	if height > 0 {
		form = form.WithHeight(height)
	}
	if width > 0 {
		form = form.WithWidth(width)
	}
	return form
}

// NewThemedFormWithHeight creates a form with the given theme and height.
// When height is set, the form becomes scrollable if content exceeds the height.
func NewThemedFormWithHeight(theme themes.Theme, height int, groups ...*huh.Group) *huh.Form {
	return huh.NewForm(groups...).WithTheme(ThemedForm(theme)).WithHeight(height)
}

// DefaultFormHeight calculates a reasonable form height based on terminal dimensions.
// It reserves space for: logo (~7 lines), breadcrumbs (~2 lines), footer (~3 lines),
// modal chrome (~4 lines), and some padding (~4 lines) = ~20 lines overhead.
// Minimum height is 10 lines to ensure usability.
func DefaultFormHeight(terminalHeight int) int {
	const overhead = 20
	const minHeight = 10

	height := terminalHeight - overhead
	if height < minHeight {
		height = minHeight
	}
	return height
}

// ConfirmButtonHeight is the space reserved for the fixed confirm button group.
const ConfirmButtonHeight = 5

// FieldsHeight calculates the height for form fields when using a fixed confirm button.
// This reserves space for the confirm button to always be visible.
func FieldsHeight(terminalHeight int) int {
	formHeight := DefaultFormHeight(terminalHeight)
	fieldsHeight := formHeight - ConfirmButtonHeight
	if fieldsHeight < 5 {
		fieldsHeight = 5
	}
	return fieldsHeight
}

// NewFormWithFixedConfirm creates a form with scrollable fields and a fixed confirm button.
// The confirm button remains visible at the bottom while fields scroll above it.
// fieldsGroup: the form fields that can scroll
// confirmValue: pointer to bool for submit confirmation
// width, height: dimensions for the form.
func NewFormWithFixedConfirm(fieldsGroup *huh.Group, confirmValue *bool, width, height int) *huh.Form {
	// Calculate height for fields group (reserve space for confirm)
	fieldsHeight := height - ConfirmButtonHeight
	if fieldsHeight < 5 {
		fieldsHeight = 5
	}

	// Create confirm group (fixed at bottom)
	confirmGroup := huh.NewGroup(
		huh.NewConfirm().
			Key("submit").
			Title("Save Changes").
			Description("Submit the form to save your changes").
			Affirmative("Submit").
			Negative("Cancel").
			Value(confirmValue),
	)

	// Apply height to fields group to make it scrollable
	fieldsGroup = fieldsGroup.WithHeight(fieldsHeight)

	form := huh.NewForm(fieldsGroup, confirmGroup).
		WithTheme(Theme()).
		WithLayout(huh.LayoutStack)

	if width > 0 {
		// Subtract 1 to compensate for LayoutStack adding an extra character
		form = form.WithWidth(width - 1)
	}

	return form
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
	Title:       lipgloss.Color("#89B4FA"),
	Description: lipgloss.Color("#94E2D5"),
	Error:       lipgloss.Color("#F38BA8"),
	Success:     lipgloss.Color("#A6E3A1"),
	Placeholder: lipgloss.Color("#6C7086"),
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
		Prompt("> ")

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
