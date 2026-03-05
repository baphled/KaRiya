package facts_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFacts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Facts Suite")
}
