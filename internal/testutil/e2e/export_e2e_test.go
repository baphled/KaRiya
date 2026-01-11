package e2e_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("E2E Export Workflow", func() {
	var env *e2e.TestEnv

	Describe("With Career Data", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 3) // 5 events, 2 bursts, 3 facts
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should have data available for export", func() {
			env.AssertEventCount(5)
			env.AssertBurstCount(2)
			env.AssertFactCount(3)
		})

		It("should allow navigating through export workflow with data", func() {
			env.SelectIntentByName("export_artifact")
			env.AssertViewContainsAny("events", "facts", "bursts", "Events", "Facts", "Bursts")
			env.Confirm() // Select type
			env.AssertViewContainsAny("Format", "JSON", "CSV")
			env.Confirm() // Select format
			env.AssertViewContainsAny("Destination", "file", "clipboard")
		})
	})

	Describe("Session Persistence", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should have data available after restart", func() {
			// Use 5 events to ensure enough for 2 bursts (each needs at least 2 events)
			env.PopulateTestData(5, 2, 2)
			env.SimulateRestart()

			env.AssertEventCount(5)
			env.AssertBurstCount(2)
			env.AssertFactCount(2)

			env.SelectIntentByName("export_artifact")
			env.AssertViewContainsAny("events", "facts", "bursts", "Type")
		})
	})

})
