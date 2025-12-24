package models

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SortModel", func() {
	Describe("Initialization", func() {
		It("should initialize with default sort (date descending)", func() {
			model := NewSortModel()
			Expect(model.GetSortBy()).To(Equal("date"))
			Expect(model.GetSortOrder()).To(Equal("desc"))
		})
	})

	Describe("Sort field selection", func() {
		It("should set sort field to date", func() {
			model := NewSortModel()
			model.SetSortBy("date")
			Expect(model.GetSortBy()).To(Equal("date"))
		})

		It("should set sort field to created_at", func() {
			model := NewSortModel()
			model.SetSortBy("created_at")
			Expect(model.GetSortBy()).To(Equal("created_at"))
		})

		It("should set sort field to text", func() {
			model := NewSortModel()
			model.SetSortBy("text")
			Expect(model.GetSortBy()).To(Equal("text"))
		})

		It("should only accept valid sort fields", func() {
			model := NewSortModel()
			// Set invalid field - should keep default
			model.SetSortBy("invalid")
			Expect(model.GetSortBy()).To(Equal("date"))
		})
	})

	Describe("Sort order selection", func() {
		It("should set sort order to ascending", func() {
			model := NewSortModel()
			model.SetSortOrder("asc")
			Expect(model.GetSortOrder()).To(Equal("asc"))
		})

		It("should set sort order to descending", func() {
			model := NewSortModel()
			model.SetSortOrder("desc")
			Expect(model.GetSortOrder()).To(Equal("desc"))
		})

		It("should only accept valid sort orders", func() {
			model := NewSortModel()
			// Set invalid order - should keep default
			model.SetSortOrder("invalid")
			Expect(model.GetSortOrder()).To(Equal("desc"))
		})
	})

	Describe("Sort options", func() {
		It("should provide list of sort fields", func() {
			model := NewSortModel()
			fields := model.GetSortFields()
			Expect(fields).To(ContainElement("date"))
			Expect(fields).To(ContainElement("created_at"))
			Expect(fields).To(ContainElement("text"))
		})

		It("should provide list of sort orders", func() {
			model := NewSortModel()
			orders := model.GetSortOrders()
			Expect(orders).To(ContainElement("asc"))
			Expect(orders).To(ContainElement("desc"))
		})
	})

	Describe("Sort comparison", func() {
		It("should identify when two events need reordering by date descending", func() {
			model := NewSortModel()
			model.SetSortBy("date")
			model.SetSortOrder("desc")

			now := time.Now()
			yesterday := now.Add(-24 * time.Hour)

			event1 := &career.CareerEvent{Text: "older", Date: yesterday}
			event2 := &career.CareerEvent{Text: "newer", Date: now}

			// In descending order, newer should come first
			shouldReorder := model.ShouldReorder(event1, event2)
			Expect(shouldReorder).To(BeTrue())
		})

		It("should keep events in correct order when sorted ascending by date", func() {
			model := NewSortModel()
			model.SetSortBy("date")
			model.SetSortOrder("asc")

			now := time.Now()
			yesterday := now.Add(-24 * time.Hour)

			event1 := &career.CareerEvent{Text: "older", Date: yesterday}
			event2 := &career.CareerEvent{Text: "newer", Date: now}

			// In ascending order, older should come first
			shouldReorder := model.ShouldReorder(event1, event2)
			Expect(shouldReorder).To(BeFalse())
		})

		It("should compare text alphabetically", func() {
			model := NewSortModel()
			model.SetSortBy("text")
			model.SetSortOrder("asc")

			event1 := &career.CareerEvent{Text: "zebra", Date: time.Now()}
			event2 := &career.CareerEvent{Text: "apple", Date: time.Now()}

			// In ascending order, apple should come first
			shouldReorder := model.ShouldReorder(event1, event2)
			Expect(shouldReorder).To(BeTrue())
		})
	})

	Describe("Reset", func() {
		It("should reset to default sort", func() {
			model := NewSortModel()
			model.SetSortBy("text")
			model.SetSortOrder("asc")

			model.Reset()

			Expect(model.GetSortBy()).To(Equal("date"))
			Expect(model.GetSortOrder()).To(Equal("desc"))
		})
	})

	Describe("Toggle sort order", func() {
		It("should toggle between asc and desc", func() {
			model := NewSortModel()
			Expect(model.GetSortOrder()).To(Equal("desc"))

			model.ToggleSortOrder()
			Expect(model.GetSortOrder()).To(Equal("asc"))

			model.ToggleSortOrder()
			Expect(model.GetSortOrder()).To(Equal("desc"))
		})
	})
})

