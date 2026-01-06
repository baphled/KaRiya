package career

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCompetencyInference(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Competency Inference Suite")
}

var _ = Describe("InferCompetencyFromCategories", func() {
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
