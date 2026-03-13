package browsetimeline

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBrowseTimeline(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Browse Timeline Suite")
}
