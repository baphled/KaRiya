package intents_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIntentsSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Intents Suite")
}
