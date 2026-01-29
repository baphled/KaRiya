// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"context"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/burstfact"
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
