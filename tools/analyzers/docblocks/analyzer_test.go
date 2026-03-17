package docblocks_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/baphled/kariya/tools/analyzers/docblocks"
)

var _ = Describe("Docblocks Analyzer", func() {
	var testdata string

	BeforeEach(func() {
		testdata = analysistest.TestData()
	})

	Describe("analysing exported functions", func() {
		It("reports missing or malformed docblocks", func() {
			results := analysistest.Run(GinkgoT(), testdata, docblocks.Analyzer, "funcs")
			Expect(results).NotTo(BeEmpty())
		})
	})

	Describe("analysing exported methods", func() {
		It("reports missing or malformed docblocks", func() {
			results := analysistest.Run(GinkgoT(), testdata, docblocks.Analyzer, "methods")
			Expect(results).NotTo(BeEmpty())
		})
	})

	Describe("analysing exported types", func() {
		It("reports missing or malformed docblocks", func() {
			results := analysistest.Run(GinkgoT(), testdata, docblocks.Analyzer, "types")
			Expect(results).NotTo(BeEmpty())
		})
	})

	Describe("analysing exported constants and variables", func() {
		It("reports missing or malformed docblocks", func() {
			results := analysistest.Run(GinkgoT(), testdata, docblocks.Analyzer, "constvars")
			Expect(results).NotTo(BeEmpty())
		})
	})

	Context("when the package is excluded", func() {
		It("skips main package files", func() {
			results := analysistest.Run(GinkgoT(), testdata, docblocks.Analyzer, "mainpkg")
			Expect(results).NotTo(BeNil())
		})
	})
})
