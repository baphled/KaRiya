package cv_test

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/cv"
)

var _ = Describe("AudienceSelectScreen", func() {
	var (
		screen    *cv.AudienceSelectScreen
		audiences []*cv.AudienceOption
	)

	BeforeEach(func() {
		audiences = []*cv.AudienceOption{
			{
				ID:          "hiring_manager",
				Name:        "Hiring Manager",
				Description: "Technical hiring managers at companies",
			},
			{
				ID:          "recruiter",
				Name:        "Recruiter",
				Description: "Internal or external recruiters",
			},
			{
				ID:          "ats",
				Name:        "ATS System",
				Description: "Applicant tracking system optimization",
			},
		}

		screen = cv.NewCVAudienceSelectScreen(audiences)
	})

	Describe("NewCVAudienceSelectScreen", func() {
		It("should create a new screen", func() {
			Expect(screen).NotTo(BeNil())
		})

		It("should accept audiences", func() {
			s := cv.NewCVAudienceSelectScreen(audiences)
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
		Context("navigation", func() {
			It("should navigate down with 'j' key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should navigate down with arrow down", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyDown})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should navigate up with 'k' key", func() {
				screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should navigate up with arrow up", func() {
				screen.Update(tea.KeyMsg{Type: tea.KeyDown})

				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyUp})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should jump to top with 'g' key", func() {
				screen.Update(tea.KeyMsg{Type: tea.KeyDown})
				screen.Update(tea.KeyMsg{Type: tea.KeyDown})

				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should jump to bottom with 'G' key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("selection", func() {
			It("should return NavigateResult on Enter key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).NotTo(BeNil())

				selectedAudience, ok := result.Data().(*cv.AudienceOption)
				Expect(ok).To(BeTrue())
				Expect(selectedAudience.ID).To(Equal("hiring_manager"))
			})

			It("should return selected audience after navigation", func() {
				screen.Update(tea.KeyMsg{Type: tea.KeyDown})

				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())

				selectedAudience, ok := result.Data().(*cv.AudienceOption)
				Expect(ok).To(BeTrue())
				Expect(selectedAudience.ID).To(Equal("recruiter"))
			})
		})

		Context("cancellation", func() {
			It("should return CancelResult on Esc key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultCancel))
			})
		})

		Context("window resize", func() {
			It("should handle window resize", func() {
				cmd, result := screen.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		It("should render audience list", func() {
			view := screen.View()

			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Hiring Manager"))
			Expect(view).To(ContainSubstring("Recruiter"))
			Expect(view).To(ContainSubstring("ATS System"))
		})

		It("should show selection indicator", func() {
			view := screen.View()

			Expect(view).To(MatchRegexp(`[▶►>]`))
		})

		It("should show audience descriptions", func() {
			view := screen.View()

			Expect(view).To(ContainSubstring("Technical hiring managers"))
			Expect(view).To(ContainSubstring("Internal or external recruiters"))
		})

		It("should show breadcrumbs", func() {
			view := screen.View()

			Expect(view).To(ContainSubstring("Main Menu"))
			Expect(view).To(ContainSubstring("Generate CV"))
			Expect(view).To(ContainSubstring("Select Audience"))
		})

		It("should show title", func() {
			view := screen.View()

			Expect(view).To(ContainSubstring("Select Target Audience"))
		})
	})

	Describe("Edge Cases", func() {
		Context("empty audience list", func() {
			It("should handle empty audience list gracefully", func() {
				emptyScreen := cv.NewCVAudienceSelectScreen([]*cv.AudienceOption{})

				view := emptyScreen.View()
				Expect(view).To(ContainSubstring("No items available"))
			})
		})

		Context("single audience", func() {
			It("should work with single audience", func() {
				singleAudience := []*cv.AudienceOption{audiences[0]}
				singleScreen := cv.NewCVAudienceSelectScreen(singleAudience)

				cmd, result := singleScreen.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
			})
		})
	})
})
