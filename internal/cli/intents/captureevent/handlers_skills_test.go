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

var _ = Describe("Skills Handling (Internal)", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentContext{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("HandleNavigate with edit_skills", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "test", "", ""),
				InferredSkills: []skillinference.SkillSuggestion{
					{Name: "Go", Confidence: 0.9},
				},
				AcceptedSkills: []*career.Skill{},
			}
		})

		It("should initialize skill modal and set editing mode", func() {
			result := &screens.NavigateResult{ResultData: "edit_skills"}
			cmd := intent.HandleNavigate(result)

			// skillModal.Init() currently returns nil, so we expect nil
			Expect(cmd).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeSkills))
			Expect(intent.reviewState.skillModal).NotTo(BeNil())
		})
	})

	Describe("updateEditingModal with EditingModeSkills", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "test", "", ""),
				AcceptedSkills: []*career.Skill{},
				EditingMode:    EditingModeSkills,
			}
			suggestions := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.9},
			}
			intent.reviewState.skillModal = modals.NewSkillSuggestionModal(suggestions, nil)

			// Mock active screen to verify it gets updated
			breadcrumbs := []string{"Test"}
			screen := captureScreens.NewEventReviewScreen(breadcrumbs, intent.reviewState.Event, nil, nil, nil)
			intent.activeScreen = screen
		})

		It("should handle modal interaction and close", func() {
			// Simulate 'a' key (accept) which should accept "Go"
			// This will accept the single item, clear the list, set visible=false
			// Then updateEditingModal will detect !IsVisible(), transfer skills, and set skillModal=nil
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			// Since modal auto-closes when empty, skillModal should be nil now
			Expect(intent.reviewState.skillModal).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))

			// Verify skills were transferred to reviewState
			Expect(intent.reviewState.AcceptedSkills).To(HaveLen(1))
			Expect(intent.reviewState.AcceptedSkills[0].Name).To(Equal("Go"))
		})
	})
})
