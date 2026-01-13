package statematrix_test

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/statematrix"
)

func TestParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StateMatrix Parser Suite")
}

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
			// Create a temporary file with no state constants
			file := filepath.Join(projectRoot, "internal", "cli", "intents", "app.go")
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
			if err != nil || len(files) == 0 {
				Skip("No screen files exist yet")
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
				classification := statematrix.ClassifyState("BrowseStateEventDetail")
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
