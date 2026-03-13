// Package types provides shared CLI type definitions.
//
// # Overview
//
// The types package contains type definitions that are shared across
// multiple CLI packages. These are primarily input/output types for
// intents and screens.
//
// # Type Categories
//
//   - Capture types: Event capture input/output
//   - Export types: CV export formats and options
//   - Common types: Shared enums and structs
//
// # Usage
//
// Import specific type files as needed:
//
//	import "github.com/baphled/kariya/internal/ui/types"
//
//	func ProcessCapture(input types.CaptureInput) types.CaptureResult
package types
