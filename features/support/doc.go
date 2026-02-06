// Package support provides infrastructure for Godog BDD tests.
//
// This package bridges the existing TestEnv infrastructure from
// internal/testutil/e2e with Godog's context-based step definitions.
// It provides:
//
//   - Context helpers for passing TestEnv between steps
//   - Scenario lifecycle hooks (BeforeScenario/AfterScenario)
//   - Gomega matcher integration for consistent assertions
//
// The package is designed to fully reuse the existing E2E test infrastructure
// without duplication, ensuring BDD scenarios benefit from the same robust
// test environment used by Ginkgo-based E2E tests.
package support
