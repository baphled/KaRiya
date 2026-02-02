package statematrix

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"
)

// generateMermaidDiagram creates a Mermaid state diagram for a component.
func generateMermaidDiagram(component ComponentInfo) string {
	if len(component.States) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("```mermaid\n")
	sb.WriteString("stateDiagram-v2\n")
	sb.WriteString("    direction TB\n")

	// Find ROOT state
	var rootState *StateInfo
	for i := range component.States {
		if component.States[i].Type == "ROOT" {
			rootState = &component.States[i]
			break
		}
	}

	// Start with ROOT state or first state
	if rootState != nil {
		sb.WriteString(fmt.Sprintf("    [*] --> %s\n", rootState.Constant))
	} else if len(component.States) > 0 {
		sb.WriteString(fmt.Sprintf("    [*] --> %s\n", component.States[0].Constant))
	}

	// Add transitions based on state types
	for i, state := range component.States {
		switch state.Type {
		case "ROOT":
			// ROOT can go to next state or cancel to main menu
			if i+1 < len(component.States) {
				nextState := component.States[i+1]
				sb.WriteString(fmt.Sprintf("    %s --> %s : proceed\n", state.Constant, nextState.Constant))
			}
			sb.WriteString(fmt.Sprintf("    %s --> [*] : cancel\n", state.Constant))

		case "Intermediate":
			// Intermediate can go forward or back
			if i+1 < len(component.States) {
				nextState := component.States[i+1]
				sb.WriteString(fmt.Sprintf("    %s --> %s : proceed\n", state.Constant, nextState.Constant))
			}
			if i > 0 {
				prevState := component.States[i-1]
				sb.WriteString(fmt.Sprintf("    %s --> %s : back\n", state.Constant, prevState.Constant))
			}

		case "Confirmation":
			// Confirmation can accept, reject, or cancel
			if i+1 < len(component.States) {
				nextState := component.States[i+1]
				sb.WriteString(fmt.Sprintf("    %s --> %s : confirm\n", state.Constant, nextState.Constant))
			}
			if i > 0 {
				prevState := component.States[i-1]
				sb.WriteString(fmt.Sprintf("    %s --> %s : cancel\n", state.Constant, prevState.Constant))
			}

		case "Async":
			// Async operations complete or error
			if i+1 < len(component.States) {
				nextState := component.States[i+1]
				sb.WriteString(fmt.Sprintf("    %s --> %s : complete\n", state.Constant, nextState.Constant))
			}
			// Look for error state
			for j := range component.States {
				if component.States[j].Type == "Error" {
					sb.WriteString(fmt.Sprintf("    %s --> %s : error\n", state.Constant, component.States[j].Constant))
					break
				}
			}

		case "Final", "Error":
			// Terminal states go to end
			sb.WriteString(fmt.Sprintf("    %s --> [*]\n", state.Constant))

		case "Modal":
			// Modals return to parent state
			if i > 0 {
				prevState := component.States[i-1]
				sb.WriteString(fmt.Sprintf("    %s --> %s : close\n", state.Constant, prevState.Constant))
			}
		}
	}

	sb.WriteString("```\n")
	return sb.String()
}

// GenerateMarkdown generates the STATE_MATRIX.md file.
//
// Expected:
//   - statematrix must be valid.
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func GenerateMarkdown(matrix *StateMatrix, path string) error {
	// Ensure matrix has timestamp
	if matrix.GeneratedAt.IsZero() {
		matrix.GeneratedAt = time.Now()
	}

	// Sort intents and screens
	sort.Slice(matrix.Intents, func(i, j int) bool {
		return matrix.Intents[i].Name < matrix.Intents[j].Name
	})
	sort.Slice(matrix.Screens, func(i, j int) bool {
		return matrix.Screens[i].Name < matrix.Screens[j].Name
	})

	tmpl := `# KaRiya TUI State Matrix

> **Auto-generated**: Do not edit manually. Run ` + "`make generate-state-matrix`" + ` to update.  
> **Generated**: {{.GeneratedAt.Format "2006-01-02T15:04:05Z07:00"}}  
> **Source**: ` + "`cmd/generate-state-matrix/main.go`" + `

## Summary

| Metric | Value |
|--------|-------|
| Total Intents | {{.TotalIntents}} |
| Total Screens | {{.TotalScreens}} |
| Total States | {{.TotalStates}} |

## State Type Definitions

| Type | Description | Escape Behavior |
|------|-------------|-----------------|
| **ROOT** | Entry state of intent/screen | Cancel intent → Main Menu |
| **Intermediate** | Middle navigation state | Back → Previous state |
| **Confirmation** | User confirmation required | Cancel → Parent state |
| **Async** | Background operation | Varies (some block, some allow back) |
| **Modal** | Overlay editor/form | Close modal → Parent state |
| **Final** | Operation complete | Deactivate intent |
| **Error** | Operation failed | Deactivate (retry often available) |

## Intent States

{{if .Intents}}
{{range .Intents}}
### {{.Name}}

**File**: ` + "`{{.File}}`" + `  
**States**: {{.StateCount}}

#### State Diagram

{{mermaidDiagram .}}

#### States Reference

| State | Type | Escape Behavior |
|-------|------|-----------------|
{{range .States}}| ` + "`{{.Constant}}`" + ` | {{.Type}} | {{.EscapeBehavior}} |
{{end}}

{{end}}
{{else}}
*No intent states found.*

{{end}}

## Screen States

{{if .Screens}}
{{range .Screens}}
### {{.Name}}

**File**: ` + "`{{.File}}`" + `  
**States**: {{.StateCount}}

| State | Type | Escape Behavior |
|-------|------|-----------------|
{{range .States}}| ` + "`{{.Constant}}`" + ` | {{.Type}} | {{.EscapeBehavior}} |
{{end}}

{{end}}
{{else}}
*No screen states found. Screens will appear here as intents are migrated to the new architecture.*

{{end}}

---

*Generated by ` + "`cmd/generate-state-matrix/main.go`" + `*
`

	// Create template with custom functions
	funcMap := template.FuncMap{
		"mermaidDiagram": generateMermaidDiagram,
	}

	t, err := template.New("matrix").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return err
	}

	// #nosec G304 - Path is safely constructed using filepath.Join from project root
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, matrix)
}

// GenerateJSON generates the state_matrix.json file.
//
// Expected:
//   - statematrix must be valid.
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func GenerateJSON(matrix *StateMatrix, path string) error {
	// Ensure matrix has timestamp
	if matrix.GeneratedAt.IsZero() {
		matrix.GeneratedAt = time.Now()
	}

	data, err := json.MarshalIndent(matrix, "", "  ")
	if err != nil {
		return err
	}

	// #nosec G306 - JSON file permissions are appropriate for generated documentation
	return os.WriteFile(path, data, 0o600)
}
