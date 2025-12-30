package components

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NavigationMenu", func() {
	var (
		menu  NavigationMenuModel
		items []MenuItem
	)

	BeforeEach(func() {
		items = []MenuItem{
			{Label: "Capture Event", Shortcut: "c", Description: "Create a new career event", Action: "capture"},
			{Label: "List Events", Shortcut: "l", Description: "View all events", Action: "list"},
			{Label: "Metadata Review", Shortcut: "m", Description: "Review event metadata", Action: "metadata"},
			{Label: "Settings", Shortcut: "s", Description: "Configure application", Action: "settings"},
			{Label: "Help", Shortcut: "?", Description: "Show help information", Action: "help"},
		}
		menu = NewNavigationMenu(items, LayoutVertical)
		menu.SetWidth(80)
		menu.SetHeight(24)
	})

	Describe("Creation", func() {
		It("creates a menu with the given items", func() {
			Expect(menu.GetItems()).To(HaveLen(5))
		})

		It("initializes with first item selected", func() {
			Expect(menu.GetSelectedIndex()).To(Equal(0))
		})

		It("sets default layout to vertical", func() {
			Expect(menu.layout).To(Equal(LayoutVertical))
		})
	})

	Describe("Navigation", func() {
		It("moves selection down with Down arrow key", func() {
			updatedModel, _ := menu.Update(tea.KeyMsg{Type: tea.KeyDown})
			menu = updatedModel.(NavigationMenuModel)
			Expect(menu.GetSelectedIndex()).To(Equal(1))
		})

		It("moves selection down with j key", func() {
			updatedModel, _ := menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			menu = updatedModel.(NavigationMenuModel)
			Expect(menu.GetSelectedIndex()).To(Equal(1))
		})

		It("moves selection up with Up arrow key", func() {
			menu.SetSelectedIndex(2)
			updatedModel, _ := menu.Update(tea.KeyMsg{Type: tea.KeyUp})
			menu = updatedModel.(NavigationMenuModel)
			Expect(menu.GetSelectedIndex()).To(Equal(1))
		})

		It("wraps selection to last item when moving up from first", func() {
			menu.SetSelectedIndex(0)
			updatedModel, _ := menu.Update(tea.KeyMsg{Type: tea.KeyUp})
			menu = updatedModel.(NavigationMenuModel)
			Expect(menu.GetSelectedIndex()).To(Equal(4))
		})

		It("wraps selection to first item when moving down from last", func() {
			menu.SetSelectedIndex(4)
			updatedModel, _ := menu.Update(tea.KeyMsg{Type: tea.KeyDown})
			menu = updatedModel.(NavigationMenuModel)
			Expect(menu.GetSelectedIndex()).To(Equal(0))
		})
	})

	Describe("Selection", func() {
		It("selects item when Enter is pressed", func() {
			updatedModel, _ := menu.Update(tea.KeyMsg{Type: tea.KeyEnter})
			menu = updatedModel.(NavigationMenuModel)
			Expect(menu.IsSelected()).To(BeTrue())
			Expect(menu.GetSelectedAction()).To(Equal("capture"))
		})

		It("deselects item when Escape is pressed", func() {
			updatedModel, _ := menu.Update(tea.KeyMsg{Type: tea.KeyEnter})
			menu = updatedModel.(NavigationMenuModel)
			updatedModel, _ = menu.Update(tea.KeyMsg{Type: tea.KeyEsc})
			menu = updatedModel.(NavigationMenuModel)
			Expect(menu.IsSelected()).To(BeFalse())
		})

		It("returns correct selected item", func() {
			menu.SetSelectedIndex(2)
			selected := menu.GetSelectedItem()
			Expect(selected).NotTo(BeNil())
			Expect(selected.Label).To(Equal("Metadata Review"))
		})
	})

	Describe("Configuration", func() {
		It("sets width", func() {
			menu.SetWidth(100)
			Expect(menu.width).To(Equal(100))
		})

		It("sets height", func() {
			menu.SetHeight(30)
			Expect(menu.height).To(Equal(30))
		})

		It("toggles border display", func() {
			menu.SetShowBorder(false)
			Expect(menu.showBorder).To(BeFalse())
		})

		It("toggles number display", func() {
			menu.SetShowNumbers(true)
			Expect(menu.showNumbers).To(BeTrue())
		})
	})

	Describe("Rendering", func() {
		It("renders menu with items", func() {
			view := menu.View()
			Expect(view).To(ContainSubstring("Capture Event"))
		})

		It("renders empty menu message when no items", func() {
			emptyMenu := NewNavigationMenu([]MenuItem{}, LayoutVertical)
			view := emptyMenu.View()
			Expect(view).To(Equal("No menu items"))
		})

		It("shows selection indicator for selected item", func() {
			menu.SetSelectedIndex(1)
			view := menu.View()
			Expect(view).To(ContainSubstring("►"))
		})
	})

	Describe("State Management", func() {
		It("resets menu to initial state", func() {
			menu.SetSelectedIndex(3)
			updatedModel, _ := menu.Update(tea.KeyMsg{Type: tea.KeyEnter})
			menu = updatedModel.(NavigationMenuModel)
			menu.Reset()
			Expect(menu.GetSelectedIndex()).To(Equal(0))
			Expect(menu.IsSelected()).To(BeFalse())
		})

		It("handles window resize message", func() {
			updatedModel, _ := menu.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			menu = updatedModel.(NavigationMenuModel)
			Expect(menu.width).To(Equal(120))
			Expect(menu.height).To(Equal(40))
		})
	})

	Describe("BubbleTea Integration", func() {
		It("implements Model interface", func() {
			var model tea.Model = menu
			Expect(model).NotTo(BeNil())
		})

		It("Init returns no command", func() {
			cmd := menu.Init()
			Expect(cmd).To(BeNil())
		})

		It("Update returns model and command", func() {
			model, cmd := menu.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(model).NotTo(BeNil())
			Expect(cmd).To(BeNil())
		})

		It("View returns string", func() {
			view := menu.View()
			Expect(view).To(BeAssignableToTypeOf(""))
		})
	})
})
