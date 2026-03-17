//nolint:errcheck // Test file - error handling for test setup is not relevant.
package skill_test

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/views/skill"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListView", func() {
	var (
		view        *skill.ListView
		skills      []display.Skill
		eventCounts map[string]int
	)

	BeforeEach(func() {
		sk1 := display.SkillFromDomain(fixtures.SkillWith("skill-1", "Go", "Languages", "Expert"))
		sk2 := display.SkillFromDomain(fixtures.SkillWith("skill-2", "Kubernetes", "DevOps", "Intermediate"))
		sk3 := display.SkillFromDomain(fixtures.SkillWith("skill-3", "React", "Frontend", "Beginner"))
		skills = []display.Skill{sk1, sk2, sk3}
		eventCounts = map[string]int{
			"skill-1": 5,
			"skill-2": 3,
			"skill-3": 1,
		}
	})

	Describe("Construction", func() {
		It("should create a skill list view", func() {
			view = skill.NewListView(skills, eventCounts)
			Expect(view).NotTo(BeNil())
		})

		It("should store skills", func() {
			view = skill.NewListView(skills, eventCounts)
			Expect(view.GetSelectedIndex()).To(Equal(0))
		})

		It("should handle empty skill list", func() {
			view = skill.NewListView([]display.Skill{}, eventCounts)
			Expect(view).NotTo(BeNil())
		})

		It("should start with first item selected", func() {
			view = skill.NewListView(skills, eventCounts)
			Expect(view.GetSelectedIndex()).To(Equal(0))
		})

		It("should store event counts", func() {
			view = skill.NewListView(skills, eventCounts)
			Expect(view).NotTo(BeNil())
		})

		It("should handle empty event counts", func() {
			view = skill.NewListView(skills, map[string]int{})
			Expect(view).NotTo(BeNil())
		})

		It("should handle skills with nil years used", func() {
			skillNoYears := display.SkillFromDomain(fixtures.SkillWith("skill-4", "Python", "Languages", "Expert"))
			skillNoYears.YearsUsed = nil
			skillsWithNil := append(skills, skillNoYears)
			view = skill.NewListView(skillsWithNil, eventCounts)
			Expect(view).NotTo(BeNil())
		})

		It("should handle skills with zero years used", func() {
			skillZeroYears := display.SkillFromDomain(fixtures.SkillWith("skill-5", "Ruby", "Languages", "Beginner"))
			zeroYears := 0
			skillZeroYears.YearsUsed = &zeroYears
			skillsWithZero := append(skills, skillZeroYears)
			view = skill.NewListView(skillsWithZero, eventCounts)
			Expect(view).NotTo(BeNil())
		})

		It("should handle skills with empty category", func() {
			skillNoCat := display.SkillFromDomain(fixtures.SkillWith("skill-6", "Java", "", "Advanced"))
			skillsNoCat := append(skills, skillNoCat)
			view = skill.NewListView(skillsNoCat, eventCounts)
			Expect(view).NotTo(BeNil())
		})

		It("should handle skills with empty level", func() {
			skillNoLevel := display.SkillFromDomain(fixtures.SkillWith("skill-7", "C++", "Systems", ""))
			skillsNoLevel := append(skills, skillNoLevel)
			view = skill.NewListView(skillsNoLevel, eventCounts)
			Expect(view).NotTo(BeNil())
		})

		It("should handle event counts with missing skills", func() {
			countsWithExtra := map[string]int{
				"skill-1":   5,
				"skill-2":   3,
				"skill-3":   1,
				"skill-999": 10,
			}
			view = skill.NewListView(skills, countsWithExtra)
			Expect(view).NotTo(BeNil())
		})

		It("should handle event counts with zero values", func() {
			countsWithZero := map[string]int{
				"skill-1": 0,
				"skill-2": 3,
				"skill-3": 0,
			}
			view = skill.NewListView(skills, countsWithZero)
			Expect(view).NotTo(BeNil())
		})

		It("should render content with all skill variations", func() {
			skillWithYears := display.SkillFromDomain(fixtures.SkillWith("skill-8", "TypeScript", "Languages", "Expert"))
			years := 8
			skillWithYears.YearsUsed = &years
			skillNoYears := display.SkillFromDomain(fixtures.SkillWith("skill-9", "Kotlin", "Languages", "Intermediate"))
			skillNoYears.YearsUsed = nil
			skillNoCat := display.SkillFromDomain(fixtures.SkillWith("skill-10", "Scala", "", "Advanced"))
			skillNoLevel := display.SkillFromDomain(fixtures.SkillWith("skill-11", "Clojure", "Functional", ""))
			allSkills := append(skills, skillWithYears, skillNoYears, skillNoCat, skillNoLevel)
			view = skill.NewListView(allSkills, eventCounts)
			view.SetTerminalInfo(120, 40)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should initialize with proper pagination prefix", func() {
			view = skill.NewListView(skills, eventCounts)
			view.SetTerminalInfo(120, 40)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should initialize with proper empty message", func() {
			view = skill.NewListView([]display.Skill{}, eventCounts)
			view.SetTerminalInfo(120, 40)
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("No skills found"))
		})
	})

	Describe("Init", func() {
		It("should return nil command", func() {
			view = skill.NewListView(skills, eventCounts)
			cmd := view.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update with WindowSizeMsg", func() {
		BeforeEach(func() {
			view = skill.NewListView(skills, eventCounts)
		})

		It("should update terminal dimensions", func() {
			msg := tea.WindowSizeMsg{Width: 120, Height: 40}
			cmd, result := view.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})
	})

	Describe("Update with KeyMsg", func() {
		BeforeEach(func() {
			view = skill.NewListView(skills, eventCounts)
			view.SetTerminalInfo(120, 40)
		})

		Describe("Escape key", func() {
			It("should return CancelViewResult", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultCancel))
				_, ok := result.(*widgets.CancelViewResult)
				Expect(ok).To(BeTrue())
			})
		})

		Describe("Navigation keys", func() {
			It("should move down with arrow key", func() {
				msg := tea.KeyMsg{Type: tea.KeyDown}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
				Expect(view.GetSelectedIndex()).To(Equal(1))
			})

			It("should move up with arrow key", func() {
				view.Update(tea.KeyMsg{Type: tea.KeyDown})
				Expect(view.GetSelectedIndex()).To(Equal(1))

				msg := tea.KeyMsg{Type: tea.KeyUp}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
				Expect(view.GetSelectedIndex()).To(Equal(0))
			})

			It("should move down with j key (vim)", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
				Expect(view.GetSelectedIndex()).To(Equal(1))
			})

			It("should move up with k key (vim)", func() {
				view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				Expect(view.GetSelectedIndex()).To(Equal(1))

				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
				Expect(view.GetSelectedIndex()).To(Equal(0))
			})
		})

		Describe("Enter key", func() {
			It("should return NavigateViewResult with selected skill", func() {
				msg := tea.KeyMsg{Type: tea.KeyEnter}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				nav, ok := navResult.ResultData.(skill.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(skill.ActionView))
				Expect(nav.Skill).NotTo(BeNil())
			})

			It("should return nil when no skills", func() {
				emptyView := skill.NewListView([]display.Skill{}, eventCounts)
				emptyView.SetTerminalInfo(120, 40)
				msg := tea.KeyMsg{Type: tea.KeyEnter}
				cmd, result := emptyView.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Describe("Action keys", func() {
			It("should return navigate with add action for 'a'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				nav, ok := navResult.ResultData.(skill.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(skill.ActionAdd))
			})

			It("should return navigate with edit action for 'e'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				nav, ok := navResult.ResultData.(skill.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(skill.ActionEdit))
				Expect(nav.Skill).NotTo(BeNil())
			})

			It("should return nil for edit when no skills", func() {
				emptyView := skill.NewListView([]display.Skill{}, eventCounts)
				emptyView.SetTerminalInfo(120, 40)
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
				cmd, result := emptyView.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should return navigate with delete action for 'd'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				nav, ok := navResult.ResultData.(skill.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(skill.ActionDelete))
				Expect(nav.Skill).NotTo(BeNil())
			})

			It("should return nil for delete when no skills", func() {
				emptyView := skill.NewListView([]display.Skill{}, eventCounts)
				emptyView.SetTerminalInfo(120, 40)
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
				cmd, result := emptyView.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should return navigate with filter action for 'f'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				nav, ok := navResult.ResultData.(skill.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(skill.ActionFilter))
			})

			It("should return navigate with sort action for 's'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				nav, ok := navResult.ResultData.(skill.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(skill.ActionSort))
			})

			It("should return navigate with search action for '/'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				nav, ok := navResult.ResultData.(skill.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(skill.ActionSearch))
			})

			It("should return navigate with infer action for 'i'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				nav, ok := navResult.ResultData.(skill.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(skill.ActionInfer))
			})

			It("should return navigate with help action for '?'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				nav, ok := navResult.ResultData.(skill.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(skill.ActionHelp))
			})
		})

		Describe("Unknown keys", func() {
			It("should return nil for unknown key", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("RenderContent", func() {
		It("should return non-empty string", func() {
			view = skill.NewListView(skills, eventCounts)
			view.SetTerminalInfo(120, 40)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should contain table content", func() {
			view = skill.NewListView(skills, eventCounts)
			view.SetTerminalInfo(120, 40)
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Go"))
		})
	})

	Describe("HelpText", func() {
		BeforeEach(func() {
			view = skill.NewListView(skills, eventCounts)
			view.SetTerminalInfo(120, 40)
		})

		It("should contain navigate hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Navigate"))
		})

		It("should contain view hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("View"))
		})

		It("should contain add hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Add"))
		})

		It("should contain edit hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Edit"))
		})

		It("should contain delete hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Delete"))
		})

		It("should contain filter hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Filter"))
		})

		It("should contain sort hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Sort"))
		})

		It("should contain search hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Search"))
		})

		It("should contain back hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Back"))
		})

		It("should contain quit hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Quit"))
		})

		It("should contain help hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Help"))
		})

		It("should return non-empty help text", func() {
			helpText := view.HelpText()
			Expect(helpText).NotTo(BeEmpty())
		})

		It("should handle nil theme gracefully", func() {
			view.SetTheme(nil)
			helpText := view.HelpText()
			Expect(helpText).NotTo(BeEmpty())
		})

		It("should work with empty skill list", func() {
			emptyView := skill.NewListView([]display.Skill{}, eventCounts)
			emptyView.SetTerminalInfo(120, 40)
			helpText := emptyView.HelpText()
			Expect(helpText).NotTo(BeEmpty())
		})

		It("should render help text with theme", func() {
			view = skill.NewListView(skills, eventCounts)
			view.SetTerminalInfo(120, 40)
			view.SetTheme(themes.NewDefaultTheme())
			helpText := view.HelpText()
			Expect(helpText).NotTo(BeEmpty())
		})

		It("should render help text multiple times", func() {
			view = skill.NewListView(skills, eventCounts)
			view.SetTerminalInfo(120, 40)
			helpText1 := view.HelpText()
			helpText2 := view.HelpText()
			Expect(helpText1).To(Equal(helpText2))
		})

		It("should render help text after terminal resize", func() {
			view = skill.NewListView(skills, eventCounts)
			view.SetTerminalInfo(120, 40)
			helpText1 := view.HelpText()
			view.SetTerminalInfo(150, 50)
			helpText2 := view.HelpText()
			Expect(helpText1).NotTo(BeEmpty())
			Expect(helpText2).NotTo(BeEmpty())
		})
	})

	Describe("GetSelectedIndex", func() {
		It("should start at 0", func() {
			view = skill.NewListView(skills, eventCounts)
			Expect(view.GetSelectedIndex()).To(Equal(0))
		})

		It("should update after navigation", func() {
			view = skill.NewListView(skills, eventCounts)
			view.SetTerminalInfo(120, 40)
			view.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(view.GetSelectedIndex()).To(Equal(1))
		})
	})

	Describe("SetItems", func() {
		It("should update skills list", func() {
			view = skill.NewListView(skills, eventCounts)
			newSkills := []display.Skill{
				display.SkillFromDomain(fixtures.SkillWith("skill-4", "Python", "Languages", "Expert")),
				display.SkillFromDomain(fixtures.SkillWith("skill-5", "Rust", "Languages", "Beginner")),
			}
			view.SetItems(newSkills)
			Expect(view.GetSelectedIndex()).To(Equal(0))
		})

		It("should handle empty skills list", func() {
			view = skill.NewListView(skills, eventCounts)
			view.SetItems([]display.Skill{})
			Expect(view).NotTo(BeNil())
		})
	})

	Describe("SetEventCounts", func() {
		It("should update event counts", func() {
			view = skill.NewListView(skills, eventCounts)
			newCounts := map[string]int{
				"skill-1": 10,
				"skill-2": 8,
				"skill-3": 2,
			}
			view.SetEventCounts(newCounts)
			Expect(view).NotTo(BeNil())
		})

		It("should handle empty event counts", func() {
			view = skill.NewListView(skills, eventCounts)
			view.SetEventCounts(map[string]int{})
			Expect(view).NotTo(BeNil())
		})

		It("should merge with existing counts", func() {
			view = skill.NewListView(skills, eventCounts)
			newCounts := map[string]int{
				"skill-1": 15,
			}
			view.SetEventCounts(newCounts)
			Expect(view).NotTo(BeNil())
		})

		It("should update counts and render content", func() {
			view = skill.NewListView(skills, eventCounts)
			view.SetTerminalInfo(120, 40)
			newCounts := map[string]int{
				"skill-1": 20,
				"skill-2": 15,
				"skill-3": 5,
			}
			view.SetEventCounts(newCounts)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should handle counts with zero values", func() {
			view = skill.NewListView(skills, eventCounts)
			view.SetTerminalInfo(120, 40)
			newCounts := map[string]int{
				"skill-1": 0,
				"skill-2": 0,
				"skill-3": 0,
			}
			view.SetEventCounts(newCounts)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should handle counts for non-existent skills", func() {
			view = skill.NewListView(skills, eventCounts)
			view.SetTerminalInfo(120, 40)
			newCounts := map[string]int{
				"skill-999":  100,
				"skill-1000": 50,
			}
			view.SetEventCounts(newCounts)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should handle large count values", func() {
			view = skill.NewListView(skills, eventCounts)
			view.SetTerminalInfo(120, 40)
			newCounts := map[string]int{
				"skill-1": 9999,
				"skill-2": 5000,
				"skill-3": 1000,
			}
			view.SetEventCounts(newCounts)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})
	})

	var _ = Describe("Skill View Boundary", func() {
		files, err := filepath.Glob("*.go")
		Expect(err).ToNot(HaveOccurred())
		for _, file := range files {
			if strings.HasSuffix(file, "_test.go") || file == "skill_suite_test.go" {
				continue
			}
			It("should not import forbidden packages in "+file, func() {
				content, err := os.ReadFile(file)
				Expect(err).ToNot(HaveOccurred())
				str := string(content)
				forbidden := []string{
					"github.com/baphled/kariya/internal/domain/career",
					"github.com/baphled/kariya/internal/service/career",
					"github.com/baphled/kariya/internal/tui/intents",
				}
				for _, imp := range forbidden {
					Expect(str).NotTo(ContainSubstring(imp), "Forbidden import: %s", imp)
				}
			})
		}
	})
})
