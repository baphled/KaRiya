package modals_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestModals(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Timeline Modals Suite")
}
