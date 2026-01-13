// Command generate-state-matrix parses intent files and generates STATE_MATRIX.md
package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// StateInfo represents a single state
type StateInfo struct {
	Constant       string `json:"constant"`
	Type           string `json:"type"`
	EscapeBehavior string `json:"escape_behavior"`
}

// IntentInfo represents an intent with its states
type IntentInfo struct {
	Name       string      `json:"name"`
	File       string      `json:"file"`
	States     []StateInfo `json:"states"`
	StateCount int         `json:"state_count"`
}

// StateMatrix is the complete output
type StateMatrix struct {
	GeneratedAt  time.Time    `json:"generated_at"`
	TotalIntents int          `json:"total_intents"`
	TotalStates  int          `json:"total_states"`
	Intents      []IntentInfo `json:"intents"`
}

func main() {
	projectRoot := findProjectRoot()
	intentsDir := filepath.Join(projectRoot, "internal", "cli", "intents")

	fmt.Println("Generating state matrix...")
	fmt.Printf("Scanning: %s\n", intentsDir)

	matrix := &StateMatrix{
		GeneratedAt: time.Now(),
		Intents:     []IntentInfo{},
	}

	// Find all intent files
	files, err := findIntentFiles(intentsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding intent files: %v\n", err)
		os.Exit(1)
	}

	// Parse each intent file and merge states by intent name
	intentMap := make(map[string]*IntentInfo)
	for _, file := range files {
		intent := parseIntentFile(file)
		if intent.StateCount > 0 {
			// Merge with existing intent if found
			if existing, ok := intentMap[intent.Name]; ok {
				// Add states that don't already exist
				for _, state := range intent.States {
					found := false
					for _, existingState := range existing.States {
						if existingState.Constant == state.Constant {
							found = true
							break
						}
					}
					if !found {
						existing.States = append(existing.States, state)
						existing.StateCount++
					}
				}
				// Update file path to prefer _intent.go files
				if strings.HasSuffix(file, "_intent.go") {
					existing.File = file
				}
			} else {
				intentMap[intent.Name] = &intent
			}
		}
	}

	// Convert map to slice
	for _, intent := range intentMap {
		matrix.Intents = append(matrix.Intents, *intent)
		matrix.TotalStates += intent.StateCount
	}

	matrix.TotalIntents = len(matrix.Intents)

	// Sort by name
	sort.Slice(matrix.Intents, func(i, j int) bool {
		return matrix.Intents[i].Name < matrix.Intents[j].Name
	})

	// Generate markdown
	mdPath := filepath.Join(projectRoot, "docs", "STATE_MATRIX.md")
	if err := generateMarkdown(matrix, mdPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating markdown: %v\n", err)
		os.Exit(1)
	}

	// Generate JSON
	jsonPath := filepath.Join(projectRoot, "docs", "state_matrix.json")
	if err := generateJSON(matrix, jsonPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Generated state matrix:\n")
	fmt.Printf("   - %d intents\n", matrix.TotalIntents)
	fmt.Printf("   - %d states\n", matrix.TotalStates)
	fmt.Printf("   - %s\n", mdPath)
	fmt.Printf("   - %s\n", jsonPath)
}

func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		dir = "."
	}
	// Walk up to find go.mod
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "." // Fallback
		}
		dir = parent
	}
}

func findIntentFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Include both *_intent.go files and model files (*.go except *_test.go)
		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func parseIntentFile(filename string) IntentInfo {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", filename, err)
		return IntentInfo{}
	}

	intent := IntentInfo{
		Name:   extractIntentName(filename),
		File:   filename,
		States: []StateInfo{},
	}

	// Find state constants
	ast.Inspect(node, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valueSpec.Names {
						if isStateConstant(name.Name) {
							state := StateInfo{
								Constant:       name.Name,
								Type:           classifyState(name.Name),
								EscapeBehavior: inferEscapeBehavior(name.Name),
							}
							intent.States = append(intent.States, state)
						}
					}
				}
			}
		}
		return true
	})

	intent.StateCount = len(intent.States)
	return intent
}

func extractIntentName(filename string) string {
	base := filepath.Base(filename)
	// Remove _intent.go or .go suffix
	name := strings.TrimSuffix(base, "_intent.go")
	name = strings.TrimSuffix(name, ".go")
	// Convert snake_case to TitleCase
	caser := cases.Title(language.English)
	parts := strings.Split(name, "_")
	for i := range parts {
		parts[i] = caser.String(parts[i])
	}
	return strings.Join(parts, "")
}

func isStateConstant(name string) bool {
	return strings.Contains(name, "State") &&
		!strings.HasSuffix(name, "Model") &&
		!strings.HasSuffix(name, "Context") &&
		!strings.HasSuffix(name, "Result")
}

func classifyState(name string) string {
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

func inferEscapeBehavior(name string) string {
	stateType := classifyState(name)
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

func generateMarkdown(matrix *StateMatrix, path string) error {
	tmpl := `# KaRiya TUI State Matrix

> **Auto-generated**: Do not edit manually. Run ` + "`make generate-state-matrix`" + ` to update.  
> **Generated**: {{.GeneratedAt.Format "2006-01-02T15:04:05Z07:00"}}  
> **Source**: ` + "`cmd/generate-state-matrix/main.go`" + `

## Summary

| Metric | Value |
|--------|-------|
| Total Intents | {{.TotalIntents}} |
| Total States | {{.TotalStates}} |

## State Type Definitions

| Type | Description | Escape Behavior |
|------|-------------|-----------------|
| **ROOT** | Entry state of intent | Cancel intent → Main Menu |
| **Intermediate** | Middle navigation state | Back → Previous state |
| **Confirmation** | User confirmation required | Cancel → Parent state |
| **Async** | Background operation | Varies (some block, some allow back) |
| **Modal** | Overlay editor/form | Close modal → Parent state |
| **Final** | Operation complete | Deactivate intent |
| **Error** | Operation failed | Deactivate (retry often available) |

## Complete State Matrix

{{range .Intents}}
### {{.Name}}

**File**: ` + "`{{.File}}`" + `  
**States**: {{.StateCount}}

| State | Type | Escape Behavior |
|-------|------|-----------------|
{{range .States}}| ` + "`{{.Constant}}`" + ` | {{.Type}} | {{.EscapeBehavior}} |
{{end}}
{{end}}

---

*Generated by ` + "`cmd/generate-state-matrix/main.go`" + `*
`

	t, err := template.New("matrix").Parse(tmpl)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, matrix)
}

func generateJSON(matrix *StateMatrix, path string) error {
	data, err := json.MarshalIndent(matrix, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
