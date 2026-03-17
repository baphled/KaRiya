// Package captureevent implements the CaptureEvent intent for capturing career events.
//
// # Overview
//
// The captureevent package provides the intent for capturing new career events
// with categorization, tagging, and fact extraction.
//
// # Architecture
//
// This package follows the standard intent subdirectory structure:
//   - context.go: IntentValidator with services
//   - intent.go: Main intent implementation
//   - types.go: Intent struct definition
//   - constants.go: State enum
//   - messages.go: tea.Msg types
//
// # Usage
//
// Create the intent:
//
//	ctx := &captureevent.IntentValidator{
//	    EventService: eventService,
//	}
//	intent, err := captureevent.NewIntent(ctx)
//
// For more details, see the intent documentation.
package captureevent
