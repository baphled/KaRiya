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
