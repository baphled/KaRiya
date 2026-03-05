package skills

import (
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// NilService creates an empty Service for testing error conditions.
// This is a helper function to avoid inline struct violations in test files.
func NilService() *careerservice.Service {
	return &careerservice.Service{}
}
