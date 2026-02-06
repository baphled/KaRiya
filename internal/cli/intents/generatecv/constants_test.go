package generatecv_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/generatecv"
)

var _ = Describe("Constants", func() {
	Describe("State", func() {
		It("should define configuring state", func() {
			Expect(string(generatecv.StateConfiguring)).To(Equal("configuring"))
		})

		It("should define extracting state", func() {
			Expect(string(generatecv.StateExtracting)).To(Equal("extracting"))
		})

		It("should define generating state", func() {
			Expect(string(generatecv.StateGenerating)).To(Equal("generating"))
		})

		It("should define review state", func() {
			Expect(string(generatecv.StateReview)).To(Equal("review"))
		})

		It("should define preview state", func() {
			Expect(string(generatecv.StatePreview)).To(Equal("preview"))
		})

		It("should define exporting state", func() {
			Expect(string(generatecv.StateExporting)).To(Equal("exporting"))
		})

		It("should define export select location state", func() {
			Expect(string(generatecv.StateExportSelectLocation)).To(Equal("export_select_location"))
		})

		It("should define export complete state", func() {
			Expect(string(generatecv.StateExportComplete)).To(Equal("export_complete"))
		})

		It("should have unique values for all states", func() {
			states := []generatecv.State{
				generatecv.StateConfiguring,
				generatecv.StateExtracting,
				generatecv.StateGenerating,
				generatecv.StateReview,
				generatecv.StatePreview,
				generatecv.StateExporting,
				generatecv.StateExportSelectLocation,
				generatecv.StateExportComplete,
			}
			seen := make(map[generatecv.State]bool)
			for _, s := range states {
				Expect(seen[s]).To(BeFalse(), "duplicate state value: %s", s)
				seen[s] = true
			}
		})
	})

	Describe("ExportFormat", func() {
		It("should define text format", func() {
			Expect(string(generatecv.ExportFormatText)).To(Equal("text"))
		})

		It("should define markdown format", func() {
			Expect(string(generatecv.ExportFormatMarkdown)).To(Equal("markdown"))
		})

		It("should define YAML format", func() {
			Expect(string(generatecv.ExportFormatYAML)).To(Equal("yaml"))
		})
	})

	Describe("ExportOption", func() {
		It("should define save to file option", func() {
			Expect(string(generatecv.ExportOptionSaveToFile)).To(Equal("save_to_file"))
		})

		It("should define clipboard option", func() {
			Expect(string(generatecv.ExportOptionClipboard)).To(Equal("clipboard"))
		})

		It("should define cancel option", func() {
			Expect(string(generatecv.ExportOptionCancel)).To(Equal("cancel"))
		})
	})
})
