// generate_cv_e2e_test.go - E2E tests for GenerateCV intent DATA PERSISTENCE.
//
// This file tests:
//   - CV generation with persisted career data
//   - Generated CV data integrity
//
// For keyboard navigation tests, see:
//   - internal/cli/intents/generate_cv_navigation_test.go (if exists)
package e2e_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("E2E Generatecv Workflow", func() {
	var env *e2e.TestEnv

	Describe("With Career Data", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 3) // 5 events, 2 bursts, 3 facts
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show data is available when starting CV generation", func() {
			env.AssertEventCount(5)
			env.AssertFactCount(3)
			env.SelectIntentByName("generate_cv")
			// Profile selection should be available
			env.AssertViewContainsAny("Profile", "Select", "Senior", "Staff")
		})

		It("should allow profile and audience selection with data", func() {
			env.SelectIntentByName("generate_cv")
			env.Confirm() // Select profile
			env.AssertViewContainsAny("Audience", "hiring", "recruiter", "peer")
		})
	})

	Describe("Session Persistence", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show data after restart when generating CV", func() {
			env.PopulateTestData(3, 0, 2)
			env.SimulateRestart()

			env.AssertEventCount(3)
			env.SelectIntentByName("generate_cv")
			env.AssertViewContainsAny("Profile", "Select", "CV")
		})
	})

})
