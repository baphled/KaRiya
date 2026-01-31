package main

import (
	"github.com/baphled/kariya/tools/analyzers/docblocks"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(docblocks.Analyzer)
}
