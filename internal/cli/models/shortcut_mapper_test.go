package models

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestShortcutMapper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ShortcutMapper Suite")
}

var _ = Describe("ShortcutMapper", func() {
	var mapper *ShortcutMapper

	BeforeEach(func() {
		mapper = NewShortcutMapper()
	})

	Describe("Global Shortcut Registration", func() {
		It("should register a global shortcut", func() {
			binding := key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save"))
			action := func() tea.Cmd { return nil }

			mapper.RegisterGlobalShortcut("save", binding, action)

			shortcut, exists := mapper.GetGlobalShortcut("save")
			Expect(exists).To(BeTrue())
			Expect(shortcut).NotTo(BeNil())
		})

		It("should register multiple global shortcuts", func() {
			binding1 := key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save"))
			binding2 := key.NewBinding(key.WithKeys("ctrl+o"), key.WithHelp("ctrl+o", "open"))

			mapper.RegisterGlobalShortcut("save", binding1, func() tea.Cmd { return nil })
			mapper.RegisterGlobalShortcut("open", binding2, func() tea.Cmd { return nil })

			shortcuts := mapper.GetAllGlobalShortcuts()
			Expect(shortcuts).To(HaveLen(2))
		})

		It("should unregister a global shortcut", func() {
			binding := key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save"))
			mapper.RegisterGlobalShortcut("save", binding, func() tea.Cmd { return nil })

			mapper.UnregisterGlobalShortcut("save")

			_, exists := mapper.GetGlobalShortcut("save")
			Expect(exists).To(BeFalse())
		})

		It("should replace existing shortcut", func() {
			binding1 := key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save"))
			binding2 := key.NewBinding(key.WithKeys("ctrl+shift+s"), key.WithHelp("ctrl+shift+s", "save"))

			mapper.RegisterGlobalShortcut("save", binding1, func() tea.Cmd { return nil })
			mapper.RegisterGlobalShortcut("save", binding2, func() tea.Cmd { return nil })

			shortcuts := mapper.GetAllGlobalShortcuts()
			Expect(shortcuts).To(HaveLen(1))
		})
	})

	Describe("Context-specific Shortcut Registration", func() {
		It("should register a context-specific shortcut", func() {
			binding := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			action := func() tea.Cmd { return nil }

			mapper.RegisterContextShortcut("list", "edit", binding, action)

			shortcut, exists := mapper.GetContextShortcut("list", "edit")
			Expect(exists).To(BeTrue())
			Expect(shortcut).NotTo(BeNil())
		})

		It("should register multiple context shortcuts for same context", func() {
			binding1 := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			binding2 := key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete"))

			mapper.RegisterContextShortcut("list", "edit", binding1, func() tea.Cmd { return nil })
			mapper.RegisterContextShortcut("list", "delete", binding2, func() tea.Cmd { return nil })

			shortcuts := mapper.GetContextShortcuts("list")
			Expect(shortcuts).To(HaveLen(2))
		})

		It("should unregister a context-specific shortcut", func() {
			binding := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			mapper.RegisterContextShortcut("list", "edit", binding, func() tea.Cmd { return nil })

			mapper.UnregisterContextShortcut("list", "edit")

			_, exists := mapper.GetContextShortcut("list", "edit")
			Expect(exists).To(BeFalse())
		})
	})

	Describe("Shortcut Conflict Detection", func() {
		It("should detect conflicting key bindings", func() {
			binding := key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit"))

			mapper.RegisterGlobalShortcut("save", binding, func() tea.Cmd { return nil })

			binding2 := key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit"))
			conflicts := mapper.CheckConflicts("quit", binding2)

			Expect(conflicts).NotTo(BeEmpty())
		})

		It("should return no conflicts for unique bindings", func() {
			binding1 := key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save"))
			binding2 := key.NewBinding(key.WithKeys("ctrl+o"), key.WithHelp("ctrl+o", "open"))

			mapper.RegisterGlobalShortcut("save", binding1, func() tea.Cmd { return nil })

			conflicts := mapper.CheckConflicts("open", binding2)
			Expect(conflicts).To(BeEmpty())
		})
	})

	Describe("Shortcut Lookup", func() {
		It("should get all global shortcuts", func() {
			binding1 := key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save"))
			binding2 := key.NewBinding(key.WithKeys("ctrl+o"), key.WithHelp("ctrl+o", "open"))

			mapper.RegisterGlobalShortcut("save", binding1, func() tea.Cmd { return nil })
			mapper.RegisterGlobalShortcut("open", binding2, func() tea.Cmd { return nil })

			shortcuts := mapper.GetAllGlobalShortcuts()
			Expect(shortcuts).To(HaveLen(2))
			Expect(shortcuts["save"]).NotTo(BeNil())
			Expect(shortcuts["open"]).NotTo(BeNil())
		})

		It("should get all context shortcuts for a context", func() {
			binding1 := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			binding2 := key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete"))

			mapper.RegisterContextShortcut("list", "edit", binding1, func() tea.Cmd { return nil })
			mapper.RegisterContextShortcut("list", "delete", binding2, func() tea.Cmd { return nil })

			shortcuts := mapper.GetContextShortcuts("list")
			Expect(shortcuts).To(HaveLen(2))
		})

		It("should return empty map for non-existent context", func() {
			shortcuts := mapper.GetContextShortcuts("nonexistent")
			Expect(shortcuts).To(BeEmpty())
		})
	})

	Describe("Shortcut Merging", func() {
		It("should merge global and context shortcuts", func() {
			globalBinding := key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save"))
			contextBinding := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))

			mapper.RegisterGlobalShortcut("save", globalBinding, func() tea.Cmd { return nil })
			mapper.RegisterContextShortcut("list", "edit", contextBinding, func() tea.Cmd { return nil })

			merged := mapper.GetMergedShortcuts("list")
			Expect(merged).To(HaveLen(2))
		})

		It("should prioritize context shortcuts over global ones", func() {
			globalBinding := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "global-edit"))
			contextBinding := key.NewBinding(key.WithKeys("E"), key.WithHelp("E", "context-edit"))

			mapper.RegisterGlobalShortcut("edit", globalBinding, func() tea.Cmd { return nil })
			mapper.RegisterContextShortcut("list", "edit", contextBinding, func() tea.Cmd { return nil })

			merged := mapper.GetMergedShortcuts("list")
			Expect(merged).To(HaveLen(1))
			Expect(merged).To(HaveKey("edit"))
		})
	})

	Describe("Action Invocation", func() {
		It("should invoke registered action", func() {
			actionCalled := false
			binding := key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save"))
			action := func() tea.Cmd {
				actionCalled = true
				return nil
			}

			mapper.RegisterGlobalShortcut("save", binding, action)
			_, exists := mapper.InvokeShortcut("save")

			Expect(exists).To(BeTrue())
			Expect(actionCalled).To(BeTrue())
		})

		It("should return false for non-existent shortcut", func() {
			_, exists := mapper.InvokeShortcut("nonexistent")
			Expect(exists).To(BeFalse())
		})

		It("should invoke context shortcut action", func() {
			actionCalled := false
			binding := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			action := func() tea.Cmd {
				actionCalled = true
				return nil
			}

			mapper.RegisterContextShortcut("list", "edit", binding, action)
			_, exists := mapper.InvokeContextShortcut("list", "edit")

			Expect(exists).To(BeTrue())
			Expect(actionCalled).To(BeTrue())
		})
	})

	Describe("Shortcut Clearing", func() {
		It("should clear all global shortcuts", func() {
			binding := key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save"))
			mapper.RegisterGlobalShortcut("save", binding, func() tea.Cmd { return nil })

			mapper.ClearGlobalShortcuts()

			shortcuts := mapper.GetAllGlobalShortcuts()
			Expect(shortcuts).To(BeEmpty())
		})

		It("should clear context shortcuts for specific context", func() {
			binding1 := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			binding2 := key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "submit"))

			mapper.RegisterContextShortcut("list", "edit", binding1, func() tea.Cmd { return nil })
			mapper.RegisterContextShortcut("form", "submit", binding2, func() tea.Cmd { return nil })

			mapper.ClearContextShortcuts("list")

			listShortcuts := mapper.GetContextShortcuts("list")
			formShortcuts := mapper.GetContextShortcuts("form")

			Expect(listShortcuts).To(BeEmpty())
			Expect(formShortcuts).NotTo(BeEmpty())
		})

		It("should clear all shortcuts", func() {
			globalBinding := key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save"))
			contextBinding := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))

			mapper.RegisterGlobalShortcut("save", globalBinding, func() tea.Cmd { return nil })
			mapper.RegisterContextShortcut("list", "edit", contextBinding, func() tea.Cmd { return nil })

			mapper.Clear()

			Expect(mapper.GetAllGlobalShortcuts()).To(BeEmpty())
			Expect(mapper.GetContextShortcuts("list")).To(BeEmpty())
		})
	})
})
