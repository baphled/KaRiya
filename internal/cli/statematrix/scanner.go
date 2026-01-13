package statematrix

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FindIntentFiles finds all Go source files in the intents directory (excluding tests)
func FindIntentFiles(intentsDir string) ([]string, error) {
	var files []string

	err := filepath.Walk(intentsDir, func(path string, info os.FileInfo, err error) error {
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

// FindScreenFiles finds all Go source files in the screens directory (excluding tests)
func FindScreenFiles(screensDir string) ([]string, error) {
	// Check if directory exists
	if _, err := os.Stat(screensDir); os.IsNotExist(err) {
		// Return empty list, not an error - screens may not exist during migration
		return []string{}, nil
	}

	var files []string

	err := filepath.Walk(screensDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Include all .go files except tests
		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// ScanAll scans both intents and screens directories and returns a complete StateMatrix
func ScanAll(intentsDir, screensDir string) (*StateMatrix, error) {
	matrix := &StateMatrix{
		Intents: []ComponentInfo{},
		Screens: []ComponentInfo{},
	}

	// Scan intents
	intentFiles, err := FindIntentFiles(intentsDir)
	if err != nil {
		return nil, err
	}

	// Parse intent files and merge by name
	intentMap := make(map[string]*ComponentInfo)
	for _, file := range intentFiles {
		component := ParseIntentFile(file)
		if component.StateCount > 0 {
			// Merge with existing component if found
			if existing, ok := intentMap[component.Name]; ok {
				// Add states that don't already exist
				for _, state := range component.States {
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
				intentMap[component.Name] = &component
			}
		}
	}

	// Convert intent map to slice
	for _, intent := range intentMap {
		matrix.Intents = append(matrix.Intents, *intent)
		matrix.TotalStates += intent.StateCount
	}
	matrix.TotalIntents = len(matrix.Intents)

	// Scan screens
	screenFiles, err := FindScreenFiles(screensDir)
	if err != nil {
		return nil, err
	}

	// Parse screen files and merge by name
	screenMap := make(map[string]*ComponentInfo)
	for _, file := range screenFiles {
		component := ParseScreenFile(file)
		if component.StateCount > 0 {
			// Merge with existing component if found
			if existing, ok := screenMap[component.Name]; ok {
				// Add states that don't already exist
				for _, state := range component.States {
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
			} else {
				screenMap[component.Name] = &component
			}
		}
	}

	// Convert screen map to slice
	for _, screen := range screenMap {
		matrix.Screens = append(matrix.Screens, *screen)
		matrix.TotalStates += screen.StateCount
	}
	matrix.TotalScreens = len(matrix.Screens)

	// Sort components by name
	sort.Slice(matrix.Intents, func(i, j int) bool {
		return matrix.Intents[i].Name < matrix.Intents[j].Name
	})
	sort.Slice(matrix.Screens, func(i, j int) bool {
		return matrix.Screens[i].Name < matrix.Screens[j].Name
	})

	return matrix, nil
}
