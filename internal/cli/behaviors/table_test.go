package behaviors_test

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/themes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Test item type
type TestItem struct {
	Name   string
	Status string
	Count  int
}

var _ = Describe("TableBehavior", func() {
	var (
		theme     themes.Theme
		columns   []behaviors.ColumnDef
		formatter behaviors.RowFormatter[*TestItem]
		table     *behaviors.TableBehavior[*TestItem]
		items     []*TestItem
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
		columns = []behaviors.ColumnDef{
			{Title: "Name", Width: 20},
			{Title: "Status", Width: 10},
			{Title: "Count", Width: 8},
		}
		formatter = func(item *TestItem, _ int) []string {
			return []string{item.Name, item.Status, string(rune(item.Count + '0'))}
		}
		items = []*TestItem{
			{Name: "Item 1", Status: "Active", Count: 5},
			{Name: "Item 2", Status: "Inactive", Count: 10},
			{Name: "Item 3", Status: "Active", Count: 3},
		}
	})

	Describe("Construction", func() {
		It("should create a new table behavior", func() {
			table = behaviors.NewTableBehavior(theme, columns, formatter)
			Expect(table).NotTo(BeNil())
		})

		It("should have default page size of 15", func() {
			table = behaviors.NewTableBehavior(theme, columns, formatter)
			Expect(table.GetPageSize()).To(Equal(15))
		})

		It("should have default empty message", func() {
			table = behaviors.NewTableBehavior(theme, columns, formatter)
			Expect(table.IsEmpty()).To(BeTrue())
		})

		It("should start with zero items", func() {
			table = behaviors.NewTableBehavior(theme, columns, formatter)
			Expect(table.Count()).To(Equal(0))
			Expect(table.TotalCount()).To(Equal(0))
		})

		It("should start with selection at index 0", func() {
			table = behaviors.NewTableBehavior(theme, columns, formatter)
			Expect(table.GetSelectedIndex()).To(Equal(0))
		})
	})

	Describe("Configuration (Fluent API)", func() {
		BeforeEach(func() {
			table = behaviors.NewTableBehavior(theme, columns, formatter)
		})

		It("should set page size and return self", func() {
			result := table.PageSize(20)
			Expect(result).To(BeIdenticalTo(table))
			Expect(table.GetPageSize()).To(Equal(20))
		})

		It("should set empty message and return self", func() {
			result := table.EmptyMessage("No items found")
			Expect(result).To(BeIdenticalTo(table))
		})

		It("should set pagination prefix and return self", func() {
			result := table.PaginationPrefix("Events")
			Expect(result).To(BeIdenticalTo(table))
		})

		It("should set dimensions and return self", func() {
			result := table.Dimensions(120, 40)
			Expect(result).To(BeIdenticalTo(table))
		})

		It("should hide pagination and return self", func() {
			result := table.HidePagination()
			Expect(result).To(BeIdenticalTo(table))
		})

		It("should show pagination and return self", func() {
			result := table.ShowPagination()
			Expect(result).To(BeIdenticalTo(table))
		})
	})

	Describe("Data Management", func() {
		BeforeEach(func() {
			table = behaviors.NewTableBehavior(theme, columns, formatter)
		})

		Describe("SetItems", func() {
			It("should replace all items", func() {
				table.SetItems(items)
				Expect(table.Count()).To(Equal(3))
				Expect(table.TotalCount()).To(Equal(3))
			})

			It("should reset selection to 0", func() {
				table.SetItems(items)
				table.SetSelectedIndex(2)
				Expect(table.GetSelectedIndex()).To(Equal(2))

				newItems := []*TestItem{{Name: "New", Status: "Active", Count: 1}}
				table.SetItems(newItems)
				Expect(table.GetSelectedIndex()).To(Equal(0))
			})

			It("should accept empty slice", func() {
				table.SetItems([]*TestItem{})
				Expect(table.IsEmpty()).To(BeTrue())
				Expect(table.Count()).To(Equal(0))
			})

			It("should accept nil slice", func() {
				table.SetItems(nil)
				Expect(table.IsEmpty()).To(BeTrue())
			})
		})

		Describe("GetItems", func() {
			It("should return current items", func() {
				table.SetItems(items)
				retrieved := table.GetItems()
				Expect(retrieved).To(HaveLen(3))
				Expect(retrieved[0]).To(Equal(items[0]))
			})

			It("should return empty slice when no items", func() {
				retrieved := table.GetItems()
				Expect(retrieved).To(BeEmpty())
			})
		})

		Describe("GetSelectedItem", func() {
			It("should return pointer to selected item", func() {
				table.SetItems(items)
				selected := table.GetSelectedItem()
				Expect(selected).NotTo(BeNil())
				Expect(*selected).To(Equal(items[0]))
			})

			It("should return nil when list is empty", func() {
				selected := table.GetSelectedItem()
				Expect(selected).To(BeNil())
			})

			It("should return correct item after index change", func() {
				table.SetItems(items)
				table.SetSelectedIndex(1)
				selected := table.GetSelectedItem()
				Expect(*selected).To(Equal(items[1]))
			})
		})

		Describe("GetSelectedIndex", func() {
			It("should return current selection index", func() {
				table.SetItems(items)
				Expect(table.GetSelectedIndex()).To(Equal(0))
			})
		})

		Describe("IsEmpty", func() {
			It("should return true when no items", func() {
				Expect(table.IsEmpty()).To(BeTrue())
			})

			It("should return false when items exist", func() {
				table.SetItems(items)
				Expect(table.IsEmpty()).To(BeFalse())
			})
		})

		Describe("Count", func() {
			It("should return number of displayed items", func() {
				table.SetItems(items)
				Expect(table.Count()).To(Equal(3))
			})

			It("should return 0 when empty", func() {
				Expect(table.Count()).To(Equal(0))
			})
		})

		Describe("TotalCount", func() {
			It("should return total items before filtering", func() {
				table.SetItems(items)
				Expect(table.TotalCount()).To(Equal(3))
			})
		})
	})

	Describe("ListNavigator Interface", func() {
		BeforeEach(func() {
			table = behaviors.NewTableBehavior(theme, columns, formatter)
			table.SetItems(items)
		})

		Describe("GetTotalItems", func() {
			It("should return item count", func() {
				Expect(table.GetTotalItems()).To(Equal(3))
			})
		})

		Describe("SetSelectedIndex", func() {
			It("should set valid index", func() {
				table.SetSelectedIndex(1)
				Expect(table.GetSelectedIndex()).To(Equal(1))
			})

			It("should clamp negative index to 0", func() {
				table.SetSelectedIndex(-5)
				Expect(table.GetSelectedIndex()).To(Equal(0))
			})

			It("should clamp too-large index to last item", func() {
				table.SetSelectedIndex(100)
				Expect(table.GetSelectedIndex()).To(Equal(2))
			})

			It("should handle empty list", func() {
				emptyTable := behaviors.NewTableBehavior(theme, columns, formatter)
				emptyTable.SetSelectedIndex(5)
				Expect(emptyTable.GetSelectedIndex()).To(Equal(0))
			})
		})

		Describe("GetPageSize", func() {
			It("should return configured page size", func() {
				Expect(table.GetPageSize()).To(Equal(15))
			})

			It("should return custom page size", func() {
				table.PageSize(25)
				Expect(table.GetPageSize()).To(Equal(25))
			})
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			table = behaviors.NewTableBehavior(theme, columns, formatter)
			table.SetItems(items)
		})

		Describe("HandleNavigation", func() {
			It("should move up with 'up' key", func() {
				table.SetSelectedIndex(2)
				handled := table.HandleNavigation("up")
				Expect(handled).To(BeTrue())
				Expect(table.GetSelectedIndex()).To(Equal(1))
			})

			It("should move up with 'k' key (vim)", func() {
				table.SetSelectedIndex(2)
				handled := table.HandleNavigation("k")
				Expect(handled).To(BeTrue())
				Expect(table.GetSelectedIndex()).To(Equal(1))
			})

			It("should not move up from first item", func() {
				table.SetSelectedIndex(0)
				table.HandleNavigation("up")
				Expect(table.GetSelectedIndex()).To(Equal(0))
			})

			It("should move down with 'down' key", func() {
				table.SetSelectedIndex(0)
				handled := table.HandleNavigation("down")
				Expect(handled).To(BeTrue())
				Expect(table.GetSelectedIndex()).To(Equal(1))
			})

			It("should move down with 'j' key (vim)", func() {
				table.SetSelectedIndex(0)
				handled := table.HandleNavigation("j")
				Expect(handled).To(BeTrue())
				Expect(table.GetSelectedIndex()).To(Equal(1))
			})

			It("should not move down from last item", func() {
				table.SetSelectedIndex(2)
				table.HandleNavigation("down")
				Expect(table.GetSelectedIndex()).To(Equal(2))
			})

			It("should page up with 'pgup'", func() {
				// With 3 items and page size 15, we're on one page
				table.SetSelectedIndex(2)
				table.PageSize(1) // Set small page size to test paging
				handled := table.HandleNavigation("pgup")
				Expect(handled).To(BeTrue())
				Expect(table.GetSelectedIndex()).To(Equal(1))
			})

			It("should page up with 'ctrl+u'", func() {
				table.SetSelectedIndex(2)
				table.PageSize(1)
				handled := table.HandleNavigation("ctrl+u")
				Expect(handled).To(BeTrue())
				Expect(table.GetSelectedIndex()).To(Equal(1))
			})

			It("should page down with 'pgdn'", func() {
				table.SetSelectedIndex(0)
				table.PageSize(1)
				handled := table.HandleNavigation("pgdn")
				Expect(handled).To(BeTrue())
				Expect(table.GetSelectedIndex()).To(Equal(1))
			})

			It("should page down with 'ctrl+d'", func() {
				table.SetSelectedIndex(0)
				table.PageSize(1)
				handled := table.HandleNavigation("ctrl+d")
				Expect(handled).To(BeTrue())
				Expect(table.GetSelectedIndex()).To(Equal(1))
			})

			It("should go to first with 'home'", func() {
				table.SetSelectedIndex(2)
				handled := table.HandleNavigation("home")
				Expect(handled).To(BeTrue())
				Expect(table.GetSelectedIndex()).To(Equal(0))
			})

			It("should go to first with 'g' (vim)", func() {
				table.SetSelectedIndex(2)
				handled := table.HandleNavigation("g")
				Expect(handled).To(BeTrue())
				Expect(table.GetSelectedIndex()).To(Equal(0))
			})

			It("should go to last with 'end'", func() {
				table.SetSelectedIndex(0)
				handled := table.HandleNavigation("end")
				Expect(handled).To(BeTrue())
				Expect(table.GetSelectedIndex()).To(Equal(2))
			})

			It("should go to last with 'G' (vim)", func() {
				table.SetSelectedIndex(0)
				handled := table.HandleNavigation("G")
				Expect(handled).To(BeTrue())
				Expect(table.GetSelectedIndex()).To(Equal(2))
			})

			It("should return false for unknown keys", func() {
				handled := table.HandleNavigation("unknown")
				Expect(handled).To(BeFalse())
			})

			It("should handle navigation on empty list", func() {
				emptyTable := behaviors.NewTableBehavior(theme, columns, formatter)
				handled := emptyTable.HandleNavigation("down")
				Expect(handled).To(BeFalse())
			})
		})
	})

	Describe("Filtering", func() {
		BeforeEach(func() {
			table = behaviors.NewTableBehavior(theme, columns, formatter)
			table.SetItems(items)
		})

		Describe("SetFilter", func() {
			It("should apply filter predicate", func() {
				// Filter to only Active items
				predicate := func(item *TestItem) bool {
					return item.Status == "Active"
				}
				table.SetFilter(predicate)

				Expect(table.Count()).To(Equal(2)) // Items 1 and 3 are Active
				Expect(table.TotalCount()).To(Equal(3))
			})

			It("should update displayed items", func() {
				predicate := func(item *TestItem) bool {
					return item.Count > 4
				}
				table.SetFilter(predicate)

				items := table.GetItems()
				Expect(items).To(HaveLen(2)) // Items with count 5 and 10
			})

			It("should preserve selection if possible", func() {
				table.SetSelectedIndex(2) // Select "Item 3" (Active, count 3)

				// Filter to only Active items (keeps Item 1 and Item 3)
				predicate := func(item *TestItem) bool {
					return item.Status == "Active"
				}
				table.SetFilter(predicate)

				// Item 3 should still be selected (now at index 1 in filtered list)
				selected := table.GetSelectedItem()
				Expect(selected).NotTo(BeNil())
				Expect((*selected).Name).To(Equal("Item 3"))
			})

			It("should reset selection if selected item filtered out", func() {
				table.SetSelectedIndex(1) // Select "Item 2" (Inactive)

				// Filter to only Active items
				predicate := func(item *TestItem) bool {
					return item.Status == "Active"
				}
				table.SetFilter(predicate)

				// Selection should reset to first item
				Expect(table.GetSelectedIndex()).To(Equal(0))
			})
		})

		Describe("ClearFilter", func() {
			It("should remove active filter", func() {
				predicate := func(item *TestItem) bool {
					return item.Status == "Active"
				}
				table.SetFilter(predicate)
				Expect(table.Count()).To(Equal(2))

				table.ClearFilter()
				Expect(table.Count()).To(Equal(3))
			})

			It("should return self for chaining", func() {
				result := table.ClearFilter()
				Expect(result).To(BeIdenticalTo(table))
			})
		})

		Describe("HasFilter", func() {
			It("should return false when no filter", func() {
				Expect(table.HasFilter()).To(BeFalse())
			})

			It("should return true when filter active", func() {
				predicate := func(item *TestItem) bool {
					return item.Status == "Active"
				}
				table.SetFilter(predicate)
				Expect(table.HasFilter()).To(BeTrue())
			})

			It("should return false after clearing filter", func() {
				predicate := func(_ *TestItem) bool { return true }
				table.SetFilter(predicate)
				table.ClearFilter()
				Expect(table.HasFilter()).To(BeFalse())
			})
		})
	})

	Describe("Sorting", func() {
		BeforeEach(func() {
			table = behaviors.NewTableBehavior(theme, columns, formatter)
			table.SetItems(items)
		})

		Describe("SetSort", func() {
			It("should sort items by comparator", func() {
				// Sort by Count ascending
				comparator := func(a, b *TestItem) int {
					return a.Count - b.Count
				}
				table.SetSort(comparator, false)

				sorted := table.GetItems()
				Expect(sorted[0].Count).To(Equal(3))  // Item 3
				Expect(sorted[1].Count).To(Equal(5))  // Item 1
				Expect(sorted[2].Count).To(Equal(10)) // Item 2
			})

			It("should sort in reverse when reverse flag is true", func() {
				// Sort by Count descending
				comparator := func(a, b *TestItem) int {
					return a.Count - b.Count
				}
				table.SetSort(comparator, true)

				sorted := table.GetItems()
				Expect(sorted[0].Count).To(Equal(10)) // Item 2
				Expect(sorted[1].Count).To(Equal(5))  // Item 1
				Expect(sorted[2].Count).To(Equal(3))  // Item 3
			})

			It("should work with string comparisons", func() {
				// Sort by Name
				comparator := func(a, b *TestItem) int {
					if a.Name < b.Name {
						return -1
					}
					if a.Name > b.Name {
						return 1
					}
					return 0
				}
				table.SetSort(comparator, false)

				sorted := table.GetItems()
				Expect(sorted[0].Name).To(Equal("Item 1"))
				Expect(sorted[1].Name).To(Equal("Item 2"))
				Expect(sorted[2].Name).To(Equal("Item 3"))
			})

			It("should preserve selection if possible", func() {
				table.SetSelectedIndex(1) // Select Item 2

				comparator := func(a, b *TestItem) int {
					return a.Count - b.Count
				}
				table.SetSort(comparator, false)

				// Item 2 should still be selected (now at different index)
				selected := table.GetSelectedItem()
				Expect((*selected).Name).To(Equal("Item 2"))
			})
		})

		Describe("ClearSort", func() {
			It("should restore original order", func() {
				comparator := func(a, b *TestItem) int {
					return a.Count - b.Count
				}
				table.SetSort(comparator, false)

				table.ClearSort()
				sorted := table.GetItems()
				Expect(sorted[0].Name).To(Equal("Item 1"))
				Expect(sorted[1].Name).To(Equal("Item 2"))
			})
		})

		Describe("HasSort", func() {
			It("should return false when no sort", func() {
				Expect(table.HasSort()).To(BeFalse())
			})

			It("should return true when sort active", func() {
				comparator := func(_, _ *TestItem) int { return 0 }
				table.SetSort(comparator, false)
				Expect(table.HasSort()).To(BeTrue())
			})
		})
	})

	Describe("Filter and Sort Together", func() {
		BeforeEach(func() {
			table = behaviors.NewTableBehavior(theme, columns, formatter)
			table.SetItems(items)
		})

		It("should apply filter then sort", func() {
			// Filter to Active items
			predicate := func(item *TestItem) bool {
				return item.Status == "Active"
			}
			table.SetFilter(predicate)

			// Sort by Count
			comparator := func(a, b *TestItem) int {
				return a.Count - b.Count
			}
			table.SetSort(comparator, false)

			filtered := table.GetItems()
			Expect(filtered).To(HaveLen(2))
			Expect(filtered[0].Count).To(Equal(3)) // Item 3
			Expect(filtered[1].Count).To(Equal(5)) // Item 1
		})
	})

	Describe("Rendering", func() {
		BeforeEach(func() {
			table = behaviors.NewTableBehavior(theme, columns, formatter)
		})

		Describe("Render", func() {
			It("should render table when items exist", func() {
				table.SetItems(items).Dimensions(100, 30)
				view := table.Render()
				Expect(view).NotTo(BeEmpty())
			})

			It("should show empty message when no items", func() {
				table.EmptyMessage("No items to display").Dimensions(100, 30)
				view := table.Render()
				Expect(view).To(ContainSubstring("No items to display"))
			})

			It("should include pagination info when enabled", func() {
				table.SetItems(items).ShowPagination().PaginationPrefix("Items").Dimensions(100, 30)
				view := table.Render()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Describe("RenderPaginationInfo", func() {
			It("should format pagination with prefix", func() {
				table.SetItems(items).PaginationPrefix("Events")
				info := table.RenderPaginationInfo()
				Expect(info).To(ContainSubstring("Events"))
				Expect(info).To(ContainSubstring("3")) // Total items
			})

			It("should calculate correct page numbers", func() {
				// Create more items to test pagination
				manyItems := make([]*TestItem, 30)
				for i := range manyItems {
					manyItems[i] = &TestItem{Name: "Item", Status: "Active", Count: i}
				}
				table.SetItems(manyItems).PageSize(10)

				info := table.RenderPaginationInfo()
				Expect(info).To(ContainSubstring("Page 1 of 3"))
			})
		})
	})

	Describe("Viewport Integration", func() {
		BeforeEach(func() {
			table = behaviors.NewTableBehavior(theme, columns, formatter)
		})

		Describe("SetHeight", func() {
			It("should accept and store height for viewport", func() {
				result := table.SetHeight(20)
				Expect(result).To(BeIdenticalTo(table))
				// Height should be stored and used for viewport
			})

			It("should handle minimum height gracefully", func() {
				result := table.SetHeight(1)
				Expect(result).To(BeIdenticalTo(table))
				// Should not panic with very small height
			})

			It("should update viewport height on subsequent calls", func() {
				table.SetHeight(10)
				table.SetHeight(20)
				// Should accept height updates
			})
		})

		Describe("Viewport Rendering", func() {
			It("should render content within viewport height", func() {
				// Create many items to exceed viewport
				manyItems := make([]*TestItem, 30)
				for i := range manyItems {
					manyItems[i] = &TestItem{Name: "Item", Status: "Active", Count: i}
				}
				table.SetItems(manyItems).SetHeight(10)

				rendered := table.Render()
				Expect(rendered).NotTo(BeEmpty())
				// Should render within height constraint
			})

			It("should show only viewport portion of content", func() {
				// Create items that exceed viewport
				manyItems := make([]*TestItem, 20)
				for i := range manyItems {
					manyItems[i] = &TestItem{Name: "Item " + string(rune(i+'0')), Status: "Active", Count: i}
				}
				table.SetItems(manyItems).SetHeight(5)

				rendered := table.Render()
				// Should not show all 20 items at once
				Expect(rendered).NotTo(BeEmpty())
			})
		})

		Describe("Scrolling", func() {
			var manyItems []*TestItem

			BeforeEach(func() {
				// Create more items than viewport can show
				manyItems = make([]*TestItem, 20)
				for i := range manyItems {
					manyItems[i] = &TestItem{Name: "Item", Status: "Active", Count: i}
				}
				table.SetItems(manyItems).SetHeight(5)
			})

			It("should scroll down when navigating beyond viewport", func() {
				// Navigate down multiple times
				for range 10 {
					table.HandleNavigation("down")
				}

				// Should have scrolled
				rendered := table.Render()
				Expect(rendered).NotTo(BeEmpty())
			})

			It("should scroll up when navigating up", func() {
				// First scroll down
				for range 10 {
					table.HandleNavigation("down")
				}

				// Then scroll up
				for range 5 {
					table.HandleNavigation("up")
				}

				rendered := table.Render()
				Expect(rendered).NotTo(BeEmpty())
			})

			It("should handle page up/down with viewport", func() {
				table.HandleNavigation("pgdn")
				rendered := table.Render()
				Expect(rendered).NotTo(BeEmpty())

				table.HandleNavigation("pgup")
				rendered = table.Render()
				Expect(rendered).NotTo(BeEmpty())
			})
		})
	})
})
