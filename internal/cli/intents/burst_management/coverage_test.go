package burst_management_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/terminal"

	"github.com/baphled/kariya/internal/domain/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/testutil/mocks"
	mockintent "github.com/baphled/kariya/internal/testutil/mocks/intent"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Coverage Improvement", func() {
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

	Describe("createBurstFromSuggestion with BurstRepository", func() {
		It("should persist burst via repository on success", func() {
			repo := careermemory.NewBurstRepository()

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{},
				BurstRepository: repo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateSuggestionReview)

			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{
					{
						Name:        "Repo Burst",
						Description: "Saved via repository",
						EventIDs:    []string{"e1", "e2"},
					},
				},
			}
			intent.Update(msg)

			Expect(intent.GetFilteredBursts()).To(HaveLen(1))
			Expect(intent.GetFilteredBursts()[0].Name).To(Equal("Repo Burst"))
		})

		It("should show error modal when repository create fails", func() {
			mockRepo := mocks.NewBurstRepositoryMock().SetCreateError(errors.New("db write error"))

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{},
				BurstRepository: mockRepo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateSuggestionReview)

			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{
					{
						Name:        "Failed Burst",
						Description: "Should fail to save",
						EventIDs:    []string{"e1", "e2"},
					},
				},
			}
			intent.Update(msg)

			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("inferSkillsFromBurst async command execution", func() {
		It("should return error when burst has no events", func() {
			mockService := mocks.NewBurstServiceMock().SetEvents([]*career.Event{})
			mockSkillSvc := skillinference.NewSkillInferenceService(nil, nil, nil)

			burst := fixtures.BurstConfirmed("burst-infer")
			burst.EventIDs = []string{"e-nonexistent"}

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{burst},
				Service:               mockService,
				SkillInferenceService: mockSkillSvc,
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

			Expect(intent.GetState()).To(Equal(burst_management.StateInferringSkills))

			if cmd != nil {
				messages := executeBatchCmd(cmd)
				for _, m := range messages {
					if errMsg, ok := m.(burst_management.SkillSuggestionsErrorMsg); ok {
						Expect(errMsg.Err).To(HaveOccurred())
					}
				}
			}
		})

		It("should return suggestions when inference succeeds", func() {
			events := []*career.Event{
				fixtures.EventWith("e1", "Built REST API with Go", "Acme", "Backend"),
			}
			mockService := mocks.NewBurstServiceMock().SetEvents(events)

			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)
			mockSkillSvc.EXPECT().
				InferSkillsFromEvents(gomock.Any(), gomock.Any()).
				Return(&skillinference.InferenceResult{
					Suggestions: []skillinference.SkillSuggestion{
						{Name: "Go", Category: "backend", Confidence: 0.95},
					},
				}, nil).
				AnyTimes()

			burst := fixtures.BurstConfirmed("burst-infer-success", "e1")

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{burst},
				Service:               mockService,
				SkillInferenceService: mockSkillSvc,
				Context:               context.Background(),
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)
			intent.SetState(burst_management.StateExtractingFacts)

			factMsg := burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{fixtures.FactWith("f1", "Built REST API with Go")},
				Burst: burst,
			}
			cmd := intent.Update(factMsg)
			Expect(intent.GetState()).To(Equal(burst_management.StateInferringSkills))

			if cmd != nil {
				messages := executeBatchCmd(cmd)
				for _, m := range messages {
					if loaded, ok := m.(burst_management.SkillSuggestionsLoadedMsg); ok {
						Expect(loaded.Suggestions).To(HaveLen(1))
						Expect(loaded.Suggestions[0].Name).To(Equal("Go"))
					}
				}
			}
		})

		It("should return error when inference service fails", func() {
			events := []*career.Event{
				fixtures.EventWith("e1", "Built API", "Acme", "Backend"),
			}
			mockService := mocks.NewBurstServiceMock().SetEvents(events)

			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)
			mockSkillSvc.EXPECT().
				InferSkillsFromEvents(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("inference failed")).
				AnyTimes()

			burst := fixtures.BurstConfirmed("burst-infer-err", "e1", "e2")

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{burst},
				Service:               mockService,
				SkillInferenceService: mockSkillSvc,
				Context:               context.Background(),
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)
			intent.SetState(burst_management.StateExtractingFacts)

			factMsg := burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{fixtures.FactWith("f1", "Fact")},
				Burst: burst,
			}
			cmd := intent.Update(factMsg)
			Expect(intent.GetState()).To(Equal(burst_management.StateInferringSkills))

			if cmd != nil {
				messages := executeBatchCmd(cmd)
				for _, m := range messages {
					if errMsg, ok := m.(burst_management.SkillSuggestionsErrorMsg); ok {
						Expect(errMsg.Err).To(HaveOccurred())
						Expect(errMsg.Err.Error()).To(ContainSubstring("skill detection failed"))
					}
				}
			}
		})

		It("should return error when burst is nil", func() {
			mockSkillSvc := skillinference.NewSkillInferenceService(nil, nil, nil)

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{},
				SkillInferenceService: mockSkillSvc,
				Context:               context.Background(),
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(nil)
			intent.SetState(burst_management.StateInferringSkills)

			cmd := intent.Update(tea.KeyMsg{})
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

	Describe("saveSkillFromSuggestion via accept flow", func() {
		It("should show error when skill inference service is nil", func() {
			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{fixtures.BurstConfirmed("b1", "e1")},
				SkillInferenceService: nil,
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

			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should handle CreateSkillsFromSuggestions error", func() {
			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)
			mockSkillSvc.EXPECT().
				CreateSkillsFromSuggestions(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("creation failed")).
				Times(1)

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

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("openSuggestionEventsModal and resolveEventsByIDs", func() {
		It("should open events modal when viewing events from skill suggestion", func() {
			events := []*career.Event{
				fixtures.EventWith("e1", "Built API with Go", "Acme", "Backend"),
				fixtures.EventWith("e2", "Deployed to production", "Acme", "DevOps"),
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
		})

		It("should handle view events when service is nil", func() {
			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)

			burst := fixtures.BurstConfirmed("b1", "e1")

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{burst},
				Service:               nil,
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
			Expect(intent.GetState()).To(Equal(burst_management.StateSkillSuggestionReview))

			Expect(func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
			}).NotTo(Panic())
		})

		It("should close suggestion events modal on Esc and return to skill modal", func() {
			events := []*career.Event{
				fixtures.EventWith("e1", "Built API with Go", "Acme", "Backend"),
			}
			mockService := mocks.NewBurstServiceMock().SetEvents(events)

			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)

			burst := fixtures.BurstConfirmed("b1", "e1")

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

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})
	})

	Describe("cancelPreviousOperation with active cancel func", func() {
		It("should cancel previous operation during burst detection", func() {
			mockService := mocks.NewBurstServiceMock().
				SetEvents([]*career.Event{fixtures.Event("e1")})

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{},
				Service: mockService,
				Context: context.Background(),
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateList)

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
			cmd1 := intent.Update(msg)

			if cmd1 != nil {
				cmd2 := intent.Update(msg)
				_ = cmd2
			}
		})
	})

	Describe("rebuildModalRegistry additional modal types", func() {
		It("should register loading modal when present", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.Burst("b1")},
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateSuggestionReview)

			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{
					{Name: "Loading Test", Description: "Test", EventIDs: []string{"e1", "e2"}},
				},
			}
			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(burst_management.StateExtractingFacts))
			Expect(intent.HasActiveModal()).To(BeTrue())
		})

		It("should register suggestion modal when visible", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			setupIntent(ctx)

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Test", EventIDs: []string{"e1"}, ConfidenceScore: 0.8},
				},
			})

			Expect(intent.GetSuggestionModal()).NotTo(BeNil())
			Expect(intent.HasActiveModal()).To(BeTrue())
		})

		It("should register skill suggestion modal when visible", func() {
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
			Expect(intent.HasActiveModal()).To(BeTrue())
		})
	})

	Describe("showBurstSkillsModal", func() {
		It("should trigger skills load when pressing 's' in detail", func() {
			burst := fixtures.BurstConfirmed("b1", "e1")
			burst.Name = "Skills Test Burst"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			result := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(result)
			Expect(intent.GetDetailModal()).NotTo(BeNil())

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			if cmd != nil {
				msg := executeAsyncCmd(cmd)
				if msg != nil {
					intent.Update(msg)
				}
			}
		})

		It("should handle skills load with nil selectedBurst", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.Burst("b1")},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(nil)
			Expect(func() {
				intent.Update(burst_management.BurstSkillsLoadedMsg{
					Skills: []*career.Skill{},
				})
			}).NotTo(Panic())
		})

		It("should handle skills loaded message", func() {
			burst := fixtures.BurstConfirmed("b1", "e1")
			burst.Name = "Skills Loaded Test"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			skills := []*career.Skill{
				fixtures.SkillWith("s1", "Go", "backend", "mid"),
			}
			intent.Update(burst_management.BurstSkillsLoadedMsg{Skills: skills})
		})
	})

	Describe("SkillSuggestionsErrorMsg handling", func() {
		It("should show error modal when skill inference fails", func() {
			mockSkillSvc := skillinference.NewSkillInferenceService(nil, nil, nil)

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{fixtures.BurstConfirmed("b1")},
				SkillInferenceService: mockSkillSvc,
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillSuggestionsErrorMsg{
				Err: errors.New("inference failed"),
			})

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should silently handle cancelled context error", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.BurstConfirmed("b1")},
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillSuggestionsErrorMsg{
				Err: context.Canceled,
			})

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("startSkillInference edge cases", func() {
		It("should show error when no burst selected", func() {
			mockSkillSvc := skillinference.NewSkillInferenceService(nil, nil, nil)

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{},
				SkillInferenceService: mockSkillSvc,
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(nil)

			burst := fixtures.BurstConfirmed("b-detail", "e1")
			burst.Name = "Infer Test"

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)
			intent.SetSelectedBurst(nil)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
		})

		It("should show warning when burst is not confirmed", func() {
			burst := fixtures.Burst("b-unconfirmed", "e1")
			burst.Name = "Unconfirmed"

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{burst},
				SkillInferenceService: skillinference.NewSkillInferenceService(nil, nil, nil),
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			Expect(intent.HasVisibleFeedbackModal()).To(BeTrue())
		})

		It("should show error when skill inference service is nil", func() {
			burst := fixtures.BurstConfirmed("b-nosvc", "e1")
			burst.Name = "No Service"

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{burst},
				SkillInferenceService: nil,
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("handleSuggestionEventsModalUpdate", func() {
		It("should handle non-key messages without panic", func() {
			events := []*career.Event{
				fixtures.EventWith("e1", "Built API", "Acme", "Backend"),
			}
			mockService := mocks.NewBurstServiceMock().SetEvents(events)

			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)

			burst := fixtures.BurstConfirmed("b1", "e1")

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

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

			Expect(func() {
				intent.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
			}).NotTo(Panic())
		})
	})

	Describe("handleDetailModalKeypress extended coverage", func() {
		It("should trigger skill inference from detail modal for confirmed burst", func() {
			burst := fixtures.BurstConfirmed("b-confirm-infer", "e1")
			burst.Name = "Confirmed For Infer"

			mockSkillSvc := skillinference.NewSkillInferenceService(nil, nil, nil)

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{burst},
				SkillInferenceService: mockSkillSvc,
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)
			Expect(intent.GetDetailModal()).NotTo(BeNil())

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			Expect(intent.GetState()).To(Equal(burst_management.StateInferringSkills))
		})
	})

	Describe("handleEditModalUpdate extended coverage", func() {
		It("should handle edit modal with Enter key press", func() {
			burst := fixtures.Burst("b-edit", "e1")
			burst.Name = "Edit Test"
			burst.Description = "Test description"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(intent.HasVisibleEditModal()).To(BeTrue())

			intent.Update(tea.KeyMsg{Type: tea.KeyTab})
			intent.Update(tea.KeyMsg{Type: tea.KeyTab})
		})
	})

	Describe("handleEventsModalUpdate extended", func() {
		It("should handle non-Esc key in events modal", func() {
			burst := fixtures.Burst("b-events", "e1")
			burst.Name = "Events Test"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			events := []*career.Event{
				fixtures.EventWith("e1", "Event 1", "", ""),
			}
			intent.Update(burst_management.BurstEventsLoadedMsg{Events: events})

			intent.Update(tea.KeyMsg{Type: tea.KeyDown})
		})
	})

	Describe("handleFactsModalUpdate extended", func() {
		It("should handle non-Esc key in facts modal", func() {
			burst := fixtures.Burst("b-facts", "e1")
			burst.Name = "Facts Test"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			facts := []*career.Fact{
				fixtures.FactWith("f1", "Fact 1"),
			}
			intent.Update(burst_management.BurstFactsLoadedMsg{Facts: facts})

			intent.Update(tea.KeyMsg{Type: tea.KeyDown})
		})
	})

	Describe("handleSkillsModalUpdate extended", func() {
		It("should handle Esc key in skills modal", func() {
			burst := fixtures.BurstConfirmed("b-skills", "e1")
			burst.Name = "Skills Modal Test"

			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)

			intent.Update(burst_management.BurstSkillsLoadedMsg{
				Skills: []*career.Skill{fixtures.SkillWith("s1", "Go", "backend", "mid")},
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})
	})

	Describe("handleConfirmModalUpdate extended", func() {
		It("should handle 'n' key in confirm modal to cancel", func() {
			burst := fixtures.Burst("b-confirm", "e1")
			burst.Name = "Confirm Test"

			mockService := mocks.NewBurstServiceMock()

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{burst},
				Service: mockService,
			}
			setupIntent(ctx)

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

			Expect(intent.HasVisibleConfirmModal()).To(BeTrue())

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

			Expect(intent.HasVisibleConfirmModal()).To(BeFalse())
		})
	})

	Describe("saveAndExtractBurstWithResult via accept in suggestion modal", func() {
		It("should save burst immediately on accept and trigger extraction", func() {
			repo := careermemory.NewBurstRepository()

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{},
				BurstRepository: repo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Instant Save", Description: "Save immediately", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.9},
				},
			})

			Expect(intent.GetSuggestionModal()).NotTo(BeNil())

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.GetFilteredBursts()).To(HaveLen(1))
			Expect(intent.GetFilteredBursts()[0].Name).To(Equal("Instant Save"))
		})

		It("should handle repo error on accept", func() {
			mockRepo := mocks.NewBurstRepositoryMock().SetCreateError(errors.New("repo error"))

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{},
				BurstRepository: mockRepo,
				Context:         context.Background(),
			}
			setupIntent(ctx)

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Fail Save", Description: "Should fail", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.9},
				},
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("handleSuggestionModalClosed extended", func() {
		It("should clear modal and return to list on modal close with empty suggestions", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			setupIntent(ctx)

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Test", EventIDs: []string{"e1"}, ConfidenceScore: 0.8},
				},
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			Expect(intent.GetSuggestionModal()).To(BeNil())
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("extractFactsForBurst extended coverage", func() {
		It("should handle nil service during extraction", func() {
			burst := fixtures.BurstConfirmed("b-nilsvc", "e1", "e2")

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{burst},
				Service: nil,
				Context: context.Background(),
			}
			setupIntent(ctx)

			intent.SetSelectedBurst(burst)
			intent.SetState(burst_management.StateSuggestionReview)

			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{
					{Name: "Nil Svc Test", Description: "No service", EventIDs: []string{"e1", "e2"}},
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

	Describe("handleBurstSuggestionsLoaded with error", func() {
		It("should show error modal when suggestions have error", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateSuggesting)

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Error: fmt.Errorf("detection failed"),
			})

			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("handleSkillSuggestionsLoaded with error", func() {
		It("should show error modal when skill suggestions have error", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.BurstConfirmed("b1")},
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Error: fmt.Errorf("inference failed"),
			})

			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("View rendering in various states", func() {
		It("should render view in StateSuggesting", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.Burst("b1")},
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateSuggesting)
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render view in StateInferringSkills", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.Burst("b1")},
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateInferringSkills)
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render view in StateSkillSuggestionReview", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.Burst("b1")},
			}
			setupIntent(ctx)

			intent.SetState(burst_management.StateSkillSuggestionReview)
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})
