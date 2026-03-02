package burst_management_test

import (
	"context"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/testutil/mocks"
	mockintent "github.com/baphled/kariya/internal/testutil/mocks/intent"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
)

var _ = Describe("Handler Coverage", func() {
	Describe("cancelAsyncOperation", func() {
		var (
			intent      *burst_management.Intent
			mockService *mocks.BurstServiceMock
		)

		BeforeEach(func() {
			events := []*career.Event{
				fixtures.EventWith("e1", "Event 1", "", ""),
				fixtures.EventWith("e2", "Event 2", "", ""),
				fixtures.EventWith("e3", "Event 3 unassigned", "", ""),
				fixtures.EventWith("e4", "Event 4 unassigned", "", ""),
			}

			suggestions := []burstfact.BurstSuggestion{
				{Name: "Suggested", EventIDs: []string{"e3", "e4"}, ConfidenceScore: 0.8},
			}

			mockService = mocks.NewBurstServiceMock().SetEvents(events).SetSuggestions(suggestions)

			ctx := &burst_management.IntentContext{
				Bursts:  []*career.Burst{},
				Service: mockService,
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should cancel async operation and reset state when Esc pressed during loading", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.GetLoadingModal()).NotTo(BeNil())

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.GetLoadingModal()).To(BeNil())
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should consume non-Esc keys while loading modal is active", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(intent.GetLoadingModal()).NotTo(BeNil())

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(cmd).NotTo(BeNil())
			Expect(intent.GetLoadingModal()).NotTo(BeNil())
		})

		It("should forward spinner tick to loading modal", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(intent.GetLoadingModal()).NotTo(BeNil())

			cmd := intent.Update(feedback.ModalSpinnerTickMsg{})

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("handleBurstSkillsLoaded", func() {
		var (
			intent *burst_management.Intent
			burst  *career.Burst
		)

		BeforeEach(func() {
			burst = fixtures.Burst("burst-1", "e1", "e2")
			burst.Name = "Skills Test Burst"
			burst.Description = "Test"

			skillRepo := careermemory.NewSkillRepository()

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				SkillRepository: skillRepo,
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			termInfo := terminal.NewInfo()
			termInfo.Width = 120
			termInfo.Height = 40
			termInfo.IsValid = true
			intent.UpdateTerminalInfo(termInfo)

			intent.Init()
		})

		It("should show skills modal on successful load", func() {
			intent.SetSelectedBurst(burst)

			skills := []*career.Skill{
				fixtures.SkillWith("s1", "Go", "Backend", "mid"),
				fixtures.SkillWith("s2", "Docker", "DevOps", "mid"),
			}

			msg := burst_management.BurstSkillsLoadedMsg{
				Skills: skills,
			}

			intent.Update(msg)

			Expect(intent.HasActiveModal()).To(BeTrue())
		})

		It("should show error modal on skills load error", func() {
			intent.SetSelectedBurst(burst)

			msg := burst_management.BurstSkillsLoadedMsg{
				Error: fmt.Errorf("failed to load skills"),
			}

			intent.Update(msg)

			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should handle nil selectedBurst gracefully", func() {
			intent.SetSelectedBurst(nil)

			msg := burst_management.BurstSkillsLoadedMsg{
				Skills: []*career.Skill{
					fixtures.SkillWith("s1", "Go", "Backend", "mid"),
				},
			}

			Expect(func() {
				intent.Update(msg)
			}).NotTo(Panic())
		})
	})

	Describe("handleSkillsModalUpdate", func() {
		var (
			intent *burst_management.Intent
			burst  *career.Burst
		)

		BeforeEach(func() {
			burst = fixtures.Burst("burst-1", "e1", "e2")
			burst.Name = "Skills Modal Test"
			burst.Description = "Test"

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				SkillRepository: careermemory.NewSkillRepository(),
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			termInfo := terminal.NewInfo()
			termInfo.Width = 120
			termInfo.Height = 40
			termInfo.IsValid = true
			intent.UpdateTerminalInfo(termInfo)

			intent.Init()
			intent.SetSelectedBurst(burst)
		})

		It("should close skills modal and return to detail on Esc", func() {
			skills := []*career.Skill{
				fixtures.SkillWith("s1", "Go", "Backend", "mid"),
			}
			intent.Update(burst_management.BurstSkillsLoadedMsg{Skills: skills})
			Expect(intent.HasActiveModal()).To(BeTrue())

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.GetDetailModal()).NotTo(BeNil())
		})

		It("should close skills modal and return to detail on Enter", func() {
			skills := []*career.Skill{
				fixtures.SkillWith("s1", "Go", "Backend", "mid"),
			}
			intent.Update(burst_management.BurstSkillsLoadedMsg{Skills: skills})
			Expect(intent.HasActiveModal()).To(BeTrue())

			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(intent.GetDetailModal()).NotTo(BeNil())
		})

		It("should return noopCmd with nil selectedBurst on close", func() {
			skills := []*career.Skill{
				fixtures.SkillWith("s1", "Go", "Backend", "mid"),
			}
			intent.Update(burst_management.BurstSkillsLoadedMsg{Skills: skills})

			intent.SetSelectedBurst(nil)

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("handleSkillSuggestionsError", func() {
		var intent *burst_management.Intent

		BeforeEach(func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.Burst("b1", "e1")},
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should show error modal and return to list state", func() {
			intent.SetState(burst_management.StateInferringSkills)

			msg := burst_management.SkillSuggestionsErrorMsg{
				Err: fmt.Errorf("skill inference failed"),
			}

			intent.Update(msg)

			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.GetLoadingModal()).To(BeNil())
		})

		It("should silently ignore cancelled operations", func() {
			intent.SetState(burst_management.StateInferringSkills)

			msg := burst_management.SkillSuggestionsErrorMsg{
				Err: context.Canceled,
			}

			intent.Update(msg)

			Expect(intent.GetFeedbackModal()).To(BeNil())
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("handleSkillSuggestionModalUpdate", func() {
		var (
			intent *burst_management.Intent
			burst  *career.Burst
		)

		BeforeEach(func() {
			burst = fixtures.BurstConfirmed("burst-1", "e1", "e2")
			burst.Name = "Skill Suggestion Test"
			burst.Description = "Test"

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{burst},
				SkillInferenceService: skillinference.NewSkillInferenceService(nil, nil, nil),
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			termInfo := terminal.NewInfo()
			termInfo.Width = 120
			termInfo.Height = 40
			termInfo.IsValid = true
			intent.UpdateTerminalInfo(termInfo)

			intent.Init()
		})

		It("should show skill suggestion modal when new suggestions exist", func() {
			intent.SetState(burst_management.StateInferringSkills)

			msg := burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
					{Name: "Docker", Category: "devops", Confidence: 0.85},
				},
				ExistingSkillNames: []string{},
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(burst_management.StateSkillSuggestionReview))
		})

		It("should close modal on Esc and return to list", func() {
			intent.SetState(burst_management.StateInferringSkills)

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

		It("should return to detail modal when inferredFromDetail and cancel pressed", func() {
			intent.SetSelectedBurst(burst)

			intent.SetState(burst_management.StateInferringSkills)
			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
				},
				ExistingSkillNames: []string{},
			})

			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			intent.SetState(burst_management.StateInferringSkills)
			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "React", Category: "frontend", Confidence: 0.90},
				},
				ExistingSkillNames: []string{},
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should handle reject action (Esc to reject last suggestion)", func() {
			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Python", Category: "backend", Confidence: 0.80},
				},
				ExistingSkillNames: []string{},
			})
			Expect(intent.GetState()).To(Equal(burst_management.StateSkillSuggestionReview))

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("handleSkillSuggestionAccept", func() {
		var (
			intent *burst_management.Intent
			burst  *career.Burst
		)

		BeforeEach(func() {
			burst = fixtures.BurstConfirmed("burst-1", "e1", "e2")
			burst.Name = "Accept Skill Test"
			burst.Description = "Test"

			ctrl := gomock.NewController(GinkgoT())
			mockSkillSvc := mockintent.NewMockBurstSkillInferenceService(ctrl)
			mockSkillSvc.EXPECT().CreateSkillsFromSuggestions(gomock.Any(), gomock.Any()).
				Return([]*career.Skill{fixtures.SkillWith("skill-1", "Go", "backend", "intermediate")}, nil).AnyTimes()

			ctx := &burst_management.IntentContext{
				Bursts:                []*career.Burst{burst},
				SkillInferenceService: mockSkillSvc,
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			termInfo := terminal.NewInfo()
			termInfo.Width = 120
			termInfo.Height = 40
			termInfo.IsValid = true
			intent.UpdateTerminalInfo(termInfo)

			intent.Init()
		})

		It("should accept skill suggestion and close modal", func() {
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
			modal := intent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalSuccess))
		})
	})

	Describe("handleScreenResult", func() {
		var intent *burst_management.Intent

		BeforeEach(func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.Burst("b1", "e1")},
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should handle nil result data gracefully", func() {
			result := &screens.NavigateResult{
				ResultData: nil,
			}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})

		It("should handle non-matching result data type gracefully", func() {
			result := &screens.NavigateResult{
				ResultData: "invalid-type",
			}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})
	})

	Describe("handleSuggestionEventsModalUpdate", func() {
		It("should return nil when no suggestion events modal exists", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.Burst("b1", "e1")},
			}
			ctx.Validate()

			intent, err := burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("handleBurstSuggestionsLoaded edge cases", func() {
		var intent *burst_management.Intent

		BeforeEach(func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			termInfo := terminal.NewInfo()
			termInfo.Width = 120
			termInfo.Height = 40
			termInfo.IsValid = true
			intent.UpdateTerminalInfo(termInfo)

			intent.Init()
		})

		It("should show warning when no suggestions found", func() {
			msg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{},
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			modal := intent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalWarning))
		})

		It("should silently ignore cancelled operations", func() {
			msg := burst_management.BurstSuggestionsLoadedMsg{
				Error: context.Canceled,
			}

			intent.Update(msg)

			Expect(intent.GetFeedbackModal()).To(BeNil())
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("FactExtractionCompleteMsg cancelled handling", func() {
		It("should silently ignore cancelled fact extraction", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{fixtures.Burst("b1", "e1")},
			}
			ctx.Validate()

			intent, err := burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			msg := burst_management.FactExtractionCompleteMsg{
				Error: context.Canceled,
			}

			intent.Update(msg)

			Expect(intent.GetFeedbackModal()).To(BeNil())
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should not show detail modal when suggestion modal is still visible during error", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			ctx.Validate()

			intent, err := burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			termInfo := terminal.NewInfo()
			termInfo.Width = 120
			termInfo.Height = 40
			termInfo.IsValid = true
			intent.UpdateTerminalInfo(termInfo)
			intent.Init()

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Test Suggestion", EventIDs: []string{"e1"}, ConfidenceScore: 0.9},
					{Name: "Test Suggestion 2", EventIDs: []string{"e2"}, ConfidenceScore: 0.8},
				},
			})
			Expect(intent.GetSuggestionModal()).NotTo(BeNil())

			msg := burst_management.FactExtractionCompleteMsg{
				Error: errors.New("extraction failed"),
			}
			intent.Update(msg)

			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should not change state when suggestion modal is visible during success", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			ctx.Validate()

			intent, err := burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			termInfo := terminal.NewInfo()
			termInfo.Width = 120
			termInfo.Height = 40
			termInfo.IsValid = true
			intent.UpdateTerminalInfo(termInfo)
			intent.Init()

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Suggestion A", EventIDs: []string{"e1"}, ConfidenceScore: 0.9},
					{Name: "Suggestion B", EventIDs: []string{"e2"}, ConfidenceScore: 0.8},
				},
			})
			Expect(intent.GetSuggestionModal()).NotTo(BeNil())

			msg := burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{fixtures.FactWith("f1", "Extracted")},
			}
			intent.Update(msg)

			Expect(intent.GetSuggestionModal().IsVisible()).To(BeTrue())
		})
	})

	Describe("Context LoadBursts with repository error", func() {
		It("should return error when repository fails", func() {
			mockRepo := mocks.NewBurstRepositoryMock().SetListError(errors.New("database unavailable"))

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{},
				BurstRepository: mockRepo,
				Context:         context.Background(),
			}

			err := ctx.LoadBursts()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("database unavailable"))
		})

		It("should update bursts list on successful load", func() {
			repo := careermemory.NewBurstRepository()
			_ = repo.Create(context.Background(), fixtures.Burst("b1", "e1"))
			_ = repo.Create(context.Background(), fixtures.Burst("b2", "e2"))

			ctx := &burst_management.IntentContext{
				Bursts:          []*career.Burst{},
				BurstRepository: repo,
				Context:         context.Background(),
			}

			err := ctx.LoadBursts()
			Expect(err).NotTo(HaveOccurred())
			Expect(ctx.Bursts).To(HaveLen(2))
		})
	})
})
