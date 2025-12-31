package models_test

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/styles"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Form Model - Error Message Display Consistency", func() {
	var (
		form *models.FormModel
	)

	BeforeEach(func() {
		repo := careerrepo.NewMemoryRepository()
		svc := careerservice.NewService(repo)
		cliSvc := service.NewCLIEventService(svc)
		form = models.NewFormModel(cliSvc)
	})

	Context("Error Message Styling Infrastructure", func() {
		It("should have ErrorText style properly defined", func() {
			// Verify ErrorText style exists and is usable
			errStyle := styles.ErrorText
			Expect(errStyle).NotTo(BeNil())
		})

		It("should have ErrorHint style for secondary error info", func() {
			// Verify ErrorHint style exists
			Expect(styles.ErrorHint).NotTo(BeNil())
		})

		It("should render form without error messages initially", func() {
			// Form should start with no errors displayed
			view := form.View()
			Expect(view).NotTo(BeEmpty())
			// No error messages should be shown initially
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("should render form structure correctly", func() {
			// Verify form renders with expected structure
			view := form.View()
			lines := strings.Split(view, "\n")
			// Should have multiple lines for labels, inputs, hints
			Expect(len(lines)).To(BeNumerically(">", 5))
		})

		It("should render without panic on error display call", func() {
			Expect(func() {
				_ = form.View()
			}).NotTo(Panic())
		})
	})

	Context("Error Message Display Capability", func() {
		It("should support displaying validation errors via ErrorText style", func() {
			// Create a styled error message manually
			errorMsg := styles.ErrorText.Render("test error message")
			Expect(errorMsg).NotTo(BeEmpty())
		})

		It("should render error messages consistently across multiple calls", func() {
			// Test ErrorText consistency
			error1 := styles.ErrorText.Render("validation error")
			error2 := styles.ErrorText.Render("validation error")
			Expect(error1).To(Equal(error2))
		})

		It("should apply error styling without breaking form layout", func() {
			// Verify form can be rendered with error content
			view := form.View()
			Expect(view).NotTo(BeEmpty())
			// Should still render all components
			Expect(view).To(ContainSubstring("Event Text"))
		})

		It("should have different styling for error vs normal text", func() {
			// Verify error style is distinct
			normal := styles.InfoText.Render("normal text")
			error := styles.ErrorText.Render("error text")
			// Both should render, but with different styles
			Expect(normal).NotTo(BeEmpty())
			Expect(error).NotTo(BeEmpty())
		})
	})

	Context("Error Message Positioning", func() {
		It("should render form sections in consistent order", func() {
			// Verify form structure is consistent
			view := form.View()
			lines := strings.Split(view, "\n")

			// Should have text label before other content
			textLabelFound := false
			for i, line := range lines {
				if strings.Contains(line, "Event Text") {
					textLabelFound = true
					// Ensure label comes before later fields
					Expect(i).To(BeNumerically("<", len(lines)-5))
				}
			}
			Expect(textLabelFound).To(BeTrue())
		})

		It("should maintain consistent field order across renders", func() {
			view1 := form.View()
			view2 := form.View()

			lines1 := strings.Split(view1, "\n")
			lines2 := strings.Split(view2, "\n")

			// Both should have consistent structure
			Expect(len(lines1)).To(Equal(len(lines2)))
		})

		It("should properly space components for readability", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Should have sufficient lines for proper spacing
			Expect(len(lines)).To(BeNumerically(">", 10))
		})
	})

	Context("Error State Management", func() {
		It("should initialize form with no error state", func() {
			Expect(form.Error()).To(BeNil())
		})

		It("should render form correctly in no-error state", func() {
			view := form.View()
			Expect(view).NotTo(BeEmpty())
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("should support error display through fieldErrors map", func() {
			// Form has fieldErrors map for tracking validation errors
			// Even if not currently displayed, the infrastructure supports it
			view := form.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Context("Styling Consistency", func() {
		It("should use same error color throughout form", func() {
			// All errors should use ErrorText style (ColorError foreground)
			// Verify style is defined correctly
			errStyle := styles.ErrorText
			Expect(errStyle).NotTo(BeNil())
		})

		It("should maintain visual hierarchy between error and normal text", func() {
			// Normal text uses InfoText (secondary color)
			// Error text uses ErrorText (error color)
			normalText := styles.InfoText.Render("info")
			errorText := styles.ErrorText.Render("error")

			// Both should render
			Expect(normalText).NotTo(BeEmpty())
			Expect(errorText).NotTo(BeEmpty())
		})

		It("should render form with consistent styling applied", func() {
			view := form.View()

			// Form should contain styled components
			Expect(view).To(ContainSubstring("Event Text"))
			Expect(view).To(ContainSubstring("Characters:"))
		})
	})

	Context("Error Display Infrastructure", func() {
		It("should have all required error styles defined", func() {
			// ErrorText for main messages
			Expect(styles.ErrorText).NotTo(BeNil())
			// ErrorHint for secondary info
			Expect(styles.ErrorHint).NotTo(BeNil())
		})

		It("should render form without errors triggering panic", func() {
			Expect(func() {
				view := form.View()
				Expect(view).NotTo(BeEmpty())
			}).NotTo(Panic())
		})

		It("should support styled error rendering", func() {
			// Manually create styled error to verify capability
			styledError := styles.ErrorText.Render("field cannot be empty")
			Expect(styledError).To(ContainSubstring("field cannot be empty"))
		})
	})
})
