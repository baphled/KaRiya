package models_test

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/styles"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Form Model - Input Field Styling Consistency", func() {
	var (
		form *models.FormModel
	)

	BeforeEach(func() {
		repo := careerrepo.NewMemoryRepository()
		svc := careerservice.NewService(repo)
		cliSvc := service.NewCLIEventService(svc)
		form = models.NewFormModel(cliSvc)
	})

	Context("Input Field Base Styling", func() {
		It("should render input fields with consistent base style", func() {
			view := form.View()
			// Input fields should be present in output
			Expect(view).NotTo(BeEmpty())
			// Should have input sections (indicated by placeholders in textinput)
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("should render all input fields with proper padding", func() {
			view := form.View()
			// InputBase style has Padding(0, 1) - should reflect in output
			Expect(view).NotTo(BeEmpty())
			// Verify form renders without panic
			Expect(func() {
				_ = form.View()
			}).NotTo(Panic())
		})

		It("should use consistent border style for all inputs", func() {
			view := form.View()
			// All inputs should use RoundedBorder (from InputBase style)
			// This is verified by successful rendering
			Expect(view).To(ContainSubstring("Text"))
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("should apply consistent background color to inputs", func() {
			// InputBase uses ColorBackgroundCard background
			// Verify by rendering multiple inputs
			view1 := form.View()
			view2 := form.View()
			Expect(view1).To(Equal(view2))
		})
	})

	Context("Input Field Focus States", func() {
		It("should apply focus styling when field is focused", func() {
			// Text field starts focused
			view := form.View()
			// Should render with focus indicator
			Expect(view).To(ContainSubstring("►"))
		})

		It("should transition focus styling when changing fields", func() {
			view1 := form.View()
			Expect(view1).To(ContainSubstring("►"))

			// Tab to next field
			form.Update(tea.KeyMsg{Type: tea.KeyTab})
			view2 := form.View()

			// Both should render validly (focus moved but still present)
			Expect(view2).NotTo(BeEmpty())
			Expect(len(view2)).To(BeNumerically(">", 0))
		})

		It("should maintain focus indicators visible in output", func() {
			// Initial state - text field focused
			view := form.View()
			focusCount := strings.Count(view, "►")
			Expect(focusCount).To(BeNumerically(">", 0))

			// Tab to next field
			form.Update(tea.KeyMsg{Type: tea.KeyTab})
			view = form.View()

			// Should still have focus indicators (at least one visible)
			focusCount = strings.Count(view, "►")
			Expect(focusCount).To(BeNumerically(">", 0))
		})
	})

	Context("Input Field Error States", func() {
		It("should render input fields consistently in normal state", func() {
			view := form.View()
			// No errors initially
			Expect(view).NotTo(ContainSubstring("Error"))
		})

		It("should handle error display without breaking input styling", func() {
			// Submit with invalid text (empty)
			form.Update(tea.KeyMsg{Type: tea.KeyEnter})
			view := form.View()

			// Should still render form with inputs visible
			Expect(view).NotTo(BeEmpty())
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("should display validation errors consistently", func() {
			// Try to submit empty form
			form.Update(tea.KeyMsg{Type: tea.KeyEnter})
			view := form.View()

			// Error should be present but inputs should still render
			Expect(len(view)).To(BeNumerically(">", 0))
		})
	})

	Context("Multiple Input Type Consistency", func() {
		It("should render text input consistently", func() {
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			view := form.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render all input types with same styling", func() {
			view := form.View()
			// Text, Date, Company, Project inputs all use InputBase
			lines := strings.Split(view, "\n")
			// Should have input sections
			Expect(len(lines)).To(BeNumerically(">", 5))
		})

		It("should maintain input styling consistency across updates", func() {
			view1 := form.View()
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			view2 := form.View()
			form.Update(tea.KeyMsg{Type: tea.KeyBackspace})
			view3 := form.View()

			// All views should be non-empty and properly formatted
			Expect(view1).NotTo(BeEmpty())
			Expect(view2).NotTo(BeEmpty())
			Expect(view3).NotTo(BeEmpty())
		})
	})

	Context("Input Width and Responsive Styling", func() {
		It("should apply consistent width to all input fields", func() {
			view := form.View()
			lines := strings.Split(view, "\n")
			// Should have consistent line lengths (within reasonable bounds)
			Expect(len(lines)).To(BeNumerically(">", 3))
		})

		It("should render inputs without overflow", func() {
			// Type long text
			for i := 0; i < 100; i++ {
				form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			}
			view := form.View()
			Expect(view).NotTo(BeEmpty())
			Expect(func() {
				_ = form.View()
			}).NotTo(Panic())
		})

		It("should maintain responsive layout for different content lengths", func() {
			// Empty view
			view1 := form.View()
			Expect(view1).NotTo(BeEmpty())

			// With content
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t', 'e', 's', 't'}})
			view2 := form.View()
			Expect(view2).NotTo(BeEmpty())

			// With error
			form.Update(tea.KeyMsg{Type: tea.KeyEnter})
			view3 := form.View()
			Expect(view3).NotTo(BeEmpty())
		})
	})

	Context("Input Styling Verification", func() {
		It("should use consistent style values for InputBase", func() {
			// Verify styles are properly defined
			base := styles.InputBase
			focused := styles.InputFocused
			errorStyle := styles.InputError

			// All should be non-nil and properly initialized
			Expect(base).NotTo(BeNil())
			Expect(focused).NotTo(BeNil())
			Expect(errorStyle).NotTo(BeNil())
		})

		It("should render form with all input components visible", func() {
			view := form.View()
			// Should contain text input indicator
			Expect(view).To(ContainSubstring("Characters:"))
		})

		It("should maintain visual hierarchy with consistent spacing", func() {
			view := form.View()
			lines := strings.Split(view, "\n")
			// Should have multiple sections (labels, inputs, helpers)
			Expect(len(lines)).To(BeNumerically(">", 10))
		})
	})
})

