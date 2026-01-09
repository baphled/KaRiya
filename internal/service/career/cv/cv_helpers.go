package cv

import (
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
)

// CVStructure represents the structure/format of a CV.
type CVStructure string

const (
	// CVStructureStandard is the traditional CV structure with Experience, Projects, Skills sections.
	CVStructureStandard CVStructure = "standard"

	// CVStructureNarrative is a language-agnostic professional format with Core Strengths, Technologies, What I Bring sections.
	CVStructureNarrative CVStructure = "narrative"

	// CVStructureConsulting emphasizes client engagements and technical capabilities for consulting roles.
	CVStructureConsulting CVStructure = "consulting"

	// CVStructureHighlights is a condensed 1-page format with key capabilities and selected highlights.
	CVStructureHighlights CVStructure = "highlights"
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
	Languages         string
	Frontend          string
	Systems           string
	ValuePropositions []string
}

// DefaultNarrativeProfile returns the default narrative profile data.
// This is used as fallback when no profile is configured.
func DefaultNarrativeProfile() *NarrativeProfileData {
	return &NarrativeProfileData{
		Name:      "Yomi Colledge",
		Role:      "Senior Software Engineer / Technical Consultant",
		Location:  "Remote (UK)",
		Email:     "yomi@boodah.net",
		GitHub:    "https://github.com/baphled",
		Portfolio: "http://boodah.net",
		CoreStrengths: []string{
			"Language-agnostic backend and systems engineering",
			"System design and architectural ownership",
			"Cross-functional collaboration and mentorship",
			"Pragmatic problem-solving with production focus",
			"Technical leadership without formal authority",
		},
		Languages: "Ruby, Go, PHP, C/C++, JavaScript, Shell",
		Frontend:  "Vue.js",
		Systems:   "Linux, SQL, APIs, CI/CD, automation",
		ValuePropositions: []string{
			"Languages as tools, not identity",
			"Calm handling of complexity",
			"Ownership without ego",
			"Clear technical communication",
			"Production-first mindset",
		},
	}
}

// NarrativeProfileFromConfig creates a NarrativeProfileData from config.ProfileConfig.
// Falls back to default values for any empty fields.
func NarrativeProfileFromConfig(cfg *config.ProfileConfig) *NarrativeProfileData {
	if cfg == nil {
		return DefaultNarrativeProfile()
	}

	defaults := DefaultNarrativeProfile()

	profile := &NarrativeProfileData{
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

	// Use defaults for empty fields
	if profile.Name == "" {
		profile.Name = defaults.Name
	}
	if profile.Role == "" {
		profile.Role = defaults.Role
	}
	if profile.Location == "" {
		profile.Location = defaults.Location
	}
	if profile.Email == "" {
		profile.Email = defaults.Email
	}
	if profile.GitHub == "" {
		profile.GitHub = defaults.GitHub
	}
	if profile.Portfolio == "" {
		profile.Portfolio = defaults.Portfolio
	}
	if len(profile.CoreStrengths) == 0 {
		profile.CoreStrengths = defaults.CoreStrengths
	}
	if profile.Languages == "" {
		profile.Languages = defaults.Languages
	}
	if profile.Frontend == "" {
		profile.Frontend = defaults.Frontend
	}
	if profile.Systems == "" {
		profile.Systems = defaults.Systems
	}
	if len(profile.ValuePropositions) == 0 {
		profile.ValuePropositions = defaults.ValuePropositions
	}

	return profile
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
