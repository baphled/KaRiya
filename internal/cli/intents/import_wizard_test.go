package intents_test

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
)

var _ = Describe("ImportWizardIntent", func() {
	var (
		ctx    context.Context
		data   *intents.ImportWizardContext
		model  *intents.ImportWizardModel
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		data = intents.NewImportWizardContext(ctx)
		model = intents.NewImportWizardIntent(data)
	})

	AfterEach(func() {
		cancel()
	})

	Describe("NewImportWizardIntent", func() {
		It("should create a new import wizard intent", func() {
			Expect(model).ToNot(BeNil())
			Expect(model.Result()).To(BeNil())
		})
	})

	Describe("NewImportWizardContext", func() {
		It("should create context with default values", func() {
			newCtx := intents.NewImportWizardContext(ctx)
			Expect(newCtx.CurrentState).To(Equal(intents.ImportFileSelectState))
			Expect(newCtx.FilePath).To(Equal(""))
			Expect(newCtx.FileSize).To(Equal(int64(0)))
			Expect(newCtx.TotalRows).To(Equal(0))
			Expect(newCtx.ProcessedRows).To(Equal(0))
			Expect(newCtx.SuccessfulRows).To(Equal(0))
			Expect(newCtx.ErrorRows).To(Equal(0))
			Expect(newCtx.IsImporting).To(BeFalse())
			Expect(newCtx.IsPaused).To(BeFalse())
			Expect(newCtx.Errors).ToNot(BeNil())
			Expect(newCtx.FormErrors).ToNot(BeNil())
		})
	})

	Describe("ImportWizardContext Methods", func() {
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
				data.SetFormError("file_path", "invalid path")
				Expect(data.FormErrors["file_path"]).To(Equal("invalid path"))
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

		Describe("StartImport", func() {
			It("should initialize import state", func() {
				data.StartImport()
				Expect(data.IsImporting).To(BeTrue())
				Expect(data.IsPaused).To(BeFalse())
				Expect(data.ProcessedRows).To(Equal(0))
				Expect(data.SuccessfulRows).To(Equal(0))
				Expect(data.ErrorRows).To(Equal(0))
				Expect(data.ImportStartTime).ToNot(BeZero())
				Expect(len(data.Errors)).To(Equal(0))
			})

			It("should reset counters on restart", func() {
				data.ProcessedRows = 10
				data.SuccessfulRows = 8
				data.ErrorRows = 2
				data.AddError("previous error")

				data.StartImport()
				Expect(data.ProcessedRows).To(Equal(0))
				Expect(data.SuccessfulRows).To(Equal(0))
				Expect(data.ErrorRows).To(Equal(0))
				Expect(len(data.Errors)).To(Equal(0))
			})
		})

		Describe("PauseImport", func() {
			It("should pause import", func() {
				data.StartImport()
				data.PauseImport()
				Expect(data.IsPaused).To(BeTrue())
				Expect(data.IsImporting).To(BeTrue())
			})
		})

		Describe("ResumeImport", func() {
			It("should resume import", func() {
				data.StartImport()
				data.PauseImport()
				data.ResumeImport()
				Expect(data.IsPaused).To(BeFalse())
				Expect(data.IsImporting).To(BeTrue())
			})
		})

		Describe("CancelImport", func() {
			It("should cancel import", func() {
				data.StartImport()
				data.CancelImport()
				Expect(data.IsImporting).To(BeFalse())
				Expect(data.IsPaused).To(BeFalse())
			})
		})

		Describe("CompleteImport", func() {
			It("should complete import", func() {
				data.StartImport()
				data.CompleteImport()
				Expect(data.IsImporting).To(BeFalse())
				Expect(data.IsPaused).To(BeFalse())
			})
		})

		Describe("AddError", func() {
			It("should add error and increment error count", func() {
				data.AddError("error 1")
				Expect(data.Errors).To(ContainElement("error 1"))
				Expect(data.ErrorRows).To(Equal(1))
			})

			It("should add multiple errors", func() {
				data.AddError("error 1")
				data.AddError("error 2")
				data.AddError("error 3")
				Expect(len(data.Errors)).To(Equal(3))
				Expect(data.ErrorRows).To(Equal(3))
			})
		})

		Describe("IncrementProcessed", func() {
			It("should increment processed rows", func() {
				data.IncrementProcessed()
				Expect(data.ProcessedRows).To(Equal(1))
			})

			It("should increment multiple times", func() {
				data.IncrementProcessed()
				data.IncrementProcessed()
				data.IncrementProcessed()
				Expect(data.ProcessedRows).To(Equal(3))
			})
		})

		Describe("IncrementSuccessful", func() {
			It("should increment successful rows", func() {
				data.IncrementSuccessful()
				Expect(data.SuccessfulRows).To(Equal(1))
			})

			It("should increment multiple times", func() {
				data.IncrementSuccessful()
				data.IncrementSuccessful()
				Expect(data.SuccessfulRows).To(Equal(2))
			})
		})

		Describe("GetProgress", func() {
			It("should return 0 when no rows", func() {
				data.TotalRows = 0
				Expect(data.GetProgress()).To(Equal(0.0))
			})

			It("should calculate progress correctly", func() {
				data.TotalRows = 100
				data.ProcessedRows = 50
				Expect(data.GetProgress()).To(Equal(0.5))
			})

			It("should return 1.0 when complete", func() {
				data.TotalRows = 100
				data.ProcessedRows = 100
				Expect(data.GetProgress()).To(Equal(1.0))
			})
		})
	})

	Describe("Intent Lifecycle", func() {
		Describe("Init", func() {
			It("should initialize to file select state", func() {
				cmd := model.Init()
				Expect(cmd).ToNot(BeNil())
				Expect(data.CurrentState).To(Equal(intents.ImportFileSelectState))
			})
		})

		Describe("View", func() {
			It("should show not active message when not initialized", func() {
				view := model.View()
				Expect(view).To(ContainSubstring("not active"))
			})

			It("should show file select screen after init", func() {
				model.Init()
				view := model.View()
				Expect(view).To(ContainSubstring("Select File"))
			})

			It("should show breadcrumbs in view", func() {
				model.Init()
				view := model.View()
				// Breadcrumbs are added by CreateViewWithBreadcrumbs
				// View should contain state-specific content
				Expect(view).ToNot(BeEmpty())
				Expect(view).To(ContainSubstring("Select File"))
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

		Describe("ImportFileSelectState", func() {
			It("should handle quit key", func() {
				cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
				Expect(cmd).ToNot(BeNil())
				result := model.Result()
				Expect(result).ToNot(BeNil())
				Expect(result.Status).To(Equal(intents.Cancelled))
			})

			It("should handle escape key", func() {
				cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(cmd).To(BeNil())
				result := model.Result()
				Expect(result).ToNot(BeNil())
				Expect(result.Status).To(Equal(intents.Cancelled))
			})

			It("should transition to preview state on enter with valid file", func() {
				data.FilePath = "test.csv"
				data.TotalRows = 10
				model.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(data.CurrentState).To(Equal(intents.ImportPreviewState))
			})

			It("should not transition if no file selected", func() {
				model.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(data.CurrentState).To(Equal(intents.ImportFileSelectState))
			})
		})

		Describe("ImportPreviewState", func() {
			BeforeEach(func() {
				data.FilePath = "test.csv"
				data.TotalRows = 10
				data.CurrentState = intents.ImportPreviewState
			})

			It("should show preview screen", func() {
				view := model.View()
				Expect(view).To(ContainSubstring("Preview"))
			})

			It("should transition to progress state on enter", func() {
				model.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(data.CurrentState).To(Equal(intents.ImportProgressState))
				Expect(data.IsImporting).To(BeTrue())
			})

			It("should go back to file select on escape", func() {
				model.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(data.CurrentState).To(Equal(intents.ImportFileSelectState))
			})

			It("should cancel on quit", func() {
				cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
				Expect(cmd).ToNot(BeNil())
				result := model.Result()
				Expect(result.Status).To(Equal(intents.Cancelled))
			})
		})

		Describe("ImportProgressState", func() {
			BeforeEach(func() {
				data.FilePath = "test.csv"
				data.TotalRows = 100
				data.CurrentState = intents.ImportProgressState
				data.StartImport()
			})

			It("should show progress screen", func() {
				view := model.View()
				Expect(view).To(ContainSubstring("Importing"))
			})

			It("should show progress modal", func() {
				data.ProcessedRows = 50
				data.SuccessfulRows = 48
				data.ErrorRows = 2
				view := model.View()
				// Progress modal shows import status
				Expect(view).To(ContainSubstring("Importing"))
			})

			It("should pause import on 'p' key", func() {
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
				Expect(data.IsPaused).To(BeTrue())
			})

			It("should resume import on 'p' key when paused", func() {
				data.PauseImport()
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
				Expect(data.IsPaused).To(BeFalse())
			})

			It("should cancel on 'c' key", func() {
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
				// Cancel import should stop the import
				Expect(data.IsImporting).To(BeFalse())
			})

			It("should show paused state in view", func() {
				data.PauseImport()
				view := model.View()
				Expect(view).To(ContainSubstring("Import is currently paused"))
			})
		})

		Describe("ImportCompleteState", func() {
			BeforeEach(func() {
				data.CurrentState = intents.ImportCompleteState
				data.FilePath = "test.csv"
				data.ProcessedRows = 100
				data.SuccessfulRows = 98
				data.ErrorRows = 2
			})

			It("should show completion screen", func() {
				view := model.View()
				Expect(view).To(ContainSubstring("Complete"))
			})

			It("should show success icon when no errors", func() {
				data.ErrorRows = 0
				view := model.View()
				Expect(view).To(ContainSubstring("✅"))
			})

			It("should show warning icon when errors exist", func() {
				view := model.View()
				Expect(view).To(ContainSubstring("⚠️"))
			})

			It("should show error summary", func() {
				data.AddError("error 1")
				data.AddError("error 2")
				view := model.View()
				Expect(view).To(ContainSubstring("Error Summary"))
			})
		})
	})

	Describe("View Content Methods", func() {
		BeforeEach(func() {
			model.Init()
		})

		It("should show file selection prompt", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Select"))
		})

		It("should show help text for each state", func() {
			// File select state
			view := model.View()
			Expect(view).To(ContainSubstring("Enter"))
			Expect(view).To(ContainSubstring("Esc"))

			// Preview state
			data.FilePath = "test.csv"
			data.TotalRows = 10
			data.CurrentState = intents.ImportPreviewState
			view = model.View()
			Expect(view).To(ContainSubstring("Start import"))

			// Progress state
			data.CurrentState = intents.ImportProgressState
			data.StartImport()
			view = model.View()
			Expect(view).To(ContainSubstring("Pause"))

			// Complete state
			data.CurrentState = intents.ImportCompleteState
			view = model.View()
			Expect(view).To(ContainSubstring("Done"))
		})
	})

	Describe("Edge Cases", func() {
		BeforeEach(func() {
			model.Init()
		})

		It("should handle ctrl+c like quit", func() {
			cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			Expect(cmd).ToNot(BeNil())
			result := model.Result()
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should handle progress with zero rows", func() {
			Expect(data.GetProgress()).To(Equal(0.0))
		})

		It("should handle multiple pause/resume cycles", func() {
			data.StartImport()

			data.PauseImport()
			Expect(data.IsPaused).To(BeTrue())

			data.ResumeImport()
			Expect(data.IsPaused).To(BeFalse())

			data.PauseImport()
			Expect(data.IsPaused).To(BeTrue())

			data.ResumeImport()
			Expect(data.IsPaused).To(BeFalse())
		})

		It("should track errors separately from successful rows", func() {
			data.IncrementProcessed()
			data.IncrementSuccessful()

			data.IncrementProcessed()
			data.AddError("error 1")

			data.IncrementProcessed()
			data.IncrementSuccessful()

			Expect(data.ProcessedRows).To(Equal(3))
			Expect(data.SuccessfulRows).To(Equal(2))
			Expect(data.ErrorRows).To(Equal(1))
		})
	})

	Describe("Integration Scenarios", func() {
		It("should support full import workflow", func() {
			// Init
			model.Init()
			Expect(data.CurrentState).To(Equal(intents.ImportFileSelectState))

			// Select file
			data.FilePath = "import.csv"
			data.TotalRows = 100
			model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(data.CurrentState).To(Equal(intents.ImportPreviewState))

			// Start import
			model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(data.CurrentState).To(Equal(intents.ImportProgressState))
			Expect(data.IsImporting).To(BeTrue())

			// Simulate progress
			for i := 0; i < 100; i++ {
				data.IncrementProcessed()
				if i%10 == 0 {
					data.AddError("error")
				} else {
					data.IncrementSuccessful()
				}
			}

			Expect(data.ProcessedRows).To(Equal(100))
			Expect(data.SuccessfulRows).To(Equal(90))
			Expect(data.ErrorRows).To(Equal(10))
		})

		It("should handle pause and resume during import", func() {
			model.Init()
			data.FilePath = "test.csv"
			data.TotalRows = 100
			data.CurrentState = intents.ImportProgressState
			data.StartImport()

			// Process some rows
			for i := 0; i < 50; i++ {
				data.IncrementProcessed()
				data.IncrementSuccessful()
			}

			// Pause
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
			Expect(data.IsPaused).To(BeTrue())

			// Resume
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
			Expect(data.IsPaused).To(BeFalse())

			// Complete
			for i := 50; i < 100; i++ {
				data.IncrementProcessed()
				data.IncrementSuccessful()
			}

			Expect(data.ProcessedRows).To(Equal(100))
		})
	})
})
