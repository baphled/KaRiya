package technology_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTechnology(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Technology Suite")
}
