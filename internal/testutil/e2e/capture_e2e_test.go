package e2e_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("E2E Capture Workflow", func() {
	var env *e2e.TestEnv

	Describe("Database Persistence", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should start with empty database", func() {
			env.AssertEventCount(0)
		})

		It("should not persist event when cancelled before submission", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm() // Select Quick strategy
			env.TypeText("Test event text")
			env.Cancel() // Cancel from form
			env.AssertEventCount(0)
		})

		It("should not persist event when returning to main menu", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm() // Select Quick strategy
			env.TypeText("Test event text")
			env.PressKeyRune('m') // Return to main menu
			env.AssertEventCount(0)
		})
	})

	Describe("Session Persistence", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show empty timeline after restart with no events", func() {
			env.SimulateRestart()
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("No events", "empty", "Timeline")
		})
	})

})
