package modals_test

import (
	"github.com/baphled/kariya/internal/cli/screens/skills/modals"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DetailModal", func() {
	var (
		modal *modals.DetailModal
		skill *career.Skill
		theme themes.Theme
	)

	BeforeEach(func() {
		years := 5
		skill = fixtures.SkillWith("skill-1", "Go Programming", "Programming Languages", "Expert")
		skill.YearsUsed = &years
		theme = themes.NewDefaultTheme()
	})

	Describe("NewDetailModal", func() {
		It("creates a modal with the skill and theme", func() {
			modal = modals.NewDetailModal(skill, theme, 10, nil)
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("creates a modal with event count", func() {
			modal = modals.NewDetailModal(skill, theme, 25, nil)
			Expect(modal).NotTo(BeNil())
		})

		It("creates a modal with last used date", func() {
			modal = modals.NewDetailModal(skill, theme, 10, nil)
			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Show/Hide", func() {
		BeforeEach(func() {
			modal = modals.NewDetailModal(skill, theme, 10, nil)
		})

		It("shows the modal", func() {
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("hides the modal", func() {
			modal.Show()
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("SetDimensions", func() {
		BeforeEach(func() {
			modal = modals.NewDetailModal(skill, theme, 10, nil)
		})

		It("updates dimensions", func() {
			modal.SetDimensions(100, 50)
			// Dimensions are used internally, just verify no panic
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = modals.NewDetailModal(skill, theme, 10, nil)
			modal.SetDimensions(120, 40)
			modal.Show()
		})

		Context("when visible", func() {
			It("closes on Escape key", func() {
				Expect(modal.IsVisible()).To(BeTrue())
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(modal.IsVisible()).To(BeFalse())
			})

			It("closes on Enter key", func() {
				Expect(modal.IsVisible()).To(BeTrue())
				modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(modal.IsVisible()).To(BeFalse())
			})

			It("closes on 'q' key", func() {
				Expect(modal.IsVisible()).To(BeTrue())
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
				Expect(modal.IsVisible()).To(BeFalse())
			})

			It("closes on backspace key", func() {
				Expect(modal.IsVisible()).To(BeTrue())
				modal.Update(tea.KeyMsg{Type: tea.KeyBackspace})
				Expect(modal.IsVisible()).To(BeFalse())
			})

			It("handles window size messages", func() {
				_, cmd := modal.Update(tea.WindowSizeMsg{Width: 150, Height: 60})
				Expect(cmd).To(BeNil())
				// Modal should still be visible
				Expect(modal.IsVisible()).To(BeTrue())
			})
		})

		Context("when not visible", func() {
			BeforeEach(func() {
				modal.Hide()
			})

			It("ignores key presses", func() {
				Expect(modal.IsVisible()).To(BeFalse())
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				// Still not visible, no change
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})
	})

	Describe("View", func() {
		BeforeEach(func() {
			modal = modals.NewDetailModal(skill, theme, 10, nil)
			modal.SetDimensions(120, 40)
		})

		Context("when visible", func() {
			BeforeEach(func() {
				modal.Show()
			})

			It("renders skill name", func() {
				view := modal.View()
				Expect(view).To(ContainSubstring("Go Programming"))
			})

			It("renders skill category", func() {
				view := modal.View()
				Expect(view).To(ContainSubstring("Programming Languages"))
			})

			It("handles nil theme by using default", func() {
				nilThemeModal := modals.NewDetailModal(skill, nil, 10, nil)
				nilThemeModal.SetDimensions(120, 40)
				nilThemeModal.Show()
				view := nilThemeModal.View()
				Expect(view).NotTo(BeEmpty())
				Expect(view).To(ContainSubstring("Go Programming"))
			})

			It("renders skill level", func() {
				view := modal.View()
				Expect(view).To(ContainSubstring("Expert"))
			})

			It("renders years used", func() {
				view := modal.View()
				Expect(view).To(ContainSubstring("5"))
			})

			It("renders event count when non-zero", func() {
				modal = modals.NewDetailModal(skill, theme, 10, nil)
				modal.SetDimensions(120, 40)
				modal.Show()
				view := modal.View()
				Expect(view).To(ContainSubstring("10"))
			})

			It("renders event count when zero", func() {
				modal = modals.NewDetailModal(skill, theme, 0, nil)
				modal.SetDimensions(120, 40)
				modal.Show()
				view := modal.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("renders with last used date", func() {
				modal = modals.NewDetailModal(skill, theme, 10, nil)
				modal.SetDimensions(120, 40)
				modal.Show()
				view := modal.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("renders action hints", func() {
				view := modal.View()
				// Check for keybind hints
				Expect(view).To(ContainSubstring("e"))
				Expect(view).To(ContainSubstring("d"))
			})

			It("renders with different widths", func() {
				modal.SetDimensions(80, 40)
				view := modal.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("renders with very narrow width", func() {
				modal.SetDimensions(40, 40)
				view := modal.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("renders footer with close hint", func() {
				view := modal.View()
				Expect(view).To(ContainSubstring("Enter/Esc"))
			})
		})

		Context("when not visible", func() {
			It("returns empty string", func() {
				view := modal.View()
				Expect(view).To(BeEmpty())
			})
		})

		Context("with nil years used", func() {
			BeforeEach(func() {
				skillWithoutYears := fixtures.SkillWith("skill-2", "Python", "Languages", "Intermediate")
				modal = modals.NewDetailModal(skillWithoutYears, theme, 5, nil)
				modal.SetDimensions(120, 40)
				modal.Show()
			})

			It("renders without years", func() {
				view := modal.View()
				Expect(view).To(ContainSubstring("Python"))
				// Should not panic and should still render
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	Describe("GetAction", func() {
		BeforeEach(func() {
			modal = modals.NewDetailModal(skill, theme, 10, nil)
			modal.SetDimensions(120, 40)
		})

		It("returns 'edit' when 'e' is pressed", func() {
			modal.Show()
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(modal.GetAction()).To(Equal("edit"))
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("returns 'events' when ctrl+e is pressed", func() {
			modal.Show()
			modal.Update(tea.KeyMsg{Type: tea.KeyCtrlE})
			Expect(modal.GetAction()).To(Equal("events"))
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("returns 'delete' when 'd' is pressed", func() {
			modal.Show()
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(modal.GetAction()).To(Equal("delete"))
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("returns empty string when closed normally", func() {
			modal.Show()
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.GetAction()).To(Equal(""))
		})
	})

	Describe("SetSkill", func() {
		BeforeEach(func() {
			modal = modals.NewDetailModal(skill, theme, 10, nil)
			modal.SetDimensions(120, 40)
		})

		It("updates the skill being displayed", func() {
			newSkill := fixtures.SkillWith("skill-new", "Rust", "Systems", "Beginner")
			modal.SetSkill(newSkill, 3, nil)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("Rust"))
			Expect(view).To(ContainSubstring("Systems"))
		})

		It("handles skill with nil years used", func() {
			newSkill := fixtures.SkillWith("skill-nil-years", "Python", "Scripting", "Intermediate")
			newSkill.YearsUsed = nil
			modal.SetSkill(newSkill, 0, nil)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("Python"))
		})

		It("handles skill with 1 year (singular)", func() {
			newSkill := fixtures.SkillWith("skill-one-year", "Ruby", "Web", "Advanced")
			oneYear := 1
			newSkill.YearsUsed = &oneYear
			modal.SetSkill(newSkill, 0, nil)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("Ruby"))
			Expect(view).To(ContainSubstring("1 year"))
		})

		It("handles skill with multiple years (plural)", func() {
			newSkill := fixtures.SkillWith("skill-many-years", "JavaScript", "Frontend", "Expert")
			manyYears := 10
			newSkill.YearsUsed = &manyYears
			modal.SetSkill(newSkill, 0, nil)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("JavaScript"))
			Expect(view).To(ContainSubstring("10 years"))
		})

		It("handles skill with empty category", func() {
			newSkill := fixtures.SkillWith("skill-no-cat", "TypeScript", "", "Intermediate")
			modal.SetSkill(newSkill, 0, nil)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("TypeScript"))
		})

		It("handles skill with empty level", func() {
			newSkill := fixtures.SkillWith("skill-no-level", "Java", "Backend", "")
			modal.SetSkill(newSkill, 0, nil)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("Java"))
		})

		It("handles skill with zero event count", func() {
			newSkill := fixtures.SkillWith("skill-no-events", "Ruby", "Web", "Advanced")
			modal.SetSkill(newSkill, 0, nil)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("Ruby"))
		})
	})

	Describe("Init", func() {
		It("returns nil command", func() {
			modal = modals.NewDetailModal(skill, theme, 10, nil)
			cmd := modal.Init()
			Expect(cmd).To(BeNil())
		})
	})
})
