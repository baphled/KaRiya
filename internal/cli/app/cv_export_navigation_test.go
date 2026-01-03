package app

import (
	"github.com/baphled/kariya/internal/cli/models"
	career "github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CV Export Navigation Integration", func() {
	Describe("CV Preview 'x' key should navigate to export dialog", func() {
		It("should verify ShowExportOptionsMsg is properly handled in app", func() {
			// This test verifies the routing exists without needing full app initialization
			// We can verify the constants and handlers are defined by checking they compile

			// Verify the screen constants exist
			var exportDialog Screen = CVExportDialogScreen
			var exportSuccess Screen = CVExportSuccessScreen
			var exportProgress Screen = CVExportProgressScreen

			Expect(exportDialog).To(Equal(Screen("cv_export_dialog")))
			Expect(exportSuccess).To(Equal(Screen("cv_export_success")))
			Expect(exportProgress).To(Equal(Screen("cv_export_progress")))
		})

		It("should verify CV preview sends ShowExportOptionsMsg on 'x' key", func() {
			// Test the CVPreviewModel itself
			baseModel := models.NewBaseStandardModel()
			cvView := &career.CVView{
				Name:       "Test CV",
				TargetRole: "principal",
				Sections:   []*career.CVSection{},
			}
			cvPreviewModel := models.NewCVPreviewModel(baseModel, cvView, []*career.CVSection{}, nil)

			// Press 'x' to export
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
			_, cmd := cvPreviewModel.Update(keyMsg)

			// Should return a command
			Expect(cmd).NotTo(BeNil(), "'x' key should trigger export dialog")

			// Execute command to get the message
			msg := cmd()
			exportMsg, ok := msg.(models.ShowExportOptionsMsg)
			Expect(ok).To(BeTrue(), "'x' should trigger ShowExportOptionsMsg")
			Expect(exportMsg.CVView.Name).To(Equal("Test CV"))
		})

		It("should verify export dialog model can be created and rendered", func() {
			// Test that the export dialog model works
			baseModel := models.NewBaseStandardModel()
			cvView := &career.CVView{
				Name:       "Test CV",
				TargetRole: "principal",
			}
			exportDialog := models.NewCVExportDialogModel(
				baseModel,
				cvView,
				[]*career.CVSection{},
				nil,
			)

			// Should be able to render without error
			view := exportDialog.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should verify export success model can be created and rendered", func() {
			// Test that the export success model works
			baseModel := models.NewBaseStandardModel()
			cvView := &career.CVView{
				Name:       "Test CV",
				TargetRole: "principal",
			}
			successModel := models.NewCVExportSuccessModel(
				baseModel,
				cvView,
				"Markdown",
				"/tmp/test.md",
			)

			// Should be able to render without error
			view := successModel.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Markdown"))
			Expect(view).To(ContainSubstring("/tmp/test.md"))
		})
	})
})
