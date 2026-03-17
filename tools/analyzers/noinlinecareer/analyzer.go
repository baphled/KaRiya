// Package noinlinecareer provides a static analyzer that forbids inline
// career struct literals in test files, enforcing fixture factory usage.
package noinlinecareer

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer forbids inline career.* struct literals in test files.
// Test files must use the fixtures package instead of constructing
// career domain types inline.
var Analyzer = &analysis.Analyzer{
	Name: "noinlinecareer",
	Doc:  "forbids inline career.* struct literals in test files; use fixtures package instead",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if isExcludedPackage(pass.Pkg.Path()) {
		return nil, nil //nolint:nilnil // go/analysis framework requires (interface{}, error) return
	}

	for _, file := range pass.Files {
		if !isTestFile(pass, file) {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			lit, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}

			named := resolveCareerType(pass, lit)
			if named == nil {
				return true
			}

			typeName := named.Obj().Name()
			pass.Reportf(lit.Pos(), "inline career.%s{} in test file; use fixtures package instead", typeName)

			return true
		})
	}

	return nil, nil //nolint:nilnil // go/analysis framework requires (interface{}, error) return
}

func resolveCareerType(pass *analysis.Pass, lit *ast.CompositeLit) *types.Named {
	named := extractNamedType(pass.TypesInfo.TypeOf(lit))
	if named == nil {
		return nil
	}

	pkg := named.Obj().Pkg()
	if pkg == nil || !named.Obj().Exported() {
		return nil
	}

	if !isCareerPackage(pkg.Path()) {
		return nil
	}

	return named
}

func extractNamedType(typ types.Type) *types.Named {
	if typ == nil {
		return nil
	}

	switch t := typ.(type) {
	case *types.Named:
		return t
	case *types.Pointer:
		return extractNamedType(t.Elem())
	default:
		return nil
	}
}

func isCareerPackage(pkgPath string) bool {
	return strings.Contains(pkgPath, "domain/career") ||
		strings.Contains(pkgPath, "model/career")
}

func isExcludedPackage(pkgPath string) bool {
	exclusions := []string{
		"domain/career",
		"model/career",
		"repository/career",
		"fixtures",
		"testutil/harness",
	}

	for _, exclusion := range exclusions {
		if strings.Contains(pkgPath, exclusion) {
			return true
		}
	}

	return false
}

func isTestFile(pass *analysis.Pass, file *ast.File) bool {
	filename := pass.Fset.Position(file.Pos()).Filename
	return strings.HasSuffix(filename, "_test.go")
}
