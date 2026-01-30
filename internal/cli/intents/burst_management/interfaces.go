// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"context"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
)

// BurstService defines the interface for burst and fact operations.
// This allows for mocking in tests.
//
//nolint:interfacebloat // Service interface groups cohesive burst operations (events, facts, suggestions).
type BurstService interface {
	// Event operations.
	GetEventByID(ctx context.Context, eventID string) (*career.Event, error)
	ListEvents(ctx context.Context, filters careerrepo.EventListFilters) ([]*career.Event, error)

	// Burst confirmation.
	ConfirmBurst(ctx context.Context, burst *career.Burst) error

	// Fact operations.
	GetFactsBySourceBurstID(ctx context.Context, burstID string) ([]*career.Fact, error)
	ExtractFactsFromBurst(ctx context.Context, burst *career.Burst) ([]career.Fact, error)
	SaveFact(ctx context.Context, fact *career.Fact) error

	// Suggestion operations.
	SuggestBursts(ctx context.Context, eventIDs []string) ([]burstfact.BurstSuggestion, error)
}

// SkillInferenceService defines the interface for skill inference operations.
type SkillInferenceService interface {
	// InferSkillsFromEvents analyzes event descriptions and infers skills with confidence scores.
	InferSkillsFromEvents(ctx context.Context, events []*career.Event) ([]skillinference.SkillSuggestion, error)

	// CreateSkillsFromSuggestions persists skill suggestions as confirmed skills.
	CreateSkillsFromSuggestions(ctx context.Context, suggestions []skillinference.SkillSuggestion) ([]*career.Skill, error)
}
