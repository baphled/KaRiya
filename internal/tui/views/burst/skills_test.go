package burst_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/testutil/fixtures"
	burstview "github.com/baphled/kariya/internal/tui/views/burst"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Skills", func() {
	var (
		modal     *burstview.Skills
		burstID   string
		burstName string
		skills    []display.Skill
		theme     themes.Theme
	)

	BeforeEach(func() {
		burstID = "test-burst-id"
		burstName = "Backend Development"
		skills = []display.Skill{
			display.SkillFromDomain(fixtures.SkillWith("s1", "Go", "Programming Languages", "advanced")),
			display.SkillFromDomain(fixtures.SkillWith("s2", "Docker", "DevOps", "intermediate")),
			display.SkillFromDomain(fixtures.SkillWith("s3", "PostgreSQL", "Databases", "advanced")),
		}
		theme = themes.NewDefaultTheme()
	})

	Describe("NewSkills", func() {
		It("creates modal with skills content", func() {
			modal = burstview.NewSkills(burstID, burstName, skills, theme)

			Expect(modal).NotTo(BeNil())
			Expect(modal.GetBurstID()).To(Equal(burstID))
		})

		It("handles nil theme with default", func() {
			modal = burstview.NewSkills(burstID, burstName, skills, nil)

			Expect(modal).NotTo(BeNil())
		})

		It("handles empty skills list", func() {
			modal = burstview.NewSkills(burstID, burstName, []display.Skill{}, theme)

			Expect(modal).NotTo(BeNil())
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("No skills"))
		})

		It("handles nil skills list", func() {
			modal = burstview.NewSkills(burstID, burstName, nil, theme)

			Expect(modal).NotTo(BeNil())
		})

		It("includes skill count in title", func() {
			modal = burstview.NewSkills(burstID, burstName, skills, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("3"))
		})

		It("includes burst name in title", func() {
			modal = burstview.NewSkills(burstID, burstName, skills, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("Backend Development"))
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			modal = burstview.NewSkills(burstID, burstName, skills, theme)
		})

		It("returns initialization command", func() {
			cmd := modal.Init()
			_ = cmd
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = burstview.NewSkills(burstID, burstName, skills, theme)
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
			modal = burstview.NewSkills(burstID, burstName, skills, theme)
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
			modal = burstview.NewSkills(burstID, burstName, []display.Skill{}, theme)
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("No skills"))
		})
	})

	Describe("Visibility methods", func() {
		BeforeEach(func() {
			modal = burstview.NewSkills(burstID, burstName, skills, theme)
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
			modal = burstview.NewSkills(burstID, burstName, skills, theme)

			modal.SetDimensions(100, 50)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("SetSkills", func() {
		It("updates the displayed skills", func() {
			modal = burstview.NewSkills(burstID, burstName, skills, theme)
			newSkills := []display.Skill{
				display.SkillFromDomain(fixtures.SkillWith("new-1", "Rust", "Programming Languages", "beginner")),
				display.SkillFromDomain(fixtures.SkillWith("new-2", "Kubernetes", "DevOps", "intermediate")),
			}

			modal.SetSkills(newSkills)
			modal.Show()
			view := modal.View()

			Expect(view).To(ContainSubstring("Rust"))
		})

		It("updates skill count in view", func() {
			modal = burstview.NewSkills(burstID, burstName, skills, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("3"))

			modal.SetSkills([]display.Skill{
				display.SkillFromDomain(fixtures.SkillWith("single", "Python", "Programming Languages", "")),
			})
			view = modal.View()
			Expect(view).To(ContainSubstring("1"))
		})
	})

	Describe("GetBurstID", func() {
		It("returns the burst ID", func() {
			modal = burstview.NewSkills(burstID, burstName, skills, theme)

			result := modal.GetBurstID()

			Expect(result).To(Equal(burstID))
		})
	})

	Describe("Skills with optional fields", func() {
		It("renders skill without category", func() {
			skillsWithoutCategory := []display.Skill{
				display.SkillFromDomain(fixtures.SkillWith("s-no-cat", "Skill without category", "", "intermediate")),
			}
			modal = burstview.NewSkills(burstID, burstName, skillsWithoutCategory, theme)
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("Skill without category"))
		})

		It("renders skill without level", func() {
			skillsWithoutLevel := []display.Skill{
				display.SkillFromDomain(fixtures.SkillWith("s-no-lvl", "Skill without level", "General", "")),
			}
			modal = burstview.NewSkills(burstID, burstName, skillsWithoutLevel, theme)
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("Skill without level"))
		})
	})

	Describe("Esc Visibility", func() {
		BeforeEach(func() {
			modal = burstview.NewSkills(burstID, burstName, skills, theme)
			modal.SetDimensions(80, 24)
		})

		It("IsVisible should be false after Esc", func() {
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())

			escMsg := tea.KeyMsg{Type: tea.KeyEsc}
			modal.Update(escMsg)

			Expect(modal.IsVisible()).To(BeFalse())
		})
	})
})
