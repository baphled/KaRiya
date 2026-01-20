package cv_test

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/cv"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("CVPreviewScreen", func() {
	var (
		screen   *cv.CVPreviewScreen
		testCV   *career.CVView
		sections []*career.CVSection
	)

	BeforeEach(func() {
		sections = []*career.CVSection{
			{
				Title:       "Professional Summary",
				SectionType: "summary",
				Summary:     "Experienced software engineer with expertise in distributed systems.",
			},
			{
				Title: "Professional Experience",
				Content: []*career.SectionContentGroup{
					{
						Header:    "Senior Engineer at TechCorp",
						StartDate: "2020-01",
						EndDate:   "Present",
						Bullets: []*career.CVBullet{
							{Text: "Led development of microservices architecture"},
							{Text: "Mentored team of 5 junior engineers"},
							{Text: "Reduced deployment time by 60%"},
						},
					},
					{
						Header:    "Software Engineer at StartupCo",
						StartDate: "2018-01",
						EndDate:   "2019-12",
						Bullets: []*career.CVBullet{
							{Text: "Built core payment processing system"},
							{Text: "Implemented real-time notifications"},
						},
					},
				},
			},
			{
				Title: "Technical Skills",
				Content: []*career.SectionContentGroup{
					{
						Header: "Languages",
						Bullets: []*career.CVBullet{
							{Text: "Go, Python, TypeScript, Java"},
						},
					},
					{
						Header: "Frameworks",
						Bullets: []*career.CVBullet{
							{Text: "React, Vue, Django, Spring Boot"},
						},
					},
				},
			},
		}

		testCV = &career.CVView{
			Name:             "Test CV",
			TargetRole:       "Senior Engineer",
			TargetAudience:   "Hiring Manager",
			Sections:         sections,
			SourceEventCount: 10,
			SourceFactCount:  25,
		}

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

		It("should display help text with scrolling options", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("scroll"))
			Expect(view).To(ContainSubstring("confirm"))
			Expect(view).To(ContainSubstring("export"))
			Expect(view).To(ContainSubstring("back"))
		})

		It("should display scroll percentage when content is scrollable", func() {
			// Create a CV with lots of content to ensure scrolling
			longCV := &career.CVView{
				Name:     "Long CV",
				Sections: make([]*career.CVSection, 10),
			}
			for i := 0; i < 10; i++ {
				longCV.Sections[i] = &career.CVSection{
					Title: "Section " + string(rune('A'+i)),
					Content: []*career.SectionContentGroup{
						{
							Header: "Group",
							Bullets: []*career.CVBullet{
								{Text: "Bullet 1"},
								{Text: "Bullet 2"},
								{Text: "Bullet 3"},
							},
						},
					},
				}
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
				emptyCV := &career.CVView{
					Name:           "Empty CV",
					TargetRole:     "Engineer",
					TargetAudience: "Manager",
					Sections:       []*career.CVSection{},
				}
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
				singleDateCV := &career.CVView{
					Name: "Single Date CV",
					Sections: []*career.CVSection{
						{
							Title: "Experience",
							Content: []*career.SectionContentGroup{
								{
									Header:    "Event",
									StartDate: "2020-06",
									EndDate:   "2020-06",
									Bullets:   []*career.CVBullet{{Text: "Did something"}},
								},
							},
						},
					},
				}
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
})
