package captureevent

import "github.com/baphled/kariya/internal/domain/career"

// Result is the output of a completed CaptureEvent intent.
//
// It contains the persisted event along with any enrichment data the user
// accepted or rejected during the review phase.
type Result struct {
	// Event is the captured or edited career event.
	Event *career.Event

	// Bursts are the activity bursts inferred from the event (may be empty).
	Bursts []*career.Burst

	// Facts are the career facts inferred from the event (may be empty).
	Facts []*career.Fact

	// AcceptedFields tracks which enrichment fields the user accepted.
	AcceptedFields map[string]bool

	// RejectedFields maps rejected enrichment item IDs to rejection reasons.
	RejectedFields map[string]string
}
