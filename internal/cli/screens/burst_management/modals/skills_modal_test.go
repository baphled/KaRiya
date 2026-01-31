package modals_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("BurstSkillsModal", func() {
	var (
		modal     *modals.BurstSkillsModal
		burstID   string
		burstName string
		skills    []*career.Skill
		theme     themes.Theme
	)

	BeforeEach(func() {
		burstID = "test-burst-id"
		burstName = "Backend Development"
		skills = []*career.Skill{
			{
				ID:       "s1",
				Name:     "Go",
				Category: "Programming Languages",
				Level:    "advanced",
			},
			{
				ID:       "s2",
				Name:     "Docker",
				Category: "DevOps",
				Level:    "intermediate",
			},
			{
				ID:       "s3",
				Name:     "PostgreSQL",
				Category: "Databases",
				Level:    "advanced",
			},
		}
		theme = themes.NewDefaultTheme()
	})

	Describe("NewBurstSkillsModal", func() {
		It("creates modal with skills content", func() {
			modal = modals.NewBurstSkillsModal(burstID, burstName, skills, theme)

			Expect(modal).NotTo(BeNil())
			Expect(modal.GetBurstID()).To(Equal(burstID))
		})

		It("handles nil theme with default", func() {
			modal = modals.NewBurstSkillsModal(burstID, burstName, skills, nil)

			Expect(modal).NotTo(BeNil())
		})

		It("handles empty skills list", func() {
			modal = modals.NewBurstSkillsModal(burstID, burstName, []*career.Skill{}, theme)

			Expect(modal).NotTo(BeNil())
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("No skills"))
		})

		It("handles nil skills list", func() {
			modal = modals.NewBurstSkillsModal(burstID, burstName, nil, theme)

			Expect(modal).NotTo(BeNil())
		})

		It("includes skill count in title", func() {
			modal = modals.NewBurstSkillsModal(burstID, burstName, skills, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("3"))
		})

		It("includes burst name in title", func() {
			modal = modals.NewBurstSkillsModal(burstID, burstName, skills, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("Backend Development"))
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			modal = modals.NewBurstSkillsModal(burstID, burstName, skills, theme)
		})

		It("returns initialization command", func() {
			cmd := modal.Init()
			_ = cmd
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = modals.NewBurstSkillsModal(burstID, burstName, skills, theme)
			modal.SetDimensions(80, 24)
			modal.Show()
		})

		It("delegates to underlying modal", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}

			model, cmd := modal.Update(msg)

			Expect(model).NotTo(BeNil())
			_ = cmd
		})

		It("handles window size message", func() {
			msg := tea.WindowSizeMsg{Width: 100, Height: 30}

			model, cmd := modal.Update(msg)

			Expect(model).NotTo(BeNil())
			_ = cmd
		})

		It("handles escape key", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}

			model, cmd := modal.Update(msg)

			Expect(model).NotTo(BeNil())
			_ = cmd
		})
	})

	Describe("View", func() {
		BeforeEach(func() {
			modal = modals.NewBurstSkillsModal(burstID, burstName, skills, theme)
			modal.SetDimensions(80, 24)
		})

		It("renders skills when visible", func() {
			modal.Show()

			view := modal.View()

			Expect(view).NotTo(BeEmpty())
		})

		It("shows skill name", func() {
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("Go"))
		})

		It("shows skill category", func() {
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("Programming Languages"))
		})

		It("shows skill level", func() {
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("advanced"))
		})

		It("numbers skills", func() {
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("1."))
			Expect(view).To(ContainSubstring("2."))
		})

		It("shows empty state message when no skills", func() {
			modal = modals.NewBurstSkillsModal(burstID, burstName, []*career.Skill{}, theme)
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("No skills"))
		})
	})

	Describe("Visibility methods", func() {
		BeforeEach(func() {
			modal = modals.NewBurstSkillsModal(burstID, burstName, skills, theme)
		})

		It("Show makes modal visible", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())

			modal.Show()

			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("Hide makes modal invisible", func() {
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())

			modal.Hide()

			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("SetDimensions", func() {
		It("sets terminal dimensions", func() {
			modal = modals.NewBurstSkillsModal(burstID, burstName, skills, theme)

			modal.SetDimensions(100, 50)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("SetSkills", func() {
		It("updates the displayed skills", func() {
			modal = modals.NewBurstSkillsModal(burstID, burstName, skills, theme)
			newSkills := []*career.Skill{
				{ID: "new-1", Name: "Rust", Category: "Programming Languages", Level: "beginner"},
				{ID: "new-2", Name: "Kubernetes", Category: "DevOps", Level: "intermediate"},
			}

			modal.SetSkills(newSkills)
			modal.Show()
			view := modal.View()

			Expect(view).To(ContainSubstring("Rust"))
		})

		It("updates skill count in view", func() {
			modal = modals.NewBurstSkillsModal(burstID, burstName, skills, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("3"))

			modal.SetSkills([]*career.Skill{
				{ID: "single", Name: "Python", Category: "Programming Languages"},
			})
			view = modal.View()
			Expect(view).To(ContainSubstring("1"))
		})
	})

	Describe("GetBurstID", func() {
		It("returns the burst ID", func() {
			modal = modals.NewBurstSkillsModal(burstID, burstName, skills, theme)

			result := modal.GetBurstID()

			Expect(result).To(Equal(burstID))
		})
	})

	Describe("Skills with optional fields", func() {
		It("renders skill without category", func() {
			skillsWithoutCategory := []*career.Skill{
				{
					ID:    "s-no-cat",
					Name:  "Skill without category",
					Level: "intermediate",
				},
			}
			modal = modals.NewBurstSkillsModal(burstID, burstName, skillsWithoutCategory, theme)
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("Skill without category"))
		})

		It("renders skill without level", func() {
			skillsWithoutLevel := []*career.Skill{
				{
					ID:       "s-no-lvl",
					Name:     "Skill without level",
					Category: "General",
				},
			}
			modal = modals.NewBurstSkillsModal(burstID, burstName, skillsWithoutLevel, theme)
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("Skill without level"))
		})
	})
})
