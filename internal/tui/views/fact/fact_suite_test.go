package fact_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFact(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fact Suite")
}
