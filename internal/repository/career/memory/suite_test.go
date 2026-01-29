package memory

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMemoryRepositorySuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Memory Repository Suite")
}
