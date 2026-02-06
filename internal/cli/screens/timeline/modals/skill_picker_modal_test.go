package modals

import (
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SkillPickerModal", func() {
	var (
		modal  *SkillPickerModal
		theme  themes.Theme
		skills []*career.Skill
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
		skills = []*career.Skill{
			fixtures.SkillWith("skill-1", "Go", "backend", "advanced"),
			fixtures.SkillWith("skill-2", "Docker", "devops", "intermediate"),
			fixtures.SkillWith("skill-3", "React", "frontend", "advanced"),
		}
		modal = NewSkillPickerModal(skills, theme)
	})

	Describe("NewSkillPickerModal", func() {
		It("should create a new modal with skills", func() {
			Expect(modal).ToNot(BeNil())
			Expect(modal.skills).To(HaveLen(3))
			Expect(modal.visible).To(BeFalse())
		})

		It("should use default theme if nil", func() {
			modal := NewSkillPickerModal(skills, nil)
			Expect(modal.theme).ToNot(BeNil())
		})

		It("should filter out nil skills", func() {
			skillsWithNil := []*career.Skill{
				fixtures.SkillWith("skill-1", "Go", "backend", "advanced"),
				nil,
				fixtures.SkillWith("skill-2", "Docker", "devops", "intermediate"),
			}
			modal := NewSkillPickerModal(skillsWithNil, theme)
			Expect(modal.skills).To(HaveLen(2))
		})

		It("should initialize table with correct columns", func() {
			Expect(modal.table).ToNot(BeNil())
		})
	})

	Describe("Visibility", func() {
		It("should start hidden", func() {
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should show when Show() is called", func() {
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should hide when Hide() is called", func() {
			modal.Show()
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should return empty view when hidden", func() {
			modal.Hide()
			view := modal.View()
			Expect(view).To(BeEmpty())
		})

		It("should reset selection when shown", func() {
			modal.selectedSkill = skills[0]
			modal.Show()
			Expect(modal.selectedSkill).To(BeNil())
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			modal.SetDimensions(100, 30)
			modal.Show()
		})

		It("should navigate down with 'j'", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(modal.table.GetSelectedIndex()).To(Equal(1))
		})

		It("should navigate up with 'k'", func() {
			modal.table.SetSelectedIndex(1)
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(modal.table.GetSelectedIndex()).To(Equal(0))
		})

		It("should navigate with arrow keys", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(modal.table.GetSelectedIndex()).To(Equal(1))

			modal.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(modal.table.GetSelectedIndex()).To(Equal(0))
		})
	})

	Describe("Selection", func() {
		BeforeEach(func() {
			modal.Show()
		})

		It("should select skill on Enter key", func() {
			Expect(modal.HasSelection()).To(BeFalse())

			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(modal.HasSelection()).To(BeTrue())
			Expect(modal.GetSelectedSkill()).To(Equal(skills[0]))
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should close without selection on Escape", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(modal.HasSelection()).To(BeFalse())
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should close without selection on 'q'", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

			Expect(modal.HasSelection()).To(BeFalse())
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should clear selection", func() {
			modal.selectedSkill = skills[0]
			modal.ClearSelection()
			Expect(modal.selectedSkill).To(BeNil())
		})
	})

	Describe("SetSkills", func() {
		It("should update the skills and reset table", func() {
			newSkills := []*career.Skill{
				fixtures.SkillWith("skill-4", "Python", "backend", "expert"),
			}
			modal.SetSkills(newSkills)

			Expect(modal.skills).To(HaveLen(1))
			Expect(modal.skills[0].Name).To(Equal("Python"))
		})

		It("should filter out nil skills", func() {
			newSkills := []*career.Skill{
				fixtures.SkillWith("skill-4", "Python", "backend", "expert"),
				nil,
			}
			modal.SetSkills(newSkills)

			Expect(modal.skills).To(HaveLen(1))
		})

		It("should clear selection when setting new skills", func() {
			modal.selectedSkill = skills[0]
			modal.SetSkills(skills)
			Expect(modal.selectedSkill).To(BeNil())
		})
	})

	Describe("View Rendering", func() {
		It("should render with skills", func() {
			modal.SetDimensions(100, 30)
			modal.Show()

			view := modal.View()
			Expect(view).ToNot(BeEmpty())
			Expect(view).To(ContainSubstring("Select Skill"))
		})

		It("should show empty message when no skills", func() {
			emptyModal := NewSkillPickerModal([]*career.Skill{}, theme)
			emptyModal.SetDimensions(100, 30)
			emptyModal.Show()

			view := emptyModal.View()
			Expect(view).To(ContainSubstring("No skills available"))
		})
	})

	Describe("Window Resize", func() {
		BeforeEach(func() {
			modal.Show()
		})

		It("should handle window resize", func() {
			modal.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

			Expect(modal.width).To(Equal(120))
			Expect(modal.height).To(Equal(40))
		})
	})
})
