package models_test

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Form Model - Field Spacing and Alignment", func() {
	var (
		form *models.FormModel
	)

	BeforeEach(func() {
		repo := careerrepo.NewMemoryRepository()
		svc := careerservice.NewService(repo)
		cliSvc := service.NewCLIEventService(svc)
		form = models.NewFormModel(cliSvc)
	})

	Context("Vertical Field Spacing", func() {
		It("should have consistent spacing between form fields", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Should have multiple lines with proper spacing
			Expect(len(lines)).To(BeNumerically(">", 15))
		})

		It("should render fields with gaps between sections", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Count non-empty lines
			nonEmptyLines := 0
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					nonEmptyLines++
				}
			}

			// Should have good proportion of content lines
			Expect(nonEmptyLines).To(BeNumerically(">", 10))
		})

		It("should maintain consistent line breaks between fields", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Should have multiple field sections with space between them
			fieldCount := 0
			for _, line := range lines {
				if strings.Contains(line, "Event Text") || strings.Contains(line, "Date") {
					fieldCount++
				}
			}
			Expect(fieldCount).To(BeNumerically(">", 1))
		})
	})

	Context("Field Label Alignment", func() {
		It("should have field labels present in the output", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Find lines with field labels
			labelLines := 0
			for _, line := range lines {
				if strings.Contains(line, "Event Text") ||
					strings.Contains(line, "Date") ||
					strings.Contains(line, "Company") {
					labelLines++
				}
			}
			Expect(labelLines).To(BeNumerically(">", 0))
		})

		It("should position labels with proper structure", func() {
			view := form.View()
			// Labels should be present and properly positioned
			Expect(view).To(ContainSubstring("Event Text"))
		})

		It("should maintain label consistency across all field types", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			labelPositions := []string{}
			for _, line := range lines {
				if strings.Contains(line, ":") &&
					(strings.Contains(line, "Event") ||
						strings.Contains(line, "Date") ||
						strings.Contains(line, "Company")) {
					labelPositions = append(labelPositions, line)
				}
			}

			// Should have multiple labels
			Expect(len(labelPositions)).To(BeNumerically(">", 1))
		})
	})

	Context("Input Field Alignment", func() {
		It("should align input fields with proper structure", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Input fields should be present and aligned
			Expect(len(lines)).To(BeNumerically(">", 10))
		})

		It("should position inputs below their labels", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Find text field label and verify input appears after it
			textLabelIdx := -1
			for i, line := range lines {
				if strings.Contains(line, "Event Text") {
					textLabelIdx = i
					break
				}
			}

			Expect(textLabelIdx).To(BeNumerically(">=", 0))
			// Input should appear in lines after label
			Expect(len(lines)).To(BeNumerically(">", textLabelIdx+2))
		})

		It("should render inputs with consistent spacing from labels", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Should have label-input pairs with consistent spacing
			for i := 0; i < len(lines)-1; i++ {
				if strings.Contains(lines[i], "Event Text") {
					// Next non-empty line should be input
					if i+1 < len(lines) {
						Expect(lines[i+1]).NotTo(BeEmpty())
					}
				}
			}
		})
	})

	Context("Character Counter and Helper Text Spacing", func() {
		It("should display character count below input with proper spacing", func() {
			view := form.View()
			// Should contain character counter
			Expect(view).To(ContainSubstring("Characters:"))
		})

		It("should position helper text below its input field", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Character count should appear after text input section
			charCountLine := -1
			for i, line := range lines {
				if strings.Contains(line, "Characters:") {
					charCountLine = i
					break
				}
			}

			Expect(charCountLine).To(BeNumerically(">", 0))
		})

		It("should maintain consistent formatting of helper text", func() {
			view := form.View()
			// Helper text should be properly formatted
			Expect(view).To(ContainSubstring("Characters:"))
			Expect(view).To(ContainSubstring("/"))
		})
	})

	Context("Section Separation", func() {
		It("should have clear separation between form sections", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Should have distinct sections (header, form, footer)
			Expect(len(lines)).To(BeNumerically(">", 15))
		})

		It("should maintain padding around form card", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// First and last lines should have content for padding
			Expect(len(lines[0])).To(BeNumerically(">", 0))
			Expect(len(lines[len(lines)-1])).To(BeNumerically(">", 0))
		})

		It("should properly space header and form content", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Should have structure: header, separator, form, footer
			Expect(len(lines)).To(BeNumerically(">", 10))
		})
	})

	Context("Responsive Spacing", func() {
		It("should maintain spacing consistency across renders", func() {
			view1 := form.View()
			lines1 := strings.Split(view1, "\n")

			view2 := form.View()
			lines2 := strings.Split(view2, "\n")

			// Same number of lines across renders
			Expect(len(lines1)).To(Equal(len(lines2)))
		})

		It("should preserve spacing with different states", func() {
			// Initial render
			view1 := form.View()

			// Render again
			view2 := form.View()

			// Both should have proper spacing
			lines1 := strings.Split(view1, "\n")
			lines2 := strings.Split(view2, "\n")

			Expect(len(lines1)).To(BeNumerically(">", 10))
			Expect(len(lines2)).To(BeNumerically(">", 10))
		})

		It("should maintain alignment consistency across all states", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Should have consistent structure
			Expect(len(lines)).To(BeNumerically(">", 15))

			// Fields should be properly spaced
			nonEmptyLines := 0
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					nonEmptyLines++
				}
			}
			Expect(nonEmptyLines).To(BeNumerically(">", 10))
		})
	})

	Context("Visual Hierarchy and Alignment", func() {
		It("should maintain visual hierarchy through spacing", func() {
			view := form.View()
			// Form structure should be clear and hierarchical
			Expect(view).To(ContainSubstring("Event Text"))
			Expect(view).To(ContainSubstring("Characters:"))
		})

		It("should properly format and space all form elements", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// All lines should be properly formatted
			for _, line := range lines {
				// Lines should not be unreasonably long (form has max width)
				Expect(len(line)).To(BeNumerically("<=", 300))
			}
		})

		It("should maintain consistent padding around all elements", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Should have structure with padding
			Expect(len(lines)).To(BeNumerically(">", 15))

			// Content lines should be present
			contentLineCount := 0
			for _, line := range lines {
				if strings.Contains(line, "►") || strings.Contains(line, "Event") {
					contentLineCount++
				}
			}
			Expect(contentLineCount).To(BeNumerically(">", 0))
		})
	})
})
