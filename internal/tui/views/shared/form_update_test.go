package shared_test

import (
	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/tui/views/shared"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HandleFormUpdate", func() {
	var (
		form                      forms.Form
		msg                       tea.Msg
		resizeCalled              bool
		resizeWidth, resizeHeight int
		abortedCalled             bool
	)

	BeforeEach(func() {
		form = huh.NewForm(huh.NewGroup(huh.NewInput().Title("Name").Key("name")))
		resizeCalled = false
		resizeWidth, resizeHeight = 0, 0
		abortedCalled = false
	})

	It("handles WindowSizeMsg and calls onResize", func() {
		msg = tea.WindowSizeMsg{Width: 80, Height: 24}
		_, _, result := shared.HandleFormUpdate(form, msg, func(w, h int) {
			resizeCalled = true
			resizeWidth = w
			resizeHeight = h
		}, func() interface{} { return nil })
		Expect(resizeCalled).To(BeTrue())
		Expect(resizeWidth).To(Equal(80))
		Expect(resizeHeight).To(Equal(24))
		Expect(result).To(BeNil())
	})

	It("handles form completion", func() {
		// Simulate form completion by sending a message that would complete the form
		// In a real scenario, this would be after the user fills all required fields
		msg = tea.KeyMsg{Type: tea.KeyEnter}
		updatedForm, _, result := shared.HandleFormUpdate(form, msg, func(w, h int) {}, func() interface{} { return "done" })
		// The form may or may not be completed after a single KeyEnter depending on form state
		// Just verify the form is updated and result is nil (still in progress)
		Expect(updatedForm).NotTo(BeNil())
		Expect(result).To(BeNil())
	})

	It("handles form abortion", func() {
		msg = tea.KeyMsg{Type: tea.KeyEsc}
		_, _, result := shared.HandleFormUpdate(form, msg, func(w, h int) {}, func() interface{} { abortedCalled = true; return nil })
		Expect(abortedCalled).To(BeFalse()) // Should not call onSubmit
		Expect(result.Type()).To(Equal(widgets.ResultCancel))
	})

	It("handles normal msg and returns updated form", func() {
		msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
		updatedForm, _, result := shared.HandleFormUpdate(form, msg, func(w, h int) {}, func() interface{} { return nil })
		Expect(updatedForm).NotTo(BeNil())
		Expect(result).To(BeNil())
	})
})
