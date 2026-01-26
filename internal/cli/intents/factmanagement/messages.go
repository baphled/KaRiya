// Package fact_management implements the FactManagement intent for managing career facts.
package factmanagement

import "github.com/baphled/kariya/internal/domain/career"

// Custom message types for FactManagement state transitions.
// ALL *Msg structs MUST be in this file.

// FactSelectedMsg indicates the user selected a fact.
type FactSelectedMsg struct {
	Fact  *career.Fact
	Index int
}

// FactSavedMsg indicates a fact was successfully saved.
type FactSavedMsg struct {
	Fact    *career.Fact
	IsNew   bool
	Message string
}

// FactDeletedMsg notifies that a fact was successfully deleted.
type FactDeletedMsg struct {
	FactID string
}

// FactsLoadedMsg indicates facts were successfully loaded.
type FactsLoadedMsg struct {
	Facts []*career.Fact
	Total int
}

// ErrorMsg indicates an error occurred during an operation.
type ErrorMsg struct {
	Code    string
	Message string
	Err     error
}
