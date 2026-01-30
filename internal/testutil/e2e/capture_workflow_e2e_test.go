package e2e_test

import (
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// createCaptureTestEvent creates a test event with the given text and date.
func createCaptureTestEvent(text string) *career.Event {
	return &career.Event{
		Text:      text,
		Date:      time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// createCaptureTestEventWithDetails creates a test event with full metadata.
func createCaptureTestEventWithDetails(text, company, project string, tags, categories []string) *career.Event {
	return &career.Event{
		Text:       text,
		Date:       time.Now(),
		Company:    company,
		Project:    project,
		Tags:       tags,
		Categories: categories,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// isAtMainMenu checks if the view is showing the main menu (not just breadcrumb containing "Main Menu").
func isAtMainMenu(view string) bool {
	return strings.Contains(view, "Capture Event") &&
		strings.Contains(view, "Browse Timeline") &&
		!strings.Contains(view, "Event Description") &&
		!strings.Contains(view, "Submit")
}

// isInFormState checks if the view is showing the capture form.
func isInFormState(view string) bool {
	return strings.Contains(view, "Event Description") ||
		strings.Contains(view, "Describe what you accomplished")
}

// isInStrategyState checks if the view is showing the strategy selection.
func isInStrategyState(view string) bool {
	return strings.Contains(view, "Choose Strategy") ||
		(strings.Contains(view, "Quick") && strings.Contains(view, "Manual") &&
			!strings.Contains(view, "Event Description"))
}

var _ = Describe("CaptureEvent E2E Workflow", func() {
	var env *e2e.TestEnv

	BeforeEach(func() {
		env = e2e.Setup(GinkgoT())
	})

	AfterEach(func() {
		env.Cleanup()
	})

	Describe("Standard Capture Workflow (Quick Mode)", func() {
		It("should navigate from main menu to form via strategy selection", func() {
			view := env.GetView()
			Expect(isAtMainMenu(view)).To(BeTrue(),
				"Should be at main menu showing menu options")

			env.SelectIntentByName("capture_event")

			view = env.GetView()
			Expect(isInStrategyState(view)).To(BeTrue(),
				"Should be in strategy selection state")

			env.Confirm()

			view = env.GetView()
			Expect(isInFormState(view)).To(BeTrue(),
				"Should be in form state")
		})

		It("should trigger submit when event is submitted from form", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEventText := "Built REST API with Go, PostgreSQL, JWT auth, and Redis caching"
			testEvent := createCaptureTestEvent(testEventText)
			env.SubmitEvent(testEvent)

			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Saving"),
				ContainSubstring("Loading"),
				ContainSubstring("Success"),
				ContainSubstring("Review"),
			), "Should be in submit or post-submit state")
		})

		It("should complete intent without re-saving in post-save review", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := createCaptureTestEvent("Test event for post-save review")
			env.SubmitEvent(testEvent)

			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Review"),
				ContainSubstring("Enrichment"),
			), "Should be on post-save review screen")

			env.Confirm()

			view = env.GetView()
			Expect(view).NotTo(ContainSubstring("Saving"),
				"Should NOT show saving modal in post-save review")
			Expect(isAtMainMenu(view)).To(BeTrue(),
				"Should return to main menu after confirming post-save review")

			env.AssertEventCount(1)
		})
	})

	Describe("Quick Submit Path (Ctrl+S from Form)", func() {
		It("should submit when using Ctrl+S from form", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			env.TypeText("Quick submit test event")

			env.PressKey(tea.KeyCtrlS)

			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Saving"),
				ContainSubstring("Success"),
				ContainSubstring("Review"),
				And(ContainSubstring("Capture Event"), ContainSubstring("Browse Timeline")),
			), "Should be in submit, post-submit, or main menu state")
		})
	})

	Describe("Cancel Behavior", func() {
		It("should not persist event when cancelled from form", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			view := env.GetView()
			Expect(isInFormState(view)).To(BeTrue(),
				"Should be in form state")

			env.Cancel()

			env.AssertEventCount(0)
		})

		It("should return to previous state when cancelled from strategy selection", func() {
			env.SelectIntentByName("capture_event")

			view := env.GetView()
			Expect(isInStrategyState(view)).To(BeTrue(),
				"Should be in strategy selection state")

			env.Cancel()

			view = env.GetView()
			Expect(isAtMainMenu(view)).To(BeTrue(),
				"Should be at main menu after cancelling strategy selection")

			env.AssertEventCount(0)
		})
	})

	Describe("Manual Strategy Workflow", func() {
		It("should support manual strategy with all fields", func() {
			env.SelectIntentByName("capture_event")

			env.NavigateDown()
			env.Confirm()

			testEvent := createCaptureTestEventWithDetails(
				"Manual capture test event with rich context",
				"Test Company Inc",
				"Test Project Alpha",
				nil,
				nil,
			)
			env.SubmitEvent(testEvent)

			maxAttempts := 10
			for range maxAttempts {
				view := env.GetView()
				if isAtMainMenu(view) {
					break
				}
				env.Confirm()
			}

			env.AssertEventCount(1)

			events := env.GetEvents()
			Expect(events).To(HaveLen(1))
			Expect(events[0].Company).To(Equal("Test Company Inc"))
			Expect(events[0].Project).To(Equal("Test Project Alpha"))
		})
	})

	Describe("Escape Key Navigation", func() {
		It("should navigate back through states correctly", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			view := env.GetView()
			Expect(isInFormState(view)).To(BeTrue(),
				"Should be in form state")

			env.Cancel()

			view = env.GetView()
			Expect(isInStrategyState(view)).To(BeTrue(),
				"Should be back in strategy selection state")

			env.Cancel()

			view = env.GetView()
			Expect(isAtMainMenu(view)).To(BeTrue(),
				"Should be at main menu")
		})
	})
})
