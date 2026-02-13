// Package features contains BDD acceptance tests and step definitions for the KaRiya application.
//
// This package implements Gherkin/Cucumber scenarios for validating application workflows
// and user interactions across all major features (capture, bursts, skills, facts, etc.).
//
// Test execution:
//   - Run: `make bdd` or `ginkgo -v ./features`
//   - Individual scenario: `ginkgo -v --focus "scenario name" ./features`
//
// Structure:
//   - *.feature: Gherkin scenario files
//   - steps/: Step implementation files
//   - support/: Test hooks and fixture setup
//   - godog_test.go: Test entry point
//
// BDD tests are designed for behaviour validation and user acceptance,
// not for coverage measurement. Coverage reports skip this package.
package features
