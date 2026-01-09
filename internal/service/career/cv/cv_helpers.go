package cv

import (
	"github.com/baphled/kariya/internal/domain/career"
)

// CVStructure represents the structure/format of a CV.
type CVStructure string

const (
	// CVStructureStandard is the traditional CV structure with Experience, Projects, Skills sections.
	CVStructureStandard CVStructure = "standard"

	// CVStructureNarrative is a language-agnostic professional format with Core Strengths, Technologies, What I Bring sections.
	CVStructureNarrative CVStructure = "narrative"
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
