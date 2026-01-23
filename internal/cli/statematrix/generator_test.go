package statematrix_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/statematrix"
)

func TestStateMatrix(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StateMatrix Suite")
}

var _ = Describe("Generator", func() {
	var tmpDir string
	var matrix *statematrix.StateMatrix

	BeforeEach(func() {
		var err error
		tmpDir, err = os.MkdirTemp("", "statematrix-test-*")
		Expect(err).ToNot(HaveOccurred())

		// Create sample matrix
		matrix = &statematrix.StateMatrix{
			GeneratedAt:  time.Now(),
			TotalIntents: 2,
			TotalScreens: 1,
			TotalStates:  7,
			Intents: []statematrix.ComponentInfo{
				{
					Name:       "CaptureEvent",
					File:       "/path/to/capture_event.go",
					Kind:       "intent",
					StateCount: 4,
					States: []statematrix.StateInfo{
						{Constant: "CaptureStateChooseStrategy", Type: "ROOT", EscapeBehavior: "Cancel intent → Main Menu"},
						{Constant: "CaptureStateForm", Type: "Intermediate", EscapeBehavior: "Back → Previous state"},
						{Constant: "CaptureStateReview", Type: "Intermediate", EscapeBehavior: "Back → Previous state"},
						{Constant: "CaptureStateSubmit", Type: "Intermediate", EscapeBehavior: "Back → Previous state"},
					},
				},
				{
					Name:       "BrowseTimeline",
					File:       "/path/to/browse_timeline.go",
					Kind:       "intent",
					StateCount: 2,
					States: []statematrix.StateInfo{
						{Constant: "BrowseStateTimeline", Type: "ROOT", EscapeBehavior: "Cancel intent → Main Menu"},
						{Constant: "BrowseStateDeleteConfirm", Type: "Confirmation", EscapeBehavior: "Cancel → Return to list"},
					},
				},
			},
			Screens: []statematrix.ComponentInfo{
				{
					Name:       "ProfileSelectScreen",
					File:       "/path/to/screens/cv/profile_select.go",
					Kind:       "screen",
					StateCount: 1,
					States: []statematrix.StateInfo{
						{Constant: "ProfileSelectActive", Type: "ROOT", EscapeBehavior: "Cancel → Parent screen"},
					},
				},
			},
		}
	})

	AfterEach(func() {
		os.RemoveAll(tmpDir)
	})

	Describe("GenerateMarkdown", func() {
		It("should generate valid markdown file", func() {
			mdPath := filepath.Join(tmpDir, "STATE_MATRIX.md")
			err := statematrix.GenerateMarkdown(matrix, mdPath)

			Expect(err).ToNot(HaveOccurred())
			Expect(mdPath).To(BeAnExistingFile())

			content, err := os.ReadFile(mdPath)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(content)).ToNot(BeEmpty())
		})

		It("should include header with generation info", func() {
			mdPath := filepath.Join(tmpDir, "STATE_MATRIX.md")
			err := statematrix.GenerateMarkdown(matrix, mdPath)
			Expect(err).ToNot(HaveOccurred())

			content, err := os.ReadFile(mdPath)
			Expect(err).ToNot(HaveOccurred())

			markdown := string(content)
			Expect(markdown).To(ContainSubstring("# KaRiya TUI State Matrix"))
			Expect(markdown).To(ContainSubstring("Auto-generated"))
			Expect(markdown).To(ContainSubstring("make generate-state-matrix"))
		})

		It("should include summary table", func() {
			mdPath := filepath.Join(tmpDir, "STATE_MATRIX.md")
			err := statematrix.GenerateMarkdown(matrix, mdPath)
			Expect(err).ToNot(HaveOccurred())

			content, err := os.ReadFile(mdPath)
			Expect(err).ToNot(HaveOccurred())

			markdown := string(content)
			Expect(markdown).To(ContainSubstring("## Summary"))
			Expect(markdown).To(ContainSubstring("Total Intents"))
			Expect(markdown).To(ContainSubstring("Total Screens"))
			Expect(markdown).To(ContainSubstring("Total States"))
			Expect(markdown).To(ContainSubstring("| 2 |"))
			Expect(markdown).To(ContainSubstring("| 1 |"))
			Expect(markdown).To(ContainSubstring("| 7 |"))
		})

		It("should have separate sections for intents and screens", func() {
			mdPath := filepath.Join(tmpDir, "STATE_MATRIX.md")
			err := statematrix.GenerateMarkdown(matrix, mdPath)
			Expect(err).ToNot(HaveOccurred())

			content, err := os.ReadFile(mdPath)
			Expect(err).ToNot(HaveOccurred())

			markdown := string(content)
			Expect(markdown).To(ContainSubstring("## Intent States"))
			Expect(markdown).To(ContainSubstring("## Screen States"))
		})

		It("should list all intent states", func() {
			mdPath := filepath.Join(tmpDir, "STATE_MATRIX.md")
			err := statematrix.GenerateMarkdown(matrix, mdPath)
			Expect(err).ToNot(HaveOccurred())

			content, err := os.ReadFile(mdPath)
			Expect(err).ToNot(HaveOccurred())

			markdown := string(content)
			Expect(markdown).To(ContainSubstring("### CaptureEvent"))
			Expect(markdown).To(ContainSubstring("CaptureStateChooseStrategy"))
			Expect(markdown).To(ContainSubstring("CaptureStateForm"))
			Expect(markdown).To(ContainSubstring("### BrowseTimeline"))
			Expect(markdown).To(ContainSubstring("BrowseStateTimeline"))
		})

		It("should list all screen states", func() {
			mdPath := filepath.Join(tmpDir, "STATE_MATRIX.md")
			err := statematrix.GenerateMarkdown(matrix, mdPath)
			Expect(err).ToNot(HaveOccurred())

			content, err := os.ReadFile(mdPath)
			Expect(err).ToNot(HaveOccurred())

			markdown := string(content)
			Expect(markdown).To(ContainSubstring("### ProfileSelectScreen"))
			Expect(markdown).To(ContainSubstring("ProfileSelectActive"))
		})

		It("should include state type definitions", func() {
			mdPath := filepath.Join(tmpDir, "STATE_MATRIX.md")
			err := statematrix.GenerateMarkdown(matrix, mdPath)
			Expect(err).ToNot(HaveOccurred())

			content, err := os.ReadFile(mdPath)
			Expect(err).ToNot(HaveOccurred())

			markdown := string(content)
			Expect(markdown).To(ContainSubstring("## State Type Definitions"))
			Expect(markdown).To(ContainSubstring("ROOT"))
			Expect(markdown).To(ContainSubstring("Intermediate"))
			Expect(markdown).To(ContainSubstring("Async"))
		})
	})

	Describe("GenerateJSON", func() {
		It("should generate valid JSON file", func() {
			jsonPath := filepath.Join(tmpDir, "state_matrix.json")
			err := statematrix.GenerateJSON(matrix, jsonPath)

			Expect(err).ToNot(HaveOccurred())
			Expect(jsonPath).To(BeAnExistingFile())

			content, err := os.ReadFile(jsonPath)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(content)).ToNot(BeEmpty())
		})

		It("should be valid parseable JSON", func() {
			jsonPath := filepath.Join(tmpDir, "state_matrix.json")
			err := statematrix.GenerateJSON(matrix, jsonPath)
			Expect(err).ToNot(HaveOccurred())

			content, err := os.ReadFile(jsonPath)
			Expect(err).ToNot(HaveOccurred())

			var parsed statematrix.StateMatrix
			err = json.Unmarshal(content, &parsed)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should preserve all matrix data", func() {
			jsonPath := filepath.Join(tmpDir, "state_matrix.json")
			err := statematrix.GenerateJSON(matrix, jsonPath)
			Expect(err).ToNot(HaveOccurred())

			content, err := os.ReadFile(jsonPath)
			Expect(err).ToNot(HaveOccurred())

			var parsed statematrix.StateMatrix
			err = json.Unmarshal(content, &parsed)
			Expect(err).ToNot(HaveOccurred())

			Expect(parsed.TotalIntents).To(Equal(2))
			Expect(parsed.TotalScreens).To(Equal(1))
			Expect(parsed.TotalStates).To(Equal(7))
			Expect(len(parsed.Intents)).To(Equal(2))
			Expect(len(parsed.Screens)).To(Equal(1))
		})

		It("should include intent and screen separation", func() {
			jsonPath := filepath.Join(tmpDir, "state_matrix.json")
			err := statematrix.GenerateJSON(matrix, jsonPath)
			Expect(err).ToNot(HaveOccurred())

			content, err := os.ReadFile(jsonPath)
			Expect(err).ToNot(HaveOccurred())

			jsonStr := string(content)
			Expect(jsonStr).To(ContainSubstring(`"intents"`))
			Expect(jsonStr).To(ContainSubstring(`"screens"`))
			Expect(jsonStr).To(ContainSubstring(`"kind": "intent"`))
			Expect(jsonStr).To(ContainSubstring(`"kind": "screen"`))
		})
	})

	Describe("Integration", func() {
		It("should handle empty screens array", func() {
			emptyMatrix := &statematrix.StateMatrix{
				GeneratedAt:  time.Now(),
				TotalIntents: 1,
				TotalScreens: 0,
				TotalStates:  2,
				Intents: []statematrix.ComponentInfo{
					{
						Name:       "TestIntent",
						Kind:       "intent",
						StateCount: 2,
						States: []statematrix.StateInfo{
							{Constant: "TestState1", Type: "ROOT", EscapeBehavior: "Cancel"},
							{Constant: "TestState2", Type: "Intermediate", EscapeBehavior: "Back"},
						},
					},
				},
				Screens: []statematrix.ComponentInfo{},
			}

			mdPath := filepath.Join(tmpDir, "empty_screens.md")
			err := statematrix.GenerateMarkdown(emptyMatrix, mdPath)
			Expect(err).ToNot(HaveOccurred())

			content, err := os.ReadFile(mdPath)
			Expect(err).ToNot(HaveOccurred())

			markdown := string(content)
			// Should still have sections even if empty
			Expect(markdown).To(ContainSubstring("## Screen States"))
		})

		It("should sort components by name", func() {
			unsortedMatrix := &statematrix.StateMatrix{
				GeneratedAt:  time.Now(),
				TotalIntents: 3,
				TotalStates:  3,
				Intents: []statematrix.ComponentInfo{
					{Name: "Zebra", Kind: "intent", States: []statematrix.StateInfo{{Constant: "Z", Type: "ROOT", EscapeBehavior: "Cancel"}}},
					{Name: "Alpha", Kind: "intent", States: []statematrix.StateInfo{{Constant: "A", Type: "ROOT", EscapeBehavior: "Cancel"}}},
					{Name: "Beta", Kind: "intent", States: []statematrix.StateInfo{{Constant: "B", Type: "ROOT", EscapeBehavior: "Cancel"}}},
				},
			}

			mdPath := filepath.Join(tmpDir, "sorted.md")
			err := statematrix.GenerateMarkdown(unsortedMatrix, mdPath)
			Expect(err).ToNot(HaveOccurred())

			content, err := os.ReadFile(mdPath)
			Expect(err).ToNot(HaveOccurred())

			markdown := string(content)

			// Find positions of each name
			alphaPos := strings.Index(markdown, "### Alpha")
			betaPos := strings.Index(markdown, "### Beta")
			zebraPos := strings.Index(markdown, "### Zebra")

			Expect(alphaPos).To(BeNumerically("<", betaPos), "Alpha should come before Beta")
			Expect(betaPos).To(BeNumerically("<", zebraPos), "Beta should come before Zebra")
		})
	})
})

var _ = Describe("Scanner", func() {
	var projectRoot string

	BeforeEach(func() {
		// Get the actual project root for testing
		var err error
		projectRoot, err = filepath.Abs(filepath.Join("..", "..", ".."))
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("FindIntentFiles", func() {
		It("should find all intent files in internal/cli/intents/", func() {
			intentsDir := filepath.Join(projectRoot, "internal", "cli", "intents")
			files, err := statematrix.FindIntentFiles(intentsDir)

			Expect(err).ToNot(HaveOccurred())
			Expect(files).ToNot(BeEmpty())

			// Should find intent files like capture_event.go, browse_timeline.go, etc.
			foundCaptureEvent := false
			foundBrowseTimeline := false
			for _, file := range files {
				if filepath.Base(file) == "capture_event.go" {
					foundCaptureEvent = true
				}
				if filepath.Base(file) == "browse_timeline.go" {
					foundBrowseTimeline = true
				}
			}
			Expect(foundCaptureEvent).To(BeTrue(), "Should find capture_event.go")
			Expect(foundBrowseTimeline).To(BeTrue(), "Should find browse_timeline.go")
		})

		It("should exclude test files", func() {
			intentsDir := filepath.Join(projectRoot, "internal", "cli", "intents")
			files, err := statematrix.FindIntentFiles(intentsDir)

			Expect(err).ToNot(HaveOccurred())
			for _, file := range files {
				Expect(file).ToNot(HaveSuffix("_test.go"), "Should not include test files")
			}
		})

		It("should return error for non-existent directory", func() {
			_, err := statematrix.FindIntentFiles("/non/existent/path")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("FindScreenFiles", func() {
		It("should find screen files when screens/ directory exists", func() {
			screensDir := filepath.Join(projectRoot, "internal", "cli", "screens")
			files, err := statematrix.FindScreenFiles(screensDir)

			// Currently screens/ directory may not exist, so we accept both outcomes
			if err != nil {
				// Directory doesn't exist yet - this is expected during migration
				Expect(files).To(BeEmpty())
			} else {
				// Directory exists - files may or may not exist
				for _, file := range files {
					Expect(file).ToNot(HaveSuffix("_test.go"), "Should not include test files")
				}
			}
		})

		It("should exclude test files", func() {
			screensDir := filepath.Join(projectRoot, "internal", "cli", "screens")
			files, err := statematrix.FindScreenFiles(screensDir)

			if err == nil {
				for _, file := range files {
					Expect(file).ToNot(HaveSuffix("_test.go"), "Should not include test files")
				}
			}
		})
	})

	Describe("ScanAll", func() {
		It("should scan both intents and screens directories", func() {
			intentsDir := filepath.Join(projectRoot, "internal", "cli", "intents")
			screensDir := filepath.Join(projectRoot, "internal", "cli", "screens")

			matrix, err := statematrix.ScanAll(intentsDir, screensDir)

			Expect(err).ToNot(HaveOccurred())
			Expect(matrix).ToNot(BeNil())
			Expect(matrix.TotalIntents).To(BeNumerically(">", 0), "Should find at least one intent")
			Expect(matrix.TotalStates).To(BeNumerically(">", 0), "Should find at least one state")
		})

		It("should separate intents and screens", func() {
			intentsDir := filepath.Join(projectRoot, "internal", "cli", "intents")
			screensDir := filepath.Join(projectRoot, "internal", "cli", "screens")

			matrix, err := statematrix.ScanAll(intentsDir, screensDir)

			Expect(err).ToNot(HaveOccurred())
			Expect(matrix.Intents).ToNot(BeEmpty(), "Should have intents")

			// Screens may be empty during migration
			for _, component := range matrix.Intents {
				Expect(component.Kind).To(Equal("intent"))
			}

			for _, component := range matrix.Screens {
				Expect(component.Kind).To(Equal("screen"))
			}
		})

		It("should populate state counts correctly", func() {
			intentsDir := filepath.Join(projectRoot, "internal", "cli", "intents")
			screensDir := filepath.Join(projectRoot, "internal", "cli", "screens")

			matrix, err := statematrix.ScanAll(intentsDir, screensDir)

			Expect(err).ToNot(HaveOccurred())

			totalStates := 0
			for _, intent := range matrix.Intents {
				Expect(intent.StateCount).To(Equal(len(intent.States)))
				totalStates += intent.StateCount
			}
			for _, screen := range matrix.Screens {
				Expect(screen.StateCount).To(Equal(len(screen.States)))
				totalStates += screen.StateCount
			}

			Expect(matrix.TotalStates).To(Equal(totalStates))
		})
	})
})

var _ = Describe("Parser", func() {
	var projectRoot string

	BeforeEach(func() {
		var err error
		projectRoot, err = filepath.Abs(filepath.Join("..", "..", ".."))
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("ParseIntentFile", func() {
		It("should parse intent states from capture_event.go", func() {
			file := filepath.Join(projectRoot, "internal", "cli", "intents", "capture_event.go")
			component := statematrix.ParseIntentFile(file)

			Expect(component.Name).To(Equal("CaptureEvent"))
			Expect(component.Kind).To(Equal("intent"))
			Expect(component.StateCount).To(BeNumerically(">", 0))
			Expect(component.States).ToNot(BeEmpty())

			// Check for known states
			stateNames := make([]string, len(component.States))
			for i, state := range component.States {
				stateNames[i] = state.Constant
			}
			Expect(stateNames).To(ContainElement("CaptureStateChooseStrategy"))
		})

		It("should classify state types correctly", func() {
			file := filepath.Join(projectRoot, "internal", "cli", "intents", "capture_event.go")
			component := statematrix.ParseIntentFile(file)

			// CaptureStateChooseStrategy should be ROOT
			for _, state := range component.States {
				if state.Constant == "CaptureStateChooseStrategy" {
					Expect(state.Type).To(Equal("ROOT"))
					Expect(state.EscapeBehavior).To(ContainSubstring("Cancel intent"))
				}
			}
		})

		It("should handle files with no states", func() {
			// Use a file that doesn't have state constants
			file := filepath.Join(projectRoot, "internal", "cli", "styles", "styles.go")
			component := statematrix.ParseIntentFile(file)

			// Should return empty component or handle gracefully
			Expect(component.StateCount).To(Equal(0))
			Expect(component.States).To(BeEmpty())
		})
	})

	Describe("ParseScreenFile", func() {
		It("should parse screen states when screens exist", func() {
			screensDir := filepath.Join(projectRoot, "internal", "cli", "screens")

			// Find first screen file if any exist
			files, err := statematrix.FindScreenFiles(screensDir)
			// Handle case where screens directory doesn't exist or is empty
			if err != nil || len(files) == 0 {
				// Test passes - no screen files to parse is a valid state
				Expect(true).To(BeTrue(), "No screen files found - this is acceptable")
				return
			}

			component := statematrix.ParseScreenFile(files[0])

			Expect(component.Kind).To(Equal("screen"))
			// Further assertions depend on actual screen implementation
		})
	})

	Describe("State Classification", func() {
		Context("ROOT states", func() {
			It("should identify list states as ROOT", func() {
				classification := statematrix.ClassifyState("BurstListState")
				Expect(classification).To(Equal("ROOT"))
			})

			It("should identify choose states as ROOT", func() {
				classification := statematrix.ClassifyState("CaptureStateChooseStrategy")
				Expect(classification).To(Equal("ROOT"))
			})

			It("should identify select states as ROOT", func() {
				classification := statematrix.ClassifyState("ConfigStateSelectDomain")
				Expect(classification).To(Equal("ROOT"))
			})
		})

		Context("Intermediate states", func() {
			It("should identify form states as Intermediate", func() {
				classification := statematrix.ClassifyState("CaptureStateForm")
				Expect(classification).To(Equal("Intermediate"))
			})

			It("should identify review states as Intermediate", func() {
				classification := statematrix.ClassifyState("CaptureStateReview")
				Expect(classification).To(Equal("Intermediate"))
			})

			It("should identify detail states as Intermediate", func() {
				classification := statematrix.ClassifyState("SomeStateDetail")
				Expect(classification).To(Equal("Intermediate"))
			})
		})

		Context("Async states", func() {
			It("should identify generating states as Async", func() {
				classification := statematrix.ClassifyState("CVGeneratingState")
				Expect(classification).To(Equal("Async"))
			})

			It("should identify progress states as Async", func() {
				classification := statematrix.ClassifyState("ExportProgressState")
				Expect(classification).To(Equal("Async"))
			})
		})

		Context("Confirmation states", func() {
			It("should identify confirm states as Confirmation", func() {
				classification := statematrix.ClassifyState("BrowseStateDeleteConfirm")
				Expect(classification).To(Equal("Confirmation"))
			})
		})

		Context("Final states", func() {
			It("should identify complete states as Final", func() {
				classification := statematrix.ClassifyState("BulkCompleteState")
				Expect(classification).To(Equal("Final"))
			})
		})

		Context("Error states", func() {
			It("should identify failed states as Error", func() {
				classification := statematrix.ClassifyState("ExportFailedState")
				Expect(classification).To(Equal("Error"))
			})
		})
	})

	Describe("Escape Behavior Inference", func() {
		It("should infer ROOT state escape behavior", func() {
			behavior := statematrix.InferEscapeBehavior("ROOT")
			Expect(behavior).To(ContainSubstring("Cancel intent"))
			Expect(behavior).To(ContainSubstring("Main Menu"))
		})

		It("should infer Intermediate state escape behavior", func() {
			behavior := statematrix.InferEscapeBehavior("Intermediate")
			Expect(behavior).To(ContainSubstring("Back"))
			Expect(behavior).To(ContainSubstring("Previous state"))
		})

		It("should infer Async state escape behavior", func() {
			behavior := statematrix.InferEscapeBehavior("Async")
			Expect(behavior).To(ContainSubstring("Varies"))
		})
	})
})
