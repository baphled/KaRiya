package skills_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/skills"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SkillDetailScreen", func() {
	var (
		screen *skills.SkillDetailScreen
		skill  *career.Skill
	)

	BeforeEach(func() {
		yearsUsed := 5
		lastUsed := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
		skill = &career.Skill{
			ID:        "skill-1",
			Name:      "Kubernetes",
			Category:  "devops",
			Level:     "advanced",
			YearsUsed: &yearsUsed,
			LastUsed:  &lastUsed,
			CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
		}
	})

	Describe("Construction", func() {
		It("should create a skill detail screen", func() {
			screen = skills.NewSkillDetailScreen(skill)
			Expect(screen).NotTo(BeNil())
			Expect(screen.GetSkill()).To(Equal(skill))
		})

		It("should use default dimensions if not set", func() {
			screen = skills.NewSkillDetailScreen(skill)
			Expect(screen.Width()).To(Equal(120))
			Expect(screen.Height()).To(Equal(40))
		})
	})

	Describe("Actions", func() {
		BeforeEach(func() {
			screen = skills.NewSkillDetailScreen(skill)
		})

		It("should return CancelResult on escape key", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})

		It("should return NavigateResult on 'e' key with 'edit' action", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			Expect(result.Data()).To(Equal("edit"))
		})

		It("should return NavigateResult on 'd' key with 'delete' action", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			Expect(result.Data()).To(Equal("delete"))
		})

		It("should return NavigateResult on enter key with 'back' action", func() {
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			Expect(result.Data()).To(Equal("back"))
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			screen = skills.NewSkillDetailScreen(skill)
			screen.SetTerminalInfo(120, 40)
		})

		It("should render a view with StandardView structure", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should include skill name in view", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Kubernetes"))
		})

		It("should include skill category in view", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("devops"))
		})

		It("should include skill level in view", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("advanced"))
		})

		It("should include years used in view", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("5 years"))
		})

		It("should include last used date in view", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("2025-12-31"))
		})

		It("should show help text in footer", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("e: Edit"))
			Expect(view).To(ContainSubstring("d: Delete"))
			Expect(view).To(ContainSubstring("Esc"))
		})

		It("should handle skill without optional fields", func() {
			minimalSkill := &career.Skill{
				ID:       "skill-2",
				Name:     "Python",
				Category: "backend",
			}
			screen = skills.NewSkillDetailScreen(minimalSkill)
			view := screen.View()

			Expect(view).To(ContainSubstring("Python"))
			Expect(view).To(ContainSubstring("backend"))
		})
	})

	Describe("Terminal Handling", func() {
		BeforeEach(func() {
			screen = skills.NewSkillDetailScreen(skill)
		})

		It("should update dimensions on SetTerminalInfo", func() {
			screen.SetTerminalInfo(100, 30)
			Expect(screen.Width()).To(Equal(100))
			Expect(screen.Height()).To(Equal(30))
		})

		It("should handle WindowSizeMsg", func() {
			msg := tea.WindowSizeMsg{Width: 80, Height: 24}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
			Expect(screen.Width()).To(Equal(80))
			Expect(screen.Height()).To(Equal(24))
		})
	})

	Describe("Screen Interface", func() {
		BeforeEach(func() {
			screen = skills.NewSkillDetailScreen(skill)
		})

		It("should implement Screen interface", func() {
			var _ screens.Screen = screen
		})

		It("should support SetTheme", func() {
			theme := "test-theme"
			screen.SetTheme(theme)
			Expect(screen.Theme()).To(Equal(theme))
		})
	})

	Describe("Edge Cases", func() {
		It("should handle very small terminal dimensions", func() {
			screen := skills.NewSkillDetailScreen(skill)
			screen.SetTerminalInfo(20, 10)

			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle long skill names", func() {
			longSkill := &career.Skill{
				ID:       "skill-3",
				Name:     "Very Long Skill Name That Might Wrap Or Be Truncated",
				Category: "backend",
			}
			screen := skills.NewSkillDetailScreen(longSkill)
			view := screen.View()

			Expect(view).NotTo(BeEmpty())
		})
	})
})
