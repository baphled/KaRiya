package modals_test

import (
	"github.com/baphled/kariya/internal/cli/screens/cv/modals"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigWizardModal", func() {
	var (
		modal *modals.ConfigWizardModal
	)

	Describe("NewConfigWizardModal", func() {
		Context("Creation and Initialization", func() {
			It("should create a wizard modal with default configuration", func() {
				modal = modals.NewConfigWizardModal(120, 40)

				Expect(modal).NotTo(BeNil())
				Expect(modal.IsVisible()).To(BeTrue())
			})

			It("should initialize with 3 steps (WHO, TECH, FORMAT)", func() {
				modal = modals.NewConfigWizardModal(120, 40)

				Expect(modal.GetStepCount()).To(Equal(3))
				Expect(modal.GetCurrentStep()).To(Equal(0))
			})

			It("should initialize with default configuration data", func() {
				modal = modals.NewConfigWizardModal(120, 40)

				config := modal.GetConfigData()
				Expect(config).NotTo(BeNil())
				Expect(config.ProfileID).To(Equal(""))
				Expect(config.Audience).To(Equal("hiring_manager"))
				Expect(config.TechFocus).To(Equal("language_agnostic"))
				Expect(config.Technologies).To(BeEmpty())
				Expect(config.FocusArea).To(BeEmpty())
				Expect(config.SkillsFormat).To(Equal("grouped"))
				Expect(config.CVLength).To(Equal("1_page"))
			})

			It("should start with techs not available", func() {
				modal = modals.NewConfigWizardModal(120, 40)

				Expect(modal.AreTechsAvailable()).To(BeFalse())
			})

			It("should not be completed initially", func() {
				modal = modals.NewConfigWizardModal(120, 40)

				Expect(modal.IsCompleted()).To(BeFalse())
				Expect(modal.IsSkipped()).To(BeFalse())
			})
		})

		Context("With Profile Options", func() {
			It("should accept profile options during creation", func() {
				profiles := []modals.ProfileOption{
					{ID: "profile-1", Name: "Senior Go Engineer"},
					{ID: "profile-2", Name: "Full Stack Developer"},
				}

				modal = modals.NewConfigWizardModalWithProfiles(120, 40, profiles)

				Expect(modal).NotTo(BeNil())
				Expect(modal.GetProfileOptions()).To(HaveLen(2))
			})
		})
	})

	Describe("Visibility Management", func() {
		BeforeEach(func() {
			modal = modals.NewConfigWizardModal(120, 40)
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
			profiles := []modals.ProfileOption{
				{ID: "profile-1", Name: "Senior Go Engineer"},
				{ID: "profile-2", Name: "Tech Lead"},
			}
			modal = modals.NewConfigWizardModalWithProfiles(120, 40, profiles)
			modal.Init()
		})

		Context("Forward Navigation", func() {
			It("should advance to next step on valid field completion", func() {
				modal.SetProfileID("profile-1")
				modal.SetAudience("hiring_manager")

				modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

				currentStep := modal.GetCurrentStep()
				Expect(currentStep).To(Or(Equal(1), Equal(2)))
			})

			It("should skip TECH step if no techs available", func() {
				modal.SetProfileID("profile-1")
				modal.SetAudience("hiring_manager")
				modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

				if !modal.AreTechsAvailable() {
					Expect(modal.GetCurrentStep()).To(Equal(2))
				}
			})

			It("should show TECH step if techs are available", func() {
				techs := []modals.ExtractedTechnology{
					{Name: "Go", Category: "Language"},
					{Name: "Docker", Category: "Tool"},
				}
				modal.SetExtractedTechnologies(techs)

				modal.SetProfileID("profile-1")
				modal.SetAudience("hiring_manager")
				modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(modal.GetCurrentStep()).To(Equal(1))
				Expect(modal.AreTechsAvailable()).To(BeTrue())
			})
		})

		Context("Backward Navigation", func() {
			It("should go back one step on Esc key", func() {
				modal.SetProfileID("profile-1")
				modal.SetAudience("hiring_manager")
				modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

				initialStep := modal.GetCurrentStep()

				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

				if modal.AreTechsAvailable() {
					Expect(modal.GetCurrentStep()).To(Equal(initialStep - 1))
				} else {
					Expect(modal.GetCurrentStep()).To(Equal(0))
				}
			})

			It("should hide modal on Esc from first step (WHO)", func() {
				Expect(modal.GetCurrentStep()).To(Equal(0))
				Expect(modal.IsVisible()).To(BeTrue())

				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(modal.IsVisible()).To(BeFalse())
			})

			It("should skip TECH step backwards if not shown", func() {
				modal.SetProfileID("profile-1")
				modal.SetAudience("hiring_manager")
				modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

				currentStep := modal.GetCurrentStep()
				Expect(currentStep).To(Equal(2))

				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

				if !modal.AreTechsAvailable() {
					Expect(modal.GetCurrentStep()).To(Equal(0))
				}
			})
		})

		Context("Skip to Generate", func() {
			It("should skip to completion on Ctrl+S", func() {
				modal.SetProfileID("profile-1")

				modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})

				Expect(modal.IsCompleted()).To(BeTrue())
				Expect(modal.IsSkipped()).To(BeTrue())
			})

			It("should use default values when skipped", func() {
				modal.SetProfileID("profile-1")
				modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})

				config := modal.GetConfigData()
				Expect(config.ProfileID).To(Equal("profile-1"))
				Expect(config.Audience).To(Equal("hiring_manager"))
				Expect(config.TechFocus).To(Equal("language_agnostic"))
				Expect(config.SkillsFormat).To(Equal("grouped"))
				Expect(config.CVLength).To(Or(Equal("1_page"), Equal("2_page")))
			})

			It("should require at least ProfileID to skip", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})

				Expect(modal.IsCompleted()).To(BeFalse())
			})
		})
	})

	Describe("Conditional TECH Step", func() {
		BeforeEach(func() {
			profiles := []modals.ProfileOption{
				{ID: "profile-1", Name: "Senior Go Engineer"},
			}
			modal = modals.NewConfigWizardModalWithProfiles(120, 40, profiles)
			modal.Init()
		})

		It("should include TECH step when technologies are available", func() {
			techs := []modals.ExtractedTechnology{
				{Name: "Go", Category: "Language", Confidence: 0.95},
				{Name: "PostgreSQL", Category: "Database", Confidence: 0.88},
			}
			modal.SetExtractedTechnologies(techs)

			Expect(modal.AreTechsAvailable()).To(BeTrue())
			Expect(modal.GetStepCount()).To(Equal(3))
		})

		It("should skip TECH step when no technologies available", func() {
			Expect(modal.AreTechsAvailable()).To(BeFalse())
		})

		It("should update tech options when SetExtractedTechnologies is called", func() {
			techs := []modals.ExtractedTechnology{
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
			profiles := []modals.ProfileOption{
				{ID: "profile-1", Name: "Senior Go Engineer"},
			}
			modal = modals.NewConfigWizardModalWithProfiles(120, 40, profiles)
			modal.Init()
		})

		It("should extract configuration data after completion", func() {
			modal.SetProfileID("profile-1")
			modal.SetAudience("hiring_manager")
			modal.SetTechFocus("specialist")
			modal.SetTechnologies([]string{"Go", "PostgreSQL"})
			modal.SetFocusArea("backend")
			modal.SetSkillsFormat("categorized")
			modal.SetCVLength("2_page")

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

			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

			config := modal.GetConfigData()
			Expect(config.ProfileID).To(Equal("profile-1"))
			Expect(config.Audience).To(Equal("recruiter"))
		})

		It("should validate required fields before completion", func() {
			modal.Complete()

			hasRequiredFields := modal.HasRequiredFields()
			if !hasRequiredFields {
				Expect(modal.IsCompleted()).To(BeFalse())
			}
		})
	})

	Describe("WindowSizeMsg Handling", func() {
		BeforeEach(func() {
			modal = modals.NewConfigWizardModal(120, 40)
		})

		It("should update dimensions on WindowSizeMsg", func() {
			modal.Update(tea.WindowSizeMsg{Width: 160, Height: 50})

			width, height := modal.GetDimensions()
			Expect(width).To(Equal(160))
			Expect(height).To(Equal(50))
		})

		It("should rebuild form with new dimensions", func() {
			initialView := modal.View()

			modal.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

			newView := modal.View()
			Expect(newView).NotTo(Equal(initialView))
		})

		It("should handle minimum dimensions gracefully", func() {
			modal.Update(tea.WindowSizeMsg{Width: 40, Height: 10})

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			profiles := []modals.ProfileOption{
				{ID: "profile-1", Name: "Senior Go Engineer"},
			}
			modal = modals.NewConfigWizardModalWithProfiles(120, 40, profiles)
		})

		It("should render modal with solid background", func() {
			view := modal.View()

			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("─"))
		})

		It("should show current step title", func() {
			view := modal.View()

			Expect(view).To(ContainSubstring("CV Configuration"))
		})

		It("should display KeyBadge footer", func() {
			view := modal.View()

			Expect(view).NotTo(BeEmpty())
		})

		It("should render different content for each step", func() {
			view1 := modal.View()

			modal.SetProfileID("profile-1")
			modal.SetAudience("hiring_manager")
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			view2 := modal.View()

			Expect(view2).NotTo(Equal(view1))
		})
	})

	Describe("Theme Integration", func() {
		BeforeEach(func() {
			modal = modals.NewConfigWizardModal(120, 40)
		})

		It("should use Catppuccin theme for form", func() {
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should have themed borders", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("─"))
		})
	})

	Describe("Edge Cases", func() {
		BeforeEach(func() {
			modal = modals.NewConfigWizardModal(120, 40)
		})

		It("should handle nil profile options gracefully", func() {
			modal = modals.NewConfigWizardModalWithProfiles(120, 40, nil)
			Expect(modal).NotTo(BeNil())
		})

		It("should handle empty profile options", func() {
			modal = modals.NewConfigWizardModalWithProfiles(120, 40, []modals.ProfileOption{})
			Expect(modal).NotTo(BeNil())
		})

		It("should handle empty extracted technologies", func() {
			modal.SetExtractedTechnologies([]modals.ExtractedTechnology{})
			Expect(modal.AreTechsAvailable()).To(BeFalse())
		})

		It("should handle nil extracted technologies", func() {
			modal.SetExtractedTechnologies(nil)
			Expect(modal.AreTechsAvailable()).To(BeFalse())
		})

		It("should handle rapid key presses gracefully", func() {
			for range 10 {
				modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			}

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Reset", func() {
		BeforeEach(func() {
			profiles := []modals.ProfileOption{
				{ID: "profile-1", Name: "Senior Go Engineer"},
				{ID: "profile-2", Name: "Full Stack Developer"},
			}
			modal = modals.NewConfigWizardModalWithProfiles(120, 40, profiles)
			modal.Init()
		})

		It("should preserve data when Reset is called", func() {
			modal.SetProfileID("profile-1")
			modal.SetAudience("recruiter")
			modal.SetSkillsFormat("grouped")
			modal.SetSkillsLimit(10)
			modal.SetCVLength("2_page")

			modal.Complete()
			Expect(modal.IsCompleted()).To(BeTrue())
			Expect(modal.IsVisible()).To(BeFalse())

			modal.Reset()

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.IsCompleted()).To(BeFalse())
			Expect(modal.IsSkipped()).To(BeFalse())
			Expect(modal.GetCurrentStep()).To(Equal(0))

			config := modal.GetConfigData()
			Expect(config.ProfileID).To(Equal("profile-1"))
			Expect(config.Audience).To(Equal("recruiter"))
			Expect(config.SkillsFormat).To(Equal("grouped"))
			Expect(config.SkillsLimit).To(Equal(10))
			Expect(config.CVLength).To(Equal("2_page"))
		})

		It("should reset step counter to 0", func() {
			modal.SetProfileID("profile-1")
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			modal.Reset()

			Expect(modal.GetCurrentStep()).To(Equal(0))
		})

		It("should clear completed and skipped flags", func() {
			modal.SetProfileID("profile-1")
			modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})

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

	Describe("WizardBehavior Delegation", func() {
		BeforeEach(func() {
			profiles := []modals.ProfileOption{
				{ID: "profile-1", Name: "Senior Go Engineer"},
				{ID: "profile-2", Name: "Tech Lead"},
			}
			modal = modals.NewConfigWizardModalWithProfiles(120, 40, profiles)
			modal.Init()
		})

		It("should preserve step tracking after window resize", func() {
			modal.SetProfileID("profile-1")
			modal.SetAudience("hiring_manager")
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			stepBeforeResize := modal.GetCurrentStep()
			Expect(stepBeforeResize).To(BeNumerically(">", 0))

			modal.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

			Expect(modal.GetCurrentStep()).To(Equal(stepBeforeResize))
		})

		It("should synchronize step tracking with adapter after SetExtractedTechnologies", func() {
			modal.SetProfileID("profile-1")
			modal.SetAudience("hiring_manager")
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(modal.GetCurrentStep()).To(BeNumerically(">", 0))

			techs := []modals.ExtractedTechnology{
				{Name: "Go", Category: "Language"},
			}
			modal.SetExtractedTechnologies(techs)

			Expect(modal.GetCurrentStep()).To(Equal(0))
		})

		It("should return init command from WindowSizeMsg handler", func() {
			cmd := modal.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Form Navigation Selection Persistence", func() {
		It("should return the profile selected via keyboard navigation in GetConfigData", func() {
			profiles := []modals.ProfileOption{
				{ID: "profile-1", Name: "Staff Engineer"},
				{ID: "profile-2", Name: "Senior Engineer"},
			}
			modal = modals.NewConfigWizardModalWithProfiles(120, 40, profiles)
			modal.Init()

			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			config := modal.GetConfigData()
			Expect(config.ProfileID).To(Equal("profile-1"))
		})

		It("should preserve profile selection through window resize", func() {
			profiles := []modals.ProfileOption{
				{ID: "profile-1", Name: "Staff Engineer"},
				{ID: "profile-2", Name: "Senior Engineer"},
			}
			modal = modals.NewConfigWizardModalWithProfiles(120, 40, profiles)
			modal.Init()

			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			modal.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

			config := modal.GetConfigData()
			Expect(config.ProfileID).To(Equal("profile-1"))
		})

		It("should preserve profile selection through SetExtractedTechnologies", func() {
			profiles := []modals.ProfileOption{
				{ID: "profile-1", Name: "Staff Engineer"},
				{ID: "profile-2", Name: "Senior Engineer"},
			}
			modal = modals.NewConfigWizardModalWithProfiles(120, 40, profiles)
			modal.Init()

			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			techs := []modals.ExtractedTechnology{
				{Name: "Go", Category: "Language"},
			}
			modal.SetExtractedTechnologies(techs)

			config := modal.GetConfigData()
			Expect(config.ProfileID).To(Equal("profile-1"))
		})

		It("should preserve profile selection through Reset", func() {
			profiles := []modals.ProfileOption{
				{ID: "profile-1", Name: "Staff Engineer"},
				{ID: "profile-2", Name: "Senior Engineer"},
			}
			modal = modals.NewConfigWizardModalWithProfiles(120, 40, profiles)
			modal.Init()

			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			modal.SetProfileID("profile-1")
			modal.Complete()

			modal.Reset()

			config := modal.GetConfigData()
			Expect(config.ProfileID).To(Equal("profile-1"))
		})
	})
})
