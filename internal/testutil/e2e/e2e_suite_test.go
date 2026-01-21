package e2e_test

import (
	"testing"

	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Test Suite")
}

var _ = BeforeSuite(func() {
	// Set up shared database once for all tests
	e2e.SetupShared()
})

var _ = AfterSuite(func() {
	// Clean up shared database
	e2e.CleanupShared()
})
