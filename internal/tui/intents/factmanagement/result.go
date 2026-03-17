// Package factmanagement implements the FactManagement intent for managing career facts.
package factmanagement

import "github.com/baphled/kariya/internal/domain/career"

// Result is the output of the FactManagement intent.
type Result struct {
	// Action describes the action performed (created, updated, deleted, none).
	Action string

	// Fact is the affected fact (may be nil if no specific fact was affected).
	Fact *career.Fact

	// Facts is the final list of facts when exiting the intent.
	Facts []*career.Fact

	// Message is a human-readable message about the operation.
	Message string

	// ScrollPosition is the final scroll position (for state restoration).
	ScrollPosition int
}
