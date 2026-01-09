package cv

import (
	"context"
	"io"
	"strings"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EnhancedBulletGenerator", func() {
	var (
		generator EnhancedBulletGenerator
		log       *logger.Logger
		ctx       context.Context
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		generator = NewEnhancedBulletGenerator(log)
		ctx = context.Background()
	})

	It("should return empty list for no input", func() {
		bullets, err := generator.GenerateBullets(ctx, nil, nil, nil, "principal", "hiring_manager")
		Expect(err).NotTo(HaveOccurred())
		Expect(bullets).To(BeEmpty())
	})

	It("should generate bullets from events", func() {
		events := []*career.CareerEvent{
			{
				ID:   "e1",
				Text: "Led team to deliver microservices",
				Date: time.Now(),
			},
		}

		bullets, err := generator.GenerateBullets(ctx, events, nil, nil, "principal", "hiring_manager")
		Expect(err).NotTo(HaveOccurred())
		Expect(len(bullets)).To(BeNumerically(">", 0))
	})

	It("should filter by role confidence", func() {
		bullets := []*EnhancedBullet{
			{ID: "b1", Confidence: 0.85},
			{ID: "b2", Confidence: 0.70},
		}

		filtered := generator.FilterByRole(bullets, "principal")
		Expect(len(filtered)).To(Equal(1))
	})

	It("should rank bullets by score", func() {
		bullets := []*EnhancedBullet{
			{ID: "b1", RoleScore: 0.8, AudienceScore: 0.7, MetricScore: 0.6, ImpactScore: 0.7, Confidence: 0.8},
			{ID: "b2", RoleScore: 0.5, AudienceScore: 0.5, MetricScore: 0.5, ImpactScore: 0.5, Confidence: 0.5},
		}

		ranked := generator.RankByRelevance(bullets, "principal", "hiring_manager")
		Expect(ranked[0].ID).To(Equal("b1"))
		Expect(ranked[1].ID).To(Equal("b2"))
	})

	It("should enhance bullet wording", func() {
		bullet := &EnhancedBullet{
			ID:   "b1",
			Text: "Worked on authentication system",
		}

		enhanced, err := generator.EnhanceWording(bullet, "principal")
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.Contains(enhanced.EnhancedText, "Led")).To(BeTrue())
	})

	Describe("Audience-Specific Filtering", func() {
		Context("when filtering facts by audience relevance", func() {
			It("should include facts relevant to hiring_manager audience", func() {
				facts := []*career.Fact{
					{
						ID:                "fact1",
						Text:              "Delivered 40% cost reduction through system optimization",
						AudienceRelevance: []string{"hiring_manager"},
						SourceEventID:     "event1",
					},
					{
						ID:                "fact2",
						Text:              "Implemented complex distributed caching algorithm",
						AudienceRelevance: []string{"peer"},
						SourceEventID:     "event2",
					},
				}

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{}, facts, nil, "principal", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				// Should only include fact1 (relevant to hiring_manager)
				// fact2 is only relevant to peer, should be excluded
				Expect(len(bullets)).To(Equal(1))
				Expect(bullets[0].Text).To(ContainSubstring("cost reduction"))
			})

			It("should include facts relevant to recruiter audience", func() {
				facts := []*career.Fact{
					{
						ID:                "fact1",
						Text:              "Senior engineer with expertise in Go, Ruby, and Python",
						AudienceRelevance: []string{"recruiter"},
						SourceEventID:     "event1",
					},
					{
						ID:                "fact2",
						Text:              "Deep expertise in consensus algorithms and CAP theorem",
						AudienceRelevance: []string{"peer"},
						SourceEventID:     "event2",
					},
				}

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{}, facts, nil, "principal", "recruiter")
				Expect(err).NotTo(HaveOccurred())
				// Should only include fact1 (relevant to recruiter)
				Expect(len(bullets)).To(Equal(1))
				Expect(bullets[0].Text).To(ContainSubstring("Go, Ruby"))
			})

			It("should include facts relevant to peer audience", func() {
				facts := []*career.Fact{
					{
						ID:                "fact1",
						Text:              "Designed novel approach to distributed consensus using Raft",
						AudienceRelevance: []string{"peer"},
						SourceEventID:     "event1",
					},
					{
						ID:                "fact2",
						Text:              "Reduced operational costs by 30%",
						AudienceRelevance: []string{"hiring_manager"},
						SourceEventID:     "event2",
					},
				}

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{}, facts, nil, "principal", "peer")
				Expect(err).NotTo(HaveOccurred())
				// Should only include fact1 (relevant to peer)
				Expect(len(bullets)).To(Equal(1))
				Expect(bullets[0].Text).To(ContainSubstring("Raft"))
			})

			It("should return all facts when audience is empty", func() {
				facts := []*career.Fact{
					{
						ID:                "fact1",
						Text:              "Business impact achievement",
						AudienceRelevance: []string{"hiring_manager"},
						SourceEventID:     "event1",
					},
					{
						ID:                "fact2",
						Text:              "Technical depth achievement",
						AudienceRelevance: []string{"peer"},
						SourceEventID:     "event2",
					},
				}

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{}, facts, nil, "principal", "")
				Expect(err).NotTo(HaveOccurred())
				// Should include both facts when no audience filter
				Expect(len(bullets)).To(Equal(2))
			})

			It("should exclude facts not relevant to selected audience", func() {
				facts := []*career.Fact{
					{
						ID:                "fact1",
						Text:              "Technical implementation detail for peers",
						AudienceRelevance: []string{"peer"},
						SourceEventID:     "event1",
					},
				}

				bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{}, facts, nil, "principal", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				// fact1 is only relevant to peer, not hiring_manager
				Expect(len(bullets)).To(Equal(0))
			})

			It("should include facts with multiple audience relevance", func() {
				facts := []*career.Fact{
					{
						ID:                "fact1",
						Text:              "Led migration that reduced costs and improved architecture",
						AudienceRelevance: []string{"hiring_manager", "peer"},
						SourceEventID:     "event1",
					},
				}

				// Should be included for hiring_manager
				bulletsMgr, err := generator.GenerateBullets(ctx, []*career.CareerEvent{}, facts, nil, "principal", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bulletsMgr)).To(Equal(1))

				// Should also be included for peer
				bulletsPeer, err := generator.GenerateBullets(ctx, []*career.CareerEvent{}, facts, nil, "principal", "peer")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bulletsPeer)).To(Equal(1))
			})
		})
	})
})
