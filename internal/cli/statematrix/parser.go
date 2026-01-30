package statematrix

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ParseIntentFile parses an intent Go file and extracts state information.
func ParseIntentFile(filename string) ComponentInfo {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		// Return empty component on parse error
		return ComponentInfo{Kind: "intent"}
	}

	component := ComponentInfo{
		Name:   extractIntentName(filename),
		File:   filename,
		Kind:   "intent",
		States: []StateInfo{},
	}

	// Find state constants
	ast.Inspect(node, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valueSpec.Names {
						if isStateConstant(name.Name) {
							stateType := ClassifyState(name.Name)
							state := StateInfo{
								Constant:       name.Name,
								Type:           stateType,
								EscapeBehavior: InferEscapeBehavior(stateType),
							}
							component.States = append(component.States, state)
						}
					}
				}
			}
		}
		return true
	})

	component.StateCount = len(component.States)
	return component
}

// ParseScreenFile parses a screen Go file and extracts state information.
func ParseScreenFile(filename string) ComponentInfo {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		// Return empty component on parse error
		return ComponentInfo{Kind: "screen"}
	}

	component := ComponentInfo{
		Name:   extractScreenName(filename),
		File:   filename,
		Kind:   "screen",
		States: []StateInfo{},
	}

	// Find state constants
	ast.Inspect(node, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valueSpec.Names {
						if isStateConstant(name.Name) {
							stateType := ClassifyState(name.Name)
							state := StateInfo{
								Constant:       name.Name,
								Type:           stateType,
								EscapeBehavior: InferEscapeBehavior(stateType),
							}
							component.States = append(component.States, state)
						}
					}
				}
			}
		}
		return true
	})

	component.StateCount = len(component.States)
	return component
}

// packageDisplayNames maps concatenated Go package names (no underscores)
// to their proper display names. Go package naming convention (ST1003)
// requires lowercase without underscores, but display names need
// word boundaries for readability in generated documentation.
var packageDisplayNames = map[string]string{
	"browsetimeline":   "BrowseTimeline",
	"captureevent":     "CaptureEvent",
	"factmanagement":   "FactManagement",
	"skillsmanagement": "SkillsManagement",
}

// extractIntentName derives the intent name from the filename.
//
// Supports both flat structure (browse_timeline.go -> BrowseTimeline)
// and subdirectory structure (captureevent/constants.go -> CaptureEvent).
// For subdirectory files (constants.go, intent.go, etc.), the parent
// directory name is used instead of the filename.
func extractIntentName(filename string) string {
	base := filepath.Base(filename)

	// For subdirectory structure, use the parent directory name when the
	// filename is a standard subdirectory file (intent.go, constants.go, etc.).
	standardFiles := map[string]bool{
		"intent.go": true, "constants.go": true, "context.go": true,
		"result.go": true, "messages.go": true, "types.go": true,
		"handlers.go": true, "helpers.go": true, "interfaces.go": true,
		"filters.go": true,
	}
	if standardFiles[base] {
		dir := filepath.Base(filepath.Dir(filename))
		if dir != "." && dir != "intents" {
			// Check known package display names for concatenated names.
			if displayName, ok := packageDisplayNames[dir]; ok {
				return displayName
			}
			base = dir
		}
	} else {
		// Remove _intent.go or .go suffix for flat files.
		base = strings.TrimSuffix(base, "_intent.go")
		base = strings.TrimSuffix(base, ".go")
	}

	// Convert snake_case to TitleCase.
	caser := cases.Title(language.English)
	parts := strings.Split(base, "_")
	for i := range parts {
		parts[i] = caser.String(parts[i])
	}
	return strings.Join(parts, "")
}

// extractScreenName derives the screen name from the filename.
func extractScreenName(filename string) string {
	base := filepath.Base(filename)
	// Remove _screen.go or .go suffix
	name := strings.TrimSuffix(base, "_screen.go")
	name = strings.TrimSuffix(name, ".go")
	// Convert snake_case to TitleCase
	caser := cases.Title(language.English)
	parts := strings.Split(name, "_")
	for i := range parts {
		parts[i] = caser.String(parts[i])
	}
	result := strings.Join(parts, "")

	// Add "Screen" suffix if not already present
	if !strings.HasSuffix(result, "Screen") {
		result += "Screen"
	}

	return result
}

// isStateConstant checks if a constant name represents a state.
func isStateConstant(name string) bool {
	return strings.Contains(name, "State") &&
		!strings.HasSuffix(name, "Model") &&
		!strings.HasSuffix(name, "Context") &&
		!strings.HasSuffix(name, "Result")
}

// ClassifyState determines the type of state based on naming patterns.
func ClassifyState(name string) string {
	nameLower := strings.ToLower(name)

	// Priority order: check most specific patterns first

	// Final states
	if strings.Contains(nameLower, "complete") || strings.Contains(nameLower, "done") {
		return "Final"
	}

	// Error states
	if strings.Contains(nameLower, "failed") || strings.Contains(nameLower, "error") {
		return "Error"
	}

	// Async states (operations in progress)
	if strings.Contains(nameLower, "progress") || strings.Contains(nameLower, "generating") ||
		strings.Contains(nameLower, "exporting") || strings.Contains(nameLower, "saving") ||
		strings.Contains(nameLower, "extracting") || strings.Contains(nameLower, "inprogress") {
		return "Async"
	}

	// Confirmation states
	if strings.Contains(nameLower, "confirm") || strings.Contains(nameLower, "deleteconfirm") {
		return "Confirmation"
	}

	// Modal states
	if strings.Contains(nameLower, "editmode") || strings.Contains(nameLower, "modal") {
		return "Modal"
	}

	// ROOT states - only first state of workflow
	// Timeline is special - it's the first state
	if strings.Contains(nameLower, "timeline") && !strings.Contains(nameLower, "detail") {
		return "ROOT"
	}
	// List states are typically ROOT (except detail lists)
	if strings.Contains(nameLower, "list") && !strings.Contains(nameLower, "detail") {
		return "ROOT"
	}
	// Choose/Select at beginning of workflow
	if strings.Contains(nameLower, "choose") || strings.Contains(nameLower, "selectop") ||
		strings.Contains(nameLower, "selectdomain") || strings.Contains(nameLower, "selecttype") ||
		strings.Contains(nameLower, "selectprofile") || strings.Contains(nameLower, "fileselect") {
		return "ROOT"
	}

	// Everything else is Intermediate
	return "Intermediate"
}

// InferEscapeBehavior infers the escape key behavior based on state type.
func InferEscapeBehavior(stateType string) string {
	switch stateType {
	case "ROOT":
		return "Cancel intent → Main Menu"
	case "Intermediate":
		return "Back → Previous state"
	case "Confirmation":
		return "Cancel → Parent state"
	case "Async":
		return "Varies (let complete or allow back)"
	case "Modal":
		return "Close modal → Parent state"
	case "Final":
		return "Deactivate intent"
	case "Error":
		return "Deactivate intent (retry available)"
	default:
		return "Back → Previous state"
	}
}
