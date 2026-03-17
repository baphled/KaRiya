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
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A ComponentInfo value.
//
// Side effects:
//   - None.
func ParseIntentFile(filename string) ComponentInfo {
	component, node, ok := parseComponentFile(filename, "intent", extractIntentName)
	if !ok {
		return component
	}

	component.States = collectStateConstants(node)
	component.StateCount = len(component.States)
	return component
}

// ParseScreenFile parses a screen Go file and extracts state information.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A ComponentInfo value.
//
// Side effects:
//   - None.
func ParseScreenFile(filename string) ComponentInfo {
	component, node, ok := parseComponentFile(filename, "screen", extractScreenName)
	if !ok {
		return component
	}

	component.States = collectStateConstants(node)
	component.StateCount = len(component.States)
	return component
}

// parseComponentFile parses a component file and initializes the component metadata.
func parseComponentFile(filename string, kind string, nameExtractor func(string) string) (ComponentInfo, *ast.File, bool) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return ComponentInfo{Kind: kind}, nil, false
	}

	component := ComponentInfo{
		Name:   nameExtractor(filename),
		File:   filename,
		Kind:   kind,
		States: []StateInfo{},
	}

	return component, node, true
}

// collectStateConstants extracts state constants from the parsed AST.
func collectStateConstants(node *ast.File) []StateInfo {
	states := []StateInfo{}
	ast.Inspect(node, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			appendStatesFromGenDecl(&states, genDecl)
		}
		return true
	})
	return states
}

// appendStatesFromGenDecl appends state constants from a constant declaration.
func appendStatesFromGenDecl(states *[]StateInfo, genDecl *ast.GenDecl) {
	for _, spec := range genDecl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		appendStatesFromValueSpec(states, valueSpec)
	}
}

// appendStatesFromValueSpec appends state constants from a value spec.
func appendStatesFromValueSpec(states *[]StateInfo, valueSpec *ast.ValueSpec) {
	for _, name := range valueSpec.Names {
		if isStateConstant(name.Name) {
			stateType := ClassifyState(name.Name)
			*states = append(*states, buildStateInfo(name.Name, stateType))
		}
	}
}

// buildStateInfo builds a StateInfo from a constant name and type.
func buildStateInfo(name, stateType string) StateInfo {
	return StateInfo{
		Constant:       name,
		Type:           stateType,
		EscapeBehavior: InferEscapeBehavior(stateType),
	}
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

	if isStandardSubdirectoryFile(base) {
		if resolved, ok := resolveSubdirectoryName(filename); ok {
			return resolved
		}
	} else {
		base = strings.TrimSuffix(base, "_intent.go")
		base = strings.TrimSuffix(base, ".go")
	}

	return snakeToTitle(base)
}

func isStandardSubdirectoryFile(base string) bool {
	standardFiles := map[string]bool{
		"intent.go": true, "constants.go": true, "context.go": true,
		"result.go": true, "messages.go": true, "types.go": true,
		"handlers.go": true, "helpers.go": true, "interfaces.go": true,
		"filters.go": true,
	}
	return standardFiles[base]
}

func resolveSubdirectoryName(filename string) (string, bool) {
	dir := filepath.Base(filepath.Dir(filename))
	if dir == "." || dir == "intents" {
		return "", false
	}
	if displayName, ok := packageDisplayNames[dir]; ok {
		return displayName, true
	}
	return snakeToTitle(dir), true
}

func snakeToTitle(s string) string {
	caser := cases.Title(language.English)
	parts := strings.Split(s, "_")
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
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
//
// containsAny returns true if any substring in subs is found in s.
func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// ClassifyState categorizes a state name into a classification type.
//
// Expected:
//   - name: state name to classify
//
// Returns:
//   - string: classification type (Final, Initial, Intermediate, etc.)
//
// Side effects:
//   - None.
func ClassifyState(name string) string {
	nameLower := strings.ToLower(name)

	// Priority order: check most specific patterns first

	// Final states
	if containsAny(nameLower, "complete", "done") {
		return "Final"
	}

	// Error states
	if containsAny(nameLower, "failed", "error") {
		return "Error"
	}

	// Async states (operations in progress)
	if containsAny(nameLower, "progress", "generating", "exporting", "saving", "extracting", "inprogress") {
		return "Async"
	}

	// Confirmation states
	if containsAny(nameLower, "confirm", "deleteconfirm") {
		return "Confirmation"
	}

	// Modal states
	if containsAny(nameLower, "editmode", "modal") {
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
	if containsAny(nameLower, "choose", "selectop", "selectdomain", "selecttype", "selectprofile", "fileselect") {
		return "ROOT"
	}

	// Everything else is Intermediate
	return "Intermediate"
}

// InferEscapeBehavior infers the escape key behavior based on state type.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
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
