package cv_test

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/cv"
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
				Title: "Professional Experience",
				Content: []*career.SectionContentGroup{
					{Header: "TechCorp"},
				},
			},
			{
				Title: "Technical Skills",
				Content: []*career.SectionContentGroup{
					{Header: "Skills"},
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

		Context("edit key", func() {
			It("should return NavigateResult with 'edit' on 'e' key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).To(Equal("edit"))
			})
		})

		Context("scrolling", func() {
			It("should scroll down with 'j' key", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should scroll down with down arrow", func() {
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyDown})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should scroll up with 'k' key", func() {
				// First scroll down
				screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should scroll up with up arrow", func() {
				screen.Update(tea.KeyMsg{Type: tea.KeyDown})

				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyUp})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should not scroll above 0", func() {
				// Try to scroll up from 0
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

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

		It("should display CV name", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Name: Test CV"))
		})

		It("should display target role", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Role: Senior Engineer"))
		})

		It("should display target audience", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Audience: Hiring Manager"))
		})

		It("should display source counts", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Events: 10"))
			Expect(view).To(ContainSubstring("Facts: 25"))
		})

		It("should display section titles", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Professional Experience"))
			Expect(view).To(ContainSubstring("Technical Skills"))
		})

		It("should display help text", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("scroll"))
			Expect(view).To(ContainSubstring("confirm"))
			Expect(view).To(ContainSubstring("edit"))
			Expect(view).To(ContainSubstring("back"))
		})
	})

	Describe("Edge Cases", func() {
		Context("nil CV", func() {
			It("should handle nil CV gracefully", func() {
				nilScreen := cv.NewCVPreviewScreen(nil)

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

				view := emptyScreen.View()
				Expect(view).To(ContainSubstring("No sections generated"))
			})
		})
	})
})
