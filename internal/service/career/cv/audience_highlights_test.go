package cv_test

import (
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Audience-Aware Highlights", func() {
	var (
		view    career.CVView
		bullets []*career.CVBullet
	)

	BeforeEach(func() {
		bulletA := fixtures.CVBulletWith("b-arch", "exp-1", "Architected microservices platform handling 10K rps")
		bulletA.Confidence = 0.90
		bulletA.AudienceRelevance = map[string]float64{
			"technical-peer":      0.95,
			"hiring-manager":      0.30,
			"recruiter":           0.40,
			"engineering-manager": 0.35,
		}

		bulletB := fixtures.CVBulletWith("b-cost", "exp-1", "Reduced infrastructure costs by 40% through capacity planning")
		bulletB.Confidence = 0.85
		bulletB.AudienceRelevance = map[string]float64{
			"technical-peer":      0.30,
			"hiring-manager":      0.95,
			"recruiter":           0.50,
			"engineering-manager": 0.60,
		}

		bulletC := fixtures.CVBulletWith("b-lead", "exp-1", "Led team of 8 engineers delivering platform migration")
		bulletC.Confidence = 0.88
		bulletC.AudienceRelevance = map[string]float64{
			"technical-peer":      0.20,
			"hiring-manager":      0.50,
			"recruiter":           0.85,
			"engineering-manager": 0.95,
		}

		bulletD := fixtures.CVBulletWith("b-ci", "exp-2", "Implemented CI/CD pipeline reducing deploy time by 60%")
		bulletD.Confidence = 0.80
		bulletD.AudienceRelevance = map[string]float64{
			"technical-peer":      0.80,
			"hiring-manager":      0.40,
			"recruiter":           0.60,
			"engineering-manager": 0.30,
		}

		bulletE := fixtures.CVBulletWith("b-mentor", "exp-2", "Mentored 5 junior engineers to mid-level promotions")
		bulletE.Confidence = 0.82
		bulletE.AudienceRelevance = map[string]float64{
			"technical-peer":      0.15,
			"hiring-manager":      0.60,
			"recruiter":           0.90,
			"engineering-manager": 0.90,
		}

		bulletF := fixtures.CVBulletWith("b-api", "exp-2", "Designed RESTful API serving 50M requests per day")
		bulletF.Confidence = 0.87
		bulletF.AudienceRelevance = map[string]float64{
			"technical-peer":      0.90,
			"hiring-manager":      0.35,
			"recruiter":           0.45,
			"engineering-manager": 0.25,
		}

		bullets = []*career.CVBullet{bulletA, bulletB, bulletC, bulletD, bulletE, bulletF}

		expSection1 := fixtures.CVSectionWith("exp-1", "", "experience", "Experience", 1)
		expSection1.Content = []*career.SectionContentGroup{
			fixtures.ContentGroupWithBullets("Acme Corp", []*career.CVBullet{bulletA, bulletB, bulletC}),
		}

		expSection2 := fixtures.CVSectionWith("exp-2", "", "experience", "Experience", 2)
		expSection2.Content = []*career.SectionContentGroup{
			fixtures.ContentGroupWithBullets("TechStart Inc", []*career.CVBullet{bulletD, bulletE, bulletF}),
		}

		summarySection := fixtures.CVSectionWithSummary("summary", "", "Experienced engineer specialising in distributed systems.")

		viewPtr := fixtures.CVViewWithSections("cv-audience-test", []*career.CVSection{summarySection, expSection1, expSection2})
		viewPtr.Name = "Test User"
		viewPtr.TargetAudience = "hiring_manager"
		viewPtr.SourceEventCount = 20
		viewPtr.SourceFactCount = 15
		view = *viewPtr
	})

	Describe("GetTopBulletsForAudience", func() {
		Context("when audience is technical-peer", func() {
			It("ranks technically relevant bullets first", func() {
				result := cv.GetTopBulletsForAudience(view, "technical-peer", 3)

				Expect(result).To(HaveLen(3))
				Expect(result[0].Text).To(ContainSubstring("Architected microservices"))
				Expect(result[1].Text).To(ContainSubstring("RESTful API"))
			})
		})

		Context("when audience is hiring-manager", func() {
			It("ranks business-impact bullets first", func() {
				result := cv.GetTopBulletsForAudience(view, "hiring-manager", 3)

				Expect(result).To(HaveLen(3))
				Expect(result[0].Text).To(ContainSubstring("Reduced infrastructure costs"))
			})
		})

		Context("when audience is recruiter", func() {
			It("ranks broadly appealing bullets first", func() {
				result := cv.GetTopBulletsForAudience(view, "recruiter", 3)

				Expect(result).To(HaveLen(3))
				Expect(result[0].Text).To(ContainSubstring("Led team of 8"))
			})
		})

		Context("when audience is engineering-manager", func() {
			It("ranks leadership and collaboration bullets first", func() {
				result := cv.GetTopBulletsForAudience(view, "engineering-manager", 3)

				Expect(result).To(HaveLen(3))
				Expect(result[0].Text).To(ContainSubstring("Led team of 8"))
			})
		})

		Context("when different audiences produce different orderings", func() {
			It("returns different top bullets for technical-peer vs hiring-manager", func() {
				techResult := cv.GetTopBulletsForAudience(view, "technical-peer", 1)
				hiringResult := cv.GetTopBulletsForAudience(view, "hiring-manager", 1)

				Expect(techResult[0].ID).NotTo(Equal(hiringResult[0].ID))
			})

			It("returns different top bullets for recruiter vs technical-peer", func() {
				recruiterResult := cv.GetTopBulletsForAudience(view, "recruiter", 1)
				techResult := cv.GetTopBulletsForAudience(view, "technical-peer", 1)

				Expect(recruiterResult[0].ID).NotTo(Equal(techResult[0].ID))
			})
		})

		Context("when bullets have no audience relevance data", func() {
			BeforeEach(func() {
				for _, b := range bullets {
					b.AudienceRelevance = nil
				}
			})

			It("falls back to confidence-only ordering", func() {
				result := cv.GetTopBulletsForAudience(view, "technical-peer", 3)

				Expect(result).To(HaveLen(3))
				Expect(result[0].Confidence).To(BeNumerically(">=", result[1].Confidence))
				Expect(result[1].Confidence).To(BeNumerically(">=", result[2].Confidence))
			})
		})

		Context("when audience relevance map is empty", func() {
			BeforeEach(func() {
				for _, b := range bullets {
					b.AudienceRelevance = map[string]float64{}
				}
			})

			It("falls back to confidence-only ordering", func() {
				result := cv.GetTopBulletsForAudience(view, "hiring-manager", 3)

				Expect(result).To(HaveLen(3))
				Expect(result[0].Confidence).To(BeNumerically(">=", result[1].Confidence))
				Expect(result[1].Confidence).To(BeNumerically(">=", result[2].Confidence))
			})
		})

		Context("when audience key is unknown", func() {
			It("falls back to confidence-only ordering", func() {
				result := cv.GetTopBulletsForAudience(view, "unknown-audience", 3)

				Expect(result).To(HaveLen(3))
				Expect(result[0].Confidence).To(BeNumerically(">=", result[1].Confidence))
				Expect(result[1].Confidence).To(BeNumerically(">=", result[2].Confidence))
			})
		})

		Context("when limit is respected", func() {
			It("returns exactly the requested number of bullets", func() {
				result := cv.GetTopBulletsForAudience(view, "technical-peer", 3)
				Expect(result).To(HaveLen(3))
			})

			It("returns fewer bullets when fewer are available", func() {
				result := cv.GetTopBulletsForAudience(view, "technical-peer", 100)
				Expect(result).To(HaveLen(len(bullets)))
			})

			It("returns one bullet when limit is 1", func() {
				result := cv.GetTopBulletsForAudience(view, "recruiter", 1)
				Expect(result).To(HaveLen(1))
			})
		})

		Context("when view has no sections", func() {
			It("returns an empty slice", func() {
				emptyViewPtr := fixtures.CVView("empty-cv")
				emptyViewPtr.Name = "Empty"
				emptyView := *emptyViewPtr
				result := cv.GetTopBulletsForAudience(emptyView, "technical-peer", 5)
				Expect(result).To(BeEmpty())
			})
		})
	})

	Describe("ExportHighlightsForAudience", func() {
		var profile config.ProfileConfig

		BeforeEach(func() {
			profile = config.ProfileConfig{
				Name:     "Test User",
				Email:    "test@example.com",
				Title:    "Staff Engineer",
				Location: "London, UK",
				GitHub:   "testuser",
				CoreStrengths: []string{
					"Distributed systems design",
					"Technical leadership",
					"Performance optimisation",
				},
				Languages: []string{"Go", "Ruby", "Python"},
				Systems:   []string{"Linux", "Kubernetes", "PostgreSQL"},
			}
		})

		Context("when audience has a known summary prefix", func() {
			It("includes the technical-peer prefix", func() {
				output := cv.ExportHighlightsForAudience(view, profile, "technical-peer")
				Expect(output).To(ContainSubstring("Technically deep engineer with"))
			})

			It("includes the hiring-manager prefix", func() {
				output := cv.ExportHighlightsForAudience(view, profile, "hiring-manager")
				Expect(output).To(ContainSubstring("Results-driven engineer delivering"))
			})

			It("includes the recruiter prefix", func() {
				output := cv.ExportHighlightsForAudience(view, profile, "recruiter")
				Expect(output).To(ContainSubstring("Versatile software engineer with"))
			})

			It("includes the engineering-manager prefix", func() {
				output := cv.ExportHighlightsForAudience(view, profile, "engineering-manager")
				Expect(output).To(ContainSubstring("Collaborative engineer who"))
			})
		})

		Context("when audience is master or unknown", func() {
			It("falls back to section summary for unknown audience", func() {
				output := cv.ExportHighlightsForAudience(view, profile, "master")
				Expect(output).NotTo(ContainSubstring("Technically deep"))
				Expect(output).NotTo(ContainSubstring("Results-driven"))
				Expect(output).To(ContainSubstring("Experienced engineer specialising"))
			})

			It("falls back to section summary for empty audience", func() {
				output := cv.ExportHighlightsForAudience(view, profile, "")
				Expect(output).To(ContainSubstring("Experienced engineer specialising"))
			})
		})

		Context("when output structure is correct", func() {
			It("includes the profile name in header", func() {
				output := cv.ExportHighlightsForAudience(view, profile, "technical-peer")
				Expect(output).To(ContainSubstring("TEST USER"))
			})

			It("includes KEY CAPABILITIES section", func() {
				output := cv.ExportHighlightsForAudience(view, profile, "hiring-manager")
				Expect(output).To(ContainSubstring("KEY CAPABILITIES"))
			})

			It("includes SELECTED HIGHLIGHTS section", func() {
				output := cv.ExportHighlightsForAudience(view, profile, "recruiter")
				Expect(output).To(ContainSubstring("SELECTED HIGHLIGHTS"))
			})

			It("includes TECHNOLOGIES section", func() {
				output := cv.ExportHighlightsForAudience(view, profile, "technical-peer")
				Expect(output).To(ContainSubstring("TECHNOLOGIES"))
				Expect(output).To(ContainSubstring("Go, Ruby, Python"))
			})

			It("includes core strengths", func() {
				output := cv.ExportHighlightsForAudience(view, profile, "hiring-manager")
				Expect(output).To(ContainSubstring("Distributed systems design"))
				Expect(output).To(ContainSubstring("Technical leadership"))
			})
		})

		Context("when audience-specific bullets are selected", func() {
			It("includes audience-relevant highlights for technical-peer", func() {
				output := cv.ExportHighlightsForAudience(view, profile, "technical-peer")
				Expect(output).To(ContainSubstring("Architected microservices"))
			})

			It("includes audience-relevant highlights for hiring-manager", func() {
				output := cv.ExportHighlightsForAudience(view, profile, "hiring-manager")
				Expect(output).To(ContainSubstring("Reduced infrastructure costs"))
			})
		})
	})
})
