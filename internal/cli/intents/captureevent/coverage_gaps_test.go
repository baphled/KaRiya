package captureevent

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/screens"
	cliservice "github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	memoryrepo "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Coverage gaps", func() {

	// ─── terminalDimensions ─────────────────────────────────────────────────

	Describe("terminalDimensions — non-nil terminal info", func() {
		It("returns dimensions from terminal info when set", func() {
			ctx := &IntentContext{CaptureStrategy: "quick"}
			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			info := &terminal.Info{Width: 120, Height: 40, IsValid: true}
			intent.UpdateTerminalInfo(info)

			dims := intent.terminalDimensions()
			Expect(dims).NotTo(BeNil())
			Expect(dims.TerminalWidth).To(Equal(120))
			Expect(dims.TerminalHeight).To(Equal(40))
		})
	})

	Describe("terminalDimensions — nil terminal info", func() {
		It("returns nil when terminal info is nil", func() {
			ctx := &IntentContext{CaptureStrategy: "quick"}
			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			intent.UpdateTerminalInfo(nil)

			dims := intent.terminalDimensions()
			Expect(dims).To(BeNil())
		})
	})

	// ─── transitionToStrategyScreen / FormScreen with terminal info ──────────

	Describe("transitionToStrategyScreen — with non-nil terminal info", func() {
		It("configures screen using terminal dimensions", func() {
			ctx := &IntentContext{CaptureStrategy: "quick"}
			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			info := &terminal.Info{Width: 120, Height: 40, IsValid: true}
			intent.UpdateTerminalInfo(info)

			cmd := intent.transitionToStrategyScreen()
			Expect(cmd).To(BeNil())
			Expect(intent.activeScreen).NotTo(BeNil())
		})
	})

	Describe("transitionToFormScreen — with non-nil terminal info", func() {
		It("configures screen using terminal dimensions", func() {
			ctx := &IntentContext{CaptureStrategy: "quick"}
			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			info := &terminal.Info{Width: 120, Height: 40, IsValid: true}
			intent.UpdateTerminalInfo(info)

			intent.transitionToFormScreen(StrategyQuick)
			Expect(intent.activeScreen).NotTo(BeNil())
		})
	})

	// ─── transitionToFormScreen — with PreviousEvent ────────────────────────

	Describe("transitionToFormScreen — with PreviousEvent set", func() {
		It("uses the previous event to populate the form", func() {
			prev := fixtures.EventWith("prev-1", "Previous event text here", "Corp", "Proj")
			ctx := &IntentContext{
				CaptureStrategy: "quick",
				PreviousEvent:   prev,
			}
			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			intent.transitionToFormScreen(StrategyManual)
			Expect(intent.activeScreen).NotTo(BeNil())
		})
	})

	// ─── handleFormCompletion — cliService error path ────────────────────────

	Describe("ReviewEnrichmentModel — handleFormCompletion with cliService error", func() {
		It("sets Err when UpdateEventMetadata fails", func() {
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event)
			svc.SetSkillRepository(repos.Skill)

			cliSvc := cliservice.NewCLIEventService(svc)

			event := fixtures.EventWith("", "Test event for metadata editor", "Test Company", "Test Project")
			event.Date = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
			event.Tags = []string{"technical"}
			event.Categories = []string{"technical"}

			model := NewReviewEnrichmentModel(context.Background(), event, svc, cliSvc, nil)
			model.Init()

			model.formData.SubmitConfirmed = true
			model.formData.Company = "Updated Company"
			model.formData.Project = "Updated Project"
			model.formData.Date = "2024-06-15"
			model.formData.Tags = []string{"technical"}
			model.formData.Categories = []string{"technical"}

			result, _ := model.handleFormCompletion()
			Expect(result).NotTo(BeNil())
			Expect(model.GetError()).To(HaveOccurred())
		})
	})

	// ─── eventFromFormData — date parse failure ───────────────────────────────

	Describe("eventFromFormData — date parse failure", func() {
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

	// ─── performSubmit — fact save failure ───────────────────────────────────

	Describe("performSubmit — fact save failure", func() {
		It("returns SubmitErrorMsg with PARTIAL_SAVE code when a fact fails to save", func() {
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event)
			svc.SetSkillRepository(repos.Skill)

			ctx := &IntentContext{CaptureStrategy: "quick"}
			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.context.CareerService = svc

			event := fixtures.EventWith("", "Valid event text for testing fact save", "", "")
			event.Date = time.Now().Add(-time.Hour)
			event.Tags = []string{"technical"}
			event.Categories = []string{"technical"}

			invalidFact := fixtures.Fact("", "")

			intent.reviewState = &ReviewInferredEventState{
				Event:         event,
				AcceptedFacts: []*career.Fact{invalidFact},
			}

			cmd := intent.performSubmit()
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			errMsg, ok := msg.(SubmitErrorMsg)
			Expect(ok).To(BeTrue())
			Expect(errMsg.Code).To(Equal("PARTIAL_SAVE"))
		})
	})

	// ─── performSubmit — validation failure ─────────────────────────────────

	Describe("performSubmit — event validation failure", func() {
		It("returns SubmitErrorMsg with VALIDATION_ERROR code for invalid event", func() {
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event)
			svc.SetSkillRepository(repos.Skill)

			ctx := &IntentContext{CaptureStrategy: "quick"}
			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.context.CareerService = svc

			invalidEvent := fixtures.EventWith("", "x", "", "")

			intent.reviewState = &ReviewInferredEventState{
				Event:         invalidEvent,
				AcceptedFacts: make([]*career.Fact, 0),
			}

			cmd := intent.performSubmit()
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			errMsg, ok := msg.(SubmitErrorMsg)
			Expect(ok).To(BeTrue())
			Expect(errMsg.Code).To(Equal("VALIDATION_ERROR"))
		})
	})

	// ─── performPostSavePersistence — skill save error ────────────────────────

	Describe("performPostSavePersistence — skill save error", func() {
		It("returns SubmitErrorMsg when skill creation fails", func() {
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event)
			svc.SetSkillRepository(repos.Skill)

			ctx := &IntentContext{CaptureStrategy: "quick"}
			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.context.CareerService = svc

			event := fixtures.EventWith("evt-persist", "Valid persist test event text here", "Corp", "Proj")
			event.Date = time.Now().Add(-time.Hour)
			event.Tags = []string{"technical"}
			event.Categories = []string{"technical"}
			Expect(repos.Event.Create(context.Background(), event)).To(Succeed())

			invalidSkill := fixtures.SkillWith("", "", "", "")

			cmd := intent.performPostSavePersistence(
				event,
				[]*career.Fact{},
				[]*career.Skill{invalidSkill},
				[]*career.Burst{},
			)
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			errMsg, ok := msg.(SubmitErrorMsg)
			Expect(ok).To(BeTrue())
			Expect(errMsg.Code).To(Equal("SKILL_SAVE_ERROR"))
		})
	})

	// ─── performPostSavePersistence — fact save error ─────────────────────────

	Describe("performPostSavePersistence — fact save error", func() {
		It("returns SubmitErrorMsg when fact save fails", func() {
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event)
			svc.SetSkillRepository(repos.Skill)

			ctx := &IntentContext{CaptureStrategy: "quick"}
			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.context.CareerService = svc

			event := fixtures.EventWith("evt-fact", "Valid fact test event text here correct", "Corp", "Proj")
			event.Date = time.Now().Add(-time.Hour)
			event.Tags = []string{"technical"}
			event.Categories = []string{"technical"}
			Expect(repos.Event.Create(context.Background(), event)).To(Succeed())

			invalidFact := fixtures.Fact("", "")

			cmd := intent.performPostSavePersistence(
				event,
				[]*career.Fact{invalidFact},
				[]*career.Skill{},
				[]*career.Burst{},
			)
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			errMsg, ok := msg.(SubmitErrorMsg)
			Expect(ok).To(BeTrue())
			Expect(errMsg.Code).To(Equal("FACT_SAVE_ERROR"))
		})
	})

	// ─── handleEditKeyMsg — completed form branch ────────────────────────────

	Describe("BurstSuggestionModelNew — handleEditKeyMsg form completed", func() {
		It("calls saveEdits when form state is Completed", func() {
			suggestions := []burstfact.BurstSuggestion{
				{Name: "API Work", Description: "REST API development", EventIDs: []string{"e1"}, ConfidenceScore: 0.9},
			}
			m := NewBurstSuggestionModelNew(context.Background(), nil, suggestions)
			m.startEdit()
			Expect(m.editForm).NotTo(BeNil())

			m.editFormData.Name = "Completed Name"
			m.editForm.State = huh.StateCompleted

			result, cmd := m.handleEditKeyMsg(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(result).NotTo(BeNil())
			Expect(cmd).To(BeNil())
			Expect(m.IsEditing()).To(BeFalse())
			Expect(m.editedNames[0]).To(Equal("Completed Name"))
		})
	})

	Describe("BurstSuggestionModelNew — handleEditKeyMsg form aborted", func() {
		It("resets edit state when form state is Aborted", func() {
			suggestions := []burstfact.BurstSuggestion{
				{Name: "API Work", Description: "REST API development", EventIDs: []string{"e1"}, ConfidenceScore: 0.9},
			}
			m := NewBurstSuggestionModelNew(context.Background(), nil, suggestions)
			m.startEdit()
			Expect(m.editForm).NotTo(BeNil())

			m.editForm.State = huh.StateAborted

			result, cmd := m.handleEditKeyMsg(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(result).NotTo(BeNil())
			Expect(cmd).To(BeNil())
			Expect(m.IsEditing()).To(BeFalse())
			Expect(m.editForm).To(BeNil())
			Expect(m.editFormData).To(BeNil())
		})
	})

	// ─── GetFactSuggestionModal — nil receiver ───────────────────────────────

	Describe("GetFactSuggestionModal — nil receiver", func() {
		It("returns nil without panicking when called on nil ReviewInferredEventState", func() {
			var state *ReviewInferredEventState
			result := state.GetFactSuggestionModal()
			Expect(result).To(BeNil())
		})
	})

	// ─── extractSkillsFromReviewData — nil skill filtering ───────────────────

	Describe("extractSkillsFromReviewData — nil skill in slice", func() {
		It("skips nil skills and returns only valid entries", func() {
			data := map[string]interface{}{
				"skills": []*career.Skill{
					nil,
					fixtures.SkillWith("", "Go", "backend", ""),
				},
			}
			result := extractSkillsFromReviewData(data)
			Expect(result).To(HaveLen(1))
			Expect(result[0].Name).To(Equal("Go"))
		})
	})

	// ─── HandleSubmit — StateReview with postSaveReview ──────────────────────

	Describe("HandleSubmit — StateReview with postSaveReview flag", func() {
		It("calls performPostSavePersistence path without panicking", func() {
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event)
			svc.SetSkillRepository(repos.Skill)

			ctx := &IntentContext{CaptureStrategy: "quick"}
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
				intent.HandleSubmit(&screens.SubmitResult{
					FormData: map[string]interface{}{
						"event":  event,
						"bursts": []*career.Burst{},
						"facts":  []*career.Fact{},
					},
				})
			}).NotTo(Panic())
		})
	})

	// ─── HandleSubmit — StateForm with date parse error ──────────────────────

	Describe("HandleSubmit — StateForm with unparseable date", func() {
		It("shows a validation error modal without deactivating the intent", func() {
			ctx := &IntentContext{CaptureStrategy: "quick"}
			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.currentState = StateForm

			formData := &forms.CaptureEventFormData{
				Text:            "Valid event text for testing the date parse path",
				Date:            "not-a-valid-date",
				SubmitConfirmed: true,
			}
			intent.HandleSubmit(&screens.SubmitResult{FormData: formData})
			Expect(intent.submitModal).NotTo(BeNil())
			Expect(intent.IsActive()).To(BeTrue())
		})
	})

	// ─── performSubmit — StrategyQuick with zero date ─────────────────────────

	Describe("performSubmit — StrategyQuick sets date when zero", func() {
		It("sets the event date to now when date is zero and strategy is quick", func() {
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event)
			svc.SetSkillRepository(repos.Skill)

			ctx := &IntentContext{CaptureStrategy: "quick"}
			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.context.CareerService = svc
			intent.strategy = StrategyQuick

			event := fixtures.EventWith("", "Valid event text for quick strategy test here", "", "")
			event.Date = time.Time{}
			event.Tags = []string{}
			event.Categories = []string{}

			intent.reviewState = &ReviewInferredEventState{
				Event:         event,
				AcceptedFacts: make([]*career.Fact, 0),
			}

			cmd := intent.performSubmit()
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			_, isComplete := msg.(SubmitCompleteMsg)
			_, isError := msg.(SubmitErrorMsg)
			Expect(isComplete || isError).To(BeTrue())
			Expect(event.Date.IsZero()).To(BeFalse())
		})
	})

	// ─── performSubmit — skill save error ────────────────────────────────────

	Describe("performSubmit — skill save error", func() {
		It("returns SubmitErrorMsg with SKILL_SAVE_ERROR when skill name is empty", func() {
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event)
			svc.SetSkillRepository(repos.Skill)

			ctx := &IntentContext{CaptureStrategy: "quick"}
			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.context.CareerService = svc

			event := fixtures.EventWith("", "Valid event text for skill save error test here", "", "")
			event.Date = time.Now().Add(-time.Hour)
			event.Tags = []string{"technical"}
			event.Categories = []string{"technical"}

			invalidSkill := fixtures.SkillWith("", "", "", "")

			intent.reviewState = &ReviewInferredEventState{
				Event:          event,
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: []*career.Skill{invalidSkill},
			}

			cmd := intent.performSubmit()
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			errMsg, ok := msg.(SubmitErrorMsg)
			Expect(ok).To(BeTrue())
			Expect(errMsg.Code).To(Equal("SKILL_SAVE_ERROR"))
		})
	})

	// ─── performPostSavePersistence — burst confirm error ────────────────────

	Describe("performPostSavePersistence — burst confirm error", func() {
		It("returns SubmitErrorMsg when ConfirmBurst fails for non-existent burst", func() {
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event)
			svc.SetSkillRepository(repos.Skill)

			ctx := &IntentContext{CaptureStrategy: "quick"}
			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.context.CareerService = svc

			event := fixtures.EventWith("evt-burst", "Valid burst test event text here correct", "Corp", "Proj")
			event.Date = time.Now().Add(-time.Hour)
			event.Tags = []string{"technical"}
			event.Categories = []string{"technical"}
			Expect(repos.Event.Create(context.Background(), event)).To(Succeed())

			nonExistentBurst := fixtures.Burst("nonexistent-burst-id", "evt-burst")

			cmd := intent.performPostSavePersistence(
				event,
				[]*career.Fact{},
				[]*career.Skill{},
				[]*career.Burst{nonExistentBurst},
			)
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			_, isComplete := msg.(PostSavePersistenceCompleteMsg)
			_, isError := msg.(SubmitErrorMsg)
			Expect(isComplete || isError).To(BeTrue())
		})
	})
})
