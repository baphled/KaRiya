package intents

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BrowseTimeline - Delete Event Workflow", func() {
	var (
		intent     *BrowseTimelineIntent
		ctx        *BrowseTimelineContext
		testEvents []*career.CareerEvent
	)

	BeforeEach(func() {
		testEvents = []*career.CareerEvent{
			{
				ID:        "event-1",
				Text:      "First event for testing delete workflow",
				Date:      time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				Company:   "Test Company A",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "event-2",
				Text:      "Second event for testing delete workflow",
				Date:      time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
				Company:   "Test Company B",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "event-3",
				Text:      "Third event for testing delete workflow",
				Date:      time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
				Company:   "Test Company C",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		ctx = &BrowseTimelineContext{
			Events: testEvents,
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
		intent.Init()
	})

	Describe("updateDeleteConfirm", func() {
		BeforeEach(func() {
			// Set up delete confirmation state
			intent.state.currentState = BrowseStateDeleteConfirm
			intent.state.selectedEvent = testEvents[0]
		})

		Context("when pressing 'y' to confirm delete", func() {
			Context("without service (nil CLIEventService)", func() {
				BeforeEach(func() {
					intent.context.CLIEventService = nil
				})

				It("should not crash when service is nil", func() {
					Expect(func() {
						intent.updateDeleteConfirm(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
					}).NotTo(Panic())
				})

				It("should remain in delete confirm state", func() {
					intent.updateDeleteConfirm(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
					Expect(intent.state.currentState).To(Equal(BrowseStateDeleteConfirm))
				})
			})
		})

		Context("when pressing 'Y' (uppercase) to confirm delete", func() {
			BeforeEach(func() {
				intent.context.CLIEventService = nil
			})

			It("should handle uppercase Y", func() {
				Expect(func() {
					intent.updateDeleteConfirm(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'Y'}})
				}).NotTo(Panic())
			})
		})

		Context("when pressing 'n' to cancel delete", func() {
			It("should transition back to event detail state", func() {
				intent.updateDeleteConfirm(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

				Expect(intent.state.currentState).To(Equal(BrowseStateEventDetail))
			})

			It("should clear any delete error", func() {
				intent.state.deleteError = errors.New("previous error")
				intent.updateDeleteConfirm(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

				Expect(intent.state.deleteError).To(BeNil())
			})

			It("should preserve selected event", func() {
				originalEvent := intent.state.selectedEvent
				intent.updateDeleteConfirm(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

				Expect(intent.state.selectedEvent).To(Equal(originalEvent))
			})
		})

		Context("when pressing 'N' (uppercase) to cancel delete", func() {
			It("should transition back to event detail state", func() {
				intent.updateDeleteConfirm(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}})

				Expect(intent.state.currentState).To(Equal(BrowseStateEventDetail))
			})
		})

		Context("when pressing Escape to cancel delete", func() {
			It("should transition back to event detail state", func() {
				intent.updateDeleteConfirm(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.state.currentState).To(Equal(BrowseStateEventDetail))
			})

			It("should clear any delete error", func() {
				intent.state.deleteError = errors.New("previous error")
				intent.updateDeleteConfirm(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.state.deleteError).To(BeNil())
			})
		})

		Context("when pressing other keys", func() {
			It("should remain in delete confirm state", func() {
				intent.updateDeleteConfirm(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

				Expect(intent.state.currentState).To(Equal(BrowseStateDeleteConfirm))
			})

			It("should not modify selected event", func() {
				originalEvent := intent.state.selectedEvent
				intent.updateDeleteConfirm(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

				Expect(intent.state.selectedEvent).To(Equal(originalEvent))
			})
		})
	})

	Describe("removeEventFromList", func() {
		BeforeEach(func() {
			intent.Init()
			// Ensure we have 3 events in filtered list
			Expect(intent.state.filteredEvents).To(HaveLen(3))
		})

		It("should remove event from context.Events", func() {
			originalCount := len(intent.context.Events)
			intent.removeEventFromList("event-1")

			Expect(intent.context.Events).To(HaveLen(originalCount - 1))
		})

		It("should remove event from filteredEvents", func() {
			originalCount := len(intent.state.filteredEvents)
			intent.removeEventFromList("event-1")

			Expect(intent.state.filteredEvents).To(HaveLen(originalCount - 1))
		})

		It("should not contain the removed event", func() {
			intent.removeEventFromList("event-2")

			for _, evt := range intent.context.Events {
				Expect(evt.ID).NotTo(Equal("event-2"))
			}
			for _, evt := range intent.state.filteredEvents {
				Expect(evt.ID).NotTo(Equal("event-2"))
			}
		})

		It("should preserve other events", func() {
			intent.removeEventFromList("event-2")

			// Check event-1 and event-3 still exist
			eventIDs := make([]string, 0)
			for _, evt := range intent.context.Events {
				eventIDs = append(eventIDs, evt.ID)
			}
			Expect(eventIDs).To(ContainElement("event-1"))
			Expect(eventIDs).To(ContainElement("event-3"))
		})

		It("should handle removing non-existent event gracefully", func() {
			originalCount := len(intent.context.Events)

			Expect(func() {
				intent.removeEventFromList("non-existent-id")
			}).NotTo(Panic())

			Expect(intent.context.Events).To(HaveLen(originalCount))
		})

		It("should handle removing from single-event list", func() {
			// Set up single event
			intent.context.Events = []*career.CareerEvent{testEvents[0]}
			intent.state.filteredEvents = []*career.CareerEvent{testEvents[0]}

			intent.removeEventFromList("event-1")

			Expect(intent.context.Events).To(HaveLen(0))
			Expect(intent.state.filteredEvents).To(HaveLen(0))
		})
	})

	Describe("getContext", func() {
		It("should return a non-nil context", func() {
			ctx := intent.getContext()
			Expect(ctx).NotTo(BeNil())
		})

		It("should return a background context", func() {
			ctx := intent.getContext()
			Expect(ctx).To(Equal(context.Background()))
		})
	})

	Describe("viewDeleteConfirm", func() {
		BeforeEach(func() {
			intent.state.currentState = BrowseStateDeleteConfirm
			intent.state.selectedEvent = testEvents[0]
		})

		It("should render event details", func() {
			view := intent.viewDeleteConfirm()

			Expect(view).To(ContainSubstring("Delete Event"))
			Expect(view).To(ContainSubstring("2025-01-15"))
			Expect(view).To(ContainSubstring("Test Company A"))
		})

		It("should render event text", func() {
			view := intent.viewDeleteConfirm()

			Expect(view).To(ContainSubstring("First event"))
		})

		It("should show cannot be undone warning", func() {
			view := intent.viewDeleteConfirm()

			Expect(view).To(ContainSubstring("cannot be undone"))
		})

		Context("with no selected event", func() {
			BeforeEach(func() {
				intent.state.selectedEvent = nil
			})

			It("should show no event selected message", func() {
				view := intent.viewDeleteConfirm()

				Expect(view).To(ContainSubstring("No event selected"))
			})
		})

		Context("with delete error", func() {
			BeforeEach(func() {
				intent.state.deleteError = errors.New("Database connection failed")
			})

			It("should display error message", func() {
				view := intent.viewDeleteConfirm()

				Expect(view).To(ContainSubstring("Error"))
				Expect(view).To(ContainSubstring("Database connection failed"))
			})
		})

		Context("with long event text", func() {
			BeforeEach(func() {
				intent.state.selectedEvent = &career.CareerEvent{
					ID:      "long-text-event",
					Text:    "This is a very long event text that should be truncated for display purposes because it exceeds the maximum length allowed in the delete confirmation dialog and we want to ensure the UI remains clean and readable",
					Date:    time.Now(),
					Company: "Test Co",
				}
			})

			It("should truncate long text with ellipsis", func() {
				view := intent.viewDeleteConfirm()

				Expect(view).To(ContainSubstring("..."))
			})
		})

		Context("with event without company", func() {
			BeforeEach(func() {
				intent.state.selectedEvent = &career.CareerEvent{
					ID:      "no-company-event",
					Text:    "Event without company",
					Date:    time.Now(),
					Company: "",
				}
			})

			It("should not show company line", func() {
				view := intent.viewDeleteConfirm()

				Expect(view).NotTo(ContainSubstring("Company:"))
			})
		})
	})

	Describe("Delete workflow integration", func() {
		BeforeEach(func() {
			intent.Init()
			// Select first event and go to delete confirm
			intent.state.selectedEvent = intent.state.filteredEvents[0]
			intent.state.currentState = BrowseStateDeleteConfirm
		})

		It("should allow cancel and return to event detail", func() {
			intent.updateDeleteConfirm(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

			Expect(intent.state.currentState).To(Equal(BrowseStateEventDetail))
			Expect(intent.state.selectedEvent).NotTo(BeNil())
		})

		It("should preserve list integrity after cancel", func() {
			originalCount := len(intent.state.filteredEvents)
			intent.updateDeleteConfirm(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

			Expect(intent.state.filteredEvents).To(HaveLen(originalCount))
		})
	})
})
