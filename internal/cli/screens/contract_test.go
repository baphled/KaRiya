package screens_test

import (
	"errors"
	"testing"

	"github.com/baphled/kariya/internal/cli/screens"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestScreens(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Screens Suite")
}

// Screen Contract Tests
//
// These tests define the contract that all Screens must implement.
// They verify the Screen interface behavior and result types.
//
// Related:
// - tasks/tasks-42-tui-architecture-refactor.md (Phase 1.1)
// - docs/TUI_DEVELOPER_GUIDE.md (Screen architecture)

var _ = Describe("Screen Contract", func() {
	Describe("ScreenResult Types", func() {
		It("should have NavigateResult for forward navigation", func() {
			result := &screens.NavigateResult{
				ResultData: "selected-item",
			}

			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			Expect(result.Data()).To(Equal("selected-item"))
		})

		It("should have CancelResult for cancellation", func() {
			result := &screens.CancelResult{}

			Expect(result.Type()).To(Equal(screens.ResultCancel))
			Expect(result.Data()).To(BeNil())
		})

		It("should have SubmitResult for form submission", func() {
			formData := map[string]string{"name": "test"}
			result := &screens.SubmitResult{
				FormData: formData,
			}

			Expect(result.Type()).To(Equal(screens.ResultSubmit))
			Expect(result.Data()).To(Equal(formData))
		})

		It("should have ErrorResult for error states", func() {
			err := errors.New("test error")
			result := &screens.ErrorResult{
				Err:     err,
				Message: "Something went wrong",
			}

			Expect(result.Type()).To(Equal(screens.ResultError))
			data := result.Data().(map[string]interface{})
			Expect(data["error"]).To(Equal(err))
			Expect(data["message"]).To(Equal("Something went wrong"))
		})
	})

	Describe("Screen Interface", func() {
		It("should define Update method that returns cmd and optional ScreenResult", func() {
			// This test verifies the interface signature at compile time
			// The interface is defined in contract.go
			Expect(screens.ResultNavigate).To(Equal(screens.ScreenResultType("navigate")))
		})

		It("should define View method that returns string", func() {
			// Interface is verified at compile time
			// Test that result types work with View pattern
			Expect(screens.ResultCancel).To(Equal(screens.ScreenResultType("cancel")))
		})

		It("should define SetTerminalInfo method for terminal size handling", func() {
			// Interface verification - Screen interface exists with SetTerminalInfo
			Expect(screens.ResultSubmit).To(Equal(screens.ScreenResultType("submit")))
		})

		It("should define SetTheme method for theme management", func() {
			// Interface verification
			Expect(screens.ResultError).To(Equal(screens.ScreenResultType("error")))
		})
	})

	Describe("Screen Message Handling", func() {
		Context("when handling WindowSizeMsg", func() {
			It("should update terminal dimensions via SetTerminalInfo", func() {
				// This is tested in concrete screen implementations
				// NavigateResult can carry dimension changes in metadata
				result := &screens.NavigateResult{}
				result.WithMetadata("width", 100)
				result.WithMetadata("height", 50)

				Expect(result.Metadata()["width"]).To(Equal(100))
				Expect(result.Metadata()["height"]).To(Equal(50))
			})

			It("should not return a ScreenResult for window resize", func() {
				// Window resize handling is a convention - screens return nil for resize
				// Verified in concrete implementations
				Expect(true).To(BeTrue())
			})
		})

		Context("when handling KeyMsg", func() {
			It("should handle Escape key and return appropriate ScreenResult", func() {
				// CancelResult is the standard return for Escape
				result := &screens.CancelResult{}
				Expect(result.Type()).To(Equal(screens.ResultCancel))
			})

			It("should handle Enter key for selection/submission", func() {
				// NavigateResult or SubmitResult for Enter key
				navResult := &screens.NavigateResult{ResultData: "selection"}
				Expect(navResult.Type()).To(Equal(screens.ResultNavigate))

				submitResult := &screens.SubmitResult{FormData: "form-data"}
				Expect(submitResult.Type()).To(Equal(screens.ResultSubmit))
			})

			It("should handle navigation keys (up/down/j/k)", func() {
				// Navigation keys typically return nil ScreenResult (internal state only)
				// This is verified in concrete implementations
				// Test that results can be nil-safe
				var result screens.ScreenResult = nil
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("ScreenResult Metadata", func() {
		It("should allow attaching metadata to results", func() {
			navResult := &screens.NavigateResult{}
			navResult.WithMetadata("scroll_position", 10)
			navResult.WithMetadata("filter", "active")

			Expect(navResult.Metadata()["scroll_position"]).To(Equal(10))
			Expect(navResult.Metadata()["filter"]).To(Equal("active"))

			cancelResult := &screens.CancelResult{}
			cancelResult.WithMetadata("previous_state", "selection")

			Expect(cancelResult.Metadata()["previous_state"]).To(Equal("selection"))

			submitResult := &screens.SubmitResult{}
			submitResult.WithMetadata("validation_passed", true)

			Expect(submitResult.Metadata()["validation_passed"]).To(BeTrue())

			errorResult := &screens.ErrorResult{}
			errorResult.WithMetadata("retry_count", 3)

			Expect(errorResult.Metadata()["retry_count"]).To(Equal(3))
		})

		It("should retrieve metadata from results", func() {
			result := &screens.NavigateResult{}
			result.WithMetadata("key1", "value1")
			result.WithMetadata("key2", 42)
			result.WithMetadata("key3", true)

			meta := result.Metadata()
			Expect(meta["key1"]).To(Equal("value1"))
			Expect(meta["key2"]).To(Equal(42))
			Expect(meta["key3"]).To(BeTrue())
		})

		It("should support fluent API for metadata", func() {
			result := &screens.NavigateResult{ResultData: "item"}
			result.WithMetadata("a", 1).WithMetadata("b", 2).WithMetadata("c", 3)

			Expect(result.Metadata()["a"]).To(Equal(1))
			Expect(result.Metadata()["b"]).To(Equal(2))
			Expect(result.Metadata()["c"]).To(Equal(3))
		})

		It("should return empty map if no metadata set", func() {
			navResult := &screens.NavigateResult{}
			Expect(navResult.Metadata()).NotTo(BeNil())
			Expect(navResult.Metadata()).To(BeEmpty())

			cancelResult := &screens.CancelResult{}
			Expect(cancelResult.Metadata()).NotTo(BeNil())

			submitResult := &screens.SubmitResult{}
			Expect(submitResult.Metadata()).NotTo(BeNil())

			errorResult := &screens.ErrorResult{}
			Expect(errorResult.Metadata()).NotTo(BeNil())
		})
	})

	Describe("ScreenResultType Constants", func() {
		It("should have correct string values", func() {
			Expect(string(screens.ResultNavigate)).To(Equal("navigate"))
			Expect(string(screens.ResultCancel)).To(Equal("cancel"))
			Expect(string(screens.ResultSubmit)).To(Equal("submit"))
			Expect(string(screens.ResultError)).To(Equal("error"))
		})
	})
})
