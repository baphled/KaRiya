package cv_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCV(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CV Views Suite")
}
