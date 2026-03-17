package generatecv

import (
	"github.com/baphled/kariya/internal/domain/career"
)

// CVGenerationCompleteMsg indicates CV generation is complete.
type CVGenerationCompleteMsg struct {
	CV    *career.CVView
	Error error
}

// TechnologiesExtractedMsg indicates technologies have been extracted from user skills.
type TechnologiesExtractedMsg struct {
	Technologies []*ExtractedTechnology
	Suggestion   *FocusAreaSuggestion
	Error        error
}

// WizardCompleteMsg indicates the configuration wizard was completed.
// Contains all configuration data from the 3-step wizard: WHO (ProfileID,
// Audience), TECH (TechFocus, Technologies, FocusArea), and FORMAT
// (SkillsFormat, SkillsLimit, CVLength). TechnologiesExtracted carries
// pre-extracted skill data when available from async extraction.
type WizardCompleteMsg struct {
	ProfileID string
	Audience  string

	TechFocus    string
	Technologies []string
	FocusArea    string

	SkillsFormat          string
	SkillsLimit           int
	CVLength              string
	TechnologiesExtracted []*ExtractedTechnology
}

// ExportCompleteMsg indicates export is complete.
type ExportCompleteMsg struct {
	Path  string
	Error error
}
