// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"context"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/burst_fact"
)

// BurstService defines the interface for burst and fact operations.
// This allows for mocking in tests.
type BurstService interface {
	// Event operations.
	GetEventByID(ctx context.Context, eventID string) (*career.CareerEvent, error)
	ListEvents(ctx context.Context, filters careerrepo.ListFilters) ([]*career.CareerEvent, error)

	// Fact operations.
	GetFactsBySourceBurstID(ctx context.Context, burstID string) ([]*career.Fact, error)
	ExtractFactsFromBurst(ctx context.Context, burst *career.Burst) ([]career.Fact, error)
	SaveFact(ctx context.Context, fact *career.Fact) error

	// Suggestion operations.
	SuggestBursts(ctx context.Context, eventIDs []string) ([]burst_fact.BurstSuggestion, error)
}
