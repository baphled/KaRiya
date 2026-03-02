package statematrix_test

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/statematrix"
)

var _ = Describe("Coverage Improvement", func() {
	Describe("GenerateMarkdown", func() {
		var tmpDir string

		BeforeEach(func() {
			var err error
			tmpDir, err = os.MkdirTemp("", "statematrix-cov-*")
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			os.RemoveAll(tmpDir)
		})

		Context("when a component has empty states", func() {
			It("should render mermaid diagram gracefully", func() {
				matrix := &statematrix.StateMatrix{
					GeneratedAt:  time.Now(),
					TotalIntents: 1,
					TotalStates:  0,
					Intents: []statematrix.ComponentInfo{
						{
							Name:       "EmptyStates",
							File:       "/path/to/empty.go",
							Kind:       "intent",
							StateCount: 0,
							States:     []statematrix.StateInfo{},
						},
					},
					Screens: []statematrix.ComponentInfo{},
				}

				mdPath := filepath.Join(tmpDir, "empty_states.md")
				err := statematrix.GenerateMarkdown(matrix, mdPath)
				Expect(err).ToNot(HaveOccurred())

				content, err := os.ReadFile(mdPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(content)).ToNot(ContainSubstring("stateDiagram-v2"))
			})
		})

		Context("when multiple screens need sorting", func() {
			It("should sort screens alphabetically", func() {
				matrix := &statematrix.StateMatrix{
					GeneratedAt:  time.Now(),
					TotalIntents: 0,
					TotalScreens: 3,
					TotalStates:  3,
					Intents:      []statematrix.ComponentInfo{},
					Screens: []statematrix.ComponentInfo{
						{
							Name:       "ZebraScreen",
							File:       "/path/to/zebra.go",
							Kind:       "screen",
							StateCount: 1,
							States: []statematrix.StateInfo{
								{Constant: "ZebraActive", Type: "ROOT", EscapeBehavior: "Cancel"},
							},
						},
						{
							Name:       "AlphaScreen",
							File:       "/path/to/alpha.go",
							Kind:       "screen",
							StateCount: 1,
							States: []statematrix.StateInfo{
								{Constant: "AlphaActive", Type: "ROOT", EscapeBehavior: "Cancel"},
							},
						},
						{
							Name:       "BetaScreen",
							File:       "/path/to/beta.go",
							Kind:       "screen",
							StateCount: 1,
							States: []statematrix.StateInfo{
								{Constant: "BetaActive", Type: "ROOT", EscapeBehavior: "Cancel"},
							},
						},
					},
				}

				mdPath := filepath.Join(tmpDir, "sorted_screens.md")
				err := statematrix.GenerateMarkdown(matrix, mdPath)
				Expect(err).ToNot(HaveOccurred())

				content, err := os.ReadFile(mdPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(content)).To(ContainSubstring("AlphaScreen"))
				Expect(string(content)).To(ContainSubstring("BetaScreen"))
				Expect(string(content)).To(ContainSubstring("ZebraScreen"))
			})
		})

		Context("when path is invalid", func() {
			It("should return error for non-writable path", func() {
				matrix := &statematrix.StateMatrix{
					GeneratedAt: time.Now(),
					Intents:     []statematrix.ComponentInfo{},
					Screens:     []statematrix.ComponentInfo{},
				}

				err := statematrix.GenerateMarkdown(matrix, "/non/existent/dir/file.md")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GenerateJSON", func() {
		Context("when path is invalid", func() {
			It("should return error for non-writable path", func() {
				matrix := &statematrix.StateMatrix{
					GeneratedAt: time.Now(),
					Intents:     []statematrix.ComponentInfo{},
					Screens:     []statematrix.ComponentInfo{},
				}

				err := statematrix.GenerateJSON(matrix, "/non/existent/dir/file.json")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("ParseIntentFile", func() {
		Context("when file cannot be parsed", func() {
			It("should return empty component for non-existent file", func() {
				component := statematrix.ParseIntentFile("/non/existent/file.go")
				Expect(component.Kind).To(Equal("intent"))
				Expect(component.StateCount).To(Equal(0))
				Expect(component.States).To(BeEmpty())
			})
		})
	})

	Describe("ParseScreenFile", func() {
		Context("when file cannot be parsed", func() {
			It("should return empty component for non-existent file", func() {
				component := statematrix.ParseScreenFile("/non/existent/file.go")
				Expect(component.Kind).To(Equal("screen"))
				Expect(component.StateCount).To(Equal(0))
				Expect(component.States).To(BeEmpty())
			})
		})
	})

	Describe("ClassifyState", func() {
		Context("Modal states", func() {
			It("should identify modal states as Modal", func() {
				classification := statematrix.ClassifyState("SomeModalState")
				Expect(classification).To(Equal("Modal"))
			})
		})

		Context("timeline detail states", func() {
			It("should classify timeline detail as Intermediate", func() {
				classification := statematrix.ClassifyState("TimelineDetailState")
				Expect(classification).To(Equal("Intermediate"))
			})
		})

		Context("done states", func() {
			It("should classify done states as Final", func() {
				classification := statematrix.ClassifyState("WorkflowDoneState")
				Expect(classification).To(Equal("Final"))
			})
		})

		Context("error states", func() {
			It("should classify error states as Error", func() {
				classification := statematrix.ClassifyState("SomeErrorState")
				Expect(classification).To(Equal("Error"))
			})
		})

		Context("saving states", func() {
			It("should classify saving states as Async", func() {
				classification := statematrix.ClassifyState("DataSavingState")
				Expect(classification).To(Equal("Async"))
			})
		})

		Context("extracting states", func() {
			It("should classify extracting states as Async", func() {
				classification := statematrix.ClassifyState("DataExtractingState")
				Expect(classification).To(Equal("Async"))
			})
		})

		Context("inprogress states", func() {
			It("should classify inprogress states as Async", func() {
				classification := statematrix.ClassifyState("TaskInProgressState")
				Expect(classification).To(Equal("Async"))
			})
		})

		Context("exporting states", func() {
			It("should classify exporting states as Async", func() {
				classification := statematrix.ClassifyState("DataExportingState")
				Expect(classification).To(Equal("Async"))
			})
		})

		Context("selectop states", func() {
			It("should classify selectop states as ROOT", func() {
				classification := statematrix.ClassifyState("SomeSelectOpState")
				Expect(classification).To(Equal("ROOT"))
			})
		})

		Context("selecttype states", func() {
			It("should classify selecttype states as ROOT", func() {
				classification := statematrix.ClassifyState("SomeSelectTypeState")
				Expect(classification).To(Equal("ROOT"))
			})
		})

		Context("selectprofile states", func() {
			It("should classify selectprofile states as ROOT", func() {
				classification := statematrix.ClassifyState("SomeSelectProfileState")
				Expect(classification).To(Equal("ROOT"))
			})
		})

		Context("fileselect states", func() {
			It("should classify fileselect states as ROOT", func() {
				classification := statematrix.ClassifyState("SomeFileSelectState")
				Expect(classification).To(Equal("ROOT"))
			})
		})
	})

	Describe("InferEscapeBehavior", func() {
		Context("Confirmation state type", func() {
			It("should return cancel to parent behavior", func() {
				behavior := statematrix.InferEscapeBehavior("Confirmation")
				Expect(behavior).To(ContainSubstring("Cancel"))
				Expect(behavior).To(ContainSubstring("Parent state"))
			})
		})

		Context("Modal state type", func() {
			It("should return close modal behavior", func() {
				behavior := statematrix.InferEscapeBehavior("Modal")
				Expect(behavior).To(ContainSubstring("Close modal"))
			})
		})

		Context("Final state type", func() {
			It("should return deactivate behavior", func() {
				behavior := statematrix.InferEscapeBehavior("Final")
				Expect(behavior).To(ContainSubstring("Deactivate intent"))
			})
		})

		Context("Error state type", func() {
			It("should return deactivate with retry behavior", func() {
				behavior := statematrix.InferEscapeBehavior("Error")
				Expect(behavior).To(ContainSubstring("Deactivate intent"))
				Expect(behavior).To(ContainSubstring("retry"))
			})
		})

		Context("unknown state type", func() {
			It("should default to back previous state behavior", func() {
				behavior := statematrix.InferEscapeBehavior("UnknownType")
				Expect(behavior).To(ContainSubstring("Back"))
				Expect(behavior).To(ContainSubstring("Previous state"))
			})
		})
	})

	Describe("FindScreenFiles", func() {
		Context("when directory does not exist", func() {
			It("should return empty list without error", func() {
				files, err := statematrix.FindScreenFiles("/non/existent/screens/path")
				Expect(err).ToNot(HaveOccurred())
				Expect(files).To(BeEmpty())
			})
		})
	})

	Describe("ScanAll", func() {
		Context("when intents directory is invalid", func() {
			It("should return error", func() {
				_, err := statematrix.ScanAll("/non/existent/intents", "/non/existent/screens")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when screen files have overlapping states", func() {
			var tmpScreensDir, tmpIntentsDir string

			BeforeEach(func() {
				var err error
				tmpIntentsDir, err = os.MkdirTemp("", "sm-cov-intents-*")
				Expect(err).ToNot(HaveOccurred())

				tmpScreensDir, err = os.MkdirTemp("", "sm-cov-screens-*")
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				os.RemoveAll(tmpIntentsDir)
				os.RemoveAll(tmpScreensDir)
			})

			It("should deduplicate screen states from multiple files", func() {
				subDir := filepath.Join(tmpScreensDir, "sub")
				err := os.MkdirAll(subDir, 0o755)
				Expect(err).ToNot(HaveOccurred())

				err = os.WriteFile(filepath.Join(tmpScreensDir, "test_screen.go"), []byte(`package screens

const (
	TestScreenStateActive = iota
	TestScreenStateLoading
)
`), 0o600)
				Expect(err).ToNot(HaveOccurred())

				err = os.WriteFile(filepath.Join(subDir, "test_screen.go"), []byte(`package screens

const (
	TestScreenStateActive = iota
	TestScreenStateError
)
`), 0o600)
				Expect(err).ToNot(HaveOccurred())

				matrix, err := statematrix.ScanAll(tmpIntentsDir, tmpScreensDir)
				Expect(err).ToNot(HaveOccurred())
				Expect(matrix.TotalScreens).To(Equal(1))

				screen := matrix.Screens[0]
				Expect(screen.StateCount).To(Equal(3))

				stateNames := make([]string, len(screen.States))
				for i, s := range screen.States {
					stateNames[i] = s.Constant
				}
				Expect(stateNames).To(ContainElement("TestScreenStateActive"))
				Expect(stateNames).To(ContainElement("TestScreenStateLoading"))
				Expect(stateNames).To(ContainElement("TestScreenStateError"))
			})
		})

		Context("when intent file has _intent.go suffix", func() {
			var tmpIntentsDir, tmpScreensDir string

			BeforeEach(func() {
				var err error
				tmpIntentsDir, err = os.MkdirTemp("", "sm-cov-intents2-*")
				Expect(err).ToNot(HaveOccurred())

				tmpScreensDir, err = os.MkdirTemp("", "sm-cov-screens2-*")
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				os.RemoveAll(tmpIntentsDir)
				os.RemoveAll(tmpScreensDir)
			})

			It("should prefer _intent.go file path", func() {
				intentDir := filepath.Join(tmpIntentsDir, "myworkflow")
				err := os.MkdirAll(intentDir, 0o755)
				Expect(err).ToNot(HaveOccurred())

				err = os.WriteFile(filepath.Join(intentDir, "constants.go"), []byte(`package myworkflow

const (
	StateListItems = iota
	StateForm
)
`), 0o600)
				Expect(err).ToNot(HaveOccurred())

				err = os.WriteFile(filepath.Join(tmpIntentsDir, "myworkflow_intent.go"), []byte(`package intents

const (
	MyworkflowStateListItems = iota
	MyworkflowStateReview
)
`), 0o600)
				Expect(err).ToNot(HaveOccurred())

				matrix, err := statematrix.ScanAll(tmpIntentsDir, tmpScreensDir)
				Expect(err).ToNot(HaveOccurred())

				found := false
				for _, intent := range matrix.Intents {
					if intent.Name == "Myworkflow" {
						found = true
						Expect(intent.File).To(HaveSuffix("_intent.go"))
					}
				}
				Expect(found).To(BeTrue(), "Should find Myworkflow intent")
			})
		})
	})
})
