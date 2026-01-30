package capture_event

import "github.com/baphled/kariya/internal/domain/career"

// Result is the output of a successful CaptureEvent intent.
type Result struct {
	// Event is the captured or edited event.
	Event *career.Event

	// Bursts are the inferred bursts from the event (may be empty).
	Bursts []*career.Burst

	// Facts are the inferred facts from the event (may be empty).
	Facts []*career.Fact

	// AcceptedFields tracks which fields were accepted by the user.
	AcceptedFields map[string]bool

	// RejectedFields tracks which inferred fields were rejected.
	RejectedFields map[string]string
}
