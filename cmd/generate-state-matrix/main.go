// Command generate-state-matrix parses intent and screen files and generates STATE_MATRIX.md
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/baphled/kariya/internal/tui/statematrix"
)

func main() {
	projectRoot := findProjectRoot()
	intentsDir := filepath.Join(projectRoot, "internal", "tui", "intents")
	screensDir := filepath.Join(projectRoot, "internal", "tui", "screens")

	fmt.Println("Generating state matrix...")
	fmt.Printf("Scanning intents: %s\n", intentsDir)
	fmt.Printf("Scanning screens: %s\n", screensDir)

	// Scan both intents and screens
	matrix, err := statematrix.ScanAll(intentsDir, screensDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning directories: %v\n", err)
		os.Exit(1)
	}

	// Set timestamp
	matrix.GeneratedAt = time.Now()

	// Generate markdown
	mdPath := filepath.Join(projectRoot, "docs", "STATE_MATRIX.md")
	if err := statematrix.GenerateMarkdown(matrix, mdPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating markdown: %v\n", err)
		os.Exit(1)
	}

	// Generate JSON
	jsonPath := filepath.Join(projectRoot, "docs", "state_matrix.json")
	if err := statematrix.GenerateJSON(matrix, jsonPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Generated state matrix:\n")
	fmt.Printf("   - %d intents\n", matrix.TotalIntents)
	fmt.Printf("   - %d screens\n", matrix.TotalScreens)
	fmt.Printf("   - %d total states\n", matrix.TotalStates)
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
