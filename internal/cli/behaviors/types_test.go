package behaviors_test

import (
	"testing"

	"github.com/baphled/kariya/internal/cli/behaviors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBehaviors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Behaviors Types Suite")
}

var _ = Describe("Behaviors Types", func() {
	Describe("ColumnDef", func() {
		It("should create a column definition", func() {
			col := behaviors.ColumnDef{
				Title: "Name",
				Width: 25,
			}

			Expect(col.Title).To(Equal("Name"))
			Expect(col.Width).To(Equal(25))
		})

		It("should allow zero width for auto-sizing", func() {
			col := behaviors.ColumnDef{
				Title: "Auto",
				Width: 0,
			}

			Expect(col.Width).To(Equal(0))
		})
	})

	Describe("CRUDMode", func() {
		It("should have List mode", func() {
			mode := behaviors.ModeList
			Expect(string(mode)).To(Equal("list"))
		})

		It("should have Create mode", func() {
			mode := behaviors.ModeCreate
			Expect(string(mode)).To(Equal("create"))
		})

		It("should have Edit mode", func() {
			mode := behaviors.ModeEdit
			Expect(string(mode)).To(Equal("edit"))
		})

		It("should have Delete mode", func() {
			mode := behaviors.ModeDelete
			Expect(string(mode)).To(Equal("delete"))
		})
	})

	Describe("MenuOption", func() {
		It("should create a menu option", func() {
			option := behaviors.MenuOption{
				Label:      "All Categories",
				Value:      "",
				IsSelected: true,
				IsDisabled: false,
			}

			Expect(option.Label).To(Equal("All Categories"))
			Expect(option.Value).To(Equal(""))
			Expect(option.IsSelected).To(BeTrue())
			Expect(option.IsDisabled).To(BeFalse())
		})

		It("should support any value type", func() {
			option := behaviors.MenuOption{
				Label: "Count",
				Value: 42,
			}

			Expect(option.Value).To(Equal(42))
		})
	})

	Describe("MenuSection", func() {
		It("should create a menu section with options", func() {
			section := behaviors.MenuSection{
				Title: "Category",
				Options: []behaviors.MenuOption{
					{Label: "All", Value: ""},
					{Label: "Backend", Value: "backend"},
				},
			}

			Expect(section.Title).To(Equal("Category"))
			Expect(section.Options).To(HaveLen(2))
			Expect(section.Options[0].Label).To(Equal("All"))
		})

		It("should allow empty options", func() {
			section := behaviors.MenuSection{
				Title:   "Empty Section",
				Options: []behaviors.MenuOption{},
			}

			Expect(section.Options).To(BeEmpty())
		})
	})

	Describe("Function Types", func() {
		Describe("RowFormatter", func() {
			It("should compile with correct signature", func() {
				// Compile-time check: this should compile if the type is defined correctly
				var formatter behaviors.RowFormatter[string]
				formatter = func(item string, index int) []string {
					return []string{item}
				}

				result := formatter("test", 0)
				Expect(result).To(Equal([]string{"test"}))
			})

			It("should work with complex types", func() {
				type TestItem struct {
					Name  string
					Count int
				}

				var formatter behaviors.RowFormatter[*TestItem]
				formatter = func(item *TestItem, index int) []string {
					return []string{item.Name, string(rune(item.Count))}
				}

				item := &TestItem{Name: "Test", Count: 42}
				result := formatter(item, 0)
				Expect(result).To(HaveLen(2))
				Expect(result[0]).To(Equal("Test"))
			})
		})

		Describe("FilterPredicate", func() {
			It("should compile with correct signature", func() {
				var predicate behaviors.FilterPredicate[int]
				predicate = func(item int) bool {
					return item > 5
				}

				Expect(predicate(10)).To(BeTrue())
				Expect(predicate(3)).To(BeFalse())
			})

			It("should work with complex types", func() {
				type TestItem struct {
					Active bool
				}

				var predicate behaviors.FilterPredicate[*TestItem]
				predicate = func(item *TestItem) bool {
					return item.Active
				}

				active := &TestItem{Active: true}
				inactive := &TestItem{Active: false}

				Expect(predicate(active)).To(BeTrue())
				Expect(predicate(inactive)).To(BeFalse())
			})
		})

		Describe("SortComparator", func() {
			It("should compile with correct signature", func() {
				var comparator behaviors.SortComparator[int]
				comparator = func(a, b int) int {
					return a - b
				}

				Expect(comparator(5, 10)).To(BeNumerically("<", 0))
				Expect(comparator(10, 5)).To(BeNumerically(">", 0))
				Expect(comparator(5, 5)).To(Equal(0))
			})

			It("should work with string comparison", func() {
				var comparator behaviors.SortComparator[string]
				comparator = func(a, b string) int {
					if a < b {
						return -1
					}
					if a > b {
						return 1
					}
					return 0
				}

				Expect(comparator("apple", "banana")).To(BeNumerically("<", 0))
				Expect(comparator("banana", "apple")).To(BeNumerically(">", 0))
				Expect(comparator("apple", "apple")).To(Equal(0))
			})
		})
	})
})
