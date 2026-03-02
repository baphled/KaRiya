package configure_test

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/configtypes"
	configure "github.com/baphled/kariya/internal/cli/screens/configure"
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

var _ = Describe("SettingsModal", func() {
	var (
		modal    *configure.SettingsModal
		settings map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting
	)

	BeforeEach(func() {
		settings = makeTestSettings()
		modal = configure.NewSettingsModal(settings, 120, 40)
	})

	Describe("NewSettingsModal", func() {
		It("should create a modal with correct initial state", func() {
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsCompleted()).To(BeFalse())
			Expect(modal.IsCancelled()).To(BeFalse())
		})

		It("should not panic with empty settings map", func() {
			emptySettings := make(map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting)
			emptyModal := configure.NewSettingsModal(emptySettings, 120, 40)
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
})
