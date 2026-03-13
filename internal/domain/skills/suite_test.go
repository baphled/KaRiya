package skills

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSkillsDomain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Skills Domain Suite")
}
