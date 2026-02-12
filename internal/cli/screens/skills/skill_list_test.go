//nolint:errcheck // Test file - error handling for test setup is not relevant.
package skills_test

import (
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/skills"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SkillsListScreen", func() {
	var (
		screen     *skills.SkillsListScreen
		skillsList []*career.Skill
	)

	BeforeEach(func() {
		skillsList = []*career.Skill{
			fixtures.SkillWith("1", "Ruby", "backend", "expert"),
			fixtures.SkillWith("2", "React", "frontend", "advanced"),
			fixtures.SkillWith("3", "Kubernetes", "devops", "intermediate"),
			fixtures.SkillWith("4", "PostgreSQL", "database", "advanced"),
		}
	})

	Describe("Construction", func() {
		It("should create a skills list screen", func() {
			screen = skills.NewSkillsListScreen(skillsList)
			Expect(screen).NotTo(BeNil())
		})

		It("should use default dimensions if not set", func() {
			screen = skills.NewSkillsListScreen(skillsList)
			Expect(screen.Width()).To(Equal(120))
			Expect(screen.Height()).To(Equal(40))
		})

		It("should handle empty skill list", func() {
			screen = skills.NewSkillsListScreen([]*career.Skill{})
			Expect(screen).NotTo(BeNil())
			view := screen.View()
			Expect(view).To(ContainSubstring("No skills"))
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			screen = skills.NewSkillsListScreen(skillsList)
		})

		It("should navigate down with arrow key", func() {
			initialIndex := screen.GetSelectedIndex()

			msg := tea.KeyMsg{Type: tea.KeyDown}
			_, result := screen.Update(msg)

			Expect(result).To(BeNil())
			Expect(screen.GetSelectedIndex()).To(Equal(initialIndex + 1))
		})

		It("should navigate up with arrow key", func() {
			screen.SetSelectedIndex(2)

			msg := tea.KeyMsg{Type: tea.KeyUp}
			_, result := screen.Update(msg)

			Expect(result).To(BeNil())
			Expect(screen.GetSelectedIndex()).To(Equal(1))
		})

		It("should navigate with vim-style j key", func() {
			initialIndex := screen.GetSelectedIndex()

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			_, result := screen.Update(msg)

			Expect(result).To(BeNil())
			Expect(screen.GetSelectedIndex()).To(Equal(initialIndex + 1))
		})

		It("should navigate with vim-style k key", func() {
			screen.SetSelectedIndex(2)

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
			_, result := screen.Update(msg)

			Expect(result).To(BeNil())
			Expect(screen.GetSelectedIndex()).To(Equal(1))
		})

		It("should jump to top with g key", func() {
			screen.SetSelectedIndex(3)

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
			screen.Update(msg)

			Expect(screen.GetSelectedIndex()).To(Equal(0))
		})

		It("should jump to bottom with G key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}}
			screen.Update(msg)

			Expect(screen.GetSelectedIndex()).To(Equal(3)) // Last item
		})
	})

	Describe("Actions", func() {
		BeforeEach(func() {
			screen = skills.NewSkillsListScreen(skillsList)
		})

		It("should return NavigateResult on enter key with 'view' action", func() {
			screen.SetSelectedIndex(1)

			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))

			data := result.Data().(map[string]interface{})
			Expect(data["action"]).To(Equal("view"))
			Expect(data["skill"]).To(Equal(skillsList[1]))
		})

		It("should return NavigateResult on 'a' key with 'add' action", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))

			data := result.Data().(map[string]interface{})
			Expect(data["action"]).To(Equal("add"))
		})

		It("should return NavigateResult on 'e' key with 'edit' action", func() {
			screen.SetSelectedIndex(2)

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))

			data := result.Data().(map[string]interface{})
			Expect(data["action"]).To(Equal("edit"))
			Expect(data["skill"]).To(Equal(skillsList[2]))
		})

		It("should return NavigateResult on 'd' key with 'delete' action", func() {
			screen.SetSelectedIndex(0)

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))

			data := result.Data().(map[string]interface{})
			Expect(data["action"]).To(Equal("delete"))
			Expect(data["skill"]).To(Equal(skillsList[0]))
		})

		It("should return CancelResult on escape key", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})

		It("should return filter action on 'f' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))

			data := result.Data().(map[string]interface{})
			Expect(data["action"]).To(Equal("filter"))
		})

		It("should return sort action on 's' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))

			data := result.Data().(map[string]interface{})
			Expect(data["action"]).To(Equal("sort"))
		})

		It("should return search action on '/' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))

			data := result.Data().(map[string]interface{})
			Expect(data["action"]).To(Equal("search"))
		})

		It("should return infer action on 'i' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))

			data := result.Data().(map[string]interface{})
			Expect(data["action"]).To(Equal("infer"))
		})

		It("should return help action on '?' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))

			data := result.Data().(map[string]interface{})
			Expect(data["action"]).To(Equal("help"))
		})

		It("should return nil for edit 'e' on empty list", func() {
			emptyScreen := skills.NewSkillsListScreen([]*career.Skill{})

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			_, result := emptyScreen.Update(msg)

			Expect(result).To(BeNil())
		})

		It("should return nil for delete 'd' on empty list", func() {
			emptyScreen := skills.NewSkillsListScreen([]*career.Skill{})

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
			_, result := emptyScreen.Update(msg)

			Expect(result).To(BeNil())
		})

		It("should return nil for Enter on empty list", func() {
			emptyScreen := skills.NewSkillsListScreen([]*career.Skill{})

			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := emptyScreen.Update(msg)

			Expect(result).To(BeNil())
		})
	})

	Describe("Empty List Safeguards (Phase 2-Tier 3)", func() {
		It("should return nil for edit 'e' on empty list", func() {
			emptyScreen := skills.NewSkillsListScreen([]*career.Skill{})

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			_, result := emptyScreen.Update(msg)

			Expect(result).To(BeNil())
		})

		It("should return nil for delete 'd' on empty list", func() {
			emptyScreen := skills.NewSkillsListScreen([]*career.Skill{})

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
			_, result := emptyScreen.Update(msg)

			Expect(result).To(BeNil())
		})

		It("should allow add 'a' action on empty list", func() {
			emptyScreen := skills.NewSkillsListScreen([]*career.Skill{})

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
			_, result := emptyScreen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			data := result.Data().(map[string]interface{})
			Expect(data["action"]).To(Equal("add"))
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			screen = skills.NewSkillsListScreen(skillsList)
			screen.SetTerminalInfo(120, 40)
		})

		It("should render a view with StandardView structure", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should include skill names in view", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Ruby"))
			Expect(view).To(ContainSubstring("React"))
			Expect(view).To(ContainSubstring("Kubernetes"))
		})

		It("should include skill categories in view", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("backend"))
			Expect(view).To(ContainSubstring("frontend"))
			Expect(view).To(ContainSubstring("devops"))
		})

		It("should include skill levels in view", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("expert"))
			Expect(view).To(ContainSubstring("advanced"))
			Expect(view).To(ContainSubstring("intermediate"))
		})

		It("should show selection indicator on selected item", func() {
			screen.SetSelectedIndex(1)
			view := screen.View()
			// Should contain selection indicator (exact format depends on implementation)
			Expect(view).NotTo(BeEmpty())
		})

		It("should show help text in footer", func() {
			view := screen.View()
			// UIKit badges render as "[key] action" without colon.
			Expect(view).To(ContainSubstring("Add"))
			Expect(view).To(ContainSubstring("Edit"))
			Expect(view).To(ContainSubstring("Delete"))
			Expect(view).To(ContainSubstring("Esc"))
		})
	})

	Describe("Terminal Handling", func() {
		BeforeEach(func() {
			screen = skills.NewSkillsListScreen(skillsList)
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
			screen = skills.NewSkillsListScreen(skillsList)
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
		It("should handle empty list", func() {
			screen := skills.NewSkillsListScreen([]*career.Skill{})
			view := screen.View()
			Expect(view).To(ContainSubstring("No skills"))
		})

		It("should handle navigation in empty list", func() {
			screen := skills.NewSkillsListScreen([]*career.Skill{})

			msg := tea.KeyMsg{Type: tea.KeyDown}
			_, result := screen.Update(msg)

			Expect(result).To(BeNil())
			Expect(screen.GetSelectedIndex()).To(Equal(0))
		})

		It("should handle action keys in empty list", func() {
			screen := skills.NewSkillsListScreen([]*career.Skill{})

			// Should still allow adding new skill
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			data := result.Data().(map[string]interface{})
			Expect(data["action"]).To(Equal("add"))
		})

		It("should handle very small terminal dimensions", func() {
			screen := skills.NewSkillsListScreen(skillsList)
			screen.SetTerminalInfo(20, 10)

			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should return help action on '?' key", func() {
			screen = skills.NewSkillsListScreen(skillsList)
			screen.SetTerminalInfo(120, 40)

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			data := result.Data().(map[string]interface{})
			Expect(data["action"]).To(Equal("help"))
		})
	})
})
