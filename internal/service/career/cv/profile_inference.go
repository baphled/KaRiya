package cv

import (
	"strings"

	"github.com/baphled/kariya/internal/constants"
	"github.com/baphled/kariya/internal/domain/career"
)

// TechnologyInference holds inferred technology slices for CV profile.
type TechnologyInference struct {
	Languages []string
	Frontend  []string
	Systems   []string
}

// ProfileInferenceService infers profile data (CoreStrengths, ValuePropositions, Technologies)
// from the user's career data (events, facts, skills).
type ProfileInferenceService struct{}

// NewProfileInferenceService creates a new profile inference service.
func NewProfileInferenceService() *ProfileInferenceService {
	return &ProfileInferenceService{}
}

// categoryStrengthMapping maps competency categories to core strength descriptions.
var categoryStrengthMapping = map[string]string{
	"technical":          "Strong technical problem-solving and implementation skills",
	"leadership":         "Technical leadership and team guidance",
	"product":            "Product-oriented thinking and delivery focus",
	"consulting":         "Client-facing consulting and advisory capabilities",
	"research":           "Research-driven approach to problem solving",
	"mentoring":          "Mentoring and knowledge sharing",
	"communication":      "Strong communication and documentation skills with stakeholder engagement",
	"collaboration":      "Effective cross-functional collaboration and team leadership",
	"problem-solving":    "Analytical and systematic problem-solving abilities",
	"project-management": "Project planning and delivery management with milestone tracking",
	"architecture":       "System architecture and technical design expertise with scalability focus",
}

// categoryValueMapping maps competency categories to value proposition descriptions.
var categoryValueMapping = map[string]string{
	"technical":          "Deep technical expertise with production focus",
	"leadership":         "Cross-functional collaboration and team leadership",
	"product":            "Outcome-driven delivery with business alignment",
	"consulting":         "Strategic consulting and stakeholder management",
	"research":           "Innovation through research and experimentation",
	"mentoring":          "Team growth through mentoring and coaching",
	"communication":      "Clear technical communication with diverse audiences",
	"collaboration":      "Effective cross-functional partnerships and teamwork",
	"problem-solving":    "Systematic debugging and root-cause analysis",
	"project-management": "Reliable project delivery with transparent planning",
	"architecture":       "Scalable system design with long-term technical vision",
}

// genericCoreStrengths provides sensible defaults when no data is available.
var genericCoreStrengths = []string{
	"Problem-solving and technical implementation",
	"Cross-functional collaboration",
	"Continuous learning and adaptation",
}

// genericValuePropositions provides sensible defaults when no data is available.
var genericValuePropositions = []string{
	"Pragmatic problem-solving",
	"Clear technical communication",
	"Ownership and accountability",
}

// InferCoreStrengths analyzes career data to infer core strengths.
// Returns a list of 3-5 strengths based on the user's competency categories.
func (s *ProfileInferenceService) InferCoreStrengths(
	events []*career.Event,
	facts []*career.Fact,
	skills []*career.Skill,
) []string {
	// Count category occurrences
	categoryCounts := s.countCategories(events, facts)

	// If no data, return generic strengths
	if len(categoryCounts) == 0 && len(skills) == 0 {
		return genericCoreStrengths
	}

	// Build strengths based on category frequency
	strengths := []string{}

	// Sort categories by count and pick top ones
	sortedCategories := s.sortByCount(categoryCounts)
	for _, cat := range sortedCategories {
		if strength, ok := categoryStrengthMapping[cat]; ok {
			strengths = append(strengths, strength)
			if len(strengths) >= 4 {
				break
			}
		}
	}

	// Add skill-based strength if we have backend skills
	if s.hasExpertiseIn(skills, string(constants.SkillCategoryBackend)) {
		strengths = append(strengths, "Backend and systems engineering expertise")
	}

	// Ensure we have at least 3 strengths
	for len(strengths) < 3 && len(genericCoreStrengths) > len(strengths) {
		strengths = append(strengths, genericCoreStrengths[len(strengths)])
	}

	// Limit to 5
	if len(strengths) > 5 {
		strengths = strengths[:5]
	}

	return strengths
}

// InferValuePropositions analyzes career data to infer value propositions.
// Returns a list of 3-5 value propositions based on the user's work style and competencies.
func (s *ProfileInferenceService) InferValuePropositions(
	events []*career.Event,
	facts []*career.Fact,
	skills []*career.Skill,
) []string {
	propositions := []string{}

	// Count categories
	categoryCounts := s.countCategories(events, facts)

	// Check for collaboration indicators
	if s.hasCollaborationIndicators(events, facts) {
		propositions = append(propositions, "Cross-functional collaboration and teamwork")
	}

	// Check for mentoring indicators
	if categoryCounts["mentoring"] > 0 {
		propositions = append(propositions, "Growing teams through mentoring and knowledge sharing")
	}

	// Check for language diversity (polyglot indicator)
	if s.hasLanguageDiversity(skills) {
		propositions = append(propositions, "Language-agnostic problem solving - adaptable to any stack")
	}

	// Add based on strong categories
	sortedCategories := s.sortByCount(categoryCounts)
	for _, cat := range sortedCategories {
		if prop, ok := categoryValueMapping[cat]; ok {
			// Avoid duplicates
			if !contains(propositions, prop) {
				propositions = append(propositions, prop)
			}
		}
		if len(propositions) >= 4 {
			break
		}
	}

	// Ensure we have at least 3 propositions
	for len(propositions) < 3 {
		// Add generic ones that aren't duplicates
		for _, generic := range genericValuePropositions {
			if !contains(propositions, generic) {
				propositions = append(propositions, generic)
				break
			}
		}
		if len(propositions) < 3 && len(genericValuePropositions) > 0 {
			// Fallback if all generics are duplicates
			propositions = append(propositions, genericValuePropositions[0])
		}
	}

	// Limit to 5
	if len(propositions) > 5 {
		propositions = propositions[:5]
	}

	return propositions
}

// InferTechnologies extracts technology information from skills.
// Returns structured technology strings for Languages, Frontend, and Systems.
func (s *ProfileInferenceService) InferTechnologies(
	_ []*career.Event,
	skills []*career.Skill,
) TechnologyInference {
	result := TechnologyInference{}

	if len(skills) == 0 {
		return result
	}

	languages := []string{}
	frontend := []string{}
	systems := []string{}

	for _, skill := range skills {
		switch skill.Category {
		case string(constants.SkillCategoryBackend):
			languages = append(languages, skill.Name)
		case string(constants.SkillCategoryFrontend):
			frontend = append(frontend, skill.Name)
		case string(constants.SkillCategoryDatabase),
			string(constants.SkillCategoryDevOps),
			string(constants.SkillCategoryCloud),
			string(constants.SkillCategoryTooling):
			systems = append(systems, skill.Name)
		}
	}

	result.Languages = languages
	result.Frontend = frontend
	result.Systems = systems

	return result
}

// countCategories counts category occurrences from events and facts.
func (s *ProfileInferenceService) countCategories(
	events []*career.Event,
	facts []*career.Fact,
) map[string]int {
	counts := make(map[string]int)

	for _, event := range events {
		for _, cat := range event.Categories {
			counts[cat]++
		}
	}

	for _, fact := range facts {
		for _, cat := range fact.CompetencyCategories {
			counts[cat]++
		}
	}

	return counts
}

// sortByCount returns categories sorted by count (highest first).
func (s *ProfileInferenceService) sortByCount(counts map[string]int) []string {
	// Simple bubble sort for small maps
	categories := make([]string, 0, len(counts))
	for cat := range counts {
		categories = append(categories, cat)
	}

	for i := 0; i < len(categories)-1; i++ {
		for j := i + 1; j < len(categories); j++ {
			if counts[categories[j]] > counts[categories[i]] {
				categories[i], categories[j] = categories[j], categories[i]
			}
		}
	}

	return categories
}

// hasExpertiseIn checks if user has skills in a specific category.
func (s *ProfileInferenceService) hasExpertiseIn(skills []*career.Skill, category string) bool {
	for _, skill := range skills {
		if skill.Category == category {
			return true
		}
	}
	return false
}

// hasCollaborationIndicators checks for collaboration-related content in events/facts.
func (s *ProfileInferenceService) hasCollaborationIndicators(
	events []*career.Event,
	facts []*career.Fact,
) bool {
	collaborationKeywords := []string{"collaborat", "cross-functional", "team", "stakeholder", "partner"}

	for _, event := range events {
		lower := strings.ToLower(event.Text)
		for _, keyword := range collaborationKeywords {
			if strings.Contains(lower, keyword) {
				return true
			}
		}
		for _, cat := range event.Categories {
			if cat == "leadership" || cat == "product" || cat == "collaboration" {
				return true
			}
		}
	}

	for _, fact := range facts {
		lower := strings.ToLower(fact.Text)
		for _, keyword := range collaborationKeywords {
			if strings.Contains(lower, keyword) {
				return true
			}
		}
	}

	return false
}

// hasLanguageDiversity checks if user knows multiple programming languages.
func (s *ProfileInferenceService) hasLanguageDiversity(skills []*career.Skill) bool {
	backendCount := 0
	for _, skill := range skills {
		if skill.Category == string(constants.SkillCategoryBackend) {
			backendCount++
		}
	}
	return backendCount >= 3
}

// contains checks if a slice contains a string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
