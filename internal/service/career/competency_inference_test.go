package career

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	career "github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("InferCompetencyFromFacts", func() {
	Describe("Weighted competency inference", func() {
		It("should return competency with highest weighted score", func() {
			facts := []*career.Fact{
				{
					Text:                 "Fact 1",
					CompetencyCategories: []string{"technical"},
					StrengthSignal:       "strong", // weight: 3.0
				},
				{
					Text:                 "Fact 2",
					CompetencyCategories: []string{"technical"},
					StrengthSignal:       "moderate", // weight: 2.0
				},
				{
					Text:                 "Fact 3",
					CompetencyCategories: []string{"leadership"},
					StrengthSignal:       "strong", // weight: 3.0
				},
			}
			// technical: 3.0 + 2.0 = 5.0, leadership: 3.0
			result := InferCompetencyFromFacts(facts)
			Expect(result).To(Equal("technical"))
		})

		It("should handle tie-breaking alphabetically", func() {
			facts := []*career.Fact{
				{
					Text:                 "Fact 1",
					CompetencyCategories: []string{"technical"},
					StrengthSignal:       "strong", // weight: 3.0
				},
				{
					Text:                 "Fact 2",
					CompetencyCategories: []string{"leadership"},
					StrengthSignal:       "strong", // weight: 3.0
				},
			}
			// technical: 3.0, leadership: 3.0 (tie, alphabetically: leadership wins)
			result := InferCompetencyFromFacts(facts)
			Expect(result).To(Equal("leadership"))
		})

		It("should handle multiple competencies per fact", func() {
			facts := []*career.Fact{
				{
					Text:                 "Fact 1",
					CompetencyCategories: []string{"technical", "leadership"},
					StrengthSignal:       "strong", // weight: 3.0 each
				},
				{
					Text:                 "Fact 2",
					CompetencyCategories: []string{"technical"},
					StrengthSignal:       "weak", // weight: 1.0
				},
			}
			// technical: 3.0 + 1.0 = 4.0, leadership: 3.0
			result := InferCompetencyFromFacts(facts)
			Expect(result).To(Equal("technical"))
		})

		It("should use default weight for facts without strength signal", func() {
			facts := []*career.Fact{
				{
					Text:                 "Fact 1",
					CompetencyCategories: []string{"technical"},
					StrengthSignal:       "", // weight: 1.0 (default)
				},
				{
					Text:                 "Fact 2",
					CompetencyCategories: []string{"technical"},
					StrengthSignal:       "", // weight: 1.0 (default)
				},
			}
			// technical: 1.0 + 1.0 = 2.0
			result := InferCompetencyFromFacts(facts)
			Expect(result).To(Equal("technical"))
		})

		It("should handle all strength levels correctly", func() {
			facts := []*career.Fact{
				{
					Text:                 "Strong fact",
					CompetencyCategories: []string{"technical"},
					StrengthSignal:       "strong", // weight: 3.0
				},
				{
					Text:                 "Moderate fact",
					CompetencyCategories: []string{"leadership"},
					StrengthSignal:       "moderate", // weight: 2.0
				},
				{
					Text:                 "Weak fact",
					CompetencyCategories: []string{"product"},
					StrengthSignal:       "weak", // weight: 1.0
				},
			}
			// technical: 3.0, leadership: 2.0, product: 1.0
			result := InferCompetencyFromFacts(facts)
			Expect(result).To(Equal("technical"))
		})

		It("should return empty string when no facts", func() {
			facts := []*career.Fact{}
			result := InferCompetencyFromFacts(facts)
			Expect(result).To(Equal(""))
		})

		It("should return empty string when nil facts", func() {
			result := InferCompetencyFromFacts(nil)
			Expect(result).To(Equal(""))
		})

		It("should return empty string when facts have no competencies", func() {
			facts := []*career.Fact{
				{
					Text:                 "Fact 1",
					CompetencyCategories: []string{},
					StrengthSignal:       "strong",
				},
			}
			result := InferCompetencyFromFacts(facts)
			Expect(result).To(Equal(""))
		})

		It("should ignore invalid competencies", func() {
			facts := []*career.Fact{
				{
					Text:                 "Fact 1",
					CompetencyCategories: []string{"technical", "invalid-category"},
					StrengthSignal:       "strong",
				},
				{
					Text:                 "Fact 2",
					CompetencyCategories: []string{"leadership"},
					StrengthSignal:       "moderate",
				},
			}
			// Only valid: technical: 3.0, leadership: 2.0
			result := InferCompetencyFromFacts(facts)
			Expect(result).To(Equal("technical"))
		})

		It("should handle case sensitivity correctly", func() {
			facts := []*career.Fact{
				{
					Text:                 "Fact 1",
					CompetencyCategories: []string{"Technical"},
					StrengthSignal:       "strong",
				},
				{
					Text:                 "Fact 2",
					CompetencyCategories: []string{"TECHNICAL"},
					StrengthSignal:       "moderate",
				},
			}
			// All normalized to "technical": 3.0 + 2.0 = 5.0
			result := InferCompetencyFromFacts(facts)
			Expect(result).To(Equal("technical"))
		})

		It("should handle whitespace in competencies", func() {
			facts := []*career.Fact{
				{
					Text:                 "Fact 1",
					CompetencyCategories: []string{" technical "},
					StrengthSignal:       "strong",
				},
				{
					Text:                 "Fact 2",
					CompetencyCategories: []string{"technical"},
					StrengthSignal:       "moderate",
				},
			}
			// After trimming: technical: 3.0 + 2.0 = 5.0
			result := InferCompetencyFromFacts(facts)
			Expect(result).To(Equal("technical"))
		})
	})
})

var _ = Describe("InferCompetencyFromCategories (deprecated)", func() {
	Describe("Single category scenarios", func() {
		It("should return the only category when all events have the same single category", func() {
			eventCategories := [][]string{
				{"technical"},
				{"technical"},
				{"technical"},
			}
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal("technical"))
		})

		It("should return empty string when all events have empty categories", func() {
			eventCategories := [][]string{
				{},
				{},
				{},
			}
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal(""))
		})

		It("should return empty string when given empty input", func() {
			eventCategories := [][]string{}
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal(""))
		})

		It("should return empty string when given nil input", func() {
			result := InferCompetencyFromCategories(nil)
			Expect(result).To(Equal(""))
		})
	})

	Describe("Multiple category scenarios", func() {
		It("should return most common category when there's a clear winner", func() {
			eventCategories := [][]string{
				{"technical"},
				{"technical", "leadership"},
				{"technical"},
				{"leadership"},
				{"product"},
			}
			// technical: 3, leadership: 2, product: 1
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal("technical"))
		})

		It("should handle multiple categories per event", func() {
			eventCategories := [][]string{
				{"technical", "leadership", "product"},
				{"leadership", "mentoring"},
				{"leadership"},
			}
			// leadership: 3, technical: 1, product: 1, mentoring: 1
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal("leadership"))
		})

		It("should break ties alphabetically", func() {
			eventCategories := [][]string{
				{"technical", "leadership"},
				{"technical", "leadership"},
			}
			// technical: 2, leadership: 2 (tie)
			// Alphabetically: leadership < technical
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal("leadership"))
		})

		It("should break three-way tie alphabetically", func() {
			eventCategories := [][]string{
				{"technical"},
				{"leadership"},
				{"product"},
			}
			// All tied at 1
			// Alphabetically: leadership < product < technical
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal("leadership"))
		})
	})

	Describe("Mixed scenarios", func() {
		It("should handle mix of events with and without categories", func() {
			eventCategories := [][]string{
				{"technical"},
				{},
				{"technical"},
				{},
				{"leadership"},
			}
			// technical: 2, leadership: 1
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal("technical"))
		})

		It("should ignore invalid categories and count valid ones", func() {
			eventCategories := [][]string{
				{"technical", "invalid-category"},
				{"technical"},
				{"leadership", "another-invalid"},
			}
			// Only valid: technical: 2, leadership: 1
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal("technical"))
		})

		It("should return empty when all categories are invalid", func() {
			eventCategories := [][]string{
				{"invalid-1"},
				{"invalid-2"},
				{"not-a-category"},
			}
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal(""))
		})

		It("should handle case sensitivity correctly", func() {
			eventCategories := [][]string{
				{"Technical"},
				{"TECHNICAL"},
				{"technical"},
			}
			// All should be normalized to "technical"
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal("technical"))
		})
	})

	Describe("All allowed categories", func() {
		It("should handle all six allowed categories", func() {
			eventCategories := [][]string{
				{"technical"},
				{"leadership"},
				{"product"},
				{"consulting"},
				{"research"},
				{"mentoring"},
			}
			// All tied at 1, alphabetically: consulting
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal("consulting"))
		})

		It("should correctly count when consulting is most common", func() {
			eventCategories := [][]string{
				{"consulting"},
				{"consulting"},
				{"technical"},
			}
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal("consulting"))
		})

		It("should correctly count when research is most common", func() {
			eventCategories := [][]string{
				{"research"},
				{"research"},
				{"research"},
				{"technical"},
			}
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal("research"))
		})

		It("should correctly count when mentoring is most common", func() {
			eventCategories := [][]string{
				{"mentoring"},
				{"mentoring"},
				{"leadership"},
			}
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal("mentoring"))
		})
	})

	Describe("Edge cases", func() {
		It("should handle single event with single category", func() {
			eventCategories := [][]string{
				{"leadership"},
			}
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal("leadership"))
		})

		It("should handle single event with multiple categories", func() {
			eventCategories := [][]string{
				{"technical", "leadership", "product"},
			}
			// All tied at 1, alphabetically: leadership
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal("leadership"))
		})

		It("should handle very large number of events", func() {
			eventCategories := make([][]string, 1000)
			for i := 0; i < 1000; i++ {
				if i%2 == 0 {
					eventCategories[i] = []string{"technical"}
				} else {
					eventCategories[i] = []string{"leadership"}
				}
			}
			// technical: 500, leadership: 500 (tie)
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal("leadership")) // alphabetical
		})

		It("should handle whitespace in categories", func() {
			eventCategories := [][]string{
				{" technical "},
				{"technical"},
				{" leadership "},
			}
			// After trimming: technical: 2, leadership: 1
			result := InferCompetencyFromCategories(eventCategories)
			Expect(result).To(Equal("technical"))
		})
	})
})
