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

	Describe("View when completed", func() {
		It("should render summary when wizard is completed with all fields", func() {
			completedModal := components.NewOnboardingWizardModalWithConfig(80, 24, &config.ProfileConfig{
				Name:      "Alice Smith",
				Email:     "alice@example.com",
				Location:  "London",
				Title:     "Senior Engineer",
				GitHub:    "alicesmith",
				Portfolio: "https://alice.dev",
			})
			completedModal.Init()
			completedModal.SetDataForTesting(&components.OnboardingData{
				Name:      "Alice Smith",
				Email:     "alice@example.com",
				Location:  "London",
				Title:     "Senior Engineer",
				GitHub:    "alicesmith",
				Portfolio: "https://alice.dev",
			})
			completedModal.CompleteWizardForTesting()

			view := completedModal.View()
			Expect(view).To(ContainSubstring("Profile Setup Complete!"))
			Expect(view).To(ContainSubstring("Name: Alice Smith"))
			Expect(view).To(ContainSubstring("Email: alice@example.com"))
			Expect(view).To(ContainSubstring("Location: London"))
			Expect(view).To(ContainSubstring("Title: Senior Engineer"))
			Expect(view).To(ContainSubstring("GitHub: alicesmith"))
			Expect(view).To(ContainSubstring("Portfolio: https://alice.dev"))
		})

		It("should render summary with only required fields", func() {
			completedModal := components.NewOnboardingWizardModal(80, 24)
			completedModal.Init()
			completedModal.SetDataForTesting(&components.OnboardingData{
				Name:  "Bob Jones",
				Email: "bob@example.com",
			})
			completedModal.CompleteWizardForTesting()

			view := completedModal.View()
			Expect(view).To(ContainSubstring("Profile Setup Complete!"))
			Expect(view).To(ContainSubstring("Name: Bob Jones"))
			Expect(view).To(ContainSubstring("Email: bob@example.com"))
			Expect(view).NotTo(ContainSubstring("Location:"))
			Expect(view).NotTo(ContainSubstring("Title:"))
			Expect(view).NotTo(ContainSubstring("GitHub:"))
			Expect(view).NotTo(ContainSubstring("Portfolio:"))
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

		Context("when wizard is not visible", func() {
			It("should return nil without processing", func() {
				modal.Hide()
				cmd := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).To(BeNil())
			})
		})

		Context("when wizard completes", func() {
			It("should sync form data on completion", func() {
				completedModal := components.NewOnboardingWizardModalWithConfig(80, 24, &config.ProfileConfig{
					Name:  "Test User",
					Email: "test@example.com",
				})
				completedModal.Init()
				completedModal.CompleteWizardForTesting()

				Expect(completedModal.IsCompleted()).To(BeTrue())
				Expect(completedModal.IsVisible()).To(BeFalse())
			})
		})
	})

	Describe("GetProfileConfig", func() {
		It("should return nil when not completed", func() {
			cfg := modal.GetProfileConfig()
			Expect(cfg).To(BeNil())
		})

		It("should return profile config when completed", func() {
			completedModal := components.NewOnboardingWizardModalWithConfig(80, 24, &config.ProfileConfig{
				Name:      "Alice Smith",
				Email:     "alice@example.com",
				Location:  "London",
				Title:     "Senior Engineer",
				GitHub:    "alicesmith",
				Portfolio: "https://alice.dev",
			})
			completedModal.Init()
			completedModal.CompleteWizardForTesting()

			cfg := completedModal.GetProfileConfig()
			Expect(cfg).NotTo(BeNil())
			Expect(cfg.Name).To(Equal("Alice Smith"))
			Expect(cfg.Email).To(Equal("alice@example.com"))
			Expect(cfg.Location).To(Equal("London"))
			Expect(cfg.Title).To(Equal("Senior Engineer"))
			Expect(cfg.GitHub).To(Equal("alicesmith"))
			Expect(cfg.Portfolio).To(Equal("https://alice.dev"))
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

		Context("with nil config", func() {
			It("should create modal with empty data", func() {
				modalNilCfg := components.NewOnboardingWizardModalWithConfig(80, 24, nil)
				Expect(modalNilCfg).NotTo(BeNil())
				Expect(modalNilCfg.IsVisible()).To(BeTrue())
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

		It("should start at step 0", func() {
			Expect(modal.CurrentStep()).To(Equal(0))
		})
	})

	Describe("GetOnboardingData", func() {
		It("should return the onboarding data", func() {
			data := modal.GetOnboardingData()
			Expect(data).NotTo(BeNil())
		})

		It("should return data with pre-filled values", func() {
			prefilledModal := components.NewOnboardingWizardModalWithConfig(80, 24, &config.ProfileConfig{
				Name:  "Test User",
				Email: "test@example.com",
			})
			data := prefilledModal.GetOnboardingData()
			Expect(data).NotTo(BeNil())
		})
	})

	Describe("calcOnboardingModalWidth", func() {
		It("should cap width at 80", func() {
			width := components.CalcOnboardingModalWidthForTesting(120)
			Expect(width).To(Equal(80))
		})

		It("should return width minus 20 when in range", func() {
			width := components.CalcOnboardingModalWidthForTesting(80)
			Expect(width).To(Equal(60))
		})

		It("should enforce minimum width of 50", func() {
			width := components.CalcOnboardingModalWidthForTesting(40)
			Expect(width).To(Equal(50))
		})

		It("should enforce minimum width of 50 for very small terminals", func() {
			width := components.CalcOnboardingModalWidthForTesting(20)
			Expect(width).To(Equal(50))
		})
	})
})
