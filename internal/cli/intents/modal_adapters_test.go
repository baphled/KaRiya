package intents_test

import (
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Modal Adapters", func() {
	Describe("ErrorModalAdapter", func() {
		var (
			modal   *feedback.Modal
			adapter *intents.ErrorModalAdapter
		)

		BeforeEach(func() {
			modal = feedback.NewErrorModal("Test Error", "Something went wrong")
			adapter = intents.NewErrorModalAdapter(modal, 120, 40, nil)
		})

		Describe("IsVisible", func() {
			It("should return true when modal is not nil", func() {
				Expect(adapter.IsVisible()).To(BeTrue())
			})

			It("should return false when modal is nil", func() {
				nilAdapter := intents.NewErrorModalAdapter(nil, 120, 40, nil)
				Expect(nilAdapter.IsVisible()).To(BeFalse())
			})
		})

		Describe("View", func() {
			It("should return rendered modal content", func() {
				view := adapter.View()
				Expect(view).NotTo(BeEmpty())
				Expect(view).To(ContainSubstring("Test Error"))
			})

			It("should return empty string when modal is nil", func() {
				nilAdapter := intents.NewErrorModalAdapter(nil, 120, 40, nil)
				Expect(nilAdapter.View()).To(BeEmpty())
			})
		})

		Describe("HandleUpdate", func() {
			It("should return closed=true on Esc key", func() {
				result := adapter.HandleUpdate(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(result.Closed).To(BeTrue())
			})

			It("should return closed=false for other keys", func() {
				result := adapter.HandleUpdate(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(result.Closed).To(BeFalse())
			})

			It("should return empty result for non-key messages", func() {
				result := adapter.HandleUpdate(tea.WindowSizeMsg{Width: 100, Height: 50})
				Expect(result.Closed).To(BeFalse())
				Expect(result.Cmd).To(BeNil())
			})
		})

		Describe("RenderOverlay", func() {
			It("should render modal over base view", func() {
				baseView := "base content"
				result := adapter.RenderOverlay(baseView)
				Expect(result).To(ContainSubstring("Test Error"))
			})

			It("should return base view when modal is nil", func() {
				nilAdapter := intents.NewErrorModalAdapter(nil, 120, 40, nil)
				baseView := "base content"
				result := nilAdapter.RenderOverlay(baseView)
				Expect(result).To(Equal(baseView))
			})
		})
	})

	Describe("FormModalAdapter", func() {
		var (
			visible      bool
			viewContent  string
			updateCalled bool
			updateResult struct {
				cmd     tea.Cmd
				applied bool
				data    *testFormData
			}
			adapter *intents.FormModalAdapter[*testFormData]
		)

		BeforeEach(func() {
			visible = true
			viewContent = "form content"
			updateCalled = false
			updateResult.cmd = nil
			updateResult.applied = false
			updateResult.data = nil

			adapter = intents.NewFormModalAdapter(
				func() bool { return visible },
				func() string { return viewContent },
				func(msg tea.Msg) (tea.Cmd, bool, *testFormData) {
					updateCalled = true
					return updateResult.cmd, updateResult.applied, updateResult.data
				},
			)
		})

		Describe("IsVisible", func() {
			It("should delegate to isVisible function", func() {
				Expect(adapter.IsVisible()).To(BeTrue())

				visible = false
				Expect(adapter.IsVisible()).To(BeFalse())
			})
		})

		Describe("View", func() {
			It("should delegate to view function", func() {
				Expect(adapter.View()).To(Equal("form content"))

				viewContent = "updated content"
				Expect(adapter.View()).To(Equal("updated content"))
			})
		})

		Describe("HandleUpdate", func() {
			It("should call update function", func() {
				adapter.HandleUpdate(tea.KeyMsg{})
				Expect(updateCalled).To(BeTrue())
			})

			It("should return Applied=true when form applied", func() {
				updateResult.applied = true
				updateResult.data = &testFormData{Value: "test"}

				result := adapter.HandleUpdate(tea.KeyMsg{})
				Expect(result.Applied).To(BeTrue())
				Expect(result.Data).To(Equal(&testFormData{Value: "test"}))
			})

			It("should detect closed state from visibility", func() {
				// Modal closes after update
				visible = false
				result := adapter.HandleUpdate(tea.KeyMsg{})
				Expect(result.Closed).To(BeTrue())
			})
		})
	})

	Describe("ConfirmModalAdapter", func() {
		var (
			modal   *feedback.ConfirmModal
			adapter *intents.ConfirmModalAdapter
		)

		BeforeEach(func() {
			modal = feedback.NewConfirmModal("Confirm", "Are you sure?")
			adapter = intents.NewConfirmModalAdapter(modal)
		})

		Describe("IsVisible", func() {
			It("should return true when modal is visible", func() {
				Expect(adapter.IsVisible()).To(BeTrue())
			})

			It("should return false when modal is nil", func() {
				nilAdapter := intents.NewConfirmModalAdapter(nil)
				Expect(nilAdapter.IsVisible()).To(BeFalse())
			})

			It("should return false when modal is hidden", func() {
				modal.Hide()
				Expect(adapter.IsVisible()).To(BeFalse())
			})
		})

		Describe("View", func() {
			It("should return modal view", func() {
				view := adapter.View()
				Expect(view).To(ContainSubstring("Confirm"))
			})

			It("should return empty string when modal is nil", func() {
				nilAdapter := intents.NewConfirmModalAdapter(nil)
				Expect(nilAdapter.View()).To(BeEmpty())
			})
		})

		Describe("HandleUpdate", func() {
			It("should return confirmed=true on 'y' key", func() {
				result := adapter.HandleUpdate(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
				Expect(result.Applied).To(BeTrue())
				Expect(result.Data).To(Equal(true))
				Expect(result.Closed).To(BeTrue())
			})

			It("should return confirmed=false on 'n' key", func() {
				result := adapter.HandleUpdate(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				Expect(result.Applied).To(BeFalse())
				Expect(result.Data).To(Equal(false))
				Expect(result.Closed).To(BeTrue())
			})

			It("should return confirmed=false on Esc key", func() {
				result := adapter.HandleUpdate(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(result.Applied).To(BeFalse())
				Expect(result.Closed).To(BeTrue())
			})
		})
	})

	Describe("ViewModalAdapter", func() {
		var (
			visible      bool
			viewContent  string
			updateCalled bool
			updateCmd    tea.Cmd
			adapter      *intents.ViewModalAdapter
		)

		BeforeEach(func() {
			visible = true
			viewContent = "detail content"
			updateCalled = false
			updateCmd = nil

			adapter = intents.NewViewModalAdapter(
				func() bool { return visible },
				func() string { return viewContent },
				func(msg tea.Msg) (tea.Model, tea.Cmd) {
					updateCalled = true
					return nil, updateCmd
				},
			)
		})

		Describe("IsVisible", func() {
			It("should delegate to isVisible function", func() {
				Expect(adapter.IsVisible()).To(BeTrue())

				visible = false
				Expect(adapter.IsVisible()).To(BeFalse())
			})
		})

		Describe("View", func() {
			It("should delegate to view function", func() {
				Expect(adapter.View()).To(Equal("detail content"))
			})
		})

		Describe("HandleUpdate", func() {
			It("should call update function", func() {
				adapter.HandleUpdate(tea.KeyMsg{})
				Expect(updateCalled).To(BeTrue())
			})

			It("should detect closed state from visibility", func() {
				visible = false
				result := adapter.HandleUpdate(tea.KeyMsg{})
				Expect(result.Closed).To(BeTrue())
			})

			It("should not set Applied for view modals", func() {
				result := adapter.HandleUpdate(tea.KeyMsg{})
				Expect(result.Applied).To(BeFalse())
			})
		})
	})
})

// testFormData is a test double for form data.
type testFormData struct {
	Value string
}
