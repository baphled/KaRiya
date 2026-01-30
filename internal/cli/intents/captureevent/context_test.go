package captureevent_test

import (
	ce "github.com/baphled/kariya/internal/cli/intents/captureevent"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IntentContext", func() {
	Describe("Validate", func() {
		It("should return error when CaptureStrategy is empty", func() {
			ctx := &ce.IntentContext{}
			err := ctx.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("CaptureStrategy"))
		})

		It("should return nil when CaptureStrategy is set", func() {
			ctx := &ce.IntentContext{CaptureStrategy: "quick"}
			err := ctx.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should initialise nil Metadata to empty map", func() {
			ctx := &ce.IntentContext{CaptureStrategy: "quick"}
			Expect(ctx.Metadata).To(BeNil())

			err := ctx.Validate()
			Expect(err).NotTo(HaveOccurred())
			Expect(ctx.Metadata).NotTo(BeNil())
			Expect(ctx.Metadata).To(BeEmpty())
		})

		It("should preserve existing Metadata", func() {
			meta := map[string]string{"key": "value"}
			ctx := &ce.IntentContext{
				CaptureStrategy: "quick",
				Metadata:        meta,
			}

			err := ctx.Validate()
			Expect(err).NotTo(HaveOccurred())
			Expect(ctx.Metadata).To(Equal(meta))
		})

		Context("when services are nil", func() {
			It("should still validate successfully", func() {
				ctx := &ce.IntentContext{
					CaptureStrategy: "manual",
					CLIEventService: nil,
					CareerService:   nil,
				}
				err := ctx.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
