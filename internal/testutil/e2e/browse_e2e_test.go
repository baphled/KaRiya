package e2e_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E Browse Workflow", func() {
	var env *e2e.TestEnv

	Describe("Empty Timeline Data State", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT()) // SQLite for persistence testing
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should verify database has no events initially", func() {
			env.AssertEventCount(0)
		})
	})

	Describe("Timeline with Persisted Events", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 0, 0) // 5 events, no bursts, no facts
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should display events from database in timeline", func() {
			env.SelectIntentByName("browse_timeline")
			env.AssertEventCount(5)
			// Should show at least some event content (company name from fixtures)
			env.AssertViewContainsAny("Acme Corp", "TechStart", "BigCorp", "Events", "Page")
		})

		It("should show event data from fixtures", func() {
			env.SelectIntentByName("browse_timeline")
			env.Confirm()
			// Events from fixtures have company and text
			env.AssertViewContainsAny("Acme Corp", "TechStart", "security", "authentication", "Date")
		})
	})

	Describe("Session Persistence", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show events after restart", func() {
			// Add event before restart
			env.PopulateTestData(3, 0, 0)
			env.AssertEventCount(3)

			// Simulate restart
			env.SimulateRestart()

			// Events should still exist
			env.AssertEventCount(3)
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("Acme Corp", "TechStart", "Events")
		})

		It("should maintain event data integrity after restart", func() {
			env.PopulateTestData(2, 0, 0)
			eventsBefore := env.GetEvents()
			Expect(len(eventsBefore)).To(Equal(2))

			env.SimulateRestart()

			eventsAfter := env.GetEvents()
			Expect(len(eventsAfter)).To(Equal(2))
			// Check that event IDs match
			Expect(eventsAfter[0].ID).To(Equal(eventsBefore[0].ID))
		})
	})

	Describe("Workflow Data Integration", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(3, 0, 0)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		// NOTE: Legacy workflow test removed - BrowseTimeline now uses screen architecture (Phase 4.2)
		// See internal/cli/screens/timeline/ for current tests
	})
})
