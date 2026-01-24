package browse_timeline_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/browse_timeline"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Messages", func() {
	Describe("EventSelectedMsg", func() {
		It("should store the selected event", func() {
			event := fixtures.EventWith("event-1", "Backend Developer at TechCorp", "TechCorp", "Platform")
			msg := browse_timeline.EventSelectedMsg{
				Event: event,
				Index: 0,
			}
			Expect(msg.Event).To(Equal(event))
			Expect(msg.Event.ID).To(Equal("event-1"))
		})

		It("should store the selection index", func() {
			msg := browse_timeline.EventSelectedMsg{
				Event: fixtures.Event("event-1"),
				Index: 5,
			}
			Expect(msg.Index).To(Equal(5))
		})

		Context("when event is nil", func() {
			It("should handle nil event", func() {
				msg := browse_timeline.EventSelectedMsg{
					Event: nil,
					Index: 0,
				}
				Expect(msg.Event).To(BeNil())
			})
		})

		Context("when index is zero", func() {
			It("should handle zero index", func() {
				msg := browse_timeline.EventSelectedMsg{
					Event: fixtures.Event("first"),
					Index: 0,
				}
				Expect(msg.Index).To(Equal(0))
			})
		})
	})

	Describe("FilterChangedMsg", func() {
		It("should store the new filters", func() {
			filters := &browse_timeline.Filters{
				SearchText: "developer",
				SortBy:     "date",
				SortOrder:  "desc",
			}
			msg := browse_timeline.FilterChangedMsg{
				Filters: filters,
			}
			Expect(msg.Filters).To(Equal(filters))
			Expect(msg.Filters.SearchText).To(Equal("developer"))
		})

		Context("when filters are empty", func() {
			It("should handle empty filters", func() {
				filters := &browse_timeline.Filters{}
				msg := browse_timeline.FilterChangedMsg{
					Filters: filters,
				}
				Expect(msg.Filters).NotTo(BeNil())
				Expect(msg.Filters.SearchText).To(BeEmpty())
			})
		})

		Context("when filters are nil", func() {
			It("should handle nil filters", func() {
				msg := browse_timeline.FilterChangedMsg{
					Filters: nil,
				}
				Expect(msg.Filters).To(BeNil())
			})
		})

		It("should store complex filter combinations", func() {
			filters := &browse_timeline.Filters{
				SearchText: "golang",
				Tags:       []string{"backend", "api"},
				Companies:  []string{"TechCorp"},
				Categories: []string{"development"},
				Projects:   []string{"API Platform"},
				DateFrom:   "2023-01-01",
				DateTo:     "2024-12-31",
				SortBy:     "text",
				SortOrder:  "asc",
			}
			msg := browse_timeline.FilterChangedMsg{
				Filters: filters,
			}
			Expect(msg.Filters.Tags).To(HaveLen(2))
			Expect(msg.Filters.Companies).To(HaveLen(1))
			Expect(msg.Filters.DateFrom).To(Equal("2023-01-01"))
		})
	})

	Describe("EventDeletedMsg", func() {
		It("should store the deleted event ID", func() {
			msg := browse_timeline.EventDeletedMsg{
				EventID: "event-123",
			}
			Expect(msg.EventID).To(Equal("event-123"))
		})

		Context("when event ID is empty", func() {
			It("should handle empty event ID", func() {
				msg := browse_timeline.EventDeletedMsg{
					EventID: "",
				}
				Expect(msg.EventID).To(BeEmpty())
			})
		})

		It("should handle UUID-style event ID", func() {
			msg := browse_timeline.EventDeletedMsg{
				EventID: "550e8400-e29b-41d4-a716-446655440000",
			}
			Expect(msg.EventID).To(Equal("550e8400-e29b-41d4-a716-446655440000"))
		})
	})

	Describe("Message Type Distinctiveness", func() {
		It("should have distinct message types", func() {
			// Verify each message type is distinguishable via type assertion.
			var msg1 interface{} = browse_timeline.EventSelectedMsg{}
			var msg2 interface{} = browse_timeline.FilterChangedMsg{}
			var msg3 interface{} = browse_timeline.EventDeletedMsg{}

			_, isEventSelected := msg1.(browse_timeline.EventSelectedMsg)
			_, isFilterChanged := msg2.(browse_timeline.FilterChangedMsg)
			_, isEventDeleted := msg3.(browse_timeline.EventDeletedMsg)

			Expect(isEventSelected).To(BeTrue())
			Expect(isFilterChanged).To(BeTrue())
			Expect(isEventDeleted).To(BeTrue())
		})

		It("should not confuse message types", func() {
			var msg interface{} = browse_timeline.EventSelectedMsg{}

			_, isFilterChanged := msg.(browse_timeline.FilterChangedMsg)
			_, isEventDeleted := msg.(browse_timeline.EventDeletedMsg)

			Expect(isFilterChanged).To(BeFalse())
			Expect(isEventDeleted).To(BeFalse())
		})
	})
})
