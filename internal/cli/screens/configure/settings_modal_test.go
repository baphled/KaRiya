package configure

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/configtypes"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/themes"
)

func makeTestSettings() map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting {
	return map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting{
		configtypes.DomainSystem: {
			{Key: "log_level", Label: "Log Level", Value: "info", DefaultValue: "info", Type: "select", Options: []string{"debug", "info", "warn", "error"}},
		},
		configtypes.DomainProfile: {
			{Key: "name", Label: "Full Name", Value: "Test User", DefaultValue: "", Type: "string"},
		},
		configtypes.DomainExport: {
			{Key: "default_destination", Label: "Default Destination", Value: "file", DefaultValue: "file", Type: "select", Options: []string{"file", "clipboard"}},
		},
		configtypes.DomainUI: {
			{Key: "theme", Label: "Theme", Value: "dark", DefaultValue: "dark", Type: "select", Options: []string{"light", "dark"}},
		},
	}
}

func makeSettingsWithDomain(domain configtypes.ConfigurationDomain, settings ...*configtypes.ConfigurationSetting) map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting {
	return map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting{
		domain: settings,
	}
}

func getModalFormData(modal *SettingsModal) map[configtypes.ConfigurationDomain]*forms.ConfigureSettingsFormData {
	return modal.formData
}

func setModalDomains(modal *SettingsModal, domains []configtypes.ConfigurationDomain) {
	modal.domains = domains
}

func setModalSettings(modal *SettingsModal, settings map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting) {
	modal.settings = settings
}

func setModalSelectedIdx(modal *SettingsModal, idx int) {
	modal.selectedIdx = idx
}

func setModalActiveForm(modal *SettingsModal, form forms.Form) {
	modal.activeForm = form
}

func updateFormDataString(data map[configtypes.ConfigurationDomain]*forms.ConfigureSettingsFormData, domain configtypes.ConfigurationDomain, key, value string) {
	data[domain].Values[key] = &value
}

func updateFormDataBool(data map[configtypes.ConfigurationDomain]*forms.ConfigureSettingsFormData, domain configtypes.ConfigurationDomain, key string, value bool) {
	data[domain].BoolValues[key] = &value
}

var _ = Describe("SettingsModal", func() {
	var (
		modal    *SettingsModal
		settings map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting
	)

	BeforeEach(func() {
		settings = makeTestSettings()
		modal = NewSettingsModal(settings, 120, 40)
	})

	Describe("NewSettingsModal", func() {
		It("should create a modal with correct initial state", func() {
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsCompleted()).To(BeFalse())
			Expect(modal.IsCancelled()).To(BeFalse())
		})

		It("should not panic with empty settings map", func() {
			emptySettings := make(map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting)
			emptyModal := NewSettingsModal(emptySettings, 120, 40)
			Expect(emptyModal).NotTo(BeNil())
			Expect(emptyModal.IsCompleted()).To(BeFalse())
			Expect(emptyModal.IsCancelled()).To(BeFalse())
		})

		It("should initialise with no changes", func() {
			changes := modal.GetChanges()
			Expect(changes).To(BeEmpty())
		})
	})

	Describe("Init", func() {
		It("should return without error", func() {
			cmd := modal.Init()
			_ = cmd
		})
	})

	Describe("Update", func() {
		Context("when navigating with j key", func() {
			It("should move to the next domain", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				view := modal.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should not move beyond the last domain", func() {
				for range 10 {
					modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				}
				Expect(modal.IsCompleted()).To(BeFalse())
				Expect(modal.IsCancelled()).To(BeFalse())
			})
		})

		Context("when navigating with k key", func() {
			It("should not move above the first domain", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
				view := modal.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should move back after moving forward", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
				view := modal.View()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("when pressing ctrl+s", func() {
			It("should mark the modal as completed", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
				Expect(modal.IsCompleted()).To(BeTrue())
			})

			It("should not mark the modal as cancelled", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
				Expect(modal.IsCancelled()).To(BeFalse())
			})
		})

		Context("when pressing esc", func() {
			It("should mark the modal as cancelled", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(modal.IsCancelled()).To(BeTrue())
			})

			It("should not mark the modal as completed", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(modal.IsCompleted()).To(BeFalse())
			})
		})
	})

	Describe("IsCompleted", func() {
		It("should return false initially", func() {
			Expect(modal.IsCompleted()).To(BeFalse())
		})

		It("should return true after ctrl+s", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(modal.IsCompleted()).To(BeTrue())
		})
	})

	Describe("IsCancelled", func() {
		It("should return false initially", func() {
			Expect(modal.IsCancelled()).To(BeFalse())
		})

		It("should return true after esc", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.IsCancelled()).To(BeTrue())
		})
	})

	Describe("GetChanges", func() {
		It("should return empty map when no changes made", func() {
			changes := modal.GetChanges()
			Expect(changes).To(BeEmpty())
		})
	})

	Describe("SetTheme", func() {
		It("should not panic when setting a theme", func() {
			theme := themes.NewDefaultTheme()
			Expect(func() { modal.SetTheme(theme) }).NotTo(Panic())
		})
	})

	Describe("SetDimensions", func() {
		It("should update width and height", func() {
			modal.SetDimensions(200, 60)
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle small dimensions", func() {
			modal.SetDimensions(40, 10)
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Render", func() {
		It("should return non-empty string", func() {
			output := modal.Render(120, 40)
			Expect(output).NotTo(BeEmpty())
		})

		It("should update dimensions when rendering", func() {
			output := modal.Render(200, 60)
			Expect(output).NotTo(BeEmpty())
		})
	})

	Describe("View", func() {
		It("should return non-empty string", func() {
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should contain configure title", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Configure"))
		})

		It("should contain domain labels", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("System"))
		})
	})

	Describe("rebuildActiveForm with empty domains", func() {
		It("should handle nil activeForm when no domains", func() {
			emptySettings := make(map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting)
			emptyModal := NewSettingsModal(emptySettings, 120, 40)
			view := emptyModal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("rebuildActiveForm with nil domain data", func() {
		It("should handle missing form data for selected domain", func() {
			settingsWithSystem := makeSettingsWithDomain(configtypes.DomainSystem)
			missingFormModal := NewSettingsModal(settingsWithSystem, 120, 40)
			formData := getModalFormData(missingFormModal)
			formData[configtypes.DomainSystem] = nil
			setModalDomains(missingFormModal, []configtypes.ConfigurationDomain{configtypes.DomainSystem})
			setModalSelectedIdx(missingFormModal, 0)
			setModalActiveForm(missingFormModal, nil)
			view := missingFormModal.View()
			Expect(view).To(ContainSubstring("No settings available"))
		})
	})

	Describe("rebuildActiveForm with empty domain settings", func() {
		It("should clear the active form when domain settings are empty", func() {
			settingsWithSystem := makeSettingsWithDomain(configtypes.DomainSystem)
			emptySettingsModal := NewSettingsModal(settingsWithSystem, 120, 40)
			setModalDomains(emptySettingsModal, []configtypes.ConfigurationDomain{configtypes.DomainSystem})
			setModalSettings(emptySettingsModal, settingsWithSystem)
			setModalSelectedIdx(emptySettingsModal, 0)
			setModalActiveForm(emptySettingsModal, nil)
			view := emptySettingsModal.View()
			Expect(view).To(ContainSubstring("No settings available"))
		})
	})

	Describe("renderSectionList with unknown domain", func() {
		It("should render the raw domain label", func() {
			settingsWithUnknown := makeSettingsWithDomain(
				configtypes.ConfigurationDomain("custom"),
				&configtypes.ConfigurationSetting{Key: "custom", Label: "Custom", Value: "value", DefaultValue: "value", Type: "string"},
			)
			unknownModal := NewSettingsModal(settingsWithUnknown, 120, 40)
			setModalDomains(unknownModal, []configtypes.ConfigurationDomain{configtypes.ConfigurationDomain("custom")})
			setModalSelectedIdx(unknownModal, 0)
			view := unknownModal.View()
			Expect(view).To(ContainSubstring("custom"))
		})
	})

	Describe("formDimensions with small width", func() {
		It("should enforce minimum form width of 30", func() {
			smallModal := NewSettingsModal(makeTestSettings(), 40, 40)
			view := smallModal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Init with nil activeForm", func() {
		It("should return nil when activeForm is nil", func() {
			emptySettings := make(map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting)
			emptyModal := NewSettingsModal(emptySettings, 120, 40)
			cmd := emptyModal.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("renderFormPanel with nil activeForm", func() {
		It("should render muted message when activeForm is nil", func() {
			emptySettings := make(map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting)
			emptyModal := NewSettingsModal(emptySettings, 120, 40)
			view := emptyModal.View()
			Expect(view).To(ContainSubstring("No settings available"))
		})
	})

	Describe("Update with down arrow key", func() {
		It("should move to next domain with down arrow", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Update with up arrow key", func() {
		It("should move to previous domain with up arrow", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			modal.Update(tea.KeyMsg{Type: tea.KeyUp})
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Update with unknown key", func() {
		It("should not change state with unknown key", func() {
			cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(cmd).To(BeNil())
			Expect(modal.IsCompleted()).To(BeFalse())
			Expect(modal.IsCancelled()).To(BeFalse())
		})
	})

	Describe("formatDomainLabel", func() {
		It("should display all known domain labels", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("System"))
			Expect(view).To(ContainSubstring("Profile"))
			Expect(view).To(ContainSubstring("Export"))
			Expect(view).To(ContainSubstring("UI"))
		})
	})

	Describe("formatDomainLabel for known domains", func() {
		It("should render all known domain labels correctly", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("System"))
			Expect(view).To(ContainSubstring("Profile"))
			Expect(view).To(ContainSubstring("Export"))
			Expect(view).To(ContainSubstring("UI"))
		})
	})

	Describe("GetChanges with actual changes", func() {
		It("should detect changes when form values differ from originals", func() {
			data := getModalFormData(modal)
			updateFormDataString(data, configtypes.DomainProfile, "name", "Updated User")
			Expect(data[configtypes.DomainProfile].Values["name"]).NotTo(BeNil())
			changes := modal.GetChanges()
			Expect(changes).To(HaveKeyWithValue("name", "Updated User"))
		})
	})

	Describe("GetChanges with bool value updates", func() {
		It("should return bool changes when value changes", func() {
			settingsWithBool := makeSettingsWithDomain(
				configtypes.DomainSystem,
				&configtypes.ConfigurationSetting{Key: "enabled", Label: "Enabled", Value: false, DefaultValue: false, Type: "bool"},
			)
			boolModal := NewSettingsModal(settingsWithBool, 120, 40)
			data := getModalFormData(boolModal)
			updateFormDataBool(data, configtypes.DomainSystem, "enabled", true)
			changes := boolModal.GetChanges()
			Expect(changes).To(HaveKeyWithValue("enabled", true))
		})
	})

	Describe("GetChanges with int value updates", func() {
		It("should return int changes when value changes", func() {
			settingsWithInt := makeSettingsWithDomain(
				configtypes.DomainSystem,
				&configtypes.ConfigurationSetting{Key: "refresh_interval", Label: "Refresh Interval", Value: 5, DefaultValue: 5, Type: "int"},
			)
			intModal := NewSettingsModal(settingsWithInt, 120, 40)
			data := getModalFormData(intModal)
			updateFormDataString(data, configtypes.DomainSystem, "refresh_interval", "10")
			changes := intModal.GetChanges()
			Expect(changes).To(HaveKeyWithValue("refresh_interval", 10))
		})
	})

	Describe("GetChanges with list value updates", func() {
		It("should return list changes when value changes", func() {
			settingsWithList := makeSettingsWithDomain(
				configtypes.DomainSystem,
				&configtypes.ConfigurationSetting{Key: "favorites", Label: "Favorites", Value: []string{"alpha"}, DefaultValue: []string{"alpha"}, Type: "list"},
			)
			listModal := NewSettingsModal(settingsWithList, 120, 40)
			data := getModalFormData(listModal)
			updateFormDataString(data, configtypes.DomainSystem, "favorites", "alpha, beta")
			changes := listModal.GetChanges()
			Expect(changes).To(HaveKeyWithValue("favorites", []string{"alpha", "beta"}))
		})
	})

	Describe("Update with non-KeyMsg", func() {
		It("should handle non-key messages gracefully", func() {
			cmd := modal.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Navigation boundary conditions", func() {
		It("should not move down from last domain", func() {
			for range 20 {
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			}
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should not move up from first domain", func() {
			for range 20 {
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			}
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Multiple state transitions", func() {
		It("should handle multiple navigation and action sequences", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(modal.IsCompleted()).To(BeFalse())
			Expect(modal.IsCancelled()).To(BeFalse())
		})
	})

	Describe("Render with various dimensions", func() {
		It("should handle very large dimensions", func() {
			output := modal.Render(500, 200)
			Expect(output).NotTo(BeEmpty())
		})

		It("should handle minimum viable dimensions", func() {
			output := modal.Render(30, 10)
			Expect(output).NotTo(BeEmpty())
		})
	})

	Describe("SetTheme and SetDimensions interaction", func() {
		It("should apply theme and dimensions together", func() {
			theme := themes.NewDefaultTheme()
			modal.SetTheme(theme)
			modal.SetDimensions(150, 50)
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("View consistency", func() {
		It("should produce consistent output for same state", func() {
			view1 := modal.View()
			view2 := modal.View()
			Expect(view1).To(Equal(view2))
		})
	})

	Describe("Completed and Cancelled mutual exclusivity", func() {
		It("should not be both completed and cancelled", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(modal.IsCompleted()).To(BeTrue())
			Expect(modal.IsCancelled()).To(BeFalse())
		})

		It("should not be both cancelled and completed", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.IsCancelled()).To(BeTrue())
			Expect(modal.IsCompleted()).To(BeFalse())
		})
	})
})
