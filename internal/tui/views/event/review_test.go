//nolint:errcheck // Test file - error handling for test setup is not relevant.
package event_test

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Review", func() {
	var (
		view              *event.Review
		testEvent         *career.Event
		testEventDisplay  display.Event
		testBursts        []*career.Burst
		testBurstsDisplay []display.Burst
		testFacts         []*career.Fact
		testFactsDisplay  []display.Fact
		testSkills        []skillinference.SkillSuggestion
		testSkillsDisplay []display.SkillSuggestion
	)

	BeforeEach(func() {
		testEvent = fixtures.EventWith("ev-1", "Built scalable API platform for distributed services", "TechCorp", "API Platform")
		testBursts = []*career.Burst{
			fixtures.Burst("burst-1", "ev-1"),
		}
		testFacts = []*career.Fact{
			fixtures.FactWith("fact-1", "Designed distributed API architecture"),
		}
		testSkills = []skillinference.SkillSuggestion{
			{Name: "Go", Category: "backend", Confidence: 0.95, EventIDs: []string{"ev-1"}},
		}
		testEventDisplay = display.EventFromDomain(testEvent)
		testBurstsDisplay = display.BurstsFromDomain(testBursts)
		testFactsDisplay = display.FactsFromDomain(testFacts)
		testSkillsDisplay = display.SkillSuggestionsFromDomain(testSkills)
	})

	Describe("Construction", func() {
		It("should create a non-nil Review with event, bursts, facts, skills", func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
			Expect(view).NotTo(BeNil())
		})

		It("should create Review with nil event without panic", func() {
			Expect(func() {
				view = event.NewReview(display.Event{}, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
			}).NotTo(Panic())
			Expect(view).NotTo(BeNil())
		})

		It("should create Review with empty slices", func() {
			view = event.NewReview(testEventDisplay, []display.Burst{}, []display.Fact{}, []display.SkillSuggestion{})
			Expect(view).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
		})

		It("should return nil command", func() {
			cmd := view.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update with WindowSizeMsg", func() {
		BeforeEach(func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
		})

		It("should update dimensions and return nil cmd and nil result", func() {
			msg := tea.WindowSizeMsg{Width: 120, Height: 40}
			cmd, result := view.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})
	})

	Describe("Update with KeyMsg", func() {
		BeforeEach(func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
			view.SetTerminalInfo(120, 40)
		})

		Describe("Enter key", func() {
			It("should return SubmitViewResult", func() {
				msg := tea.KeyMsg{Type: tea.KeyEnter}
				_, result := view.Update(msg)
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultSubmit))
			})

			It("should have FormData map with event, bursts, facts, skills keys", func() {
				msg := tea.KeyMsg{Type: tea.KeyEnter}
				_, result := view.Update(msg)
				submitResult, ok := result.(*widgets.SubmitViewResult)
				Expect(ok).To(BeTrue())
				formData, ok := submitResult.FormData.(event.ReviewResult)
				Expect(ok).To(BeTrue())
				Expect(formData.Event).NotTo(BeZero())
				Expect(formData.Bursts).To(BeEmpty())
				Expect(formData.Facts).To(BeEmpty())
				Expect(formData.Skills).To(BeEmpty())
			})

			It("should have nil accepted items initially", func() {
				msg := tea.KeyMsg{Type: tea.KeyEnter}
				_, result := view.Update(msg)
				submitResult, ok := result.(*widgets.SubmitViewResult)
				Expect(ok).To(BeTrue())
				formData, ok := submitResult.FormData.(event.ReviewResult)
				Expect(ok).To(BeTrue())
				Expect(formData.Bursts).To(BeNil())
				Expect(formData.Facts).To(BeNil())
				Expect(formData.Skills).To(BeNil())
			})
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

		Describe("Rune keys", func() {
			It("should return NavigateViewResult with edit_metadata for 'e'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				Expect(navResult.ResultData).To(Equal("edit_metadata"))
			})

			It("should return NavigateViewResult with suggest_bursts for 'b'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				Expect(navResult.ResultData).To(Equal("suggest_bursts"))
			})

			It("should return NavigateViewResult with suggest_facts for 'f'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				Expect(navResult.ResultData).To(Equal("suggest_facts"))
			})

			It("should return NavigateViewResult with suggest_skills for 's'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				Expect(navResult.ResultData).To(Equal("suggest_skills"))
			})
		})

		Describe("Scroll keys", func() {
			It("should return nil result for up arrow", func() {
				msg := tea.KeyMsg{Type: tea.KeyUp}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should return nil result for down arrow", func() {
				msg := tea.KeyMsg{Type: tea.KeyDown}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should return nil result for j key (vim)", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should return nil result for k key (vim)", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should return nil result for page up", func() {
				msg := tea.KeyMsg{Type: tea.KeyPgUp}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should return nil result for page down", func() {
				msg := tea.KeyMsg{Type: tea.KeyPgDown}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Describe("Unknown key", func() {
			It("should return nil cmd and nil result", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("RenderContent", func() {
		BeforeEach(func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
			view.SetTerminalInfo(120, 40)
		})

		It("should return non-empty string", func() {
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should contain Event Details text", func() {
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Event Details"))
		})

		It("should contain event text when event is provided", func() {
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Built scalable API platform for distributed services"))
		})

		It("should contain No event data when event is nil", func() {
			nilView := event.NewReview(display.Event{}, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
			nilView.SetTerminalInfo(120, 40)
			content := nilView.RenderContent()
			Expect(content).To(ContainSubstring("No event data"))
		})

		It("should contain No bursts detected when bursts is empty", func() {
			emptyView := event.NewReview(testEventDisplay, []display.Burst{}, testFactsDisplay, testSkillsDisplay)
			emptyView.SetTerminalInfo(120, 40)
			content := emptyView.RenderContent()
			Expect(content).To(ContainSubstring("No bursts detected"))
		})

		It("should contain No facts detected when facts is empty", func() {
			emptyView := event.NewReview(testEventDisplay, testBurstsDisplay, []display.Fact{}, testSkillsDisplay)
			emptyView.SetTerminalInfo(120, 40)
			content := emptyView.RenderContent()
			Expect(content).To(ContainSubstring("No facts detected"))
		})

		It("should contain No skills detected when skills is empty", func() {
			emptyView := event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, []display.SkillSuggestion{})
			emptyView.SetTerminalInfo(120, 40)
			content := emptyView.RenderContent()
			Expect(content).To(ContainSubstring("No skills detected"))
		})

		It("should contain burst name when bursts are provided", func() {
			content := view.RenderContent()
			Expect(content).To(ContainSubstring(testBursts[0].Name))
		})

		It("should contain fact text when facts are provided", func() {
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Designed distributed API architecture"))
		})

		It("should contain skill name when skills are provided", func() {
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Go"))
		})
	})

	Describe("HelpText", func() {
		BeforeEach(func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
		})

		It("should return non-empty string", func() {
			Expect(view.HelpText()).NotTo(BeEmpty())
		})

		It("should contain Edit metadata from e badge", func() {
			Expect(view.HelpText()).To(ContainSubstring("Edit metadata"))
		})

		It("should contain Edit bursts from b badge", func() {
			Expect(view.HelpText()).To(ContainSubstring("Edit bursts"))
		})

		It("should contain Edit facts from f badge", func() {
			Expect(view.HelpText()).To(ContainSubstring("Edit facts"))
		})

		It("should contain Scroll text", func() {
			Expect(view.HelpText()).To(ContainSubstring("Scroll"))
		})

		It("should contain Edit skills when skills are provided", func() {
			Expect(view.HelpText()).To(ContainSubstring("Edit skills"))
		})

		It("should not contain Edit skills when skills are empty", func() {
			emptyView := event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, []display.SkillSuggestion{})
			Expect(emptyView.HelpText()).NotTo(ContainSubstring("Edit skills"))
		})
	})

	Describe("SetAcceptedSkills", func() {
		BeforeEach(func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
			view.SetTerminalInfo(120, 40)
		})

		It("should update accepted skills", func() {
			newSkills := []skillinference.SkillSuggestion{
				{Name: "Rust", Category: "backend", Confidence: 0.85, EventIDs: []string{"ev-1"}},
			}
			view.SetAcceptedSkills(display.SkillSuggestionsFromDomain(newSkills))
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := view.Update(msg)
			submitResult, ok := result.(*widgets.SubmitViewResult)
			Expect(ok).To(BeTrue())
			formData, ok := submitResult.FormData.(event.ReviewResult)
			Expect(ok).To(BeTrue())
			Expect(formData.Skills).To(HaveLen(1))
			Expect(formData.Skills[0].Name).To(Equal("Rust"))
		})
	})

	Describe("SetAcceptedFacts", func() {
		BeforeEach(func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
			view.SetTerminalInfo(120, 40)
		})

		It("should update accepted facts", func() {
			newFacts := []*career.Fact{
				fixtures.FactWith("fact-2", "Implemented caching layer"),
			}
			view.SetAcceptedFacts(display.FactsFromDomain(newFacts))
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := view.Update(msg)
			submitResult, ok := result.(*widgets.SubmitViewResult)
			Expect(ok).To(BeTrue())
			formData, ok := submitResult.FormData.(event.ReviewResult)
			Expect(ok).To(BeTrue())
			Expect(formData.Facts).To(HaveLen(1))
			Expect(formData.Facts[0].Text).To(Equal("Implemented caching layer"))
		})
	})

	Describe("SetAcceptedBursts", func() {
		BeforeEach(func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
			view.SetTerminalInfo(120, 40)
		})

		It("should update accepted bursts", func() {
			newBursts := []*career.Burst{
				fixtures.Burst("burst-2", "ev-1"),
			}
			view.SetAcceptedBursts(display.BurstsFromDomain(newBursts))
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := view.Update(msg)
			submitResult, ok := result.(*widgets.SubmitViewResult)
			Expect(ok).To(BeTrue())
			formData, ok := submitResult.FormData.(event.ReviewResult)
			Expect(ok).To(BeTrue())
			Expect(formData.Bursts).To(HaveLen(1))
			Expect(formData.Bursts[0].ID).To(Equal("burst-2"))
		})
	})

	Describe("SetSuggestedSkills", func() {
		BeforeEach(func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
			view.SetTerminalInfo(120, 40)
		})

		It("should update suggested skills", func() {
			newSkills := []skillinference.SkillSuggestion{
				{Name: "Python", Category: "backend", Confidence: 0.75, EventIDs: []string{"ev-1"}},
			}
			view.SetSuggestedSkills(display.SkillSuggestionsFromDomain(newSkills))
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Python"))
		})
	})

	Describe("SetSuggestedBursts", func() {
		BeforeEach(func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
			view.SetTerminalInfo(120, 40)
		})

		It("should update suggested bursts", func() {
			newBursts := []*career.Burst{
				fixtures.Burst("burst-3", "ev-1"),
			}
			view.SetSuggestedBursts(display.BurstsFromDomain(newBursts))
			content := view.RenderContent()
			Expect(content).To(ContainSubstring(newBursts[0].Name))
		})
	})

	Describe("SetSuggestedFacts", func() {
		BeforeEach(func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
			view.SetTerminalInfo(120, 40)
		})

		It("should update suggested facts", func() {
			newFacts := []*career.Fact{
				fixtures.FactWith("fact-3", "Optimized database queries"),
			}
			view.SetSuggestedFacts(display.FactsFromDomain(newFacts))
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Optimized database queries"))
		})
	})

	Describe("GetSuggestedSkills", func() {
		BeforeEach(func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
		})

		It("should return the skills passed to constructor", func() {
			skills := view.GetSuggestedSkills()
			Expect(skills).To(Equal(testSkillsDisplay))
		})
	})

	Describe("Accepted items rendering", func() {
		BeforeEach(func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
			view.SetTerminalInfo(120, 40)
		})

		It("should show accepted indicator for accepted bursts", func() {
			view.SetAcceptedBursts(testBurstsDisplay)
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("●"))
		})

		It("should show accepted indicator for accepted facts", func() {
			view.SetAcceptedFacts(testFactsDisplay)
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("●"))
		})

		It("should show accepted indicator for accepted skills", func() {
			view.SetAcceptedSkills(testSkillsDisplay)
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("●"))
		})

		It("should show unaccepted indicator when no items accepted", func() {
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("○"))
		})
	})

	Describe("Viewport centering", func() {
		It("should handle zero-width viewport", func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should handle narrow viewport", func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
			view.SetTerminalInfo(40, 20)
			msg := tea.WindowSizeMsg{Width: 40, Height: 20}
			view.Update(msg)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should handle wide viewport", func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
			view.SetTerminalInfo(200, 60)
			msg := tea.WindowSizeMsg{Width: 200, Height: 60}
			view.Update(msg)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})
	})

	Describe("Theme resolution", func() {
		BeforeEach(func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
			view.SetTerminalInfo(120, 40)
		})

		It("should use default theme when no theme is set", func() {
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should render with custom theme set", func() {
			view.SetTheme(nil)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should use an explicitly configured theme for content and help text", func() {
			view.SetTheme(themes.NewDefaultTheme())

			content := view.RenderContent()
			help := view.HelpText()

			Expect(content).NotTo(BeEmpty())
			Expect(help).To(ContainSubstring("Edit metadata"))
		})
	})

	Describe("Update with non-key non-window messages", func() {
		BeforeEach(func() {
			view = event.NewReview(testEventDisplay, testBurstsDisplay, testFactsDisplay, testSkillsDisplay)
			view.SetTerminalInfo(120, 40)
		})

		It("should return nil for unhandled message types", func() {
			cmd, result := view.Update(nil)
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})
	})
})
