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

func TestGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StateMatrix Generator Suite")
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
						{Constant: "BrowseStateEventDetail", Type: "Intermediate", EscapeBehavior: "Back → Previous state"},
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
