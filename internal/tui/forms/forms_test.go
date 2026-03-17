package forms_test

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/ui/themes"
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

var _ = Describe("Update", func() {
	It("should update form with a message", func() {
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

		updatedForm, _ := forms.Update(form, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
		Expect(updatedForm).NotTo(BeNil())
	})

	It("should return original form if update fails", func() {
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

		updatedForm, _ := forms.Update(form, nil)
		Expect(updatedForm).To(Equal(form))
	})

	It("should handle window size message", func() {
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

		updatedForm, _ := forms.Update(form, tea.WindowSizeMsg{Width: 100, Height: 30})
		Expect(updatedForm).NotTo(BeNil())
	})

	It("should handle key messages", func() {
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

		updatedForm, _ := forms.Update(form, tea.KeyMsg{Type: tea.KeyEnter})
		Expect(updatedForm).NotTo(BeNil())
	})
})

var _ = Describe("WithCustomSubmitKey", func() {
	It("should configure custom submit key", func() {
		var value string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("name").
					Title("Name").
					Value(&value),
			),
		).WithWidth(50)

		customForm := forms.WithCustomSubmitKey(form, "ctrl+s")
		Expect(customForm).NotTo(BeNil())
	})

	It("should accept multiple key bindings", func() {
		var value string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("name").
					Title("Name").
					Value(&value),
			),
		).WithWidth(50)

		customForm := forms.WithCustomSubmitKey(form, "ctrl+s", "enter")
		Expect(customForm).NotTo(BeNil())
	})
})

var _ = Describe("NewMultiSelect", func() {
	It("should create a multi-select field with options", func() {
		options := []forms.SelectOption{
			{Key: "a", Value: "Option A"},
			{Key: "b", Value: "Option B"},
		}

		field := forms.NewMultiSelect("tags", "Tags", "Select tags", options, 0)
		Expect(field).NotTo(BeNil())
	})

	It("should create a multi-select field with limit", func() {
		options := []forms.SelectOption{
			{Key: "a", Value: "Option A"},
			{Key: "b", Value: "Option B"},
			{Key: "c", Value: "Option C"},
		}

		field := forms.NewMultiSelect("tags", "Tags", "Select tags", options, 2)
		Expect(field).NotTo(BeNil())
	})

	It("should create a multi-select field without description", func() {
		options := []forms.SelectOption{
			{Key: "a", Value: "Option A"},
		}

		field := forms.NewMultiSelect("tags", "Tags", "", options, 0)
		Expect(field).NotTo(BeNil())
	})
})

var _ = Describe("CaptureEventFormData", func() {
	Describe("ConfirmSubmit", func() {
		It("should set SubmitConfirmed to true", func() {
			data := &forms.CaptureEventFormData{
				Text:            "Test event",
				SubmitConfirmed: false,
			}

			data.ConfirmSubmit()
			Expect(data.SubmitConfirmed).To(BeTrue())
		})
	})

	Describe("NewCaptureEventFormData", func() {
		It("should create form data with default values", func() {
			data := forms.NewCaptureEventFormData()

			Expect(data).NotTo(BeNil())
			Expect(data.Text).To(Equal(""))
			Expect(data.Date).To(Equal(""))
			Expect(data.Company).To(Equal(""))
			Expect(data.Project).To(Equal(""))
			Expect(data.Tags).To(BeEmpty())
			Expect(data.Categories).To(BeEmpty())
			Expect(data.SubmitConfirmed).To(BeFalse())
		})
	})

	Describe("NewCaptureEventForm", func() {
		It("should create a form with quick strategy", func() {
			data := forms.NewCaptureEventFormData()
			form := forms.NewCaptureEventForm(data, "quick", 80, 20)

			Expect(form).NotTo(BeNil())
		})

		It("should create a form with manual strategy", func() {
			data := forms.NewCaptureEventFormData()
			form := forms.NewCaptureEventForm(data, "manual", 80, 20)

			Expect(form).NotTo(BeNil())
		})

		It("should create a form with custom dimensions", func() {
			data := forms.NewCaptureEventFormData()
			form := forms.NewCaptureEventForm(data, "quick", 100, 30)

			Expect(form).NotTo(BeNil())
		})
	})

	Describe("NewCaptureEventFormForModal", func() {
		It("should create a modal form with quick strategy", func() {
			data := forms.NewCaptureEventFormData()
			form := forms.NewCaptureEventFormForModal(data, "quick", 80, 20)

			Expect(form).NotTo(BeNil())
		})

		It("should create a modal form with manual strategy", func() {
			data := forms.NewCaptureEventFormData()
			form := forms.NewCaptureEventFormForModal(data, "manual", 80, 20)

			Expect(form).NotTo(BeNil())
		})

		It("should create a modal form with custom dimensions", func() {
			data := forms.NewCaptureEventFormData()
			form := forms.NewCaptureEventFormForModal(data, "quick", 100, 30)

			Expect(form).NotTo(BeNil())
		})
	})

	Describe("GetCaptureEventFormData", func() {
		It("should extract form data from event", func() {
			eventDate, _ := time.Parse("2006-01-02", "2024-01-15")
			event := fixtures.Event("test-id")
			event.Text = "Test event"
			event.Date = eventDate
			event.Company = "TechCorp"
			event.Project = "Project X"
			event.Tags = []string{"technical", "achievement"}
			event.Categories = []string{"leadership"}

			data := forms.GetCaptureEventFormData(event)

			Expect(data).NotTo(BeNil())
			Expect(data.Text).To(Equal("Test event"))
			Expect(data.Company).To(Equal("TechCorp"))
			Expect(data.Project).To(Equal("Project X"))
			Expect(data.Tags).To(ContainElements("technical", "achievement"))
			Expect(data.Categories).To(ContainElement("leadership"))
		})

		It("should handle nil tags and categories", func() {
			eventDate, _ := time.Parse("2006-01-02", "2024-01-15")
			event := fixtures.Event("test-id")
			event.Text = "Test event"
			event.Date = eventDate
			event.Company = "TechCorp"
			event.Project = "Project X"
			event.Tags = nil
			event.Categories = nil

			data := forms.GetCaptureEventFormData(event)

			Expect(data).NotTo(BeNil())
			Expect(data.Tags).To(BeEmpty())
			Expect(data.Categories).To(BeEmpty())
		})
	})
})

var _ = Describe("EditorFields", func() {
	Describe("EditorUpdate", func() {
		It("should handle window size message", func() {
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

			fe := &forms.EditorFields{
				Form:   form,
				Width:  50,
				Height: 20,
			}

			model := &mockModel{}
			msg := tea.WindowSizeMsg{Width: 100, Height: 30}

			_, _ = forms.EditorUpdate(fe, model, msg, func() (tea.Model, tea.Cmd) {
				return model, nil
			}, func() tea.Msg {
				return nil
			})

			Expect(fe.Width).To(Equal(100))
			Expect(fe.Height).To(Equal(30))
		})

		It("should handle escape key", func() {
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

			fe := &forms.EditorFields{
				Form:      form,
				Cancelled: false,
			}

			model := &mockModel{}
			msg := tea.KeyMsg{Type: tea.KeyEscape}

			_, _ = forms.EditorUpdate(fe, model, msg, func() (tea.Model, tea.Cmd) {
				return model, nil
			}, func() tea.Msg {
				return nil
			})

			Expect(fe.Cancelled).To(BeTrue())
		})

		It("should handle q key", func() {
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

			fe := &forms.EditorFields{
				Form: form,
			}

			model := &mockModel{}
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

			_, cmd := forms.EditorUpdate(fe, model, msg, func() (tea.Model, tea.Cmd) {
				return model, nil
			}, func() tea.Msg {
				return "quit"
			})

			Expect(cmd).NotTo(BeNil())
		})

		It("should handle ctrl+c key", func() {
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

			fe := &forms.EditorFields{
				Form: form,
			}

			model := &mockModel{}
			msg := tea.KeyMsg{Type: tea.KeyCtrlC}

			_, cmd := forms.EditorUpdate(fe, model, msg, func() (tea.Model, tea.Cmd) {
				return model, nil
			}, func() tea.Msg {
				return "quit"
			})

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("RenderEditorView", func() {
		It("should render editor view with title and help", func() {
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

			fe := &forms.EditorFields{
				Form:   form,
				Width:  80,
				Height: 20,
				Theme:  themes.NewDefaultTheme(),
			}

			view := forms.RenderEditorView(fe, "Test Editor", "test-context")
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Test Editor"))
		})

		It("should render editor view with error", func() {
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

			fe := &forms.EditorFields{
				Form:   form,
				Width:  80,
				Height: 20,
				Err:    forms.ErrRequired,
				Theme:  themes.NewDefaultTheme(),
			}

			view := forms.RenderEditorView(fe, "Test Editor", "test-context")
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("required"))
		})
	})

})

type mockModel struct{}

func (m *mockModel) Init() tea.Cmd {
	return nil
}

func (m *mockModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *mockModel) View() string {
	return ""
}
