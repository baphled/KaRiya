package e2e_test

import (
	"context"

	burstmgmt "github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/domain/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E Skill Inference Workflow", func() {
	var intent *burstmgmt.Intent
	var ctx *burstmgmt.IntentContext
	var mockService *mocks.BurstServiceMock
	var burst *career.Burst

	BeforeEach(func() {
		burstRepo := careermemory.NewBurstRepository()
		mockService = mocks.NewBurstServiceMock()

		// Create repositories for skill inference
		skillRepo := careermemory.NewSkillRepository()
		eventRepo := careermemory.NewEventRepository()

		// Create skill inference service
		skillInferenceService := skillinference.NewSkillInferenceService(skillRepo, eventRepo)

		// Create a burst with events
		burst = &career.Burst{
			ID:          "burst-1",
			Name:        "API Development",
			Description: "Built REST API with Go and PostgreSQL",
			EventIDs:    []string{"event-1", "event-2"},
			Confirmed:   true,
		}

		ctx = &burstmgmt.IntentContext{
			Bursts:                []*career.Burst{burst},
			Service:               mockService,
			BurstRepository:       burstRepo,
			SkillInferenceService: skillInferenceService,
		}
		ctx.Validate()

		var err error
		intent, err = burstmgmt.NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("Auto Skill Inference Trigger", func() {
		It("should automatically infer skills after fact extraction completes on a confirmed burst", func() {
			// Given: A confirmed burst that was selected for fact extraction
			burst.Confirmed = true

			// Simulate burst selection by navigating to it (this sets selectedBurst)
			// This happens when user selects a burst from the list screen
			intent.HandleNavigate(&screens.NavigateResult{
				ResultData: burst,
			})

			// Create some extracted facts
			facts := []*career.Fact{
				{
					ID:                   "fact-1",
					Text:                 "Expert in Go programming",
					SourceEventID:        burst.EventIDs[0],
					CompetencyCategories: []string{"technical"},
				},
				{
					ID:                   "fact-2",
					Text:                 "PostgreSQL database design",
					SourceEventID:        burst.EventIDs[1],
					CompetencyCategories: []string{"technical"},
				},
			}

			// When: Fact extraction completes successfully
			intent.Update(burstmgmt.FactExtractionCompleteMsg{
				Facts: facts,
				Error: nil,
			})

			// Then: State should transition to StateInferringSkills
			Expect(intent.GetState()).To(Equal(burstmgmt.StateInferringSkills))

			// And: View should show skill inference in progress
			view := intent.View()
			Expect(view).To(ContainSubstring("Detecting skills from burst events"))
		})
	})

	Describe("Skill Inference Message Handling", func() {
		It("should display skill suggestions modal when skills are detected", func() {
			// Given: Skill inference returns suggestions
			suggestions := []skillinference.SkillSuggestion{
				{
					Name:       "Go",
					Category:   "Backend",
					Confidence: 0.95,
					EventIDs:   []string{burst.EventIDs[0]},
					Contexts:   []string{"Built API with Go"},
				},
				{
					Name:       "PostgreSQL",
					Category:   "Database",
					Confidence: 0.85,
					EventIDs:   []string{burst.EventIDs[0]},
					Contexts:   []string{"Integrated PostgreSQL database"},
				},
			}

			// When: SkillSuggestionsLoadedMsg is sent
			intent.Update(burstmgmt.SkillSuggestionsLoadedMsg{
				Suggestions: suggestions,
				Error:       nil,
			})

			// Then: State should be StateSkillSuggestionReview
			Expect(intent.GetState()).To(Equal(burstmgmt.StateSkillSuggestionReview))

			// And: View should show skill suggestion modal
			view := intent.View()
			Expect(view).To(ContainSubstring("Skill Suggestions"))
			Expect(view).To(ContainSubstring("Go"))
			Expect(view).To(ContainSubstring("PostgreSQL"))
			Expect(view).To(ContainSubstring("Backend"))
			Expect(view).To(ContainSubstring("Database"))
		})

		It("should show error modal when skill inference fails", func() {
			// When: SkillSuggestionsErrorMsg is sent with an error
			intent.Update(burstmgmt.SkillSuggestionsErrorMsg{
				Err: context.DeadlineExceeded,
			})

			// Then: State should return to list
			Expect(intent.GetState()).To(Equal(burstmgmt.StateList))

			// And: Error modal should be shown
			view := intent.View()
			Expect(view).To(ContainSubstring("Skill Inference Failed"))
		})

		It("should handle no skills detected gracefully", func() {
			// When: SkillSuggestionsLoadedMsg is sent with empty suggestions
			intent.Update(burstmgmt.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{},
				Error:       nil,
			})

			// Then: State should return to list
			Expect(intent.GetState()).To(Equal(burstmgmt.StateList))

			// And: Informative error should be shown
			view := intent.View()
			Expect(view).To(ContainSubstring("No skills were detected"))
		})

		It("should show success message when skills are created", func() {
			// Given: Skills have been created from suggestions
			createdSkills := []*career.Skill{
				{
					ID:       "skill-1",
					Name:     "Go",
					Category: "Backend",
				},
				{
					ID:       "skill-2",
					Name:     "PostgreSQL",
					Category: "Database",
				},
			}

			// When: SkillsCreatedMsg is sent
			intent.Update(burstmgmt.SkillsCreatedMsg{
				Skills: createdSkills,
				Error:  nil,
			})

			// Then: State should return to list
			Expect(intent.GetState()).To(Equal(burstmgmt.StateList))

			// And: Success message should be shown
			view := intent.View()
			Expect(view).To(ContainSubstring("Skills Created"))
			Expect(view).To(ContainSubstring("Successfully created 2 skill"))
		})

		It("should handle skill creation errors gracefully", func() {
			// When: SkillsCreatedMsg is sent with an error
			intent.Update(burstmgmt.SkillsCreatedMsg{
				Skills: nil,
				Error:  context.DeadlineExceeded,
			})

			// Then: State should return to list
			Expect(intent.GetState()).To(Equal(burstmgmt.StateList))

			// And: Error modal should be shown
			view := intent.View()
			Expect(view).To(ContainSubstring("Skill Creation Failed"))
		})

		It("should handle user cancellation silently", func() {
			// When: User cancels skill inference (context.Canceled)
			intent.Update(burstmgmt.SkillSuggestionsErrorMsg{
				Err: context.Canceled,
			})

			// Then: State should return to list silently
			Expect(intent.GetState()).To(Equal(burstmgmt.StateList))

			// And: No error modal should be shown for cancelled operations
			view := intent.View()
			Expect(view).NotTo(ContainSubstring("Skill Inference Failed"))
		})
	})
})
