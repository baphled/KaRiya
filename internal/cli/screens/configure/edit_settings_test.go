package configure_test

import (
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/configure"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EditSettingsScreen", func() {
	var (
		screen   *configure.EditSettingsScreen
		settings []*intents.ConfigurationSetting
	)

	BeforeEach(func() {
		// Create sample settings with different types
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
			// No changes yet, so should be empty
			Expect(changes).To(BeEmpty())
		})
	})

	Describe("Init", func() {
		It("returns a command for form initialization", func() {
			cmd := screen.Init()
			// Init may return cmd for form initialization or nil
			// Just verify it doesn't panic
			_ = cmd
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

				// Tab should be handled by form, not return result
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
			// System domain should be shown
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
			// View should still render correctly
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("SetTheme", func() {
		It("accepts theme interface", func() {
			// Should not panic with nil theme
			screen.SetTheme(nil)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("GetChanges", func() {
		It("returns empty map when no changes made", func() {
			changes := screen.GetChanges()
			Expect(changes).To(BeEmpty())
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
				map[string]string{}, // no labels
			)
			view := screenWithoutLabels.View()
			Expect(view).To(ContainSubstring("some_key"))
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
		It("accepts theme interface", func() {
			screen.SetTheme(nil)
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

		Context("when pressing 'n' for no", func() {
			It("returns NavigateResult with false", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).To(BeFalse())
			})
		})

		Context("when toggling selection", func() {
			It("toggles with left/right keys", func() {
				// Start at No (false)
				Expect(screen.GetSelection()).To(BeFalse())

				// Toggle to Yes
				screen.Update(tea.KeyMsg{Type: tea.KeyLeft})
				Expect(screen.GetSelection()).To(BeTrue())

				// Toggle back to No
				screen.Update(tea.KeyMsg{Type: tea.KeyRight})
				Expect(screen.GetSelection()).To(BeFalse())
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

	Describe("GetErrorMessage", func() {
		It("returns the error message", func() {
			Expect(screen.GetErrorMessage()).To(Equal("Connection timed out"))
		})
	})
})
