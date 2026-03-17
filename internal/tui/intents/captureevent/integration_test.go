package captureevent

import (
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	memoryrepo "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/tui/intents"
	"github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// CaptureEvent Integration Tests verify the bridge between Intent and Views.
//
// These tests ensure that ViewResults are correctly routed through the Intent's
// handlers, that state transitions work as expected, and that the full workflow
// from strategy selection through submission functions correctly.
var _ = Describe("CaptureEvent Integration Tests", func() {
	var (
		intent          *Intent
		svc             *careerservice.Service
		ctx             *IntentValidator
		cliEventService *service.CLIEventService
	)

	BeforeEach(func() {
		repos := memoryrepo.NewRepositories()
		svc = careerservice.NewService(repos.Event)
		svc.SetSkillRepository(repos.Skill)
		svc.SetBurstRepository(repos.Burst)
		svc.SetFactRepository(repos.Fact)

		cliEventService = service.NewCLIEventService(svc)

		ctx = &IntentValidator{
			CaptureStrategy: "quick",
			CLIEventService: cliEventService,
			CareerService:   svc,
		}

		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(intent).NotTo(BeNil())

		intent.Init()
	})

	Describe("handleViewResult routing", func() {
		It("routes NavigateViewResult to HandleNavigate", func() {
			result := &widgets.NavigateViewResult{
				ResultData: CaptureStrategy("quick"),
			}

			_ = intent.handleViewResult(result)
			Expect(intent.currentState).To(Equal(StateForm))
			Expect(intent.strategy).To(Equal(CaptureStrategy("quick")))
		})

		It("routes SubmitViewResult to HandleSubmit when in StateForm", func() {
			intent.currentState = StateForm
			intent.strategy = "quick"

			evt := fixtures.EventWith("", "Integration test event description here", "TestCo", "TestProj")
			formData := forms.GetCaptureEventFormData(evt)
			formData.SubmitConfirmed = true

			result := &widgets.SubmitViewResult{
				FormData: formData,
			}

			_ = intent.handleViewResult(result)
			Expect(intent.reviewState).NotTo(BeNil())
			Expect(intent.reviewState.Event).NotTo(BeNil())
		})

		It("routes CancelViewResult to HandleCancel from StateChooseStrategy", func() {
			result := &widgets.CancelViewResult{}

			_ = intent.handleViewResult(result)
			Expect(intent.result.Status).To(Equal(intents.Cancelled))
		})

		It("routes ErrorViewResult to HandleError", func() {
			result := &widgets.ErrorViewResult{
				Err:     nil,
				Message: "Test error message",
			}

			_ = intent.handleViewResult(result)
			Expect(intent.result.Status).To(Equal(intents.Failed))
		})
	})

	Describe("Full workflow: StateChooseStrategy → StateForm → StateReview → StateSubmit", func() {
		It("completes strategy selection and transitions to form", func() {
			Expect(intent.currentState).To(Equal(StateChooseStrategy))
			Expect(intent.activeView).NotTo(BeNil())
			_, ok := intent.activeView.(*event.StrategySelect)
			Expect(ok).To(BeTrue())

			result := &widgets.NavigateViewResult{
				ResultData: CaptureStrategy("quick"),
			}
			_ = intent.handleViewResult(result)

			Expect(intent.currentState).To(Equal(StateForm))
			Expect(intent.strategy).To(Equal(CaptureStrategy("quick")))
		})

		It("transitions from form to review after submission", func() {
			intent.currentState = StateForm
			intent.strategy = "quick"

			Expect(intent.activeView).NotTo(BeNil())

			evt := fixtures.EventWith("", "Test event for submission details", "TestCorp", "ProjectA")
			formData := forms.GetCaptureEventFormData(evt)
			formData.SubmitConfirmed = true

			result := &widgets.SubmitViewResult{
				FormData: formData,
			}

			_ = intent.handleViewResult(result)

			Expect(intent.reviewState).NotTo(BeNil())
			Expect(intent.reviewState.Event).NotTo(BeNil())
			Expect(intent.reviewState.Event.Text).To(Equal("Test event for submission details"))
		})

		It("handles submit completion and transitions to review", func() {
			intent.currentState = StateSubmit
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("", "Reviewed event", "TestCorp", "ProjectB"),
				InferredBursts: make([]*career.Burst, 0),
				InferredFacts:  make([]*career.Fact, 0),
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: make([]*career.Skill, 0),
				RejectedItems:  make(map[string]string),
			}

			submitMsg := SubmitCompleteMsg{}
			_ = intent.Update(submitMsg)

			Expect(intent.submitModal).NotTo(BeNil())
		})

		It("handles modal dismissal and returns to review state", func() {
			intent.currentState = StateSubmit
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("", "Event for review", "TestCorp", "ProjectC"),
				InferredBursts: make([]*career.Burst, 0),
				InferredFacts:  make([]*career.Fact, 0),
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: make([]*career.Skill, 0),
				RejectedItems:  make(map[string]string),
			}

			submitMsg := SubmitCompleteMsg{}
			intent.Update(submitMsg)

			dismissMsg := DismissModalMsg{}
			intent.Update(dismissMsg)

			Expect(intent.currentState).To(Equal(StateReview))
			Expect(intent.submitModal).To(BeNil())
			_, ok := intent.activeView.(*event.Review)
			Expect(ok).To(BeTrue())
		})

		It("completes the workflow on review submission", func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("", "Final event", "TestCorp", "ProjectD"),
				InferredBursts: make([]*career.Burst, 0),
				InferredFacts:  make([]*career.Fact, 0),
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: make([]*career.Skill, 0),
				RejectedItems:  make(map[string]string),
			}

			reviewData := event.ReviewResult{
				Event:  display.EventFromDomain(intent.reviewState.Event),
				Bursts: []display.Burst{},
				Facts:  []display.Fact{},
				Skills: []display.SkillSuggestion{},
			}

			result := &widgets.SubmitViewResult{
				FormData: reviewData,
			}

			_ = intent.handleViewResult(result)
			Expect(intent.currentState).To(Equal(StateSubmit))
		})
	})

	Describe("View produces ViewResult → Intent responds correctly", func() {
		It("processes KeyMsg from strategy view and triggers transition", func() {
			Expect(intent.currentState).To(Equal(StateChooseStrategy))

			keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
			_ = intent.Update(keyMsg)

			Expect(intent.currentState).NotTo(Equal(StateChooseStrategy))
		})

		It("handles cancellation from strategy view and marks intent as cancelled", func() {
			Expect(intent.currentState).To(Equal(StateChooseStrategy))

			keyMsg := tea.KeyMsg{Type: tea.KeyEsc}
			_ = intent.Update(keyMsg)

			Expect(intent.result.Status).To(Equal(intents.Cancelled))
		})

		It("handles cancellation from form view without previous event", func() {
			intent.currentState = StateForm
			intent.strategy = "quick"
			intent.context.PreviousEvent = nil

			result := &widgets.CancelViewResult{}
			_ = intent.handleViewResult(result)

			Expect(intent.currentState).To(Equal(StateChooseStrategy))
		})

		It("handles cancellation from form view with previous event", func() {
			previousEvent := fixtures.EventWith("prev-1", "Previous event", "OldCorp", "OldProject")
			intent.context.PreviousEvent = previousEvent
			intent.currentState = StateForm
			intent.strategy = "quick"

			result := &widgets.CancelViewResult{}
			_ = intent.handleViewResult(result)

			Expect(intent.result.Status).To(Equal(intents.Cancelled))
		})
	})

	Describe("View() renders correctly through state transitions", func() {
		It("renders strategy-related content after init", func() {
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders form-related content after strategy transition", func() {
			result := &widgets.NavigateViewResult{
				ResultData: CaptureStrategy("quick"),
			}
			_ = intent.handleViewResult(result)

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			Expect(intent.currentState).To(Equal(StateForm))
		})

		It("renders review-related content after form submission", func() {
			intent.currentState = StateForm
			intent.strategy = "quick"

			evt := fixtures.EventWith("", "Event for review render here", "TestCorp", "ProjectE")
			formData := forms.GetCaptureEventFormData(evt)
			formData.SubmitConfirmed = true

			result := &widgets.SubmitViewResult{
				FormData: formData,
			}
			_ = intent.handleViewResult(result)

			Expect(intent.reviewState).NotTo(BeNil())
			Expect(intent.reviewState.Event).NotTo(BeNil())
		})

		It("returns meaningful message when intent is inactive", func() {
			intent.active = false
			view := intent.View()
			Expect(view).To(ContainSubstring("not active"))
		})
	})

	Describe("State transitions maintain consistency", func() {
		It("prevents invalid state transitions via proper error handling", func() {
			intent.currentState = StateReview
			result := &widgets.ErrorViewResult{
				Err:     nil,
				Message: "Test validation error",
			}

			_ = intent.handleViewResult(result)
			Expect(intent.result.Status).To(Equal(intents.Failed))
		})

		It("preserves activeView during valid transitions", func() {
			initialView := intent.activeView
			Expect(initialView).NotTo(BeNil())

			result := &widgets.NavigateViewResult{
				ResultData: CaptureStrategy("quick"),
			}
			_ = intent.handleViewResult(result)

			Expect(intent.activeView).NotTo(BeNil())
			Expect(intent.currentState).To(Equal(StateForm))
		})

		It("initializes reviewState with correct defaults when form is submitted", func() {
			intent.currentState = StateForm
			intent.strategy = "quick"

			evt := fixtures.EventWith("", "Event with defaults", "DefaultCorp", "DefaultProj")
			formData := forms.GetCaptureEventFormData(evt)
			formData.SubmitConfirmed = true

			result := &widgets.SubmitViewResult{
				FormData: formData,
			}
			_ = intent.handleViewResult(result)

			Expect(intent.reviewState).NotTo(BeNil())
			Expect(intent.reviewState.AcceptedBursts).To(BeEmpty())
			Expect(intent.reviewState.AcceptedFacts).To(BeEmpty())
			Expect(intent.reviewState.RejectedItems).To(BeEmpty())
		})
	})

	Describe("Intent results are properly populated", func() {
		It("completes with correct result data after workflow", func() {
			intent.currentState = StateSubmit
			evt := fixtures.EventWith("final-event", "Final workflow event", "FinalCorp", "FinalProj")
			intent.reviewState = &ReviewInferredEventState{
				Event:          evt,
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: make([]*career.Skill, 0),
				RejectedItems:  make(map[string]string),
			}

			submitData := event.ReviewResult{
				Event:  display.EventFromDomain(evt),
				Bursts: display.BurstsFromDomain(make([]*career.Burst, 0)),
				Facts:  display.FactsFromDomain(make([]*career.Fact, 0)),
				Skills: []display.SkillSuggestion{},
			}

			result := &widgets.SubmitViewResult{
				FormData: submitData,
			}
			_ = intent.handleViewResult(result)

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Status).To(Equal(intents.Completed))
			Expect(intent.result.Data).NotTo(BeNil())
			Expect(intent.result.Data.Event).NotTo(BeNil())
			Expect(intent.result.Data.Event.Text).To(Equal("Final workflow event"))
		})

		It("marks as cancelled with proper status", func() {
			result := &widgets.CancelViewResult{}
			_ = intent.handleViewResult(result)

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Status).To(Equal(intents.Cancelled))
		})

		It("marks as failed with error information", func() {
			errResult := &widgets.ErrorViewResult{
				Err:     nil,
				Message: "Integration test error",
			}
			_ = intent.handleViewResult(errResult)

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Status).To(Equal(intents.Failed))
		})
	})

	Describe("HandleNavigate with metadata editing", func() {
		It("opens metadata editing modal when edit_metadata action is triggered", func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("", "Event for metadata edit", "TestCorp", "TestProj"),
				InferredBursts: make([]*career.Burst, 0),
				InferredFacts:  make([]*career.Fact, 0),
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: make([]*career.Skill, 0),
				RejectedItems:  make(map[string]string),
			}

			result := &widgets.NavigateViewResult{
				ResultData: "edit_metadata",
			}
			_ = intent.handleViewResult(result)

			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeMetadata))
			Expect(intent.reviewState.metadataModal).NotTo(BeNil())
		})

		It("fails gracefully when edit_metadata is triggered without event", func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          nil,
				InferredBursts: make([]*career.Burst, 0),
				InferredFacts:  make([]*career.Fact, 0),
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: make([]*career.Skill, 0),
				RejectedItems:  make(map[string]string),
			}

			result := &widgets.NavigateViewResult{
				ResultData: "edit_metadata",
			}
			_ = intent.handleViewResult(result)

			Expect(intent.result.Status).To(Equal(intents.Failed))
		})
	})

	Describe("HandleNavigate with burst suggestions", func() {
		It("opens burst suggestion modal when suggest_bursts action is triggered", func() {
			intent.currentState = StateReview
			burst := fixtures.Burst("burst-1", "event-1")
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("", "Event with bursts", "TestCorp", "TestProj"),
				InferredBursts: []*career.Burst{burst},
				InferredFacts:  make([]*career.Fact, 0),
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: make([]*career.Skill, 0),
				RejectedItems:  make(map[string]string),
			}

			result := &widgets.NavigateViewResult{
				ResultData: "suggest_bursts",
			}
			_ = intent.handleViewResult(result)

			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeBursts))
			Expect(intent.reviewState.burstModal).NotTo(BeNil())
		})
	})

	Describe("HandleNavigate with fact suggestions", func() {
		It("opens fact suggestion modal when suggest_facts action is triggered", func() {
			intent.currentState = StateReview
			fact := fixtures.FactWith("fact-1", "Test fact description")
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("", "Event with facts", "TestCorp", "TestProj"),
				InferredBursts: make([]*career.Burst, 0),
				InferredFacts:  []*career.Fact{fact},
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: make([]*career.Skill, 0),
				RejectedItems:  make(map[string]string),
			}

			result := &widgets.NavigateViewResult{
				ResultData: "suggest_facts",
			}
			_ = intent.handleViewResult(result)

			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeFacts))
			Expect(intent.reviewState.factSuggestionModal).NotTo(BeNil())
		})

		It("handles nil facts gracefully in fact suggestions", func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("", "Event with nil facts", "TestCorp", "TestProj"),
				InferredBursts: make([]*career.Burst, 0),
				InferredFacts:  []*career.Fact{nil, fixtures.FactWith("fact-1", "Test")},
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: make([]*career.Skill, 0),
				RejectedItems:  make(map[string]string),
			}

			result := &widgets.NavigateViewResult{
				ResultData: "suggest_facts",
			}
			_ = intent.handleViewResult(result)

			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeFacts))
			Expect(intent.reviewState.factSuggestionModal).NotTo(BeNil())
		})
	})

	Describe("HandleNavigate with skill suggestions", func() {
		It("opens skill suggestion modal when suggest_skills action is triggered", func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("", "Event with skills", "TestCorp", "TestProj"),
				InferredBursts: make([]*career.Burst, 0),
				InferredFacts:  make([]*career.Fact, 0),
				InferredSkills: []skillinference.SkillSuggestion{
					{
						Name:       "Go",
						Category:   "backend",
						Confidence: 0.95,
						EventIDs:   []string{"event-1"},
						Contexts:   []string{"Used Go for microservices"},
					},
				},
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: make([]*career.Skill, 0),
				RejectedItems:  make(map[string]string),
			}

			result := &widgets.NavigateViewResult{
				ResultData: "suggest_skills",
			}
			_ = intent.handleViewResult(result)

			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeSkills))
			Expect(intent.reviewState.skillModal).NotTo(BeNil())
		})
	})

	Describe("HandleNavigate with invalid actions", func() {
		It("fails when unknown navigation action is provided", func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("", "Event for invalid action", "TestCorp", "TestProj"),
				InferredBursts: make([]*career.Burst, 0),
				InferredFacts:  make([]*career.Fact, 0),
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: make([]*career.Skill, 0),
				RejectedItems:  make(map[string]string),
			}

			result := &widgets.NavigateViewResult{
				ResultData: "unknown_action",
			}
			_ = intent.handleViewResult(result)

			Expect(intent.result.Status).To(Equal(intents.Failed))
		})

		It("fails when navigate result has invalid data type", func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("", "Event for invalid type", "TestCorp", "TestProj"),
				InferredBursts: make([]*career.Burst, 0),
				InferredFacts:  make([]*career.Fact, 0),
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: make([]*career.Skill, 0),
				RejectedItems:  make(map[string]string),
			}

			result := &widgets.NavigateViewResult{
				ResultData: 12345, // Invalid type
			}
			_ = intent.handleViewResult(result)

			Expect(intent.result.Status).To(Equal(intents.Failed))
		})
	})

	Describe("Intent results are properly populated", func() {
		It("completes with correct result data after workflow", func() {
			intent.currentState = StateSubmit
			evt := fixtures.EventWith("final-event", "Final workflow event", "FinalCorp", "FinalProj")
			intent.reviewState = &ReviewInferredEventState{
				Event:          evt,
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: make([]*career.Skill, 0),
				RejectedItems:  make(map[string]string),
			}

			submitData := event.ReviewResult{
				Event:  display.EventFromDomain(evt),
				Bursts: display.BurstsFromDomain(make([]*career.Burst, 0)),
				Facts:  display.FactsFromDomain(make([]*career.Fact, 0)),
				Skills: []display.SkillSuggestion{},
			}

			result := &widgets.SubmitViewResult{
				FormData: submitData,
			}
			_ = intent.handleViewResult(result)

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Status).To(Equal(intents.Completed))
			Expect(intent.result.Data).NotTo(BeNil())
			Expect(intent.result.Data.Event).NotTo(BeNil())
			Expect(intent.result.Data.Event.Text).To(Equal("Final workflow event"))
		})

		It("marks as cancelled with proper status", func() {
			result := &widgets.CancelViewResult{}
			_ = intent.handleViewResult(result)

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Status).To(Equal(intents.Cancelled))
		})

		It("marks as failed with error information", func() {
			errResult := &widgets.ErrorViewResult{
				Err:     nil,
				Message: "Integration test error",
			}
			_ = intent.handleViewResult(errResult)

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Status).To(Equal(intents.Failed))
		})
	})
})
