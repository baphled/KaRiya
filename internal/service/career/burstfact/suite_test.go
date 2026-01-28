package burstfact

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBurstFactSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Burst Fact Suite")
}
