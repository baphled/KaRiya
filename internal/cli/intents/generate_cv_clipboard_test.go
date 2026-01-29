package intents

import (
	"context"
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GenerateCV Clipboard Export Error Handling", func() {
	var (
		intent *GenerateCVIntent
		ctx    context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()

		// Create minimal profiles for testing (required by NewGenerateCVIntent)
		profiles := []*CVProfile{
			{
				ID:          "test-profile",
				Name:        "Test Profile",
				Description: "Test profile for clipboard export tests",
			},
		}

		// Create minimal events for testing (required by NewGenerateCVIntent)
		events := []*career.Event{
			{
				ID:      "event-1",
				Text:    "Test event for clipboard export tests",
				Date:    time.Now(),
				Company: "Test Company",
			},
		}

		// Create minimal context for testing
		genCtx := &GenerateCVContext{
			AppContext:        ctx,
			AvailableProfiles: profiles,
			DefaultProfile:    profiles[0],
			Events:            events,
			Facts:             make([]*career.Fact, 0),
		}
		var err error
		intent, err = NewGenerateCVIntent(genCtx)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("CVExportCompleteMsg error handling", func() {
		It("should transition to ExportComplete state when clipboard export fails", func() {
			// Setup: Set the intent to exporting state
			intent.state.currentState = GenerateCVStateExporting
			intent.state.isExporting = true
			intent.state.selectedExportOption = CVExportOptionClipboard

			// Simulate clipboard export failure
			errorMsg := CVExportCompleteMsg{
				Path:  "",
				Error: fmt.Errorf("failed to copy to clipboard: exit status 1"),
			}

			// Process the error message
			intent.updateExporting(errorMsg)

			// Verify: Should transition to ExportComplete state (not ExportSelectLocation)
			Expect(intent.state.currentState).To(Equal(GenerateCVStateExportComplete))
			Expect(intent.state.exportError).NotTo(BeNil())
			Expect(intent.state.exportError.Error()).To(ContainSubstring("clipboard"))
			Expect(intent.state.isExporting).To(BeFalse())
		})

		It("should show error in viewExportComplete when clipboard fails", func() {
			// Setup: Set error state
			intent.state.currentState = GenerateCVStateExportComplete
			intent.state.exportError = fmt.Errorf("failed to copy to clipboard: exit status 1")
			intent.state.selectedExportOption = CVExportOptionClipboard

			// Get the view
			view := intent.viewExportComplete()

			// Verify: Error should be displayed (UIKit uses styled text without emoji)
			Expect(view).To(ContainSubstring("Export Failed"))
			Expect(view).To(ContainSubstring("Error:"))
			Expect(view).To(ContainSubstring("clipboard"))
			Expect(view).To(ContainSubstring("Try a different location or format"))
		})

		It("should NOT show error when export succeeds", func() {
			// Setup: Successful export
			intent.state.currentState = GenerateCVStateExportComplete
			intent.state.exportError = nil
			intent.state.exportedPath = "clipboard"
			intent.state.selectedExportOption = CVExportOptionClipboard

			// Get the view
			view := intent.viewExportComplete()

			// Verify: Success message, no error (UIKit uses styled text without emoji)
			Expect(view).To(ContainSubstring("Export Complete"))
			Expect(view).NotTo(ContainSubstring("Export Failed"))
			Expect(view).To(ContainSubstring("Location: Clipboard"))
			Expect(view).To(ContainSubstring("paste the CV anywhere"))
		})
	})
})
