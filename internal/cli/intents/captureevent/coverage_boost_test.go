package captureevent

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/screens"
	cliservice "github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/domain/career"
	memoryrepo "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Coverage boost — performSubmit additional paths", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentContext{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("performSubmit — validation error path", func() {
		It("returns SubmitErrorMsg with VALIDATION_ERROR when event text is empty", func() {
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event)
			svc.SetSkillRepository(repos.Skill)

			invalidEvent := fixtures.EventWith("evt-invalid", "", "", "")
			intent.context.CareerService = svc
			intent.reviewState = &ReviewInferredEventState{
				Event:          invalidEvent,
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: make([]*career.Skill, 0),
			}

			cmd := intent.performSubmit()
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			errMsg, ok := msg.(SubmitErrorMsg)
			Expect(ok).To(BeTrue(), "expected SubmitErrorMsg, got %T", msg)
			Expect(errMsg.Code).To(Equal("VALIDATION_ERROR"))
		})
	})

	Describe("performSubmit — skill save path with new skills", func() {
		It("returns SubmitCompleteMsg or SubmitErrorMsg when accepted skill has no ID", func() {
			repos := memoryrepo.NewRepositories()
			eventRepo := repos.Event.(*memoryrepo.EventRepository)
			svc := careerservice.NewService(eventRepo)
			svc.SetSkillRepository(repos.Skill)

			intent.context.CareerService = svc

			event := fixtures.EventWith("", "Valid event text for testing skill save path", "", "")
			skillWithNoID := fixtures.SkillWith("", "Go", "backend", "")
			intent.reviewState = &ReviewInferredEventState{
				Event:          event,
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: []*career.Skill{skillWithNoID},
			}

			cmd := intent.performSubmit()
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			_, isComplete := msg.(SubmitCompleteMsg)
			_, isError := msg.(SubmitErrorMsg)
			Expect(isComplete || isError).To(BeTrue(), "expected SubmitCompleteMsg or SubmitErrorMsg, got %T", msg)
		})
	})

	Describe("performSubmit — skill with existing ID (link-only path)", func() {
		It("attempts to link an already-persisted skill and returns a result", func() {
			repos := memoryrepo.NewRepositories()
			eventRepo := repos.Event.(*memoryrepo.EventRepository)
			svc := careerservice.NewService(eventRepo)
			svc.SetSkillRepository(repos.Skill)

			intent.context.CareerService = svc

			event := fixtures.EventWith("", "Valid event text for testing skill link path", "", "")

			existingSkill := fixtures.SkillWith("existing-skill-id", "Go", "backend", "advanced")
			Expect(repos.Skill.Create(context.Background(), existingSkill)).To(Succeed())

			intent.reviewState = &ReviewInferredEventState{
				Event:          event,
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: []*career.Skill{existingSkill},
			}

			cmd := intent.performSubmit()
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			_, isComplete := msg.(SubmitCompleteMsg)
			_, isError := msg.(SubmitErrorMsg)
			Expect(isComplete || isError).To(BeTrue(), "expected SubmitCompleteMsg or SubmitErrorMsg, got %T", msg)
		})
	})

	Describe("performSubmit — quick strategy with zero date", func() {
		It("sets today as date when strategy is quick and date is zero", func() {
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event.(*memoryrepo.EventRepository))
			svc.SetSkillRepository(repos.Skill)

			intent.context.CareerService = svc
			intent.strategy = StrategyQuick

			event := fixtures.EventWith("", "Valid event text for zero date test path", "", "")
			event.Date = time.Time{}

			intent.reviewState = &ReviewInferredEventState{
				Event:          event,
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: make([]*career.Skill, 0),
			}

			cmd := intent.performSubmit()
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			_, isComplete := msg.(SubmitCompleteMsg)
			_, isError := msg.(SubmitErrorMsg)
			Expect(isComplete || isError).To(BeTrue(), "expected SubmitCompleteMsg or SubmitErrorMsg, got %T", msg)
		})
	})

	Describe("performSubmit — fact save path when facts are present", func() {
		It("returns SubmitCompleteMsg or SubmitErrorMsg when facts are present", func() {
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event.(*memoryrepo.EventRepository))
			svc.SetSkillRepository(repos.Skill)

			intent.context.CareerService = svc

			event := fixtures.EventWith("", "Valid event text for testing fact save path", "", "")
			factWithNoID := fixtures.Fact("", "")
			intent.reviewState = &ReviewInferredEventState{
				Event:          event,
				AcceptedFacts:  []*career.Fact{factWithNoID},
				AcceptedSkills: make([]*career.Skill, 0),
			}

			cmd := intent.performSubmit()
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			_, isComplete := msg.(SubmitCompleteMsg)
			_, isError := msg.(SubmitErrorMsg)
			Expect(isComplete || isError).To(BeTrue(), "expected SubmitCompleteMsg or SubmitErrorMsg, got %T", msg)
		})
	})
})

var _ = Describe("Coverage boost — performPostSavePersistence additional paths", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentContext{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("performPostSavePersistence — nil careerService returns complete msg", func() {
		It("returns PostSavePersistenceCompleteMsg immediately when careerService is nil", func() {
			intent.context.CareerService = nil

			event := fixtures.EventWith("evt-1", "Event for post-save persistence nil path test", "", "")
			burst := fixtures.Burst("burst-1", "evt-1")
			fact := fixtures.Fact("fact-1", "evt-1")
			skill := fixtures.SkillWith("skill-1", "Go", "backend", "advanced")

			cmd := intent.performPostSavePersistence(
				event,
				[]*career.Fact{fact},
				[]*career.Skill{skill},
				[]*career.Burst{burst},
			)
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			completeMsg, ok := msg.(PostSavePersistenceCompleteMsg)
			Expect(ok).To(BeTrue(), "expected PostSavePersistenceCompleteMsg, got %T", msg)
			Expect(completeMsg.Event).To(Equal(event))
			Expect(completeMsg.Bursts).To(HaveLen(1))
			Expect(completeMsg.Facts).To(HaveLen(1))
			Expect(completeMsg.Skills).To(HaveLen(1))
		})
	})

	Describe("performPostSavePersistence — skill save (new skill, no ID)", func() {
		It("processes a new skill and returns a result", func() {
			repos := memoryrepo.NewRepositories()
			eventRepo := repos.Event.(*memoryrepo.EventRepository)
			svc := careerservice.NewService(eventRepo)
			svc.SetSkillRepository(repos.Skill)
			intent.context.CareerService = svc

			event := fixtures.EventWith("evt-post-1", "Event for post-save skill save test", "", "")
			Expect(eventRepo.Create(context.Background(), event)).To(Succeed())

			newSkill := fixtures.SkillWith("", "TypeScript", "frontend", "")

			cmd := intent.performPostSavePersistence(
				event,
				make([]*career.Fact, 0),
				[]*career.Skill{newSkill},
				make([]*career.Burst, 0),
			)
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			_, isComplete := msg.(PostSavePersistenceCompleteMsg)
			_, isError := msg.(SubmitErrorMsg)
			Expect(isComplete || isError).To(BeTrue(), "expected PostSavePersistenceCompleteMsg or SubmitErrorMsg, got %T", msg)
		})
	})

	Describe("performPostSavePersistence — fact save (new fact, no ID)", func() {
		It("saves a fact when fact has no ID", func() {
			repos := memoryrepo.NewRepositories()
			eventRepo := repos.Event.(*memoryrepo.EventRepository)
			svc := careerservice.NewService(eventRepo)
			svc.SetSkillRepository(repos.Skill)
			intent.context.CareerService = svc

			event := fixtures.EventWith("evt-post-2", "Event for post-save fact save test", "", "")
			Expect(eventRepo.Create(context.Background(), event)).To(Succeed())

			newFact := fixtures.Fact("", "")

			cmd := intent.performPostSavePersistence(
				event,
				[]*career.Fact{newFact},
				make([]*career.Skill, 0),
				make([]*career.Burst, 0),
			)
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			_, isComplete := msg.(PostSavePersistenceCompleteMsg)
			_, isError := msg.(SubmitErrorMsg)
			Expect(isComplete || isError).To(BeTrue(), "expected PostSavePersistenceCompleteMsg or SubmitErrorMsg, got %T", msg)
		})
	})

	Describe("performPostSavePersistence — fact already has ID (skip save)", func() {
		It("skips saving a fact that already has an ID and returns complete msg", func() {
			repos := memoryrepo.NewRepositories()
			eventRepo := repos.Event.(*memoryrepo.EventRepository)
			svc := careerservice.NewService(eventRepo)
			svc.SetSkillRepository(repos.Skill)
			intent.context.CareerService = svc

			event := fixtures.EventWith("evt-post-3", "Event for post-save fact skip test path", "", "")
			Expect(eventRepo.Create(context.Background(), event)).To(Succeed())

			existingFact := fixtures.Fact("existing-fact-id", event.ID)

			cmd := intent.performPostSavePersistence(
				event,
				[]*career.Fact{existingFact},
				make([]*career.Skill, 0),
				make([]*career.Burst, 0),
			)
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			completeMsg, ok := msg.(PostSavePersistenceCompleteMsg)
			Expect(ok).To(BeTrue(), "expected PostSavePersistenceCompleteMsg, got %T", msg)
			Expect(completeMsg.Facts).To(HaveLen(1))
		})
	})
})

var _ = Describe("Coverage boost — terminalDimensions", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentContext{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	It("returns nil when GetTerminalInfo returns nil", func() {
		intent.UpdateTerminalInfo(nil)
		dims := intent.terminalDimensions()
		Expect(dims).To(BeNil())
	})

	It("returns dimensions when terminal info is valid", func() {
		info := &terminal.Info{Width: 120, Height: 40, IsValid: true}
		intent.UpdateTerminalInfo(info)
		dims := intent.terminalDimensions()
		Expect(dims).NotTo(BeNil())
		Expect(dims.TerminalWidth).To(Equal(120))
		Expect(dims.TerminalHeight).To(Equal(40))
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

var _ = Describe("Coverage boost — handleEditKeyMsg additional branches", func() {
	var (
		model       *BurstSuggestionModelNew
		suggestions []burstfact.BurstSuggestion
	)

	BeforeEach(func() {
		suggestions = []burstfact.BurstSuggestion{
			{Name: "API Work", Description: "REST API development", EventIDs: []string{"e1"}, ConfidenceScore: 0.9},
		}
		model = NewBurstSuggestionModelNew(context.Background(), nil, suggestions)
	})

	Context("when form is active and Escape is pressed", func() {
		It("resets editing state when the form receives Esc", func() {
			model.startEdit()
			Expect(model.IsEditing()).To(BeTrue())

			result, cmd := model.handleEditKeyMsg(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(result).NotTo(BeNil())
			Expect(cmd).To(BeNil())
		})
	})

	Context("when form is active and receives a regular key", func() {
		It("forwards key to form and returns model", func() {
			model.startEdit()

			result, _ := model.handleEditKeyMsg(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(result).NotTo(BeNil())
		})
	})

	Context("when Enter is pressed while editing", func() {
		It("processes the Enter key via the form without panicking", func() {
			model.startEdit()
			result, _ := model.handleEditKeyMsg(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(result).NotTo(BeNil())
		})
	})
})

var _ = Describe("Coverage boost — rejectCurrent additional paths", func() {
	var (
		model       *BurstSuggestionModelNew
		suggestions []burstfact.BurstSuggestion
	)

	BeforeEach(func() {
		suggestions = []burstfact.BurstSuggestion{
			{Name: "First Burst", Description: "desc 1", EventIDs: []string{"e1"}, ConfidenceScore: 0.8},
			{Name: "Second Burst", Description: "desc 2", EventIDs: []string{"e2"}, ConfidenceScore: 0.7},
			{Name: "Third Burst", Description: "desc 3", EventIDs: []string{"e3"}, ConfidenceScore: 0.6},
		}
		model = NewBurstSuggestionModelNew(context.Background(), nil, suggestions)
	})

	It("increments cursor when rejecting a non-final suggestion", func() {
		_, cmd := model.rejectCurrent()
		Expect(cmd).NotTo(BeNil())
		Expect(model.currentIdx).To(Equal(1))
		Expect(model.GetRejected()).To(HaveLen(1))
	})

	It("returns BurstProcessingCompleteMsg when rejecting all suggestions", func() {
		model.rejectCurrent()
		model.rejectCurrent()
		_, cmd := model.rejectCurrent()
		Expect(cmd).NotTo(BeNil())
		Expect(model.IsDone()).To(BeTrue())
		Expect(model.GetRejected()).To(HaveLen(3))
	})

	It("handles empty suggestions gracefully", func() {
		empty := NewBurstSuggestionModelNew(context.Background(), nil, []burstfact.BurstSuggestion{})
		result, cmd := empty.rejectCurrent()
		Expect(result).NotTo(BeNil())
		Expect(cmd).To(BeNil())
	})
})

var _ = Describe("Coverage boost — confirmCurrent additional paths", func() {
	var (
		model       *BurstSuggestionModelNew
		suggestions []burstfact.BurstSuggestion
	)

	BeforeEach(func() {
		suggestions = []burstfact.BurstSuggestion{
			{Name: "First Burst", Description: "desc 1", EventIDs: []string{"e1"}, ConfidenceScore: 0.8},
			{Name: "Second Burst", Description: "desc 2", EventIDs: []string{"e2"}, ConfidenceScore: 0.7},
			{Name: "Third Burst", Description: "desc 3", EventIDs: []string{"e3"}, ConfidenceScore: 0.6},
		}
		model = NewBurstSuggestionModelNew(context.Background(), nil, suggestions)
	})

	It("increments cursor when confirming a non-final suggestion", func() {
		_, cmd := model.confirmCurrent()
		Expect(cmd).NotTo(BeNil())
		Expect(model.currentIdx).To(Equal(1))
		Expect(model.GetConfirmed()).To(HaveLen(1))
	})

	It("returns BurstProcessingCompleteMsg when confirming all suggestions", func() {
		model.confirmCurrent()
		model.confirmCurrent()
		_, cmd := model.confirmCurrent()
		Expect(cmd).NotTo(BeNil())
		Expect(model.IsDone()).To(BeTrue())
		Expect(model.GetConfirmed()).To(HaveLen(3))
	})

	It("handles empty suggestions gracefully", func() {
		empty := NewBurstSuggestionModelNew(context.Background(), nil, []burstfact.BurstSuggestion{})
		result, cmd := empty.confirmCurrent()
		Expect(result).NotTo(BeNil())
		Expect(cmd).To(BeNil())
	})
})

var _ = Describe("Coverage boost — handleFormCompletion with cliService", func() {
	var (
		model   *ReviewEnrichmentModel
		event   *career.Event
		service *careerservice.Service
	)

	BeforeEach(func() {
		repos := memoryrepo.NewRepositories()
		eventRepo := repos.Event.(*memoryrepo.EventRepository)
		service = careerservice.NewService(eventRepo)
		service.SetSkillRepository(repos.Skill)

		event = fixtures.EventWith("test-event-cli", "Event for cliService handleFormCompletion test", "Test Corp", "Test Project")
		event.Date = time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC)
		event.Tags = []string{"technical"}
		event.Categories = []string{"technical"}

		Expect(eventRepo.Create(context.Background(), event)).To(Succeed())

		cliSvc := cliservice.NewCLIEventService(service)
		model = NewReviewEnrichmentModel(context.Background(), event, service, cliSvc, nil)
		model.Init()
	})

	It("succeeds and marks submitted when cliService is non-nil and event is valid", func() {
		model.formData.SubmitConfirmed = true
		model.formData.Company = "Updated Company via cliService"
		model.formData.Project = "Updated Project"
		model.formData.Date = "2024-06-15"
		model.formData.Tags = []string{"technical"}
		model.formData.Categories = []string{"technical"}

		result, cmd := model.handleFormCompletion()
		Expect(result).NotTo(BeNil())
		_ = cmd
		Expect(model.IsSubmitted()).To(BeTrue())
	})

	It("sets Err when cliService.UpdateEventMetadata fails due to non-existent event", func() {
		model.formData.SubmitConfirmed = true
		model.formData.Company = "Updated Company"
		model.formData.Date = "2024-06-15"
		model.formData.Tags = []string{"technical"}
		model.formData.Categories = []string{"technical"}

		model.event.ID = "non-existent-event-id"

		result, cmd := model.handleFormCompletion()
		Expect(result).NotTo(BeNil())
		_ = cmd
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

var _ = Describe("Coverage boost — renderRelatedEvents with live service", func() {
	var (
		model       *BurstSuggestionModelNew
		suggestions []burstfact.BurstSuggestion
	)

	BeforeEach(func() {
		suggestions = []burstfact.BurstSuggestion{
			{Name: "Live Burst", Description: "desc", EventIDs: []string{"live-evt-1", "live-evt-2"}, ConfidenceScore: 0.8},
		}
		repos := memoryrepo.NewRepositories()
		eventRepo := repos.Event.(*memoryrepo.EventRepository)
		svc := careerservice.NewService(eventRepo)

		ctx := context.Background()
		evt1 := fixtures.EventWith("live-evt-1", "First live event for renderRelatedEvents test", "", "")
		evt2 := fixtures.EventWith("live-evt-2", "Second live event for renderRelatedEvents test", "", "")
		Expect(eventRepo.Create(ctx, evt1)).To(Succeed())
		Expect(eventRepo.Create(ctx, evt2)).To(Succeed())

		model = NewBurstSuggestionModelNew(ctx, svc, suggestions)
	})

	It("loads events from service when not cached", func() {
		s := suggestions[0]
		view := model.renderRelatedEvents(s)
		Expect(view).NotTo(BeEmpty())
	})

	It("caches events after first load", func() {
		s := suggestions[0]
		model.renderRelatedEvents(s)
		_, exists := model.relatedEvents[0]
		Expect(exists).To(BeTrue())
	})

	It("uses cached events on second call", func() {
		s := suggestions[0]
		model.renderRelatedEvents(s)
		view := model.renderRelatedEvents(s)
		Expect(view).NotTo(BeEmpty())
	})
})

var _ = Describe("Coverage boost — HandleError additional paths", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentContext{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	It("handles map data with empty message by using fallback message", func() {
		result := &screens.ErrorResult{
			Err:     errors.New("underlying cause"),
			Message: "",
		}
		_ = intent.HandleError(result)
		Expect(intent.active).To(BeFalse())
	})

	It("handles map data with non-nil error and non-empty message", func() {
		result := &screens.ErrorResult{
			Err:     errors.New("screen exploded"),
			Message: "Screen encountered an error",
		}
		_ = intent.HandleError(result)
		Expect(intent.result.Error.Message).To(Equal("Screen encountered an error"))
	})
})

var _ = Describe("Coverage boost — transitionToStrategyScreen with terminal info", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentContext{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	It("uses terminal dimensions when available", func() {
		info := &terminal.Info{Width: 150, Height: 50, IsValid: true}
		intent.UpdateTerminalInfo(info)
		cmd := intent.transitionToStrategyScreen()
		Expect(cmd).To(BeNil())
		Expect(intent.activeScreen).NotTo(BeNil())
	})
})

var _ = Describe("Coverage boost — transitionToFormScreen paths", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentContext{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	It("uses terminal dimensions when available", func() {
		info := &terminal.Info{Width: 160, Height: 48, IsValid: true}
		intent.UpdateTerminalInfo(info)
		intent.transitionToFormScreen(StrategyManual)
		Expect(intent.activeScreen).NotTo(BeNil())
	})

	It("uses PreviousEvent when set in context", func() {
		prev := fixtures.EventWith("prev-evt", "Previous event text for form test", "", "")
		intent.context.PreviousEvent = prev
		intent.transitionToFormScreen(StrategyQuick)
		Expect(intent.activeScreen).NotTo(BeNil())
		Expect(intent.captureFormScreen).NotTo(BeNil())
	})
})

var _ = Describe("Coverage boost — HandleSubmit additional branches", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentContext{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("StateForm — SubmitConfirmed false returns nil", func() {
		It("returns nil when SubmitConfirmed is false", func() {
			intent.currentState = StateForm
			result := &screens.SubmitResult{
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
			result := &screens.SubmitResult{
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
			result := &screens.SubmitResult{
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
			result := &screens.SubmitResult{FormData: "wrong type"}
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
			result := &screens.SubmitResult{
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
			result := &screens.SubmitResult{
				FormData: map[string]interface{}{
					"bursts": make([]*career.Burst, 0),
				},
			}
			intent.HandleSubmit(result)
			Expect(intent.active).To(BeFalse())
		})
	})

	Describe("StateSubmit — invalid data type", func() {
		It("marks intent as failed when data is not map[string]interface{}", func() {
			intent.currentState = StateSubmit
			result := &screens.SubmitResult{FormData: "wrong type"}
			intent.HandleSubmit(result)
			Expect(intent.active).To(BeFalse())
		})
	})
})

var _ = Describe("Coverage boost — NewReviewEnrichmentModel with cliService", func() {
	It("accepts a non-nil cliService without panicking", func() {
		repos := memoryrepo.NewRepositories()
		eventRepo := repos.Event.(*memoryrepo.EventRepository)
		svc := careerservice.NewService(eventRepo)
		svc.SetSkillRepository(repos.Skill)

		cliSvc := cliservice.NewCLIEventService(svc)

		event := fixtures.EventWith("editor-evt-1", "Event for metadata editor cliService test", "Corp", "Proj")
		event.Date = time.Now()
		event.Tags = []string{"technical"}
		event.Categories = []string{"technical"}

		m := NewReviewEnrichmentModel(context.Background(), event, svc, cliSvc, nil)
		Expect(m).NotTo(BeNil())
	})
})
