package display

import (
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
)

// BurstSuggestion is a presentation-only view of a burst suggestion.
type BurstSuggestion struct {
	EventIDs        []string
	ConfidenceScore float64
	Name            string
	Description     string
}

// SkillSuggestion is a presentation-only view of a skill suggestion.
type SkillSuggestion struct {
	Name       string
	Category   string
	Confidence float64
	EventIDs   []string
	Contexts   []string
}

// BurstSuggestionFromDomain converts a domain burst suggestion to a display burst suggestion.
//
// Expected:
//   - burstsuggestion must be valid.
//
// Returns:
//   - A BurstSuggestion value.
//
// Side effects:
//   - None.
func BurstSuggestionFromDomain(s burstfact.BurstSuggestion) BurstSuggestion {
	return BurstSuggestion{
		EventIDs:        append([]string(nil), s.EventIDs...),
		ConfidenceScore: s.ConfidenceScore,
		Name:            s.Name,
		Description:     s.Description,
	}
}

// BurstSuggestionsFromDomain converts domain burst suggestions to display burst suggestions.
//
// Expected:
//   - burstsuggestion must be valid.
//
// Returns:
//   - A []BurstSuggestion value.
//
// Side effects:
//   - None.
func BurstSuggestionsFromDomain(suggestions []burstfact.BurstSuggestion) []BurstSuggestion {
	if suggestions == nil {
		return nil
	}

	result := make([]BurstSuggestion, len(suggestions))
	for i, suggestion := range suggestions {
		result[i] = BurstSuggestionFromDomain(suggestion)
	}

	return result
}

// SkillSuggestionFromDomain converts a domain skill suggestion to a display skill suggestion.
//
// Expected:
//   - skillsuggestion must be valid.
//
// Returns:
//   - A SkillSuggestion value.
//
// Side effects:
//   - None.
func SkillSuggestionFromDomain(s skillinference.SkillSuggestion) SkillSuggestion {
	return SkillSuggestion{
		Name:       s.Name,
		Category:   s.Category,
		Confidence: s.Confidence,
		EventIDs:   append([]string(nil), s.EventIDs...),
		Contexts:   append([]string(nil), s.Contexts...),
	}
}

// SkillSuggestionsFromDomain converts domain skill suggestions to display skill suggestions.
//
// Expected:
//   - skillsuggestion must be valid.
//
// Returns:
//   - A []SkillSuggestion value.
//
// Side effects:
//   - None.
func SkillSuggestionsFromDomain(suggestions []skillinference.SkillSuggestion) []SkillSuggestion {
	if suggestions == nil {
		return nil
	}

	result := make([]SkillSuggestion, len(suggestions))
	for i, suggestion := range suggestions {
		result[i] = SkillSuggestionFromDomain(suggestion)
	}

	return result
}
