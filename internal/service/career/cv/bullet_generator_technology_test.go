package cv

import (
	"io"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultBulletGenerator - Technology Filtering", func() {
	var (
		generator *DefaultBulletGenerator
		log       *logger.Logger
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		generator = &DefaultBulletGenerator{
			eventRepo: nil,
			factRepo:  nil,
			logger:    log,
		}
	})

	Describe("FilterByTechnologies", func() {
		var (
			baseBullets []*career.CVBullet
			events      []*career.CareerEvent
		)

		BeforeEach(func() {
			// Create test events with different skills
			events = []*career.CareerEvent{
				{
					ID:     "event1",
					Text:   "Built API with Go",
					Date:   time.Now().AddDate(0, -1, 0),
					Skills: []string{"Go", "REST API"},
				},
				{
					ID:     "event2",
					Text:   "Developed frontend with React",
					Date:   time.Now().AddDate(0, -2, 0),
					Skills: []string{"React", "TypeScript"},
				},
				{
					ID:     "event3",
					Text:   "Led team retrospectives and planning sessions",
					Date:   time.Now().AddDate(0, -3, 0),
					Skills: []string{}, // No skills - leadership event
				},
				{
					ID:     "event4",
					Text:   "Architected Python microservices",
					Date:   time.Now().AddDate(0, -4, 0),
					Skills: []string{"Python", "Docker"},
				},
			}

			// Create corresponding bullets with base scores
			baseBullets = []*career.CVBullet{
				{
					ID:              "bullet1",
					Text:            "Built API with Go",
					SourceEventIDs:  []string{"event1"},
					InclusionReason: "event_direct",
					Confidence:      0.7,
					Rank:            0.85, // Base score from existing scoring
				},
				{
					ID:              "bullet2",
					Text:            "Developed frontend with React",
					SourceEventIDs:  []string{"event2"},
					InclusionReason: "event_direct",
					Confidence:      0.7,
					Rank:            0.75, // Lower base score
				},
				{
					ID:              "bullet3",
					Text:            "Led team retrospectives and planning sessions",
					SourceEventIDs:  []string{"event3"},
					InclusionReason: "event_direct",
					Confidence:      0.8,
					Rank:            0.90, // High-quality, no skills
				},
				{
					ID:              "bullet4",
					Text:            "Architected Python microservices",
					SourceEventIDs:  []string{"event4"},
					InclusionReason: "event_direct",
					Confidence:      0.7,
					Rank:            0.75, // Same as bullet2
				},
			}
		})

		Context("Language Agnostic focus", func() {
			It("should return all bullets without filtering", func() {
				filtered := generator.FilterByTechnologies(
					baseBullets,
					events,
					TechnologyFocusLanguageAgnostic,
					[]string{},
				)

				Expect(filtered).To(HaveLen(4))
				// Scores should remain unchanged
				Expect(filtered[0].Rank).To(Equal(0.85))
				Expect(filtered[1].Rank).To(Equal(0.75))
				Expect(filtered[2].Rank).To(Equal(0.90))
				Expect(filtered[3].Rank).To(Equal(0.75))
			})

			It("should include events without skills", func() {
				filtered := generator.FilterByTechnologies(
					baseBullets,
					events,
					TechnologyFocusLanguageAgnostic,
					[]string{},
				)

				// Bullet 3 has no skills
				noSkillsBullet := filtered[2]
				Expect(noSkillsBullet.ID).To(Equal("bullet3"))
				Expect(noSkillsBullet.Rank).To(Equal(0.90)) // No penalty
			})
		})

		Context("Events without skills are not penalized", func() {
			It("should maintain baseline score for events without skills", func() {
				// Test with Specialist focus to see skill bonus in action
				filtered := generator.FilterByTechnologies(
					baseBullets,
					events,
					TechnologyFocusSpecialist,
					[]string{"Go"},
				)

				// Find bullet3 (no skills)
				var noSkillsBullet *career.CVBullet
				for _, b := range filtered {
					if b.ID == "bullet3" {
						noSkillsBullet = b
						break
					}
				}

				Expect(noSkillsBullet).NotTo(BeNil())
				// Should keep original score (no penalty for missing skills)
				Expect(noSkillsBullet.Rank).To(Equal(0.90))
			})
		})

		Context("Events with selected technology get skill match bonus", func() {
			It("should add +0.15 bonus for matching technology", func() {
				filtered := generator.FilterByTechnologies(
					baseBullets,
					events,
					TechnologyFocusSpecialist,
					[]string{"Go"},
				)

				// Find bullet1 (has Go skill)
				var goBullet *career.CVBullet
				for _, b := range filtered {
					if b.ID == "bullet1" {
						goBullet = b
						break
					}
				}

				Expect(goBullet).NotTo(BeNil())
				// 0.85 base + 0.15 bonus = 1.00
				Expect(goBullet.Rank).To(Equal(1.00))
			})

			It("should cap score at 1.0 even with bonus", func() {
				// Create a bullet with high base score
				highScoreBullet := &career.CVBullet{
					ID:              "bullet5",
					Text:            "Critical Go implementation",
					SourceEventIDs:  []string{"event1"},
					InclusionReason: "fact_extraction",
					Confidence:      0.9,
					Rank:            0.95, // Very high base score
				}

				bullets := append(baseBullets, highScoreBullet)

				filtered := generator.FilterByTechnologies(
					bullets,
					events,
					TechnologyFocusSpecialist,
					[]string{"Go"},
				)

				// Find the high score bullet
				var highBullet *career.CVBullet
				for _, b := range filtered {
					if b.ID == "bullet5" {
						highBullet = b
						break
					}
				}

				Expect(highBullet).NotTo(BeNil())
				// 0.95 + 0.15 = 1.10, but should cap at 1.00
				Expect(highBullet.Rank).To(Equal(1.00))
			})
		})

		Context("High-quality events without skills can rank higher than low-quality events with skills", func() {
			It("should allow quality bullets without skills to outrank weak bullets with skills", func() {
				filtered := generator.FilterByTechnologies(
					baseBullets,
					events,
					TechnologyFocusSpecialist,
					[]string{"Go"},
				)

				// Find bullet2 (React, 0.75 base, no bonus) and bullet3 (no skills, 0.90 base)
				var reactBullet, leadershipBullet *career.CVBullet
				for _, b := range filtered {
					if b.ID == "bullet2" {
						reactBullet = b
					}
					if b.ID == "bullet3" {
						leadershipBullet = b
					}
				}

				Expect(reactBullet).NotTo(BeNil())
				Expect(leadershipBullet).NotTo(BeNil())

				// Leadership bullet (0.90 base, no bonus) should beat React bullet (0.75 base, no bonus)
				Expect(leadershipBullet.Rank).To(BeNumerically(">", reactBullet.Rank))
			})
		})

		Context("Specialist filters to selected tech", func() {
			It("should boost bullets with selected technology", func() {
				filtered := generator.FilterByTechnologies(
					baseBullets,
					events,
					TechnologyFocusSpecialist,
					[]string{"Go"},
				)

				// All bullets should be included (we don't filter out, just boost)
				Expect(filtered).To(HaveLen(4))

				// Find Go bullet and check boost
				var goBullet *career.CVBullet
				for _, b := range filtered {
					if b.ID == "bullet1" {
						goBullet = b
						break
					}
				}

				Expect(goBullet).NotTo(BeNil())
				Expect(goBullet.Rank).To(Equal(1.00)) // 0.85 + 0.15 = 1.00
			})

			It("should sort boosted bullets to top", func() {
				filtered := generator.FilterByTechnologies(
					baseBullets,
					events,
					TechnologyFocusSpecialist,
					[]string{"Go"},
				)

				// After sorting, Go bullet should be first (highest rank)
				Expect(filtered[0].ID).To(Equal("bullet1")) // Go bullet with 1.00 score
			})
		})

		Context("Edge cases", func() {
			It("should handle bullets without source events gracefully", func() {
				// Create bullet without source events
				orphanBullet := &career.CVBullet{
					ID:              "orphan",
					Text:            "Some achievement",
					SourceEventIDs:  []string{},
					InclusionReason: "fact_extraction",
					Confidence:      0.8,
					Rank:            0.80,
				}

				bulletsWithOrphan := append(baseBullets, orphanBullet)

				filtered := generator.FilterByTechnologies(
					bulletsWithOrphan,
					events,
					TechnologyFocusSpecialist,
					[]string{"Go"},
				)

				// Orphan should still be included
				Expect(filtered).To(HaveLen(5))

				// Find orphan and verify it kept original score
				var foundOrphan *career.CVBullet
				for _, b := range filtered {
					if b.ID == "orphan" {
						foundOrphan = b
						break
					}
				}

				Expect(foundOrphan).NotTo(BeNil())
				Expect(foundOrphan.Rank).To(Equal(0.80)) // No change
			})

			It("should handle events without matching event IDs", func() {
				// Create bullet pointing to non-existent event
				invalidBullet := &career.CVBullet{
					ID:              "invalid",
					Text:            "Some work",
					SourceEventIDs:  []string{"nonexistent"},
					InclusionReason: "event_direct",
					Confidence:      0.7,
					Rank:            0.70,
				}

				bulletsWithInvalid := append(baseBullets, invalidBullet)

				filtered := generator.FilterByTechnologies(
					bulletsWithInvalid,
					events,
					TechnologyFocusSpecialist,
					[]string{"Go"},
				)

				// Invalid should still be included
				Expect(filtered).To(HaveLen(5))

				// Find invalid and verify it kept original score
				var foundInvalid *career.CVBullet
				for _, b := range filtered {
					if b.ID == "invalid" {
						foundInvalid = b
						break
					}
				}

				Expect(foundInvalid).NotTo(BeNil())
				Expect(foundInvalid.Rank).To(Equal(0.70)) // No change
			})
		})

		Context("Generalist boosts selected techs", func() {
			It("should boost all selected technologies", func() {
				filtered := generator.FilterByTechnologies(
					baseBullets,
					events,
					TechnologyFocusGeneralist,
					[]string{"Go", "React"},
				)

				// Find Go and React bullets
				var goBullet, reactBullet *career.CVBullet
				for _, b := range filtered {
					if b.ID == "bullet1" {
						goBullet = b
					}
					if b.ID == "bullet2" {
						reactBullet = b
					}
				}

				Expect(goBullet).NotTo(BeNil())
				Expect(reactBullet).NotTo(BeNil())

				// Both should get +0.15 bonus
				Expect(goBullet.Rank).To(Equal(1.00))    // 0.85 + 0.15
				Expect(reactBullet.Rank).To(Equal(0.90)) // 0.75 + 0.15
			})

			It("should keep non-selected bullets with baseline score", func() {
				filtered := generator.FilterByTechnologies(
					baseBullets,
					events,
					TechnologyFocusGeneralist,
					[]string{"Go", "React"},
				)

				// Find Python bullet (not selected)
				var pythonBullet *career.CVBullet
				for _, b := range filtered {
					if b.ID == "bullet4" {
						pythonBullet = b
						break
					}
				}

				Expect(pythonBullet).NotTo(BeNil())
				// Should keep original score (no penalty)
				Expect(pythonBullet.Rank).To(Equal(0.75))
			})
		})
	})
})
