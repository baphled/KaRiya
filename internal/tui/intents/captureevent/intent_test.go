package captureevent

import (
	"errors"

	"github.com/baphled/kariya/internal/domain/career"
	memoryrepo "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/intents"
	burstviews "github.com/baphled/kariya/internal/tui/views/burst"
	eventview "github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/terminal"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Intent", func() {
	Describe("NewIntent", func() {
		It("should create an intent with a valid context", func() {
			ctx := &IntentValidator{CaptureStrategy: "quick"}
			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent).NotTo(BeNil())
		})

		It("should fail with an empty CaptureStrategy", func() {
			ctx := &IntentValidator{}
			intent, err := NewIntent(ctx)
			Expect(err).To(HaveOccurred())
			Expect(intent).To(BeNil())
		})

		It("should start in StateChooseStrategy", func() {
			ctx := &IntentValidator{CaptureStrategy: "quick"}
			intent, _ := NewIntent(ctx)
			Expect(intent.GetState()).To(Equal("choose_strategy"))
		})

		It("should be active after creation", func() {
			ctx := &IntentValidator{CaptureStrategy: "quick"}
			intent, _ := NewIntent(ctx)
			Expect(intent.IsActive()).To(BeTrue())
		})

		It("should initialise review state with empty slices", func() {
			ctx := &IntentValidator{CaptureStrategy: "quick"}
			intent, _ := NewIntent(ctx)
			review := intent.GetReviewState()
			Expect(review).NotTo(BeNil())
			Expect(review.AcceptedBursts).To(BeEmpty())
			Expect(review.AcceptedFacts).To(BeEmpty())
			Expect(review.RejectedItems).To(BeEmpty())
		})
	})

	Describe("Init", func() {
		var intent *Intent

		BeforeEach(func() {
			ctx := &IntentValidator{CaptureStrategy: "quick"}
			var err error
			intent, err = NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not panic", func() {
			Expect(func() { intent.Init() }).NotTo(Panic())
		})
	})

	Describe("Update", func() {
		var intent *Intent

		BeforeEach(func() {
			ctx := &IntentValidator{CaptureStrategy: "quick"}
			var err error
			intent, err = NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		Context("when intent is not active", func() {
			BeforeEach(func() {
				// Cancel to deactivate.
				intent.HandleCancel(nil)
			})

			It("should return nil", func() {
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).To(BeNil())
			})
		})

		Context("when receiving SubmitCompleteMsg", func() {
			It("should not panic", func() {
				Expect(func() {
					intent.Update(SubmitCompleteMsg{})
				}).NotTo(Panic())
			})
		})

		Context("when receiving SubmitErrorMsg", func() {
			It("should not panic", func() {
				Expect(func() {
					intent.Update(SubmitErrorMsg{
						Code:    "TEST_ERROR",
						Message: "test",
					})
				}).NotTo(Panic())
			})
		})

		Context("when receiving DismissModalMsg", func() {
			It("should not panic when no submit modal exists", func() {
				Expect(func() {
					intent.Update(DismissModalMsg{})
				}).NotTo(Panic())
			})
		})

		Context("global key handling", func() {
			It("should return quit command on '?' key without panic", func() {
				Expect(func() {
					intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
				}).NotTo(Panic())
			})
		})

		Context("screen result dispatch", func() {
			It("should process key events without panic", func() {
				Expect(func() {
					intent.Update(tea.KeyMsg{Type: tea.KeyDown})
					intent.Update(tea.KeyMsg{Type: tea.KeyUp})
					intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				}).NotTo(Panic())
			})
		})
	})

	Describe("View", func() {
		var intent *Intent

		BeforeEach(func() {
			ctx := &IntentValidator{CaptureStrategy: "quick"}
			var err error
			intent, err = NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should return a non-empty view", func() {
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should not contain panic text", func() {
			view := intent.View()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		Context("when intent is inactive", func() {
			BeforeEach(func() {
				intent.HandleCancel(nil)
			})

			It("should return an inactive message", func() {
				view := intent.View()
				Expect(view).To(ContainSubstring("not active"))
			})
		})
	})

	Describe("Result", func() {
		var intent *Intent

		BeforeEach(func() {
			ctx := &IntentValidator{CaptureStrategy: "quick"}
			var err error
			intent, err = NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should return nil before completion", func() {
			Expect(intent.Result()).To(BeNil())
		})

		Context("after cancellation", func() {
			It("should return a cancelled result", func() {
				intent.HandleCancel(nil)
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Cancelled))
			})
		})

		Context("after submit completion", func() {
			It("should return completed result with event data", func() {
				event := fixtures.EventWith("", "Test event", "", "")

				// Transition to review state with data.
				intent.SetStateForTesting(StateSubmit)
				intent.GetReviewState().Event = event

				submitData := eventview.ReviewResult{
					Event:  display.EventFromDomain(event),
					Bursts: display.BurstsFromDomain([]*career.Burst{}),
					Facts:  display.FactsFromDomain([]*career.Fact{}),
					Skills: []display.SkillSuggestion{},
				}
				intent.HandleSubmit(&widgets.SubmitViewResult{FormData: submitData})

				Expect(intent.IsActive()).To(BeFalse())
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Completed))
			})
		})
	})
})

var _ = Describe("Intent — additional coverage", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("Init — terminal dimensions branch", func() {
		It("uses terminal info dimensions when available without panicking", func() {
			info := &terminal.Info{Width: 160, Height: 50, IsValid: true}
			intent.UpdateTerminalInfo(info)
			Expect(func() { intent.Init() }).NotTo(Panic())
		})
	})

	Describe("View — submit modal branch", func() {
		It("renders modal overlay when submitModal is set", func() {
			intent.submitModal = feedback.NewErrorModal("Test", "something broke")
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders modal overlay with terminal dimensions set", func() {
			info := &terminal.Info{Width: 120, Height: 40, IsValid: true}
			intent.UpdateTerminalInfo(info)
			intent.submitModal = feedback.NewErrorModal("Test", "something broke")
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("View — editing mode overlays", func() {
		BeforeEach(func() {
			evt := fixtures.EventWith("evt-1", "Test event for view overlay", "", "")
			intent.reviewState = &ReviewInferredEventState{
				Event:         evt,
				AcceptedFacts: make([]*career.Fact, 0),
			}
			intent.activeView = eventview.NewReview(display.EventFromDomain(evt), nil, nil, nil)
		})

		It("renders burst modal overlay when EditingModeBursts is active", func() {
			suggestions := []burstfact.BurstSuggestion{
				{Name: "API Work", Description: "REST APIs", EventIDs: []string{"e1"}},
			}
			intent.reviewState.EditingMode = EditingModeBursts
			intent.reviewState.burstModal = burstviews.NewSuggestionReview(display.BurstSuggestionsFromDomain(suggestions), nil)
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders skills modal overlay when EditingModeSkills is active", func() {
			skills := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.9},
			}
			intent.reviewState.EditingMode = EditingModeSkills
			intent.reviewState.skillModal = burstviews.NewSkillSuggestion(display.SkillSuggestionsFromDomain(skills), nil)
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders facts suggestion modal overlay when EditingModeFacts is active", func() {
			facts := []career.Fact{
				*fixtures.FactWith("f1", "Led API design for teams"),
			}
			displayFacts := make([]display.Fact, len(facts))
			for idx := range facts {
				displayFacts[idx] = display.FactFromDomain(&facts[idx])
			}
			intent.reviewState.EditingMode = EditingModeFacts
			intent.reviewState.factSuggestionModal = burstviews.NewFactSuggestion(displayFacts, nil)
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders metadata modal when editing mode is metadata with nil service", func() {
			intent.reviewState.EditingMode = EditingModeMetadata
			intent.reviewState.metadataModal = nil
			intent.context.CareerService = nil

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Update — ModalAutoDismissMsg", func() {
		It("produces DismissModalMsg command", func() {
			cmd := intent.Update(feedback.ModalAutoDismissMsg{})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			_, ok := msg.(DismissModalMsg)
			Expect(ok).To(BeTrue())
		})
	})

	Describe("Update — DismissModalMsg with active submitModal", func() {
		It("transitions to review state and clears the modal after a successful save", func() {
			intent.submitModal = feedback.NewSuccessModal("Event saved!")
			event := fixtures.EventWith("evt-1", "Test event text for the modal", "", "")
			intent.reviewState = &ReviewInferredEventState{
				Event:          event,
				InferredBursts: make([]*career.Burst, 0),
				InferredFacts:  make([]*career.Fact, 0),
				InferredSkills: nil,
			}

			intent.Update(DismissModalMsg{})
			Expect(intent.submitModal).To(BeNil())
			Expect(intent.currentState).To(Equal(StateReview))
		})

		It("returns to form state when dismissing an error modal", func() {
			intent.currentState = StateForm
			intent.submitModal = feedback.NewErrorModal("Save Failed", "something went wrong")

			intent.Update(DismissModalMsg{})
			Expect(intent.submitModal).To(BeNil())
			Expect(intent.currentState).To(Equal(StateForm))
		})

		It("sets up active screen for review after a successful save", func() {
			intent.submitModal = feedback.NewSuccessModal("Event saved!")
			event := fixtures.EventWith("evt-1", "Test event text for the modal", "", "")
			intent.reviewState = &ReviewInferredEventState{
				Event: event,
			}

			intent.Update(DismissModalMsg{})
			Expect(intent.activeView).NotTo(BeNil())
		})

		It("uses terminal dimensions when available during dismiss", func() {
			info := &terminal.Info{Width: 120, Height: 40, IsValid: true}
			intent.UpdateTerminalInfo(info)
			intent.submitModal = feedback.NewSuccessModal("Event saved!")
			event := fixtures.EventWith("evt-1", "Test event text for the modal", "", "")
			intent.reviewState = &ReviewInferredEventState{
				Event: event,
			}

			Expect(func() { intent.Update(DismissModalMsg{}) }).NotTo(Panic())
		})
	})

	Describe("Update — SubmitMsg handling", func() {
		BeforeEach(func() {
			intent.currentState = StateForm
		})

		It("shows validation error modal when SubmitMsg carries an error", func() {
			err := errors.New("event text is required")
			intent.Update(SubmitMsg{Error: err})
			Expect(intent.submitModal).NotTo(BeNil())
		})

		It("processes event from SubmitMsg when event is non-nil", func() {
			event := fixtures.EventWith("", "Valid event text for testing SubmitMsg processing", "", "")
			cmd := intent.Update(SubmitMsg{Event: event})
			Expect(cmd).NotTo(BeNil())
		})

		It("returns nil when not in StateForm", func() {
			intent.currentState = StateReview
			event := fixtures.EventWith("", "Test", "", "")
			cmd := intent.Update(SubmitMsg{Event: event})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update — submitModal key handling", func() {
		BeforeEach(func() {
			intent.submitModal = feedback.NewErrorModal("Error", "something failed")
		})

		It("returns nil for non-Esc key when modal is active", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(cmd).To(BeNil())
		})

		It("returns DismissModalMsg command on Esc when modal is not loading", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			_, ok := msg.(DismissModalMsg)
			Expect(ok).To(BeTrue())
		})

		It("does not dismiss loading modal on Esc", func() {
			intent.submitModal = feedback.NewLoadingModal("Saving...", false)
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
		})

		It("handles ModalSpinnerTickMsg for loading modal", func() {
			intent.submitModal = feedback.NewLoadingModal("Saving...", false)
			intent.submitModal.Init()
			cmd := intent.Update(feedback.ModalSpinnerTickMsg{})
			_ = cmd
		})

		It("handles ModalCountdownTickMsg for success modal", func() {
			intent.submitModal = feedback.NewSuccessModal("Saved!")
			cmd := intent.Update(feedback.ModalCountdownTickMsg{})
			_ = cmd
		})

		It("returns nil for unrecognised message types when modal is active", func() {
			cmd := intent.Update("unknown message type")
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update — editing modal routing", func() {
		BeforeEach(func() {
			evt := fixtures.EventWith("evt-1", "Test event for editing modal", "", "")
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event)
			svc.SetSkillRepository(repos.Skill)
			intent.context.CareerService = svc

			intent.activeView = eventview.NewReview(display.EventFromDomain(evt), nil, nil, nil)
			intent.reviewState = &ReviewInferredEventState{
				Event:         evt,
				EditingMode:   EditingModeMetadata,
				AcceptedFacts: make([]*career.Fact, 0),
			}
			intent.reviewState.metadataModal = eventview.NewReviewEnrichment(display.EventFromDomain(evt), eventview.ReviewEnrichmentConfig{})
		})

		It("routes message through editing modal when EditingMode is not None", func() {
			Expect(func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyTab})
			}).NotTo(Panic())
		})
	})

	Describe("Update — screen active path with result", func() {
		It("handles screen update with no result", func() {
			Expect(func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			}).NotTo(Panic())
		})
	})

	Describe("Update — InferenceCompleteMsg", func() {
		It("updates review state with inferred data", func() {
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("", "Test event", "", ""),
			}
			screen := eventview.NewReview(
				display.EventFromDomain(intent.reviewState.Event),
				nil, nil, nil,
			)
			intent.activeView = screen

			skills := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.9},
			}
			cmd := intent.Update(InferenceCompleteMsg{
				InferredSkills: skills,
				InferredBursts: make([]*career.Burst, 0),
				InferredFacts:  make([]*career.Fact, 0),
			})
			Expect(cmd).To(BeNil())
			Expect(intent.reviewState.InferredSkills).To(HaveLen(1))
		})
	})
})

var _ = Describe("CaptureEvent Global Keys Enforcement", func() {

	Context("Root State Escape Behavior", func() {
		It("should cancel intent on escape from root state", func() {
			intent, err := NewIntent(&IntentValidator{
				CaptureStrategy: "manual",
				Metadata:        make(map[string]string),
			})
			Expect(err).NotTo(HaveOccurred())

			intent.Init()
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := intent.Result()
			Expect(result).NotTo(BeNil(),
				"CaptureEvent: Root state escape should produce result")
		})
	})

	Context("Quit Key Behavior", func() {
		It("should ignore 'q' key in root state (quit only from main menu)", func() {
			intent, err := NewIntent(&IntentValidator{
				CaptureStrategy: "manual",
				Metadata:        make(map[string]string),
			})
			Expect(err).NotTo(HaveOccurred())

			intent.Init()
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

			Expect(cmd).To(BeNil(),
				"CaptureEvent: 'q' key should be ignored within intent (quit only from main menu)")
		})
	})

	Context("Help Key Behavior", func() {
		It("should toggle help on '?' in root state", func() {
			intent, err := NewIntent(&IntentValidator{
				CaptureStrategy: "manual",
				Metadata:        make(map[string]string),
			})
			Expect(err).NotTo(HaveOccurred())

			intent.Init()

			baseIntent := intent.BaseIntent
			if baseIntent != nil {
				helpBefore := baseIntent.IsHelpVisible()
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
				Expect(baseIntent.IsHelpVisible()).NotTo(Equal(helpBefore),
					"CaptureEvent: Help key should toggle help modal")
			}
		})
	})

	Context("Edit vs New - Context-Aware Navigation", func() {
		It("should cancel when editing existing event (PreviousEvent != nil)", func() {
			existingEvent := fixtures.EventWith(uuid.New().String(), "Existing event", "", "")

			ctx := &IntentValidator{
				CaptureStrategy: "manual",
				PreviousEvent:   existingEvent,
				Metadata:        make(map[string]string),
			}

			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			// Transition to form state.
			intent.SetStateForTesting(StateForm)

			// Press escape - should cancel (return to caller like BrowseTimeline).
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should go back when creating new event (PreviousEvent == nil)", func() {
			ctx := &IntentValidator{
				CaptureStrategy: "manual",
				PreviousEvent:   nil,
				Metadata:        make(map[string]string),
			}

			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			// Transition to form state.
			intent.SetStateForTesting(StateForm)

			// Press escape - should go back to strategy selection.
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.GetState()).To(Equal(string(StateChooseStrategy)))
			Expect(intent.Result()).To(BeNil())
		})
	})
})
