package generatecv_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGenerateCV(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GenerateCV Suite")
}
