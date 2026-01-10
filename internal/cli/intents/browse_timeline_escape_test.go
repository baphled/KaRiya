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

		// Note: 'h' key is vim-style left navigation, not home
		// Going home is done by pressing Esc (KeyBack) from root state

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := browseIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

			// q should return tea.Quit, not cancel the intent
			Expect(cmd).ToNot(BeNil())
			// The intent result should be nil (not cancelled, app is quitting)
			result := browseIntent.Result()
			Expect(result).To(BeNil())
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

		// Note: 'h' key is vim-style left navigation, not home
		// Going home from detail is done by pressing Esc twice (back to timeline, then to main menu)

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := browseIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

			// q should return tea.Quit, not cancel the intent
			Expect(cmd).ToNot(BeNil())
			// The intent result should be nil (not cancelled, app is quitting)
			result := browseIntent.Result()
			Expect(result).To(BeNil())
		})

		It("should render event detail view correctly", func() {
			// Footer is now rendered by StandardView, not by View() method
			// This test verifies the view renders correctly
			view := browseIntent.View()
			Expect(view).NotTo(BeEmpty()) // Verify view renders
			Expect(view).To(ContainSubstring("Event Details"))
		})
	})
})
