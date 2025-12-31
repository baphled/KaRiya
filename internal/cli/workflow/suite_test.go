package workflow_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWorkflowSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Workflow Package Suite")
}

