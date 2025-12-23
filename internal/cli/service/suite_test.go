package service

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCLIEventService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CLI Event Service Suite")
}
