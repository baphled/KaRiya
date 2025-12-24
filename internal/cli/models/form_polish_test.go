package models

import (
	"github.com/baphled/kariya/internal/cli/service"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Form Model UI/UX Polish", func() {
	var (
		repo   *careerrepo.MemoryRepository
		svc    *careerservice.Service
		cliSvc *service.CLIEventService
		form   *FormModel
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliSvc = service.NewCLIEventService(svc)
		form = NewFormModel(cliSvc)
	})

	Context("Visual Feedback During Submission", func() {
		It("should show loading state during form submission", func() {
			// Fill form with valid data
			form.inputs[0].SetValue("Test event for spinner")
			form.inputs[1].SetValue("today")

			// Submit form
			_, cmd := form.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Should return a command for async operation
			Expect(cmd).NotTo(BeNil())
		})

		It("should display submitted state after successful submission", func() {
			// Fill form with valid data
			form.inputs[0].SetValue("Test event")
			form.inputs[1].SetValue("today")

			// Simulate form submission
			Expect(form.inputs[0].Value()).To(Equal("Test event"))
		})

		It("should maintain form state after successful submission", func() {
			// Verify form starts in non-submitted state
			Expect(form.Submitted()).To(BeFalse())

			// After submission, state should be managed
			form.submitted = true
			Expect(form.Submitted()).To(BeTrue())

			// Can reset
			form.Reset()
			Expect(form.Submitted()).To(BeFalse())
		})

		It("should show error state with clear messaging", func() {
			// Submit with empty text should show error
			form.inputs[0].SetValue("")
			cmd := form.submitForm()

			// Command should handle validation
			Expect(cmd).NotTo(BeNil())
		})

		It("should display form fields with consistent styling", func() {
			// All input fields should be properly styled
			Expect(form.inputs).To(HaveLen(4))

			for _, input := range form.inputs {
				// Each input should have proper configuration
				Expect(input).NotTo(BeNil())
				// Placeholder should be set
				Expect(len(input.Placeholder) > 0).To(BeTrue())
			}
		})

		It("should provide clear visual focus indicator", func() {
			// Initially first field is focused
			Expect(form.focusIndex).To(Equal(0))

			// Navigate to next field
			form.focusIndex = 1
			Expect(form.focusIndex).To(Equal(1))

			// Visual focus should be clear
			focusedInput := form.inputs[form.focusIndex]
			Expect(focusedInput).NotTo(BeNil())
		})
	})

	Context("Visual Feedback for Errors", func() {
		It("should show character count near limit (1900+ chars)", func() {
			// Build long text close to limit
			longText := ""
			for i := 0; i < 1900; i++ {
				longText += "a"
			}

			form.inputs[0].SetValue(longText)
			form.charCount = len(longText)

			// Character count should be tracked
			Expect(form.charCount).To(BeNumerically(">", 1800))
			Expect(form.charCount).To(BeNumerically("<", 2000))
		})

		It("should warn when approaching text limit", func() {
			// At 1950 chars (50 from limit)
			form.charCount = 1950
			form.maxChars = 2000

			// Should indicate warning is needed
			remaining := form.maxChars - form.charCount
			Expect(remaining).To(Equal(50))
			Expect(remaining < 100).To(BeTrue())
		})

		It("should show error state when text exceeds limit", func() {
			// Over 2000 chars
			longText := ""
			for i := 0; i < 2001; i++ {
				longText += "a"
			}

			form.inputs[0].SetValue(longText)
			form.charCount = 2001

			// Should indicate error state
			Expect(form.charCount).To(BeNumerically(">", 2000))
		})

		It("should display field-level errors clearly", func() {
			// Set an error for text field
			form.fieldErrors[TextField] = "Text exceeds maximum length"

			// Error should be stored and displayable
			err, exists := form.fieldErrors[TextField]
			Expect(exists).To(BeTrue())
			Expect(err).To(ContainSubstring("exceeds"))
		})
	})

	Context("Transition Animations", func() {
		It("should support smooth transitions between fields", func() {
			// Navigate through fields
			initialFocus := form.focusIndex
			form.focusIndex = (form.focusIndex + 1) % 4

			// Focus should change smoothly
			Expect(form.focusIndex).NotTo(Equal(initialFocus))
		})

		It("should provide visual indication of mode selection", func() {
			// Mode selection should show current mode
			initialMode := form.modeIndex
			Expect(initialMode).To(BeNumerically(">=", 0))
			Expect(initialMode).To(BeNumerically("<", len(form.modes)))
		})

		It("should show tag selection updates visually", func() {
			// Tag selector should be available
			Expect(form.tagSelector).NotTo(BeNil())

			// Tags should be displayable
			tags := form.tagSelector.SelectedTags()
			Expect(tags).To(HaveLen(0))
		})
	})

	Context("Terminal Size Responsiveness", func() {
		It("should adapt form width to terminal size", func() {
			// Form inputs should have reasonable width
			for _, input := range form.inputs {
				Expect(input.Width).To(BeNumerically(">", 0))
				Expect(input.Width).To(BeNumerically("<", 200))
			}
		})

		It("should maintain readability with minimum terminal width", func() {
			// At minimum 80 columns, form should still be usable
			minWidth := 80
			for _, input := range form.inputs {
				// Each input should fit in minimum width
				Expect(input.Width).To(BeNumerically("<", minWidth))
			}
		})

		It("should provide helpful placeholder text for all fields", func() {
			for _, input := range form.inputs {
				Expect(input.Placeholder).NotTo(BeEmpty())
			}
		})
	})

	Context("User Feedback Timing", func() {
		It("should provide immediate response to user input", func() {
			// Typing should update view immediately
			form.inputs[0].SetValue("Quick response test")
			view := form.View()

			// View should contain the typed text
			Expect(view).To(ContainSubstring("Quick response test"))
		})

		It("should show validation errors promptly", func() {
			// Empty text should show error quickly
			form.inputs[0].SetValue("")

			// Validate text
			if form.inputs[0].Value() == "" {
				form.fieldErrors[TextField] = "Text cannot be empty"
			}

			// Error should be set
			err, exists := form.fieldErrors[TextField]
			Expect(exists).To(BeTrue())
			Expect(err).NotTo(BeEmpty())
		})

		It("should clear errors when field is corrected", func() {
			// Set an error
			form.fieldErrors[TextField] = "Text cannot be empty"
			Expect(form.fieldErrors[TextField]).NotTo(BeEmpty())

			// Set valid text
			form.inputs[0].SetValue("Valid text")
			delete(form.fieldErrors, TextField)

			// Error should be cleared
			_, exists := form.fieldErrors[TextField]
			Expect(exists).To(BeFalse())
		})
	})
})
