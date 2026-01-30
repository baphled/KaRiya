package skillinference

import (
	"context"
	"regexp"
	"strings"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/technology"
)

// DefaultSkillInferenceService implements SkillInferenceService using
// keyword-based detection with word boundary regex matching.
type DefaultSkillInferenceService struct {
	keywordMap map[string]technology.TechnologyKeyword
}

// NewSkillInferenceService creates a new skill inference service.
// Uses the technology keyword dictionary for detection.
func NewSkillInferenceService() SkillInferenceService {
	return &DefaultSkillInferenceService{
		keywordMap: technology.GetKeywordMap(),
	}
}

// InferSkillsFromEvents analyzes all events for technology mentions.
// Returns skill suggestions sorted by confidence (highest first).
//
// Algorithm:
// 1. Scan each event's text for keyword matches (word boundary regex)
// 2. Extract context snippets (~80 chars around match)
// 3. Deduplicate same skill across events
// 4. Limit contexts to 3 per skill (for UI display)
//
// Returns empty slice if no skills detected (not an error).
func (s *DefaultSkillInferenceService) InferSkillsFromEvents(
	ctx context.Context,
	events []*career.Event,
) ([]SkillSuggestion, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if len(events) == 0 {
		return []SkillSuggestion{}, nil
	}

	// Map to collect suggestions by canonical skill name
	suggestionMap := make(map[string]*SkillSuggestion)

	for _, event := range events {
		// Detect skills in this event
		detectedSkills := s.detectSkillsInText(event.Text, event.ID)

		// Merge into suggestion map
		for _, detected := range detectedSkills {
			if existing, found := suggestionMap[detected.Name]; found {
				// Merge with existing suggestion
				existing.EventIDs = append(existing.EventIDs, detected.EventIDs...)
				existing.Contexts = append(existing.Contexts, detected.Contexts...)

				// Keep highest confidence
				if detected.Confidence > existing.Confidence {
					existing.Confidence = detected.Confidence
				}

				// Limit contexts to 3
				if len(existing.Contexts) > 3 {
					existing.Contexts = existing.Contexts[:3]
				}
			} else {
				// New suggestion
				suggestionMap[detected.Name] = detected
			}
		}
	}

	// Convert map to slice
	suggestions := make([]SkillSuggestion, 0, len(suggestionMap))
	for _, suggestion := range suggestionMap {
		suggestions = append(suggestions, *suggestion)
	}

	return suggestions, nil
}

// InferSkillsFromBurst analyzes events within a specific burst.
// Delegates to InferSkillsFromEvents after filtering events.
func (s *DefaultSkillInferenceService) InferSkillsFromBurst(
	ctx context.Context,
	burst *career.Burst,
	events []*career.Event,
) ([]SkillSuggestion, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if burst == nil || len(burst.EventIDs) == 0 {
		return []SkillSuggestion{}, nil
	}

	// Filter events to only those in the burst
	burstEventMap := make(map[string]bool)
	for _, id := range burst.EventIDs {
		burstEventMap[id] = true
	}

	var burstEvents []*career.Event
	for _, event := range events {
		if burstEventMap[event.ID] {
			burstEvents = append(burstEvents, event)
		}
	}

	// Delegate to InferSkillsFromEvents
	return s.InferSkillsFromEvents(ctx, burstEvents)
}

// CreateSkillsFromSuggestions will be implemented in Phase 5.
// Stub for now to satisfy interface.
func (s *DefaultSkillInferenceService) CreateSkillsFromSuggestions(
	_ context.Context,
	_ []SkillSuggestion,
) ([]*career.Skill, error) {
	// TODO: Implement in Phase 5
	return []*career.Skill{}, nil
}

// detectSkillsInText scans text for technology keywords using word boundary regex.
// Returns a slice of SkillSuggestion, one per detected keyword.
//
// Algorithm:
// 1. Convert text to lowercase for matching
// 2. For each keyword in dictionary, check word boundary match
// 3. If match found, extract context and create suggestion
// 4. Return all detected skills (deduplication happens in caller)
func (s *DefaultSkillInferenceService) detectSkillsInText(
	text string,
	eventID string,
) []*SkillSuggestion {
	lowerText := strings.ToLower(text)
	detected := []*SkillSuggestion{}

	// Check each keyword in dictionary
	for keyword, tech := range s.keywordMap {
		// Use word boundary regex to avoid partial matches
		// e.g., "goal" won't match "go", "going" won't match "go"
		pattern := `\b` + regexp.QuoteMeta(keyword) + `\b`
		re := regexp.MustCompile(pattern)

		if re.MatchString(lowerText) {
			// Extract context (up to 80 chars around match)
			contextSnippet := s.extractContext(text, keyword)

			// Calculate confidence based on usage patterns
			confidence := s.calculateConfidence(text, keyword)

			detected = append(detected, &SkillSuggestion{
				Name:       tech.Skill,
				Category:   tech.Category,
				Confidence: confidence,
				EventIDs:   []string{eventID},
				Contexts:   []string{contextSnippet},
			})
		}
	}

	return detected
}

// extractContext extracts a snippet of text around the keyword.
// Returns ~80 characters total (40 before + keyword + 40 after).
// Adds ellipsis (...) when truncated.
func (s *DefaultSkillInferenceService) extractContext(text string, keyword string) string {
	lowerText := strings.ToLower(text)
	keywordIndex := strings.Index(lowerText, keyword)

	if keywordIndex == -1 {
		return ""
	}

	// Extract 40 chars before and after keyword
	start := keywordIndex - 40
	if start < 0 {
		start = 0
	}

	end := keywordIndex + len(keyword) + 40
	if end > len(text) {
		end = len(text)
	}

	context := text[start:end]

	// Trim to complete words (don't cut mid-word)
	context = strings.TrimSpace(context)

	// Add ellipsis if truncated
	if start > 0 {
		context = "..." + context
	}
	if end < len(text) {
		context = context + "..."
	}

	return context
}

// calculateConfidence scores skill detection based on usage patterns in text.
// Returns confidence score: 0.95 (high), 0.75 (medium), or 0.5 (low).
//
// High Confidence (0.95): Active usage patterns
//   - "built" + ("with"|"using") + keyword nearby
//   - "developed" + ("in"|"using") + keyword nearby
//   - "implemented" + ("in"|"using") + keyword nearby
//   - "wrote" + keyword, "using" + keyword
//   - keyword + ("developer"|"engineer")
//   - ("expert"|"proficient") + "in" + keyword
//
// Medium Confidence (0.75): Passive or project-related patterns
//   - "worked with" + keyword
//   - "experience with" + keyword
//   - keyword + ("project"|"system"|"application"|"service")
//   - ("migrated to"|"integrated") + keyword
//
// Low Confidence (0.5): Simple keyword presence without context
func (s *DefaultSkillInferenceService) calculateConfidence(text string, keyword string) float64 {
	lowerText := strings.ToLower(text)

	// High confidence patterns (0.95)
	// Check for action verbs + keyword in flexible positions
	if s.containsPattern(lowerText, []string{"built", "with", keyword}) ||
		s.containsPattern(lowerText, []string{"built", "using", keyword}) ||
		s.containsPattern(lowerText, []string{"developed", "in", keyword}) ||
		s.containsPattern(lowerText, []string{"developed", "using", keyword}) ||
		s.containsPattern(lowerText, []string{"implemented", "in", keyword}) ||
		s.containsPattern(lowerText, []string{"implemented", "using", keyword}) ||
		strings.Contains(lowerText, "wrote "+keyword) ||
		strings.Contains(lowerText, "using "+keyword) ||
		strings.Contains(lowerText, keyword+" developer") ||
		strings.Contains(lowerText, keyword+" engineer") ||
		s.containsPattern(lowerText, []string{"expert", "in", keyword}) ||
		s.containsPattern(lowerText, []string{"proficient", "in", keyword}) {
		return 0.95
	}

	// Medium confidence patterns (0.75)
	if s.containsPattern(lowerText, []string{"worked", "with", keyword}) ||
		s.containsPattern(lowerText, []string{"working", "with", keyword}) ||
		s.containsPattern(lowerText, []string{"experience", "with", keyword}) ||
		strings.Contains(lowerText, keyword+" project") ||
		strings.Contains(lowerText, keyword+" system") ||
		strings.Contains(lowerText, keyword+" application") ||
		strings.Contains(lowerText, keyword+" service") ||
		strings.Contains(lowerText, "migrated to "+keyword) ||
		strings.Contains(lowerText, "integrated "+keyword) {
		return 0.75
	}

	// Low confidence - simple presence (0.5)
	return 0.5
}

// containsPattern checks if text contains all words in the pattern (in order, with reasonable proximity).
// Words must appear within maxWordsApart (default 3) of each other to match.
// Example: containsPattern("built API with Go", ["built", "with", "go"]) returns true
// Example: containsPattern("built API with PostgreSQL after working with Go", ["built", "with", "go"]) returns false (too far apart)
func (s *DefaultSkillInferenceService) containsPattern(text string, words []string) bool {
	const maxWordsApart = 3 // Maximum number of words allowed between pattern words

	lastIndex := -1
	for i, word := range words {
		searchStart := lastIndex + 1
		index := strings.Index(text[searchStart:], word)
		if index == -1 {
			return false
		}

		// Check proximity: count words between last match and current match
		if i > 0 {
			betweenText := text[lastIndex+len(words[i-1]) : searchStart+index]
			wordsBetween := len(strings.Fields(betweenText))
			if wordsBetween > maxWordsApart {
				return false
			}
		}

		lastIndex = searchStart + index
	}
	return true
}
