package docblocks_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDocblocks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Docblocks Analyzer Suite")
}
