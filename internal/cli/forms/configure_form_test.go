package forms_test

import (
	"github.com/baphled/kariya/internal/cli/configtypes"
	"github.com/baphled/kariya/internal/cli/forms"
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
})
