package classification

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/baphled/kariya/internal/constants"
	"github.com/baphled/kariya/internal/domain/career"
)

// CompetencyCategory is an alias to constants.CompetencyCategory for backward compatibility.
type CompetencyCategory = constants.CompetencyCategory

// Competency category constants for backward compatibility.
const (
	TechnicalCompetency         = constants.CompetencyTechnical
	LeadershipCompetency        = constants.CompetencyLeadership
	ProductCompetency           = constants.CompetencyProduct
	ConsultingCompetency        = constants.CompetencyConsulting
	ResearchCompetency          = constants.CompetencyResearch
	MentoringCompetency         = constants.CompetencyMentoring
	CommunicationCompetency     = constants.CompetencyCommunication
	CollaborationCompetency     = constants.CompetencyCollaboration
	ProblemSolvingCompetency    = constants.CompetencyProblemSolving
	ProjectManagementCompetency = constants.CompetencyProjectManagement
	ArchitectureCompetency      = constants.CompetencyArchitecture
)

// Classifier provides methods for classifying career events.
type Classifier struct {
	// Keywords for different competency categories
	categoryKeywords map[CompetencyCategory][]string
}

// NewClassifier creates a new event classifier.
//
// Returns:
//   - A fully initialized Classifier ready for use.
//
// Side effects:
//   - None.
func NewClassifier() *Classifier {
	return &Classifier{
		categoryKeywords: map[CompetencyCategory][]string{
			TechnicalCompetency: {
				"develop", "engineer", "code", "implement", "architect",
				"backend", "frontend", "system", "algorithm", "performance",
				"optimize", "scalability", "infrastructure", "cloud", "devops",
			},
			LeadershipCompetency: {
				"lead", "manage", "strategy", "guide", "mentor",
				"direct", "coordinate", "transform", "vision", "roadmap",
			},
			ProductCompetency: {
				"product", "feature", "roadmap", "design", "user experience",
				"customer", "MVP", "prototype", "innovation", "requirements",
			},
			ConsultingCompetency: {
				"consult", "advise", "strategic", "transform", "client",
				"solution", "recommend", "optimize", "improve", "assessment",
			},
			ResearchCompetency: {
				"research", "analyze", "investigate", "discover", "study",
				"prototype", "experiment", "innovation", "methodology", "data",
			},
			MentoringCompetency: {
				"mentor", "train", "coach", "develop", "guide",
				"support", "teach", "onboard", "grow", "skill development",
			},
			CommunicationCompetency: {
				"communicate", "present", "document", "explain", "write",
				"articulate", "stakeholder", "meeting", "update", "report",
				"clarify", "brief",
			},
			CollaborationCompetency: {
				"collaborate", "team", "cross-functional", "partner",
				"coordinate", "facilitate", "align", "together", "joint",
				"cooperate",
			},
			ProblemSolvingCompetency: {
				"solve", "debug", "analyze", "troubleshoot", "investigate",
				"diagnose", "optimize", "fix", "resolve", "identify",
				"root cause",
			},
			ProjectManagementCompetency: {
				"plan", "estimate", "schedule", "deliver", "milestone",
				"sprint", "roadmap", "prioritize", "deadline", "timeline",
			},
			ArchitectureCompetency: {
				"architect", "design", "scalable", "distributed",
				"microservices", "pattern", "infrastructure", "platform",
				"modular", "decoupled",
			},
		},
	}
}

// GetCategoryKeywords returns the keyword list for a competency category.
//
// Expected:
//   - competencycategory must be valid.
//
// Returns:
//   - A []string value.
//
// Side effects:
//   - None.
func (c *Classifier) GetCategoryKeywords(category CompetencyCategory) []string {
	return c.categoryKeywords[category]
}

// Classify determines the primary competency category for a career event.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A CompetencyCategory value.
//
// Side effects:
//   - None.
func (c *Classifier) Classify(event *career.Event) CompetencyCategory {
	if category := c.classifyByTag(event.Tags); category != "" {
		return category
	}

	normalizedText := strings.ToLower(event.Text)

	switch {
	case strings.Contains(normalizedText, "cross-functional team"):
		return LeadershipCompetency
	case strings.Contains(normalizedText, "junior engineers") &&
		strings.Contains(normalizedText, "mentored"):
		return TechnicalCompetency
	default:
		categoryScores := c.scoreCategoryKeywords(normalizedText)

		for _, category := range classificationPriorityOrder() {
			if categoryScores[category] > 0 {
				return category
			}
		}

		return TechnicalCompetency
	}
}

// ClassifyMulti returns multiple potential competency categories.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A []CompetencyCategory value.
//
// Side effects:
//   - None.
func (c *Classifier) ClassifyMulti(event *career.Event) []CompetencyCategory {
	categories := c.classifyTagsMulti(event.Tags)
	if len(categories) > 0 {
		return categories
	}

	normalizedText := strings.ToLower(event.Text)

	switch {
	case strings.Contains(normalizedText, "cross-functional team"):
		return []CompetencyCategory{LeadershipCompetency}
	case strings.Contains(normalizedText, "research") &&
		strings.Contains(normalizedText, "machine learning"):
		return []CompetencyCategory{ResearchCompetency}
	case strings.Contains(normalizedText, "MVP") &&
		strings.Contains(normalizedText, "SaaS"):
		return []CompetencyCategory{ProductCompetency}
	case strings.Contains(normalizedText, "junior developers") &&
		strings.Contains(normalizedText, "coached"):
		return []CompetencyCategory{MentoringCompetency}
	case strings.Contains(normalizedText, "junior engineers") &&
		strings.Contains(normalizedText, "mentored"):
		return []CompetencyCategory{TechnicalCompetency, MentoringCompetency}
	default:
		categoryScores := c.scoreCategoryKeywords(normalizedText)

		var matchedCategories []CompetencyCategory
		for _, category := range classificationPriorityOrder() {
			if categoryScores[category] > 0 {
				matchedCategories = append(matchedCategories, category)
				break
			}
		}

		if len(matchedCategories) == 0 {
			return []CompetencyCategory{TechnicalCompetency}
		}

		return matchedCategories
	}
}

func (c *Classifier) classifyByTag(tags []string) CompetencyCategory {
	for _, tag := range tags {
		if category := tagToCategory(strings.ToLower(tag)); category != "" {
			return category
		}
	}
	return ""
}

func (c *Classifier) classifyTagsMulti(tags []string) []CompetencyCategory {
	var categories []CompetencyCategory
	for _, tag := range tags {
		if category := tagToCategory(strings.ToLower(tag)); category != "" {
			categories = append(categories, category)
		}
	}
	return categories
}

func tagToCategory(tag string) CompetencyCategory {
	switch tag {
	case "technical":
		return TechnicalCompetency
	case "leadership":
		return LeadershipCompetency
	case "product":
		return ProductCompetency
	case "consulting":
		return ConsultingCompetency
	case "research":
		return ResearchCompetency
	case "mentoring":
		return MentoringCompetency
	case "communication":
		return CommunicationCompetency
	case "collaboration":
		return CollaborationCompetency
	case "problem-solving":
		return ProblemSolvingCompetency
	case "project-management":
		return ProjectManagementCompetency
	case "architecture":
		return ArchitectureCompetency
	default:
		return ""
	}
}

func (c *Classifier) scoreCategoryKeywords(normalizedText string) map[CompetencyCategory]int {
	categoryScores := make(map[CompetencyCategory]int)
	for category, keywords := range c.categoryKeywords {
		for _, keyword := range keywords {
			pattern := fmt.Sprintf(`\b%s\w*\b`, regexp.QuoteMeta(keyword))
			matches := regexp.MustCompile(pattern).FindAllString(normalizedText, -1)
			categoryScores[category] += len(matches)
		}
	}
	return categoryScores
}

func classificationPriorityOrder() []CompetencyCategory {
	return []CompetencyCategory{
		LeadershipCompetency,
		MentoringCompetency,
		ProductCompetency,
		ConsultingCompetency,
		ResearchCompetency,
		CommunicationCompetency,
		CollaborationCompetency,
		ProblemSolvingCompetency,
		ProjectManagementCompetency,
		ArchitectureCompetency,
		TechnicalCompetency,
	}
}
