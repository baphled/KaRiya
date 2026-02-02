// Package burst_management implements the BurstManagement intent for managing career bursts.
//
// # Overview
//
// The burst_management package provides the intent for managing career bursts,
// including detection, review, and organization of burst periods.
//
// # Architecture
//
// This package follows the standard intent subdirectory structure:
//   - context.go: IntentContext with services
//   - intent.go: Main intent implementation
//   - types.go: Intent struct definition
//   - handlers.go: Message handlers
//   - helpers.go: Utility functions
//   - interfaces.go: Service interfaces
//
// # Usage
//
// Create the intent:
//
//	ctx := &burst_management.IntentContext{
//	    EventService: eventService,
//	}
//	intent, err := burst_management.NewIntent(ctx)
//
// For more details, see the intent documentation.
package burst_management
