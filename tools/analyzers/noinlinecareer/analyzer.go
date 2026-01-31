package noinlinecareer

import (
	"fmt"
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
	typ := pass.TypesInfo.TypeOf(lit)
	if typ == nil {
		return nil
	}

	named := extractNamed(typ)
	if named == nil {
		return nil
	}

	pkg := named.Obj().Pkg()
	if pkg == nil {
		return nil
	}

	if !isCareerPackage(pkg.Path()) {
		return nil
	}

	if !named.Obj().Exported() {
		return nil
	}

	return named
}

func extractNamed(typ types.Type) *types.Named {
	switch t := typ.(type) {
	case *types.Named:
		return t
	case *types.Pointer:
		return extractNamed(t.Elem())
	default:
		return nil
	}
}

func isCareerPackage(pkgPath string) bool {
	return pkgPath == "career" ||
		strings.HasSuffix(pkgPath, "/career") ||
		strings.Contains(pkgPath, "domain/career")
}

func isExcludedPackage(pkgPath string) bool {
	exclusions := []string{
		"fixtures",
		"testutil/e2e",
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

// FormatDiagnostic creates a consistent diagnostic message for a given career type.
//
// Expected:
//   - typeName is the name of the career struct type (e.g. "Event").
//
// Returns:
//   - A formatted string describing the violation.
//
// Side effects:
//   - None.
func FormatDiagnostic(typeName string) string {
	return fmt.Sprintf("inline career.%s{} in test file; use fixtures package instead", typeName)
}
