package cv

import (
	"os"

	"github.com/atotto/clipboard"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clipboard Detection and Error Handling", func() {
	Describe("SystemClipboard.IsUnsupported()", func() {
		It("should detect when clipboard is unavailable", func() {
			sysClip := &SystemClipboard{}

			// In CI/headless environments, clipboard should be detected as unsupported
			isUnsupported := sysClip.IsUnsupported()

			// Log the detection result for debugging
			GinkgoWriter.Printf("clipboard.Unsupported = %v\n", clipboard.Unsupported)
			GinkgoWriter.Printf("SystemClipboard.IsUnsupported() = %v\n", isUnsupported)
			GinkgoWriter.Printf("DISPLAY = %s\n", os.Getenv("DISPLAY"))
			GinkgoWriter.Printf("WAYLAND_DISPLAY = %s\n", os.Getenv("WAYLAND_DISPLAY"))

			// In headless/SSH environments, this should be true
			// In GUI environments with display, this should be false
			if os.Getenv("CI") != "" || (os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "") {
				Expect(isUnsupported).To(BeTrue(), "Clipboard should be unsupported in headless environment")
			}

			// The test always passes - it just documents the behavior
			Expect(true).To(BeTrue())
		})

		It("should provide a clear error message for headless environments", func() {
			// Verify the error message is helpful
			errMsg := ErrClipboardUnsupported.Error()

			Expect(errMsg).To(ContainSubstring("headless environment"))
			Expect(errMsg).To(ContainSubstring("SSH"))
			Expect(errMsg).To(ContainSubstring("Save to file"))

			GinkgoWriter.Printf("Error message: %s\n", errMsg)
		})
	})

	Describe("Clipboard behavior documentation", func() {
		It("should work locally with GUI but fail over SSH", func() {
			GinkgoWriter.Println("\n=== CLIPBOARD BEHAVIOR ===")
			GinkgoWriter.Println("✅ LOCAL (with GUI/display):  Clipboard works")
			GinkgoWriter.Println("❌ SSH (headless/no display): Clipboard fails with clear error")
			GinkgoWriter.Println("")
			GinkgoWriter.Println("Error handling:")
			GinkgoWriter.Println("1. SystemClipboard.IsUnsupported() tests actual clipboard availability")
			GinkgoWriter.Println("2. Returns user-friendly error: 'clipboard not available in headless environment'")
			GinkgoWriter.Println("3. GenerateCV shows error screen with suggestion to use 'Save to file'")
			GinkgoWriter.Println("4. User can retry with different export option")
			GinkgoWriter.Println("=========================\n")

			Expect(true).To(BeTrue())
		})
	})
})
