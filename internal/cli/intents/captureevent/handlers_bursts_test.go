package captureevent

import (
	"context"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	captureScreens "github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/domain/career"
	memoryrepo "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Burst Acceptance in CaptureEvent", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentContext{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("when user accepts bursts from modal", func() {
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
			intent.reviewState.burstModal = NewBurstSuggestionModelNew(context.Background(), nil, suggestions)

			breadcrumbs := []string{"Test"}
			screen := captureScreens.NewEventReviewScreen(breadcrumbs, intent.reviewState.Event, nil, nil, nil)
			intent.activeScreen = screen
		})

		It("transfers accepted bursts to review state", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			Expect(intent.reviewState.AcceptedBursts).To(HaveLen(1))
			Expect(intent.reviewState.AcceptedBursts[0].Name).To(Equal("API Development"))
			Expect(intent.reviewState.AcceptedBursts[0].Description).To(Equal("Built REST APIs"))
			Expect(intent.reviewState.AcceptedBursts[0].EventIDs).To(Equal([]string{"evt-1"}))
		})

		It("closes the modal and resets editing mode", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			Expect(intent.reviewState.burstModal).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
		})

		It("returns a command from the modal update", func() {
			cmd := intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("when user rejects all bursts from modal", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("evt-1", "test event", "", ""),
				AcceptedBursts: make([]*career.Burst, 0),
				EditingMode:    EditingModeBursts,
			}
			suggestions := []burstfact.BurstSuggestion{
				{Name: "Rejected Burst", Description: "desc", EventIDs: []string{"evt-1"}, ConfidenceScore: 0.5},
			}
			intent.reviewState.burstModal = NewBurstSuggestionModelNew(context.Background(), nil, suggestions)

			breadcrumbs := []string{"Test"}
			screen := captureScreens.NewEventReviewScreen(breadcrumbs, intent.reviewState.Event, nil, nil, nil)
			intent.activeScreen = screen
		})

		It("does not add rejected bursts to accepted list", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

			Expect(intent.reviewState.AcceptedBursts).To(BeEmpty())
		})

		It("closes the modal and resets editing mode", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

			Expect(intent.reviewState.burstModal).To(BeNil())
			Expect(intent.reviewState.EditingMode).To(Equal(EditingModeNone))
		})
	})

	Describe("when user accepts some and rejects some bursts", func() {
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
			intent.reviewState.burstModal = NewBurstSuggestionModelNew(context.Background(), nil, suggestions)

			breadcrumbs := []string{"Test"}
			screen := captureScreens.NewEventReviewScreen(breadcrumbs, intent.reviewState.Event, nil, nil, nil)
			intent.activeScreen = screen
		})

		It("only includes confirmed bursts in accepted list", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

			Expect(intent.reviewState.AcceptedBursts).To(HaveLen(1))
			Expect(intent.reviewState.AcceptedBursts[0].Name).To(Equal("Keep This"))
		})
	})

	Describe("when burst suggestion has empty name", func() {
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
			intent.reviewState.burstModal = NewBurstSuggestionModelNew(context.Background(), nil, suggestions)

			breadcrumbs := []string{"Test"}
			screen := captureScreens.NewEventReviewScreen(breadcrumbs, intent.reviewState.Event, nil, nil, nil)
			intent.activeScreen = screen
		})

		It("generates a name for the burst", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			Expect(intent.reviewState.AcceptedBursts).To(HaveLen(1))
			Expect(intent.reviewState.AcceptedBursts[0].Name).To(Equal("Burst of 2 events"))
		})
	})

	Describe("when bursts are appended to existing accepted bursts", func() {
		BeforeEach(func() {
			intent.currentState = StateReview
			existingBurst := fixtures.BurstConfirmed("burst-1")
			existingBurst.Name = "Existing Burst"
			existingBurst.Description = "already accepted"
			intent.reviewState = &ReviewInferredEventState{
				Event: fixtures.EventWith("evt-1", "test event", "", ""),
				AcceptedBursts: []*career.Burst{
					existingBurst,
				},
				EditingMode: EditingModeBursts,
			}
			suggestions := []burstfact.BurstSuggestion{
				{Name: "New Burst", Description: "newly confirmed", EventIDs: []string{"evt-1"}, ConfidenceScore: 0.8},
			}
			intent.reviewState.burstModal = NewBurstSuggestionModelNew(context.Background(), nil, suggestions)

			breadcrumbs := []string{"Test"}
			screen := captureScreens.NewEventReviewScreen(breadcrumbs, intent.reviewState.Event, nil, nil, nil)
			intent.activeScreen = screen
		})

		It("appends new bursts to existing accepted list", func() {
			intent.updateEditingModal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			Expect(intent.reviewState.AcceptedBursts).To(HaveLen(2))
			Expect(intent.reviewState.AcceptedBursts[0].Name).To(Equal("Existing Burst"))
			Expect(intent.reviewState.AcceptedBursts[1].Name).To(Equal("New Burst"))
		})
	})

	Describe("burst persistence in postSaveReview", func() {
		var (
			svc       *careerservice.Service
			burstRepo *memoryrepo.BurstRepository
		)

		BeforeEach(func() {
			repos := memoryrepo.NewRepositories()
			burstRepo = repos.Burst.(*memoryrepo.BurstRepository)

			svc = careerservice.NewService(repos.Event)
			svc.SetBurstRepository(burstRepo)
			svc.SetSkillRepository(repos.Skill)
			svc.SetFactRepository(repos.Fact)

			intent.context.CareerService = svc
		})

		Context("with accepted bursts that exist in the repository", func() {
			It("calls ConfirmBurst to mark them as confirmed", func() {
				event := fixtures.EventWith("evt-saved-1", "Reviewed event with bursts", "", "")
				burst := fixtures.Burst("burst-to-confirm", "evt-saved-1", "evt-other-1")

				ctx := context.Background()
				Expect(burstRepo.Create(ctx, burst)).To(Succeed())
				Expect(burst.Confirmed).To(BeFalse())

				intent.currentState = StateReview
				intent.postSaveReview = true
				intent.reviewState = &ReviewInferredEventState{
					Event: event,
				}

				reviewData := map[string]interface{}{
					"event":  event,
					"bursts": []*career.Burst{burst},
					"facts":  []*career.Fact{},
				}
				cmd := intent.HandleSubmit(&screens.SubmitResult{FormData: reviewData})
				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				intent.Update(msg)

				Expect(intent.result).NotTo(BeNil())
				Expect(intent.result.Status).To(Equal(intents.Completed))

				confirmed, err := burstRepo.GetByID(ctx, "burst-to-confirm")
				Expect(err).NotTo(HaveOccurred())
				Expect(confirmed.Confirmed).To(BeTrue())
				Expect(confirmed.ConfirmedAt).NotTo(BeNil())
			})

			It("preserves EventIDs through ConfirmBurst", func() {
				event := fixtures.EventWith("evt-saved-2", "Reviewed event", "", "")
				burst := fixtures.Burst("burst-preserve-ids", "evt-saved-2", "evt-other-2")

				ctx := context.Background()
				Expect(burstRepo.Create(ctx, burst)).To(Succeed())

				intent.currentState = StateReview
				intent.postSaveReview = true
				intent.reviewState = &ReviewInferredEventState{
					Event: event,
				}

				reviewData := map[string]interface{}{
					"event":  event,
					"bursts": []*career.Burst{burst},
					"facts":  []*career.Fact{},
				}
				cmd := intent.HandleSubmit(&screens.SubmitResult{FormData: reviewData})
				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				intent.Update(msg)

				Expect(intent.result).NotTo(BeNil())
				Expect(intent.result.Status).To(Equal(intents.Completed))

				confirmed, err := burstRepo.GetByID(ctx, "burst-preserve-ids")
				Expect(err).NotTo(HaveOccurred())
				Expect(confirmed.EventIDs).To(Equal([]string{"evt-saved-2", "evt-other-2"}))
			})
		})

		Context("when CareerService is nil", func() {
			It("completes without confirming bursts", func() {
				event := fixtures.EventWith("evt-saved-3", "Reviewed event", "", "")
				burst := fixtures.Burst("burst-no-svc", "evt-saved-3", "evt-other-3")

				intent.context.CareerService = nil
				intent.currentState = StateReview
				intent.postSaveReview = true
				intent.reviewState = &ReviewInferredEventState{
					Event: event,
				}

				reviewData := map[string]interface{}{
					"event":  event,
					"bursts": []*career.Burst{burst},
					"facts":  []*career.Fact{},
				}
				cmd := intent.HandleSubmit(&screens.SubmitResult{FormData: reviewData})
				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				intent.Update(msg)

				Expect(intent.result).NotTo(BeNil())
				Expect(intent.result.Status).To(Equal(intents.Completed))
				Expect(burst.Confirmed).To(BeFalse())
			})
		})
	})
})
