// Package generatecv implements the GenerateCV intent for multi-step CV generation.
//
// # Overview
//
// The generatecv package provides the intent for generating CVs through a
// wizard-based workflow. Users configure profile, audience, technology focus,
// and format options via a modal wizard, then review and export the generated CV.
//
// # Architecture
//
// This package follows the standard intent subdirectory structure:
//   - context.go: IntentValidator with profiles, events, facts, and service dependencies
//   - result.go: Result struct returned on completion
//   - constants.go: State enum and export-related enums
//   - messages.go: All Msg types for async operations
//   - intent.go: Main intent lifecycle (New, Init, Update, View, Result)
//   - types.go: Intent struct definition
//   - handlers.go: All handle* methods for screen results and async messages
//   - helpers.go: Async operations, view helpers, and state management utilities
//
// # Usage
//
// Create the intent:
//
//	ctx := &generatecv.IntentValidator{
//	    AvailableProfiles: profiles,
//	    Events:            events,
//	}
//	intent, err := generatecv.NewIntent(ctx)
//
// For more details, see the intent documentation.
package generatecv
