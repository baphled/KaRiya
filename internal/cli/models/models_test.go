package models

import (
	"context"
	"errors"
	"testing"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// testContextKey is a custom type for context keys to avoid collisions
type testContextKey string

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CLI Models Suite")
}

var _ = Describe("BaseStandardModel", func() {
	var model *BaseStandardModel

	BeforeEach(func() {
		model = NewBaseStandardModel()
	})

	Describe("Context Management", func() {
		It("should initialize with a background context", func() {
			Expect(model.GetContext()).NotTo(BeNil())
			Expect(model.GetContext()).To(Equal(context.Background()))
		})

		It("should set and retrieve context", func() {
			ctx := context.WithValue(context.Background(), testContextKey("test"), "value")
			model.SetContext(ctx)
			Expect(model.GetContext()).To(Equal(ctx))
			Expect(model.GetContext().Value(testContextKey("test"))).To(Equal("value"))
		})

		It("should initialize context metadata with empty data map", func() {
			metadata := model.GetContextMetadata()
			Expect(metadata).NotTo(BeNil())
			Expect(metadata.Data).NotTo(BeNil())
			Expect(len(metadata.Data)).To(Equal(0))
		})

		It("should set and retrieve context metadata", func() {
			metadata := &ContextMetadata{
				ScreenID: "test-screen",
				Data: map[string]interface{}{
					"key": "value",
				},
			}
			model.SetContextMetadata(metadata)
			retrieved := model.GetContextMetadata()
			Expect(retrieved.ScreenID).To(Equal("test-screen"))
			Expect(retrieved.Data["key"]).To(Equal("value"))
		})
	})

	Describe("Breadcrumb Management", func() {
		It("should start with empty breadcrumbs", func() {
			Expect(model.GetBreadcrumbs()).To(HaveLen(0))
		})

		It("should add breadcrumbs", func() {
			item1 := BreadcrumbItem{Label: "Home", ID: "home"}
			item2 := BreadcrumbItem{Label: "Events", ID: "events"}

			model.AddBreadcrumb(item1)
			model.AddBreadcrumb(item2)

			breadcrumbs := model.GetBreadcrumbs()
			Expect(breadcrumbs).To(HaveLen(2))
			Expect(breadcrumbs[0].Label).To(Equal("Home"))
			Expect(breadcrumbs[1].Label).To(Equal("Events"))
		})

		It("should pop breadcrumbs", func() {
			item := BreadcrumbItem{Label: "Test", ID: "test"}
			model.AddBreadcrumb(item)
			popped := model.PopBreadcrumb()

			Expect(popped.Label).To(Equal("Test"))
			Expect(model.GetBreadcrumbs()).To(HaveLen(0))
		})

		It("should return empty breadcrumb when popping from empty list", func() {
			popped := model.PopBreadcrumb()
			Expect(popped.Label).To(Equal(""))
			Expect(popped.ID).To(Equal(""))
		})

		It("should peek at the last breadcrumb without removing it", func() {
			item := BreadcrumbItem{Label: "Test", ID: "test"}
			model.AddBreadcrumb(item)

			peeked := model.PeekBreadcrumb()
			Expect(peeked).NotTo(BeNil())
			Expect(peeked.Label).To(Equal("Test"))
			Expect(model.GetBreadcrumbs()).To(HaveLen(1))
		})

		It("should return nil when peeking empty breadcrumbs", func() {
			peeked := model.PeekBreadcrumb()
			Expect(peeked).To(BeNil())
		})

		It("should generate correct breadcrumb path", func() {
			model.AddBreadcrumb(BreadcrumbItem{Label: "Home", ID: "home"})
			model.AddBreadcrumb(BreadcrumbItem{Label: "Events", ID: "events"})
			model.AddBreadcrumb(BreadcrumbItem{Label: "Details", ID: "details"})

			path := model.GetBreadcrumbPath()
			Expect(path).To(Equal("/Home/Events/Details/"))
		})

		It("should return root path for empty breadcrumbs", func() {
			path := model.GetBreadcrumbPath()
			Expect(path).To(Equal("/"))
		})

		It("should navigate to specific breadcrumb", func() {
			model.AddBreadcrumb(BreadcrumbItem{Label: "Home", ID: "home"})
			model.AddBreadcrumb(BreadcrumbItem{Label: "Events", ID: "events"})
			model.AddBreadcrumb(BreadcrumbItem{Label: "Details", ID: "details"})

			err := model.NavigateToBreadcrumb(1)
			Expect(err).NotTo(HaveOccurred())
			Expect(model.GetBreadcrumbs()).To(HaveLen(2))
			Expect(model.GetBreadcrumbs()[1].Label).To(Equal("Events"))
		})

		It("should return error for invalid breadcrumb index", func() {
			model.AddBreadcrumb(BreadcrumbItem{Label: "Home", ID: "home"})

			err := model.NavigateToBreadcrumb(5)
			Expect(err).To(Equal(ErrInvalidBreadcrumbIndex))
		})

		It("should clear all breadcrumbs", func() {
			model.AddBreadcrumb(BreadcrumbItem{Label: "Home", ID: "home"})
			model.AddBreadcrumb(BreadcrumbItem{Label: "Events", ID: "events"})

			model.ClearBreadcrumbs()
			Expect(model.GetBreadcrumbs()).To(HaveLen(0))
		})

		It("should automatically initialize metadata for breadcrumb items", func() {
			item := BreadcrumbItem{Label: "Test", ID: "test"}
			model.AddBreadcrumb(item)

			breadcrumb := model.GetBreadcrumbs()[0]
			Expect(breadcrumb.Metadata).NotTo(BeNil())
		})
	})

	Describe("Navigation History", func() {
		It("should start with empty navigation history", func() {
			Expect(model.GetNavigationHistory()).To(HaveLen(0))
		})

		It("should push to navigation history", func() {
			model.PushNavigationHistory("screen1", "state1")
			model.PushNavigationHistory("screen2", "state2")

			history := model.GetNavigationHistory()
			Expect(history).To(HaveLen(2))
			Expect(history[0]).To(Equal("screen1"))
			Expect(history[1]).To(Equal("screen2"))
		})

		It("should pop from navigation history", func() {
			model.PushNavigationHistory("screen1", "state1")
			model.PushNavigationHistory("screen2", "state2")

			label, state, err := model.PopNavigationHistory()
			Expect(err).NotTo(HaveOccurred())
			Expect(label).To(Equal("screen2"))
			Expect(state).To(Equal("state2"))
			Expect(model.GetNavigationHistory()).To(HaveLen(1))
		})

		It("should return error when popping empty history", func() {
			_, _, err := model.PopNavigationHistory()
			Expect(err).To(Equal(ErrEmptyNavigationHistory))
		})
	})

	Describe("Shortcut Management", func() {
		It("should register and retrieve shortcuts", func() {
			shortcuts := map[string]key.Binding{
				"quit": key.NewBinding(
					key.WithKeys("ctrl+c"),
					key.WithHelp("ctrl+c", "quit"),
				),
			}

			model.RegisterShortcuts(shortcuts)
			retrieved := model.GetShortcuts()

			Expect(retrieved).To(HaveLen(1))
			Expect(retrieved["quit"]).NotTo(BeNil())
		})

		It("should handle shortcuts with default implementation", func() {
			binding := key.NewBinding(key.WithKeys("ctrl+c"))
			result, cmd := model.HandleShortcut(binding)

			Expect(result).To(Equal(model))
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Error Handling", func() {
		It("should start with no error", func() {
			Expect(model.GetLastError()).To(BeNil())
		})

		It("should set and retrieve errors", func() {
			err := errors.New("test error")
			model.SetError(err)

			Expect(model.GetLastError()).To(Equal(err))
		})

		It("should clear errors", func() {
			model.SetError(errors.New("test error"))
			model.ClearError()

			Expect(model.GetLastError()).To(BeNil())
		})
	})

	Describe("State Management", func() {
		It("should set and retrieve state", func() {
			state := map[string]interface{}{"key": "value"}
			model.SetState(state)

			retrieved := model.GetState()
			Expect(retrieved).To(Equal(state))
		})

		It("should handle nil state", func() {
			model.SetState(nil)
			Expect(model.GetState()).To(BeNil())
		})
	})

	Describe("Validation and Reset", func() {
		It("should validate with default no-op implementation", func() {
			err := model.Validate()
			Expect(err).To(BeNil())
		})

		It("should reset all state", func() {
			model.AddBreadcrumb(BreadcrumbItem{Label: "Test", ID: "test"})
			model.PushNavigationHistory("screen", "state")
			model.SetError(errors.New("test"))
			model.SetState("test-state")

			model.Reset()

			Expect(model.GetBreadcrumbs()).To(HaveLen(0))
			Expect(model.GetNavigationHistory()).To(HaveLen(0))
			Expect(model.GetLastError()).To(BeNil())
			Expect(model.GetState()).To(BeNil())
			Expect(model.GetContextMetadata().Data).NotTo(BeNil())
		})
	})
})

var _ = Describe("ShortcutHandler", func() {
	var handler *ShortcutHandler

	BeforeEach(func() {
		handler = NewShortcutHandler()
	})

	Describe("Shortcut Registration", func() {
		It("should register a shortcut with action", func() {
			binding := key.NewBinding(key.WithKeys("ctrl+c"))
			action := func() tea.Cmd { return nil }

			handler.RegisterShortcut("quit", binding, action)

			shortcut, exists := handler.GetShortcut("quit")
			Expect(exists).To(BeTrue())
			Expect(shortcut).NotTo(BeNil())
		})

		It("should unregister a shortcut", func() {
			binding := key.NewBinding(key.WithKeys("ctrl+c"))
			action := func() tea.Cmd { return nil }

			handler.RegisterShortcut("quit", binding, action)
			handler.UnregisterShortcut("quit")

			_, exists := handler.GetShortcut("quit")
			Expect(exists).To(BeFalse())
		})

		It("should retrieve non-existent shortcut as not found", func() {
			_, exists := handler.GetShortcut("nonexistent")
			Expect(exists).To(BeFalse())
		})
	})

	Describe("Shortcut Retrieval", func() {
		It("should get all shortcuts", func() {
			quit := key.NewBinding(key.WithKeys("ctrl+c"))
			back := key.NewBinding(key.WithKeys("esc"))

			handler.RegisterShortcut("quit", quit, func() tea.Cmd { return nil })
			handler.RegisterShortcut("back", back, func() tea.Cmd { return nil })

			shortcuts := handler.GetAllShortcuts()
			Expect(shortcuts).To(HaveLen(2))
			Expect(shortcuts["quit"]).NotTo(BeNil())
			Expect(shortcuts["back"]).NotTo(BeNil())
		})
	})

	Describe("Message Handling", func() {
		It("should not handle non-matching key message", func() {
			binding := key.NewBinding(key.WithKeys("ctrl+c"))
			handler.RegisterShortcut("quit", binding, func() tea.Cmd { return nil })

			msg := tea.KeyMsg{Type: tea.KeyEscape}
			_, handled := handler.HandleMsg(msg)

			Expect(handled).To(BeFalse())
		})

		It("should ignore non-key messages", func() {
			msg := tea.WindowSizeMsg{Width: 100, Height: 50}
			_, handled := handler.HandleMsg(msg)

			Expect(handled).To(BeFalse())
		})

		It("should register action and not error on handle", func() {
			binding := key.NewBinding(key.WithKeys("ctrl+c"))
			action := func() tea.Cmd { return nil }

			handler.RegisterShortcut("quit", binding, action)
			Expect(handler.GetAllShortcuts()).To(HaveLen(1))
		})
	})

	Describe("Shortcut Clearing", func() {
		It("should clear all shortcuts", func() {
			handler.RegisterShortcut("quit", key.NewBinding(key.WithKeys("ctrl+c")), func() tea.Cmd { return nil })
			handler.RegisterShortcut("back", key.NewBinding(key.WithKeys("esc")), func() tea.Cmd { return nil })

			handler.ClearShortcuts()

			Expect(handler.GetAllShortcuts()).To(HaveLen(0))
		})
	})
})

var _ = Describe("CommonShortcuts", func() {
	It("should create common shortcuts with default bindings", func() {
		common := NewCommonShortcuts()

		Expect(common.Quit).NotTo(BeNil())
		Expect(common.Back).NotTo(BeNil())
		Expect(common.Help).NotTo(BeNil())
		Expect(common.Enter).NotTo(BeNil())
		Expect(common.Up).NotTo(BeNil())
		Expect(common.Down).NotTo(BeNil())
		Expect(common.PageUp).NotTo(BeNil())
		Expect(common.PageDown).NotTo(BeNil())
	})
})

var _ = Describe("ErrorHandler", func() {
	var handler *ErrorHandler

	BeforeEach(func() {
		handler = NewErrorHandler()
		Expect(handler).NotTo(BeNil())
		Expect(handler.logger).NotTo(BeNil())
	})

	Describe("Error Logging", func() {
		It("should log an error with severity", func() {
			err := errors.New("test error")
			handler.LogError(SeverityError, "test message", err, "screen1")

			Expect(handler.GetErrorHistory()).To(HaveLen(1))
			lastError := handler.GetLastError()
			Expect(lastError).NotTo(BeNil())
			Expect(lastError.Message).To(Equal("test message"))
			Expect(lastError.Severity).To(Equal(SeverityError))
		})

		It("should maintain error timestamp", func() {
			err := errors.New("test error")
			handler.LogError(SeverityWarning, "warning message", err, "screen1")

			lastError := handler.GetLastError()
			Expect(lastError.Timestamp).NotTo(BeZero())
		})

		It("should log multiple errors in order", func() {
			handler.LogError(SeverityInfo, "info", nil, "screen1")
			handler.LogError(SeverityWarning, "warning", errors.New("warn"), "screen1")
			handler.LogError(SeverityError, "error", errors.New("err"), "screen1")

			history := handler.GetErrorHistory()
			Expect(history).To(HaveLen(3))
			Expect(history[0].Severity).To(Equal(SeverityInfo))
			Expect(history[1].Severity).To(Equal(SeverityWarning))
			Expect(history[2].Severity).To(Equal(SeverityError))
		})
	})

	Describe("Error Retrieval", func() {
		It("should get last error", func() {
			handler.LogError(SeverityInfo, "first", nil, "screen1")
			handler.LogError(SeverityError, "second", errors.New("err"), "screen1")

			lastError := handler.GetLastError()
			Expect(lastError.Message).To(Equal("second"))
		})

		It("should return nil for last error when no errors exist", func() {
			lastError := handler.GetLastError()
			Expect(lastError).To(BeNil())
		})

		It("should filter errors by screen ID", func() {
			handler.LogError(SeverityError, "screen1 error", errors.New("err1"), "screen1")
			handler.LogError(SeverityError, "screen2 error", errors.New("err2"), "screen2")
			handler.LogError(SeverityError, "screen1 error2", errors.New("err3"), "screen1")

			screen1Errors := handler.GetErrorsByScreenID("screen1")
			Expect(screen1Errors).To(HaveLen(2))
			Expect(screen1Errors[0].Message).To(Equal("screen1 error"))
			Expect(screen1Errors[1].Message).To(Equal("screen1 error2"))
		})

		It("should filter errors by severity", func() {
			handler.LogError(SeverityInfo, "info", nil, "screen1")
			handler.LogError(SeverityWarning, "warning", nil, "screen1")
			handler.LogError(SeverityError, "error", errors.New("err"), "screen1")
			handler.LogError(SeverityWarning, "warning2", nil, "screen1")

			warnings := handler.GetErrorsBySeverity(SeverityWarning)
			Expect(warnings).To(HaveLen(2))
		})
	})

	Describe("Error Clearing", func() {
		It("should clear all errors", func() {
			handler.LogError(SeverityError, "error1", errors.New("err1"), "screen1")
			handler.LogError(SeverityError, "error2", errors.New("err2"), "screen1")

			handler.ClearErrors()

			Expect(handler.GetErrorHistory()).To(HaveLen(0))
			Expect(handler.GetLastError()).To(BeNil())
		})

		It("should clear errors by screen ID", func() {
			handler.LogError(SeverityError, "screen1 error", errors.New("err1"), "screen1")
			handler.LogError(SeverityError, "screen2 error", errors.New("err2"), "screen2")
			handler.LogError(SeverityError, "screen1 error2", errors.New("err3"), "screen1")

			handler.ClearErrorsByScreenID("screen1")

			remaining := handler.GetErrorHistory()
			Expect(remaining).To(HaveLen(1))
			Expect(remaining[0].ScreenID).To(Equal("screen2"))
		})
	})

	Describe("Suggestion Management", func() {
		It("should register and retrieve suggestions", func() {
			err := errors.New("database connection failed")
			suggestion := "Check your database connection settings"

			handler.RegisterSuggestion(err.Error(), suggestion)
			retrieved := handler.GetSuggestion(err)

			Expect(retrieved).To(Equal(suggestion))
		})

		It("should return empty string for unknown error suggestion", func() {
			err := errors.New("unknown error")
			suggestion := handler.GetSuggestion(err)

			Expect(suggestion).To(Equal(""))
		})

		It("should attach suggestion to logged error", func() {
			err := errors.New("connection timeout")
			suggestion := "Try increasing the timeout value"

			handler.RegisterSuggestion(err.Error(), suggestion)
			handler.LogError(SeverityError, "connection timeout", err, "screen1")

			lastError := handler.GetLastError()
			Expect(lastError.Suggestion).To(Equal(suggestion))
		})
	})

	Describe("Error Limiting", func() {
		It("should limit stored errors to maxErrors", func() {
			// Log more than maxErrors (default 100)
			for i := 0; i < 150; i++ {
				handler.LogError(SeverityInfo, "test", nil, "screen1")
			}

			Expect(handler.GetErrorHistory()).To(HaveLen(100))
		})
	})
})

var _ = Describe("ErrorSeverity String Representation", func() {
	It("should return correct string for each severity level", func() {
		Expect(SeverityInfo.String()).To(Equal("INFO"))
		Expect(SeverityWarning.String()).To(Equal("WARNING"))
		Expect(SeverityError.String()).To(Equal("ERROR"))
		Expect(SeverityCritical.String()).To(Equal("CRITICAL"))
	})

	It("should handle unknown severity level", func() {
		unknownSeverity := ErrorSeverity(999)
		Expect(unknownSeverity.String()).To(Equal("UNKNOWN"))
	})
})

var _ = Describe("Error Message Formatting", func() {
	It("should format error message with severity", func() {
		message := FormatErrorMessage(SeverityError, "An error occurred", "Try again")

		Expect(message).To(ContainSubstring("[ERROR]"))
		Expect(message).To(ContainSubstring("An error occurred"))
		Expect(message).To(ContainSubstring("Suggestion"))
		Expect(message).To(ContainSubstring("Try again"))
	})

	It("should format error message without suggestion", func() {
		message := FormatErrorMessage(SeverityWarning, "A warning", "")

		Expect(message).To(ContainSubstring("[WARNING]"))
		Expect(message).To(ContainSubstring("A warning"))
		Expect(message).NotTo(ContainSubstring("Suggestion"))
	})
})

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

var _ = Describe("ContextShortcutHandler", func() {
	var handler *ContextShortcutHandler

	BeforeEach(func() {
		handler = NewContextShortcutHandler()
	})

	Describe("Context Switching", func() {
		It("should set current context", func() {
			handler.SetCurrentContext("list")
			Expect(handler.GetCurrentContext()).To(Equal("list"))
		})

		It("should default to empty context", func() {
			Expect(handler.GetCurrentContext()).To(Equal(""))
		})

		It("should switch between contexts", func() {
			handler.SetCurrentContext("list")
			Expect(handler.GetCurrentContext()).To(Equal("list"))

			handler.SetCurrentContext("form")
			Expect(handler.GetCurrentContext()).To(Equal("form"))
		})
	})

	Describe("Context Stack Management", func() {
		It("should push context to stack", func() {
			handler.PushContext("list")
			Expect(handler.GetCurrentContext()).To(Equal("list"))
			Expect(handler.GetContextStack()).To(HaveLen(1))
		})

		It("should pop context from stack", func() {
			handler.PushContext("list")
			handler.PushContext("form")

			popped := handler.PopContext()
			Expect(popped).To(Equal("form"))
			Expect(handler.GetCurrentContext()).To(Equal("list"))
		})

		It("should return to previous context after pop", func() {
			handler.PushContext("list")
			handler.PushContext("form")
			handler.PushContext("modal")

			handler.PopContext()
			Expect(handler.GetCurrentContext()).To(Equal("form"))

			handler.PopContext()
			Expect(handler.GetCurrentContext()).To(Equal("list"))
		})

		It("should handle pop on empty stack", func() {
			popped := handler.PopContext()
			Expect(popped).To(Equal(""))
		})
	})

	Describe("Context-aware Shortcut Resolution", func() {
		It("should resolve shortcut in current context", func() {
			mapper := NewShortcutMapper()
			binding := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			mapper.RegisterContextShortcut("list", "edit", binding, func() tea.Cmd { return nil })

			handler.SetMapper(mapper)
			handler.SetCurrentContext("list")

			shortcuts := handler.GetActiveShortcuts()
			Expect(shortcuts).To(HaveKey("edit"))
		})

		It("should merge global and context shortcuts", func() {
			mapper := NewShortcutMapper()
			globalBinding := key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save"))
			contextBinding := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))

			mapper.RegisterGlobalShortcut("save", globalBinding, func() tea.Cmd { return nil })
			mapper.RegisterContextShortcut("list", "edit", contextBinding, func() tea.Cmd { return nil })

			handler.SetMapper(mapper)
			handler.SetCurrentContext("list")

			shortcuts := handler.GetActiveShortcuts()
			Expect(shortcuts).To(HaveLen(2))
			Expect(shortcuts).To(HaveKey("save"))
			Expect(shortcuts).To(HaveKey("edit"))
		})

		It("should prioritize context shortcuts over global", func() {
			mapper := NewShortcutMapper()
			globalBinding := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "expand"))
			contextBinding := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))

			mapper.RegisterGlobalShortcut("edit", globalBinding, func() tea.Cmd { return nil })
			mapper.RegisterContextShortcut("list", "edit", contextBinding, func() tea.Cmd { return nil })

			handler.SetMapper(mapper)
			handler.SetCurrentContext("list")

			shortcuts := handler.GetActiveShortcuts()
			Expect(shortcuts).To(HaveLen(1))
			Expect(shortcuts).To(HaveKey("edit"))
		})
	})

	Describe("Global Shortcut Availability", func() {
		It("should always include global shortcuts", func() {
			mapper := NewShortcutMapper()
			binding := key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit"))
			mapper.RegisterGlobalShortcut("quit", binding, func() tea.Cmd { return nil })

			handler.SetMapper(mapper)

			shortcuts := handler.GetActiveShortcuts()
			Expect(shortcuts).To(HaveKey("quit"))
		})

		It("should have global shortcuts in any context", func() {
			mapper := NewShortcutMapper()
			binding := key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit"))
			mapper.RegisterGlobalShortcut("quit", binding, func() tea.Cmd { return nil })

			handler.SetMapper(mapper)
			handler.SetCurrentContext("list")

			shortcuts := handler.GetActiveShortcuts()
			Expect(shortcuts).To(HaveKey("quit"))

			handler.SetCurrentContext("form")
			shortcuts = handler.GetActiveShortcuts()
			Expect(shortcuts).To(HaveKey("quit"))
		})
	})

	Describe("Shortcut Lookup", func() {
		It("should look up shortcut in current context", func() {
			mapper := NewShortcutMapper()
			binding := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			mapper.RegisterContextShortcut("list", "edit", binding, func() tea.Cmd { return nil })

			handler.SetMapper(mapper)
			handler.SetCurrentContext("list")

			shortcut, exists := handler.LookupShortcut("edit")
			Expect(exists).To(BeTrue())
			Expect(shortcut).NotTo(BeNil())
		})

		It("should fall back to global shortcut", func() {
			mapper := NewShortcutMapper()
			binding := key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save"))
			mapper.RegisterGlobalShortcut("save", binding, func() tea.Cmd { return nil })

			handler.SetMapper(mapper)
			handler.SetCurrentContext("list")

			_, exists := handler.LookupShortcut("save")
			Expect(exists).To(BeTrue())
		})

		It("should return not found for non-existent shortcut", func() {
			mapper := NewShortcutMapper()
			handler.SetMapper(mapper)

			_, exists := handler.LookupShortcut("nonexistent")
			Expect(exists).To(BeFalse())
		})
	})
})

var _ = Describe("ShortcutHelpSystem", func() {
	var helpSystem *ShortcutHelpSystem

	BeforeEach(func() {
		helpSystem = NewShortcutHelpSystem()
	})

	Describe("Shortcut Registration for Help", func() {
		It("should register shortcut with help information", func() {
			binding := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			helpSystem.RegisterShortcut("edit", binding, "Edit the selected event", []string{"list", "form"})

			info, exists := helpSystem.GetShortcutInfo("edit")
			Expect(exists).To(BeTrue())
			Expect(info.Description).To(Equal("Edit the selected event"))
		})

		It("should include context information", func() {
			binding := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			helpSystem.RegisterShortcut("edit", binding, "Edit the selected event", []string{"list"})

			info, _ := helpSystem.GetShortcutInfo("edit")
			Expect(info.Contexts).To(ContainElement("list"))
		})

		It("should update shortcut information", func() {
			binding1 := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			helpSystem.RegisterShortcut("edit", binding1, "Old description", []string{"list"})

			binding2 := key.NewBinding(key.WithKeys("E"), key.WithHelp("E", "edit"))
			helpSystem.RegisterShortcut("edit", binding2, "New description", []string{"list", "form"})

			info, _ := helpSystem.GetShortcutInfo("edit")
			Expect(info.Description).To(Equal("New description"))
			Expect(info.Contexts).To(HaveLen(2))
		})
	})

	Describe("Shortcut Lookup", func() {
		It("should retrieve all registered shortcuts", func() {
			binding1 := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			binding2 := key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete"))

			helpSystem.RegisterShortcut("edit", binding1, "Edit item", []string{"list"})
			helpSystem.RegisterShortcut("delete", binding2, "Delete item", []string{"list"})

			shortcuts := helpSystem.GetAllShortcuts()
			Expect(shortcuts).To(HaveLen(2))
		})

		It("should get shortcuts for specific context", func() {
			binding1 := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			binding2 := key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete"))
			binding3 := key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "submit"))

			helpSystem.RegisterShortcut("edit", binding1, "Edit item", []string{"list"})
			helpSystem.RegisterShortcut("delete", binding2, "Delete item", []string{"list"})
			helpSystem.RegisterShortcut("submit", binding3, "Submit form", []string{"form"})

			listShortcuts := helpSystem.GetShortcutsForContext("list")
			Expect(listShortcuts).To(HaveLen(2))
			Expect(listShortcuts).To(HaveKey("edit"))
			Expect(listShortcuts).To(HaveKey("delete"))
		})

		It("should return empty map for unknown context", func() {
			shortcuts := helpSystem.GetShortcutsForContext("unknown")
			Expect(shortcuts).To(BeEmpty())
		})
	})

	Describe("Category Management", func() {
		It("should categorize shortcuts", func() {
			binding := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			helpSystem.RegisterShortcut("edit", binding, "Edit item", []string{"list"})
			helpSystem.CategorizeShortcut("edit", "Editing")

			category, exists := helpSystem.GetShortcutCategory("edit")
			Expect(exists).To(BeTrue())
			Expect(category).To(Equal("Editing"))
		})

		It("should get shortcuts by category", func() {
			binding1 := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			binding2 := key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete"))
			binding3 := key.NewBinding(key.WithKeys("j"), key.WithHelp("j", "down"))

			helpSystem.RegisterShortcut("edit", binding1, "Edit item", []string{"list"})
			helpSystem.RegisterShortcut("delete", binding2, "Delete item", []string{"list"})
			helpSystem.RegisterShortcut("down", binding3, "Navigate down", []string{"list"})

			helpSystem.CategorizeShortcut("edit", "Editing")
			helpSystem.CategorizeShortcut("delete", "Editing")
			helpSystem.CategorizeShortcut("down", "Navigation")

			editingShortcuts := helpSystem.GetShortcutsByCategory("Editing")
			Expect(editingShortcuts).To(HaveLen(2))

			navigationShortcuts := helpSystem.GetShortcutsByCategory("Navigation")
			Expect(navigationShortcuts).To(HaveLen(1))
		})
	})

	Describe("Help Text Generation", func() {
		It("should generate help text for shortcut", func() {
			binding := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			helpSystem.RegisterShortcut("edit", binding, "Edit the selected item", []string{"list"})

			helpText := helpSystem.GenerateHelpText("edit")
			Expect(helpText).To(ContainSubstring("e"))
			Expect(helpText).To(ContainSubstring("edit"))
			Expect(helpText).To(ContainSubstring("Edit the selected item"))
		})

		It("should generate help for context", func() {
			binding1 := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			binding2 := key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete"))

			helpSystem.RegisterShortcut("edit", binding1, "Edit item", []string{"list"})
			helpSystem.RegisterShortcut("delete", binding2, "Delete item", []string{"list"})

			helpText := helpSystem.GenerateContextHelp("list")
			Expect(helpText).NotTo(BeEmpty())
			Expect(helpText).To(ContainSubstring("edit"))
			Expect(helpText).To(ContainSubstring("delete"))
		})
	})

	Describe("Shortcut Search", func() {
		It("should search shortcuts by description", func() {
			binding1 := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			binding2 := key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete"))

			helpSystem.RegisterShortcut("edit", binding1, "Edit the selected item", []string{"list"})
			helpSystem.RegisterShortcut("delete", binding2, "Remove the selected item", []string{"list"})

			results := helpSystem.SearchShortcuts("selected")
			Expect(results).To(HaveLen(2))
		})

		It("should search shortcuts by key", func() {
			binding1 := key.NewBinding(key.WithKeys("ctrl+e"), key.WithHelp("ctrl+e", "edit"))
			binding2 := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "expand"))

			helpSystem.RegisterShortcut("edit", binding1, "Edit item", []string{"list"})
			helpSystem.RegisterShortcut("expand", binding2, "Expand item", []string{"list"})

			results := helpSystem.SearchShortcuts("ctrl+e")
			Expect(results).To(HaveLen(1))
			Expect(results[0].ID).To(Equal("edit"))
		})

		It("should return empty results for no matches", func() {
			results := helpSystem.SearchShortcuts("nonexistent")
			Expect(results).To(BeEmpty())
		})
	})

	Describe("Discovery Features", func() {
		It("should get all available contexts", func() {
			binding1 := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			binding2 := key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "submit"))

			helpSystem.RegisterShortcut("edit", binding1, "Edit item", []string{"list", "form"})
			helpSystem.RegisterShortcut("submit", binding2, "Submit form", []string{"form"})

			contexts := helpSystem.GetAllContexts()
			Expect(contexts).To(HaveLen(2))
			Expect(contexts).To(ContainElement("list"))
			Expect(contexts).To(ContainElement("form"))
		})

		It("should get all categories", func() {
			binding1 := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
			binding2 := key.NewBinding(key.WithKeys("j"), key.WithHelp("j", "down"))

			helpSystem.RegisterShortcut("edit", binding1, "Edit item", []string{"list"})
			helpSystem.RegisterShortcut("down", binding2, "Navigate down", []string{"list"})

			helpSystem.CategorizeShortcut("edit", "Editing")
			helpSystem.CategorizeShortcut("down", "Navigation")

			categories := helpSystem.GetAllCategories()
			Expect(categories).To(HaveLen(2))
		})
	})
})
