package event_test

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/tui/views/shared"
	"github.com/baphled/kariya/internal/ui/uikit/theme"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter", func() {
	var (
		view       *event.Filter
		fields     []event.FilterFieldConfig
		sortFields []shared.SortFieldOption
	)

	BeforeEach(func() {
		fields = []event.FilterFieldConfig{
			{
				Key:   "company",
				Title: "Filter by Company",
				Options: []forms.SelectOption{
					{Key: "acme", Value: "Acme"},
					{Key: "globex", Value: "Globex"},
				},
			},
			{
				Key:   "category",
				Title: "Filter by Category",
				Options: []forms.SelectOption{
					{Key: "tech", Value: "Tech"},
					{Key: "ops", Value: "Ops"},
				},
			},
		}
		sortFields = []shared.SortFieldOption{
			{Key: "date", Label: "Date"},
			{Key: "text", Label: "Text"},
		}
	})

	Describe("Construction", func() {
		It("should create a filter view", func() {
			view = event.NewFilter(fields, sortFields, nil)
			Expect(view).NotTo(BeNil())
		})

		It("should accept defaults", func() {
			defaults := &event.FilterFormData{
				Selections: map[string][]string{"company": {"acme"}},
				SortBy:     "date",
				SortOrder:  "desc",
			}
			view = event.NewFilter(fields, sortFields, defaults)
			Expect(view).NotTo(BeNil())
		})

		It("should accept nil fields", func() {
			view = event.NewFilter(nil, sortFields, nil)
			Expect(view).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		It("should return a command from the form", func() {
			view = event.NewFilter(fields, sortFields, nil)
			cmd := view.Init()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			view = event.NewFilter(fields, sortFields, nil)
		})

		Context("with WindowSizeMsg", func() {
			It("should return nil result", func() {
				msg := tea.WindowSizeMsg{Width: 120, Height: 40}
				_, result := view.Update(msg)
				Expect(result).To(BeNil())
			})
		})

		Context("with Escape key", func() {
			It("should return CancelViewResult", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				_, result := view.Update(msg)
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultCancel))
				_, ok := result.(*widgets.CancelViewResult)
				Expect(ok).To(BeTrue())
			})
		})

		Context("when form completes", func() {
			var processViewCmd func(*event.Filter, tea.Cmd) widgets.ViewResult

			BeforeEach(func() {
				processViewCmd = func(v *event.Filter, cmd tea.Cmd) widgets.ViewResult {
					if cmd == nil {
						return nil
					}
					msg := cmd()
					if msg == nil {
						return nil
					}
					if batch, ok := msg.(tea.BatchMsg); ok {
						for _, c := range batch {
							if r := processViewCmd(v, c); r != nil {
								return r
							}
						}
						return nil
					}
					nextCmd, result := v.Update(msg)
					if result != nil {
						return result
					}
					return processViewCmd(v, nextCmd)
				}
			})

			It("should return SubmitViewResult with FilterFormData", func() {
				view = event.NewFilter(nil, sortFields, nil)

				processViewCmd(view, view.Init())

				var result widgets.ViewResult
				for attempts := 0; attempts < 20 && result == nil; attempts++ {
					cmd, r := view.Update(tea.KeyMsg{Type: tea.KeyEnter})
					if r != nil {
						result = r
						break
					}
					result = processViewCmd(view, cmd)
				}

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultSubmit))
				submitResult, ok := result.(*widgets.SubmitViewResult)
				Expect(ok).To(BeTrue())
				_, ok = submitResult.FormData.(event.FilterFormData)
				Expect(ok).To(BeTrue())
			})

			It("should sync multi-select values on submit", func() {
				view = event.NewFilter(fields, sortFields, nil)

				processViewCmd(view, view.Init())

				var result widgets.ViewResult
				for attempts := 0; attempts < 30 && result == nil; attempts++ {
					cmd, r := view.Update(tea.KeyMsg{Type: tea.KeyEnter})
					if r != nil {
						result = r
						break
					}
					result = processViewCmd(view, cmd)
				}

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultSubmit))
				submitResult, ok := result.(*widgets.SubmitViewResult)
				Expect(ok).To(BeTrue())
				formData, ok := submitResult.FormData.(event.FilterFormData)
				Expect(ok).To(BeTrue())
				Expect(formData.Selections).NotTo(BeNil())
				Expect(formData.SortBy).NotTo(BeEmpty())
				Expect(formData.SortOrder).NotTo(BeEmpty())
			})
		})

		Context("when form is aborted", func() {
			It("should return CancelViewResult", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				_, result := view.Update(msg)
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultCancel))
			})
		})
	})

	Describe("RenderContent", func() {
		It("should return non-empty string", func() {
			view = event.NewFilter(fields, sortFields, nil)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})
	})

	Describe("HelpText", func() {
		BeforeEach(func() {
			view = event.NewFilter(fields, sortFields, nil)
		})

		It("should return non-empty string", func() {
			Expect(view.HelpText()).NotTo(BeEmpty())
		})

		It("should contain Tab hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Tab"))
		})

		It("should contain Enter hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Enter"))
		})

		It("should contain Esc hint", func() {
			Expect(view.HelpText()).To(ContainSubstring("Esc"))
		})

		It("should render help text with terminal info", func() {
			view.SetTerminalInfo(80, 24)
			helpText := view.HelpText()
			Expect(helpText).NotTo(BeEmpty())
			Expect(helpText).To(ContainSubstring("Tab"))
			Expect(helpText).To(ContainSubstring("Enter"))
			Expect(helpText).To(ContainSubstring("Esc"))
		})

		It("should render help text with theme set", func() {
			view.SetTerminalInfo(80, 24)
			view.SetTheme(theme.Default())
			helpText := view.HelpText()
			Expect(helpText).NotTo(BeEmpty())
			Expect(helpText).To(ContainSubstring("Tab"))
			Expect(helpText).To(ContainSubstring("Enter"))
			Expect(helpText).To(ContainSubstring("Esc"))
		})
	})
})
