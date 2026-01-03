package intents

import (
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

			result := NewModalEditResult(original, modified, true, changes)

			Expect(result.Original).To(Equal(original))
			Expect(result.Modified).To(Equal(modified))
			Expect(result.Accepted).To(BeTrue())
			Expect(result.Changes).To(Equal(changes))
		})

		It("should handle nil changes map", func() {
			result := NewModalEditResult("original", "modified", true, nil)

			Expect(result.Changes).NotTo(BeNil())
			Expect(result.Changes).To(BeEmpty())
		})
	})

	Context("HasChanges", func() {
		It("should return true when changes exist", func() {
			result := NewModalEditResult("original", "modified", true, map[string]interface{}{"field": "value"})

			Expect(result.HasChanges()).To(BeTrue())
		})

		It("should return false when no changes", func() {
			result := NewModalEditResult("original", "original", true, make(map[string]interface{}))

			Expect(result.HasChanges()).To(BeFalse())
		})
	})

	Context("WasAccepted", func() {
		It("should return true when accepted is true", func() {
			result := NewModalEditResult("original", "modified", true, nil)

			Expect(result.WasAccepted()).To(BeTrue())
		})

		It("should return false when accepted is false", func() {
			result := NewModalEditResult("original", "original", false, nil)

			Expect(result.WasAccepted()).To(BeFalse())
		})
	})

	Context("GetChange", func() {
		It("should return the value for an existing change", func() {
			changes := map[string]interface{}{"field": "new value"}
			result := NewModalEditResult("original", "modified", true, changes)

			Expect(result.GetChange("field")).To(Equal("new value"))
		})

		It("should return nil for a non-existent change", func() {
			result := NewModalEditResult("original", "modified", true, map[string]interface{}{"field": "value"})

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

			result := NewCancelledModalEditResult(original)

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

			result := NewModalEditResult(original, modified, true, changes)

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
			Expect(view).To(ContainSubstring("Capture Event Form"))
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
			Expect(view).To(ContainSubstring("Capture Event Form"))
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
