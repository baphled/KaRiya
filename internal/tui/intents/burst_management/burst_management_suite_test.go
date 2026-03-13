package burst_management_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBurstManagement(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BurstManagement Suite")
}
