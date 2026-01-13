// Package statematrix provides types and utilities for generating state matrix documentation.
package statematrix

import "time"

// StateInfo represents a single state in the TUI
type StateInfo struct {
	Constant       string `json:"constant"`
	Type           string `json:"type"`
	EscapeBehavior string `json:"escape_behavior"`
}

// ComponentInfo represents either an Intent or a Screen with its states
type ComponentInfo struct {
	Name       string      `json:"name"`
	File       string      `json:"file"`
	Kind       string      `json:"kind"` // "intent" or "screen"
	States     []StateInfo `json:"states"`
	StateCount int         `json:"state_count"`
}

// StateMatrix is the complete output containing both intent and screen states
type StateMatrix struct {
	GeneratedAt  time.Time       `json:"generated_at"`
	TotalIntents int             `json:"total_intents"`
	TotalScreens int             `json:"total_screens"`
	TotalStates  int             `json:"total_states"`
	Intents      []ComponentInfo `json:"intents"`
	Screens      []ComponentInfo `json:"screens"`
	LegacyStates []ComponentInfo `json:"legacy_states,omitempty"` // States pending migration
}
