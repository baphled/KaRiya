package components

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Header", func() {
	var header HeaderModel

	BeforeEach(func() {
		header = NewHeader("Main Title", 80)
	})

	Describe("Creation", func() {
		It("creates a header with title", func() {
			Expect(header.GetTitle()).To(Equal("Main Title"))
		})

		It("initializes with empty subtitle", func() {
			Expect(header.GetSubtitle()).To(Equal(""))
		})
	})

	Describe("Subtitle", func() {
		It("sets subtitle", func() {
			header.SetSubtitle("Secondary text")
			Expect(header.GetSubtitle()).To(Equal("Secondary text"))
		})

		It("renders with subtitle", func() {
			header.SetSubtitle("Event details view")
			view := header.View()
			Expect(view).To(ContainSubstring("Event details view"))
		})
	})

	Describe("Configuration", func() {
		It("sets width", func() {
			header.SetWidth(100)
			Expect(header.width).To(Equal(100))
		})

		It("sets height", func() {
			header.SetHeight(3)
			Expect(header.height).To(Equal(3))
		})

		It("toggles border display", func() {
			header.SetShowBorder(true)
			Expect(header.showBorder).To(BeTrue())
		})
	})

	Describe("Rendering", func() {
		It("renders title", func() {
			view := header.View()
			Expect(view).To(ContainSubstring("Main Title"))
		})

		It("renders with border when enabled", func() {
			header.SetShowBorder(true)
			view := header.View()
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("handles empty header width", func() {
			header.SetWidth(0)
			view := header.View()
			Expect(view).To(Equal(""))
		})

	})

	Describe("Edge Cases", func() {
		It("handles empty title", func() {
			emptyHeader := NewHeader("", 80)
			view := emptyHeader.View()
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("handles unicode in title", func() {
			unicodeHeader := NewHeader("événement", 80)
			view := unicodeHeader.View()
			Expect(view).To(ContainSubstring("événement"))
		})
	})

})
