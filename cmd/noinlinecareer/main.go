// Package main provides the noinlinecareer analyzer CLI entry point.
package main

import (
	"github.com/baphled/kariya/tools/analyzers/noinlinecareer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(noinlinecareer.Analyzer)
}
