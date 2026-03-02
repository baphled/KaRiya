// Package intents (not intents_test) is used intentionally here.
//
// This suite file uses `package intents` instead of `package intents_test` because
// Ginkgo's BeforeSuite/AfterSuite hooks only apply to tests within the same package.
// The intents directory contains both internal tests (package intents) and external
// tests (package intents_test). Using `package intents` ensures that internal tests
// like configure_system_test.go get the config isolation set up by BeforeSuite.
//
// External tests (package intents_test) that use harness.Setup() get their own config
// isolation through SwapConfigPathForTesting, which preserves and restores the
// suite-level path set here.
//
// See BUG-007 documentation for details on why config isolation is required.
package intents

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/baphled/kariya/internal/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var testConfigDir string

var _ = BeforeSuite(func() {
	var err error
	testConfigDir, err = os.MkdirTemp("", "intents_test_*")
	Expect(err).NotTo(HaveOccurred())

	config.SetConfigPathForTesting(filepath.Join(testConfigDir, "config.yaml"))
})

var _ = AfterSuite(func() {
	config.ResetConfigPath()
	if testConfigDir != "" {
		_ = os.RemoveAll(testConfigDir)
	}
})

func TestIntentsSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Intents Suite")
}
