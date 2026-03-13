// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"github.com/baphled/kariya/internal/domain/career"
)

// Result is the result returned when BurstManagement intent completes.
type Result struct {
	// Action performed: "selected", "created", "updated", "deleted", "confirmed", "cancelled"
	Action string

	// The burst that was affected (if applicable).
	Burst *career.Burst

	// All bursts (for list view).
	Bursts []*career.Burst

	// ViewedBursts tracks bursts viewed during the session.
	ViewedBursts []*career.Burst

	// SelectedIndex is the index of the selected burst.
	SelectedIndex int
}
