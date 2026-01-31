package skillinference

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	career_repo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/technology"
)

// SkillRepository provides data access for skill records.
type SkillRepository interface {
	Create(ctx context.Context, skill *career.Skill) error
	Update(ctx context.Context, skill *career.Skill) error
	GetByName(ctx context.Context, name string) (*career.Skill, error)
	GetByID(ctx context.Context, id string) (*career.Skill, error)
}

// EventRepository provides data access for event-skill linking.
type EventRepository interface {
	GetByID(ctx context.Context, id string) (*career.Event, error)
	Update(ctx context.Context, event *career.Event) error
	LinkSkill(ctx context.Context, eventID string, skillID string) error
}

// DefaultSkillInferenceService implements SkillInferenceService using
// keyword-based detection with word boundary regex matching.
type DefaultSkillInferenceService struct {
	skillRepo  SkillRepository
	eventRepo  EventRepository
	keywordMap map[string]technology.TechnologyKeyword
}

// NewSkillInferenceService creates a new skill inference service.
// Requires repositories for skill persistence and event linking.
func NewSkillInferenceService(skillRepo SkillRepository, eventRepo EventRepository) SkillInferenceService {
	return &DefaultSkillInferenceService{
		skillRepo:  skillRepo,
		eventRepo:  eventRepo,
		keywordMap: technology.GetKeywordMap(),
	}
}

// InferSkillsFromEvents analyzes all events for technology mentions.
func (s *DefaultSkillInferenceService) InferSkillsFromEvents(
	ctx context.Context,
	events []*career.Event,
) (*InferenceResult, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if len(events) == 0 {
		return &InferenceResult{Suggestions: []SkillSuggestion{}}, nil
	}

	suggestionMap := make(map[string]*SkillSuggestion)

	for _, event := range events {
		detectedSkills := s.detectSkillsInText(event.Text, event.ID)

		for _, detected := range detectedSkills {
			if existing, found := suggestionMap[detected.Name]; found {
				existing.EventIDs = append(existing.EventIDs, detected.EventIDs...)
				existing.Contexts = append(existing.Contexts, detected.Contexts...)

				if detected.Confidence > existing.Confidence {
					existing.Confidence = detected.Confidence
				}

				if len(existing.Contexts) > 3 {
					existing.Contexts = existing.Contexts[:3]
				}
			} else {
				suggestionMap[detected.Name] = detected
			}
		}
	}

	var existingNames []string

	if s.skillRepo != nil {
		for name := range suggestionMap {
			existing, _ := s.skillRepo.GetByName(ctx, name)
			if existing != nil {
				existingNames = append(existingNames, name)
			}
		}
	}

	suggestions := make([]SkillSuggestion, 0, len(suggestionMap))
	for _, suggestion := range suggestionMap {
		suggestions = append(suggestions, *suggestion)
	}

	return &InferenceResult{
		Suggestions:        suggestions,
		ExistingSkillNames: existingNames,
	}, nil
}

// InferSkillsFromBurst analyzes events within a specific burst.
// Delegates to InferSkillsFromEvents after filtering events.
func (s *DefaultSkillInferenceService) InferSkillsFromBurst(
	ctx context.Context,
	burst *career.Burst,
	events []*career.Event,
) (*InferenceResult, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if burst == nil || len(burst.EventIDs) == 0 {
		return &InferenceResult{Suggestions: []SkillSuggestion{}}, nil
	}

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

	return s.InferSkillsFromEvents(ctx, burstEvents)
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

// CreateSkillsFromSuggestions persists accepted suggestions as skills and links them to events.
//
// Algorithm:
//  1. Check context cancellation
//  2. Deduplicate suggestions (same name → merge event IDs)
//  3. For each unique suggestion:
//     a. Check if skill exists (case-insensitive): skillRepo.GetByName()
//     b. Create new skill if not exists, or update existing skill
//     c. Update skill.LastUsed to most recent event date
//     d. Link skill to events via eventRepo.LinkSkill()
//  4. Return created/updated skills
//
// Returns empty slice if no suggestions provided (not an error).
func (s *DefaultSkillInferenceService) CreateSkillsFromSuggestions(
	ctx context.Context,
	suggestions []SkillSuggestion,
) ([]*career.Skill, error) {
	// Check context cancellation
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if len(suggestions) == 0 {
		return []*career.Skill{}, nil
	}

	// Deduplicate suggestions by canonical name (case-insensitive)
	dedupedSuggestions := s.dedupeSuggestions(suggestions)

	var skills []*career.Skill

	for _, suggestion := range dedupedSuggestions {
		// Check context cancellation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		// Find most recent event date for LastUsed
		lastUsed, err := s.getMostRecentEventDate(ctx, suggestion.EventIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get event dates for skill %s: %w", suggestion.Name, err)
		}

		existingSkill, err := s.skillRepo.GetByName(ctx, suggestion.Name)
		if err != nil && !errors.Is(err, career_repo.ErrSkillNotFound) {
			return nil, fmt.Errorf("failed to check existing skill %s: %w", suggestion.Name, err)
		}

		var skill *career.Skill

		if existingSkill != nil {
			// Update existing skill
			existingSkill.LastUsed = lastUsed

			if err := s.skillRepo.Update(ctx, existingSkill); err != nil {
				return nil, fmt.Errorf("failed to update skill %s: %w", suggestion.Name, err)
			}

			skill = existingSkill
		} else {
			// Create new skill
			newSkill := &career.Skill{
				Name:      suggestion.Name,
				Category:  suggestion.Category,
				LastUsed:  lastUsed,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := s.skillRepo.Create(ctx, newSkill); err != nil {
				return nil, fmt.Errorf("failed to create skill %s: %w", suggestion.Name, err)
			}

			skill = newSkill
		}

		// Link skill to events
		for _, eventID := range suggestion.EventIDs {
			if err := s.eventRepo.LinkSkill(ctx, eventID, skill.ID); err != nil {
				return nil, fmt.Errorf("failed to link skill %s to event %s: %w", skill.Name, eventID, err)
			}
		}

		skills = append(skills, skill)
	}

	return skills, nil
}

// dedupeSuggestions merges suggestions with the same name (case-insensitive).
// Merges event IDs and uses the most recent event date.
func (s *DefaultSkillInferenceService) dedupeSuggestions(suggestions []SkillSuggestion) []SkillSuggestion {
	suggestionMap := make(map[string]*SkillSuggestion)

	for _, suggestion := range suggestions {
		key := strings.ToLower(suggestion.Name)

		existing, exists := suggestionMap[key]
		if !exists {
			// First occurrence - create new entry
			suggestionCopy := suggestion // Copy to avoid mutation
			suggestionMap[key] = &suggestionCopy
		} else {
			// Merge event IDs (dedupe using map)
			eventIDSet := make(map[string]bool)
			for _, id := range existing.EventIDs {
				eventIDSet[id] = true
			}
			for _, id := range suggestion.EventIDs {
				eventIDSet[id] = true
			}

			// Rebuild event IDs slice
			existing.EventIDs = make([]string, 0, len(eventIDSet))
			for id := range eventIDSet {
				existing.EventIDs = append(existing.EventIDs, id)
			}
		}
	}

	// Convert map to slice
	result := make([]SkillSuggestion, 0, len(suggestionMap))
	for _, suggestion := range suggestionMap {
		result = append(result, *suggestion)
	}

	// Sort by name for consistent ordering
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

// getMostRecentEventDate finds the most recent event date from a list of event IDs.
// Returns pointer to time (matching Skill.LastUsed type).
func (s *DefaultSkillInferenceService) getMostRecentEventDate(ctx context.Context, eventIDs []string) (*time.Time, error) {
	if len(eventIDs) == 0 {
		return nil, nil
	}

	var mostRecent time.Time
	foundAny := false

	for _, eventID := range eventIDs {
		event, err := s.eventRepo.GetByID(ctx, eventID)
		if err != nil {
			return nil, err
		}

		if !foundAny || event.Date.After(mostRecent) {
			mostRecent = event.Date
			foundAny = true
		}
	}

	return &mostRecent, nil
}
