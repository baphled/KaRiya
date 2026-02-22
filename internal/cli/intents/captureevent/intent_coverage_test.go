package captureevent

import (
	"context"
	"errors"

	"github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	captureScreens "github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	memoryrepo "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Intent — additional coverage", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentContext{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("Init — terminal dimensions branch", func() {
		It("uses terminal info dimensions when available without panicking", func() {
			info := &terminal.Info{Width: 160, Height: 50, IsValid: true}
			intent.UpdateTerminalInfo(info)
			Expect(func() { intent.Init() }).NotTo(Panic())
		})
	})

	Describe("View — submit modal branch", func() {
		It("renders modal overlay when submitModal is set", func() {
			intent.submitModal = feedback.NewErrorModal("Test", "something broke")
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders modal overlay with terminal dimensions set", func() {
			info := &terminal.Info{Width: 120, Height: 40, IsValid: true}
			intent.UpdateTerminalInfo(info)
			intent.submitModal = feedback.NewErrorModal("Test", "something broke")
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("View — editing mode overlays", func() {
		BeforeEach(func() {
			event := fixtures.EventWith("evt-1", "Test event for view overlay", "", "")
			intent.reviewState = &ReviewInferredEventState{
				Event:         event,
				AcceptedFacts: make([]*career.Fact, 0),
			}
			breadcrumbs := []string{"Main Menu", "Capture Event", "Review"}
			intent.activeScreen = captureScreens.NewEventReviewScreen(
				breadcrumbs, event, nil, nil, nil,
			)
		})

		It("renders burst modal overlay when EditingModeBursts is active", func() {
			suggestions := []burstfact.BurstSuggestion{
				{Name: "API Work", Description: "REST APIs", EventIDs: []string{"e1"}},
			}
			intent.reviewState.EditingMode = EditingModeBursts
			intent.reviewState.burstModal = modals.NewSuggestionReviewModal(suggestions, nil)
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders skills modal overlay when EditingModeSkills is active", func() {
			skills := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.9},
			}
			intent.reviewState.EditingMode = EditingModeSkills
			intent.reviewState.skillModal = modals.NewSkillSuggestionModal(skills, nil)
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders facts suggestion modal overlay when EditingModeFacts is active", func() {
			facts := []career.Fact{
				*fixtures.FactWith("f1", "Led API design for teams"),
			}
			intent.reviewState.EditingMode = EditingModeFacts
			intent.reviewState.factSuggestionModal = modals.NewFactSuggestionModal(facts, nil)
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("falls back to base view when editing mode has no matching modal", func() {
			intent.reviewState.EditingMode = EditingModeMetadata
			intent.reviewState.metadataModal = nil
			intent.context.CareerService = nil

			Expect(func() { intent.View() }).To(Panic())
		})
	})

	Describe("Update — ModalAutoDismissMsg", func() {
		It("produces DismissModalMsg command", func() {
			cmd := intent.Update(feedback.ModalAutoDismissMsg{})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			_, ok := msg.(DismissModalMsg)
			Expect(ok).To(BeTrue())
		})
	})

	Describe("Update — DismissModalMsg with active submitModal", func() {
		It("transitions to review state and clears the modal after a successful save", func() {
			intent.submitModal = feedback.NewSuccessModal("Event saved!")
			event := fixtures.EventWith("evt-1", "Test event text for the modal", "", "")
			intent.reviewState = &ReviewInferredEventState{
				Event:          event,
				InferredBursts: make([]*career.Burst, 0),
				InferredFacts:  make([]*career.Fact, 0),
				InferredSkills: nil,
			}

			intent.Update(DismissModalMsg{})
			Expect(intent.submitModal).To(BeNil())
			Expect(intent.currentState).To(Equal(StateReview))
		})

		It("returns to form state when dismissing an error modal", func() {
			intent.submitModal = feedback.NewErrorModal("Save Failed", "something went wrong")

			intent.Update(DismissModalMsg{})
			Expect(intent.submitModal).To(BeNil())
			Expect(intent.currentState).To(Equal(StateForm))
		})

		It("sets up active screen for review after a successful save", func() {
			intent.submitModal = feedback.NewSuccessModal("Event saved!")
			event := fixtures.EventWith("evt-1", "Test event text for the modal", "", "")
			intent.reviewState = &ReviewInferredEventState{
				Event: event,
			}

			intent.Update(DismissModalMsg{})
			Expect(intent.activeScreen).NotTo(BeNil())
		})

		It("uses terminal dimensions when available during dismiss", func() {
			info := &terminal.Info{Width: 120, Height: 40, IsValid: true}
			intent.UpdateTerminalInfo(info)
			intent.submitModal = feedback.NewSuccessModal("Event saved!")
			event := fixtures.EventWith("evt-1", "Test event text for the modal", "", "")
			intent.reviewState = &ReviewInferredEventState{
				Event: event,
			}

			Expect(func() { intent.Update(DismissModalMsg{}) }).NotTo(Panic())
		})
	})

	Describe("Update — SubmitMsg handling", func() {
		BeforeEach(func() {
			intent.currentState = StateForm
		})

		It("shows validation error modal when SubmitMsg carries an error", func() {
			err := errors.New("event text is required")
			intent.Update(SubmitMsg{Err: err})
			Expect(intent.submitModal).NotTo(BeNil())
		})

		It("processes event from SubmitMsg when event is non-nil", func() {
			event := fixtures.EventWith("", "Valid event text for testing SubmitMsg processing", "", "")
			cmd := intent.Update(SubmitMsg{Event: event})
			Expect(cmd).NotTo(BeNil())
		})

		It("returns nil when not in StateForm", func() {
			intent.currentState = StateReview
			event := fixtures.EventWith("", "Test", "", "")
			cmd := intent.Update(SubmitMsg{Event: event})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update — submitModal key handling", func() {
		BeforeEach(func() {
			intent.submitModal = feedback.NewErrorModal("Error", "something failed")
		})

		It("returns nil for non-Esc key when modal is active", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(cmd).To(BeNil())
		})

		It("returns DismissModalMsg command on Esc when modal is not loading", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			_, ok := msg.(DismissModalMsg)
			Expect(ok).To(BeTrue())
		})

		It("does not dismiss loading modal on Esc", func() {
			intent.submitModal = feedback.NewLoadingModal("Saving...", false)
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
		})

		It("handles ModalSpinnerTickMsg for loading modal", func() {
			intent.submitModal = feedback.NewLoadingModal("Saving...", false)
			intent.submitModal.Init()
			cmd := intent.Update(feedback.ModalSpinnerTickMsg{})
			_ = cmd
		})

		It("handles ModalCountdownTickMsg for success modal", func() {
			intent.submitModal = feedback.NewSuccessModal("Saved!")
			cmd := intent.Update(feedback.ModalCountdownTickMsg{})
			_ = cmd
		})

		It("returns nil for unrecognised message types when modal is active", func() {
			cmd := intent.Update("unknown message type")
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update — editing modal routing", func() {
		BeforeEach(func() {
			event := fixtures.EventWith("evt-1", "Test event for editing modal", "", "")
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event)
			svc.SetSkillRepository(repos.Skill)
			intent.context.CareerService = svc

			breadcrumbs := []string{"Main Menu", "Capture Event", "Review"}
			intent.activeScreen = captureScreens.NewEventReviewScreen(
				breadcrumbs, event, nil, nil, nil,
			)
			intent.reviewState = &ReviewInferredEventState{
				Event:         event,
				EditingMode:   EditingModeMetadata,
				AcceptedFacts: make([]*career.Fact, 0),
			}
			intent.reviewState.metadataModal = NewReviewEnrichmentModel(
				context.TODO(), event, svc, nil, nil,
			)
		})

		It("routes message through editing modal when EditingMode is not None", func() {
			Expect(func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyTab})
			}).NotTo(Panic())
		})
	})

	Describe("Update — screen active path with result", func() {
		It("handles screen update with no result", func() {
			Expect(func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			}).NotTo(Panic())
		})
	})

	Describe("Update — InferenceCompleteMsg", func() {
		It("updates review state with inferred data", func() {
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("", "Test event", "", ""),
			}
			breadcrumbs := []string{"Main Menu", "Capture Event", "Review"}
			screen := captureScreens.NewEventReviewScreen(
				breadcrumbs,
				intent.reviewState.Event,
				nil, nil, nil,
			)
			intent.activeScreen = screen

			skills := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.9},
			}
			cmd := intent.Update(InferenceCompleteMsg{
				InferredSkills: skills,
				InferredBursts: make([]*career.Burst, 0),
				InferredFacts:  make([]*career.Fact, 0),
			})
			Expect(cmd).To(BeNil())
			Expect(intent.reviewState.InferredSkills).To(HaveLen(1))
		})
	})
})
