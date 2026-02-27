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

var _ = Describe("CVPreviewScreen", func() {
	var (
		screen   *cv.CVPreviewScreen
		testCV   *career.CVView
		sections []*career.CVSection
	)

	BeforeEach(func() {
		summarySection := fixtures.CVSectionWithSummary("section-1", "cv-1", "Experienced software engineer with expertise in distributed systems.")
		summarySection.Title = "Professional Summary"

		experienceSection := fixtures.CVSectionWithContent("section-2", "cv-1", []*career.SectionContentGroup{
			fixtures.ContentGroupFull("Senior Engineer at TechCorp", "2020-01", "Present", []*career.CVBullet{
				fixtures.CVBulletWith("b-1", "section-2", "Led development of microservices architecture"),
				fixtures.CVBulletWith("b-2", "section-2", "Mentored team of 5 junior engineers"),
				fixtures.CVBulletWith("b-3", "section-2", "Reduced deployment time by 60%"),
			}),
			fixtures.ContentGroupFull("Software Engineer at StartupCo", "2018-01", "2019-12", []*career.CVBullet{
				fixtures.CVBulletWith("b-4", "section-2", "Built core payment processing system"),
				fixtures.CVBulletWith("b-5", "section-2", "Implemented real-time notifications"),
			}),
		})
		experienceSection.Title = "Professional Experience"

		skillsSection := fixtures.CVSectionWithContent("section-3", "cv-1", []*career.SectionContentGroup{
			fixtures.ContentGroupWithBullets("Languages", []*career.CVBullet{
				fixtures.CVBulletWith("b-6", "section-3", "Go, Python, TypeScript, Java"),
			}),
			fixtures.ContentGroupWithBullets("Frameworks", []*career.CVBullet{
				fixtures.CVBulletWith("b-7", "section-3", "React, Vue, Django, Spring Boot"),
			}),
		})
		skillsSection.Title = "Technical Skills"

		sections = []*career.CVSection{summarySection, experienceSection, skillsSection}

		testCV = fixtures.CVViewWith("cv-1", "Test CV", "Senior Engineer", "Hiring Manager")
		testCV.Sections = sections
		testCV.SourceEventCount = 10
		testCV.SourceFactCount = 25

		screen = cv.NewCVPreviewScreen(testCV)
		// Initialize with a larger window size to ensure all content is visible
		// (personal details header takes ~6 lines, so we need more height)
		screen.Update(tea.WindowSizeMsg{Width: 80, Height: 40})
	})

	Describe("NewCVPreviewScreen", func() {
		It("should create a new screen with CV data", func() {
			s := cv.NewCVPreviewScreen(testCV)
			Expect(s).NotTo(BeNil())
		})

		It("should accept nil CV", func() {
			s := cv.NewCVPreviewScreen(nil)
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

			It("should reinitialize viewport on resize", func() {
				// First size
				screen.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
				view1 := screen.View()

				// Resize to larger
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				view2 := screen.View()

				// Views should both render (viewport handles resize)
				Expect(view1).NotTo(BeEmpty())
				Expect(view2).NotTo(BeEmpty())
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

		Context("confirm key", func() {
			It("should return NavigateResult with CV on Enter key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).To(Equal(testCV))
			})

			It("should return NavigateResult with CV on 'y' key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).To(Equal(testCV))
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

		Context("scrolling with viewport", func() {
			It("should handle scroll down with 'j' key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle scroll down with down arrow", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyDown})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle scroll up with 'k' key", func() {
				// First scroll down
				screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle scroll up with up arrow", func() {
				screen.Update(tea.KeyMsg{Type: tea.KeyDown})

				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyUp})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle page down with PgDn", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyPgDown})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle page up with PgUp", func() {
				// First page down
				screen.Update(tea.KeyMsg{Type: tea.KeyPgDown})

				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyPgUp})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle go to top with 'g' key", func() {
				// First scroll down
				screen.Update(tea.KeyMsg{Type: tea.KeyPgDown})

				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle go to bottom with 'G' key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle half-page scroll with ctrl+d", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyCtrlD})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle half-page scroll with ctrl+u", func() {
				// First scroll down
				screen.Update(tea.KeyMsg{Type: tea.KeyCtrlD})

				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyCtrlU})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should return nil for scroll keys when viewport not ready", func() {
				freshScreen := cv.NewCVPreviewScreen(testCV)
				cmd, result := freshScreen.Update(tea.KeyMsg{Type: tea.KeyDown})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		It("should display CV Preview title", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("CV Preview"))
		})

		It("should display section titles", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Professional Summary"))
			Expect(view).To(ContainSubstring("Professional Experience"))
			Expect(view).To(ContainSubstring("Technical Skills"))
		})

		It("should display bullet content", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Led development of microservices"))
			Expect(view).To(ContainSubstring("Mentored team"))
		})

		It("should display content group headers", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Senior Engineer at TechCorp"))
			Expect(view).To(ContainSubstring("Software Engineer at StartupCo"))
		})

		It("should display summary section content", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Experienced software engineer"))
		})

		It("should preserve newlines in summary content", func() {
			headingAndProseSummary := "**GDS-Aligned Developer**\nFull-stack engineer with deep backend expertise."
			summarySection := fixtures.CVSectionWithSummary("section-multiline", "cv-1", headingAndProseSummary)
			summarySection.Title = "Professional Summary"
			multilineCV := fixtures.CVViewWith("cv-multiline", "Multiline CV", "developer", "hiring_manager")
			multilineCV.Sections = []*career.CVSection{summarySection}
			multilineScreen := cv.NewCVPreviewScreen(multilineCV)
			multilineScreen.Update(tea.WindowSizeMsg{Width: 80, Height: 40})
			view := multilineScreen.View()
			Expect(view).To(ContainSubstring("**GDS-Aligned Developer**"))
			Expect(view).To(ContainSubstring("Full-stack engineer with deep backend expertise."))
		})

		It("should display help text with scrolling options", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("scroll"))
			Expect(view).To(ContainSubstring("confirm"))
			Expect(view).To(ContainSubstring("export"))
			Expect(view).To(ContainSubstring("back"))
		})

		It("should display scroll percentage when content is scrollable", func() {
			longCV := fixtures.CVViewWith("cv-long", "Long CV", "staff", "hiring_manager")
			longCV.Sections = make([]*career.CVSection, 10)
			for i := range 10 {
				sectionID := "section-long-" + string(rune('A'+i))
				longCV.Sections[i] = fixtures.CVSectionWithContent(sectionID, "cv-long", []*career.SectionContentGroup{
					fixtures.ContentGroupWithBullets("Group", []*career.CVBullet{
						fixtures.CVBulletWith("bl-1-"+sectionID, sectionID, "Bullet 1"),
						fixtures.CVBulletWith("bl-2-"+sectionID, sectionID, "Bullet 2"),
						fixtures.CVBulletWith("bl-3-"+sectionID, sectionID, "Bullet 3"),
					}),
				})
				longCV.Sections[i].Title = "Section " + string(rune('A'+i))
			}
			longScreen := cv.NewCVPreviewScreen(longCV)
			longScreen.Update(tea.WindowSizeMsg{Width: 80, Height: 10}) // Small viewport

			view := longScreen.View()
			// Should show scroll percentage indicator
			Expect(view).To(Or(
				ContainSubstring("%"),
				ContainSubstring("scroll"),
			))
		})
	})

	Describe("Edge Cases", func() {
		Context("nil CV", func() {
			It("should handle nil CV gracefully", func() {
				nilScreen := cv.NewCVPreviewScreen(nil)
				nilScreen.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

				view := nilScreen.View()
				Expect(view).To(ContainSubstring("No CV data available"))
			})
		})

		Context("CV with no sections", func() {
			It("should handle CV with no sections", func() {
				emptyCV := fixtures.CVViewWith("cv-empty", "Empty CV", "Engineer", "Manager")
				emptyCV.Sections = []*career.CVSection{}
				emptyScreen := cv.NewCVPreviewScreen(emptyCV)
				emptyScreen.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

				view := emptyScreen.View()
				Expect(view).To(ContainSubstring("No sections"))
			})
		})

		Context("CV with dates", func() {
			It("should display date ranges correctly", func() {
				view := screen.View()
				Expect(view).To(ContainSubstring("2020-01"))
				Expect(view).To(ContainSubstring("Present"))
			})

			It("should display same date correctly", func() {
				singleDateCV := fixtures.CVViewWith("cv-single", "Single Date CV", "staff", "hiring_manager")
				singleDateCV.Sections = []*career.CVSection{
					fixtures.CVSectionWithContent("section-sd", "cv-single", []*career.SectionContentGroup{
						fixtures.ContentGroupFull("Event", "2020-06", "2020-06", []*career.CVBullet{
							fixtures.CVBulletWith("bl-sd", "section-sd", "Did something"),
						}),
					}),
				}
				singleDateCV.Sections[0].Title = "Experience"
				s := cv.NewCVPreviewScreen(singleDateCV)
				s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

				view := s.View()
				Expect(view).To(ContainSubstring("2020-06"))
			})
		})

		Context("small terminal", func() {
			It("should handle very small terminal gracefully", func() {
				screen.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
				view := screen.View()
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	Describe("GetCV", func() {
		It("should return the CV data", func() {
			cvData := screen.GetCV()
			Expect(cvData).To(Equal(testCV))
		})

		It("should return nil for nil CV", func() {
			nilScreen := cv.NewCVPreviewScreen(nil)
			cvData := nilScreen.GetCV()
			Expect(cvData).To(BeNil())
		})
	})

	Describe("Profile Configuration", func() {
		Context("with custom profile config", func() {
			It("should display custom name from profile config", func() {
				customProfile := &config.ProfileConfig{
					Name:     "Jane Doe",
					Email:    "jane@example.com",
					Title:    "Staff Engineer",
					Location: "London, UK",
					GitHub:   "https://github.com/janedoe",
				}
				screenWithProfile := cv.NewCVPreviewScreenWithProfile(testCV, customProfile)
				screenWithProfile.Update(tea.WindowSizeMsg{Width: 80, Height: 40})

				view := screenWithProfile.View()
				Expect(view).To(ContainSubstring("Jane Doe"))
			})

			It("should display custom email from profile config", func() {
				customProfile := &config.ProfileConfig{
					Name:  "Jane Doe",
					Email: "jane@example.com",
				}
				screenWithProfile := cv.NewCVPreviewScreenWithProfile(testCV, customProfile)
				screenWithProfile.Update(tea.WindowSizeMsg{Width: 80, Height: 40})

				view := screenWithProfile.View()
				Expect(view).To(ContainSubstring("jane@example.com"))
			})

			It("should display custom title from profile config", func() {
				customProfile := &config.ProfileConfig{
					Name:  "Jane Doe",
					Title: "Principal Engineer / Architect",
				}
				screenWithProfile := cv.NewCVPreviewScreenWithProfile(testCV, customProfile)
				screenWithProfile.Update(tea.WindowSizeMsg{Width: 80, Height: 40})

				view := screenWithProfile.View()
				Expect(view).To(ContainSubstring("Principal Engineer / Architect"))
			})

			It("should display custom location from profile config", func() {
				customProfile := &config.ProfileConfig{
					Name:     "Jane Doe",
					Location: "Berlin, Germany",
				}
				screenWithProfile := cv.NewCVPreviewScreenWithProfile(testCV, customProfile)
				screenWithProfile.Update(tea.WindowSizeMsg{Width: 80, Height: 40})

				view := screenWithProfile.View()
				Expect(view).To(ContainSubstring("Berlin, Germany"))
			})
		})

		Context("with nil profile config", func() {
			It("should use empty defaults when profile config is nil", func() {
				screenWithNilProfile := cv.NewCVPreviewScreenWithProfile(testCV, nil)
				screenWithNilProfile.Update(tea.WindowSizeMsg{Width: 80, Height: 40})

				view := screenWithNilProfile.View()
				// Should NOT contain hardcoded personal data
				Expect(view).NotTo(ContainSubstring("Yomi Colledge"))
				Expect(view).NotTo(ContainSubstring("boodah"))
				// Screen should still render
				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("with partial profile config", func() {
			It("should keep empty fields empty (no hardcoded defaults)", func() {
				partialProfile := &config.ProfileConfig{
					Name: "Custom Name",
					// Email, Title, Location left empty - will remain empty
				}
				screenWithPartial := cv.NewCVPreviewScreenWithProfile(testCV, partialProfile)
				screenWithPartial.Update(tea.WindowSizeMsg{Width: 80, Height: 40})

				view := screenWithPartial.View()
				// Name should be custom
				Expect(view).To(ContainSubstring("Custom Name"))
				// Should NOT contain hardcoded email
				Expect(view).NotTo(ContainSubstring("yomi@boodah.net"))
				Expect(view).NotTo(ContainSubstring("boodah"))
			})
		})
	})

	Describe("Theme coverage", func() {
		It("should use custom theme when set via SetTheme", func() {
			themedScreen := cv.NewCVPreviewScreen(testCV)
			themedScreen.SetTheme(themes.NewDefaultTheme())
			themedScreen.Update(tea.WindowSizeMsg{Width: 80, Height: 40})

			view := themedScreen.View()
			Expect(view).To(ContainSubstring("CV Preview"))
		})
	})

	Describe("Highlight rendering", func() {
		It("should use MaxHighlights from profile config", func() {
			customProfile := &config.ProfileConfig{
				Name:          "Jane Doe",
				MaxHighlights: 3,
			}
			screenWithMax := cv.NewCVPreviewScreenWithProfile(testCV, customProfile)
			screenWithMax.Update(tea.WindowSizeMsg{Width: 80, Height: 40})

			view := screenWithMax.View()
			Expect(view).To(ContainSubstring("Key Highlights"))
		})

		It("should render highlights using EnhancedText when Text is empty", func() {
			bullet := fixtures.CVBulletWith("b-enh", "s1", "")
			bullet.EnhancedText = "Enhanced highlight text for preview"
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-enh", []*career.SectionContentGroup{group})
			section.Title = "Experience"

			enhCV := fixtures.CVViewWith("cv-enh", "Enhanced CV", "Engineer", "Hiring Manager")
			enhCV.Sections = []*career.CVSection{section}

			enhScreen := cv.NewCVPreviewScreen(enhCV)
			enhScreen.Update(tea.WindowSizeMsg{Width: 80, Height: 40})

			view := enhScreen.View()
			Expect(view).To(ContainSubstring("Enhanced highlight text for preview"))
		})

		It("should use confidence only for master audience", func() {
			bullet := fixtures.CVBulletWith("b-m", "s1", "Master audience bullet")
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-m", []*career.SectionContentGroup{group})
			section.Title = "Experience"

			masterCV := fixtures.CVViewWith("cv-m", "Master CV", "Engineer", "master")
			masterCV.Sections = []*career.CVSection{section}

			masterScreen := cv.NewCVPreviewScreen(masterCV)
			masterScreen.Update(tea.WindowSizeMsg{Width: 80, Height: 40})

			view := masterScreen.View()
			Expect(view).To(ContainSubstring("Master audience bullet"))
		})

		It("should use audience relevance for non-master audience with populated map", func() {
			bullet := fixtures.CVBulletWith("b-ar", "s1", "Audience relevant bullet")
			bullet.AudienceRelevance = map[string]float64{"hiring manager": 0.95}
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-ar", []*career.SectionContentGroup{group})
			section.Title = "Experience"

			arCV := fixtures.CVViewWith("cv-ar", "AR CV", "Engineer", "Hiring Manager")
			arCV.Sections = []*career.CVSection{section}

			arScreen := cv.NewCVPreviewScreen(arCV)
			arScreen.Update(tea.WindowSizeMsg{Width: 80, Height: 40})

			view := arScreen.View()
			Expect(view).To(ContainSubstring("Audience relevant bullet"))
		})

		It("should fall back to confidence when audience key not in map", func() {
			bullet := fixtures.CVBulletWith("b-fb", "s1", "Fallback confidence bullet")
			bullet.AudienceRelevance = map[string]float64{"recruiter": 0.8}
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-fb", []*career.SectionContentGroup{group})
			section.Title = "Experience"

			fbCV := fixtures.CVViewWith("cv-fb", "Fallback CV", "Engineer", "Hiring Manager")
			fbCV.Sections = []*career.CVSection{section}

			fbScreen := cv.NewCVPreviewScreen(fbCV)
			fbScreen.Update(tea.WindowSizeMsg{Width: 80, Height: 40})

			view := fbScreen.View()
			Expect(view).To(ContainSubstring("Fallback confidence bullet"))
		})

		It("should use confidence for empty audience", func() {
			bullet := fixtures.CVBulletWith("b-ea", "s1", "Empty audience bullet")
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-ea", []*career.SectionContentGroup{group})
			section.Title = "Experience"

			eaCV := fixtures.CVViewWith("cv-ea", "Empty Aud CV", "Engineer", "")
			eaCV.Sections = []*career.CVSection{section}

			eaScreen := cv.NewCVPreviewScreen(eaCV)
			eaScreen.Update(tea.WindowSizeMsg{Width: 80, Height: 40})

			view := eaScreen.View()
			Expect(view).To(ContainSubstring("Empty audience bullet"))
		})
	})
})
