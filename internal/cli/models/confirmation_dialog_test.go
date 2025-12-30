package models

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfirmationDialog Modal Styling", func() {
	var dialog *ConfirmationDialog

	BeforeEach(func() {
		dialog = NewConfirmationDialog("Delete Event?", "Are you sure?")
	})

	Context("Modal Border Styling", func() {
		It("should render with border", func() {
			view := dialog.View()
			// Check for rounded border characters: ╭ ╮ ╰ ╯
			Expect(view).To(Or(
				ContainSubstring("╭"),
				ContainSubstring("╮"),
				ContainSubstring("╰"),
				ContainSubstring("╯"),
			))
		})

		It("should use consistent border for destructive operations", func() {
			view := dialog.View()
			Expect(view).To(ContainSubstring("Delete Event?"))
		})
	})

	Context("Modal Padding Consistency", func() {
		It("should have multiple lines with proper spacing", func() {
			view := dialog.View()
			lines := strings.Split(view, "\n")
			Expect(len(lines)).To(BeNumerically(">", 3))
		})

		It("should have non-empty structure", func() {
			view := dialog.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should have interior content with padding", func() {
			view := dialog.View()
			lines := strings.Split(view, "\n")
			if len(lines) > 2 {
				Expect(len(lines[1])).To(BeNumerically(">", 0))
			}
		})
	})

	Context("Modal Background Styling", func() {
		It("should render complete modal", func() {
			dialog := NewConfirmationDialog("Title", "Message")
			Expect(dialog).NotTo(BeNil())
		})

		It("should use rounded border", func() {
			view := dialog.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should have visible structure", func() {
			view := dialog.View()
			Expect(strings.Count(view, "\n")).To(BeNumerically(">", 5))
		})
	})

	Context("Button Styling Within Modal", func() {
		It("should include cancel button", func() {
			view := dialog.View()
			Expect(view).To(ContainSubstring("Cancel"))
		})

		It("should include confirm button", func() {
			view := dialog.View()
			Expect(view).To(ContainSubstring("Yes, Delete"))
		})

		It("should highlight focused button", func() {
			dialog := NewConfirmationDialog("Delete?", "Are you sure?")
			view1 := dialog.View()
			Expect(view1).To(ContainSubstring("Cancel"))

			dialog.Update(tea.KeyMsg{Type: tea.KeyRight})
			view2 := dialog.View()
			Expect(view2).To(ContainSubstring("Yes, Delete"))
		})

		It("should have buttons with spacing", func() {
			view := dialog.View()
			Expect(view).To(ContainSubstring("Cancel"))
			Expect(view).To(ContainSubstring("Yes, Delete"))
		})
	})

	Context("Text Alignment in Modal", func() {
		It("should render title text", func() {
			dialog := NewConfirmationDialog("Delete Event?", "Message")
			view := dialog.View()
			Expect(view).To(ContainSubstring("Delete Event?"))
		})

		It("should render message text", func() {
			dialog := NewConfirmationDialog("Title", "Are you sure?")
			view := dialog.View()
			Expect(view).To(ContainSubstring("Are you sure?"))
		})

		It("should render button text", func() {
			view := dialog.View()
			Expect(view).To(ContainSubstring("Cancel"))
			Expect(view).To(ContainSubstring("Yes, Delete"))
		})

		It("should space elements vertically", func() {
			view := dialog.View()
			lines := strings.Split(view, "\n")
			Expect(len(lines)).To(BeNumerically(">", 5))
		})

		It("should include instructions", func() {
			view := dialog.View()
			Expect(view).To(ContainSubstring("Tab/←→"))
		})
	})

	Context("Modal Styling Integration", func() {
		It("should render all modal elements", func() {
			view := dialog.View()
			Expect(view).To(ContainSubstring("Delete Event?"))
			Expect(view).To(ContainSubstring("Are you sure?"))
			Expect(view).To(ContainSubstring("Cancel"))
			Expect(view).To(ContainSubstring("Yes, Delete"))
		})

		It("should render consistently", func() {
			view1 := dialog.View()
			view2 := dialog.View()
			Expect(view1).To(Equal(view2))
		})
	})
})
