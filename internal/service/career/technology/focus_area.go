package technology

import "github.com/baphled/kariya/internal/constants"

// FocusArea is an alias to constants.FocusArea for backward compatibility.
type FocusArea = constants.FocusArea

// Focus area constants for backward compatibility.
const (
	FocusAreaBackend   = constants.FocusAreaBackend
	FocusAreaFrontend  = constants.FocusAreaFrontend
	FocusAreaFullstack = constants.FocusAreaFullstack
	FocusAreaDevOps    = constants.FocusAreaDevOps
)

// FocusAreaSuggestion represents a suggested focus area with confidence and evidence.
type FocusAreaSuggestion struct {
	Area       FocusArea      // Suggested focus area
	Confidence float64        // Confidence level (0.0-1.0)
	Evidence   map[string]int // Category counts: {"backend": 12, "frontend": 3}
}

// Analyzer analyzes skill distributions to suggest focus areas.
type Analyzer struct{}

// AnalyzeSkills suggests focus area from skill categories.
// Analysis is weighted by event count - skills used more frequently have greater influence.
func (a *Analyzer) AnalyzeSkills(techs []*ExtractedTechnology) *FocusAreaSuggestion {
	// Handle empty case
	if len(techs) == 0 {
		return &FocusAreaSuggestion{
			Area:       FocusAreaBackend,
			Confidence: 0.0,
			Evidence:   map[string]int{},
		}
	}

	// 1. Count skills by category (for evidence)
	categoryCounts := make(map[string]int)
	for _, tech := range techs {
		categoryCounts[tech.Category]++
	}

	// 2. Weight by event count - sum event counts per category
	categoryWeights := make(map[string]int)
	totalEvents := 0
	for _, tech := range techs {
		categoryWeights[tech.Category] += tech.EventCount
		totalEvents += tech.EventCount
	}

	// 3. Calculate percentages for key categories (weighted by events)
	backendWeight := float64(categoryWeights["backend"])
	frontendWeight := float64(categoryWeights["frontend"])
	devopsWeight := float64(categoryWeights["devops"])
	totalEventsFloat := float64(totalEvents)

	backendPct := backendWeight / totalEventsFloat
	frontendPct := frontendWeight / totalEventsFloat
	devopsPct := devopsWeight / totalEventsFloat

	// 3. Determine focus area
	var area FocusArea
	var confidence float64

	if backendPct >= 0.70 {
		// Backend dominant
		area = FocusAreaBackend
		confidence = backendPct
	} else if frontendPct >= 0.70 {
		// Frontend dominant
		area = FocusAreaFrontend
		confidence = frontendPct
	} else if devopsPct >= 0.70 {
		// DevOps dominant
		area = FocusAreaDevOps
		confidence = devopsPct
	} else if categoryWeights["backend"] > 0 && categoryWeights["frontend"] > 0 {
		// Mixed backend and frontend = Fullstack
		area = FocusAreaFullstack
		// Confidence is based on how balanced they are
		// More balanced = higher confidence
		smaller := backendPct
		if frontendPct < backendPct {
			smaller = frontendPct
		}
		// If both are ~50%, confidence is high
		// If one is much larger, confidence is lower
		confidence = smaller * 2 // Scale to 0-1 range
	} else {
		// Default to backend if no clear pattern
		area = FocusAreaBackend
		confidence = backendPct
	}

	return &FocusAreaSuggestion{
		Area:       area,
		Confidence: confidence,
		Evidence:   categoryCounts,
	}
}
