package e2e_test

import (
	"os"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E Onboarding Wizard Workflow", func() {
	var env *e2e.TestEnv

	Describe("Onboarding Initialization", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnvWithOnboarding(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show onboarding wizard on fresh startup", func() {
			Expect(env.IsInOnboardingState()).To(BeTrue(), "Should be in onboarding state")
		})

		It("should show Profile Setup title", func() {
			env.AssertViewContains("Profile Setup")
		})

		It("should show Step 1 of 3", func() {
			env.AssertViewContains("Step 1 of 3")
		})

		It("should show Welcome message", func() {
			env.InitModel()
			env.AssertViewContains("Welcome to KaRiya")
		})

		It("should show Name field", func() {
			env.InitModel()
			env.AssertViewContains("Your Name")
		})

		It("should show keyboard shortcuts", func() {
			env.AssertViewContainsAny("tab", "enter", "next")
		})
	})

	Describe("Onboarding Step Navigation", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnvWithOnboarding(GinkgoT())
			// Initialize the model to set up huh forms properly
			env.InitModel()
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should remain on Step 1 without entering name", func() {
			// Try to advance without entering name - validation should block
			env.PressEnterWithFormProcessing()
			// Should still show step 1 due to validation
			env.AssertViewContains("Step 1 of 3")
		})

		It("should accept typed name text", func() {
			env.TypeText("Test User")
			// The input should show the typed text
			env.AssertViewContains("Test User")
		})

		It("should advance to Step 2 after entering valid name", func() {
			env.TypeText("Test User")
			env.PressEnterWithFormProcessing()
			// Should now show Step 2
			env.AssertViewContains("Step 2 of 3")
			env.AssertViewContains("Email")
		})

		It("should show Location field on Step 2", func() {
			env.TypeText("Test User")
			env.PressEnterWithFormProcessing()
			env.AssertViewContains("Location")
		})

		It("should advance to Step 3 after entering valid email", func() {
			// Step 1: Name
			env.TypeText("Test User")
			env.PressEnterWithFormProcessing()
			// Step 2: Email (required) and Location (optional)
			env.TypeText("test@example.com")
			env.PressEnterWithFormProcessing() // Accept email, move to Location
			env.PressEnterWithFormProcessing() // Accept empty Location, advance to Step 3
			// Should now show Step 3
			env.AssertViewContains("Step 3 of 3")
		})

		It("should show Professional Details on Step 3", func() {
			// Navigate to Step 3
			env.TypeText("Test User")
			env.PressEnterWithFormProcessing()
			// Step 2: Email and Location
			env.TypeText("test@example.com")
			env.PressEnterWithFormProcessing() // Accept email, move to Location
			env.PressEnterWithFormProcessing() // Accept Location, advance to Step 3
			// Check Step 3 content
			env.AssertViewContains("Professional")
			env.AssertViewContainsAny("Title", "GitHub", "Portfolio")
		})
	})

	Describe("Onboarding Completion", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnvWithOnboarding(GinkgoT())
			// InitModel is called by CompleteOnboarding
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should complete onboarding and show main menu", func() {
			// Complete all steps
			env.CompleteOnboarding("Test User", "test@example.com")

			// Should now be at main menu
			Expect(env.IsInOnboardingState()).To(BeFalse(), "Should not be in onboarding state after completion")
			Expect(env.IsInMenuState()).To(BeTrue(), "Should be in menu state after completion")
		})

		It("should show menu items after completion", func() {
			env.CompleteOnboarding("Test User", "test@example.com")
			env.AssertViewContains("Capture Event")
			env.AssertViewContains("Browse Timeline")
		})
	})

	Describe("Onboarding Escape Key Behavior", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnvWithOnboarding(GinkgoT())
			// Initialize the model to set up huh forms properly
			env.InitModel()
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should NOT cancel onboarding on Esc key (mandatory wizard)", func() {
			// Press Esc
			env.Cancel()
			// Should still be in onboarding
			Expect(env.IsInOnboardingState()).To(BeTrue(), "Onboarding should not be cancelled with Esc")
			env.AssertViewContains("Profile Setup")
		})

		It("should NOT allow skipping after entering partial data", func() {
			env.TypeText("Test User")
			env.Cancel()
			// Should still be in onboarding
			Expect(env.IsInOnboardingState()).To(BeTrue())
		})
	})

	Describe("Onboarding with Existing E2E Setup", func() {
		BeforeEach(func() {
			// Use regular Setup (which skips onboarding)
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should skip onboarding with standard Setup", func() {
			// Regular E2E Setup should skip onboarding
			Expect(env.IsInOnboardingState()).To(BeFalse(), "Regular Setup should skip onboarding")
		})

		It("should go directly to menu with standard Setup", func() {
			// Should be at menu immediately
			Expect(env.IsInMenuState()).To(BeTrue(), "Should start at menu with standard Setup")
		})
	})

	Describe("BUG-007 Regression: Config File Isolation", func() {
		var (
			realConfigPath  string
			originalContent []byte
			originalExists  bool
		)

		BeforeEach(func() {
			// Get the real config path BEFORE any test setup
			var err error
			realConfigPath, err = config.GetConfigPath()
			Expect(err).NotTo(HaveOccurred())

			// Read original content if it exists
			originalContent, err = os.ReadFile(realConfigPath)
			if err == nil {
				originalExists = true
			} else if os.IsNotExist(err) {
				originalExists = false
			} else {
				Fail("Failed to read original config: " + err.Error())
			}

			// Now set up the test environment with onboarding
			env = e2e.SetupWithOnboarding(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("BUG-007: should NOT write to user's real config file during onboarding", func() {
			// Complete onboarding with test values
			env.CompleteOnboarding("BUG007 Test User", "bug007@example.com")

			// Verify onboarding completed
			Expect(env.IsInOnboardingState()).To(BeFalse(), "Onboarding should be complete")
			Expect(env.IsInMenuState()).To(BeTrue(), "Should be at menu after completion")

			// Now verify the real config file was NOT modified
			if originalExists {
				currentContent, err := os.ReadFile(realConfigPath)
				Expect(err).NotTo(HaveOccurred(), "Should be able to read real config")
				Expect(currentContent).To(Equal(originalContent),
					"Real config file should NOT have been modified by test")
			} else {
				// If config didn't exist before, it should still not exist
				_, err := os.Stat(realConfigPath)
				Expect(os.IsNotExist(err)).To(BeTrue(),
					"Real config file should NOT have been created by test")
			}
		})

		It("BUG-007: should NOT contain test values in real config after onboarding", func() {
			// Complete onboarding with distinctive test values
			env.CompleteOnboarding("BUG007 Unique Name", "bug007unique@test.com")

			// Read the real config file
			if originalExists {
				currentContent, err := os.ReadFile(realConfigPath)
				Expect(err).NotTo(HaveOccurred())

				contentStr := string(currentContent)
				Expect(contentStr).NotTo(ContainSubstring("BUG007 Unique Name"),
					"Real config should NOT contain test name")
				Expect(contentStr).NotTo(ContainSubstring("bug007unique@test.com"),
					"Real config should NOT contain test email")
			}
		})
	})
})
