package forms

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/baphled/kariya/internal/ui/configtypes"
)

// ConfigureSettingsFormData holds bound form values for configuration settings.
type ConfigureSettingsFormData struct {
	Values          map[string]*string
	BoolValues      map[string]*bool
	SubmitConfirmed bool
}

// NewConfigureSettingsFormData initialises form data from a slice of settings.
//
// Expected:
//   - settings must be a valid slice of ConfigurationSetting pointers.
//
// Returns:
//   - A fully initialised ConfigureSettingsFormData ready for form binding.
//
// Side effects:
//   - None.
func NewConfigureSettingsFormData(settings []*configtypes.ConfigurationSetting) *ConfigureSettingsFormData {
	data := &ConfigureSettingsFormData{
		Values:     make(map[string]*string),
		BoolValues: make(map[string]*bool),
	}
	for _, s := range settings {
		switch s.Type {
		case "bool":
			boolVal := SettingToBool(s.Value)
			data.BoolValues[s.Key] = &boolVal
		default:
			strVal := SettingToString(s.Value)
			data.Values[s.Key] = &strVal
		}
	}
	return data
}

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
		cfg := FieldConfig{
			Key:         setting.Key,
			Title:       setting.Label,
			Description: setting.Description,
		}
		switch setting.Key {
		case "email":
			cfg.Validate = ValidateEmail
		case "github", "portfolio":
			cfg.Validate = ValidateURL
		}
		return NewInput(cfg).Value(values[setting.Key])

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

// ValidateEmail checks that a string value is either empty or contains an "@" character.
//
// Expected: val is any string value.
//
// Returns: nil if val is empty or contains "@", otherwise an error.
//
// Side effects: none.
func ValidateEmail(val string) error {
	if val == "" {
		return nil
	}
	if !strings.Contains(val, "@") {
		return errors.New("must be a valid email address")
	}
	return nil
}

// ValidateURL checks that a string value is either empty or starts with "http://" or "https://".
//
// Expected: val is any string value.
//
// Returns: nil if val is empty or starts with a valid URL scheme, otherwise an error.
//
// Side effects: none.
func ValidateURL(val string) error {
	if val == "" {
		return nil
	}
	if !strings.HasPrefix(val, "http://") && !strings.HasPrefix(val, "https://") {
		return errors.New("must be a valid URL starting with http:// or https://")
	}
	return nil
}

// GetConfigureSettingChanges returns only the settings that changed from their original values.
//
// Expected:
//   - settings must be a valid slice of ConfigurationSetting pointers.
//   - data must be a valid ConfigureSettingsFormData pointer.
//   - originalValues maps setting keys to their original string representations.
//
// Returns:
//   - A map of changed keys to their type-converted new values.
//
// Side effects:
//   - None.
func GetConfigureSettingChanges(
	settings []*configtypes.ConfigurationSetting,
	data *ConfigureSettingsFormData,
	originalValues map[string]string,
) map[string]interface{} {
	changes := make(map[string]interface{})
	for _, s := range settings {
		switch s.Type {
		case "bool":
			AppendBoolChange(changes, s.Key, data.BoolValues, originalValues)
		case "int":
			AppendIntChange(changes, s.Key, data.Values, originalValues)
		case "list":
			AppendListChange(changes, s.Key, data.Values, originalValues)
		default:
			AppendStringChange(changes, s.Key, data.Values, originalValues)
		}
	}
	return changes
}

// SettingToString converts a configuration setting value to its string representation.
//
// Expected: val must be a valid configuration value (string, bool, int, or []string).
//
// Returns: The string representation of val.
//
// Side effects: None.
func SettingToString(val interface{}) string {
	if val == nil {
		return ""
	}
	if sl, ok := val.([]string); ok {
		return strings.Join(sl, ", ")
	}
	return fmt.Sprintf("%v", val)
}

// SettingToBool converts a configuration setting value to a boolean.
//
// Expected: val must be a valid configuration value.
//
// Returns: The boolean representation of val, or false if not a bool.
//
// Side effects: None.
func SettingToBool(val interface{}) bool {
	b, ok := val.(bool)
	return ok && b
}

// AppendBoolChange appends a bool setting change to changes if the value differs from the original.
//
// Expected: changes, boolValues, and originalValues must be non-nil maps.
//
// Returns: Nothing.
//
// Side effects: May add an entry to changes.
func AppendBoolChange(changes map[string]interface{}, key string, boolValues map[string]*bool, originalValues map[string]string) {
	boolPtr, ok := boolValues[key]
	if !ok || boolPtr == nil {
		return
	}
	current := strconv.FormatBool(*boolPtr)
	if current != originalValues[key] {
		changes[key] = *boolPtr
	}
}

// AppendIntChange appends an int setting change to changes if the value differs from the original.
//
// Expected: changes, values, and originalValues must be non-nil maps.
//
// Returns: Nothing.
//
// Side effects: May add an entry to changes.
func AppendIntChange(changes map[string]interface{}, key string, values map[string]*string, originalValues map[string]string) {
	strPtr, ok := values[key]
	if !ok || strPtr == nil {
		return
	}
	if *strPtr == originalValues[key] {
		return
	}
	iv, err := strconv.Atoi(*strPtr)
	if err != nil {
		return
	}
	changes[key] = iv
}

// AppendStringChange appends a string setting change to changes if the value differs from the original.
//
// Expected: changes, values, and originalValues must be non-nil maps.
//
// Returns: Nothing.
//
// Side effects: May add an entry to changes.
func AppendStringChange(changes map[string]interface{}, key string, values map[string]*string, originalValues map[string]string) {
	strPtr, ok := values[key]
	if !ok || strPtr == nil {
		return
	}
	if *strPtr != originalValues[key] {
		changes[key] = *strPtr
	}
}

// AppendListChange appends a list setting change to changes if the value differs from the original.
//
// Expected: changes, values, and originalValues must be non-nil maps.
//
// Returns: Nothing.
//
// Side effects: May add an entry to changes.
func AppendListChange(changes map[string]interface{}, key string, values map[string]*string, originalValues map[string]string) {
	strPtr, ok := values[key]
	if !ok || strPtr == nil {
		return
	}
	if *strPtr == originalValues[key] {
		return
	}
	parts := strings.Split(*strPtr, ",")
	var result []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	changes[key] = result
}
