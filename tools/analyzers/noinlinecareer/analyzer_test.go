package noinlinecareer_test

import (
	. "github.com/onsi/ginkgo/v2"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/baphled/kariya/tools/analyzers/noinlinecareer"
)

var _ = Describe("Noinlinecareer Analyzer", func() {
	var testdata string

	BeforeEach(func() {
		testdata = analysistest.TestData()
	})

	It("detects violations", func() {
		analysistest.Run(GinkgoT(), testdata, noinlinecareer.Analyzer, "violations")
	})

	It("passes clean files", func() {
		analysistest.Run(GinkgoT(), testdata, noinlinecareer.Analyzer, "clean")
	})

	It("skips non-test files", func() {
		analysistest.Run(GinkgoT(), testdata, noinlinecareer.Analyzer, "nontest")
	})

	It("excludes fixture packages", func() {
		analysistest.Run(GinkgoT(), testdata, noinlinecareer.Analyzer, "fixtures")
	})

	It("excludes domain career package", func() {
		analysistest.Run(GinkgoT(), testdata, noinlinecareer.Analyzer, "fake/domain/career")
	})
})
