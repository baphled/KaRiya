package shared_test

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/tui/views/shared"
	"github.com/baphled/kariya/internal/ui/uikit/theme"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SortView", func() {
	var view *shared.SortView

	defaultFields := []shared.SortFieldOption{
		{Key: "name", Label: "Name"},
		{Key: "category", Label: "Category"},
		{Key: "level", Label: "Level"},
	}

	Describe("Construction", func() {
		It("should create a sort view", func() {
			view = shared.NewSortView(defaultFields, "name", "asc")
			Expect(view).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		It("should return a command from the form", func() {
			view = shared.NewSortView(defaultFields, "name", "asc")
			cmd := view.Init()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			view = shared.NewSortView(defaultFields, "name", "asc")
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
			It("should return SubmitViewResult with SortFormData", func() {
				var processViewCmd func(*shared.SortView, tea.Cmd) widgets.ViewResult
				processViewCmd = func(v *shared.SortView, cmd tea.Cmd) widgets.ViewResult {
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

				processViewCmd(view, view.Init())

				enterMsg := tea.KeyMsg{Type: tea.KeyEnter}

				cmd, result := view.Update(enterMsg)
				if result == nil {
					result = processViewCmd(view, cmd)
				}

				if result == nil {
					cmd, result = view.Update(enterMsg)
					if result == nil {
						result = processViewCmd(view, cmd)
					}
				}

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultSubmit))
				submitResult, ok := result.(*widgets.SubmitViewResult)
				Expect(ok).To(BeTrue())
				formData, ok := submitResult.FormData.(shared.SortFormData)
				Expect(ok).To(BeTrue())
				Expect(formData.SortBy).To(Equal("name"))
				Expect(formData.SortOrder).To(Equal("asc"))
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
			view = shared.NewSortView(defaultFields, "name", "asc")
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})
	})

	Describe("HelpText", func() {
		BeforeEach(func() {
			view = shared.NewSortView(defaultFields, "name", "asc")
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
