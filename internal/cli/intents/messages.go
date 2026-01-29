// Package intents provides shared messages for cross-intent coordination.
package intents

import "github.com/baphled/kariya/internal/domain/career"

// RequestEditEventMsg is sent to request editing an event.
// This message is handled by the app router to open the CaptureEvent intent.
type RequestEditEventMsg struct {
	Event *career.Event
}

// RequestAddEventMsg is sent to request adding a new event.
// This message is handled by the app router to open the CaptureEvent intent.
type RequestAddEventMsg struct{}
