package configure_test

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/configtypes"
	"github.com/baphled/kariya/internal/cli/screens/configure"
	"github.com/baphled/kariya/internal/cli/themes"
)

type mockLogoRenderer struct{}

func (m *mockLogoRenderer) ViewStatic() string { return "LOGO" }
func (m *mockLogoRenderer) SetWidth(_ int)     {}

var _ = Describe("ConfirmModal", func() {
	var modal *configure.ConfirmModal

	BeforeEach(func() {
		modal = configure.NewConfirmModal("Save?", "Are you sure?", 80, 24)
	})

	Describe("NewConfirmModal", func() {
		It("should create a visible modal", func() {
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should not be confirmed or cancelled initially", func() {
			Expect(modal.IsConfirmed()).To(BeFalse())
			Expect(modal.IsCancelled()).To(BeFalse())
		})
	})

	Describe("Init", func() {
		It("should return nil", func() {
			Expect(modal.Init()).To(BeNil())
		})
	})

	Describe("Update", func() {
		It("should handle window size message", func() {
			cmd := modal.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
			Expect(cmd).To(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should confirm on enter", func() {
			cmd := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())
			Expect(modal.IsConfirmed()).To(BeTrue())
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should confirm on y", func() {
			cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
			Expect(cmd).To(BeNil())
			Expect(modal.IsConfirmed()).To(BeTrue())
		})

		It("should cancel on esc", func() {
			cmd := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
			Expect(modal.IsCancelled()).To(BeTrue())
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should cancel on n", func() {
			cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
			Expect(cmd).To(BeNil())
			Expect(modal.IsCancelled()).To(BeTrue())
		})

		It("should cancel on q", func() {
			cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
			Expect(cmd).To(BeNil())
			Expect(modal.IsCancelled()).To(BeTrue())
		})

		It("should ignore messages when not visible", func() {
			modal.Hide()
			cmd := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())
			Expect(modal.IsConfirmed()).To(BeFalse())
		})

		It("should return nil for unhandled key messages", func() {
			cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
			Expect(cmd).To(BeNil())
			Expect(modal.IsConfirmed()).To(BeFalse())
			Expect(modal.IsCancelled()).To(BeFalse())
		})
	})

	Describe("View", func() {
		It("should render content when visible", func() {
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should return empty when not visible", func() {
			modal.Hide()
			Expect(modal.View()).To(BeEmpty())
		})

		It("should handle narrow width", func() {
			narrow := configure.NewConfirmModal("Save?", "Sure?", 40, 24)
			view := narrow.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Render", func() {
		It("should render at specified dimensions", func() {
			result := modal.Render(100, 50)
			Expect(result).NotTo(BeEmpty())
		})
	})

	Describe("SetTheme", func() {
		It("should accept a theme", func() {
			theme := themes.NewDefaultTheme()
			modal.SetTheme(theme)
			Expect(modal.View()).NotTo(BeEmpty())
		})
	})

	Describe("Show", func() {
		It("should make modal visible and reset state", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.IsConfirmed()).To(BeTrue())
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.IsConfirmed()).To(BeFalse())
			Expect(modal.IsCancelled()).To(BeFalse())
		})
	})

	Describe("Hide", func() {
		It("should make modal not visible", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})
})

var _ = Describe("ReviewChangesModal", func() {
	var modal *configure.ReviewChangesModal

	BeforeEach(func() {
		changes := map[string]interface{}{
			"log_level": "debug",
			"db_path":   "/tmp/test.db",
		}
		modal = configure.NewReviewChangesModal(configtypes.DomainSystem, changes, 80, 24)
	})

	Describe("NewReviewChangesModal", func() {
		It("should create a visible modal", func() {
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should store changes", func() {
			Expect(modal.GetChanges()).To(HaveLen(2))
		})
	})

	Describe("Init", func() {
		It("should return nil", func() {
			Expect(modal.Init()).To(BeNil())
		})
	})

	Describe("Update", func() {
		It("should handle window size", func() {
			cmd := modal.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
			Expect(cmd).To(BeNil())
		})

		It("should confirm on enter", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.IsConfirmed()).To(BeTrue())
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should confirm on y", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
			Expect(modal.IsConfirmed()).To(BeTrue())
		})

		It("should cancel on esc", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.IsCancelled()).To(BeTrue())
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should cancel on q", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
			Expect(modal.IsCancelled()).To(BeTrue())
		})

		It("should cancel on n", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
			Expect(modal.IsCancelled()).To(BeTrue())
		})

		It("should ignore messages when not visible", func() {
			modal.Hide()
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.IsConfirmed()).To(BeFalse())
		})

		It("should return nil for unhandled key messages", func() {
			cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
			Expect(cmd).To(BeNil())
			Expect(modal.IsConfirmed()).To(BeFalse())
			Expect(modal.IsCancelled()).To(BeFalse())
		})
	})

	Describe("View", func() {
		It("should render changes when visible", func() {
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Review"))
		})

		It("should return empty when not visible", func() {
			modal.Hide()
			Expect(modal.View()).To(BeEmpty())
		})

		It("should handle no changes", func() {
			empty := configure.NewReviewChangesModal(configtypes.DomainSystem, map[string]interface{}{}, 80, 24)
			view := empty.View()
			Expect(view).To(ContainSubstring("No changes"))
		})

		It("should handle narrow width", func() {
			narrow := configure.NewReviewChangesModal(configtypes.DomainSystem, map[string]interface{}{"k": "v"}, 30, 24)
			view := narrow.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Render", func() {
		It("should render at specified dimensions", func() {
			result := modal.Render(100, 50)
			Expect(result).NotTo(BeEmpty())
		})
	})

	Describe("SetTheme", func() {
		It("should accept a theme", func() {
			theme := themes.NewDefaultTheme()
			modal.SetTheme(theme)
			Expect(modal.View()).NotTo(BeEmpty())
		})
	})

	Describe("Show", func() {
		It("should reset state and make visible", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.IsConfirmed()).To(BeFalse())
			Expect(modal.IsCancelled()).To(BeFalse())
		})
	})

	Describe("Hide", func() {
		It("should make modal not visible", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("GetChanges", func() {
		It("should return the changes map", func() {
			changes := modal.GetChanges()
			Expect(changes).To(HaveKey("log_level"))
			Expect(changes).To(HaveKey("db_path"))
		})
	})
})

var _ = Describe("EditSettingsModal", func() {
	var (
		modal    *configure.EditSettingsModal
		settings []*configtypes.ConfigurationSetting
	)

	BeforeEach(func() {
		settings = []*configtypes.ConfigurationSetting{
			{Key: "name", Label: "Name", Value: "Alice", Type: "string", Description: "User name"},
			{Key: "age", Label: "Age", Value: 30, Type: "int", Description: "User age"},
			{Key: "active", Label: "Active", Value: true, Type: "bool", Description: "Is active"},
			{Key: "theme", Label: "Theme", Value: "dark", Type: "select", Options: []string{"dark", "light"}, Description: "UI theme"},
		}
		modal = configure.NewEditSettingsModal(configtypes.DomainProfile, settings, 80, 24)
	})

	Describe("NewEditSettingsModal", func() {
		It("should create a visible modal", func() {
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should not be completed or cancelled", func() {
			Expect(modal.IsCompleted()).To(BeFalse())
			Expect(modal.IsCancelled()).To(BeFalse())
		})

		It("should handle empty settings", func() {
			empty := configure.NewEditSettingsModal(configtypes.DomainSystem, nil, 80, 24)
			Expect(empty).NotTo(BeNil())
			Expect(empty.IsVisible()).To(BeTrue())
		})

		It("should handle select settings with options", func() {
			noOpts := []*configtypes.ConfigurationSetting{
				{Key: "sel", Label: "Sel", Value: "a", Type: "select", Options: []string{"a", "b"}},
			}
			m := configure.NewEditSettingsModal(configtypes.DomainSystem, noOpts, 80, 24)
			Expect(m).NotTo(BeNil())
		})

		It("should handle default type settings", func() {
			def := []*configtypes.ConfigurationSetting{
				{Key: "custom", Label: "Custom", Value: "val", Type: "unknown_type"},
			}
			m := configure.NewEditSettingsModal(configtypes.DomainSystem, def, 80, 24)
			Expect(m).NotTo(BeNil())
		})

		It("should handle very wide modal", func() {
			wide := configure.NewEditSettingsModal(configtypes.DomainSystem, settings, 200, 24)
			Expect(wide).NotTo(BeNil())
			Expect(wide.View()).NotTo(BeEmpty())
		})

		It("should handle very narrow modal", func() {
			narrow := configure.NewEditSettingsModal(configtypes.DomainSystem, settings, 40, 24)
			Expect(narrow).NotTo(BeNil())
			Expect(narrow.View()).NotTo(BeEmpty())
		})

		It("should handle very short modal", func() {
			short := configure.NewEditSettingsModal(configtypes.DomainSystem, settings, 80, 10)
			Expect(short).NotTo(BeNil())
			Expect(short.View()).NotTo(BeEmpty())
		})

		It("should handle select with empty options", func() {
			emptyOpts := []*configtypes.ConfigurationSetting{
				{Key: "name", Label: "Name", Value: "test", Type: "string"},
				{Key: "sel", Label: "Sel", Value: "a", Type: "select", Options: []string{}},
			}
			m := configure.NewEditSettingsModal(configtypes.DomainSystem, emptyOpts, 80, 24)
			Expect(m).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		It("should return a command when form exists", func() {
			cmd := modal.Init()
			_ = cmd
		})

		It("should return nil when no form", func() {
			empty := configure.NewEditSettingsModal(configtypes.DomainSystem, nil, 80, 24)
			Expect(empty.Init()).To(BeNil())
		})
	})

	Describe("Update", func() {
		It("should handle window size", func() {
			cmd := modal.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
			Expect(cmd).To(BeNil())
		})

		It("should cancel on esc", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.IsCancelled()).To(BeTrue())
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should complete on ctrl+s", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(modal.IsCompleted()).To(BeTrue())
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should ignore messages when not visible", func() {
			modal.Hide()
			cmd := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
			Expect(modal.IsCancelled()).To(BeFalse())
		})

		It("should delegate to form for regular keys", func() {
			cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
			_ = cmd
		})

		It("should handle window resize to very small", func() {
			cmd := modal.Update(tea.WindowSizeMsg{Width: 30, Height: 10})
			Expect(cmd).To(BeNil())
		})

		It("should handle window resize to very large", func() {
			cmd := modal.Update(tea.WindowSizeMsg{Width: 200, Height: 100})
			Expect(cmd).To(BeNil())
		})

		It("should delegate non-key messages to form", func() {
			cmd := modal.Update(tea.MouseMsg{})
			_ = cmd
		})

		It("should return nil for non-key messages when no form", func() {
			empty := configure.NewEditSettingsModal(configtypes.DomainSystem, nil, 80, 24)
			cmd := empty.Update(tea.MouseMsg{})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("View", func() {
		It("should render form when visible", func() {
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Edit"))
		})

		It("should return empty when not visible", func() {
			modal.Hide()
			Expect(modal.View()).To(BeEmpty())
		})

		It("should show no settings message", func() {
			empty := configure.NewEditSettingsModal(configtypes.DomainSystem, nil, 80, 24)
			view := empty.View()
			Expect(view).To(ContainSubstring("No settings"))
		})
	})

	Describe("Render", func() {
		It("should render at specified dimensions", func() {
			result := modal.Render(100, 50)
			Expect(result).NotTo(BeEmpty())
		})
	})

	Describe("GetChanges", func() {
		It("should return empty when no changes made", func() {
			changes := modal.GetChanges()
			Expect(changes).To(BeEmpty())
		})

		It("should detect string changes", func() {
			formData := modal.GetFormData()
			if ptr, ok := formData.Values["name"]; ok {
				*ptr = "Bob"
			}
			changes := modal.GetChanges()
			Expect(changes).To(HaveKey("name"))
			Expect(changes["name"]).To(Equal("Bob"))
		})

		It("should detect int changes", func() {
			formData := modal.GetFormData()
			if ptr, ok := formData.Values["age"]; ok {
				*ptr = "42"
			}
			changes := modal.GetChanges()
			Expect(changes).To(HaveKey("age"))
			Expect(changes["age"]).To(Equal(42))
		})

		It("should detect bool changes", func() {
			formData := modal.GetFormData()
			if ptr, ok := formData.BoolValues["active"]; ok {
				*ptr = false
			}
			changes := modal.GetChanges()
			Expect(changes).To(HaveKey("active"))
			Expect(changes["active"]).To(BeFalse())
		})

		It("should detect select changes", func() {
			formData := modal.GetFormData()
			if ptr, ok := formData.Values["theme"]; ok {
				*ptr = "light"
			}
			changes := modal.GetChanges()
			Expect(changes).To(HaveKey("theme"))
			Expect(changes["theme"]).To(Equal("light"))
		})

		It("should handle invalid int value gracefully", func() {
			formData := modal.GetFormData()
			if ptr, ok := formData.Values["age"]; ok {
				*ptr = "not_a_number"
			}
			changes := modal.GetChanges()
			Expect(changes).NotTo(HaveKey("age"))
		})
	})

	Describe("SetTheme", func() {
		It("should accept a theme and rebuild form", func() {
			theme := themes.NewDefaultTheme()
			modal.SetTheme(theme)
			Expect(modal.View()).NotTo(BeEmpty())
		})
	})

	Describe("Show", func() {
		It("should reset state and make visible", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.IsCompleted()).To(BeFalse())
			Expect(modal.IsCancelled()).To(BeFalse())
		})
	})

	Describe("Hide", func() {
		It("should make modal not visible", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})
})

var _ = Describe("CompleteScreen", func() {
	var screen *configure.CompleteScreen

	BeforeEach(func() {
		screen = configure.NewCompleteScreen(configtypes.DomainSystem, 3)
		screen.SetTerminalInfo(120, 40)
	})

	Describe("Init", func() {
		It("should return nil", func() {
			Expect(screen.Init()).To(BeNil())
		})
	})

	Describe("SetTheme", func() {
		It("should accept a valid theme", func() {
			screen.SetTheme(themes.NewDefaultTheme())
			view := screen.View()
			Expect(view).To(ContainSubstring("Configuration Updated"))
		})

		It("should ignore non-theme values", func() {
			screen.SetTheme("not a theme")
			view := screen.View()
			Expect(view).To(ContainSubstring("Configuration Updated"))
		})
	})

	Describe("Update", func() {
		It("should handle enter key", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(result).NotTo(BeNil())
		})

		It("should handle esc key", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(result).NotTo(BeNil())
		})

		It("should handle q key", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
			Expect(result).NotTo(BeNil())
		})

		It("should handle window size", func() {
			cmd, result := screen.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
			Expect(result).To(BeNil())
			_ = cmd
		})

		It("should return nil for unhandled keys", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
			Expect(result).To(BeNil())
		})

		It("should return nil for non-key messages", func() {
			_, result := screen.Update(tea.MouseMsg{})
			Expect(result).To(BeNil())
		})
	})

	Describe("View", func() {
		It("should show no changes message when count is zero", func() {
			zeroScreen := configure.NewCompleteScreen(configtypes.DomainProfile, 0)
			zeroScreen.SetTerminalInfo(120, 40)
			view := zeroScreen.View()
			Expect(view).To(ContainSubstring("No changes"))
		})

		It("should show change count when changes exist", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("3"))
		})

		It("should render for each domain type", func() {
			for _, domain := range []configtypes.ConfigurationDomain{
				configtypes.DomainSystem,
				configtypes.DomainProfile,
				configtypes.DomainExport,
				configtypes.DomainUI,
			} {
				s := configure.NewCompleteScreen(domain, 1)
				s.SetTerminalInfo(120, 40)
				Expect(s.View()).NotTo(BeEmpty())
			}
		})
	})
})

var _ = Describe("ConfirmScreen", func() {
	var screen *configure.ConfirmScreen

	BeforeEach(func() {
		screen = configure.NewConfirmScreen(configtypes.DomainExport, 5)
		screen.SetTerminalInfo(120, 40)
	})

	Describe("Init", func() {
		It("should return nil", func() {
			Expect(screen.Init()).To(BeNil())
		})
	})

	Describe("SetTheme", func() {
		It("should accept a valid theme", func() {
			screen.SetTheme(themes.NewDefaultTheme())
			view := screen.View()
			Expect(view).To(ContainSubstring("Confirm"))
		})

		It("should ignore non-theme values", func() {
			screen.SetTheme(42)
			view := screen.View()
			Expect(view).To(ContainSubstring("Confirm"))
		})
	})

	Describe("Update", func() {
		It("should handle esc key", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(result).NotTo(BeNil())
		})

		It("should handle enter key", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(result).NotTo(BeNil())
		})

		It("should handle y key", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
			Expect(result).NotTo(BeNil())
		})

		It("should handle Y key", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Y")})
			Expect(result).NotTo(BeNil())
		})

		It("should handle n key", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
			Expect(result).NotTo(BeNil())
		})

		It("should handle N key", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("N")})
			Expect(result).NotTo(BeNil())
		})

		It("should handle tab for button toggle", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyTab})
			Expect(result).To(BeNil())
		})

		It("should handle left/right for button toggle", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyLeft})
			Expect(result).To(BeNil())
			_, result = screen.Update(tea.KeyMsg{Type: tea.KeyRight})
			Expect(result).To(BeNil())
		})

		It("should handle h/l for button toggle", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
			Expect(result).To(BeNil())
			_, result = screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
			Expect(result).To(BeNil())
		})

		It("should handle window size", func() {
			cmd, result := screen.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
			Expect(result).To(BeNil())
			_ = cmd
		})

		It("should return nil for unhandled keys", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
			Expect(result).To(BeNil())
		})

		It("should return nil for non-key messages", func() {
			_, result := screen.Update(tea.MouseMsg{})
			Expect(result).To(BeNil())
		})
	})

	Describe("GetSelection", func() {
		It("should return false initially (cancel focused)", func() {
			Expect(screen.GetSelection()).To(BeFalse())
		})
	})

	Describe("SetSelection", func() {
		It("should set to yes", func() {
			screen.SetSelection(true)
			Expect(screen.GetSelection()).To(BeTrue())
		})

		It("should set to no", func() {
			screen.SetSelection(false)
			Expect(screen.GetSelection()).To(BeFalse())
		})
	})
})

var _ = Describe("FailedScreen", func() {
	var screen *configure.FailedScreen

	BeforeEach(func() {
		screen = configure.NewFailedScreen(configtypes.DomainUI, "save error")
		screen.SetTerminalInfo(120, 40)
	})

	Describe("Init", func() {
		It("should return nil", func() {
			Expect(screen.Init()).To(BeNil())
		})
	})

	Describe("SetTheme", func() {
		It("should accept a valid theme", func() {
			screen.SetTheme(themes.NewDefaultTheme())
			view := screen.View()
			Expect(view).To(ContainSubstring("Configuration Failed"))
		})

		It("should ignore non-theme values", func() {
			screen.SetTheme("invalid")
			view := screen.View()
			Expect(view).To(ContainSubstring("Configuration Failed"))
		})
	})

	Describe("Update", func() {
		It("should return CancelResult on esc", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(result).NotTo(BeNil())
		})

		It("should return NavigateResult on enter", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(result).NotTo(BeNil())
		})

		It("should return NavigateResult on r", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
			Expect(result).NotTo(BeNil())
		})

		It("should return CancelResult with metadata on q", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
			Expect(result).NotTo(BeNil())
		})

		It("should handle window size", func() {
			cmd, result := screen.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
			Expect(result).To(BeNil())
			_ = cmd
		})

		It("should return nil for unhandled keys", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
			Expect(result).To(BeNil())
		})

		It("should return nil for non-key messages", func() {
			_, result := screen.Update(tea.MouseMsg{})
			Expect(result).To(BeNil())
		})
	})

	Describe("NewFailedScreen", func() {
		It("should use default message when empty", func() {
			s := configure.NewFailedScreen(configtypes.DomainSystem, "")
			Expect(s.GetErrorMessage()).To(ContainSubstring("unknown error"))
		})
	})

	Describe("GetErrorMessage", func() {
		It("should return the error message", func() {
			Expect(screen.GetErrorMessage()).To(Equal("save error"))
		})
	})
})

var _ = Describe("DomainSelectScreen Init", func() {
	It("should return nil from Init", func() {
		screen := configure.NewDomainSelectScreen([]configtypes.ConfigurationDomain{
			configtypes.DomainSystem,
		})
		Expect(screen.Init()).To(BeNil())
	})
})

var _ = Describe("EditSettingsScreen", func() {
	var screen *configure.EditSettingsScreen

	BeforeEach(func() {
		settings := []*configtypes.ConfigurationSetting{
			{Key: "name", Label: "Name", Value: "Bob", Type: "string"},
			{Key: "count", Label: "Count", Value: 10, Type: "int"},
			{Key: "enabled", Label: "Enabled", Value: false, Type: "bool"},
			{Key: "format", Label: "Format", Value: "json", Type: "select", Options: []string{"json", "yaml"}},
			{Key: "custom", Label: "Custom", Value: "x", Type: "other"},
		}
		screen = configure.NewEditSettingsScreen(configtypes.DomainExport, settings)
		screen.SetTerminalInfo(120, 40)
	})

	Describe("Init", func() {
		It("should return a command for form init", func() {
			cmd := screen.Init()
			_ = cmd
		})

		It("should return nil when no form", func() {
			empty := configure.NewEditSettingsScreen(configtypes.DomainSystem, nil)
			Expect(empty.Init()).To(BeNil())
		})
	})

	Describe("SetTheme", func() {
		It("should accept a valid theme", func() {
			screen.SetTheme(themes.NewDefaultTheme())
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should ignore non-theme values", func() {
			screen.SetTheme("not a theme")
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("SetLogo", func() {
		It("should ignore non-logo values", func() {
			screen.SetLogo("not a logo", 0)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should accept a valid LogoRenderer", func() {
			screen.SetLogo(&mockLogoRenderer{}, 2)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Update", func() {
		It("should handle window size message", func() {
			cmd, result := screen.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
			Expect(result).To(BeNil())
			Expect(cmd).To(BeNil())
		})

		It("should handle esc key", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(result).NotTo(BeNil())
		})

		It("should handle ctrl+s key", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(result).NotTo(BeNil())
		})

		It("should delegate to form for regular keys", func() {
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
			_ = cmd
			_ = result
		})

		It("should handle non-key messages with form", func() {
			cmd, result := screen.Update(tea.MouseMsg{})
			_ = cmd
			_ = result
		})

		It("should handle window resize to very small height", func() {
			cmd, result := screen.Update(tea.WindowSizeMsg{Width: 80, Height: 15})
			Expect(result).To(BeNil())
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update without form", func() {
		It("should return nil for non-key messages", func() {
			empty := configure.NewEditSettingsScreen(configtypes.DomainSystem, nil)
			cmd, result := empty.Update(tea.MouseMsg{})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should still handle esc without form", func() {
			empty := configure.NewEditSettingsScreen(configtypes.DomainSystem, nil)
			_, result := empty.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(result).NotTo(BeNil())
		})
	})

	Describe("View", func() {
		It("should show no settings message", func() {
			empty := configure.NewEditSettingsScreen(configtypes.DomainSystem, nil)
			empty.SetTerminalInfo(120, 40)
			view := empty.View()
			Expect(view).To(ContainSubstring("No settings"))
		})

		It("should render form content", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render with logo", func() {
			screen.SetLogo(&mockLogoRenderer{}, 2)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("GetChanges", func() {
		It("should return empty when no changes made", func() {
			changes := screen.GetChanges()
			Expect(changes).To(BeEmpty())
		})

		It("should detect string changes", func() {
			formData := screen.GetFormData()
			if ptr, ok := formData.Values["name"]; ok {
				*ptr = "Alice"
			}
			changes := screen.GetChanges()
			Expect(changes).To(HaveKey("name"))
			Expect(changes["name"]).To(Equal("Alice"))
		})

		It("should detect int changes", func() {
			formData := screen.GetFormData()
			if ptr, ok := formData.Values["count"]; ok {
				*ptr = "42"
			}
			changes := screen.GetChanges()
			Expect(changes).To(HaveKey("count"))
			Expect(changes["count"]).To(Equal(42))
		})

		It("should detect bool changes", func() {
			formData := screen.GetFormData()
			if ptr, ok := formData.BoolValues["enabled"]; ok {
				*ptr = true
			}
			changes := screen.GetChanges()
			Expect(changes).To(HaveKey("enabled"))
			Expect(changes["enabled"]).To(BeTrue())
		})

		It("should detect default type changes", func() {
			formData := screen.GetFormData()
			if ptr, ok := formData.Values["custom"]; ok {
				*ptr = "y"
			}
			changes := screen.GetChanges()
			Expect(changes).To(HaveKey("custom"))
			Expect(changes["custom"]).To(Equal("y"))
		})

		It("should handle invalid int value gracefully", func() {
			formData := screen.GetFormData()
			if ptr, ok := formData.Values["count"]; ok {
				*ptr = "not_a_number"
			}
			changes := screen.GetChanges()
			Expect(changes).NotTo(HaveKey("count"))
		})
	})

	Describe("GetFormData", func() {
		It("should return form data", func() {
			fd := screen.GetFormData()
			Expect(fd).NotTo(BeNil())
			Expect(fd.Values).To(HaveKey("name"))
			Expect(fd.BoolValues).To(HaveKey("enabled"))
		})
	})

	Describe("SetTerminalInfo", func() {
		It("should update terminal dimensions", func() {
			screen.SetTerminalInfo(80, 24)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("createFieldForSetting via select with empty options", func() {
		It("should handle select with valid options", func() {
			noOpts := []*configtypes.ConfigurationSetting{
				{Key: "sel", Label: "Sel", Value: "a", Type: "select", Options: []string{"a", "b"}},
			}
			s := configure.NewEditSettingsScreen(configtypes.DomainSystem, noOpts)
			s.SetTerminalInfo(120, 40)
			Expect(s).NotTo(BeNil())
		})

		It("should handle select with empty options", func() {
			emptyOpts := []*configtypes.ConfigurationSetting{
				{Key: "name", Label: "Name", Value: "test", Type: "string"},
				{Key: "sel", Label: "Sel", Value: "a", Type: "select", Options: []string{}},
			}
			s := configure.NewEditSettingsScreen(configtypes.DomainSystem, emptyOpts)
			s.SetTerminalInfo(120, 40)
			Expect(s).NotTo(BeNil())
		})
	})

	Describe("rebuildForm edge cases", func() {
		It("should handle very small terminal height", func() {
			screen.SetTerminalInfo(120, 5)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})

var _ = Describe("ReviewChangesScreen", func() {
	var screen *configure.ReviewChangesScreen

	BeforeEach(func() {
		changes := map[string]interface{}{
			"log_level":                         "debug",
			"a_very_long_setting_key_name_here": "a_very_long_value_that_should_be_truncated_in_the_table_display",
		}
		screen = configure.NewReviewChangesScreen(
			configtypes.DomainSystem,
			changes,
			map[string]interface{}{"log_level": "info", "a_very_long_setting_key_name_here": "old_value"},
			map[string]string{"log_level": "Log Level"},
		)
		screen.SetTerminalInfo(120, 40)
	})

	Describe("Init", func() {
		It("should return nil", func() {
			Expect(screen.Init()).To(BeNil())
		})
	})

	Describe("Update", func() {
		It("should handle window size", func() {
			cmd, result := screen.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
			Expect(result).To(BeNil())
			Expect(cmd).To(BeNil())
		})

		It("should handle esc key", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(result).NotTo(BeNil())
		})

		It("should handle enter key", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(result).NotTo(BeNil())
		})

		It("should handle c key", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
			Expect(result).NotTo(BeNil())
		})

		It("should handle y key", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
			Expect(result).NotTo(BeNil())
		})

		It("should return nil for unhandled keys", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
			Expect(result).To(BeNil())
		})

		It("should return nil for non-key messages", func() {
			_, result := screen.Update(tea.MouseMsg{})
			Expect(result).To(BeNil())
		})
	})

	Describe("SetTheme", func() {
		It("should accept a valid theme", func() {
			screen.SetTheme(themes.NewDefaultTheme())
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should ignore non-theme values", func() {
			screen.SetTheme(false)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("SetLogo", func() {
		It("should ignore non-logo values", func() {
			screen.SetLogo("not a logo", 0)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should accept a valid LogoRenderer", func() {
			screen.SetLogo(&mockLogoRenderer{}, 2)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("View", func() {
		It("should render changes table", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should show no changes message", func() {
			noChanges := configure.NewReviewChangesScreen(
				configtypes.DomainUI,
				map[string]interface{}{},
				map[string]interface{}{},
				map[string]string{},
			)
			noChanges.SetTerminalInfo(120, 40)
			view := noChanges.View()
			Expect(view).To(ContainSubstring("No changes"))
		})

		It("should use key as label when label is empty", func() {
			s := configure.NewReviewChangesScreen(
				configtypes.DomainSystem,
				map[string]interface{}{"unlabeled_key": "new_val"},
				map[string]interface{}{"unlabeled_key": "old_val"},
				map[string]string{},
			)
			s.SetTerminalInfo(120, 40)
			view := s.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render with logo", func() {
			screen.SetLogo(&mockLogoRenderer{}, 2)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("GetChanges", func() {
		It("should return the changes map", func() {
			changes := screen.GetChanges()
			Expect(changes).To(HaveKey("log_level"))
		})
	})

	Describe("HasChanges", func() {
		It("should return true when changes exist", func() {
			Expect(screen.HasChanges()).To(BeTrue())
		})

		It("should return false when no changes", func() {
			noChanges := configure.NewReviewChangesScreen(
				configtypes.DomainUI,
				map[string]interface{}{},
				map[string]interface{}{},
				map[string]string{},
			)
			Expect(noChanges.HasChanges()).To(BeFalse())
		})
	})

	Describe("SetTerminalInfo", func() {
		It("should update terminal dimensions", func() {
			screen.SetTerminalInfo(80, 24)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})

var _ = Describe("formatDomainLabel via DomainSelectScreen", func() {
	It("should handle unknown domain", func() {
		screen := configure.NewDomainSelectScreen([]configtypes.ConfigurationDomain{
			configtypes.ConfigurationDomain("custom_domain"),
		})
		screen.SetTerminalInfo(120, 40)
		view := screen.View()
		Expect(view).To(ContainSubstring("custom_domain"))
	})
})
