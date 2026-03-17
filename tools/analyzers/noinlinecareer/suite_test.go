package noinlinecareer_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNoinlinecareer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Noinlinecareer Analyzer Suite")
}
