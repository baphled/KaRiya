package captureevent

import (
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	captureScreens "github.com/baphled/kariya/internal/cli/screens/capture"
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

var _ = Describe("Burst SuggestionReviewModal in CaptureEvent", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentContext{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("HandleNavigate suggest_bursts", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				InferredBursts: func() []*career.Burst {
					b := fixtures.Burst("", "evt-1", "evt-2")
					b.Name = "API Work"
					b.Description = "REST API development"
					return []*career.Burst{b}
				}(),
				AcceptedBursts: make([]*career.Burst, 0),
			}
		})

		Context("with CareerService nil", func() {
			It("marks intent as failed", func() {
				intent.context.CareerService = nil
				result := &screens.NavigateResult{ResultData: "suggest_bursts"}
				intent.HandleNavigate(result)

				Expect(intent.active).To(BeFalse())
			})
		})

		Context("with CareerService set", func() {
			BeforeEach(func() {
				repos := memoryrepo.NewRepositories()
				svc := careerservice.NewService(repos.Event)
				intent.context.CareerService = svc
				intent.reviewState.EditingMode = EditingModeNone
			})

			It("sets EditingMode to EditingModeBursts", func() {
				result := &screens.NavigateResult{ResultData: "suggest_bursts"}
				intent.HandleNavigate(result)

				Expect(intent.reviewState.EditingMode).To(Equal(EditingModeBursts))
			})

			It("creates a SuggestionReviewModal on burstModal", func() {
				result := &screens.NavigateResult{ResultData: "suggest_bursts"}
				intent.HandleNavigate(result)

				Expect(intent.reviewState.burstModal).NotTo(BeNil())
			})

			It("populates EventIDs from InferredBursts", func() {
				result := &screens.NavigateResult{ResultData: "suggest_bursts"}
				intent.HandleNavigate(result)

				Expect(intent.reviewState.burstModal).NotTo(BeNil())
				Expect(intent.reviewState.burstModal.GetSuggestionsCount()).To(Equal(1))
				suggestion := intent.reviewState.burstModal.GetCurrentSuggestion()
				Expect(suggestion).NotTo(BeNil())
				Expect(suggestion.EventIDs).To(Equal([]string{"evt-1", "evt-2"}))
			})
		})
	})

	Describe("updateEditingModal with EditingModeBursts using SuggestionReviewModal", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				AcceptedBursts: make([]*career.Burst, 0),
				EditingMode:    EditingModeBursts,
			}
			suggestions := []burstfact.BurstSuggestion{
				{Name: "API Development", Description: "Built REST APIs", EventIDs: []string{"evt-1"}, ConfidenceScore: 0.9},
			}
			intent.reviewState.burstModal = modals.NewSuggestionReviewModal(suggestions, nil)

			breadcrumbs := []string{"Test"}
			screen := captureScreens.NewEventReviewScreen(breadcrumbs, intent.reviewState.Event, nil, nil, nil)
			intent.activeScreen = screen
		})

		It("auto-closes modal and resets mode when no suggestions remain", func() {
			intent.reviewState.burstModal = modals.NewSuggestionReviewModal([]burstfact.BurstSuggestion{}, nil)

			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyDown})

			Expect(intent.reviewState.burstModal).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
		})

		It("transfers accepted bursts to review state on accept key", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.AcceptedBursts).To(HaveLen(1))
			Expect(intent.reviewState.AcceptedBursts[0].Name).To(Equal("API Development"))
			Expect(intent.reviewState.AcceptedBursts[0].Description).To(Equal("Built REST APIs"))
			Expect(intent.reviewState.AcceptedBursts[0].EventIDs).To(Equal([]string{"evt-1"}))
		})

		It("closes modal and resets editing mode after all accepted", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.burstModal).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
		})

		It("closes modal on escape without adding bursts", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.reviewState.burstModal).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
			Expect(intent.reviewState.AcceptedBursts).To(BeEmpty())
		})

		It("does not add rejected bursts to accepted list", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			Expect(intent.reviewState.AcceptedBursts).To(BeEmpty())
			Expect(intent.reviewState.burstModal).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
		})
	})

	Describe("burst suggestion with empty name", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "test event", "", ""),
				AcceptedBursts: make([]*career.Burst, 0),
				EditingMode:    EditingModeBursts,
			}
			suggestions := []burstfact.BurstSuggestion{
				{Name: "", Description: "unnamed burst", EventIDs: []string{"evt-1", "evt-2"}, ConfidenceScore: 0.7},
			}
			intent.reviewState.burstModal = modals.NewSuggestionReviewModal(suggestions, nil)

			breadcrumbs := []string{"Test"}
			screen := captureScreens.NewEventReviewScreen(breadcrumbs, intent.reviewState.Event, nil, nil, nil)
			intent.activeScreen = screen
		})

		It("generates a name based on event count when name is empty", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.AcceptedBursts).To(HaveLen(1))
			Expect(intent.reviewState.AcceptedBursts[0].Name).To(Equal("Burst of 2 events"))
		})
	})

	Describe("multiple burst suggestions: accept some reject some", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "test event", "", ""),
				AcceptedBursts: make([]*career.Burst, 0),
				EditingMode:    EditingModeBursts,
			}
			suggestions := []burstfact.BurstSuggestion{
				{Name: "Keep This", Description: "accepted", EventIDs: []string{"evt-1"}, ConfidenceScore: 0.9},
				{Name: "Skip This", Description: "rejected", EventIDs: []string{"evt-2"}, ConfidenceScore: 0.3},
			}
			intent.reviewState.burstModal = modals.NewSuggestionReviewModal(suggestions, nil)

			breadcrumbs := []string{"Test"}
			screen := captureScreens.NewEventReviewScreen(breadcrumbs, intent.reviewState.Event, nil, nil, nil)
			intent.activeScreen = screen
		})

		It("only includes accepted bursts in AcceptedBursts", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			Expect(intent.reviewState.AcceptedBursts).To(HaveLen(1))
			Expect(intent.reviewState.AcceptedBursts[0].Name).To(Equal("Keep This"))
		})
	})

	Describe("accepted burst preserves original ID from InferredBursts", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				InferredBursts: func() []*career.Burst {
					b := fixtures.Burst("inferred-burst-uuid-1", "evt-1", "evt-2")
					b.Name = "API Development"
					b.Description = "Built REST APIs"
					return []*career.Burst{b}
				}(),
				AcceptedBursts: make([]*career.Burst, 0),
				EditingMode:    EditingModeBursts,
			}
			suggestions := []burstfact.BurstSuggestion{
				{Name: "API Development", Description: "Built REST APIs", EventIDs: []string{"evt-1", "evt-2"}, ConfidenceScore: 0.9},
			}
			intent.reviewState.burstModal = modals.NewSuggestionReviewModal(suggestions, nil)

			breadcrumbs := []string{"Test"}
			screen := captureScreens.NewEventReviewScreen(breadcrumbs, intent.reviewState.Event, nil, nil, nil)
			intent.activeScreen = screen
		})

		It("uses the original burst with DB-assigned ID when name matches an InferredBurst", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.AcceptedBursts).To(HaveLen(1))
			Expect(intent.reviewState.AcceptedBursts[0].ID).To(Equal("inferred-burst-uuid-1"))
			Expect(intent.reviewState.AcceptedBursts[0].Name).To(Equal("API Development"))
		})

		It("assigns a fresh UUID when the suggestion name does not match any InferredBurst", func() {
			intent.reviewState.InferredBursts = []*career.Burst{}
			intent.reviewState.burstModal = modals.NewSuggestionReviewModal([]burstfact.BurstSuggestion{
				{Name: "Unknown Burst", Description: "no match", EventIDs: []string{"evt-3"}, ConfidenceScore: 0.5},
			}, nil)

			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.AcceptedBursts).To(HaveLen(1))
			Expect(intent.reviewState.AcceptedBursts[0].ID).NotTo(BeEmpty())
			Expect(intent.reviewState.AcceptedBursts[0].Name).To(Equal("Unknown Burst"))
		})
	})

	Describe("appending to existing accepted bursts", func() {
		BeforeEach(func() {
			existingBurst := fixtures.BurstConfirmed("burst-1")
			existingBurst.Name = "Existing Burst"
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "test event", "", ""),
				AcceptedBursts: []*career.Burst{existingBurst},
				EditingMode:    EditingModeBursts,
			}
			suggestions := []burstfact.BurstSuggestion{
				{Name: "New Burst", Description: "newly confirmed", EventIDs: []string{"evt-1"}, ConfidenceScore: 0.8},
			}
			intent.reviewState.burstModal = modals.NewSuggestionReviewModal(suggestions, nil)

			breadcrumbs := []string{"Test"}
			screen := captureScreens.NewEventReviewScreen(breadcrumbs, intent.reviewState.Event, nil, nil, nil)
			intent.activeScreen = screen
		})

		It("appends new bursts to existing accepted list", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(intent.reviewState.AcceptedBursts).To(HaveLen(2))
			Expect(intent.reviewState.AcceptedBursts[0].Name).To(Equal("Existing Burst"))
			Expect(intent.reviewState.AcceptedBursts[1].Name).To(Equal("New Burst"))
		})
	})

	Describe("View with EditingModeBursts", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "Built services", "", ""),
				AcceptedBursts: make([]*career.Burst, 0),
				EditingMode:    EditingModeBursts,
			}
			suggestions := []burstfact.BurstSuggestion{
				{Name: "API Work", Description: "REST API development", EventIDs: []string{"evt-1"}, ConfidenceScore: 0.9},
			}
			intent.reviewState.burstModal = modals.NewSuggestionReviewModal(suggestions, nil)

			breadcrumbs := []string{"Test"}
			screen := captureScreens.NewEventReviewScreen(breadcrumbs, intent.reviewState.Event, nil, nil, nil)
			intent.activeScreen = screen
		})

		It("renders the burst modal overlay when EditingModeBursts is active", func() {
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("does not use getEditingModalContent for bursts", func() {
			view := intent.View()
			Expect(view).NotTo(BeNil())
		})
	})

	Describe("HandleNavigate suggest_skills", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				InferredSkills: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
				},
				AcceptedSkills: make([]*career.Skill, 0),
			}
		})

		It("sets EditingMode to EditingModeSkills", func() {
			result := &screens.NavigateResult{ResultData: "suggest_skills"}
			intent.HandleNavigate(result)

			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeSkills))
		})

		It("creates a skillModal", func() {
			result := &screens.NavigateResult{ResultData: "suggest_skills"}
			intent.HandleNavigate(result)

			Expect(intent.reviewState.skillModal).NotTo(BeNil())
		})
	})

	Describe("updateEditingModal with EditingModeSkills", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "Built services in Go", "", ""),
				AcceptedSkills: make([]*career.Skill, 0),
				EditingMode:    EditingModeSkills,
			}
			suggestions := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.95},
			}
			intent.reviewState.skillModal = modals.NewSkillSuggestionModal(suggestions, nil)

			breadcrumbs := []string{"Test"}
			screen := captureScreens.NewEventReviewScreen(breadcrumbs, intent.reviewState.Event, nil, nil, nil)
			intent.activeScreen = screen
		})

		It("auto-closes modal and resets mode when no suggestions remain", func() {
			intent.reviewState.skillModal = modals.NewSkillSuggestionModal([]skillinference.SkillSuggestion{}, nil)

			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyDown})

			Expect(intent.reviewState.skillModal).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
		})

		It("closes modal on escape without adding skills", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
			Expect(intent.reviewState.AcceptedSkills).To(BeEmpty())
		})
	})
})
