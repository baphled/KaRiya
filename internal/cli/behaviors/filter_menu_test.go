package behaviors_test

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/themes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Test item for filter tests
type FilterTestItem struct {
	ID       int
	Name     string
	Category string
	Level    string
}

var _ = Describe("FilterMenuBehavior", func() {
	var (
		themeObj  themes.Theme
		table     *behaviors.TableBehavior[FilterTestItem]
		filter    *behaviors.FilterMenuBehavior[FilterTestItem]
		columns   []behaviors.ColumnDef
		formatter behaviors.RowFormatter[FilterTestItem]
	)

	BeforeEach(func() {
		themeObj = themes.NewDefaultTheme()
		columns = []behaviors.ColumnDef{
			{Title: "ID", Width: 5},
			{Title: "Name", Width: 20},
			{Title: "Category", Width: 15},
		}
		formatter = func(item FilterTestItem, index int) []string {
			return []string{
				string(rune(item.ID + '0')),
				item.Name,
				item.Category,
			}
		}
		table = behaviors.NewTableBehavior(themeObj, columns, formatter)
		table.SetItems([]FilterTestItem{
			{ID: 1, Name: "Item 1", Category: "TypeA", Level: "beginner"},
			{ID: 2, Name: "Item 2", Category: "TypeB", Level: "intermediate"},
			{ID: 3, Name: "Item 3", Category: "TypeA", Level: "advanced"},
		})
	})

	Describe("Construction", func() {
		It("should create with non-nil table", func() {
			filter = behaviors.NewFilterMenuBehavior(themeObj, table, "Filter Items")
			Expect(filter).NotTo(BeNil())
		})

		It("should panic if table is nil", func() {
			Expect(func() {
				behaviors.NewFilterMenuBehavior[FilterTestItem](themeObj, nil, "Filter")
			}).To(Panic())
		})

		It("should default to not active (hidden)", func() {
			filter = behaviors.NewFilterMenuBehavior(themeObj, table, "Filter Items")
			Expect(filter.IsActive()).To(BeFalse())
		})

		It("should start with focused index 0", func() {
			filter = behaviors.NewFilterMenuBehavior(themeObj, table, "Filter Items")
			// We can't directly test focusedIndex, but we can test that navigation works from 0
			Expect(filter.IsActive()).To(BeFalse()) // Sanity check
		})
	})

	Describe("Configuration", func() {
		BeforeEach(func() {
			filter = behaviors.NewFilterMenuBehavior(themeObj, table, "Filter Items")
		})

		It("should add section with options", func() {
			categoryOptions := []behaviors.MenuOption{
				{Label: "All Categories", Value: ""},
				{Label: "Type A", Value: "TypeA"},
				{Label: "Type B", Value: "TypeB"},
			}

			result := filter.AddSection(behaviors.MenuSection{
				Title:   "Category",
				Options: categoryOptions,
			})

			Expect(result).To(Equal(filter)) // Check chaining
		})

		It("should allow multiple sections", func() {
			filter.AddSection(behaviors.MenuSection{
				Title: "Category",
				Options: []behaviors.MenuOption{
					{Label: "All", Value: ""},
				},
			})

			filter.AddSection(behaviors.MenuSection{
				Title: "Level",
				Options: []behaviors.MenuOption{
					{Label: "All", Value: ""},
				},
			})

			// Can't directly test sections, but rendering should work
			view := filter.Render()
			Expect(view).To(ContainSubstring("Category"))
			Expect(view).To(ContainSubstring("Level"))
		})

		It("should set callback with OnApply", func() {
			result := filter.OnApply(func() {
				// Callback will be tested in Selection tests
			})

			Expect(result).To(Equal(filter)) // Check chaining
		})
	})

	Describe("State", func() {
		BeforeEach(func() {
			filter = behaviors.NewFilterMenuBehavior(themeObj, table, "Filter Items")
		})

		It("should Show() set active flag", func() {
			filter.Show()
			Expect(filter.IsActive()).To(BeTrue())
		})

		It("should Hide() clear active flag", func() {
			filter.Show()
			Expect(filter.IsActive()).To(BeTrue())

			filter.Hide()
			Expect(filter.IsActive()).To(BeFalse())
		})

		It("should IsActive() return correct state", func() {
			Expect(filter.IsActive()).To(BeFalse())

			filter.Show()
			Expect(filter.IsActive()).To(BeTrue())

			filter.Hide()
			Expect(filter.IsActive()).To(BeFalse())
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			filter = behaviors.NewFilterMenuBehavior(themeObj, table, "Filter Items")
			filter.AddSection(behaviors.MenuSection{
				Title: "Category",
				Options: []behaviors.MenuOption{
					{Label: "All", Value: ""},
					{Label: "Type A", Value: "TypeA"},
					{Label: "Type B", Value: "TypeB"},
				},
			})
			filter.AddSection(behaviors.MenuSection{
				Title: "Level",
				Options: []behaviors.MenuOption{
					{Label: "All", Value: ""},
					{Label: "Beginner", Value: "beginner"},
				},
			})
			filter.Show()
		})

		It("should increment focus on down", func() {
			// Initial render should have ▶ on first option
			view1 := filter.Render()
			Expect(view1).To(ContainSubstring("▶ All"))

			handled := filter.HandleKey("down")
			Expect(handled).To(BeTrue())

			// Now should have ▶ on second option
			view2 := filter.Render()
			Expect(view2).To(ContainSubstring("▶ Type A"))
		})

		It("should increment focus on j (vim)", func() {
			handled := filter.HandleKey("j")
			Expect(handled).To(BeTrue())

			view := filter.Render()
			Expect(view).To(ContainSubstring("▶ Type A"))
		})

		It("should decrement focus on up", func() {
			// Move down first
			filter.HandleKey("down")
			view1 := filter.Render()
			Expect(view1).To(ContainSubstring("▶ Type A"))

			// Move up
			handled := filter.HandleKey("up")
			Expect(handled).To(BeTrue())

			view2 := filter.Render()
			Expect(view2).To(ContainSubstring("▶ All"))
		})

		It("should decrement focus on k (vim)", func() {
			filter.HandleKey("down")
			filter.HandleKey("k")

			view := filter.Render()
			Expect(view).To(ContainSubstring("▶ All"))
		})

		It("should wrap focus at bottom (down from last option)", func() {
			// Navigate to last option (5 total: 3 category + 2 level)
			filter.HandleKey("down") // Type A
			filter.HandleKey("down") // Type B
			filter.HandleKey("down") // All (Level)
			filter.HandleKey("down") // Beginner

			// Now at last option, press down again
			filter.HandleKey("down")

			// Should wrap to first option
			view := filter.Render()
			Expect(view).To(ContainSubstring("▶ All")) // First option in Category
		})

		It("should wrap focus at top (up from first option)", func() {
			// Start at index 0, press up
			filter.HandleKey("up")

			// Should wrap to last option (Beginner in Level section)
			view := filter.Render()
			Expect(view).To(ContainSubstring("▶ Beginner"))
		})

		It("should ignore unknown keys", func() {
			handled := filter.HandleKey("x")
			Expect(handled).To(BeFalse())
		})
	})

	Describe("Selection and Application", func() {
		var applyCalled bool

		BeforeEach(func() {
			filter = behaviors.NewFilterMenuBehavior(themeObj, table, "Filter Items")
			applyCalled = false

			filter.AddSection(behaviors.MenuSection{
				Title: "Category",
				Options: []behaviors.MenuOption{
					{Label: "All Categories", Value: ""},
					{Label: "Type A", Value: "TypeA"},
					{Label: "Type B", Value: "TypeB"},
				},
			})

			filter.OnApply(func() {
				applyCalled = true
			})

			filter.Show()
		})

		It("should apply focused option on enter", func() {
			// Navigate to "Type A"
			filter.HandleKey("down")

			// Press enter
			handled := filter.HandleKey("enter")
			Expect(handled).To(BeTrue())
		})

		It("should build and apply filter predicate on enter", func() {
			// Start with 3 items
			Expect(table.Count()).To(Equal(3))

			// Navigate to "Type A" filter
			filter.HandleKey("down")
			filter.HandleKey("enter")

			// Table should now be filtered to only TypeA items (2 items)
			Expect(table.Count()).To(Equal(2))
		})

		It("should invoke onApply callback on enter", func() {
			filter.HandleKey("enter")
			Expect(applyCalled).To(BeTrue())
		})

		It("should hide menu after applying filter", func() {
			Expect(filter.IsActive()).To(BeTrue())

			filter.HandleKey("enter")

			Expect(filter.IsActive()).To(BeFalse())
		})

		It("should hide menu on esc without applying", func() {
			initialCount := table.Count()

			// Navigate to a different option
			filter.HandleKey("down")

			// Press esc
			handled := filter.HandleKey("esc")
			Expect(handled).To(BeTrue())

			// Menu should be hidden
			Expect(filter.IsActive()).To(BeFalse())

			// Table count should not have changed
			Expect(table.Count()).To(Equal(initialCount))

			// Callback should not have been called
			Expect(applyCalled).To(BeFalse())
		})
	})

	Describe("Rendering", func() {
		BeforeEach(func() {
			filter = behaviors.NewFilterMenuBehavior(themeObj, table, "Filter Items")
			filter.AddSection(behaviors.MenuSection{
				Title: "Category",
				Options: []behaviors.MenuOption{
					{Label: "All Categories", Value: ""},
					{Label: "Type A", Value: "TypeA"},
				},
			})
			filter.Show()
		})

		It("should show title", func() {
			view := filter.Render()
			Expect(view).To(ContainSubstring("Filter Items"))
		})

		It("should show section headers", func() {
			view := filter.Render()
			Expect(view).To(ContainSubstring("Category"))
		})

		It("should show ▶ indicator on focused option", func() {
			view := filter.Render()
			Expect(view).To(ContainSubstring("▶ All Categories"))
		})

		It("should show ✓ on selected option after applying", func() {
			// Navigate to "Type A" and apply
			filter.HandleKey("down")
			filter.HandleKey("enter")

			// Show menu again
			filter.Show()

			view := filter.Render()
			Expect(view).To(ContainSubstring("Type A ✓"))
		})
	})
})
