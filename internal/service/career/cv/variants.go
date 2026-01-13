package cv

// TechnologyFocus represents the technology presentation style for a CV.
type TechnologyFocus string

const (
	// TechnologyFocusLanguageAgnostic - Technology-agnostic narrative style
	TechnologyFocusLanguageAgnostic TechnologyFocus = "language_agnostic"

	// TechnologyFocusGeneralist - Generalist style highlighting 2-5 technologies
	TechnologyFocusGeneralist TechnologyFocus = "generalist"

	// TechnologyFocusSpecialist - Specialist style focused on 1 technology
	TechnologyFocusSpecialist TechnologyFocus = "specialist"
)

// FocusArea represents a career focus area.
type FocusArea string

const (
	// FocusAreaBackend - Backend development focus
	FocusAreaBackend FocusArea = "backend"

	// FocusAreaFrontend - Frontend development focus
	FocusAreaFrontend FocusArea = "frontend"

	// FocusAreaFullstack - Fullstack development focus
	FocusAreaFullstack FocusArea = "fullstack"

	// FocusAreaDevOps - DevOps/infrastructure focus
	FocusAreaDevOps FocusArea = "devops"
)

// LengthFormat represents the CV length/detail level.
type LengthFormat string

const (
	// LengthUltraShort - Ultra-short format (1 page, highlights only)
	LengthUltraShort LengthFormat = "ultra_short"

	// LengthShort - Short format (1-2 pages)
	LengthShort LengthFormat = "short"

	// LengthStandard - Standard format (2-3 pages)
	LengthStandard LengthFormat = "standard"

	// LengthFull - Full format (3+ pages, comprehensive)
	LengthFull LengthFormat = "full"
)

// CVStructure represents the CV organizational structure.
type CVStructure string

const (
	// CVStructureHighlights - Highlights-focused structure (for UltraShort)
	CVStructureHighlights CVStructure = "highlights"

	// CVStructureStandard - Standard chronological structure
	CVStructureStandard CVStructure = "standard"

	// CVStructureNarrative - Narrative/story-based structure
	CVStructureNarrative CVStructure = "narrative"
)

// CVVariant represents a specific CV configuration.
type CVVariant struct {
	ID              string          // Variant identifier (e.g., "agnostic_backend_full")
	Name            string          // Human-readable name
	Description     string          // Variant description
	TechnologyFocus TechnologyFocus // Technology presentation style
	FocusArea       FocusArea       // Career focus area
	LengthFormat    LengthFormat    // CV length/detail level
	Technologies    []string        // Selected technologies (skill IDs) - for Generalist/Specialist
	BaseStructure   CVStructure     // CV organizational structure
}

// BuildVariantID generates a variant ID from selections.
//
// Format examples:
//   - Language Agnostic: "agnostic_backend_full"
//   - Generalist: "generalist_fullstack_standard"
//   - Specialist: "specialist_ruby_backend_full"
func BuildVariantID(techFocus TechnologyFocus, focusArea FocusArea, length LengthFormat, technologies []string) string {
	// Determine focus prefix
	var focusPrefix string
	switch techFocus {
	case TechnologyFocusLanguageAgnostic:
		focusPrefix = "agnostic"
	case TechnologyFocusGeneralist:
		focusPrefix = "generalist"
	case TechnologyFocusSpecialist:
		focusPrefix = "specialist"
	default:
		focusPrefix = string(techFocus)
	}

	// Build ID based on technology focus
	if techFocus == TechnologyFocusSpecialist && len(technologies) > 0 {
		// Specialist includes technology name
		return focusPrefix + "_" + technologies[0] + "_" + string(focusArea) + "_" + string(length)
	}

	// Standard format for Language Agnostic and Generalist
	return focusPrefix + "_" + string(focusArea) + "_" + string(length)
}

// DetermineStructure determines the appropriate CV structure based on technology focus and length.
//
// Rules:
//   - UltraShort always uses Highlights structure
//   - Language Agnostic uses Narrative structure (emphasizes adaptability)
//   - Generalist and Specialist use Standard structure (traditional format)
func DetermineStructure(techFocus TechnologyFocus, length LengthFormat) CVStructure {
	// UltraShort always uses Highlights regardless of focus
	if length == LengthUltraShort {
		return CVStructureHighlights
	}

	// Technology focus determines structure for other lengths
	switch techFocus {
	case TechnologyFocusLanguageAgnostic:
		return CVStructureNarrative

	case TechnologyFocusGeneralist, TechnologyFocusSpecialist:
		return CVStructureStandard

	default:
		return CVStructureStandard
	}
}
