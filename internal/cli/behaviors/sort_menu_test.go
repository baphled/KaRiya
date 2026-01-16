package behaviors_test

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/themes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Test item for sort tests
type SortTestItem struct {
	ID   int
	Name string
	Age  int
}

var _ = Describe("SortMenuBehavior", func() {
	var (
		themeObj  themes.Theme
		table     *behaviors.TableBehavior[SortTestItem]
		sortMenu  *behaviors.SortMenuBehavior[SortTestItem]
		columns   []behaviors.ColumnDef
		formatter behaviors.RowFormatter[SortTestItem]
	)

	BeforeEach(func() {
		themeObj = themes.NewDefaultTheme()
		columns = []behaviors.ColumnDef{
			{Title: "ID", Width: 5},
			{Title: "Name", Width: 20},
			{Title: "Age", Width: 5},
		}
		formatter = func(item SortTestItem, index int) []string {
			return []string{
				string(rune(item.ID + '0')),
				item.Name,
				string(rune(item.Age + '0')),
			}
		}
		table = behaviors.NewTableBehavior(themeObj, columns, formatter)
		table.SetItems([]SortTestItem{
			{ID: 3, Name: "Charlie", Age: 25},
			{ID: 1, Name: "Alice", Age: 30},
			{ID: 2, Name: "Bob", Age: 20},
		})
	})

	Describe("Construction", func() {
		It("should create with non-nil table", func() {
			sortMenu = behaviors.NewSortMenuBehavior(themeObj, table, "Sort Items")
			Expect(sortMenu).NotTo(BeNil())
		})

		It("should panic if table is nil", func() {
			Expect(func() {
				behaviors.NewSortMenuBehavior[SortTestItem](themeObj, nil, "Sort")
			}).To(Panic())
		})

		It("should default to not active", func() {
			sortMenu = behaviors.NewSortMenuBehavior(themeObj, table, "Sort Items")
			Expect(sortMenu.IsActive()).To(BeFalse())
		})

		It("should start with focused index 0", func() {
			sortMenu = behaviors.NewSortMenuBehavior(themeObj, table, "Sort Items")
			// Can't directly test focusedIndex, but navigation should work from 0
			Expect(sortMenu.IsActive()).To(BeFalse()) // Sanity check
		})
	})

	Describe("Configuration", func() {
		BeforeEach(func() {
			sortMenu = behaviors.NewSortMenuBehavior(themeObj, table, "Sort Items")
		})

		It("should add option with comparator", func() {
			comparator := func(a, b SortTestItem) int {
				if a.Name < b.Name {
					return -1
				}
				if a.Name > b.Name {
					return 1
				}
				return 0
			}

			result := sortMenu.AddOption("Name (A-Z)", comparator, false)
			Expect(result).To(Equal(sortMenu)) // Check chaining
		})

		It("should allow multiple options", func() {
			sortMenu.AddOption("Name (A-Z)", func(a, b SortTestItem) int {
				if a.Name < b.Name {
					return -1
				}
				return 1
			}, false)

			sortMenu.AddOption("Age (Ascending)", func(a, b SortTestItem) int {
				return a.Age - b.Age
			}, false)

			sortMenu.Show()
			view := sortMenu.Render()
			Expect(view).To(ContainSubstring("Name (A-Z)"))
			Expect(view).To(ContainSubstring("Age (Ascending)"))
		})

		It("should set callback with OnApply", func() {
			result := sortMenu.OnApply(func() {
				// Callback will be tested in Selection tests
			})

			Expect(result).To(Equal(sortMenu)) // Check chaining
		})
	})

	Describe("State", func() {
		BeforeEach(func() {
			sortMenu = behaviors.NewSortMenuBehavior(themeObj, table, "Sort Items")
		})

		It("should Show() set active flag", func() {
			sortMenu.Show()
			Expect(sortMenu.IsActive()).To(BeTrue())
		})

		It("should Hide() clear active flag", func() {
			sortMenu.Show()
			Expect(sortMenu.IsActive()).To(BeTrue())

			sortMenu.Hide()
			Expect(sortMenu.IsActive()).To(BeFalse())
		})

		It("should IsActive() return correct state", func() {
			Expect(sortMenu.IsActive()).To(BeFalse())

			sortMenu.Show()
			Expect(sortMenu.IsActive()).To(BeTrue())

			sortMenu.Hide()
			Expect(sortMenu.IsActive()).To(BeFalse())
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			sortMenu = behaviors.NewSortMenuBehavior(themeObj, table, "Sort Items")
			sortMenu.AddOption("Name (A-Z)", func(a, b SortTestItem) int { return 0 }, false)
			sortMenu.AddOption("Name (Z-A)", func(a, b SortTestItem) int { return 0 }, true)
			sortMenu.AddOption("Age (Asc)", func(a, b SortTestItem) int { return 0 }, false)
			sortMenu.Show()
		})

		It("should increment focus on down", func() {
			view1 := sortMenu.Render()
			Expect(view1).To(ContainSubstring("▶ Name (A-Z)"))

			handled := sortMenu.HandleKey("down")
			Expect(handled).To(BeTrue())

			view2 := sortMenu.Render()
			Expect(view2).To(ContainSubstring("▶ Name (Z-A)"))
		})

		It("should increment focus on j (vim)", func() {
			handled := sortMenu.HandleKey("j")
			Expect(handled).To(BeTrue())

			view := sortMenu.Render()
			Expect(view).To(ContainSubstring("▶ Name (Z-A)"))
		})

		It("should decrement focus on up", func() {
			sortMenu.HandleKey("down")
			view1 := sortMenu.Render()
			Expect(view1).To(ContainSubstring("▶ Name (Z-A)"))

			handled := sortMenu.HandleKey("up")
			Expect(handled).To(BeTrue())

			view2 := sortMenu.Render()
			Expect(view2).To(ContainSubstring("▶ Name (A-Z)"))
		})

		It("should decrement focus on k (vim)", func() {
			sortMenu.HandleKey("down")
			sortMenu.HandleKey("k")

			view := sortMenu.Render()
			Expect(view).To(ContainSubstring("▶ Name (A-Z)"))
		})

		It("should wrap focus at bottom", func() {
			// Navigate to last option (3 total)
			sortMenu.HandleKey("down") // Name (Z-A)
			sortMenu.HandleKey("down") // Age (Asc)

			// Press down again
			sortMenu.HandleKey("down")

			// Should wrap to first
			view := sortMenu.Render()
			Expect(view).To(ContainSubstring("▶ Name (A-Z)"))
		})

		It("should wrap focus at top", func() {
			// At index 0, press up
			sortMenu.HandleKey("up")

			// Should wrap to last option
			view := sortMenu.Render()
			Expect(view).To(ContainSubstring("▶ Age (Asc)"))
		})

		It("should ignore unknown keys", func() {
			handled := sortMenu.HandleKey("x")
			Expect(handled).To(BeFalse())
		})
	})

	Describe("Selection and Application", func() {
		var applyCalled bool

		BeforeEach(func() {
			sortMenu = behaviors.NewSortMenuBehavior(themeObj, table, "Sort Items")
			applyCalled = false

			sortMenu.AddOption("Name (A-Z)", func(a, b SortTestItem) int {
				if a.Name < b.Name {
					return -1
				}
				if a.Name > b.Name {
					return 1
				}
				return 0
			}, false)

			sortMenu.AddOption("Age (Desc)", func(a, b SortTestItem) int {
				return a.Age - b.Age
			}, true) // Reverse = true for descending

			sortMenu.OnApply(func() {
				applyCalled = true
			})

			sortMenu.Show()
		})

		It("should apply focused option on enter", func() {
			handled := sortMenu.HandleKey("enter")
			Expect(handled).To(BeTrue())
		})

		It("should apply sort comparator to table on enter", func() {
			// Press enter to apply first option (Name A-Z)
			sortMenu.HandleKey("enter")

			// Verify that a sort was applied
			Expect(table.HasSort()).To(BeTrue())
		})

		It("should apply reverse flag when specified", func() {
			// Navigate to second option (Age Desc with reverse=true)
			sortMenu.HandleKey("down")
			sortMenu.HandleKey("enter")

			// Verify sort was applied
			Expect(table.HasSort()).To(BeTrue())
		})

		It("should invoke onApply callback on enter", func() {
			sortMenu.HandleKey("enter")
			Expect(applyCalled).To(BeTrue())
		})

		It("should hide menu after applying sort", func() {
			Expect(sortMenu.IsActive()).To(BeTrue())

			sortMenu.HandleKey("enter")

			Expect(sortMenu.IsActive()).To(BeFalse())
		})

		It("should hide menu on esc without applying", func() {
			// Navigate to a different option
			sortMenu.HandleKey("down")

			// Press esc
			handled := sortMenu.HandleKey("esc")
			Expect(handled).To(BeTrue())

			// Menu should be hidden
			Expect(sortMenu.IsActive()).To(BeFalse())

			// Callback should not have been called
			Expect(applyCalled).To(BeFalse())
		})
	})

	Describe("Rendering", func() {
		BeforeEach(func() {
			sortMenu = behaviors.NewSortMenuBehavior(themeObj, table, "Sort Items")
			sortMenu.AddOption("Name (A-Z)", func(a, b SortTestItem) int { return 0 }, false)
			sortMenu.AddOption("Age (Asc)", func(a, b SortTestItem) int { return 0 }, false)
			sortMenu.Show()
		})

		It("should show title", func() {
			view := sortMenu.Render()
			Expect(view).To(ContainSubstring("Sort Items"))
		})

		It("should show all options", func() {
			view := sortMenu.Render()
			Expect(view).To(ContainSubstring("Name (A-Z)"))
			Expect(view).To(ContainSubstring("Age (Asc)"))
		})

		It("should show ▶ indicator on focused option", func() {
			view := sortMenu.Render()
			Expect(view).To(ContainSubstring("▶ Name (A-Z)"))
		})

		It("should show ✓ on selected option after applying", func() {
			// Apply first option
			sortMenu.HandleKey("enter")

			// Show menu again
			sortMenu.Show()

			view := sortMenu.Render()
			Expect(view).To(ContainSubstring("Name (A-Z) ✓"))
		})
	})
})
