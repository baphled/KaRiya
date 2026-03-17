package docblocks_test

import (
	. "github.com/onsi/ginkgo/v2"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/baphled/kariya/tools/analyzers/docblocks"
)

var _ = Describe("Docblocks Analyzer", func() {
	var testdata string

	BeforeEach(func() {
		testdata = analysistest.TestData()
	})

	It("validates function docblocks", func() {
		analysistest.Run(GinkgoT(), testdata, docblocks.Analyzer, "funcs")
	})

	It("validates method docblocks", func() {
		analysistest.Run(GinkgoT(), testdata, docblocks.Analyzer, "methods")
	})

	It("validates type docblocks", func() {
		analysistest.Run(GinkgoT(), testdata, docblocks.Analyzer, "types")
	})

	It("validates const and var docblocks", func() {
		analysistest.Run(GinkgoT(), testdata, docblocks.Analyzer, "constvars")
	})

	It("handles exclusions in main package", func() {
		analysistest.Run(GinkgoT(), testdata, docblocks.Analyzer, "mainpkg")
	})
})
