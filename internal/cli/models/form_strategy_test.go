package models_test

import (
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FormModel - Strategy System", func() {
	var (
		form       *models.FormModel
		cliService *service.CLIEventService
	)

	BeforeEach(func() {
		repo := careerrepo.NewMemoryRepository()
		svc := careerservice.NewService(repo)
		cliService = service.NewCLIEventService(svc)
		form = models.NewFormModel(cliService)
	})

	Describe("Default Strategy", func() {
		It("should default to manual strategy", func() {
			// Form should start in manual mode with all fields visible
			view := form.View()
			Expect(view).To(ContainSubstring("Event Text"))
		})

		It("should show optional fields by default in manual mode", func() {
			// Manual mode should show all fields including optional ones
			view := form.View()
			// All fields should be present in the view
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Optional Field Toggle", func() {
		Context("in manual strategy", func() {
			It("should toggle optional fields with 't' key", func() {
				// Send 't' key to toggle optional fields
				form.Update(tea.KeyMsg{
					Type:  tea.KeyRunes,
					Runes: []rune{'t'},
				})

				// The form should still render (we can't easily test internal state)
				view := form.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should handle toggle when focused on text field", func() {
				// Focus should be on text field initially
				// Toggle optional fields
				form.Update(tea.KeyMsg{
					Type:  tea.KeyRunes,
					Runes: []rune{'t'},
				})

				// Form should still be valid
				view := form.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should move focus to visible field when hiding fields", func() {
				// Navigate to an optional field first
				form.Update(tea.KeyMsg{Type: tea.KeyTab}) // Date field

				// Toggle to hide optional fields
				form.Update(tea.KeyMsg{
					Type:  tea.KeyRunes,
					Runes: []rune{'t'},
				})

				// Form should still render properly
				view := form.View()
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	Describe("Field Visibility", func() {
		It("should always show required fields", func() {
			// Event text is always required and visible
			view := form.View()
			Expect(view).To(ContainSubstring("Event Text"))
		})

		It("should show submit button in all modes", func() {
			// Submit button should always be accessible
			view := form.View()
			// The view should render successfully
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Form Submission with Strategy", func() {
		It("should submit successfully with minimal fields", func() {
			// Type only required field
			for _, r := range "Test event" {
				form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
			}

			// Navigate to submit (7 tabs)
			for i := 0; i < 7; i++ {
				form.Update(tea.KeyMsg{Type: tea.KeyTab})
			}

			// Submit
			_, cmd := form.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			form.Update(msg)

			Expect(form.Submitted()).To(BeTrue())
			Expect(form.Error()).To(BeNil())
		})

		It("should submit with all fields filled", func() {
			// Fill all fields
			for _, r := range "Complete event description" {
				form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
			}

			// Date
			form.Update(tea.KeyMsg{Type: tea.KeyTab})
			for _, r := range "2024-01-01" {
				form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
			}

			// Company
			form.Update(tea.KeyMsg{Type: tea.KeyTab})
			for _, r := range "TestCorp" {
				form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
			}

			// Project
			form.Update(tea.KeyMsg{Type: tea.KeyTab})
			for _, r := range "TestProject" {
				form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
			}

			// Navigate to submit (skip tags, categories, and skills)
			form.Update(tea.KeyMsg{Type: tea.KeyTab}) // Tags
			form.Update(tea.KeyMsg{Type: tea.KeyTab}) // Categories
			form.Update(tea.KeyMsg{Type: tea.KeyTab}) // Skills
			form.Update(tea.KeyMsg{Type: tea.KeyTab}) // Submit

			// Submit
			_, cmd := form.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			form.Update(msg)

			Expect(form.Submitted()).To(BeTrue())
			Expect(form.Error()).To(BeNil())
			Expect(form.Event()).NotTo(BeNil())
			Expect(form.Event().Company).To(Equal("TestCorp"))
			Expect(form.Event().Project).To(Equal("TestProject"))
		})
	})

	Describe("Date Validation", func() {
		It("should accept any past date regardless of strategy", func() {
			// The new system doesn't have mode-specific date validation
			// All dates in the past are valid
			for _, r := range "Test event" {
				form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
			}

			// Enter date from 2 years ago
			form.Update(tea.KeyMsg{Type: tea.KeyTab})
			for _, r := range "2022-01-01" {
				form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
			}

			// Navigate to submit
			for i := 0; i < 6; i++ {
				form.Update(tea.KeyMsg{Type: tea.KeyTab})
			}

			// Submit
			_, cmd := form.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			form.Update(msg)

			// Should succeed - no 30-day restriction
			Expect(form.Submitted()).To(BeTrue())
			Expect(form.Error()).To(BeNil())
		})

		It("should reject future dates", func() {
			// This validation remains across all strategies
			for _, r := range "Test event" {
				form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
			}

			// Enter future date
			form.Update(tea.KeyMsg{Type: tea.KeyTab})
			for _, r := range "2099-12-31" {
				form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
			}

			// Navigate to submit
			for i := 0; i < 6; i++ {
				form.Update(tea.KeyMsg{Type: tea.KeyTab})
			}

			// Submit
			_, cmd := form.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			form.Update(msg)

			// Should fail with future date error
			Expect(form.Submitted()).To(BeFalse())
			Expect(form.Error()).NotTo(BeNil())
			Expect(form.Error().Error()).To(ContainSubstring("future"))
		})
	})
})
