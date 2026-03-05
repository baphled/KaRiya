package harness

import (
	"context"
	"sort"
	"strings"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
)

// StubSkillInferenceService provides deterministic skill inference for BDD tests.
// It avoids async operations and returns predictable results based on event text.
type StubSkillInferenceService struct {
	skillRepo       careerrepo.SkillRepository
	eventRepo       careerrepo.EventRepository
	inferenceCalled bool
}

// NewStubSkillInferenceService creates a deterministic stub for testing.
//
// Expected:
//   - skillRepo must be a valid SkillRepository implementation.
//   - eventRepo must be a valid EventRepository implementation.
//
// Returns:
//   - A configured StubSkillInferenceService instance.
//
// Side effects:
//   - None.
func NewStubSkillInferenceService(
	skillRepo careerrepo.SkillRepository,
	eventRepo careerrepo.EventRepository,
) *StubSkillInferenceService {
	return &StubSkillInferenceService{
		skillRepo: skillRepo,
		eventRepo: eventRepo,
	}
}

type skillDef struct {
	key        string
	name       string
	category   string
	confidence float64
	context    string
}

func buildSkillMap() map[string]skillinference.SkillSuggestion {
	defs := []skillDef{
		{"go", "Go", "backend", 0.95, "Built service in Go"},
		{"postgresql", "PostgreSQL", "database", 0.90, "Used PostgreSQL for storage"},
		{"postgres", "PostgreSQL", "database", 0.90, "Used PostgreSQL for storage"},
		{"python", "Python", "backend", 0.92, "Wrote Python scripts"},
		{"javascript", "JavaScript", "frontend", 0.88, "Implemented features in JavaScript"},
		{"react", "React", "frontend", 0.93, "Built UI with React"},
		{"kubernetes", "Kubernetes", "devops", 0.91, "Deployed on Kubernetes"},
		{"docker", "Docker", "devops", 0.89, "Containerized with Docker"},
		{"aws", "AWS", "cloud", 0.94, "Deployed to AWS"},
		{"api", "API Development", "backend", 0.87, "Implemented API platform features"},
		{"backend", "Backend Development", "backend", 0.85, "Built backend services"},
	}

	skillMap := make(map[string]skillinference.SkillSuggestion, len(defs))
	for _, d := range defs {
		skillMap[d.key] = skillinference.SkillSuggestion{
			Name:       d.name,
			Category:   d.category,
			Confidence: d.confidence,
			Contexts:   []string{d.context},
		}
	}
	return skillMap
}

// InferSkillsFromEvents returns deterministic skill suggestions based on event text.
//
// Expected:
//   - ctx must be a valid context.
//   - events may be empty; an empty slice returns empty suggestions.
//
// Returns:
//   - An InferenceResult with skill suggestions, or an error.
//
// Side effects:
//   - Sets inferenceCalled to true for test assertions.
func (s *StubSkillInferenceService) InferSkillsFromEvents(
	ctx context.Context,
	events []*career.Event,
) (*skillinference.InferenceResult, error) {
	s.inferenceCalled = true

	if len(events) == 0 {
		return &skillinference.InferenceResult{
			Suggestions:        []skillinference.SkillSuggestion{},
			ExistingSkillNames: []string{},
		}, nil
	}

	skillMap := buildSkillMap()

	var suggestions []skillinference.SkillSuggestion
	seen := make(map[string]bool)

	for _, event := range events {
		text := strings.ToLower(event.Text)
		for keyword, suggestion := range skillMap {
			if strings.Contains(text, keyword) && !seen[suggestion.Name] {
				seen[suggestion.Name] = true
				sugg := suggestion
				sugg.EventIDs = []string{event.ID}
				suggestions = append(suggestions, sugg)
			}
		}
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Name < suggestions[j].Name
	})

	return &skillinference.InferenceResult{
		Suggestions:        suggestions,
		ExistingSkillNames: []string{},
	}, nil
}

// InferSkillsFromBurst returns deterministic skill suggestions for a burst.
//
// Expected:
//   - ctx must be a valid context.
//   - burst must be a valid Burst; currently unused but required for interface.
//   - events may be empty; an empty slice returns empty suggestions.
//
// Returns:
//   - An InferenceResult with skill suggestions, or an error.
//
// Side effects:
//   - Delegates to InferSkillsFromEvents, setting inferenceCalled to true.
func (s *StubSkillInferenceService) InferSkillsFromBurst(
	ctx context.Context,
	burst *career.Burst,
	events []*career.Event,
) (*skillinference.InferenceResult, error) {
	return s.InferSkillsFromEvents(ctx, events)
}

// CreateSkillsFromSuggestions persists accepted suggestions as skills.
//
// Expected:
//   - ctx must be a valid context.
//   - suggestions may be empty; an empty slice returns no skills.
//
// Returns:
//   - A slice of created Skill pointers, or an error.
//
// Side effects:
//   - Creates skills in the skill repository.
//   - Links events to skills via the event repository.
func (s *StubSkillInferenceService) CreateSkillsFromSuggestions(
	ctx context.Context,
	suggestions []skillinference.SkillSuggestion,
) ([]*career.Skill, error) {
	var skills []*career.Skill

	for _, sugg := range suggestions {
		skill := &career.Skill{
			Name:     sugg.Name,
			Category: sugg.Category,
		}

		if err := s.skillRepo.Create(ctx, skill); err != nil {
			return nil, err
		}

		for _, eventID := range sugg.EventIDs {
			if err := s.eventRepo.LinkSkill(ctx, eventID, skill.ID); err != nil {
				return nil, err
			}
		}

		skills = append(skills, skill)
	}

	return skills, nil
}

// InferenceCalled returns whether InferSkillsFromEvents was called for assertions.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (s *StubSkillInferenceService) InferenceCalled() bool {
	return s.inferenceCalled
}
