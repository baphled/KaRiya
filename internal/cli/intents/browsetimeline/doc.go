// Package browsetimeline implements the BrowseTimeline intent for browsing career events.
//
// # Overview
//
// The browsetimeline package provides the intent for browsing career events
// in chronological order with filtering and search capabilities.
//
// # Architecture
//
// This package follows the standard intent subdirectory structure:
//   - context.go: IntentContext with events and filters
//   - intent.go: Main intent implementation
//   - types.go: Intent struct definition
//   - helpers.go: Utility functions
//
// # Usage
//
// Create the intent:
//
//	ctx := &browsetimeline.IntentContext{
//	    Events: events,
//	}
//	intent, err := browsetimeline.NewIntent(ctx)
//
// For more details, see the intent documentation.
package browsetimeline
