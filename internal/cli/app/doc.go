// Package app provides the main application orchestration and routing for the KaRiya CLI.
//
// # Overview
//
// The app package implements the top-level Bubble Tea application that coordinates
// all intents, handles global key bindings, and manages the application lifecycle.
// It serves as the entry point for the interactive TUI experience.
//
// # Architecture
//
// The App struct embeds the tea.Model interface and implements:
//   - Intent routing and navigation
//   - Global keyboard handling (quit, help, etc.)
//   - Theme management and initialization
//   - Service coordination
//
// # Usage
//
// Create and run the application:
//
//	app := app.New(cfg, db)
//	p := tea.NewProgram(app, tea.WithAltScreen())
//	_, err := p.Run()
//
// For testing and development, the app can be initialized with different
// configurations and service implementations.
package app
