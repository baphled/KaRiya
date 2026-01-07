package intents_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("BrowseTimeline - Escape Key Behavior", func() {
	var (
		browseContext *intents.BrowseTimelineContext
		browseIntent  *intents.BrowseTimelineIntent
		sampleEvents  []*career.CareerEvent
	)

	BeforeEach(func() {
		// Create sample events for testing
		now := time.Now()
		sampleEvents = []*career.CareerEvent{
			{
				ID:         "evt-1",
				Text:       "Completed API integration",
				Date:       now.AddDate(0, 0, -1),
				Company:    "Tech Corp",
				CreatedAt:  now.AddDate(0, 0, -1),
				Categories: []string{"technical"},
				Tags:       []string{"backend"},
			},
			{
				ID:         "evt-2",
				Text:       "Led team standup",
				Date:       now.AddDate(0, 0, -2),
				Company:    "Tech Corp",
				CreatedAt:  now.AddDate(0, 0, -2),
				Categories: []string{"leadership"},
				Tags:       []string{"management"},
			},
		}

		browseContext = &intents.BrowseTimelineContext{
			Events: sampleEvents,
			InitialFilters: &intents.TimelineFilters{
				SortBy:    "date",
				SortOrder: "desc",
			},
		}

		var err error
		browseIntent, err = intents.NewBrowseTimelineIntent(browseContext)
		Expect(err).ToNot(HaveOccurred())
		browseIntent.Init()
	})

	Describe("Timeline View (Root State)", func() {
		It("should cancel intent when escape is pressed", func() {
			browseIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := browseIntent.Result()
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should cancel intent when 'm' is pressed", func() {
			browseIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			result := browseIntent.Result()
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should cancel intent when 'q' is pressed", func() {
			browseIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

			result := browseIntent.Result()
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})
	})

	Describe("Event Detail View", func() {
		BeforeEach(func() {
			// Navigate to event detail
			browseIntent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})

		It("should go back to timeline when escape is pressed", func() {
			// Verify we're in event detail by checking the view
			view := browseIntent.View()
			Expect(view).To(ContainSubstring("Event Details"))

			// Press escape
			browseIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Verify we're back in timeline view
			view = browseIntent.View()
			Expect(view).ToNot(ContainSubstring("Event Details"))
			Expect(view).To(Or(
				ContainSubstring("Browse Timeline"),
				ContainSubstring("Completed API integration"),
			))

			// Intent should not be cancelled yet
			result := browseIntent.Result()
			Expect(result).To(BeNil())
		})

		It("should cancel intent when 'm' is pressed", func() {
			// Verify we're in event detail
			view := browseIntent.View()
			Expect(view).To(ContainSubstring("Event Details"))

			// Press 'm'
			browseIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			result := browseIntent.Result()
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should cancel intent when 'q' is pressed", func() {
			browseIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

			result := browseIntent.Result()
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should show updated footer with 'm' option in event detail view", func() {
			view := browseIntent.View()
			Expect(view).To(ContainSubstring("m: Main menu"))
		})
	})
})
