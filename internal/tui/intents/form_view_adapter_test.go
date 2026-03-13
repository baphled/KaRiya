package intents_test

import (
	"github.com/baphled/kariya/internal/tui/intents"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockView struct {
	content string
	result  widgets.ViewResult
	cmd     tea.Cmd
}

func (m *mockView) RenderContent() string                          { return m.content }
func (m *mockView) HelpText() string                               { return "" }
func (m *mockView) Init() tea.Cmd                                  { return nil }
func (m *mockView) Update(_ tea.Msg) (tea.Cmd, widgets.ViewResult) { return m.cmd, m.result }

var _ = Describe("FormViewAdapter", func() {
	var (
		view    *mockView
		adapter *intents.FormViewAdapter
	)

	BeforeEach(func() {
		view = &mockView{content: "form content"}
		adapter = intents.NewFormViewAdapter(view)
	})

	Describe("implements ManagedModal", func() {
		It("returns false for IsVisible by default", func() {
			Expect(adapter.IsVisible()).To(BeFalse())
		})

		It("returns true for IsVisible after Show", func() {
			adapter.Show()
			Expect(adapter.IsVisible()).To(BeTrue())
		})

		It("returns false for IsVisible after Hide", func() {
			adapter.Show()
			adapter.Hide()
			Expect(adapter.IsVisible()).To(BeFalse())
		})

		It("stores dimensions via SetDimensions", func() {
			adapter.SetDimensions(120, 40)
			Expect(adapter.IsVisible()).To(BeFalse())
		})
	})

	Describe("View", func() {
		It("delegates to wrapped view's RenderContent", func() {
			Expect(adapter.View()).To(Equal("form content"))
		})

		It("returns empty string when view is nil", func() {
			nilAdapter := intents.NewFormViewAdapter(nil)
			Expect(nilAdapter.View()).To(BeEmpty())
		})
	})

	Describe("HandleUpdate", func() {
		It("returns ModalUpdateResult with Cmd when view returns nil result", func() {
			testCmd := func() tea.Msg { return nil }
			view.cmd = testCmd
			view.result = nil

			result := adapter.HandleUpdate(tea.KeyMsg{})
			Expect(result.Cmd).NotTo(BeNil())
			Expect(result.Closed).To(BeFalse())
			Expect(result.Applied).To(BeFalse())
		})

		It("maps ResultSubmit to Applied=true, Closed=true, hides adapter", func() {
			adapter.Show()
			view.result = &widgets.SubmitViewResult{FormData: "submitted-data"}

			result := adapter.HandleUpdate(tea.KeyMsg{})
			Expect(result.Applied).To(BeTrue())
			Expect(result.Closed).To(BeTrue())
			Expect(result.Data).To(Equal("submitted-data"))
			Expect(adapter.IsVisible()).To(BeFalse())
		})

		It("maps ResultCancel to Closed=true, Applied=false, hides adapter", func() {
			adapter.Show()
			view.result = &widgets.CancelViewResult{}

			result := adapter.HandleUpdate(tea.KeyMsg{})
			Expect(result.Closed).To(BeTrue())
			Expect(result.Applied).To(BeFalse())
			Expect(adapter.IsVisible()).To(BeFalse())
		})

		It("maps ResultNavigate to not-closed, passes Data", func() {
			adapter.Show()
			view.result = &widgets.NavigateViewResult{ResultData: "nav-target"}

			result := adapter.HandleUpdate(tea.KeyMsg{})
			Expect(result.Closed).To(BeFalse())
			Expect(result.Applied).To(BeFalse())
			Expect(result.Data).To(Equal("nav-target"))
			Expect(adapter.IsVisible()).To(BeTrue())
		})

		It("maps ResultError to not-closed, passes Data", func() {
			adapter.Show()
			view.result = &widgets.ErrorViewResult{Message: "something failed"}

			result := adapter.HandleUpdate(tea.KeyMsg{})
			Expect(result.Closed).To(BeFalse())
			Expect(result.Applied).To(BeFalse())
			Expect(result.Data).NotTo(BeNil())
			Expect(adapter.IsVisible()).To(BeTrue())
		})

		It("returns empty result when view is nil", func() {
			nilAdapter := intents.NewFormViewAdapter(nil)
			result := nilAdapter.HandleUpdate(tea.KeyMsg{})
			Expect(result.Cmd).To(BeNil())
			Expect(result.Closed).To(BeFalse())
			Expect(result.Applied).To(BeFalse())
		})
	})

	Describe("visibility after result", func() {
		BeforeEach(func() {
			adapter.Show()
		})

		It("is hidden after Submit", func() {
			view.result = &widgets.SubmitViewResult{FormData: "data"}
			adapter.HandleUpdate(tea.KeyMsg{})
			Expect(adapter.IsVisible()).To(BeFalse())
		})

		It("is hidden after Cancel", func() {
			view.result = &widgets.CancelViewResult{}
			adapter.HandleUpdate(tea.KeyMsg{})
			Expect(adapter.IsVisible()).To(BeFalse())
		})

		It("is NOT hidden after Navigate", func() {
			view.result = &widgets.NavigateViewResult{ResultData: "target"}
			adapter.HandleUpdate(tea.KeyMsg{})
			Expect(adapter.IsVisible()).To(BeTrue())
		})

		It("is NOT hidden after Error", func() {
			view.result = &widgets.ErrorViewResult{Message: "err"}
			adapter.HandleUpdate(tea.KeyMsg{})
			Expect(adapter.IsVisible()).To(BeTrue())
		})
	})
})
