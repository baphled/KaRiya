package e2e_test

import (
	"os"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// NOTE: These tests need to be rewritten to use bootstrap.OnboardingTestModel
// since onboarding is now a separate pre-app phase that runs before the main app.
// The main app no longer has an "onboarding state" - it starts directly in menu state.
//
// TODO: Create new test infrastructure for testing the onboarding Bubble Tea program.
// See: internal/cli/bootstrap/testing.go for OnboardingTestModel

var _ = Describe("E2E Onboarding Wizard Workflow", func() {
	var env *e2e.TestEnv

	Describe("Onboarding Initialization", Pending, func() {
		// These tests need to be rewritten using bootstrap.NewOnboardingTestModel()
		BeforeEach(func() {
			env = e2e.GetSharedEnvWithOnboarding(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show onboarding wizard on fresh startup", func() {
			// TODO: Test using OnboardingTestModel directly
			Skip("Onboarding is now a separate pre-app phase - test needs rewrite")
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

	Describe("Onboarding Step Navigation", Pending, func() {
		// These tests need to be rewritten using bootstrap.NewOnboardingTestModel()
		BeforeEach(func() {
			env = e2e.GetSharedEnvWithOnboarding(GinkgoT())
			env.InitModel()
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should remain on Step 1 without entering name", func() {
			env.PressEnterWithFormProcessing()
			env.AssertViewContains("Step 1 of 3")
		})

		It("should accept typed name text", func() {
			env.TypeText("Test User")
			env.AssertViewContains("Test User")
		})

		It("should advance to Step 2 after entering valid name", func() {
			env.TypeText("Test User")
			env.PressEnterWithFormProcessing()
			env.AssertViewContains("Step 2 of 3")
			env.AssertViewContains("Email")
		})

		It("should show Location field on Step 2", func() {
			env.TypeText("Test User")
			env.PressEnterWithFormProcessing()
			env.AssertViewContains("Location")
		})

		It("should advance to Step 3 after entering valid email", func() {
			env.TypeText("Test User")
			env.PressEnterWithFormProcessing()
			env.TypeText("test@example.com")
			env.PressEnterWithFormProcessing()
			env.PressEnterWithFormProcessing()
			env.AssertViewContains("Step 3 of 3")
		})

		It("should show Professional Details on Step 3", func() {
			env.TypeText("Test User")
			env.PressEnterWithFormProcessing()
			env.TypeText("test@example.com")
			env.PressEnterWithFormProcessing()
			env.PressEnterWithFormProcessing()
			env.AssertViewContains("Professional")
			env.AssertViewContainsAny("Title", "GitHub", "Portfolio")
		})
	})

	Describe("Onboarding Completion", Pending, func() {
		// These tests need to be rewritten using bootstrap.NewOnboardingTestModel()
		BeforeEach(func() {
			env = e2e.GetSharedEnvWithOnboarding(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should complete onboarding and return profile config", func() {
			// TODO: Test using OnboardingTestModel directly
			// Should verify that completing onboarding returns a valid ProfileConfig
			Skip("Onboarding is now a separate pre-app phase - test needs rewrite")
		})
	})

	Describe("Onboarding Escape Key Behavior", Pending, func() {
		// These tests need to be rewritten using bootstrap.NewOnboardingTestModel()
		BeforeEach(func() {
			env = e2e.GetSharedEnvWithOnboarding(GinkgoT())
			env.InitModel()
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should NOT cancel onboarding on Esc key (mandatory wizard)", func() {
			env.Cancel()
			// Onboarding is mandatory - should still show wizard
			env.AssertViewContains("Profile Setup")
		})

		It("should NOT allow skipping after entering partial data", func() {
			env.TypeText("Test User")
			env.Cancel()
			// Should still be in wizard
			env.AssertViewContains("Profile Setup")
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
			// Regular E2E Setup should skip onboarding (via bootstrap.SkipOnboarding)
			// The app starts directly in menu state
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

			// Now set up the test environment
			env = e2e.SetupWithOnboarding(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("BUG-007: config file isolation should work in test setup", func() {
			// Verify we're using an isolated config path
			testConfigPath, err := config.GetConfigPath()
			Expect(err).NotTo(HaveOccurred())
			Expect(testConfigPath).NotTo(Equal(realConfigPath),
				"Test should use isolated config path, not real config")
		})

		It("BUG-007: should NOT write to user's real config file", func() {
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
	})
})
