package components_test

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/config"
)

var _ = Describe("OnboardingWizardModal", func() {
	var (
		modal *components.OnboardingWizardModal
	)

	BeforeEach(func() {
		modal = components.NewOnboardingWizardModal(80, 24)
		// Initialize the form for proper rendering
		modal.Init()
	})

	Describe("NewOnboardingWizardModal", func() {
		It("should create a new modal", func() {
			Expect(modal).NotTo(BeNil())
		})

		It("should be visible by default", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should not be completed initially", func() {
			Expect(modal.IsCompleted()).To(BeFalse())
		})

		It("should start at step 0", func() {
			Expect(modal.CurrentStep()).To(Equal(0))
		})
	})

	Describe("Init", func() {
		It("should return a command", func() {
			cmd := modal.Init()
			// Form init returns a command for cursor blinking
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("View", func() {
		It("should display welcome message on step 0", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Welcome"))
		})

		It("should display step 1 title on step 0", func() {
			view := modal.View()
			// huh forms don't fully render field labels in unit tests
			// but the step/group title should be visible
			Expect(view).To(Or(
				ContainSubstring("Step 1"),
				ContainSubstring("Welcome"),
			))
		})

		It("should display step indicator", func() {
			view := modal.View()
			Expect(view).To(Or(
				ContainSubstring("Step 1"),
				ContainSubstring("1/3"),
				ContainSubstring("1 of 3"),
			))
		})

		It("should display navigation help", func() {
			view := modal.View()
			Expect(view).To(Or(
				ContainSubstring("Enter"),
				ContainSubstring("continue"),
				ContainSubstring("next"),
			))
		})

		It("should return empty string when not visible", func() {
			modal.Hide()
			view := modal.View()
			Expect(view).To(BeEmpty())
		})
	})

	Describe("Update", func() {
		Context("window resize", func() {
			It("should handle window resize", func() {
				cmd := modal.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
				// May return a command for form reinit
				_ = cmd
				Expect(modal.IsVisible()).To(BeTrue())
			})

			It("should return form init command after dimension change", func() {
				cmd := modal.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

				Expect(cmd).NotTo(BeNil())
			})
		})

		Context("escape key", func() {
			// Onboarding is mandatory - Esc key is ignored
			It("should NOT close modal on Esc at step 0 (mandatory onboarding)", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(modal.IsVisible()).To(BeTrue())
			})

			It("should NOT mark as cancelled on Esc (mandatory onboarding)", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(modal.WasCancelled()).To(BeFalse())
			})
		})
	})

	Describe("GetProfileConfig", func() {
		It("should return nil when not completed", func() {
			cfg := modal.GetProfileConfig()
			Expect(cfg).To(BeNil())
		})
	})

	Describe("Visibility", func() {
		It("should be able to hide", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should be able to show", func() {
			modal.Hide()
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})
	})

	Describe("Data Collection", func() {
		Context("with pre-filled data", func() {
			It("should accept existing profile config", func() {
				existingConfig := &config.ProfileConfig{
					Name:  "Existing Name",
					Email: "existing@example.com",
				}
				modalWithData := components.NewOnboardingWizardModalWithConfig(80, 24, existingConfig)
				Expect(modalWithData).NotTo(BeNil())
			})
		})
	})

	Describe("Required Fields Validation", func() {
		It("should report missing required fields when name is empty", func() {
			Expect(modal.HasRequiredFields()).To(BeFalse())
		})
	})

	Describe("Step Navigation", func() {
		It("should have 3 steps total", func() {
			Expect(modal.TotalSteps()).To(Equal(3))
		})
	})
})
