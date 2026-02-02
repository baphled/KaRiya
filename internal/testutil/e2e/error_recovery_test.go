package e2e_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E Error Recovery", func() {
	var env *e2e.TestEnv

	Describe("Empty State Handling", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should handle Browse with no events gracefully", func() {
			env.AssertEventCount(0)
			env.SelectIntentByName("browse_timeline")
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle FactManagement with no facts gracefully", func() {
			env.AssertFactCount(0)
			env.SelectIntentByName("fact_management")
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle BurstManagement with no bursts gracefully", func() {
			env.AssertBurstCount(0)
			env.SelectIntentByName("burst_management")
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle GenerateCV with no facts gracefully", func() {
			env.SelectIntentByName("generate_cv")
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Boundary Navigation", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(3, 0, 3)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should handle navigating past list start in Browse", func() {
			env.SelectIntentByName("browse_timeline")
			// Press up multiple times at the start
			for i := range 5 {
				_ = i
				env.PressKey(tea.KeyUp)
			}
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle navigating past list end in Browse", func() {
			env.SelectIntentByName("browse_timeline")
			for range 10 {
				env.PressKey(tea.KeyDown)
			}
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle vim navigation past boundaries in FactManagement", func() {
			env.SelectIntentByName("fact_management")
			for range 5 {
				env.PressKeyRune('k')
			}
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))

			for range 10 {
				env.PressKeyRune('j')
			}
			view = env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})
	})

	Describe("Rapid Key Presses", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(5, 2, 3)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should handle rapid navigation in Browse", func() {
			env.SelectIntentByName("browse_timeline")
			for range 20 {
				env.PressKeyRune('j')
				env.PressKeyRune('k')
			}
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle rapid Escape presses", func() {
			env.SelectIntentByName("browse_timeline")
			for range 5 {
				env.Cancel()
			}
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			env.AssertViewContains("Capture Event")
		})

		It("should handle rapid intent switching", func() {
			intents := []string{
				"browse_timeline",
				"generate_cv",
				"configure_system",
			}
			for _, intent := range intents {
				env.SelectIntentByName(intent)
				env.Cancel()
			}
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Cancel Recovery", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should recover to menu after cancelling Browse", func() {
			env.SelectIntentByName("browse_timeline")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

		It("should recover to menu after cancelling GenerateCV", func() {
			env.SelectIntentByName("generate_cv")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

		It("should recover to menu after cancelling Configure", func() {
			env.SelectIntentByName("configure_system")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

		It("should recover to menu after cancelling FactManagement", func() {
			env.SelectIntentByName("fact_management")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

		It("should recover to menu after cancelling BurstManagement", func() {
			env.SelectIntentByName("burst_management")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Quit Key Recovery", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should handle quit from main menu", func() {
			env.Quit()
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should handle quit from Browse", func() {
			env.SelectIntentByName("browse_timeline")
			env.Quit()
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should handle quit from FactManagement", func() {
			env.SelectIntentByName("fact_management")
			env.Quit()
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})
	})

	Describe("State Consistency After Errors", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should maintain data after navigation errors", func() {
			env.PopulateTestData(3, 0, 2)
			env.AssertEventCount(3)
			env.AssertFactCount(2)

			env.SelectIntentByName("browse_timeline")
			for range 20 {
				env.PressKey(tea.KeyUp)
			}
			env.Cancel()

			// Data should still be consistent
			env.AssertEventCount(3)
			env.AssertFactCount(2)
		})

		It("should maintain data after rapid cancel", func() {
			env.PopulateTestData(5, 2, 3)

			for range 5 {
				env.SelectIntentByName("browse_timeline")
				env.Cancel()
			}

			// Data should still be consistent
			env.AssertEventCount(5)
			env.AssertBurstCount(2)
			env.AssertFactCount(3)
		})
	})

	Describe("Restart Recovery", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should recover state after restart", func() {
			env.PopulateTestData(3, 0, 2)

			env.SimulateRestart()

			env.AssertEventCount(3)
			env.AssertFactCount(2)
		})

		It("should allow navigation after restart", func() {
			env.PopulateTestData(3, 0, 0)

			env.SimulateRestart()

			env.SelectIntentByName("browse_timeline")
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should maintain menu state after restart", func() {
			env.SimulateRestart()

			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Multiple Cancel in Nested Views", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(5, 2, 3)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should handle cancel from GenerateCV audience selection", func() {
			env.SelectIntentByName("generate_cv")
			env.Confirm() // Select profile
			env.Cancel()  // Back
			env.Cancel()  // Back to menu
			env.AssertViewContains("Capture Event")
		})

		It("should handle cancel from Configure settings", func() {
			env.SelectIntentByName("configure_system")
			env.Confirm() // Select category
			env.Cancel()  // Back
			env.Cancel()  // Back to menu
			env.AssertViewContains("Capture Event")
		})
	})

})
