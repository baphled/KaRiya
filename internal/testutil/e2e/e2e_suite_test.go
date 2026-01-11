// Package e2e_test contains end-to-end tests focused on DATA PERSISTENCE.
//
// Scope:
//   - Database read/write operations across intent boundaries
//   - Data integrity after complex workflows
//   - Repository integration verification
//   - Persistence across simulated app restarts
//
// Out of Scope (use internal/cli/intents/*_navigation_test.go instead):
//   - Keyboard navigation (j/k, arrow keys, Enter, Esc)
//   - View rendering and layout
//   - State transitions without persistence concerns
//   - UI component interactions
//
// Setup Functions:
//   - e2e.Setup(GinkgoT()) - Use for persistence tests (SQLite)
//   - e2e.SetupWithMemory(GinkgoT()) - Use for navigation tests (in intents package)
//
// See also: docs/TESTING_PATTERNS.md for comprehensive testing guide.
package e2e_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Test Suite")
}
