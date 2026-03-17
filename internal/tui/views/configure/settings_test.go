package configure

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/ui/configtypes"
	"github.com/baphled/kariya/internal/ui/themes"
) // LSP: removed unused imports after boundary test block deletion

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

func getModalFormData(modal *Settings) map[configtypes.ConfigurationDomain]*forms.ConfigureSettingsFormData {
	return modal.formData
}

func setModalDomains(modal *Settings, domains []configtypes.ConfigurationDomain) {
	modal.domains = domains
}

func setModalSettings(modal *Settings, settings map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting) {
	modal.settings = settings
}

func setModalSelectedIdx(modal *Settings, idx int) {
	modal.selectedIdx = idx
}

func setModalActiveForm(modal *Settings, form forms.Form) {
	modal.activeForm = form
}

func updateFormDataString(data map[configtypes.ConfigurationDomain]*forms.ConfigureSettingsFormData, domain configtypes.ConfigurationDomain, key, value string) {
	data[domain].Values[key] = &value
}

func updateFormDataBool(data map[configtypes.ConfigurationDomain]*forms.ConfigureSettingsFormData, domain configtypes.ConfigurationDomain, key string, value bool) {
	data[domain].BoolValues[key] = &value
}

var _ = Describe("Settings", func() {
	var (
		modal    *Settings
		settings map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting
	)

	BeforeEach(func() {
		settings = makeTestSettings()
		modal = NewSettings(settings, 120, 40)
	})

	Describe("NewSettings", func() {
		It("should create a modal with correct initial state", func() {
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsCompleted()).To(BeFalse())
			Expect(modal.IsCancelled()).To(BeFalse())
		})

		It("should not panic with empty settings map", func() {
			emptySettings := make(map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting)
			emptyModal := NewSettings(emptySettings, 120, 40)
			Expect(emptyModal).NotTo(BeNil())
			Expect(emptyModal.IsCompleted()).To(BeFalse())
			Expect(emptyModal.IsCancelled()).To(BeFalse())
		})

		It("should initialise with no changes", func() {
			changes := modal.GetChanges()
			Expect(changes.Changes).To(BeEmpty())
		})
	})

	Describe("Init", func() {
		It("should return without error", func() {
			cmd := modal.Init()
			_ = cmd
		})
	})

	Describe("Update", func() {
		Context("when pressing j key (vim-style navigation)", func() {
			var selectOnlyModal *Settings

			BeforeEach(func() {
				selectSettings := map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting{
					configtypes.DomainSystem: {
						{Key: "log_level", Label: "Log Level", Value: "info", DefaultValue: "info", Type: "select", Options: []string{"debug", "info", "warn", "error"}},
					},
					configtypes.DomainExport: {
						{Key: "default_destination", Label: "Default Destination", Value: "file", DefaultValue: "file", Type: "select", Options: []string{"file", "clipboard"}},
					},
					configtypes.DomainUI: {
						{Key: "theme", Label: "Theme", Value: "dark", DefaultValue: "dark", Type: "select", Options: []string{"light", "dark"}},
					},
				}
				selectOnlyModal = NewSettings(selectSettings, 120, 40)
			})

			It("should move to the next domain", func() {
				initialIdx := selectOnlyModal.selectedIdx
				selectOnlyModal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				Expect(selectOnlyModal.selectedIdx).To(Equal(initialIdx + 1))
			})

			It("should not move beyond the last domain", func() {
				for range 10 {
					selectOnlyModal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				}
				Expect(selectOnlyModal.selectedIdx).To(Equal(len(selectOnlyModal.domains) - 1))
			})
		})

		Context("when pressing k key (vim-style navigation)", func() {
			var selectOnlyModal *Settings

			BeforeEach(func() {
				selectSettings := map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting{
					configtypes.DomainSystem: {
						{Key: "log_level", Label: "Log Level", Value: "info", DefaultValue: "info", Type: "select", Options: []string{"debug", "info", "warn", "error"}},
					},
					configtypes.DomainExport: {
						{Key: "default_destination", Label: "Default Destination", Value: "file", DefaultValue: "file", Type: "select", Options: []string{"file", "clipboard"}},
					},
					configtypes.DomainUI: {
						{Key: "theme", Label: "Theme", Value: "dark", DefaultValue: "dark", Type: "select", Options: []string{"light", "dark"}},
					},
				}
				selectOnlyModal = NewSettings(selectSettings, 120, 40)
			})

			It("should move to the previous domain", func() {
				selectOnlyModal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				initialIdx := selectOnlyModal.selectedIdx
				selectOnlyModal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
				Expect(selectOnlyModal.selectedIdx).To(Equal(initialIdx - 1))
			})

			It("should not move above the first domain", func() {
				for range 10 {
					selectOnlyModal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
				}
				Expect(selectOnlyModal.selectedIdx).To(Equal(0))
			})
		})

		Context("when navigating with down arrow", func() {
			It("should move to the next domain", func() {
				initialIdx := modal.selectedIdx
				modal.Update(tea.KeyMsg{Type: tea.KeyDown})
				Expect(modal.selectedIdx).To(Equal(initialIdx + 1))
			})

			It("should not move beyond the last domain", func() {
				for range 10 {
					modal.Update(tea.KeyMsg{Type: tea.KeyDown})
				}
				Expect(modal.selectedIdx).To(Equal(len(modal.domains) - 1))
			})
		})

		Context("when navigating with up arrow", func() {
			It("should move to the previous domain", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyDown})
				initialIdx := modal.selectedIdx
				modal.Update(tea.KeyMsg{Type: tea.KeyUp})
				Expect(modal.selectedIdx).To(Equal(initialIdx - 1))
			})

			It("should not move above the first domain", func() {
				for range 10 {
					modal.Update(tea.KeyMsg{Type: tea.KeyUp})
				}
				Expect(modal.selectedIdx).To(Equal(0))
			})
		})

		Context("when pressing tab", func() {
			It("should forward to form for field navigation", func() {
				initialIdx := modal.selectedIdx
				modal.Update(tea.KeyMsg{Type: tea.KeyTab})
				Expect(modal.selectedIdx).To(Equal(initialIdx))
			})
		})

		Context("when pressing shift+tab", func() {
			It("should forward to form for reverse field navigation", func() {
				initialIdx := modal.selectedIdx
				modal.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
				Expect(modal.selectedIdx).To(Equal(initialIdx))
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
			Expect(changes.Changes).To(BeEmpty())
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

		It("should not rebuild form when dimensions unchanged", func() {
			initialForm := modal.activeForm
			modal.SetDimensions(120, 40)
			Expect(modal.activeForm).To(BeIdenticalTo(initialForm))
		})

		It("should rebuild form when width changes", func() {
			initialForm := modal.activeForm
			modal.SetDimensions(150, 40)
			Expect(modal.activeForm).NotTo(BeIdenticalTo(initialForm))
		})

		It("should rebuild form when height changes", func() {
			initialForm := modal.activeForm
			modal.SetDimensions(120, 50)
			Expect(modal.activeForm).NotTo(BeIdenticalTo(initialForm))
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
			emptyModal := NewSettings(emptySettings, 120, 40)
			view := emptyModal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("rebuildActiveForm with nil domain data", func() {
		It("should handle missing form data for selected domain", func() {
			settingsWithSystem := makeSettingsWithDomain(configtypes.DomainSystem)
			missingFormModal := NewSettings(settingsWithSystem, 120, 40)
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
			emptySettings := NewSettings(settingsWithSystem, 120, 40)
			setModalDomains(emptySettings, []configtypes.ConfigurationDomain{configtypes.DomainSystem})
			setModalSettings(emptySettings, settingsWithSystem)
			setModalSelectedIdx(emptySettings, 0)
			setModalActiveForm(emptySettings, nil)
			view := emptySettings.View()
			Expect(view).To(ContainSubstring("No settings available"))
		})
	})

	Describe("renderSectionList with unknown domain", func() {
		It("should render the raw domain label", func() {
			settingsWithUnknown := makeSettingsWithDomain(
				configtypes.ConfigurationDomain("custom"),
				&configtypes.ConfigurationSetting{Key: "custom", Label: "Custom", Value: "value", DefaultValue: "value", Type: "string"},
			)
			unknownModal := NewSettings(settingsWithUnknown, 120, 40)
			setModalDomains(unknownModal, []configtypes.ConfigurationDomain{configtypes.ConfigurationDomain("custom")})
			setModalSelectedIdx(unknownModal, 0)
			view := unknownModal.View()
			Expect(view).To(ContainSubstring("custom"))
		})
	})

	Describe("formDimensions with small width", func() {
		It("should enforce minimum form width of 30", func() {
			smallModal := NewSettings(makeTestSettings(), 40, 40)
			view := smallModal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Init with nil activeForm", func() {
		It("should return nil when activeForm is nil", func() {
			emptySettings := make(map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting)
			emptyModal := NewSettings(emptySettings, 120, 40)
			cmd := emptyModal.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("renderFormPanel with nil activeForm", func() {
		It("should render muted message when activeForm is nil", func() {
			emptySettings := make(map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting)
			emptyModal := NewSettings(emptySettings, 120, 40)
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
			Expect(changes.Changes).To(HaveKeyWithValue("name", "Updated User"))
		})
	})

	Describe("GetChanges with bool value updates", func() {
		It("should return bool changes when value changes", func() {
			settingsWithBool := makeSettingsWithDomain(
				configtypes.DomainSystem,
				&configtypes.ConfigurationSetting{Key: "enabled", Label: "Enabled", Value: false, DefaultValue: false, Type: "bool"},
			)
			boolModal := NewSettings(settingsWithBool, 120, 40)
			data := getModalFormData(boolModal)
			updateFormDataBool(data, configtypes.DomainSystem, "enabled", true)
			changes := boolModal.GetChanges()
			Expect(changes.Changes).To(HaveKeyWithValue("enabled", true))
		})
	})

	Describe("GetChanges with int value updates", func() {
		It("should return int changes when value changes", func() {
			settingsWithInt := makeSettingsWithDomain(
				configtypes.DomainSystem,
				&configtypes.ConfigurationSetting{Key: "refresh_interval", Label: "Refresh Interval", Value: 5, DefaultValue: 5, Type: "int"},
			)
			intModal := NewSettings(settingsWithInt, 120, 40)
			data := getModalFormData(intModal)
			updateFormDataString(data, configtypes.DomainSystem, "refresh_interval", "10")
			changes := intModal.GetChanges()
			Expect(changes.Changes).To(HaveKeyWithValue("refresh_interval", 10))
		})
	})

	Describe("GetChanges with list value updates", func() {
		It("should return list changes when value changes", func() {
			settingsWithList := makeSettingsWithDomain(
				configtypes.DomainSystem,
				&configtypes.ConfigurationSetting{Key: "favorites", Label: "Favorites", Value: []string{"alpha"}, DefaultValue: []string{"alpha"}, Type: "list"},
			)
			listModal := NewSettings(settingsWithList, 120, 40)
			data := getModalFormData(listModal)
			updateFormDataString(data, configtypes.DomainSystem, "favorites", "alpha, beta")
			changes := listModal.GetChanges()
			Expect(changes.Changes).To(HaveKeyWithValue("favorites", []string{"alpha", "beta"}))
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
				modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
			Expect(modal.selectedIdx).To(Equal(len(modal.domains) - 1))
		})

		It("should not move up from first domain", func() {
			for range 20 {
				modal.Update(tea.KeyMsg{Type: tea.KeyUp})
			}
			Expect(modal.selectedIdx).To(Equal(0))
		})
	})

	Describe("Multiple state transitions", func() {
		It("should handle multiple navigation and action sequences", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			modal.Update(tea.KeyMsg{Type: tea.KeyUp})
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

	Describe("Enter key on input field", func() {
		It("should not complete the modal when pressing Enter on an input field", func() {
			profileSettings := makeSettingsWithDomain(
				configtypes.DomainProfile,
				&configtypes.ConfigurationSetting{Key: "name", Label: "Full Name", Value: "Test User", DefaultValue: "", Type: "string"},
				&configtypes.ConfigurationSetting{Key: "email", Label: "Email", Value: "test@example.com", DefaultValue: "", Type: "string"},
			)
			inputModal := NewSettings(profileSettings, 120, 40)
			inputModal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(inputModal.IsCompleted()).To(BeFalse())
		})

		It("should forward Enter key to the form without completing", func() {
			profileSettings := makeSettingsWithDomain(
				configtypes.DomainProfile,
				&configtypes.ConfigurationSetting{Key: "name", Label: "Full Name", Value: "Test User", DefaultValue: "", Type: "string"},
			)
			inputModal := NewSettings(profileSettings, 120, 40)
			initialCompleted := inputModal.IsCompleted()
			inputModal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(inputModal.IsCompleted()).To(Equal(initialCompleted))
			Expect(inputModal.IsCompleted()).To(BeFalse())
		})

		It("should not complete form after multiple Enter presses", func() {
			profileSettings := makeSettingsWithDomain(
				configtypes.DomainProfile,
				&configtypes.ConfigurationSetting{Key: "name", Label: "Full Name", Value: "Test User", DefaultValue: "", Type: "string"},
				&configtypes.ConfigurationSetting{Key: "email", Label: "Email", Value: "test@example.com", DefaultValue: "", Type: "string"},
			)
			inputModal := NewSettings(profileSettings, 120, 40)
			inputModal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			inputModal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			inputModal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(inputModal.IsCompleted()).To(BeFalse())
		})

		It("should complete modal when pressing Ctrl+S", func() {
			profileSettings := makeSettingsWithDomain(
				configtypes.DomainProfile,
				&configtypes.ConfigurationSetting{Key: "name", Label: "Full Name", Value: "Test User", DefaultValue: "", Type: "string"},
			)
			inputModal := NewSettings(profileSettings, 120, 40)
			Expect(inputModal.IsCompleted()).To(BeFalse())
			inputModal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(inputModal.IsCompleted()).To(BeTrue())
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

	Describe("Domain navigation with arrow keys", func() {
		It("should return Init command when switching domain with down arrow", func() {
			cmd := modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(cmd).NotTo(BeNil())
		})

		It("should return Init command when switching domain with up arrow", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			cmd := modal.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(cmd).NotTo(BeNil())
		})

		It("should maintain form focus after domain switch", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Arrow keys always navigate regardless of text input focus", func() {
		Context("when a text input field is focused", func() {
			var textInputModal *Settings

			BeforeEach(func() {
				profileSettings := makeSettingsWithDomain(
					configtypes.DomainProfile,
					&configtypes.ConfigurationSetting{Key: "name", Label: "Full Name", Value: "Test User", DefaultValue: "", Type: "string"},
				)
				textInputModal = NewSettings(profileSettings, 120, 40)
			})

			It("should navigate down with down arrow even in text input", func() {
				Expect(textInputModal.selectedIdx).To(Equal(0))
				textInputModal.Update(tea.KeyMsg{Type: tea.KeyDown})
				view := textInputModal.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should navigate up with up arrow even in text input", func() {
				view := textInputModal.View()
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	Describe("j/k keys respect text input focus state", func() {
		Context("when a text input field is focused", func() {
			var textInputModal *Settings

			BeforeEach(func() {
				profileSettings := map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting{
					configtypes.DomainProfile: {
						{Key: "name", Label: "Full Name", Value: "Test User", DefaultValue: "", Type: "string"},
					},
					configtypes.DomainSystem: {
						{Key: "log_level", Label: "Log Level", Value: "info", DefaultValue: "info", Type: "select", Options: []string{"debug", "info", "warn", "error"}},
					},
				}
				textInputModal = NewSettings(profileSettings, 120, 40)
				textInputModal.Update(tea.KeyMsg{Type: tea.KeyDown})
			})

			It("should NOT navigate with j key when text input is focused", func() {
				initialIdx := textInputModal.selectedIdx
				textInputModal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				Expect(textInputModal.selectedIdx).To(Equal(initialIdx))
			})

			It("should NOT navigate with k key when text input is focused", func() {
				initialIdx := textInputModal.selectedIdx
				textInputModal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
				Expect(textInputModal.selectedIdx).To(Equal(initialIdx))
			})
		})

		Context("when a non-text field is focused (select)", func() {
			var selectModal *Settings

			BeforeEach(func() {
				selectSettings := map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting{
					configtypes.DomainSystem: {
						{Key: "log_level", Label: "Log Level", Value: "info", DefaultValue: "info", Type: "select", Options: []string{"debug", "info", "warn", "error"}},
					},
					configtypes.DomainUI: {
						{Key: "theme", Label: "Theme", Value: "dark", DefaultValue: "dark", Type: "select", Options: []string{"light", "dark"}},
					},
				}
				selectModal = NewSettings(selectSettings, 120, 40)
			})

			It("should navigate with j key when select field is focused", func() {
				initialIdx := selectModal.selectedIdx
				selectModal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				Expect(selectModal.selectedIdx).To(Equal(initialIdx + 1))
			})

			It("should navigate with k key when select field is focused", func() {
				selectModal.Update(tea.KeyMsg{Type: tea.KeyDown})
				initialIdx := selectModal.selectedIdx
				selectModal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
				Expect(selectModal.selectedIdx).To(Equal(initialIdx - 1))
			})
		})
	})

})
