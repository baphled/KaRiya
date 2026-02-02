// Package intents provides workflow orchestration and state management for the TUI.
//
// # Overview
//
// The intents package implements the core workflow engine using the Bubble Tea
// architecture. Each intent represents a complete user workflow with multiple
// states, screens, and modals.
//
// # Architecture
//
// Intents follow a strict subdirectory structure:
//
//	intents/
//	├── context.go          # Input parameters and validation
//	├── result.go           # Output types
//	├── constants.go        # State enums
//	├── messages.go         # tea.Msg types
//	├── intent.go           # Main implementation
//	├── types.go            # Intent struct (optional)
//	├── handlers.go         # Message handlers (optional)
//	└── helpers.go          # Utilities (optional)
//
// # Intent Requirements
//
// All intents must:
//   - Embed *BaseIntent
//   - Define a typed state enum
//   - Have explicit screen fields
//   - Use ScreenResult for communication
//   - Follow the subdirectory structure
//
// For detailed documentation, see docs/INTENT_ARCHITECTURE_GUIDE.md.
package intents
