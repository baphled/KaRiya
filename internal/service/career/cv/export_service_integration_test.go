package cv

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/atotto/clipboard"
	"github.com/baphled/kariya/internal/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExportService Integration Tests", func() {
	var (
		service *ExportService
		log     *logger.Logger
		ctx     context.Context
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		service = NewExportService(log)
		ctx = context.Background()
	})

	Describe("CopyToClipboard Integration", func() {
		Context("when clipboard is supported", func() {
			BeforeEach(func() {
				// Skip if clipboard is unsupported or in CI/headless environment
				if clipboard.Unsupported {
					Skip("Skipping integration test - clipboard not supported in this environment")
				}
				if os.Getenv("CI") != "" {
					Skip("Skipping integration test on CI - no clipboard utilities available")
				}
				if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" && os.Getenv("TERM_PROGRAM") != "Apple_Terminal" {
					Skip("Skipping integration test - no display available (headless environment)")
				}
			})

			It("should copy content to system clipboard and verify with ReadAll", func() {
				testContent := fmt.Sprintf("Integration Test: KaRiya Clipboard Content %d", GinkgoRandomSeed())

				// Copy to clipboard using our service
				err := service.CopyToClipboard(ctx, testContent)
				Expect(err).NotTo(HaveOccurred(), "CopyToClipboard should succeed on supported platforms")

				// Verify by reading directly from clipboard
				clipboardContent, err := clipboard.ReadAll()
				Expect(err).NotTo(HaveOccurred(), "Reading from clipboard should succeed")
				Expect(clipboardContent).To(Equal(testContent), "Clipboard content should match what was written")
			})

			It("should handle multiple sequential clipboard operations", func() {
				seed := GinkgoRandomSeed()
				firstContent := fmt.Sprintf("First content %d", seed)
				secondContent := fmt.Sprintf("Second content %d", seed+1)

				// First copy
				err := service.CopyToClipboard(ctx, firstContent)
				Expect(err).NotTo(HaveOccurred())

				content, err := clipboard.ReadAll()
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(Equal(firstContent))

				// Second copy (should overwrite)
				err = service.CopyToClipboard(ctx, secondContent)
				Expect(err).NotTo(HaveOccurred())

				content, err = clipboard.ReadAll()
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(Equal(secondContent))
			})

			It("should handle large content", func() {
				// Generate a large string (1MB)
				largeContent := make([]byte, 1024*1024)
				for i := range largeContent {
					largeContent[i] = byte('A' + (i % 26))
				}

				err := service.CopyToClipboard(ctx, string(largeContent))
				Expect(err).NotTo(HaveOccurred(), "Should handle large clipboard content")

				content, err := clipboard.ReadAll()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(content)).To(Equal(len(largeContent)), "Large content should be preserved")
			})
		})

		Context("when clipboard is unsupported", func() {
			It("should return ErrClipboardUnsupported on Linux without utilities", func() {
				// This test documents expected behavior
				// On Linux without clipboard utilities, clipboard.Unsupported will be true
				// and CopyToClipboard should return ErrClipboardUnsupported

				if !clipboard.Unsupported {
					Skip("This test only runs on systems without clipboard support")
				}

				err := service.CopyToClipboard(ctx, "test content")
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(ErrClipboardUnsupported))
			})
		})
	})
})
