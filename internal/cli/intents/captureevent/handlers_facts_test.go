package captureevent

import (
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	captureScreens "github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fact Suggestion in CaptureEvent", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentContext{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("smart routing via HandleNavigate", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "Led API design for payments platform", "", ""),
				InferredFacts: []*career.Fact{
					fixtures.FactWith("f1", "Led cross-team API design"),
					fixtures.FactWith("f2", "Improved system reliability"),
				},
				AcceptedFacts: []*career.Fact{},
			}
		})

		Context("when inferred facts exist", func() {
			It("creates a fact suggestion modal", func() {
				result := &screens.NavigateResult{ResultData: "suggest_facts"}
				intent.HandleNavigate(result)

				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeFacts))
				Expect(intent.reviewState.factSuggestionModal).NotTo(BeNil())
			})

			It("does not create the old fact editor modal", func() {
				result := &screens.NavigateResult{ResultData: "suggest_facts"}
				intent.HandleNavigate(result)

				Expect(intent.reviewState.factModal).To(BeNil())
			})

			It("sets terminal dimensions on the modal", func() {
				result := &screens.NavigateResult{ResultData: "suggest_facts"}
				intent.HandleNavigate(result)

				Expect(intent.reviewState.factSuggestionModal).NotTo(BeNil())
			})
		})

		Context("when no inferred facts exist", func() {
			BeforeEach(func() {
				intent.reviewState.InferredFacts = []*career.Fact{}
			})

			It("returns early without creating a modal", func() {
				result := &screens.NavigateResult{ResultData: "suggest_facts"}
				cmd := intent.HandleNavigate(result)

				Expect(cmd).To(BeNil())
				Expect(intent.reviewState.factSuggestionModal).To(BeNil())
				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
			})

			It("keeps intent active", func() {
				result := &screens.NavigateResult{ResultData: "suggest_facts"}
				intent.HandleNavigate(result)

				Expect(intent.IsActive()).To(BeTrue())
			})
		})

		Context("when InferredFacts is nil", func() {
			BeforeEach(func() {
				intent.reviewState.InferredFacts = nil
			})

			It("returns early without creating a modal", func() {
				result := &screens.NavigateResult{ResultData: "suggest_facts"}
				cmd := intent.HandleNavigate(result)

				Expect(cmd).To(BeNil())
				Expect(intent.reviewState.factSuggestionModal).To(BeNil())
			})
		})
	})

	Describe("fact acceptance via updateEditingModal", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:         fixtures.EventWith("evt-1", "Led API design", "", ""),
				AcceptedFacts: []*career.Fact{},
				EditingMode:   EditingModeFacts,
			}

			facts := []career.Fact{
				*fixtures.FactWith("f1", "Led cross-team API design"),
			}
			intent.reviewState.factSuggestionModal = modals.NewFactSuggestionModal(facts, nil)

			breadcrumbs := []string{"Test"}
			screen := captureScreens.NewEventReviewScreen(breadcrumbs, intent.reviewState.Event, nil, nil, nil)
			intent.activeScreen = screen
		})

		It("transfers accepted facts to review state as pointers", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.AcceptedFacts).To(HaveLen(1))
			Expect(intent.reviewState.AcceptedFacts[0].Text).To(Equal("Led cross-team API design"))
		})

		It("closes the modal and resets editing mode", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.factSuggestionModal).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
		})

		It("calls SetAcceptedFacts on the review screen", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			_, ok := intent.activeScreen.(*captureScreens.EventReviewScreen)
			Expect(ok).To(BeTrue())
			Expect(intent.reviewState.AcceptedFacts).To(HaveLen(1))
		})

		Context("with multiple facts", func() {
			BeforeEach(func() {
				facts := []career.Fact{
					*fixtures.FactWith("f1", "Led cross-team API design"),
					*fixtures.FactWith("f2", "Improved system reliability"),
				}
				intent.reviewState.factSuggestionModal = modals.NewFactSuggestionModal(facts, nil)
			})

			It("accepts all facts in sequence", func() {
				intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
				intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

				Expect(intent.reviewState.AcceptedFacts).To(HaveLen(2))
				Expect(intent.reviewState.AcceptedFacts[0].Text).To(Equal("Led cross-team API design"))
				Expect(intent.reviewState.AcceptedFacts[1].Text).To(Equal("Improved system reliability"))
			})
		})

		Context("when all facts are rejected", func() {
			It("clears modal without adding facts", func() {
				intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

				Expect(intent.reviewState.factSuggestionModal).To(BeNil())
				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
				Expect(intent.reviewState.AcceptedFacts).To(BeEmpty())
			})
		})

		Context("when escape is pressed", func() {
			It("clears the fact suggestion modal", func() {
				intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.reviewState.factSuggestionModal).To(BeNil())
				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
			})
		})
	})
})
