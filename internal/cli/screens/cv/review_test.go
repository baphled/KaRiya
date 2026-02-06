package cv_test

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/cv"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("ReviewScreen", func() {
	var (
		screen   *cv.ReviewScreen
		testCV   *career.CVView
		sections []*career.CVSection
	)

	BeforeEach(func() {
		summarySection := fixtures.CVSectionWithSummary("section-1", "cv-1", "Experienced engineer with 10+ years...")
		summarySection.Title = "Professional Summary"

		experienceSection := fixtures.CVSectionWithContent("section-2", "cv-1", []*career.SectionContentGroup{
			fixtures.ContentGroupFull("Senior Engineer at TechCorp", "2020-01", "Present", []*career.CVBullet{
				fixtures.CVBulletWith("b-1", "section-2", "Led team of 5 engineers"),
				fixtures.CVBulletWith("b-2", "section-2", "Delivered critical features"),
			}),
			fixtures.ContentGroupFull("Engineer at StartupCo", "2018-01", "2019-12", []*career.CVBullet{
				fixtures.CVBulletWith("b-3", "section-2", "Built core platform"),
			}),
		})
		experienceSection.Title = "Professional Experience"

		skillsSection := fixtures.CVSectionWithContent("section-3", "cv-1", []*career.SectionContentGroup{
			fixtures.ContentGroupWithBullets("Languages", []*career.CVBullet{
				fixtures.CVBulletWith("b-4", "section-3", "Go, Python, TypeScript"),
			}),
		})
		skillsSection.Title = "Technical Skills"

		sections = []*career.CVSection{summarySection, experienceSection, skillsSection}

		testCV = fixtures.CVViewWith("cv-1", "Software Engineer CV", "Senior Engineer", "Hiring Manager")
		testCV.Sections = sections
		testCV.SourceEventCount = 15
		testCV.SourceFactCount = 42

		screen = cv.NewCVReviewScreen(testCV)
		screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	})

	Describe("NewCVReviewScreen", func() {
		It("should create a new screen with CV data", func() {
			s := cv.NewCVReviewScreen(testCV)
			Expect(s).NotTo(BeNil())
		})

		It("should accept nil CV", func() {
			s := cv.NewCVReviewScreen(nil)
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

		Context("preview navigation", func() {
			It("should navigate to preview on Enter key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).To(Equal("preview"))
			})

			It("should navigate to preview on 'p' key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})

				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).To(Equal("preview"))
			})
		})

		Context("export key", func() {
			It("should return NavigateResult with 'export' on 'x' key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).To(Equal("export"))
			})
		})

		Context("edit key", func() {
			It("should return NavigateResult with 'edit' on 'e' key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).To(Equal("edit"))
			})
		})

		Context("viewport navigation", func() {
			It("should go to top on 'g' key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should go to bottom on 'G' key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle 'k' key for scrolling up", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

				Expect(result).To(BeNil())
				_ = cmd
			})

			It("should handle 'j' key for scrolling down", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

				Expect(result).To(BeNil())
				_ = cmd
			})

			It("should handle up arrow key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyUp})

				Expect(result).To(BeNil())
				_ = cmd
			})

			It("should handle down arrow key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyDown})

				Expect(result).To(BeNil())
				_ = cmd
			})

			It("should handle page up key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyPgUp})

				Expect(result).To(BeNil())
				_ = cmd
			})

			It("should handle page down key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyPgDown})

				Expect(result).To(BeNil())
				_ = cmd
			})

			It("should handle ctrl+u for half page up", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyCtrlU})

				Expect(result).To(BeNil())
				_ = cmd
			})

			It("should handle ctrl+d for half page down", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyCtrlD})

				Expect(result).To(BeNil())
				_ = cmd
			})
		})

		Context("unhandled messages", func() {
			It("should return nil for unhandled key messages", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should return nil for mouse messages", func() {
				cmd, result := screen.Update(tea.MouseMsg{})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		It("should display CV Review title", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("CV Review"))
		})

		It("should display CV name", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Software Engineer CV"))
		})

		It("should display target role", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Senior Engineer"))
		})

		It("should display target audience", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Hiring Manager"))
		})

		It("should display source counts", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("15")) // event count
			Expect(view).To(ContainSubstring("42")) // fact count
		})

		It("should display section count", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("3")) // 3 sections
		})

		It("should display section titles as summary", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Professional Summary"))
			Expect(view).To(ContainSubstring("Professional Experience"))
			Expect(view).To(ContainSubstring("Technical Skills"))
		})

		It("should display bullet count per section", func() {
			view := screen.View()
			// Experience has 3 bullets, Skills has 1 bullet
			Expect(view).To(ContainSubstring("3 bullets"))
			Expect(view).To(ContainSubstring("1 bullet"))
		})

		It("should display help text with preview option", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("preview"))
			Expect(view).To(ContainSubstring("export"))
			Expect(view).To(ContainSubstring("back"))
		})
	})

	Describe("Edge Cases", func() {
		Context("nil CV", func() {
			It("should handle nil CV gracefully", func() {
				nilScreen := cv.NewCVReviewScreen(nil)

				view := nilScreen.View()
				Expect(view).To(ContainSubstring("No CV data available"))
			})
		})

		Context("CV with no sections", func() {
			It("should handle CV with no sections", func() {
				emptyCV := fixtures.CVViewWith("cv-empty", "Empty CV", "Engineer", "Manager")
				emptyCV.Sections = []*career.CVSection{}
				emptyScreen := cv.NewCVReviewScreen(emptyCV)
				emptyScreen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

				view := emptyScreen.View()
				Expect(view).To(ContainSubstring("0 sections"))
			})
		})

		Context("CV with empty sections", func() {
			It("should handle sections with no content", func() {
				emptySection := fixtures.CVSection("section-empty", "cv-empty-s")
				emptySection.Title = "Empty Section"
				emptySection.Content = []*career.SectionContentGroup{}
				cvWithEmptySections := fixtures.CVViewWith("cv-empty-s", "CV", "Engineer", "Manager")
				cvWithEmptySections.Sections = []*career.CVSection{emptySection}
				s := cv.NewCVReviewScreen(cvWithEmptySections)
				s.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

				view := s.View()
				Expect(view).To(ContainSubstring("Empty Section"))
				Expect(view).To(ContainSubstring("0 bullets"))
			})
		})
	})

	Describe("GetCV", func() {
		It("should return the CV data", func() {
			cvData := screen.GetCV()
			Expect(cvData).To(Equal(testCV))
		})

		It("should return nil for nil CV", func() {
			nilScreen := cv.NewCVReviewScreen(nil)
			cvData := nilScreen.GetCV()
			Expect(cvData).To(BeNil())
		})
	})

	Describe("Theme Support", func() {
		var view string

		Context("with default theme", func() {
			BeforeEach(func() {
				view = screen.View()
			})

			Context("section headers", func() {
				It("should display CV Review header", func() {
					Expect(view).To(ContainSubstring("CV Review"))
				})

				It("should display Statistics header", func() {
					Expect(view).To(ContainSubstring("Statistics"))
				})

				It("should display Sections header", func() {
					Expect(view).To(ContainSubstring("Sections"))
				})
			})

			Context("metadata labels", func() {
				It("should display Name label", func() {
					Expect(view).To(ContainSubstring("Name:"))
				})

				It("should display Role label", func() {
					Expect(view).To(ContainSubstring("Role:"))
				})

				It("should display Audience label", func() {
					Expect(view).To(ContainSubstring("Audience:"))
				})
			})

			Context("statistics labels", func() {
				It("should display Source Events label", func() {
					Expect(view).To(ContainSubstring("Source Events:"))
				})

				It("should display Source Facts label", func() {
					Expect(view).To(ContainSubstring("Source Facts:"))
				})

				It("should display Total Bullets label", func() {
					Expect(view).To(ContainSubstring("Total Bullets:"))
				})
			})
		})

		Context("with custom theme", func() {
			BeforeEach(func() {
				theme := themes.NewDefaultTheme()
				screen.SetTheme(theme)
				view = screen.View()
			})

			It("should render with provided theme", func() {
				Expect(view).To(ContainSubstring("CV Review"))
			})
		})
	})

	Describe("Profile Configuration", func() {
		Context("with custom profile config", func() {
			It("should display custom name from profile config", func() {
				customProfile := &config.ProfileConfig{
					Name:     "Test User",
					Email:    "test@example.com",
					Title:    "Staff Engineer",
					Location: "Test City, TC",
				}
				screenWithProfile := cv.NewCVReviewScreenWithProfile(testCV, customProfile)

				view := screenWithProfile.View()
				Expect(view).To(ContainSubstring("Test User"))
			})

			It("should display custom email from profile config", func() {
				customProfile := &config.ProfileConfig{
					Name:  "Test User",
					Email: "test@example.com",
				}
				screenWithProfile := cv.NewCVReviewScreenWithProfile(testCV, customProfile)

				view := screenWithProfile.View()
				Expect(view).To(ContainSubstring("test@example.com"))
			})

			It("should display custom location from profile config", func() {
				customProfile := &config.ProfileConfig{
					Name:     "Test User",
					Location: "Test City, TC",
				}
				screenWithProfile := cv.NewCVReviewScreenWithProfile(testCV, customProfile)

				view := screenWithProfile.View()
				Expect(view).To(ContainSubstring("Test City, TC"))
			})
		})

		Context("with nil profile config", func() {
			It("should render without error when profile config is nil", func() {
				screenWithNilProfile := cv.NewCVReviewScreenWithProfile(testCV, nil)

				view := screenWithNilProfile.View()
				Expect(view).NotTo(BeEmpty())
				Expect(view).To(ContainSubstring("CV Review"))
			})
		})

		Context("with partial profile config", func() {
			It("should display provided fields only", func() {
				partialProfile := &config.ProfileConfig{
					Name: "Partial User",
				}
				screenWithPartial := cv.NewCVReviewScreenWithProfile(testCV, partialProfile)

				view := screenWithPartial.View()
				Expect(view).To(ContainSubstring("Partial User"))
			})
		})
	})
})
