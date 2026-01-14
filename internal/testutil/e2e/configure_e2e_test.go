package e2e_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("E2E Configure Workflow", func() {
	var env *e2e.TestEnv

	Describe("Session Persistence", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show configuration options after restart", func() {
			env.SimulateRestart()
			env.SelectIntentByName("configure_system")
			env.AssertViewContainsAny("System", "Profile", "Export", "UI", "Domain")
		})
	})

})
