// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"context"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
)

// BurstService defines the interface for burst operations.
// This allows for mocking in tests.
type BurstService interface {
	// Burst CRUD operations.
	GetBurst(ctx context.Context, id string) (*career.Burst, error)
	ListBursts(ctx context.Context, filters careerrepo.BurstListFilters) ([]*career.Burst, error)
	CreateBurst(ctx context.Context, burst *career.Burst) error
	UpdateBurst(ctx context.Context, burst *career.Burst) error
	DeleteBurst(ctx context.Context, id string) error

	// Event operations.
	GetEventByID(ctx context.Context, eventID string) (*career.CareerEvent, error)
	ListEvents(ctx context.Context, filters *careerrepo.ListFilters) ([]*career.CareerEvent, error)

	// Fact operations.
	GetFactsBySourceBurstID(ctx context.Context, burstID string) ([]*career.Fact, error)
	ExtractFactsFromBurst(ctx context.Context, burst *career.Burst) ([]career.Fact, error)
	SaveFact(ctx context.Context, fact *career.Fact) error
}
