// Package factmanagement implements the FactManagement intent for managing career facts.
//
// # Overview
//
// The factmanagement package provides the intent for managing career facts
// extracted from events, including editing, categorization, and deletion.
//
// # Architecture
//
// This package follows the standard intent subdirectory structure:
//   - context.go: IntentValidator with services
//   - intent.go: Main intent implementation
//   - constants.go: State enum
//
// # Usage
//
// Create the intent:
//
//	ctx := &factmanagement.IntentValidator{
//	    FactService: factService,
//	}
//	intent, err := factmanagement.NewIntent(ctx)
//
// For more details, see the intent documentation.
package factmanagement
