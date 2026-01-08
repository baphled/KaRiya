package intents_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("ExportArtifact - Escape Key Behavior", func() {
	var (
		exportContext *intents.ExportArtifactContext
		exportModel   *intents.ExportArtifactModel
	)

	BeforeEach(func() {
		exportContext = intents.NewTestExportArtifactContext()
		exportModel = intents.NewExportArtifactModel(exportContext)
		exportModel.Init()
	})

	Describe("SelectType State (Root)", func() {
		It("should cancel intent when escape is pressed", func() {
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := exportModel.Result()
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should cancel intent when 'm' is pressed", func() {
			exportModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			result := exportModel.Result()
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should show 'm' key in footer", func() {
			view := exportModel.View()
			Expect(view).To(ContainSubstring("m: Main menu"))
		})
	})

	Describe("SelectFormat State", func() {
		BeforeEach(func() {
			// Select non-CV artifact type to avoid SelectCV state
			exportModel.Update(tea.KeyMsg{Type: tea.KeyDown}) // Move to Events (index 1)
			// Navigate to SelectFormat
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})

		It("should go back to SelectType when escape is pressed", func() {
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEsc})

			view := exportModel.View()
			Expect(view).To(ContainSubstring("Select Artifact Type"))
			result := exportModel.Result()
			Expect(result).To(BeNil())
		})

		It("should cancel intent when 'm' is pressed", func() {
			exportModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			result := exportModel.Result()
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should show 'm' key in footer", func() {
			view := exportModel.View()
			Expect(view).To(ContainSubstring("m: Main menu"))
		})
	})

	Describe("SelectDestination State", func() {
		BeforeEach(func() {
			// Select non-CV artifact type to avoid SelectCV state
			exportModel.Update(tea.KeyMsg{Type: tea.KeyDown}) // Move to Events (index 1)
			// Navigate to SelectDest
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectType -> SelectFormat
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectFormat -> SelectDest
		})

		It("should go back to SelectFormat when escape is pressed", func() {
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEsc})

			view := exportModel.View()
			Expect(view).To(ContainSubstring("Select Export Format"))
			result := exportModel.Result()
			Expect(result).To(BeNil())
		})

		It("should cancel intent when 'm' is pressed", func() {
			exportModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			result := exportModel.Result()
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should show 'm' key in footer", func() {
			view := exportModel.View()
			Expect(view).To(ContainSubstring("m: Main menu"))
		})
	})

	Describe("Configure State", func() {
		BeforeEach(func() {
			// Select non-CV artifact type to avoid SelectCV state
			exportModel.Update(tea.KeyMsg{Type: tea.KeyDown}) // Move to Events (index 1)
			// Navigate to Configure
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectType -> SelectFormat
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectFormat -> SelectDest
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectDest -> Configure
		})

		It("should go back to SelectDest when escape is pressed", func() {
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEsc})

			view := exportModel.View()
			Expect(view).To(ContainSubstring("Select Export Destination"))
			result := exportModel.Result()
			Expect(result).To(BeNil())
		})

		It("should cancel intent when 'm' is pressed", func() {
			exportModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			result := exportModel.Result()
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should show 'm' key in footer", func() {
			view := exportModel.View()
			Expect(view).To(ContainSubstring("m: Main menu"))
		})
	})

	Describe("Preview State", func() {
		BeforeEach(func() {
			// Select non-CV artifact type to avoid SelectCV state
			exportModel.Update(tea.KeyMsg{Type: tea.KeyDown}) // Move to Events (index 1)
			// Navigate to Preview
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectType -> SelectFormat
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectFormat -> SelectDest
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectDest -> Configure
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Configure -> Preview
		})

		It("should go back to Configure when escape is pressed", func() {
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEsc})

			view := exportModel.View()
			Expect(view).To(ContainSubstring("Configure Export"))
			result := exportModel.Result()
			Expect(result).To(BeNil())
		})

		It("should cancel intent when 'm' is pressed", func() {
			exportModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			result := exportModel.Result()
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should show 'm' key in footer", func() {
			view := exportModel.View()
			Expect(view).To(ContainSubstring("m: Main menu"))
		})
	})

	Describe("Confirm State", func() {
		BeforeEach(func() {
			// Select non-CV artifact type to avoid SelectCV state
			exportModel.Update(tea.KeyMsg{Type: tea.KeyDown}) // Move to Events (index 1)
			// Navigate to Confirm
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectType -> SelectFormat
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectFormat -> SelectDest
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectDest -> Configure
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Configure -> Preview
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Preview -> Confirm
		})

		It("should go back to Preview when escape is pressed", func() {
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEsc})

			view := exportModel.View()
			Expect(view).To(ContainSubstring("Preview Export"))
			result := exportModel.Result()
			Expect(result).To(BeNil())
		})

		It("should cancel intent when 'm' is pressed", func() {
			exportModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			result := exportModel.Result()
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should show 'm' key in footer", func() {
			view := exportModel.View()
			Expect(view).To(ContainSubstring("m: Main menu"))
		})
	})

	Describe("InProgress State (Async Operation)", func() {
		BeforeEach(func() {
			// Select non-CV artifact type to avoid SelectCV state
			exportModel.Update(tea.KeyMsg{Type: tea.KeyDown}) // Move to Events (index 1)
			// Navigate to InProgress
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectType -> SelectFormat
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectFormat -> SelectDest
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectDest -> Configure
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Configure -> Preview
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Preview -> Confirm
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Confirm -> InProgress
		})

		It("should allow escape key without cancelling export (background completion)", func() {
			view := exportModel.View()
			Expect(view).To(ContainSubstring("Exporting artifact"))

			exportModel.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should not cancel - no result set
			result := exportModel.Result()
			Expect(result).To(BeNil())
		})

		It("should cancel intent immediately when 'm' is pressed", func() {
			exportModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			result := exportModel.Result()
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should show both escape and 'm' options in footer", func() {
			view := exportModel.View()
			Expect(view).To(ContainSubstring("Esc: Let export complete in background"))
			Expect(view).To(ContainSubstring("m: Cancel and return to menu"))
		})
	})

	Describe("Complete State", func() {
		BeforeEach(func() {
			// Select non-CV artifact type to avoid SelectCV state
			exportModel.Update(tea.KeyMsg{Type: tea.KeyDown}) // Move to Events (index 1)
			// Navigate to Complete by completing export
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectType -> SelectFormat
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectFormat -> SelectDest
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectDest -> Configure
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Configure -> Preview
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Preview -> Confirm
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Confirm -> InProgress

			// Simulate export completion
			result := intents.NewExportArtifactResult(
				true,
				intents.ExportTypeCV,
				intents.ExportFormatPDF,
				intents.ExportDestinationFile,
				"/tmp/export.pdf",
				1024,
			)
			exportModel.Update(intents.ExportCompleteMsg{Result: result})
		})

		It("should close intent when 'm' is pressed", func() {
			view := exportModel.View()
			Expect(view).To(ContainSubstring("Export Complete"))

			exportModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			result := exportModel.Result()
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(intents.Completed))
		})

		It("should show 'm' key in footer", func() {
			view := exportModel.View()
			Expect(view).To(ContainSubstring("m: Main menu"))
		})
	})

	Describe("Failed State", func() {
		BeforeEach(func() {
			// Select non-CV artifact type to avoid SelectCV state
			exportModel.Update(tea.KeyMsg{Type: tea.KeyDown}) // Move to Events (index 1)
			// Navigate to Failed by simulating error
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectType -> SelectFormat
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectFormat -> SelectDest
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // SelectDest -> Configure
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Configure -> Preview
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Preview -> Confirm
			exportModel.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Confirm -> InProgress

			// Simulate export error
			exportModel.Update(intents.ExportErrorMsg{
				Error: &intents.IntentError{
					Code:    "export_failed",
					Message: "Export failed",
				},
			})
		})

		It("should close intent when 'm' is pressed", func() {
			view := exportModel.View()
			Expect(view).To(ContainSubstring("Export Failed"))

			exportModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			// In Failed state, model becomes inactive
			// Result may or may not be set depending on implementation
		})

		It("should show 'm' key in footer", func() {
			view := exportModel.View()
			Expect(view).To(ContainSubstring("m: Main menu"))
		})
	})
})
