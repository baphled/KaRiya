package behaviors_test

import (
	"github.com/baphled/kariya/internal/ui/behaviors"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListBehavior", func() {
	type testItem struct {
		Name string
	}

	formatter := func(item testItem, _ int) []string {
		return []string{item.Name}
	}

	columns := []behaviors.ColumnDef{
		{Title: "Name", Width: 20},
	}

	downKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	enterKey := tea.KeyMsg{Type: tea.KeyEnter}

	Describe("Construction", func() {
		It("returns non-nil", func() {
			w := behaviors.NewListBehavior(columns, formatter, nil)
			Expect(w).NotTo(BeNil())
		})

		It("starts with selection at index 0", func() {
			w := behaviors.NewListBehavior(columns, formatter, nil)
			items := []testItem{{Name: "A"}, {Name: "B"}}
			w.SetItems(items)
			Expect(w.GetSelectedIndex()).To(Equal(0))
		})
	})

	Describe("HandleKey", func() {
		It("consumes navigation keys and moves selection", func() {
			w := behaviors.NewListBehavior(columns, formatter, nil)
			items := []testItem{{Name: "A"}, {Name: "B"}, {Name: "C"}}
			w.SetItems(items)
			consumed := w.HandleKey(downKey)
			Expect(consumed).To(BeTrue())
			Expect(w.GetSelectedIndex()).To(Equal(1))
		})

		It("consumes select key and returns true", func() {
			w := behaviors.NewListBehavior(columns, formatter, nil)
			items := []testItem{{Name: "A"}, {Name: "B"}}
			w.SetItems(items)
			consumed := w.HandleKey(enterKey)
			Expect(consumed).To(BeTrue())
		})

		It("returns selected item after navigation", func() {
			w := behaviors.NewListBehavior(columns, formatter, nil)
			items := []testItem{{Name: "A"}, {Name: "B"}}
			w.SetItems(items)
			w.HandleKey(downKey)
			selected := w.GetSelectedItem()
			Expect(selected).NotTo(BeNil())
			Expect(selected.Name).To(Equal("B"))
		})

		It("returns false for unrecognised keys", func() {
			w := behaviors.NewListBehavior(columns, formatter, nil)
			w.SetItems([]testItem{{Name: "A"}})
			unknown := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("z")}
			Expect(w.HandleKey(unknown)).To(BeFalse())
		})
	})

	Describe("Rendering", func() {
		It("RenderContent returns non-empty string after SetItems", func() {
			w := behaviors.NewListBehavior(columns, formatter, nil)
			items := []testItem{{Name: "Alpha"}}
			w.SetItems(items)
			Expect(w.RenderContent()).NotTo(BeEmpty())
		})

		It("HelpText returns empty string with no badge funcs", func() {
			w := behaviors.NewListBehavior(columns, formatter, nil)
			Expect(w.HelpText()).To(BeEmpty())
		})

		It("HelpText with nil theme does not panic", func() {
			w := behaviors.NewListBehavior(columns, formatter, nil)
			Expect(func() { w.HelpText() }).NotTo(Panic())
		})

		It("SetTheme does not panic", func() {
			w := behaviors.NewListBehavior(columns, formatter, nil)
			Expect(func() { w.SetTheme(nil) }).NotTo(Panic())
		})
	})
})
