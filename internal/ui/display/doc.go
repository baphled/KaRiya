// Package display provides presentation-only types for the TUI layer.
//
// Display types mirror domain types but contain no business logic,
// no validation, and no domain-specific behaviour. They exist to
// decouple screens and modals from the domain layer.
//
// Conversion functions (e.g. EventFromDomain) translate domain types
// into display types at the intent boundary.
package display
