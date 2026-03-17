package intents_test

import (
	"github.com/baphled/kariya/internal/tui/intents"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// mockModal is a test double implementing ManagedModal.
type mockModal struct {
	visible      bool
	viewContent  string
	updateResult intents.ModalUpdateResult
}

func (m *mockModal) IsVisible() bool {
	return m.visible
}

func (m *mockModal) View() string {
	return m.viewContent
}

func (m *mockModal) HandleUpdate(msg tea.Msg) intents.ModalUpdateResult {
	return m.updateResult
}

var _ = Describe("ModalRegistry", func() {
	var registry *intents.ModalRegistry

	BeforeEach(func() {
		registry = intents.NewModalRegistry()
	})

	Describe("NewModalRegistry", func() {
		It("should create an empty registry", func() {
			Expect(registry).NotTo(BeNil())
			Expect(registry.HasVisibleModal()).To(BeFalse())
		})
	})

	Describe("Register", func() {
		It("should add a modal to the registry", func() {
			modal := &mockModal{visible: true, viewContent: "test"}
			registry.Register(modal)
			Expect(registry.HasVisibleModal()).To(BeTrue())
		})

		It("should not add nil modals", func() {
			registry.Register(nil)
			Expect(registry.HasVisibleModal()).To(BeFalse())
		})

		It("should maintain priority order (first registered = highest priority)", func() {
			highPriority := &mockModal{visible: true, viewContent: "high"}
			lowPriority := &mockModal{visible: true, viewContent: "low"}

			registry.Register(highPriority)
			registry.Register(lowPriority)

			// GetFirstVisible should return highest priority (first registered)
			first := registry.GetFirstVisible()
			Expect(first).To(Equal(highPriority))
		})
	})

	Describe("Clear", func() {
		It("should remove all modals from the registry", func() {
			modal := &mockModal{visible: true, viewContent: "test"}
			registry.Register(modal)
			Expect(registry.HasVisibleModal()).To(BeTrue())

			registry.Clear()
			Expect(registry.HasVisibleModal()).To(BeFalse())
		})
	})

	Describe("GetFirstVisible", func() {
		It("should return nil when no modals are registered", func() {
			Expect(registry.GetFirstVisible()).To(BeNil())
		})

		It("should return nil when no modals are visible", func() {
			modal := &mockModal{visible: false, viewContent: "hidden"}
			registry.Register(modal)
			Expect(registry.GetFirstVisible()).To(BeNil())
		})

		It("should return the first visible modal", func() {
			hidden := &mockModal{visible: false, viewContent: "hidden"}
			visible := &mockModal{visible: true, viewContent: "visible"}

			registry.Register(hidden)
			registry.Register(visible)

			first := registry.GetFirstVisible()
			Expect(first).To(Equal(visible))
		})

		It("should return highest priority visible modal", func() {
			high := &mockModal{visible: true, viewContent: "high"}
			low := &mockModal{visible: true, viewContent: "low"}

			registry.Register(high)
			registry.Register(low)

			first := registry.GetFirstVisible()
			Expect(first).To(Equal(high))
		})
	})

	Describe("HasVisibleModal", func() {
		It("should return false when registry is empty", func() {
			Expect(registry.HasVisibleModal()).To(BeFalse())
		})

		It("should return false when no modals are visible", func() {
			modal := &mockModal{visible: false}
			registry.Register(modal)
			Expect(registry.HasVisibleModal()).To(BeFalse())
		})

		It("should return true when at least one modal is visible", func() {
			hidden := &mockModal{visible: false}
			visible := &mockModal{visible: true}

			registry.Register(hidden)
			registry.Register(visible)
			Expect(registry.HasVisibleModal()).To(BeTrue())
		})
	})

	Describe("HandleUpdate", func() {
		It("should return nil when no modal is visible", func() {
			result := registry.HandleUpdate(tea.KeyMsg{})
			Expect(result).To(BeNil())
		})

		It("should update the first visible modal", func() {
			expectedResult := intents.ModalUpdateResult{
				Cmd:     nil,
				Closed:  true,
				Applied: true,
				Data:    "test-data",
			}

			modal := &mockModal{
				visible:      true,
				updateResult: expectedResult,
			}
			registry.Register(modal)

			result := registry.HandleUpdate(tea.KeyMsg{})
			Expect(result).NotTo(BeNil())
			Expect(result.Closed).To(BeTrue())
			Expect(result.Applied).To(BeTrue())
			Expect(result.Data).To(Equal("test-data"))
		})

		It("should only update the highest priority visible modal", func() {
			highResult := intents.ModalUpdateResult{Data: "high"}
			lowResult := intents.ModalUpdateResult{Data: "low"}

			high := &mockModal{visible: true, updateResult: highResult}
			low := &mockModal{visible: true, updateResult: lowResult}

			registry.Register(high)
			registry.Register(low)

			result := registry.HandleUpdate(tea.KeyMsg{})
			Expect(result.Data).To(Equal("high"))
		})
	})

	Describe("RenderOverlay", func() {
		It("should return base view when no modal is visible", func() {
			baseView := "base content"
			result := registry.RenderOverlay(baseView)
			Expect(result).To(Equal(baseView))
		})

		It("should render the first visible modal over base view", func() {
			modal := &mockModal{visible: true, viewContent: "modal content"}
			registry.Register(modal)

			baseView := "base content"
			result := registry.RenderOverlay(baseView)

			// The overlay should contain the modal content
			Expect(result).To(ContainSubstring("modal content"))
		})
	})
})

var _ = Describe("ModalUpdateResult", func() {
	Describe("Zero value", func() {
		It("should have sensible defaults", func() {
			var result intents.ModalUpdateResult
			Expect(result.Cmd).To(BeNil())
			Expect(result.Closed).To(BeFalse())
			Expect(result.Applied).To(BeFalse())
			Expect(result.Data).To(BeNil())
		})
	})
})
