package intents

import (
	"time"

	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CaptureEvent - Escape Key Behavior", func() {
	var (
		intent         *CaptureEventIntent
		testCLIService *service.CLIEventService
	)

	BeforeEach(func() {
		// Create minimal test services
		testCLIService = &service.CLIEventService{}

		// Setup intent with test context
		ctx := &CaptureEventContext{
			CLIEventService: testCLIService,
			CareerService:   &careerservice.Service{},
			CaptureStrategy: "manual",
			PreviousEvent:   nil,
		}

		var err error
		intent, err = NewCaptureEventIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(intent).NotTo(BeNil())
		intent.Init()
	})

	Describe("ChooseStrategy State", func() {
		BeforeEach(func() {
			intent.state.currentState = CaptureStateChooseStrategy
		})

		It("should cancel intent when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.active).To(BeFalse())
			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

	})

	Describe("Form State", func() {
		BeforeEach(func() {
			intent.state.currentState = CaptureStateForm
		})

		It("should go back to ChooseStrategy when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.state.currentState).To(Equal(CaptureStateChooseStrategy))
			Expect(intent.active).To(BeTrue())
		})
	})

	Describe("Review State", func() {
		BeforeEach(func() {
			intent.state.currentState = CaptureStateReview
			intent.state.reviewState.Event = &career.CareerEvent{
				Text: "Test event",
				Date: time.Now(),
			}
		})

		It("should go back to Form when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.state.currentState).To(Equal(CaptureStateForm))
			Expect(intent.active).To(BeTrue())
		})
	})

	Describe("Submit State", func() {
		BeforeEach(func() {
			intent.state.currentState = CaptureStateSubmit
			intent.state.error = &IntentError{
				Code:    "TEST_ERROR",
				Message: "Submission failed",
			}
		})

		It("should go back to Review when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.state.currentState).To(Equal(CaptureStateReview))
			Expect(intent.active).To(BeTrue())
		})

		It("should keep error visible when going back", func() {
			originalError := intent.state.error

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.state.error).To(Equal(originalError))
			Expect(intent.state.error.Code).To(Equal("TEST_ERROR"))
		})
	})

	Describe("View Methods", func() {
		It("should show 'esc' and 'm' in ChooseStrategy footer", func() {
			intent.state.currentState = CaptureStateChooseStrategy
			view := intent.View()

			Expect(view).To(ContainSubstring("Esc"))
			Expect(view).To(ContainSubstring("Main Menu"))
		})

		It("should show 'esc' and 'm' in Form footer", func() {
			intent.state.currentState = CaptureStateForm
			view := intent.View()

			Expect(view).To(ContainSubstring("Esc"))
			Expect(view).To(ContainSubstring("Main Menu"))
		})

		It("should show 'esc' and 'm' in Review footer", func() {
			intent.state.currentState = CaptureStateReview
			intent.state.reviewState.Event = &career.CareerEvent{
				Text: "Test",
				Date: time.Now(),
			}
			view := intent.View()

			Expect(view).To(ContainSubstring("Esc"))
			Expect(view).To(ContainSubstring("Main Menu"))
		})

		It("should show 'esc', 'm', and 'r' in Submit footer", func() {
			intent.state.currentState = CaptureStateSubmit
			intent.state.result = &CaptureEventResult{
				Event: &career.CareerEvent{
					Text: "Test",
					Date: time.Now(),
				},
			}
			view := intent.View()

			Expect(view).To(ContainSubstring("Esc"))
			Expect(view).To(ContainSubstring("Main Menu"))
			Expect(view).To(ContainSubstring("Back"))
		})
	})
})
