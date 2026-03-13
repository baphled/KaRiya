package statematrix

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FindIntentFiles finds all Go source files in the intents directory (excluding tests).
//
// Expected:
//   - intentsdir must be a valid directory path.
//
// Returns:
//   - A []string value containing file paths.
//   - An error value if walking the directory failed.
//
// Side effects:
//   - None.
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

// FindScreenFiles finds all Go source files in the screens directory (excluding tests).
//
// Expected:
//   - screensdir must be a valid directory path.
//
// Returns:
//   - A []string value containing file paths (empty if directory doesn't exist).
//   - An error value if walking the directory failed.
//
// Side effects:
//   - None.
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

// ScanAll scans both intents and screens directories and returns a complete StateMatrix.
//
// Expected:
//   - intentsdir must be a valid directory path.
//   - screensdir must be a valid directory path.
//
// Returns:
//   - A fully initialized StateMatrix ready for use.
//   - An error value if scanning failed.
//
// Side effects:
//   - None.
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
	intentMap := mergeIntentFiles(intentFiles)
	appendIntents(matrix, intentMap)

	// Scan screens
	screenFiles, err := FindScreenFiles(screensDir)
	if err != nil {
		return nil, err
	}

	// Parse screen files and merge by name
	screenMap := mergeScreenFiles(screenFiles)
	appendScreens(matrix, screenMap)

	// Sort components by name
	sort.Slice(matrix.Intents, func(i, j int) bool {
		return matrix.Intents[i].Name < matrix.Intents[j].Name
	})
	sort.Slice(matrix.Screens, func(i, j int) bool {
		return matrix.Screens[i].Name < matrix.Screens[j].Name
	})

	return matrix, nil
}

// mergeIntentFiles parses intent files and merges state data by component name.
func mergeIntentFiles(intentFiles []string) map[string]*ComponentInfo {
	intentMap := make(map[string]*ComponentInfo)
	for _, file := range intentFiles {
		component := ParseIntentFile(file)
		if component.StateCount == 0 {
			continue
		}
		mergeComponent(intentMap, component, prefersIntentFile(file))
	}
	return intentMap
}

// mergeScreenFiles parses screen files and merges state data by component name.
func mergeScreenFiles(screenFiles []string) map[string]*ComponentInfo {
	screenMap := make(map[string]*ComponentInfo)
	for _, file := range screenFiles {
		component := ParseScreenFile(file)
		if component.StateCount == 0 {
			continue
		}
		mergeComponent(screenMap, component, false)
	}
	return screenMap
}

// mergeComponent merges component state data into the component map.
func mergeComponent(componentMap map[string]*ComponentInfo, component ComponentInfo, preferFile bool) {
	if existing, ok := componentMap[component.Name]; ok {
		appendMissingStates(existing, component.States)
		if preferFile {
			existing.File = component.File
		}
		return
	}
	componentMap[component.Name] = &component
}

// appendMissingStates appends states that are not already present.
func appendMissingStates(existing *ComponentInfo, states []StateInfo) {
	for _, state := range states {
		if !containsState(existing.States, state.Constant) {
			existing.States = append(existing.States, state)
			existing.StateCount++
		}
	}
}

// containsState reports whether the states slice already has the constant.
func containsState(states []StateInfo, constant string) bool {
	for _, state := range states {
		if state.Constant == constant {
			return true
		}
	}
	return false
}

// prefersIntentFile reports whether the file should be preferred for intent metadata.
func prefersIntentFile(path string) bool {
	return strings.HasSuffix(path, "_intent.go")
}

// appendIntents adds merged intents to the matrix and updates totals.
func appendIntents(matrix *StateMatrix, intentMap map[string]*ComponentInfo) {
	for _, intent := range intentMap {
		matrix.Intents = append(matrix.Intents, *intent)
		matrix.TotalStates += intent.StateCount
	}
	matrix.TotalIntents = len(matrix.Intents)
}

// appendScreens adds merged screens to the matrix and updates totals.
func appendScreens(matrix *StateMatrix, screenMap map[string]*ComponentInfo) {
	for _, screen := range screenMap {
		matrix.Screens = append(matrix.Screens, *screen)
		matrix.TotalStates += screen.StateCount
	}
	matrix.TotalScreens = len(matrix.Screens)
}
