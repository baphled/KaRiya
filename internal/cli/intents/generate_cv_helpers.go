package intents

import (
	"github.com/baphled/kariya/internal/domain/career"
)

// NarrativeProfileData contains hardcoded profile data for narrative CV structure.
// In Phase 5, this will be configurable via Configure System → Profile.
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
// This is hardcoded for Phase 3 and will be made configurable in Phase 5.
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

// MinConfidenceForNarrative is the minimum confidence score for bullets to appear in narrative CV.
const MinConfidenceForNarrative = 0.75

// extractStrengthsFromSections attempts to extract core strengths from CV sections.
// Falls back to default strengths if none can be extracted.
func extractStrengthsFromSections(sections []*career.CVSection) []string {
	// TODO: Extract from facts in Phase 3
	// For now, return default strengths
	return DefaultNarrativeProfile().CoreStrengths
}

// extractTechnologiesFromSections attempts to extract technologies from CV sections.
// Returns languages, frontend, and systems strings.
func extractTechnologiesFromSections(sections []*career.CVSection) (languages, frontend, systems string) {
	// TODO: Extract from top bullets in Phase 3
	// For now, return defaults
	profile := DefaultNarrativeProfile()
	return profile.Languages, profile.Frontend, profile.Systems
}

// extractValuePropositions attempts to extract value propositions from CV sections.
// Falls back to default propositions if none can be extracted.
func extractValuePropositions(sections []*career.CVSection) []string {
	// TODO: Add email format validation
	// For now, return default value propositions
	return DefaultNarrativeProfile().ValuePropositions
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
// Returns empty string if no summary section exists.
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
