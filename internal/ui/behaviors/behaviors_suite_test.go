package behaviors_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBehaviorsSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Behaviors Suite")
}
