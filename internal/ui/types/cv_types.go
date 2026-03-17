// Package types provides shared type definitions for the CLI layer.
//
//revive:disable:var-naming Package name is intentionally generic for shared CLI types.
package types

// CVProfile represents a CV generation profile with targeting information.
// This type is shared between intents and screens for CV generation workflows.
type CVProfile struct {
	ID             string
	Name           string
	TargetRole     string // principal, staff, em, senior_ic
	TargetAudience string // hiring_manager, recruiter, peer
	Description    string
}
