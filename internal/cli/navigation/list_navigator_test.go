package navigation_test

import (
	"github.com/baphled/kariya/internal/cli/navigation"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Mock implementation of ListNavigator for testing.
type mockListNavigator struct {
	totalItems    int
	selectedIndex int
	pageSize      int
}

func newMockListNavigator(totalItems, pageSize int) *mockListNavigator {
	return &mockListNavigator{
		totalItems:    totalItems,
		selectedIndex: 0,
		pageSize:      pageSize,
	}
}

func (m *mockListNavigator) GetTotalItems() int {
	return m.totalItems
}

func (m *mockListNavigator) GetSelectedIndex() int {
	return m.selectedIndex
}

func (m *mockListNavigator) SetSelectedIndex(idx int) {
	if idx < 0 {
		idx = 0
	}
	if idx >= m.totalItems {
		idx = m.totalItems - 1
	}
	m.selectedIndex = idx
}

func (m *mockListNavigator) GetPageSize() int {
	return m.pageSize
}

var _ = Describe("ListNavigationHandler", func() {
	var (
		handler   *navigation.ListNavigationHandler
		navigator *mockListNavigator
	)

	Context("with 10 items and page size of 5", func() {
		BeforeEach(func() {
			navigator = newMockListNavigator(10, 5)
			handler = navigation.NewListNavigationHandler(navigator)
		})

		Describe("HandleKey", func() {
			Context("basic navigation", func() {
				It("should move down by 1 with 'down' key", func() {
					handled := handler.HandleKey("down")
					Expect(handled).To(BeTrue())
					Expect(navigator.GetSelectedIndex()).To(Equal(1))
				})

				It("should move down by 1 with 'j' key", func() {
					handled := handler.HandleKey("j")
					Expect(handled).To(BeTrue())
					Expect(navigator.GetSelectedIndex()).To(Equal(1))
				})

				It("should move up by 1 with 'up' key", func() {
					navigator.SetSelectedIndex(5)
					handled := handler.HandleKey("up")
					Expect(handled).To(BeTrue())
					Expect(navigator.GetSelectedIndex()).To(Equal(4))
				})

				It("should move up by 1 with 'k' key", func() {
					navigator.SetSelectedIndex(5)
					handled := handler.HandleKey("k")
					Expect(handled).To(BeTrue())
					Expect(navigator.GetSelectedIndex()).To(Equal(4))
				})
			})

			Context("boundary conditions", func() {
				It("should stay at 0 when moving up from first item", func() {
					navigator.SetSelectedIndex(0)
					handled := handler.HandleKey("up")
					Expect(handled).To(BeTrue())
					Expect(navigator.GetSelectedIndex()).To(Equal(0))
				})

				It("should stay at last item when moving down from last item", func() {
					navigator.SetSelectedIndex(9) // Last item
					handled := handler.HandleKey("down")
					Expect(handled).To(BeTrue())
					Expect(navigator.GetSelectedIndex()).To(Equal(9))
				})
			})

			Context("page navigation", func() {
				It("should move down by page size with 'pgdn' key", func() {
					navigator.SetSelectedIndex(0)
					handled := handler.HandleKey("pgdn")
					Expect(handled).To(BeTrue())
					Expect(navigator.GetSelectedIndex()).To(Equal(5))
				})

				It("should move down by page size with 'ctrl+d' key", func() {
					navigator.SetSelectedIndex(0)
					handled := handler.HandleKey("ctrl+d")
					Expect(handled).To(BeTrue())
					Expect(navigator.GetSelectedIndex()).To(Equal(5))
				})

				It("should move up by page size with 'pgup' key", func() {
					navigator.SetSelectedIndex(7)
					handled := handler.HandleKey("pgup")
					Expect(handled).To(BeTrue())
					Expect(navigator.GetSelectedIndex()).To(Equal(2))
				})

				It("should move up by page size with 'ctrl+u' key", func() {
					navigator.SetSelectedIndex(7)
					handled := handler.HandleKey("ctrl+u")
					Expect(handled).To(BeTrue())
					Expect(navigator.GetSelectedIndex()).To(Equal(2))
				})

				It("should stop at last item when paging down beyond end", func() {
					navigator.SetSelectedIndex(7)
					handled := handler.HandleKey("pgdn")
					Expect(handled).To(BeTrue())
					Expect(navigator.GetSelectedIndex()).To(Equal(9))
				})

				It("should stop at first item when paging up beyond start", func() {
					navigator.SetSelectedIndex(2)
					handled := handler.HandleKey("pgup")
					Expect(handled).To(BeTrue())
					Expect(navigator.GetSelectedIndex()).To(Equal(0))
				})
			})

			Context("jump navigation", func() {
				It("should go to first item with 'home' key", func() {
					navigator.SetSelectedIndex(5)
					handled := handler.HandleKey("home")
					Expect(handled).To(BeTrue())
					Expect(navigator.GetSelectedIndex()).To(Equal(0))
				})

				It("should go to first item with 'g' key", func() {
					navigator.SetSelectedIndex(5)
					handled := handler.HandleKey("g")
					Expect(handled).To(BeTrue())
					Expect(navigator.GetSelectedIndex()).To(Equal(0))
				})

				It("should go to last item with 'end' key", func() {
					navigator.SetSelectedIndex(0)
					handled := handler.HandleKey("end")
					Expect(handled).To(BeTrue())
					Expect(navigator.GetSelectedIndex()).To(Equal(9))
				})

				It("should go to last item with 'G' key", func() {
					navigator.SetSelectedIndex(0)
					handled := handler.HandleKey("G")
					Expect(handled).To(BeTrue())
					Expect(navigator.GetSelectedIndex()).To(Equal(9))
				})
			})

			Context("unhandled keys", func() {
				It("should return false for unhandled keys", func() {
					handled := handler.HandleKey("x")
					Expect(handled).To(BeFalse())
					Expect(navigator.GetSelectedIndex()).To(Equal(0))
				})

				It("should return false for 'enter' key", func() {
					handled := handler.HandleKey("enter")
					Expect(handled).To(BeFalse())
				})

				It("should return false for 'esc' key", func() {
					handled := handler.HandleKey("esc")
					Expect(handled).To(BeFalse())
				})
			})
		})

		Describe("FormatRowText", func() {
			It("should add indicator to selected row", func() {
				navigator.SetSelectedIndex(3)
				text := handler.FormatRowText(3, "Item 3")
				Expect(text).To(Equal("▶ Item 3"))
			})

			It("should add spaces to non-selected rows", func() {
				navigator.SetSelectedIndex(3)
				text := handler.FormatRowText(2, "Item 2")
				Expect(text).To(Equal("  Item 2"))
			})

			It("should update indicator when selection changes", func() {
				navigator.SetSelectedIndex(0)
				text0 := handler.FormatRowText(0, "Item 0")
				Expect(text0).To(Equal("▶ Item 0"))

				handler.HandleKey("down")
				text0After := handler.FormatRowText(0, "Item 0")
				text1After := handler.FormatRowText(1, "Item 1")
				Expect(text0After).To(Equal("  Item 0"))
				Expect(text1After).To(Equal("▶ Item 1"))
			})
		})
	})

	Context("with empty list", func() {
		BeforeEach(func() {
			navigator = newMockListNavigator(0, 5)
			handler = navigation.NewListNavigationHandler(navigator)
		})

		It("should return false for any navigation key", func() {
			handled := handler.HandleKey("down")
			Expect(handled).To(BeFalse())

			handled = handler.HandleKey("up")
			Expect(handled).To(BeFalse())
		})
	})

	Context("with single item", func() {
		BeforeEach(func() {
			navigator = newMockListNavigator(1, 5)
			handler = navigation.NewListNavigationHandler(navigator)
		})

		It("should stay at 0 when trying to move down", func() {
			handled := handler.HandleKey("down")
			Expect(handled).To(BeTrue())
			Expect(navigator.GetSelectedIndex()).To(Equal(0))
		})

		It("should stay at 0 when trying to move up", func() {
			handled := handler.HandleKey("up")
			Expect(handled).To(BeTrue())
			Expect(navigator.GetSelectedIndex()).To(Equal(0))
		})
	})
})
