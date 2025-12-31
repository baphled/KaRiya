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

		It("initializes with empty breadcrumbs", func() {
			Expect(header.GetBreadcrumbs()).To(HaveLen(0))
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

	Describe("Breadcrumbs", func() {
		It("sets breadcrumbs", func() {
			breadcrumbs := []string{"Home", "List", "Event"}
			header.SetBreadcrumbs(breadcrumbs)
			Expect(header.GetBreadcrumbs()).To(Equal(breadcrumbs))
		})

		It("adds single breadcrumb", func() {
			header.AddBreadcrumb("Home")
			Expect(header.GetBreadcrumbs()).To(Equal([]string{"Home"}))
		})

		It("clears breadcrumbs", func() {
			header.SetBreadcrumbs([]string{"Home", "List"})
			header.ClearBreadcrumbs()
			Expect(header.GetBreadcrumbs()).To(HaveLen(0))
		})

		It("renders breadcrumb navigation", func() {
			header.SetBreadcrumbs([]string{"Home", "List", "Details"})
			view := header.View()
			Expect(view).To(ContainSubstring("Home"))
			Expect(view).To(ContainSubstring("Details"))
		})

		It("includes separator between breadcrumbs", func() {
			header.SetBreadcrumbs([]string{"Home", "List"})
			view := header.View()
			Expect(view).To(ContainSubstring(">"))
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

		It("renders breadcrumbs and title", func() {
			header.SetBreadcrumbs([]string{"Home", "List"})
			view := header.View()
			Expect(view).To(ContainSubstring("Home"))
			Expect(view).To(ContainSubstring("Main Title"))
		})

		It("renders breadcrumbs, title, and subtitle", func() {
			header.SetBreadcrumbs([]string{"Home", "List"})
			header.SetSubtitle("Details view")
			view := header.View()
			Expect(view).To(ContainSubstring("Home"))
			Expect(view).To(ContainSubstring("Main Title"))
			Expect(view).To(ContainSubstring("Details view"))
		})
	})

	Describe("Edge Cases", func() {
		It("handles empty title", func() {
			emptyHeader := NewHeader("", 80)
			view := emptyHeader.View()
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("handles single breadcrumb", func() {
			header.SetBreadcrumbs([]string{"Home"})
			view := header.View()
			Expect(view).To(ContainSubstring("Home"))
		})

		It("handles many breadcrumbs", func() {
			many := []string{"Home", "List", "Filter", "Event", "Details", "Edit"}
			header.SetBreadcrumbs(many)
			view := header.View()
			Expect(view).To(ContainSubstring("Home"))
			Expect(view).To(ContainSubstring("Edit"))
		})

		It("handles unicode in title", func() {
			unicodeHeader := NewHeader("événement", 80)
			view := unicodeHeader.View()
			Expect(view).To(ContainSubstring("événement"))
		})
	})

	Describe("Breadcrumb Click Detection", func() {
		It("detects click on first breadcrumb", func() {
			header := NewHeader("Main Title", 80)
			header.SetBreadcrumbs([]string{"Home", "List", "Details"})

			// Simulate click at position where "Home" breadcrumb is rendered
			clickIndex := header.GetClickedBreadcrumbIndex(0, 0)
			Expect(clickIndex).To(Equal(0)) // Should detect "Home"
		})

		It("detects click on middle breadcrumb", func() {
			header := NewHeader("Main Title", 80)
			header.SetBreadcrumbs([]string{"Home", "List", "Details"})

			// Simulate click at position where "List" breadcrumb is rendered
			// "Home" (4 chars) + " > " (3 chars) = 7 chars offset
			clickIndex := header.GetClickedBreadcrumbIndex(7, 0)
			Expect(clickIndex).To(Equal(1)) // Should detect "List"
		})

		It("detects click on last breadcrumb", func() {
			header := NewHeader("Main Title", 80)
			header.SetBreadcrumbs([]string{"Home", "List", "Details"})

			// Simulate click at position where "Details" breadcrumb is rendered
			// "Home" (4) + " > " (3) + "List" (4) + " > " (3) = 14 chars offset
			clickIndex := header.GetClickedBreadcrumbIndex(14, 0)
			Expect(clickIndex).To(Equal(2)) // Should detect "Details"
		})

		It("returns -1 when click is not on any breadcrumb", func() {
			header := NewHeader("Main Title", 80)
			header.SetBreadcrumbs([]string{"Home", "List"})

			// Click far outside breadcrumb area
			clickIndex := header.GetClickedBreadcrumbIndex(100, 0)
			Expect(clickIndex).To(Equal(-1))
		})

		It("returns -1 when no breadcrumbs exist", func() {
			header := NewHeader("Main Title", 80)

			clickIndex := header.GetClickedBreadcrumbIndex(0, 0)
			Expect(clickIndex).To(Equal(-1))
		})
	})
})
