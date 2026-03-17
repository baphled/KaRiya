package noinlinecareer_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/baphled/kariya/tools/analyzers/noinlinecareer"
)

var _ = Describe("Noinlinecareer Analyzer", func() {
	var testdata string

	BeforeEach(func() {
		testdata = analysistest.TestData()
	})

	Context("when files contain inline career types", func() {
		It("reports violations", func() {
			results := analysistest.Run(GinkgoT(), testdata, noinlinecareer.Analyzer, "violations")
			Expect(results).NotTo(BeEmpty())
		})
	})

	Context("when files use proper imports", func() {
		It("passes without issues", func() {
			results := analysistest.Run(GinkgoT(), testdata, noinlinecareer.Analyzer, "clean")
			Expect(results).NotTo(BeNil())
		})
	})

	Context("when files are not test files", func() {
		It("skips analysis", func() {
			results := analysistest.Run(GinkgoT(), testdata, noinlinecareer.Analyzer, "nontest")
			Expect(results).NotTo(BeNil())
		})
	})

	Context("when the package is a fixture package", func() {
		It("excludes from analysis", func() {
			results := analysistest.Run(GinkgoT(), testdata, noinlinecareer.Analyzer, "fixtures")
			Expect(results).NotTo(BeNil())
		})
	})

	Context("when the package is domain/career", func() {
		It("excludes from analysis", func() {
			results := analysistest.Run(GinkgoT(), testdata, noinlinecareer.Analyzer, "fake/domain/career")
			Expect(results).NotTo(BeNil())
		})
	})
})
