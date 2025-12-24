package models

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FilterModel", func() {
	var (
		model *FilterModel
	)

	BeforeEach(func() {
		model = NewFilterModel()
	})

	Describe("Initialization", func() {
		It("should initialize with empty filters", func() {
			Expect(model.GetTags()).To(BeEmpty())
			Expect(model.GetStartDate()).To(BeNil())
			Expect(model.GetEndDate()).To(BeNil())
		})

		It("should initialize with default focus index", func() {
			Expect(model.focusIndex).To(Equal(0))
		})
	})

	Describe("Tag selection", func() {
		It("should add a tag when selected", func() {
			model.ToggleTag("technical")
			Expect(model.GetTags()).To(ContainElement("technical"))
		})

		It("should remove a tag when deselected", func() {
			model.ToggleTag("technical")
			model.ToggleTag("technical")
			Expect(model.GetTags()).NotTo(ContainElement("technical"))
		})

		It("should not add duplicate tags", func() {
			model.ToggleTag("technical")
			model.ToggleTag("technical")
			model.ToggleTag("technical")
			count := 0
			for _, tag := range model.GetTags() {
				if tag == "technical" {
					count++
				}
			}
			Expect(count).To(Equal(1))
		})

		It("should support multiple tags", func() {
			model.ToggleTag("technical")
			model.ToggleTag("leadership")
			model.ToggleTag("product")
			Expect(model.GetTags()).To(HaveLen(3))
			Expect(model.GetTags()).To(ContainElement("technical"))
			Expect(model.GetTags()).To(ContainElement("leadership"))
			Expect(model.GetTags()).To(ContainElement("product"))
		})

		It("should check if tag is selected", func() {
			model.ToggleTag("technical")
			Expect(model.IsTagSelected("technical")).To(BeTrue())
			Expect(model.IsTagSelected("leadership")).To(BeFalse())
		})
	})

	Describe("Date range filter", func() {
		It("should set start date", func() {
			startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			model.SetStartDate(&startDate)
			Expect(model.GetStartDate()).NotTo(BeNil())
			Expect(*model.GetStartDate()).To(Equal(startDate))
		})

		It("should set end date", func() {
			endDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
			model.SetEndDate(&endDate)
			Expect(model.GetEndDate()).NotTo(BeNil())
			Expect(*model.GetEndDate()).To(Equal(endDate))
		})

		It("should clear start date when set to nil", func() {
			startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			model.SetStartDate(&startDate)
			model.SetStartDate(nil)
			Expect(model.GetStartDate()).To(BeNil())
		})

		It("should clear end date when set to nil", func() {
			endDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
			model.SetEndDate(&endDate)
			model.SetEndDate(nil)
			Expect(model.GetEndDate()).To(BeNil())
		})
	})

	Describe("Filter conversion to repository ListFilters", func() {
		It("should convert to empty ListFilters when no filters set", func() {
			filters := model.ToListFilters()
			Expect(filters.Tags).To(BeEmpty())
			Expect(filters.StartDate).To(BeNil())
			Expect(filters.EndDate).To(BeNil())
		})

		It("should convert tags to ListFilters", func() {
			model.ToggleTag("technical")
			model.ToggleTag("leadership")
			filters := model.ToListFilters()
			Expect(filters.Tags).To(HaveLen(2))
			Expect(filters.Tags).To(ContainElement("technical"))
			Expect(filters.Tags).To(ContainElement("leadership"))
		})

		It("should convert date range to ListFilters", func() {
			startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
			model.SetStartDate(&startDate)
			model.SetEndDate(&endDate)
			filters := model.ToListFilters()
			Expect(filters.StartDate).NotTo(BeNil())
			Expect(filters.EndDate).NotTo(BeNil())
			Expect(*filters.StartDate).To(Equal(startDate))
			Expect(*filters.EndDate).To(Equal(endDate))
		})

		It("should convert all filters together", func() {
			model.ToggleTag("technical")
			startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			model.SetStartDate(&startDate)
			filters := model.ToListFilters()
			Expect(filters.Tags).To(ContainElement("technical"))
			Expect(filters.StartDate).NotTo(BeNil())
		})
	})

	Describe("Reset filters", func() {
		It("should clear all filters", func() {
			model.ToggleTag("technical")
			startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			model.SetStartDate(&startDate)

			model.Reset()

			Expect(model.GetTags()).To(BeEmpty())
			Expect(model.GetStartDate()).To(BeNil())
			Expect(model.GetEndDate()).To(BeNil())
		})
	})

	Describe("Navigation", func() {
		It("should navigate forward through fields", func() {
			initialIndex := model.focusIndex
			model.Next()
			Expect(model.focusIndex).To(Equal(initialIndex + 1))
		})

		It("should navigate backward through fields", func() {
			model.focusIndex = 2
			model.Previous()
			Expect(model.focusIndex).To(Equal(1))
		})

		It("should not go below zero on previous", func() {
			model.focusIndex = 0
			model.Previous()
			Expect(model.focusIndex).To(Equal(0))
		})
	})

	Describe("IsActive", func() {
		It("should return false when no filters set", func() {
			Expect(model.IsActive()).To(BeFalse())
		})

		It("should return true when tags filter set", func() {
			model.ToggleTag("technical")
			Expect(model.IsActive()).To(BeTrue())
		})

		It("should return true when date filter set", func() {
			startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			model.SetStartDate(&startDate)
			Expect(model.IsActive()).To(BeTrue())
		})
	})

	Describe("Error handling", func() {
		It("should set and get errors", func() {
			testErr := errors.New("test error")
			model.SetError(testErr)
			Expect(model.GetError()).To(Equal(testErr))
		})

		It("should clear errors on reset", func() {
			testErr := errors.New("test error")
			model.SetError(testErr)
			model.Reset()
			Expect(model.GetError()).To(BeNil())
		})
	})
})
