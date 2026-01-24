package browse_timeline_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/browse_timeline"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Result", func() {
	Describe("Construction", func() {
		It("should create an empty result", func() {
			result := &browse_timeline.Result{}
			Expect(result).NotTo(BeNil())
		})

		It("should store selected event", func() {
			event := fixtures.EventWith("event-1", "Backend Developer", "TechCorp", "Platform")
			result := &browse_timeline.Result{
				SelectedEvent: event,
			}
			Expect(result.SelectedEvent).To(Equal(event))
			Expect(result.SelectedEvent.ID).To(Equal("event-1"))
		})

		It("should store final filters", func() {
			filters := &browse_timeline.Filters{
				SearchText: "golang",
				SortBy:     "date",
				SortOrder:  "desc",
			}
			result := &browse_timeline.Result{
				FinalFilters: filters,
			}
			Expect(result.FinalFilters).To(Equal(filters))
			Expect(result.FinalFilters.SearchText).To(Equal("golang"))
		})

		It("should store viewed events", func() {
			events := []*career.CareerEvent{
				fixtures.Event("event-1"),
				fixtures.Event("event-2"),
				fixtures.Event("event-3"),
			}
			result := &browse_timeline.Result{
				ViewedEvents: events,
			}
			Expect(result.ViewedEvents).To(HaveLen(3))
		})

		It("should store selected facts", func() {
			facts := []*career.Fact{
				fixtures.Fact("fact-1", "event-1"),
				fixtures.Fact("fact-2", "event-1"),
			}
			result := &browse_timeline.Result{
				SelectedFacts: facts,
			}
			Expect(result.SelectedFacts).To(HaveLen(2))
		})
	})

	Describe("Nil Handling", func() {
		Context("when fields are nil", func() {
			It("should handle nil selected event", func() {
				result := &browse_timeline.Result{
					SelectedEvent: nil,
				}
				Expect(result.SelectedEvent).To(BeNil())
			})

			It("should handle nil final filters", func() {
				result := &browse_timeline.Result{
					FinalFilters: nil,
				}
				Expect(result.FinalFilters).To(BeNil())
			})

			It("should handle nil viewed events", func() {
				result := &browse_timeline.Result{
					ViewedEvents: nil,
				}
				Expect(result.ViewedEvents).To(BeNil())
			})

			It("should handle nil selected facts", func() {
				result := &browse_timeline.Result{
					SelectedFacts: nil,
				}
				Expect(result.SelectedFacts).To(BeNil())
			})
		})
	})

	Describe("Empty Slices", func() {
		It("should handle empty viewed events", func() {
			result := &browse_timeline.Result{
				ViewedEvents: []*career.CareerEvent{},
			}
			Expect(result.ViewedEvents).NotTo(BeNil())
			Expect(result.ViewedEvents).To(BeEmpty())
		})

		It("should handle empty selected facts", func() {
			result := &browse_timeline.Result{
				SelectedFacts: []*career.Fact{},
			}
			Expect(result.SelectedFacts).NotTo(BeNil())
			Expect(result.SelectedFacts).To(BeEmpty())
		})
	})

	Describe("Complete Result", func() {
		It("should store all fields together", func() {
			selectedEvent := fixtures.EventWith("selected-1", "Selected Event", "TechCorp", "Platform")
			filters := &browse_timeline.Filters{
				SearchText: "backend",
				Tags:       []string{"golang", "api"},
				SortBy:     "date",
				SortOrder:  "desc",
			}
			viewedEvents := []*career.CareerEvent{
				fixtures.Event("viewed-1"),
				fixtures.Event("viewed-2"),
			}
			selectedFacts := []*career.Fact{
				fixtures.Fact("fact-1", "selected-1"),
			}

			result := &browse_timeline.Result{
				SelectedEvent: selectedEvent,
				FinalFilters:  filters,
				ViewedEvents:  viewedEvents,
				SelectedFacts: selectedFacts,
			}

			Expect(result.SelectedEvent.ID).To(Equal("selected-1"))
			Expect(result.FinalFilters.SearchText).To(Equal("backend"))
			Expect(result.ViewedEvents).To(HaveLen(2))
			Expect(result.SelectedFacts).To(HaveLen(1))
		})
	})

	Describe("Cancelled Result Pattern", func() {
		It("should represent cancellation with nil selected event", func() {
			// When user cancels, SelectedEvent is nil but filters may be preserved.
			result := &browse_timeline.Result{
				SelectedEvent: nil,
				FinalFilters: &browse_timeline.Filters{
					SearchText: "previous search",
				},
			}
			Expect(result.SelectedEvent).To(BeNil())
			Expect(result.FinalFilters).NotTo(BeNil())
		})
	})

	Describe("Analytics Support", func() {
		It("should track viewing history for analytics", func() {
			// ViewedEvents supports tracking which events were viewed during session.
			event := fixtures.Event("event-1")
			result := &browse_timeline.Result{
				ViewedEvents: []*career.CareerEvent{
					event,
					fixtures.Event("event-2"),
					event, // Can have duplicates (viewed same event twice).
				},
			}
			Expect(result.ViewedEvents).To(HaveLen(3))
		})
	})
})
