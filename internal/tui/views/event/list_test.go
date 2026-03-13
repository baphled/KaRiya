//nolint:errcheck // Test file - error handling for test setup is not relevant.
package event_test

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/theme"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListView", func() {
	var (
		view   *event.ListView
		events []display.Event
	)

	BeforeEach(func() {
		ev1 := fixtures.EventWith("event-1", "Backend Developer at TechCorp - Built scalable APIs", "TechCorp", "")
		ev1.Date = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		ev2 := fixtures.EventWith("event-2", "DevOps Engineer at CloudInc - Managed Kubernetes clusters", "CloudInc", "")
		ev2.Date = time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)
		ev3 := fixtures.EventWith("event-3", "Frontend Developer at WebSolutions - Created React applications", "WebSolutions", "")
		ev3.Date = time.Date(2022, 3, 10, 0, 0, 0, 0, time.UTC)
		events = display.EventsFromDomain([]*career.Event{ev1, ev2, ev3})
	})

	Describe("Construction", func() {
		It("should create an event list view", func() {
			view = event.NewListView(events)
			Expect(view).NotTo(BeNil())
		})

		It("should store events", func() {
			view = event.NewListView(events)
			Expect(view.GetEvents()).To(Equal(events))
		})

		It("should handle empty event list", func() {
			view = event.NewListView([]display.Event{})
			Expect(view).NotTo(BeNil())
			Expect(view.GetEvents()).To(BeEmpty())
		})

		It("should start with first item selected", func() {
			view = event.NewListView(events)
			Expect(view.GetSelectedIndex()).To(Equal(0))
		})
	})

	Describe("Init", func() {
		It("should return nil command", func() {
			view = event.NewListView(events)
			cmd := view.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update with WindowSizeMsg", func() {
		BeforeEach(func() {
			view = event.NewListView(events)
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
			view = event.NewListView(events)
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
			It("should return NavigateViewResult with selected event", func() {
				msg := tea.KeyMsg{Type: tea.KeyEnter}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				nav, ok := navResult.ResultData.(event.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(event.ActionView))
				Expect(nav.Event).To(Equal(events[0]))
			})

			It("should return nil when no events", func() {
				emptyView := event.NewListView([]display.Event{})
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
				nav, ok := navResult.ResultData.(event.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(event.ActionAdd))
			})

			It("should return navigate with edit action for 'e'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				nav, ok := navResult.ResultData.(event.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(event.ActionEdit))
				Expect(nav.Event).To(Equal(events[0]))
			})

			It("should return nil for edit when no events", func() {
				emptyView := event.NewListView([]display.Event{})
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
				nav, ok := navResult.ResultData.(event.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(event.ActionDelete))
				Expect(nav.Event).To(Equal(events[0]))
			})

			It("should return nil for delete when no events", func() {
				emptyView := event.NewListView([]display.Event{})
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
				nav, ok := navResult.ResultData.(event.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(event.ActionFilter))
			})

			It("should return navigate with sort action for 's'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				nav, ok := navResult.ResultData.(event.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(event.ActionSort))
			})

			It("should return navigate with search action for '/'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				nav, ok := navResult.ResultData.(event.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(event.ActionSearch))
			})

			It("should return navigate with clear action for 'x'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				nav, ok := navResult.ResultData.(event.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(event.ActionClear))
			})

			It("should return navigate with help action for '?'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				nav, ok := navResult.ResultData.(event.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(event.ActionHelp))
			})

			// --- Skill action tests inserted here ---
			It("should return EventSkillNav with show skills action for 'K'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'K'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				nav, ok := navResult.ResultData.(event.SkillNav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(event.ActionShowEventSkills))
				Expect(nav.Event).To(Equal(events[0]))
			})

			It("should return EventSkillNav with infer skills action for 'I'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'I'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				nav, ok := navResult.ResultData.(event.SkillNav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(event.ActionInferSkills))
				Expect(nav.Event).To(Equal(events[0]))
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
			view = event.NewListView(events)
			view.SetTerminalInfo(120, 40)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should contain table content", func() {
			view = event.NewListView(events)
			view.SetTerminalInfo(120, 40)
			content := view.RenderContent()
			// Should contain at least one event's date
			Expect(content).To(ContainSubstring("2024-01-01"))
		})
	})

	Describe("HelpText", func() {
		BeforeEach(func() {
			view = event.NewListView(events)
		})

		It("should return non-empty string", func() {
			Expect(view.HelpText()).NotTo(BeEmpty())
		})

		It("should contain navigation hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Navigate"))
		})

		It("should contain page hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Page"))
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

		It("should contain back hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Back"))
		})

		It("should contain quit hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Quit"))
		})

		It("should contain help hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Help"))
		})
	})

	Describe("GetEvents", func() {
		It("should return the events passed to constructor", func() {
			view = event.NewListView(events)
			Expect(view.GetEvents()).To(Equal(events))
		})
	})

	Describe("GetSelectedIndex", func() {
		It("should start at 0", func() {
			view = event.NewListView(events)
			Expect(view.GetSelectedIndex()).To(Equal(0))
		})
	})

	Describe("Skill actions with empty event list", func() {
		BeforeEach(func() {
			view = event.NewListView([]display.Event{})
			view.SetTerminalInfo(120, 40)
		})

		It("should return nil for show skills when no events", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'K'}}
			_, result := view.Update(msg)
			Expect(result).To(BeNil())
		})

		It("should return nil for infer skills when no events", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'I'}}
			_, result := view.Update(msg)
			Expect(result).To(BeNil())
		})

		It("should return nil for edit when no events", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			_, result := view.Update(msg)
			Expect(result).To(BeNil())
		})

		It("should return nil for delete when no events", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
			_, result := view.Update(msg)
			Expect(result).To(BeNil())
		})
	})

	Describe("Update with non-key messages", func() {
		BeforeEach(func() {
			view = event.NewListView(events)
			view.SetTerminalInfo(120, 40)
		})

		It("should return nil for unhandled message types", func() {
			cmd, result := view.Update(nil)
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})
	})

	Describe("HelpText with theme", func() {
		It("should return help text without theme set", func() {
			view = event.NewListView(events)
			Expect(view.HelpText()).NotTo(BeEmpty())
		})

		It("should return help text with nil theme", func() {
			view = event.NewListView(events)
			view.SetTheme(nil)
			Expect(view.HelpText()).NotTo(BeEmpty())
		})

		It("should return help text with a concrete theme", func() {
			view = event.NewListView(events)
			view.SetTheme(theme.Default())

			help := view.HelpText()
			Expect(help).To(ContainSubstring("Navigate"))
			Expect(help).To(ContainSubstring("Help"))
		})
	})

})
