package cv_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cv "github.com/baphled/kariya/internal/tui/views/cv"
)

var _ = Describe("Helpers", func() {
	Describe("WordWrap", func() {
		It("wraps text at width boundary", func() {
			text := "hello world"
			wrapped := cv.WordWrap(text, 10)
			Expect(wrapped).To(Equal("hello\nworld"))
		})

		It("returns text unchanged when width is zero", func() {
			text := "line one\nline two"
			wrapped := cv.WordWrap(text, 0)
			Expect(wrapped).To(Equal(text))
		})

		It("returns text unchanged when width is negative", func() {
			text := "line one\nline two"
			wrapped := cv.WordWrap(text, -5)
			Expect(wrapped).To(Equal(text))
		})

		It("preserves explicit newlines", func() {
			text := "alpha beta\ngamma delta"
			wrapped := cv.WordWrap(text, 7)
			Expect(wrapped).To(Equal("alpha\nbeta\ngamma\ndelta"))
		})

		It("preserves empty paragraph lines", func() {
			text := "alpha\n\nbeta"
			wrapped := cv.WordWrap(text, 10)
			Expect(wrapped).To(Equal("alpha\n\nbeta"))
		})

		It("keeps single long word on one line", func() {
			text := "supercalifragilisticexpialidocious"
			wrapped := cv.WordWrap(text, 5)
			Expect(wrapped).To(Equal(text))
		})

		It("keeps short text on one line", func() {
			text := "short text"
			wrapped := cv.WordWrap(text, 50)
			Expect(wrapped).To(Equal(text))
		})
	})

	Describe("WrapParagraph", func() {
		It("wraps words when they exceed width", func() {
			text := "alpha beta gamma"
			wrapped := cv.WrapParagraph(text, 10)
			Expect(wrapped).To(Equal("alpha beta\ngamma"))
		})

		It("returns text unchanged when width is zero", func() {
			text := "alpha beta"
			wrapped := cv.WrapParagraph(text, 0)
			Expect(wrapped).To(Equal(text))
		})

		It("places long words on their own line", func() {
			text := "superlongword other"
			wrapped := cv.WrapParagraph(text, 5)
			Expect(wrapped).To(Equal("superlongword\nother"))
		})

		It("keeps multiple words on one line when they fit", func() {
			text := "a bb ccc"
			wrapped := cv.WrapParagraph(text, 20)
			Expect(wrapped).To(Equal(text))
		})

		It("returns empty text unchanged", func() {
			wrapped := cv.WrapParagraph("", 10)
			Expect(wrapped).To(Equal(""))
		})
	})

	Describe("MinInt", func() {
		It("returns a when a is less than b", func() {
			Expect(cv.MinInt(3, 5)).To(Equal(3))
		})

		It("returns b when a is greater than b", func() {
			Expect(cv.MinInt(10, 4)).To(Equal(4))
		})

		It("returns a when a equals b", func() {
			Expect(cv.MinInt(7, 7)).To(Equal(7))
		})

		It("handles negative numbers", func() {
			Expect(cv.MinInt(-5, 2)).To(Equal(-5))
			Expect(cv.MinInt(-1, -3)).To(Equal(-3))
		})
	})
})
