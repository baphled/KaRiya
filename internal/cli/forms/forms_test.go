package forms_test

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/forms"
)

var _ = Describe("Forms", func() {
	Describe("NewGroup", func() {
		It("creates a group from a single field", func() {
			field := huh.NewInput().Key("test").Title("Test")
			group := forms.NewGroup(field)

			Expect(group).NotTo(BeNil())
		})

		It("creates a group from multiple fields", func() {
			field1 := huh.NewInput().Key("first").Title("First")
			field2 := huh.NewInput().Key("second").Title("Second")
			group := forms.NewGroup(field1, field2)

			Expect(group).NotTo(BeNil())
		})
	})

	Describe("NewMultiSelect", func() {
		It("creates a multi-select with options", func() {
			options := []forms.SelectOption{
				{Key: "go", Value: "Go"},
				{Key: "python", Value: "Python"},
				{Key: "rust", Value: "Rust"},
			}

			multi := forms.NewMultiSelect("languages", "Languages", "Select languages", options, 3)

			Expect(multi).NotTo(BeNil())
		})

		It("creates a multi-select without description", func() {
			options := []forms.SelectOption{
				{Key: "one", Value: "One"},
			}

			multi := forms.NewMultiSelect("items", "Items", "", options, 0)

			Expect(multi).NotTo(BeNil())
		})

		It("creates a multi-select without limit", func() {
			options := []forms.SelectOption{
				{Key: "a", Value: "A"},
				{Key: "b", Value: "B"},
			}

			multi := forms.NewMultiSelect("letters", "Letters", "Pick letters", options, 0)

			Expect(multi).NotTo(BeNil())
		})
	})

	Describe("NewConfirm", func() {
		It("creates a confirm field with all parameters", func() {
			confirm := forms.NewConfirm("save", "Save Changes", "Are you sure?", "Yes", "No")

			Expect(confirm).NotTo(BeNil())
		})

		It("creates a confirm field without description", func() {
			confirm := forms.NewConfirm("delete", "Delete", "", "Yes", "No")

			Expect(confirm).NotTo(BeNil())
		})

		It("creates a confirm field without affirmative/negative", func() {
			confirm := forms.NewConfirm("confirm", "Confirm", "Please confirm", "", "")

			Expect(confirm).NotTo(BeNil())
		})
	})

	Describe("NewFormWithDimensions", func() {
		It("creates a form with specified dimensions", func() {
			group := huh.NewGroup(huh.NewInput().Key("test").Title("Test"))
			form := forms.NewFormWithDimensions(80, 24, group)

			Expect(form).NotTo(BeNil())
		})

		It("creates a form with zero height", func() {
			group := huh.NewGroup(huh.NewInput().Key("test").Title("Test"))
			form := forms.NewFormWithDimensions(80, 0, group)

			Expect(form).NotTo(BeNil())
		})

		It("creates a form with zero width", func() {
			group := huh.NewGroup(huh.NewInput().Key("test").Title("Test"))
			form := forms.NewFormWithDimensions(0, 24, group)

			Expect(form).NotTo(BeNil())
		})

		It("creates a form with both dimensions zero", func() {
			group := huh.NewGroup(huh.NewInput().Key("test").Title("Test"))
			form := forms.NewFormWithDimensions(0, 0, group)

			Expect(form).NotTo(BeNil())
		})

		It("creates a form with multiple groups", func() {
			group1 := huh.NewGroup(huh.NewInput().Key("name").Title("Name"))
			group2 := huh.NewGroup(huh.NewInput().Key("email").Title("Email"))
			form := forms.NewFormWithDimensions(80, 24, group1, group2)

			Expect(form).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		It("returns the updated form and command", func() {
			group := huh.NewGroup(huh.NewInput().Key("test").Title("Test"))
			form := huh.NewForm(group).WithTheme(huh.ThemeCatppuccin())

			updatedForm, cmd := forms.Update(form, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(updatedForm).NotTo(BeNil())
			_ = cmd
		})
	})
})
