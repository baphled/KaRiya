package intents

import (
	"fmt"
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
			// After sorting by date desc, event2 (2025-01-02) is first
			Expect(intent.state.selectedEvent.ID).To(Equal("event2"))
		})

		It("should sort events chronologically (latest first) on Init", func() {
			// Create events in non-chronological order
			events := []*career.CareerEvent{
				{
					ID:        "old",
					Text:      "Old event",
					Date:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					Tags:      []string{"test"},
				},
				{
					ID:        "new",
					Text:      "New event",
					Date:      time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
					Tags:      []string{"test"},
				},
				{
					ID:        "middle",
					Text:      "Middle event",
					Date:      time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
					Tags:      []string{"test"},
				},
			}

			ctx := &BrowseTimelineContext{Events: events}
			newIntent, err := NewBrowseTimelineIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			newIntent.Init()

			// Verify latest event is first (desc order)
			Expect(newIntent.state.filteredEvents[0].ID).To(Equal("new"))
			Expect(newIntent.state.filteredEvents[1].ID).To(Equal("middle"))
			Expect(newIntent.state.filteredEvents[2].ID).To(Equal("old"))
		})

		It("should use CreatedAt as secondary sort for same-date events", func() {
			sameDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
			events := []*career.CareerEvent{
				{
					ID:        "first-created",
					Text:      "First created",
					Date:      sameDate,
					CreatedAt: time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC),
					Tags:      []string{"test"},
				},
				{
					ID:        "last-created",
					Text:      "Last created",
					Date:      sameDate,
					CreatedAt: time.Date(2025, 1, 15, 17, 0, 0, 0, time.UTC),
					Tags:      []string{"test"},
				},
			}

			ctx := &BrowseTimelineContext{Events: events}
			newIntent, err := NewBrowseTimelineIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			newIntent.Init()

			// Latest created should be first when dates are equal
			Expect(newIntent.state.filteredEvents[0].ID).To(Equal("last-created"))
			Expect(newIntent.state.filteredEvents[1].ID).To(Equal("first-created"))
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
			// Note: "Browse Timeline" title is now in StandardView breadcrumbs, not View() output
		})

		It("should show events in timeline", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("First event text"))
			Expect(view).To(ContainSubstring("Second event text"))
		})

		It("should show selected event with marker", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("▶"))
		})

		It("should render event detail view", func() {
			intent.state.currentState = BrowseStateEventDetail
			intent.state.selectedEvent = intent.state.filteredEvents[0]
			view := intent.View()
			Expect(view).To(ContainSubstring("Event Details"))
			// After sorting by date desc, event2 (Company B) is first
			Expect(view).To(ContainSubstring("Company B"))
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

	Describe("Pagination", func() {
		var (
			manyEventsIntent *BrowseTimelineIntent
			manyEventsCtx    *BrowseTimelineContext
		)

		BeforeEach(func() {
			// Create 35 events to span multiple pages (pageSize = 15)
			events := make([]*career.CareerEvent, 35)
			for i := 0; i < 35; i++ {
				dateDay := (i % 28) + 1 // Ensure valid day
				events[i] = &career.CareerEvent{
					ID:         fmt.Sprintf("event%02d", i),
					Text:       fmt.Sprintf("Event %02d", i+1),
					Date:       time.Date(2025, 1, dateDay, 0, 0, 0, 0, time.UTC),
					Company:    fmt.Sprintf("Company %02d", i+1),
					CreatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					Tags:       []string{"test"},
					Categories: []string{"testing"},
				}
			}

			manyEventsCtx = &BrowseTimelineContext{
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
			manyEventsIntent, err = NewBrowseTimelineIntent(manyEventsCtx)
			Expect(err).NotTo(HaveOccurred())
			manyEventsIntent.Init()
		})

		It("should display correct events on first page", func() {
			// Verify we're on page 1
			view := manyEventsIntent.View()
			Expect(view).To(ContainSubstring("Page 1 of 3"))

			// Verify table shows first 15 events
			rows := manyEventsIntent.table.Rows()
			Expect(len(rows)).To(Equal(15))
		})

		It("should update table rows when navigating to next page", func() {
			// Navigate to page 2 using ctrl+d key
			manyEventsIntent.Update(tea.KeyMsg{Type: tea.KeyCtrlD})

			// Verify we're on page 2
			view := manyEventsIntent.View()
			Expect(view).To(ContainSubstring("Page 2 of 3"))

			// Verify the selected index is now in the second page range
			Expect(manyEventsIntent.state.selectedIndex).To(Equal(15))

			// FAILING TEST: Verify table shows the correct events for page 2
			// Currently the table shows ALL events (all 35 rows) instead of just the current page (15 rows)
			rows := manyEventsIntent.table.Rows()
			Expect(len(rows)).To(Equal(15), "Table should show only 15 events for page 2, but shows %d events", len(rows))
		})

		It("should update table rows when navigating to last page", func() {
			// Navigate to page 3 using ctrl+d twice
			manyEventsIntent.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
			manyEventsIntent.Update(tea.KeyMsg{Type: tea.KeyCtrlD})

			// Verify we're on page 3
			view := manyEventsIntent.View()
			Expect(view).To(ContainSubstring("Page 3 of 3"))

			// Verify the selected index is now in the third page range
			Expect(manyEventsIntent.state.selectedIndex).To(Equal(30))

			// FAILING TEST: Verify table shows the correct events for page 3
			// Currently the table shows ALL events (all 35 rows) instead of just the current page (5 rows)
			rows := manyEventsIntent.table.Rows()
			Expect(len(rows)).To(Equal(5), "Table should show only 5 events for page 3, but shows %d events", len(rows))
		})
	})
})
