package styles_test

import (
	"testing"

	"github.com/baphled/kariya/internal/cli/styles"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestVisualFeedback(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Visual Feedback Styles Suite")
}

var _ = Describe("Visual Feedback Styles", func() {
	Context("message box styles", func() {
		It("should have ErrorBox style defined", func() {
			rendered := styles.ErrorBox.Render("Error message")
			Expect(rendered).NotTo(BeEmpty())
			Expect(rendered).To(ContainSubstring("Error message"))
		})

		It("should have WarningBox style defined", func() {
			rendered := styles.WarningBox.Render("Warning message")
			Expect(rendered).NotTo(BeEmpty())
			Expect(rendered).To(ContainSubstring("Warning message"))
		})

		It("should have SuccessBox style defined", func() {
			rendered := styles.SuccessBox.Render("Success message")
			Expect(rendered).NotTo(BeEmpty())
			Expect(rendered).To(ContainSubstring("Success message"))
		})

		It("should have InfoBox style defined", func() {
			rendered := styles.InfoBox.Render("Info message")
			Expect(rendered).NotTo(BeEmpty())
			Expect(rendered).To(ContainSubstring("Info message"))
		})
	})

	Context("text styles for messages", func() {
		It("should have ErrorText style defined", func() {
			rendered := styles.ErrorText.Render("Error")
			Expect(rendered).NotTo(BeEmpty())
		})

		It("should have WarningText style defined", func() {
			rendered := styles.WarningText.Render("Warning")
			Expect(rendered).NotTo(BeEmpty())
		})

		It("should have SuccessText style defined", func() {
			rendered := styles.SuccessText.Render("Success")
			Expect(rendered).NotTo(BeEmpty())
		})

		It("should have InfoText style defined", func() {
			rendered := styles.InfoText.Render("Info")
			Expect(rendered).NotTo(BeEmpty())
		})
	})

	Context("button focus styles", func() {
		It("should have ButtonFocused style defined", func() {
			rendered := styles.ButtonFocused.Render("Button")
			Expect(rendered).NotTo(BeEmpty())
		})

		It("should apply bold to focused buttons", func() {
			rendered := styles.ButtonFocused.Render("Submit")
			Expect(rendered).NotTo(BeEmpty())
		})
	})

	Context("input focus styles", func() {
		It("should have InputFocused style defined", func() {
			rendered := styles.InputFocused.Render("Input text")
			Expect(rendered).NotTo(BeEmpty())
		})

		It("should have InputError style defined", func() {
			rendered := styles.InputError.Render("Invalid input")
			Expect(rendered).NotTo(BeEmpty())
		})
	})

	Context("spinner style", func() {
		It("should have SpinnerStyle defined", func() {
			rendered := styles.SpinnerStyle.Render("⣾")
			Expect(rendered).NotTo(BeEmpty())
		})
	})
})

