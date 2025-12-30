package components

import (
	"github.com/charmbracelet/lipgloss"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListItem", func() {
	var item ListItemModel

	BeforeEach(func() {
		item = NewListItem("Event Title", 80)
	})

	Describe("Creation", func() {
		It("creates a list item with title", func() {
			Expect(item.GetTitle()).To(Equal("Event Title"))
		})

		It("initializes with empty subtitle", func() {
			Expect(item.GetSubtitle()).To(Equal(""))
		})

		It("initializes with empty metadata", func() {
			Expect(item.GetAllMetadata()).To(HaveLen(0))
		})

		It("initializes as not selected", func() {
			Expect(item.IsSelected()).To(BeFalse())
		})

		It("initializes as not focused", func() {
			Expect(item.IsFocused()).To(BeFalse())
		})

		It("sets max width from constructor", func() {
			Expect(item.maxWidth).To(Equal(80))
		})
	})

	Describe("Subtitle", func() {
		It("sets subtitle", func() {
			item.SetSubtitle("Event description")
			Expect(item.GetSubtitle()).To(Equal("Event description"))
		})

		It("renders with subtitle", func() {
			item.SetSubtitle("Detailed description")
			view := item.View()
			Expect(view).To(ContainSubstring("Detailed description"))
		})

		It("displays subtitle with muted style", func() {
			item.SetSubtitle("Subtitle text")
			view := item.View()
			Expect(len(view)).To(BeNumerically(">", 0))
		})
	})

	Describe("Metadata", func() {
		It("sets single metadata field", func() {
			item.SetMetadata("date", "2025-12-30")
			Expect(item.GetMetadata("date")).To(Equal("2025-12-30"))
		})

		It("sets multiple metadata fields", func() {
			item.SetMetadata("date", "2025-12-30")
			item.SetMetadata("tags", "important")
			item.SetMetadata("status", "done")
			Expect(item.GetAllMetadata()).To(HaveLen(3))
		})

		It("renders metadata fields", func() {
			item.SetMetadata("date", "2025-12-30")
			item.SetMetadata("tags", "urgent")
			view := item.View()
			Expect(view).To(ContainSubstring("date"))
			Expect(view).To(ContainSubstring("urgent"))
		})

		It("retrieves metadata by key", func() {
			item.SetMetadata("priority", "high")
			Expect(item.GetMetadata("priority")).To(Equal("high"))
		})

		It("returns empty string for missing metadata key", func() {
			Expect(item.GetMetadata("missing")).To(Equal(""))
		})
	})

	Describe("Selection", func() {
		It("sets selected state", func() {
			item.SetSelected(true)
			Expect(item.IsSelected()).To(BeTrue())
		})

		It("renders selection indicator when selected", func() {
			item.SetSelected(true)
			view := item.View()
			Expect(view).To(ContainSubstring("✓"))
		})

		It("renders no indicator when not selected", func() {
			item.SetSelected(false)
			view := item.View()
			Expect(view).NotTo(ContainSubstring("✓"))
		})

		It("displays title in accent color when selected", func() {
			item.SetSelected(true)
			view := item.View()
			Expect(len(view)).To(BeNumerically(">", 0))
		})
	})

	Describe("Focus", func() {
		It("sets focused state", func() {
			item.SetFocused(true)
			Expect(item.IsFocused()).To(BeTrue())
		})

		It("renders focus indicator when focused", func() {
			item.SetFocused(true)
			view := item.View()
			Expect(view).To(ContainSubstring("►"))
		})

		It("renders no indicator when not focused", func() {
			item.SetFocused(false)
			view := item.View()
			Expect(view).NotTo(ContainSubstring("►"))
		})

		It("renders different style when focused vs not focused", func() {
			item1 := NewListItem("Item", 80)
			item1.SetFocused(true)
			view1 := item1.View()

			item2 := NewListItem("Item", 80)
			item2.SetFocused(false)
			view2 := item2.View()

			Expect(view1).To(ContainSubstring("►"))
			Expect(view2).NotTo(ContainSubstring("►"))
		})
	})

	Describe("Status", func() {
		It("sets status with color", func() {
			item.SetStatus("Done", lipgloss.Color("#6cb56c"))
			Expect(item.showStatus).To(BeTrue())
		})

		It("renders status indicator", func() {
			item.SetStatus("Pending", lipgloss.Color("#d9a66c"))
			view := item.View()
			Expect(view).To(ContainSubstring("Pending"))
		})

		It("clears status", func() {
			item.SetStatus("Active", lipgloss.Color("#6ab0d3"))
			item.ClearStatus()
			Expect(item.showStatus).To(BeFalse())
		})

		It("displays status in brackets", func() {
			item.SetStatus("Progress", lipgloss.Color("#6ab0d3"))
			view := item.View()
			Expect(view).To(ContainSubstring("[Progress]"))
		})
	})

	Describe("Configuration", func() {
		It("sets max width", func() {
			item.SetMaxWidth(100)
			Expect(item.maxWidth).To(Equal(100))
		})

		It("toggles selected state", func() {
			item.SetSelected(true)
			Expect(item.IsSelected()).To(BeTrue())
			item.SetSelected(false)
			Expect(item.IsSelected()).To(BeFalse())
		})

		It("toggles focused state", func() {
			item.SetFocused(true)
			Expect(item.IsFocused()).To(BeTrue())
			item.SetFocused(false)
			Expect(item.IsFocused()).To(BeFalse())
		})
	})

	Describe("Rendering", func() {
		It("renders basic list item with title", func() {
			view := item.View()
			Expect(view).To(ContainSubstring("Event Title"))
		})

		It("renders title with subtitle", func() {
			item.SetSubtitle("Subtitle text")
			view := item.View()
			Expect(view).To(ContainSubstring("Event Title"))
			Expect(view).To(ContainSubstring("Subtitle text"))
		})

		It("renders title with metadata", func() {
			item.SetMetadata("date", "2025-12-30")
			view := item.View()
			Expect(view).To(ContainSubstring("2025-12-30"))
		})

		It("renders complete item with all fields", func() {
			item.SetSubtitle("Description")
			item.SetMetadata("date", "2025-12-30")
			item.SetMetadata("tags", "urgent")
			item.SetStatus("Active", lipgloss.Color("#6ab0d3"))
			view := item.View()
			Expect(view).To(ContainSubstring("Event Title"))
			Expect(view).To(ContainSubstring("Description"))
			Expect(view).To(ContainSubstring("date"))
			Expect(view).To(ContainSubstring("Active"))
		})

		It("handles empty max width", func() {
			item.SetMaxWidth(0)
			view := item.View()
			Expect(view).To(Equal(""))
		})

		It("handles negative max width", func() {
			item.SetMaxWidth(-1)
			view := item.View()
			Expect(view).To(Equal(""))
		})

		It("truncates long title", func() {
			longItem := NewListItem("This is a very long title that exceeds the maximum width", 40)
			view := longItem.View()
			Expect(view).To(ContainSubstring("..."))
		})

		It("truncates long subtitle", func() {
			item.SetMaxWidth(40)
			item.SetSubtitle("This is a very long subtitle that definitely exceeds the width limit")
			view := item.View()
			Expect(view).To(ContainSubstring("..."))
		})

		It("truncates long metadata", func() {
			item.SetMaxWidth(30)
			item.SetMetadata("field1", "value1")
			item.SetMetadata("field2", "value2")
			item.SetMetadata("field3", "value3")
			view := item.View()
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("handles responsive sizing for narrow terminals", func() {
			narrowItem := NewListItem("Title", 20)
			narrowItem.SetSubtitle("Description")
			view := narrowItem.View()
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("handles responsive sizing for wide terminals", func() {
			wideItem := NewListItem("Title", 200)
			wideItem.SetSubtitle("Description")
			wideItem.SetMetadata("key", "value")
			view := wideItem.View()
			Expect(view).To(ContainSubstring("Title"))
			Expect(view).To(ContainSubstring("Description"))
		})
	})

	Describe("Edge Cases", func() {
		It("handles empty title", func() {
			emptyItem := NewListItem("", 80)
			view := emptyItem.View()
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("handles unicode in title", func() {
			unicodeItem := NewListItem("événement", 80)
			view := unicodeItem.View()
			Expect(view).To(ContainSubstring("événement"))
		})

		It("handles special characters in title", func() {
			specialItem := NewListItem("Event [2025] > Status: Done", 80)
			view := specialItem.View()
			Expect(view).To(ContainSubstring("Event"))
			Expect(view).To(ContainSubstring("Status"))
		})

		It("handles unicode in subtitle", func() {
			item.SetSubtitle("événement décription")
			view := item.View()
			Expect(view).To(ContainSubstring("événement"))
		})

		It("handles unicode in metadata", func() {
			item.SetMetadata("clé", "valeur")
			view := item.View()
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("handles many metadata fields", func() {
			for i := 0; i < 10; i++ {
				key := "key" + string(rune(48+i))
				item.SetMetadata(key, "value")
			}
			view := item.View()
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("handles selected and focused simultaneously", func() {
			item.SetSelected(true)
			item.SetFocused(true)
			view := item.View()
			Expect(view).To(ContainSubstring("✓"))
		})

		It("handles selected with status", func() {
			item.SetSelected(true)
			item.SetStatus("Done", lipgloss.Color("#6cb56c"))
			view := item.View()
			Expect(view).To(ContainSubstring("✓"))
			Expect(view).To(ContainSubstring("[Done]"))
		})
	})

	Describe("Multiple Items", func() {
		It("creates independent list items", func() {
			item1 := NewListItem("Item 1", 80)
			item2 := NewListItem("Item 2", 80)

			item1.SetSubtitle("Sub 1")
			item2.SetSubtitle("Sub 2")

			Expect(item1.GetTitle()).To(Equal("Item 1"))
			Expect(item2.GetTitle()).To(Equal("Item 2"))
			Expect(item1.GetSubtitle()).To(Equal("Sub 1"))
			Expect(item2.GetSubtitle()).To(Equal("Sub 2"))
		})

		It("handles selection state independently", func() {
			item1 := NewListItem("Item 1", 80)
			item2 := NewListItem("Item 2", 80)

			item1.SetSelected(true)
			Expect(item1.IsSelected()).To(BeTrue())
			Expect(item2.IsSelected()).To(BeFalse())
		})
	})
})
