// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burst_fact"
)

// Custom message types for BurstManagement state transitions.
// ALL *Msg structs MUST be in this file.

// BurstSelectedMsg indicates the user selected a burst.
type BurstSelectedMsg struct {
	Burst *career.Burst
	Index int
}

// BurstEventsLoadedMsg is sent when events for a burst are loaded.
type BurstEventsLoadedMsg struct {
	Events []*career.CareerEvent
	Error  error
}

// BurstFactsLoadedMsg is sent when facts for a burst are loaded.
type BurstFactsLoadedMsg struct {
	Facts []*career.Fact
	Error error
}

// BurstEditCompleteMsg is sent when burst editing is complete.
type BurstEditCompleteMsg struct {
	Burst     *career.Burst
	Cancelled bool
	Error     error
}

// BurstDeletedMsg is sent when a burst is deleted.
type BurstDeletedMsg struct {
	BurstID string
	Error   error
}

// BurstConfirmedMsg is sent when a burst is confirmed.
type BurstConfirmedMsg struct {
	Burst *career.Burst
	Error error
}

// FactExtractionCompleteMsg is sent when fact extraction is complete.
type FactExtractionCompleteMsg struct {
	Facts []*career.Fact
	Error error
}

// BurstSuggestionsLoadedMsg is sent when burst detection completes.
type BurstSuggestionsLoadedMsg struct {
	Suggestions []burst_fact.BurstSuggestion
	Error       error
}

// SuggestionReviewCompleteMsg is sent when the suggestion review modal closes.
type SuggestionReviewCompleteMsg struct {
	AcceptedSuggestions []burst_fact.BurstSuggestion
	Cancelled           bool
}

// EditBurstMsg is sent when the edit modal completes successfully with updated burst data.
type EditBurstMsg struct {
	BurstID     string
	Name        string
	Description string
}
