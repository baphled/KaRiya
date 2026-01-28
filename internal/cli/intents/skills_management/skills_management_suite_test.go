package skills_management_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestManageSkills(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Manage Skills Suite")
}
