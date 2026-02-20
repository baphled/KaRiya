package contract_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRepositoryContractSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Contract Suite")
}
