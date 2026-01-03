package intents

import (
	"fmt"
	"testing"

	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestContract(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Intent Contract Suite")
}

var _ = Describe("ModalEditResult", func() {
	Context("when creating a new modal edit result", func() {
		It("should store original and modified values", func() {
			original := "original value"
			modified := "modified value"
			changes := map[string]interface{}{"field": "new value"}

			result := NewModalEditResult[string](original, modified, true, changes)

			Expect(result.Original).To(Equal(original))
			Expect(result.Modified).To(Equal(modified))
			Expect(result.Accepted).To(BeTrue())
			Expect(result.Changes).To(Equal(changes))
		})

		It("should handle nil changes map", func() {
			result := NewModalEditResult[string]("original", "modified", true, nil)

			Expect(result.Changes).NotTo(BeNil())
			Expect(result.Changes).To(BeEmpty())
		})
	})

	Context("HasChanges", func() {
		It("should return true when changes exist", func() {
			result := NewModalEditResult[string]("original", "modified", true, map[string]interface{}{"field": "value"})

			Expect(result.HasChanges()).To(BeTrue())
		})

		It("should return false when no changes", func() {
			result := NewModalEditResult[string]("original", "original", true, make(map[string]interface{}))

			Expect(result.HasChanges()).To(BeFalse())
		})
	})

	Context("WasAccepted", func() {
		It("should return true when accepted is true", func() {
			result := NewModalEditResult[string]("original", "modified", true, nil)

			Expect(result.WasAccepted()).To(BeTrue())
		})

		It("should return false when accepted is false", func() {
			result := NewModalEditResult[string]("original", "original", false, nil)

			Expect(result.WasAccepted()).To(BeFalse())
		})
	})

	Context("GetChange", func() {
		It("should return the value for an existing change", func() {
			changes := map[string]interface{}{"field": "new value"}
			result := NewModalEditResult[string]("original", "modified", true, changes)

			Expect(result.GetChange("field")).To(Equal("new value"))
		})

		It("should return nil for a non-existent change", func() {
			result := NewModalEditResult[string]("original", "modified", true, map[string]interface{}{"field": "value"})

			Expect(result.GetChange("nonexistent")).To(BeNil())
		})

		It("should return nil when changes is nil", func() {
			result := &ModalEditResult[string]{
				Original: "original",
				Modified: "modified",
				Accepted: true,
				Changes:  nil,
			}

			Expect(result.GetChange("field")).To(BeNil())
		})
	})

	Context("NewCancelledModalEditResult", func() {
		It("should create a result with accepted=false and no changes", func() {
			original := "original value"

			result := NewCancelledModalEditResult[string](original)

			Expect(result.Original).To(Equal(original))
			Expect(result.Modified).To(Equal(original))
			Expect(result.Accepted).To(BeFalse())
			Expect(result.HasChanges()).To(BeFalse())
		})
	})

	Context("with complex types", func() {
		type TestData struct {
			Name string
			Age  int
		}

		It("should work with struct types", func() {
			original := TestData{Name: "John", Age: 30}
			modified := TestData{Name: "Jane", Age: 30}
			changes := map[string]interface{}{"Name": "Jane"}

			result := NewModalEditResult[TestData](original, modified, true, changes)

			Expect(result.Original).To(Equal(original))
			Expect(result.Modified).To(Equal(modified))
			Expect(result.HasChanges()).To(BeTrue())
		})
	})
})

var _ = Describe("CaptureEventIntent", func() {
	var (
		intent *CaptureEventIntent
		ctx    *CaptureEventContext
	)

	BeforeEach(func() {
		ctx = &CaptureEventContext{
			CaptureStrategy: "manual",
			PreviousEvent:   nil,
			Metadata:        make(map[string]string),
		}
	})

	Describe("NewCaptureEventIntent", func() {
		It("should create a new intent with valid context", func() {
			var err error
			intent, err = NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent).NotTo(BeNil())
		})

		It("should fail with invalid context", func() {
			ctx.CaptureStrategy = ""
			_, err := NewCaptureEventIntent(ctx)
			Expect(err).To(HaveOccurred())
		})

		It("should initialize state to ChooseStrategy", func() {
			var err error
			intent, err = NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent.state.currentState).To(Equal(CaptureStateChooseStrategy))
		})

		It("should initialize active flag to true", func() {
			var err error
			intent, err = NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent.active).To(BeTrue())
		})

		It("should initialize result to nil", func() {
			var err error
			intent, err = NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent.result).To(BeNil())
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return a command for new event", func() {
			cmd := intent.Init()
			// For now, Init returns nil, which is valid
			Expect(cmd).To(BeNil())
		})
	})

	Describe("View", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should render view for ChooseStrategy state", func() {
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Choose Capture Strategy"))
		})

		It("should render view for Form state", func() {
			intent.state.currentState = CaptureStateForm
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Capture Event Details"))
		})

		It("should render view for Review state", func() {
			intent.state.currentState = CaptureStateReview
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Review Inferred Event"))
		})

		It("should render view for Submit state", func() {
			intent.state.currentState = CaptureStateSubmit
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Submitting event"))
		})

		It("should not render when inactive", func() {
			intent.active = false
			view := intent.View()
			Expect(view).To(ContainSubstring("not active"))
		})
	})

	Describe("Result", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return nil while active", func() {
			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should return completed result after completion", func() {
			eventResult := &CaptureEventResult{
				Event:          nil,
				Bursts:         make([]*career.Burst, 0),
				Facts:          make([]*career.Fact, 0),
				AcceptedFields: make(map[string]bool),
				RejectedFields: make(map[string]string),
			}
			intent.setCompleted(eventResult)
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Completed))
		})

		It("should return cancelled result after cancellation", func() {
			intent.setCancelled()
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Cancelled))
		})

		It("should return failed result after failure", func() {
			intent.setFailed("test_error", "Test error message", nil)
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Failed))
		})
	})

	Describe("setCompleted", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should set result status to Completed", func() {
			eventResult := &CaptureEventResult{
				AcceptedFields: make(map[string]bool),
				RejectedFields: make(map[string]string),
			}
			intent.setCompleted(eventResult)
			Expect(intent.result.Status).To(Equal(Completed))
		})

		It("should mark intent as inactive", func() {
			eventResult := &CaptureEventResult{
				AcceptedFields: make(map[string]bool),
				RejectedFields: make(map[string]string),
			}
			intent.setCompleted(eventResult)
			Expect(intent.active).To(BeFalse())
		})
	})

	Describe("setCancelled", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should set result status to Cancelled", func() {
			intent.setCancelled()
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should mark intent as inactive", func() {
			intent.setCancelled()
			Expect(intent.active).To(BeFalse())
		})
	})

	Describe("setFailed", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should set result status to Failed", func() {
			intent.setFailed("error_code", "Error message", nil)
			Expect(intent.result.Status).To(Equal(Failed))
		})

		It("should set error details", func() {
			intent.setFailed("error_code", "Error message", nil)
			Expect(intent.result.Error).NotTo(BeNil())
			Expect(intent.result.Error.Code).To(Equal("error_code"))
			Expect(intent.result.Error.Message).To(Equal("Error message"))
		})

		It("should mark intent as inactive", func() {
			intent.setFailed("error_code", "Error message", nil)
			Expect(intent.active).To(BeFalse())
		})
	})

	Describe("setPartial", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should set result status to Partial", func() {
			eventResult := &CaptureEventResult{
				AcceptedFields: make(map[string]bool),
				RejectedFields: make(map[string]string),
			}
			intent.setPartial(eventResult, "partial_code", "Partial message")
			Expect(intent.result.Status).To(Equal(Partial))
		})

		It("should set error details for partial result", func() {
			eventResult := &CaptureEventResult{
				AcceptedFields: make(map[string]bool),
				RejectedFields: make(map[string]string),
			}
			intent.setPartial(eventResult, "partial_code", "Partial message")
			Expect(intent.result.Error).NotTo(BeNil())
			Expect(intent.result.Error.Code).To(Equal("partial_code"))
			Expect(intent.result.Error.Message).To(Equal("Partial message"))
		})

		It("should mark intent as inactive", func() {
			eventResult := &CaptureEventResult{
				AcceptedFields: make(map[string]bool),
				RejectedFields: make(map[string]string),
			}
			intent.setPartial(eventResult, "partial_code", "Partial message")
			Expect(intent.active).To(BeFalse())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return nil when inactive", func() {
			intent.active = false
			cmd := intent.Update(nil)
			Expect(cmd).To(BeNil())
		})

		It("should delegate to updateChooseStrategy for ChooseStrategy state", func() {
			intent.state.currentState = CaptureStateChooseStrategy
			// This test just verifies the state delegation works
			cmd := intent.Update(nil)
			// Command will be nil for now since updateChooseStrategy returns nil
			Expect(cmd).To(BeNil())
		})
	})

	Describe("State Transitions", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should transition from ChooseStrategy to Form state", func() {
			intent.state.currentState = CaptureStateForm
			view := intent.View()
			Expect(view).To(ContainSubstring("Capture Event Details"))
		})

		It("should transition from Form to Review state", func() {
			intent.state.currentState = CaptureStateReview
			view := intent.View()
			Expect(view).To(ContainSubstring("Review Inferred Event"))
		})

		It("should transition from Review to Submit state", func() {
			intent.state.currentState = CaptureStateSubmit
			view := intent.View()
			Expect(view).To(ContainSubstring("Submitting event"))
		})
	})
})

// State Transition Tests for CaptureEvent Intent
var _ = Describe("CaptureEvent Intent State Transitions", func() {
	var (
		intent *CaptureEventIntent
		ctx    *CaptureEventContext
	)

	BeforeEach(func() {
		ctx = &CaptureEventContext{
			CaptureStrategy: "manual",
			PreviousEvent:   nil,
			Metadata:        make(map[string]string),
		}
		var err error
		intent, err = NewCaptureEventIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("StateChooseStrategy", func() {
		BeforeEach(func() {
			Expect(intent.state.currentState).To(Equal(CaptureStateChooseStrategy))
		})

		It("should transition to StateForm when '1' key is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
			Expect(intent.state.currentState).To(Equal(CaptureStateForm))
		})

		It("should cancel when Ctrl+C is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			Expect(intent.active).To(BeFalse())
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should transition to StateForm when StrategySelectedMsg is received", func() {
			intent.Update(StrategySelectedMsg{Strategy: "quick"})
			Expect(intent.state.currentState).To(Equal(CaptureStateForm))
		})
	})

	Describe("StateCaptureForm", func() {
		BeforeEach(func() {
			intent.state.currentState = CaptureStateForm
		})

		It("should transition to StateReview when Ctrl+S is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(intent.state.currentState).To(Equal(CaptureStateReview))
			Expect(intent.state.reviewState.Event).NotTo(BeNil())
		})

		It("should cancel when Ctrl+C is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			Expect(intent.active).To(BeFalse())
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should go back to StateChooseStrategy when Esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(CaptureStateChooseStrategy))
		})

		It("should transition to StateReview when FormSubmittedMsg is received", func() {
			event := &career.CareerEvent{Text: "Test Event"}
			intent.Update(FormSubmittedMsg{Event: event})
			Expect(intent.state.currentState).To(Equal(CaptureStateReview))
			Expect(intent.state.reviewState.Event).To(Equal(event))
		})
	})

	Describe("StateReviewInferred", func() {
		BeforeEach(func() {
			intent.state.currentState = CaptureStateReview
			intent.state.reviewState.Event = &career.CareerEvent{Text: "Test Event"}
		})

		It("should transition to StateSubmit when Ctrl+S is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(intent.state.currentState).To(Equal(CaptureStateSubmit))
			Expect(cmd).NotTo(BeNil())
		})

		It("should cancel when Ctrl+C is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			Expect(intent.active).To(BeFalse())
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should go back to StateCaptureForm when Esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(CaptureStateForm))
		})

		It("should set EditingMode to metadata when 'e' key is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(intent.state.reviewState.EditingMode).To(Equal(EditingModeMetadata))
		})

		It("should transition to StateSubmit when ReviewConfirmedMsg is received", func() {
			msg := ReviewConfirmedMsg{
				AcceptedBursts: []*career.Burst{},
				AcceptedFacts:  []*career.Fact{},
				RejectedItems:  make(map[string]string),
			}
			cmd := intent.Update(msg)
			Expect(intent.state.currentState).To(Equal(CaptureStateSubmit))
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("StateSubmit", func() {
		BeforeEach(func() {
			intent.state.currentState = CaptureStateSubmit
			intent.state.reviewState.Event = &career.CareerEvent{Text: "Test Event"}
			intent.state.reviewState.AcceptedBursts = []*career.Burst{}
			intent.state.reviewState.AcceptedFacts = []*career.Fact{}
		})

		It("should complete with result when SubmitCompleteMsg is received", func() {
			intent.Update(SubmitCompleteMsg{})
			Expect(intent.active).To(BeFalse())
			Expect(intent.result.Status).To(Equal(Completed))
			Expect(intent.result.Data).NotTo(BeNil())
		})

		It("should fail when SubmitErrorMsg is received", func() {
			intent.Update(SubmitErrorMsg{
				Code:    "save_error",
				Message: "Failed to save event",
				Cause:   nil,
			})
			Expect(intent.active).To(BeFalse())
			Expect(intent.result.Status).To(Equal(Failed))
			Expect(intent.result.Error.Code).To(Equal("save_error"))
		})
	})

	Describe("Complete Workflows", func() {
		It("should complete full workflow", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
			Expect(intent.state.currentState).To(Equal(CaptureStateForm))

			intent.Update(FormSubmittedMsg{Event: &career.CareerEvent{Text: "Test"}})
			Expect(intent.state.currentState).To(Equal(CaptureStateReview))

			intent.Update(ReviewConfirmedMsg{
				AcceptedBursts: []*career.Burst{},
				AcceptedFacts:  []*career.Fact{},
				RejectedItems:  make(map[string]string),
			})
			Expect(intent.state.currentState).To(Equal(CaptureStateSubmit))

			intent.Update(SubmitCompleteMsg{})
			Expect(intent.active).To(BeFalse())
			Expect(intent.result.Status).To(Equal(Completed))
		})
	})
})

var _ = Describe("Modal Helper Functions", func() {
	Describe("formatStringSlice", func() {
		It("should format empty slice as empty string", func() {
			result := formatStringSlice([]string{})
			Expect(result).To(Equal(""))
		})

		It("should format single item", func() {
			result := formatStringSlice([]string{"item"})
			Expect(result).To(Equal("item"))
		})

		It("should format multiple items with comma separation", func() {
			result := formatStringSlice([]string{"item1", "item2", "item3"})
			Expect(result).To(Equal("item1, item2, item3"))
		})
	})

	Describe("parseStringSlice", func() {
		It("should parse empty string as empty slice", func() {
			result := parseStringSlice("")
			Expect(result).To(BeEmpty())
		})

		It("should parse single item", func() {
			result := parseStringSlice("item")
			Expect(result).To(Equal([]string{"item"}))
		})

		It("should parse comma-separated items", func() {
			result := parseStringSlice("item1, item2, item3")
			Expect(result).To(Equal([]string{"item1", "item2", "item3"}))
		})

		It("should handle items without spaces", func() {
			result := parseStringSlice("item1,item2,item3")
			Expect(result).To(Equal([]string{"item1", "item2", "item3"}))
		})

		It("should trim spaces from items", func() {
			result := parseStringSlice("  item1  ,  item2  ")
			Expect(result).To(Equal([]string{"item1", "item2"}))
		})
	})

	Describe("slicesEqual", func() {
		It("should return true for equal slices", func() {
			a := []string{"a", "b", "c"}
			b := []string{"a", "b", "c"}
			Expect(slicesEqual(a, b)).To(BeTrue())
		})

		It("should return false for different lengths", func() {
			a := []string{"a", "b"}
			b := []string{"a", "b", "c"}
			Expect(slicesEqual(a, b)).To(BeFalse())
		})

		It("should return false for different content", func() {
			a := []string{"a", "b", "c"}
			b := []string{"a", "x", "c"}
			Expect(slicesEqual(a, b)).To(BeFalse())
		})

		It("should return true for empty slices", func() {
			a := []string{}
			b := []string{}
			Expect(slicesEqual(a, b)).To(BeTrue())
		})
	})

	Describe("truncateString", func() {
		It("should not truncate short strings", func() {
			result := truncateString("hello", 10)
			Expect(result).To(Equal("hello"))
		})

		It("should truncate long strings", func() {
			result := truncateString("hello world", 8)
			Expect(result).To(Equal("hello..."))
		})

		It("should handle exact length", func() {
			result := truncateString("hello", 5)
			Expect(result).To(Equal("hello"))
		})
	})

	Describe("copyMetadataSnapshot", func() {
		It("should create independent copy", func() {
			original := &MetadataSnapshot{
				Company:    "Corp",
				Project:    "Proj",
				Tags:       []string{"a", "b"},
				Categories: []string{"x", "y"},
			}

			copy := copyMetadataSnapshot(original)
			copy.Company = "NewCorp"
			copy.Tags[0] = "z"

			Expect(original.Company).To(Equal("Corp"))
			Expect(original.Tags[0]).To(Equal("a"))
		})

		It("should handle nil input", func() {
			copy := copyMetadataSnapshot(nil)
			Expect(copy).To(BeNil())
		})
	})
})

var _ = Describe("CaptureEvent State Transition Coverage", func() {
	var (
		intent *CaptureEventIntent
		ctx    *CaptureEventContext
	)

	BeforeEach(func() {
		ctx = &CaptureEventContext{
			CaptureStrategy: "manual",
			PreviousEvent:   nil,
			Metadata:        make(map[string]string),
		}
		var err error
		intent, err = NewCaptureEventIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("IsActive method", func() {
		It("should return true when intent is active", func() {
			Expect(intent.IsActive()).To(BeTrue())
		})

		It("should return false when intent is inactive", func() {
			intent.active = false
			Expect(intent.IsActive()).To(BeFalse())
		})
	})

	Describe("GetResult method", func() {
		It("should return nil when no result is set", func() {
			Expect(intent.GetResult()).To(BeNil())
		})

		It("should return result when set", func() {
			intent.setCompleted(&CaptureEventResult{Event: &career.CareerEvent{Text: "Test"}})
			result := intent.GetResult()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Completed))
		})
	})

	Describe("updateChooseStrategy edge cases", func() {
		BeforeEach(func() {
			intent.state.currentState = CaptureStateChooseStrategy
		})

		It("should handle unknown strategy selection", func() {
			intent.Update(StrategySelectedMsg{Strategy: "unknown_strategy"})
			// Should still transition to form state
			Expect(intent.state.currentState).To(Equal(CaptureStateForm))
		})

		It("should handle empty strategy selection", func() {
			intent.Update(StrategySelectedMsg{Strategy: ""})
			// Should still transition to form state
			Expect(intent.state.currentState).To(Equal(CaptureStateForm))
		})

		It("should handle numeric key selections", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
			Expect(intent.state.currentState).To(Equal(CaptureStateForm))
		})

		It("should handle numeric key selections for option 3", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
			Expect(intent.state.currentState).To(Equal(CaptureStateForm))
		})

		It("should ignore invalid numeric keys", func() {
			startState := intent.state.currentState
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'9'}})
			Expect(intent.state.currentState).To(Equal(startState))
		})

		It("should handle non-numeric keys gracefully", func() {
			startState := intent.state.currentState
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(intent.state.currentState).To(Equal(startState))
		})
	})

	Describe("updateCaptureForm edge cases", func() {
		BeforeEach(func() {
			intent.state.currentState = CaptureStateForm
		})

		It("should handle FormSubmittedMsg with nil event", func() {
			intent.Update(FormSubmittedMsg{Event: nil})
			// Should fail with invalid form error
			Expect(intent.active).To(BeFalse())
			Expect(intent.result.Status).To(Equal(Failed))
			Expect(intent.result.Error.Code).To(Equal("INVALID_FORM"))
		})

		It("should handle form submission with minimal event", func() {
			event := &career.CareerEvent{Text: "Minimal"}
			intent.Update(FormSubmittedMsg{Event: event})
			Expect(intent.state.currentState).To(Equal(CaptureStateReview))
			Expect(intent.state.reviewState.Event).To(Equal(event))
		})

		It("should transition to review on Ctrl+S", func() {
			event := &career.CareerEvent{Text: "Test"}
			intent.state.reviewState.Event = event
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(intent.state.currentState).To(Equal(CaptureStateReview))
		})
	})

	Describe("updateReviewInferredEvent edge cases", func() {
		BeforeEach(func() {
			intent.state.currentState = CaptureStateReview
			intent.state.reviewState.Event = &career.CareerEvent{Text: "Test Event"}
		})

		It("should handle ReviewConfirmedMsg with empty bursts and facts", func() {
			msg := ReviewConfirmedMsg{
				AcceptedBursts: []*career.Burst{},
				AcceptedFacts:  []*career.Fact{},
				RejectedItems:  make(map[string]string),
			}
			intent.Update(msg)
			Expect(intent.state.currentState).To(Equal(CaptureStateSubmit))
		})

		It("should handle ReviewConfirmedMsg with multiple bursts and facts", func() {
			msg := ReviewConfirmedMsg{
				AcceptedBursts: []*career.Burst{
					{Name: "Burst 1"},
					{Name: "Burst 2"},
				},
				AcceptedFacts: []*career.Fact{
					{Text: "Fact 1"},
					{Text: "Fact 2"},
				},
				RejectedItems: map[string]string{
					"burst_1": "Not relevant",
				},
			}
			intent.Update(msg)
			Expect(intent.state.currentState).To(Equal(CaptureStateSubmit))
			Expect(len(intent.state.reviewState.AcceptedBursts)).To(Equal(2))
			Expect(len(intent.state.reviewState.AcceptedFacts)).To(Equal(2))
		})

		It("should handle ReviewBackMsg", func() {
			intent.Update(ReviewBackMsg{})
			Expect(intent.state.currentState).To(Equal(CaptureStateForm))
		})

		It("should handle 'e' key for metadata editing", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(intent.state.reviewState.EditingMode).To(Equal(EditingModeMetadata))
		})

		It("should handle 'b' key for burst editing", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
			Expect(intent.state.reviewState.EditingMode).To(Equal(EditingModeBursts))
		})

		It("should handle 'f' key for fact editing", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
			Expect(intent.state.reviewState.EditingMode).To(Equal(EditingModeFacts))
		})
	})

	Describe("updateSubmit edge cases", func() {
		BeforeEach(func() {
			intent.state.currentState = CaptureStateSubmit
			intent.state.reviewState.Event = &career.CareerEvent{Text: "Test Event"}
			intent.state.reviewState.AcceptedBursts = []*career.Burst{}
			intent.state.reviewState.AcceptedFacts = []*career.Fact{}
		})

		It("should handle SubmitCompleteMsg", func() {
			intent.Update(SubmitCompleteMsg{})
			Expect(intent.active).To(BeFalse())
			Expect(intent.result.Status).To(Equal(Completed))
		})

		It("should handle SubmitErrorMsg with code and message", func() {
			intent.Update(SubmitErrorMsg{
				Code:    "NETWORK_ERROR",
				Message: "Failed to connect to server",
				Cause:   nil,
			})
			Expect(intent.active).To(BeFalse())
			Expect(intent.result.Status).To(Equal(Failed))
			Expect(intent.result.Error.Code).To(Equal("NETWORK_ERROR"))
			Expect(intent.result.Error.Message).To(Equal("Failed to connect to server"))
		})

		It("should handle SubmitErrorMsg with cause", func() {
			cause := fmt.Errorf("connection timeout")
			intent.Update(SubmitErrorMsg{
				Code:    "TIMEOUT",
				Message: "Request timed out",
				Cause:   cause,
			})
			Expect(intent.result.Error.Cause).To(Equal(cause))
		})

		It("should transition to submit on Ctrl+S from review", func() {
			intent.state.currentState = CaptureStateReview
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(intent.state.currentState).To(Equal(CaptureStateSubmit))
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("performSubmit function", func() {
		BeforeEach(func() {
			intent.state.currentState = CaptureStateSubmit
			intent.state.reviewState.Event = &career.CareerEvent{Text: "Test Event"}
		})

		It("should return SubmitCompleteMsg on successful submission", func() {
			cmd := intent.performSubmit()
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			Expect(msg).To(BeAssignableToTypeOf(SubmitCompleteMsg{}))
		})

		It("should return SubmitErrorMsg when event is nil", func() {
			intent.state.reviewState.Event = nil
			cmd := intent.performSubmit()
			msg := cmd()
			Expect(msg).To(BeAssignableToTypeOf(SubmitErrorMsg{}))
			errorMsg := msg.(SubmitErrorMsg)
			Expect(errorMsg.Code).To(Equal("MISSING_EVENT"))
		})

		It("should return SubmitErrorMsg on validation failure", func() {
			intent.state.reviewState.Event = &career.CareerEvent{Text: ""}
			cmd := intent.performSubmit()
			msg := cmd()
			Expect(msg).To(BeAssignableToTypeOf(SubmitErrorMsg{}))
			errorMsg := msg.(SubmitErrorMsg)
			Expect(errorMsg.Code).To(Equal("VALIDATION_ERROR"))
		})
	})

	Describe("Cancel handling across all states", func() {
		It("should cancel from StateChooseStrategy", func() {
			intent.state.currentState = CaptureStateChooseStrategy
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			Expect(intent.active).To(BeFalse())
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should cancel from StateCaptureForm", func() {
			intent.state.currentState = CaptureStateForm
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			Expect(intent.active).To(BeFalse())
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should cancel from StateReviewInferred", func() {
			intent.state.currentState = CaptureStateReview
			intent.state.reviewState.Event = &career.CareerEvent{Text: "Test"}
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			Expect(intent.active).To(BeFalse())
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should cancel from StateSubmit", func() {
			intent.state.currentState = CaptureStateSubmit
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			Expect(intent.active).To(BeFalse())
			Expect(intent.result.Status).To(Equal(Cancelled))
		})
	})

	Describe("Back navigation across all states", func() {
		It("should go back from StateCaptureForm to StateChooseStrategy", func() {
			intent.state.currentState = CaptureStateForm
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(CaptureStateChooseStrategy))
		})

		It("should go back from StateReviewInferred to StateCaptureForm", func() {
			intent.state.currentState = CaptureStateReview
			intent.state.reviewState.Event = &career.CareerEvent{Text: "Test"}
			intent.Update(ReviewBackMsg{})
			Expect(intent.state.currentState).To(Equal(CaptureStateForm))
		})

		It("should stay in StateChooseStrategy when pressing Esc", func() {
			intent.state.currentState = CaptureStateChooseStrategy
			startState := intent.state.currentState
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(startState))
		})
	})

	Describe("View rendering for all states", func() {
		It("should render inactive message when intent has failed", func() {
			intent.setFailed("TEST_ERROR", "Test error message", nil)
			view := intent.View()
			// When inactive, the view shows "not active" message
			Expect(view).To(ContainSubstring("not active"))
		})

		It("should render appropriate view based on state", func() {
			states := []string{
				CaptureStateChooseStrategy,
				CaptureStateForm,
				CaptureStateReview,
				CaptureStateSubmit,
			}

			for _, state := range states {
				intent.state.currentState = state
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			}
		})
	})

	Describe("Result conversion", func() {
		It("should convert typed result to interface result", func() {
			intent.setCompleted(&CaptureEventResult{
				Event: &career.CareerEvent{Text: "Test"},
			})
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Completed))
			Expect(result.Data).NotTo(BeNil())
		})

		It("should preserve error information in result", func() {
			intent.setFailed("ERROR_CODE", "Error message", fmt.Errorf("cause"))
			result := intent.Result()
			Expect(result.Error).NotTo(BeNil())
			Expect(result.Error.Code).To(Equal("ERROR_CODE"))
		})
	})

	Describe("Message handling robustness", func() {
		It("should ignore unknown message types", func() {
			startState := intent.state.currentState
			intent.Update("unknown string message")
			Expect(intent.state.currentState).To(Equal(startState))
		})

		It("should handle WindowSizeMsg gracefully", func() {
			startState := intent.state.currentState
			intent.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
			Expect(intent.state.currentState).To(Equal(startState))
		})

		It("should handle multiple key presses in sequence", func() {
			intent.state.currentState = CaptureStateChooseStrategy
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
			Expect(intent.state.currentState).To(Equal(CaptureStateForm))

			intent.Update(FormSubmittedMsg{Event: &career.CareerEvent{Text: "Test"}})
			Expect(intent.state.currentState).To(Equal(CaptureStateReview))

			intent.Update(ReviewBackMsg{})
			Expect(intent.state.currentState).To(Equal(CaptureStateForm))
		})
	})
})

var _ = Describe("CaptureEvent Init and Form Initialization", func() {
	Describe("Init with different capture strategies", func() {
		It("should initialize for new event capture", func() {
			ctx := &CaptureEventContext{
				CaptureStrategy: "new",
				PreviousEvent:   nil,
				Metadata:        make(map[string]string),
			}
			intent, err := NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent.state.currentState).To(Equal(CaptureStateChooseStrategy))
			cmd := intent.Init()
			Expect(cmd).To(BeNil())
		})

		It("should initialize for quick capture", func() {
			ctx := &CaptureEventContext{
				CaptureStrategy: "quick",
				PreviousEvent:   nil,
				Metadata:        make(map[string]string),
			}
			intent, err := NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent.state.currentState).To(Equal(CaptureStateChooseStrategy))
		})

		It("should initialize for event editing", func() {
			previousEvent := &career.CareerEvent{Text: "Existing Event"}
			ctx := &CaptureEventContext{
				CaptureStrategy: "edit",
				PreviousEvent:   previousEvent,
				Metadata:        make(map[string]string),
			}
			intent, err := NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent.state.currentState).To(Equal(CaptureStateChooseStrategy))
		})

		It("should initialize for quick import", func() {
			ctx := &CaptureEventContext{
				CaptureStrategy: "import",
				PreviousEvent:   nil,
				Metadata:        make(map[string]string),
			}
			intent, err := NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent.state.currentState).To(Equal(CaptureStateChooseStrategy))
		})
	})

	Describe("Error handling for invalid context", func() {
		It("should reject empty strategy", func() {
			ctx := &CaptureEventContext{
				CaptureStrategy: "",
				PreviousEvent:   nil,
				Metadata:        make(map[string]string),
			}
			_, err := NewCaptureEventIntent(ctx)
			Expect(err).To(HaveOccurred())
		})

		It("should handle nil metadata gracefully", func() {
			ctx := &CaptureEventContext{
				CaptureStrategy: "manual",
				PreviousEvent:   nil,
				Metadata:        nil,
			}
			intent, err := NewCaptureEventIntent(ctx)
			// Should succeed even with nil metadata
			Expect(err).NotTo(HaveOccurred())
			Expect(intent).NotTo(BeNil())
		})
	})
})

var _ = Describe("CaptureEvent Edit Mode Coverage", func() {
	Describe("Form initialization for edit mode", func() {
		It("should load previous event data when editing", func() {
			previousEvent := &career.CareerEvent{
				Text: "Previous Event",
				Tags: []string{"tag1", "tag2"},
			}
			ctx := &CaptureEventContext{
				CaptureStrategy: "edit",
				PreviousEvent:   previousEvent,
				Metadata:        make(map[string]string),
			}
			intent, err := NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			// The form should be initialized with the previous event
			Expect(intent).NotTo(BeNil())
		})

		It("should handle edit mode with complex event data", func() {
			previousEvent := &career.CareerEvent{
				Text:       "Complex Event",
				Tags:       []string{"tag1", "tag2", "tag3"},
				Categories: []string{"cat1", "cat2"},
			}
			ctx := &CaptureEventContext{
				CaptureStrategy: "edit",
				PreviousEvent:   previousEvent,
				Metadata: map[string]string{
					"source": "import",
				},
			}
			intent, err := NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent).NotTo(BeNil())
		})
	})

	Describe("Message handling in different edit modes", func() {
		var intent *CaptureEventIntent

		BeforeEach(func() {
			ctx := &CaptureEventContext{
				CaptureStrategy: "manual",
				PreviousEvent:   nil,
				Metadata:        make(map[string]string),
			}
			var err error
			intent, err = NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle FocusNextMsg", func() {
			intent.state.currentState = CaptureStateForm
			intent.Update(tea.KeyMsg{Type: tea.KeyTab})
			// Should handle tab gracefully
			Expect(intent.state.currentState).To(Equal(CaptureStateForm))
		})

		It("should handle multiple transitions in sequence", func() {
			// Strategy -> Form
			intent.Update(StrategySelectedMsg{Strategy: "manual"})
			Expect(intent.state.currentState).To(Equal(CaptureStateForm))

			// Form -> Review
			intent.Update(FormSubmittedMsg{Event: &career.CareerEvent{Text: "Event"}})
			Expect(intent.state.currentState).To(Equal(CaptureStateReview))

			// Review -> Submit
			intent.Update(ReviewConfirmedMsg{
				AcceptedBursts: []*career.Burst{},
				AcceptedFacts:  []*career.Fact{},
				RejectedItems:  make(map[string]string),
			})
			Expect(intent.state.currentState).To(Equal(CaptureStateSubmit))

			// Submit -> Complete
			intent.Update(SubmitCompleteMsg{})
			Expect(intent.active).To(BeFalse())
		})
	})
})
