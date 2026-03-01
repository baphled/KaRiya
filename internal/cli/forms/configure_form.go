package forms

import (
	"errors"
	"strconv"

	"github.com/baphled/kariya/internal/cli/configtypes"
)

// NewConfigureSettingField creates the appropriate form field for a configuration setting based on its type.
//
// Expected: setting must be non-nil with a valid Type field; values and boolValues must contain the setting key.
//
// Returns: the appropriate Field for the setting type, or nil for select with no options.
//
// Side effects: none.
func NewConfigureSettingField(
	setting *configtypes.ConfigurationSetting,
	values map[string]*string,
	boolValues map[string]*bool,
) Field {
	switch setting.Type {
	case "string":
		return NewInput(FieldConfig{
			Key:         setting.Key,
			Title:       setting.Label,
			Description: setting.Description,
		}).Value(values[setting.Key])

	case "int":
		return NewInput(FieldConfig{
			Key:         setting.Key,
			Title:       setting.Label,
			Description: setting.Description,
			Validate:    ValidateInteger,
		}).Value(values[setting.Key])

	case "bool":
		return NewConfirm(
			setting.Key,
			setting.Label,
			setting.Description,
			"Yes",
			"No",
		).Value(boolValues[setting.Key])

	case "select":
		if len(setting.Options) == 0 {
			return nil
		}
		options := buildSettingOptions(setting.Options)
		return NewSelect(setting.Key, setting.Label, setting.Description, options).
			Value(values[setting.Key])

	default:
		return NewInput(FieldConfig{
			Key:         setting.Key,
			Title:       setting.Label,
			Description: setting.Description,
		}).Value(values[setting.Key])
	}
}

// NewConfigureSettingsForm builds a complete form for all configuration settings with the given dimensions.
//
// Expected: settings must be a valid slice; values and boolValues must contain keys for each setting.
//
// Returns: a fully initialized Form, or nil if no valid fields are produced.
//
// Side effects: none.
func NewConfigureSettingsForm(
	settings []*configtypes.ConfigurationSetting,
	values map[string]*string,
	boolValues map[string]*bool,
	width, height int,
) Form {
	fields := BuildConfigureFields(settings, values, boolValues)
	if len(fields) == 0 {
		return nil
	}
	group := NewGroup(fields...)
	return NewFormWithDimensions(width, height, group)
}

// BuildConfigureFields creates form fields for a slice of configuration settings, filtering out nil results.
//
// Expected: settings must be a valid slice; values and boolValues must contain keys for each setting.
//
// Returns: a slice of non-nil Field values, one per valid setting.
//
// Side effects: none.
func BuildConfigureFields(
	settings []*configtypes.ConfigurationSetting,
	values map[string]*string,
	boolValues map[string]*bool,
) []Field {
	var fields []Field
	for _, setting := range settings {
		field := NewConfigureSettingField(setting, values, boolValues)
		if field != nil {
			fields = append(fields, field)
		}
	}
	return fields
}

func buildSettingOptions(opts []string) []SelectOption {
	options := make([]SelectOption, len(opts))
	for i, opt := range opts {
		options[i] = SelectOption{Key: opt, Value: opt}
	}
	return options
}

// ValidateInteger checks that a string value is either empty or a valid integer.
//
// Expected: val is any string value.
//
// Returns: nil if val is empty or a valid integer, otherwise an error.
//
// Side effects: none.
func ValidateInteger(val string) error {
	if val == "" {
		return nil
	}
	_, err := strconv.Atoi(val)
	if err != nil {
		return errors.New("must be a number")
	}
	return nil
}
