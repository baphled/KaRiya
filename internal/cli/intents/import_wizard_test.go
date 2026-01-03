package intents_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
)

func TestImportWizard(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ImportWizard Intent Suite")
}

var _ = Describe("ImportWizard Intent", func() {
	var (
		model *intents.ImportWizardModel
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		data := intents.NewImportWizardContext(ctx)
		model = intents.NewImportWizardIntent(data)
	})

	Describe("Initialization", func() {
		It("should initialize with file select state", func() {
			cmd := model.Init(ctx)
			Expect(cmd).To(BeNil())
			Expect(model.View()).NotTo(BeEmpty())
		})
	})

	Describe("File Select State", func() {
		It("should render file select view", func() {
			model.Init(ctx)
			view := model.View()
			Expect(view).To(ContainSubstring("Select File"))
		})

		It("should show no file selected initially", func() {
			model.Init(ctx)
			view := model.View()
			Expect(view).To(ContainSubstring("none selected"))
		})
	})

	Describe("Context Operations", func() {
		It("should track form errors", func() {
			data := intents.NewImportWizardContext(ctx)
			data.SetFormError("file", "File not found")
			Expect(data.HasFormErrors()).To(BeTrue())
		})

		It("should clear form errors", func() {
			data := intents.NewImportWizardContext(ctx)
			data.SetFormError("file", "File not found")
			data.ClearFormErrors()
			Expect(data.HasFormErrors()).To(BeFalse())
		})

		It("should start import", func() {
			data := intents.NewImportWizardContext(ctx)
			data.TotalRows = 100
			data.StartImport()
			Expect(data.IsImporting).To(BeTrue())
			Expect(data.IsPaused).To(BeFalse())
		})

		It("should pause and resume import", func() {
			data := intents.NewImportWizardContext(ctx)
			data.StartImport()
			data.PauseImport()
			Expect(data.IsPaused).To(BeTrue())
			data.ResumeImport()
			Expect(data.IsPaused).To(BeFalse())
		})

		It("should cancel import", func() {
			data := intents.NewImportWizardContext(ctx)
			data.StartImport()
			data.CancelImport()
			Expect(data.IsImporting).To(BeFalse())
		})

		It("should track processed rows", func() {
			data := intents.NewImportWizardContext(ctx)
			data.TotalRows = 100
			data.IncrementProcessed()
			data.IncrementSuccessful()
			Expect(data.ProcessedRows).To(Equal(1))
			Expect(data.SuccessfulRows).To(Equal(1))
		})

		It("should add errors", func() {
			data := intents.NewImportWizardContext(ctx)
			data.AddError("Row 1: Invalid data")
			Expect(data.ErrorRows).To(Equal(1))
			Expect(data.Errors).To(HaveLen(1))
		})

		It("should calculate progress", func() {
			data := intents.NewImportWizardContext(ctx)
			data.TotalRows = 100
			data.ProcessedRows = 50
			progress := data.GetProgress()
			Expect(progress).To(Equal(0.5))
		})

		It("should handle zero total rows", func() {
			data := intents.NewImportWizardContext(ctx)
			progress := data.GetProgress()
			Expect(progress).To(Equal(0.0))
		})
	})

	Describe("View Methods", func() {
		It("should render view without error", func() {
			model.Init(ctx)
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Result Handling", func() {
		It("should return result", func() {
			result := model.Result()
			Expect(result).NotTo(BeNil())
		})
	})
})
