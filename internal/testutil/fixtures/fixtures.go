// Package fixtures provides test data factories for KaRiya domain objects.
// It uses factory-go for the factory pattern and gofakeit for realistic fake data.
//
// Usage:
//
//	// Simple creation with defaults
//	event := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
//
//	// With overrides
//	event := fixtures.EventFactory.MustCreateWithOption(map[string]interface{}{
//	    "Company": "TechCorp",
//	    "Project": "Platform",
//	}).(*career.CareerEvent)
//
//	// Quick helpers for minimal valid objects
//	event := fixtures.Event("my-id")
//	burst := fixtures.Burst("burst-id", "evt-1", "evt-2")
//	fact := fixtures.Fact("fact-id", "evt-1")
package fixtures

import (
	"github.com/brianvoe/gofakeit/v7"
)

func init() {
	// Seed gofakeit with 0 for consistent default behavior across test runs.
	// Use SetSeed() to set a specific seed for deterministic test data.
	if err := gofakeit.Seed(0); err != nil {
		// Seeding should not fail in normal circumstances
		_ = err
	}
}

// SetSeed sets the random seed for reproducible test data.
// Use this at the start of tests that need deterministic data.
func SetSeed(seed int64) {
	if err := gofakeit.Seed(seed); err != nil {
		// Seeding should not fail in normal circumstances
		_ = err
	}
}
