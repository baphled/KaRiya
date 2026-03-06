package testutil

import (
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// NilService creates an empty Service for testing error conditions.
// This is a test helper to avoid inline struct violations in test files.
//
// Returns:
//   - An empty careerservice.Service with no configured repositories
//
// Side effects:
//   - None
func NilService() *careerservice.Service {
	return &careerservice.Service{}
}
