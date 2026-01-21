package e2e_test

import (
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E Burst Management Workflow", func() {
	var env *e2e.TestEnv

	Describe("Empty Burst List Data State", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should verify database has no bursts initially", func() {
			env.AssertBurstCount(0)
		})
	})

	Describe("Burst List with Persisted Data", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(5, 2, 0) // 5 events, 2 bursts, 0 facts
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should display bursts from database in list", func() {
			env.SelectIntentByName("burst_management")
			env.AssertBurstCount(2)
			env.AssertViewContainsAny("Authentication", "Mentoring", "Burst", "Name")
		})

		It("should show burst data from fixtures", func() {
			env.SelectIntentByName("burst_management")
			env.Confirm()
			env.AssertViewContainsAny("Authentication", "Mentoring", "Description")
		})
	})

	Describe("Session Persistence", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show bursts after restart", func() {
			env.PopulateTestData(5, 2, 0)
			env.AssertBurstCount(2)

			env.SimulateRestart()

			env.AssertBurstCount(2)
			env.SelectIntentByName("burst_management")
			env.AssertViewContainsAny("Authentication", "Mentoring", "Burst")
		})

		It("should maintain burst data integrity after restart", func() {
			env.PopulateTestData(5, 2, 0)
			burstsBefore := env.GetBursts()
			Expect(len(burstsBefore)).To(Equal(2))

			env.SimulateRestart()

			burstsAfter := env.GetBursts()
			Expect(len(burstsAfter)).To(Equal(2))
			Expect(burstsAfter[0].ID).To(Equal(burstsBefore[0].ID))
		})
	})

	Describe("Burst Edit Data Persistence", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(5, 2, 0)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should persist burst name changes to database", func() {
			originalBursts := env.GetBursts()
			Expect(len(originalBursts)).To(BeNumerically(">=", 1))
			originalName := originalBursts[0].Name
			burstID := originalBursts[0].ID

			env.SelectIntentByName("burst_management")
			env.Confirm()
			env.PressKeyRune('e')
			env.AssertViewContainsAny("Edit Burst", "Burst Name", "Name")

			activeIntent := env.Model.GetActiveIntent()
			Expect(activeIntent).NotTo(BeNil())

			burstIntent, ok := activeIntent.(*intents.BurstManagementIntent)
			Expect(ok).To(BeTrue(), "Active intent should be BurstManagementIntent")

			newName := "UPDATED_BURST_NAME_E2E_TEST"
			modifiedBurst := &career.Burst{
				ID:          burstID,
				Name:        newName,
				Description: originalBursts[0].Description,
				EventIDs:    originalBursts[0].EventIDs,
				CreatedAt:   originalBursts[0].CreatedAt,
				UpdatedAt:   originalBursts[0].UpdatedAt,
			}

			burstIntent.SetTestModalResult(&intents.ModalEditResult[*career.Burst]{
				Original: originalBursts[0],
				Modified: modifiedBurst,
				Accepted: true,
				Changes:  map[string]interface{}{"name": newName},
			})

			env.Model.Update(nil)

			updatedBursts := env.GetBursts()
			var foundBurst *career.Burst
			for _, b := range updatedBursts {
				if b.ID == burstID {
					foundBurst = b
					break
				}
			}

			Expect(foundBurst).NotTo(BeNil(), "Should find the updated burst in database")
			Expect(foundBurst.Name).To(Equal(newName), "Burst name should be updated in database")
			Expect(foundBurst.Name).NotTo(Equal(originalName), "Burst name should differ from original")
		})

		It("should persist burst description changes to database", func() {
			originalBursts := env.GetBursts()
			Expect(len(originalBursts)).To(BeNumerically(">=", 1))
			originalDesc := originalBursts[0].Description
			burstID := originalBursts[0].ID

			env.SelectIntentByName("burst_management")
			env.Confirm()
			env.PressKeyRune('e')

			activeIntent := env.Model.GetActiveIntent()
			burstIntent := activeIntent.(*intents.BurstManagementIntent)

			newDesc := "UPDATED_DESCRIPTION_E2E_TEST"
			modifiedBurst := &career.Burst{
				ID:          burstID,
				Name:        originalBursts[0].Name,
				Description: newDesc,
				EventIDs:    originalBursts[0].EventIDs,
				CreatedAt:   originalBursts[0].CreatedAt,
				UpdatedAt:   originalBursts[0].UpdatedAt,
			}

			burstIntent.SetTestModalResult(&intents.ModalEditResult[*career.Burst]{
				Original: originalBursts[0],
				Modified: modifiedBurst,
				Accepted: true,
				Changes:  map[string]interface{}{"description": newDesc},
			})

			env.Model.Update(nil)

			updatedBursts := env.GetBursts()
			var foundBurst *career.Burst
			for _, b := range updatedBursts {
				if b.ID == burstID {
					foundBurst = b
					break
				}
			}

			Expect(foundBurst).NotTo(BeNil())
			Expect(foundBurst.Description).To(Equal(newDesc))
			Expect(foundBurst.Description).NotTo(Equal(originalDesc))
		})

		It("should NOT persist changes when edit is cancelled", func() {
			originalBursts := env.GetBursts()
			Expect(len(originalBursts)).To(BeNumerically(">=", 1))
			originalName := originalBursts[0].Name
			burstID := originalBursts[0].ID

			env.SelectIntentByName("burst_management")
			env.Confirm()
			env.PressKeyRune('e')

			activeIntent := env.Model.GetActiveIntent()
			burstIntent := activeIntent.(*intents.BurstManagementIntent)

			burstIntent.SetTestModalResult(&intents.ModalEditResult[*career.Burst]{
				Original: originalBursts[0],
				Modified: originalBursts[0],
				Accepted: false,
				Changes:  map[string]interface{}{},
			})

			env.Model.Update(nil)

			unchangedBursts := env.GetBursts()
			var foundBurst *career.Burst
			for _, b := range unchangedBursts {
				if b.ID == burstID {
					foundBurst = b
					break
				}
			}

			Expect(foundBurst).NotTo(BeNil())
			Expect(foundBurst.Name).To(Equal(originalName), "Burst name should NOT change when cancelled")
		})

		It("should persist changes and survive application restart", func() {
			originalBursts := env.GetBursts()
			Expect(len(originalBursts)).To(BeNumerically(">=", 1))
			burstID := originalBursts[0].ID

			env.SelectIntentByName("burst_management")
			env.Confirm()
			env.PressKeyRune('e')

			activeIntent := env.Model.GetActiveIntent()
			burstIntent := activeIntent.(*intents.BurstManagementIntent)

			newName := "PERSISTED_ACROSS_RESTART"
			modifiedBurst := &career.Burst{
				ID:          burstID,
				Name:        newName,
				Description: originalBursts[0].Description,
				EventIDs:    originalBursts[0].EventIDs,
				CreatedAt:   originalBursts[0].CreatedAt,
				UpdatedAt:   originalBursts[0].UpdatedAt,
			}

			burstIntent.SetTestModalResult(&intents.ModalEditResult[*career.Burst]{
				Original: originalBursts[0],
				Modified: modifiedBurst,
				Accepted: true,
				Changes:  map[string]interface{}{"name": newName},
			})

			env.Model.Update(nil)

			env.SimulateRestart()

			persistedBursts := env.GetBursts()
			var foundBurst *career.Burst
			for _, b := range persistedBursts {
				if b.ID == burstID {
					foundBurst = b
					break
				}
			}

			Expect(foundBurst).NotTo(BeNil(), "Burst should exist after restart")
			Expect(foundBurst.Name).To(Equal(newName), "Burst name should persist after restart")
		})
	})
})
