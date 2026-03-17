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

	rootState := findRootState(component.States)

	appendInitialTransition(&sb, rootState, component.States)

	// Add transitions based on state types
	appendStateTransitions(&sb, component.States)

	sb.WriteString("```\n")
	return sb.String()
}

// findRootState locates the first ROOT state in the list.
func findRootState(states []StateInfo) *StateInfo {
	for i := range states {
		if states[i].Type == "ROOT" {
			return &states[i]
		}
	}
	return nil
}

// appendInitialTransition writes the initial transition to the first or ROOT state.
func appendInitialTransition(sb *strings.Builder, rootState *StateInfo, states []StateInfo) {
	if rootState != nil {
		fmt.Fprintf(sb, "    [*] --> %s\n", rootState.Constant)
		return
	}
	if len(states) > 0 {
		fmt.Fprintf(sb, "    [*] --> %s\n", states[0].Constant)
	}
}

// appendStateTransitions writes state transitions for all component states.
func appendStateTransitions(sb *strings.Builder, states []StateInfo) {
	for i, state := range states {
		appendStateTransition(sb, states, i, state)
	}
}

// appendStateTransition writes transitions for a single state.
func appendStateTransition(sb *strings.Builder, states []StateInfo, index int, state StateInfo) {
	switch state.Type {
	case "ROOT":
		appendRootTransitions(sb, states, index, state)
	case "Intermediate":
		appendIntermediateTransitions(sb, states, index, state)
	case "Confirmation":
		appendConfirmationTransitions(sb, states, index, state)
	case "Async":
		appendAsyncTransitions(sb, states, index, state)
	case "Final", "Error":
		appendTerminalTransition(sb, state)
	case "Modal":
		appendModalTransitions(sb, states, index, state)
	}
}

// appendRootTransitions writes transitions for ROOT states.
func appendRootTransitions(sb *strings.Builder, states []StateInfo, index int, state StateInfo) {
	if nextState, ok := nextState(states, index); ok {
		fmt.Fprintf(sb, "    %s --> %s : proceed\n", state.Constant, nextState.Constant)
	}
	fmt.Fprintf(sb, "    %s --> [*] : cancel\n", state.Constant)
}

// appendIntermediateTransitions writes transitions for intermediate states.
func appendIntermediateTransitions(sb *strings.Builder, states []StateInfo, index int, state StateInfo) {
	if nextState, ok := nextState(states, index); ok {
		fmt.Fprintf(sb, "    %s --> %s : proceed\n", state.Constant, nextState.Constant)
	}
	if prevState, ok := previousState(states, index); ok {
		fmt.Fprintf(sb, "    %s --> %s : back\n", state.Constant, prevState.Constant)
	}
}

// appendConfirmationTransitions writes transitions for confirmation states.
func appendConfirmationTransitions(sb *strings.Builder, states []StateInfo, index int, state StateInfo) {
	if nextState, ok := nextState(states, index); ok {
		fmt.Fprintf(sb, "    %s --> %s : confirm\n", state.Constant, nextState.Constant)
	}
	if prevState, ok := previousState(states, index); ok {
		fmt.Fprintf(sb, "    %s --> %s : cancel\n", state.Constant, prevState.Constant)
	}
}

// appendAsyncTransitions writes transitions for async states.
func appendAsyncTransitions(sb *strings.Builder, states []StateInfo, index int, state StateInfo) {
	if nextState, ok := nextState(states, index); ok {
		fmt.Fprintf(sb, "    %s --> %s : complete\n", state.Constant, nextState.Constant)
	}
	if errorState, ok := findErrorState(states); ok {
		fmt.Fprintf(sb, "    %s --> %s : error\n", state.Constant, errorState.Constant)
	}
}

// appendTerminalTransition writes transitions for final/error states.
func appendTerminalTransition(sb *strings.Builder, state StateInfo) {
	fmt.Fprintf(sb, "    %s --> [*]\n", state.Constant)
}

// appendModalTransitions writes transitions for modal states.
func appendModalTransitions(sb *strings.Builder, states []StateInfo, index int, state StateInfo) {
	if prevState, ok := previousState(states, index); ok {
		fmt.Fprintf(sb, "    %s --> %s : close\n", state.Constant, prevState.Constant)
	}
}

// nextState returns the next state if it exists.
func nextState(states []StateInfo, index int) (StateInfo, bool) {
	if index+1 < len(states) {
		return states[index+1], true
	}
	return StateInfo{}, false
}

// previousState returns the previous state if it exists.
func previousState(states []StateInfo, index int) (StateInfo, bool) {
	if index > 0 {
		return states[index-1], true
	}
	return StateInfo{}, false
}

// findErrorState returns the first Error state if present.
func findErrorState(states []StateInfo) (StateInfo, bool) {
	for i := range states {
		if states[i].Type == "Error" {
			return states[i], true
		}
	}
	return StateInfo{}, false
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
