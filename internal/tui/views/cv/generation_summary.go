package cv

import "github.com/baphled/kariya/internal/ui/types"

// GenerationSummary captures the selections and derived counts used during CV generation.
// This struct is intentionally a plain data holder used by the review view.
type GenerationSummary struct {
	SelectedProfile  *types.CVProfile
	SelectedAudience string
	TechnologyFocus  string
	Technologies     []string
	FocusArea        string
	SkillsFormat     string
	SkillsLimit      int
	CVLength         string
	SourceEventCount int
	SourceFactCount  int
	SectionCount     int
	TotalBullets     int
}
