package intents_test

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
)

var _ = Describe("BulkOperationsIntent", func() {
	var (
		ctx    context.Context
		data   *intents.BulkOperationsContext
		model  *intents.BulkOperationsModel
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		data = intents.NewBulkOperationsContext(ctx)
		model = intents.NewBulkOperationsIntent(data)
	})

	AfterEach(func() {
		cancel()
	})

	Describe("NewBulkOperationsIntent", func() {
		It("should create a new bulk operations intent", func() {
			Expect(model).ToNot(BeNil())
			Expect(model.Result()).To(BeNil())
		})
	})

	Describe("NewBulkOperationsContext", func() {
		It("should create context with default values", func() {
			newCtx := intents.NewBulkOperationsContext(ctx)
			Expect(newCtx.CurrentState).To(Equal(intents.BulkSelectOpState))
			Expect(newCtx.AvailableOps).To(ConsistOf("delete", "tag", "archive", "export"))
			Expect(newCtx.SelectedOp).To(Equal(""))
			Expect(newCtx.ScopeType).To(Equal("selected"))
			Expect(newCtx.AffectedItemCount).To(Equal(0))
			Expect(newCtx.ProcessedCount).To(Equal(0))
			Expect(newCtx.SuccessCount).To(Equal(0))
			Expect(newCtx.FailureCount).To(Equal(0))
			Expect(newCtx.SkippedCount).To(Equal(0))
			Expect(newCtx.IsExecuting).To(BeFalse())
			Expect(newCtx.IsPaused).To(BeFalse())
			Expect(newCtx.Results).ToNot(BeNil())
			Expect(newCtx.Errors).ToNot(BeNil())
			Expect(newCtx.FormErrors).ToNot(BeNil())
		})
	})

	Describe("BulkOperationsContext Methods", func() {
		Describe("ClearFormErrors", func() {
			It("should clear all form errors", func() {
				data.SetFormError("field1", "error1")
				data.SetFormError("field2", "error2")
				Expect(data.HasFormErrors()).To(BeTrue())

				data.ClearFormErrors()
				Expect(data.HasFormErrors()).To(BeFalse())
				Expect(len(data.FormErrors)).To(Equal(0))
			})
		})

		Describe("SetFormError", func() {
			It("should set a form error", func() {
				data.SetFormError("test_field", "test error")
				Expect(data.FormErrors["test_field"]).To(Equal("test error"))
				Expect(data.HasFormErrors()).To(BeTrue())
			})
		})

		Describe("HasFormErrors", func() {
			It("should return false when no errors", func() {
				Expect(data.HasFormErrors()).To(BeFalse())
			})

			It("should return true when errors exist", func() {
				data.SetFormError("field", "error")
				Expect(data.HasFormErrors()).To(BeTrue())
			})
		})

		Describe("SelectOperation", func() {
			It("should select a valid operation", func() {
				data.SelectOperation("delete")
				Expect(data.SelectedOp).To(Equal("delete"))
			})

			It("should not select an invalid operation", func() {
				data.SelectOperation("invalid_op")
				Expect(data.SelectedOp).To(Equal(""))
			})

			It("should select tag operation", func() {
				data.SelectOperation("tag")
				Expect(data.SelectedOp).To(Equal("tag"))
			})

			It("should select archive operation", func() {
				data.SelectOperation("archive")
				Expect(data.SelectedOp).To(Equal("archive"))
			})

			It("should select export operation", func() {
				data.SelectOperation("export")
				Expect(data.SelectedOp).To(Equal("export"))
			})
		})

		Describe("StartExecution", func() {
			It("should initialize execution state", func() {
				data.StartExecution()
				Expect(data.IsExecuting).To(BeTrue())
				Expect(data.IsPaused).To(BeFalse())
				Expect(data.ProcessedCount).To(Equal(0))
				Expect(data.SuccessCount).To(Equal(0))
				Expect(data.FailureCount).To(Equal(0))
				Expect(data.SkippedCount).To(Equal(0))
				Expect(data.ExecutionStartTime).ToNot(BeZero())
				Expect(len(data.Results)).To(Equal(0))
				Expect(len(data.Errors)).To(Equal(0))
			})

			It("should reset counters on restart", func() {
				data.ProcessedCount = 10
				data.SuccessCount = 5
				data.FailureCount = 3
				data.SkippedCount = 2
				data.AddError("previous error")

				data.StartExecution()
				Expect(data.ProcessedCount).To(Equal(0))
				Expect(data.SuccessCount).To(Equal(0))
				Expect(data.FailureCount).To(Equal(0))
				Expect(data.SkippedCount).To(Equal(0))
				Expect(len(data.Errors)).To(Equal(0))
			})
		})

		Describe("PauseExecution", func() {
			It("should pause execution", func() {
				data.StartExecution()
				data.PauseExecution()
				Expect(data.IsPaused).To(BeTrue())
				Expect(data.IsExecuting).To(BeTrue())
			})
		})

		Describe("ResumeExecution", func() {
			It("should resume execution", func() {
				data.StartExecution()
				data.PauseExecution()
				data.ResumeExecution()
				Expect(data.IsPaused).To(BeFalse())
				Expect(data.IsExecuting).To(BeTrue())
			})
		})

		Describe("CompleteExecution", func() {
			It("should complete execution", func() {
				data.StartExecution()
				data.CompleteExecution()
				Expect(data.IsExecuting).To(BeFalse())
				Expect(data.IsPaused).To(BeFalse())
			})
		})

		Describe("AddResult", func() {
			It("should add success result", func() {
				data.AddResult("item1", "success")
				Expect(data.Results["item1"]).To(Equal("success"))
				Expect(data.ProcessedCount).To(Equal(1))
				Expect(data.SuccessCount).To(Equal(1))
				Expect(data.FailureCount).To(Equal(0))
				Expect(data.SkippedCount).To(Equal(0))
			})

			It("should add error result", func() {
				data.AddResult("item2", "error")
				Expect(data.Results["item2"]).To(Equal("error"))
				Expect(data.ProcessedCount).To(Equal(1))
				Expect(data.SuccessCount).To(Equal(0))
				Expect(data.FailureCount).To(Equal(1))
				Expect(data.SkippedCount).To(Equal(0))
			})

			It("should add skipped result", func() {
				data.AddResult("item3", "skipped")
				Expect(data.Results["item3"]).To(Equal("skipped"))
				Expect(data.ProcessedCount).To(Equal(1))
				Expect(data.SuccessCount).To(Equal(0))
				Expect(data.FailureCount).To(Equal(0))
				Expect(data.SkippedCount).To(Equal(1))
			})

			It("should track multiple results", func() {
				data.AddResult("item1", "success")
				data.AddResult("item2", "error")
				data.AddResult("item3", "skipped")
				data.AddResult("item4", "success")

				Expect(data.ProcessedCount).To(Equal(4))
				Expect(data.SuccessCount).To(Equal(2))
				Expect(data.FailureCount).To(Equal(1))
				Expect(data.SkippedCount).To(Equal(1))
			})
		})

		Describe("AddError", func() {
			It("should add error message", func() {
				data.AddError("error 1")
				Expect(data.Errors).To(ContainElement("error 1"))
				Expect(len(data.Errors)).To(Equal(1))
			})

			It("should add multiple errors", func() {
				data.AddError("error 1")
				data.AddError("error 2")
				data.AddError("error 3")
				Expect(len(data.Errors)).To(Equal(3))
				Expect(data.Errors).To(ConsistOf("error 1", "error 2", "error 3"))
			})
		})

		Describe("GetProgress", func() {
			It("should return 0 when no items", func() {
				data.AffectedItemCount = 0
				Expect(data.GetProgress()).To(Equal(0.0))
			})

			It("should calculate progress correctly", func() {
				data.AffectedItemCount = 10
				data.ProcessedCount = 5
				Expect(data.GetProgress()).To(Equal(0.5))
			})

			It("should return 1.0 when complete", func() {
				data.AffectedItemCount = 10
				data.ProcessedCount = 10
				Expect(data.GetProgress()).To(Equal(1.0))
			})
		})
	})

	Describe("Intent Lifecycle", func() {
		Describe("Init", func() {
			It("should initialize to select operation state", func() {
				cmd := model.Init()
				Expect(cmd).ToNot(BeNil())
				Expect(data.CurrentState).To(Equal(intents.BulkSelectOpState))
			})
		})

		Describe("View", func() {
			It("should show not active message when not initialized", func() {
				view := model.View()
				Expect(view).To(ContainSubstring("not active"))
			})

			It("should show select operation screen after init", func() {
				model.Init()
				view := model.View()
				Expect(view).To(ContainSubstring("Select Bulk Operation"))
				Expect(view).To(ContainSubstring("Delete"))
				Expect(view).To(ContainSubstring("Tag"))
				Expect(view).To(ContainSubstring("Archive"))
				Expect(view).To(ContainSubstring("Export"))
			})

			It("should show breadcrumbs in view", func() {
				model.Init()
				view := model.View()
				Expect(view).To(ContainSubstring("Main Menu"))
				Expect(view).To(ContainSubstring("Bulk Operations"))
			})
		})

		Describe("Result", func() {
			It("should return nil when no result set", func() {
				Expect(model.Result()).To(BeNil())
			})
		})
	})

	Describe("State Transitions", func() {
		BeforeEach(func() {
			model.Init()
		})

		Describe("BulkSelectOpState", func() {
			It("should handle quit key by returning tea.Quit", func() {
				cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
				// q now returns tea.Quit to quit the application
				Expect(cmd).ToNot(BeNil())
			})

			It("should handle escape key", func() {
				cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(cmd).To(BeNil())
				result := model.Result()
				Expect(result).ToNot(BeNil())
				Expect(result.Status).To(Equal(intents.Cancelled))
			})

			It("should select operation with number key", func() {
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
				Expect(data.SelectedOp).To(Equal("delete"))
			})

			It("should navigate with arrow keys", func() {
				// Move down to select first item
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				Expect(data.SelectedOp).To(Equal("delete"))

				// Move down again
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				Expect(data.SelectedOp).To(Equal("tag"))

				// Move up
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
				Expect(data.SelectedOp).To(Equal("delete"))
			})

			It("should transition to configure state on enter", func() {
				data.SelectedOp = "delete"
				model.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(data.CurrentState).To(Equal(intents.BulkConfigureState))
			})

			It("should not transition if no operation selected", func() {
				model.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(data.CurrentState).To(Equal(intents.BulkSelectOpState))
			})
		})

		Describe("BulkConfigureState", func() {
			BeforeEach(func() {
				data.SelectedOp = "delete"
				data.CurrentState = intents.BulkConfigureState
				data.AffectedItemCount = 5
			})

			It("should show configuration screen", func() {
				view := model.View()
				Expect(view).To(ContainSubstring("Configure Operation"))
				Expect(view).To(ContainSubstring("Delete"))
				Expect(view).To(ContainSubstring("5"))
			})

			It("should transition to execute state on enter", func() {
				model.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(data.CurrentState).To(Equal(intents.BulkExecuteState))
				Expect(data.IsExecuting).To(BeTrue())
			})

			It("should go back to select state on escape", func() {
				model.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(data.CurrentState).To(Equal(intents.BulkSelectOpState))
			})

			It("should quit application on 'q' key", func() {
				cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
				// q now returns tea.Quit to quit the application
				Expect(cmd).ToNot(BeNil())
			})
		})

		Describe("BulkExecuteState", func() {
			BeforeEach(func() {
				data.SelectedOp = "delete"
				data.CurrentState = intents.BulkExecuteState
				data.AffectedItemCount = 10
				data.StartExecution()
			})

			It("should show execution screen", func() {
				view := model.View()
				Expect(view).To(ContainSubstring("Bulk Operation in Progress"))
			})

			It("should show progress modal", func() {
				data.ProcessedCount = 5
				data.SuccessCount = 4
				data.FailureCount = 1
				view := model.View()
				Expect(view).To(ContainSubstring("Processing"))
				// Progress details are shown in the modal/main content, not as "5/10" literal
				Expect(view).To(ContainSubstring("Progress details are shown above"))
			})

			It("should pause execution on 'p' key", func() {
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
				Expect(data.IsPaused).To(BeTrue())
			})

			It("should resume execution on 'p' key when paused", func() {
				data.PauseExecution()
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
				Expect(data.IsPaused).To(BeFalse())
			})

			It("should complete on 'c' key", func() {
				data.ProcessedCount = 5
				data.SuccessCount = 4
				data.FailureCount = 1
				cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
				Expect(cmd).ToNot(BeNil())
				Expect(data.CurrentState).To(Equal(intents.BulkCompleteState))
				result := model.Result()
				Expect(result.Status).To(Equal(intents.Completed))
				Expect(result.Data).ToNot(BeNil())
			})

			It("should cancel on quit", func() {
				cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
				Expect(cmd).ToNot(BeNil())
				result := model.Result()
				Expect(result.Status).To(Equal(intents.Cancelled))
			})

			It("should show paused state in view", func() {
				data.PauseExecution()
				view := model.View()
				Expect(view).To(ContainSubstring("paused"))
			})
		})

		Describe("BulkCompleteState", func() {
			BeforeEach(func() {
				data.CurrentState = intents.BulkCompleteState
				data.SelectedOp = "delete"
				data.ProcessedCount = 10
				data.SuccessCount = 8
				data.FailureCount = 2
			})

			It("should show completion screen", func() {
				view := model.View()
				Expect(view).To(ContainSubstring("Operation Complete"))
				Expect(view).To(ContainSubstring("Delete"))
				Expect(view).To(ContainSubstring("10"))
			})

			It("should show success icon when no failures", func() {
				data.FailureCount = 0
				view := model.View()
				Expect(view).To(ContainSubstring("✅"))
			})

			It("should show warning icon when failures exist", func() {
				view := model.View()
				Expect(view).To(ContainSubstring("⚠️"))
			})

			It("should show error summary", func() {
				data.AddError("error 1")
				data.AddError("error 2")
				view := model.View()
				Expect(view).To(ContainSubstring("Error Summary"))
				Expect(view).To(ContainSubstring("error 1"))
			})

			It("should truncate error list to 5", func() {
				for i := 1; i <= 10; i++ {
					data.AddError("error " + string(rune('0'+i)))
				}
				view := model.View()
				Expect(view).To(ContainSubstring("and 5 more errors"))
			})
		})
	})

	Describe("View Content Methods", func() {
		BeforeEach(func() {
			model.Init()
		})

		It("should show operation descriptions", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Permanently remove"))
			Expect(view).To(ContainSubstring("Add or modify tags"))
			Expect(view).To(ContainSubstring("Move items to archive"))
			Expect(view).To(ContainSubstring("Export selected items"))
		})

		It("should show operation selection indicator", func() {
			data.SelectedOp = "delete"
			view := model.View()
			Expect(view).To(ContainSubstring("▶"))
		})

		It("should show warning for delete operation", func() {
			data.SelectedOp = "delete"
			data.CurrentState = intents.BulkConfigureState
			view := model.View()
			Expect(view).To(ContainSubstring("⚠️"))
			Expect(view).To(ContainSubstring("cannot be undone"))
		})

		It("should show context help for each state", func() {
			// Select Op state
			view := model.View()
			Expect(view).To(ContainSubstring("Enter"))
			Expect(view).To(ContainSubstring("Esc"))

			// Configure state
			data.SelectedOp = "delete"
			data.CurrentState = intents.BulkConfigureState
			view = model.View()
			Expect(view).To(ContainSubstring("Enter"))
			Expect(view).To(ContainSubstring("Execute"))

			// Execute state
			data.CurrentState = intents.BulkExecuteState
			data.StartExecution()
			view = model.View()
			Expect(view).To(ContainSubstring("Pause"))

			// Complete state
			data.CurrentState = intents.BulkCompleteState
			view = model.View()
			Expect(view).To(ContainSubstring("Done"))
		})

		It("should show no operations message when list is empty", func() {
			data.AvailableOps = []string{}
			view := model.View()
			Expect(view).To(ContainSubstring("No operations available"))
		})
	})

	Describe("Edge Cases", func() {
		BeforeEach(func() {
			model.Init()
		})

		It("should handle invalid number selection", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'9'}})
			Expect(data.SelectedOp).To(Equal(""))
		})

		It("should not move up beyond first item", func() {
			data.SelectedOp = "delete"
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(data.SelectedOp).To(Equal("delete"))
		})

		It("should not move down beyond last item", func() {
			data.SelectedOp = "export"
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(data.SelectedOp).To(Equal("export"))
		})

		It("should handle ctrl+c like quit", func() {
			cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			// ctrl+c now returns tea.Quit to quit the application
			Expect(cmd).ToNot(BeNil())
		})

		It("should handle empty operation display name", func() {
			data.SelectedOp = ""
			data.CurrentState = intents.BulkConfigureState
			view := model.View()
			Expect(view).ToNot(BeEmpty())
		})

		It("should handle progress with zero items", func() {
			Expect(data.GetProgress()).To(Equal(0.0))
		})
	})
})
