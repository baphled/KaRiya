package capture_event

import (
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CaptureEvent - Escape Key Behavior", func() {
	var (
		intent         *Intent
		testCLIService *service.CLIEventService
	)

	BeforeEach(func() {
		// Create minimal test services
		testCLIService = &service.CLIEventService{}

		// Setup intent with test context
		ctx := &IntentContext{
			CLIEventService: testCLIService,
			CareerService:   &careerservice.Service{},
			CaptureStrategy: "manual",
			PreviousEvent:   nil,
		}

		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(intent).NotTo(BeNil())
		intent.Init()
	})

	Describe("ChooseStrategy State", func() {
		BeforeEach(func() {
			intent.state.currentState = StateChooseStrategy
		})

		It("should cancel intent when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.active).To(BeFalse())
			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Status).To(Equal(intents.Cancelled))
		})

	})

	Describe("Form State", func() {
		BeforeEach(func() {
			intent.state.currentState = StateForm
		})

		It("should go back to ChooseStrategy when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.state.currentState).To(Equal(StateChooseStrategy))
			Expect(intent.active).To(BeTrue())
		})
	})

	Describe("Review State", func() {
		BeforeEach(func() {
			intent.state.currentState = StateReview
			intent.state.reviewState.Event = &career.Event{
				Text: "Test event",
				Date: time.Now(),
			}
		})

		It("should go back to Form when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.state.currentState).To(Equal(StateForm))
			Expect(intent.active).To(BeTrue())
		})
	})

	Describe("Error Modal (after submission failure)", func() {
		BeforeEach(func() {
			// Set up a realistic scenario: form submitted, error occurred, modal showing
			intent.state.currentState = StateForm
			intent.state.reviewState = &ReviewInferredEventState{
				Event: &career.Event{
					Text: "Test event",
					Date: time.Now(),
				},
			}
			// Simulate error modal being shown after failed submission
			intent.Update(SubmitErrorMsg{Message: "Submission failed"})
		})

		It("should dismiss error modal when esc is pressed", func() {
			// Verify modal is showing
			Expect(intent.state.submitModal).NotTo(BeNil())

			// Press Esc to dismiss modal
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Modal should be dismissed
			Expect(intent.state.submitModal).To(BeNil())
			// Should stay in Form state (not transition back)
			Expect(intent.state.currentState).To(Equal(StateForm))
			Expect(intent.active).To(BeTrue())
		})

		It("should allow retry after dismissing error modal", func() {
			// Dismiss the error modal
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.submitModal).To(BeNil())

			// User should be able to retry submission
			// (In real workflow, they would modify form and resubmit)
			Expect(intent.state.reviewState).NotTo(BeNil())
			Expect(intent.active).To(BeTrue())
		})
	})

	Describe("View Methods", func() {
		It("should show 'esc' and 'm' in ChooseStrategy footer", func() {
			intent.state.currentState = StateChooseStrategy
			view := intent.View()

			Expect(view).To(ContainSubstring("Esc"))
			Expect(view).To(ContainSubstring("Main Menu"))
		})

		It("should show 'esc' and 'm' in Form footer", func() {
			intent.state.currentState = StateForm
			view := intent.View()

			Expect(view).To(ContainSubstring("Esc"))
			Expect(view).To(ContainSubstring("Main Menu"))
		})

		It("should show 'esc' and 'm' in Review footer", func() {
			intent.state.currentState = StateReview
			intent.state.reviewState.Event = &career.Event{
				Text: "Test",
				Date: time.Now(),
			}
			view := intent.View()

			Expect(view).To(ContainSubstring("Esc"))
			Expect(view).To(ContainSubstring("Main Menu"))
		})

		It("should show error modal overlay when submission fails", func() {
			// Set up error modal scenario
			intent.state.currentState = StateForm
			intent.state.submitModal = feedback.NewErrorModal("Save Failed", "Submission failed")

			view := intent.View()

			// Should show error modal content
			Expect(view).To(ContainSubstring("Save Failed"))
			Expect(view).To(ContainSubstring("Submission failed"))
			Expect(view).To(ContainSubstring("Esc")) // Esc to dismiss modal
		})
	})
})
