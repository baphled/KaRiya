package display_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/ui/display"
)

var _ = Describe("Event display type", func() {
	Describe("EventFromDomain", func() {
		It("converts a fully populated event", func() {
			createdAt := time.Date(2024, time.January, 10, 8, 30, 0, 0, time.UTC)
			updatedAt := time.Date(2024, time.January, 12, 9, 45, 0, 0, time.UTC)
			date := time.Date(2023, time.December, 1, 0, 0, 0, 0, time.UTC)

			event := fixtures.EventWith("event-1", "Led migration to Go services", "Acme", "Platform")
			event.Date = date
			event.Tags = []string{"technical", "leadership"}
			event.Categories = []string{"technical", "architecture"}
			event.Skills = []string{"Go", "Kubernetes"}
			event.CreatedAt = createdAt
			event.UpdatedAt = updatedAt

			result := display.EventFromDomain(event)

			Expect(result).To(Equal(display.Event{
				ID:         "event-1",
				Text:       "Led migration to Go services",
				Date:       date,
				Company:    "Acme",
				Project:    "Platform",
				Tags:       []string{"technical", "leadership"},
				Categories: []string{"technical", "architecture"},
				Skills:     []string{"Go", "Kubernetes"},
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
			}))
		})

		It("handles nil input gracefully", func() {
			result := display.EventFromDomain(nil)

			Expect(result).To(Equal(display.Event{}))
		})
	})

	Describe("EventsFromDomain", func() {
		It("converts a slice of events", func() {
			e1 := fixtures.Event("event-1")
			e1.Text = "First event"
			e2 := fixtures.Event("event-2")
			e2.Text = "Second event"
			events := []*career.Event{e1, e2}

			result := display.EventsFromDomain(events)

			Expect(result).To(HaveLen(2))
			Expect(result[0].ID).To(Equal("event-1"))
			Expect(result[0].Text).To(Equal("First event"))
			Expect(result[1].ID).To(Equal("event-2"))
			Expect(result[1].Text).To(Equal("Second event"))
		})

		It("returns nil for nil input", func() {
			Expect(display.EventsFromDomain(nil)).To(BeNil())
		})

		It("returns empty slice for empty input", func() {
			result := display.EventsFromDomain([]*career.Event{})

			Expect(result).To(Equal([]display.Event{}))
		})
	})
})
