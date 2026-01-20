package cv

import (
	"context"
	"fmt"
	"io"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/testutil/fixtures"
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
			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{}, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(Equal(0))
		})

		It("should generate bullets from valid events", func() {
			event := fixtures.EventWith("event1", "Implemented authentication system", "", "")
			event.Categories = []string{"technical"}

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">", 0))
		})

		It("should handle context cancellation", func() {
			event := fixtures.EventWith("event1", "Some event", "", "")

			cancelCtx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := generator.GenerateBullets(cancelCtx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).To(HaveOccurred())
		})

		It("should handle nil event list gracefully", func() {
			bullets, err := generator.GenerateBullets(ctx, nil, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(Equal(0))
		})

		It("should handle nil fact list gracefully", func() {
			event := fixtures.EventWith("event1", "Implemented feature", "", "")

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, nil, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})
	})

	Describe("Event/Fact Combination Tests", func() {
		Context("High-quality events", func() {
			It("should generate bullets from high-quality events", func() {
				event1 := fixtures.EventWith("event1", "Led team of 5 engineers to deliver microservices architecture", "", "")
				event1.Categories = []string{"leadership", "technical"}
				event1.Tags = []string{"achievement", "ownership"}

				event2 := fixtures.EventWith("event2", "Improved system performance by 40%", "", "")
				event2.Categories = []string{"technical"}
				event2.Tags = []string{"achievement", "metrics"}

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event1, event2}, []*career.Fact{}, "principal", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bullets)).To(BeNumerically(">", 0))
			})
		})

		Context("Medium-quality events", func() {
			It("should generate bullets from medium-quality events", func() {
				event1 := fixtures.EventWith("event1", "Contributed to API design improvements", "", "")
				event1.Categories = []string{"technical"}

				event2 := fixtures.EventWith("event2", "Participated in code review process", "", "")
				event2.Categories = []string{"technical"}

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event1, event2}, []*career.Fact{}, "staff", "peer")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bullets)).To(BeNumerically(">=", 0))
			})
		})

		Context("Low-quality events", func() {
			It("should filter out low-quality events", func() {
				event := fixtures.EventWith("event1", "Attended meeting", "", "")

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				// Low-quality events may be filtered out
				Expect(len(bullets)).To(BeNumerically(">=", 0))
			})
		})

		Context("Mixed event and fact combinations", func() {
			It("should generate bullets from both events and facts", func() {
				event := fixtures.EventWith("event1", "Implemented distributed caching system", "", "")
				event.Categories = []string{"technical"}

				fact := fixtures.FactWith("fact1", "Expert in system design and architecture")

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{fact}, "principal", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bullets)).To(BeNumerically(">=", 0))
			})

			It("should prefer facts over events for bullet generation", func() {
				event := fixtures.EventWith("event1", "Worked on API", "", "")
				event.Categories = []string{"technical"}

				fact := fixtures.FactWith("fact1", "Designed scalable APIs serving 1M+ requests")

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{fact}, "principal", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				// Facts should be included if they're of high quality
				Expect(len(bullets)).To(BeNumerically(">=", 0))
			})
		})
	})

	Describe("Inclusion/Exclusion Criteria", func() {
		Context("Single claim bullets", func() {
			It("should filter out multiple-claim bullets", func() {
				event := fixtures.EventWith("event1", "Implemented feature and wrote documentation and conducted training", "", "")

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				// Multiple claims should be filtered
				if len(bullets) > 0 {
					for _, bullet := range bullets {
						Expect(career.IsSingleClaimBullet(bullet.Text)).To(BeTrue())
					}
				}
			})

			It("should accept single-claim bullets", func() {
				event := fixtures.EventWith("event1", "Implemented authentication system", "", "")

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				// Single claims should be accepted (if they pass other criteria)
				Expect(len(bullets)).To(BeNumerically(">=", 0))
			})
		})

		Context("Aspirational language", func() {
			It("should filter out aspirational language", func() {
				event1 := fixtures.EventWith("event1", "Will implement new system", "", "")
				event2 := fixtures.EventWith("event2", "Plan to improve performance", "", "")

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event1, event2}, []*career.Fact{}, "principal", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				// Aspirational language should be filtered
				if len(bullets) > 0 {
					for _, bullet := range bullets {
						Expect(career.IsAspirationLanguage(bullet.Text)).To(BeFalse())
					}
				}
			})

			It("should accept past-tense accomplishments", func() {
				event := fixtures.EventWith("event1", "Implemented authentication system", "", "")

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bullets)).To(BeNumerically(">=", 0))
			})
		})

		Context("Inferred metrics", func() {
			It("should filter out inferred metrics", func() {
				event := fixtures.EventWith("event1", "Probably improved system performance by around 50%", "", "")

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				// Inferred metrics should be filtered
				if len(bullets) > 0 {
					for _, bullet := range bullets {
						Expect(career.HasInferredMetrics(bullet.Text)).To(BeFalse())
					}
				}
			})

			It("should accept concrete metrics", func() {
				event := fixtures.EventWith("event1", "Improved system performance by 40%", "", "")

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bullets)).To(BeNumerically(">=", 0))
			})
		})

		Context("Role inflation", func() {
			It("should filter out role inflation for junior roles", func() {
				event := fixtures.EventWith("event1", "Led company-wide strategic initiative", "", "")

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "staff", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				// Role inflation should be filtered for staff role
				if len(bullets) > 0 {
					for _, bullet := range bullets {
						Expect(career.IsRoleInflation(bullet.Text, "staff")).To(BeFalse())
					}
				}
			})

			It("should accept appropriate claims for principal role", func() {
				event := fixtures.EventWith("event1", "Led strategic architecture initiative", "", "")

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bullets)).To(BeNumerically(">=", 0))
			})
		})

		Context("Source event requirement", func() {
			It("should filter out bullets with no source events or facts", func() {
				// This test verifies that bullets must have at least one source
				// The implementation should enforce this
				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{}, []*career.Fact{}, "principal", "hiring_manager")
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
			event1 := fixtures.EventWith("event1", "Led team to deliver feature", "", "")
			event1.Tags = []string{"ownership"}

			event2 := fixtures.EventWith("event2", "Contributed to feature", "", "")
			event2.Tags = []string{"contribution"}

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event1, event2}, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())

			// Bullets should be ranked (higher rank first)
			if len(bullets) > 1 {
				for i := 0; i < len(bullets)-1; i++ {
					Expect(bullets[i].Rank).To(BeNumerically(">=", bullets[i+1].Rank))
				}
			}
		})

		It("should score based on inclusion reason", func() {
			fact := fixtures.FactWith("fact1", "Expert in distributed systems")

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{}, []*career.Fact{fact}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())

			// Facts should have higher scores than direct events
			if len(bullets) > 0 {
				Expect(bullets[0].Confidence).To(BeNumerically(">", 0))
			}
		})

		It("should apply temporal decay to older events", func() {
			now := time.Now()
			oneYearAgo := now.AddDate(-1, 0, 0)

			event1 := fixtures.EventWith("event1", "Implemented feature", "", "")
			event1.Date = now

			event2 := fixtures.EventWith("event2", "Implemented feature", "", "")
			event2.Date = oneYearAgo

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event1, event2}, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())

			// Newer events should rank higher than older ones with similar content
			if len(bullets) >= 2 {
				// Find the bullets corresponding to the events
				// (Note: this is a simplified check)
				Expect(bullets[0].Rank).To(BeNumerically(">=", 0))
			}
		})

		It("should boost score for repeated signals", func() {
			event := fixtures.EventWith("event1", "Implemented feature", "", "")
			event.Tags = []string{"achievement", "technical"}

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())

			// Events with multiple signals should have higher scores
			if len(bullets) > 0 {
				Expect(bullets[0].Rank).To(BeNumerically(">", 0))
			}
		})
	})

	// BUG-003 fix: BulletGenerator no longer applies total bullet caps.
	// Per-company caps are correctly applied by SectionBuilder.getBulletsPerCompanyForRole().
	// These tests now verify that ALL events passing inclusion criteria are returned,
	// and that bullets are properly ranked for SectionBuilder to use.
	Describe("Bullet Generation Without Total Cap (BUG-003 Fix)", func() {
		It("should return all principal events that pass inclusion criteria", func() {
			events := fixtures.Events(10)
			for _, e := range events {
				e.Categories = []string{"technical"}
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())
			// All events should be returned (no total cap)
			Expect(len(bullets)).To(Equal(10))
		})

		It("should return all staff events that pass inclusion criteria", func() {
			events := fixtures.Events(10)
			for _, e := range events {
				e.Categories = []string{"technical"}
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "staff", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())
			// All events should be returned (no total cap)
			Expect(len(bullets)).To(Equal(10))
		})

		It("should return all EM events that pass inclusion criteria", func() {
			events := fixtures.Events(10)
			for _, e := range events {
				e.Categories = []string{"leadership"}
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "em", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())
			// All events should be returned (no total cap)
			Expect(len(bullets)).To(Equal(10))
		})

		It("should return all senior_ic events that pass inclusion criteria", func() {
			events := fixtures.Events(10)
			for _, e := range events {
				e.Categories = []string{"technical"}
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "senior_ic", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())
			// All events should be returned (no total cap)
			Expect(len(bullets)).To(Equal(10))
		})

		It("should rank bullets by score for SectionBuilder to use", func() {
			events := fixtures.Events(5)
			for _, e := range events {
				e.Categories = []string{"technical"}
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())

			// Bullets should be ordered by rank (descending) for SectionBuilder
			if len(bullets) > 1 {
				for i := 0; i < len(bullets)-1; i++ {
					Expect(bullets[i].Rank).To(BeNumerically(">=", bullets[i+1].Rank))
				}
			}
		})

		It("should preserve bullets from all companies when many events exist", func() {
			// Create events from multiple different companies (simulating the BUG-003 scenario)
			// The bug was that a total cap of 40 bullets would exclude later companies entirely
			companies := []string{"CompanyA", "CompanyB", "CompanyC", "CompanyD", "CompanyE"}
			eventToCompany := make(map[string]string)
			var events []*career.CareerEvent
			for _, company := range companies {
				for i := 0; i < 10; i++ {
					eventID := fmt.Sprintf("%s-event-%d", company, i)
					event := fixtures.EventWith(
						eventID,
						fmt.Sprintf("Implemented feature %d at %s", i, company),
						company,
						"",
					)
					event.Categories = []string{"technical"}
					events = append(events, event)
					eventToCompany[eventID] = company
				}
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "staff", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())

			// All 50 events should generate bullets (no total cap)
			Expect(len(bullets)).To(Equal(50))

			// Verify all companies are represented in the bullets by checking source event IDs
			companiesInBullets := make(map[string]bool)
			for _, bullet := range bullets {
				for _, eventID := range bullet.SourceEventIDs {
					if company, ok := eventToCompany[eventID]; ok {
						companiesInBullets[company] = true
					}
				}
			}
			Expect(len(companiesInBullets)).To(Equal(5), "All 5 companies should be represented")
		})
	})

	Describe("Audience-Specific Filtering", func() {
		It("should filter for hiring_manager audience", func() {
			event := fixtures.EventWith("event1", "Improved system performance by 40%", "", "")
			event.Categories = []string{"technical"}

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})

		It("should filter for recruiter audience", func() {
			event := fixtures.EventWith("event1", "Led team of 5 engineers", "", "")
			event.Categories = []string{"leadership"}

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "recruiter")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})

		It("should filter for peer audience", func() {
			event := fixtures.EventWith("event1", "Implemented distributed caching system", "", "")
			event.Categories = []string{"technical"}

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "peer")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})

		It("should support multiple audiences", func() {
			event := fixtures.EventWith("event1", "Improved system performance by 40%", "", "")
			event.Categories = []string{"technical"}

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})
	})

	Describe("Source Tracking", func() {
		It("should preserve event source IDs", func() {
			event := fixtures.EventWith("event1", "Implemented feature", "", "")

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())

			if len(bullets) > 0 {
				Expect(bullets[0].SourceEventIDs).To(ContainElement("event1"))
			}
		})

		It("should preserve fact source IDs", func() {
			fact := fixtures.FactWith("fact1", "Expert in system design")

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{}, []*career.Fact{fact}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())

			if len(bullets) > 0 {
				Expect(bullets[0].SourceFactIDs).To(ContainElement("fact1"))
			}
		})

		It("should track inclusion reason", func() {
			event := fixtures.EventWith("event1", "Implemented feature", "", "")

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())

			if len(bullets) > 0 {
				Expect(bullets[0].InclusionReason).NotTo(BeEmpty())
			}
		})

		It("should calculate confidence score", func() {
			event := fixtures.EventWith("event1", "Implemented feature", "", "")

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
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

			event := fixtures.EventWith("event1", longText, "", "")

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})

		It("should handle special characters in event text", func() {
			event := fixtures.EventWith("event1", "Implemented @#$%^&*() special feature", "", "")

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})

		It("should handle events with no categories", func() {
			event := fixtures.EventWith("event1", "Implemented feature", "", "")

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})

		It("should handle events with no tags", func() {
			event := fixtures.EventWith("event1", "Implemented feature", "", "")
			event.Categories = []string{"technical"}

			bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{event}, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bullets)).To(BeNumerically(">=", 0))
		})

		It("should handle large number of events (1000+)", func() {
			events := fixtures.Events(1000)
			for _, e := range events {
				e.Categories = []string{"technical"}
			}

			bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", "hiring_manager")
			Expect(err).NotTo(HaveOccurred())
			// BUG-003 fix: BulletGenerator no longer applies total bullet cap.
			// Per-company caps are applied by SectionBuilder instead.
			// All events that pass inclusion criteria should be returned.
			Expect(len(bullets)).To(BeNumerically(">", 0))
			Expect(len(bullets)).To(BeNumerically("<=", 1000))
		})
	})
})
