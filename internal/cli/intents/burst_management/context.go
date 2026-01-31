// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"context"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
)

// IntentContext holds the input parameters, dependencies, and mutable editing state
// needed by the BurstManagement intent throughout its lifecycle.
type IntentContext struct {
	// Bursts is the list of bursts to manage.
	Bursts []*career.Burst

	// Service provides burst and fact operations.
	Service BurstService

	// SkillInferenceService provides skill detection operations.
	SkillInferenceService SkillInferenceService

	// BurstRepository for direct burst CRUD operations.
	BurstRepository careerrepo.BurstRepository

	// SkillRepository for loading skills associated with burst events.
	SkillRepository careerrepo.SkillRepository

	// Context for service calls.
	Context context.Context

	// IsNewBurst indicates if we're creating a new burst.
	IsNewBurst bool

	// EditingBurst is the burst being edited (if any).
	EditingBurst *career.Burst
}

// Validate ensures the context has all required fields initialized with safe defaults.
//
// Returns:
//   - Always nil; missing optional fields are initialized with defaults rather than rejected.
//
// Side effects:
//   - Initializes Bursts to an empty slice if nil.
//   - Initializes Context to context.Background() if nil.
func (c *IntentContext) Validate() error {
	if c.Bursts == nil {
		c.Bursts = make([]*career.Burst, 0)
	}
	if c.Context == nil {
		c.Context = context.Background()
	}
	return nil
}

// LoadBursts fetches all bursts from the repository and replaces the in-memory burst list.
//
// Returns:
//   - An error if the repository query fails, or nil on success.
//   - Nil when BurstRepository is not configured (no-op).
//
// Side effects:
//   - Overwrites c.Bursts with the full set of bursts from the repository.
func (c *IntentContext) LoadBursts() error {
	if c.BurstRepository == nil {
		return nil
	}

	bursts, err := c.BurstRepository.List(c.Context, careerrepo.BurstListFilters{})
	if err != nil {
		return err
	}

	c.Bursts = bursts
	return nil
}

// CreateBurst validates and persists a new burst to the repository.
//
// Expected:
//   - burst must be non-nil and pass domain validation.
//
// Returns:
//   - A validation or repository error, or nil on success.
//   - Nil when BurstRepository is not configured (no-op).
//
// Side effects:
//   - Persists the burst to the underlying repository, which may assign an ID.
func (c *IntentContext) CreateBurst(burst *career.Burst) error {
	if c.BurstRepository == nil {
		return nil
	}

	if err := burst.Validate(); err != nil {
		return err
	}

	return c.BurstRepository.Create(c.Context, burst)
}

// UpdateBurst validates and persists changes to an existing burst in the repository.
//
// Expected:
//   - burst must be non-nil, have a valid ID, and pass domain validation.
//
// Returns:
//   - A validation or repository error, or nil on success.
//   - Nil when BurstRepository is not configured (no-op).
//
// Side effects:
//   - Overwrites the stored burst record with the provided values.
func (c *IntentContext) UpdateBurst(burst *career.Burst) error {
	if c.BurstRepository == nil {
		return nil
	}

	if err := burst.Validate(); err != nil {
		return err
	}

	return c.BurstRepository.Update(c.Context, burst)
}

// DeleteBurst removes a burst from the repository by its identifier.
//
// Expected:
//   - burstID must be a non-empty identifier of an existing burst.
//
// Returns:
//   - A repository error if deletion fails, or nil on success.
//   - Nil when BurstRepository is not configured (no-op).
//
// Side effects:
//   - Permanently removes the burst record from the repository.
func (c *IntentContext) DeleteBurst(burstID string) error {
	if c.BurstRepository == nil {
		return nil
	}

	return c.BurstRepository.Delete(c.Context, burstID)
}

// StartNewBurst initializes a blank burst template and marks the context for creation mode.
//
// Side effects:
//   - Assigns a new empty Burst to EditingBurst.
//   - Sets IsNewBurst to true.
func (c *IntentContext) StartNewBurst() {
	c.EditingBurst = &career.Burst{
		ID:          "",
		Name:        "",
		Description: "",
		EventIDs:    make([]string, 0),
	}
	c.IsNewBurst = true
}

// CancelEdit discards the in-progress burst edit and resets editing state.
//
// Side effects:
//   - Clears EditingBurst to nil.
//   - Sets IsNewBurst to false.
func (c *IntentContext) CancelEdit() {
	c.EditingBurst = nil
	c.IsNewBurst = false
}
