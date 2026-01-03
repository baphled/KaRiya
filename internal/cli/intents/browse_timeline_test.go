package intents

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BrowseTimelineIntent", func() {
	var (
		intent *BrowseTimelineIntent
		ctx    *BrowseTimelineContext
	)

	BeforeEach(func() {
		// Create test events.
		events := []*career.CareerEvent{
			{
				ID:         "event1",
				Text:       "First event text",
				Date:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Company:    "Company A",
				CreatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Tags:       []string{"go", "backend"},
				Categories: []string{"technical"},
			},
			{
				ID:         "event2",
				Text:       "Second event text",
				Date:       time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				Company:    "Company B",
				CreatedAt:  time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				Tags:       []string{"react", "frontend"},
				Categories: []string{"technical"},
			},
		}

		ctx = &BrowseTimelineContext{
			Events: events,
			InitialFilters: &TimelineFilters{
				Tags:       make([]string, 0),
				Companies:  make([]string, 0),
				Categories: make([]string, 0),
				SortBy:     "date",
				SortOrder:  "desc",
			},
		}

		var err error
		intent, err = NewBrowseTimelineIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(intent).NotTo(BeNil())
	})

	Describe("Creation", func() {
		It("should create a valid intent", func() {
			Expect(intent).NotTo(BeNil())
			Expect(intent.active).To(BeTrue())
		})

		It("should initialize with timeline state", func() {
			Expect(intent.state.currentState).To(Equal(BrowseStateTimeline))
		})

		It("should have all events in filtered list", func() {
			intent.Init()
			Expect(len(intent.state.filteredEvents)).To(Equal(2))
		})
	})

	Describe("Init", func() {
		It("should initialize filtered events", func() {
			cmd := intent.Init()
			Expect(cmd).To(BeNil())
			Expect(len(intent.state.filteredEvents)).To(Equal(2))
		})

		It("should select first event", func() {
			intent.Init()
			Expect(intent.state.selectedEvent).NotTo(BeNil())
			Expect(intent.state.selectedEvent.ID).To(Equal("event1"))
		})
	})

	Describe("Update - Timeline View", func() {
		BeforeEach(func() {
			intent.Init()
		})

		It("should move selection down with down arrow", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(intent.state.selectedIndex).To(Equal(1))
		})

		It("should move selection down with j key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(intent.state.selectedIndex).To(Equal(1))
		})

		It("should move selection up with up arrow", func() {
			intent.state.selectedIndex = 1
			intent.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(intent.state.selectedIndex).To(Equal(0))
		})

		It("should move selection up with k key", func() {
			intent.state.selectedIndex = 1
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(intent.state.selectedIndex).To(Equal(0))
		})

		It("should not move selection below first item", func() {
			intent.state.selectedIndex = 0
			intent.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(intent.state.selectedIndex).To(Equal(0))
		})

		It("should not move selection above last item", func() {
			intent.state.selectedIndex = 1
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(intent.state.selectedIndex).To(Equal(1))
		})

		It("should transition to event detail on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(BrowseStateEventDetail))
		})

		It("should add event to viewed events on selection", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(len(intent.state.viewedEvents)).To(Equal(1))
		})

		It("should cancel on q key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should cancel on ctrl+c", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})
	})

	Describe("Update - Event Detail View", func() {
		BeforeEach(func() {
			intent.Init()
			intent.state.currentState = BrowseStateEventDetail
			intent.state.selectedEvent = intent.state.filteredEvents[0]
		})

		It("should complete on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.result.Status).To(Equal(Completed))
		})

		It("should go back to timeline on esc", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(BrowseStateTimeline))
		})

		It("should cancel on q key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			intent.Init()
		})

		It("should render timeline view", func() {
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Browse Timeline"))
		})

		It("should show events in timeline", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("First event text"))
			Expect(view).To(ContainSubstring("Second event text"))
		})

		It("should show selected event with marker", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring(">"))
		})

		It("should render event detail view", func() {
			intent.state.currentState = BrowseStateEventDetail
			intent.state.selectedEvent = intent.state.filteredEvents[0]
			view := intent.View()
			Expect(view).To(ContainSubstring("Event Details"))
			Expect(view).To(ContainSubstring("Company A"))
		})
	})

	Describe("Filtering", func() {
		BeforeEach(func() {
			intent.Init()
		})

		It("should filter by tag", func() {
			intent.state.filters.Tags = []string{"go"}
			intent.applyFilters()
			Expect(len(intent.state.filteredEvents)).To(Equal(1))
			Expect(intent.state.filteredEvents[0].ID).To(Equal("event1"))
		})

		It("should filter by search text", func() {
			intent.state.filters.SearchText = "Second"
			intent.applyFilters()
			Expect(len(intent.state.filteredEvents)).To(Equal(1))
			Expect(intent.state.filteredEvents[0].ID).To(Equal("event2"))
		})

		It("should sort by date ascending", func() {
			intent.state.filters.SortBy = "date"
			intent.state.filters.SortOrder = "asc"
			intent.applyFilters()
			Expect(intent.state.filteredEvents[0].ID).To(Equal("event1"))
			Expect(intent.state.filteredEvents[1].ID).To(Equal("event2"))
		})

		It("should sort by date descending", func() {
			intent.state.filters.SortBy = "date"
			intent.state.filters.SortOrder = "desc"
			intent.applyFilters()
			Expect(intent.state.filteredEvents[0].ID).To(Equal("event2"))
			Expect(intent.state.filteredEvents[1].ID).To(Equal("event1"))
		})

		It("should sort by text", func() {
			intent.state.filters.SortBy = "text"
			intent.state.filters.SortOrder = "asc"
			intent.applyFilters()
			Expect(intent.state.filteredEvents[0].Text).To(Equal("First event text"))
			Expect(intent.state.filteredEvents[1].Text).To(Equal("Second event text"))
		})
	})

	Describe("Result Handling", func() {
		BeforeEach(func() {
			intent.Init()
		})

		It("should return completed result on success", func() {
			intent.state.currentState = BrowseStateEventDetail
			intent.state.selectedEvent = intent.state.filteredEvents[0]
			intent.setCompleted()

			result := intent.Result()
			Expect(result.Status).To(Equal(Completed))
		})

		It("should include selected event in result", func() {
			intent.state.selectedEvent = intent.state.filteredEvents[0]
			intent.setCompleted()

			Expect(intent.result.Data.SelectedEvent).To(Equal(intent.state.filteredEvents[0]))
		})

		It("should return cancelled result", func() {
			intent.setCancelled()
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should return failed result", func() {
			intent.setFailed("TEST_ERROR", "Test error message", nil)
			Expect(intent.result.Status).To(Equal(Failed))
			Expect(intent.result.Error.Code).To(Equal("TEST_ERROR"))
		})

		It("should include metadata in result", func() {
			intent.state.selectedEvent = intent.state.filteredEvents[0]
			intent.setCompleted()

			Expect(intent.result.Metadata).To(HaveKey("selected_index"))
			Expect(intent.result.Metadata).To(HaveKey("timestamp"))
		})
	})

	Describe("Message Handling", func() {
		BeforeEach(func() {
			intent.Init()
		})

		It("should handle EventSelectedMsg", func() {
			event := intent.state.filteredEvents[1]
			intent.Update(EventSelectedMsg{
				Event: event,
				Index: 1,
			})

			Expect(intent.state.currentState).To(Equal(BrowseStateEventDetail))
			Expect(intent.state.selectedEvent).To(Equal(event))
		})

		It("should handle FilterChangedMsg", func() {
			newFilters := &TimelineFilters{
				SearchText: "First",
				SortBy:     "date",
				SortOrder:  "asc",
			}

			intent.Update(FilterChangedMsg{
				Filters: newFilters,
			})

			Expect(intent.state.filters).To(Equal(newFilters))
		})
	})

	Describe("Edge Cases", func() {
		It("should handle empty events list", func() {
			emptyCtx := &BrowseTimelineContext{
				Events: make([]*career.CareerEvent, 0),
			}
			_ = emptyCtx.Validate()

			emptyIntent, err := NewBrowseTimelineIntent(emptyCtx)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(emptyIntent.state.filteredEvents)).To(Equal(0))
		})

		It("should handle inactive intent", func() {
			intent.active = false
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())
		})

		It("should handle nil selected event in detail view", func() {
			intent.state.currentState = BrowseStateEventDetail
			intent.state.selectedEvent = nil
			view := intent.View()
			Expect(view).To(ContainSubstring("No event selected"))
		})
	})
})
