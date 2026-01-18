package cv_test

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/cv"
)

var _ = Describe("CVGeneratingScreen", func() {
	var screen *cv.CVGeneratingScreen

	BeforeEach(func() {
		screen = cv.NewCVGeneratingScreen("Senior Engineer", "Hiring Manager")
	})

	Describe("NewCVGeneratingScreen", func() {
		It("should create a new screen with profile and audience", func() {
			s := cv.NewCVGeneratingScreen("Staff Engineer", "Recruiter")
			Expect(s).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		It("should return nil command", func() {
			cmd := screen.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update", func() {
		Context("window resize", func() {
			It("should handle window resize", func() {
				cmd, result := screen.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("escape key", func() {
			It("should return CancelResult on Esc key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultCancel))
			})
		})

		Context("other keys", func() {
			It("should ignore other key presses", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		It("should display spinner", func() {
			view := screen.View()
			Expect(view).To(MatchRegexp("[⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏]"))
		})

		It("should display generating message", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Generating CV"))
		})

		It("should display profile", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Profile: Senior Engineer"))
		})

		It("should display audience", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Audience: Hiring Manager"))
		})

		It("should display cancel hint", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("esc"))
		})
	})
})
