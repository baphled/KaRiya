package career

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCareerModelsSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Career Models Suite")
}
