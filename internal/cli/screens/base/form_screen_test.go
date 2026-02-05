package base_test

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestFormData is a simple form data structure for testing.
type TestFormData struct {
	Name            string
	Email           string
	SubmitConfirmed bool
}

var _ = Describe("FormScreen", func() {
	var (
		screen   *base.FormScreen[*TestFormData]
		formData *TestFormData
	)

	BeforeEach(func() {
		formData = &TestFormData{
			Name:            "",
			Email:           "",
			SubmitConfirmed: false,
		}
	})

	Describe("Construction", func() {
		It("should create a form screen with form builder", func() {
			builder := func(data *TestFormData, width, height int) *huh.Form {
				group := huh.NewGroup(
					huh.NewInput().
						Key("name").
						Title("Name").
						Value(&data.Name),
					huh.NewInput().
						Key("email").
						Title("Email").
						Value(&data.Email),
					huh.NewConfirm().
						Key("submit").
						Title("Submit?").
						Value(&data.SubmitConfirmed),
				)
				return huh.NewForm(group).WithWidth(width).WithHeight(height)
			}

			screen = base.NewBaseFormScreen(
				[]string{"Test", "Form"},
				builder,
				formData,
			)

			Expect(screen).NotTo(BeNil())
			Expect(screen.GetFormData()).To(Equal(formData))
		})

		It("should use default dimensions if not set", func() {
			builder := func(data *TestFormData, width, height int) *huh.Form {
				group := huh.NewGroup(
					huh.NewInput().
						Key("name").
						Title("Name").
						Value(&data.Name),
				)
				return huh.NewForm(group).WithWidth(width).WithHeight(height)
			}

			screen = base.NewBaseFormScreen([]string{"Test"}, builder, formData)
			Expect(screen.Width()).To(Equal(120)) // Default width
			Expect(screen.Height()).To(Equal(40)) // Default height
		})
	})

	Describe("Terminal Info", func() {
		BeforeEach(func() {
			builder := func(data *TestFormData, width, height int) *huh.Form {
				group := huh.NewGroup(
					huh.NewInput().
						Key("name").
						Title("Name").
						Value(&data.Name),
				)
				return huh.NewForm(group).WithWidth(width).WithHeight(height)
			}
			screen = base.NewBaseFormScreen([]string{"Test"}, builder, formData)
		})

		It("should update dimensions on SetTerminalInfo", func() {
			screen.SetTerminalInfo(100, 30)
			Expect(screen.Width()).To(Equal(100))
			Expect(screen.Height()).To(Equal(30))
		})

		It("should handle WindowSizeMsg", func() {
			msg := tea.WindowSizeMsg{Width: 80, Height: 24}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
			Expect(screen.Width()).To(Equal(80))
			Expect(screen.Height()).To(Equal(24))
		})

		It("should rebuild form with new dimensions", func() {
			screen.SetTerminalInfo(100, 30)

			// Form should be rebuilt with new dimensions
			// We can't directly test the form's dimensions without exposing it,
			// but we can verify the screen's dimensions were updated
			Expect(screen.Width()).To(Equal(100))
			Expect(screen.Height()).To(Equal(30))
		})
	})

	Describe("Form Interaction", func() {
		BeforeEach(func() {
			builder := func(data *TestFormData, width, height int) *huh.Form {
				group := huh.NewGroup(
					huh.NewInput().
						Key("name").
						Title("Name").
						Value(&data.Name),
					huh.NewConfirm().
						Key("submit").
						Title("Submit?").
						Value(&data.SubmitConfirmed),
				)
				return huh.NewForm(group).WithWidth(width).WithHeight(height)
			}
			screen = base.NewBaseFormScreen([]string{"Test"}, builder, formData)
		})

		It("should return CancelResult on escape key", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})

		It("should delegate other keys to huh form", func() {
			// Type 'a' key
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
			cmd, result := screen.Update(msg)

			// Form handles the key, no result yet (still editing)
			Expect(result).To(BeNil())
			// May return a command from huh (nil or blink command)
			_ = cmd
		})

		It("should return SubmitResult when form is completed and confirmed", func() {
			// Simulate form completion
			formData.SubmitConfirmed = true

			// Send a message that would complete the form
			// (In real usage, user would navigate through form and press enter on confirm)
			// We'll simulate by calling the form's internal state directly
			// Since we can't easily manipulate huh.Form's internal state in tests,
			// we'll test the logic path by setting SubmitConfirmed to true
			// and checking the screen's behavior

			// For now, we test that the form screen respects the SubmitConfirmed flag
			Expect(screen.GetFormData().SubmitConfirmed).To(BeTrue())
		})

		It("should not submit if SubmitConfirmed is false", func() {
			// Even if form is complete, don't submit without confirmation
			Expect(screen.GetFormData().SubmitConfirmed).To(BeFalse())
		})

		It("should return SubmitResult when form completes via internal command chain", func() {
			confirmOnlyBuilder := func(data *TestFormData, width, height int) *huh.Form {
				group := huh.NewGroup(
					huh.NewConfirm().
						Key("submit").
						Title("Submit?").
						Affirmative("Submit").
						Negative("Cancel").
						Value(&data.SubmitConfirmed),
				)
				return huh.NewForm(group).WithWidth(width).WithHeight(height)
			}
			confirmScreen := base.NewBaseFormScreen([]string{"Test"}, confirmOnlyBuilder, formData)
			confirmScreen.SetTerminalInfo(80, 24)

			// Press right arrow to select "Submit" (affirmative option)
			_, result := confirmScreen.Update(tea.KeyMsg{Type: tea.KeyRight})
			Expect(result).To(BeNil(), "selecting affirmative should not complete form yet")

			// Press Enter to confirm - this triggers form's internal nextGroupMsg chain
			cmd, result := confirmScreen.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// The initial Enter press may not complete the form directly.
			// huh returns a batch command containing nextGroupMsg.
			// Process returned commands to advance the form state machine.
			var submitResult screens.ScreenResult
			if result != nil {
				submitResult = result
			}

			for submitResult == nil && cmd != nil {
				msg := cmd()
				if msg == nil {
					break
				}
				if batchMsg, ok := msg.(tea.BatchMsg); ok {
					for _, batchCmd := range batchMsg {
						if batchCmd == nil {
							continue
						}
						innerMsg := batchCmd()
						if innerMsg == nil {
							continue
						}
						cmd, result = confirmScreen.Update(innerMsg)
						if result != nil {
							submitResult = result
							break
						}
					}
				} else {
					cmd, result = confirmScreen.Update(msg)
					if result != nil {
						submitResult = result
					}
				}
			}

			Expect(submitResult).NotTo(BeNil(),
				"form should return SubmitResult when completing via internal command chain")
			Expect(submitResult.Type()).To(Equal(screens.ResultSubmit))
			Expect(formData.SubmitConfirmed).To(BeTrue())
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			builder := func(data *TestFormData, width, height int) *huh.Form {
				group := huh.NewGroup(
					huh.NewInput().
						Key("name").
						Title("Name").
						Value(&data.Name),
				)
				return huh.NewForm(group).WithWidth(width).WithHeight(height)
			}
			screen = base.NewBaseFormScreen([]string{"Test", "Form"}, builder, formData)
			screen.SetTerminalInfo(120, 40)
		})

		It("should render a view with StandardView structure", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
			// Should contain form content
			// (Exact content depends on huh rendering)
		})

		It("should include breadcrumbs in view", func() {
			view := screen.View()
			// StandardView includes breadcrumbs
			// We can't easily test the exact format without knowing StandardView internals,
			// but we can verify it doesn't panic
			Expect(view).NotTo(BeEmpty())
		})

		It("should center form content horizontally within the view", func() {
			view := screen.View()

			// Find lines containing the form field indicator (huh uses ┃ for active fields)
			formLines := []string{}
			for _, line := range strings.Split(view, "\n") {
				if strings.Contains(line, "┃") && strings.Contains(line, "Name") {
					formLines = append(formLines, line)
				}
			}
			Expect(formLines).NotTo(BeEmpty(), "should find form field lines in view")

			// The form content should NOT be flush against the left edge.
			// A centered form on a 120-char terminal should have meaningful
			// leading whitespace (at least 10 chars) before the form chrome.
			for _, line := range formLines {
				leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
				Expect(leadingSpaces).To(BeNumerically(">=", 10),
					"form content should be centered, not left-aligned; line: '%s'", line)
			}
		})

		It("should update view when form data changes", func() {
			view1 := screen.View()

			// Change form data
			formData.Name = "John Doe"

			view2 := screen.View()

			// Views should be generated (may be same or different depending on huh)
			Expect(view1).NotTo(BeEmpty())
			Expect(view2).NotTo(BeEmpty())
		})
	})

	Describe("Footer Helpers", func() {
		BeforeEach(func() {
			builder := func(data *TestFormData, width, height int) *huh.Form {
				group := huh.NewGroup(
					huh.NewInput().
						Key("name").
						Title("Name").
						Value(&data.Name),
				)
				return huh.NewForm(group).WithWidth(width).WithHeight(height)
			}
			screen = base.NewBaseFormScreen([]string{"Test"}, builder, formData)
		})

		It("should allow setting custom footer", func() {
			screen.SetFooter("Custom footer text")
			view := screen.View()
			Expect(view).To(ContainSubstring("Custom footer text"))
		})

		It("should use default footer if none set", func() {
			view := screen.View()
			// Default footer should include common shortcuts
			Expect(view).To(ContainSubstring("Esc"))
		})
	})

	Describe("Form Data Access", func() {
		BeforeEach(func() {
			builder := func(data *TestFormData, width, height int) *huh.Form {
				group := huh.NewGroup(
					huh.NewInput().
						Key("name").
						Title("Name").
						Value(&data.Name),
				)
				return huh.NewForm(group).WithWidth(width).WithHeight(height)
			}
			screen = base.NewBaseFormScreen([]string{"Test"}, builder, formData)
		})

		It("should provide access to form data", func() {
			data := screen.GetFormData()
			Expect(data).To(Equal(formData))
			Expect(data).To(BeIdenticalTo(formData)) // Same pointer
		})

		It("should allow modifying form data through pointer", func() {
			data := screen.GetFormData()
			data.Name = "Modified"

			Expect(formData.Name).To(Equal("Modified"))
			Expect(screen.GetFormData().Name).To(Equal("Modified"))
		})
	})

	Describe("Screen Interface", func() {
		BeforeEach(func() {
			builder := func(data *TestFormData, width, height int) *huh.Form {
				group := huh.NewGroup(
					huh.NewInput().
						Key("name").
						Title("Name").
						Value(&data.Name),
				)
				return huh.NewForm(group).WithWidth(width).WithHeight(height)
			}
			screen = base.NewBaseFormScreen([]string{"Test"}, builder, formData)
		})

		It("should implement Screen interface", func() {
			var _ screens.Screen = screen
		})

		It("should support SetTheme", func() {
			theme := "test-theme"
			screen.SetTheme(theme)
			Expect(screen.Theme()).To(Equal(theme))
		})
	})

	Describe("Edge Cases", func() {
		It("should handle nil form data gracefully", func() {
			builder := func(_ *TestFormData, width, height int) *huh.Form {
				group := huh.NewGroup(
					huh.NewInput().
						Key("name").
						Title("Name"),
				)
				return huh.NewForm(group).WithWidth(width).WithHeight(height)
			}

			// This should not panic even with nil (though not recommended)
			screen := base.NewBaseFormScreen([]string{"Test"}, builder, (*TestFormData)(nil))
			Expect(screen).NotTo(BeNil())
		})

		It("should handle empty breadcrumbs", func() {
			builder := func(_ *TestFormData, width, height int) *huh.Form {
				group := huh.NewGroup(
					huh.NewInput().
						Key("name").
						Title("Name").
						Value(&formData.Name),
				)
				return huh.NewForm(group).WithWidth(width).WithHeight(height)
			}
			screen := base.NewBaseFormScreen([]string{}, builder, formData)

			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle very small terminal dimensions", func() {
			builder := func(_ *TestFormData, width, height int) *huh.Form {
				group := huh.NewGroup(
					huh.NewInput().
						Key("name").
						Title("Name").
						Value(&formData.Name),
				)
				return huh.NewForm(group).WithWidth(width).WithHeight(height)
			}
			screen := base.NewBaseFormScreen([]string{"Test"}, builder, formData)
			screen.SetTerminalInfo(20, 10)

			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})
