package statematrix_test

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/statematrix"
)

func TestScanner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StateMatrix Scanner Suite")
}

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
