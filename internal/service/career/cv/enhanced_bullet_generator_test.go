package cv

import (
	"context"
	"io"
	"strings"

	"github.com/baphled/kariya/internal/constants"
	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/testutil/fixtures"
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
			fixtures.EventWith("e1", "Led team to deliver microservices", "", ""),
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

// Task 44: Conversion helper tests
var _ = Describe("EnhancedBullet Conversion", func() {
	Describe("ToCVBullet", func() {
		It("should convert all fields correctly", func() {
			enhanced := &EnhancedBullet{
				ID:              "bullet-1",
				Text:            "Led team to deliver microservices",
				EnhancedText:    "Spearheaded cross-functional team to architect and deliver microservices platform",
				SourceEventIDs:  []string{"event-1", "event-2"},
				SourceFactIDs:   []string{"fact-1"},
				Confidence:      0.9,
				RoleScore:       0.85,
				AudienceScore:   0.80,
				MetricScore:     0.75,
				ImpactScore:     0.95,
				ImpactLevel:     "high",
				KeywordMatches:  []string{"leadership", "architecture"},
				InclusionReason: string(constants.InclusionReasonAchievementExtraction),
				Rank:            0.88,
			}

			cvBullet := enhanced.ToCVBullet()

			Expect(cvBullet.ID).To(Equal("bullet-1"))
			Expect(cvBullet.Text).To(Equal("Spearheaded cross-functional team to architect and deliver microservices platform"))
			Expect(cvBullet.EnhancedText).To(Equal("Spearheaded cross-functional team to architect and deliver microservices platform"))
			Expect(cvBullet.SourceEventIDs).To(Equal([]string{"event-1", "event-2"}))
			Expect(cvBullet.SourceFactIDs).To(Equal([]string{"fact-1"}))
			Expect(cvBullet.Confidence).To(Equal(0.9))
			Expect(cvBullet.RoleScore).To(Equal(0.85))
			Expect(cvBullet.AudienceScore).To(Equal(0.80))
			Expect(cvBullet.MetricScore).To(Equal(0.75))
			Expect(cvBullet.ImpactScore).To(Equal(0.95))
			Expect(cvBullet.ImpactLevel).To(Equal("high"))
			Expect(cvBullet.KeywordMatches).To(Equal([]string{"leadership", "architecture"}))
			Expect(cvBullet.InclusionReason).To(Equal(string(constants.InclusionReasonAchievementExtraction)))
			Expect(cvBullet.Rank).To(Equal(0.88))
		})

		It("should use original Text when EnhancedText is empty", func() {
			enhanced := &EnhancedBullet{
				ID:           "bullet-2",
				Text:         "Original bullet text",
				EnhancedText: "",
			}

			cvBullet := enhanced.ToCVBullet()

			Expect(cvBullet.Text).To(Equal("Original bullet text"))
			Expect(cvBullet.EnhancedText).To(Equal(""))
		})

		It("should use EnhancedText as Text when available", func() {
			enhanced := &EnhancedBullet{
				ID:           "bullet-3",
				Text:         "Original text",
				EnhancedText: "Enhanced and improved text",
			}

			cvBullet := enhanced.ToCVBullet()

			Expect(cvBullet.Text).To(Equal("Enhanced and improved text"))
		})
	})

	Describe("ConvertBullets", func() {
		It("should convert a slice of EnhancedBullets to CVBullets", func() {
			enhanced := []*EnhancedBullet{
				{ID: "b1", Text: "Bullet 1", Confidence: 0.8},
				{ID: "b2", Text: "Bullet 2", Confidence: 0.9},
				{ID: "b3", Text: "Bullet 3", Confidence: 0.7},
			}

			cvBullets := ConvertBullets(enhanced)

			Expect(len(cvBullets)).To(Equal(3))
			Expect(cvBullets[0].ID).To(Equal("b1"))
			Expect(cvBullets[1].ID).To(Equal("b2"))
			Expect(cvBullets[2].ID).To(Equal("b3"))
		})

		It("should return nil for nil input", func() {
			cvBullets := ConvertBullets(nil)
			Expect(cvBullets).To(BeNil())
		})

		It("should return empty slice for empty input", func() {
			cvBullets := ConvertBullets([]*EnhancedBullet{})
			Expect(cvBullets).To(BeEmpty())
		})
	})
})

// BUG-008: Role-based CV differentiation tests
var _ = Describe("BUG-008: Role-based scoring", func() {
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

	Describe("Category propagation", func() {
		Context("from CareerEvent to EnhancedBullet", func() {
			It("should propagate primary category from event", func() {
				events := []*career.CareerEvent{
					fixtures.EventWithCategories("evt-1", "Led team migration to Kubernetes", []string{"leadership", "technical"}),
				}

				bullets, err := generator.GenerateBullets(ctx, events, nil, nil, "principal", "")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bullets)).To(BeNumerically(">", 0))
				// Primary category (first in list) should be propagated
				Expect(bullets[0].Category).To(Equal(constants.CompetencyLeadership))
			})

			It("should propagate technical category from event", func() {
				events := []*career.CareerEvent{
					fixtures.EventWithCategories("evt-2", "Built real-time data pipeline", []string{"technical"}),
				}

				bullets, err := generator.GenerateBullets(ctx, events, nil, nil, "senior_ic", "")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bullets)).To(BeNumerically(">", 0))
				Expect(bullets[0].Category).To(Equal(constants.CompetencyTechnical))
			})
		})

		Context("from Fact to EnhancedBullet", func() {
			It("should propagate primary competency category from fact", func() {
				facts := []*career.Fact{
					fixtures.FactWithCategories("fact-1", "Mentored 4 junior engineers", "evt-1",
						[]string{"mentoring", "leadership"}, []string{"hiring_manager"}),
				}

				bullets, err := generator.GenerateBullets(ctx, nil, facts, nil, "em", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bullets)).To(BeNumerically(">", 0))
				Expect(bullets[0].Category).To(Equal(constants.CompetencyMentoring))
			})
		})

		Context("ToCVBullet conversion", func() {
			It("should include category in CVBullet", func() {
				enhanced := &EnhancedBullet{
					ID:       "bullet-1",
					Text:     "Technical achievement",
					Category: constants.CompetencyTechnical,
				}

				cvBullet := enhanced.ToCVBullet()
				Expect(cvBullet.Category).To(Equal(constants.CompetencyTechnical))
			})
		})
	})

	Describe("Role-based scoring", func() {
		var defaultGen *DefaultEnhancedBulletGenerator

		BeforeEach(func() {
			defaultGen = generator.(*DefaultEnhancedBulletGenerator)
		})

		Context("calculateRoleScore with categories", func() {
			It("should score leadership bullets higher for principal than senior_ic", func() {
				bullet := &EnhancedBullet{
					Category:   constants.CompetencyLeadership,
					Confidence: 0.7,
				}

				principalScore := defaultGen.calculateRoleScore(bullet, "principal")
				seniorScore := defaultGen.calculateRoleScore(bullet, "senior_ic")

				Expect(principalScore).To(BeNumerically(">", seniorScore))
			})

			It("should score technical bullets higher for senior_ic than principal", func() {
				bullet := &EnhancedBullet{
					Category:   constants.CompetencyTechnical,
					Confidence: 0.7,
				}

				seniorScore := defaultGen.calculateRoleScore(bullet, "senior_ic")
				principalScore := defaultGen.calculateRoleScore(bullet, "principal")

				Expect(seniorScore).To(BeNumerically(">", principalScore))
			})

			It("should score mentoring bullets higher for em than senior_ic", func() {
				bullet := &EnhancedBullet{
					Category:   constants.CompetencyMentoring,
					Confidence: 0.7,
				}

				emScore := defaultGen.calculateRoleScore(bullet, "em")
				seniorScore := defaultGen.calculateRoleScore(bullet, "senior_ic")

				Expect(emScore).To(BeNumerically(">", seniorScore))
			})

			It("should give primary category bullets a strong boost", func() {
				// Technical is primary for senior_ic
				bullet := &EnhancedBullet{
					Category:   constants.CompetencyTechnical,
					Confidence: 0.7,
				}

				score := defaultGen.calculateRoleScore(bullet, "senior_ic")
				expectedMin := roleScoreBase + roleScorePrimaryCategoryBoost
				Expect(score).To(BeNumerically(">=", expectedMin))
			})

			It("should give secondary category bullets a medium boost", func() {
				// Leadership is secondary for senior_ic
				bullet := &EnhancedBullet{
					Category:   constants.CompetencyLeadership,
					Confidence: 0.7,
				}

				score := defaultGen.calculateRoleScore(bullet, "senior_ic")
				expectedMin := roleScoreBase + roleScoreSecondaryCategoryBoost
				expectedMax := roleScoreBase + roleScorePrimaryCategoryBoost
				Expect(score).To(BeNumerically(">=", expectedMin))
				// But less than primary boost
				Expect(score).To(BeNumerically("<", expectedMax))
			})
		})

		Context("end-to-end role differentiation", func() {
			It("should produce different rankings for different roles", func() {
				events := []*career.CareerEvent{
					fixtures.EventWithCategories("1", "Built data pipeline", []string{"technical"}),
					fixtures.EventWithCategories("2", "Led architecture redesign", []string{"leadership"}),
					fixtures.EventWithCategories("3", "Mentored junior engineers", []string{"mentoring"}),
				}

				seniorBullets, err := generator.GenerateBullets(ctx, events, nil, nil, "senior_ic", "")
				Expect(err).NotTo(HaveOccurred())

				principalBullets, err := generator.GenerateBullets(ctx, events, nil, nil, "principal", "")
				Expect(err).NotTo(HaveOccurred())

				// Technical bullet should rank higher for senior_ic
				Expect(seniorBullets[0].Category).To(Equal(constants.CompetencyTechnical))

				// Leadership bullet should rank higher for principal
				Expect(principalBullets[0].Category).To(Equal(constants.CompetencyLeadership))
			})
		})
	})
})
