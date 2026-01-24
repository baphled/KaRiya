// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"context"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// IntentContext is the minimal context passed to BurstManagement intent.
// It contains only what's necessary to start the intent.
type IntentContext struct {
	// Bursts is the list of bursts to manage.
	Bursts []*career.Burst

	// Service provides burst and fact operations.
	Service *careerservice.Service

	// BurstRepository for direct burst CRUD operations.
	BurstRepository careerrepo.BurstRepository

	// Context for service calls.
	Context context.Context

	// IsNewBurst indicates if we're creating a new burst.
	IsNewBurst bool

	// EditingBurst is the burst being edited (if any).
	EditingBurst *career.Burst
}

// Validate ensures the context is complete.
func (c *IntentContext) Validate() error {
	if c.Bursts == nil {
		c.Bursts = make([]*career.Burst, 0)
	}
	if c.Context == nil {
		c.Context = context.Background()
	}
	return nil
}

// LoadBursts loads all bursts from the repository.
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

// CreateBurst creates a new burst in the repository.
func (c *IntentContext) CreateBurst(burst *career.Burst) error {
	if c.BurstRepository == nil {
		return nil
	}

	if err := burst.Validate(); err != nil {
		return err
	}

	return c.BurstRepository.Create(c.Context, burst)
}

// UpdateBurst updates an existing burst in the repository.
func (c *IntentContext) UpdateBurst(burst *career.Burst) error {
	if c.BurstRepository == nil {
		return nil
	}

	if err := burst.Validate(); err != nil {
		return err
	}

	return c.BurstRepository.Update(c.Context, burst)
}

// DeleteBurst deletes a burst from the repository.
func (c *IntentContext) DeleteBurst(burstID string) error {
	if c.BurstRepository == nil {
		return nil
	}

	return c.BurstRepository.Delete(c.Context, burstID)
}

// StartNewBurst initializes a new burst for editing.
func (c *IntentContext) StartNewBurst() {
	c.EditingBurst = &career.Burst{
		ID:          "",
		Name:        "",
		Description: "",
		EventIDs:    make([]string, 0),
	}
	c.IsNewBurst = true
}

// CancelEdit cancels the current edit without saving.
func (c *IntentContext) CancelEdit() {
	c.EditingBurst = nil
	c.IsNewBurst = false
}
