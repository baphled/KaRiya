package importer

import (
	"strings"
)

// CategoryMapper maps CSV categories to KaRiya's supported competency categories
type CategoryMapper struct {
	// Mapping of keywords to competency categories
	technicalKeywords  []string
	leadershipKeywords []string
	productKeywords    []string
	consultingKeywords []string
	researchKeywords   []string
	mentoringKeywords  []string
}

// NewCategoryMapper creates a new category mapper
func NewCategoryMapper() *CategoryMapper {
	return &CategoryMapper{
		technicalKeywords: []string{
			"technical", "backend", "frontend", "architecture", "development",
			"devops", "infrastructure", "system", "performance", "database",
			"api", "integration", "automation", "deployment", "testing",
			"quality", "debugging", "coding", "engineering", "firmware",
			"hardware", "embedded", "sysadmin", "admin",
		},
		leadershipKeywords: []string{
			"leadership", "lead", "lead-engineer", "lead-developer", "senior",
			"management", "manager", "director", "strategy", "vision",
			"team-lead", "tech-lead", "stakeholders", "decision",
		},
		productKeywords: []string{
			"product", "features", "ux", "ui", "design", "customer",
			"user-experience", "product-thinking", "roadmap", "innovation",
		},
		consultingKeywords: []string{
			"consulting", "consultant", "advisory", "advice", "advise",
			"consultation", "client", "external",
		},
		researchKeywords: []string{
			"research", "analysis", "investigate", "study", "experiment",
			"prototype", "prototyping", "data", "analytics", "insight",
		},
		mentoringKeywords: []string{
			"mentoring", "mentor", "training", "enablement", "teach",
			"onboarding", "support", "guidance", "coaching", "development",
			"growth", "learning", "knowledge-sharing", "knowledge-base",
		},
	}
}

// MapCategories converts CSV categories to KaRiya competency categories
// It analyzes all provided category keywords and returns the most appropriate category
func (m *CategoryMapper) MapCategories(csvCategories []string) []string {
	if len(csvCategories) == 0 {
		return []string{}
	}

	// Convert all categories to lowercase for comparison
	lowerCategories := make([]string, len(csvCategories))
	for i, cat := range csvCategories {
		lowerCategories[i] = strings.ToLower(cat)
	}

	// Count matches for each competency category
	scores := map[string]int{
		"technical":  0,
		"leadership": 0,
		"product":    0,
		"consulting": 0,
		"research":   0,
		"mentoring":  0,
	}

	// Score each category
	combinedText := strings.Join(lowerCategories, " ")

	for _, keyword := range m.technicalKeywords {
		if strings.Contains(combinedText, keyword) {
			scores["technical"]++
		}
	}

	for _, keyword := range m.leadershipKeywords {
		if strings.Contains(combinedText, keyword) {
			scores["leadership"]++
		}
	}

	for _, keyword := range m.productKeywords {
		if strings.Contains(combinedText, keyword) {
			scores["product"]++
		}
	}

	for _, keyword := range m.consultingKeywords {
		if strings.Contains(combinedText, keyword) {
			scores["consulting"]++
		}
	}

	for _, keyword := range m.researchKeywords {
		if strings.Contains(combinedText, keyword) {
			scores["research"]++
		}
	}

	for _, keyword := range m.mentoringKeywords {
		if strings.Contains(combinedText, keyword) {
			scores["mentoring"]++
		}
	}

	// Get all categories with non-zero score, sorted by priority
	var result []string
	priorities := []string{"leadership", "mentoring", "product", "consulting", "research", "technical"}

	for _, priority := range priorities {
		if scores[priority] > 0 {
			result = append(result, priority)
		}
	}

	// If no matches found, default to technical
	if len(result) == 0 {
		result = []string{"technical"}
	}

	// Return only top 2 categories to avoid bloat
	if len(result) > 2 {
		result = result[:2]
	}

	return result
}

// TagMapper maps CSV tags to KaRiya's supported tags
type TagMapper struct {
	allowedTags []string
	// Mapping of common domain tags to allowed tags
	tagMappings map[string]string
}

// NewTagMapper creates a new tag mapper
func NewTagMapper() *TagMapper {
	return &TagMapper{
		allowedTags: []string{
			"project", "achievement", "leadership", "technical",
			"consulting", "research", "product", "mentoring",
		},
		tagMappings: map[string]string{
			// Technical-related
			"backend":      "technical",
			"frontend":     "technical",
			"php":          "technical",
			"ruby":         "technical",
			"python":       "technical",
			"javascript":   "technical",
			"golang":       "technical",
			"go":           "technical",
			"c++":          "technical",
			"c#":           "technical",
			"java":         "technical",
			"rust":         "technical",
			"devops":       "technical",
			"infrastructure": "technical",
			"sysadmin":     "technical",
			"linux":        "technical",
			"docker":       "technical",
			"kubernetes":   "technical",
			"database":     "technical",
			"sql":          "technical",
			"mongodb":      "technical",
			"api":          "technical",
			"rest":         "technical",
			"soap":         "technical",
			"grpc":         "technical",
			"ci-cd":        "technical",
			"deployment":   "technical",
			"automation":   "technical",
			"testing":      "technical",
			"quality":      "technical",
			"performance":  "technical",
			"debugging":    "technical",
			"architecture": "technical",
			"design":       "technical",
			"firmware":     "technical",
			"embedded":     "technical",
			"iot":          "technical",
			"hardware":     "technical",

			// Leadership-related
			"leadership": "leadership",
			"lead":       "leadership",
			"management": "leadership",
			"manager":    "leadership",
			"strategy":   "leadership",
			"decision":   "leadership",
			"team-lead":  "leadership",
			"tech-lead":  "leadership",

			// Product-related
			"product":    "product",
			"ux":         "product",
			"ui":         "product",
			"customer":   "product",
			"innovation": "product",
			"roadmap":    "product",

			// Consulting-related
			"consulting": "consulting",
			"advisory":   "consulting",
			"clients":    "consulting",
			"agency":     "consulting",

			// Research-related
			"research": "research",
			"analysis": "research",
			"data":     "research",
			"analytics": "research",

			// Mentoring-related
			"mentoring":      "mentoring",
			"training":       "mentoring",
			"enablement":     "mentoring",
			"onboarding":     "mentoring",
			"coaching":       "mentoring",
			"knowledge-base": "mentoring",

			// Achievement-related
			"achievement":  "achievement",
			"delivered":    "achievement",
			"completed":    "achievement",
			"successful":   "achievement",
			"milestone":    "achievement",

			// Project-related
			"project": "project",
			"contract": "project",
			"freelance": "project",
		},
	}
}

// MapTags converts CSV tags to KaRiya's supported tags
// It filters and maps domain-specific tags to allowed tags
func (m *TagMapper) MapTags(csvTags []string) []string {
	if len(csvTags) == 0 {
		return []string{}
	}

	result := make(map[string]bool) // Use map to avoid duplicates

	for _, tag := range csvTags {
		tag = strings.TrimSpace(strings.ToLower(tag))
		if tag == "" {
			continue
		}

		// Check if tag is already in allowed tags
		if m.isAllowedTag(tag) {
			result[tag] = true
			continue
		}

		// Check if tag maps to an allowed tag
		if mappedTag, exists := m.tagMappings[tag]; exists {
			result[mappedTag] = true
		}
	}

	// Convert map to slice
	var tags []string
	for tag := range result {
		tags = append(tags, tag)
	}

	// Limit to 8 tags maximum
	if len(tags) > 8 {
		tags = tags[:8]
	}

	return tags
}

// isAllowedTag checks if a tag is in the allowed tags list
func (m *TagMapper) isAllowedTag(tag string) bool {
	tag = strings.ToLower(tag)
	for _, allowed := range m.allowedTags {
		if tag == allowed {
			return true
		}
	}
	return false
}

// GetAllowedTags returns the list of allowed tags
func (m *TagMapper) GetAllowedTags() []string {
	return m.allowedTags
}

