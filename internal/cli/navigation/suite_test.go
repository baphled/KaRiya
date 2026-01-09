package navigation

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNavigation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Navigation Suite")
}
