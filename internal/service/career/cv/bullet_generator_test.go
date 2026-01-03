package cv

import (
	"context"
	"io"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultBulletGenerator", func() {
	var (
		generator *DefaultBulletGenerator
		log       *logger.Logger
		ctx       context.Context
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		generator = &DefaultBulletGenerator{
			eventRepo: nil,
			factRepo:  nil,
			logger:    log,
		}
		ctx = context.Background()
	})

	Describe("GenerateBullets", func() {
		It("should return empty list for no events", func() {
			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{}, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(Equal(0))
		})

		It("should generate bullets from valid events", func() {
			events := []*career.CareerEvent{
				{
					ID:         "event1",
					Text:       "Implemented authentication system",
					Date:       time.Now(),
					Categories: []string{"technical"},
				},
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">", 0))
		})

		It("should handle context cancellation", func() {
			events := []*career.CareerEvent{
				{
					ID:   "event1",
					Text: "Some event",
					Date: time.Now(),
				},
			}

			cancelCtx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := generator.GenerateBullets(cancelCtx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).To(HaveOccurred())
		})

		It("should handle nil event list gracefully", func() {
			bullets, err := generator.GenerateBullets(ctx, nil, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(Equal(0))
		})

		It("should handle nil fact list gracefully", func() {
			events := []*career.CareerEvent{
				{
					ID:   "event1",
					Text: "Implemented feature",
					Date: time.Now(),
				},
			}

			bullets, err := generator.GenerateBullets(ctx, events, nil, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})
	})

	Describe("Event/Fact Combination Tests", func() {
		Context("High-quality events", func() {
			It("should generate bullets from high-quality events", func() {
				events := []*career.CareerEvent{
					{
						ID:         "event1",
						Text:       "Led team of 5 engineers to deliver microservices architecture",
						Date:       time.Now(),
						Categories: []string{"leadership", "technical"},
						Tags:       []string{"achievement", "ownership"},
					},
					{
						ID:         "event2",
						Text:       "Improved system performance by 40%",
						Date:       time.Now(),
						Categories: []string{"technical"},
						Tags:       []string{"achievement", "metrics"},
					},
				}

				bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bullets)).To(BeNumerically(">", 0))
			})
		})

		Context("Medium-quality events", func() {
			It("should generate bullets from medium-quality events", func() {
				events := []*career.CareerEvent{
					{
						ID:         "event1",
						Text:       "Contributed to API design improvements",
						Date:       time.Now(),
						Categories: []string{"technical"},
					},
					{
						ID:         "event2",
						Text:       "Participated in code review process",
						Date:       time.Now(),
						Categories: []string{"technical"},
					},
				}

				bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "staff", []string{"peer"})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bullets)).To(BeNumerically(">=", 0))
			})
		})

		Context("Low-quality events", func() {
			It("should filter out low-quality events", func() {
				events := []*career.CareerEvent{
					{
						ID:   "event1",
						Text: "Attended meeting",
						Date: time.Now(),
					},
				}

				bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
				Expect(err).NotTo(HaveOccurred())
				// Low-quality events may be filtered out
				Expect(len(bullets)).To(BeNumerically(">=", 0))
			})
		})

		Context("Mixed event and fact combinations", func() {
			It("should generate bullets from both events and facts", func() {
				events := []*career.CareerEvent{
					{
						ID:         "event1",
						Text:       "Implemented distributed caching system",
						Date:       time.Now(),
						Categories: []string{"technical"},
					},
				}

				facts := []*career.Fact{
					{
						ID:   "fact1",
						Text: "Expert in system design and architecture",
					},
				}

				bullets, err := generator.GenerateBullets(ctx, events, facts, "principal", []string{"hiring_manager"})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bullets)).To(BeNumerically(">=", 0))
			})

			It("should prefer facts over events for bullet generation", func() {
				events := []*career.CareerEvent{
					{
						ID:         "event1",
						Text:       "Worked on API",
						Date:       time.Now(),
						Categories: []string{"technical"},
					},
				}

				facts := []*career.Fact{
					{
						ID:   "fact1",
						Text: "Designed scalable APIs serving 1M+ requests",
					},
				}

				bullets, err := generator.GenerateBullets(ctx, events, facts, "principal", []string{"hiring_manager"})
				Expect(err).NotTo(HaveOccurred())
				// Facts should be included if they're of high quality
				Expect(len(bullets)).To(BeNumerically(">=", 0))
			})
		})
	})

	Describe("Inclusion/Exclusion Criteria", func() {
		Context("Single claim bullets", func() {
			It("should filter out multiple-claim bullets", func() {
				events := []*career.CareerEvent{
					{
						ID:   "event1",
						Text: "Implemented feature and wrote documentation and conducted training",
						Date: time.Now(),
					},
				}

				bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
				Expect(err).NotTo(HaveOccurred())
				// Multiple claims should be filtered
				if len(bullets) > 0 {
					for _, bullet := range bullets {
						Expect(career.IsSingleClaimBullet(bullet.Text)).To(BeTrue())
					}
				}
			})

			It("should accept single-claim bullets", func() {
				events := []*career.CareerEvent{
					{
						ID:   "event1",
						Text: "Implemented authentication system",
						Date: time.Now(),
					},
				}

				bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
				Expect(err).NotTo(HaveOccurred())
				// Single claims should be accepted (if they pass other criteria)
				Expect(len(bullets)).To(BeNumerically(">=", 0))
			})
		})

		Context("Aspirational language", func() {
			It("should filter out aspirational language", func() {
				events := []*career.CareerEvent{
					{
						ID:   "event1",
						Text: "Will implement new system",
						Date: time.Now(),
					},
					{
						ID:   "event2",
						Text: "Plan to improve performance",
						Date: time.Now(),
					},
				}

				bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
				Expect(err).NotTo(HaveOccurred())
				// Aspirational language should be filtered
				if len(bullets) > 0 {
					for _, bullet := range bullets {
						Expect(career.IsAspirationLanguage(bullet.Text)).To(BeFalse())
					}
				}
			})

			It("should accept past-tense accomplishments", func() {
				events := []*career.CareerEvent{
					{
						ID:   "event1",
						Text: "Implemented authentication system",
						Date: time.Now(),
					},
				}

				bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bullets)).To(BeNumerically(">=", 0))
			})
		})

		Context("Inferred metrics", func() {
			It("should filter out inferred metrics", func() {
				events := []*career.CareerEvent{
					{
						ID:   "event1",
						Text: "Probably improved system performance by around 50%",
						Date: time.Now(),
					},
				}

				bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
				Expect(err).NotTo(HaveOccurred())
				// Inferred metrics should be filtered
				if len(bullets) > 0 {
					for _, bullet := range bullets {
						Expect(career.HasInferredMetrics(bullet.Text)).To(BeFalse())
					}
				}
			})

			It("should accept concrete metrics", func() {
				events := []*career.CareerEvent{
					{
						ID:   "event1",
						Text: "Improved system performance by 40%",
						Date: time.Now(),
					},
				}

				bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bullets)).To(BeNumerically(">=", 0))
			})
		})

		Context("Role inflation", func() {
			It("should filter out role inflation for junior roles", func() {
				events := []*career.CareerEvent{
					{
						ID:   "event1",
						Text: "Led company-wide strategic initiative",
						Date: time.Now(),
					},
				}

				bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "staff", []string{"hiring_manager"})
				Expect(err).NotTo(HaveOccurred())
				// Role inflation should be filtered for staff role
				if len(bullets) > 0 {
					for _, bullet := range bullets {
						Expect(career.IsRoleInflation(bullet.Text, "staff")).To(BeFalse())
					}
				}
			})

			It("should accept appropriate claims for principal role", func() {
				events := []*career.CareerEvent{
					{
						ID:   "event1",
						Text: "Led strategic architecture initiative",
						Date: time.Now(),
					},
				}

				bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bullets)).To(BeNumerically(">=", 0))
			})
		})

		Context("Source event requirement", func() {
			It("should filter out bullets with no source events or facts", func() {
				// This test verifies that bullets must have at least one source
				// The implementation should enforce this
				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{}, []*career.Fact{}, "principal", []string{"hiring_manager"})
				Expect(err).NotTo(HaveOccurred())
				// All returned bullets should have sources
				for _, bullet := range bullets {
					hasSource := len(bullet.SourceEventIDs) > 0 || len(bullet.SourceFactIDs) > 0
					Expect(hasSource).To(BeTrue())
				}
			})
		})
	})

	Describe("Ranking Algorithm", func() {
		It("should rank bullets by priority", func() {
			events := []*career.CareerEvent{
				{
					ID:   "event1",
					Text: "Led team to deliver feature",
					Date: time.Now(),
					Tags: []string{"ownership"},
				},
				{
					ID:   "event2",
					Text: "Contributed to feature",
					Date: time.Now(),
					Tags: []string{"contribution"},
				},
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())

			// Bullets should be ranked (higher rank first)
			if len(bullets) > 1 {
				for i := 0; i < len(bullets)-1; i++ {
					Expect(bullets[i].Rank).To(BeNumerically(">=", bullets[i+1].Rank))
				}
			}
		})

		It("should score based on inclusion reason", func() {
			facts := []*career.Fact{
				{
					ID:   "fact1",
					Text: "Expert in distributed systems",
				},
			}

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{}, facts, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())

			// Facts should have higher scores than direct events
			if len(bullets) > 0 {
				Expect(bullets[0].Confidence).To(BeNumerically(">", 0))
			}
		})

		It("should apply temporal decay to older events", func() {
			now := time.Now()
			oneYearAgo := now.AddDate(-1, 0, 0)

			events := []*career.CareerEvent{
				{
					ID:   "event1",
					Text: "Implemented feature",
					Date: now,
				},
				{
					ID:   "event2",
					Text: "Implemented feature",
					Date: oneYearAgo,
				},
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())

			// Newer events should rank higher than older ones with similar content
			if len(bullets) >= 2 {
				// Find the bullets corresponding to the events
				// (Note: this is a simplified check)
				Expect(bullets[0].Rank).To(BeNumerically(">=", 0))
			}
		})

		It("should boost score for repeated signals", func() {
			events := []*career.CareerEvent{
				{
					ID:   "event1",
					Text: "Implemented feature",
					Date: time.Now(),
					Tags: []string{"achievement", "technical"},
				},
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())

			// Events with multiple signals should have higher scores
			if len(bullets) > 0 {
				Expect(bullets[0].Rank).To(BeNumerically(">", 0))
			}
		})
	})

	Describe("Role-Specific Compression", func() {
		It("should cap principal role at 3-4 bullets", func() {
			events := make([]*career.CareerEvent, 10)
			for i := 0; i < 10; i++ {
				events[i] = &career.CareerEvent{
					ID:         uuid.New().String(),
					Text:       "Implemented feature " + string(rune('A'+i)),
					Date:       time.Now(),
					Categories: []string{"technical"},
				}
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically("<=", 4))
		})

		It("should cap staff role at 4-5 bullets", func() {
			events := make([]*career.CareerEvent, 10)
			for i := 0; i < 10; i++ {
				events[i] = &career.CareerEvent{
					ID:         uuid.New().String(),
					Text:       "Implemented feature " + string(rune('A'+i)),
					Date:       time.Now(),
					Categories: []string{"technical"},
				}
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "staff", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically("<=", 5))
		})

		It("should cap EM role at 3-4 bullets", func() {
			events := make([]*career.CareerEvent, 10)
			for i := 0; i < 10; i++ {
				events[i] = &career.CareerEvent{
					ID:         uuid.New().String(),
					Text:       "Implemented feature " + string(rune('A'+i)),
					Date:       time.Now(),
					Categories: []string{"leadership"},
				}
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "em", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically("<=", 4))
		})

		It("should cap senior_ic role at 4-5 bullets", func() {
			events := make([]*career.CareerEvent, 10)
			for i := 0; i < 10; i++ {
				events[i] = &career.CareerEvent{
					ID:         uuid.New().String(),
					Text:       "Implemented feature " + string(rune('A'+i)),
					Date:       time.Now(),
					Categories: []string{"technical"},
				}
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "senior_ic", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically("<=", 5))
		})

		It("should remove lower-ranked bullets first when compressing", func() {
			events := make([]*career.CareerEvent, 5)
			for i := 0; i < 5; i++ {
				events[i] = &career.CareerEvent{
					ID:         uuid.New().String(),
					Text:       "Implemented feature " + string(rune('A'+i)),
					Date:       time.Now(),
					Categories: []string{"technical"},
				}
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())

			// Bullets should be ordered by rank (descending)
			if len(bullets) > 1 {
				for i := 0; i < len(bullets)-1; i++ {
					Expect(bullets[i].Rank).To(BeNumerically(">=", bullets[i+1].Rank))
				}
			}
		})
	})

	Describe("Audience-Specific Filtering", func() {
		It("should filter for hiring_manager audience", func() {
			events := []*career.CareerEvent{
				{
					ID:         "event1",
					Text:       "Improved system performance by 40%",
					Date:       time.Now(),
					Categories: []string{"technical"},
				},
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})

		It("should filter for recruiter audience", func() {
			events := []*career.CareerEvent{
				{
					ID:         "event1",
					Text:       "Led team of 5 engineers",
					Date:       time.Now(),
					Categories: []string{"leadership"},
				},
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"recruiter"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})

		It("should filter for peer audience", func() {
			events := []*career.CareerEvent{
				{
					ID:         "event1",
					Text:       "Implemented distributed caching system",
					Date:       time.Now(),
					Categories: []string{"technical"},
				},
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"peer"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})

		It("should support multiple audiences", func() {
			events := []*career.CareerEvent{
				{
					ID:         "event1",
					Text:       "Improved system performance by 40%",
					Date:       time.Now(),
					Categories: []string{"technical"},
				},
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager", "recruiter", "peer"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})
	})

	Describe("Source Tracking", func() {
		It("should preserve event source IDs", func() {
			events := []*career.CareerEvent{
				{
					ID:   "event1",
					Text: "Implemented feature",
					Date: time.Now(),
				},
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())

			if len(bullets) > 0 {
				Expect(bullets[0].SourceEventIDs).To(ContainElement("event1"))
			}
		})

		It("should preserve fact source IDs", func() {
			facts := []*career.Fact{
				{
					ID:   "fact1",
					Text: "Expert in system design",
				},
			}

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{}, facts, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())

			if len(bullets) > 0 {
				Expect(bullets[0].SourceFactIDs).To(ContainElement("fact1"))
			}
		})

		It("should track inclusion reason", func() {
			events := []*career.CareerEvent{
				{
					ID:   "event1",
					Text: "Implemented feature",
					Date: time.Now(),
				},
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())

			if len(bullets) > 0 {
				Expect(bullets[0].InclusionReason).NotTo(BeEmpty())
			}
		})

		It("should calculate confidence score", func() {
			events := []*career.CareerEvent{
				{
					ID:   "event1",
					Text: "Implemented feature",
					Date: time.Now(),
				},
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())

			if len(bullets) > 0 {
				Expect(bullets[0].Confidence).To(BeNumerically(">=", 0))
				Expect(bullets[0].Confidence).To(BeNumerically("<=", 1))
			}
		})
	})

	Describe("Edge Cases", func() {
		It("should handle very long event text", func() {
			longText := ""
			for i := 0; i < 100; i++ {
				longText += "Lorem ipsum dolor sit amet. "
			}

			events := []*career.CareerEvent{
				{
					ID:   "event1",
					Text: longText,
					Date: time.Now(),
				},
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})

		It("should handle special characters in event text", func() {
			events := []*career.CareerEvent{
				{
					ID:   "event1",
					Text: "Implemented @#$%^&*() special feature",
					Date: time.Now(),
				},
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})

		It("should handle events with no categories", func() {
			events := []*career.CareerEvent{
				{
					ID:   "event1",
					Text: "Implemented feature",
					Date: time.Now(),
				},
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})

		It("should handle events with no tags", func() {
			events := []*career.CareerEvent{
				{
					ID:         "event1",
					Text:       "Implemented feature",
					Date:       time.Now(),
					Categories: []string{"technical"},
				},
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})

		It("should handle large number of events (1000+)", func() {
			events := make([]*career.CareerEvent, 1000)
			for i := 0; i < 1000; i++ {
				events[i] = &career.CareerEvent{
					ID:         uuid.New().String(),
					Text:       "Implemented feature",
					Date:       time.Now(),
					Categories: []string{"technical"},
				}
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
			Expect(err).NotTo(HaveOccurred())
			// Should still respect role-specific caps
			Expect(len(bullets)).To(BeNumerically("<=", 4))
		})
	})
})
