//nolint:errcheck // Test file - error handling for test setup is not relevant.
package cv

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/constants"
	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BulletGenerator", func() {
	var (
		generator BulletGenerator
		log       *logger.Logger
		ctx       context.Context
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		generator = NewBulletGenerator(log, nil) // nil uses default scoring config
		ctx = context.Background()
	})

	It("should return empty list for no input", func() {
		bullets, err := generator.GenerateBullets(ctx, nil, nil, nil, "principal", "hiring_manager")
		Expect(err).NotTo(HaveOccurred())
		Expect(bullets).To(BeEmpty())
	})

	It("should generate bullets from events", func() {
		events := []*career.Event{
			fixtures.EventWith("e1", "Led team to deliver microservices", "", ""),
		}

		bullets, err := generator.GenerateBullets(ctx, events, nil, nil, "principal", "hiring_manager")
		Expect(err).NotTo(HaveOccurred())
		Expect(bullets).ToNot(BeEmpty())
	})

	It("should filter by role confidence", func() {
		bullets := []*Bullet{
			{ID: "b1", Confidence: 0.85},
			{ID: "b2", Confidence: 0.70},
		}

		filtered := generator.FilterByRole(bullets, "principal")
		Expect(filtered).To(HaveLen(1))
	})

	It("should rank bullets by score", func() {
		bullets := []*Bullet{
			{ID: "b1", RoleScore: 0.8, AudienceScore: 0.7, MetricScore: 0.6, ImpactScore: 0.7, Confidence: 0.8},
			{ID: "b2", RoleScore: 0.5, AudienceScore: 0.5, MetricScore: 0.5, ImpactScore: 0.5, Confidence: 0.5},
		}

		ranked := generator.RankByRelevance(bullets, "principal", "hiring_manager")
		Expect(ranked[0].ID).To(Equal("b1"))
		Expect(ranked[1].ID).To(Equal("b2"))
	})

	It("should enhance bullet wording", func() {
		bullet := &Bullet{
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
					fixtures.FactWithCategories("fact1", "Delivered 40% cost reduction through system optimization", "event1",
						[]string{"technical"}, []string{"hiring_manager"}),
					fixtures.FactWithCategories("fact2", "Implemented complex distributed caching algorithm", "event2",
						[]string{"technical"}, []string{"peer"}),
				}

				bullets, err := generator.GenerateBullets(ctx, []*career.Event{}, facts, nil, "principal", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				// Should only include fact1 (relevant to hiring_manager)
				// fact2 is only relevant to peer, should be excluded
				Expect(bullets).To(HaveLen(1))
				Expect(bullets[0].Text).To(ContainSubstring("cost reduction"))
			})

			It("should include facts relevant to recruiter audience", func() {
				facts := []*career.Fact{
					fixtures.FactWithCategories("fact1", "Senior engineer with expertise in Go, Ruby, and Python", "event1",
						[]string{"technical"}, []string{"recruiter"}),
					fixtures.FactWithCategories("fact2", "Deep expertise in consensus algorithms and CAP theorem", "event2",
						[]string{"technical"}, []string{"peer"}),
				}

				bullets, err := generator.GenerateBullets(ctx, []*career.Event{}, facts, nil, "principal", "recruiter")
				Expect(err).NotTo(HaveOccurred())
				// Should only include fact1 (relevant to recruiter)
				Expect(bullets).To(HaveLen(1))
				Expect(bullets[0].Text).To(ContainSubstring("Go, Ruby"))
			})

			It("should include facts relevant to peer audience", func() {
				facts := []*career.Fact{
					fixtures.FactWithCategories("fact1", "Designed novel approach to distributed consensus using Raft", "event1",
						[]string{"technical"}, []string{"peer"}),
					fixtures.FactWithCategories("fact2", "Reduced operational costs by 30%", "event2",
						[]string{"technical"}, []string{"hiring_manager"}),
				}

				bullets, err := generator.GenerateBullets(ctx, []*career.Event{}, facts, nil, "principal", "peer")
				Expect(err).NotTo(HaveOccurred())
				// Should only include fact1 (relevant to peer)
				Expect(bullets).To(HaveLen(1))
				Expect(bullets[0].Text).To(ContainSubstring("Raft"))
			})

			It("should return all facts when audience is empty", func() {
				facts := []*career.Fact{
					fixtures.FactWithCategories("fact1", "Business impact achievement", "event1",
						[]string{"technical"}, []string{"hiring_manager"}),
					fixtures.FactWithCategories("fact2", "Technical depth achievement", "event2",
						[]string{"technical"}, []string{"peer"}),
				}

				bullets, err := generator.GenerateBullets(ctx, []*career.Event{}, facts, nil, "principal", "")
				Expect(err).NotTo(HaveOccurred())
				// Should include both facts when no audience filter
				Expect(bullets).To(HaveLen(2))
			})

			It("should exclude facts not relevant to selected audience", func() {
				facts := []*career.Fact{
					fixtures.FactWithCategories("fact1", "Technical implementation detail for peers", "event1",
						[]string{"technical"}, []string{"peer"}),
				}

				bullets, err := generator.GenerateBullets(ctx, []*career.Event{}, facts, nil, "principal", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				// fact1 is only relevant to peer, not hiring_manager
				Expect(bullets).To(BeEmpty())
			})

			It("should include facts with multiple audience relevance", func() {
				facts := []*career.Fact{
					fixtures.FactWithCategories("fact1", "Led migration that reduced costs and improved architecture", "event1",
						[]string{"technical"}, []string{"hiring_manager", "peer"}),
				}

				// Should be included for hiring_manager
				bulletsMgr, err := generator.GenerateBullets(ctx, []*career.Event{}, facts, nil, "principal", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				Expect(bulletsMgr).To(HaveLen(1))

				// Should also be included for peer
				bulletsPeer, err := generator.GenerateBullets(ctx, []*career.Event{}, facts, nil, "principal", "peer")
				Expect(err).NotTo(HaveOccurred())
				Expect(bulletsPeer).To(HaveLen(1))
			})
		})
	})
})

// Task 44: Conversion helper tests.
var _ = Describe("Bullet Conversion", func() {
	Describe("ToCVBullet", func() {
		It("should convert all fields correctly", func() {
			enhanced := &Bullet{
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
			enhanced := &Bullet{
				ID:           "bullet-2",
				Text:         "Original bullet text",
				EnhancedText: "",
			}

			cvBullet := enhanced.ToCVBullet()

			Expect(cvBullet.Text).To(Equal("Original bullet text"))
			Expect(cvBullet.EnhancedText).To(Equal(""))
		})

		It("should use EnhancedText as Text when available", func() {
			enhanced := &Bullet{
				ID:           "bullet-3",
				Text:         "Original text",
				EnhancedText: "Enhanced and improved text",
			}

			cvBullet := enhanced.ToCVBullet()

			Expect(cvBullet.Text).To(Equal("Enhanced and improved text"))
		})

		It("should convert AudienceRelevance slice to map with 1.0 scores", func() {
			enhanced := &Bullet{
				ID:                "bullet-4",
				Text:              "Led cross-functional initiative",
				AudienceRelevance: []string{"recruiter", "hiring_manager"},
			}

			cvBullet := enhanced.ToCVBullet()

			Expect(cvBullet.AudienceRelevance).To(Equal(map[string]float64{
				"recruiter":      1.0,
				"hiring_manager": 1.0,
			}))
		})

		It("should produce nil AudienceRelevance when source slice is nil", func() {
			enhanced := &Bullet{
				ID:                "bullet-5",
				Text:              "Some achievement",
				AudienceRelevance: nil,
			}

			cvBullet := enhanced.ToCVBullet()

			Expect(cvBullet.AudienceRelevance).To(BeNil())
		})

		It("should produce nil AudienceRelevance when source slice is empty", func() {
			enhanced := &Bullet{
				ID:                "bullet-6",
				Text:              "Another achievement",
				AudienceRelevance: []string{},
			}

			cvBullet := enhanced.ToCVBullet()

			Expect(cvBullet.AudienceRelevance).To(BeNil())
		})
	})

	Describe("ConvertBullets", func() {
		It("should convert bullets to domain CVBullets", func() {
			enhanced := []*Bullet{
				{ID: "b1", Text: "Bullet 1", Confidence: 0.8},
				{ID: "b2", Text: "Bullet 2", Confidence: 0.9},
				{ID: "b3", Text: "Bullet 3", Confidence: 0.7},
			}

			cvBullets := ConvertBullets(enhanced)

			Expect(cvBullets).To(HaveLen(3))
			Expect(cvBullets[0].ID).To(Equal("b1"))
			Expect(cvBullets[1].ID).To(Equal("b2"))
			Expect(cvBullets[2].ID).To(Equal("b3"))
		})

		It("should return nil for nil input", func() {
			cvBullets := ConvertBullets(nil)
			Expect(cvBullets).To(BeNil())
		})

		It("should return empty slice for empty input", func() {
			cvBullets := ConvertBullets([]*Bullet{})
			Expect(cvBullets).To(BeEmpty())
		})
	})
})

// BUG-008: Role-based CV differentiation tests.
var _ = Describe("BUG-008: Role-based scoring", func() {
	var (
		generator BulletGenerator
		log       *logger.Logger
		ctx       context.Context
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		generator = NewBulletGenerator(log, nil) // nil uses default scoring config
		ctx = context.Background()
	})

	Describe("Category propagation", func() {
		Context("from Event to Bullet", func() {
			It("should propagate primary category from event", func() {
				events := []*career.Event{
					fixtures.EventWithCategories("evt-1", "Led team migration to Kubernetes", []string{"leadership", "technical"}),
				}

				bullets, err := generator.GenerateBullets(ctx, events, nil, nil, "principal", "")
				Expect(err).NotTo(HaveOccurred())
				Expect(bullets).ToNot(BeEmpty())
				// Primary category (first in list) should be propagated
				Expect(bullets[0].Category).To(Equal(constants.CompetencyLeadership))
			})

			It("should propagate technical category from event", func() {
				events := []*career.Event{
					fixtures.EventWithCategories("evt-2", "Built real-time data pipeline", []string{"technical"}),
				}

				bullets, err := generator.GenerateBullets(ctx, events, nil, nil, "senior_ic", "")
				Expect(err).NotTo(HaveOccurred())
				Expect(bullets).ToNot(BeEmpty())
				Expect(bullets[0].Category).To(Equal(constants.CompetencyTechnical))
			})
		})

		Context("from Fact to Bullet", func() {
			It("should propagate primary competency category from fact", func() {
				facts := []*career.Fact{
					fixtures.FactWithCategories("fact-1", "Mentored 4 junior engineers", "evt-1",
						[]string{"mentoring", "leadership"}, []string{"hiring_manager"}),
				}

				bullets, err := generator.GenerateBullets(ctx, nil, facts, nil, "em", "hiring_manager")
				Expect(err).NotTo(HaveOccurred())
				Expect(bullets).ToNot(BeEmpty())
				Expect(bullets[0].Category).To(Equal(constants.CompetencyMentoring))
			})
		})

		Context("ToCVBullet conversion", func() {
			It("should include category in CVBullet", func() {
				enhanced := &Bullet{
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
		var defaultGen *DefaultBulletGenerator

		BeforeEach(func() {
			defaultGen = generator.(*DefaultBulletGenerator)
		})

		Context("calculateRoleScore with categories", func() {
			It("should score leadership bullets higher for principal than senior_ic", func() {
				bullet := &Bullet{
					Category:   constants.CompetencyLeadership,
					Confidence: 0.7,
				}

				principalScore := defaultGen.calculateRoleScore(bullet, "principal")
				seniorScore := defaultGen.calculateRoleScore(bullet, "senior_ic")

				Expect(principalScore).To(BeNumerically(">", seniorScore))
			})

			It("should score technical bullets higher for senior_ic than principal", func() {
				bullet := &Bullet{
					Category:   constants.CompetencyTechnical,
					Confidence: 0.7,
				}

				seniorScore := defaultGen.calculateRoleScore(bullet, "senior_ic")
				principalScore := defaultGen.calculateRoleScore(bullet, "principal")

				Expect(seniorScore).To(BeNumerically(">", principalScore))
			})

			It("should score mentoring bullets higher for em than senior_ic", func() {
				bullet := &Bullet{
					Category:   constants.CompetencyMentoring,
					Confidence: 0.7,
				}

				emScore := defaultGen.calculateRoleScore(bullet, "em")
				seniorScore := defaultGen.calculateRoleScore(bullet, "senior_ic")

				Expect(emScore).To(BeNumerically(">", seniorScore))
			})

			It("should give primary category bullets a strong boost", func() {
				// Technical is primary for senior_ic
				bullet := &Bullet{
					Category:   constants.CompetencyTechnical,
					Confidence: 0.7,
				}

				score := defaultGen.calculateRoleScore(bullet, "senior_ic")
				expectedMin := roleScoreBase + roleScorePrimaryCategoryBoost
				Expect(score).To(BeNumerically(">=", expectedMin))
			})

			It("should give secondary category bullets a medium boost", func() {
				// Leadership is secondary for senior_ic
				bullet := &Bullet{
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
				events := []*career.Event{
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

// BUG-008: ScoringConfig integration tests.
var _ = Describe("BUG-008: ScoringConfig integration", func() {
	var log *logger.Logger

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
	})

	Describe("getWeights", func() {
		Context("with nil config", func() {
			It("should return default weights", func() {
				gen := NewBulletGenerator(log, nil).(*DefaultBulletGenerator)

				weights := gen.getWeights()

				Expect(weights.RoleScore).To(Equal(0.25))
				Expect(weights.AudienceScore).To(Equal(0.20))
				Expect(weights.MetricScore).To(Equal(0.20))
				Expect(weights.ImpactScore).To(Equal(0.20))
				Expect(weights.Confidence).To(Equal(0.15))
			})

			It("should return weights that sum to 1.0", func() {
				gen := NewBulletGenerator(log, nil).(*DefaultBulletGenerator)

				weights := gen.getWeights()
				sum := weights.RoleScore + weights.AudienceScore + weights.MetricScore +
					weights.ImpactScore + weights.Confidence

				Expect(sum).To(BeNumerically("~", 1.0, 0.001))
			})
		})

		Context("with config provided", func() {
			It("should return config weights", func() {
				scoringCfg := &config.ScoringConfig{
					Weights: config.ScoringWeights{
						RoleScore:     0.30,
						AudienceScore: 0.25,
						MetricScore:   0.15,
						ImpactScore:   0.15,
						Confidence:    0.15,
					},
				}
				gen := NewBulletGenerator(log, scoringCfg).(*DefaultBulletGenerator)

				weights := gen.getWeights()

				Expect(weights.RoleScore).To(Equal(0.30))
				Expect(weights.AudienceScore).To(Equal(0.25))
				Expect(weights.MetricScore).To(Equal(0.15))
				Expect(weights.ImpactScore).To(Equal(0.15))
				Expect(weights.Confidence).To(Equal(0.15))
			})

			It("should handle zero weights", func() {
				scoringCfg := &config.ScoringConfig{
					Weights: config.ScoringWeights{
						RoleScore:     0.0,
						AudienceScore: 0.0,
						MetricScore:   0.0,
						ImpactScore:   0.0,
						Confidence:    1.0, // Only confidence matters
					},
				}
				gen := NewBulletGenerator(log, scoringCfg).(*DefaultBulletGenerator)

				weights := gen.getWeights()

				Expect(weights.RoleScore).To(Equal(0.0))
				Expect(weights.Confidence).To(Equal(1.0))
			})

			It("should handle weights that sum to more than 1.0", func() {
				// Config allows weights > 1.0 (validation is separate)
				scoringCfg := &config.ScoringConfig{
					Weights: config.ScoringWeights{
						RoleScore:     0.50,
						AudienceScore: 0.50,
						MetricScore:   0.50,
						ImpactScore:   0.50,
						Confidence:    0.50,
					},
				}
				gen := NewBulletGenerator(log, scoringCfg).(*DefaultBulletGenerator)

				weights := gen.getWeights()

				// Should return as-is (validation is done by config.ValidateWeights)
				Expect(weights.RoleScore).To(Equal(0.50))
			})

			It("should handle partial zero weights", func() {
				scoringCfg := &config.ScoringConfig{
					Weights: config.ScoringWeights{
						RoleScore:     0.50,
						AudienceScore: 0.0, // Zero
						MetricScore:   0.25,
						ImpactScore:   0.0, // Zero
						Confidence:    0.25,
					},
				}
				gen := NewBulletGenerator(log, scoringCfg).(*DefaultBulletGenerator)

				weights := gen.getWeights()

				Expect(weights.AudienceScore).To(Equal(0.0))
				Expect(weights.ImpactScore).To(Equal(0.0))
			})
		})
	})

	Describe("getMinConfidenceForRole", func() {
		Context("with nil config", func() {
			It("should return default for principal", func() {
				gen := NewBulletGenerator(log, nil).(*DefaultBulletGenerator)
				Expect(gen.getMinConfidenceForRole("principal")).To(Equal(0.80))
			})

			It("should return default for staff", func() {
				gen := NewBulletGenerator(log, nil).(*DefaultBulletGenerator)
				Expect(gen.getMinConfidenceForRole("staff")).To(Equal(0.75))
			})

			It("should return default for em", func() {
				gen := NewBulletGenerator(log, nil).(*DefaultBulletGenerator)
				Expect(gen.getMinConfidenceForRole("em")).To(Equal(0.75))
			})

			It("should return default for senior_ic", func() {
				gen := NewBulletGenerator(log, nil).(*DefaultBulletGenerator)
				Expect(gen.getMinConfidenceForRole("senior_ic")).To(Equal(0.75))
			})

			It("should return lowest default for unknown role", func() {
				gen := NewBulletGenerator(log, nil).(*DefaultBulletGenerator)
				Expect(gen.getMinConfidenceForRole("unknown")).To(Equal(0.70))
			})

			It("should return lowest default for empty role string", func() {
				gen := NewBulletGenerator(log, nil).(*DefaultBulletGenerator)
				Expect(gen.getMinConfidenceForRole("")).To(Equal(0.70))
			})
		})

		Context("with config provided", func() {
			It("should return config values for known roles", func() {
				scoringCfg := &config.ScoringConfig{
					RoleSettings: map[string]config.RoleScoringCfg{
						"principal": {MinConfidence: 0.90},
						"staff":     {MinConfidence: 0.85},
						"em":        {MinConfidence: 0.85},
						"senior_ic": {MinConfidence: 0.80},
					},
				}
				gen := NewBulletGenerator(log, scoringCfg).(*DefaultBulletGenerator)

				Expect(gen.getMinConfidenceForRole("principal")).To(Equal(0.90))
				Expect(gen.getMinConfidenceForRole("staff")).To(Equal(0.85))
				Expect(gen.getMinConfidenceForRole("em")).To(Equal(0.85))
				Expect(gen.getMinConfidenceForRole("senior_ic")).To(Equal(0.80))
			})

			It("should fall back to default for role not in config", func() {
				scoringCfg := &config.ScoringConfig{
					RoleSettings: map[string]config.RoleScoringCfg{
						"principal": {MinConfidence: 0.90},
					},
				}
				gen := NewBulletGenerator(log, scoringCfg).(*DefaultBulletGenerator)

				// staff not in config, should use default
				Expect(gen.getMinConfidenceForRole("staff")).To(Equal(0.75))
			})

			It("should handle empty RoleSettings map", func() {
				scoringCfg := &config.ScoringConfig{
					RoleSettings: map[string]config.RoleScoringCfg{},
				}
				gen := NewBulletGenerator(log, scoringCfg).(*DefaultBulletGenerator)

				// All should use defaults
				Expect(gen.getMinConfidenceForRole("principal")).To(Equal(0.80))
				Expect(gen.getMinConfidenceForRole("staff")).To(Equal(0.75))
			})

			It("should handle nil RoleSettings map", func() {
				scoringCfg := &config.ScoringConfig{
					RoleSettings: nil,
				}
				gen := NewBulletGenerator(log, scoringCfg).(*DefaultBulletGenerator)

				Expect(gen.getMinConfidenceForRole("principal")).To(Equal(0.80))
			})
		})

		Context("role name case sensitivity", func() {
			It("should normalize uppercase role to lowercase", func() {
				gen := NewBulletGenerator(log, nil).(*DefaultBulletGenerator)
				Expect(gen.getMinConfidenceForRole("PRINCIPAL")).To(Equal(0.80))
			})

			It("should normalize mixed case role to lowercase", func() {
				gen := NewBulletGenerator(log, nil).(*DefaultBulletGenerator)
				Expect(gen.getMinConfidenceForRole("Principal")).To(Equal(0.80))
			})

			It("should match config with lowercase lookup", func() {
				scoringCfg := &config.ScoringConfig{
					RoleSettings: map[string]config.RoleScoringCfg{
						"principal": {MinConfidence: 0.95},
					},
				}
				gen := NewBulletGenerator(log, scoringCfg).(*DefaultBulletGenerator)

				// Uppercase input should match lowercase key
				Expect(gen.getMinConfidenceForRole("PRINCIPAL")).To(Equal(0.95))
			})
		})

		Context("edge cases", func() {
			It("should handle zero MinConfidence in config", func() {
				scoringCfg := &config.ScoringConfig{
					RoleSettings: map[string]config.RoleScoringCfg{
						"principal": {MinConfidence: 0.0},
					},
				}
				gen := NewBulletGenerator(log, scoringCfg).(*DefaultBulletGenerator)

				Expect(gen.getMinConfidenceForRole("principal")).To(Equal(0.0))
			})

			It("should handle MinConfidence of 1.0 in config", func() {
				scoringCfg := &config.ScoringConfig{
					RoleSettings: map[string]config.RoleScoringCfg{
						"principal": {MinConfidence: 1.0},
					},
				}
				gen := NewBulletGenerator(log, scoringCfg).(*DefaultBulletGenerator)

				Expect(gen.getMinConfidenceForRole("principal")).To(Equal(1.0))
			})
		})
	})

	Describe("calculateFinalScore", func() {
		Context("with config weights", func() {
			It("should use config weights in calculation", func() {
				scoringCfg := &config.ScoringConfig{
					Weights: config.ScoringWeights{
						RoleScore:     0.50, // High weight
						AudienceScore: 0.10,
						MetricScore:   0.10,
						ImpactScore:   0.10,
						Confidence:    0.20,
					},
				}
				gen := NewBulletGenerator(log, scoringCfg).(*DefaultBulletGenerator)

				bullet := &Bullet{
					RoleScore:     1.0,
					AudienceScore: 0.5,
					MetricScore:   0.5,
					ImpactScore:   0.5,
					Confidence:    0.5,
				}

				score := gen.calculateFinalScore(bullet)

				// Expected: 0.50*1.0 + 0.10*0.5 + 0.10*0.5 + 0.10*0.5 + 0.20*0.5 = 0.75
				Expect(score).To(BeNumerically("~", 0.75, 0.01))
			})

			It("should use default weights when config is nil", func() {
				gen := NewBulletGenerator(log, nil).(*DefaultBulletGenerator)

				bullet := &Bullet{
					RoleScore:     1.0,
					AudienceScore: 1.0,
					MetricScore:   1.0,
					ImpactScore:   1.0,
					Confidence:    1.0,
				}

				score := gen.calculateFinalScore(bullet)

				// Expected: 0.25 + 0.20 + 0.20 + 0.20 + 0.15 = 1.0
				Expect(score).To(BeNumerically("~", 1.0, 0.01))
			})
		})

		Context("boundary conditions", func() {
			It("should return 0 when all scores are 0", func() {
				gen := NewBulletGenerator(log, nil).(*DefaultBulletGenerator)

				bullet := &Bullet{
					RoleScore:     0.0,
					AudienceScore: 0.0,
					MetricScore:   0.0,
					ImpactScore:   0.0,
					Confidence:    0.0,
				}

				score := gen.calculateFinalScore(bullet)

				Expect(score).To(Equal(0.0))
			})

			It("should cap score at 1.0 when weighted sum exceeds 1.0", func() {
				// Use weights that sum to > 1.0 with max scores
				scoringCfg := &config.ScoringConfig{
					Weights: config.ScoringWeights{
						RoleScore:     0.50,
						AudienceScore: 0.50,
						MetricScore:   0.50,
						ImpactScore:   0.50,
						Confidence:    0.50,
					},
				}
				gen := NewBulletGenerator(log, scoringCfg).(*DefaultBulletGenerator)

				bullet := &Bullet{
					RoleScore:     1.0,
					AudienceScore: 1.0,
					MetricScore:   1.0,
					ImpactScore:   1.0,
					Confidence:    1.0,
				}

				score := gen.calculateFinalScore(bullet)

				// Should be capped at 1.0
				Expect(score).To(Equal(1.0))
			})

			It("should handle scores at threshold values", func() {
				gen := NewBulletGenerator(log, nil).(*DefaultBulletGenerator)

				bullet := &Bullet{
					RoleScore:     0.5,
					AudienceScore: 0.5,
					MetricScore:   0.5,
					ImpactScore:   0.5,
					Confidence:    0.5,
				}

				score := gen.calculateFinalScore(bullet)

				// Expected: 0.25*0.5 + 0.20*0.5 + 0.20*0.5 + 0.20*0.5 + 0.15*0.5 = 0.5
				Expect(score).To(BeNumerically("~", 0.5, 0.01))
			})
		})

		Context("with zero weights", func() {
			It("should ignore score components with zero weight", func() {
				scoringCfg := &config.ScoringConfig{
					Weights: config.ScoringWeights{
						RoleScore:     0.0, // Ignored
						AudienceScore: 0.0, // Ignored
						MetricScore:   0.0, // Ignored
						ImpactScore:   0.0, // Ignored
						Confidence:    1.0, // Only this matters
					},
				}
				gen := NewBulletGenerator(log, scoringCfg).(*DefaultBulletGenerator)

				bullet := &Bullet{
					RoleScore:     1.0, // Should be ignored
					AudienceScore: 1.0, // Should be ignored
					MetricScore:   1.0, // Should be ignored
					ImpactScore:   1.0, // Should be ignored
					Confidence:    0.8,
				}

				score := gen.calculateFinalScore(bullet)

				// Only confidence contributes: 1.0 * 0.8 = 0.8
				Expect(score).To(BeNumerically("~", 0.8, 0.01))
			})
		})

		Context("precision and rounding", func() {
			It("should handle fractional weights correctly", func() {
				scoringCfg := &config.ScoringConfig{
					Weights: config.ScoringWeights{
						RoleScore:     0.333,
						AudienceScore: 0.333,
						MetricScore:   0.334,
						ImpactScore:   0.0,
						Confidence:    0.0,
					},
				}
				gen := NewBulletGenerator(log, scoringCfg).(*DefaultBulletGenerator)

				bullet := &Bullet{
					RoleScore:     1.0,
					AudienceScore: 1.0,
					MetricScore:   1.0,
					ImpactScore:   0.0,
					Confidence:    0.0,
				}

				score := gen.calculateFinalScore(bullet)

				// Expected: 0.333 + 0.333 + 0.334 = 1.0
				Expect(score).To(BeNumerically("~", 1.0, 0.01))
			})
		})
	})

	Describe("FilterByRole with config MinConfidence", func() {
		Context("basic filtering", func() {
			It("should filter bullets below threshold", func() {
				scoringCfg := &config.ScoringConfig{
					RoleSettings: map[string]config.RoleScoringCfg{
						"principal": {MinConfidence: 0.85},
					},
				}
				gen := NewBulletGenerator(log, scoringCfg)

				bullets := []*Bullet{
					{ID: "b1", Confidence: 0.90}, // Above
					{ID: "b2", Confidence: 0.80}, // Below
					{ID: "b3", Confidence: 0.85}, // At threshold (passes)
				}

				filtered := gen.FilterByRole(bullets, "principal")

				Expect(filtered).To(HaveLen(2))
				Expect(filtered[0].ID).To(Equal("b1"))
				Expect(filtered[1].ID).To(Equal("b3"))
			})

			It("should use default threshold when config is nil", func() {
				gen := NewBulletGenerator(log, nil)

				bullets := []*Bullet{
					{ID: "b1", Confidence: 0.85}, // Above 0.80 default
					{ID: "b2", Confidence: 0.75}, // Below 0.80 default
					{ID: "b3", Confidence: 0.80}, // At threshold
				}

				filtered := gen.FilterByRole(bullets, "principal")

				Expect(filtered).To(HaveLen(2))
			})
		})

		Context("edge cases", func() {
			It("should return empty slice for empty input", func() {
				gen := NewBulletGenerator(log, nil)

				filtered := gen.FilterByRole([]*Bullet{}, "principal")

				Expect(filtered).To(BeEmpty())
			})

			It("should return empty slice when all bullets below threshold", func() {
				scoringCfg := &config.ScoringConfig{
					RoleSettings: map[string]config.RoleScoringCfg{
						"principal": {MinConfidence: 0.95},
					},
				}
				gen := NewBulletGenerator(log, scoringCfg)

				bullets := []*Bullet{
					{ID: "b1", Confidence: 0.90},
					{ID: "b2", Confidence: 0.80},
					{ID: "b3", Confidence: 0.70},
				}

				filtered := gen.FilterByRole(bullets, "principal")

				Expect(filtered).To(BeEmpty())
			})

			It("should return all bullets when all above threshold", func() {
				scoringCfg := &config.ScoringConfig{
					RoleSettings: map[string]config.RoleScoringCfg{
						"principal": {MinConfidence: 0.50},
					},
				}
				gen := NewBulletGenerator(log, scoringCfg)

				bullets := []*Bullet{
					{ID: "b1", Confidence: 0.90},
					{ID: "b2", Confidence: 0.80},
					{ID: "b3", Confidence: 0.70},
				}

				filtered := gen.FilterByRole(bullets, "principal")

				Expect(filtered).To(HaveLen(3))
			})

			It("should handle MinConfidence of 0 (all pass)", func() {
				scoringCfg := &config.ScoringConfig{
					RoleSettings: map[string]config.RoleScoringCfg{
						"principal": {MinConfidence: 0.0},
					},
				}
				gen := NewBulletGenerator(log, scoringCfg)

				bullets := []*Bullet{
					{ID: "b1", Confidence: 0.10},
					{ID: "b2", Confidence: 0.01},
					{ID: "b3", Confidence: 0.0},
				}

				filtered := gen.FilterByRole(bullets, "principal")

				Expect(filtered).To(HaveLen(3))
			})

			It("should handle MinConfidence of 1.0 (only perfect pass)", func() {
				scoringCfg := &config.ScoringConfig{
					RoleSettings: map[string]config.RoleScoringCfg{
						"principal": {MinConfidence: 1.0},
					},
				}
				gen := NewBulletGenerator(log, scoringCfg)

				bullets := []*Bullet{
					{ID: "b1", Confidence: 1.0},  // Exactly 1.0
					{ID: "b2", Confidence: 0.99}, // Just under
					{ID: "b3", Confidence: 0.999},
				}

				filtered := gen.FilterByRole(bullets, "principal")

				Expect(filtered).To(HaveLen(1))
				Expect(filtered[0].ID).To(Equal("b1"))
			})

			It("should handle nil input slice", func() {
				gen := NewBulletGenerator(log, nil)

				filtered := gen.FilterByRole(nil, "principal")

				Expect(filtered).To(BeNil())
			})
		})

		Context("role name handling", func() {
			It("should handle uppercase role name", func() {
				scoringCfg := &config.ScoringConfig{
					RoleSettings: map[string]config.RoleScoringCfg{
						"principal": {MinConfidence: 0.90},
					},
				}
				gen := NewBulletGenerator(log, scoringCfg)

				bullets := []*Bullet{
					{ID: "b1", Confidence: 0.95},
					{ID: "b2", Confidence: 0.85},
				}

				filtered := gen.FilterByRole(bullets, "PRINCIPAL")

				Expect(filtered).To(HaveLen(1))
			})

			It("should use default for unknown role", func() {
				scoringCfg := &config.ScoringConfig{
					RoleSettings: map[string]config.RoleScoringCfg{
						"principal": {MinConfidence: 0.90},
					},
				}
				gen := NewBulletGenerator(log, scoringCfg)

				bullets := []*Bullet{
					{ID: "b1", Confidence: 0.75},
					{ID: "b2", Confidence: 0.65},
				}

				// Unknown role uses default 0.70
				filtered := gen.FilterByRole(bullets, "unknown_role")

				Expect(filtered).To(HaveLen(1))
				Expect(filtered[0].ID).To(Equal("b1"))
			})

			It("should use lowest default for empty role string", func() {
				gen := NewBulletGenerator(log, nil)

				bullets := []*Bullet{
					{ID: "b1", Confidence: 0.75},
					{ID: "b2", Confidence: 0.65},
				}

				// Empty role uses default 0.70
				filtered := gen.FilterByRole(bullets, "")

				Expect(filtered).To(HaveLen(1))
			})
		})
	})

	Describe("Bullet creation from Achievement with Category", func() {
		It("should propagate Category from Achievement to Bullet", func() {
			gen := NewBulletGenerator(log, nil).(*DefaultBulletGenerator)
			ctx := context.Background()

			achievements := []*Achievement{
				{
					ID:          "ach-1",
					Description: "Led team migration",
					EventID:     "evt-1",
					Confidence:  0.9,
					Category:    constants.CompetencyLeadership,
				},
			}

			bullets, err := gen.GenerateBullets(ctx, nil, nil, achievements, "principal", "")
			Expect(err).NotTo(HaveOccurred())
			Expect(bullets).To(HaveLen(1))
			Expect(bullets[0].Category).To(Equal(constants.CompetencyLeadership))
		})

		It("should handle Achievement with empty Category", func() {
			gen := NewBulletGenerator(log, nil).(*DefaultBulletGenerator)
			ctx := context.Background()

			achievements := []*Achievement{
				{
					ID:          "ach-1",
					Description: "Did something",
					EventID:     "evt-1",
					Confidence:  0.9,
					Category:    "", // Empty
				},
			}

			bullets, err := gen.GenerateBullets(ctx, nil, nil, achievements, "principal", "")
			Expect(err).NotTo(HaveOccurred())
			Expect(bullets).To(HaveLen(1))
			Expect(bullets[0].Category).To(Equal(constants.CompetencyCategory("")))
		})

		It("should propagate different Categories for multiple Achievements", func() {
			gen := NewBulletGenerator(log, nil).(*DefaultBulletGenerator)
			ctx := context.Background()

			achievements := []*Achievement{
				{
					ID:          "ach-1",
					Description: "Led team",
					EventID:     "evt-1",
					Confidence:  0.9,
					Category:    constants.CompetencyLeadership,
				},
				{
					ID:          "ach-2",
					Description: "Built system",
					EventID:     "evt-2",
					Confidence:  0.85,
					Category:    constants.CompetencyTechnical,
				},
			}

			bullets, err := gen.GenerateBullets(ctx, nil, nil, achievements, "principal", "")
			Expect(err).NotTo(HaveOccurred())
			Expect(bullets).To(HaveLen(2))

			// Find bullets by their text to check categories (order may vary due to scoring)
			var leadershipBullet, technicalBullet *Bullet
			for _, b := range bullets {
				if strings.Contains(b.Text, "Led team") {
					leadershipBullet = b
				}
				if strings.Contains(b.Text, "Built system") {
					technicalBullet = b
				}
			}

			Expect(leadershipBullet).NotTo(BeNil())
			Expect(leadershipBullet.Category).To(Equal(constants.CompetencyLeadership))
			Expect(technicalBullet).NotTo(BeNil())
			Expect(technicalBullet.Category).To(Equal(constants.CompetencyTechnical))
		})
	})
})

// BUG-013: Bullet deduplication must be company-aware.
var _ = Describe("BUG-013: Company-aware bullet deduplication", func() {
	var (
		generator BulletGenerator
		log       *logger.Logger
		ctx       context.Context
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		generator = NewBulletGenerator(log, nil)
		ctx = context.Background()
	})

	It("should NOT merge bullets with identical text from different companies", func() {
		events := []*career.Event{
			fixtures.EventWith("e-friday", "Acted as senior stabilising engineer during late-stage delivery pressure", "We Are Friday", ""),
			fixtures.EventWith("e-beis", "Acted as senior stabilising engineer during late-stage delivery pressure", "BEIS", ""),
		}

		bullets, err := generator.GenerateBullets(ctx, events, nil, nil, "", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(bullets).To(HaveLen(2), "identical text at different companies must produce separate bullets")
	})

	It("should merge bullets with identical text from the same company", func() {
		events := []*career.Event{
			fixtures.EventWith("e1", "Built scalable backend services", "Acme Corp", "Project A"),
			fixtures.EventWith("e2", "Built scalable backend services", "Acme Corp", "Project B"),
		}

		bullets, err := generator.GenerateBullets(ctx, events, nil, nil, "", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(bullets).To(HaveLen(1), "identical text at the same company should be merged")
	})

	It("should produce separate bullets for cross-cutting entries across many companies", func() {
		companies := []string{"Company A", "Company B", "Company C", "Company D", "Company E"}
		events := make([]*career.Event, len(companies))
		for i, company := range companies {
			events[i] = fixtures.EventWith(
				fmt.Sprintf("e-%d", i+1),
				"Designed and delivered scalable backend services using Ruby on Rails",
				company,
				"",
			)
		}

		bullets, err := generator.GenerateBullets(ctx, events, nil, nil, "", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(bullets).To(HaveLen(len(companies)),
			"cross-cutting entries across %d companies must produce %d separate bullets", len(companies), len(companies))
	})

	It("should keep SourceEventIDs scoped to the same company after dedup", func() {
		events := []*career.Event{
			fixtures.EventWith("e-friday-1", "Led delivery of key features", "We Are Friday", ""),
			fixtures.EventWith("e-friday-2", "Led delivery of key features", "We Are Friday", ""),
			fixtures.EventWith("e-beis-1", "Led delivery of key features", "BEIS", ""),
		}

		bullets, err := generator.GenerateBullets(ctx, events, nil, nil, "", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(bullets).To(HaveLen(2), "should produce one bullet per company")

		// Build a lookup of company per event for verification.
		eventCompany := map[string]string{
			"e-friday-1": "We Are Friday",
			"e-friday-2": "We Are Friday",
			"e-beis-1":   "BEIS",
		}

		for _, bullet := range bullets {
			companies := map[string]bool{}
			for _, eid := range bullet.SourceEventIDs {
				companies[eventCompany[eid]] = true
			}
			Expect(companies).To(HaveLen(1),
				"SourceEventIDs should only reference events from one company, got %v", bullet.SourceEventIDs)
		}
	})

	It("should resolve primary company deterministically when counts are tied", func() {
		// Build an event map where a bullet has source events from two
		// companies with equal counts - a genuine tie scenario.
		eventMap := map[string]*career.Event{
			"e-zebra": fixtures.EventWith("e-zebra", "work", "Zebra Inc", ""),
			"e-alpha": fixtures.EventWith("e-alpha", "work", "Alpha Corp", ""),
		}
		sourceIDs := []string{"e-zebra", "e-alpha"}

		// Run multiple times to verify determinism.
		first := resolvePrimaryCompany(sourceIDs, eventMap)
		Expect(first).To(Equal("Alpha Corp"),
			"lexicographically smallest company should win on tie")

		for i := range 20 {
			result := resolvePrimaryCompany(sourceIDs, eventMap)
			Expect(result).To(Equal(first),
				"iteration %d: tie-breaking must be deterministic", i)
		}
	})

	It("should not merge bullets from events missing from the event map", func() {
		// Events with IDs that won't resolve to a company should not merge
		// with each other under an empty key (BUG-015 defence-in-depth).
		events := []*career.Event{
			fixtures.EventWith("e1", "Built monitoring dashboards", "", "Project X"),
			fixtures.EventWith("e2", "Built monitoring dashboards", "", "Project Y"),
		}

		bullets, err := generator.GenerateBullets(ctx, events, nil, nil, "", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(bullets).To(HaveLen(2),
			"bullets from events with no company should not merge under an empty key")
	})

	It("should not merge orphaned bullets with no SourceEventIDs", func() {
		// Bullets with no SourceEventIDs and no company should use their
		// bullet ID as fallback key to prevent incorrect merging.
		bg := generator.(*DefaultBulletGenerator)
		bullets := []*Bullet{
			{ID: "b1", Text: "Identical orphaned text", SourceEventIDs: nil},
			{ID: "b2", Text: "Identical orphaned text", SourceEventIDs: nil},
		}
		eventMap := map[string]*career.Event{}

		result := bg.deduplicateBullets(bullets, eventMap)
		Expect(result).To(HaveLen(2),
			"orphaned bullets with distinct IDs must not merge under an empty key")
	})
})
