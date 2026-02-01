package burst_management_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var (
	errTest = errors.New("test error")
)

var _ = Describe("Messages", func() {
	Describe("BurstSelectedMsg", func() {
		It("should store the selected burst", func() {
			burst := fixtures.Burst("burst-1")
			msg := burst_management.BurstSelectedMsg{
				Burst: burst,
				Index: 0,
			}
			Expect(msg.Burst).To(Equal(burst))
			Expect(msg.Burst.ID).To(Equal("burst-1"))
		})

		It("should store the selection index", func() {
			msg := burst_management.BurstSelectedMsg{
				Burst: fixtures.Burst("burst-1"),
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
			events := []*career.Event{
				fixtures.EventWith("event-1", "Event 1", "", ""),
				fixtures.EventWith("event-2", "Event 2", "", ""),
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
				fixtures.FactWith("fact-1", "Fact 1"),
				fixtures.FactWith("fact-2", "Fact 2"),
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
			burst := fixtures.Burst("burst-1")
			burst.Name = "Updated Burst"
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
			burst := fixtures.BurstConfirmed("burst-1")
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
				fixtures.FactWith("fact-1", "Extracted fact 1"),
				fixtures.FactWith("fact-2", "Extracted fact 2"),
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

		It("should store associated burst", func() {
			burst := fixtures.BurstConfirmed("burst-1")
			burst.Name = "Test Burst"
			msg := burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{fixtures.FactWith("f1", "Fact")},
				Burst: burst,
			}
			Expect(msg.Burst).To(Equal(burst))
			Expect(msg.Burst.ID).To(Equal("burst-1"))
		})

		It("should handle nil burst", func() {
			msg := burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{fixtures.FactWith("f1", "Fact")},
				Burst: nil,
			}
			Expect(msg.Burst).To(BeNil())
		})
	})

	Describe("BurstSuggestionsLoadedMsg", func() {
		It("should store burst suggestions", func() {
			suggestions := []burstfact.BurstSuggestion{
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

	Describe("SuggestionReviewCompleteMsg", func() {
		It("should store accepted suggestions", func() {
			suggestions := []burstfact.BurstSuggestion{
				{
					EventIDs:        []string{"event-1", "event-2"},
					ConfidenceScore: 0.9,
					Name:            "Accepted Burst 1",
				},
				{
					EventIDs:        []string{"event-3", "event-4"},
					ConfidenceScore: 0.85,
					Name:            "Accepted Burst 2",
				},
			}
			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
				Cancelled:           false,
			}
			Expect(msg.AcceptedSuggestions).To(HaveLen(2))
			Expect(msg.AcceptedSuggestions[0].Name).To(Equal("Accepted Burst 1"))
			Expect(msg.Cancelled).To(BeFalse())
		})

		It("should indicate cancellation with empty suggestions", func() {
			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: nil,
				Cancelled:           true,
			}
			Expect(msg.Cancelled).To(BeTrue())
			Expect(msg.AcceptedSuggestions).To(BeNil())
		})

		It("should handle completion with no accepted suggestions", func() {
			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{},
				Cancelled:           false,
			}
			Expect(msg.AcceptedSuggestions).To(HaveLen(0))
			Expect(msg.Cancelled).To(BeFalse())
		})
	})

	Describe("EditBurstMsg", func() {
		It("should store edit data", func() {
			msg := burst_management.EditBurstMsg{
				BurstID:     "burst-1",
				Name:        "Updated Name",
				Description: "Updated description",
			}
			Expect(msg.BurstID).To(Equal("burst-1"))
			Expect(msg.Name).To(Equal("Updated Name"))
			Expect(msg.Description).To(Equal("Updated description"))
		})

		It("should handle empty description", func() {
			msg := burst_management.EditBurstMsg{
				BurstID: "burst-1",
				Name:    "Name Only",
			}
			Expect(msg.Description).To(BeEmpty())
		})
	})

	Describe("BurstSkillsLoadedMsg", func() {
		It("should store loaded skills", func() {
			skills := []*career.Skill{
				fixtures.SkillWith("s1", "Go", "Backend", "mid"),
				fixtures.SkillWith("s2", "Docker", "DevOps", "mid"),
			}
			msg := burst_management.BurstSkillsLoadedMsg{
				Skills: skills,
				Error:  nil,
			}
			Expect(msg.Skills).To(HaveLen(2))
			Expect(msg.Skills[0].Name).To(Equal("Go"))
			Expect(msg.Error).To(BeNil())
		})

		It("should store error when loading fails", func() {
			msg := burst_management.BurstSkillsLoadedMsg{
				Skills: nil,
				Error:  errTest,
			}
			Expect(msg.Skills).To(BeNil())
			Expect(msg.Error).To(Equal(errTest))
		})
	})

	Describe("SkillSuggestionsLoadedMsg", func() {
		It("should store skill suggestions", func() {
			suggestions := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.95},
				{Name: "Docker", Category: "devops", Confidence: 0.85},
			}
			msg := burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: suggestions,
				Error:       nil,
			}
			Expect(msg.Suggestions).To(HaveLen(2))
			Expect(msg.Suggestions[0].Name).To(Equal("Go"))
			Expect(msg.Suggestions[0].Confidence).To(Equal(0.95))
			Expect(msg.Error).To(BeNil())
		})

		It("should store existing skill names", func() {
			msg := burst_management.SkillSuggestionsLoadedMsg{
				Suggestions:        []skillinference.SkillSuggestion{{Name: "Go", Category: "backend", Confidence: 0.9}},
				ExistingSkillNames: []string{"Docker", "Kubernetes"},
			}
			Expect(msg.ExistingSkillNames).To(HaveLen(2))
			Expect(msg.ExistingSkillNames).To(ContainElements("Docker", "Kubernetes"))
		})

		It("should store error on inference failure", func() {
			msg := burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: nil,
				Error:       errTest,
			}
			Expect(msg.Suggestions).To(BeNil())
			Expect(msg.Error).To(Equal(errTest))
		})

		It("should handle empty suggestions", func() {
			msg := burst_management.SkillSuggestionsLoadedMsg{
				Suggestions:        []skillinference.SkillSuggestion{},
				ExistingSkillNames: []string{},
			}
			Expect(msg.Suggestions).To(BeEmpty())
			Expect(msg.ExistingSkillNames).To(BeEmpty())
		})
	})

	Describe("SkillSuggestionsErrorMsg", func() {
		It("should store the error", func() {
			msg := burst_management.SkillSuggestionsErrorMsg{
				Err: errTest,
			}
			Expect(msg.Err).To(Equal(errTest))
		})

		It("should handle nil error", func() {
			msg := burst_management.SkillSuggestionsErrorMsg{
				Err: nil,
			}
			Expect(msg.Err).To(BeNil())
		})
	})

	Describe("SkillsCreatedMsg", func() {
		It("should store created skills", func() {
			skills := []*career.Skill{
				{ID: "s1", Name: "Go"},
				{ID: "s2", Name: "Docker"},
			}
			msg := burst_management.SkillsCreatedMsg{
				Skills: skills,
				Error:  nil,
			}
			Expect(msg.Skills).To(HaveLen(2))
			Expect(msg.Skills[0].Name).To(Equal("Go"))
			Expect(msg.Error).To(BeNil())
		})

		It("should store error on creation failure", func() {
			msg := burst_management.SkillsCreatedMsg{
				Skills: nil,
				Error:  errTest,
			}
			Expect(msg.Skills).To(BeNil())
			Expect(msg.Error).To(Equal(errTest))
		})

		It("should handle empty skills list on success", func() {
			msg := burst_management.SkillsCreatedMsg{
				Skills: []*career.Skill{},
				Error:  nil,
			}
			Expect(msg.Skills).To(BeEmpty())
			Expect(msg.Error).To(BeNil())
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
