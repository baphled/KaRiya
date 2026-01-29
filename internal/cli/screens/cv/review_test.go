package cv_test

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/cv"
	"github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("ReviewScreen", func() {
	var (
		screen   *cv.ReviewScreen
		testCV   *career.CVView
		sections []*career.CVSection
	)

	BeforeEach(func() {
		sections = []*career.CVSection{
			{
				Title:       "Professional Summary",
				SectionType: "summary",
				Summary:     "Experienced engineer with 10+ years...",
			},
			{
				Title: "Professional Experience",
				Content: []*career.SectionContentGroup{
					{
						Header:    "Senior Engineer at TechCorp",
						StartDate: "2020-01",
						EndDate:   "Present",
						Bullets: []*career.CVBullet{
							{Text: "Led team of 5 engineers"},
							{Text: "Delivered critical features"},
						},
					},
					{
						Header:    "Engineer at StartupCo",
						StartDate: "2018-01",
						EndDate:   "2019-12",
						Bullets: []*career.CVBullet{
							{Text: "Built core platform"},
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
							{Text: "Go, Python, TypeScript"},
						},
					},
				},
			},
		}

		testCV = &career.CVView{
			Name:             "Software Engineer CV",
			TargetRole:       "Senior Engineer",
			TargetAudience:   "Hiring Manager",
			Sections:         sections,
			SourceEventCount: 15,
			SourceFactCount:  42,
		}

		screen = cv.NewCVReviewScreen(testCV)
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
				emptyCV := &career.CVView{
					Name:           "Empty CV",
					TargetRole:     "Engineer",
					TargetAudience: "Manager",
					Sections:       []*career.CVSection{},
				}
				emptyScreen := cv.NewCVReviewScreen(emptyCV)

				view := emptyScreen.View()
				Expect(view).To(ContainSubstring("0 sections"))
			})
		})

		Context("CV with empty sections", func() {
			It("should handle sections with no content", func() {
				cvWithEmptySections := &career.CVView{
					Name:           "CV",
					TargetRole:     "Engineer",
					TargetAudience: "Manager",
					Sections: []*career.CVSection{
						{
							Title:   "Empty Section",
							Content: []*career.SectionContentGroup{},
						},
					},
				}
				s := cv.NewCVReviewScreen(cvWithEmptySections)

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
})
