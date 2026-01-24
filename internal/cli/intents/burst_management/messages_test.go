package burst_management_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burst_fact"
)

var (
	errTest = errors.New("test error")
)

var _ = Describe("Messages", func() {
	Describe("BurstSelectedMsg", func() {
		It("should store the selected burst", func() {
			burst := &career.Burst{
				ID:   "burst-1",
				Name: "Q1 2024 Backend Work",
			}
			msg := burst_management.BurstSelectedMsg{
				Burst: burst,
				Index: 0,
			}
			Expect(msg.Burst).To(Equal(burst))
			Expect(msg.Burst.ID).To(Equal("burst-1"))
		})

		It("should store the selection index", func() {
			msg := burst_management.BurstSelectedMsg{
				Burst: &career.Burst{ID: "burst-1"},
				Index: 3,
			}
			Expect(msg.Index).To(Equal(3))
		})

		It("should handle nil burst", func() {
			msg := burst_management.BurstSelectedMsg{
				Burst: nil,
				Index: 0,
			}
			Expect(msg.Burst).To(BeNil())
		})
	})

	Describe("BurstEventsLoadedMsg", func() {
		It("should store loaded events", func() {
			events := []*career.CareerEvent{
				{ID: "event-1", Text: "Event 1"},
				{ID: "event-2", Text: "Event 2"},
			}
			msg := burst_management.BurstEventsLoadedMsg{
				Events: events,
				Error:  nil,
			}
			Expect(msg.Events).To(HaveLen(2))
			Expect(msg.Error).To(BeNil())
		})

		It("should store error when loading fails", func() {
			msg := burst_management.BurstEventsLoadedMsg{
				Events: nil,
				Error:  errTest,
			}
			Expect(msg.Events).To(BeNil())
			Expect(msg.Error).To(Equal(errTest))
		})
	})

	Describe("BurstFactsLoadedMsg", func() {
		It("should store loaded facts", func() {
			facts := []*career.Fact{
				{ID: "fact-1", Text: "Fact 1"},
				{ID: "fact-2", Text: "Fact 2"},
			}
			msg := burst_management.BurstFactsLoadedMsg{
				Facts: facts,
				Error: nil,
			}
			Expect(msg.Facts).To(HaveLen(2))
			Expect(msg.Error).To(BeNil())
		})

		It("should store error when loading fails", func() {
			msg := burst_management.BurstFactsLoadedMsg{
				Facts: nil,
				Error: errTest,
			}
			Expect(msg.Facts).To(BeNil())
			Expect(msg.Error).NotTo(BeNil())
		})
	})

	Describe("BurstEditCompleteMsg", func() {
		It("should store edited burst on success", func() {
			burst := &career.Burst{
				ID:   "burst-1",
				Name: "Updated Burst",
			}
			msg := burst_management.BurstEditCompleteMsg{
				Burst:     burst,
				Cancelled: false,
				Error:     nil,
			}
			Expect(msg.Burst).To(Equal(burst))
			Expect(msg.Cancelled).To(BeFalse())
			Expect(msg.Error).To(BeNil())
		})

		It("should indicate cancellation", func() {
			msg := burst_management.BurstEditCompleteMsg{
				Burst:     nil,
				Cancelled: true,
				Error:     nil,
			}
			Expect(msg.Cancelled).To(BeTrue())
		})

		It("should store error on failure", func() {
			msg := burst_management.BurstEditCompleteMsg{
				Burst:     nil,
				Cancelled: false,
				Error:     errTest,
			}
			Expect(msg.Error).To(Equal(errTest))
		})
	})

	Describe("BurstDeletedMsg", func() {
		It("should store deleted burst ID", func() {
			msg := burst_management.BurstDeletedMsg{
				BurstID: "burst-123",
				Error:   nil,
			}
			Expect(msg.BurstID).To(Equal("burst-123"))
			Expect(msg.Error).To(BeNil())
		})

		It("should store error on delete failure", func() {
			msg := burst_management.BurstDeletedMsg{
				BurstID: "burst-123",
				Error:   errTest,
			}
			Expect(msg.Error).NotTo(BeNil())
		})
	})

	Describe("BurstConfirmedMsg", func() {
		It("should store confirmed burst", func() {
			now := time.Now()
			burst := &career.Burst{
				ID:          "burst-1",
				Name:        "Confirmed Burst",
				Confirmed:   true,
				ConfirmedAt: &now,
			}
			msg := burst_management.BurstConfirmedMsg{
				Burst: burst,
				Error: nil,
			}
			Expect(msg.Burst.Confirmed).To(BeTrue())
			Expect(msg.Burst.ConfirmedAt).NotTo(BeNil())
		})

		It("should store error on confirmation failure", func() {
			msg := burst_management.BurstConfirmedMsg{
				Burst: nil,
				Error: errTest,
			}
			Expect(msg.Error).NotTo(BeNil())
		})
	})

	Describe("FactExtractionCompleteMsg", func() {
		It("should store extracted facts", func() {
			facts := []*career.Fact{
				{ID: "fact-1", Text: "Extracted fact 1"},
				{ID: "fact-2", Text: "Extracted fact 2"},
			}
			msg := burst_management.FactExtractionCompleteMsg{
				Facts: facts,
				Error: nil,
			}
			Expect(msg.Facts).To(HaveLen(2))
			Expect(msg.Error).To(BeNil())
		})

		It("should store error on extraction failure", func() {
			msg := burst_management.FactExtractionCompleteMsg{
				Facts: nil,
				Error: errTest,
			}
			Expect(msg.Error).NotTo(BeNil())
		})
	})

	Describe("BurstSuggestionsLoadedMsg", func() {
		It("should store burst suggestions", func() {
			suggestions := []burst_fact.BurstSuggestion{
				{
					EventIDs:        []string{"event-1", "event-2"},
					ConfidenceScore: 0.85,
					Name:            "Q1 Backend Work",
				},
			}
			msg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: suggestions,
				Error:       nil,
			}
			Expect(msg.Suggestions).To(HaveLen(1))
			Expect(msg.Suggestions[0].ConfidenceScore).To(Equal(0.85))
			Expect(msg.Error).To(BeNil())
		})

		It("should store error when detection fails", func() {
			msg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: nil,
				Error:       errTest,
			}
			Expect(msg.Suggestions).To(BeNil())
			Expect(msg.Error).NotTo(BeNil())
		})
	})

	Describe("BurstSuggestionAcceptedMsg", func() {
		It("should store accepted suggestion and created burst", func() {
			suggestion := burst_fact.BurstSuggestion{
				EventIDs:        []string{"event-1", "event-2"},
				ConfidenceScore: 0.9,
				Name:            "Accepted Burst",
			}
			burst := &career.Burst{
				ID:   "burst-new",
				Name: "Accepted Burst",
			}
			msg := burst_management.BurstSuggestionAcceptedMsg{
				Suggestion: suggestion,
				Burst:      burst,
				Error:      nil,
			}
			Expect(msg.Suggestion.Name).To(Equal("Accepted Burst"))
			Expect(msg.Burst.ID).To(Equal("burst-new"))
			Expect(msg.Error).To(BeNil())
		})

		It("should store error on acceptance failure", func() {
			suggestion := burst_fact.BurstSuggestion{
				EventIDs: []string{"event-1"},
			}
			msg := burst_management.BurstSuggestionAcceptedMsg{
				Suggestion: suggestion,
				Burst:      nil,
				Error:      errTest,
			}
			Expect(msg.Error).NotTo(BeNil())
		})
	})

	Describe("Message Type Distinctiveness", func() {
		It("should have distinct message types", func() {
			// Verify each message type is distinguishable via type assertion.
			var msg1 interface{} = burst_management.BurstSelectedMsg{}
			var msg2 interface{} = burst_management.BurstEventsLoadedMsg{}
			var msg3 interface{} = burst_management.BurstDeletedMsg{}

			_, isSelected := msg1.(burst_management.BurstSelectedMsg)
			_, isEventsLoaded := msg2.(burst_management.BurstEventsLoadedMsg)
			_, isDeleted := msg3.(burst_management.BurstDeletedMsg)

			Expect(isSelected).To(BeTrue())
			Expect(isEventsLoaded).To(BeTrue())
			Expect(isDeleted).To(BeTrue())
		})

		It("should not confuse message types", func() {
			var msg interface{} = burst_management.BurstSelectedMsg{}

			_, isEventsLoaded := msg.(burst_management.BurstEventsLoadedMsg)
			_, isDeleted := msg.(burst_management.BurstDeletedMsg)

			Expect(isEventsLoaded).To(BeFalse())
			Expect(isDeleted).To(BeFalse())
		})
	})
})
