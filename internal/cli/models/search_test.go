package models

import (
	"errors"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SearchModel", func() {
	var (
		model *SearchModel
	)

	BeforeEach(func() {
		model = NewSearchModel()
	})

	Describe("Initialization", func() {
		It("should initialize with empty query", func() {
			Expect(model.GetQuery()).To(Equal(""))
		})

		It("should initialize with debounce delay", func() {
			Expect(model.GetDebounceDelay()).To(BeNumerically(">=", 0))
		})

		It("should initialize with no error", func() {
			Expect(model.GetError()).To(BeNil())
		})
	})

	Describe("Query management", func() {
		It("should set search query", func() {
			model.SetQuery("technical achievement")
			Expect(model.GetQuery()).To(Equal("technical achievement"))
		})

		It("should clear query", func() {
			model.SetQuery("test")
			model.SetQuery("")
			Expect(model.GetQuery()).To(Equal(""))
		})

		It("should handle whitespace in query", func() {
			model.SetQuery("  multiple   spaces  ")
			Expect(model.GetQuery()).To(Equal("  multiple   spaces  "))
		})

		It("should be case-insensitive for matching", func() {
			model.SetQuery("TECHNICAL")
			event := &career.CareerEvent{
				Text: "implemented technical solution",
				Date: time.Now(),
			}
			Expect(model.Matches(event)).To(BeTrue())
		})
	})

	Describe("Event matching", func() {
		It("should match event text", func() {
			model.SetQuery("feature")
			event := &career.CareerEvent{
				Text: "implemented new feature",
				Date: time.Now(),
			}
			Expect(model.Matches(event)).To(BeTrue())
		})

		It("should not match non-matching text", func() {
			model.SetQuery("feature")
			event := &career.CareerEvent{
				Text: "code review session",
				Date: time.Now(),
			}
			Expect(model.Matches(event)).To(BeFalse())
		})

		It("should match partial words", func() {
			model.SetQuery("tech")
			event := &career.CareerEvent{
				Text: "technical leadership",
				Date: time.Now(),
			}
			Expect(model.Matches(event)).To(BeTrue())
		})

		It("should match company name", func() {
			model.SetQuery("acme")
			event := &career.CareerEvent{
				Text:    "led project",
				Date:    time.Now(),
				Company: "Acme Corp",
			}
			Expect(model.Matches(event)).To(BeTrue())
		})

		It("should match project name", func() {
			model.SetQuery("platform")
			event := &career.CareerEvent{
				Text:    "migration work",
				Date:    time.Now(),
				Project: "Platform Upgrade",
			}
			Expect(model.Matches(event)).To(BeTrue())
		})

		It("should return true for empty query", func() {
			model.SetQuery("")
			event := &career.CareerEvent{
				Text: "any text",
				Date: time.Now(),
			}
			Expect(model.Matches(event)).To(BeTrue())
		})

		It("should match all searchable fields", func() {
			model.SetQuery("xyz")
			event := &career.CareerEvent{
				Text:    "XYZ corporation project",
				Date:    time.Now(),
				Company: "XYZ Inc",
				Project: "XYZ Platform",
			}
			Expect(model.Matches(event)).To(BeTrue())
		})
	})

	Describe("Search state", func() {
		It("should indicate search is active when query is not empty", func() {
			model.SetQuery("test")
			Expect(model.IsActive()).To(BeTrue())
		})

		It("should indicate search is inactive when query is empty", func() {
			model.SetQuery("")
			Expect(model.IsActive()).To(BeFalse())
		})
	})

	Describe("Highlighting", func() {
		It("should mark search positions for highlighting", func() {
			model.SetQuery("test")
			text := "this is a test string"
			positions := model.GetHighlightPositions(text)
			Expect(positions).NotTo(BeEmpty())
		})

		It("should find multiple occurrences", func() {
			model.SetQuery("a")
			text := "a cat and a dog and a bird"
			positions := model.GetHighlightPositions(text)
			// Multiple 'a' instances should be found
			Expect(len(positions)).To(BeNumerically(">", 0))
		})

		It("should return empty for non-matching text", func() {
			model.SetQuery("xyz")
			text := "no match here"
			positions := model.GetHighlightPositions(text)
			Expect(positions).To(BeEmpty())
		})

		It("should be case-insensitive for highlighting", func() {
			model.SetQuery("TEST")
			text := "test string"
			positions := model.GetHighlightPositions(text)
			Expect(positions).NotTo(BeEmpty())
		})
	})

	Describe("Error handling", func() {
		It("should set and get errors", func() {
			testErr := errors.New("search error")
			model.SetError(testErr)
			Expect(model.GetError()).To(Equal(testErr))
		})

		It("should clear errors", func() {
			testErr := errors.New("search error")
			model.SetError(testErr)
			model.ClearError()
			Expect(model.GetError()).To(BeNil())
		})

		It("should clear error when query changes", func() {
			testErr := errors.New("search error")
			model.SetError(testErr)
			model.SetQuery("new query")
			Expect(model.GetError()).To(BeNil())
		})
	})

	Describe("Reset", func() {
		It("should clear all search state", func() {
			model.SetQuery("test")
			testErr := errors.New("error")
			model.SetError(testErr)

			model.Reset()

			Expect(model.GetQuery()).To(Equal(""))
			Expect(model.GetError()).To(BeNil())
			Expect(model.IsActive()).To(BeFalse())
		})
	})

	Describe("Debounce", func() {
		It("should track last update time", func() {
			before := time.Now()
			model.SetQuery("test")
			after := time.Now()

			lastUpdate := model.GetLastUpdate()
			Expect(lastUpdate).To(BeTemporally(">=", before))
			Expect(lastUpdate).To(BeTemporally("<=", after))
		})

		It("should indicate if debounce delay has passed", func() {
			model.SetQuery("test")
			// Wait should not have passed immediately
			Expect(model.ShouldSearch()).To(BeFalse())
		})
	})
})
