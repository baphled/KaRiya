package career

import (
	"sort"
	"strings"

	"github.com/baphled/kariya/internal/domain/career"
)

// InferCompetencyFromCategories analyzes event categories and returns the most common valid category.
// If multiple categories are tied for most common, returns the alphabetically first one.
// Returns empty string if no valid categories are found.
//
// Algorithm:
// 1. Count occurrences of each valid category across all events
// 2. Find the category with the highest count
// 3. If there's a tie, select alphabetically first category
// 4. Return empty string if no valid categories exist
//
// Valid categories are defined in career.AllowedCategories
func InferCompetencyFromCategories(eventCategories [][]string) string {
	if eventCategories == nil || len(eventCategories) == 0 {
		return ""
	}

	// Count valid categories across all events
	counts := countValidCategories(eventCategories)

	if len(counts) == 0 {
		return ""
	}

	// Find most common category, breaking ties alphabetically
	return findMostCommonCategory(counts)
}

// countValidCategories counts occurrences of each valid category across all events.
// Categories are normalized (trimmed and lowercased) before counting.
// Invalid categories are ignored.
func countValidCategories(eventCategories [][]string) map[string]int {
	counts := make(map[string]int)

	for _, categories := range eventCategories {
		for _, category := range categories {
			// Normalize: trim whitespace and lowercase
			normalized := strings.ToLower(strings.TrimSpace(category))

			// Skip empty or invalid categories
			if normalized == "" || !career.AllowedCategories[normalized] {
				continue
			}

			counts[normalized]++
		}
	}

	return counts
}

// findMostCommonCategory returns the category with the highest count.
// If multiple categories are tied, returns the alphabetically first one.
func findMostCommonCategory(counts map[string]int) string {
	if len(counts) == 0 {
		return ""
	}

	// Find maximum count
	maxCount := 0
	for _, count := range counts {
		if count > maxCount {
			maxCount = count
		}
	}

	// Collect all categories with max count
	candidates := make([]string, 0)
	for category, count := range counts {
		if count == maxCount {
			candidates = append(candidates, category)
		}
	}

	// Sort alphabetically and return first
	sort.Strings(candidates)
	return candidates[0]
}
