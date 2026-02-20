package configure_test

import (
	"github.com/baphled/kariya/internal/cli/configtypes"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/configure"
	"github.com/baphled/kariya/internal/cli/themes"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockLogoRenderer struct{}

func (m *mockLogoRenderer) ViewStatic() string { return "LOGO" }
func (m *mockLogoRenderer) SetWidth(_ int)     {}

var _ = Describe("EditSettingsScreen", func() {
	var (
		screen   *configure.EditSettingsScreen
		settings []*intents.ConfigurationSetting
	)

	BeforeEach(func() {
		settings = []*intents.ConfigurationSetting{
			{
				Key:          "log_level",
				Label:        "Log Level",
				Value:        "info",
				DefaultValue: "info",
				Type:         "select",
				Options:      []string{"debug", "info", "warn", "error"},
				Description:  "Logging verbosity",
			},
			{
				Key:          "max_retries",
				Label:        "Max Retries",
				Value:        3,
				DefaultValue: 3,
				Type:         "int",
				Description:  "Maximum retry attempts",
			},
			{
				Key:          "auto_save",
				Label:        "Auto Save",
				Value:        true,
				DefaultValue: true,
				Type:         "bool",
				Description:  "Auto save on exit",
			},
			{
				Key:          "username",
				Label:        "Username",
				Value:        "testuser",
				DefaultValue: "",
				Type:         "string",
				Description:  "Your username",
			},
		}
		screen = configure.NewEditSettingsScreen(intents.DomainSystem, settings)
	})

	Describe("NewEditSettingsScreen", func() {
		It("creates a screen successfully", func() {
			Expect(screen).NotTo(BeNil())
		})

		It("handles empty settings list", func() {
			emptyScreen := configure.NewEditSettingsScreen(intents.DomainProfile, []*intents.ConfigurationSetting{})
			Expect(emptyScreen).NotTo(BeNil())
		})

		It("initializes form data with original values", func() {
			changes := screen.GetChanges()
			Expect(changes).To(BeEmpty())
		})
	})

	Describe("Init", func() {
		It("returns a command for form initialization", func() {
			cmd := screen.Init()
			_ = cmd
		})

		It("returns nil for empty settings", func() {
			emptyScreen := configure.NewEditSettingsScreen(intents.DomainProfile, []*intents.ConfigurationSetting{})
			cmd := emptyScreen.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update", func() {
		Context("when pressing Escape", func() {
			It("returns CancelResult", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultCancel))
			})
		})

		Context("when pressing Ctrl+S to save", func() {
			It("returns SubmitResult", func() {
				msg := tea.KeyMsg{Type: tea.KeyCtrlS}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultSubmit))
			})

			It("returns changes in SubmitResult data", func() {
				msg := tea.KeyMsg{Type: tea.KeyCtrlS}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				data := result.Data()
				Expect(data).To(BeAssignableToTypeOf(map[string]interface{}{}))
			})
		})

		Context("when handling window resize", func() {
			It("updates terminal info without returning result", func() {
				msg := tea.WindowSizeMsg{Width: 100, Height: 50}
				_, result := screen.Update(msg)

				Expect(result).To(BeNil())
			})
		})

		Context("when navigating with Tab", func() {
			It("delegates to form without returning result", func() {
				msg := tea.KeyMsg{Type: tea.KeyTab}
				_, result := screen.Update(msg)

				Expect(result).To(BeNil())
			})
		})

		Context("when pressing Enter to advance form", func() {
			It("delegates Enter to form as submit confirmation", func() {
				msg := tea.KeyMsg{Type: tea.KeyEnter}
				_, result := screen.Update(msg)
				_ = result
			})
		})

		Context("with empty settings (no form)", func() {
			It("returns nil for Tab key", func() {
				emptyScreen := configure.NewEditSettingsScreen(intents.DomainProfile, []*intents.ConfigurationSetting{})
				msg := tea.KeyMsg{Type: tea.KeyTab}
				_, result := emptyScreen.Update(msg)
				Expect(result).To(BeNil())
			})

			It("returns nil for non-key messages", func() {
				emptyScreen := configure.NewEditSettingsScreen(intents.DomainProfile, []*intents.ConfigurationSetting{})
				msg := tea.MouseMsg{}
				_, result := emptyScreen.Update(msg)
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		It("renders the screen", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("shows the domain name in breadcrumbs", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("System"))
		})

		It("renders empty state message for empty settings", func() {
			emptyScreen := configure.NewEditSettingsScreen(intents.DomainProfile, []*intents.ConfigurationSetting{})
			view := emptyScreen.View()
			Expect(view).To(ContainSubstring("No settings available"))
		})
	})

	Describe("SetTerminalInfo", func() {
		It("updates terminal dimensions", func() {
			screen.SetTerminalInfo(200, 100)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("SetTheme", func() {
		It("accepts valid theme", func() {
			screen.SetTheme(themes.NewDefaultTheme())
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("ignores invalid theme type", func() {
			screen.SetTheme("not-a-theme")
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("accepts nil theme without panic", func() {
			screen.SetTheme(nil)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("SetLogo", func() {
		It("accepts valid LogoRenderer", func() {
			screen.SetLogo(&mockLogoRenderer{}, 2)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("ignores nil logo", func() {
			screen.SetLogo(nil, 2)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("ignores invalid logo type", func() {
			screen.SetLogo("not-a-logo", 2)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("GetChanges", func() {
		It("returns empty map when no changes made", func() {
			changes := screen.GetChanges()
			Expect(changes).To(BeEmpty())
		})

		It("detects string value changes", func() {
			formData := screen.GetFormData()
			newVal := "newuser"
			formData.Values["username"] = &newVal
			changes := screen.GetChanges()
			Expect(changes).To(HaveKeyWithValue("username", "newuser"))
		})

		It("detects int value changes", func() {
			formData := screen.GetFormData()
			newVal := "10"
			formData.Values["max_retries"] = &newVal
			changes := screen.GetChanges()
			Expect(changes).To(HaveKeyWithValue("max_retries", 10))
		})

		It("detects bool value changes", func() {
			formData := screen.GetFormData()
			newVal := false
			formData.BoolValues["auto_save"] = &newVal
			changes := screen.GetChanges()
			Expect(changes).To(HaveKeyWithValue("auto_save", false))
		})

		It("handles int parse error gracefully", func() {
			formData := screen.GetFormData()
			newVal := "not-a-number"
			formData.Values["max_retries"] = &newVal
			changes := screen.GetChanges()
			Expect(changes).NotTo(HaveKey("max_retries"))
		})

		It("handles nil string pointer", func() {
			formData := screen.GetFormData()
			formData.Values["username"] = nil
			changes := screen.GetChanges()
			Expect(changes).To(HaveKey("username"))
		})

		It("handles nil bool pointer", func() {
			formData := screen.GetFormData()
			formData.BoolValues["auto_save"] = nil
			changes := screen.GetChanges()
			Expect(changes).To(HaveKey("auto_save"))
		})
	})

	Describe("GetFormData", func() {
		It("returns form data", func() {
			formData := screen.GetFormData()
			Expect(formData).NotTo(BeNil())
			Expect(formData.Values).NotTo(BeNil())
			Expect(formData.BoolValues).NotTo(BeNil())
		})
	})

	Describe("createFieldForSetting coverage", func() {
		It("handles unknown field type via default branch", func() {
			unknownSettings := []*intents.ConfigurationSetting{
				{
					Key:          "custom_field",
					Label:        "Custom",
					Value:        "val",
					DefaultValue: "val",
					Type:         "unknown_type",
					Description:  "Unknown type field",
				},
			}
			s := configure.NewEditSettingsScreen(intents.DomainSystem, unknownSettings)
			Expect(s).NotTo(BeNil())
			view := s.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("handles select with no options", func() {
			selectSettings := []*intents.ConfigurationSetting{
				{
					Key:          "empty_select",
					Label:        "Empty Select",
					Value:        "",
					DefaultValue: "",
					Type:         "select",
					Options:      []string{},
					Description:  "Select with no options",
				},
			}
			s := configure.NewEditSettingsScreen(intents.DomainSystem, selectSettings)
			Expect(s).NotTo(BeNil())
		})
	})

	Describe("Configure Screen Escape/Help - EditSettings (Phase 2-Tier 3)", func() {
		BeforeEach(func() {
			screen = configure.NewEditSettingsScreen(intents.DomainSystem, settings)
		})

		It("should return CancelResult on Escape key", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})

		It("should handle '?' key without panic", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
			_, _ = screen.Update(msg)

			Expect(screen).NotTo(BeNil())
		})
	})
})

var _ = Describe("ReviewChangesScreen", func() {
	var (
		screen   *configure.ReviewChangesScreen
		changes  map[string]interface{}
		original map[string]interface{}
		labels   map[string]string
	)

	BeforeEach(func() {
		changes = map[string]interface{}{
			"log_level":   "debug",
			"max_retries": 5,
		}
		original = map[string]interface{}{
			"log_level":   "info",
			"max_retries": 3,
		}
		labels = map[string]string{
			"log_level":   "Log Level",
			"max_retries": "Max Retries",
		}
		screen = configure.NewReviewChangesScreen(intents.DomainSystem, changes, original, labels)
	})

	Describe("NewReviewChangesScreen", func() {
		It("creates a screen successfully", func() {
			Expect(screen).NotTo(BeNil())
		})

		It("handles empty changes", func() {
			emptyScreen := configure.NewReviewChangesScreen(
				intents.DomainProfile,
				map[string]interface{}{},
				map[string]interface{}{},
				map[string]string{},
			)
			Expect(emptyScreen).NotTo(BeNil())
		})

		It("handles nil maps gracefully", func() {
			nilScreen := configure.NewReviewChangesScreen(
				intents.DomainSystem,
				nil,
				nil,
				nil,
			)
			Expect(nilScreen).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		It("returns nil command", func() {
			cmd := screen.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update", func() {
		Context("when pressing Escape", func() {
			It("returns CancelResult", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultCancel))
			})
		})

		Context("when pressing Enter", func() {
			It("returns NavigateResult to confirm", func() {
				msg := tea.KeyMsg{Type: tea.KeyEnter}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).To(Equal("confirm"))
			})
		})

		Context("when pressing 'c' shortcut", func() {
			It("returns NavigateResult to confirm", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).To(Equal("confirm"))
			})
		})

		Context("when pressing 'y' shortcut", func() {
			It("returns NavigateResult to confirm", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).To(Equal("confirm"))
			})
		})

		Context("when handling window resize", func() {
			It("updates terminal info without returning result", func() {
				msg := tea.WindowSizeMsg{Width: 100, Height: 50}
				_, result := screen.Update(msg)

				Expect(result).To(BeNil())
			})
		})

		Context("when pressing other keys", func() {
			It("returns nil result", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
				_, result := screen.Update(msg)

				Expect(result).To(BeNil())
			})
		})

		Context("with non-key messages", func() {
			It("returns nil result", func() {
				msg := tea.MouseMsg{}
				_, result := screen.Update(msg)
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		It("renders the screen", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("shows the domain name in breadcrumbs", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("System"))
		})

		It("shows 'Review Changes' in breadcrumbs", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Review Changes"))
		})

		It("shows number of changes", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("2 setting(s) changed"))
		})

		It("shows table headers", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Setting"))
			Expect(view).To(ContainSubstring("Original"))
			Expect(view).To(ContainSubstring("New Value"))
		})

		It("shows setting labels in table", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Log Level"))
			Expect(view).To(ContainSubstring("Max Retries"))
		})

		It("shows no changes message when empty", func() {
			emptyScreen := configure.NewReviewChangesScreen(
				intents.DomainProfile,
				map[string]interface{}{},
				map[string]interface{}{},
				map[string]string{},
			)
			view := emptyScreen.View()
			Expect(view).To(ContainSubstring("No changes to review"))
		})

		It("uses key as label when label not provided", func() {
			screenWithoutLabels := configure.NewReviewChangesScreen(
				intents.DomainSystem,
				map[string]interface{}{"some_key": "value"},
				map[string]interface{}{"some_key": "old"},
				map[string]string{},
			)
			view := screenWithoutLabels.View()
			Expect(view).To(ContainSubstring("some_key"))
		})

		It("truncates long values in table", func() {
			longScreen := configure.NewReviewChangesScreen(
				intents.DomainSystem,
				map[string]interface{}{"k": "this is a very long value that should be truncated in the table display"},
				map[string]interface{}{"k": "another very long original value that should also be truncated"},
				map[string]string{"k": "A Very Long Setting Label That Exceeds Column Width"},
			)
			view := longScreen.View()
			Expect(view).To(ContainSubstring("..."))
		})
	})

	Describe("SetTerminalInfo", func() {
		It("updates terminal dimensions", func() {
			screen.SetTerminalInfo(200, 100)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("SetTheme", func() {
		It("accepts valid theme", func() {
			screen.SetTheme(themes.NewDefaultTheme())
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("ignores invalid theme type", func() {
			screen.SetTheme("not-a-theme")
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("accepts nil theme without panic", func() {
			screen.SetTheme(nil)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("SetLogo", func() {
		It("accepts valid LogoRenderer", func() {
			screen.SetLogo(&mockLogoRenderer{}, 2)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("ignores nil logo", func() {
			screen.SetLogo(nil, 2)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("ignores invalid logo type", func() {
			screen.SetLogo("not-a-logo", 2)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("GetChanges", func() {
		It("returns the changes map", func() {
			result := screen.GetChanges()
			Expect(result).To(HaveLen(2))
			Expect(result["log_level"]).To(Equal("debug"))
			Expect(result["max_retries"]).To(Equal(5))
		})

		It("returns nil when changes is nil", func() {
			nilScreen := configure.NewReviewChangesScreen(
				intents.DomainSystem,
				nil,
				nil,
				nil,
			)
			result := nilScreen.GetChanges()
			Expect(result).To(BeNil())
		})
	})

	Describe("HasChanges", func() {
		It("returns true when there are changes", func() {
			Expect(screen.HasChanges()).To(BeTrue())
		})

		It("returns false when no changes", func() {
			emptyScreen := configure.NewReviewChangesScreen(
				intents.DomainProfile,
				map[string]interface{}{},
				map[string]interface{}{},
				map[string]string{},
			)
			Expect(emptyScreen.HasChanges()).To(BeFalse())
		})

		It("returns false when changes is nil", func() {
			nilScreen := configure.NewReviewChangesScreen(
				intents.DomainSystem,
				nil,
				nil,
				nil,
			)
			Expect(nilScreen.HasChanges()).To(BeFalse())
		})
	})
})

var _ = Describe("ConfirmScreen", func() {
	var screen *configure.ConfirmScreen

	BeforeEach(func() {
		screen = configure.NewConfirmScreen(intents.DomainSystem, 3)
	})

	Describe("NewConfirmScreen", func() {
		It("creates a screen successfully", func() {
			Expect(screen).NotTo(BeNil())
		})

		It("defaults to No selection (safer)", func() {
			Expect(screen.GetSelection()).To(BeFalse())
		})
	})

	Describe("Init", func() {
		It("returns nil command", func() {
			cmd := screen.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update", func() {
		Context("when pressing Escape", func() {
			It("returns CancelResult", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultCancel))
			})
		})

		Context("when pressing Enter", func() {
			It("returns NavigateResult with current selection", func() {
				msg := tea.KeyMsg{Type: tea.KeyEnter}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
			})
		})

		Context("when pressing 'y' for yes", func() {
			It("returns NavigateResult with true", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).To(BeTrue())
			})
		})

		Context("when pressing 'Y' for yes", func() {
			It("returns NavigateResult with true", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'Y'}}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).To(BeTrue())
			})
		})

		Context("when pressing 'n' for no", func() {
			It("returns NavigateResult with false", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).To(BeFalse())
			})
		})

		Context("when pressing 'N' for no", func() {
			It("returns NavigateResult with false", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).To(BeFalse())
			})
		})

		Context("when handling window resize", func() {
			It("updates dimensions without returning result", func() {
				msg := tea.WindowSizeMsg{Width: 100, Height: 50}
				_, result := screen.Update(msg)
				Expect(result).To(BeNil())
			})
		})

		Context("when toggling selection", func() {
			It("toggles with left/right keys", func() {
				Expect(screen.GetSelection()).To(BeFalse())

				screen.Update(tea.KeyMsg{Type: tea.KeyLeft})
				Expect(screen.GetSelection()).To(BeTrue())

				screen.Update(tea.KeyMsg{Type: tea.KeyRight})
				Expect(screen.GetSelection()).To(BeFalse())
			})

			It("toggles with h/l keys", func() {
				screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
				Expect(screen.GetSelection()).To(BeTrue())

				screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
				Expect(screen.GetSelection()).To(BeFalse())
			})

			It("toggles with tab key", func() {
				screen.Update(tea.KeyMsg{Type: tea.KeyTab})
				Expect(screen.GetSelection()).To(BeTrue())
			})
		})

		Context("when pressing unhandled keys", func() {
			It("returns nil result", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
				_, result := screen.Update(msg)
				Expect(result).To(BeNil())
			})
		})

		Context("with non-key messages", func() {
			It("returns nil result", func() {
				msg := tea.MouseMsg{}
				_, result := screen.Update(msg)
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		It("renders the screen", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("shows confirmation title", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Confirm"))
		})

		It("shows change count", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("3"))
		})
	})

	Describe("SetTheme", func() {
		It("accepts valid theme", func() {
			screen.SetTheme(themes.NewDefaultTheme())
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("ignores nil theme", func() {
			screen.SetTheme(nil)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("ignores invalid theme type", func() {
			screen.SetTheme("not-a-theme")
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("SetSelection", func() {
		It("sets selection to yes", func() {
			screen.SetSelection(true)
			Expect(screen.GetSelection()).To(BeTrue())
		})

		It("sets selection to no", func() {
			screen.SetSelection(false)
			Expect(screen.GetSelection()).To(BeFalse())
		})
	})
})

var _ = Describe("CompleteScreen", func() {
	var screen *configure.CompleteScreen

	BeforeEach(func() {
		screen = configure.NewCompleteScreen(intents.DomainProfile, 5)
	})

	Describe("NewCompleteScreen", func() {
		It("creates a screen successfully", func() {
			Expect(screen).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		It("returns nil command", func() {
			cmd := screen.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update", func() {
		Context("when pressing Enter", func() {
			It("returns SubmitResult", func() {
				msg := tea.KeyMsg{Type: tea.KeyEnter}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultSubmit))
			})
		})

		Context("when pressing Escape", func() {
			It("returns SubmitResult", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultSubmit))
			})
		})

		Context("when pressing q", func() {
			It("returns SubmitResult", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultSubmit))
			})
		})

		Context("when handling window resize", func() {
			It("updates dimensions without returning result", func() {
				msg := tea.WindowSizeMsg{Width: 100, Height: 50}
				_, result := screen.Update(msg)
				Expect(result).To(BeNil())
			})
		})

		Context("when pressing unhandled keys", func() {
			It("returns nil result", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
				_, result := screen.Update(msg)
				Expect(result).To(BeNil())
			})
		})

		Context("with non-key messages", func() {
			It("returns nil result", func() {
				msg := tea.MouseMsg{}
				_, result := screen.Update(msg)
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		It("renders the screen", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("shows success message", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Updated"))
		})

		It("shows change count", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("5"))
		})

		It("shows no changes message when count is zero", func() {
			emptyScreen := configure.NewCompleteScreen(intents.DomainSystem, 0)
			view := emptyScreen.View()
			Expect(view).To(ContainSubstring("No changes"))
		})

		It("renders for all domains", func() {
			for _, domain := range []configtypes.ConfigurationDomain{
				intents.DomainSystem,
				intents.DomainProfile,
				intents.DomainExport,
				intents.DomainUI,
				configtypes.ConfigurationDomain("custom"),
			} {
				s := configure.NewCompleteScreen(domain, 1)
				view := s.View()
				Expect(view).NotTo(BeEmpty())
			}
		})
	})

	Describe("SetTheme", func() {
		It("accepts valid theme", func() {
			screen.SetTheme(themes.NewDefaultTheme())
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("ignores nil theme", func() {
			screen.SetTheme(nil)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("ignores invalid theme type", func() {
			screen.SetTheme("not-a-theme")
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})

var _ = Describe("FailedScreen", func() {
	var screen *configure.FailedScreen

	BeforeEach(func() {
		screen = configure.NewFailedScreen(intents.DomainExport, "Connection timed out")
	})

	Describe("NewFailedScreen", func() {
		It("creates a screen successfully", func() {
			Expect(screen).NotTo(BeNil())
		})

		It("uses default error message when empty", func() {
			emptyScreen := configure.NewFailedScreen(intents.DomainSystem, "")
			Expect(emptyScreen.GetErrorMessage()).To(ContainSubstring("unknown error"))
		})
	})

	Describe("Init", func() {
		It("returns nil command", func() {
			cmd := screen.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update", func() {
		Context("when pressing Escape", func() {
			It("returns CancelResult", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultCancel))
			})
		})

		Context("when pressing Enter for retry", func() {
			It("returns NavigateResult with retry", func() {
				msg := tea.KeyMsg{Type: tea.KeyEnter}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).To(Equal("retry"))
			})
		})

		Context("when pressing 'r' for retry", func() {
			It("returns NavigateResult with retry", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).To(Equal("retry"))
			})
		})

		Context("when pressing q to quit", func() {
			It("returns CancelResult with main_menu metadata", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultCancel))
			})
		})

		Context("when handling window resize", func() {
			It("updates dimensions without returning result", func() {
				msg := tea.WindowSizeMsg{Width: 100, Height: 50}
				_, result := screen.Update(msg)
				Expect(result).To(BeNil())
			})
		})

		Context("when pressing unhandled keys", func() {
			It("returns nil result", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
				_, result := screen.Update(msg)
				Expect(result).To(BeNil())
			})
		})

		Context("with non-key messages", func() {
			It("returns nil result", func() {
				msg := tea.MouseMsg{}
				_, result := screen.Update(msg)
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		It("renders the screen", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("shows error title", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Failed"))
		})

		It("shows error message", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Connection timed out"))
		})
	})

	Describe("SetTheme", func() {
		It("accepts valid theme", func() {
			screen.SetTheme(themes.NewDefaultTheme())
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("ignores nil theme", func() {
			screen.SetTheme(nil)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("ignores invalid theme type", func() {
			screen.SetTheme("not-a-theme")
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("GetErrorMessage", func() {
		It("returns the error message", func() {
			Expect(screen.GetErrorMessage()).To(Equal("Connection timed out"))
		})
	})
})

var _ = Describe("DomainSelectScreen", func() {
	Describe("NewDomainSelectScreen", func() {
		It("creates a screen with all domains", func() {
			domains := []configtypes.ConfigurationDomain{
				intents.DomainSystem,
				intents.DomainProfile,
				intents.DomainExport,
				intents.DomainUI,
			}
			screen := configure.NewDomainSelectScreen(domains)
			Expect(screen).NotTo(BeNil())
		})

		It("creates a screen with empty domains", func() {
			screen := configure.NewDomainSelectScreen([]configtypes.ConfigurationDomain{})
			Expect(screen).NotTo(BeNil())
		})

		It("creates a screen with custom domain", func() {
			domains := []configtypes.ConfigurationDomain{
				configtypes.ConfigurationDomain("custom_domain"),
			}
			screen := configure.NewDomainSelectScreen(domains)
			Expect(screen).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		It("returns nil command", func() {
			domains := []configtypes.ConfigurationDomain{intents.DomainSystem}
			screen := configure.NewDomainSelectScreen(domains)
			cmd := screen.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("View", func() {
		It("renders the screen", func() {
			domains := []configtypes.ConfigurationDomain{
				intents.DomainSystem,
				intents.DomainProfile,
			}
			screen := configure.NewDomainSelectScreen(domains)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
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
			{
				Key:          "username",
				Label:        "Username",
				Value:        "testuser",
				DefaultValue: "",
				Type:         "string",
				Description:  "Your username",
			},
			{
				Key:          "max_retries",
				Label:        "Max Retries",
				Value:        3,
				DefaultValue: 3,
				Type:         "int",
				Description:  "Maximum retry attempts",
			},
			{
				Key:          "auto_save",
				Label:        "Auto Save",
				Value:        true,
				DefaultValue: true,
				Type:         "bool",
				Description:  "Auto save on exit",
			},
			{
				Key:          "log_level",
				Label:        "Log Level",
				Value:        "info",
				DefaultValue: "info",
				Type:         "select",
				Options:      []string{"debug", "info", "warn", "error"},
				Description:  "Logging verbosity",
			},
		}
		modal = configure.NewEditSettingsModal(intents.DomainSystem, settings, 120, 40)
	})

	Describe("NewEditSettingsModal", func() {
		It("creates a modal successfully", func() {
			Expect(modal).NotTo(BeNil())
		})

		It("starts visible", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("starts not completed", func() {
			Expect(modal.IsCompleted()).To(BeFalse())
		})

		It("starts not cancelled", func() {
			Expect(modal.IsCancelled()).To(BeFalse())
		})

		It("handles empty settings", func() {
			emptyModal := configure.NewEditSettingsModal(intents.DomainProfile, []*configtypes.ConfigurationSetting{}, 120, 40)
			Expect(emptyModal).NotTo(BeNil())
			Expect(emptyModal.IsVisible()).To(BeTrue())
		})

		It("handles select with no options", func() {
			noOptSettings := []*configtypes.ConfigurationSetting{
				{
					Key:     "empty_select",
					Label:   "Empty",
					Value:   "",
					Type:    "select",
					Options: []string{},
				},
			}
			m := configure.NewEditSettingsModal(intents.DomainSystem, noOptSettings, 120, 40)
			Expect(m).NotTo(BeNil())
		})

		It("handles unknown field type", func() {
			unknownSettings := []*configtypes.ConfigurationSetting{
				{
					Key:   "custom",
					Label: "Custom",
					Value: "val",
					Type:  "unknown_type",
				},
			}
			m := configure.NewEditSettingsModal(intents.DomainSystem, unknownSettings, 120, 40)
			Expect(m).NotTo(BeNil())
		})

		It("handles small width", func() {
			m := configure.NewEditSettingsModal(intents.DomainSystem, settings, 40, 40)
			Expect(m).NotTo(BeNil())
		})

		It("handles large width", func() {
			m := configure.NewEditSettingsModal(intents.DomainSystem, settings, 200, 40)
			Expect(m).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		It("returns a command for form initialization", func() {
			cmd := modal.Init()
			_ = cmd
		})

		It("returns nil for empty settings", func() {
			emptyModal := configure.NewEditSettingsModal(intents.DomainProfile, []*configtypes.ConfigurationSetting{}, 120, 40)
			cmd := emptyModal.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update", func() {
		Context("when hidden", func() {
			It("returns nil", func() {
				modal.Hide()
				cmd := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).To(BeNil())
			})
		})

		Context("when pressing Escape", func() {
			It("cancels the modal", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(modal.IsCancelled()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when pressing Ctrl+S", func() {
			It("completes the modal", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
				Expect(modal.IsCompleted()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when handling window resize", func() {
			It("updates dimensions", func() {
				cmd := modal.Update(tea.WindowSizeMsg{Width: 200, Height: 80})
				Expect(cmd).To(BeNil())
				Expect(modal.IsVisible()).To(BeTrue())
			})
		})

		Context("when pressing Tab", func() {
			It("delegates to form", func() {
				cmd := modal.Update(tea.KeyMsg{Type: tea.KeyTab})
				_ = cmd
				Expect(modal.IsVisible()).To(BeTrue())
			})
		})

		Context("with non-key messages", func() {
			It("delegates to form", func() {
				cmd := modal.Update(tea.MouseMsg{})
				_ = cmd
			})
		})

		Context("with empty settings (no form)", func() {
			It("returns nil for non-special keys", func() {
				emptyModal := configure.NewEditSettingsModal(intents.DomainProfile, []*configtypes.ConfigurationSetting{}, 120, 40)
				cmd := emptyModal.Update(tea.KeyMsg{Type: tea.KeyTab})
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		It("renders when visible", func() {
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("returns empty when hidden", func() {
			modal.Hide()
			view := modal.View()
			Expect(view).To(BeEmpty())
		})

		It("shows empty state for no settings", func() {
			emptyModal := configure.NewEditSettingsModal(intents.DomainProfile, []*configtypes.ConfigurationSetting{}, 120, 40)
			view := emptyModal.View()
			Expect(view).To(ContainSubstring("No settings available"))
		})

		It("shows domain name in title", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("System"))
		})
	})

	Describe("Render", func() {
		It("renders at specified dimensions", func() {
			view := modal.Render(100, 50)
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("GetChanges", func() {
		It("returns empty map when no changes", func() {
			changes := modal.GetChanges()
			Expect(changes).To(BeEmpty())
		})

		It("detects int value changes", func() {
			formData := modal.GetFormData()
			newVal := "10"
			formData.Values["max_retries"] = &newVal
			changes := modal.GetChanges()
			Expect(changes).To(HaveKeyWithValue("max_retries", 10))
		})

		It("detects bool value changes", func() {
			formData := modal.GetFormData()
			newVal := false
			formData.BoolValues["auto_save"] = &newVal
			changes := modal.GetChanges()
			Expect(changes).To(HaveKeyWithValue("auto_save", false))
		})

		It("detects string value changes", func() {
			formData := modal.GetFormData()
			newVal := "newuser"
			formData.Values["username"] = &newVal
			changes := modal.GetChanges()
			Expect(changes).To(HaveKeyWithValue("username", "newuser"))
		})

		It("handles int parse error gracefully", func() {
			formData := modal.GetFormData()
			newVal := "not-a-number"
			formData.Values["max_retries"] = &newVal
			changes := modal.GetChanges()
			Expect(changes).NotTo(HaveKey("max_retries"))
		})

		It("handles nil string pointer", func() {
			formData := modal.GetFormData()
			formData.Values["username"] = nil
			changes := modal.GetChanges()
			Expect(changes).To(HaveKey("username"))
		})

		It("handles nil bool pointer", func() {
			formData := modal.GetFormData()
			formData.BoolValues["auto_save"] = nil
			changes := modal.GetChanges()
			Expect(changes).To(HaveKey("auto_save"))
		})
	})

	Describe("GetFormData", func() {
		It("returns form data", func() {
			formData := modal.GetFormData()
			Expect(formData).NotTo(BeNil())
			Expect(formData.Values).NotTo(BeNil())
			Expect(formData.BoolValues).NotTo(BeNil())
		})
	})

	Describe("SetTheme", func() {
		It("updates the theme", func() {
			modal.SetTheme(themes.NewDefaultTheme())
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Show", func() {
		It("makes modal visible and resets state", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())

			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.IsCompleted()).To(BeFalse())
			Expect(modal.IsCancelled()).To(BeFalse())
		})
	})

	Describe("Hide", func() {
		It("makes modal invisible", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})
})

var _ = Describe("ConfirmModal", func() {
	var modal *configure.ConfirmModal

	BeforeEach(func() {
		modal = configure.NewConfirmModal("Confirm Action", "Are you sure?", 120, 40)
	})

	Describe("NewConfirmModal", func() {
		It("creates a modal successfully", func() {
			Expect(modal).NotTo(BeNil())
		})

		It("starts visible", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("starts not confirmed", func() {
			Expect(modal.IsConfirmed()).To(BeFalse())
		})

		It("starts not cancelled", func() {
			Expect(modal.IsCancelled()).To(BeFalse())
		})
	})

	Describe("Init", func() {
		It("returns nil command", func() {
			cmd := modal.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update", func() {
		Context("when hidden", func() {
			It("returns nil", func() {
				modal.Hide()
				cmd := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).To(BeNil())
			})
		})

		Context("when pressing Escape", func() {
			It("cancels the modal", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(modal.IsCancelled()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when pressing 'n'", func() {
			It("cancels the modal", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				Expect(modal.IsCancelled()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when pressing 'q'", func() {
			It("cancels the modal", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
				Expect(modal.IsCancelled()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when pressing Enter", func() {
			It("confirms the modal", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(modal.IsConfirmed()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when pressing 'y'", func() {
			It("confirms the modal", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
				Expect(modal.IsConfirmed()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when handling window resize", func() {
			It("updates dimensions", func() {
				cmd := modal.Update(tea.WindowSizeMsg{Width: 200, Height: 80})
				Expect(cmd).To(BeNil())
				Expect(modal.IsVisible()).To(BeTrue())
			})
		})

		Context("when pressing unhandled keys", func() {
			It("returns nil", func() {
				cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
				Expect(cmd).To(BeNil())
				Expect(modal.IsVisible()).To(BeTrue())
			})
		})

		Context("with non-key messages", func() {
			It("returns nil", func() {
				cmd := modal.Update(tea.MouseMsg{})
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		It("renders when visible", func() {
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("returns empty when hidden", func() {
			modal.Hide()
			view := modal.View()
			Expect(view).To(BeEmpty())
		})

		It("renders narrow modal for small width", func() {
			narrowModal := configure.NewConfirmModal("Title", "Message", 50, 30)
			view := narrowModal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Render", func() {
		It("renders at specified dimensions", func() {
			view := modal.Render(100, 50)
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("SetTheme", func() {
		It("updates the theme", func() {
			modal.SetTheme(themes.NewDefaultTheme())
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Show", func() {
		It("makes modal visible and resets state", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.IsCancelled()).To(BeTrue())

			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.IsConfirmed()).To(BeFalse())
			Expect(modal.IsCancelled()).To(BeFalse())
		})
	})

	Describe("Hide", func() {
		It("makes modal invisible", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})
})

var _ = Describe("ReviewChangesModal", func() {
	var modal *configure.ReviewChangesModal

	BeforeEach(func() {
		changes := map[string]interface{}{
			"log_level":   "debug",
			"max_retries": 5,
		}
		modal = configure.NewReviewChangesModal(intents.DomainSystem, changes, 120, 40)
	})

	Describe("NewReviewChangesModal", func() {
		It("creates a modal successfully", func() {
			Expect(modal).NotTo(BeNil())
		})

		It("starts visible", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("starts not confirmed", func() {
			Expect(modal.IsConfirmed()).To(BeFalse())
		})

		It("starts not cancelled", func() {
			Expect(modal.IsCancelled()).To(BeFalse())
		})
	})

	Describe("Init", func() {
		It("returns nil command", func() {
			cmd := modal.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update", func() {
		Context("when hidden", func() {
			It("returns nil", func() {
				modal.Hide()
				cmd := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).To(BeNil())
			})
		})

		Context("when pressing Escape", func() {
			It("cancels the modal", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(modal.IsCancelled()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when pressing 'q'", func() {
			It("cancels the modal", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
				Expect(modal.IsCancelled()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when pressing 'n'", func() {
			It("cancels the modal", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				Expect(modal.IsCancelled()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when pressing Enter", func() {
			It("confirms the modal", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(modal.IsConfirmed()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when pressing 'y'", func() {
			It("confirms the modal", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
				Expect(modal.IsConfirmed()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when handling window resize", func() {
			It("updates dimensions", func() {
				cmd := modal.Update(tea.WindowSizeMsg{Width: 200, Height: 80})
				Expect(cmd).To(BeNil())
				Expect(modal.IsVisible()).To(BeTrue())
			})
		})

		Context("when pressing unhandled keys", func() {
			It("returns nil", func() {
				cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
				Expect(cmd).To(BeNil())
				Expect(modal.IsVisible()).To(BeTrue())
			})
		})

		Context("with non-key messages", func() {
			It("returns nil", func() {
				cmd := modal.Update(tea.MouseMsg{})
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		It("renders when visible with changes", func() {
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("returns empty when hidden", func() {
			modal.Hide()
			view := modal.View()
			Expect(view).To(BeEmpty())
		})

		It("shows empty state for no changes", func() {
			emptyModal := configure.NewReviewChangesModal(intents.DomainProfile, map[string]interface{}{}, 120, 40)
			view := emptyModal.View()
			Expect(view).To(ContainSubstring("No changes"))
		})

		It("renders narrow modal for small width", func() {
			narrowModal := configure.NewReviewChangesModal(intents.DomainSystem, map[string]interface{}{"k": "v"}, 50, 30)
			view := narrowModal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Render", func() {
		It("renders at specified dimensions", func() {
			view := modal.Render(100, 50)
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("GetChanges", func() {
		It("returns the changes map", func() {
			changes := modal.GetChanges()
			Expect(changes).To(HaveLen(2))
			Expect(changes["log_level"]).To(Equal("debug"))
			Expect(changes["max_retries"]).To(Equal(5))
		})

		It("returns empty map for no changes", func() {
			emptyModal := configure.NewReviewChangesModal(intents.DomainSystem, map[string]interface{}{}, 120, 40)
			changes := emptyModal.GetChanges()
			Expect(changes).To(BeEmpty())
		})
	})

	Describe("SetTheme", func() {
		It("updates the theme", func() {
			modal.SetTheme(themes.NewDefaultTheme())
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Show", func() {
		It("makes modal visible and resets state", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.IsCancelled()).To(BeTrue())

			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.IsConfirmed()).To(BeFalse())
			Expect(modal.IsCancelled()).To(BeFalse())
		})
	})

	Describe("Hide", func() {
		It("makes modal invisible", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})
})
