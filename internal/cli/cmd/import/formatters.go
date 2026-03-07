// Pure formatter functions - extracted from display functions for testability

package importcmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/baphled/kariya/internal/cli/importer"
)

// FormatFailedRows generates formatted output for failed import rows.
//
// Expected:
//   - result: ImportResult with FailedCount and optional FailedRows
//
// Returns:
//   - Empty string if no failures
//   - Formatted string showing count and first 5 failed rows
//
// Side effects:
//   - None (pure function)
func FormatFailedRows(result *importer.ImportResult) string {
	if result.FailedCount == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Failed: %d\n", result.FailedCount))

	if len(result.FailedRows) == 0 {
		return b.String()
	}

	b.WriteString("\nFailed rows (first 5):\n")
	maxDisplay := 5
	if len(result.FailedRows) < maxDisplay {
		maxDisplay = len(result.FailedRows)
	}

	for i := range maxDisplay {
		if result.FailedRows[i].Event != nil {
			b.WriteString(fmt.Sprintf("  - %s\n", result.FailedRows[i].Event.Text))
		}
	}

	return b.String()
}

// FormatBurstSuggestions generates formatted output for burst suggestions.
//
// Expected:
//   - result: ImportResult with BurstSuggestions slice
//
// Returns:
//   - Empty string if no suggestions
//   - Formatted string with numbered list of burst suggestions
//
// Side effects:
//   - None (pure function)
func FormatBurstSuggestions(result *importer.ImportResult) string {
	if len(result.BurstSuggestions) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("\n=== Burst Suggestions ===\n")
	b.WriteString(fmt.Sprintf("Detected %d potential bursts from imported events:\n\n", len(result.BurstSuggestions)))

	for i, burst := range result.BurstSuggestions {
		burstName := burst.Name
		if burstName == "" {
			burstName = fmt.Sprintf("Burst %d", i+1)
		}
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, burstName))
		b.WriteString(fmt.Sprintf("   Events: %d | Confidence: %.1f%%\n", len(burst.EventIDs), burst.ConfidenceScore*100))
	}
	b.WriteString("\nThese bursts represent potential project groupings or themes.\n")

	return b.String()
}

// FormatFactExtraction generates formatted output for extracted facts.
//
// Expected:
//   - result: ImportResult with ExtractedFactsCount and FactsByCompetency
//
// Returns:
//   - Empty string if no facts extracted
//   - Formatted string with fact count and competency breakdown (alphabetically sorted)
//
// Side effects:
//   - None (pure function)
func FormatFactExtraction(result *importer.ImportResult) string {
	if result.ExtractedFactsCount == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("\n=== Fact Extraction ===\n")
	b.WriteString(fmt.Sprintf("Extracted %d facts from %d events\n", result.ExtractedFactsCount, len(result.CreatedEvents)))

	if len(result.FactsByCompetency) == 0 {
		return b.String()
	}

	b.WriteString("\nCompetency breakdown:\n")
	competencies := make([]string, 0, len(result.FactsByCompetency))
	for c := range result.FactsByCompetency {
		competencies = append(competencies, c)
	}
	sort.Strings(competencies)

	for _, c := range competencies {
		b.WriteString(fmt.Sprintf("  - %s: %d facts\n", c, result.FactsByCompetency[c]))
	}
	b.WriteString("\nThese facts highlight key competencies and achievements.\n")

	return b.String()
}

// CalculateExitCode determines the exit code based on import results.
//
// Expected:
//   - result: ImportResult with SuccessCount and SkippedCount
//
// Returns:
//   - 0 if any events were successfully imported or skipped
//   - 1 if no events were processed
//
// Side effects:
//   - None (pure function)
func CalculateExitCode(result *importer.ImportResult) int {
	if result.SuccessCount > 0 {
		return 0
	}

	if result.SkippedCount > 0 {
		return 0
	}

	return 1
}
