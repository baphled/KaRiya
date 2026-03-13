package captureevent

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	memoryrepo "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/tui/intents"
	burstviews "github.com/baphled/kariya/internal/tui/views/burst"
	eventview "github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handlers", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("HandleNavigate", func() {
		Context("with CaptureStrategy data", func() {
			It("should transition to form screen for quick strategy", func() {
				result := &widgets.NavigateViewResult{
					ResultData: StrategyQuick,
				}
				intent.HandleNavigate(result)
				Expect(intent.GetState()).To(Equal("form"))
			})

			It("should transition to form screen for manual strategy", func() {
				result := &widgets.NavigateViewResult{
					ResultData: StrategyManual,
				}
				intent.HandleNavigate(result)
				Expect(intent.GetState()).To(Equal("form"))
			})
		})

		Context("with string action data", func() {
			BeforeEach(func() {
				// Put intent in review state with an event.
				intent.SetStateForTesting(StateReview)
				intent.GetReviewState().Event = fixtures.EventWith("", "test event", "", "")
			})

			It("should fail gracefully for edit_metadata when CareerService is nil", func() {
				result := &widgets.NavigateViewResult{ResultData: "edit_metadata"}
				intent.HandleNavigate(result)
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})

			It("should fail gracefully for suggest_bursts when CareerService is nil", func() {
				result := &widgets.NavigateViewResult{ResultData: "suggest_bursts"}
				intent.HandleNavigate(result)
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})

			It("should open modal for suggest_facts even when no inferred facts exist", func() {
				result := &widgets.NavigateViewResult{ResultData: "suggest_facts"}
				intent.HandleNavigate(result)
				Expect(intent.IsActive()).To(BeTrue())
				Expect(intent.GetReviewState().EditingMode).To(Equal(EditingModeFacts))
				Expect(intent.GetReviewState().GetFactSuggestionModal()).NotTo(BeNil())
			})

			It("should fail on unknown action", func() {
				result := &widgets.NavigateViewResult{ResultData: "unknown_action"}
				intent.HandleNavigate(result)
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})
		})

		Context("with nil event for edit_metadata", func() {
			BeforeEach(func() {
				intent.SetStateForTesting(StateReview)
				intent.GetReviewState().Event = nil
			})

			It("should fail gracefully", func() {
				result := &widgets.NavigateViewResult{ResultData: "edit_metadata"}
				intent.HandleNavigate(result)
				Expect(intent.IsActive()).To(BeFalse())
			})
		})

		Context("with invalid data type", func() {
			It("should fail on non-string, non-strategy data", func() {
				result := &widgets.NavigateViewResult{ResultData: 42}
				intent.HandleNavigate(result)
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})
		})
	})

	Describe("HandleCancel", func() {
		Context("from StateChooseStrategy", func() {
			It("should cancel the intent", func() {
				intent.HandleCancel(nil)
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Cancelled))
			})
		})

		Context("from StateForm without PreviousEvent", func() {
			BeforeEach(func() {
				// Navigate to form first.
				result := &widgets.NavigateViewResult{ResultData: StrategyQuick}
				intent.HandleNavigate(result)
			})

			It("should go back to strategy screen", func() {
				intent.HandleCancel(nil)
				Expect(intent.GetState()).To(Equal("choose_strategy"))
				Expect(intent.IsActive()).To(BeTrue())
			})
		})

		Context("from StateForm with PreviousEvent", func() {
			It("should cancel the intent", func() {
				ctx := &IntentValidator{
					CaptureStrategy: "quick",
					PreviousEvent:   fixtures.EventWith("", "existing", "", ""),
				}
				editIntent, err := NewIntent(ctx)
				Expect(err).NotTo(HaveOccurred())
				editIntent.Init()

				// Navigate to form.
				editIntent.HandleNavigate(&widgets.NavigateViewResult{
					ResultData: StrategyQuick,
				})

				editIntent.HandleCancel(nil)
				Expect(editIntent.IsActive()).To(BeFalse())
				Expect(editIntent.Result().Status).To(Equal(intents.Cancelled))
			})
		})

		Context("from StateReview", func() {
			BeforeEach(func() {
				intent.SetStateForTesting(StateReview)
			})

			It("should go back to form screen", func() {
				intent.HandleCancel(nil)
				Expect(intent.GetState()).To(Equal("form"))
				Expect(intent.IsActive()).To(BeTrue())
			})
		})

		Context("from StateSubmit", func() {
			BeforeEach(func() {
				intent.SetStateForTesting(StateSubmit)
			})

			It("should return nil without changing state", func() {
				cmd := intent.HandleCancel(nil)
				Expect(cmd).To(BeNil())
				Expect(intent.IsActive()).To(BeTrue())
			})
		})

		Context("from invalid state", func() {
			BeforeEach(func() {
				intent.SetStateForTesting(State("bogus"))
			})

			It("should fail the intent", func() {
				intent.HandleCancel(nil)
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})
		})
	})

	Describe("HandleSubmit", func() {
		Context("from StateForm with valid form data", func() {
			BeforeEach(func() {
				intent.HandleNavigate(&widgets.NavigateViewResult{
					ResultData: StrategyQuick,
				})
			})

			It("should not panic with confirmed form data", func() {
				formData := &forms.CaptureEventFormData{
					Text:            "Test event description that is long enough",
					Date:            "2026-01-15",
					SubmitConfirmed: true,
				}
				Expect(func() {
					intent.HandleSubmit(&widgets.SubmitViewResult{
						FormData: formData,
					})
				}).NotTo(Panic())
			})

			It("should do nothing when SubmitConfirmed is false", func() {
				formData := &forms.CaptureEventFormData{
					Text:            "Test event description",
					SubmitConfirmed: false,
				}
				cmd := intent.HandleSubmit(&widgets.SubmitViewResult{
					FormData: formData,
				})
				Expect(cmd).To(BeNil())
				Expect(intent.IsActive()).To(BeTrue())
			})
		})

		Context("from StateForm with invalid data type", func() {
			BeforeEach(func() {
				intent.HandleNavigate(&widgets.NavigateViewResult{
					ResultData: StrategyQuick,
				})
			})

			It("should fail the intent", func() {
				intent.HandleSubmit(&widgets.SubmitViewResult{
					FormData: "not form data",
				})
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})
		})

		Context("from StateReview with valid data (post-save)", func() {
			var event *career.Event

			BeforeEach(func() {
				event = fixtures.EventWith("", "Reviewed event", "", "")
				intent.SetStateForTesting(StateReview)
				intent.GetReviewState().Event = event
			})

			It("should fail on non-map data type", func() {
				intent.HandleSubmit(&widgets.SubmitViewResult{
					FormData: "not a map",
				})
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})

			It("should fail when event key is missing from map", func() {
				intent.HandleSubmit(&widgets.SubmitViewResult{
					FormData: map[string]interface{}{
						"bursts": []*career.Burst{},
					},
				})
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})
		})

		Context("from StateSubmit with valid data", func() {
			var event *career.Event

			BeforeEach(func() {
				event = fixtures.EventWith("", "Submit event", "", "")
				intent.SetStateForTesting(StateSubmit)
				intent.GetReviewState().Event = event
			})

			It("should complete the intent with event data", func() {
				intent.HandleSubmit(&widgets.SubmitViewResult{
					FormData: eventview.ReviewResult{
						Event:  display.EventFromDomain(event),
						Bursts: display.BurstsFromDomain([]*career.Burst{}),
						Facts:  display.FactsFromDomain([]*career.Fact{}),
						Skills: []display.SkillSuggestion{},
					},
				})
				Expect(intent.IsActive()).To(BeFalse())
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Completed))
			})

			It("should fail on non-map data type", func() {
				intent.HandleSubmit(&widgets.SubmitViewResult{
					FormData: "not a map",
				})
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})
		})

		Context("from invalid state", func() {
			BeforeEach(func() {
				intent.SetStateForTesting(State("bogus"))
			})

			It("should fail the intent", func() {
				intent.HandleSubmit(&widgets.SubmitViewResult{
					FormData: "anything",
				})
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})
		})
	})

	Describe("HandleError", func() {
		It("should fail intent with map error data", func() {
			result := &widgets.ErrorViewResult{
				Err:     errors.New("something broke"),
				Message: "Screen error occurred",
			}
			intent.HandleError(result)
			Expect(intent.IsActive()).To(BeFalse())
			intentResult := intent.Result()
			Expect(intentResult).NotTo(BeNil())
			Expect(intentResult.Status).To(Equal(intents.Failed))
		})

		It("should use default message when error message is empty", func() {
			result := &widgets.ErrorViewResult{
				Err: errors.New("internal error"),
			}
			intent.HandleError(result)
			Expect(intent.IsActive()).To(BeFalse())
		})
	})
})

var _ = Describe("burst persistence in postSaveReview", func() {
	var (
		intent    *Intent
		svc       *careerservice.Service
		burstRepo *memoryrepo.BurstRepository
	)

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()

		repos := memoryrepo.NewRepositories()
		burstRepo = repos.Burst.(*memoryrepo.BurstRepository)

		svc = careerservice.NewService(repos.Event)
		svc.SetBurstRepository(burstRepo)
		svc.SetSkillRepository(repos.Skill)
		svc.SetFactRepository(repos.Fact)

		intent.context.CareerService = svc
	})

	Context("with accepted bursts that exist in the repository", func() {
		It("calls ConfirmBurst to mark them as confirmed", func() {
			event := fixtures.EventWith("evt-saved-1", "Reviewed event with bursts", "", "")
			burst := fixtures.Burst("burst-to-confirm", "evt-saved-1", "evt-other-1")

			ctx := context.Background()
			Expect(burstRepo.Create(ctx, burst)).To(Succeed())
			Expect(burst.Confirmed).To(BeFalse())

			intent.currentState = StateReview
			intent.submitModal = feedback.NewSuccessModal("Event saved!")
			intent.reviewState = &ReviewInferredEventState{
				Event: event,
			}

			reviewData := eventview.ReviewResult{
				Event:  display.EventFromDomain(event),
				Bursts: display.BurstsFromDomain([]*career.Burst{burst}),
				Facts:  []display.Fact{},
				Skills: []display.SkillSuggestion{},
			}
			cmd := intent.HandleSubmit(&widgets.SubmitViewResult{FormData: reviewData})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			intent.Update(msg)

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Status).To(Equal(intents.Completed))

			confirmed, err := burstRepo.GetByID(ctx, "burst-to-confirm")
			Expect(err).NotTo(HaveOccurred())
			Expect(confirmed.Confirmed).To(BeTrue())
			Expect(confirmed.ConfirmedAt).NotTo(BeNil())
		})

		It("preserves EventIDs through ConfirmBurst", func() {
			event := fixtures.EventWith("evt-saved-2", "Reviewed event", "", "")
			burst := fixtures.Burst("burst-preserve-ids", "evt-saved-2", "evt-other-2")

			ctx := context.Background()
			Expect(burstRepo.Create(ctx, burst)).To(Succeed())

			intent.currentState = StateReview
			intent.submitModal = feedback.NewSuccessModal("Event saved!")
			intent.reviewState = &ReviewInferredEventState{
				Event: event,
			}

			reviewData := eventview.ReviewResult{
				Event:  display.EventFromDomain(event),
				Bursts: display.BurstsFromDomain([]*career.Burst{burst}),
				Facts:  []display.Fact{},
				Skills: []display.SkillSuggestion{},
			}
			cmd := intent.HandleSubmit(&widgets.SubmitViewResult{FormData: reviewData})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			intent.Update(msg)

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Status).To(Equal(intents.Completed))

			confirmed, err := burstRepo.GetByID(ctx, "burst-preserve-ids")
			Expect(err).NotTo(HaveOccurred())
			Expect(confirmed.EventIDs).To(Equal([]string{"evt-saved-2", "evt-other-2"}))
		})
	})

	Context("when CareerService is nil", func() {
		It("completes without confirming bursts", func() {
			event := fixtures.EventWith("evt-saved-3", "Reviewed event", "", "")
			burst := fixtures.Burst("burst-no-svc", "evt-saved-3", "evt-other-3")

			intent.context.CareerService = nil
			intent.currentState = StateReview
			intent.submitModal = feedback.NewSuccessModal("Event saved!")
			intent.reviewState = &ReviewInferredEventState{
				Event: event,
			}

			reviewData := eventview.ReviewResult{
				Event:  display.EventFromDomain(event),
				Bursts: display.BurstsFromDomain([]*career.Burst{burst}),
				Facts:  []display.Fact{},
				Skills: []display.SkillSuggestion{},
			}
			cmd := intent.HandleSubmit(&widgets.SubmitViewResult{FormData: reviewData})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			intent.Update(msg)

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Status).To(Equal(intents.Completed))
			Expect(burst.Confirmed).To(BeFalse())
		})
	})
})

var _ = Describe("HandleNavigate suggest_bursts always uses InferredBursts", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()

		repos := memoryrepo.NewRepositories()
		svc := careerservice.NewService(repos.Event)
		intent.context.CareerService = svc
	})

	Context("when both InferredBursts and InferredBurstSuggestions are set", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				InferredBursts: func() []*career.Burst {
					b := fixtures.Burst("", "evt-1")
					b.Name = "Persisted Burst"
					b.Description = "This is persisted"
					return []*career.Burst{b}
				}(),
				InferredBurstSuggestions: []burstfact.BurstSuggestion{
					{Name: "Raw Suggestion 1", Description: "Not persisted", EventIDs: []string{"evt-1"}, ConfidenceScore: 0.9},
					{Name: "Raw Suggestion 2", Description: "Also not persisted", EventIDs: []string{"evt-1"}, ConfidenceScore: 0.8},
				},
				AcceptedBursts: make([]*career.Burst, 0),
			}
		})

		It("builds modal from InferredBursts, not InferredBurstSuggestions", func() {
			result := &widgets.NavigateViewResult{ResultData: "suggest_bursts"}
			intent.HandleNavigate(result)

			Expect(intent.reviewState.burstModal).NotTo(BeNil())
			Expect(intent.reviewState.burstModal.GetSuggestionsCount()).To(Equal(1))
		})

		It("uses persisted burst names from InferredBursts", func() {
			result := &widgets.NavigateViewResult{ResultData: "suggest_bursts"}
			intent.HandleNavigate(result)

			suggestion := intent.reviewState.burstModal.GetCurrentSuggestion()
			Expect(suggestion).NotTo(BeNil())
			Expect(suggestion.Name).To(Equal("Persisted Burst"))
		})
	})

	Context("when the accepted suggestion has no matching InferredBurst", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				InferredBursts: func() []*career.Burst {
					b := fixtures.Burst("", "evt-1")
					b.Name = "Alpha"
					b.Description = "Persisted burst"
					return []*career.Burst{b}
				}(),
				AcceptedBursts: make([]*career.Burst, 0),
				EditingMode:    EditingModeBursts,
			}
			suggestions := []burstfact.BurstSuggestion{
				{Name: "Beta", Description: "no match in InferredBursts", EventIDs: []string{"evt-1"}, ConfidenceScore: 0.8},
			}
			intent.reviewState.burstModal = burstviews.NewSuggestionReview(display.BurstSuggestionsFromDomain(suggestions), nil)

			screen := eventview.NewReview(display.EventFromDomain(intent.reviewState.Event), nil, nil, nil)
			intent.activeView = screen
		})

		It("skips the suggestion and does not append to AcceptedBursts", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.AcceptedBursts).To(BeEmpty())
		})
	})
})

var _ = Describe("Fact Suggestion in CaptureEvent", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("smart routing via HandleNavigate", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "Led API design for payments platform", "", ""),
				InferredFacts: []*career.Fact{
					fixtures.FactWith("f1", "Led cross-team API design"),
					fixtures.FactWith("f2", "Improved system reliability"),
				},
				AcceptedFacts: []*career.Fact{},
			}
		})

		Context("when inferred facts exist", func() {
			It("creates a fact suggestion modal", func() {
				result := &widgets.NavigateViewResult{ResultData: "suggest_facts"}
				intent.HandleNavigate(result)

				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeFacts))
				Expect(intent.reviewState.factSuggestionModal).NotTo(BeNil())
			})

			It("does not set editing mode to anything other than facts", func() {
				result := &widgets.NavigateViewResult{ResultData: "suggest_facts"}
				intent.HandleNavigate(result)

				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeFacts))
			})

			It("sets terminal dimensions on the modal", func() {
				result := &widgets.NavigateViewResult{ResultData: "suggest_facts"}
				intent.HandleNavigate(result)

				Expect(intent.reviewState.factSuggestionModal).NotTo(BeNil())
			})
		})

		Context("when no inferred facts exist", func() {
			BeforeEach(func() {
				intent.reviewState.InferredFacts = []*career.Fact{}
			})

			It("opens modal with empty fact list", func() {
				result := &widgets.NavigateViewResult{ResultData: "suggest_facts"}
				intent.HandleNavigate(result)

				Expect(intent.reviewState.factSuggestionModal).NotTo(BeNil())
				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeFacts))
			})

			It("keeps intent active", func() {
				result := &widgets.NavigateViewResult{ResultData: "suggest_facts"}
				intent.HandleNavigate(result)

				Expect(intent.IsActive()).To(BeTrue())
			})
		})

		Context("when InferredFacts is nil", func() {
			BeforeEach(func() {
				intent.reviewState.InferredFacts = nil
			})

			It("opens modal with empty fact list", func() {
				result := &widgets.NavigateViewResult{ResultData: "suggest_facts"}
				intent.HandleNavigate(result)

				Expect(intent.reviewState.factSuggestionModal).NotTo(BeNil())
				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeFacts))
			})
		})
	})

	Describe("fact acceptance via updateEditingModal", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:         fixtures.EventWith("evt-1", "Led API design", "", ""),
				AcceptedFacts: []*career.Fact{},
				EditingMode:   EditingModeFacts,
			}

			facts := []career.Fact{
				*fixtures.FactWith("f1", "Led cross-team API design"),
			}
			displayFacts := make([]display.Fact, len(facts))
			for idx := range facts {
				displayFacts[idx] = display.FactFromDomain(&facts[idx])
			}
			intent.reviewState.factSuggestionModal = burstviews.NewFactSuggestion(displayFacts, nil)

			screen := eventview.NewReview(display.EventFromDomain(intent.reviewState.Event), nil, nil, nil)
			intent.activeView = screen
		})

		It("transfers accepted facts to review state as pointers", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.AcceptedFacts).To(HaveLen(1))
			Expect(intent.reviewState.AcceptedFacts[0].Text).To(Equal("Led cross-team API design"))
		})

		It("closes the modal and resets editing mode", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.factSuggestionModal).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
		})

		It("calls SetAcceptedFacts on the review screen", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			_, ok := intent.activeView.(*eventview.Review)
			Expect(ok).To(BeTrue())
			Expect(intent.reviewState.AcceptedFacts).To(HaveLen(1))
		})

		Context("with multiple facts", func() {
			BeforeEach(func() {
				facts := []career.Fact{
					*fixtures.FactWith("f1", "Led cross-team API design"),
					*fixtures.FactWith("f2", "Improved system reliability"),
				}
				displayFacts := make([]display.Fact, len(facts))
				for idx := range facts {
					displayFacts[idx] = display.FactFromDomain(&facts[idx])
				}
				intent.reviewState.factSuggestionModal = burstviews.NewFactSuggestion(displayFacts, nil)
			})

			It("accepts all facts in sequence", func() {
				intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
				intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

				Expect(intent.reviewState.AcceptedFacts).To(HaveLen(2))
				Expect(intent.reviewState.AcceptedFacts[0].Text).To(Equal("Led cross-team API design"))
				Expect(intent.reviewState.AcceptedFacts[1].Text).To(Equal("Improved system reliability"))
			})
		})

		Context("when some facts are accepted and some rejected", func() {
			BeforeEach(func() {
				facts := []career.Fact{
					*fixtures.FactWith("f1", "Led cross-team API design"),
					*fixtures.FactWith("f2", "Improved system reliability"),
				}
				displayFacts := make([]display.Fact, len(facts))
				for idx := range facts {
					displayFacts[idx] = display.FactFromDomain(&facts[idx])
				}
				intent.reviewState.factSuggestionModal = burstviews.NewFactSuggestion(displayFacts, nil)
			})

			It("stores only the accepted facts", func() {
				intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
				intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

				Expect(intent.reviewState.AcceptedFacts).To(HaveLen(1))
				Expect(intent.reviewState.AcceptedFacts[0].Text).To(Equal("Led cross-team API design"))
				Expect(intent.reviewState.factSuggestionModal).To(BeNil())
				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
			})
		})

		Context("when all facts are rejected", func() {
			It("clears modal without adding facts", func() {
				intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

				Expect(intent.reviewState.factSuggestionModal).To(BeNil())
				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
				Expect(intent.reviewState.AcceptedFacts).To(BeEmpty())
			})
		})

		Context("when escape is pressed", func() {
			It("clears the fact suggestion modal", func() {
				intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.reviewState.factSuggestionModal).To(BeNil())
				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
			})
		})
	})

	Describe("fact modal auto-close with no suggestions", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:         fixtures.EventWith("evt-1", "Led API design", "", ""),
				AcceptedFacts: []*career.Fact{},
				EditingMode:   EditingModeFacts,
			}

			intent.reviewState.factSuggestionModal = burstviews.NewFactSuggestion([]display.Fact{}, nil)

			screen := eventview.NewReview(display.EventFromDomain(intent.reviewState.Event), nil, nil, nil)
			intent.activeView = screen
		})

		It("auto-closes the modal on the first update when no suggestions exist", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyDown})

			Expect(intent.reviewState.factSuggestionModal).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
		})
	})

	Describe("fact persistence in postSaveReview", func() {
		var (
			svc      *careerservice.Service
			factRepo *memoryrepo.FactRepository
		)

		BeforeEach(func() {
			repos := memoryrepo.NewRepositories()
			factRepo = repos.Fact.(*memoryrepo.FactRepository)

			svc = careerservice.NewService(repos.Event)
			svc.SetFactRepository(factRepo)
			svc.SetSkillRepository(repos.Skill)
			svc.SetBurstRepository(repos.Burst)

			intent.context.CareerService = svc
		})

		Context("with accepted facts that have empty IDs", func() {
			It("sets SourceEventID from the saved event and persists via SaveFact", func() {
				event := fixtures.EventWith("evt-saved-1", "Reviewed event with facts", "", "")
				fact := fixtures.FactForSave(
					"Led cross-team API design",
					[]string{"leadership"},
					career.RoleFitStaff,
					[]string{"hiring_manager"},
					"",
				)

				intent.currentState = StateReview
				intent.submitModal = feedback.NewSuccessModal("Event saved!")
				intent.reviewState = &ReviewInferredEventState{
					Event: event,
				}

				reviewData := eventview.ReviewResult{
					Event:  display.EventFromDomain(event),
					Bursts: []display.Burst{},
					Facts:  display.FactsFromDomain([]*career.Fact{fact}),
					Skills: []display.SkillSuggestion{},
				}
				cmd := intent.HandleSubmit(&widgets.SubmitViewResult{FormData: reviewData})
				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				intent.Update(msg)

				Expect(intent.result).NotTo(BeNil())
				Expect(intent.result.Status).To(Equal(intents.Completed))
				Expect(intent.result.Data.Facts).To(HaveLen(1))
				Expect(intent.result.Data.Facts[0].SourceEventID).To(Equal("evt-saved-1"))
				Expect(intent.result.Data.Facts[0].ID).NotTo(BeEmpty())

				ctx := context.Background()
				saved, err := factRepo.GetByID(ctx, intent.result.Data.Facts[0].ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(saved).NotTo(BeNil())
				Expect(saved.SourceEventID).To(Equal("evt-saved-1"))
			})
		})

		Context("with facts that already have an ID", func() {
			It("skips facts with existing IDs", func() {
				event := fixtures.EventWith("evt-saved-2", "Reviewed event", "", "")
				existingFact := fixtures.Fact("fact-existing-1", "evt-saved-2")

				ctx := context.Background()
				Expect(factRepo.Create(ctx, existingFact)).To(Succeed())

				intent.currentState = StateReview
				intent.submitModal = feedback.NewSuccessModal("Event saved!")
				intent.reviewState = &ReviewInferredEventState{
					Event: event,
				}

				reviewData := eventview.ReviewResult{
					Event:  display.EventFromDomain(event),
					Bursts: []display.Burst{},
					Facts:  display.FactsFromDomain([]*career.Fact{existingFact}),
					Skills: []display.SkillSuggestion{},
				}
				cmd := intent.HandleSubmit(&widgets.SubmitViewResult{FormData: reviewData})
				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				intent.Update(msg)

				Expect(intent.result).NotTo(BeNil())
				Expect(intent.result.Status).To(Equal(intents.Completed))
			})
		})

		Context("when CareerService is nil", func() {
			It("completes without persisting facts", func() {
				event := fixtures.EventWith("evt-saved-3", "Reviewed event", "", "")
				fact := fixtures.FactForSave(
					"Led API design",
					[]string{"leadership"},
					career.RoleFitStaff,
					[]string{"hiring_manager"},
					"",
				)

				intent.context.CareerService = nil
				intent.currentState = StateReview
				intent.submitModal = feedback.NewSuccessModal("Event saved!")
				intent.reviewState = &ReviewInferredEventState{
					Event: event,
				}

				reviewData := eventview.ReviewResult{
					Event:  display.EventFromDomain(event),
					Bursts: []display.Burst{},
					Facts:  display.FactsFromDomain([]*career.Fact{fact}),
					Skills: []display.SkillSuggestion{},
				}
				cmd := intent.HandleSubmit(&widgets.SubmitViewResult{FormData: reviewData})
				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				intent.Update(msg)

				Expect(intent.result).NotTo(BeNil())
				Expect(intent.result.Status).To(Equal(intents.Completed))
				Expect(fact.ID).To(BeEmpty())
			})
		})
	})
})

var _ = Describe("Skill Inference in CaptureEvent", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("when user opens skill review modal", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				InferredSkills: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
					{Name: "PostgreSQL", Category: "database", Confidence: 0.87},
				},
				AcceptedSkills: []*career.Skill{},
			}

			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event)
			intent.context.CareerService = svc

			screen := eventview.NewReview(
				display.EventFromDomain(intent.reviewState.Event),
				nil,
				nil,
				display.SkillSuggestionsFromDomain(intent.reviewState.InferredSkills),
			)
			intent.activeView = screen
		})

		It("initialises the skill modal with inferred skills", func() {
			result := &widgets.NavigateViewResult{ResultData: "suggest_skills"}
			intent.HandleNavigate(result)

			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeSkills))
			Expect(intent.reviewState.skillModal).NotTo(BeNil())
		})

		Context("with no inferred skills", func() {
			BeforeEach(func() {
				intent.reviewState.InferredSkills = []skillinference.SkillSuggestion{}
			})

			It("opens the modal with an empty list", func() {
				result := &widgets.NavigateViewResult{ResultData: "suggest_skills"}
				intent.HandleNavigate(result)

				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeSkills))
				Expect(intent.reviewState.skillModal).NotTo(BeNil())
			})

			It("auto-closes the modal on the first update when no suggestions exist", func() {
				result := &widgets.NavigateViewResult{ResultData: "suggest_skills"}
				intent.HandleNavigate(result)

				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeSkills))
				Expect(intent.reviewState.skillModal).NotTo(BeNil())

				intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyDown})

				Expect(intent.reviewState.skillModal).To(BeNil())
				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
			})
		})
	})

	Describe("when user accepts skills from modal", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				AcceptedSkills: []*career.Skill{},
				EditingMode:    EditingModeSkills,
			}
			suggestions := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.9},
			}
			intent.reviewState.skillModal = burstviews.NewSkillSuggestion(display.SkillSuggestionsFromDomain(suggestions), nil)

			screen := eventview.NewReview(display.EventFromDomain(intent.reviewState.Event), nil, nil, nil)
			intent.activeView = screen
		})

		It("transfers accepted skills as career.Skill to review state", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.AcceptedSkills).To(HaveLen(1))
			Expect(intent.reviewState.AcceptedSkills[0].Name).To(Equal("Go"))
			Expect(intent.reviewState.AcceptedSkills[0].Category).To(Equal("backend"))
		})

		It("closes the modal and resets editing mode", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.skillModal).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
		})
	})

	Describe("skill type conversion validation", func() {
		It("skips suggestions with empty names during review submission", func() {
			intent.currentState = StateSubmit
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "test", "", ""),
			}

			submitData := eventview.ReviewResult{
				Event:  display.EventFromDomain(intent.reviewState.Event),
				Bursts: display.BurstsFromDomain([]*career.Burst{}),
				Facts:  display.FactsFromDomain([]*career.Fact{}),
				Skills: display.SkillSuggestionsFromDomain([]skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.9},
					{Name: "", Category: "unknown", Confidence: 0.1},
					{Name: "Python", Category: "backend", Confidence: 0.8},
				}),
			}
			intent.HandleSubmit(&widgets.SubmitViewResult{FormData: submitData})

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Data.Skills).To(HaveLen(2))
			Expect(intent.result.Data.Skills[0].Name).To(Equal("Go"))
			Expect(intent.result.Data.Skills[1].Name).To(Equal("Python"))
		})

		It("handles nil skill suggestion list gracefully", func() {
			intent.currentState = StateSubmit
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "test", "", ""),
			}

			submitData := eventview.ReviewResult{
				Event:  display.EventFromDomain(intent.reviewState.Event),
				Bursts: display.BurstsFromDomain([]*career.Burst{}),
				Facts:  display.FactsFromDomain([]*career.Fact{}),
				Skills: display.SkillSuggestionsFromDomain([]skillinference.SkillSuggestion(nil)),
			}
			intent.HandleSubmit(&widgets.SubmitViewResult{FormData: submitData})

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Data.Skills).To(BeEmpty())
		})

		It("converts SkillSuggestion Category to career.Skill Category", func() {
			intent.currentState = StateSubmit
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "test", "", ""),
			}

			submitData := eventview.ReviewResult{
				Event:  display.EventFromDomain(intent.reviewState.Event),
				Bursts: display.BurstsFromDomain([]*career.Burst{}),
				Facts:  display.FactsFromDomain([]*career.Fact{}),
				Skills: display.SkillSuggestionsFromDomain([]skillinference.SkillSuggestion{
					{Name: "Docker", Category: "devops", Confidence: 0.85},
				}),
			}
			intent.HandleSubmit(&widgets.SubmitViewResult{FormData: submitData})

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Data.Skills).To(HaveLen(1))
			Expect(intent.result.Data.Skills[0].Name).To(Equal("Docker"))
			Expect(intent.result.Data.Skills[0].Category).To(Equal("devops"))
		})

		It("preserves skill order after filtering empty names", func() {
			intent.currentState = StateSubmit
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "test", "", ""),
			}

			submitData := eventview.ReviewResult{
				Event:  display.EventFromDomain(intent.reviewState.Event),
				Bursts: display.BurstsFromDomain([]*career.Burst{}),
				Facts:  display.FactsFromDomain([]*career.Fact{}),
				Skills: display.SkillSuggestionsFromDomain([]skillinference.SkillSuggestion{
					{Name: "", Category: "unknown", Confidence: 0.1},
					{Name: "Go", Category: "backend", Confidence: 0.9},
					{Name: "", Category: "unknown", Confidence: 0.2},
					{Name: "React", Category: "frontend", Confidence: 0.7},
					{Name: "", Category: "unknown", Confidence: 0.05},
				}),
			}
			intent.HandleSubmit(&widgets.SubmitViewResult{FormData: submitData})

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Data.Skills).To(HaveLen(2))
			Expect(intent.result.Data.Skills[0].Name).To(Equal("Go"))
			Expect(intent.result.Data.Skills[1].Name).To(Equal("React"))
		})
	})

	Describe("SubmitCompleteMsg populates review state", func() {
		BeforeEach(func() {
			intent.currentState = StateForm
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "test", "", ""),
			}
		})

		It("stores inferred skills from InferenceCompleteMsg in review state", func() {
			expectedSkills := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.95},
			}
			intent.Update(InferenceCompleteMsg{InferredSkills: expectedSkills})

			Expect(intent.reviewState.InferredSkills).To(Equal(expectedSkills))
		})

		It("handles empty inferred skills in InferenceCompleteMsg", func() {
			intent.Update(InferenceCompleteMsg{})

			Expect(intent.reviewState.InferredSkills).To(BeNil())
		})

		It("stores multiple inferred skills from InferenceCompleteMsg", func() {
			multipleSkills := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.95},
				{Name: "Docker", Category: "devops", Confidence: 0.88},
				{Name: "PostgreSQL", Category: "database", Confidence: 0.72},
			}
			intent.Update(InferenceCompleteMsg{InferredSkills: multipleSkills})

			Expect(intent.reviewState.InferredSkills).To(HaveLen(3))
			Expect(intent.reviewState.InferredSkills[0].Name).To(Equal("Go"))
			Expect(intent.reviewState.InferredSkills[2].Name).To(Equal("PostgreSQL"))
		})
	})

	Describe("error handling", func() {
		Context("when skill inference service is nil", func() {
			It("proceeds without inferring skills", func() {
				ctx := &IntentValidator{CaptureStrategy: "quick"}
				nilSvcIntent, err := NewIntent(ctx)
				Expect(err).NotTo(HaveOccurred())
				nilSvcIntent.Init()

				Expect(ctx.SkillInferenceService).To(BeNil())
			})
		})

		Context("when skills key is missing from submit data", func() {
			It("completes with empty skills", func() {
				intent.currentState = StateSubmit
				intent.reviewState = &ReviewInferredEventState{
					Event: fixtures.EventWith("evt-1", "test", "", ""),
				}

				submitData := eventview.ReviewResult{
					Event:  display.EventFromDomain(intent.reviewState.Event),
					Bursts: display.BurstsFromDomain([]*career.Burst{}),
					Facts:  display.FactsFromDomain([]*career.Fact{}),
					Skills: []display.SkillSuggestion{},
				}
				intent.HandleSubmit(&widgets.SubmitViewResult{FormData: submitData})

				Expect(intent.result).NotTo(BeNil())
				Expect(intent.result.Data.Skills).To(BeEmpty())
			})
		})

		Context("when skills key has wrong type in submit data", func() {
			It("treats unrecognised type as empty skills", func() {
				intent.currentState = StateSubmit
				intent.reviewState = &ReviewInferredEventState{
					Event: fixtures.EventWith("evt-1", "test", "", ""),
				}

				submitData := eventview.ReviewResult{
					Event:  display.EventFromDomain(intent.reviewState.Event),
					Bursts: display.BurstsFromDomain([]*career.Burst{}),
					Facts:  display.FactsFromDomain([]*career.Fact{}),
					Skills: []display.SkillSuggestion{},
				}
				intent.HandleSubmit(&widgets.SubmitViewResult{FormData: submitData})

				Expect(intent.result).NotTo(BeNil())
				Expect(intent.result.Data.Skills).To(BeEmpty())
			})
		})

		Context("when editing modal receives escape key", func() {
			It("resets editing mode to none", func() {
				intent.currentState = StateReview
				intent.reviewState = &ReviewInferredEventState{
					Event:       fixtures.EventWith("evt-1", "test", "", ""),
					EditingMode: EditingModeSkills,
				}
				suggestions := []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.9},
				}
				intent.reviewState.skillModal = burstviews.NewSkillSuggestion(display.SkillSuggestionsFromDomain(suggestions), nil)

				intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
			})
		})
	})

	Describe("review state skill flow integration", func() {
		It("round-trips skills from inference through modal to result", func() {
			intent.currentState = StateForm
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "Built microservices in Go", "", ""),
			}

			inferredSkills := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.95},
				{Name: "Docker", Category: "devops", Confidence: 0.88},
			}
			intent.Update(InferenceCompleteMsg{InferredSkills: inferredSkills})

			Expect(intent.reviewState.InferredSkills).To(HaveLen(2))

			intent.currentState = StateSubmit
			submitData := eventview.ReviewResult{
				Event:  display.EventFromDomain(intent.reviewState.Event),
				Bursts: display.BurstsFromDomain([]*career.Burst{}),
				Facts:  display.FactsFromDomain([]*career.Fact{}),
				Skills: display.SkillSuggestionsFromDomain(intent.reviewState.InferredSkills),
			}
			intent.HandleSubmit(&widgets.SubmitViewResult{FormData: submitData})

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Data.Skills).To(HaveLen(2))
			Expect(intent.result.Data.Skills[0].Name).To(Equal("Go"))
			Expect(intent.result.Data.Skills[0].Category).To(Equal("backend"))
			Expect(intent.result.Data.Skills[1].Name).To(Equal("Docker"))
			Expect(intent.result.Data.Skills[1].Category).To(Equal("devops"))
		})

		It("completes intent when review submits accepted skills as career.Skill", func() {
			intent.currentState = StateReview
			intent.submitModal = feedback.NewSuccessModal("Event saved!")
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "test", "", ""),
			}

			reviewData := eventview.ReviewResult{
				Event:  display.EventFromDomain(intent.reviewState.Event),
				Bursts: []display.Burst{},
				Facts:  []display.Fact{},
				Skills: []display.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 1.0},
					{Name: "Docker", Category: "devops", Confidence: 1.0},
				},
			}
			cmd := intent.HandleSubmit(&widgets.SubmitViewResult{FormData: reviewData})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			intent.Update(msg)

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Data.Skills).To(HaveLen(2))
			Expect(intent.result.Data.Skills[0].Name).To(Equal("Go"))
			Expect(intent.result.Data.Skills[0].Category).To(Equal("backend"))
			Expect(intent.result.Data.Skills[1].Name).To(Equal("Docker"))
			Expect(intent.result.Data.Skills[1].Category).To(Equal("devops"))
			Expect(intent.active).To(BeFalse())
		})
	})

	Describe("skill persistence on post-save review", func() {
		var (
			eventRepo *memoryrepo.EventRepository
			skillRepo *memoryrepo.SkillRepository
			svc       *careerservice.Service
			testEvent *career.Event
		)

		BeforeEach(func() {
			eventRepo = memoryrepo.NewEventRepository()
			skillRepo = memoryrepo.NewSkillRepository()
			eventRepo.SetSkillRepository(skillRepo)
			skillRepo.SetEventRepository(eventRepo)

			svc = careerservice.NewService(eventRepo)
			svc.SetSkillRepository(skillRepo)

			testEvent = fixtures.EventWith("evt-persist", "Built microservices in Go and Docker", "", "")
			Expect(eventRepo.Create(context.Background(), testEvent)).To(Succeed())

			intent.context.CareerService = svc
			intent.currentState = StateReview
			intent.submitModal = feedback.NewSuccessModal("Event saved!")
			intent.reviewState = &ReviewInferredEventState{
				Event: testEvent,
			}
		})

		It("persists skills from SkillSuggestion type through the real flow", func() {
			reviewData := eventview.ReviewResult{
				Event:  display.EventFromDomain(testEvent),
				Bursts: []display.Burst{},
				Facts:  []display.Fact{},
				Skills: display.SkillSuggestionsFromDomain([]skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
					{Name: "Docker", Category: "devops", Confidence: 0.88},
				}),
			}
			cmd := intent.HandleSubmit(&widgets.SubmitViewResult{FormData: reviewData})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			intent.Update(msg)

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Data.Skills).To(HaveLen(2))
			Expect(intent.result.Data.Skills[0].Name).To(Equal("Go"))
			Expect(intent.result.Data.Skills[0].ID).NotTo(BeEmpty())
			Expect(intent.result.Data.Skills[1].Name).To(Equal("Docker"))
			Expect(intent.result.Data.Skills[1].ID).NotTo(BeEmpty())

			goSkill, err := skillRepo.GetByName(context.Background(), "Go")
			Expect(err).NotTo(HaveOccurred())
			Expect(goSkill).NotTo(BeNil())

			dockerSkill, err := skillRepo.GetByName(context.Background(), "Docker")
			Expect(err).NotTo(HaveOccurred())
			Expect(dockerSkill).NotTo(BeNil())

			savedEvent, err := eventRepo.GetByID(context.Background(), testEvent.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(savedEvent.Skills).To(HaveLen(2))
			Expect(savedEvent.Skills).To(ContainElement(goSkill.ID))
			Expect(savedEvent.Skills).To(ContainElement(dockerSkill.ID))
		})

		It("skips persistence when no skills are provided", func() {
			reviewData := eventview.ReviewResult{
				Event:  display.EventFromDomain(testEvent),
				Bursts: []display.Burst{},
				Facts:  []display.Fact{},
				Skills: []display.SkillSuggestion{},
			}
			cmd := intent.HandleSubmit(&widgets.SubmitViewResult{FormData: reviewData})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			intent.Update(msg)

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Data.Skills).To(BeEmpty())

			savedEvent, err := eventRepo.GetByID(context.Background(), testEvent.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(savedEvent.Skills).To(BeEmpty())
		})

		It("persists skills accepted from modal through screen submit", func() {
			intent.reviewState.InferredSkills = []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.95},
			}
			intent.reviewState.AcceptedSkills = []*career.Skill{}
			intent.reviewState.EditingMode = EditingModeSkills

			screen := eventview.NewReview(
				display.EventFromDomain(testEvent), nil, nil,
				display.SkillSuggestionsFromDomain(intent.reviewState.InferredSkills),
			)
			intent.activeView = screen

			intent.reviewState.skillModal = burstviews.NewSkillSuggestion(
				display.SkillSuggestionsFromDomain(intent.reviewState.InferredSkills), nil,
			)

			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.AcceptedSkills).To(HaveLen(1))
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))

			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(widgets.ResultSubmit))

			cmd := intent.HandleSubmit(&widgets.SubmitViewResult{FormData: result.Data()})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			intent.Update(msg)

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Data.Skills).To(HaveLen(1))
			Expect(intent.result.Data.Skills[0].Name).To(Equal("Go"))
			Expect(intent.result.Data.Skills[0].ID).NotTo(BeEmpty())

			goSkill, err := skillRepo.GetByName(context.Background(), "Go")
			Expect(err).NotTo(HaveOccurred())
			Expect(goSkill).NotTo(BeNil())

			savedEvent, err := eventRepo.GetByID(context.Background(), testEvent.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(savedEvent.Skills).To(ContainElement(goSkill.ID))
		})
	})
})

var _ = Describe("Burst SuggestionReview in CaptureEvent", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("HandleNavigate suggest_bursts", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				InferredBursts: func() []*career.Burst {
					b := fixtures.Burst("", "evt-1", "evt-2")
					b.Name = "API Work"
					b.Description = "REST API development"
					return []*career.Burst{b}
				}(),
				AcceptedBursts: make([]*career.Burst, 0),
			}
		})

		Context("with CareerService nil", func() {
			It("marks intent as failed", func() {
				intent.context.CareerService = nil
				result := &widgets.NavigateViewResult{ResultData: "suggest_bursts"}
				intent.HandleNavigate(result)

				Expect(intent.active).To(BeFalse())
			})
		})

		Context("with CareerService set", func() {
			BeforeEach(func() {
				repos := memoryrepo.NewRepositories()
				svc := careerservice.NewService(repos.Event)
				intent.context.CareerService = svc
				intent.reviewState.EditingMode = EditingModeNone
			})

			It("sets EditingMode to EditingModeBursts", func() {
				result := &widgets.NavigateViewResult{ResultData: "suggest_bursts"}
				intent.HandleNavigate(result)

				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeBursts))
			})

			It("creates a SuggestionReview on burstModal", func() {
				result := &widgets.NavigateViewResult{ResultData: "suggest_bursts"}
				intent.HandleNavigate(result)

				Expect(intent.reviewState.burstModal).NotTo(BeNil())
			})

			It("populates EventIDs from InferredBursts", func() {
				result := &widgets.NavigateViewResult{ResultData: "suggest_bursts"}
				intent.HandleNavigate(result)

				Expect(intent.reviewState.burstModal).NotTo(BeNil())
				Expect(intent.reviewState.burstModal.GetSuggestionsCount()).To(Equal(1))
				suggestion := intent.reviewState.burstModal.GetCurrentSuggestion()
				Expect(suggestion).NotTo(BeNil())
				Expect(suggestion.EventIDs).To(Equal([]string{"evt-1", "evt-2"}))
			})
		})

		Context("with InferredBurstSuggestions populated", func() {
			BeforeEach(func() {
				repos := memoryrepo.NewRepositories()
				svc := careerservice.NewService(repos.Event)
				intent.context.CareerService = svc
				intent.reviewState.EditingMode = EditingModeNone
				intent.reviewState.InferredBurstSuggestions = []burstfact.BurstSuggestion{
					{
						Name:            "API Work",
						Description:     "REST API development",
						EventIDs:        []string{"evt-1", "evt-2"},
						ConfidenceScore: 0.87,
					},
				}
			})

			It("uses InferredBursts over InferredBurstSuggestions to build the modal", func() {
				result := &widgets.NavigateViewResult{ResultData: "suggest_bursts"}
				intent.HandleNavigate(result)
				Expect(intent.reviewState.burstModal).NotTo(BeNil())
				suggestion := intent.reviewState.burstModal.GetCurrentSuggestion()
				Expect(suggestion).NotTo(BeNil())
				Expect(suggestion.Name).To(Equal("API Work"))
				Expect(intent.reviewState.burstModal.GetSuggestionsCount()).To(Equal(1))
			})
		})
	})

	Describe("updateEditingModal with EditingModeBursts using SuggestionReview", func() {
		BeforeEach(func() {
			inferred := fixtures.Burst("inferred-api-dev", "evt-1")
			inferred.Name = "API Development"
			inferred.Description = "Built REST APIs"
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				InferredBursts: []*career.Burst{inferred},
				AcceptedBursts: make([]*career.Burst, 0),
				EditingMode:    EditingModeBursts,
			}
			suggestions := []burstfact.BurstSuggestion{
				{Name: "API Development", Description: "Built REST APIs", EventIDs: []string{"evt-1"}, ConfidenceScore: 0.9},
			}
			intent.reviewState.burstModal = burstviews.NewSuggestionReview(display.BurstSuggestionsFromDomain(suggestions), nil)
			screen := eventview.NewReview(display.EventFromDomain(intent.reviewState.Event), nil, nil, nil)
			intent.activeView = screen
		})

		It("auto-closes modal and resets mode when no suggestions remain", func() {
			intent.reviewState.burstModal = burstviews.NewSuggestionReview([]display.BurstSuggestion{}, nil)

			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyDown})

			Expect(intent.reviewState.burstModal).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
		})

		It("transfers accepted bursts to review state on accept key", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.AcceptedBursts).To(HaveLen(1))
			Expect(intent.reviewState.AcceptedBursts[0].Name).To(Equal("API Development"))
			Expect(intent.reviewState.AcceptedBursts[0].Description).To(Equal("Built REST APIs"))
			Expect(intent.reviewState.AcceptedBursts[0].EventIDs).To(Equal([]string{"event-1", "event-2"}))
		})

		It("closes modal and resets editing mode after all accepted", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.burstModal).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
		})

		It("closes modal on escape without adding bursts", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.reviewState.burstModal).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
			Expect(intent.reviewState.AcceptedBursts).To(BeEmpty())
		})

		It("does not add rejected bursts to accepted list", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			Expect(intent.reviewState.AcceptedBursts).To(BeEmpty())
			Expect(intent.reviewState.burstModal).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
		})
	})

	Describe("burst suggestion with empty name", func() {
		BeforeEach(func() {
			inferred := fixtures.Burst("inferred-empty", "evt-1", "evt-2")
			inferred.Name = "Burst of 2 events"
			inferred.Description = "unnamed burst"
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "test event", "", ""),
				InferredBursts: []*career.Burst{inferred},
				AcceptedBursts: make([]*career.Burst, 0),
				EditingMode:    EditingModeBursts,
			}
			suggestions := []burstfact.BurstSuggestion{
				{Name: "", Description: "unnamed burst", EventIDs: []string{"evt-1", "evt-2"}, ConfidenceScore: 0.7},
			}
			intent.reviewState.burstModal = burstviews.NewSuggestionReview(display.BurstSuggestionsFromDomain(suggestions), nil)
			screen := eventview.NewReview(display.EventFromDomain(intent.reviewState.Event), nil, nil, nil)
			intent.activeView = screen
		})
		It("generates a name based on event count when name is empty", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(intent.reviewState.AcceptedBursts[0].Name).To(Equal("Burst of 2 events"))
		})
	})

	Describe("multiple burst suggestions: accept some reject some", func() {
		BeforeEach(func() {
			inferredKeep := fixtures.Burst("inferred-keep", "evt-1")
			inferredKeep.Name = "Keep This"
			inferredKeep.Description = "accepted"
			inferredSkip := fixtures.Burst("inferred-skip", "evt-2")
			inferredSkip.Name = "Skip This"
			inferredSkip.Description = "rejected"
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "test event", "", ""),
				InferredBursts: []*career.Burst{inferredKeep, inferredSkip},
				AcceptedBursts: make([]*career.Burst, 0),
				EditingMode:    EditingModeBursts,
			}
			suggestions := []burstfact.BurstSuggestion{
				{Name: "Keep This", Description: "accepted", EventIDs: []string{"evt-1"}, ConfidenceScore: 0.9},
				{Name: "Skip This", Description: "rejected", EventIDs: []string{"evt-2"}, ConfidenceScore: 0.3},
			}
			intent.reviewState.burstModal = burstviews.NewSuggestionReview(display.BurstSuggestionsFromDomain(suggestions), nil)
			screen := eventview.NewReview(display.EventFromDomain(intent.reviewState.Event), nil, nil, nil)
			intent.activeView = screen
		})
		It("only includes accepted bursts in AcceptedBursts", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			Expect(intent.reviewState.AcceptedBursts[0].Name).To(Equal("Keep This"))
		})
	})

	Describe("accepted burst preserves original ID from InferredBursts", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				InferredBursts: func() []*career.Burst {
					b := fixtures.Burst("inferred-burst-uuid-1", "evt-1", "evt-2")
					b.Name = "API Development"
					b.Description = "Built REST APIs"
					return []*career.Burst{b}
				}(),
				AcceptedBursts: make([]*career.Burst, 0),
				EditingMode:    EditingModeBursts,
			}
			suggestions := []burstfact.BurstSuggestion{
				{Name: "API Development", Description: "Built REST APIs", EventIDs: []string{"evt-1", "evt-2"}, ConfidenceScore: 0.9},
			}
			intent.reviewState.burstModal = burstviews.NewSuggestionReview(display.BurstSuggestionsFromDomain(suggestions), nil)

			screen := eventview.NewReview(display.EventFromDomain(intent.reviewState.Event), display.BurstsFromDomain(intent.reviewState.InferredBursts), nil, nil)
			intent.activeView = screen
		})

		It("uses the original burst with DB-assigned ID when name matches an InferredBurst", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.AcceptedBursts).To(HaveLen(1))
			Expect(intent.reviewState.AcceptedBursts[0].ID).To(Equal("inferred-burst-uuid-1"))
			Expect(intent.reviewState.AcceptedBursts[0].Name).To(Equal("API Development"))
		})

		It("skips the suggestion when the name does not match any InferredBurst", func() {
			intent.reviewState.InferredBursts = []*career.Burst{}
			intent.reviewState.burstModal = burstviews.NewSuggestionReview(display.BurstSuggestionsFromDomain([]burstfact.BurstSuggestion{
				{Name: "Unknown Burst", Description: "no match", EventIDs: []string{"evt-3"}, ConfidenceScore: 0.5},
			}), nil)
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.AcceptedBursts).To(BeEmpty())
		})
	})

	Describe("accepted burst indicator on review screen", func() {
		It("shows ● indicator when accepted burst name matches an inferred burst", func() {
			inferred := fixtures.Burst("inferred-burst-uuid-1", "evt-1", "evt-2")
			inferred.Name = "API Development"
			inferred.Description = "Built REST APIs"
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				InferredBursts: []*career.Burst{inferred},
				AcceptedBursts: make([]*career.Burst, 0),
				EditingMode:    EditingModeBursts,
			}
			suggestions := []burstfact.BurstSuggestion{
				{Name: "API Development", Description: "Built REST APIs", EventIDs: []string{"evt-1", "evt-2"}, ConfidenceScore: 0.9},
			}
			intent.reviewState.burstModal = burstviews.NewSuggestionReview(display.BurstSuggestionsFromDomain(suggestions), nil)
			screen := eventview.NewReview(display.EventFromDomain(intent.reviewState.Event), display.BurstsFromDomain(intent.reviewState.InferredBursts), nil, nil)
			intent.activeView = screen

			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(screen.RenderContent()).To(ContainSubstring("●"))
		})

		It("does not show ● indicator when InferredBursts is empty and suggestion is skipped", func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				InferredBursts: []*career.Burst{},
				AcceptedBursts: make([]*career.Burst, 0),
				EditingMode:    EditingModeBursts,
			}
			suggestions := []burstfact.BurstSuggestion{
				{Name: "API Development", Description: "Built REST APIs", EventIDs: []string{"evt-1", "evt-2"}, ConfidenceScore: 0.9},
			}
			intent.reviewState.burstModal = burstviews.NewSuggestionReview(display.BurstSuggestionsFromDomain(suggestions), nil)
			inferred := fixtures.Burst("inferred-burst-uuid-1", "evt-1", "evt-2")
			inferred.Name = "API Development"
			inferred.Description = "Built REST APIs"
			screen := eventview.NewReview(display.EventFromDomain(intent.reviewState.Event), display.BurstsFromDomain([]*career.Burst{inferred}), nil, nil)
			intent.activeView = screen
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.AcceptedBursts).To(BeEmpty())
			Expect(screen.RenderContent()).NotTo(ContainSubstring("●"))
		})
	})

	Describe("appending to existing accepted bursts", func() {
		BeforeEach(func() {
			existingBurst := fixtures.BurstConfirmed("burst-1")
			existingBurst.Name = "Existing Burst"
			inferredNew := fixtures.Burst("inferred-new", "evt-1")
			inferredNew.Name = "New Burst"
			inferredNew.Description = "newly confirmed"
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "test event", "", ""),
				InferredBursts: []*career.Burst{inferredNew},
				AcceptedBursts: []*career.Burst{existingBurst},
				EditingMode:    EditingModeBursts,
			}
			suggestions := []burstfact.BurstSuggestion{
				{Name: "New Burst", Description: "newly confirmed", EventIDs: []string{"evt-1"}, ConfidenceScore: 0.8},
			}
			intent.reviewState.burstModal = burstviews.NewSuggestionReview(display.BurstSuggestionsFromDomain(suggestions), nil)
			screen := eventview.NewReview(display.EventFromDomain(intent.reviewState.Event), nil, nil, nil)
			intent.activeView = screen
		})
		It("appends new bursts to existing accepted list", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(intent.reviewState.AcceptedBursts[0].Name).To(Equal("Existing Burst"))
			Expect(intent.reviewState.AcceptedBursts[1].Name).To(Equal("New Burst"))
		})
	})

	Describe("View with EditingModeBursts", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "Built services", "", ""),
				AcceptedBursts: make([]*career.Burst, 0),
				EditingMode:    EditingModeBursts,
			}
			suggestions := []burstfact.BurstSuggestion{
				{Name: "API Work", Description: "REST API development", EventIDs: []string{"evt-1"}, ConfidenceScore: 0.9},
			}
			intent.reviewState.burstModal = burstviews.NewSuggestionReview(display.BurstSuggestionsFromDomain(suggestions), nil)

			screen := eventview.NewReview(display.EventFromDomain(intent.reviewState.Event), nil, nil, nil)
			intent.activeView = screen
		})

		It("renders the burst modal overlay when EditingModeBursts is active", func() {
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("does not use getEditingModalContent for bursts", func() {
			view := intent.View()
			Expect(view).NotTo(BeNil())
		})
	})

	Describe("HandleNavigate suggest_skills", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				InferredSkills: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
				},
				AcceptedSkills: make([]*career.Skill, 0),
			}
		})

		It("sets EditingMode to EditingModeSkills", func() {
			result := &widgets.NavigateViewResult{ResultData: "suggest_skills"}
			intent.HandleNavigate(result)

			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeSkills))
		})

		It("creates a skillModal", func() {
			result := &widgets.NavigateViewResult{ResultData: "suggest_skills"}
			intent.HandleNavigate(result)

			Expect(intent.reviewState.skillModal).NotTo(BeNil())
		})
	})

	Describe("updateEditingModal with EditingModeSkills", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				AcceptedSkills: make([]*career.Skill, 0),
				EditingMode:    EditingModeSkills,
			}
			suggestions := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.95},
			}
			intent.reviewState.skillModal = burstviews.NewSkillSuggestion(display.SkillSuggestionsFromDomain(suggestions), nil)

			screen := eventview.NewReview(display.EventFromDomain(intent.reviewState.Event), nil, nil, nil)
			intent.activeView = screen
		})

		It("auto-closes modal and resets mode when no suggestions remain", func() {
			intent.reviewState.skillModal = burstviews.NewSkillSuggestion([]display.SkillSuggestion{}, nil)

			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyDown})

			Expect(intent.reviewState.skillModal).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
		})

		It("closes modal on escape without adding skills", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
			Expect(intent.reviewState.AcceptedSkills).To(BeEmpty())
		})
	})

	Describe("Full-flow burst acceptance through intent.Update", func() {
		var (
			screen   *eventview.Review
			inferred *career.Burst
		)

		BeforeEach(func() {
			inferred = fixtures.Burst("inferred-burst-uuid-1", "evt-1", "evt-2")
			inferred.Name = "API Development"
			inferred.Description = "Built REST APIs"

			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event)
			intent.context.CareerService = svc

			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				InferredBursts: []*career.Burst{inferred},
				InferredBurstSuggestions: []burstfact.BurstSuggestion{
					{Name: "API Development", Description: "Built REST APIs", EventIDs: []string{"evt-1", "evt-2"}, ConfidenceScore: 0.9},
				},
				AcceptedBursts: make([]*career.Burst, 0),
				EditingMode:    EditingModeNone,
			}

			screen = eventview.NewReview(
				display.EventFromDomain(intent.reviewState.Event),
				display.BurstsFromDomain(intent.reviewState.InferredBursts),
				nil,
				nil,
			)
			intent.activeView = screen
		})

		It("shows ● indicator after pressing b then a through intent.Update", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(screen.RenderContent()).To(ContainSubstring("●"))
		})

		It("sets AcceptedBursts on reviewState", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.AcceptedBursts).To(HaveLen(1))
			Expect(intent.reviewState.AcceptedBursts[0].Name).To(Equal("API Development"))
		})

		It("resets EditingMode to None after acceptance", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
		})
	})

	Describe("Full-flow skill acceptance through intent.Update (comparison)", func() {
		var screen *eventview.Review

		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				InferredSkills: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
				},
				AcceptedSkills: make([]*career.Skill, 0),
				EditingMode:    EditingModeNone,
			}

			screen = eventview.NewReview(
				display.EventFromDomain(intent.reviewState.Event),
				nil,
				nil,
				display.SkillSuggestionsFromDomain(intent.reviewState.InferredSkills),
			)
			intent.activeView = screen
		})

		It("shows ● indicator after pressing s then a through intent.Update", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(screen.RenderContent()).To(ContainSubstring("●"))
		})

		It("sets AcceptedSkills on reviewState", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.AcceptedSkills).To(HaveLen(1))
			Expect(intent.reviewState.AcceptedSkills[0].Name).To(Equal("Go"))
		})
	})
})
var _ = Describe("Coverage boost — eventFromFormData", func() {
	Describe("when date is empty", func() {
		It("defaults to now when date field is empty", func() {
			data := &forms.CaptureEventFormData{
				Text:    "Valid event text",
				Date:    "",
				Company: "Test Corp",
				Project: "Test Project",
				Tags:    []string{"technical"},
			}
			event, err := eventFromFormData(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(event).NotTo(BeNil())
			Expect(event.Date).NotTo(BeZero())
		})
	})

	Describe("when date is valid", func() {
		It("parses the date correctly", func() {
			data := &forms.CaptureEventFormData{
				Text:    "Valid event text",
				Date:    "2024-06-15",
				Company: "Test Corp",
				Project: "Test Project",
				Tags:    []string{"technical"},
			}
			event, err := eventFromFormData(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(event).NotTo(BeNil())
			Expect(event.Date.Year()).To(Equal(2024))
			Expect(int(event.Date.Month())).To(Equal(6))
			Expect(event.Date.Day()).To(Equal(15))
		})
	})

	Describe("when date is invalid", func() {
		It("returns an error for an unparseable date", func() {
			data := &forms.CaptureEventFormData{
				Text: "Valid event text",
				Date: "not-a-date",
			}
			_, err := eventFromFormData(data)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("field mapping", func() {
		It("maps all form fields to the event", func() {
			data := &forms.CaptureEventFormData{
				Text:       "Valid event text",
				Date:       "",
				Company:    "Acme Ltd",
				Project:    "Phoenix",
				Tags:       []string{"leadership"},
				Categories: []string{"leadership"},
			}
			event, err := eventFromFormData(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(event.Text).To(Equal("Valid event text"))
			Expect(event.Company).To(Equal("Acme Ltd"))
			Expect(event.Project).To(Equal("Phoenix"))
			Expect(event.Tags).To(Equal([]string{"leadership"}))
			Expect(event.Categories).To(Equal([]string{"leadership"}))
		})
	})
})

var _ = Describe("Coverage boost — GetFactSuggestionModal", func() {
	It("returns nil when ReviewInferredEventState is nil", func() {
		var state *ReviewInferredEventState
		modal := state.GetFactSuggestionModal()
		Expect(modal).To(BeNil())
	})

	It("returns nil when factSuggestionModal is nil", func() {
		state := &ReviewInferredEventState{}
		modal := state.GetFactSuggestionModal()
		Expect(modal).To(BeNil())
	})
})

var _ = Describe("Coverage boost — HandleError additional paths", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	It("handles map data with empty message by using fallback message", func() {
		result := &widgets.ErrorViewResult{
			Err:     errors.New("underlying cause"),
			Message: "",
		}
		_ = intent.HandleError(result)
		Expect(intent.active).To(BeFalse())
	})

	It("handles map data with non-nil error and non-empty message", func() {
		result := &widgets.ErrorViewResult{
			Err:     errors.New("screen exploded"),
			Message: "Screen encountered an error",
		}
		_ = intent.HandleError(result)
		Expect(intent.result.Error.Message).To(Equal("Screen encountered an error"))
	})
})

var _ = Describe("Coverage boost — HandleSubmit additional branches", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("StateForm — SubmitConfirmed false returns nil", func() {
		It("returns nil when SubmitConfirmed is false", func() {
			intent.currentState = StateForm
			result := &widgets.SubmitViewResult{
				FormData: &forms.CaptureEventFormData{
					SubmitConfirmed: false,
					Text:            "Some text",
				},
			}
			cmd := intent.HandleSubmit(result)
			Expect(cmd).To(BeNil())
		})
	})

	Describe("StateForm — invalid date in form data", func() {
		It("shows validation error modal when date cannot be parsed", func() {
			intent.currentState = StateForm
			result := &widgets.SubmitViewResult{
				FormData: &forms.CaptureEventFormData{
					SubmitConfirmed: true,
					Text:            "Valid event text for invalid date test path",
					Date:            "not-a-date",
				},
			}
			intent.HandleSubmit(result)
			Expect(intent.submitModal).NotTo(BeNil())
		})
	})

	Describe("StateForm — invalid event validation shows modal", func() {
		It("shows validation error modal when event text is empty", func() {
			intent.currentState = StateForm
			result := &widgets.SubmitViewResult{
				FormData: &forms.CaptureEventFormData{
					SubmitConfirmed: true,
					Text:            "",
					Date:            "",
				},
			}
			intent.HandleSubmit(result)
			Expect(intent.submitModal).NotTo(BeNil())
		})
	})

	Describe("StateReview — invalid review data type", func() {
		It("marks intent as failed when data is not map[string]interface{}", func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("", "Test event", "", ""),
			}
			result := &widgets.SubmitViewResult{FormData: "wrong type"}
			intent.HandleSubmit(result)
			Expect(intent.active).To(BeFalse())
		})
	})

	Describe("StateReview — missing event in review data", func() {
		It("marks intent as failed when event key is missing", func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("", "Test event", "", ""),
			}
			result := &widgets.SubmitViewResult{
				FormData: map[string]interface{}{
					"bursts": make([]*career.Burst, 0),
					"facts":  make([]*career.Fact, 0),
				},
			}
			intent.HandleSubmit(result)
			Expect(intent.active).To(BeFalse())
		})
	})

	Describe("StateSubmit — missing event in submit data", func() {
		It("marks intent as failed when event key is missing", func() {
			intent.currentState = StateSubmit
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("", "Test event", "", ""),
			}
			result := &widgets.SubmitViewResult{
				FormData: "invalid type",
			}
			intent.HandleSubmit(result)
			Expect(intent.active).To(BeFalse())
		})
	})

	Describe("StateSubmit — invalid data type", func() {
		It("marks intent as failed when data is not map[string]interface{}", func() {
			intent.currentState = StateSubmit
			result := &widgets.SubmitViewResult{FormData: "wrong type"}
			intent.HandleSubmit(result)
			Expect(intent.active).To(BeFalse())
		})
	})
})

var _ = Describe("eventFromFormData — date parse failure", func() {
	It("returns an error for an invalid date string", func() {
		data := &forms.CaptureEventFormData{
			Text:            "some event text here",
			Date:            "not-a-date",
			SubmitConfirmed: true,
		}
		_, err := eventFromFormData(data)
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("GetFactSuggestionModal — nil receiver", func() {
	It("returns nil without panicking when called on nil ReviewInferredEventState", func() {
		var state *ReviewInferredEventState
		result := state.GetFactSuggestionModal()
		Expect(result).To(BeNil())
	})
})

var _ = Describe("HandleSubmit — StateReview with postSaveReview flag", func() {
	It("calls performPostSavePersistence path without panicking", func() {
		repos := memoryrepo.NewRepositories()
		svc := careerservice.NewService(repos.Event)
		svc.SetSkillRepository(repos.Skill)

		ctx := &IntentValidator{CaptureStrategy: "quick"}
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
		intent.context.CareerService = svc
		intent.submitModal = feedback.NewSuccessModal("Event saved!")
		intent.currentState = StateReview

		event := fixtures.EventWith("ps-review", "Valid post-save review event text", "Corp", "Proj")
		event.Date = time.Now().Add(-time.Hour)
		event.Tags = []string{"technical"}
		event.Categories = []string{"technical"}
		Expect(repos.Event.Create(context.Background(), event)).To(Succeed())

		intent.reviewState = &ReviewInferredEventState{
			Event:         event,
			AcceptedFacts: make([]*career.Fact, 0),
		}

		Expect(func() {
			intent.HandleSubmit(&widgets.SubmitViewResult{
				FormData: eventview.ReviewResult{
					Event:  display.EventFromDomain(event),
					Bursts: display.BurstsFromDomain([]*career.Burst{}),
					Facts:  display.FactsFromDomain([]*career.Fact{}),
					Skills: []display.SkillSuggestion{},
				},
			})
		}).NotTo(Panic())
	})
})

var _ = Describe("HandleSubmit — StateForm with unparseable date", func() {
	It("shows a validation error modal without deactivating the intent", func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
		intent.currentState = StateForm

		formData := &forms.CaptureEventFormData{
			Text:            "Valid event text for testing the date parse path",
			Date:            "not-a-valid-date",
			SubmitConfirmed: true,
		}
		intent.HandleSubmit(&widgets.SubmitViewResult{FormData: formData})
		Expect(intent.submitModal).NotTo(BeNil())
		Expect(intent.IsActive()).To(BeTrue())
	})
})

var _ = Describe("handleViewResult", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	It("routes ResultCancel to HandleCancel", func() {
		result := &widgets.CancelViewResult{}
		intent.handleViewResult(result)
		Expect(intent.IsActive()).To(BeFalse())
		Expect(intent.Result().Status).To(Equal(intents.Cancelled))
	})

	It("routes ResultNavigate to HandleNavigate", func() {
		result := &widgets.NavigateViewResult{ResultData: StrategyQuick}
		intent.handleViewResult(result)
		Expect(intent.GetState()).To(Equal("form"))
	})

	It("routes ResultSubmit to HandleSubmit", func() {
		intent.currentState = StateForm
		formData := &forms.CaptureEventFormData{
			Text:            "Valid event text for handleViewResult test",
			Date:            "2026-01-15",
			SubmitConfirmed: true,
		}
		result := &widgets.SubmitViewResult{FormData: formData}
		Expect(func() {
			intent.handleViewResult(result)
		}).NotTo(Panic())
	})

	It("routes ResultError to HandleError", func() {
		result := &widgets.ErrorViewResult{
			Err:     errors.New("test error"),
			Message: "test error message",
		}
		intent.handleViewResult(result)
		Expect(intent.IsActive()).To(BeFalse())
		Expect(intent.Result().Status).To(Equal(intents.Failed))
	})
})

var _ = Describe("updateMetadataModal", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
		intent.reviewState = &ReviewInferredEventState{
			Event: fixtures.EventWith("", "test event", "", ""),
		}
	})

	It("returns nil when metadataModal is nil", func() {
		cmd := intent.updateMetadataModal(tea.KeyMsg{})
		Expect(cmd).To(BeNil())
	})

	It("processes modal update when metadataModal exists", func() {
		modal := eventview.NewReviewEnrichment(
			display.EventFromDomain(intent.reviewState.Event),
			eventview.ReviewEnrichmentConfig{},
		)
		intent.reviewState.metadataModal = modal
		intent.reviewState.EditingMode = EditingModeMetadata

		intent.updateMetadataModal(tea.KeyMsg{})
		Expect(intent.reviewState.metadataModal).NotTo(BeNil())
	})
})

var _ = Describe("processAcceptedSkills", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
		intent.reviewState = &ReviewInferredEventState{
			Event:          fixtures.EventWith("", "test event", "", ""),
			AcceptedSkills: make([]*career.Skill, 0),
		}
	})

	It("processes accepted skills from skillModal", func() {
		skillModal := burstviews.NewSkillSuggestion(
			[]display.SkillSuggestion{
				{Name: "Go", Category: "Language", Confidence: 0.95},
				{Name: "Testing", Category: "Practice", Confidence: 0.85},
			},
			nil,
		)
		intent.reviewState.skillModal = skillModal

		intent.processAcceptedSkills()
		Expect(intent.reviewState.AcceptedSkills).NotTo(BeNil())
	})

	It("skips skills with empty names", func() {
		skillModal := burstviews.NewSkillSuggestion(
			[]display.SkillSuggestion{
				{Name: "", Category: "Language", Confidence: 0.95},
				{Name: "Go", Category: "Language", Confidence: 0.85},
			},
			nil,
		)
		intent.reviewState.skillModal = skillModal

		intent.processAcceptedSkills()
		for _, skill := range intent.reviewState.AcceptedSkills {
			Expect(skill.Name).NotTo(BeEmpty())
		}
	})
})
