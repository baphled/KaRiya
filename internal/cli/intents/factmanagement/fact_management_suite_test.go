package factmanagement

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFactManagement(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fact Management Suite")
}
