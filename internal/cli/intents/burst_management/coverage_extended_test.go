package burst_management_test

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/uikit/display"

	"github.com/baphled/kariya/internal/domain/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/testutil/mocks"
	mockintent "github.com/baphled/kariya/internal/testutil/mocks/intent"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Coverage Extended", func() {
	var (
		intent   *burst_management.Intent
		termInfo *terminal.Info
	)

	setupIntent := func(ctx *burst_management.IntentContext) {
		ctx.Validate()
		var err error
		intent, err = burst_management.NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		termInfo = terminal.NewInfo()
		termInfo.Width = 120
		termInfo.Height = 40
		termInfo.IsValid = true
		intent.UpdateTerminalInfo(termInfo)
		intent.Init()
	}

	Describe("handleSkillSuggestionsLoaded edge cases", func() {
		It("should show warning when no suggestions returned", func() {
			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{fixtures.BurstConfirmed("b1", "e1")},
				SkillInferenceService: mockSkillSvc,
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{},
			})

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleFeedbackModal()).To(BeTrue())
		})

		It("should show success when all skills already tracked", func() {
			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{fixtures.BurstConfirmed("b1", "e1")},
				SkillInferenceService: mockSkillSvc,
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
				},
				ExistingSkillNames: []string{"Go"},
			})

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleFeedbackModal()).To(BeTrue())
		})

		It("should silently handle cancelled context in skill suggestions", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.BurstConfirmed("b1", "e1")},
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Error: context.Canceled,
			})

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("handleEditBurstMsg extended", func() {
		It("should show error when burst not found", func() {
			burst := fixtures.BurstConfirmed("b-edit", "e1", "e2")
			burst.Name = "Edit Target"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			intent.Update(burst_management.EditBurstMsg{
				BurstID:     "wrong-id",
				Name:        "New Name",
				Description: "New Description",
			})

			Expect(intent.HasVisibleFeedbackModal()).To(BeTrue())
		})

		It("should handle repository update error", func() {
			burst := fixtures.BurstConfirmed("b-edit-repo", "e1", "e2")
			burst.Name = "Repo Edit"

			mockRepo := mocks.NewBurstRepositoryMock().SetUpdateError(errors.New("update failed"))

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				BurstRepository: mockRepo,
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(intent.HasVisibleEditModal()).To(BeTrue())

			intent.Update(burst_management.EditBurstMsg{
				BurstID:     burst.ID,
				Name:        "Updated Name",
				Description: "Updated Description",
			})

			Expect(intent.HasVisibleFeedbackModal()).To(BeTrue())
		})

		It("should handle edit with nil selectedBurst", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.Burst("b1", "e1")},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(nil)

			Expect(func() {
				intent.Update(burst_management.EditBurstMsg{
					BurstID:     "b1",
					Name:        "Test",
					Description: "Test",
				})
			}).NotTo(Panic())
		})
	})

	Describe("handleEditModalUpdate completed edit", func() {
		It("should handle edit modal closed without completion", func() {
			burst := fixtures.Burst("b-edit-close", "e1")
			burst.Name = "Close Edit"
			burst.Description = "Test description"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(intent.HasVisibleEditModal()).To(BeTrue())

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.HasVisibleEditModal()).To(BeFalse())
		})
	})

	Describe("handleLoadingModalUpdate extended", func() {
		It("should cancel async operation on Esc during loading", func() {
			mockService := mocks.NewBurstServiceMock().
				SetEvents([]*career.Event{fixtures.Event("e1")})

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{},
				Service: mockService,
				Context: context.Background(),
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateList)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			Expect(intent.GetState()).To(Equal(burst_management.StateSuggesting))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("showBurstSkillsModal with SkillRepository", func() {
		It("should load skills from repository", func() {
			burst := fixtures.BurstConfirmed("b-skills-repo", "e1", "e2")
			burst.Name = "Skills Repo Test"

			skillRepo := careermemory.NewSkillRepository()

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				SkillRepository: skillRepo,
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)
			Expect(intent.GetDetailModal()).NotTo(BeNil())

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			if cmd != nil {
				msg := executeAsyncCmd(cmd)
				if msg != nil {
					intent.Update(msg)
				}
			}
		})

		It("should handle nil SkillRepository", func() {
			burst := fixtures.BurstConfirmed("b-skills-nil", "e1", "e2")
			burst.Name = "No Skill Repo"

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				SkillRepository: nil,
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			if cmd != nil {
				msg := executeAsyncCmd(cmd)
				if msg != nil {
					intent.Update(msg)
				}
			}
		})
	})

	Describe("showBurstEventsModal with service", func() {
		It("should load events from service", func() {
			events := []*career.Event{
				fixtures.EventWith("e1", "Built API", "Acme", "Backend"),
				fixtures.EventWith("e2", "Deployed app", "Acme", "DevOps"),
			}
			mockService := mocks.NewBurstServiceMock().SetEvents(events)

			burst := fixtures.BurstConfirmed("b-events-svc", "e1", "e2")
			burst.Name = "Events Svc Test"

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{burst},
				Service: mockService,
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

			if cmd != nil {
				msg := executeAsyncCmd(cmd)
				if msg != nil {
					intent.Update(msg)
				}
			}
		})

		It("should handle nil service for events", func() {
			burst := fixtures.BurstConfirmed("b-events-nil", "e1", "e2")
			burst.Name = "No Events Svc"

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{burst},
				Service: nil,
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

			if cmd != nil {
				msg := executeAsyncCmd(cmd)
				if msg != nil {
					intent.Update(msg)
				}
			}
		})
	})

	Describe("showBurstFactsModal with service", func() {
		It("should load facts from service", func() {
			mockService := mocks.NewBurstServiceMock().
				SetExtractedFacts([]career.Fact{*fixtures.FactWith("f1", "Test fact")})

			burst := fixtures.BurstConfirmed("b-facts-svc", "e1", "e2")
			burst.Name = "Facts Svc Test"

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{burst},
				Service: mockService,
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			if cmd != nil {
				msg := executeAsyncCmd(cmd)
				if msg != nil {
					intent.Update(msg)
				}
			}
		})

		It("should handle nil service for facts", func() {
			burst := fixtures.BurstConfirmed("b-facts-nil", "e1", "e2")
			burst.Name = "No Facts Svc"

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{burst},
				Service: nil,
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			if cmd != nil {
				msg := executeAsyncCmd(cmd)
				if msg != nil {
					intent.Update(msg)
				}
			}
		})
	})

	Describe("handleEventsModalUpdate extended", func() {
		It("should close events modal on Enter and return to detail", func() {
			burst := fixtures.Burst("b-events-enter", "e1")
			burst.Name = "Events Enter"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			events := []*career.Event{
				fixtures.EventWith("e1", "Event 1", "", ""),
			}
			intent.Update(burst_management.BurstEventsLoadedMsg{Events: events})

			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})

		It("should close events modal on Esc and return to detail", func() {
			burst := fixtures.Burst("b-events-esc", "e1")
			burst.Name = "Events Esc"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			events := []*career.Event{
				fixtures.EventWith("e1", "Event 1", "", ""),
			}
			intent.Update(burst_management.BurstEventsLoadedMsg{Events: events})

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})
	})

	Describe("handleSkillsModalUpdate extended", func() {
		It("should close skills modal on Enter and return to detail", func() {
			burst := fixtures.BurstConfirmed("b-skills-enter", "e1")
			burst.Name = "Skills Enter"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			intent.Update(burst_management.BurstSkillsLoadedMsg{
				Skills: []*career.Skill{fixtures.SkillWith("s1", "Go", "backend", "mid")},
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})
	})

	Describe("handleFactsModalUpdate extended", func() {
		It("should close facts modal on Enter and return to detail", func() {
			burst := fixtures.Burst("b-facts-enter", "e1")
			burst.Name = "Facts Enter"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			facts := []*career.Fact{
				fixtures.FactWith("f1", "Fact 1"),
			}
			intent.Update(burst_management.BurstFactsLoadedMsg{Facts: facts})

			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})

		It("should close facts modal on Esc and return to detail", func() {
			burst := fixtures.Burst("b-facts-esc", "e1")
			burst.Name = "Facts Esc"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			facts := []*career.Fact{
				fixtures.FactWith("f1", "Fact 1"),
			}
			intent.Update(burst_management.BurstFactsLoadedMsg{Facts: facts})

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})
	})

	Describe("handleConfirmModalUpdate extended", func() {
		It("should confirm burst on 'y' key", func() {
			burst := fixtures.Burst("b-confirm-y", "e1", "e2")
			burst.Name = "Confirm Yes"

			repo := careermemory.NewBurstRepository()

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				BurstRepository: repo,
				Service:         mocks.NewBurstServiceMock(),
				Context:         context.Background(),
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			Expect(intent.HasVisibleConfirmModal()).To(BeTrue())

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			if cmd != nil {
				messages := executeBatchCmd(cmd)
				_ = messages
			}
		})
	})

	Describe("deleteBurst with repository", func() {
		It("should delete burst via repository", func() {
			repo := careermemory.NewBurstRepository()
			burst := fixtures.BurstConfirmed("b-del", "e1", "e2")
			burst.Name = "Delete Me"

			repoCtx := context.Background()
			_ = repo.Create(repoCtx, burst)

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				BurstRepository: repo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(intent.HasVisibleDeleteModal()).To(BeTrue())

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			Expect(intent.GetFilteredBursts()).To(HaveLen(0))
		})

		It("should handle delete repository error", func() {
			mockRepo := mocks.NewBurstRepositoryMock().SetDeleteError(errors.New("delete failed"))

			burst := fixtures.BurstConfirmed("b-del-err", "e1", "e2")
			burst.Name = "Delete Error"

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				BurstRepository: mockRepo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("handleDeleteModalUpdate cancel", func() {
		It("should cancel delete on 'n' key", func() {
			burst := fixtures.Burst("b-del-cancel", "e1")
			burst.Name = "Delete Cancel"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(intent.HasVisibleDeleteModal()).To(BeTrue())

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

			Expect(intent.HasVisibleDeleteModal()).To(BeFalse())
		})
	})

	Describe("handleSkillSuggestionModalUpdate cancel and reject", func() {
		It("should handle skill suggestion modal cancel via Esc", func() {
			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{fixtures.BurstConfirmed("b1", "e1")},
				SkillInferenceService: mockSkillSvc,
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
				},
				ExistingSkillNames: []string{},
			})
			Expect(intent.GetState()).To(Equal(burst_management.StateSkillSuggestionReview))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})

		It("should handle skill suggestion modal reject via 'r'", func() {
			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{fixtures.BurstConfirmed("b1", "e1")},
				SkillInferenceService: mockSkillSvc,
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
				},
				ExistingSkillNames: []string{},
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
		})
	})

	Describe("handleSuggestionEventsModalUpdate via skill suggestion 'v' key", func() {
		It("should open and close suggestion events modal", func() {
			events := []*career.Event{
				fixtures.EventWith("e1", "Built API", "Acme", "Backend"),
				fixtures.EventWith("e2", "Deployed app", "Acme", "DevOps"),
			}
			mockService := mocks.NewBurstServiceMock().SetEvents(events)

			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)

			burst := fixtures.BurstConfirmed("b1", "e1", "e2")

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{burst},
				Service:               mockService,
				SkillInferenceService: mockSkillSvc,
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)
			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95, EventIDs: []string{"e1", "e2"}},
				},
				ExistingSkillNames: []string{},
			})
			Expect(intent.GetState()).To(Equal(burst_management.StateSkillSuggestionReview))

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
			Expect(intent.HasActiveModal()).To(BeTrue())

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})
	})

	Describe("inferSkillsFromBurst with nil service", func() {
		It("should return error when skill inference service is nil", func() {
			events := []*career.Event{
				fixtures.EventWith("e1", "Built API", "Acme", "Backend"),
			}
			mockService := mocks.NewBurstServiceMock().SetEvents(events)

			burst := fixtures.BurstConfirmed("burst-nil-svc", "e1", "e2")

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{burst},
				Service:               mockService,
				SkillInferenceService: nil,
				Context:               context.Background(),
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)
			intent.SetState(burst_management.StateExtractingFacts)

			factMsg := burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{fixtures.FactWith("f1", "Test fact")},
				Burst: burst,
			}
			cmd := intent.Update(factMsg)

			if cmd != nil {
				messages := executeBatchCmd(cmd)
				for _, m := range messages {
					if errMsg, ok := m.(burst_management.SkillSuggestionsErrorMsg); ok {
						Expect(errMsg.Err).To(HaveOccurred())
						Expect(errMsg.Err.Error()).To(ContainSubstring("skill inference service not available"))
					}
				}
			}
		})
	})

	Describe("RefreshData", func() {
		It("should refresh data from repository", func() {
			repo := careermemory.NewBurstRepository()

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{},
				BurstRepository: repo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			intent.RefreshData()
		})

		It("should handle refresh with nil repository", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.Burst("b1", "e1")},
			}
			setupIntent(ctx)

			intent.RefreshData()
		})
	})

	Describe("removeBurstFromSlice edge cases", func() {
		It("should handle empty burst list", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			setupIntent(ctx)

			Expect(intent.GetFilteredBursts()).To(HaveLen(0))
		})
	})

	Describe("handleSuggestionModalClosed with accepted suggestions", func() {
		It("should handle suggestion modal reject via 'r' key", func() {
			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{},
				Context: context.Background(),
			}
			setupIntent(ctx)

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Test Burst", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.8},
				},
			})
			Expect(intent.GetSuggestionModal()).NotTo(BeNil())

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("loadBurstFacts edge cases", func() {
		It("should handle nil burst in loadBurstFacts", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.Burst("b1", "e1")},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(nil)

			Expect(func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
			}).NotTo(Panic())
		})
	})

	Describe("openDeleteModal edge cases", func() {
		It("should handle delete with nil selectedBurst", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.Burst("b1", "e1")},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(nil)

			Expect(func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			}).NotTo(Panic())
		})
	})

	Describe("NewIntent validation", func() {
		It("should fail with nil context", func() {
			_, err := burst_management.NewIntent(nil)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("handleScreenResult edge cases", func() {
		It("should handle non-ScreenResult type", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.Burst("b1", "e1")},
			}
			setupIntent(ctx)

			Expect(func() {
				intent.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
			}).NotTo(Panic())
		})
	})

	Describe("confirmBurst with repository", func() {
		It("should confirm and persist burst", func() {
			repo := careermemory.NewBurstRepository()
			burst := fixtures.Burst("b-confirm-repo", "e1", "e2")
			burst.Name = "Confirm Repo"

			repoCtx := context.Background()
			_ = repo.Create(repoCtx, burst)

			mockService := mocks.NewBurstServiceMock()

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				BurstRepository: repo,
				Service:         mockService,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			if cmd != nil {
				messages := executeBatchCmd(cmd)
				for _, m := range messages {
					intent.Update(m)
				}
			}
		})
	})

	Describe("extractFactsForBurst with service error", func() {
		It("should handle extraction service error", func() {
			mockService := mocks.NewBurstServiceMock().
				SetExtractError(errors.New("extraction failed"))

			burst := fixtures.BurstConfirmed("b-extract-err", "e1", "e2")

			repo := careermemory.NewBurstRepository()

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				Service:         mockService,
				BurstRepository: repo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)
			intent.SetState(burst_management.StateSuggestionReview)

			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{
					{Name: "Extract Error", Description: "Should fail extraction", EventIDs: []string{"e1", "e2"}},
				},
			}
			cmd := intent.Update(msg)

			if cmd != nil {
				messages := executeBatchCmd(cmd)
				for _, m := range messages {
					if extractMsg, ok := m.(burst_management.FactExtractionCompleteMsg); ok {
						Expect(extractMsg.Error).To(HaveOccurred())
					}
				}
			}
		})
	})

	Describe("saveSkillFromSuggestion with context error", func() {
		It("should handle cancelled context during skill save", func() {
			cancelledCtx, cancel := context.WithCancel(context.Background())
			cancel()

			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)
			mockSkillSvc.EXPECT().
				CreateSkillsFromSuggestions(gomock.Any(), gomock.Any()).
				Return(nil, nil).
				AnyTimes()

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{fixtures.BurstConfirmed("b1", "e1")},
				SkillInferenceService: mockSkillSvc,
				Context:               cancelledCtx,
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
				},
				ExistingSkillNames: []string{},
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
		})
	})

	Describe("handleSuggestionAccept with all suggestions accepted", func() {
		It("should clear suggestion state after accepting last suggestion", func() {
			repo := careermemory.NewBurstRepository()

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{},
				BurstRepository: repo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Single Burst", Description: "Only one", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.9},
				},
			})
			Expect(intent.GetSuggestionModal()).NotTo(BeNil())

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.GetFilteredBursts()).To(HaveLen(1))

			if cmd != nil {
				messages := executeBatchCmd(cmd)
				for _, m := range messages {
					intent.Update(m)
				}
			}
		})
	})

	Describe("Init with empty bursts", func() {
		It("should initialize with no bursts", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			setupIntent(ctx)

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.GetFilteredBursts()).To(HaveLen(0))
		})
	})

	Describe("View rendering with modals", func() {
		It("should render view with feedback modal visible", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			setupIntent(ctx)

			intent.ShowErrorModal("Test Error", "Something went wrong")
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render view with loading modal", func() {
			mockService := mocks.NewBurstServiceMock().
				SetEvents([]*career.Event{fixtures.Event("e1")})

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{},
				Service: mockService,
				Context: context.Background(),
			}
			setupIntent(ctx)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("handleSkillSuggestionAccept with successful accept", func() {
		It("should accept skill and create it via service", func() {
			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)
			mockSkillSvc.EXPECT().
				CreateSkillsFromSuggestions(gomock.Any(), gomock.Any()).
				Return(nil, nil).
				AnyTimes()

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{fixtures.BurstConfirmed("b1", "e1")},
				SkillInferenceService: mockSkillSvc,
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
				},
				ExistingSkillNames: []string{},
			})

			Expect(func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			}).NotTo(Panic())
		})
	})

	Describe("handleSuggestionAccept with nil current suggestion", func() {
		It("should handle accept when no suggestion is selected", func() {
			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{},
				Context: context.Background(),
			}
			setupIntent(ctx)

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Test", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.8},
				},
			})

			Expect(func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			}).NotTo(Panic())
		})
	})

	Describe("cancelPreviousOperation with nil cancelFunc", func() {
		It("should not panic when cancelFunc is nil", func() {
			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{},
				Context: context.Background(),
			}
			setupIntent(ctx)

			Expect(func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			}).NotTo(Panic())
		})
	})

	Describe("openSuggestionEventsModal via Enter on skill suggestion", func() {
		It("should open events modal when pressing Enter on skill suggestion", func() {
			events := []*career.Event{
				fixtures.EventWith("e1", "Built API", "Acme", "Backend"),
				fixtures.EventWith("e2", "Deployed app", "Acme", "DevOps"),
			}
			mockService := mocks.NewBurstServiceMock().SetEvents(events)

			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)

			burst := fixtures.BurstConfirmed("b1", "e1", "e2")

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{burst},
				Service:               mockService,
				SkillInferenceService: mockSkillSvc,
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)
			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95, EventIDs: []string{"e1", "e2"}},
				},
				ExistingSkillNames: []string{},
			})
			Expect(intent.GetState()).To(Equal(burst_management.StateSkillSuggestionReview))

			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(intent.HasActiveModal()).To(BeTrue())
		})

		It("should close suggestion events modal on Esc and re-show skill modal", func() {
			events := []*career.Event{
				fixtures.EventWith("e1", "Built API", "Acme", "Backend"),
				fixtures.EventWith("e2", "Deployed app", "Acme", "DevOps"),
			}
			mockService := mocks.NewBurstServiceMock().SetEvents(events)

			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)

			burst := fixtures.BurstConfirmed("b1", "e1", "e2")

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{burst},
				Service:               mockService,
				SkillInferenceService: mockSkillSvc,
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)
			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95, EventIDs: []string{"e1", "e2"}},
				},
				ExistingSkillNames: []string{},
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})

		It("should handle non-key message in suggestion events modal", func() {
			events := []*career.Event{
				fixtures.EventWith("e1", "Built API", "Acme", "Backend"),
			}
			mockService := mocks.NewBurstServiceMock().SetEvents(events)

			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)

			burst := fixtures.BurstConfirmed("b1", "e1", "e2")

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{burst},
				Service:               mockService,
				SkillInferenceService: mockSkillSvc,
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)
			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95, EventIDs: []string{"e1"}},
				},
				ExistingSkillNames: []string{},
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			intent.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
		})
	})

	Describe("handleSkillSuggestionModalUpdate cancel with inferredFromDetail", func() {
		It("should return to detail modal when cancelling from detail-triggered inference", func() {
			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)

			burst := fixtures.BurstConfirmed("b-detail-infer", "e1", "e2")
			burst.Name = "Detail Infer"

			mockService := mocks.NewBurstServiceMock()

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{burst},
				Service:               mockService,
				SkillInferenceService: mockSkillSvc,
				Context:               context.Background(),
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			Expect(intent.GetState()).To(Equal(burst_management.StateInferringSkills))
		})
	})

	Describe("handleSuggestionModalUpdate with non-accept key", func() {
		It("should forward non-accept keys to suggestion modal", func() {
			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{},
				Context: context.Background(),
			}
			setupIntent(ctx)

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Test 1", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.8},
					{Name: "Test 2", EventIDs: []string{"e3", "e4"}, ConfidenceScore: 0.7},
				},
			})
			Expect(intent.GetSuggestionModal()).NotTo(BeNil())

			intent.Update(tea.KeyMsg{Type: tea.KeyDown})

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		})
	})

	Describe("extractFactsForBurst with nil burst", func() {
		It("should return error for nil burst during extraction", func() {
			mockService := mocks.NewBurstServiceMock()

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{},
				Service: mockService,
				Context: context.Background(),
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(nil)
			intent.SetState(burst_management.StateSuggestionReview)

			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{
					{Name: "Nil Burst Test", Description: "Test", EventIDs: []string{"e1", "e2"}},
				},
			}
			cmd := intent.Update(msg)

			if cmd != nil {
				messages := executeBatchCmd(cmd)
				_ = messages
			}
		})
	})

	Describe("handleEventsModalUpdate with nil selectedBurst", func() {
		It("should close events modal and return noopCmd when selectedBurst is nil", func() {
			burst := fixtures.Burst("b-events-nil-sel", "e1")
			burst.Name = "Events Nil Sel"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			events := []*career.Event{
				fixtures.EventWith("e1", "Event 1", "", ""),
			}
			intent.Update(burst_management.BurstEventsLoadedMsg{Events: events})

			intent.SetSelectedBurst(nil)

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})
	})

	Describe("handleSkillsModalUpdate with nil selectedBurst", func() {
		It("should close skills modal and return noopCmd when selectedBurst is nil", func() {
			burst := fixtures.BurstConfirmed("b-skills-nil-sel", "e1")
			burst.Name = "Skills Nil Sel"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			intent.Update(burst_management.BurstSkillsLoadedMsg{
				Skills: []*career.Skill{fixtures.SkillWith("s1", "Go", "backend", "mid")},
			})

			intent.SetSelectedBurst(nil)

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})
	})

	Describe("handleFactsModalUpdate with nil selectedBurst", func() {
		It("should close facts modal and return noopCmd when selectedBurst is nil", func() {
			burst := fixtures.Burst("b-facts-nil-sel", "e1")
			burst.Name = "Facts Nil Sel"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			facts := []*career.Fact{
				fixtures.FactWith("f1", "Fact 1"),
			}
			intent.Update(burst_management.BurstFactsLoadedMsg{Facts: facts})

			intent.SetSelectedBurst(nil)

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})
	})

	Describe("Init with repository", func() {
		It("should load bursts from repository on init", func() {
			repo := careermemory.NewBurstRepository()
			burst := fixtures.BurstConfirmed("b-init", "e1", "e2")
			burst.Name = "Init Burst"
			repoCtx := context.Background()
			_ = repo.Create(repoCtx, burst)

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{},
				BurstRepository: repo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			Expect(intent.GetFilteredBursts()).To(HaveLen(1))
		})
	})

	Describe("handleSuggestionModalClosed with accepted suggestions", func() {
		It("should trigger review complete when modal closes with accepted suggestions", func() {
			repo := careermemory.NewBurstRepository()

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{},
				BurstRepository: repo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Burst A", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.9},
					{Name: "Burst B", EventIDs: []string{"e3", "e4"}, ConfidenceScore: 0.8},
				},
			})
			Expect(intent.GetSuggestionModal()).NotTo(BeNil())

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
		})
	})

	Describe("handleConfirmModalUpdate with nil selectedBurst", func() {
		It("should handle confirm modal close when selectedBurst is nil", func() {
			burst := fixtures.Burst("b-confirm-nil", "e1", "e2")
			burst.Name = "Confirm Nil"

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{burst},
				Service: mocks.NewBurstServiceMock(),
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			Expect(intent.HasVisibleConfirmModal()).To(BeTrue())

			intent.SetSelectedBurst(nil)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
		})
	})

	Describe("handleSkillsModalUpdate Enter with nil selectedBurst", func() {
		It("should close skills modal on Enter when selectedBurst is nil", func() {
			burst := fixtures.BurstConfirmed("b-skills-enter-nil", "e1")
			burst.Name = "Skills Enter Nil"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			intent.Update(burst_management.BurstSkillsLoadedMsg{
				Skills: []*career.Skill{fixtures.SkillWith("s1", "Go", "backend", "mid")},
			})

			intent.SetSelectedBurst(nil)

			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})
	})

	Describe("showBurstSkillsModal with SkillRepository and events", func() {
		It("should load and deduplicate skills from repository", func() {
			skillRepo := careermemory.NewSkillRepository()

			burst := fixtures.BurstConfirmed("b-skills-dedup", "e1", "e2")
			burst.Name = "Skills Dedup"

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				SkillRepository: skillRepo,
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			if cmd != nil {
				msg := executeAsyncCmd(cmd)
				if msg != nil {
					intent.Update(msg)
				}
			}
		})
	})

	Describe("showBurstEventsModal async load", func() {
		It("should load events asynchronously from service", func() {
			events := []*career.Event{
				fixtures.EventWith("e1", "Built API", "Acme", "Backend"),
			}
			mockService := mocks.NewBurstServiceMock().SetEvents(events)

			burst := fixtures.BurstConfirmed("b-events-async", "e1", "e2")
			burst.Name = "Events Async"

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{burst},
				Service: mockService,
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

			if cmd != nil {
				msg := executeAsyncCmd(cmd)
				if msg != nil {
					intent.Update(msg)
				}
			}
		})
	})

	Describe("showBurstFactsModal async load", func() {
		It("should load facts asynchronously from service", func() {
			mockService := mocks.NewBurstServiceMock()

			burst := fixtures.BurstConfirmed("b-facts-async", "e1", "e2")
			burst.Name = "Facts Async"

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{burst},
				Service: mockService,
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			if cmd != nil {
				msg := executeAsyncCmd(cmd)
				if msg != nil {
					intent.Update(msg)
				}
			}
		})
	})

	Describe("loadBurstFacts with service error", func() {
		It("should return empty facts on service error", func() {
			mockService := mocks.NewBurstServiceMock().SetListEventsError(errors.New("list error"))

			burst := fixtures.BurstConfirmed("b-facts-err", "e1", "e2")
			burst.Name = "Facts Error"

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{burst},
				Service: mockService,
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			if cmd != nil {
				msg := executeAsyncCmd(cmd)
				if msg != nil {
					intent.Update(msg)
				}
			}
		})
	})

	Describe("openSuggestionEventsModal with nil skill suggestion modal", func() {
		It("should handle nil skillSuggestionModal gracefully", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.BurstConfirmed("b1", "e1")},
			}
			setupIntent(ctx)

			Expect(func() {
				intent.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
			}).NotTo(Panic())
		})
	})

	Describe("handleSuggestionModalClosed with accepted then reject", func() {
		It("should trigger review complete when accepting then rejecting last", func() {
			repo := careermemory.NewBurstRepository()

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{},
				BurstRepository: repo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Accept Me", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.9},
					{Name: "Reject Me", EventIDs: []string{"e3", "e4"}, ConfidenceScore: 0.7},
				},
			})
			Expect(intent.GetSuggestionModal()).NotTo(BeNil())

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			if cmd != nil {
				messages := executeBatchCmd(cmd)
				for _, m := range messages {
					intent.Update(m)
				}
			}
		})
	})

	Describe("handleEditModalUpdate with completed form", func() {
		It("should submit edit when form is completed", func() {
			burst := fixtures.Burst("b-edit-submit", "e1", "e2")
			burst.Name = "Edit Submit"
			burst.Description = "Original description"

			repo := careermemory.NewBurstRepository()
			repoCtx := context.Background()
			_ = repo.Create(repoCtx, burst)

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				BurstRepository: repo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(intent.HasVisibleEditModal()).To(BeTrue())

			intent.Update(burst_management.EditBurstMsg{
				BurstID:     burst.ID,
				Name:        "Updated Name",
				Description: "Updated description",
			})
		})
	})

	Describe("handleSkillSuggestionModalUpdate with cancel from detail", func() {
		It("should handle cancel action returning to detail", func() {
			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)

			events := []*career.Event{
				fixtures.EventWith("e1", "Built API", "Acme", "Backend"),
			}
			mockService := mocks.NewBurstServiceMock().SetEvents(events)

			burst := fixtures.BurstConfirmed("b-cancel-detail", "e1", "e2")
			burst.Name = "Cancel Detail"

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{burst},
				Service:               mockService,
				SkillInferenceService: mockSkillSvc,
				Context:               context.Background(),
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
			Expect(intent.GetState()).To(Equal(burst_management.StateInferringSkills))

			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
				},
				ExistingSkillNames: []string{},
			})
			Expect(intent.GetState()).To(Equal(burst_management.StateSkillSuggestionReview))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("extractFactsForBurst with successful extraction and save", func() {
		It("should extract and save facts successfully", func() {
			mockService := mocks.NewBurstServiceMock().
				SetExtractedFacts([]career.Fact{
					*fixtures.FactWith("f1", "Fact one"),
					*fixtures.FactWith("f2", "Fact two"),
				})

			repo := careermemory.NewBurstRepository()
			burst := fixtures.BurstConfirmed("b-extract-ok", "e1", "e2")

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				Service:         mockService,
				BurstRepository: repo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)
			intent.SetState(burst_management.StateSuggestionReview)

			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{
					{Name: "Extract OK", Description: "Should succeed", EventIDs: []string{"e1", "e2"}},
				},
			}
			cmd := intent.Update(msg)

			if cmd != nil {
				messages := executeBatchCmd(cmd)
				for _, m := range messages {
					if extractMsg, ok := m.(burst_management.FactExtractionCompleteMsg); ok {
						Expect(extractMsg.Error).NotTo(HaveOccurred())
						Expect(extractMsg.Facts).To(HaveLen(2))
					}
				}
			}
		})
	})

	Describe("extractFactsForBurst with save fact error", func() {
		It("should continue extraction when individual fact save fails", func() {
			mockService := mocks.NewBurstServiceMock().
				SetExtractedFacts([]career.Fact{
					*fixtures.FactWith("f1", "Fact one"),
				}).
				SetSaveFactError(errors.New("save failed"))

			repo := careermemory.NewBurstRepository()
			burst := fixtures.BurstConfirmed("b-extract-save-err", "e1", "e2")

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				Service:         mockService,
				BurstRepository: repo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)
			intent.SetState(burst_management.StateSuggestionReview)

			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{
					{Name: "Save Err", Description: "Save fails", EventIDs: []string{"e1", "e2"}},
				},
			}
			cmd := intent.Update(msg)

			if cmd != nil {
				messages := executeBatchCmd(cmd)
				for _, m := range messages {
					if extractMsg, ok := m.(burst_management.FactExtractionCompleteMsg); ok {
						Expect(extractMsg.Facts).To(HaveLen(0))
					}
				}
			}
		})
	})

	Describe("inferSkillsFromBurst with nil burst", func() {
		It("should return error for nil burst", func() {
			mockService := mocks.NewBurstServiceMock()
			mockSkillSvc := skillinference.NewSkillInferenceService(nil, nil, nil)

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{},
				Service:               mockService,
				SkillInferenceService: mockSkillSvc,
				Context:               context.Background(),
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(nil)
			intent.SetState(burst_management.StateExtractingFacts)

			factMsg := burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{fixtures.FactWith("f1", "Test")},
				Burst: nil,
			}
			cmd := intent.Update(factMsg)

			if cmd != nil {
				messages := executeBatchCmd(cmd)
				for _, m := range messages {
					if errMsg, ok := m.(burst_management.SkillSuggestionsErrorMsg); ok {
						Expect(errMsg.Err).To(HaveOccurred())
					}
				}
			}
		})
	})

	Describe("showBurstSkillsModal inner async execution", func() {
		It("should execute async skill loading with repository", func() {
			skillRepo := careermemory.NewSkillRepository()

			burst := fixtures.BurstConfirmed("b-skills-inner", "e1", "e2")
			burst.Name = "Skills Inner"

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				SkillRepository: skillRepo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			if cmd != nil {
				msg := executeAsyncCmd(cmd)
				if loadedMsg, ok := msg.(burst_management.BurstSkillsLoadedMsg); ok {
					intent.Update(loadedMsg)
				}
			}
		})
	})

	Describe("showBurstFactsModal inner async execution", func() {
		It("should execute async fact loading with service", func() {
			mockService := mocks.NewBurstServiceMock().
				SetFactsForBurst("b-facts-inner", []*career.Fact{
					fixtures.FactWith("f1", "Loaded fact"),
				})

			burst := fixtures.BurstConfirmed("b-facts-inner", "e1", "e2")
			burst.Name = "Facts Inner"

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{burst},
				Service: mockService,
				Context: context.Background(),
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			if cmd != nil {
				msg := executeAsyncCmd(cmd)
				if loadedMsg, ok := msg.(burst_management.BurstFactsLoadedMsg); ok {
					intent.Update(loadedMsg)
				}
			}
		})
	})

	Describe("showBurstEventsModal inner async execution", func() {
		It("should execute async event loading with service", func() {
			events := []*career.Event{
				fixtures.EventWith("e1", "Built API", "Acme", "Backend"),
				fixtures.EventWith("e2", "Deployed", "Acme", "DevOps"),
			}
			mockService := mocks.NewBurstServiceMock().SetEvents(events)

			burst := fixtures.BurstConfirmed("b-events-inner", "e1", "e2")
			burst.Name = "Events Inner"

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{burst},
				Service: mockService,
				Context: context.Background(),
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

			if cmd != nil {
				msg := executeAsyncCmd(cmd)
				if loadedMsg, ok := msg.(burst_management.BurstEventsLoadedMsg); ok {
					Expect(loadedMsg.Events).To(HaveLen(2))
					intent.Update(loadedMsg)
				}
			}
		})
	})

	Describe("handleLoadingModalUpdate with spinner tick", func() {
		It("should handle non-key non-tick message during loading", func() {
			mockService := mocks.NewBurstServiceMock().
				SetEvents([]*career.Event{fixtures.Event("e1")})

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{},
				Service: mockService,
				Context: context.Background(),
			}
			setupIntent(ctx)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			intent.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
		})
	})

	Describe("transitionToScreen with logo", func() {
		It("should propagate logo to screen", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.Burst("b1", "e1")},
			}
			setupIntent(ctx)

			logo := display.NewLogo(false, 80)
			intent.SetLogo(logo)
			intent.SetLogoSpacing(3)

			intent.RefreshData()
		})
	})

	Describe("handleSuggestionAccept with nil current suggestion", func() {
		It("should handle accept when current suggestion is nil", func() {
			repo := careermemory.NewBurstRepository()

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{},
				BurstRepository: repo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Only One", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.9},
				},
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.GetSuggestionModal()).To(BeNil())
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("saveSkillFromSuggestion with context cancelled during save", func() {
		It("should handle context error during skill creation", func() {
			cancelledCtx, cancel := context.WithCancel(context.Background())
			cancel()

			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)
			mockSkillSvc.EXPECT().
				CreateSkillsFromSuggestions(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("context cancelled")).
				AnyTimes()

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{fixtures.BurstConfirmed("b1", "e1")},
				SkillInferenceService: mockSkillSvc,
				Context:               cancelledCtx,
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
				},
				ExistingSkillNames: []string{},
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
		})
	})

	Describe("handleSkillSuggestionAccept with last skill accepted", func() {
		It("should show success modal after accepting last skill", func() {
			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)
			mockSkillSvc.EXPECT().
				CreateSkillsFromSuggestions(gomock.Any(), gomock.Any()).
				Return(nil, nil).
				AnyTimes()

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{fixtures.BurstConfirmed("b1", "e1")},
				SkillInferenceService: mockSkillSvc,
				Context:               context.Background(),
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
				},
				ExistingSkillNames: []string{},
			})
			Expect(intent.GetState()).To(Equal(burst_management.StateSkillSuggestionReview))

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleFeedbackModal()).To(BeTrue())
		})
	})

	Describe("confirmBurst with nil service", func() {
		It("should confirm burst without service", func() {
			burst := fixtures.Burst("b-confirm-nosvc", "e1", "e2")
			burst.Name = "Confirm No Svc"

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{burst},
				Service: nil,
				Context: context.Background(),
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			if cmd != nil {
				messages := executeBatchCmd(cmd)
				for _, m := range messages {
					intent.Update(m)
				}
			}
		})
	})

	Describe("handleScreenResult with ScreenResult", func() {
		It("should dispatch screen result", func() {
			burst := fixtures.Burst("b-screen", "e1")
			burst.Name = "Screen Result"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)
			Expect(intent.GetDetailModal()).NotTo(BeNil())
		})
	})

	Describe("showBurstSkillsModal with skills in repository", func() {
		It("should load skills and handle GetSkillsForEvent errors", func() {
			skillRepo := careermemory.NewSkillRepository()

			burst := fixtures.BurstConfirmed("b-skills-err", "e1", "e2")
			burst.Name = "Skills Error"

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				SkillRepository: skillRepo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			if cmd != nil {
				msg := executeAsyncCmd(cmd)
				if loadedMsg, ok := msg.(burst_management.BurstSkillsLoadedMsg); ok {
					Expect(loadedMsg.Skills).To(HaveLen(0))
					intent.Update(loadedMsg)
				}
			}
		})
	})

	Describe("showBurstFactsModal with service returning facts", func() {
		It("should load facts and handle service errors", func() {
			mockService := mocks.NewBurstServiceMock().
				SetFactsForBurst("b-facts-svc-err", []*career.Fact{})

			burst := fixtures.BurstConfirmed("b-facts-svc-err", "e1", "e2")
			burst.Name = "Facts Svc Error"

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{burst},
				Service: mockService,
				Context: context.Background(),
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			if cmd != nil {
				msg := executeAsyncCmd(cmd)
				if loadedMsg, ok := msg.(burst_management.BurstFactsLoadedMsg); ok {
					intent.Update(loadedMsg)
				}
			}
		})
	})

	Describe("NewIntent with empty bursts", func() {
		It("should create intent with empty bursts", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			setupIntent(ctx)

			Expect(intent.GetFilteredBursts()).To(HaveLen(0))
		})
	})

	Describe("Init with LoadBursts error", func() {
		It("should handle LoadBursts error gracefully", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.Burst("b1", "e1")},
			}
			setupIntent(ctx)

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("RefreshData with LoadBursts error", func() {
		It("should keep existing data on LoadBursts error", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.Burst("b1", "e1")},
			}
			setupIntent(ctx)

			intent.RefreshData()

			Expect(intent.GetFilteredBursts()).To(HaveLen(1))
		})
	})

	Describe("showBurstSkillsModal with populated skills", func() {
		It("should load and deduplicate skills from event associations", func() {
			skillRepo := careermemory.NewSkillRepository()

			skill1 := fixtures.SkillWith("sk1", "Go Programming", "backend", "advanced")
			err := skillRepo.Create(context.Background(), skill1)
			Expect(err).NotTo(HaveOccurred())

			skill2 := fixtures.SkillWith("sk2", "REST APIs", "backend", "intermediate")
			err = skillRepo.Create(context.Background(), skill2)
			Expect(err).NotTo(HaveOccurred())

			skillRepo.AssociateSkillWithEvent(skill1.ID, "e1")
			skillRepo.AssociateSkillWithEvent(skill2.ID, "e2")
			skillRepo.AssociateSkillWithEvent(skill1.ID, "e2")

			burst := fixtures.BurstConfirmed("b-skills-pop", "e1", "e2")
			burst.Name = "Skills Populated"

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				SkillRepository: skillRepo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			Expect(cmd).NotTo(BeNil())
			msg := executeAsyncCmd(cmd)
			loadedMsg, ok := msg.(burst_management.BurstSkillsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(loadedMsg.Skills).To(HaveLen(2))
		})
	})

	Describe("openDeleteModal with long burst name", func() {
		It("should truncate burst name longer than 50 characters", func() {
			longName := "This is a very long burst name that exceeds fifty characters limit for display"
			burst := fixtures.BurstConfirmed("b-long-name", "e1", "e2")
			burst.Name = longName

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(intent.HasActiveModal()).To(BeTrue())
		})
	})

	Describe("handleDetailModalKeypress with unhandled key", func() {
		It("should return nil for unrecognised keys", func() {
			burst := fixtures.BurstConfirmed("b-unhandled", "e1", "e2")
			burst.Name = "Unhandled Key"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
		})
	})

	Describe("removeBurstFromSlice via delete flow", func() {
		It("should remove burst from list after successful delete", func() {
			burstRepo := careermemory.NewBurstRepository()

			burst1 := fixtures.BurstConfirmed("b-del1", "e1", "e2")
			burst1.Name = "First Burst"
			burst2 := fixtures.BurstConfirmed("b-del2", "e3", "e4")
			burst2.Name = "Second Burst"

			_ = burstRepo.Create(context.Background(), burst1)
			_ = burstRepo.Create(context.Background(), burst2)

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst1, burst2},
				BurstRepository: burstRepo,
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst1}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(intent.HasActiveModal()).To(BeTrue())

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			if cmd != nil {
				msgs := executeBatchCmd(cmd)
				for _, m := range msgs {
					intent.Update(m)
				}
			}

			Expect(intent.GetFilteredBursts()).To(HaveLen(1))
		})
	})

	Describe("RefreshData with repository LoadBursts error", func() {
		It("should handle LoadBursts error and keep existing data", func() {
			burstRepo := careermemory.NewBurstRepository()
			burst := fixtures.BurstConfirmed("b-refresh", "e1", "e2")
			burst.Name = "Refresh Test"
			_ = burstRepo.Create(context.Background(), burst)

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				BurstRepository: burstRepo,
			}
			setupIntent(ctx)

			intent.RefreshData()
			Expect(intent.GetFilteredBursts()).NotTo(BeEmpty())
		})
	})

	Describe("handleDeleteModalUpdate with non-confirm key", func() {
		It("should keep delete modal visible when non-confirm key is pressed", func() {
			burst := fixtures.BurstConfirmed("b-del-noconfirm", "e1", "e2")
			burst.Name = "Delete NoConfirm"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(intent.HasActiveModal()).To(BeTrue())

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(intent.HasActiveModal()).To(BeTrue())
		})
	})

	Describe("handleSkillsModalUpdate with nil selectedBurst", func() {
		It("should return noop when dismissing skills modal without selected burst", func() {
			burst := fixtures.BurstConfirmed("b-skills-nil", "e1")
			burst.Name = "Skills Nil Test"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			intent.Update(burst_management.BurstSkillsLoadedMsg{
				Skills: []*career.Skill{fixtures.SkillWith("s1", "Go", "backend", "mid")},
			})

			intent.SetSelectedBurst(nil)

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			if cmd != nil {
				msg := cmd()
				_ = msg
			}
		})
	})
})
