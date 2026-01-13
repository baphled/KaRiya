package technology

// FocusArea represents a career focus area based on skill categories.
type FocusArea string

const (
	// FocusAreaBackend indicates backend development focus
	FocusAreaBackend FocusArea = "backend"

	// FocusAreaFrontend indicates frontend development focus
	FocusAreaFrontend FocusArea = "frontend"

	// FocusAreaFullstack indicates fullstack development focus
	FocusAreaFullstack FocusArea = "fullstack"

	// FocusAreaDevOps indicates DevOps/infrastructure focus
	FocusAreaDevOps FocusArea = "devops"
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
// TODO: Implement analysis logic (RED phase stub)
func (a *Analyzer) AnalyzeSkills(techs []*ExtractedTechnology) *FocusAreaSuggestion {
	return nil
}
