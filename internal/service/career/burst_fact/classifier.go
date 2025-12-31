package burst_fact

import (
	"strings"

	"github.com/baphled/kariya/internal/domain/career"
)

// Classifier provides classification and inference for facts
type Classifier struct{}

// NewClassifier creates a new classifier
func NewClassifier() *Classifier {
	return &Classifier{}
}

// ClassifyRoleFit determines the role fit for a fact based on keywords
func (c *Classifier) ClassifyRoleFit(text string) career.RoleFit {
	lowerText := strings.ToLower(text)

	// Principal indicators
	principalKeywords := []string{
		"principal", "architect", "vision", "strategy", "roadmap",
		"company-wide", "enterprise", "organization",
		"technical direction", "founding", "founder",
	}

	// EM indicators
	emKeywords := []string{
		"manager", "director", "head", "vp", "vice president",
		"management", "people management", "hiring", "team",
	}

	// Staff indicators
	staffKeywords := []string{
		"staff engineer", "principal engineer", "deep expertise",
		"complex", "difficult", "systems", "architecture design",
	}

	// Score each role fit
	principalScore := scoreText(lowerText, principalKeywords)
	emScore := scoreText(lowerText, emKeywords)
	staffScore := scoreText(lowerText, staffKeywords)

	// Return highest scoring role fit (principal > em > staff > senior_ic)
	if principalScore > 0 && principalScore >= emScore && principalScore >= staffScore {
		return career.RoleFitPrincipal
	}
	if emScore > 0 && emScore > staffScore {
		return career.RoleFitEM
	}
	if staffScore > 0 {
		return career.RoleFitStaff
	}
	return career.RoleFitSeniorIC
}

// ClassifyAudienceRelevance determines the audience relevance for a fact
func (c *Classifier) ClassifyAudienceRelevance(text string, roleFit career.RoleFit) []string {
	var audiences []string

	// All facts are relevant to peers
	audiences = append(audiences, "peer")

	// Determine relevance based on role fit and content
	lowerText := strings.ToLower(text)

	// Leadership/management facts relevant to hiring managers and recruiters
	if strings.Contains(lowerText, "lead") || strings.Contains(lowerText, "manage") ||
		strings.Contains(lowerText, "team") || strings.Contains(lowerText, "mentor") {
		audiences = append(audiences, "hiring_manager")
		audiences = append(audiences, "recruiter")
	} else {
		// Technical facts relevant to hiring managers
		audiences = append(audiences, "hiring_manager")
	}

	return audiences
}

// ExtractStrengthSignal extracts the strength signal from text
func (c *Classifier) ExtractStrengthSignal(text string) string {
	// Look for impact indicators
	impactKeywords := map[string]string{
		"delivered":   "delivery capability",
		"shipped":     "execution excellence",
		"led":         "leadership",
		"managed":     "management",
		"architected": "technical architecture",
		"designed":    "design thinking",
		"optimized":   "optimization",
		"improved":    "improvement mindset",
		"reduced":     "efficiency focus",
		"increased":   "growth orientation",
		"scaled":      "scalability expertise",
		"mentored":    "mentoring ability",
		"built":       "building capability",
	}

	lowerText := strings.ToLower(text)
	words := strings.Fields(lowerText)

	// Find the strongest signal
	for _, word := range words {
		cleanWord := strings.Trim(word, ".,!?;:")
		if signal, found := impactKeywords[cleanWord]; found {
			return signal
		}
	}

	// Default signals based on content
	if strings.Contains(lowerText, "technical") {
		return "technical expertise"
	}
	if strings.Contains(lowerText, "team") {
		return "team collaboration"
	}
	if strings.Contains(lowerText, "project") {
		return "project execution"
	}

	return "professional accomplishment"
}

// InferCompetencies infers competency categories from event/burst
func (c *Classifier) InferCompetencies(text string, tags []string) []string {
	lowerText := strings.ToLower(text)
	competencies := make(map[string]bool)

	// Check explicit tags first
	for _, tag := range tags {
		lowerTag := strings.ToLower(tag)
		if strings.Contains(lowerTag, "leadership") {
			competencies["leadership"] = true
		}
		if strings.Contains(lowerTag, "technical") {
			competencies["technical"] = true
		}
		if strings.Contains(lowerTag, "product") {
			competencies["product"] = true
		}
		if strings.Contains(lowerTag, "consulting") {
			competencies["consulting"] = true
		}
		if strings.Contains(lowerTag, "research") {
			competencies["research"] = true
		}
		if strings.Contains(lowerTag, "mentoring") {
			competencies["mentoring"] = true
		}
	}

	// Keyword-based inference
	if strings.Contains(lowerText, "led") || strings.Contains(lowerText, "lead") || strings.Contains(lowerText, "manage") ||
		strings.Contains(lowerText, "directed") || strings.Contains(lowerText, "strategy") {
		competencies["leadership"] = true
	}

	if strings.Contains(lowerText, "code") || strings.Contains(lowerText, "engineer") ||
		strings.Contains(lowerText, "technical") || strings.Contains(lowerText, "architect") ||
		strings.Contains(lowerText, "system") || strings.Contains(lowerText, "backend") ||
		strings.Contains(lowerText, "frontend") || strings.Contains(lowerText, "database") {
		competencies["technical"] = true
	}

	if strings.Contains(lowerText, "product") || strings.Contains(lowerText, "feature") ||
		strings.Contains(lowerText, "user") || strings.Contains(lowerText, "customer") {
		competencies["product"] = true
	}

	if strings.Contains(lowerText, "consult") || strings.Contains(lowerText, "advise") ||
		strings.Contains(lowerText, "client") || strings.Contains(lowerText, "solution") {
		competencies["consulting"] = true
	}

	if strings.Contains(lowerText, "research") || strings.Contains(lowerText, "analyze") ||
		strings.Contains(lowerText, "investigate") || strings.Contains(lowerText, "study") {
		competencies["research"] = true
	}

	if strings.Contains(lowerText, "mentor") || strings.Contains(lowerText, "coach") ||
		strings.Contains(lowerText, "train") || strings.Contains(lowerText, "teach") ||
		strings.Contains(lowerText, "develop") || strings.Contains(lowerText, "grow") {
		competencies["mentoring"] = true
	}

	// Default to technical if nothing found
	if len(competencies) == 0 {
		competencies["technical"] = true
	}

	// Convert map to slice
	var result []string
	for comp := range competencies {
		result = append(result, comp)
	}

	return result
}

// scoreText returns a score based on keyword matches
func scoreText(text string, keywords []string) int {
	score := 0
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			score++
		}
	}
	return score
}
