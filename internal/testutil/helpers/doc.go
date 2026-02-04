// Package helpers provides reusable test utilities for KaRiya tests.
//
// # Overview
//
// The helpers package extracts common test operations from the E2E test
// helpers into a reusable package. This avoids duplication across test suites.
//
// # Categories
//
//   - setup.go: Test environment setup and cleanup
//   - navigation.go: TUI navigation helpers (keys, form interaction)
//   - assertions.go: View and data verification
//   - data.go: Test data population
//   - commands.go: Bubble Tea command execution
//
// # Usage
//
// Most helpers are methods on TestEnv for a fluent API:
//
//	env := helpers.Setup(t)
//	defer env.Cleanup()
//	env.NavigateDown().Confirm()
//	env.AssertEventCount(1)
package helpers
