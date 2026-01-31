package main

import (
	"github.com/baphled/kariya/tools/analyzers/noinlinecareer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(noinlinecareer.Analyzer)
}
