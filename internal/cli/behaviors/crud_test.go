package behaviors_test

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/themes"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Test item type for CRUD tests
type CRUDTestItem struct {
	ID   int
	Name string
}

var _ = Describe("CRUDBehavior", func() {
	var (
		themeObj  themes.Theme
		table     *behaviors.TableBehavior[CRUDTestItem]
		crud      *behaviors.CRUDBehavior[CRUDTestItem]
		columns   []behaviors.ColumnDef
		formatter behaviors.RowFormatter[CRUDTestItem]
	)

	BeforeEach(func() {
		themeObj = themes.NewDefaultTheme()
		columns = []behaviors.ColumnDef{
			{Title: "ID", Width: 5},
			{Title: "Name", Width: 20},
		}
		formatter = func(item CRUDTestItem, index int) []string {
			return []string{
				string(rune(item.ID + '0')),
				item.Name,
			}
		}
		table = behaviors.NewTableBehavior(themeObj, columns, formatter)
		table.SetItems([]CRUDTestItem{
			{ID: 1, Name: "Item 1"},
			{ID: 2, Name: "Item 2"},
			{ID: 3, Name: "Item 3"},
		})
	})

	Describe("Construction", func() {
		It("should create with non-nil table", func() {
			crud = behaviors.NewCRUDBehavior(themeObj, table)
			Expect(crud).NotTo(BeNil())
		})

		It("should panic if table is nil", func() {
			Expect(func() {
				behaviors.NewCRUDBehavior[CRUDTestItem](themeObj, nil)
			}).To(Panic())
		})

		It("should default to ModeList", func() {
			crud = behaviors.NewCRUDBehavior(themeObj, table)
			Expect(crud.Mode()).To(Equal(behaviors.ModeList))
		})

		It("should have all callbacks nil by default", func() {
			crud = behaviors.NewCRUDBehavior(themeObj, table)
			// We can't directly test callbacks are nil, but we can test that operations are disabled
			cmd, handled := crud.HandleKey("n")
			Expect(cmd).To(BeNil())
			Expect(handled).To(BeFalse())
		})
	})

	Describe("Configuration", func() {
		BeforeEach(func() {
			crud = behaviors.NewCRUDBehavior(themeObj, table)
		})

		It("should set item namer", func() {
			namer := func(item CRUDTestItem) string { return item.Name }
			result := crud.ItemNamer(namer)
			Expect(result).To(Equal(crud)) // Check chaining
		})

		It("should enable create with callback", func() {
			called := false
			onCreate := func() tea.Cmd {
				called = true
				return nil
			}
			result := crud.OnCreate(onCreate)
			Expect(result).To(Equal(crud)) // Check chaining

			// Now create should work
			crud.HandleKey("n")
			Expect(called).To(BeTrue())
		})

		It("should enable edit with callback", func() {
			var receivedItem *CRUDTestItem
			onEdit := func(item CRUDTestItem) tea.Cmd {
				receivedItem = &item
				return nil
			}
			result := crud.OnEdit(onEdit)
			Expect(result).To(Equal(crud)) // Check chaining

			// Now edit should work
			crud.HandleKey("e")
			Expect(receivedItem).NotTo(BeNil())
			Expect(receivedItem.Name).To(Equal("Item 1"))
		})

		It("should enable delete with callback", func() {
			onDelete := func(item CRUDTestItem) tea.Cmd {
				return nil
			}
			result := crud.OnDelete(onDelete)
			Expect(result).To(Equal(crud)) // Check chaining

			// Delete shows modal, so we need to confirm it
			crud.HandleKey("d")
			Expect(crud.Mode()).To(Equal(behaviors.ModeDelete))
		})

		It("should support chaining", func() {
			result := crud.
				ItemNamer(func(item CRUDTestItem) string { return item.Name }).
				OnCreate(func() tea.Cmd { return nil }).
				OnEdit(func(item CRUDTestItem) tea.Cmd { return nil }).
				OnDelete(func(item CRUDTestItem) tea.Cmd { return nil })

			Expect(result).To(Equal(crud))
		})
	})

	Describe("Mode Management", func() {
		BeforeEach(func() {
			crud = behaviors.NewCRUDBehavior(themeObj, table)
		})

		It("should return current mode", func() {
			Expect(crud.Mode()).To(Equal(behaviors.ModeList))
		})

		It("should detect list mode", func() {
			Expect(crud.IsInListMode()).To(BeTrue())
			crud.OnCreate(func() tea.Cmd { return nil })
			crud.HandleKey("n")
			Expect(crud.IsInListMode()).To(BeFalse())
		})

		It("should return to list mode", func() {
			crud.OnCreate(func() tea.Cmd { return nil })
			crud.HandleKey("n")
			Expect(crud.Mode()).To(Equal(behaviors.ModeCreate))

			crud.ReturnToList()
			Expect(crud.Mode()).To(Equal(behaviors.ModeList))
			Expect(crud.IsInListMode()).To(BeTrue())
		})
	})

	Describe("Key Handling - Create", func() {
		var createCalled bool
		var createCmd tea.Cmd

		BeforeEach(func() {
			crud = behaviors.NewCRUDBehavior(themeObj, table)
			createCalled = false
			createCmd = func() tea.Msg { return "create" }
		})

		It("should handle 'n' when create enabled", func() {
			crud.OnCreate(func() tea.Cmd {
				createCalled = true
				return createCmd
			})

			cmd, handled := crud.HandleKey("n")
			Expect(handled).To(BeTrue())
			Expect(createCalled).To(BeTrue())
			Expect(cmd).NotTo(BeNil()) // Cmd was returned
			Expect(crud.Mode()).To(Equal(behaviors.ModeCreate))
		})

		It("should ignore 'n' when create disabled", func() {
			cmd, handled := crud.HandleKey("n")
			Expect(handled).To(BeFalse())
			Expect(cmd).To(BeNil())
			Expect(createCalled).To(BeFalse())
		})
	})

	Describe("Key Handling - Edit", func() {
		var editCalled bool
		var editCmd tea.Cmd
		var receivedItem *CRUDTestItem

		BeforeEach(func() {
			crud = behaviors.NewCRUDBehavior(themeObj, table)
			editCalled = false
			editCmd = func() tea.Msg { return "edit" }
			receivedItem = nil
		})

		It("should handle 'e' when edit enabled and item selected", func() {
			crud.OnEdit(func(item CRUDTestItem) tea.Cmd {
				editCalled = true
				receivedItem = &item
				return editCmd
			})

			cmd, handled := crud.HandleKey("e")
			Expect(handled).To(BeTrue())
			Expect(editCalled).To(BeTrue())
			Expect(cmd).NotTo(BeNil()) // Cmd was returned
			Expect(crud.Mode()).To(Equal(behaviors.ModeEdit))
			Expect(receivedItem).NotTo(BeNil())
			Expect(receivedItem.Name).To(Equal("Item 1"))
		})

		It("should ignore 'e' when no item selected", func() {
			crud.OnEdit(func(item CRUDTestItem) tea.Cmd {
				editCalled = true
				return editCmd
			})

			// Empty table
			table.SetItems([]CRUDTestItem{})

			cmd, handled := crud.HandleKey("e")
			Expect(handled).To(BeFalse())
			Expect(cmd).To(BeNil())
			Expect(editCalled).To(BeFalse())
		})

		It("should ignore 'e' when edit disabled", func() {
			cmd, handled := crud.HandleKey("e")
			Expect(handled).To(BeFalse())
			Expect(cmd).To(BeNil())
			Expect(editCalled).To(BeFalse())
		})
	})

	Describe("Key Handling - Delete", func() {
		var deleteCalled bool
		var deleteCmd tea.Cmd

		BeforeEach(func() {
			crud = behaviors.NewCRUDBehavior(themeObj, table)
			deleteCalled = false
			deleteCmd = func() tea.Msg { return "delete" }
		})

		It("should show confirm modal when delete enabled", func() {
			crud.OnDelete(func(item CRUDTestItem) tea.Cmd {
				deleteCalled = true
				return deleteCmd
			})

			cmd, handled := crud.HandleKey("d")
			Expect(handled).To(BeTrue())
			Expect(cmd).To(BeNil()) // Modal shown, no cmd yet
			Expect(crud.Mode()).To(Equal(behaviors.ModeDelete))
			Expect(deleteCalled).To(BeFalse()) // Not called until confirmed
		})

		It("should use ItemNamer for confirmation message", func() {
			crud.ItemNamer(func(item CRUDTestItem) string {
				return "Custom: " + item.Name
			})
			crud.OnDelete(func(item CRUDTestItem) tea.Cmd {
				return nil
			})

			crud.HandleKey("d")
			view := crud.RenderDeleteConfirm()
			Expect(view).To(ContainSubstring("Custom: Item 1"))
		})

		It("should use 'this item' when ItemNamer not set", func() {
			crud.OnDelete(func(item CRUDTestItem) tea.Cmd {
				return nil
			})

			crud.HandleKey("d")
			view := crud.RenderDeleteConfirm()
			Expect(view).To(ContainSubstring("this item"))
		})

		It("should ignore 'd' when no item selected", func() {
			crud.OnDelete(func(item CRUDTestItem) tea.Cmd {
				deleteCalled = true
				return deleteCmd
			})

			// Empty table
			table.SetItems([]CRUDTestItem{})

			cmd, handled := crud.HandleKey("d")
			Expect(handled).To(BeFalse())
			Expect(cmd).To(BeNil())
			Expect(deleteCalled).To(BeFalse())
		})

		It("should ignore 'd' when delete disabled", func() {
			cmd, handled := crud.HandleKey("d")
			Expect(handled).To(BeFalse())
			Expect(cmd).To(BeNil())
			Expect(deleteCalled).To(BeFalse())
		})
	})

	Describe("Update - Delete Confirmation", func() {
		var deleteCalled bool
		var deleteCmd tea.Cmd
		var receivedItem *CRUDTestItem

		BeforeEach(func() {
			crud = behaviors.NewCRUDBehavior(themeObj, table)
			deleteCalled = false
			deleteCmd = func() tea.Msg { return "delete" }
			receivedItem = nil

			crud.OnDelete(func(item CRUDTestItem) tea.Cmd {
				deleteCalled = true
				receivedItem = &item
				return deleteCmd
			})
		})

		It("should invoke onDelete when user confirms", func() {
			// Show delete modal
			crud.HandleKey("d")
			Expect(crud.Mode()).To(Equal(behaviors.ModeDelete))

			// Simulate user confirming with 'y' key
			cmd, handled := crud.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(handled).To(BeTrue())
			Expect(deleteCalled).To(BeTrue())
			Expect(cmd).NotTo(BeNil()) // Cmd was returned
			Expect(receivedItem).NotTo(BeNil())
			Expect(receivedItem.Name).To(Equal("Item 1"))
			Expect(crud.Mode()).To(Equal(behaviors.ModeList)) // Should return to list
		})

		It("should return to list when user cancels", func() {
			// Show delete modal
			crud.HandleKey("d")
			Expect(crud.Mode()).To(Equal(behaviors.ModeDelete))

			// Simulate user canceling with Esc
			cmd, handled := crud.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(handled).To(BeTrue())
			Expect(deleteCalled).To(BeFalse())
			Expect(cmd).To(BeNil())
			Expect(crud.Mode()).To(Equal(behaviors.ModeList))
		})

		It("should pass selected item to onDelete callback", func() {
			// Select second item
			table.SetSelectedIndex(1)

			// Show delete modal
			crud.HandleKey("d")

			// Confirm
			crud.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(receivedItem).NotTo(BeNil())
			Expect(receivedItem.Name).To(Equal("Item 2"))
		})

		It("should return false when not in ModeDelete", func() {
			cmd, handled := crud.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(handled).To(BeFalse())
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Rendering", func() {
		BeforeEach(func() {
			crud = behaviors.NewCRUDBehavior(themeObj, table)
		})

		It("should render delete confirmation modal", func() {
			crud.OnDelete(func(item CRUDTestItem) tea.Cmd { return nil })
			crud.HandleKey("d")

			view := crud.RenderDeleteConfirm()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Delete"))
		})

		It("should include 'n New' when create enabled", func() {
			crud.OnCreate(func() tea.Cmd { return nil })
			help := crud.RenderHelpKeys()
			Expect(help).To(ContainSubstring("n"))
			Expect(help).To(ContainSubstring("New"))
		})

		It("should include 'e Edit' when edit enabled", func() {
			crud.OnEdit(func(item CRUDTestItem) tea.Cmd { return nil })
			help := crud.RenderHelpKeys()
			Expect(help).To(ContainSubstring("e"))
			Expect(help).To(ContainSubstring("Edit"))
		})

		It("should include 'd Delete' when delete enabled", func() {
			crud.OnDelete(func(item CRUDTestItem) tea.Cmd { return nil })
			help := crud.RenderHelpKeys()
			Expect(help).To(ContainSubstring("d"))
			Expect(help).To(ContainSubstring("Delete"))
		})

		It("should exclude disabled operations", func() {
			// Only enable create
			crud.OnCreate(func() tea.Cmd { return nil })

			help := crud.RenderHelpKeys()
			Expect(help).To(ContainSubstring("n New"))
			Expect(help).NotTo(ContainSubstring("Edit"))
			Expect(help).NotTo(ContainSubstring("Delete"))
		})

		It("should return empty string when no operations enabled", func() {
			help := crud.RenderHelpKeys()
			Expect(help).To(Equal(""))
		})
	})
})
