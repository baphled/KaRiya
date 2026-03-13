package burst_management

import (
	"context"
	"errors"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/testutil/mocks"
	burstviews "github.com/baphled/kariya/internal/tui/views/burst"
	skillviews "github.com/baphled/kariya/internal/tui/views/skill"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/terminal"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type unknownViewResult struct{}

func (unknownViewResult) Type() widgets.ViewResultType     { return "unknown" }
func (unknownViewResult) Data() interface{}                { return nil }
func (unknownViewResult) Metadata() map[string]interface{} { return nil }
func (unknownViewResult) WithMetadata(string, interface{}) widgets.ViewResult {
	return unknownViewResult{}
}

type testViewNoConfig struct{}

func (testViewNoConfig) RenderContent() string { return "content" }
func (testViewNoConfig) HelpText() string      { return "help" }
func (testViewNoConfig) Init() tea.Cmd         { return nil }
func (testViewNoConfig) Update(tea.Msg) (tea.Cmd, widgets.ViewResult) {
	return nil, nil
}

var _ = Describe("Handlers", func() {
	Describe("noopCmd", func() {
		It("returns nil tea.Msg", func() {
			msg := noopCmd()
			Expect(msg).To(BeNil())
		})
	})

	Describe("handleViewResult", func() {
		var intent *Intent
		BeforeEach(func() {
			ctx := &IntentValidator{Bursts: nil}
			ctx.Validate()
			var err error
			intent, err = NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("dispatches ResultSubmit to HandleSubmit and returns nil", func() {
			result := &widgets.SubmitViewResult{FormData: map[string]interface{}{}}
			cmd := intent.handleViewResult(result)
			Expect(cmd).To(BeNil())
		})

		It("dispatches ResultError to HandleError and returns nil", func() {
			err := errors.New("test error")
			result := &widgets.ErrorViewResult{Err: err, Message: "fail"}
			cmd := intent.handleViewResult(result)
			Expect(cmd).To(BeNil())
			Expect(intent.deleteError).To(Equal(err))
			Expect(intent.feedbackModal).NotTo(BeNil())
		})

		It("dispatches ResultCancel to HandleCancel and returns nil", func() {
			result := &widgets.CancelViewResult{}
			cmd := intent.handleViewResult(result)
			Expect(cmd).To(BeNil())
		})

		It("dispatches ResultNavigate to HandleNavigate and returns nil", func() {
			result := &widgets.NavigateViewResult{}
			cmd := intent.handleViewResult(result)
			Expect(cmd).To(BeNil())
		})

		It("returns nil for unknown result type", func() {
			cmd := intent.handleViewResult(unknownViewResult{})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("HandleError", func() {
		var intent *Intent
		BeforeEach(func() {
			ctx := &IntentValidator{Bursts: nil}
			ctx.Validate()
			var err error
			intent, err = NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("sets deleteError and feedbackModal for valid ErrorViewResult", func() {
			err := errors.New("error")
			result := &widgets.ErrorViewResult{Err: err, Message: "fail"}
			cmd := intent.HandleError(result)
			Expect(cmd).To(BeNil())
			Expect(intent.deleteError).To(Equal(err))
			Expect(intent.feedbackModal).NotTo(BeNil())
		})

		It("does not set deleteError for non-error result type", func() {
			cmd := intent.HandleError(unknownViewResult{})
			Expect(cmd).To(BeNil())
			Expect(intent.deleteError).ToNot(HaveOccurred())
		})

		It("does not set deleteError when error is nil", func() {
			result := &widgets.ErrorViewResult{Err: nil, Message: "no error"}
			cmd := intent.HandleError(result)
			Expect(cmd).To(BeNil())
			Expect(intent.deleteError).ToNot(HaveOccurred())
		})
	})
})

var _ = Describe("Internal helpers", func() {
	var intent *Intent

	setupIntent := func(ctx *IntentValidator) {
		ctx.Validate()
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		termInfo := terminal.NewInfo()
		termInfo.Width = 120
		termInfo.Height = 40
		termInfo.IsValid = true
		intent.UpdateTerminalInfo(termInfo)
		intent.Init()
	}

	testTheme := themes.NewDefaultTheme()

	Describe("loadBurstEvents", func() {
		It("returns empty slice when burst is nil", func() {
			setupIntent(&IntentValidator{})
			events := intent.loadBurstEvents(context.Background(), nil)
			Expect(events).To(BeEmpty())
		})

		It("returns empty slice when service is nil", func() {
			setupIntent(&IntentValidator{})
			burst := fixtures.Burst("load-ev-1")
			events := intent.loadBurstEvents(context.Background(), burst)
			Expect(events).To(BeEmpty())
		})

		It("loads events matching burst event IDs", func() {
			event1 := fixtures.Event("ev-1")
			event2 := fixtures.Event("ev-2")
			mockSvc := mocks.NewBurstServiceMock().SetEvents([]*career.Event{event1, event2})
			setupIntent(&IntentValidator{Service: mockSvc})
			burst := fixtures.Burst("load-ev-2", "ev-1", "ev-2")
			events := intent.loadBurstEvents(context.Background(), burst)
			Expect(events).To(HaveLen(2))
		})

		It("skips events that are not found", func() {
			event1 := fixtures.Event("ev-1")
			mockSvc := mocks.NewBurstServiceMock().SetEvents([]*career.Event{event1})
			setupIntent(&IntentValidator{Service: mockSvc})
			burst := fixtures.Burst("load-ev-3", "ev-1", "ev-missing")
			events := intent.loadBurstEvents(context.Background(), burst)
			Expect(events).To(HaveLen(1))
		})
	})

	Describe("inferSkillsFromBurst", func() {
		It("returns error msg when burst is nil", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.inferSkillsFromBurst(nil)
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			errMsg, ok := msg.(SkillSuggestionsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(errMsg.Error).To(MatchError("no burst provided"))
		})

		It("returns error msg when service is nil", func() {
			mockSvc := mocks.NewBurstServiceMock()
			setupIntent(&IntentValidator{Service: mockSvc})
			burst := fixtures.Burst("infer-1")
			cmd := intent.inferSkillsFromBurst(burst)
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			errMsg, ok := msg.(SkillSuggestionsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(errMsg.Error).To(MatchError("skill inference service not available"))
		})

		It("returns error msg when no events found for burst", func() {
			mockSvc := mocks.NewBurstServiceMock()
			setupIntent(&IntentValidator{
				Service:               mockSvc,
				SkillInferenceService: &stubSkillInferenceService{},
			})
			burst := fixtures.Burst("infer-2", "missing-1", "missing-2")
			cmd := intent.inferSkillsFromBurst(burst)
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			errMsg, ok := msg.(SkillSuggestionsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(errMsg.Error).To(MatchError("no events found for burst"))
		})

		It("returns error msg when context is cancelled before inner call", func() {
			ev1 := fixtures.Event("ev-1")
			ev2 := fixtures.Event("ev-2")
			mockSvc := mocks.NewBurstServiceMock().SetEvents([]*career.Event{ev1, ev2})
			setupIntent(&IntentValidator{
				Service:               mockSvc,
				SkillInferenceService: &stubSkillInferenceService{},
			})
			burst := fixtures.Burst("infer-3", "ev-1", "ev-2")
			cmd := intent.inferSkillsFromBurst(burst)
			intent.cancelFunc()
			msg := cmd()
			errMsg, ok := msg.(SkillSuggestionsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(errMsg.Error).To(MatchError(context.Canceled))
		})

		It("returns error msg when inference fails", func() {
			ev1 := fixtures.Event("ev-1")
			ev2 := fixtures.Event("ev-2")
			mockSvc := mocks.NewBurstServiceMock().SetEvents([]*career.Event{ev1, ev2})
			setupIntent(&IntentValidator{
				Service:               mockSvc,
				SkillInferenceService: &stubSkillInferenceService{inferErr: errors.New("inference boom")},
			})
			burst := fixtures.Burst("infer-4", "ev-1", "ev-2")
			cmd := intent.inferSkillsFromBurst(burst)
			msg := cmd()
			errMsg, ok := msg.(SkillSuggestionsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(errMsg.Error.Error()).To(ContainSubstring("skill detection failed"))
		})

		It("returns suggestions on success", func() {
			ev1 := fixtures.Event("ev-1")
			ev2 := fixtures.Event("ev-2")
			mockSvc := mocks.NewBurstServiceMock().SetEvents([]*career.Event{ev1, ev2})
			setupIntent(&IntentValidator{
				Service: mockSvc,
				SkillInferenceService: &stubSkillInferenceService{
					result: &skillinference.InferenceResult{
						Suggestions:        []skillinference.SkillSuggestion{{Name: "Go", Category: "backend", Confidence: 0.9}},
						ExistingSkillNames: []string{"Python"},
					},
				},
			})
			burst := fixtures.Burst("infer-5", "ev-1", "ev-2")
			cmd := intent.inferSkillsFromBurst(burst)
			msg := cmd()
			loadedMsg, ok := msg.(SkillSuggestionsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(loadedMsg.Suggestions).To(HaveLen(1))
			Expect(loadedMsg.ExistingSkillNames).To(ConsistOf("Python"))
		})

		It("cancels previous operation", func() {
			setupIntent(&IntentValidator{Service: mocks.NewBurstServiceMock()})
			called := false
			intent.cancelFunc = func() { called = true }
			burst := fixtures.Burst("infer-cancel")
			intent.inferSkillsFromBurst(burst)
			Expect(called).To(BeTrue())
		})
	})

	Describe("showSkills", func() {
		It("returns nil when selectedBurst is nil", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.showSkills()
			Expect(cmd).To(BeNil())
		})

		It("returns empty skills when SkillRepository is nil", func() {
			setupIntent(&IntentValidator{})
			intent.selectedBurst = fixtures.Burst("skill-1")
			cmd := intent.showSkills()
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			loadedMsg, ok := msg.(BurstSkillsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(loadedMsg.Skills).To(BeEmpty())
		})

		It("loads and deduplicates skills from repository", func() {
			mockRepo := &stubSkillRepository{
				skillsByEvent: map[string][]*career.Skill{
					"ev-1": {fixtures.SkillWith("s1", "Go", "backend", "advanced")},
					"ev-2": {
						fixtures.SkillWith("s1", "Go", "backend", "advanced"),
						fixtures.SkillWith("s2", "Docker", "devops", "intermediate"),
					},
				},
			}
			setupIntent(&IntentValidator{SkillRepository: mockRepo})
			intent.selectedBurst = fixtures.Burst("skill-2", "ev-1", "ev-2")
			cmd := intent.showSkills()
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			loadedMsg, ok := msg.(BurstSkillsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(loadedMsg.Skills).To(HaveLen(2))
		})

		It("skips events with repository errors", func() {
			mockRepo := &stubSkillRepository{
				skillsByEvent: map[string][]*career.Skill{
					"ev-1": {fixtures.SkillWith("s1", "Go", "backend", "advanced")},
				},
				errorEvents: map[string]bool{"ev-2": true},
			}
			setupIntent(&IntentValidator{SkillRepository: mockRepo})
			intent.selectedBurst = fixtures.Burst("skill-3", "ev-1", "ev-2")
			cmd := intent.showSkills()
			msg := cmd()
			loadedMsg, ok := msg.(BurstSkillsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(loadedMsg.Skills).To(HaveLen(1))
		})
	})

	Describe("saveSkillFromSuggestion", func() {
		It("returns nil when service is nil", func() {
			setupIntent(&IntentValidator{})
			suggestion := skillinference.SkillSuggestion{Name: "Go"}
			cmd := intent.saveSkillFromSuggestion(suggestion)
			Expect(cmd).To(BeNil())
		})

		It("returns async command that creates skill", func() {
			setupIntent(&IntentValidator{
				SkillInferenceService: &stubSkillInferenceService{},
			})
			suggestion := skillinference.SkillSuggestion{Name: "Go"}
			cmd := intent.saveSkillFromSuggestion(suggestion)
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			createdMsg, ok := msg.(SkillsCreatedMsg)
			Expect(ok).To(BeTrue())
			Expect(createdMsg.Error).ToNot(HaveOccurred())
		})

		It("returns error msg when CreateSkillsFromSuggestions fails", func() {
			setupIntent(&IntentValidator{
				SkillInferenceService: &stubSkillInferenceService{
					createErr: errors.New("create boom"),
				},
			})
			suggestion := skillinference.SkillSuggestion{Name: "Go"}
			cmd := intent.saveSkillFromSuggestion(suggestion)
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			createdMsg, ok := msg.(SkillsCreatedMsg)
			Expect(ok).To(BeTrue())
			Expect(createdMsg.Error).To(HaveOccurred())
		})
	})

	Describe("cancelPreviousOperation", func() {
		It("does nothing when cancelFunc is nil", func() {
			setupIntent(&IntentValidator{})
			intent.cancelFunc = nil
			Expect(func() { intent.cancelPreviousOperation() }).NotTo(Panic())
		})

		It("calls cancelFunc when set", func() {
			setupIntent(&IntentValidator{})
			called := false
			intent.cancelFunc = func() { called = true }
			intent.cancelPreviousOperation()
			Expect(called).To(BeTrue())
		})
	})

	Describe("createBurstFromSuggestion", func() {
		It("creates burst without repository", func() {
			setupIntent(&IntentValidator{})
			suggestion := burstfact.BurstSuggestion{
				Name:        "Test Burst",
				Description: "A test",
				EventIDs:    []string{"ev-1", "ev-2"},
			}
			burst := intent.createBurstFromSuggestion(suggestion)
			Expect(burst).NotTo(BeNil())
			Expect(burst.Name).To(Equal("Test Burst"))
			Expect(intent.filteredBursts).To(ContainElement(burst))
		})

		It("saves to repository on success", func() {
			mockRepo := mocks.NewBurstRepositoryMock()
			setupIntent(&IntentValidator{BurstRepository: mockRepo})
			suggestion := burstfact.BurstSuggestion{
				Name:        "Saved Burst",
				Description: "Saved",
				EventIDs:    []string{"ev-1", "ev-2"},
			}
			burst := intent.createBurstFromSuggestion(suggestion)
			Expect(burst).NotTo(BeNil())
		})

		It("shows error and returns nil on repository failure", func() {
			mockRepo := mocks.NewBurstRepositoryMock().SetCreateError(errors.New("create fail"))
			setupIntent(&IntentValidator{BurstRepository: mockRepo})
			suggestion := burstfact.BurstSuggestion{
				Name:        "Fail Burst",
				Description: "Fail",
				EventIDs:    []string{"ev-1", "ev-2"},
			}
			burst := intent.createBurstFromSuggestion(suggestion)
			Expect(burst).To(BeNil())
			Expect(intent.feedbackModal).NotTo(BeNil())
		})
	})

	Describe("transitionToView", func() {
		It("sets activeView", func() {
			setupIntent(&IntentValidator{})
			view := burstviews.NewList(nil)
			intent.transitionToView(view)
			Expect(intent.activeView).To(Equal(view))
		})

		It("propagates terminal info to view", func() {
			setupIntent(&IntentValidator{})
			view := burstviews.NewList(nil)
			intent.transitionToView(view)
			Expect(intent.activeView).NotTo(BeNil())
		})

		It("handles nil terminal info gracefully", func() {
			ctx := &IntentValidator{}
			ctx.Validate()
			var err error
			intent, err = NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			view := burstviews.NewList(nil)
			intent.transitionToView(view)
			Expect(intent.activeView).To(Equal(view))
		})

		It("handles view without configurer interface", func() {
			setupIntent(&IntentValidator{})
			view := testViewNoConfig{}
			intent.transitionToView(view)
			Expect(intent.activeView).To(Equal(view))
		})
	})

	Describe("extractFactsForBurst closure execution", func() {
		It("returns error when context is cancelled", func() {
			mockSvc := mocks.NewBurstServiceMock()
			setupIntent(&IntentValidator{Service: mockSvc})
			burst := fixtures.Burst("ef-1")
			cmd := intent.extractFactsForBurst(burst)
			intent.cancelFunc()
			msg := cmd()
			completeMsg, ok := msg.(FactExtractionCompleteMsg)
			Expect(ok).To(BeTrue())
			Expect(completeMsg.Error).To(MatchError(context.Canceled))
		})

		It("returns error when service is nil", func() {
			setupIntent(&IntentValidator{})
			burst := fixtures.Burst("ef-2")
			cmd := intent.extractFactsForBurst(burst)
			msg := cmd()
			completeMsg, ok := msg.(FactExtractionCompleteMsg)
			Expect(ok).To(BeTrue())
			Expect(completeMsg.Error).To(MatchError("service not available"))
		})

		It("returns error when extraction fails", func() {
			mockSvc := mocks.NewBurstServiceMock().SetExtractError(errors.New("extract boom"))
			setupIntent(&IntentValidator{Service: mockSvc})
			burst := fixtures.Burst("ef-3")
			cmd := intent.extractFactsForBurst(burst)
			msg := cmd()
			completeMsg, ok := msg.(FactExtractionCompleteMsg)
			Expect(ok).To(BeTrue())
			Expect(completeMsg.Error).To(MatchError("extract boom"))
		})

		It("saves extracted facts and returns them", func() {
			mockSvc := mocks.NewBurstServiceMock().SetExtractedFacts([]career.Fact{
				*fixtures.FactWith("f1", "Fact one"),
				*fixtures.FactWith("f2", "Fact two"),
			})
			setupIntent(&IntentValidator{Service: mockSvc})
			burst := fixtures.Burst("ef-4")
			cmd := intent.extractFactsForBurst(burst)
			msg := cmd()
			completeMsg, ok := msg.(FactExtractionCompleteMsg)
			Expect(ok).To(BeTrue())
			Expect(completeMsg.Error).ToNot(HaveOccurred())
			Expect(completeMsg.Facts).To(HaveLen(2))
		})

		It("skips facts that fail to save", func() {
			mockSvc := mocks.NewBurstServiceMock().
				SetExtractedFacts([]career.Fact{*fixtures.FactWith("f1", "Fact one")}).
				SetSaveFactError(errors.New("save fail"))
			setupIntent(&IntentValidator{Service: mockSvc})
			burst := fixtures.Burst("ef-5")
			cmd := intent.extractFactsForBurst(burst)
			msg := cmd()
			completeMsg, ok := msg.(FactExtractionCompleteMsg)
			Expect(ok).To(BeTrue())
			Expect(completeMsg.Facts).To(BeEmpty())
		})

		It("returns error for nil burst", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.extractFactsForBurst(nil)
			msg := cmd()
			completeMsg, ok := msg.(FactExtractionCompleteMsg)
			Expect(ok).To(BeTrue())
			Expect(completeMsg.Error).To(MatchError("no burst provided"))
		})

		It("checks context cancellation after extraction", func() {
			mockSvc := mocks.NewBurstServiceMock().SetExtractedFacts([]career.Fact{
				*fixtures.FactWith("f1", "Fact one"),
			})
			setupIntent(&IntentValidator{Service: mockSvc})
			burst := fixtures.Burst("ef-6")
			cmd := intent.extractFactsForBurst(burst)
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			completeMsg, ok := msg.(FactExtractionCompleteMsg)
			Expect(ok).To(BeTrue())
			Expect(completeMsg.Error).ToNot(HaveOccurred())
		})
	})

	Describe("startSkillInference", func() {
		It("shows error when selectedBurst is nil", func() {
			setupIntent(&IntentValidator{})
			intent.selectedBurst = nil
			cmd := intent.startSkillInference()
			Expect(cmd).To(BeNil())
			Expect(intent.feedbackModal).NotTo(BeNil())
			Expect(intent.state).To(Equal(StateList))
		})

		It("shows error when SkillInferenceService is nil", func() {
			setupIntent(&IntentValidator{})
			intent.selectedBurst = fixtures.Burst("si-1")
			cmd := intent.startSkillInference()
			Expect(cmd).To(BeNil())
			Expect(intent.feedbackModal).NotTo(BeNil())
			Expect(intent.state).To(Equal(StateList))
		})

		It("starts inference and returns batch command", func() {
			mockSvc := mocks.NewBurstServiceMock()
			setupIntent(&IntentValidator{
				Service:               mockSvc,
				SkillInferenceService: &stubSkillInferenceService{},
			})
			intent.selectedBurst = fixtures.Burst("si-2")
			cmd := intent.startSkillInference()
			Expect(cmd).NotTo(BeNil())
			Expect(intent.state).To(Equal(StateInferringSkills))
			Expect(intent.inferringSkills).To(BeTrue())
			Expect(intent.loadingModal).NotTo(BeNil())
		})
	})

	Describe("startBurstDetection", func() {
		It("returns nil when service is nil", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.startBurstDetection()
			Expect(cmd).To(BeNil())
		})

		It("starts detection and returns batch command", func() {
			ev := fixtures.Event("ev-1")
			mockSvc := mocks.NewBurstServiceMock().SetEvents([]*career.Event{ev})
			setupIntent(&IntentValidator{Service: mockSvc})
			cmd := intent.startBurstDetection()
			Expect(cmd).NotTo(BeNil())
			Expect(intent.suggestionsLoading).To(BeTrue())
			Expect(intent.loadingModal).NotTo(BeNil())
		})
	})

	Describe("handleEventsModalUpdate with nil selectedBurst", func() {
		It("returns noopCmd when Esc is pressed and selectedBurst is nil", func() {
			setupIntent(&IntentValidator{})
			intent.eventsModal = burstviews.NewEvents("b1", "test", nil, testTheme)
			cmd := intent.handleEventsModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.eventsModal).To(BeNil())
		})
	})

	Describe("handleSkillsModalUpdate with nil selectedBurst", func() {
		It("returns noopCmd when Esc is pressed and selectedBurst is nil", func() {
			setupIntent(&IntentValidator{})
			intent.skillsModal = burstviews.NewSkills("b1", "test", nil, testTheme)
			cmd := intent.handleSkillsModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.skillsModal).To(BeNil())
		})
	})

	Describe("handleFactsModalUpdate with nil selectedBurst", func() {
		It("returns noopCmd when Esc is pressed and selectedBurst is nil", func() {
			setupIntent(&IntentValidator{})
			intent.factsModal = burstviews.NewFacts("b1", "test", nil, testTheme)
			cmd := intent.handleFactsModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.factsModal).To(BeNil())
		})
	})

	Describe("handleEditModalUpdate", func() {
		It("returns nil when editModal is nil", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.handleEditModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
		})

		It("closes and returns noop when edit modal is cancelled", func() {
			setupIntent(&IntentValidator{})
			intent.selectedBurst = fixtures.Burst("edit-cancel")
			intent.openEditModal(intent.selectedBurst)
			cmd := intent.handleEditModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.editModal).To(BeNil())
			Expect(intent.state).To(Equal(StateList))
		})

		It("returns edit message when form completes", func() {
			setupIntent(&IntentValidator{})
			intent.selectedBurst = fixtures.Burst("edit-complete")
			intent.openEditModal(intent.selectedBurst)
			cmd := intent.handleEditModalUpdate(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			editMsg, ok := msg.(EditBurstMsg)
			Expect(ok).To(BeTrue())
			Expect(editMsg.BurstID).To(Equal(intent.selectedBurst.ID))
		})
	})

	Describe("handleSuggestionModalClosed", func() {
		It("returns to list state when cancelled with no accepted suggestions", func() {
			setupIntent(&IntentValidator{})
			intent.state = StateSuggestionReview
			intent.suggestionModal = burstviews.NewSuggestionReview(
				display.BurstSuggestionsFromDomain([]burstfact.BurstSuggestion{
					{Name: "SugA", Description: "A", EventIDs: []string{"ev-1", "ev-2"}},
				}),
				testTheme,
			)
			cmd := intent.handleSuggestionModalClosed(noopCmd)
			Expect(cmd).NotTo(BeNil())
			Expect(intent.state).To(Equal(StateList))
		})

		It("closes suggestion modal on Esc via Update", func() {
			setupIntent(&IntentValidator{})
			intent.state = StateSuggestionReview
			intent.suggestionModal = burstviews.NewSuggestionReview(
				display.BurstSuggestionsFromDomain([]burstfact.BurstSuggestion{{
					Name: "SugB", Description: "B", EventIDs: []string{"ev-1", "ev-2"},
				}}),
				testTheme,
			)
			intent.suggestionModal.Show()
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.suggestionModal).To(BeNil())
			Expect(intent.state).To(Equal(StateList))
		})
	})

	Describe("removeBurstFromSlice", func() {
		It("removes matching burst by ID", func() {
			setupIntent(&IntentValidator{})
			b1 := fixtures.Burst("b1")
			b2 := fixtures.Burst("b2")
			result := intent.removeBurstFromSlice([]*career.Burst{b1, b2}, "b1")
			Expect(result).To(HaveLen(1))
			Expect(result[0].ID).To(Equal("b2"))
		})

		It("returns same slice when ID not found", func() {
			setupIntent(&IntentValidator{})
			b1 := fixtures.Burst("b1")
			result := intent.removeBurstFromSlice([]*career.Burst{b1}, "missing")
			Expect(result).To(HaveLen(1))
		})
	})

	Describe("RefreshData", func() {
		It("reloads bursts from repository", func() {
			mockRepo := mocks.NewBurstRepositoryMock()
			setupIntent(&IntentValidator{BurstRepository: mockRepo})
			err := intent.RefreshData()
			Expect(err).To(BeNil())
		})
	})

	Describe("startFactExtractionForBursts", func() {
		It("returns nil and sets StateList when no bursts", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.startFactExtractionForBursts(nil)
			Expect(cmd).To(BeNil())
			Expect(intent.state).To(Equal(StateList))
		})

		It("extracts facts for multiple bursts", func() {
			mockSvc := mocks.NewBurstServiceMock()
			setupIntent(&IntentValidator{Service: mockSvc})
			b1 := fixtures.Burst("b1")
			b2 := fixtures.Burst("b2")
			cmd := intent.startFactExtractionForBursts([]*career.Burst{b1, b2})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.extractingFacts).To(BeTrue())
			Expect(intent.state).To(Equal(StateExtractingFacts))
		})
	})

	Describe("handleConfirmModalUpdate", func() {
		It("returns nil when confirmModal is nil", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.handleConfirmModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
		})

		It("handles Esc to dismiss confirm modal", func() {
			setupIntent(&IntentValidator{})
			intent.confirmModal = feedback.NewConfirmModal("Test", "Confirm?")
			intent.confirmModal.Init()
			intent.state = StateConfirm
			cmd := intent.handleConfirmModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("handleDeleteModalUpdate", func() {
		It("returns nil when deleteModal is nil", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.handleDeleteModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("handleLoadingModalUpdate", func() {
		It("returns nil when loadingModal is nil", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.handleLoadingModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
		})

		It("cancels on Esc", func() {
			setupIntent(&IntentValidator{})
			intent.loadingModal = feedback.NewLoadingModal("Loading...", true)
			cancelled := false
			intent.cancelFunc = func() { cancelled = true }
			cmd := intent.handleLoadingModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
			Expect(cancelled).To(BeTrue())
		})

		It("consumes non-Esc keys", func() {
			setupIntent(&IntentValidator{})
			intent.loadingModal = feedback.NewLoadingModal("Loading...", true)
			cmd := intent.handleLoadingModalUpdate(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("handleSuggestionModalUpdate", func() {
		It("returns nil when suggestionModal is nil", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.handleSuggestionModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("handleActionData", func() {
		It("handles unknown action", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.handleActionData(map[string]interface{}{"action": "unknown"})
			Expect(cmd).To(BeNil())
		})

		It("handles view action with invalid burst type", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.handleActionData(map[string]interface{}{"action": "view", "burst": "not-a-burst"})
			Expect(cmd).To(BeNil())
		})

		It("handles view_events action with invalid burst type", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.handleActionData(map[string]interface{}{"action": "view_events", "burst": nil})
			Expect(cmd).To(BeNil())
		})

		It("handles view_facts action with invalid burst type", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.handleActionData(map[string]interface{}{"action": "view_facts", "burst": nil})
			Expect(cmd).To(BeNil())
		})

		It("handles confirm action with invalid burst type", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.handleActionData(map[string]interface{}{"action": "confirm", "burst": nil})
			Expect(cmd).To(BeNil())
		})

		It("handles edit action with invalid burst type", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.handleActionData(map[string]interface{}{"action": "edit", "burst": nil})
			Expect(cmd).To(BeNil())
		})

		It("handles delete action with invalid burst type", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.handleActionData(map[string]interface{}{"action": "delete", "burst": nil})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("rebuildModalRegistry", func() {
		It("rebuilds with all modal types", func() {
			setupIntent(&IntentValidator{})
			intent.feedbackModal = feedback.NewErrorModal("Err", "error")
			intent.rebuildModalRegistry()
			Expect(intent.modalRegistry).NotTo(BeNil())
		})
	})

	Describe("showEvents", func() {
		It("returns cmd for loading events", func() {
			ev1 := fixtures.Event("ev-1")
			mockSvc := mocks.NewBurstServiceMock().SetEvents([]*career.Event{ev1})
			setupIntent(&IntentValidator{Service: mockSvc})
			intent.selectedBurst = fixtures.Burst("ev-modal", "ev-1", "ev-2")
			cmd := intent.showEvents()
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			evMsg, ok := msg.(BurstEventsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(evMsg.Events).To(HaveLen(1))
		})
	})

	Describe("showFacts", func() {
		It("returns cmd for loading facts", func() {
			mockSvc := mocks.NewBurstServiceMock()
			setupIntent(&IntentValidator{Service: mockSvc})
			intent.selectedBurst = fixtures.Burst("fact-modal")
			cmd := intent.showFacts()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("openDeleteModal", func() {
		It("creates delete modal for burst", func() {
			setupIntent(&IntentValidator{})
			burst := fixtures.Burst("del-1")
			intent.openDeleteModal(burst)
			Expect(intent.deleteModal).NotTo(BeNil())
			Expect(intent.selectedBurst).To(Equal(burst))
		})
	})

	Describe("NewIntent", func() {
		It("returns error for nil context", func() {
			_, err := NewIntent(nil)
			Expect(err).To(MatchError(ErrInvalidContext))
		})
	})

	Describe("Init", func() {
		It("initializes activeView", func() {
			setupIntent(&IntentValidator{})
			Expect(intent.activeView).NotTo(BeNil())
		})
	})

	Describe("openSuggestionEventsModal", func() {
		It("returns noopCmd when skillSuggestionModal is nil", func() {
			setupIntent(&IntentValidator{})
			intent.selectedBurst = fixtures.Burst("ose-1")
			intent.skillSuggestionModal = nil
			cmd := intent.openSuggestionEventsModal()
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			Expect(msg).To(BeNil())
		})
	})

	Describe("getUnassignedEventIDs", func() {
		It("returns unassigned event IDs", func() {
			ev1 := fixtures.Event("ua-1")
			ev2 := fixtures.Event("ua-2")
			mockSvc := mocks.NewBurstServiceMock().SetEvents([]*career.Event{ev1, ev2})
			setupIntent(&IntentValidator{Service: mockSvc})
			ids, err := intent.getUnassignedEventIDs(
				context.Background(),
				mockSvc,
				[]*career.Burst{},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(ids).To(HaveLen(2))
		})
	})

	Describe("handleSkillSuggestionModalUpdate", func() {
		It("returns nil when skillSuggestionModal is nil", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.handleSkillSuggestionModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("handleSuggestionEventsModalUpdate", func() {
		It("returns nil when suggestionEventsModal is nil", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.handleSuggestionEventsModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("handleEventsModalUpdate", func() {
		It("returns nil when eventsModal is nil", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.handleEventsModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
		})

		It("closes events modal and returns detail modal on Esc", func() {
			ev := fixtures.Event("ev-modal-1")
			mockSvc := mocks.NewBurstServiceMock().SetEvents([]*career.Event{ev})
			setupIntent(&IntentValidator{Service: mockSvc})
			intent.selectedBurst = fixtures.Burst("events-1", "ev-modal-1", "ev-2")
			loaded := intent.Update(BurstEventsLoadedMsg{Events: []*career.Event{ev}})
			Expect(loaded).To(BeNil())
			cmd := intent.handleEventsModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
			Expect(intent.eventsModal).To(BeNil())
			Expect(intent.detailModal).NotTo(BeNil())
		})
	})

	Describe("handleDetailModalUpdate", func() {
		It("returns nil when detailModal is nil", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.handleDetailModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("confirmBurst", func() {
		It("returns nil when selectedBurst is nil", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.confirmBurst()
			Expect(cmd).To(BeNil())
		})

		It("returns async command that confirms burst with service", func() {
			mockSvc := mocks.NewBurstServiceMock()
			setupIntent(&IntentValidator{Service: mockSvc})
			intent.selectedBurst = fixtures.Burst("confirm-1")
			cmd := intent.confirmBurst()
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			confirmedMsg, ok := msg.(BurstConfirmedMsg)
			Expect(ok).To(BeTrue())
			Expect(confirmedMsg.Error).ToNot(HaveOccurred())
			Expect(mockSvc.GetConfirmCallCount()).To(Equal(1))
		})
	})

	Describe("showConfirmBurstModal", func() {
		It("returns nil when selectedBurst is nil", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.showConfirmBurstModal()
			Expect(cmd).To(BeNil())
		})

		It("returns async command that loads facts", func() {
			setupIntent(&IntentValidator{})
			intent.selectedBurst = fixtures.Burst("scm-1")
			cmd := intent.showConfirmBurstModal()
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			factsMsg, ok := msg.(ConfirmBurstFactsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(factsMsg.Facts).To(BeEmpty())
		})
	})

	Describe("handleModalUpdates routing", func() {
		It("returns nil when no modal is active", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.handleModalUpdates(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("transitionToView with theme manager", func() {
		It("propagates theme manager to configurer view", func() {
			setupIntent(&IntentValidator{})
			tm := themes.NewThemeManager()
			intent.SetThemeManager(tm)
			view := burstviews.NewList(nil)
			intent.transitionToView(view)
			Expect(intent.activeView).To(Equal(view))
		})
	})

	Describe("handleEventsModalUpdate with non-Esc key", func() {
		It("forwards non-Esc key to modal update", func() {
			setupIntent(&IntentValidator{})
			intent.eventsModal = burstviews.NewEvents("b1", "test", nil, testTheme)
			intent.eventsModal.Show()
			intent.handleEventsModalUpdate(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(intent.eventsModal.IsVisible()).To(BeTrue())
		})
	})

	Describe("handleSkillsModalUpdate with non-Esc key", func() {
		It("forwards non-Esc key to modal update", func() {
			setupIntent(&IntentValidator{})
			intent.skillsModal = burstviews.NewSkills("b1", "test", nil, testTheme)
			intent.skillsModal.Show()
			intent.handleSkillsModalUpdate(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(intent.skillsModal.IsVisible()).To(BeTrue())
		})
	})

	Describe("handleFactsModalUpdate with non-Esc key", func() {
		It("forwards non-Esc key to modal update", func() {
			setupIntent(&IntentValidator{})
			intent.factsModal = burstviews.NewFacts("b1", "test", nil, testTheme)
			intent.factsModal.Show()
			intent.handleFactsModalUpdate(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(intent.factsModal.IsVisible()).To(BeTrue())
		})
	})

	Describe("handleSuggestionEventsModalUpdate with non-Esc key", func() {
		It("forwards non-Esc key to modal update", func() {
			setupIntent(&IntentValidator{})
			ev := fixtures.Event("semu-1")
			intent.suggestionEventsModal = skillviews.NewEvents("s1", "Go", []*career.Event{ev}, testTheme)
			intent.suggestionEventsModal.Show()
			intent.handleSuggestionEventsModalUpdate(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(intent.suggestionEventsModal.IsVisible()).To(BeTrue())
		})
	})

	Describe("handleLoadingModalUpdate spinner tick", func() {
		It("forwards spinner tick msg to loading modal", func() {
			setupIntent(&IntentValidator{})
			intent.loadingModal = feedback.NewLoadingModal("Loading...", true)
			initCmd := intent.loadingModal.Init()
			Expect(initCmd).NotTo(BeNil())
			tickMsg := initCmd()
			cmd := intent.handleLoadingModalUpdate(tickMsg)
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("deleteBurst", func() {
		It("returns nil when burst is nil", func() {
			setupIntent(&IntentValidator{})
			cmd := intent.deleteBurst(nil)
			Expect(cmd).To(BeNil())
		})

		It("returns async command that deletes burst", func() {
			mockRepo := mocks.NewBurstRepositoryMock()
			burst := fixtures.Burst("del-burst")
			mockRepo.AddBurst(burst)
			setupIntent(&IntentValidator{
				Bursts:          []*career.Burst{burst},
				BurstRepository: mockRepo,
			})
			intent.selectedBurst = burst
			cmd := intent.deleteBurst(burst)
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			deletedMsg, ok := msg.(BurstDeletedMsg)
			Expect(ok).To(BeTrue())
			Expect(deletedMsg.Error).ToNot(HaveOccurred())
			Expect(deletedMsg.BurstID).To(Equal(burst.ID))
		})

		It("returns error msg on repository failure", func() {
			mockRepo := mocks.NewBurstRepositoryMock()
			mockRepo.SetDeleteError(errors.New("delete fail"))
			burst := fixtures.Burst("del-fail")
			setupIntent(&IntentValidator{
				Bursts:          []*career.Burst{burst},
				BurstRepository: mockRepo,
			})
			cmd := intent.deleteBurst(burst)
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			deletedMsg, ok := msg.(BurstDeletedMsg)
			Expect(ok).To(BeTrue())
			Expect(deletedMsg.Error).To(HaveOccurred())
		})
	})

})

type stubSkillInferenceService struct {
	result    *skillinference.InferenceResult
	inferErr  error
	createErr error
}

func (s *stubSkillInferenceService) InferSkillsFromEvents(_ context.Context, _ []*career.Event) (*skillinference.InferenceResult, error) {
	if s.inferErr != nil {
		return nil, s.inferErr
	}
	if s.result != nil {
		return s.result, nil
	}
	return &skillinference.InferenceResult{}, nil
}

func (s *stubSkillInferenceService) CreateSkillsFromSuggestions(_ context.Context, suggestions []skillinference.SkillSuggestion) ([]*career.Skill, error) {
	if s.createErr != nil {
		return nil, s.createErr
	}
	skills := make([]*career.Skill, len(suggestions))
	for i, sug := range suggestions {
		skills[i] = fixtures.SkillWith(fmt.Sprintf("skill-%d", i), sug.Name, "General", "mid")
	}
	return skills, nil
}

type stubSkillRepository struct {
	skillsByEvent map[string][]*career.Skill
	errorEvents   map[string]bool
}

func (r *stubSkillRepository) GetSkillsForEvent(_ context.Context, eventID string) ([]*career.Skill, error) {
	if r.errorEvents != nil && r.errorEvents[eventID] {
		return nil, errors.New("repo error")
	}
	if skills, ok := r.skillsByEvent[eventID]; ok {
		return skills, nil
	}
	return []*career.Skill{}, nil
}

func (r *stubSkillRepository) Create(_ context.Context, _ *career.Skill) error { return nil }
func (r *stubSkillRepository) Update(_ context.Context, _ *career.Skill) error { return nil }
func (r *stubSkillRepository) Delete(_ context.Context, _ string) error        { return nil }
func (r *stubSkillRepository) GetByID(_ context.Context, _ string) (*career.Skill, error) {
	return fixtures.Skill("stub"), nil
}
func (r *stubSkillRepository) GetByName(_ context.Context, _ string) (*career.Skill, error) {
	return fixtures.Skill("stub"), nil
}
func (r *stubSkillRepository) GetByCategory(_ context.Context, _ string) ([]*career.Skill, error) {
	return []*career.Skill{}, nil
}
func (r *stubSkillRepository) List(_ context.Context, _ *careerrepo.SkillListFilters) ([]*career.Skill, error) {
	return []*career.Skill{}, nil
}
func (r *stubSkillRepository) GetSkillsForEvents(_ context.Context, _ []string) ([]*career.Skill, error) {
	return []*career.Skill{}, nil
}
func (r *stubSkillRepository) GetEventsUsingSkill(_ context.Context, _ string) ([]*career.Event, error) {
	return []*career.Event{}, nil
}
func (r *stubSkillRepository) GetEventCountsForSkills(_ context.Context) (map[string]int, error) {
	return map[string]int{}, nil
}
func (r *stubSkillRepository) GetLastUsedForSkills(_ context.Context) (map[string]time.Time, error) {
	return map[string]time.Time{}, nil
}

var _ = Describe("HandleNavigate", func() {
	It("returns nil when burst not found", func() {
		ctx := &IntentValidator{Bursts: []*career.Burst{}}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		nav := burstviews.Nav{
			Action: burstviews.ActionView,
			Burst:  display.Burst{ID: "missing"},
		}
		result := &widgets.NavigateViewResult{ResultData: nav}
		cmd := intent.HandleNavigate(result)

		Expect(cmd).To(BeNil())
	})

	It("handles ActionEdit navigation", func() {
		burst := fixtures.Burst("nav-edit-1")
		ctx := &IntentValidator{Bursts: []*career.Burst{burst}}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		nav := burstviews.Nav{
			Action: burstviews.ActionEdit,
			Burst:  display.BurstFromDomain(burst),
		}
		result := &widgets.NavigateViewResult{ResultData: nav}
		cmd := intent.HandleNavigate(result)

		Expect(cmd).NotTo(BeNil())
	})

	It("handles ActionDelete navigation", func() {
		burst := fixtures.Burst("nav-delete-1")
		ctx := &IntentValidator{Bursts: []*career.Burst{burst}}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		nav := burstviews.Nav{
			Action: burstviews.ActionDelete,
			Burst:  display.BurstFromDomain(burst),
		}
		result := &widgets.NavigateViewResult{ResultData: nav}
		cmd := intent.HandleNavigate(result)

		Expect(cmd).To(BeNil())
	})

	It("handles ActionViewEvents navigation", func() {
		burst := fixtures.Burst("nav-events-1")
		ctx := &IntentValidator{Bursts: []*career.Burst{burst}}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		nav := burstviews.Nav{
			Action: burstviews.ActionViewEvents,
			Burst:  display.BurstFromDomain(burst),
		}
		result := &widgets.NavigateViewResult{ResultData: nav}
		cmd := intent.HandleNavigate(result)

		Expect(cmd).NotTo(BeNil())
	})

	It("handles ActionViewFacts navigation", func() {
		burst := fixtures.Burst("nav-facts-1")
		ctx := &IntentValidator{Bursts: []*career.Burst{burst}}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		nav := burstviews.Nav{
			Action: burstviews.ActionViewFacts,
			Burst:  display.BurstFromDomain(burst),
		}
		result := &widgets.NavigateViewResult{ResultData: nav}
		cmd := intent.HandleNavigate(result)

		Expect(cmd).NotTo(BeNil())
	})

	It("handles ActionViewSkills navigation", func() {
		burst := fixtures.Burst("nav-skills-1")
		ctx := &IntentValidator{Bursts: []*career.Burst{burst}}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		nav := burstviews.Nav{
			Action: burstviews.ActionViewSkills,
			Burst:  display.BurstFromDomain(burst),
		}
		result := &widgets.NavigateViewResult{ResultData: nav}
		cmd := intent.HandleNavigate(result)

		Expect(cmd).To(BeNil())
	})
})

var _ = Describe("showConfirmBurstModal", func() {
	It("returns nil when selectedBurst is nil", func() {
		ctx := &IntentValidator{Bursts: []*career.Burst{}}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		cmd := intent.showConfirmBurstModal()

		Expect(cmd).To(BeNil())
	})

	It("returns command that loads facts", func() {
		burst := fixtures.Burst("confirm-1")
		mockSvc := mocks.NewBurstServiceMock()
		ctx := &IntentValidator{
			Bursts:  []*career.Burst{burst},
			Service: mockSvc,
		}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.selectedBurst = burst

		cmd := intent.showConfirmBurstModal()

		Expect(cmd).NotTo(BeNil())
		msg := cmd()
		loadedMsg, ok := msg.(ConfirmBurstFactsLoadedMsg)
		Expect(ok).To(BeTrue())
		Expect(loadedMsg.Facts).NotTo(BeNil())
	})

	It("returns empty facts when service is nil", func() {
		burst := fixtures.Burst("confirm-2")
		ctx := &IntentValidator{
			Bursts:  []*career.Burst{burst},
			Service: nil,
		}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.selectedBurst = burst

		cmd := intent.showConfirmBurstModal()

		Expect(cmd).NotTo(BeNil())
		msg := cmd()
		loadedMsg, ok := msg.(ConfirmBurstFactsLoadedMsg)
		Expect(ok).To(BeTrue())
		Expect(loadedMsg.Facts).To(BeEmpty())
	})
})

var _ = Describe("HandleNavigate with display.Burst", func() {
	It("navigates to detail view when display.Burst is provided", func() {
		burst := fixtures.Burst("nav-display-burst-1")
		ctx := &IntentValidator{Bursts: []*career.Burst{burst}}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		displayBurst := display.BurstFromDomain(burst)
		result := &widgets.NavigateViewResult{ResultData: displayBurst}
		cmd := intent.HandleNavigate(result)

		Expect(cmd).To(BeNil())
		Expect(intent.selectedBurst).To(Equal(burst))
		Expect(intent.detailModal).NotTo(BeNil())
	})

	It("returns nil when display.Burst ID not found", func() {
		burst := fixtures.Burst("nav-display-burst-2")
		ctx := &IntentValidator{Bursts: []*career.Burst{burst}}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		displayBurst := display.Burst{ID: "nonexistent-id", Name: "Missing"}
		result := &widgets.NavigateViewResult{ResultData: displayBurst}
		cmd := intent.HandleNavigate(result)

		Expect(cmd).To(BeNil())
	})
})

var _ = Describe("handleSuggestionModalClosed with accepted suggestions", func() {
	It("returns SuggestionReviewCompleteMsg when suggestions are accepted", func() {
		setupIntent := func(ctx *IntentValidator) *Intent {
			ctx.Validate()
			var err error
			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			termInfo := terminal.NewInfo()
			termInfo.Width = 120
			termInfo.Height = 40
			termInfo.IsValid = true
			intent.UpdateTerminalInfo(termInfo)
			intent.Init()
			return intent
		}

		ctx := &IntentValidator{}
		intent := setupIntent(ctx)

		suggestion := burstfact.BurstSuggestion{
			Name:        "AcceptedSug",
			Description: "Test suggestion",
			EventIDs:    []string{"ev-1"},
		}
		intent.burstSuggestions = []burstfact.BurstSuggestion{suggestion}
		intent.state = StateSuggestionReview
		intent.suggestionModal = burstviews.NewSuggestionReview(
			display.BurstSuggestionsFromDomain([]burstfact.BurstSuggestion{suggestion}),
			themes.NewDefaultTheme(),
		)
		intent.suggestionModal.Show()

		intent.suggestionModal.Update(tea.KeyMsg{Type: tea.KeyEnter})

		cmd := intent.handleSuggestionModalClosed(noopCmd)

		Expect(cmd).NotTo(BeNil())
	})
})

var _ = Describe("findBurstSuggestionByName", func() {
	It("returns nil when name is empty", func() {
		ctx := &IntentValidator{}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		result := intent.findBurstSuggestionByName("")

		Expect(result).To(BeNil())
	})

	It("returns nil when suggestion not found", func() {
		ctx := &IntentValidator{}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		intent.burstSuggestions = []burstfact.BurstSuggestion{
			{Name: "ExistingSug", Description: "Test", EventIDs: []string{"ev-1"}},
		}

		result := intent.findBurstSuggestionByName("NonexistentSug")

		Expect(result).To(BeNil())
	})

	It("returns suggestion when found", func() {
		ctx := &IntentValidator{}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		suggestion := burstfact.BurstSuggestion{
			Name:        "FoundSug",
			Description: "Test",
			EventIDs:    []string{"ev-1"},
		}
		intent.burstSuggestions = []burstfact.BurstSuggestion{suggestion}

		result := intent.findBurstSuggestionByName("FoundSug")

		Expect(result).NotTo(BeNil())
		Expect(result.Name).To(Equal("FoundSug"))
	})
})

var _ = Describe("findSkillSuggestionByName", func() {
	It("returns nil when name is empty", func() {
		ctx := &IntentValidator{}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		result := intent.findSkillSuggestionByName("")

		Expect(result).To(BeNil())
	})

	It("returns nil when suggestion not found", func() {
		ctx := &IntentValidator{}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		intent.skillSuggestions = []skillinference.SkillSuggestion{
			{Name: "ExistingSkill", Category: "backend", EventIDs: []string{"ev-1"}},
		}

		result := intent.findSkillSuggestionByName("NonexistentSkill")

		Expect(result).To(BeNil())
	})

	It("returns suggestion when found", func() {
		ctx := &IntentValidator{}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		suggestion := skillinference.SkillSuggestion{
			Name:       "FoundSkill",
			Category:   "backend",
			EventIDs:   []string{"ev-1"},
			Confidence: 0.9,
		}
		intent.skillSuggestions = []skillinference.SkillSuggestion{suggestion}

		result := intent.findSkillSuggestionByName("FoundSkill")

		Expect(result).NotTo(BeNil())
		Expect(result.Name).To(Equal("FoundSkill"))
	})
})

var _ = Describe("showFacts", func() {
	It("returns nil when selectedBurst is nil", func() {
		ctx := &IntentValidator{}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		cmd := intent.showFacts()

		Expect(cmd).To(BeNil())
	})

	It("returns empty facts when service is nil", func() {
		burst := fixtures.Burst("facts-1")
		ctx := &IntentValidator{
			Bursts:  []*career.Burst{burst},
			Service: nil,
		}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.selectedBurst = burst

		cmd := intent.showFacts()

		Expect(cmd).NotTo(BeNil())
		msg := cmd()
		factsMsg, ok := msg.(BurstFactsLoadedMsg)
		Expect(ok).To(BeTrue())
		Expect(factsMsg.Facts).To(BeEmpty())
	})
})

var _ = Describe("findBurstByID", func() {
	It("returns nil when burst not found", func() {
		ctx := &IntentValidator{}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		result := intent.findBurstByID("nonexistent")

		Expect(result).To(BeNil())
	})

	It("returns burst when found", func() {
		burst := fixtures.Burst("find-1")
		ctx := &IntentValidator{
			Bursts: []*career.Burst{burst},
		}
		ctx.Validate()
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		result := intent.findBurstByID(burst.ID)

		Expect(result).NotTo(BeNil())
		Expect(result.ID).To(Equal(burst.ID))
	})
})
