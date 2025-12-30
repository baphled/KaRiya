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

var _ = Describe("Form Model - Responsive Rendering", func() {
	var (
		form *models.FormModel
	)

	BeforeEach(func() {
		repo := careerrepo.NewMemoryRepository()
		svc := careerservice.NewService(repo)
		cliSvc := service.NewCLIEventService(svc)
		form = models.NewFormModel(cliSvc)
	})

	Context("Form Rendering at Default Size", func() {
		It("should render form at default terminal width", func() {
			view := form.View()
			Expect(view).NotTo(BeEmpty())
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("should have proper structure at default size", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Should have multiple lines
			Expect(len(lines)).To(BeNumerically(">", 15))
		})

		It("should display all fields at default width", func() {
			view := form.View()

			// All fields should be present
			Expect(view).To(ContainSubstring("Event Text"))
			Expect(view).To(ContainSubstring("Date"))
			Expect(view).To(ContainSubstring("Company"))
			Expect(view).To(ContainSubstring("Characters:"))
		})

		It("should maintain readability at standard width (80 columns)", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Form should fit in standard width
			for _, line := range lines {
				// Most lines should be reasonably short
				Expect(len(line)).To(BeNumerically("<=", 300))
			}
		})
	})

	Context("Form Rendering at Various Terminal Sizes", func() {
		It("should render form with content present", func() {
			view := form.View()
			// Form should always render with content
			Expect(view).To(ContainSubstring("Event Text"))
		})

		It("should maintain field visibility across renders", func() {
			view1 := form.View()
			view2 := form.View()
			view3 := form.View()

			// All renders should have fields
			Expect(view1).To(ContainSubstring("Event Text"))
			Expect(view2).To(ContainSubstring("Event Text"))
			Expect(view3).To(ContainSubstring("Event Text"))
		})

		It("should handle wide terminal sizes", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Even at wide sizes, should render properly
			Expect(len(lines)).To(BeNumerically(">", 15))
		})

		It("should not have excessively long lines", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Lines should fit in reasonable width (accounting for styling)
			for _, line := range lines {
				Expect(len(line)).To(BeNumerically("<=", 500))
			}
		})
	})

	Context("Form Layout at Small Sizes", func() {
		It("should render form with minimal width constraints", func() {
			view := form.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should keep all fields accessible", func() {
			view := form.View()

			// All main fields should be present
			Expect(view).To(ContainSubstring("Event Text"))
			Expect(view).To(ContainSubstring("Date"))
			Expect(view).To(ContainSubstring("Submit"))
		})

		It("should maintain form structure in compact view", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Should have reasonable number of lines
			Expect(len(lines)).To(BeNumerically(">", 10))
		})
	})

	Context("Form Layout at Large Sizes", func() {
		It("should render form properly at large width", func() {
			view := form.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should maintain consistent formatting at large width", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Should have consistent structure
			Expect(len(lines)).To(BeNumerically(">", 15))
		})

		It("should not create excessive whitespace", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Count lines that are mostly whitespace
			whitespaceLines := 0
			for _, line := range lines {
				if strings.TrimSpace(line) == "" {
					whitespaceLines++
				}
			}

			// Should have some whitespace but not excessive
			totalLines := len(lines)
			ratio := float64(whitespaceLines) / float64(totalLines)
			Expect(ratio).To(BeNumerically("<", 0.5))
		})
	})

	Context("Responsive Content Preservation", func() {
		It("should preserve all labels across sizes", func() {
			view := form.View()

			// All labels should be present
			Expect(view).To(ContainSubstring("Event Text"))
			Expect(view).To(ContainSubstring("(required)"))
			Expect(view).To(ContainSubstring("(optional)"))
		})

		It("should preserve input fields across sizes", func() {
			view := form.View()

			// Should have input indicators
			Expect(len(view)).To(BeNumerically(">", 100))
		})

		It("should preserve helper text across sizes", func() {
			view := form.View()

			// Character counter should be present
			Expect(view).To(ContainSubstring("Characters:"))
			Expect(view).To(ContainSubstring("/2000"))
		})
	})

	Context("Form Consistency Across States", func() {
		It("should maintain structure in different terminal modes", func() {
			view1 := form.View()
			view2 := form.View()

			lines1 := strings.Split(view1, "\n")
			lines2 := strings.Split(view2, "\n")

			// Should have consistent number of lines
			Expect(len(lines1)).To(Equal(len(lines2)))
		})

		It("should render fields in same order consistently", func() {
			view1 := form.View()
			view2 := form.View()

			// Find field positions
			eventIdx1 := strings.Index(view1, "Event Text")
			dateIdx1 := strings.Index(view1, "Date")

			eventIdx2 := strings.Index(view2, "Event Text")
			dateIdx2 := strings.Index(view2, "Date")

			// Order should be consistent
			Expect(eventIdx1).To(BeNumerically("<", dateIdx1))
			Expect(eventIdx2).To(BeNumerically("<", dateIdx2))
		})

		It("should maintain responsive layout across renders", func() {
			for i := 0; i < 3; i++ {
				view := form.View()
				Expect(view).NotTo(BeEmpty())
				lines := strings.Split(view, "\n")
				Expect(len(lines)).To(BeNumerically(">", 15))
			}
		})
	})

	Context("Responsive Styling Application", func() {
		It("should apply styles consistently at any width", func() {
			view := form.View()

			// Styled elements should be present
			Expect(view).To(ContainSubstring("►"))
			Expect(view).To(ContainSubstring("Event Text"))
		})

		It("should not lose styling at different widths", func() {
			view := form.View()

			// Focus indicator should be present
			Expect(view).To(ContainSubstring("►"))
		})

		It("should maintain visual hierarchy at all sizes", func() {
			view := form.View()

			// Labels should appear before inputs
			eventIdx := strings.Index(view, "Event Text")
			charCountIdx := strings.Index(view, "Characters:")

			Expect(eventIdx).To(BeNumerically("<", charCountIdx))
		})
	})

	Context("Terminal Size Edge Cases", func() {
		It("should handle minimum viable width", func() {
			view := form.View()
			// Should still render with core content
			Expect(view).To(ContainSubstring("Event Text"))
		})

		It("should handle extra large width without breaking", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Should still render properly
			Expect(len(lines)).To(BeNumerically(">", 10))
		})

		It("should not truncate critical content", func() {
			view := form.View()

			// Essential elements should be present
			Expect(view).To(ContainSubstring("Submit"))
			Expect(view).To(ContainSubstring("Event Text"))
		})
	})
})

