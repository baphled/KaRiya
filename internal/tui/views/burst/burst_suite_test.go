package burst_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBurst(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Burst Views Suite")
}
