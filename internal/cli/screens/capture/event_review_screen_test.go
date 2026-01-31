//nolint:errcheck // Test file - error handling for test setup is not relevant.
package capture_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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

		testEvent = &career.Event{
			ID:      "evt-1",
			Text:    "Implemented authentication system",
			Date:    now,
			Company: "TechCorp",
			Project: "Auth Service",
		}

		testBursts = []*career.Burst{
			{
				ID:          "burst-1",
				Name:        "OAuth2 Integration",
				Description: "OAuth2 with multiple providers",
				EventIDs:    []string{"evt-1"},
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		testFacts = []*career.Fact{
			{
				ID:            "fact-1",
				Text:          "Reduced login time by 50%",
				SourceEventID: "evt-1",
				CreatedAt:     now,
				UpdatedAt:     now,
			},
		}

		breadcrumbs = []string{"Main Menu", "Capture Event", "Review"}
		screen = capture.NewEventReviewScreen(breadcrumbs, testEvent, testBursts, testFacts)
	})

	Describe("Creation", func() {
		It("should create with event, bursts, and facts", func() {
			Expect(screen).NotTo(BeNil())
		})

		It("should handle nil event", func() {
			screen = capture.NewEventReviewScreen(breadcrumbs, nil, testBursts, testFacts)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle empty bursts and facts", func() {
			screen = capture.NewEventReviewScreen(breadcrumbs, testEvent, nil, nil)
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

		It("should display title using UIKit primitives", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Review Enrichment Results"))
		})

		It("should display section headers using UIKit primitives", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Event Details"))
			Expect(view).To(ContainSubstring("Inferred Bursts"))
			Expect(view).To(ContainSubstring("Inferred Facts"))
		})

		Context("with nil event", func() {
			It("should show no event data message", func() {
				screen = capture.NewEventReviewScreen(breadcrumbs, nil, testBursts, testFacts)
				view := screen.View()
				Expect(view).To(ContainSubstring("No event data"))
			})
		})

		Context("with empty bursts", func() {
			It("should show no bursts detected message", func() {
				screen = capture.NewEventReviewScreen(breadcrumbs, testEvent, nil, testFacts)
				view := screen.View()
				Expect(view).To(ContainSubstring("No bursts detected"))
			})
		})

		Context("with empty facts", func() {
			It("should show no facts detected message", func() {
				screen = capture.NewEventReviewScreen(breadcrumbs, testEvent, testBursts, nil)
				view := screen.View()
				Expect(view).To(ContainSubstring("No facts detected"))
			})
		})

		Context("with event missing optional fields", func() {
			It("should omit company when empty", func() {
				noCompanyEvent := &career.Event{
					ID:   "evt-2",
					Text: "Simple event",
					Date: testEvent.Date,
				}
				screen = capture.NewEventReviewScreen(breadcrumbs, noCompanyEvent, nil, nil)
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

		It("should hide edit bursts badge when no bursts", func() {
			screen = capture.NewEventReviewScreen(breadcrumbs, testEvent, nil, testFacts)
			view := screen.View()
			Expect(view).NotTo(ContainSubstring("Edit bursts"))
		})

		It("should hide edit facts badge when no facts", func() {
			screen = capture.NewEventReviewScreen(breadcrumbs, testEvent, testBursts, nil)
			view := screen.View()
			Expect(view).NotTo(ContainSubstring("Edit facts"))
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

		It("should return NavigateResult for 'b' (edit bursts)", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
			navResult := result.(*screens.NavigateResult)
			Expect(navResult.Data()).To(Equal("edit_bursts"))
		})

		It("should return NavigateResult for 'f' (edit facts)", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
			navResult := result.(*screens.NavigateResult)
			Expect(navResult.Data()).To(Equal("edit_facts"))
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
})
