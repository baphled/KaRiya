package cv

import (
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
)

// Structure represents the structure/format of a CV.
type Structure string

const (
	// CVStructureStandard is the traditional CV structure with Experience, Projects, Skills sections.
	CVStructureStandard Structure = "standard"

	// CVStructureNarrative is a language-agnostic professional format with Core Strengths, Technologies, What I Bring sections.
	CVStructureNarrative Structure = "narrative"

	// CVStructureConsulting emphasizes client engagements and technical capabilities for consulting roles.
	CVStructureConsulting Structure = "consulting"

	// CVStructureHighlights is a condensed 1-page format with key capabilities and selected highlights.
	CVStructureHighlights Structure = "highlights"
)

// MinConfidenceForNarrative is the minimum confidence score for bullets to appear in narrative CV.
const MinConfidenceForNarrative = 0.75

// NarrativeProfileData contains hardcoded profile data for narrative CV structure.
type NarrativeProfileData struct {
	Name              string
	Role              string
	Location          string
	Email             string
	GitHub            string
	Portfolio         string
	CoreStrengths     []string
	Languages         []string
	Frontend          []string
	Systems           []string
	ValuePropositions []string
}

// DefaultNarrativeProfile returns an empty narrative profile data.
//
// Returns:
//   - A fully initialized NarrativeProfileData ready for use.
//
// Side effects:
//   - None.
func DefaultNarrativeProfile() *NarrativeProfileData {
	return &NarrativeProfileData{
		Name:              "",
		Role:              "",
		Location:          "",
		Email:             "",
		GitHub:            "",
		Portfolio:         "",
		CoreStrengths:     []string{},
		Languages:         []string{},
		Frontend:          []string{},
		Systems:           []string{},
		ValuePropositions: []string{},
	}
}

// NarrativeProfileFromConfig creates a NarrativeProfileData from config.ProfileConfig.
//
// Expected:
//   - config must be a valid configuration object.
//
// Returns:
//   - A fully initialized NarrativeProfileData ready for use.
//
// Side effects:
//   - None.
func NarrativeProfileFromConfig(cfg *config.ProfileConfig) *NarrativeProfileData {
	if cfg == nil {
		return DefaultNarrativeProfile()
	}

	return &NarrativeProfileData{
		Name:              cfg.Name,
		Role:              cfg.Title,
		Location:          cfg.Location,
		Email:             cfg.Email,
		GitHub:            cfg.GitHub,
		Portfolio:         cfg.Portfolio,
		CoreStrengths:     cfg.CoreStrengths,
		Languages:         cfg.Languages,
		Frontend:          cfg.Frontend,
		Systems:           cfg.Systems,
		ValuePropositions: cfg.WhatIBring,
	}
}

// filterBulletsByConfidence filters bullets to only include those with confidence >= threshold.
func filterBulletsByConfidence(bullets []*career.CVBullet, threshold float64) []*career.CVBullet {
	filtered := make([]*career.CVBullet, 0)
	for _, bullet := range bullets {
		if bullet.Confidence >= threshold {
			filtered = append(filtered, bullet)
		}
	}
	return filtered
}

// getSummaryFromSections extracts the summary text from CV sections.
func getSummaryFromSections(sections []*career.CVSection) string {
	for _, section := range sections {
		if section.SectionType == "summary" && section.Summary != "" {
			return section.Summary
		}
	}
	return ""
}

// getExperienceSections returns only experience-type sections from CV.
func getExperienceSections(sections []*career.CVSection) []*career.CVSection {
	experience := make([]*career.CVSection, 0)
	for _, section := range sections {
		if section.SectionType == "experience" {
			experience = append(experience, section)
		}
	}
	return experience
}
