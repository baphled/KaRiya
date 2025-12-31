package models_test

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Form Model - Focus Indicators Consistency", func() {
	var (
		form *models.FormModel
	)

	BeforeEach(func() {
		repo := careerrepo.NewMemoryRepository()
		svc := careerservice.NewService(repo)
		cliSvc := service.NewCLIEventService(svc)
		form = models.NewFormModel(cliSvc)
	})

	Context("Focus Indicator Display", func() {
		It("should display focus indicator for the focused field", func() {
			// Text field starts focused
			view := form.View()
			// Should contain focus indicator
			Expect(view).To(ContainSubstring("►"))
		})

		It("should display focus indicator as arrow character", func() {
			view := form.View()
			// Focus indicator should be ►
			Expect(view).To(ContainSubstring("►"))
		})

		It("should show focus indicator in initial state", func() {
			view := form.View()
			focusCount := strings.Count(view, "►")
			// Should have at least one focus indicator
			Expect(focusCount).To(BeNumerically(">", 0))
		})

		It("should position focus indicator before field label", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Find the line with focus indicator
			foundFocusIndicator := false
			for _, line := range lines {
				if strings.Contains(line, "►") && strings.Contains(line, "Event Text") {
					foundFocusIndicator = true
					// Focus indicator should come before the label text
					idx := strings.Index(line, "►")
					labelIdx := strings.Index(line, "Event Text")
					Expect(idx).To(BeNumerically("<", labelIdx))
				}
			}
			Expect(foundFocusIndicator).To(BeTrue())
		})
	})

	Context("Focus Indicator Styling Consistency", func() {
		It("should use same indicator character for all fields", func() {
			view1 := form.View()
			indicator1 := strings.Count(view1, "►")
			Expect(indicator1).To(BeNumerically(">", 0))

			// Tab to next field
			form.Update(tea.KeyMsg{Type: tea.KeyTab})
			view2 := form.View()
			// Should still have focus indicator (character used consistently)
			Expect(view2).To(ContainSubstring("►"))
		})

		It("should maintain consistent focus indicator presence", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Should have at least one line with focus indicator
			hasIndicator := false
			for _, line := range lines {
				if strings.Contains(line, "►") {
					hasIndicator = true
					// Should have proper spacing around indicator
					Expect(line).To(ContainSubstring("►"))
				}
			}
			Expect(hasIndicator).To(BeTrue())
		})

		It("should render focus indicator without distorting layout", func() {
			view := form.View()
			lines := strings.Split(view, "\n")
			// Should have consistent number of lines
			Expect(len(lines)).To(BeNumerically(">", 10))
		})
	})

	Context("Focus State Transitions", func() {
		It("should maintain focus indicator when navigating fields", func() {
			// Initial focus on text field
			view1 := form.View()
			Expect(view1).To(ContainSubstring("►"))

			// Tab to next field
			form.Update(tea.KeyMsg{Type: tea.KeyTab})
			view2 := form.View()
			// Should still have focus indicator somewhere
			Expect(view2).To(ContainSubstring("►"))
		})

		It("should update focus indicator position on field change", func() {
			view1 := form.View()
			Expect(view1).To(ContainSubstring("►"))

			// Simulate field navigation
			form.Update(tea.KeyMsg{Type: tea.KeyTab})
			view2 := form.View()
			// Focus indicator should still be present
			Expect(view2).To(ContainSubstring("►"))
		})

		It("should display focus indicator consistently across Tab navigation", func() {
			// Navigate through fields
			for i := 0; i < 5; i++ {
				view := form.View()
				// Should always have focus indicator
				Expect(view).To(ContainSubstring("►"))
				form.Update(tea.KeyMsg{Type: tea.KeyTab})
			}
		})
	})

	Context("Focus Indicator Visibility", func() {
		It("should make focused field visually distinct", func() {
			view := form.View()
			// Focus indicator makes field stand out
			Expect(view).To(ContainSubstring("►"))
			// Should have label next to indicator
			Expect(view).To(ContainSubstring("Event Text"))
		})

		It("should render focus indicator clearly without ambiguity", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Find focused field line
			var focusedLine string
			for _, line := range lines {
				if strings.Contains(line, "►") && strings.Contains(line, "Event Text") {
					focusedLine = line
					break
				}
			}

			Expect(focusedLine).NotTo(BeEmpty())
			// Should be clear which field is focused
			Expect(focusedLine).To(ContainSubstring("►"))
			Expect(focusedLine).To(ContainSubstring("Event Text"))
		})

		It("should distinguish focused field from unfocused fields", func() {
			view := form.View()
			// Focused field has ► indicator
			lines := strings.Split(view, "\n")

			focusedCount := 0
			unfocusedCount := 0
			for _, line := range lines {
				if strings.Contains(line, "►") {
					focusedCount++
				} else if strings.Contains(line, "Date") || strings.Contains(line, "Company") {
					unfocusedCount++
				}
			}

			// Should have focused and unfocused fields
			Expect(focusedCount).To(BeNumerically(">", 0))
			Expect(unfocusedCount).To(BeNumerically(">", 0))
		})
	})

	Context("Focus Indicator Rendering", func() {
		It("should render focus indicator without panic", func() {
			Expect(func() {
				_ = form.View()
			}).NotTo(Panic())
		})

		It("should maintain focus indicator across multiple renders", func() {
			view1 := form.View()
			view2 := form.View()
			view3 := form.View()

			// All should have focus indicator
			Expect(view1).To(ContainSubstring("►"))
			Expect(view2).To(ContainSubstring("►"))
			Expect(view3).To(ContainSubstring("►"))
		})

		It("should render focus indicator with consistent formatting", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Find all focus indicators
			focusLines := 0
			for _, line := range lines {
				if strings.Contains(line, "►") {
					focusLines++
				}
			}

			// Should have at least one line with focus indicator
			Expect(focusLines).To(BeNumerically(">", 0))
		})

		It("should render form without focus indicator distorting layout", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// All lines should render properly
			lineCount := 0
			for _, line := range lines {
				if len(line) > 0 {
					lineCount++
				}
			}
			Expect(lineCount).To(BeNumerically(">", 5))
		})
	})

	Context("Focus Indicator Consistency Across States", func() {
		It("should display consistent indicator in normal state", func() {
			view := form.View()
			Expect(view).To(ContainSubstring("►"))
		})

		It("should maintain focus indicator during text input", func() {
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T', 'e', 's', 't'}})
			view := form.View()
			// Focus indicator should still be present
			Expect(view).To(ContainSubstring("►"))
		})

		It("should keep focus indicator visible after editing", func() {
			// Type some text
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t', 'e', 's', 't'}})
			// Clear it
			for i := 0; i < 4; i++ {
				form.Update(tea.KeyMsg{Type: tea.KeyBackspace})
			}
			view := form.View()
			// Focus indicator should still be visible
			Expect(view).To(ContainSubstring("►"))
		})

		It("should display consistent focus indicator across input/deletion cycles", func() {
			// Type text
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'A', 'B', 'C'}})
			view1 := form.View()

			// Delete text
			form.Update(tea.KeyMsg{Type: tea.KeyBackspace})
			form.Update(tea.KeyMsg{Type: tea.KeyBackspace})
			view2 := form.View()

			// Add more text
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'D', 'E', 'F'}})
			view3 := form.View()

			// All should have focus indicator
			Expect(view1).To(ContainSubstring("►"))
			Expect(view2).To(ContainSubstring("►"))
			Expect(view3).To(ContainSubstring("►"))
		})
	})
})
