package display_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/constants"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/ui/display"
)

var _ = Describe("CV display types", func() {
	Describe("CVViewFromDomain", func() {
		It("converts a fully populated CV view with nested content", func() {
			generatedAt := time.Date(2024, time.May, 3, 15, 0, 0, 0, time.UTC)

			bullet := fixtures.CVBulletWithScores("bullet-1", "section-1", "Led platform improvements", 0.95, 0.9, 0.8, 0.7)
			bullet.SourceEventIDs = []string{"event-1"}
			bullet.SourceFactIDs = []string{"fact-1"}
			bullet.InclusionReason = "impact"
			bullet.EnhancedText = "Led major platform improvements"
			bullet.Category = constants.CompetencyTechnical
			bullet.MetricScore = 0.6
			bullet.ImpactScore = 0.9
			bullet.ImpactLevel = "high"
			bullet.KeywordMatches = []string{"platform", "leadership"}
			bullet.AudienceRelevance = map[string]float64{"hiring-manager": 0.9}

			contentGroup := fixtures.ContentGroupWith("Acme", "2022", "2024")
			contentGroup.Bullets = []*career.CVBullet{bullet}

			section := fixtures.CVSectionWith("section-1", "cv-1", "experience", "Experience", 1)
			section.Summary = "Selected experience"
			section.Content = []*career.SectionContentGroup{contentGroup}

			view := fixtures.CVViewWith("cv-1", "Staff CV", "staff", "hiring-manager")
			view.EventFilters = map[string]interface{}{"company": "Acme", "year": 2024}
			view.GeneratedAt = generatedAt
			view.SourceEventCount = 5
			view.SourceFactCount = 3
			view.Sections = []*career.CVSection{section}

			result := display.CVViewFromDomain(view)

			Expect(result).To(Equal(display.CVView{
				ID:               "cv-1",
				Name:             "Staff CV",
				TargetRole:       "staff",
				TargetAudience:   "hiring-manager",
				EventFilters:     map[string]interface{}{"company": "Acme", "year": 2024},
				GeneratedAt:      generatedAt,
				SourceEventCount: 5,
				SourceFactCount:  3,
				Sections: []display.CVSection{
					{
						ID:          "section-1",
						CVViewID:    "cv-1",
						SectionType: "experience",
						Title:       "Experience",
						Order:       1,
						Summary:     "Selected experience",
						Content: []display.SectionContentGroup{
							{
								Header:    "Acme",
								StartDate: "2022",
								EndDate:   "2024",
								Bullets: []display.CVBullet{
									{
										ID:                "bullet-1",
										SectionID:         "section-1",
										Text:              "Led platform improvements",
										SourceEventIDs:    []string{"event-1"},
										SourceFactIDs:     []string{"fact-1"},
										Rank:              0.95,
										InclusionReason:   "impact",
										Confidence:        0.9,
										EnhancedText:      "Led major platform improvements",
										Category:          "technical",
										RoleScore:         0.8,
										AudienceScore:     0.7,
										MetricScore:       0.6,
										ImpactScore:       0.9,
										ImpactLevel:       "high",
										KeywordMatches:    []string{"platform", "leadership"},
										AudienceRelevance: map[string]float64{"hiring-manager": 0.9},
									},
								},
							},
						},
					},
				},
			}))
		})

		It("handles nil input gracefully", func() {
			Expect(display.CVViewFromDomain(nil)).To(Equal(display.CVView{}))
		})
	})

	Describe("CVViewsFromDomain", func() {
		It("converts a slice of CV views", func() {
			v1 := fixtures.CVView("cv-1")
			v2 := fixtures.CVView("cv-2")
			views := []*career.CVView{v1, v2}

			result := display.CVViewsFromDomain(views)

			Expect(result).To(HaveLen(2))
			Expect(result[0].ID).To(Equal("cv-1"))
			Expect(result[1].ID).To(Equal("cv-2"))
		})

		It("returns nil for nil input", func() {
			Expect(display.CVViewsFromDomain(nil)).To(BeNil())
		})

		It("returns empty slice for empty input", func() {
			Expect(display.CVViewsFromDomain([]*career.CVView{})).To(Equal([]display.CVView{}))
		})
	})

	Describe("CVSectionFromDomain", func() {
		It("converts a fully populated CV section", func() {
			contentGroup := fixtures.ContentGroup("Acme")
			section := fixtures.CVSectionWith("section-1", "cv-1", "experience", "Experience", 2)
			section.Summary = "Summary"
			section.Content = []*career.SectionContentGroup{contentGroup}

			Expect(display.CVSectionFromDomain(section)).To(Equal(display.CVSection{
				ID:          "section-1",
				CVViewID:    "cv-1",
				SectionType: "experience",
				Title:       "Experience",
				Order:       2,
				Summary:     "Summary",
				Content: []display.SectionContentGroup{
					{Header: "Acme"},
				},
			}))
		})

		It("handles nil input gracefully", func() {
			Expect(display.CVSectionFromDomain(nil)).To(Equal(display.CVSection{}))
		})
	})

	Describe("CVSectionsFromDomain", func() {
		It("returns nil for nil input", func() {
			Expect(display.CVSectionsFromDomain(nil)).To(BeNil())
		})

		It("returns empty slice for empty input", func() {
			Expect(display.CVSectionsFromDomain([]*career.CVSection{})).To(Equal([]display.CVSection{}))
		})
	})

	Describe("SectionContentGroupFromDomain", func() {
		It("converts a fully populated content group", func() {
			bullet := fixtures.CVBullet("bullet-1", "section-1")
			bullet.Text = "Delivered migration"
			group := fixtures.ContentGroupWith("Acme", "2021", "2024")
			group.Bullets = []*career.CVBullet{bullet}

			result := display.SectionContentGroupFromDomain(group)

			Expect(result.Header).To(Equal("Acme"))
			Expect(result.StartDate).To(Equal("2021"))
			Expect(result.EndDate).To(Equal("2024"))
			Expect(result.Bullets).To(HaveLen(1))
			Expect(result.Bullets[0].ID).To(Equal("bullet-1"))
			Expect(result.Bullets[0].Text).To(Equal("Delivered migration"))
		})

		It("handles nil input gracefully", func() {
			Expect(display.SectionContentGroupFromDomain(nil)).To(Equal(display.SectionContentGroup{}))
		})
	})

	Describe("SectionContentGroupsFromDomain", func() {
		It("returns nil for nil input", func() {
			Expect(display.SectionContentGroupsFromDomain(nil)).To(BeNil())
		})

		It("returns empty slice for empty input", func() {
			Expect(display.SectionContentGroupsFromDomain([]*career.SectionContentGroup{})).To(Equal([]display.SectionContentGroup{}))
		})
	})

	Describe("CVBulletFromDomain", func() {
		It("converts a fully populated bullet", func() {
			bullet := fixtures.CVBulletWithScores("bullet-1", "section-1", "Led migration", 0.7, 0.8, 0.4, 0.5)
			bullet.SourceEventIDs = []string{"event-1"}
			bullet.SourceFactIDs = []string{"fact-1"}
			bullet.InclusionReason = "impact"
			bullet.EnhancedText = "Led a major migration"
			bullet.Category = constants.CompetencyLeadership
			bullet.MetricScore = 0.6
			bullet.ImpactScore = 0.9
			bullet.ImpactLevel = "high"
			bullet.KeywordMatches = []string{"migration"}
			bullet.AudienceRelevance = map[string]float64{"technical-peer": 0.75}

			Expect(display.CVBulletFromDomain(bullet)).To(Equal(display.CVBullet{
				ID:                "bullet-1",
				SectionID:         "section-1",
				Text:              "Led migration",
				SourceEventIDs:    []string{"event-1"},
				SourceFactIDs:     []string{"fact-1"},
				Rank:              0.7,
				InclusionReason:   "impact",
				Confidence:        0.8,
				EnhancedText:      "Led a major migration",
				Category:          "leadership",
				RoleScore:         0.4,
				AudienceScore:     0.5,
				MetricScore:       0.6,
				ImpactScore:       0.9,
				ImpactLevel:       "high",
				KeywordMatches:    []string{"migration"},
				AudienceRelevance: map[string]float64{"technical-peer": 0.75},
			}))
		})

		It("handles nil input gracefully", func() {
			Expect(display.CVBulletFromDomain(nil)).To(Equal(display.CVBullet{}))
		})
	})

	Describe("CVBulletsFromDomain", func() {
		It("returns nil for nil input", func() {
			Expect(display.CVBulletsFromDomain(nil)).To(BeNil())
		})

		It("returns empty slice for empty input", func() {
			Expect(display.CVBulletsFromDomain([]*career.CVBullet{})).To(Equal([]display.CVBullet{}))
		})
	})

	Describe("CVConfigFromDomain", func() {
		It("converts a fully populated CV config", func() {
			createdAt := time.Date(2024, time.June, 1, 9, 0, 0, 0, time.UTC)
			updatedAt := time.Date(2024, time.June, 2, 9, 0, 0, 0, time.UTC)

			config := fixtures.CVConfigWith("Staff CV", "staff", "technical-peer")
			config.EventFilters = map[string]interface{}{"tag": "technical"}
			config.TechnologyFocus = "specialist"
			config.SelectedTechnologies = []string{"go", "kubernetes"}
			config.FocusArea = "backend"
			config.LengthFormat = "standard"
			config.SkillsFormat = "grouped"
			config.SkillsLimit = 8
			config.SummaryHeading = "Profile"
			config.ProfileTitle = "Staff Engineer"
			config.WhatIBring = []string{"Systems thinking"}
			config.CoreStrengths = []string{"Architecture"}
			config.CreatedAt = createdAt
			config.UpdatedAt = updatedAt

			Expect(display.CVConfigFromDomain(config)).To(Equal(display.CVConfig{
				Name:                 "Staff CV",
				TargetRole:           "staff",
				TargetAudience:       "technical-peer",
				EventFilters:         map[string]interface{}{"tag": "technical"},
				TechnologyFocus:      "specialist",
				SelectedTechnologies: []string{"go", "kubernetes"},
				FocusArea:            "backend",
				LengthFormat:         "standard",
				SkillsFormat:         "grouped",
				SkillsLimit:          8,
				SummaryHeading:       "Profile",
				ProfileTitle:         "Staff Engineer",
				WhatIBring:           []string{"Systems thinking"},
				CoreStrengths:        []string{"Architecture"},
				CreatedAt:            createdAt,
				UpdatedAt:            updatedAt,
			}))
		})

		It("handles nil input gracefully", func() {
			Expect(display.CVConfigFromDomain(nil)).To(Equal(display.CVConfig{}))
		})
	})

	Describe("CVConfigsFromDomain", func() {
		It("converts a slice of CV configs", func() {
			c1 := fixtures.CVConfig("One")
			c2 := fixtures.CVConfig("Two")
			configs := []*career.CVConfig{c1, c2}

			result := display.CVConfigsFromDomain(configs)

			Expect(result).To(HaveLen(2))
			Expect(result[0].Name).To(Equal("One"))
			Expect(result[1].Name).To(Equal("Two"))
		})

		It("returns nil for nil input", func() {
			Expect(display.CVConfigsFromDomain(nil)).To(BeNil())
		})

		It("returns empty slice for empty input", func() {
			Expect(display.CVConfigsFromDomain([]*career.CVConfig{})).To(Equal([]display.CVConfig{}))
		})
	})
})
