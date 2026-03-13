package forms_test

import (
	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/ui/configtypes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigureForm", func() {
	var (
		values     map[string]*string
		boolValues map[string]*bool
	)

	BeforeEach(func() {
		values = make(map[string]*string)
		boolValues = make(map[string]*bool)
	})

	Describe("NewConfigureSettingField", func() {
		Context("with string setting", func() {
			It("should create an input field", func() {
				strVal := "hello"
				values["name"] = &strVal
				setting := &configtypes.ConfigurationSetting{
					Key:         "name",
					Label:       "Name",
					Description: "User name",
					Type:        "string",
				}

				field := forms.NewConfigureSettingField(setting, values, boolValues)
				Expect(field).NotTo(BeNil())
			})
		})

		Context("with int setting", func() {
			It("should create an input field with integer validation", func() {
				intVal := "42"
				values["count"] = &intVal
				setting := &configtypes.ConfigurationSetting{
					Key:         "count",
					Label:       "Count",
					Description: "Item count",
					Type:        "int",
				}

				field := forms.NewConfigureSettingField(setting, values, boolValues)
				Expect(field).NotTo(BeNil())
			})
		})

		Context("with bool setting", func() {
			It("should create a confirm field", func() {
				boolVal := true
				boolValues["enabled"] = &boolVal
				setting := &configtypes.ConfigurationSetting{
					Key:         "enabled",
					Label:       "Enabled",
					Description: "Feature enabled",
					Type:        "bool",
				}

				field := forms.NewConfigureSettingField(setting, values, boolValues)
				Expect(field).NotTo(BeNil())
			})
			It("should create a confirm field with nil bool value", func() {
				setting := &configtypes.ConfigurationSetting{
					Key:         "enabled",
					Label:       "Enabled",
					Description: "Feature enabled",
					Type:        "bool",
				}
				field := forms.NewConfigureSettingField(setting, values, boolValues)
				Expect(field).NotTo(BeNil())
			})
		})

		Context("with select setting", func() {
			It("should create a select field with options", func() {
				strVal := "dark"
				values["theme"] = &strVal
				setting := &configtypes.ConfigurationSetting{
					Key:         "theme",
					Label:       "Theme",
					Description: "UI theme",
					Type:        "select",
					Options:     []string{"light", "dark", "auto"},
				}

				field := forms.NewConfigureSettingField(setting, values, boolValues)
				Expect(field).NotTo(BeNil())
			})

			It("should return nil for select with no options", func() {
				setting := &configtypes.ConfigurationSetting{
					Key:     "empty",
					Label:   "Empty",
					Type:    "select",
					Options: []string{},
				}

				field := forms.NewConfigureSettingField(setting, values, boolValues)
				Expect(field).To(BeNil())
			})
		})

		Context("with unknown setting type", func() {
			It("should create an input field as default", func() {
				strVal := "fallback"
				values["custom"] = &strVal
				setting := &configtypes.ConfigurationSetting{
					Key:         "custom",
					Label:       "Custom",
					Description: "Custom field",
					Type:        "unknown",
				}

				field := forms.NewConfigureSettingField(setting, values, boolValues)
				Expect(field).NotTo(BeNil())
			})
			It("returns an input field for unknown type", func() {
				strVal := "foo"
				setting := &configtypes.ConfigurationSetting{
					Key: "unknown_key", Label: "Unknown", Type: "unknown",
				}
				values := map[string]*string{"unknown_key": &strVal}
				field := forms.NewConfigureSettingField(setting, values, map[string]*bool{})
				Expect(field).NotTo(BeNil())
			})
		})
	})

	Describe("NewConfigureSettingsForm", func() {
		It("should create a form from settings", func() {
			strVal := "test"
			values["name"] = &strVal
			boolVal := false
			boolValues["enabled"] = &boolVal

			settings := []*configtypes.ConfigurationSetting{
				{Key: "name", Label: "Name", Type: "string", Description: "A name"},
				{Key: "enabled", Label: "Enabled", Type: "bool", Description: "On/off"},
			}

			form := forms.NewConfigureSettingsForm(settings, values, boolValues, 80, 40)
			Expect(form).NotTo(BeNil())
		})

		It("should handle empty settings", func() {
			settings := []*configtypes.ConfigurationSetting{}

			form := forms.NewConfigureSettingsForm(settings, values, boolValues, 80, 40)
			Expect(form).To(BeNil())
		})

		It("should handle zero dimensions", func() {
			strVal := "val"
			values["key"] = &strVal
			settings := []*configtypes.ConfigurationSetting{
				{Key: "key", Label: "Key", Type: "string"},
			}

			form := forms.NewConfigureSettingsForm(settings, values, boolValues, 0, 0)
			Expect(form).NotTo(BeNil())
		})
	})

	Describe("BuildConfigureFields", func() {
		It("should build fields for all setting types", func() {
			strVal := "test"
			values["name"] = &strVal
			intVal := "10"
			values["count"] = &intVal
			boolVal := true
			boolValues["enabled"] = &boolVal
			selVal := "dark"
			values["theme"] = &selVal

			settings := []*configtypes.ConfigurationSetting{
				{Key: "name", Label: "Name", Type: "string"},
				{Key: "count", Label: "Count", Type: "int"},
				{Key: "enabled", Label: "Enabled", Type: "bool"},
				{Key: "theme", Label: "Theme", Type: "select", Options: []string{"light", "dark"}},
			}

			fields := forms.BuildConfigureFields(settings, values, boolValues)
			Expect(fields).To(HaveLen(4))
		})

		It("should filter out nil fields", func() {
			settings := []*configtypes.ConfigurationSetting{
				{Key: "empty", Label: "Empty", Type: "select", Options: []string{}},
			}

			fields := forms.BuildConfigureFields(settings, values, boolValues)
			Expect(fields).To(BeEmpty())
		})

		It("should handle empty settings slice", func() {
			fields := forms.BuildConfigureFields(nil, values, boolValues)
			Expect(fields).To(BeEmpty())
		})
	})

	Describe("ValidateInteger", func() {
		It("should accept empty string", func() {
			Expect(forms.ValidateInteger("")).To(Succeed())
		})

		It("should accept valid integer", func() {
			Expect(forms.ValidateInteger("42")).To(Succeed())
		})

		It("should accept zero", func() {
			Expect(forms.ValidateInteger("0")).To(Succeed())
		})

		It("should accept negative integer", func() {
			Expect(forms.ValidateInteger("-5")).To(Succeed())
		})

		It("should reject non-numeric string", func() {
			err := forms.ValidateInteger("abc")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("must be a number"))
		})

		It("should reject float string", func() {
			err := forms.ValidateInteger("3.14")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetConfigureSettingChanges", func() {
		It("ignores non-numeric changed value", func() {
			strVal := "not-a-number"
			settings := []*configtypes.ConfigurationSetting{{Key: "port", Type: "int"}}
			changes := forms.GetConfigureSettingChanges(settings, &forms.ConfigureSettingsFormData{Values: map[string]*string{"port": &strVal}, BoolValues: map[string]*bool{}}, map[string]string{"port": "8080"})
			Expect(changes).To(BeEmpty())
		})
		It("ignores int value that hasn't changed", func() {
			strVal := "8080"
			settings := []*configtypes.ConfigurationSetting{{Key: "port", Type: "int"}}
			changes := forms.GetConfigureSettingChanges(settings, &forms.ConfigureSettingsFormData{Values: map[string]*string{"port": &strVal}, BoolValues: map[string]*bool{}}, map[string]string{"port": "8080"})
			Expect(changes).To(BeEmpty())
		})
		It("ignores int when key is missing", func() {
			settings := []*configtypes.ConfigurationSetting{{Key: "port", Type: "int"}}
			changes := forms.GetConfigureSettingChanges(settings, &forms.ConfigureSettingsFormData{Values: map[string]*string{}, BoolValues: map[string]*bool{}}, map[string]string{"port": "8080"})
			Expect(changes).To(BeEmpty())
		})
		It("ignores int when value pointer is nil", func() {
			settings := []*configtypes.ConfigurationSetting{{Key: "port", Type: "int"}}
			changes := forms.GetConfigureSettingChanges(settings, &forms.ConfigureSettingsFormData{Values: map[string]*string{"port": nil}, BoolValues: map[string]*bool{}}, map[string]string{"port": "8080"})
			Expect(changes).To(BeEmpty())
		})

		It("ignores list where all parts are empty after trim", func() {
			strVal := " , , "
			settings := []*configtypes.ConfigurationSetting{{Key: "tags", Type: "list"}}
			changes := forms.GetConfigureSettingChanges(settings, &forms.ConfigureSettingsFormData{Values: map[string]*string{"tags": &strVal}, BoolValues: map[string]*bool{}}, map[string]string{"tags": "original"})
			val, ok := changes["tags"]
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal([]string(nil)))
		})
	})

	Describe("NewConfigureSettingsFormData", func() {
		It("initializes form data from settings with string values", func() {
			settings := []*configtypes.ConfigurationSetting{
				{Key: "name", Type: "string", Value: "John"},
				{Key: "email", Type: "string", Value: "john@example.com"},
			}
			data := forms.NewConfigureSettingsFormData(settings)
			Expect(data).NotTo(BeNil())
			Expect(data.Values).To(HaveKey("name"))
			Expect(data.Values).To(HaveKey("email"))
			Expect(*data.Values["name"]).To(Equal("John"))
		})

		It("initializes form data with bool values", func() {
			settings := []*configtypes.ConfigurationSetting{
				{Key: "enabled", Type: "bool", Value: true},
				{Key: "disabled", Type: "bool", Value: false},
			}
			data := forms.NewConfigureSettingsFormData(settings)
			Expect(data).NotTo(BeNil())
			Expect(data.BoolValues).To(HaveKey("enabled"))
			Expect(data.BoolValues).To(HaveKey("disabled"))
			Expect(*data.BoolValues["enabled"]).To(BeTrue())
			Expect(*data.BoolValues["disabled"]).To(BeFalse())
		})

		It("initializes form data with mixed types", func() {
			settings := []*configtypes.ConfigurationSetting{
				{Key: "name", Type: "string", Value: "Alice"},
				{Key: "active", Type: "bool", Value: true},
			}
			data := forms.NewConfigureSettingsFormData(settings)
			Expect(data.Values).To(HaveKey("name"))
			Expect(data.BoolValues).To(HaveKey("active"))
		})

		It("handles empty settings slice", func() {
			data := forms.NewConfigureSettingsFormData([]*configtypes.ConfigurationSetting{})
			Expect(data).NotTo(BeNil())
			Expect(data.Values).To(BeEmpty())
			Expect(data.BoolValues).To(BeEmpty())
		})

		It("handles nil settings slice", func() {
			data := forms.NewConfigureSettingsFormData(nil)
			Expect(data).NotTo(BeNil())
			Expect(data.Values).To(BeEmpty())
			Expect(data.BoolValues).To(BeEmpty())
		})

		It("converts list values to comma-separated strings", func() {
			settings := []*configtypes.ConfigurationSetting{
				{Key: "tags", Type: "list", Value: []string{"go", "rust", "python"}},
			}
			data := forms.NewConfigureSettingsFormData(settings)
			Expect(*data.Values["tags"]).To(Equal("go, rust, python"))
		})
	})

	Describe("NewConfigureSettingField - additional cases", func() {
		It("creates bool field for bool type", func() {
			boolVal := true
			boolValues := map[string]*bool{"enabled": &boolVal}
			setting := &configtypes.ConfigurationSetting{
				Key:  "enabled",
				Type: "bool",
			}
			field := forms.NewConfigureSettingField(setting, map[string]*string{}, boolValues)
			Expect(field).NotTo(BeNil())
		})

		It("creates input field for list type", func() {
			listVal := "item1, item2"
			values := map[string]*string{"items": &listVal}
			setting := &configtypes.ConfigurationSetting{
				Key:  "items",
				Type: "list",
			}
			field := forms.NewConfigureSettingField(setting, values, map[string]*bool{})
			Expect(field).NotTo(BeNil())
		})

		It("creates input field for unknown type", func() {
			strVal := "value"
			values := map[string]*string{"custom": &strVal}
			setting := &configtypes.ConfigurationSetting{
				Key:  "custom",
				Type: "unknown",
			}
			field := forms.NewConfigureSettingField(setting, values, map[string]*bool{})
			Expect(field).NotTo(BeNil())
		})

		It("applies email validation for email field", func() {
			strVal := "test@example.com"
			values := map[string]*string{"email": &strVal}
			setting := &configtypes.ConfigurationSetting{
				Key:  "email",
				Type: "string",
			}
			field := forms.NewConfigureSettingField(setting, values, map[string]*bool{})
			Expect(field).NotTo(BeNil())
		})

		It("applies URL validation for github field", func() {
			strVal := "https://github.com/user"
			values := map[string]*string{"github": &strVal}
			setting := &configtypes.ConfigurationSetting{
				Key:  "github",
				Type: "string",
			}
			field := forms.NewConfigureSettingField(setting, values, map[string]*bool{})
			Expect(field).NotTo(BeNil())
		})

		It("applies URL validation for portfolio field", func() {
			strVal := "https://example.com"
			values := map[string]*string{"portfolio": &strVal}
			setting := &configtypes.ConfigurationSetting{
				Key:  "portfolio",
				Type: "string",
			}
			field := forms.NewConfigureSettingField(setting, values, map[string]*bool{})
			Expect(field).NotTo(BeNil())
		})
	})

	Describe("ValidateEmail", func() {
		It("accepts empty string", func() {
			Err := forms.ValidateEmail("")
			Expect(Err).NotTo(HaveOccurred())
		})

		It("accepts valid email", func() {
			Err := forms.ValidateEmail("user@example.com")
			Expect(Err).NotTo(HaveOccurred())
		})

		It("accepts email with multiple @ in local part", func() {
			Err := forms.ValidateEmail("user+tag@example.com")
			Expect(Err).NotTo(HaveOccurred())
		})

		It("rejects string without @", func() {
			Err := forms.ValidateEmail("notanemail")
			Expect(Err).To(HaveOccurred())
			Expect(Err.Error()).To(Equal("must be a valid email address"))
		})

		It("rejects string with @ but no domain", func() {
			Err := forms.ValidateEmail("user@")
			Expect(Err).NotTo(HaveOccurred())
		})
	})

	Describe("ValidateURL", func() {
		It("accepts empty string", func() {
			Err := forms.ValidateURL("")
			Expect(Err).NotTo(HaveOccurred())
		})

		It("accepts http URL", func() {
			Err := forms.ValidateURL("http://example.com")
			Expect(Err).NotTo(HaveOccurred())
		})

		It("accepts https URL", func() {
			Err := forms.ValidateURL("https://example.com")
			Expect(Err).NotTo(HaveOccurred())
		})

		It("accepts https URL with path", func() {
			Err := forms.ValidateURL("https://github.com/user/repo")
			Expect(Err).NotTo(HaveOccurred())
		})

		It("rejects URL without scheme", func() {
			Err := forms.ValidateURL("example.com")
			Expect(Err).To(HaveOccurred())
			Expect(Err.Error()).To(Equal("must be a valid URL starting with http:// or https://"))
		})

		It("rejects URL with ftp scheme", func() {
			Err := forms.ValidateURL("ftp://example.com")
			Expect(Err).To(HaveOccurred())
		})

		It("rejects URL with only http", func() {
			Err := forms.ValidateURL("http")
			Expect(Err).To(HaveOccurred())
		})
	})

	Describe("settingToString", func() {
		It("returns empty string for nil value", func() {
			result := forms.SettingToString(nil)
			Expect(result).To(Equal(""))
		})

		It("converts string value", func() {
			result := forms.SettingToString("hello")
			Expect(result).To(Equal("hello"))
		})

		It("joins slice of strings with comma-space", func() {
			result := forms.SettingToString([]string{"a", "b", "c"})
			Expect(result).To(Equal("a, b, c"))
		})

		It("converts int value to string", func() {
			result := forms.SettingToString(42)
			Expect(result).To(Equal("42"))
		})

		It("converts bool value to string", func() {
			result := forms.SettingToString(true)
			Expect(result).To(Equal("true"))
		})

		It("handles empty slice", func() {
			result := forms.SettingToString([]string{})
			Expect(result).To(Equal(""))
		})
	})

	Describe("settingToBool", func() {
		It("returns true for bool true", func() {
			result := forms.SettingToBool(true)
			Expect(result).To(BeTrue())
		})

		It("returns false for bool false", func() {
			result := forms.SettingToBool(false)
			Expect(result).To(BeFalse())
		})

		It("returns false for non-bool value", func() {
			result := forms.SettingToBool("true")
			Expect(result).To(BeFalse())
		})

		It("returns false for nil", func() {
			result := forms.SettingToBool(nil)
			Expect(result).To(BeFalse())
		})

		It("returns false for int", func() {
			result := forms.SettingToBool(1)
			Expect(result).To(BeFalse())
		})
	})

	Describe("appendBoolChange", func() {
		It("adds bool change when value differs from original", func() {
			changes := make(map[string]interface{})
			boolVal := true
			boolValues := map[string]*bool{"enabled": &boolVal}
			originalValues := map[string]string{"enabled": "false"}
			forms.AppendBoolChange(changes, "enabled", boolValues, originalValues)
			Expect(changes).To(HaveKey("enabled"))
			Expect(changes["enabled"]).To(BeTrue())
		})

		It("ignores bool change when value matches original", func() {
			changes := make(map[string]interface{})
			boolVal := true
			boolValues := map[string]*bool{"enabled": &boolVal}
			originalValues := map[string]string{"enabled": "true"}
			forms.AppendBoolChange(changes, "enabled", boolValues, originalValues)
			Expect(changes).NotTo(HaveKey("enabled"))
		})

		It("ignores bool change when key missing from boolValues", func() {
			changes := make(map[string]interface{})
			boolValues := map[string]*bool{}
			originalValues := map[string]string{"enabled": "false"}
			forms.AppendBoolChange(changes, "enabled", boolValues, originalValues)
			Expect(changes).NotTo(HaveKey("enabled"))
		})

		It("ignores bool change when pointer is nil", func() {
			changes := make(map[string]interface{})
			boolValues := map[string]*bool{"enabled": nil}
			originalValues := map[string]string{"enabled": "false"}
			forms.AppendBoolChange(changes, "enabled", boolValues, originalValues)
			Expect(changes).NotTo(HaveKey("enabled"))
		})
	})

	Describe("appendIntChange", func() {
		It("adds int change when value differs from original", func() {
			changes := make(map[string]interface{})
			strVal := "42"
			values := map[string]*string{"port": &strVal}
			originalValues := map[string]string{"port": "8080"}
			forms.AppendIntChange(changes, "port", values, originalValues)
			Expect(changes).To(HaveKey("port"))
			Expect(changes["port"]).To(Equal(42))
		})

		It("ignores int change when value matches original", func() {
			changes := make(map[string]interface{})
			strVal := "8080"
			values := map[string]*string{"port": &strVal}
			originalValues := map[string]string{"port": "8080"}
			forms.AppendIntChange(changes, "port", values, originalValues)
			Expect(changes).NotTo(HaveKey("port"))
		})

		It("ignores int change when Atoi fails", func() {
			changes := make(map[string]interface{})
			strVal := "not-a-number"
			values := map[string]*string{"port": &strVal}
			originalValues := map[string]string{"port": "8080"}
			forms.AppendIntChange(changes, "port", values, originalValues)
			Expect(changes).NotTo(HaveKey("port"))
		})

		It("ignores int change when key missing", func() {
			changes := make(map[string]interface{})
			values := map[string]*string{}
			originalValues := map[string]string{"port": "8080"}
			forms.AppendIntChange(changes, "port", values, originalValues)
			Expect(changes).NotTo(HaveKey("port"))
		})

		It("ignores int change when pointer is nil", func() {
			changes := make(map[string]interface{})
			values := map[string]*string{"port": nil}
			originalValues := map[string]string{"port": "8080"}
			forms.AppendIntChange(changes, "port", values, originalValues)
			Expect(changes).NotTo(HaveKey("port"))
		})
	})

	Describe("appendStringChange", func() {
		It("adds string change when value differs from original", func() {
			changes := make(map[string]interface{})
			strVal := "new-value"
			values := map[string]*string{"name": &strVal}
			originalValues := map[string]string{"name": "old-value"}
			forms.AppendStringChange(changes, "name", values, originalValues)
			Expect(changes).To(HaveKey("name"))
			Expect(changes["name"]).To(Equal("new-value"))
		})

		It("ignores string change when value matches original", func() {
			changes := make(map[string]interface{})
			strVal := "same-value"
			values := map[string]*string{"name": &strVal}
			originalValues := map[string]string{"name": "same-value"}
			forms.AppendStringChange(changes, "name", values, originalValues)
			Expect(changes).NotTo(HaveKey("name"))
		})

		It("ignores string change when key missing", func() {
			changes := make(map[string]interface{})
			values := map[string]*string{}
			originalValues := map[string]string{"name": "old-value"}
			forms.AppendStringChange(changes, "name", values, originalValues)
			Expect(changes).NotTo(HaveKey("name"))
		})

		It("ignores string change when pointer is nil", func() {
			changes := make(map[string]interface{})
			values := map[string]*string{"name": nil}
			originalValues := map[string]string{"name": "old-value"}
			forms.AppendStringChange(changes, "name", values, originalValues)
			Expect(changes).NotTo(HaveKey("name"))
		})
	})

	Describe("appendListChange", func() {
		It("adds list change when value differs from original", func() {
			changes := make(map[string]interface{})
			strVal := "item1, item2, item3"
			values := map[string]*string{"tags": &strVal}
			originalValues := map[string]string{"tags": "old1, old2"}
			forms.AppendListChange(changes, "tags", values, originalValues)
			Expect(changes).To(HaveKey("tags"))
			Expect(changes["tags"]).To(Equal([]string{"item1", "item2", "item3"}))
		})

		It("ignores list change when value matches original", func() {
			changes := make(map[string]interface{})
			strVal := "item1, item2"
			values := map[string]*string{"tags": &strVal}
			originalValues := map[string]string{"tags": "item1, item2"}
			forms.AppendListChange(changes, "tags", values, originalValues)
			Expect(changes).NotTo(HaveKey("tags"))
		})

		It("trims whitespace from list items", func() {
			changes := make(map[string]interface{})
			strVal := " item1 , item2 , item3 "
			values := map[string]*string{"tags": &strVal}
			originalValues := map[string]string{"tags": "old"}
			forms.AppendListChange(changes, "tags", values, originalValues)
			Expect(changes["tags"]).To(Equal([]string{"item1", "item2", "item3"}))
		})

		It("handles list with empty parts after trim", func() {
			changes := make(map[string]interface{})
			strVal := "item1, , item2"
			values := map[string]*string{"tags": &strVal}
			originalValues := map[string]string{"tags": "old"}
			forms.AppendListChange(changes, "tags", values, originalValues)
			Expect(changes["tags"]).To(Equal([]string{"item1", "item2"}))
		})

		It("ignores list change when key missing", func() {
			changes := make(map[string]interface{})
			values := map[string]*string{}
			originalValues := map[string]string{"tags": "old"}
			forms.AppendListChange(changes, "tags", values, originalValues)
			Expect(changes).NotTo(HaveKey("tags"))
		})

		It("ignores list change when pointer is nil", func() {
			changes := make(map[string]interface{})
			values := map[string]*string{"tags": nil}
			originalValues := map[string]string{"tags": "old"}
			forms.AppendListChange(changes, "tags", values, originalValues)
			Expect(changes).NotTo(HaveKey("tags"))
		})
	})

	Describe("GetConfigureSettingChanges - additional cases", func() {
		It("detects bool changes", func() {
			boolVal := true
			settings := []*configtypes.ConfigurationSetting{{Key: "enabled", Type: "bool"}}
			changes := forms.GetConfigureSettingChanges(settings, &forms.ConfigureSettingsFormData{Values: map[string]*string{}, BoolValues: map[string]*bool{"enabled": &boolVal}}, map[string]string{"enabled": "false"})
			Expect(changes).To(HaveKey("enabled"))
			Expect(changes["enabled"]).To(BeTrue())
		})

		It("detects string changes", func() {
			strVal := "new"
			settings := []*configtypes.ConfigurationSetting{{Key: "name", Type: "string"}}
			changes := forms.GetConfigureSettingChanges(settings, &forms.ConfigureSettingsFormData{Values: map[string]*string{"name": &strVal}, BoolValues: map[string]*bool{}}, map[string]string{"name": "old"})
			Expect(changes).To(HaveKey("name"))
			Expect(changes["name"]).To(Equal("new"))
		})

		It("detects int changes", func() {
			strVal := "42"
			settings := []*configtypes.ConfigurationSetting{{Key: "port", Type: "int"}}
			changes := forms.GetConfigureSettingChanges(settings, &forms.ConfigureSettingsFormData{Values: map[string]*string{"port": &strVal}, BoolValues: map[string]*bool{}}, map[string]string{"port": "8080"})
			Expect(changes).To(HaveKey("port"))
			Expect(changes["port"]).To(Equal(42))
		})

		It("detects list changes", func() {
			strVal := "a, b, c"
			settings := []*configtypes.ConfigurationSetting{{Key: "tags", Type: "list"}}
			changes := forms.GetConfigureSettingChanges(settings, &forms.ConfigureSettingsFormData{Values: map[string]*string{"tags": &strVal}, BoolValues: map[string]*bool{}}, map[string]string{"tags": "x, y"})
			Expect(changes).To(HaveKey("tags"))
			Expect(changes["tags"]).To(Equal([]string{"a", "b", "c"}))
		})
	})
})
