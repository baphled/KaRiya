package forms_test

import (
	"github.com/charmbracelet/huh"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/forms"
)

var _ = Describe("IsTextInputFocused", func() {
	Describe("when form is nil", func() {
		It("should return false", func() {
			Expect(forms.IsTextInputFocused(nil)).To(BeFalse())
		})
	})

	Describe("when form has an Input field focused", func() {
		It("should return true", func() {
			var value string
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Key("name").
						Title("Name").
						Value(&value),
				),
			).WithWidth(50)
			form.Init()

			Expect(forms.IsTextInputFocused(form)).To(BeTrue())
		})
	})

	Describe("when form has a Text field focused", func() {
		It("should return true", func() {
			var value string
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewText().
						Key("description").
						Title("Description").
						Value(&value),
				),
			).WithWidth(50)
			form.Init()

			Expect(forms.IsTextInputFocused(form)).To(BeTrue())
		})
	})

	Describe("when form has a Select field focused", func() {
		It("should return false", func() {
			var value string
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Key("theme").
						Title("Theme").
						Options(
							huh.NewOption("Light", "light"),
							huh.NewOption("Dark", "dark"),
						).
						Value(&value),
				),
			).WithWidth(50)
			form.Init()

			Expect(forms.IsTextInputFocused(form)).To(BeFalse())
		})
	})

	Describe("when form has a Confirm field focused", func() {
		It("should return false", func() {
			var value bool
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Key("confirm").
						Title("Proceed?").
						Value(&value),
				),
			).WithWidth(50)
			form.Init()

			Expect(forms.IsTextInputFocused(form)).To(BeFalse())
		})
	})

	Describe("when form has a MultiSelect field focused", func() {
		It("should return false", func() {
			var value []string
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewMultiSelect[string]().
						Key("options").
						Title("Options").
						Options(
							huh.NewOption("A", "a"),
							huh.NewOption("B", "b"),
						).
						Value(&value),
				),
			).WithWidth(50)
			form.Init()

			Expect(forms.IsTextInputFocused(form)).To(BeFalse())
		})
	})

	Describe("when navigating between fields", func() {
		It("should return false when moving from Input to Select", func() {
			var name string
			var theme string
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Key("name").
						Title("Name").
						Value(&name),
					huh.NewSelect[string]().
						Key("theme").
						Title("Theme").
						Options(
							huh.NewOption("Light", "light"),
							huh.NewOption("Dark", "dark"),
						).
						Value(&theme),
				),
			).WithWidth(50)
			form.Init()
			Expect(forms.IsTextInputFocused(form)).To(BeTrue())
			form.NextField()
			Expect(forms.IsTextInputFocused(form)).To(BeFalse())
		})
	})
})
