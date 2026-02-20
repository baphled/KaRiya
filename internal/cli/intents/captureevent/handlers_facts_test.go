package captureevent

import (
	"context"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	captureScreens "github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/domain/career"
	memoryrepo "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
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

	Describe("fact persistence in postSaveReview", func() {
		var (
			svc      *careerservice.Service
			factRepo *memoryrepo.FactRepository
		)

		BeforeEach(func() {
			repos := memoryrepo.NewRepositories()
			factRepo = repos.Fact.(*memoryrepo.FactRepository)

			svc = careerservice.NewService(repos.Event)
			svc.SetFactRepository(factRepo)
			svc.SetSkillRepository(repos.Skill)
			svc.SetBurstRepository(repos.Burst)

			intent.context.CareerService = svc
		})

		Context("with accepted facts that have empty IDs", func() {
			It("sets SourceEventID from the saved event and persists via SaveFact", func() {
				event := fixtures.EventWith("evt-saved-1", "Reviewed event with facts", "", "")
				fact := fixtures.FactForSave(
					"Led cross-team API design",
					[]string{"leadership"},
					career.RoleFitStaff,
					[]string{"hiring_manager"},
					"",
				)

				intent.currentState = StateReview
				intent.postSaveReview = true
				intent.reviewState = &ReviewInferredEventState{
					Event: event,
				}

				reviewData := map[string]interface{}{
					"event":  event,
					"bursts": []*career.Burst{},
					"facts":  []*career.Fact{fact},
				}
				intent.HandleSubmit(&screens.SubmitResult{FormData: reviewData})

				Expect(intent.result).NotTo(BeNil())
				Expect(intent.result.Status).To(Equal(intents.Completed))
				Expect(fact.SourceEventID).To(Equal("evt-saved-1"))
				Expect(fact.ID).NotTo(BeEmpty())

				ctx := context.Background()
				saved, err := factRepo.GetByID(ctx, fact.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(saved).NotTo(BeNil())
				Expect(saved.SourceEventID).To(Equal("evt-saved-1"))
			})
		})

		Context("with facts that already have an ID", func() {
			It("skips facts with existing IDs", func() {
				event := fixtures.EventWith("evt-saved-2", "Reviewed event", "", "")
				existingFact := fixtures.Fact("fact-existing-1", "evt-saved-2")

				ctx := context.Background()
				Expect(factRepo.Create(ctx, existingFact)).To(Succeed())

				intent.currentState = StateReview
				intent.postSaveReview = true
				intent.reviewState = &ReviewInferredEventState{
					Event: event,
				}

				reviewData := map[string]interface{}{
					"event":  event,
					"bursts": []*career.Burst{},
					"facts":  []*career.Fact{existingFact},
				}
				intent.HandleSubmit(&screens.SubmitResult{FormData: reviewData})

				Expect(intent.result).NotTo(BeNil())
				Expect(intent.result.Status).To(Equal(intents.Completed))
			})
		})

		Context("when CareerService is nil", func() {
			It("completes without persisting facts", func() {
				event := fixtures.EventWith("evt-saved-3", "Reviewed event", "", "")
				fact := fixtures.FactForSave(
					"Led API design",
					[]string{"leadership"},
					career.RoleFitStaff,
					[]string{"hiring_manager"},
					"",
				)

				intent.context.CareerService = nil
				intent.currentState = StateReview
				intent.postSaveReview = true
				intent.reviewState = &ReviewInferredEventState{
					Event: event,
				}

				reviewData := map[string]interface{}{
					"event":  event,
					"bursts": []*career.Burst{},
					"facts":  []*career.Fact{fact},
				}
				intent.HandleSubmit(&screens.SubmitResult{FormData: reviewData})

				Expect(intent.result).NotTo(BeNil())
				Expect(intent.result.Status).To(Equal(intents.Completed))
				Expect(fact.ID).To(BeEmpty())
			})
		})
	})
})
