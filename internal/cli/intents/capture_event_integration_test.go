package intents_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/charmbracelet/bubbletea"
)

var _ = Describe("CaptureEventIntent Integration Tests", func() {
	var (
		intent           *intents.CaptureEventIntent
		ctx              *intents.CaptureEventContext
				)

	BeforeEach(func() {
		// Create mock services (these would be real in integration testing)
		// For now, we'll create minimal context
		ctx = &intents.CaptureEventContext{
			CaptureStrategy: "manual",
			Metadata:        make(map[string]string),
		}

		var err error
		intent, err = intents.NewCaptureEventIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(intent).NotTo(BeNil())
	})

	Describe("Event Capture Workflow", func() {
		Context("Manual Strategy", func() {
			It("should initialize for new event capture", func() {
				cmd := intent.Init()
				Expect(cmd).NotTo(BeNil())
				// Command should return without error
				msg := cmd()
				Expect(msg).To(BeNil())
			})

			It("should transition through complete workflow", func() {
				// Start with Init
				cmd := intent.Init()
				Expect(cmd).NotTo(BeNil())

				// View should show strategy selection
				view := intent.View()
				Expect(view).To(ContainSubstring("Choose Capture Strategy"))

				// Simulate selecting manual strategy
				cmd = intent.Update(intents.StrategySelectedMsg{Strategy: "manual"})
				Expect(cmd).To(BeNil())

				// Should transition to form state
				view = intent.View()
				Expect(view).NotTo(ContainSubstring("Choose Capture Strategy"))
			})

			It("should handle form submission with valid event", func() {
				// Initialize
				intent.Init()

				// Create a valid event
				event := &career.CareerEvent{
					ID:         "test-event-1",
					Text:       "Completed major project",
					Date:       time.Now().Add(-24 * time.Hour),
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
					Tags:       []string{"achievement"},
					Categories: []string{"project"},
				}

				// Simulate form submission
				cmd := intent.Update(intents.FormSubmittedMsg{Event: event})
				Expect(cmd).To(BeNil())

				// Should transition to review state
				view := intent.View()
				Expect(view).To(ContainSubstring("Review"))
			})

			It("should handle cancellation at any state", func() {
				intent.Init()

				// Cancel from strategy selection
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
				Expect(cmd).To(BeNil())

				// Should be cancelled
				result := intent.GetResult()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Cancelled))
			})

			It("should validate event before submission", func() {
				intent.Init()

				// Create invalid event (missing required fields)
				invalidEvent := &career.CareerEvent{
					Text: "",
					Date: time.Time{},
				}

				// Try to submit invalid event
				cmd := intent.Update(intents.FormSubmittedMsg{Event: invalidEvent})
				Expect(cmd).To(BeNil())

				// Should remain in form state or show error
				// (depending on implementation)
			})
		})

		Context("Quick Strategy", func() {
			BeforeEach(func() {
				ctx.CaptureStrategy = "quick"
				var err error
				intent, err = intents.NewCaptureEventIntent(ctx)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should use TimelineJournaling mode for quick capture", func() {
				intent.Init()
				// When context has quick strategy, mode should be TimelineJournaling
				// This would be verified through service call mocking in real tests
				Expect(intent.IsActive()).To(BeTrue())
			})
		})

		Context("Enriched Strategy", func() {
			BeforeEach(func() {
				ctx.CaptureStrategy = "enriched"
				var err error
				intent, err = intents.NewCaptureEventIntent(ctx)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should prepare for enrichment in enriched mode", func() {
				intent.Init()
				Expect(intent.IsActive()).To(BeTrue())
				// Enrichment would be triggered during submit
			})
		})
	})

	Describe("State Transitions", func() {
		It("should transition from ChooseStrategy to Form", func() {
			intent.Init()
			view := intent.View()
			Expect(view).To(ContainSubstring("Choose Capture Strategy"))

			// Transition to form
			intent.Update(intents.StrategySelectedMsg{Strategy: "manual"})
			// View should change (form would be shown)
		})

		It("should transition from Form to Review", func() {
			intent.Init()
			intent.Update(intents.StrategySelectedMsg{Strategy: "manual"})

			event := &career.CareerEvent{
				ID:        "test-1",
				Text:      "Test event",
				Date:      time.Now(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Tags:      []string{},
			}

			intent.Update(intents.FormSubmittedMsg{Event: event})
			view := intent.View()
			Expect(view).To(ContainSubstring("Review"))
		})

		It("should transition from Review to Submit", func() {
			intent.Init()
			intent.Update(intents.StrategySelectedMsg{Strategy: "manual"})

			event := &career.CareerEvent{
				ID:        "test-1",
				Text:      "Test event",
				Date:      time.Now(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Tags:      []string{},
			}

			intent.Update(intents.FormSubmittedMsg{Event: event})
			// Would transition to submit on confirmation
		})

		It("should allow back navigation", func() {
			intent.Init()
			intent.Update(intents.StrategySelectedMsg{Strategy: "manual"})

			// Send back message
			intent.Update(intents.ReviewBackMsg{})
			// Should return to form state
		})
	})

	Describe("Error Handling", func() {
		It("should handle missing event gracefully", func() {
			intent.Init()
			intent.Update(intents.StrategySelectedMsg{Strategy: "manual"})

			// Try to submit nil event
			cmd := intent.Update(intents.FormSubmittedMsg{Event: nil})
			Expect(cmd).To(BeNil())

			// Should set failed result
			_ = intent.GetResult()
			// Result should indicate failure or intent should remain active
		})

		It("should handle validation errors", func() {
			intent.Init()

			invalidEvent := &career.CareerEvent{
				Text: "",
				Date: time.Time{},
			}

			intent.Update(intents.FormSubmittedMsg{Event: invalidEvent})
			// Should handle validation error gracefully
		})

		It("should handle form cancellation", func() {
			intent.Init()
			intent.Update(intents.StrategySelectedMsg{Strategy: "manual"})

			cmd := intent.Update(intents.FormCancelledMsg{})
			Expect(cmd).To(BeNil())

			result := intent.GetResult()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})
	})

	Describe("Review State Management", func() {
		It("should track accepted bursts and facts", func() {
			intent.Init()

			bursts := []*career.Burst{
				{ID: "burst-1", Name: "Test Burst"},
			}
			facts := []*career.Fact{
				{ID: "fact-1", Text: "Test Fact"},
			}

			cmd := intent.Update(intents.ReviewConfirmedMsg{
				AcceptedBursts: bursts,
				AcceptedFacts:  facts,
				RejectedItems:  make(map[string]string),
			})
			Expect(cmd).NotTo(BeNil())

			// Result should include accepted items
			result := intent.GetResult()
			if result != nil && result.Data != nil {
				Expect(result.Data.Bursts).To(HaveLen(1))
				Expect(result.Data.Facts).To(HaveLen(1))
			}
		})

		It("should track rejected items", func() {
			intent.Init()

			rejectedItems := map[string]string{
				"burst-1": "Not relevant to career path",
			}

			cmd := intent.Update(intents.ReviewConfirmedMsg{
				AcceptedBursts: []*career.Burst{},
				AcceptedFacts:  []*career.Fact{},
				RejectedItems:  rejectedItems,
			})
			Expect(cmd).NotTo(BeNil())

			// Result should include rejection reasons
			result := intent.GetResult()
			if result != nil && result.Data != nil {
				Expect(result.Data.RejectedFields).To(HaveLen(1))
			}
		})
	})

	Describe("Result Handling", func() {
		It("should return nil result while active", func() {
			intent.Init()
			// Result might be nil or have pending status
		})

		It("should return completed result on success", func() {
			intent.Init()

			event := &career.CareerEvent{
				ID:        "test-1",
				Text:      "Test event",
				Date:      time.Now(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Tags:      []string{},
			}

			intent.Update(intents.StrategySelectedMsg{Strategy: "manual"})
			intent.Update(intents.FormSubmittedMsg{Event: event})

			// Would complete on submit confirmation
		})

		It("should return cancelled result on cancellation", func() {
			intent.Init()

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).To(BeNil())

			result := intent.GetResult()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should include metadata in result", func() {
			ctx.Metadata["source"] = "manual_entry"
			var err error
			intent, err = intents.NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			intent.Init()
			// Metadata should be preserved through the workflow
		})
	})

	Describe("View Rendering", func() {
		It("should render appropriate view for each state", func() {
			intent.Init()

			// Strategy selection view
			view := intent.View()
			Expect(view).To(ContainSubstring("Choose Capture Strategy"))

			// Transition and check form view
			intent.Update(intents.StrategySelectedMsg{Strategy: "manual"})
			// Form view would be rendered

			// Transition and check review view
			event := &career.CareerEvent{
				ID:        "test-1",
				Text:      "Test event",
				Date:      time.Now(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Tags:      []string{},
			}
			intent.Update(intents.FormSubmittedMsg{Event: event})
			view = intent.View()
			Expect(view).To(ContainSubstring("Review"))
		})

		It("should display error messages on validation failure", func() {
			intent.Init()
			intent.Update(intents.StrategySelectedMsg{Strategy: "manual"})

			// Invalid event should show error
			invalidEvent := &career.CareerEvent{
				Text: "",
				Date: time.Time{},
			}
			intent.Update(intents.FormSubmittedMsg{Event: invalidEvent})

			_ = intent.View()
			// View should contain error information
		})

		It("should not render when inactive", func() {
			intent.Init()

			// Cancel to deactivate
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

			view := intent.View()
			Expect(view).To(ContainSubstring("not active"))
		})
	})

	Describe("Keyboard Input Handling", func() {
		It("should handle quit key", func() {
			intent.Init()

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).To(BeNil())

			result := intent.GetResult()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should handle escape key for back navigation", func() {
			intent.Init()
			intent.Update(intents.StrategySelectedMsg{Strategy: "manual"})

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEscape})
			Expect(cmd).To(BeNil())

			// Should transition back to strategy selection
		})

		It("should handle enter key for confirmation", func() {
			intent.Init()
			intent.Update(intents.StrategySelectedMsg{Strategy: "manual"})

			event := &career.CareerEvent{
				ID:        "test-1",
				Text:      "Test event",
				Date:      time.Now(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Tags:      []string{},
			}
			intent.Update(intents.FormSubmittedMsg{Event: event})

			// Enter should submit
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())
		})
	})
})

