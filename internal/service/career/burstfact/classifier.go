package burstfact

import (
	"strings"

	"github.com/baphled/kariya/internal/domain/career"
)

// Classifier provides classification and inference for facts.
type Classifier struct{}

// NewClassifier creates a new classifier.
func NewClassifier() *Classifier {
	return &Classifier{}
}

// ClassifyRoleFit determines the role fit for a fact based on keywords.
//
// Deprecated: Use ClassifyRoleFitWithCategories instead for more accurate classification.
func (c *Classifier) ClassifyRoleFit(text string) career.RoleFit {
	return c.ClassifyRoleFitWithCategories(text, nil)
}

// ClassifyRoleFitWithCategories determines role fit using both text keywords and event categories.
// Categories provide a strong signal that should take precedence over keyword matching.
// The logic is:
//   - "leadership" + "technical" categories → staff (technical leadership)
//   - "leadership" only → em (people management)
//   - "technical" only → senior_ic (individual contributor)
//   - Keywords in text can elevate to principal if strong signals present
func (c *Classifier) ClassifyRoleFitWithCategories(text string, categories []string) career.RoleFit {
	lowerText := strings.ToLower(text)

	// First, determine base role from categories (strong signal)
	hasLeadership := false
	hasTechnical := false
	for _, cat := range categories {
		lowerCat := strings.ToLower(cat)
		if lowerCat == "leadership" {
			hasLeadership = true
		}
		if lowerCat == "technical" || lowerCat == "architecture" {
			hasTechnical = true
		}
	}

	// Category-based role determination (primary signal)
	var categoryBasedRole career.RoleFit
	if hasLeadership && hasTechnical {
		categoryBasedRole = career.RoleFitStaff
	} else if hasLeadership {
		categoryBasedRole = career.RoleFitEM
	} else {
		categoryBasedRole = career.RoleFitSeniorIC
	}

	// Only elevate to principal if STRONG principal indicators in text
	// These are indicators that suggest company-wide or strategic impact
	strongPrincipalKeywords := []string{
		"principal", "company-wide", "enterprise", "organization-wide",
		"technical direction", "founding", "founder", "cto", "chief",
	}

	principalScore := scoreText(lowerText, strongPrincipalKeywords)
	if principalScore >= 2 {
		// Multiple strong principal signals → elevate to principal
		return career.RoleFitPrincipal
	}

	// If categories provided a signal, use it
	if len(categories) > 0 {
		return categoryBasedRole
	}

	// Fallback to keyword-based classification when no categories available
	// (backwards compatibility for facts without source events)
	principalKeywords := []string{
		"principal", "architect", "vision", "strategy", "roadmap",
		"company-wide", "enterprise", "organization",
		"technical direction", "founding", "founder",
	}

	emKeywords := []string{
		"manager", "managed", "director", "head of", "vp", "vice president",
		"management", "people management", "hiring", "team lead", "team of",
	}

	staffKeywords := []string{
		"staff engineer", "principal engineer", "deep expertise",
		"complex systems", "architecture design", "cross-team", "complex",
	}

	principalScore = scoreText(lowerText, principalKeywords)
	emScore := scoreText(lowerText, emKeywords)
	staffScore := scoreText(lowerText, staffKeywords)

	// Return highest scoring role fit
	// For backwards compatibility, accept single keyword matches when no categories
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

// ClassifyAudienceRelevance determines the audience relevance for a fact.
func (c *Classifier) ClassifyAudienceRelevance(text string, _ career.RoleFit) []string {
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

// ExtractStrengthSignal extracts the strength signal from text.
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

// InferCompetencies infers competency categories from event/burst.
func (c *Classifier) InferCompetencies(text string, tags []string) []string {
	lowerText := strings.ToLower(text)
	competencies := make(map[string]bool)

	c.inferFromTags(tags, competencies)
	c.inferFromText(lowerText, competencies)

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

func (c *Classifier) inferFromTags(tags []string, competencies map[string]bool) {
	tagMap := map[string]string{
		"leadership":         "leadership",
		"technical":          "technical",
		"product":            "product",
		"consulting":         "consulting",
		"research":           "research",
		"mentoring":          "mentoring",
		"communication":      "communication",
		"collaboration":      "collaboration",
		"problem-solving":    "problem-solving",
		"project-management": "project-management",
		"architecture":       "architecture",
	}
	for _, tag := range tags {
		lowerTag := strings.ToLower(tag)
		if comp, ok := tagMap[lowerTag]; ok {
			competencies[comp] = true
		}
	}
}

func (c *Classifier) inferFromText(lowerText string, competencies map[string]bool) {
	keywordGroups := map[string][]string{
		"leadership": {
			"led", "lead", "manage", "directed", "strategy",
		},
		"technical": {
			"code", "engineer", "technical", "system",
			"backend", "frontend", "database",
		},
		"product": {
			"product", "feature", "user", "customer",
		},
		"consulting": {
			"consult", "advise", "client", "solution",
		},
		"research": {
			"research", "analyze", "investigate", "study",
		},
		"mentoring": {
			"mentor", "coach", "train", "teach", "develop", "grow",
		},
		"communication": {
			"present", "document", "explain", "stakeholder",
			"articulate", "communicate", "clarify", "brief",
		},
		"collaboration": {
			"collaborate", "cross-functional", "partner",
			"facilitate", "align", "together", "joint", "cooperate",
		},
		"problem-solving": {
			"debug", "troubleshoot", "diagnose",
			"root cause", "resolve",
		},
		"project-management": {
			"plan", "estimate", "schedule", "milestone",
			"sprint", "roadmap", "prioritize", "deadline",
		},
		"architecture": {
			"architect", "distributed", "microservices",
			"scalable", "platform", "modular", "decoupled",
		},
	}

	for comp, keywords := range keywordGroups {
		for _, kw := range keywords {
			if strings.Contains(lowerText, kw) {
				competencies[comp] = true
				break
			}
		}
	}
}

// scoreText returns a score based on keyword matches.
func scoreText(text string, keywords []string) int {
	score := 0
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			score++
		}
	}
	return score
}
