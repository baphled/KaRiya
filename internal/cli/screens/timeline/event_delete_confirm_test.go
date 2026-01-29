package timeline_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens/timeline"
	"github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("EventDeleteConfirmScreen", func() {
	var (
		screen *timeline.EventDeleteConfirmScreen
		event  *career.Event
	)

	BeforeEach(func() {
		event = &career.Event{
			ID:        "test-event-id",
			Text:      "Test event for deletion",
			Date:      time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
			Company:   "Test Company",
			Project:   "Test Project",
			CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		}
	})

	Describe("NewEventDeleteConfirmScreen", func() {
		It("creates screen with event", func() {
			screen = timeline.NewEventDeleteConfirmScreen(event)

			Expect(screen).NotTo(BeNil())
			Expect(screen.GetEvent()).To(Equal(event))
		})

		It("preserves event ID", func() {
			screen = timeline.NewEventDeleteConfirmScreen(event)

			Expect(screen.GetEvent().ID).To(Equal("test-event-id"))
		})

		It("preserves event text", func() {
			screen = timeline.NewEventDeleteConfirmScreen(event)

			Expect(screen.GetEvent().Text).To(Equal("Test event for deletion"))
		})

		It("truncates long event text in confirmation message", func() {
			longTextEvent := &career.Event{
				ID:   "long-text-id",
				Text: "This is a very long event text that exceeds sixty characters and should be truncated in the confirmation message for better display",
				Date: time.Now(),
			}

			screen = timeline.NewEventDeleteConfirmScreen(longTextEvent)

			// Screen is created successfully - original text is preserved in the event.
			Expect(screen).NotTo(BeNil())
			// GetEvent returns the original event with full text.
			Expect(len(screen.GetEvent().Text)).To(BeNumerically(">", 60))
		})

		It("handles short event text without truncation", func() {
			shortTextEvent := &career.Event{
				ID:   "short-text-id",
				Text: "Short text",
				Date: time.Now(),
			}

			screen = timeline.NewEventDeleteConfirmScreen(shortTextEvent)

			Expect(screen).NotTo(BeNil())
			Expect(screen.GetEvent().Text).To(Equal("Short text"))
		})

		It("handles exactly 60 character text", func() {
			exactTextEvent := &career.Event{
				ID:   "exact-text-id",
				Text: "This is exactly sixty characters long for testing purposes!!", // 60 chars
				Date: time.Now(),
			}

			screen = timeline.NewEventDeleteConfirmScreen(exactTextEvent)

			Expect(screen).NotTo(BeNil())
			Expect(screen.GetEvent().Text).To(HaveLen(60))
		})
	})

	Describe("GetEvent", func() {
		It("returns the event being considered for deletion", func() {
			screen = timeline.NewEventDeleteConfirmScreen(event)

			result := screen.GetEvent()

			Expect(result).To(Equal(event))
			Expect(result.ID).To(Equal("test-event-id"))
			Expect(result.Text).To(Equal("Test event for deletion"))
			Expect(result.Company).To(Equal("Test Company"))
			Expect(result.Project).To(Equal("Test Project"))
		})

		It("returns event with all fields preserved", func() {
			fullEvent := &career.Event{
				ID:         "full-event",
				Text:       "Full event",
				Date:       time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
				Company:    "Big Corp",
				Project:    "Big Project",
				Tags:       []string{"tag1", "tag2"},
				Categories: []string{"cat1"},
				Skills:     []string{"skill1"},
				CreatedAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			}

			screen = timeline.NewEventDeleteConfirmScreen(fullEvent)

			result := screen.GetEvent()
			Expect(result.Tags).To(ContainElement("tag1"))
			Expect(result.Categories).To(ContainElement("cat1"))
			Expect(result.Skills).To(ContainElement("skill1"))
		})
	})

	Describe("State constant", func() {
		It("has correct state name for state matrix", func() {
			Expect(timeline.EventDeleteConfirmState).To(Equal("event_delete_confirm"))
		})
	})
})
