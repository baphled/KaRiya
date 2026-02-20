package captureevent

import (
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	captureScreens "github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Skill Inference in CaptureEvent", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentContext{CaptureStrategy: "quick"}
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
		})

		It("initialises the skill modal with inferred skills", func() {
			result := &screens.NavigateResult{ResultData: "suggest_skills"}
			intent.HandleNavigate(result)

			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeSkills))
			Expect(intent.reviewState.skillModal).NotTo(BeNil())
		})

		Context("with no inferred skills", func() {
			BeforeEach(func() {
				intent.reviewState.InferredSkills = []skillinference.SkillSuggestion{}
			})

			It("opens the modal with an empty list", func() {
				result := &screens.NavigateResult{ResultData: "suggest_skills"}
				intent.HandleNavigate(result)

				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeSkills))
				Expect(intent.reviewState.skillModal).NotTo(BeNil())
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
			intent.reviewState.skillModal = modals.NewSkillSuggestionModal(suggestions, nil)

			breadcrumbs := []string{"Test"}
			screen := captureScreens.NewEventReviewScreen(breadcrumbs, intent.reviewState.Event, nil, nil, nil)
			intent.activeScreen = screen
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

			submitData := map[string]interface{}{
				"event":  intent.reviewState.Event,
				"bursts": []*career.Burst{},
				"facts":  []*career.Fact{},
				"skills": []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.9},
					{Name: "", Category: "unknown", Confidence: 0.1},
					{Name: "Python", Category: "backend", Confidence: 0.8},
				},
			}
			intent.HandleSubmit(&screens.SubmitResult{FormData: submitData})

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

			submitData := map[string]interface{}{
				"event":  intent.reviewState.Event,
				"bursts": []*career.Burst{},
				"facts":  []*career.Fact{},
				"skills": []skillinference.SkillSuggestion(nil),
			}
			intent.HandleSubmit(&screens.SubmitResult{FormData: submitData})

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Data.Skills).To(BeEmpty())
		})

		It("converts SkillSuggestion Category to career.Skill Category", func() {
			intent.currentState = StateSubmit
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "test", "", ""),
			}

			submitData := map[string]interface{}{
				"event":  intent.reviewState.Event,
				"bursts": []*career.Burst{},
				"facts":  []*career.Fact{},
				"skills": []skillinference.SkillSuggestion{
					{Name: "Docker", Category: "devops", Confidence: 0.85},
				},
			}
			intent.HandleSubmit(&screens.SubmitResult{FormData: submitData})

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

			submitData := map[string]interface{}{
				"event":  intent.reviewState.Event,
				"bursts": []*career.Burst{},
				"facts":  []*career.Fact{},
				"skills": []skillinference.SkillSuggestion{
					{Name: "", Category: "unknown", Confidence: 0.1},
					{Name: "Go", Category: "backend", Confidence: 0.9},
					{Name: "", Category: "unknown", Confidence: 0.2},
					{Name: "React", Category: "frontend", Confidence: 0.7},
					{Name: "", Category: "unknown", Confidence: 0.05},
				},
			}
			intent.HandleSubmit(&screens.SubmitResult{FormData: submitData})

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

		It("stores inferred skills from SubmitCompleteMsg in review state", func() {
			expectedSkills := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.95},
			}
			intent.Update(SubmitCompleteMsg{InferredSkills: expectedSkills})

			Expect(intent.reviewState.InferredSkills).To(Equal(expectedSkills))
		})

		It("handles empty inferred skills in SubmitCompleteMsg", func() {
			intent.Update(SubmitCompleteMsg{})

			Expect(intent.reviewState.InferredSkills).To(BeNil())
		})

		It("stores multiple inferred skills from SubmitCompleteMsg", func() {
			multipleSkills := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.95},
				{Name: "Docker", Category: "devops", Confidence: 0.88},
				{Name: "PostgreSQL", Category: "database", Confidence: 0.72},
			}
			intent.Update(SubmitCompleteMsg{InferredSkills: multipleSkills})

			Expect(intent.reviewState.InferredSkills).To(HaveLen(3))
			Expect(intent.reviewState.InferredSkills[0].Name).To(Equal("Go"))
			Expect(intent.reviewState.InferredSkills[2].Name).To(Equal("PostgreSQL"))
		})
	})

	Describe("error handling", func() {
		Context("when skill inference service is nil", func() {
			It("proceeds without inferring skills", func() {
				ctx := &IntentContext{CaptureStrategy: "quick"}
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

				submitData := map[string]interface{}{
					"event":  intent.reviewState.Event,
					"bursts": []*career.Burst{},
					"facts":  []*career.Fact{},
				}
				intent.HandleSubmit(&screens.SubmitResult{FormData: submitData})

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

				submitData := map[string]interface{}{
					"event":  intent.reviewState.Event,
					"bursts": []*career.Burst{},
					"facts":  []*career.Fact{},
					"skills": "invalid-type",
				}
				intent.HandleSubmit(&screens.SubmitResult{FormData: submitData})

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
				intent.reviewState.skillModal = modals.NewSkillSuggestionModal(suggestions, nil)

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
			intent.Update(SubmitCompleteMsg{InferredSkills: inferredSkills})

			Expect(intent.reviewState.InferredSkills).To(HaveLen(2))

			intent.currentState = StateSubmit
			submitData := map[string]interface{}{
				"event":  intent.reviewState.Event,
				"bursts": []*career.Burst{},
				"facts":  []*career.Fact{},
				"skills": intent.reviewState.InferredSkills,
			}
			intent.HandleSubmit(&screens.SubmitResult{FormData: submitData})

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Data.Skills).To(HaveLen(2))
			Expect(intent.result.Data.Skills[0].Name).To(Equal("Go"))
			Expect(intent.result.Data.Skills[0].Category).To(Equal("backend"))
			Expect(intent.result.Data.Skills[1].Name).To(Equal("Docker"))
			Expect(intent.result.Data.Skills[1].Category).To(Equal("devops"))
		})

		It("completes intent when review submits accepted skills as career.Skill", func() {
			intent.currentState = StateReview
			intent.postSaveReview = true
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "test", "", ""),
			}

			reviewData := map[string]interface{}{
				"event":  intent.reviewState.Event,
				"bursts": []*career.Burst{},
				"facts":  []*career.Fact{},
				"skills": []*career.Skill{
					fixtures.SkillWith("s-1", "Go", "backend", "advanced"),
					fixtures.SkillWith("s-2", "Docker", "devops", "intermediate"),
				},
			}
			intent.HandleSubmit(&screens.SubmitResult{FormData: reviewData})

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Data.Skills).To(HaveLen(2))
			Expect(intent.result.Data.Skills[0].Name).To(Equal("Go"))
			Expect(intent.result.Data.Skills[0].Category).To(Equal("backend"))
			Expect(intent.result.Data.Skills[0].Level).To(Equal("advanced"))
			Expect(intent.result.Data.Skills[1].Name).To(Equal("Docker"))
			Expect(intent.result.Data.Skills[1].Category).To(Equal("devops"))
			Expect(intent.result.Data.Skills[1].Level).To(Equal("intermediate"))
			Expect(intent.active).To(BeFalse())
		})
	})
})
