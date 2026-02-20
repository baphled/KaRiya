package captureevent

import (
	"context"

	captureScreens "github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/domain/career"
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
})
