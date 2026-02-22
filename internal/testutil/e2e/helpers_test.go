package e2e_test

import (
	"github.com/baphled/kariya/internal/cli/intents/captureevent"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E Test Helpers", func() {
	Describe("Setup", func() {
		It("should create a complete test environment with SQLite", func() {
			env := e2e.Setup(GinkgoT())
			defer env.Cleanup()

			Expect(env.Model).NotTo(BeNil())
			Expect(env.DB).NotTo(BeNil())
			Expect(env.Service).NotTo(BeNil())
		})
	})

	Describe("SetupWithMemory", func() {
		It("should create a test environment with in-memory repositories", func() {
			env := e2e.SetupWithMemory(GinkgoT())
			defer env.Cleanup()

			Expect(env.Model).NotTo(BeNil())
			Expect(env.DB).To(BeNil())
			Expect(env.Service).NotTo(BeNil())
		})
	})

	Describe("Navigation Helpers", func() {
		var env *e2e.TestEnv

		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should start in menu state", func() {
			Expect(env.IsInMenuState()).To(BeTrue())
		})

		It("should remain in menu after navigation", func() {
			env.NavigateDown().NavigateDown()
			Expect(env.IsInMenuState()).To(BeTrue())
		})
	})

	Describe("SelectIntent", func() {
		var env *e2e.TestEnv

		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should select intent by index", func() {
			env.SelectIntent(0)
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("SelectIntentByName", func() {
		var env *e2e.TestEnv

		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should select intent by name", func() {
			env.SelectIntentByName("browse_timeline")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("View Assertions", func() {
		var env *e2e.TestEnv

		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should assert view contains expected text", func() {
			env.AssertViewContains("Career Event Management System")
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Data Population", func() {
		var env *e2e.TestEnv

		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should start with empty database", func() {
			env.AssertEventCount(0)
			env.AssertBurstCount(0)
			env.AssertFactCount(0)
		})

		It("should populate test data", func() {
			env.PopulateTestData(5, 2, 3)
			env.AssertEventCount(5)
			env.AssertBurstCount(2)
			env.AssertFactCount(3)
		})
	})

	Describe("SimulateRestart", func() {
		var env *e2e.TestEnv

		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should preserve data after restart", func() {
			event := e2e.CreateMinimalEvent("test_event_1")
			env.AddEvent(event)
			env.AssertEventCount(1)

			env.SimulateRestart()

			env.AssertEventCount(1)
			Expect(env.IsInMenuState()).To(BeTrue())
		})
	})

	Describe("Fixture Generators", func() {
		Describe("CreateSampleEvents", func() {
			It("should create valid sample events", func() {
				events := e2e.CreateSampleEvents(10)

				Expect(events).To(HaveLen(10))
				for _, event := range events {
					Expect(event.ID).NotTo(BeEmpty())
					Expect(event.Text).NotTo(BeEmpty())
					Expect(event.Date.IsZero()).To(BeFalse())
				}
			})
		})

		Describe("CreateSampleBursts", func() {
			It("should create valid sample bursts with event references", func() {
				events := e2e.CreateSampleEvents(5)
				bursts := e2e.CreateSampleBursts(3, events)

				Expect(bursts).To(HaveLen(3))
				for _, burst := range bursts {
					Expect(burst.ID).NotTo(BeEmpty())
					Expect(burst.Name).NotTo(BeEmpty())
					Expect(burst.EventIDs).NotTo(BeEmpty())
				}
			})
		})

		Describe("CreateSampleFacts", func() {
			It("should create valid sample facts with event references", func() {
				events := e2e.CreateSampleEvents(5)
				facts := e2e.CreateSampleFacts(5, events)

				Expect(facts).To(HaveLen(5))
				for _, fact := range facts {
					Expect(fact.ID).NotTo(BeEmpty())
					Expect(fact.Text).NotTo(BeEmpty())
					Expect(fact.CompetencyCategories).NotTo(BeEmpty())
				}
			})
		})

		Describe("CreateSampleProfiles", func() {
			It("should create valid sample profiles", func() {
				profiles := e2e.CreateSampleProfiles()

				Expect(len(profiles)).To(BeNumerically(">=", 3))
				for _, profile := range profiles {
					Expect(profile.ID).NotTo(BeEmpty())
					Expect(profile.Name).NotTo(BeEmpty())
					Expect(profile.TargetRole).NotTo(BeEmpty())
				}
			})
		})
	})

	Describe("processCmdResult PostSavePersistenceCompleteMsg", func() {
		var env *e2e.TestEnv

		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("completes the capture event intent when received", func() {
			env.SelectIntentByName("capture_event")

			event := e2e.CreateMinimalEvent("evt-post-save")
			env.SendMessage(captureevent.PostSavePersistenceCompleteMsg{
				Event:  event,
				Bursts: []*career.Burst{},
				Facts:  []*career.Fact{},
				Skills: []*career.Skill{},
			})

			view := env.GetView()
			Expect(view).To(ContainSubstring("Capture Event"))
		})
	})
})
