package classification

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestClassificationSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Classification Suite")
}
