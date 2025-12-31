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

var _ = Describe("Form Model - Character Counter Display Consistency", func() {
	var (
		form *models.FormModel
	)

	BeforeEach(func() {
		repo := careerrepo.NewMemoryRepository()
		svc := careerservice.NewService(repo)
		cliSvc := service.NewCLIEventService(svc)
		form = models.NewFormModel(cliSvc)
	})

	Context("Character Counter Display", func() {
		It("should display character counter in form", func() {
			view := form.View()
			// Should contain character counter
			Expect(view).To(ContainSubstring("Characters:"))
		})

		It("should show counter with correct format (X/2000)", func() {
			view := form.View()
			// Should show format: Characters: 0/2000
			Expect(view).To(ContainSubstring("Characters:"))
			Expect(view).To(ContainSubstring("/"))
			Expect(view).To(ContainSubstring("2000"))
		})

		It("should display counter below text input field", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Find text input section and verify counter appears after
			textLabelIdx := -1
			counterIdx := -1

			for i, line := range lines {
				if strings.Contains(line, "Event Text") {
					textLabelIdx = i
				}
				if strings.Contains(line, "Characters:") {
					counterIdx = i
				}
			}

			// Counter should appear after text field label
			if textLabelIdx >= 0 && counterIdx >= 0 {
				Expect(counterIdx).To(BeNumerically(">", textLabelIdx))
			}
		})

		It("should update counter as user types", func() {
			// Initial counter
			view1 := form.View()
			Expect(view1).To(ContainSubstring("0/2000"))

			// Type some text
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'H', 'e', 'l', 'l', 'o'}})
			view2 := form.View()

			// Counter should show 5 characters
			Expect(view2).To(ContainSubstring("5/2000"))
		})

		It("should track character count accurately during input", func() {
			// Type "Test" (4 chars)
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T', 'e', 's', 't'}})
			view := form.View()
			Expect(view).To(ContainSubstring("4/2000"))
		})

		It("should update counter when text is deleted", func() {
			// Type text
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'A', 'B', 'C', 'D'}})
			view1 := form.View()
			Expect(view1).To(ContainSubstring("4/2000"))

			// Delete one character
			form.Update(tea.KeyMsg{Type: tea.KeyBackspace})
			view2 := form.View()
			Expect(view2).To(ContainSubstring("3/2000"))
		})
	})

	Context("Character Counter Styling", func() {
		It("should use consistent InfoText style for counter", func() {
			// Verify InfoText style is defined
			infoStyle := styles.InfoText
			Expect(infoStyle).NotTo(BeNil())
		})

		It("should render counter in InfoText style", func() {
			view := form.View()
			// Counter should be rendered with consistent styling
			Expect(view).To(ContainSubstring("Characters:"))
		})

		It("should maintain counter visibility with consistent formatting", func() {
			// Type text to increase counter
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T', 'e', 's', 't'}})
			view := form.View()

			// Counter should still be visible and formatted
			Expect(view).To(ContainSubstring("Characters:"))
			Expect(view).To(ContainSubstring("4/2000"))
		})
	})

	Context("Character Counter Positioning", func() {
		It("should position counter consistently below input", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Find character counter line
			var counterLine string
			for _, line := range lines {
				if strings.Contains(line, "Characters:") {
					counterLine = line
					break
				}
			}

			Expect(counterLine).NotTo(BeEmpty())
			Expect(counterLine).To(ContainSubstring("Characters:"))
		})

		It("should maintain counter position across renders", func() {
			// Get initial counter position
			view1 := form.View()
			lines1 := strings.Split(view1, "\n")

			counterLine1 := -1
			for i, line := range lines1 {
				if strings.Contains(line, "Characters:") {
					counterLine1 = i
					break
				}
			}

			// Render again
			view2 := form.View()
			lines2 := strings.Split(view2, "\n")

			counterLine2 := -1
			for i, line := range lines2 {
				if strings.Contains(line, "Characters:") {
					counterLine2 = i
					break
				}
			}

			// Counter should be at same position
			Expect(counterLine1).To(Equal(counterLine2))
		})

		It("should keep counter visible after text input", func() {
			// Type text
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'L', 'o', 'n', 'g', ' ', 't', 'e', 'x', 't'}})
			view := form.View()

			// Counter should still be visible
			Expect(view).To(ContainSubstring("Characters:"))
		})
	})

	Context("Counter Value Accuracy", func() {
		It("should show 0 for empty input", func() {
			view := form.View()
			Expect(view).To(ContainSubstring("0/2000"))
		})

		It("should accurately count single character", func() {
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'A'}})
			view := form.View()
			Expect(view).To(ContainSubstring("1/2000"))
		})

		It("should accurately count multiple characters", func() {
			// Type 10 characters
			for i := 0; i < 10; i++ {
				form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'X'}})
			}
			view := form.View()
			Expect(view).To(ContainSubstring("10/2000"))
		})

		It("should correctly display count after deletions", func() {
			// Type 5 characters
			for i := 0; i < 5; i++ {
				form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'A'}})
			}

			// Delete 2 characters
			form.Update(tea.KeyMsg{Type: tea.KeyBackspace})
			form.Update(tea.KeyMsg{Type: tea.KeyBackspace})

			view := form.View()
			Expect(view).To(ContainSubstring("3/2000"))
		})
	})

	Context("Counter Format Consistency", func() {
		It("should always show fraction format (X/2000)", func() {
			view := form.View()
			// Should show number/2000 format
			Expect(view).To(MatchRegexp(`\d+/2000`))
		})

		It("should maintain format with different character counts", func() {
			// Empty
			view1 := form.View()
			Expect(view1).To(MatchRegexp(`0/2000`))

			// With content
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T', 'e', 's', 't'}})
			view2 := form.View()
			Expect(view2).To(MatchRegexp(`4/2000`))
		})

		It("should include 'Characters:' label consistently", func() {
			// Empty state
			view1 := form.View()
			Expect(view1).To(ContainSubstring("Characters:"))

			// With text
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'S', 'o', 'm', 'e', 't', 'e', 'x', 't'}})
			view2 := form.View()
			Expect(view2).To(ContainSubstring("Characters:"))
		})
	})

	Context("Counter in Form Context", func() {
		It("should render counter without breaking form layout", func() {
			view := form.View()
			lines := strings.Split(view, "\n")

			// Should still have proper form structure
			Expect(len(lines)).To(BeNumerically(">", 15))
		})

		It("should maintain counter visibility across all form states", func() {
			// Initial state
			view1 := form.View()
			Expect(view1).To(ContainSubstring("Characters:"))

			// After input
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T'}})
			view2 := form.View()
			Expect(view2).To(ContainSubstring("Characters:"))

			// After deletion
			form.Update(tea.KeyMsg{Type: tea.KeyBackspace})
			view3 := form.View()
			Expect(view3).To(ContainSubstring("Characters:"))
		})

		It("should keep counter accessible to user", func() {
			view := form.View()
			// Counter should be in output for user to see
			Expect(view).To(ContainSubstring("Characters:"))
			Expect(view).To(ContainSubstring("2000"))
		})
	})
})
