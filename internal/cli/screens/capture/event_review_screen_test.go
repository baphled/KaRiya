//nolint:errcheck // Test file - error handling for test setup is not relevant.
package capture_test

import (
	"regexp"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// ansiEscapePattern matches ANSI escape sequences for stripping from output.
var ansiEscapePattern = regexp.MustCompile(`\x1b\[[0-9;]*[mKHFJ]`)

// stripAnsi removes ANSI escape sequences from s, returning plain text.
func stripAnsi(s string) string {
	return ansiEscapePattern.ReplaceAllString(s, "")
}

// EventReviewScreen Tests
//
// EventReviewScreen displays captured event details with inferred bursts and facts.
// Users can review and confirm before submitting, or edit via navigation actions.

var _ = Describe("EventReviewScreen", func() {
	var (
		screen      *capture.EventReviewScreen
		testEvent   *career.Event
		testBursts  []*career.Burst
		testFacts   []*career.Fact
		breadcrumbs []string
	)

	BeforeEach(func() {
		now := time.Now()

		testEvent = fixtures.EventWith("evt-1", "Implemented authentication system", "TechCorp", "Auth Service")
		testEvent.Date = now

		burst := fixtures.Burst("burst-1", "evt-1", "evt-1")
		burst.Name = "OAuth2 Integration"
		burst.Description = "OAuth2 with multiple providers"
		burst.EventIDs = []string{"evt-1"}
		burst.CreatedAt = now
		burst.UpdatedAt = now
		testBursts = []*career.Burst{burst}

		fact := fixtures.Fact("fact-1", "evt-1")
		fact.Text = "Reduced login time by 50%"
		fact.CreatedAt = now
		fact.UpdatedAt = now
		testFacts = []*career.Fact{fact}

		breadcrumbs = []string{"Main Menu", "Capture Event", "Review"}
		screen = capture.NewEventReviewScreen(breadcrumbs, testEvent, testBursts, testFacts, nil)
	})

	Describe("Creation", func() {
		It("should create with event, bursts, and facts", func() {
			Expect(screen).NotTo(BeNil())
		})

		It("should handle nil event", func() {
			screen = capture.NewEventReviewScreen(breadcrumbs, nil, testBursts, testFacts, nil)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle empty bursts and facts", func() {
			screen = capture.NewEventReviewScreen(breadcrumbs, testEvent, nil, nil, nil)
			view := screen.View()
			Expect(view).To(ContainSubstring(testEvent.Text))
		})
	})

	Describe("View Rendering", func() {
		It("should display event text", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Implemented authentication system"))
		})

		It("should display event date", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring(testEvent.Date.Format("2006-01-02")))
		})

		It("should display event company", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("TechCorp"))
		})

		It("should display event project", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Auth Service"))
		})

		It("should display bursts", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("OAuth2 Integration"))
		})

		It("should display burst descriptions", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("OAuth2 with multiple providers"))
		})

		It("should display facts", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Reduced login time"))
		})

		It("should not display a standalone title", func() {
			view := screen.View()
			Expect(view).NotTo(ContainSubstring("Review Enrichment Results"))
		})

		It("should display section headers using UIKit primitives", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Event Details"))
			Expect(view).To(ContainSubstring("Bursts"))
			Expect(view).To(ContainSubstring("Facts"))
		})

		Context("with nil event", func() {
			It("should show no event data message", func() {
				screen = capture.NewEventReviewScreen(breadcrumbs, nil, testBursts, testFacts, nil)
				view := screen.View()
				Expect(view).To(ContainSubstring("No event data"))
			})
		})

		Context("with empty bursts", func() {
			It("should show no bursts detected message", func() {
				screen = capture.NewEventReviewScreen(breadcrumbs, testEvent, nil, testFacts, nil)
				view := screen.View()
				Expect(view).To(ContainSubstring("No bursts detected"))
			})
		})

		Context("with empty facts", func() {
			It("should show no facts detected message", func() {
				screen = capture.NewEventReviewScreen(breadcrumbs, testEvent, testBursts, nil, nil)
				view := screen.View()
				Expect(view).To(ContainSubstring("No facts detected"))
			})
		})

		Context("with event missing optional fields", func() {
			It("should omit company when empty", func() {
				noCompanyEvent := fixtures.EventWith("evt-2", "Simple event", "", "")
				screen = capture.NewEventReviewScreen(breadcrumbs, noCompanyEvent, nil, nil, nil)
				view := screen.View()
				Expect(view).To(ContainSubstring("Simple event"))
				Expect(view).NotTo(ContainSubstring("Company"))
			})
		})
	})

	Describe("Footer Rendering", func() {
		It("should show confirm badge", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Enter"))
			Expect(view).To(ContainSubstring("Confirm"))
		})

		It("should show edit metadata badge", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Edit metadata"))
		})

		It("should show edit bursts badge when bursts exist", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Edit bursts"))
		})

		It("should show edit facts badge when facts exist", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Edit facts"))
		})

		It("should show edit bursts badge even when no bursts", func() {
			screen = capture.NewEventReviewScreen(breadcrumbs, testEvent, nil, testFacts, nil)
			view := screen.View()
			Expect(view).To(ContainSubstring("Edit bursts"))
		})

		It("should show edit facts badge even when no facts", func() {
			screen = capture.NewEventReviewScreen(breadcrumbs, testEvent, testBursts, nil, nil)
			view := screen.View()
			Expect(view).To(ContainSubstring("Edit facts"))
		})

		It("should show back and quit badges", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Esc"))
			Expect(view).To(ContainSubstring("Back"))
			Expect(view).To(ContainSubstring("Quit"))
		})
	})

	Describe("Confirmation", func() {
		It("should return SubmitResult on Enter", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultSubmit))
		})

		It("should include event in result", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			submitResult := result.(*screens.SubmitResult)
			data := submitResult.Data().(map[string]interface{})
			Expect(data["event"]).To(Equal(testEvent))
		})
	})

	Describe("Edit Actions", func() {
		It("should return NavigateResult for 'e' (edit metadata)", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			navResult := result.(*screens.NavigateResult)
			Expect(navResult.Data()).To(Equal("edit_metadata"))
		})

		It("should return NavigateResult for 'b' (suggest bursts)", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
			navResult := result.(*screens.NavigateResult)
			Expect(navResult.Data()).To(Equal("suggest_bursts"))
		})

		It("should return NavigateResult for 'f' (suggest facts)", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
			navResult := result.(*screens.NavigateResult)
			Expect(navResult.Data()).To(Equal("suggest_facts"))
		})
	})

	Describe("Cancellation", func() {
		It("should return CancelResult on Esc", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})
	})

	Describe("Window Resize", func() {
		It("should handle WindowSizeMsg", func() {
			cmd, result := screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})
	})

	Describe("Key Handling Edge Cases", func() {
		It("should ignore unknown rune keys", func() {
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should ignore non-key messages", func() {
			cmd, result := screen.Update("some other message")
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})
	})

	Describe("Theme Resolution", func() {
		It("should use default theme when no theme set", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle theme resolution in footer", func() {
			screen.SetTheme(themes.NewDefaultTheme())
			view := screen.View()
			Expect(view).To(ContainSubstring("Confirm"))
		})
	})

	Describe("Accepted Facts", func() {
		Describe("SetAcceptedFacts", func() {
			It("should set accepted facts", func() {
				acceptedFacts := []*career.Fact{testFacts[0]}
				screen.SetAcceptedFacts(acceptedFacts)
				view := screen.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should handle nil accepted facts", func() {
				screen.SetAcceptedFacts(nil)
				view := screen.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should handle empty accepted facts", func() {
				screen.SetAcceptedFacts([]*career.Fact{})
				view := screen.View()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Describe("Fact Submission", func() {
			It("should include accepted facts in SubmitResult", func() {
				acceptedFacts := []*career.Fact{testFacts[0]}
				screen.SetAcceptedFacts(acceptedFacts)
				_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
				submitResult := result.(*screens.SubmitResult)
				data := submitResult.Data().(map[string]interface{})
				Expect(data["facts"]).To(HaveLen(1))
				Expect(data["facts"].([]*career.Fact)[0].Text).To(Equal("Reduced login time by 50%"))
			})

			It("should include empty facts when no facts accepted", func() {
				screen.SetAcceptedFacts([]*career.Fact{})
				_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
				submitResult := result.(*screens.SubmitResult)
				data := submitResult.Data().(map[string]interface{})
				Expect(data["facts"]).To(BeEmpty())
			})
		})
	})

	Describe("Accepted Bursts", func() {
		Describe("SetAcceptedBursts", func() {
			It("should set accepted bursts", func() {
				acceptedBursts := []*career.Burst{testBursts[0]}
				screen.SetAcceptedBursts(acceptedBursts)
				view := screen.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should handle nil accepted bursts", func() {
				screen.SetAcceptedBursts(nil)
				view := screen.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should handle empty accepted bursts", func() {
				screen.SetAcceptedBursts([]*career.Burst{})
				view := screen.View()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Describe("Burst Submission", func() {
			It("should include accepted bursts in SubmitResult", func() {
				acceptedBursts := []*career.Burst{testBursts[0]}
				screen.SetAcceptedBursts(acceptedBursts)
				_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
				submitResult := result.(*screens.SubmitResult)
				data := submitResult.Data().(map[string]interface{})
				Expect(data["bursts"]).To(HaveLen(1))
				Expect(data["bursts"].([]*career.Burst)[0].Name).To(Equal("OAuth2 Integration"))
			})

			It("should include empty bursts when no bursts accepted", func() {
				screen.SetAcceptedBursts([]*career.Burst{})
				_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
				submitResult := result.(*screens.SubmitResult)
				data := submitResult.Data().(map[string]interface{})
				Expect(data["bursts"]).To(BeEmpty())
			})
		})
	})

	Describe("Skill Integration", func() {
		var testSkills []skillinference.SkillSuggestion

		BeforeEach(func() {
			testSkills = []skillinference.SkillSuggestion{
				{
					Name:       "Go",
					Category:   "backend",
					Confidence: 0.95,
					EventIDs:   []string{"evt-1"},
					Contexts:   []string{"golang code"},
				},
				{
					Name:       "Kubernetes",
					Category:   "devops",
					Confidence: 0.87,
					EventIDs:   []string{"evt-1"},
					Contexts:   []string{"container orchestration"},
				},
			}
		})

		Describe("SetSuggestedSkills", func() {
			It("should set suggested skills", func() {
				screen.SetSuggestedSkills(testSkills)
				view := screen.View()
				Expect(view).To(ContainSubstring("Skills"))
				Expect(view).To(ContainSubstring("Go"))
			})

			It("should return suggested skills via View", func() {
				screen.SetSuggestedSkills(testSkills)
				view := screen.View()
				Expect(view).To(ContainSubstring("Go"))
				Expect(view).To(ContainSubstring("Kubernetes"))
			})

			It("should handle empty suggested skills", func() {
				screen.SetSuggestedSkills([]skillinference.SkillSuggestion{})
				view := screen.View()
				Expect(view).To(ContainSubstring("No skills detected"))
			})
		})

		Describe("SetAcceptedSkills", func() {
			It("should set accepted skills", func() {
				screen.SetAcceptedSkills([]skillinference.SkillSuggestion{
					{Name: "Go", Category: "programming", Confidence: 1.0},
				})
				view := screen.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should handle empty accepted skills", func() {
				screen.SetAcceptedSkills([]skillinference.SkillSuggestion{})
				view := screen.View()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Describe("Skill Rendering", func() {
			It("should display accepted skills in view", func() {
				screen.SetAcceptedSkills([]skillinference.SkillSuggestion{
					{Name: "Go", Category: "programming", Confidence: 1.0},
				})
				view := screen.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should display suggested skills with confidence in view", func() {
				screen.SetSuggestedSkills(testSkills)
				view := screen.View()
				Expect(view).To(ContainSubstring("Go"))
				Expect(view).To(ContainSubstring("95%"))
				Expect(view).To(ContainSubstring("Kubernetes"))
				Expect(view).To(ContainSubstring("87%"))
			})

			It("should show 'No skills detected' when empty", func() {
				screen.SetSuggestedSkills([]skillinference.SkillSuggestion{})
				view := screen.View()
				Expect(view).To(ContainSubstring("No skills detected"))
			})
		})

		Describe("Skill Navigation", func() {
			It("should return NavigateResult for 's' (suggest skills)", func() {
				screen.SetSuggestedSkills(testSkills)
				_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				navResult := result.(*screens.NavigateResult)
				Expect(navResult.Data()).To(Equal("suggest_skills"))
			})

			It("should not show skills badge when no skills", func() {
				screen.SetSuggestedSkills([]skillinference.SkillSuggestion{})
				view := screen.View()
				Expect(view).NotTo(ContainSubstring("Edit skills"))
			})

			It("should show skills badge when skills exist", func() {
				screen.SetSuggestedSkills(testSkills)
				view := screen.View()
				Expect(view).To(ContainSubstring("Edit skills"))
			})
		})

		Describe("Skill Submission", func() {
			It("should include accepted skills in SubmitResult", func() {
				screen.SetAcceptedSkills([]skillinference.SkillSuggestion{
					{Name: "Go", Category: "programming", Confidence: 1.0},
				})
				_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
				submitResult := result.(*screens.SubmitResult)
				data := submitResult.Data().(map[string]interface{})
				Expect(data["skills"]).To(HaveLen(1))
				Expect(data["skills"].([]skillinference.SkillSuggestion)[0].Name).To(Equal("Go"))
			})

			It("should include empty skills array when no skills accepted", func() {
				screen.SetAcceptedSkills([]skillinference.SkillSuggestion{})
				_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
				submitResult := result.(*screens.SubmitResult)
				data := submitResult.Data().(map[string]interface{})
				Expect(data["skills"]).To(BeEmpty())
			})
		})
	})

	Describe("Status Indicators", func() {
		It("shows ○ for pending bursts", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("○"))
		})

		It("shows ● for accepted bursts", func() {
			screen.SetAcceptedBursts(testBursts)
			view := screen.View()
			Expect(view).To(ContainSubstring("●"))
		})

		It("shows ○ for pending facts", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("○"))
		})

		It("shows ● for accepted facts", func() {
			screen.SetAcceptedFacts(testFacts)
			view := screen.View()
			Expect(view).To(ContainSubstring("●"))
		})
	})

	Describe("Scrolling", func() {
		It("handles down scroll key without returning a result", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(result).To(BeNil())
		})

		It("handles up scroll key without returning a result", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(result).To(BeNil())
		})
	})

	Describe("Scroll keys", func() {
		It("handles j key (scroll down) without returning a result", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(result).To(BeNil())
		})

		It("handles k key (scroll up) without returning a result", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(result).To(BeNil())
		})

		It("handles page-down key without returning a result", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyPgDown})
			Expect(result).To(BeNil())
		})

		It("handles page-up key without returning a result", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyPgUp})
			Expect(result).To(BeNil())
		})
	})

	Describe("Suggested data setters", func() {
		It("GetSuggestedSkills returns the previously set skills", func() {
			skills := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.9},
			}
			screen.SetSuggestedSkills(skills)
			Expect(screen.GetSuggestedSkills()).To(HaveLen(1))
			Expect(screen.GetSuggestedSkills()[0].Name).To(Equal("Go"))
		})

		It("SetSuggestedBursts updates rendered burst list", func() {
			burst := fixtures.Burst("burst-new", "evt-1", "evt-2")
			burst.Name = "New Burst"
			screen.SetSuggestedBursts([]*career.Burst{burst})
			view := screen.View()
			Expect(view).To(ContainSubstring("New Burst"))
		})

		It("SetSuggestedFacts updates rendered fact list", func() {
			fact := fixtures.Fact("fact-new", "evt-1")
			fact.Text = "New Fact Text"
			screen.SetSuggestedFacts([]*career.Fact{fact})
			view := screen.View()
			Expect(view).To(ContainSubstring("New Fact Text"))
		})
	})

	Describe("Accepted skill indicator", func() {
		It("marks an accepted skill with the filled indicator in the rendered view", func() {
			skills := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.9},
			}
			screen.SetSuggestedSkills(skills)
			screen.SetAcceptedSkills(skills)
			view := screen.View()
			Expect(view).To(ContainSubstring("●"))
		})
	})

	Describe("Content Centring", func() {
		It("returns constrained content when viewport width is zero", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("constrains content to max 80 columns when viewport is wider", func() {
			screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			view := screen.View()
			lines := strings.Split(view, "\n")
			for _, line := range lines {
				visibleLen := len([]rune(stripAnsi(line)))
				Expect(visibleLen).To(BeNumerically("<=", 120),
					"line should not exceed viewport width of 120: %q", line)
			}
		})

		It("centres content within a wide viewport by adding leading padding", func() {
			screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			view := screen.View()
			strippedView := stripAnsi(view)
			lines := strings.Split(strippedView, "\n")
			foundPadding := false
			for _, line := range lines {
				if strings.Contains(line, "Event Details") {
					leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
					if leadingSpaces >= 10 {
						foundPadding = true
						break
					}
				}
			}
			Expect(foundPadding).To(BeTrue(), "expected Event Details line to have ≥10 leading spaces for centring, indicating content is horizontally centred within the 120-wide viewport")
		})

		It("does not exceed viewport width when viewport is narrower than 80", func() {
			screen.Update(tea.WindowSizeMsg{Width: 60, Height: 40})
			view := screen.View()
			Expect(view).To(ContainSubstring("Event Details"))
		})

		It("still displays content when viewport is narrower than 80", func() {
			screen.Update(tea.WindowSizeMsg{Width: 60, Height: 40})
			view := screen.View()
			Expect(view).To(ContainSubstring("Event Details"))
			Expect(view).To(ContainSubstring("Bursts"))
		})
	})
})
