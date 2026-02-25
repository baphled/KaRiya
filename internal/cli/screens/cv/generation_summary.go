package cv

import "github.com/baphled/kariya/internal/cli/types"

// GenerationSummary captures the selections and derived counts used during CV generation.
// This struct is intentionally a plain data holder used by the review/export screens.
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
