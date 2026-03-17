// Package skills provides pure domain functions for skill add/edit workflows.
//
// These functions contain business logic extracted from TUI modals and UI layers to enable:
//   - Direct testing without Bubble Tea event loops
//   - BDD step definitions that bypass TUI form interactions
//   - Clear separation between domain logic and presentation
//
// All functions in this package are pure: they have no Bubble Tea, Huh, or TUI dependencies.
// They accept plain structs and domain types, and return results without side effects
// beyond modifying their input parameters.
package skills
