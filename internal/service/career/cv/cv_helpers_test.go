package cv_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/service/career/cv"
)

var _ = Describe("CVHelpers", func() {
	Describe("ApplyProfileOverride", func() {
		var baseCfg *config.ProfileConfig

		BeforeEach(func() {
			baseCfg = &config.ProfileConfig{
				Name:          "Test User",
				Email:         "test@example.com",
				Title:         "Senior Engineer",
				CoreStrengths: []string{"Go", "Python"},
				WhatIBring:    []string{"Leadership", "Mentoring"},
			}
		})

		Context("when override is nil", func() {
			It("should return the original config unchanged", func() {
				result := cv.ApplyProfileOverride(baseCfg, nil)
				Expect(result).To(Equal(baseCfg))
			})
		})

		Context("when config is nil", func() {
			It("should create a config with override values", func() {
				title := "Consulting Engineer"
				override := &cv.ProfileOverride{
					ProfessionalTitle: &title,
					CoreStrengths:     []string{"Consulting", "Architecture"},
				}

				result := cv.ApplyProfileOverride(nil, override)

				Expect(result).NotTo(BeNil())
				Expect(result.Title).To(Equal("Consulting Engineer"))
				Expect(result.CoreStrengths).To(Equal([]string{"Consulting", "Architecture"}))
			})
		})

		Context("when override has ProfessionalTitle", func() {
			It("should override the title", func() {
				title := "Staff Engineer"
				override := &cv.ProfileOverride{
					ProfessionalTitle: &title,
				}

				result := cv.ApplyProfileOverride(baseCfg, override)

				Expect(result.Title).To(Equal("Staff Engineer"))
				// Other fields unchanged
				Expect(result.Name).To(Equal("Test User"))
				Expect(result.CoreStrengths).To(Equal([]string{"Go", "Python"}))
			})
		})

		Context("when override has CoreStrengths", func() {
			It("should override core strengths", func() {
				override := &cv.ProfileOverride{
					CoreStrengths: []string{"Architecture", "Team Leadership", "Mentoring"},
				}

				result := cv.ApplyProfileOverride(baseCfg, override)

				Expect(result.CoreStrengths).To(Equal([]string{"Architecture", "Team Leadership", "Mentoring"}))
				// Other fields unchanged
				Expect(result.Title).To(Equal("Senior Engineer"))
			})
		})

		Context("when override has CareerDifferentiators", func() {
			It("should override WhatIBring", func() {
				override := &cv.ProfileOverride{
					CareerDifferentiators: []string{"Cross-functional leadership", "Technical vision"},
				}

				result := cv.ApplyProfileOverride(baseCfg, override)

				Expect(result.WhatIBring).To(Equal([]string{"Cross-functional leadership", "Technical vision"}))
			})
		})

		Context("when override has multiple fields", func() {
			It("should apply all overrides", func() {
				title := "Principal Engineer"
				override := &cv.ProfileOverride{
					ProfessionalTitle:     &title,
					CoreStrengths:         []string{"System Design", "Technical Strategy"},
					CareerDifferentiators: []string{"Org-wide impact", "Technical roadmaps"},
				}

				result := cv.ApplyProfileOverride(baseCfg, override)

				Expect(result.Title).To(Equal("Principal Engineer"))
				Expect(result.CoreStrengths).To(Equal([]string{"System Design", "Technical Strategy"}))
				Expect(result.WhatIBring).To(Equal([]string{"Org-wide impact", "Technical roadmaps"}))
				// Original fields unchanged
				Expect(result.Name).To(Equal("Test User"))
				Expect(result.Email).To(Equal("test@example.com"))
			})
		})

		It("should not modify the original config", func() {
			title := "Staff Engineer"
			override := &cv.ProfileOverride{
				ProfessionalTitle: &title,
			}

			_ = cv.ApplyProfileOverride(baseCfg, override)

			// Original should be unchanged
			Expect(baseCfg.Title).To(Equal("Senior Engineer"))
		})
	})

	Describe("ApplyProfileOverrideToNarrative", func() {
		var baseProfile *cv.NarrativeProfileData

		BeforeEach(func() {
			baseProfile = &cv.NarrativeProfileData{
				Name:              "Test User",
				Role:              "Senior Engineer",
				Location:          "Remote",
				Email:             "test@example.com",
				GitHub:            "github.com/test",
				Portfolio:         "test.com",
				CoreStrengths:     []string{"Go", "Python"},
				Languages:         "Go, Python",
				Frontend:          "React",
				Systems:           "Linux",
				ValuePropositions: []string{"Leadership", "Mentoring"},
			}
		})

		Context("when override is nil", func() {
			It("should return the original profile unchanged", func() {
				result := cv.ApplyProfileOverrideToNarrative(baseProfile, nil)
				Expect(result).To(Equal(baseProfile))
			})
		})

		Context("when profile is nil", func() {
			It("should return nil", func() {
				override := &cv.ProfileOverride{}
				result := cv.ApplyProfileOverrideToNarrative(nil, override)
				Expect(result).To(BeNil())
			})
		})

		Context("when override has ProfessionalTitle", func() {
			It("should override the role", func() {
				title := "Staff Engineer"
				override := &cv.ProfileOverride{
					ProfessionalTitle: &title,
				}

				result := cv.ApplyProfileOverrideToNarrative(baseProfile, override)

				Expect(result.Role).To(Equal("Staff Engineer"))
				// Other fields unchanged
				Expect(result.Name).To(Equal("Test User"))
			})
		})

		Context("when override has CoreStrengths", func() {
			It("should override core strengths", func() {
				override := &cv.ProfileOverride{
					CoreStrengths: []string{"Architecture", "Team Leadership"},
				}

				result := cv.ApplyProfileOverrideToNarrative(baseProfile, override)

				Expect(result.CoreStrengths).To(Equal([]string{"Architecture", "Team Leadership"}))
			})
		})

		Context("when override has CareerDifferentiators", func() {
			It("should override value propositions", func() {
				override := &cv.ProfileOverride{
					CareerDifferentiators: []string{"Cross-functional leadership"},
				}

				result := cv.ApplyProfileOverrideToNarrative(baseProfile, override)

				Expect(result.ValuePropositions).To(Equal([]string{"Cross-functional leadership"}))
			})
		})

		It("should not modify the original profile", func() {
			title := "Staff Engineer"
			override := &cv.ProfileOverride{
				ProfessionalTitle: &title,
				CoreStrengths:     []string{"New Strength"},
			}

			_ = cv.ApplyProfileOverrideToNarrative(baseProfile, override)

			// Original should be unchanged
			Expect(baseProfile.Role).To(Equal("Senior Engineer"))
			Expect(baseProfile.CoreStrengths).To(Equal([]string{"Go", "Python"}))
		})
	})
})
