package cv

import (
	"context"
	"errors"
	"io"

	"github.com/baphled/kariya/internal/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Note: MockClipboard is defined in export_service_test.go
// This file uses that shared mock for clipboard testing

var _ = Describe("ExportService Clipboard Tests", func() {
	var (
		service       *ExportService
		mockClipboard *MockClipboard
		log           *logger.Logger
		ctx           context.Context
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		mockClipboard = &MockClipboard{}
		service = NewExportServiceWithClipboard(log, mockClipboard)
		ctx = context.Background()
	})

	Describe("CopyToClipboard", func() {
		Context("when clipboard is supported", func() {
			BeforeEach(func() {
				mockClipboard.Unsupported = false
			})

			It("should copy content to clipboard successfully", func() {
				testContent := "Test CV content for clipboard"

				err := service.CopyToClipboard(ctx, testContent)
				Expect(err).NotTo(HaveOccurred())
				Expect(mockClipboard.Content).To(Equal(testContent))
			})

			It("should handle multiple sequential clipboard operations", func() {
				firstContent := "First content"
				secondContent := "Second content"

				// First copy
				err := service.CopyToClipboard(ctx, firstContent)
				Expect(err).NotTo(HaveOccurred())
				Expect(mockClipboard.Content).To(Equal(firstContent))

				// Second copy (should overwrite)
				err = service.CopyToClipboard(ctx, secondContent)
				Expect(err).NotTo(HaveOccurred())
				Expect(mockClipboard.Content).To(Equal(secondContent))
			})

			It("should handle large content", func() {
				// Generate a large string (1MB)
				largeContent := make([]byte, 1024*1024)
				for i := range largeContent {
					largeContent[i] = byte('A' + (i % 26))
				}

				err := service.CopyToClipboard(ctx, string(largeContent))
				Expect(err).NotTo(HaveOccurred())
				Expect(len(mockClipboard.Content)).To(Equal(len(largeContent)))
			})

			It("should return error for empty content", func() {
				err := service.CopyToClipboard(ctx, "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("empty"))
			})

			It("should propagate clipboard write errors", func() {
				mockClipboard.WriteError = errors.New("clipboard write failed")

				err := service.CopyToClipboard(ctx, "test content")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("clipboard"))
			})
		})

		Context("when clipboard is unsupported", func() {
			BeforeEach(func() {
				mockClipboard.Unsupported = true
			})

			It("should return ErrClipboardUnsupported", func() {
				err := service.CopyToClipboard(ctx, "test content")
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(ErrClipboardUnsupported))
			})

			It("should not attempt to write when unsupported", func() {
				mockClipboard.WriteError = errors.New("should not be called")

				err := service.CopyToClipboard(ctx, "test content")
				Expect(err).To(Equal(ErrClipboardUnsupported))
				Expect(mockClipboard.Content).To(BeEmpty())
			})
		})
	})
})
