// Package steps provides Godog step definitions for BDD tests.
//
// This package contains high-level, business-focused step definitions
// that map Gherkin scenarios to actions on the test environment.
// Steps are designed to be readable and reusable across feature files.
//
// Step definitions are organized by workflow:
//   - common_steps.go: Shared navigation and assertion steps
//   - onboarding_steps.go: Onboarding wizard specific steps
//
// All steps use Gomega matchers for consistent assertions with
// the existing Ginkgo-based E2E tests.
package steps
