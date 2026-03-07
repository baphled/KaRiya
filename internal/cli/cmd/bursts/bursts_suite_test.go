package bursts_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBursts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bursts Suite")
}
