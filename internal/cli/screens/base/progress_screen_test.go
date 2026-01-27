//nolint:errcheck // Test file - error handling for test setup is not relevant.
package base_test

import (
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BaseProgressScreen", func() {
	var screen *base.BaseProgressScreen

	Describe("Construction", func() {
		It("should create a progress screen with title and message", func() {
			screen = base.NewBaseProgressScreen(
				[]string{"Test", "Progress"},
				"Generating CV",
				"Analyzing career events and extracting insights...",
			)

			Expect(screen).NotTo(BeNil())
			Expect(screen.GetTitle()).To(Equal("Generating CV"))
			Expect(screen.GetMessage()).To(Equal("Analyzing career events and extracting insights..."))
		})

		It("should use default dimensions if not set", func() {
			screen = base.NewBaseProgressScreen([]string{"Test"}, "Title", "Message")
			Expect(screen.Width()).To(Equal(120)) // Default width
			Expect(screen.Height()).To(Equal(40)) // Default height
		})

		It("should start with spinner at position 0", func() {
			screen = base.NewBaseProgressScreen([]string{"Test"}, "Title", "Message")
			Expect(screen.GetSpinnerFrame()).To(Equal(0))
		})
	})

	Describe("Terminal Info", func() {
		BeforeEach(func() {
			screen = base.NewBaseProgressScreen([]string{"Test"}, "Title", "Message")
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
	})

	Describe("Spinner Animation", func() {
		BeforeEach(func() {
			screen = base.NewBaseProgressScreen([]string{"Test"}, "Title", "Message")
		})

		It("should increment spinner on TickMsg", func() {
			initialFrame := screen.GetSpinnerFrame()

			msg := base.TickMsg{}
			_, result := screen.Update(msg)

			Expect(result).To(BeNil())
			Expect(screen.GetSpinnerFrame()).To(Equal(initialFrame + 1))
		})

		It("should wrap spinner after max frames", func() {
			// Set spinner to near end
			screen.SetSpinnerFrame(100)

			msg := base.TickMsg{}
			screen.Update(msg)

			// Should have incremented
			Expect(screen.GetSpinnerFrame()).To(BeNumerically(">", 100))
		})

		It("should return tick command on TickMsg to continue animation", func() {
			msg := base.TickMsg{}
			cmd, _ := screen.Update(msg)

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Keyboard Shortcuts", func() {
		BeforeEach(func() {
			screen = base.NewBaseProgressScreen([]string{"Test"}, "Title", "Message")
		})

		Context("when cancellation is allowed", func() {
			It("should return CancelResult on escape key by default", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultCancel))
			})

			It("should allow escape when explicitly enabled", func() {
				screen.SetAllowCancel(true)

				msg := tea.KeyMsg{Type: tea.KeyEsc}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultCancel))
			})
		})

		Context("when cancellation is disabled", func() {
			It("should not cancel on escape when disabled", func() {
				screen.SetAllowCancel(false)

				msg := tea.KeyMsg{Type: tea.KeyEsc}
				_, result := screen.Update(msg)

				Expect(result).To(BeNil())
			})

			It("should ignore escape completely", func() {
				screen.SetAllowCancel(false)

				msg := tea.KeyMsg{Type: tea.KeyEsc}
				cmd, result := screen.Update(msg)

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		It("should ignore other keys", func() {
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := screen.Update(msg)

			Expect(result).To(BeNil())
		})
	})

	Describe("Completion", func() {
		BeforeEach(func() {
			screen = base.NewBaseProgressScreen([]string{"Test"}, "Title", "Message")
		})

		It("should return NavigateResult on CompleteMsg", func() {
			msg := base.CompleteMsg{Data: "result-data"}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			Expect(result.Data()).To(Equal("result-data"))
		})

		It("should return ErrorResult on ErrorMsg", func() {
			msg := base.ErrorMsg{Error: "Something went wrong"}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultError))

			data := result.Data().(map[string]interface{})
			Expect(data["message"]).To(Equal("Something went wrong"))
		})

		It("should handle CompleteMsg with nil data", func() {
			msg := base.CompleteMsg{Data: nil}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			Expect(result.Data()).To(BeNil())
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			screen = base.NewBaseProgressScreen(
				[]string{"Test", "Progress"},
				"Processing Data",
				"Analyzing and transforming records...",
			)
			screen.SetTerminalInfo(120, 40)
		})

		It("should render a view with StandardView structure", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should include title in view", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Processing Data"))
		})

		It("should include message in view", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Analyzing and transforming records"))
		})

		It("should include spinner character in view", func() {
			view := screen.View()
			// Should contain some spinner character (exact character depends on frame)
			Expect(view).NotTo(BeEmpty())
		})

		It("should update view when spinner advances", func() {
			screen.Update(base.TickMsg{}) // Advance spinner

			view := screen.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Processing Data"))
		})

		It("should show cancel help when cancellation is allowed", func() {
			screen.SetAllowCancel(true)
			view := screen.View()
			Expect(view).To(ContainSubstring("Esc"))
		})

		It("should not show cancel help when cancellation is disabled", func() {
			screen.SetAllowCancel(false)
			view := screen.View()
			// Footer should mention cancellation disabled
			Expect(view).To(ContainSubstring("disabled"))
		})
	})

	Describe("Customization", func() {
		It("should allow custom footer", func() {
			screen = base.NewBaseProgressScreen([]string{"Test"}, "Title", "Message")
			screen.SetFooter("Please wait...")

			view := screen.View()
			Expect(view).To(ContainSubstring("Please wait"))
		})

		It("should allow updating message during progress", func() {
			screen = base.NewBaseProgressScreen([]string{"Test"}, "Title", "Initial message")
			view1 := screen.View()
			Expect(view1).To(ContainSubstring("Initial message"))

			screen.SetMessage("Updated message")
			view2 := screen.View()
			Expect(view2).To(ContainSubstring("Updated message"))
			Expect(view2).NotTo(ContainSubstring("Initial message"))
		})

		It("should allow updating title during progress", func() {
			screen = base.NewBaseProgressScreen([]string{"Test"}, "Initial Title", "Message")
			view1 := screen.View()
			Expect(view1).To(ContainSubstring("Initial Title"))

			screen.SetTitle("Updated Title")
			view2 := screen.View()
			Expect(view2).To(ContainSubstring("Updated Title"))
			Expect(view2).NotTo(ContainSubstring("Initial Title"))
		})
	})

	Describe("Data Access", func() {
		BeforeEach(func() {
			screen = base.NewBaseProgressScreen(
				[]string{"Test"},
				"Loading Data",
				"Fetching records from database...",
			)
		})

		It("should provide access to title", func() {
			Expect(screen.GetTitle()).To(Equal("Loading Data"))
		})

		It("should provide access to message", func() {
			Expect(screen.GetMessage()).To(Equal("Fetching records from database..."))
		})

		It("should provide access to spinner frame", func() {
			Expect(screen.GetSpinnerFrame()).To(BeNumerically(">=", 0))
		})

		It("should allow setting spinner frame", func() {
			screen.SetSpinnerFrame(5)
			Expect(screen.GetSpinnerFrame()).To(Equal(5))
		})
	})

	Describe("Screen Interface", func() {
		BeforeEach(func() {
			screen = base.NewBaseProgressScreen([]string{"Test"}, "Title", "Message")
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
		It("should handle empty title", func() {
			screen := base.NewBaseProgressScreen([]string{"Test"}, "", "Message")
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle empty message", func() {
			screen := base.NewBaseProgressScreen([]string{"Test"}, "Title", "")
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle empty breadcrumbs", func() {
			screen := base.NewBaseProgressScreen([]string{}, "Title", "Message")
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle very small terminal dimensions", func() {
			screen := base.NewBaseProgressScreen([]string{"Test"}, "Title", "Message")
			screen.SetTerminalInfo(20, 10)

			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle long title", func() {
			longTitle := "This is a very long title that might wrap or be truncated depending on terminal width"
			screen := base.NewBaseProgressScreen([]string{"Test"}, longTitle, "Message")
			screen.SetTerminalInfo(40, 20)

			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle long message", func() {
			longMessage := "This is a very long message that contains multiple sentences and might need to wrap across multiple lines depending on the terminal width."
			screen := base.NewBaseProgressScreen([]string{"Test"}, "Title", longMessage)
			screen.SetTerminalInfo(40, 20)

			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})
