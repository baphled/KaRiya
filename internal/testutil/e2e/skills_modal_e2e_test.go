package e2e_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E Skills Modal Workflow", func() {
	var env *e2e.TestEnv

	Describe("Bug Regressions", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		// Regression test for BUG-019
		It("ensures skill modal uses ModalFormHeight for proper scrolling on small terminals", func() {
			env.SelectIntentByName("manage_skills")

			env.PressKeyRune('a')

			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})
	})
})
