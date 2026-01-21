package e2e_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("E2E Chained Workflows", func() {
	var env *e2e.TestEnv

	Describe("Browse After Data Population", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show events in Browse after populating data", func() {
			// Populate data (simulates capture)
			env.PopulateTestData(3, 0, 0)
			env.AssertEventCount(3)

			// Browse timeline
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("Acme Corp", "TechStart", "Events", "Timeline")
		})

		It("should show facts in FactManagement after populating data", func() {
			env.PopulateTestData(3, 0, 2)
			env.AssertFactCount(2)

			env.SelectIntentByName("fact_management")
			env.AssertViewContainsAny("Reduced", "Mentored", "Fact")
		})

		It("should show bursts in BurstManagement after populating data", func() {
			env.PopulateTestData(5, 2, 0)
			env.AssertBurstCount(2)

			env.SelectIntentByName("burst_management")
			env.AssertViewContainsAny("Authentication", "Mentoring", "Burst")
		})
	})

	Describe("Multi-Intent Navigation", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 3)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should navigate from Browse to Generate CV", func() {
			// Start in Browse
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("Timeline", "Events")

			// Return to menu
			env.Cancel()
			env.AssertViewContains("Capture Event")

			// Go to Generate CV
			env.SelectIntentByName("generate_cv")
			env.AssertViewContainsAny("Profile", "Select", "CV")
		})

		It("should navigate through multiple intents in sequence", func() {
			// Browse
			env.SelectIntentByName("browse_timeline")
			env.Cancel()

			// Generate CV
			env.SelectIntentByName("generate_cv")
			env.Cancel()

			// Configure
			env.SelectIntentByName("configure_system")
			env.Cancel()

			// Should be back at menu
			env.AssertViewContains("Capture Event")
		})

		It("should navigate from BurstManagement to FactManagement", func() {
			env.SelectIntentByName("burst_management")
			env.AssertViewContainsAny("Burst", "Authentication")

			env.Cancel()

			env.SelectIntentByName("fact_management")
			env.AssertViewContainsAny("Fact", "Reduced", "Mentored")
		})
	})

	Describe("Session Persistence Across Intents", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should persist data across intent navigation", func() {
			// Add data
			env.PopulateTestData(3, 0, 2)

			// Navigate to Browse and verify
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("Acme Corp", "TechStart")
			env.Cancel()

			// Navigate to FactManagement and verify same data
			env.SelectIntentByName("fact_management")
			env.AssertViewContainsAny("Reduced", "Mentored")
			env.Cancel()

			// Data should still be there
			env.AssertEventCount(3)
			env.AssertFactCount(2)
		})

		It("should persist data across simulated restart and intent navigation", func() {
			env.PopulateTestData(4, 0, 0)

			// Simulate restart
			env.SimulateRestart()

			// Data should persist
			env.AssertEventCount(4)

			// Navigate to browse
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("Acme Corp", "TechStart", "Events")
		})
	})

	Describe("Data Visibility Across Intents", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should see events in Browse", func() {
			env.PopulateTestData(3, 0, 0)

			// Check Browse
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("Events", "Acme Corp")
			env.Cancel()
		})

		It("should see facts in FactManagement and GenerateCV context", func() {
			env.PopulateTestData(3, 0, 3)
			env.AssertFactCount(3)

			// Check FactManagement
			env.SelectIntentByName("fact_management")
			env.AssertViewContainsAny("Reduced", "Mentored", "Fact")
			env.Cancel()

			// GenerateCV should be available (facts provide CV content)
			env.SelectIntentByName("generate_cv")
			env.AssertViewContainsAny("Profile", "Select", "CV")
		})
	})

	Describe("Workflow: Browse then Configure", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(3, 0, 0)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should allow browsing data then configuring system", func() {
			// Browse first
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("Timeline", "Events")
			env.Cancel()

			// Then configure
			env.SelectIntentByName("configure_system")
			env.AssertViewContainsAny("System", "Profile", "Domain")
			env.Cancel()

			// Data should still exist
			env.AssertEventCount(3)
		})
	})

	Describe("Return to Menu Between Intents", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should return to menu after each intent", func() {
			// Each intent should return to menu on cancel
			intents := []string{
				"browse_timeline",
				"generate_cv",
				"configure_system",
			}

			for _, intent := range intents {
				env.SelectIntentByName(intent)
				env.Cancel()
				env.AssertViewContains("Capture Event")
			}
		})
	})
})
