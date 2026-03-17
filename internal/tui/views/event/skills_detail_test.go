package event

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SkillsDetail", func() {
	var (
		modal  *SkillsDetail
		theme  themes.Theme
		skills []display.Skill
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
		skills = display.SkillsFromDomain([]*career.Skill{
			fixtures.SkillWith("skill-1", "Go", "backend", "advanced"),
			fixtures.SkillWith("skill-2", "Docker", "devops", "intermediate"),
			fixtures.SkillWith("skill-3", "React", "frontend", "advanced"),
		})
		modal = NewSkillsDetail("event-1", skills, theme)
	})

	Describe("NewSkillsDetail", func() {
		It("should create a new modal with skills", func() {
			Expect(modal).ToNot(BeNil())
			Expect(modal.eventID).To(Equal("event-1"))
			Expect(modal.skills).To(HaveLen(3))
			Expect(modal.visible).To(BeFalse())
		})

		It("should use default theme if nil", func() {
			modal := NewSkillsDetail("event-1", skills, nil)
			Expect(modal.theme).ToNot(BeNil())
		})

		It("should filter out nil skills", func() {
			skillsWithNil := display.SkillsFromDomain([]*career.Skill{
				fixtures.SkillWith("skill-1", "Go", "backend", "advanced"),
				nil,
				fixtures.SkillWith("skill-2", "Docker", "devops", "intermediate"),
			})
			modal := NewSkillsDetail("event-1", skillsWithNil, theme)
			Expect(modal.skills).To(HaveLen(2))
		})

		It("should keep partially populated skills for placeholder rendering", func() {
			placeholderSkill := display.Skill{ID: "skill-4", CreatedAt: skills[0].CreatedAt}
			modal := NewSkillsDetail("event-1", []display.Skill{placeholderSkill}, theme)

			Expect(modal.skills).To(HaveLen(1))
			modal.Show()
			Expect(modal.View()).To(ContainSubstring("(Unnamed)"))
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
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			modal.SetDimensions(100, 30)
			modal.Show()
		})

		It("should navigate down with 'j'", func() {
			Expect(modal.GetSelectedIndex()).To(Equal(0))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

			Expect(modal.GetSelectedIndex()).To(Equal(1))
		})

		It("should navigate up with 'k'", func() {
			modal.table.SetSelectedIndex(1)

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

			Expect(modal.GetSelectedIndex()).To(Equal(0))
		})

		It("should navigate with arrow keys", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(modal.GetSelectedIndex()).To(Equal(1))

			modal.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(modal.GetSelectedIndex()).To(Equal(0))
		})
	})

	Describe("Action Keys", func() {
		BeforeEach(func() {
			modal.Show()
		})

		It("should close on Escape", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should close on q", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should detect 'a' key for add existing", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.GetAction()).To(Equal(ActionAddExisting))
		})

		It("should detect 'n' key for add new", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(modal.GetAction()).To(Equal(ActionAddNew))
		})

		It("should detect 'i' key for infer", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
			Expect(modal.GetAction()).To(Equal(ActionInfer))
		})

		It("should detect 'd' key for remove", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(modal.GetAction()).To(Equal(ActionRemove))
			Expect(modal.GetSelectedSkill()).To(Equal(skills[0]))
		})
	})

	Describe("SetSkills", func() {
		It("should update the skills and reset table", func() {
			newSkills := display.SkillsFromDomain([]*career.Skill{
				fixtures.SkillWith("skill-4", "Python", "backend", "expert"),
			})
			modal.SetSkills(newSkills)

			Expect(modal.skills).To(HaveLen(1))
			Expect(modal.skills[0].Name).To(Equal("Python"))
		})

		It("should filter out nil skills", func() {
			newSkills := display.SkillsFromDomain([]*career.Skill{
				fixtures.SkillWith("skill-4", "Python", "backend", "expert"),
				nil,
			})
			modal.SetSkills(newSkills)

			Expect(modal.skills).To(HaveLen(1))
		})
	})

	Describe("GetSelectedSkill", func() {
		It("should return nil when no skills", func() {
			emptyModal := NewSkillsDetail("event-1", []display.Skill{}, theme)
			Expect(emptyModal.GetSelectedSkill()).To(Equal(display.Skill{}))
		})

		It("should return currently selected skill", func() {
			modal.Show()
			modal.table.SetSelectedIndex(1)

			selected := modal.GetSelectedSkill()
			Expect(selected).ToNot(Equal(display.Skill{}))
			Expect(selected.Name).To(Equal("Docker"))
		})
	})

	Describe("ClearAction", func() {
		It("should clear the action state", func() {
			modal.Show()
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(modal.GetAction()).To(Equal(ActionAddExisting))

			modal.ClearAction()
			Expect(modal.GetAction()).To(Equal(ActionNone))
		})

		It("should clear the selected skill when clearing the action", func() {
			modal.Show()
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

			Expect(modal.GetSelectedSkill()).To(Equal(skills[0]))

			modal.ClearAction()
			Expect(modal.GetSelectedSkill()).To(Equal(skills[0]))
			Expect(modal.GetAction()).To(Equal(ActionNone))
		})
	})
})
