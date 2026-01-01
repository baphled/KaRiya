package navigation

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Navigation Constants", func() {
	Describe("NavigationKey type", func() {
		It("should define KeyBack as Escape", func() {
			Expect(KeyBack).To(Equal(NavigationKey("Esc")))
		})

		It("should define KeyUp as up arrow/k", func() {
			Expect(KeyUp).To(Equal(NavigationKey("↑/k")))
		})

		It("should define KeyDown as down arrow/j", func() {
			Expect(KeyDown).To(Equal(NavigationKey("↓/j")))
		})

		It("should define KeyLeft as left arrow/h", func() {
			Expect(KeyLeft).To(Equal(NavigationKey("←/h")))
		})

		It("should define KeyRight as right arrow/l", func() {
			Expect(KeyRight).To(Equal(NavigationKey("→/l")))
		})

		It("should define KeySelect as Enter", func() {
			Expect(KeySelect).To(Equal(NavigationKey("Enter")))
		})

		It("should define KeyToggle as Space", func() {
			Expect(KeyToggle).To(Equal(NavigationKey("Space")))
		})

		It("should define KeyFilter as f", func() {
			Expect(KeyFilter).To(Equal(NavigationKey("f")))
		})

		It("should define KeySort as s", func() {
			Expect(KeySort).To(Equal(NavigationKey("s")))
		})

		It("should define KeySearch as /", func() {
			Expect(KeySearch).To(Equal(NavigationKey("/")))
		})

		It("should define KeyEdit as e", func() {
			Expect(KeyEdit).To(Equal(NavigationKey("e")))
		})

		It("should define KeyDelete as d", func() {
			Expect(KeyDelete).To(Equal(NavigationKey("d")))
		})

		It("should define KeyHelp as ?", func() {
			Expect(KeyHelp).To(Equal(NavigationKey("?")))
		})

		It("should define KeyHome as h", func() {
			Expect(KeyHome).To(Equal(NavigationKey("h")))
		})

		It("should define KeyQuit as q", func() {
			Expect(KeyQuit).To(Equal(NavigationKey("q")))
		})

		It("should define KeyBulk as b", func() {
			Expect(KeyBulk).To(Equal(NavigationKey("b")))
		})

		It("should define KeyCapture as c", func() {
			Expect(KeyCapture).To(Equal(NavigationKey("c")))
		})

		It("should define KeyList as l", func() {
			Expect(KeyList).To(Equal(NavigationKey("l")))
		})

		It("should define KeyMetadata as m", func() {
			Expect(KeyMetadata).To(Equal(NavigationKey("m")))
		})
	})

	Describe("AllNavigationKeys", func() {
		It("should return all 20 navigation keys", func() {
			keys := AllNavigationKeys()
			Expect(keys).To(HaveLen(20))
		})

		It("should include all primary navigation keys", func() {
			keys := AllNavigationKeys()
			Expect(keys).To(ContainElement(KeyBack))
			Expect(keys).To(ContainElement(KeyUp))
			Expect(keys).To(ContainElement(KeyDown))
			Expect(keys).To(ContainElement(KeyLeft))
			Expect(keys).To(ContainElement(KeyRight))
			Expect(keys).To(ContainElement(KeySelect))
			Expect(keys).To(ContainElement(KeyToggle))
		})

		It("should include all action keys", func() {
			keys := AllNavigationKeys()
			Expect(keys).To(ContainElement(KeyFilter))
			Expect(keys).To(ContainElement(KeySort))
			Expect(keys).To(ContainElement(KeySearch))
			Expect(keys).To(ContainElement(KeyEdit))
			Expect(keys).To(ContainElement(KeyDelete))
			Expect(keys).To(ContainElement(KeyHelp))
			Expect(keys).To(ContainElement(KeyHome))
			Expect(keys).To(ContainElement(KeyQuit))
			Expect(keys).To(ContainElement(KeyBulk))
			Expect(keys).To(ContainElement(KeyCapture))
			Expect(keys).To(ContainElement(KeyList))
			Expect(keys).To(ContainElement(KeyMetadata))
		})

		It("should contain no duplicates", func() {
			keys := AllNavigationKeys()
			uniqueKeys := make(map[NavigationKey]bool)
			for _, key := range keys {
				Expect(uniqueKeys[key]).To(BeFalse(), "duplicate key found: %s", key)
				uniqueKeys[key] = true
			}
		})
	})

	Describe("KeyDescription", func() {
		It("should have description for every navigation key", func() {
			keys := AllNavigationKeys()
			for _, key := range keys {
				description, exists := KeyDescription[key]
				Expect(exists).To(BeTrue(), "missing description for key: %s", key)
				Expect(description).NotTo(BeEmpty(), "empty description for key: %s", key)
			}
		})

		It("should provide description for KeyBack", func() {
			description := KeyDescription[KeyBack]
			Expect(description).To(Equal("Go back to previous screen"))
		})

		It("should provide description for KeySelect", func() {
			description := KeyDescription[KeySelect]
			Expect(description).To(Equal("Confirm selection or submit"))
		})

		It("should provide description for KeyHelp", func() {
			description := KeyDescription[KeyHelp]
			Expect(description).To(Equal("Show help information"))
		})

		It("should have exactly 19 descriptions matching 19 keys", func() {
			keys := AllNavigationKeys()
			Expect(len(KeyDescription)).To(Equal(len(keys)))
		})
	})
})
