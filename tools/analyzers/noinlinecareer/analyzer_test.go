package noinlinecareer_test

import (
	"testing"

	"github.com/baphled/kariya/tools/analyzers/noinlinecareer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestViolationsDetected(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, noinlinecareer.Analyzer, "violations")
}

func TestCleanFilesPass(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, noinlinecareer.Analyzer, "clean")
}

func TestNonTestFilesSkipped(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, noinlinecareer.Analyzer, "nontest")
}

func TestFixturePackagesExcluded(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, noinlinecareer.Analyzer, "fixtures")
}

func TestUnexportedTypesSkipped(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, noinlinecareer.Analyzer, "career")
}
