package event_test

import (
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/types"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Form", func() {
	var view *event.Form

	Describe("Construction", func() {
		It("should create a non-nil Form with nil event and StrategyQuick", func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
			Expect(view).NotTo(BeNil())
		})

		It("should create Form with existing event and StrategyManual", func() {
			testEvent := display.EventFromDomain(fixtures.EventWith("ev-1", "Built scalable API platform", "TechCorp", "API Platform"))
			view = event.NewForm(testEvent, types.StrategyManual)
			Expect(view).NotTo(BeNil())
		})

		It("should create Form with nil event and StrategyManual", func() {
			view = event.NewForm(display.Event{}, types.StrategyManual)
			Expect(view).NotTo(BeNil())
		})

		It("should store the strategy correctly", func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
			Expect(view.GetStrategy()).To(Equal(types.StrategyQuick))
		})

		It("should store strategy as StrategyManual when provided", func() {
			view = event.NewForm(display.Event{}, types.StrategyManual)
			Expect(view.GetStrategy()).To(Equal(types.StrategyManual))
		})

		It("should create form data for nil event", func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
			formData := view.GetFormData()
			Expect(formData).NotTo(BeNil())
		})

		It("should create form data for existing event", func() {
			testEvent := display.EventFromDomain(fixtures.EventWith("ev-1", "Built scalable API platform", "TechCorp", "API Platform"))
			view = event.NewForm(testEvent, types.StrategyManual)
			formData := view.GetFormData()
			Expect(formData).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
		})

		It("should return nil command", func() {
			cmd := view.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update with WindowSizeMsg", func() {
		BeforeEach(func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
		})

		It("should update dimensions and return nil cmd and nil result", func() {
			msg := tea.WindowSizeMsg{Width: 120, Height: 40}
			cmd, result := view.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})
	})

	Describe("Update with KeyMsg", func() {
		BeforeEach(func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
			view.SetTerminalInfo(120, 40)
		})

		Describe("Esc key", func() {
			It("should return CancelViewResult with tea.KeyEsc type", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				_, result := view.Update(msg)
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultCancel))
				_, ok := result.(*widgets.CancelViewResult)
				Expect(ok).To(BeTrue())
			})

			It("should return CancelViewResult with esc string", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				_, result := view.Update(msg)
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultCancel))
				_, ok := result.(*widgets.CancelViewResult)
				Expect(ok).To(BeTrue())
			})
		})

		Describe("Ctrl+S key", func() {
			It("should return SubmitViewResult when formData exists", func() {
				msg := tea.KeyMsg{Type: tea.KeyCtrlS}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultSubmit))
			})

			It("should return SubmitViewResult with FormData as *forms.CaptureEventFormData", func() {
				msg := tea.KeyMsg{Type: tea.KeyCtrlS}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				submitResult, ok := result.(*widgets.SubmitViewResult)
				Expect(ok).To(BeTrue())
				formData, ok := submitResult.FormData.(*forms.CaptureEventFormData)
				Expect(ok).To(BeTrue())
				Expect(formData).NotTo(BeNil())
			})

			It("should set SubmitConfirmed to true after Ctrl+S", func() {
				msg := tea.KeyMsg{Type: tea.KeyCtrlS}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				submitResult, ok := result.(*widgets.SubmitViewResult)
				Expect(ok).To(BeTrue())
				formData, ok := submitResult.FormData.(*forms.CaptureEventFormData)
				Expect(ok).To(BeTrue())
				Expect(formData.SubmitConfirmed).To(BeTrue())
			})
		})

		Describe("Unknown key", func() {
			It("should return nil result for non-esc, non-ctrl+s key", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}}
				_, result := view.Update(msg)
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("RenderContent", func() {
		It("should return non-empty string when form has been built", func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
			view.SetTerminalInfo(120, 40)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should return empty string when form is nil", func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
			view = &event.Form{}
			content := view.RenderContent()
			Expect(content).To(BeEmpty())
		})
	})

	Describe("HelpText", func() {
		BeforeEach(func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
		})

		It("should return default footer", func() {
			helpText := view.HelpText()
			Expect(helpText).To(Equal("Esc: Back  Tab: Next field  Enter: Select"))
		})

		It("should return custom footer after SetFooter", func() {
			view.SetFooter("Custom footer")
			Expect(view.HelpText()).To(Equal("Custom footer"))
		})
	})

	Describe("GetFormData as input source", func() {
		It("should return empty fields for new form", func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
			input := view.GetFormData()
			Expect(input.Text).To(Equal(""))
			Expect(input.Date).To(Equal(""))
			Expect(input.Company).To(Equal(""))
			Expect(input.Project).To(Equal(""))
		})

		It("should return empty slices for Tags and Categories", func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
			input := view.GetFormData()
			Expect(input.Tags).To(BeEmpty())
			Expect(input.Categories).To(BeEmpty())
		})

		It("should return values from existing event", func() {
			testEvent := display.EventFromDomain(fixtures.EventWith("ev-1", "Built scalable API platform", "TechCorp", "API Platform"))
			view = event.NewForm(testEvent, types.StrategyManual)
			input := view.GetFormData()
			Expect(input.Text).To(Equal("Built scalable API platform"))
			Expect(input.Company).To(Equal("TechCorp"))
			Expect(input.Project).To(Equal("API Platform"))
		})
	})

	Describe("GetStrategy", func() {
		It("should return types.StrategyQuick when constructed with StrategyQuick", func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
			Expect(view.GetStrategy()).To(Equal(types.StrategyQuick))
		})

		It("should return types.StrategyManual when constructed with StrategyManual", func() {
			view = event.NewForm(display.Event{}, types.StrategyManual)
			Expect(view.GetStrategy()).To(Equal(types.StrategyManual))
		})
	})

	Describe("GetFormData", func() {
		It("should return non-nil *forms.CaptureEventFormData", func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
			formData := view.GetFormData()
			Expect(formData).NotTo(BeNil())
		})

		It("should return FormData with empty initial values for new form", func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
			formData := view.GetFormData()
			Expect(formData.Text).To(Equal(""))
			Expect(formData.Date).To(Equal(""))
			Expect(formData.Company).To(Equal(""))
			Expect(formData.Project).To(Equal(""))
		})

		It("should return FormData with values from existing event", func() {
			testEvent := display.EventFromDomain(fixtures.EventWith("ev-1", "Built scalable API platform", "TechCorp", "API Platform"))
			view = event.NewForm(testEvent, types.StrategyManual)
			formData := view.GetFormData()
			Expect(formData.Text).To(Equal("Built scalable API platform"))
			Expect(formData.Company).To(Equal("TechCorp"))
			Expect(formData.Project).To(Equal("API Platform"))
		})
	})

	Describe("SetFooter", func() {
		BeforeEach(func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
		})

		It("should update the help text returned by HelpText()", func() {
			view.SetFooter("Custom footer")
			Expect(view.HelpText()).To(Equal("Custom footer"))
		})

		It("should allow changing footer multiple times", func() {
			view.SetFooter("First footer")
			Expect(view.HelpText()).To(Equal("First footer"))
			view.SetFooter("Second footer")
			Expect(view.HelpText()).To(Equal("Second footer"))
		})
	})

	Describe("Update with non-key messages", func() {
		BeforeEach(func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
			view.SetTerminalInfo(120, 40)
		})

		It("should handle cursor blink message without error", func() {
			cmd, result := view.Update(nil)
			_ = cmd
			Expect(result).To(BeNil())
		})
	})

	Describe("Update with WindowSizeMsg triggers form rebuild", func() {
		It("should rebuild form with wide terminal", func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
			msg := tea.WindowSizeMsg{Width: 200, Height: 60}
			cmd, result := view.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should rebuild form with narrow terminal", func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
			msg := tea.WindowSizeMsg{Width: 20, Height: 10}
			cmd, result := view.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should rebuild form with tiny terminal below minimum", func() {
			view = event.NewForm(display.Event{}, types.StrategyQuick)
			msg := tea.WindowSizeMsg{Width: 10, Height: 5}
			cmd, result := view.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})
	})

	Describe("GetFormData with nil formData", func() {
		It("should return nil for empty view", func() {
			view = &event.Form{}
			Expect(view.GetFormData()).To(BeNil())
		})
	})

	Describe("CtrlS with nil formData", func() {
		It("should return nil result", func() {
			view = &event.Form{}
			msg := tea.KeyMsg{Type: tea.KeyCtrlS}
			_, result := view.Update(msg)
			Expect(result).To(BeNil())
		})
	})
})
