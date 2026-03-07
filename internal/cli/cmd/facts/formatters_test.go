package facts_test

import (
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/cmd/facts"
	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Facts Formatters", func() {
	Describe("FormatExtractionResults", func() {
		Context("when no facts extracted", func() {
			It("should return header with zero count", func() {
				output := facts.FormatExtractionResults(0, 10, map[string]int{})
				Expect(output).To(ContainSubstring("=== Fact Extraction Results ==="))
				Expect(output).To(ContainSubstring("Extracted 0 facts from 10 events"))
				Expect(output).NotTo(ContainSubstring("Competency breakdown"))
			})
		})

		Context("when facts extracted without competencies", func() {
			It("should show count only without breakdown", func() {
				output := facts.FormatExtractionResults(5, 10, map[string]int{})
				Expect(output).To(ContainSubstring("=== Fact Extraction Results ==="))
				Expect(output).To(ContainSubstring("Extracted 5 facts from 10 events"))
				Expect(output).NotTo(ContainSubstring("Competency breakdown"))
			})
		})

		Context("when facts have competencies", func() {
			It("should show sorted competency breakdown", func() {
				competencies := map[string]int{
					"Leadership":   5,
					"Technical":    10,
					"Architecture": 3,
				}
				output := facts.FormatExtractionResults(18, 20, competencies)

				Expect(output).To(ContainSubstring("=== Fact Extraction Results ==="))
				Expect(output).To(ContainSubstring("Extracted 18 facts from 20 events"))
				Expect(output).To(ContainSubstring("Competency breakdown"))
				Expect(output).To(ContainSubstring("Architecture: 3 facts"))
				Expect(output).To(ContainSubstring("Leadership: 5 facts"))
				Expect(output).To(ContainSubstring("Technical: 10 facts"))

				// Verify alphabetical order
				archIdx := indexOf(output, "Architecture")
				leadIdx := indexOf(output, "Leadership")
				techIdx := indexOf(output, "Technical")
				Expect(archIdx).To(BeNumerically("<", leadIdx))
				Expect(leadIdx).To(BeNumerically("<", techIdx))
			})
		})

		Context("with single competency", func() {
			It("should format single competency correctly", func() {
				competencies := map[string]int{"Backend": 7}
				output := facts.FormatExtractionResults(7, 5, competencies)

				Expect(output).To(ContainSubstring("Competency breakdown"))
				Expect(output).To(ContainSubstring("Backend: 7 facts"))
			})
		})

		Context("with empty competency map", func() {
			It("should not show breakdown section", func() {
				output := facts.FormatExtractionResults(10, 15, map[string]int{})
				Expect(output).NotTo(ContainSubstring("Competency breakdown"))
			})
		})
	})

	Describe("FormatFactList", func() {
		Context("with empty list", func() {
			It("should show zero count", func() {
				output := facts.FormatFactList([]*domain.Fact{})
				Expect(output).To(ContainSubstring("Existing Facts"))
				Expect(output).To(ContainSubstring("Total facts: 0"))
			})
		})

		Context("with facts having all fields populated", func() {
			It("should format complete details", func() {
				createdAt := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
				fact := fixtures.FactWithCategories(
					"fact-1",
					"Implemented authentication system",
					"event-1",
					[]string{"Security", "Backend"},
					[]string{"Technical", "Leadership"},
				)
				fact.RoleFit = "Senior Engineer"
				fact.CreatedAt = createdAt
				factList := []*domain.Fact{fact}

				output := facts.FormatFactList(factList)

				Expect(output).To(ContainSubstring("Existing Facts"))
				Expect(output).To(ContainSubstring("Total facts: 1"))
				Expect(output).To(ContainSubstring("1. Implemented authentication system"))
				Expect(output).To(ContainSubstring("ID: fact-1"))
				Expect(output).To(ContainSubstring("Source Event: event-1"))
				Expect(output).To(ContainSubstring("Competencies: Security, Backend"))
				Expect(output).To(ContainSubstring("Role Fit: Senior Engineer"))
				Expect(output).To(ContainSubstring("Target Audiences: Technical, Leadership"))
				Expect(output).To(ContainSubstring("Created: 2024-01-15 10:00:00"))
			})
		})

		Context("with facts missing optional fields", func() {
			It("should omit empty fields from output", func() {
				fact := fixtures.Fact("fact-1", "event-1")
				fact.Text = "Simple fact"
				fact.CompetencyCategories = []string{}
				fact.RoleFit = ""
				fact.AudienceRelevance = []string{}
				factList := []*domain.Fact{fact}

				output := facts.FormatFactList(factList)

				Expect(output).To(ContainSubstring("Simple fact"))
				Expect(output).NotTo(ContainSubstring("Competencies:"))
				Expect(output).NotTo(ContainSubstring("Role Fit:"))
				Expect(output).NotTo(ContainSubstring("Target Audiences:"))
			})
		})

		Context("with facts without competencies", func() {
			It("should skip competency section", func() {
				fact := fixtures.Fact("fact-1", "event-1")
				fact.Text = "Fact without competencies"
				fact.CompetencyCategories = []string{}
				factList := []*domain.Fact{fact}

				output := facts.FormatFactList(factList)
				Expect(output).NotTo(ContainSubstring("Competencies:"))
			})
		})

		Context("with multiple facts", func() {
			It("should format all facts with numbering", func() {
				fact1 := fixtures.FactWithCategories("fact-1", "First fact", "event-1", []string{"Technical"}, []string{})
				fact2 := fixtures.FactWithCategories("fact-2", "Second fact", "event-2", []string{"Leadership"}, []string{})
				factList := []*domain.Fact{fact1, fact2}

				output := facts.FormatFactList(factList)

				Expect(output).To(ContainSubstring("Total facts: 2"))
				Expect(output).To(ContainSubstring("1. First fact"))
				Expect(output).To(ContainSubstring("2. Second fact"))
				Expect(output).To(ContainSubstring("ID: fact-1"))
				Expect(output).To(ContainSubstring("ID: fact-2"))
			})
		})

		Context("with partial audience relevance", func() {
			It("should show only provided audiences", func() {
				fact := fixtures.FactWithCategories("fact-1", "Fact with single audience", "event-1", []string{}, []string{"Technical"})
				factList := []*domain.Fact{fact}

				output := facts.FormatFactList(factList)
				Expect(output).To(ContainSubstring("Target Audiences: Technical"))
			})
		})
	})
})

// Helper function to find index of substring.
func indexOf(s, substr string) int {
	return strings.Index(s, substr)
}
