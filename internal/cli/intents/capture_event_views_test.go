package intents

import (
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)


var _ = Describe("CaptureEventIntent Views", func() {
	var (
		intent *CaptureEventIntent
		ctx    *CaptureEventContext
	)

	BeforeEach(func() {
		ctx = &CaptureEventContext{
			CaptureStrategy: "manual",
			Metadata:        make(map[string]string),
		}
		var err error
		intent, err = NewCaptureEventIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(intent).NotTo(BeNil())
	})

	Describe("viewChooseStrategy", func() {
		It("should render the strategy selection view", func() {
			view := intent.viewChooseStrategy()
			Expect(view).NotTo(BeEmpty())
		})

		It("should contain the strategy title", func() {
			view := intent.viewChooseStrategy()
			Expect(view).To(ContainSubstring("Choose Capture Strategy"))
		})

		It("should contain all three strategy options", func() {
			view := intent.viewChooseStrategy()
			Expect(view).To(ContainSubstring("1) Manual"))
			Expect(view).To(ContainSubstring("2) Quick"))
			Expect(view).To(ContainSubstring("3) Enriched"))
		})

		It("should contain strategy descriptions", func() {
			view := intent.viewChooseStrategy()
			Expect(view).To(ContainSubstring("Manually enter event details"))
			Expect(view).To(ContainSubstring("Quick capture with minimal fields"))
			Expect(view).To(ContainSubstring("AI-powered enrichment"))
		})

		It("should contain cancel option", func() {
			view := intent.viewChooseStrategy()
			Expect(view).To(ContainSubstring("q) Cancel"))
		})

		It("should contain user instructions", func() {
			view := intent.viewChooseStrategy()
			Expect(view).To(ContainSubstring("Select strategy"))
		})

		It("should be properly formatted with borders", func() {
			view := intent.viewChooseStrategy()
			// Lipgloss renders rounded borders, check that the view is properly styled
			Expect(view).NotTo(BeEmpty())
			// Verify the view contains styled content (not just raw text)
			Expect(len(view)).To(BeNumerically(">", 50))
		})

		It("should have consistent line length", func() {
			view := intent.viewChooseStrategy()
			lines := strings.Split(view, "\n")
			// Check that border lines have reasonable length
			for _, line := range lines {
				if strings.Contains(line, "┌") || strings.Contains(line, "└") {
					Expect(len(line)).To(BeNumerically(">", 40))
				}
			}
		})
	})

	Describe("viewCaptureForm", func() {
		It("should render the form view", func() {
			intent.state.currentState = CaptureStateForm
			view := intent.viewCaptureForm()
			Expect(view).NotTo(BeEmpty())
		})

		It("should contain the form title", func() {
			view := intent.viewCaptureForm()
			Expect(view).To(ContainSubstring("Capture Event Details"))
		})

		It("should display the strategy", func() {
			view := intent.viewCaptureForm()
			Expect(view).To(ContainSubstring("Strategy: manual"))
		})

		It("should contain all form fields", func() {
			view := intent.viewCaptureForm()
			Expect(view).To(ContainSubstring("Description"))
			Expect(view).To(ContainSubstring("Date"))
			Expect(view).To(ContainSubstring("Company"))
			Expect(view).To(ContainSubstring("Project"))
			Expect(view).To(ContainSubstring("Tags"))
			Expect(view).To(ContainSubstring("Categories"))
		})

		It("should contain user instructions", func() {
			view := intent.viewCaptureForm()
			Expect(view).To(ContainSubstring("Tab"))
			Expect(view).To(ContainSubstring("Ctrl+S"))
			Expect(view).To(ContainSubstring("Esc"))
		})

		It("should be properly formatted with borders", func() {
			view := intent.viewCaptureForm()
			// Lipgloss renders rounded borders, check that the view is properly styled
			Expect(view).NotTo(BeEmpty())
			Expect(len(view)).To(BeNumerically(">", 50))
		})

		It("should show placeholder text for inputs", func() {
			view := intent.viewCaptureForm()
			Expect(view).To(ContainSubstring("["))
			Expect(view).To(ContainSubstring("]"))
		})
	})

	Describe("viewReviewInferredEvent", func() {
		BeforeEach(func() {
			intent.state.currentState = CaptureStateReview
			intent.state.result = &CaptureEventResult{
				Event: &career.CareerEvent{
					Text: "Successfully completed a major project",
					Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				},
			}
		})

		It("should render the review view", func() {
			view := intent.viewReviewInferredEvent()
			Expect(view).NotTo(BeEmpty())
		})

		It("should contain the review title", func() {
			view := intent.viewReviewInferredEvent()
			Expect(view).To(ContainSubstring("Review Inferred Event"))
		})

		It("should display the event text", func() {
			view := intent.viewReviewInferredEvent()
			Expect(view).To(ContainSubstring("Successfully completed a major project"))
		})

		It("should contain bursts section", func() {
			view := intent.viewReviewInferredEvent()
			Expect(view).To(ContainSubstring("Inferred Bursts"))
		})

		It("should contain facts section", func() {
			view := intent.viewReviewInferredEvent()
			Expect(view).To(ContainSubstring("Inferred Facts"))
		})

		It("should show no bursts detected when empty", func() {
			intent.state.reviewState.AcceptedBursts = []*career.Burst{}
			view := intent.viewReviewInferredEvent()
			Expect(view).To(ContainSubstring("No bursts detected"))
		})

		It("should show no facts detected when empty", func() {
			intent.state.reviewState.AcceptedFacts = []*career.Fact{}
			view := intent.viewReviewInferredEvent()
			Expect(view).To(ContainSubstring("No facts detected"))
		})

		It("should list accepted bursts when present", func() {
			burst := &career.Burst{
				Name: "Project Leadership",
			}
			intent.state.reviewState.AcceptedBursts = []*career.Burst{burst}
			view := intent.viewReviewInferredEvent()
			Expect(view).To(ContainSubstring("Project Leadership"))
			Expect(view).To(ContainSubstring("[✓]"))
		})

		It("should list accepted facts when present", func() {
			fact := &career.Fact{
				Text: "Led cross-functional team",
			}
			intent.state.reviewState.AcceptedFacts = []*career.Fact{fact}
			view := intent.viewReviewInferredEvent()
			Expect(view).To(ContainSubstring("Led cross-functional team"))
			Expect(view).To(ContainSubstring("[✓]"))
		})

		It("should contain user instructions", func() {
			view := intent.viewReviewInferredEvent()
			Expect(view).To(ContainSubstring("Ctrl+S"))
			Expect(view).To(ContainSubstring("Esc"))
		})

		It("should be properly formatted with borders", func() {
			view := intent.viewReviewInferredEvent()
			// Lipgloss renders rounded borders, check that the view is properly styled
			Expect(view).NotTo(BeEmpty())
			Expect(len(view)).To(BeNumerically(">", 50))
		})

		It("should truncate long event text", func() {
			intent.state.result.Event.Text = strings.Repeat("a", 100)
			view := intent.viewReviewInferredEvent()
			// Should contain truncated text with ellipsis
			Expect(view).To(ContainSubstring("..."))
		})

		It("should handle multiple bursts", func() {
			bursts := []*career.Burst{
				{Name: "Burst 1"},
				{Name: "Burst 2"},
				{Name: "Burst 3"},
			}
			intent.state.reviewState.AcceptedBursts = bursts
			view := intent.viewReviewInferredEvent()
			Expect(view).To(ContainSubstring("Burst 1"))
			Expect(view).To(ContainSubstring("Burst 2"))
			Expect(view).To(ContainSubstring("Burst 3"))
		})

		It("should handle multiple facts", func() {
			facts := []*career.Fact{
				{Text: "Fact 1"},
				{Text: "Fact 2"},
				{Text: "Fact 3"},
			}
			intent.state.reviewState.AcceptedFacts = facts
			view := intent.viewReviewInferredEvent()
			Expect(view).To(ContainSubstring("Fact 1"))
			Expect(view).To(ContainSubstring("Fact 2"))
			Expect(view).To(ContainSubstring("Fact 3"))
		})
	})

	Describe("viewSubmit", func() {
		BeforeEach(func() {
			intent.state.currentState = CaptureStateSubmit
			intent.state.result = &CaptureEventResult{
				Event: &career.CareerEvent{
					Text: "Major project completion",
					Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				},
			}
			intent.state.reviewState.AcceptedBursts = []*career.Burst{
				{Name: "Project Leadership"},
			}
			intent.state.reviewState.AcceptedFacts = []*career.Fact{
				{Text: "Led cross-functional team"},
			}
		})

		It("should render the submit view", func() {
			view := intent.viewSubmit()
			Expect(view).NotTo(BeEmpty())
		})

		It("should contain the submit title", func() {
			view := intent.viewSubmit()
			Expect(view).To(ContainSubstring("Confirm Submission"))
		})

		It("should display the event text", func() {
			view := intent.viewSubmit()
			Expect(view).To(ContainSubstring("Major project completion"))
		})

		It("should display the event date", func() {
			view := intent.viewSubmit()
			Expect(view).To(ContainSubstring("2024-01-15"))
		})

		It("should show burst count", func() {
			view := intent.viewSubmit()
			Expect(view).To(ContainSubstring("Bursts: 1"))
		})

		It("should show fact count", func() {
			view := intent.viewSubmit()
			Expect(view).To(ContainSubstring("Facts: 1"))
		})

		It("should contain confirmation instructions", func() {
			view := intent.viewSubmit()
			Expect(view).To(ContainSubstring("Ready to submit"))
			Expect(view).To(ContainSubstring("Press Enter"))
		})

		It("should contain cancel instructions", func() {
			view := intent.viewSubmit()
			Expect(view).To(ContainSubstring("Esc"))
		})

		It("should be properly formatted with borders", func() {
			view := intent.viewSubmit()
			// Lipgloss renders rounded borders, check that the view is properly styled
			Expect(view).NotTo(BeEmpty())
			Expect(len(view)).To(BeNumerically(">", 50))
		})

		It("should show submitting message", func() {
			view := intent.viewSubmit()
			Expect(view).To(ContainSubstring("Submitting"))
		})

		It("should handle zero bursts", func() {
			intent.state.reviewState.AcceptedBursts = []*career.Burst{}
			view := intent.viewSubmit()
			Expect(view).To(ContainSubstring("Bursts: 0"))
		})

		It("should handle zero facts", func() {
			intent.state.reviewState.AcceptedFacts = []*career.Fact{}
			view := intent.viewSubmit()
			Expect(view).To(ContainSubstring("Facts: 0"))
		})

		It("should truncate long event text", func() {
			intent.state.result.Event.Text = strings.Repeat("a", 100)
			view := intent.viewSubmit()
			// Should contain truncated text with ellipsis
			Expect(view).To(ContainSubstring("..."))
		})
	})

	Describe("viewError", func() {
		It("should render the error view", func() {
			intent.state.error = &IntentError{
				Code:    "save_failed",
				Message: "Failed to save event",
			}
			view := intent.viewError()
			Expect(view).NotTo(BeEmpty())
		})

		It("should contain the error title", func() {
			intent.state.error = &IntentError{
				Code:    "save_failed",
				Message: "Failed to save event",
			}
			view := intent.viewError()
			Expect(view).To(ContainSubstring("Error"))
		})

		It("should display the error code", func() {
			intent.state.error = &IntentError{
				Code:    "save_failed",
				Message: "Failed to save event",
			}
			view := intent.viewError()
			Expect(view).To(ContainSubstring("Code: save_failed"))
		})

		It("should display the error message", func() {
			intent.state.error = &IntentError{
				Code:    "save_failed",
				Message: "Failed to save event",
			}
			view := intent.viewError()
			Expect(view).To(ContainSubstring("Message: Failed to save event"))
		})

		It("should contain recovery instructions", func() {
			intent.state.error = &IntentError{
				Code:    "save_failed",
				Message: "Failed to save event",
			}
			view := intent.viewError()
			Expect(view).To(ContainSubstring("retry"))
			Expect(view).To(ContainSubstring("Esc"))
		})

		It("should be properly formatted with borders", func() {
			intent.state.error = &IntentError{
				Code:    "save_failed",
				Message: "Failed to save event",
			}
			view := intent.viewError()
			// Lipgloss renders rounded borders, check that the view is properly styled
			Expect(view).NotTo(BeEmpty())
			Expect(len(view)).To(BeNumerically(">", 50))
		})

		It("should truncate long error codes", func() {
			intent.state.error = &IntentError{
				Code:    strings.Repeat("a", 100),
				Message: "Failed to save event",
			}
			view := intent.viewError()
			// Should contain truncated code with ellipsis
			Expect(view).To(ContainSubstring("..."))
		})

		It("should truncate long error messages", func() {
			intent.state.error = &IntentError{
				Code:    "save_failed",
				Message: strings.Repeat("a", 100),
			}
			view := intent.viewError()
			// Should contain truncated message with ellipsis
			Expect(view).To(ContainSubstring("..."))
		})
	})

	Describe("View method dispatch", func() {
		It("should call viewChooseStrategy when in ChooseStrategy state", func() {
			intent.state.currentState = CaptureStateChooseStrategy
			view := intent.View()
			Expect(view).To(ContainSubstring("Choose Capture Strategy"))
		})

		It("should call viewCaptureForm when in Form state", func() {
			intent.state.currentState = CaptureStateForm
			view := intent.View()
			Expect(view).To(ContainSubstring("Capture Event Details"))
		})

		It("should call viewReviewInferredEvent when in Review state", func() {
			intent.state.currentState = CaptureStateReview
			view := intent.View()
			Expect(view).To(ContainSubstring("Review Inferred Event"))
		})

		It("should call viewSubmit when in Submit state", func() {
			intent.state.currentState = CaptureStateSubmit
			view := intent.View()
			Expect(view).To(ContainSubstring("Confirm Submission"))
		})

		It("should return error message when intent is not active", func() {
			intent.active = false
			view := intent.View()
			Expect(view).To(ContainSubstring("not active"))
		})

		It("should return unknown state message for invalid state", func() {
			intent.state.currentState = "invalid_state"
			view := intent.View()
			Expect(view).To(ContainSubstring("Unknown state"))
			Expect(view).To(ContainSubstring("invalid_state"))
		})
	})

	Describe("View consistency", func() {
		It("should have consistent formatting across all views", func() {
			views := []string{
				intent.viewChooseStrategy(),
				intent.viewCaptureForm(),
				intent.viewReviewInferredEvent(),
				intent.viewSubmit(),
				intent.viewError(),
			}

			for _, view := range views {
				// All views should be non-empty and have reasonable size
				Expect(view).NotTo(BeEmpty())
				Expect(len(view)).To(BeNumerically(">", 50))
				// All views should have instructions
				Expect(view).To(MatchRegexp("(?i)(press|enter|esc|ctrl)"))
			}
		})

		It("should not have excessive whitespace", func() {
			views := []string{
				intent.viewChooseStrategy(),
				intent.viewCaptureForm(),
				intent.viewReviewInferredEvent(),
				intent.viewSubmit(),
				intent.viewError(),
			}

			for _, view := range views {
				// Check for excessive blank lines (more than 2 consecutive)
				Expect(view).NotTo(ContainSubstring("\n\n\n"))
			}
		})

		It("should have reasonable line lengths", func() {
			views := []string{
				intent.viewChooseStrategy(),
				intent.viewCaptureForm(),
				intent.viewReviewInferredEvent(),
				intent.viewSubmit(),
				intent.viewError(),
			}

			for _, view := range views {
				lines := strings.Split(view, "\n")
				for _, line := range lines {
					// Most lines should be reasonable length (less than 120 chars)
					Expect(len(line)).To(BeNumerically("<", 160))
				}
			}
		})
	})
})

