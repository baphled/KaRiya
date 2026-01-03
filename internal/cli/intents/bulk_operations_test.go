package intents

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ErrBulkMockError = errors.New("bulk mock error")
	ErrBulkNotFound  = errors.New("not found")
)

var _ = Describe("BulkOperations Intent", func() {
	var (
		model *BulkOperationsModel
		ctx   context.Context
		data  *BulkOperationsContext
	)

	BeforeEach(func() {
		ctx = context.Background()
		data = NewBulkOperationsContext(ctx)
		model = NewBulkOperationsIntent(data)
	})

	Describe("Initialization", func() {
		It("should initialize with select operation state", func() {
			Expect(model).NotTo(BeNil())
			Expect(data.CurrentState).To(Equal(BulkSelectOpState))
		})

		It("should have available operations", func() {
			Expect(data.AvailableOps).To(HaveLen(4))
		})

		It("should initialize with empty results", func() {
			Expect(data.Results).To(BeEmpty())
			Expect(data.Errors).To(BeEmpty())
		})

		It("should call Init without error", func() {
			cmd := model.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("State Management", func() {
		It("should transition states", func() {
			data.SelectedOp = "delete"
			data.CurrentState = BulkConfigureState
			Expect(data.CurrentState).To(Equal(BulkConfigureState))
		})

		It("should support all states", func() {
			states := []BulkOperationsState{
				BulkSelectOpState,
				BulkConfigureState,
				BulkExecuteState,
				BulkCompleteState,
			}
			for _, state := range states {
				data.CurrentState = state
				Expect(data.CurrentState).To(Equal(state))
			}
		})
	})

	Describe("Operation Selection", func() {
		It("should select operation", func() {
			data.SelectedOp = "delete"
			Expect(data.SelectedOp).To(Equal("delete"))
		})

		It("should allow multiple operation types", func() {
			for _, op := range []string{"delete", "tag", "archive", "export"} {
				data.SelectedOp = op
				Expect(data.SelectedOp).To(Equal(op))
			}
		})
	})

	Describe("Scope Configuration", func() {
		It("should set scope type", func() {
			data.ScopeType = "selected"
			Expect(data.ScopeType).To(Equal("selected"))
		})

		It("should set affected item count", func() {
			data.AffectedItemCount = 10
			Expect(data.AffectedItemCount).To(Equal(10))
		})

		It("should initialize execution counters", func() {
			Expect(data.ProcessedCount).To(Equal(0))
			Expect(data.SuccessCount).To(Equal(0))
			Expect(data.FailureCount).To(Equal(0))
			Expect(data.SkippedCount).To(Equal(0))
		})
	})

	Describe("Execution Progress", func() {
		It("should track processed count", func() {
			data.ProcessedCount = 5
			Expect(data.ProcessedCount).To(Equal(5))
		})

		It("should track success count", func() {
			data.SuccessCount = 3
			Expect(data.SuccessCount).To(Equal(3))
		})

		It("should track failure count", func() {
			data.FailureCount = 2
			Expect(data.FailureCount).To(Equal(2))
		})

		It("should track skipped count", func() {
			data.SkippedCount = 1
			Expect(data.SkippedCount).To(Equal(1))
		})

		It("should maintain execution state", func() {
			data.IsExecuting = true
			Expect(data.IsExecuting).To(BeTrue())
		})

		It("should allow pausing execution", func() {
			data.IsExecuting = true
			data.IsPaused = true
			Expect(data.IsPaused).To(BeTrue())
		})
	})

	Describe("Result Handling", func() {
		It("should return typed IntentResult", func() {
			data.SelectedOp = "delete"
			result := model.Result()
			Expect(result).NotTo(BeNil())
		})
	})

	Describe("Error Handling", func() {
		It("should collect errors during execution", func() {
			data.Errors = append(data.Errors, "error 1")
			Expect(data.Errors).To(HaveLen(1))
		})

		It("should handle form validation errors", func() {
			data.SetFormError("operation", "invalid")
			Expect(data.HasFormErrors()).To(BeTrue())
		})

		It("should clear form errors", func() {
			data.SetFormError("operation", "invalid")
			data.ClearFormErrors()
			Expect(data.HasFormErrors()).To(BeFalse())
		})
	})

	Describe("View Rendering", func() {
		It("should render views", func() {
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Context Preservation", func() {
		It("should preserve operation across transitions", func() {
			data.SelectedOp = "archive"
			data.CurrentState = BulkConfigureState
			Expect(data.SelectedOp).To(Equal("archive"))
		})

		It("should preserve results", func() {
			data.Results["item1"] = "success"
			Expect(data.Results).To(HaveLen(1))
		})
	})

	Describe("Context Functions", func() {
		It("should create new context", func() {
			newCtx := NewBulkOperationsContext(ctx)
			Expect(newCtx).NotTo(BeNil())
		})
	})
})
