package skill_test

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/tui/views/skill"
	"github.com/baphled/kariya/internal/ui/uikit/theme"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter", func() {
	var (
		view     *skill.Filter
		fields   []skill.FilterFieldConfig
		defaults *skill.FilterFormData
	)

	BeforeEach(func() {
		fields = []skill.FilterFieldConfig{
			{
				Key:   "categories",
				Title: "Filter by Category",
				Options: []forms.SelectOption{
					{Key: "cat1", Value: "cat1"},
					{Key: "cat2", Value: "cat2"},
				},
			},
			{
				Key:   "levels",
				Title: "Filter by Level",
				Options: []forms.SelectOption{
					{Key: "junior", Value: "junior"},
					{Key: "senior", Value: "senior"},
				},
			},
		}
		defaults = &skill.FilterFormData{
			Selections: map[string][]string{
				"categories": {"cat1"},
				"levels":     {"junior"},
			},
			MinYearsStr: "",
			MaxYearsStr: "",
		}
	})

	Describe("Construction", func() {
		It("should create a filter view", func() {
			view = skill.NewFilter(fields, defaults)
			Expect(view).NotTo(BeNil())
		})

		It("should create a filter view with nil defaults", func() {
			view = skill.NewFilter(fields, nil)
			Expect(view).NotTo(BeNil())
		})

		It("should skip fields with empty options", func() {
			fieldsWithEmpty := []skill.FilterFieldConfig{
				{Key: "empty", Title: "Empty", Options: nil},
				{Key: "cats", Title: "Cats", Options: []forms.SelectOption{{Key: "a", Value: "a"}}},
			}
			view = skill.NewFilter(fieldsWithEmpty, nil)
			Expect(view).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		It("should return a command from the form", func() {
			view = skill.NewFilter(fields, defaults)
			cmd := view.Init()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			view = skill.NewFilter(fields, defaults)
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
			var (
				execCmd        func(tea.Cmd) tea.Msg
				processViewCmd func(*skill.Filter, tea.Cmd) widgets.ViewResult
			)

			BeforeEach(func() {
				execCmd = func(cmd tea.Cmd) tea.Msg {
					if cmd == nil {
						return nil
					}
					ch := make(chan tea.Msg, 1)
					go func() { ch <- cmd() }()
					select {
					case msg := <-ch:
						return msg
					case <-time.After(50 * time.Millisecond):
						return nil
					}
				}

				processViewCmd = func(v *skill.Filter, cmd tea.Cmd) widgets.ViewResult {
					msg := execCmd(cmd)
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
				emptyView := skill.NewFilter(nil, nil)

				processViewCmd(emptyView, emptyView.Init())

				var result widgets.ViewResult
				for attempts := 0; attempts < 10 && result == nil; attempts++ {
					cmd, r := emptyView.Update(tea.KeyMsg{Type: tea.KeyEnter})
					if r != nil {
						result = r
						break
					}
					result = processViewCmd(emptyView, cmd)
				}

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultSubmit))
				submitResult, ok := result.(*widgets.SubmitViewResult)
				Expect(ok).To(BeTrue())
				formData, ok := submitResult.FormData.(skill.FilterFormData)
				Expect(ok).To(BeTrue())
				Expect(formData.Selections).NotTo(BeNil())
				Expect(formData.MinYearsStr).To(Equal(""))
				Expect(formData.MaxYearsStr).To(Equal(""))
			})

			It("should sync multi-select selections on submit", func() {
				filterView := skill.NewFilter(fields, defaults)

				processViewCmd(filterView, filterView.Init())

				var result widgets.ViewResult
				for attempts := 0; attempts < 20 && result == nil; attempts++ {
					cmd, r := filterView.Update(tea.KeyMsg{Type: tea.KeyEnter})
					if r != nil {
						result = r
						break
					}
					result = processViewCmd(filterView, cmd)
				}

				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultSubmit))
				submitResult, ok := result.(*widgets.SubmitViewResult)
				Expect(ok).To(BeTrue())
				formData, ok := submitResult.FormData.(skill.FilterFormData)
				Expect(ok).To(BeTrue())
				Expect(formData.Selections).To(HaveKey("categories"))
				Expect(formData.Selections).To(HaveKey("levels"))
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
			view = skill.NewFilter(fields, defaults)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})
	})

	Describe("HelpText", func() {
		BeforeEach(func() {
			view = skill.NewFilter(fields, defaults)
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

	Describe("ValidateYearsInput", func() {
		It("should reject non-numeric input", func() {
			err := skill.ValidateYearsInput("abc")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("number"))
		})

		It("should reject negative input", func() {
			err := skill.ValidateYearsInput("-5")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("positive"))
		})

		It("should allow empty input", func() {
			err := skill.ValidateYearsInput("")
			Expect(err).ToNot(HaveOccurred())
		})

		It("should allow positive numeric input", func() {
			err := skill.ValidateYearsInput("3")
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
