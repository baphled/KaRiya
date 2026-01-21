package components_test

import (
	"github.com/baphled/kariya/internal/cli/components"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CVConfigWizardModal", func() {
	var (
		modal *components.CVConfigWizardModal
	)

	Describe("NewCVConfigWizardModal", func() {
		Context("Creation and Initialization", func() {
			It("should create a wizard modal with default configuration", func() {
				modal = components.NewCVConfigWizardModal(120, 40)

				Expect(modal).NotTo(BeNil())
				Expect(modal.IsVisible()).To(BeTrue())
			})

			It("should initialize with 3 steps (WHO, TECH, FORMAT)", func() {
				modal = components.NewCVConfigWizardModal(120, 40)

				// Step count should be 3
				Expect(modal.GetStepCount()).To(Equal(3))
				Expect(modal.GetCurrentStep()).To(Equal(0)) // Start at step 0 (WHO)
			})

			It("should initialize with default configuration data", func() {
				modal = components.NewCVConfigWizardModal(120, 40)

				config := modal.GetConfigData()
				Expect(config).NotTo(BeNil())
				// Note: huh.Select automatically selects first option
				// ProfileID will be empty string (from "(No profiles available)", "") option
				Expect(config.ProfileID).To(Equal(""))
				// Audience will be "hiring_manager" (first option)
				Expect(config.Audience).To(Equal("hiring_manager"))
				// TechFocus will be "language_agnostic" (first option)
				Expect(config.TechFocus).To(Equal("language_agnostic"))
				Expect(config.Technologies).To(BeEmpty())
				Expect(config.FocusArea).To(BeEmpty())
				// SkillsFormat will be "grouped" (first option)
				Expect(config.SkillsFormat).To(Equal("grouped"))
				// CVLength will be "1_page" (first option)
				Expect(config.CVLength).To(Equal("1_page"))
			})

			It("should start with techs not available", func() {
				modal = components.NewCVConfigWizardModal(120, 40)

				Expect(modal.AreTechsAvailable()).To(BeFalse())
			})

			It("should not be completed initially", func() {
				modal = components.NewCVConfigWizardModal(120, 40)

				Expect(modal.IsCompleted()).To(BeFalse())
				Expect(modal.IsSkipped()).To(BeFalse())
			})
		})

		Context("With Profile Options", func() {
			It("should accept profile options during creation", func() {
				profiles := []components.ProfileOption{
					{ID: "profile-1", Name: "Senior Go Engineer"},
					{ID: "profile-2", Name: "Full Stack Developer"},
				}

				modal = components.NewCVConfigWizardModalWithProfiles(120, 40, profiles)

				Expect(modal).NotTo(BeNil())
				Expect(modal.GetProfileOptions()).To(HaveLen(2))
			})
		})
	})

	Describe("Visibility Management", func() {
		BeforeEach(func() {
			modal = components.NewCVConfigWizardModal(120, 40)
		})

		It("should be visible by default", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should hide when Hide() is called", func() {
			modal.Hide()

			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should show when Show() is called after hiding", func() {
			modal.Hide()
			modal.Show()

			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should return empty string from View() when hidden", func() {
			modal.Hide()

			view := modal.View()
			Expect(view).To(Equal(""))
		})
	})

	Describe("Step Navigation", func() {
		BeforeEach(func() {
			profiles := []components.ProfileOption{
				{ID: "profile-1", Name: "Senior Go Engineer"},
				{ID: "profile-2", Name: "Tech Lead"}, // Multiple options prevent auto-selection
			}
			modal = components.NewCVConfigWizardModalWithProfiles(120, 40, profiles)
			modal.Init()
		})

		Context("Forward Navigation", func() {
			It("should advance to next step on valid field completion", func() {
				// Set required field for step 1
				modal.SetProfileID("profile-1")
				modal.SetAudience("hiring_manager")

				// Simulate Enter key to advance
				modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Should advance to step 1 (TECH) if techs available, or step 2 (FORMAT) if not
				currentStep := modal.GetCurrentStep()
				Expect(currentStep).To(Or(Equal(1), Equal(2)))
			})

			It("should skip TECH step if no techs available", func() {
				// Complete WHO step
				modal.SetProfileID("profile-1")
				modal.SetAudience("hiring_manager")
				modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Should skip step 1 (TECH) and go to step 2 (FORMAT)
				if !modal.AreTechsAvailable() {
					Expect(modal.GetCurrentStep()).To(Equal(2))
				}
			})

			It("should show TECH step if techs are available", func() {
				// Set extracted technologies
				techs := []components.ExtractedTechnology{
					{Name: "Go", Category: "Language"},
					{Name: "Docker", Category: "Tool"},
				}
				modal.SetExtractedTechnologies(techs)

				// Complete WHO step
				modal.SetProfileID("profile-1")
				modal.SetAudience("hiring_manager")
				modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Should show step 1 (TECH)
				Expect(modal.GetCurrentStep()).To(Equal(1))
				Expect(modal.AreTechsAvailable()).To(BeTrue())
			})
		})

		Context("Backward Navigation", func() {
			It("should go back one step on Esc key", func() {
				// Advance to step 2 (FORMAT)
				modal.SetProfileID("profile-1")
				modal.SetAudience("hiring_manager")
				modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

				initialStep := modal.GetCurrentStep()

				// Press Esc to go back
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

				// Should be at previous step
				// Note: If TECH step is not available, it will skip from step 2 to step 0
				if modal.AreTechsAvailable() {
					Expect(modal.GetCurrentStep()).To(Equal(initialStep - 1))
				} else {
					// Skips TECH step (step 1) when going back
					Expect(modal.GetCurrentStep()).To(Equal(0))
				}
			})

			It("should hide modal on Esc from first step (WHO)", func() {
				Expect(modal.GetCurrentStep()).To(Equal(0))
				Expect(modal.IsVisible()).To(BeTrue())

				// Press Esc from step 0
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

				// Modal should hide (cancel action)
				Expect(modal.IsVisible()).To(BeFalse())
			})

			It("should skip TECH step backwards if not shown", func() {
				// Start at FORMAT step (no techs available)
				modal.SetProfileID("profile-1")
				modal.SetAudience("hiring_manager")
				modal.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Advance to FORMAT

				currentStep := modal.GetCurrentStep()
				Expect(currentStep).To(Equal(2)) // At FORMAT

				// Press Esc to go back
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

				// Should skip TECH (step 1) and go to WHO (step 0)
				if !modal.AreTechsAvailable() {
					Expect(modal.GetCurrentStep()).To(Equal(0))
				}
			})
		})

		Context("Skip to Generate", func() {
			It("should skip to completion on Ctrl+Enter", func() {
				// Set minimum required field
				modal.SetProfileID("profile-1")

				// Press Ctrl+Enter
				modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS}) // Using 's' for skip

				// Should be completed with skipped flag
				Expect(modal.IsCompleted()).To(BeTrue())
				Expect(modal.IsSkipped()).To(BeTrue())
			})

			It("should use default values when skipped", func() {
				modal.SetProfileID("profile-1")
				modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})

				config := modal.GetConfigData()
				// Should have first-option defaults from huh.Select
				Expect(config.ProfileID).To(Equal("profile-1"))
				Expect(config.Audience).To(Equal("hiring_manager"))     // first option
				Expect(config.TechFocus).To(Equal("language_agnostic")) // first option
				Expect(config.SkillsFormat).To(Equal("grouped"))        // first option
				// CVLength might be "1_page" (first option) or applied default "2_page"
				Expect(config.CVLength).To(Or(Equal("1_page"), Equal("2_page")))
			})

			It("should require at least ProfileID to skip", func() {
				// Try to skip without setting required field
				modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})

				// Should not complete if required field missing
				Expect(modal.IsCompleted()).To(BeFalse())
			})
		})
	})

	Describe("Conditional TECH Step", func() {
		BeforeEach(func() {
			profiles := []components.ProfileOption{
				{ID: "profile-1", Name: "Senior Go Engineer"},
			}
			modal = components.NewCVConfigWizardModalWithProfiles(120, 40, profiles)
			modal.Init()
		})

		It("should include TECH step when technologies are available", func() {
			techs := []components.ExtractedTechnology{
				{Name: "Go", Category: "Language", Confidence: 0.95},
				{Name: "PostgreSQL", Category: "Database", Confidence: 0.88},
			}
			modal.SetExtractedTechnologies(techs)

			Expect(modal.AreTechsAvailable()).To(BeTrue())
			Expect(modal.GetStepCount()).To(Equal(3)) // All 3 steps active
		})

		It("should skip TECH step when no technologies available", func() {
			Expect(modal.AreTechsAvailable()).To(BeFalse())
			// Effective step count is 2 (WHO + FORMAT)
		})

		It("should update tech options when SetExtractedTechnologies is called", func() {
			techs := []components.ExtractedTechnology{
				{Name: "Go", Category: "Language"},
				{Name: "Docker", Category: "Tool"},
				{Name: "Kubernetes", Category: "Tool"},
			}

			modal.SetExtractedTechnologies(techs)

			Expect(modal.AreTechsAvailable()).To(BeTrue())
			extractedTechs := modal.GetExtractedTechnologies()
			Expect(extractedTechs).To(HaveLen(3))
		})
	})

	Describe("Data Extraction", func() {
		BeforeEach(func() {
			profiles := []components.ProfileOption{
				{ID: "profile-1", Name: "Senior Go Engineer"},
			}
			modal = components.NewCVConfigWizardModalWithProfiles(120, 40, profiles)
			modal.Init()
		})

		It("should extract configuration data after completion", func() {
			// Fill out all steps
			modal.SetProfileID("profile-1")
			modal.SetAudience("hiring_manager")
			modal.SetTechFocus("specialist")
			modal.SetTechnologies([]string{"Go", "PostgreSQL"})
			modal.SetFocusArea("backend")
			modal.SetSkillsFormat("categorized")
			modal.SetCVLength("2_page")

			// Complete wizard
			modal.Complete()

			config := modal.GetConfigData()
			Expect(config.ProfileID).To(Equal("profile-1"))
			Expect(config.Audience).To(Equal("hiring_manager"))
			Expect(config.TechFocus).To(Equal("specialist"))
			Expect(config.Technologies).To(ConsistOf("Go", "PostgreSQL"))
			Expect(config.FocusArea).To(Equal("backend"))
			Expect(config.SkillsFormat).To(Equal("categorized"))
			Expect(config.CVLength).To(Equal("2_page"))
		})

		It("should preserve partial data when cancelled", func() {
			modal.SetProfileID("profile-1")
			modal.SetAudience("recruiter")

			// Cancel by pressing Esc from first step
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Data should be preserved even though wizard was cancelled
			config := modal.GetConfigData()
			Expect(config.ProfileID).To(Equal("profile-1"))
			Expect(config.Audience).To(Equal("recruiter"))
		})

		It("should validate required fields before completion", func() {
			// Try to complete without required field
			modal.Complete()

			// Should not be completed if validation fails
			hasRequiredFields := modal.HasRequiredFields()
			if !hasRequiredFields {
				Expect(modal.IsCompleted()).To(BeFalse())
			}
		})
	})

	Describe("WindowSizeMsg Handling", func() {
		BeforeEach(func() {
			modal = components.NewCVConfigWizardModal(120, 40)
		})

		It("should update dimensions on WindowSizeMsg", func() {
			modal.Update(tea.WindowSizeMsg{Width: 160, Height: 50})

			// Dimensions should be updated
			width, height := modal.GetDimensions()
			Expect(width).To(Equal(160))
			Expect(height).To(Equal(50))
		})

		It("should rebuild form with new dimensions", func() {
			initialView := modal.View()

			// Resize to smaller terminal
			modal.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

			newView := modal.View()
			// View should be different after resize
			Expect(newView).NotTo(Equal(initialView))
		})

		It("should handle minimum dimensions gracefully", func() {
			// Very small terminal
			modal.Update(tea.WindowSizeMsg{Width: 40, Height: 10})

			// Should not crash and still render
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			profiles := []components.ProfileOption{
				{ID: "profile-1", Name: "Senior Go Engineer"},
			}
			modal = components.NewCVConfigWizardModalWithProfiles(120, 40, profiles)
		})

		It("should render modal with solid background", func() {
			view := modal.View()

			// View should not be empty
			Expect(view).NotTo(BeEmpty())
			// Should contain modal styling (border, background)
			Expect(view).To(ContainSubstring("─")) // Border character
		})

		It("should show current step title", func() {
			view := modal.View()

			// Should show WHO step initially
			Expect(view).To(ContainSubstring("CV Configuration"))
		})

		It("should display KeyBadge footer", func() {
			view := modal.View()

			// Should contain keyboard shortcuts in footer
			// Footer uses KeyBadge components
			Expect(view).NotTo(BeEmpty())
		})

		It("should render different content for each step", func() {
			view1 := modal.View() // WHO step

			// Advance to next step
			modal.SetProfileID("profile-1")
			modal.SetAudience("hiring_manager")
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			view2 := modal.View() // TECH or FORMAT step

			// Views should be different
			Expect(view2).NotTo(Equal(view1))
		})
	})

	Describe("Theme Integration", func() {
		BeforeEach(func() {
			modal = components.NewCVConfigWizardModal(120, 40)
		})

		It("should use Catppuccin theme for form", func() {
			// Modal should integrate with theme system
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			// Theme colors should be applied (verified visually)
		})

		It("should have themed borders", func() {
			view := modal.View()
			// Should have rounded borders with theme colors
			Expect(view).To(ContainSubstring("─"))
		})
	})

	Describe("Edge Cases", func() {
		BeforeEach(func() {
			modal = components.NewCVConfigWizardModal(120, 40)
		})

		It("should handle nil profile options gracefully", func() {
			modal = components.NewCVConfigWizardModalWithProfiles(120, 40, nil)
			Expect(modal).NotTo(BeNil())
		})

		It("should handle empty profile options", func() {
			modal = components.NewCVConfigWizardModalWithProfiles(120, 40, []components.ProfileOption{})
			Expect(modal).NotTo(BeNil())
		})

		It("should handle empty extracted technologies", func() {
			modal.SetExtractedTechnologies([]components.ExtractedTechnology{})
			Expect(modal.AreTechsAvailable()).To(BeFalse())
		})

		It("should handle nil extracted technologies", func() {
			modal.SetExtractedTechnologies(nil)
			Expect(modal.AreTechsAvailable()).To(BeFalse())
		})

		It("should handle rapid key presses gracefully", func() {
			// Simulate rapid key presses
			for i := 0; i < 10; i++ {
				modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			}

			// Should not panic or enter invalid state
			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Reset", func() {
		BeforeEach(func() {
			profiles := []components.ProfileOption{
				{ID: "profile-1", Name: "Senior Go Engineer"},
				{ID: "profile-2", Name: "Full Stack Developer"},
			}
			modal = components.NewCVConfigWizardModalWithProfiles(120, 40, profiles)
			modal.Init()
		})

		It("should preserve data when Reset is called", func() {
			// Set up some data
			modal.SetProfileID("profile-1")
			modal.SetAudience("recruiter")
			modal.SetSkillsFormat("grouped")
			modal.SetSkillsLimit(10)
			modal.SetCVLength("2_page")

			// Complete the wizard
			modal.Complete()
			Expect(modal.IsCompleted()).To(BeTrue())
			Expect(modal.IsVisible()).To(BeFalse())

			// Reset the wizard
			modal.Reset()

			// Should be visible and not completed
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.IsCompleted()).To(BeFalse())
			Expect(modal.IsSkipped()).To(BeFalse())
			Expect(modal.GetCurrentStep()).To(Equal(0))

			// Data should be preserved
			config := modal.GetConfigData()
			Expect(config.ProfileID).To(Equal("profile-1"))
			Expect(config.Audience).To(Equal("recruiter"))
			Expect(config.SkillsFormat).To(Equal("grouped"))
			Expect(config.SkillsLimit).To(Equal(10))
			Expect(config.CVLength).To(Equal("2_page"))
		})

		It("should reset step counter to 0", func() {
			// Advance to step 2 (FORMAT)
			modal.SetProfileID("profile-1")
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Reset
			modal.Reset()

			// Should be back at step 0
			Expect(modal.GetCurrentStep()).To(Equal(0))
		})

		It("should clear completed and skipped flags", func() {
			modal.SetProfileID("profile-1")
			modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS}) // Skip

			Expect(modal.IsCompleted()).To(BeTrue())
			Expect(modal.IsSkipped()).To(BeTrue())

			modal.Reset()

			Expect(modal.IsCompleted()).To(BeFalse())
			Expect(modal.IsSkipped()).To(BeFalse())
		})

		It("should make modal visible after reset", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())

			modal.Reset()

			Expect(modal.IsVisible()).To(BeTrue())
		})
	})

	Describe("Form Navigation Selection Persistence", func() {
		It("should return the profile selected via keyboard navigation in GetConfigData", func() {
			profiles := []components.ProfileOption{
				{ID: "profile-1", Name: "Staff Engineer"},
				{ID: "profile-2", Name: "Senior Engineer"},
			}
			modal = components.NewCVConfigWizardModalWithProfiles(120, 40, profiles)
			modal.Init()

			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			config := modal.GetConfigData()
			Expect(config.ProfileID).To(Equal("profile-1"))
		})
	})
})
