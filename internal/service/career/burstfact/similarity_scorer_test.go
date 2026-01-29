package burstfact

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SimilarityScorer", func() {
	var scorer *SimilarityScorer

	BeforeEach(func() {
		scorer = NewSimilarityScorer()
	})

	Describe("TextSimilarity", func() {
		Context("when comparing identical text", func() {
			It("should return 1.0 (perfect match)", func() {
				text1 := "Led cross-functional team to deliver critical project"
				text2 := "Led cross-functional team to deliver critical project"

				score := scorer.TextSimilarity(text1, text2)

				Expect(score).To(BeNumerically("==", 1.0))
			})
		})

		Context("when comparing very similar text", func() {
			It("should return high score (>0.7)", func() {
				text1 := "Led cross-functional team to deliver critical project"
				text2 := "Led cross-functional team to deliver critical feature"

				score := scorer.TextSimilarity(text1, text2)

				// Jaccard: 6 matching tokens out of 8 union = 0.75
				Expect(score).To(BeNumerically(">", 0.7))
			})
		})

		Context("when comparing partially similar text", func() {
			It("should return moderate score (0-1)", func() {
				text1 := "Led cross-functional team to deliver critical project"
				text2 := "Worked on backend optimization for performance"

				score := scorer.TextSimilarity(text1, text2)

				Expect(score).To(BeNumerically(">=", 0.0))
				Expect(score).To(BeNumerically("<=", 1.0))
			})
		})

		Context("when comparing completely different text", func() {
			It("should return low score (close to 0)", func() {
				text1 := "Led cross-functional team to deliver critical project"
				text2 := "Organized office party for team building"

				score := scorer.TextSimilarity(text1, text2)

				Expect(score).To(BeNumerically("<=", 0.5))
			})
		})

		Context("when comparing empty strings", func() {
			It("should return 1.0 (both empty)", func() {
				score := scorer.TextSimilarity("", "")

				Expect(score).To(BeNumerically("==", 1.0))
			})
		})

		Context("when one string is empty", func() {
			It("should return 0", func() {
				score := scorer.TextSimilarity("some text", "")

				Expect(score).To(Equal(0.0))
			})
		})
	})

	Describe("KeywordOverlapScore", func() {
		Context("when both events have identical keywords", func() {
			It("should return 1.0 (perfect match)", func() {
				keywords1 := []string{"leadership", "technical", "backend"}
				keywords2 := []string{"leadership", "technical", "backend"}

				score := scorer.KeywordOverlapScore(keywords1, keywords2)

				Expect(score).To(BeNumerically("==", 1.0))
			})
		})

		Context("when events have partial keyword overlap", func() {
			It("should return score based on overlap ratio", func() {
				keywords1 := []string{"leadership", "technical", "backend"}
				keywords2 := []string{"leadership", "technical", "frontend"}

				score := scorer.KeywordOverlapScore(keywords1, keywords2)

				// 2 matching, union = 3+3-2 = 4, Jaccard = 2/4 = 0.5
				Expect(score).To(BeNumerically(">", 0.4))
				Expect(score).To(BeNumerically("<", 0.6))
			})
		})

		Context("when events have no keyword overlap", func() {
			It("should return 0", func() {
				keywords1 := []string{"leadership", "technical"}
				keywords2 := []string{"frontend", "design"}

				score := scorer.KeywordOverlapScore(keywords1, keywords2)

				Expect(score).To(Equal(0.0))
			})
		})

		Context("when one keyword list is empty", func() {
			It("should return 0", func() {
				keywords1 := []string{"leadership", "technical"}
				keywords2 := []string{}

				score := scorer.KeywordOverlapScore(keywords1, keywords2)

				Expect(score).To(Equal(0.0))
			})
		})

		Context("when both keyword lists are empty", func() {
			It("should return 1.0 (both empty)", func() {
				keywords1 := []string{}
				keywords2 := []string{}

				score := scorer.KeywordOverlapScore(keywords1, keywords2)

				Expect(score).To(BeNumerically("==", 1.0))
			})
		})
	})

	Describe("CompanyMatchScore", func() {
		Context("when companies are identical", func() {
			It("should return 1.0", func() {
				score := scorer.CompanyMatchScore("TechCorp Inc.", "TechCorp Inc.")

				Expect(score).To(Equal(1.0))
			})
		})

		Context("when companies are different", func() {
			It("should return 0", func() {
				score := scorer.CompanyMatchScore("TechCorp Inc.", "OtherCorp Ltd.")

				Expect(score).To(Equal(0.0))
			})
		})

		Context("when one company is empty", func() {
			It("should return 0", func() {
				score := scorer.CompanyMatchScore("TechCorp Inc.", "")

				Expect(score).To(Equal(0.0))
			})
		})

		Context("when both companies are empty", func() {
			It("should return 1.0 (both empty)", func() {
				score := scorer.CompanyMatchScore("", "")

				Expect(score).To(Equal(1.0))
			})
		})

		Context("when comparing case-insensitively", func() {
			It("should match regardless of case", func() {
				score := scorer.CompanyMatchScore("techcorp inc.", "TECHCORP INC.")

				Expect(score).To(Equal(1.0))
			})
		})
	})

	Describe("ProjectMatchScore", func() {
		Context("when projects are identical", func() {
			It("should return 1.0", func() {
				score := scorer.ProjectMatchScore("Platform Migration", "Platform Migration")

				Expect(score).To(Equal(1.0))
			})
		})

		Context("when projects are different", func() {
			It("should return 0", func() {
				score := scorer.ProjectMatchScore("Platform Migration", "API Redesign")

				Expect(score).To(Equal(0.0))
			})
		})

		Context("when one project is empty", func() {
			It("should return 0", func() {
				score := scorer.ProjectMatchScore("Platform Migration", "")

				Expect(score).To(Equal(0.0))
			})
		})

		Context("when both projects are empty", func() {
			It("should return 1.0 (both empty)", func() {
				score := scorer.ProjectMatchScore("", "")

				Expect(score).To(Equal(1.0))
			})
		})

		Context("when comparing case-insensitively", func() {
			It("should match regardless of case", func() {
				score := scorer.ProjectMatchScore("platform migration", "PLATFORM MIGRATION")

				Expect(score).To(Equal(1.0))
			})
		})
	})

	Describe("CombinedSimilarityScore", func() {
		Context("when all components match perfectly", func() {
			It("should return 1.0", func() {
				event1 := EventSimilarityInput{
					Text:     "Led team to deliver project",
					Keywords: []string{"leadership", "technical"},
					Company:  "TechCorp Inc.",
					Project:  "Platform Migration",
				}
				event2 := EventSimilarityInput{
					Text:     "Led team to deliver project",
					Keywords: []string{"leadership", "technical"},
					Company:  "TechCorp Inc.",
					Project:  "Platform Migration",
				}

				score := scorer.CombinedSimilarityScore(event1, event2)

				Expect(score).To(BeNumerically("==", 1.0))
			})
		})

		Context("when text matches but other components differ", func() {
			It("should return score weighted by text similarity", func() {
				event1 := EventSimilarityInput{
					Text:     "Led team to deliver project",
					Keywords: []string{"leadership", "technical"},
					Company:  "TechCorp Inc.",
					Project:  "Platform Migration",
				}
				event2 := EventSimilarityInput{
					Text:     "Led team to deliver project",
					Keywords: []string{"frontend", "design"},
					Company:  "OtherCorp Ltd.",
					Project:  "API Redesign",
				}

				score := scorer.CombinedSimilarityScore(event1, event2)

				// Text matches (1.0), others don't (0.0)
				// Weighted: 0.5*1.0 + 0.2*0.0 + 0.15*0.0 + 0.15*0.0 = 0.5
				Expect(score).To(BeNumerically(">=", 0.4))
				Expect(score).To(BeNumerically("<=", 0.6))
			})
		})

		Context("when components are completely different", func() {
			It("should return low score", func() {
				event1 := EventSimilarityInput{
					Text:     "Led team to deliver project",
					Keywords: []string{"leadership"},
					Company:  "TechCorp Inc.",
					Project:  "Platform Migration",
				}
				event2 := EventSimilarityInput{
					Text:     "Organized office party",
					Keywords: []string{"social"},
					Company:  "OtherCorp Ltd.",
					Project:  "Team Building",
				}

				score := scorer.CombinedSimilarityScore(event1, event2)

				Expect(score).To(BeNumerically("<", 0.5))
			})
		})

		Context("when one event has empty optional fields", func() {
			It("should still compute valid similarity", func() {
				event1 := EventSimilarityInput{
					Text:     "Led team to deliver",
					Keywords: []string{"leadership"},
					Company:  "TechCorp Inc.",
					Project:  "",
				}
				event2 := EventSimilarityInput{
					Text:     "Led team to deliver",
					Keywords: []string{"leadership"},
					Company:  "",
					Project:  "",
				}

				score := scorer.CombinedSimilarityScore(event1, event2)

				Expect(score).To(BeNumerically(">", 0.0))
				Expect(score).To(BeNumerically("<=", 1.0))
			})
		})
	})

	Describe("Edge cases", func() {
		Context("when comparing very long text", func() {
			It("should handle efficiently", func() {
				longText := "This is a very long text that describes " +
					"a complex project involving multiple teams " +
					"and spanning several months with many deliverables " +
					"and achievements that resulted in significant impact"

				score := scorer.TextSimilarity(longText, longText)

				Expect(score).To(Equal(1.0))
			})
		})

		Context("when comparing text with special characters", func() {
			It("should handle correctly", func() {
				text1 := "Led team on C++ optimization (performance: +50%)"
				text2 := "Led team on C++ optimization (performance: +50%)"

				score := scorer.TextSimilarity(text1, text2)

				Expect(score).To(Equal(1.0))
			})
		})

		Context("when comparing text with unicode characters", func() {
			It("should handle correctly", func() {
				text1 := "Led team on café project with résumé review"
				text2 := "Led team on café project with résumé review"

				score := scorer.TextSimilarity(text1, text2)

				Expect(score).To(Equal(1.0))
			})
		})
	})
})
