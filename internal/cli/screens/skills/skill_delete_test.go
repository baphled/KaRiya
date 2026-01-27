package skills_test

import (
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/skills"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SkillDeleteConfirmScreen", func() {
	var (
		screen *skills.SkillDeleteConfirmScreen
		skill  *career.Skill
	)

	BeforeEach(func() {
		skill = &career.Skill{
			ID:       "skill-1",
			Name:     "Kubernetes",
			Category: "devops",
			Level:    "advanced",
		}
	})

	Describe("Construction", func() {
		It("should create a delete confirmation screen", func() {
			screen = skills.NewSkillDeleteConfirmScreen(skill)
			Expect(screen).NotTo(BeNil())
			Expect(screen.GetSkill()).To(Equal(skill))
		})

		It("should default to No selection", func() {
			screen = skills.NewSkillDeleteConfirmScreen(skill)
			Expect(screen.GetSelection()).To(BeFalse())
		})
	})

	Describe("Selection Toggle", func() {
		BeforeEach(func() {
			screen = skills.NewSkillDeleteConfirmScreen(skill)
		})

		It("should toggle selection with arrow keys", func() {
			Expect(screen.GetSelection()).To(BeFalse())

			msg := tea.KeyMsg{Type: tea.KeyLeft}
			screen.Update(msg)
			Expect(screen.GetSelection()).To(BeTrue())

			msg = tea.KeyMsg{Type: tea.KeyRight}
			screen.Update(msg)
			Expect(screen.GetSelection()).To(BeFalse())
		})

		It("should toggle with h/l keys", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
			screen.Update(msg)
			Expect(screen.GetSelection()).To(BeTrue())
		})
	})

	Describe("Confirmation", func() {
		BeforeEach(func() {
			screen = skills.NewSkillDeleteConfirmScreen(skill)
		})

		It("should return NavigateResult with false on No confirmation", func() {
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			Expect(result.Data()).To(BeFalse())
		})

		It("should return NavigateResult with true on Yes confirmation", func() {
			screen.SetSelection(true)

			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			Expect(result.Data()).To(BeTrue())
		})

		It("should return true on 'y' key press", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Data()).To(BeTrue())
		})

		It("should return false on 'n' key press", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Data()).To(BeFalse())
		})

		It("should return CancelResult on escape", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			screen = skills.NewSkillDeleteConfirmScreen(skill)
			screen.SetTerminalInfo(120, 40)
		})

		It("should render a view with StandardView structure", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should include skill name in confirmation message", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Kubernetes"))
		})

		It("should show delete warning message", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("delete"))
		})

		It("should show Delete/Cancel buttons", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Delete"))
			Expect(view).To(ContainSubstring("Cancel"))
		})

		It("should show help text in footer", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("y/n"))
			Expect(view).To(ContainSubstring("Esc"))
		})
	})

	Describe("Terminal Handling", func() {
		BeforeEach(func() {
			screen = skills.NewSkillDeleteConfirmScreen(skill)
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
			screen = skills.NewSkillDeleteConfirmScreen(skill)
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
})
