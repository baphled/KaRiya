package facts

import (
	"fmt"
	"sort"
	"strings"

	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	domain "github.com/baphled/kariya/internal/domain/career"
)

// FormatExtractionResults generates formatted output for fact extraction results.
//
// Expected:
//   - factCount: Number of facts extracted
//   - eventCount: Number of events processed
//   - competencyCount: Map of competency categories to fact counts
//
// Returns:
//   - Formatted string with extraction summary and competency breakdown
//   - If competencyCount is empty, returns summary without breakdown section
//
// Side effects:
//   - None (pure function)
func FormatExtractionResults(factCount, eventCount int, competencyCount map[string]int) string {
	var b strings.Builder

	b.WriteString("\n=== Fact Extraction Results ===\n")
	fmt.Fprintf(&b, "Extracted %d facts from %d events\n\n", factCount, eventCount)

	if len(competencyCount) == 0 {
		return b.String()
	}

	b.WriteString("Competency breakdown:\n")
	competencies := make([]string, 0, len(competencyCount))
	for c := range competencyCount {
		competencies = append(competencies, c)
	}
	sort.Strings(competencies)

	for _, c := range competencies {
		fmt.Fprintf(&b, "  - %s: %d facts\n", c, competencyCount[c])
	}
	b.WriteString("\n")

	return b.String()
}

// FormatFactList generates formatted output for listing existing facts.
//
// Expected:
//   - facts: Slice of domain.Fact pointers to format
//
// Returns:
//   - Formatted string with title, count, and detailed fact information
//   - Each fact includes: text, ID, source event, competencies, role fit, audiences, creation time
//
// Side effects:
//   - None (pure function)
func FormatFactList(facts []*domain.Fact) string {
	var b strings.Builder

	th := theme.Default()
	b.WriteString(primitives.Title("Existing Facts", th).Render())
	b.WriteString("\n")
	fmt.Fprintf(&b, "Total facts: %d\n\n", len(facts))

	for i, fact := range facts {
		fmt.Fprintf(&b, "%d. %s\n", i+1, fact.Text)
		fmt.Fprintf(&b, "   ID: %s\n", fact.ID)
		fmt.Fprintf(&b, "   Source Event: %s\n", fact.SourceEventID)
		if len(fact.CompetencyCategories) > 0 {
			fmt.Fprintf(&b, "   Competencies: %s\n", strings.Join(fact.CompetencyCategories, ", "))
		}
		if fact.RoleFit != "" {
			fmt.Fprintf(&b, "   Role Fit: %s\n", fact.RoleFit)
		}
		if len(fact.AudienceRelevance) > 0 {
			fmt.Fprintf(&b, "   Target Audiences: %s\n", strings.Join(fact.AudienceRelevance, ", "))
		}
		fmt.Fprintf(&b, "   Created: %s\n", fact.CreatedAt.Format("2006-01-02 15:04:05"))
		b.WriteString("\n")
	}

	return b.String()
}
