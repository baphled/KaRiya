package career

import (
	"sort"
	"strings"

	"github.com/baphled/kariya/internal/domain/career"
)

// StrengthWeight maps fact strength signals to numeric weights
var StrengthWeight = map[string]float64{
	"strong":   3.0,
	"moderate": 2.0,
	"weak":     1.0,
	"":         1.0, // Default weight for facts without strength signal
}

// InferCompetencyFromFacts analyzes facts and returns the competency with highest total weight.
// Weights are based on fact strength signals (strong=3.0, moderate=2.0, weak=1.0).
// If multiple competencies are tied, returns the alphabetically first one.
// Returns empty string if no valid competencies are found.
//
// Algorithm:
// 1. For each fact, extract competency categories
// 2. Apply weight based on strength signal
// 3. Sum weights per competency
// 4. Return competency with highest total weight
// 5. Break ties alphabetically
func InferCompetencyFromFacts(facts []*career.Fact) string {
	if facts == nil || len(facts) == 0 {
		return ""
	}

	// Calculate weighted sums per competency
	weights := calculateCompetencyWeights(facts)

	if len(weights) == 0 {
		return ""
	}

	// Find competency with highest weight, breaking ties alphabetically
	return findHighestWeightedCompetency(weights)
}

// InferCompetencyFromCategories analyzes event categories and returns the most common valid category.
// DEPRECATED: Use InferCompetencyFromFacts for better accuracy based on fact weights.
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

// calculateCompetencyWeights calculates total weight for each competency across all facts.
// Weights are based on fact strength signals.
func calculateCompetencyWeights(facts []*career.Fact) map[string]float64 {
	weights := make(map[string]float64)

	for _, fact := range facts {
		// Get weight for this fact based on strength signal
		weight := StrengthWeight[strings.ToLower(fact.StrengthSignal)]

		// Add weight to each competency category in this fact
		for _, category := range fact.CompetencyCategories {
			// Normalize: trim whitespace and lowercase
			normalized := strings.ToLower(strings.TrimSpace(category))

			// Skip empty or invalid categories
			if normalized == "" || !career.AllowedCategories[normalized] {
				continue
			}

			weights[normalized] += weight
		}
	}

	return weights
}

// findHighestWeightedCompetency returns the competency with the highest total weight.
// If multiple competencies are tied, returns the alphabetically first one.
func findHighestWeightedCompetency(weights map[string]float64) string {
	if len(weights) == 0 {
		return ""
	}

	// Find maximum weight
	maxWeight := 0.0
	for _, weight := range weights {
		if weight > maxWeight {
			maxWeight = weight
		}
	}

	// Collect all competencies with max weight
	candidates := make([]string, 0)
	for competency, weight := range weights {
		if weight == maxWeight {
			candidates = append(candidates, competency)
		}
	}

	// Sort alphabetically and return first
	sort.Strings(candidates)
	return candidates[0]
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
