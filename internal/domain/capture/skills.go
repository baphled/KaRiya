package capture

import (
	"strings"
)

// SkillSuggestion represents an inferred skill suggestion.
// This is a domain-level representation decoupled from the service layer.
type SkillSuggestion struct {
	Name     string
	EventIDs []string
}

// FilterNewSkillSuggestions removes suggestions whose names already exist (case-insensitive).
//
// Expected: Suggestions and existingNames may be empty.
// Returns: A filtered slice containing only suggestions not in existingNames.
// Side effects: None.
func FilterNewSkillSuggestions(suggestions []SkillSuggestion, existingNames []string) []SkillSuggestion {
	if len(existingNames) == 0 {
		return suggestions
	}

	existingSet := make(map[string]bool, len(existingNames))
	for _, name := range existingNames {
		existingSet[strings.ToLower(name)] = true
	}

	filtered := make([]SkillSuggestion, 0, len(suggestions))
	for _, s := range suggestions {
		if !existingSet[strings.ToLower(s.Name)] {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
