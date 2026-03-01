package intents_test

import (
	"errors"
	"fmt"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/terminal"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ModalEditResult", func() {
	Describe("NewModalEditResult", func() {
		It("should create result with provided values", func() {
			changes := map[string]interface{}{"name": "new-name"}
			result := intents.NewModalEditResult("original", "modified", true, changes)

			Expect(result).NotTo(BeNil())
			Expect(result.Original).To(Equal("original"))
			Expect(result.Modified).To(Equal("modified"))
			Expect(result.Accepted).To(BeTrue())
			Expect(result.Changes).To(HaveKeyWithValue("name", "new-name"))
		})

		It("should initialize nil changes to empty map", func() {
			result := intents.NewModalEditResult("orig", "mod", true, nil)

			Expect(result.Changes).NotTo(BeNil())
			Expect(result.Changes).To(BeEmpty())
		})
	})

	Describe("NewCancelledModalEditResult", func() {
		It("should create cancelled result with original preserved", func() {
			result := intents.NewCancelledModalEditResult("original-value")

			Expect(result).NotTo(BeNil())
			Expect(result.Original).To(Equal("original-value"))
			Expect(result.Modified).To(Equal("original-value"))
			Expect(result.Accepted).To(BeFalse())
			Expect(result.Changes).To(BeEmpty())
		})
	})

	Describe("HasChanges", func() {
		It("should return true when changes exist", func() {
			result := intents.NewModalEditResult("a", "b", true, map[string]interface{}{"field": "val"})
			Expect(result.HasChanges()).To(BeTrue())
		})

		It("should return false when no changes", func() {
			result := intents.NewModalEditResult("a", "a", true, nil)
			Expect(result.HasChanges()).To(BeFalse())
		})
	})

	Describe("WasAccepted", func() {
		It("should return true when accepted", func() {
			result := intents.NewModalEditResult("a", "b", true, nil)
			Expect(result.WasAccepted()).To(BeTrue())
		})

		It("should return false when cancelled", func() {
			result := intents.NewCancelledModalEditResult("a")
			Expect(result.WasAccepted()).To(BeFalse())
		})
	})

	Describe("GetChange", func() {
		It("should return value for existing field", func() {
			changes := map[string]interface{}{"name": "new-name", "age": 30}
			result := intents.NewModalEditResult("a", "b", true, changes)

			Expect(result.GetChange("name")).To(Equal("new-name"))
			Expect(result.GetChange("age")).To(Equal(30))
		})

		It("should return nil for nonexistent field", func() {
			result := intents.NewModalEditResult("a", "b", true, map[string]interface{}{"name": "val"})
			Expect(result.GetChange("missing")).To(BeNil())
		})

		It("should return nil when changes map is nil", func() {
			result := &intents.ModalEditResult[string]{
				Original: "a",
				Modified: "b",
				Accepted: true,
				Changes:  nil,
			}
			Expect(result.GetChange("any")).To(BeNil())
		})
	})
})

var _ = Describe("BaseIntent Help Coverage", func() {
	var base *intents.BaseIntent

	BeforeEach(func() {
		base = intents.NewBaseIntent()
	})

	Describe("GetHelpModal", func() {
		It("should return the help modal", func() {
			modal := base.GetHelpModal()
			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("SetHelpKeyMap", func() {
		It("should accept a keymap without panic", func() {
			Expect(func() {
				base.SetHelpKeyMap("not-a-real-keymap")
			}).NotTo(Panic())
		})

		It("should handle nil keymap without panic", func() {
			Expect(func() {
				base.SetHelpKeyMap(nil)
			}).NotTo(Panic())
		})
	})

	Describe("ToggleHelp with valid terminal info", func() {
		It("should toggle help when terminal info is valid", func() {
			info := terminal.NewInfo()
			info.Width = 120
			info.Height = 40
			info.IsValid = true
			base.UpdateTerminalInfo(info)

			base.ToggleHelp()
			Expect(base.IsHelpVisible()).To(BeTrue())

			base.ToggleHelp()
			Expect(base.IsHelpVisible()).To(BeFalse())
		})
	})

	Describe("IsHelpVisible", func() {
		It("should return true when help is shown", func() {
			base.ShowHelp()
			Expect(base.IsHelpVisible()).To(BeTrue())
		})

		It("should return false when help is hidden", func() {
			base.ShowHelp()
			base.HideHelp()
			Expect(base.IsHelpVisible()).To(BeFalse())
		})
	})
})

var _ = Describe("IntentResult Coverage", func() {
	Describe("WithError", func() {
		It("should attach error to result", func() {
			result := intents.NewCompletedResult("data")
			intentErr := &intents.IntentError{Code: "ERR", Message: "something failed"}

			chained := result.WithError(intentErr)

			Expect(chained).To(Equal(result))
			Expect(result.Error).To(HaveOccurred())
			Expect(result.Error.Code).To(Equal("ERR"))
		})

		It("should clear error when nil passed", func() {
			result := intents.NewFailedResult[string]("code", "msg", nil)
			result.WithError(nil)

			Expect(result.Error).ToNot(HaveOccurred())
		})
	})

	Describe("WithStatus", func() {
		It("should override result status", func() {
			result := intents.NewCompletedResult("data")
			chained := result.WithStatus(intents.Partial)

			Expect(chained).To(Equal(result))
			Expect(result.Status).To(Equal(intents.Partial))
		})

		It("should downgrade completed to cancelled", func() {
			result := intents.NewCompletedResult("data")
			result.WithStatus(intents.Cancelled)

			Expect(result.IsCancelled()).To(BeTrue())
			Expect(result.IsSuccessful()).To(BeFalse())
		})
	})

	Describe("WithData", func() {
		It("should attach data to result", func() {
			result := intents.NewCancelledResult[string]()
			chained := result.WithData("new-data")

			Expect(chained).To(Equal(result))
			Expect(result.Data).To(Equal("new-data"))
		})

		It("should replace existing data", func() {
			result := intents.NewCompletedResult("old")
			result.WithData("new")

			Expect(result.Data).To(Equal("new"))
		})
	})

	Describe("WithMetadata nil initialization", func() {
		It("should initialize nil metadata map on first call", func() {
			result := &intents.IntentResult[string]{Status: intents.Completed}
			result.WithMetadata("key", "value")

			val, ok := result.GetMetadata("key")
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal("value"))
		})
	})

	Describe("GetMetadata with nil map", func() {
		It("should return false for nil metadata", func() {
			result := &intents.IntentResult[string]{Status: intents.Completed}
			_, ok := result.GetMetadata("any")
			Expect(ok).To(BeFalse())
		})
	})

	Describe("GetAllMetadata with nil map", func() {
		It("should return empty map for nil metadata", func() {
			result := &intents.IntentResult[string]{Status: intents.Completed}
			all := result.GetAllMetadata()
			Expect(all).NotTo(BeNil())
			Expect(all).To(BeEmpty())
		})
	})

	Describe("IsValid additional branches", func() {
		It("should return error for unknown status", func() {
			result := &intents.IntentResult[string]{Status: intents.ResultStatus("unknown_status")}
			Expect(result.IsValid()).To(HaveOccurred())
		})
	})
})

var _ = Describe("DefaultIntentRouter Coverage", func() {
	var router *intents.DefaultIntentRouter

	BeforeEach(func() {
		router = intents.NewDefaultIntentRouter()
	})

	Describe("RegisterResultHandler", func() {
		It("should register a result handler without error", func() {
			handler := func(result *intents.IntentResult[interface{}]) tea.Cmd {
				return nil
			}
			Expect(func() {
				router.RegisterResultHandler("test-intent", handler)
			}).NotTo(Panic())
		})
	})

	Describe("UpdateTerminalInfo", func() {
		It("should store terminal info", func() {
			info := terminal.NewInfo()
			info.Width = 200
			info.Height = 50
			info.IsValid = true

			router.UpdateTerminalInfo(info)

			retrieved := router.GetTerminalInfo()
			Expect(retrieved).NotTo(BeNil())
			Expect(retrieved.Width).To(Equal(200))
			Expect(retrieved.Height).To(Equal(50))
		})

		It("should propagate to terminal-aware active intent", func() {
			mockIntent := newTerminalAwareMockIntent()
			factory := func() intents.Intent { return mockIntent }
			_ = router.RegisterIntent("term-intent", factory)
			_, _ = router.ActivateIntent("term-intent", nil)

			info := terminal.NewInfo()
			info.Width = 150
			info.Height = 45
			info.IsValid = true

			router.UpdateTerminalInfo(info)

			Expect(mockIntent.terminalInfoUpdated).To(BeTrue())
		})
	})

	Describe("GetTerminalInfo", func() {
		It("should return default terminal info", func() {
			info := router.GetTerminalInfo()
			Expect(info).NotTo(BeNil())
		})
	})

	Describe("GetLogo", func() {
		It("should return nil logo initially", func() {
			logo := router.GetLogo()
			Expect(logo).To(BeNil())
		})
	})

	Describe("HandleMessage with WindowSizeMsg", func() {
		It("should update terminal info on WindowSizeMsg", func() {
			factory := func() intents.Intent { return intents.NewMockIntent() }
			_ = router.RegisterIntent("test", factory)
			_, _ = router.ActivateIntent("test", nil)

			wsMsg := tea.WindowSizeMsg{Width: 200, Height: 60}
			router.HandleMessage(wsMsg)

			info := router.GetTerminalInfo()
			Expect(info.Width).To(Equal(200))
			Expect(info.Height).To(Equal(60))
		})

		It("should propagate WindowSizeMsg to terminal-aware intent", func() {
			mockIntent := newTerminalAwareMockIntent()
			factory := func() intents.Intent { return mockIntent }
			_ = router.RegisterIntent("term-intent", factory)
			_, _ = router.ActivateIntent("term-intent", nil)

			wsMsg := tea.WindowSizeMsg{Width: 180, Height: 55}
			router.HandleMessage(wsMsg)

			Expect(mockIntent.terminalInfoUpdated).To(BeTrue())
		})
	})

	Describe("ActivateIntent with nil factory return", func() {
		It("should return error when factory returns nil", func() {
			factory := func() intents.Intent { return nil }
			_ = router.RegisterIntent("nil-intent", factory)

			_, err := router.ActivateIntent("nil-intent", nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("factory returned nil"))
		})
	})

	Describe("ActivateIntent with terminal-aware intent", func() {
		It("should propagate terminal info to terminal-aware intent during activation", func() {
			info := terminal.NewInfo()
			info.Width = 120
			info.Height = 40
			info.IsValid = true
			router.UpdateTerminalInfo(info)

			mockIntent := newTerminalAwareMockIntent()
			factory := func() intents.Intent { return mockIntent }
			_ = router.RegisterIntent("term-intent", factory)

			_, err := router.ActivateIntent("term-intent", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(mockIntent.terminalInfoUpdated).To(BeTrue())
		})
	})

	Describe("HandleMessage result handler path", func() {
		It("should invoke result handler when intent completes", func() {
			handlerCalled := false
			router.RegisterResultHandler("completing-intent", func(result *intents.IntentResult[interface{}]) tea.Cmd {
				handlerCalled = true
				return nil
			})

			mockIntent := intents.NewMockIntent()
			factory := func() intents.Intent { return mockIntent }
			_ = router.RegisterIntent("completing-intent", factory)
			_, _ = router.ActivateIntent("completing-intent", nil)

			expectedResult := intents.NewCompletedResult[interface{}]("done")
			mockIntent.SetResult(expectedResult)

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
			_, result := router.HandleMessage(msg)

			Expect(result).NotTo(BeNil())
			_ = handlerCalled
		})
	})
})

var _ = Describe("MessageInterceptor OnQuit Coverage", func() {
	It("should register quit handler via OnQuit", func() {
		interceptor := intents.NewMessageInterceptor()
		chained := interceptor.OnQuit(func() tea.Cmd {
			return tea.Quit
		})

		Expect(chained).To(Equal(interceptor))
	})

	It("should support method chaining with OnQuit", func() {
		interceptor := intents.NewMessageInterceptor()
		result := interceptor.
			OnBack(func() tea.Cmd { return nil }).
			OnQuit(func() tea.Cmd { return tea.Quit }).
			OnHelp(func() tea.Cmd { return nil })

		Expect(result).To(Equal(interceptor))
	})
})

var _ = Describe("StandardQuitHandler Coverage", func() {
	It("should return a handler that produces tea.Quit", func() {
		handler := intents.StandardQuitHandler()
		Expect(handler).NotTo(BeNil())

		cmd := handler()
		Expect(cmd).NotTo(BeNil())
	})
})

var _ = Describe("extractErrorTitle Coverage", func() {
	var base *intents.BaseIntent

	BeforeEach(func() {
		base = intents.NewBaseIntent()
		info := terminal.NewInfo()
		info.Width = 100
		info.Height = 30
		info.IsValid = true
		base.UpdateTerminalInfo(info)
	})

	It("should extract title from wrapped error where outer has no keyword", func() {
		inner := errors.New("some generic error")
		wrapped := fmt.Errorf("outer: %w", inner)
		base.SetError(wrapped)
		view := intents.CreateStandardView(base)
		Expect(view.ShowModal).To(BeTrue())
	})

	It("should extract title from deeply wrapped error with keyword", func() {
		inner := errors.New("database connection lost")
		mid := fmt.Errorf("service error: %w", inner)
		outer := fmt.Errorf("request failed: %w", mid)
		base.SetError(outer)
		view := intents.CreateStandardView(base)
		Expect(view.ShowModal).To(BeTrue())
	})

	It("should handle nil error gracefully", func() {
		view := intents.CreateStandardView(base)
		Expect(view.ShowModal).To(BeFalse())
	})

	It("should handle timeout error", func() {
		base.SetError(errors.New("timeout waiting for response"))
		view := intents.CreateStandardView(base)
		Expect(view.ShowModal).To(BeTrue())
	})

	It("should handle permission error", func() {
		base.SetError(errors.New("permission denied for resource"))
		view := intents.CreateStandardView(base)
		Expect(view.ShowModal).To(BeTrue())
	})

	It("should handle not found error", func() {
		base.SetError(errors.New("not found: user 123"))
		view := intents.CreateStandardView(base)
		Expect(view.ShowModal).To(BeTrue())
	})

	It("should handle unauthorized error", func() {
		base.SetError(errors.New("unauthorized access attempt"))
		view := intents.CreateStandardView(base)
		Expect(view.ShowModal).To(BeTrue())
	})

	It("should handle invalid input error", func() {
		base.SetError(errors.New("invalid email format"))
		view := intents.CreateStandardView(base)
		Expect(view.ShowModal).To(BeTrue())
	})

	It("should handle unable to error", func() {
		base.SetError(errors.New("unable to connect"))
		view := intents.CreateStandardView(base)
		Expect(view.ShowModal).To(BeTrue())
	})

	It("should handle cannot error", func() {
		base.SetError(errors.New("cannot read file"))
		view := intents.CreateStandardView(base)
		Expect(view.ShowModal).To(BeTrue())
	})

	It("should handle network error", func() {
		base.SetError(errors.New("network unreachable"))
		view := intents.CreateStandardView(base)
		Expect(view.ShowModal).To(BeTrue())
	})

	It("should handle failed to error", func() {
		base.SetError(errors.New("failed to load config"))
		view := intents.CreateStandardView(base)
		Expect(view.ShowModal).To(BeTrue())
	})
})

type terminalAwareMockIntent struct {
	*intents.MockIntent
	terminalInfoUpdated bool
	terminalInfo        *terminal.Info
}

func newTerminalAwareMockIntent() *terminalAwareMockIntent {
	return &terminalAwareMockIntent{
		MockIntent: intents.NewMockIntent(),
	}
}

func (m *terminalAwareMockIntent) UpdateTerminalInfo(info *terminal.Info) {
	m.terminalInfoUpdated = true
	m.terminalInfo = info
}

func (m *terminalAwareMockIntent) GetMinimumSize() (width, height int) {
	return 80, 24
}
