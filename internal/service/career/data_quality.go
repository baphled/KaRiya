package career

import (
	"strings"

	"github.com/baphled/kariya/internal/domain/career"
)

// QualityLevel represents the quality level of an event's metadata
type QualityLevel string

const (
	// QualityIncomplete: 0-25 points - Missing critical fields
	QualityIncomplete QualityLevel = "Incomplete"
	// QualityBasic: 26-50 points - Has text and date, missing optional fields
	QualityBasic QualityLevel = "Basic"
	// QualityEnriched: 51-75 points - Has text, date, and some optional fields
	QualityEnriched QualityLevel = "Enriched"
	// QualityComplete: 76-100 points - Has all or most fields populated
	QualityComplete QualityLevel = "Complete"
)

// QualityScore represents the calculated quality metrics for a CareerEvent
type QualityScore struct {
	Score           int          // Total score 0-100
	Level           QualityLevel // Quality level based on score
	TextScore       int          // Text field score (0-20)
	DateScore       int          // Date field score (0-20)
	CompanyScore    int          // Company field score (0-10)
	ProjectScore    int          // Project field score (0-10)
	TagsScore       int          // Tags field score (0-15)
	CategoriesScore int          // Categories field score (0-15)
	MatchScore      int          // Text-category match score (0-10)
	MissingFields   []string     // List of missing optional fields
}

// DataQualityCalculator provides methods for calculating event metadata quality
type DataQualityCalculator struct {
}

// NewDataQualityCalculator creates a new data quality calculator
func NewDataQualityCalculator() *DataQualityCalculator {
	return &DataQualityCalculator{}
}

// CalculateQuality calculates the overall quality score for a CareerEvent
// Scoring breakdown:
// - Text (required): +20 if present and non-empty
// - Date (required): +20 if present and valid
// - Company (optional): +10 if present
// - Project (optional): +10 if present
// - Tags (optional): +15 if present (2+ tags)
// - Categories (optional): +15 if present (1+ categories)
// - Match (optional): +10 if text matches categories
func (dqc *DataQualityCalculator) CalculateQuality(event *career.CareerEvent) QualityScore {
	score := QualityScore{
		MissingFields: []string{},
	}

	// Text score (0-20): Required field
	if strings.TrimSpace(event.Text) != "" {
		score.TextScore = 20
	}

	// Date score (0-20): Required field
	// Date is always set by the service, so we give full points if it's not zero
	if !event.Date.IsZero() {
		score.DateScore = 20
	}

	// Company score (0-10): Optional field
	if strings.TrimSpace(event.Company) != "" {
		score.CompanyScore = 10
	} else {
		score.MissingFields = append(score.MissingFields, "company")
	}

	// Project score (0-10): Optional field
	if strings.TrimSpace(event.Project) != "" {
		score.ProjectScore = 10
	} else {
		score.MissingFields = append(score.MissingFields, "project")
	}

	// Tags score (0-15): Optional field
	// Give full points if 2+ tags (more descriptive)
	// Give partial points if 1 tag
	if len(event.Tags) >= 2 {
		score.TagsScore = 15
	} else if len(event.Tags) == 1 {
		score.TagsScore = 7
	} else {
		score.MissingFields = append(score.MissingFields, "tags")
	}

	// Categories score (0-15): Optional field
	if len(event.Categories) > 0 {
		score.CategoriesScore = 15
	} else {
		score.MissingFields = append(score.MissingFields, "categories")
	}

	// Match score (0-10): Text matches categories
	if len(event.Categories) > 0 && len(event.Tags) > 0 {
		// Check if text contains keywords from categories
		textLower := strings.ToLower(event.Text)
		categoryKeywords := map[string][]string{
			"technical":  {"develop", "engineer", "code", "implement", "architect", "backend", "frontend", "system", "algorithm", "build", "technical"},
			"leadership": {"lead", "manage", "strategy", "guide", "mentor", "direct", "coordinate", "transform", "vision", "roadmap", "lead"},
			"product":    {"product", "feature", "roadmap", "design", "user experience", "customer", "mvp", "prototype", "innovation", "product"},
			"consulting": {"consult", "advise", "strategic", "transform", "client", "solution", "recommend", "optimize", "advise"},
			"research":   {"research", "analyze", "investigate", "discover", "study", "prototype", "experiment", "innovation", "methodology"},
			"mentoring":  {"mentor", "train", "coach", "develop", "guide", "support", "teach", "onboard", "grow", "skill"},
		}

		matchFound := false
		for _, category := range event.Categories {
			if keywords, ok := categoryKeywords[strings.ToLower(category)]; ok {
				for _, keyword := range keywords {
					if strings.Contains(textLower, keyword) {
						matchFound = true
						break
					}
				}
				if matchFound {
					break
				}
			}
		}

		if matchFound {
			score.MatchScore = 10
		} else {
			score.MatchScore = 5
		}
	}

	// Calculate total score
	score.Score = score.TextScore + score.DateScore + score.CompanyScore + score.ProjectScore +
		score.TagsScore + score.CategoriesScore + score.MatchScore

	// Ensure score doesn't exceed 100
	if score.Score > 100 {
		score.Score = 100
	}

	// Determine quality level
	score.Level = dqc.GetQualityLevel(score.Score)

	return score
}

// GetQualityLevel returns the QualityLevel for a given score
func (dqc *DataQualityCalculator) GetQualityLevel(score int) QualityLevel {
	switch {
	case score >= 76:
		return QualityComplete
	case score >= 51:
		return QualityEnriched
	case score >= 26:
		return QualityBasic
	default:
		return QualityIncomplete
	}
}

// GetMissingFieldsSuggestion returns a user-friendly suggestion for improving quality
func (dqc *DataQualityCalculator) GetMissingFieldsSuggestion(score QualityScore) string {
	if len(score.MissingFields) == 0 {
		return "All fields are complete!"
	}

	if len(score.MissingFields) == 1 {
		return "Add " + score.MissingFields[0] + " to improve quality"
	}

	// Join all but last field with ", " and add "and" before last field
	fields := strings.Join(score.MissingFields[:len(score.MissingFields)-1], ", ")
	return "Add " + fields + " and " + score.MissingFields[len(score.MissingFields)-1] + " to improve quality"
}
