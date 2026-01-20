package components_test

import (
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/types"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExportConfigWizardModal", func() {
	var modal *components.ExportConfigWizardModal

	BeforeEach(func() {
		modal = components.NewExportConfigWizardModal(100, 40)
	})

	Describe("NewExportConfigWizardModal", func() {
		It("creates a modal with default dimensions", func() {
			Expect(modal).NotTo(BeNil())
			width, height := modal.GetDimensions()
			Expect(width).To(Equal(100))
			Expect(height).To(Equal(40))
		})

		It("starts visible", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("starts not completed", func() {
			Expect(modal.IsCompleted()).To(BeFalse())
		})

		It("starts not cancelled", func() {
			Expect(modal.IsCancelled()).To(BeFalse())
		})

		It("starts not skipped", func() {
			Expect(modal.IsSkipped()).To(BeFalse())
		})

		It("starts at step 0", func() {
			Expect(modal.GetCurrentStep()).To(Equal(0))
		})

		It("has 2 steps", func() {
			Expect(modal.GetStepCount()).To(Equal(2))
		})

		It("initializes with default config data", func() {
			data := modal.GetConfigData()
			Expect(data).NotTo(BeNil())
			Expect(data.ArtifactType).To(Equal("events"))
			Expect(data.Format).To(Equal("json"))
			Expect(data.Destination).To(Equal("file"))
		})
	})

	Describe("Init", func() {
		It("returns a command", func() {
			cmd := modal.Init()
			// Init may return a command to start the form
			// We just verify it doesn't panic
			_ = cmd
		})
	})

	Describe("Visibility", func() {
		It("can be hidden", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("can be shown", func() {
			modal.Hide()
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("returns empty view when hidden", func() {
			modal.Hide()
			Expect(modal.View()).To(BeEmpty())
		})
	})

	Describe("View", func() {
		It("renders the modal when visible", func() {
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("contains the title", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Export Configuration"))
		})

		It("shows step 1 title initially", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("What to Export"))
		})

		It("contains step 1 group title", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Step 1: What to Export"))
		})
	})

	Describe("Update", func() {
		Context("when modal is hidden", func() {
			It("returns nil command", func() {
				modal.Hide()
				cmd := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).To(BeNil())
			})
		})

		Context("with Escape key", func() {
			It("cancels the wizard when at step 0", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(modal.IsCancelled()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("with Ctrl+S", func() {
			It("skips the wizard and applies defaults", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
				Expect(modal.IsSkipped()).To(BeTrue())
				Expect(modal.IsCompleted()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})

			It("preserves default values when skipping", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
				data := modal.GetConfigData()
				Expect(data.ArtifactType).To(Equal("events"))
				Expect(data.Format).To(Equal("json"))
				Expect(data.Destination).To(Equal("file"))
			})
		})

		Context("with WindowSizeMsg", func() {
			It("updates dimensions", func() {
				modal.Update(tea.WindowSizeMsg{Width: 120, Height: 50})
				width, height := modal.GetDimensions()
				Expect(width).To(Equal(120))
				Expect(height).To(Equal(50))
			})
		})
	})

	Describe("Reset", func() {
		It("resets completion state", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(modal.IsCompleted()).To(BeTrue())

			modal.Reset()
			Expect(modal.IsCompleted()).To(BeFalse())
		})

		It("resets cancelled state", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.IsCancelled()).To(BeTrue())

			modal.Reset()
			Expect(modal.IsCancelled()).To(BeFalse())
		})

		It("resets skipped state", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(modal.IsSkipped()).To(BeTrue())

			modal.Reset()
			Expect(modal.IsSkipped()).To(BeFalse())
		})

		It("makes modal visible again", func() {
			modal.Hide()
			modal.Reset()
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("resets to step 0", func() {
			// Note: Step progression is tracked internally
			modal.Reset()
			Expect(modal.GetCurrentStep()).To(Equal(0))
		})
	})

	Describe("Setters", func() {
		It("sets artifact type", func() {
			modal.SetArtifactType("facts")
			Expect(modal.GetConfigData().ArtifactType).To(Equal("facts"))
		})

		It("sets format", func() {
			modal.SetFormat("yaml")
			Expect(modal.GetConfigData().Format).To(Equal("yaml"))
		})

		It("sets destination", func() {
			modal.SetDestination("clipboard")
			Expect(modal.GetConfigData().Destination).To(Equal("clipboard"))
		})
	})

	Describe("Type conversions", func() {
		It("returns artifact type as ExportArtifactType", func() {
			modal.SetArtifactType("events")
			Expect(modal.GetArtifactType()).To(Equal(types.ExportTypeEvents))
		})

		It("returns format as ExportFormat", func() {
			modal.SetFormat("json")
			Expect(modal.GetFormat()).To(Equal(types.ExportFormatJSON))
		})

		It("returns destination as ExportDestination", func() {
			modal.SetDestination("file")
			Expect(modal.GetDestination()).To(Equal(types.ExportDestinationFile))
		})

		It("handles all artifact types", func() {
			testCases := []struct {
				value    string
				expected types.ExportArtifactType
			}{
				{"events", types.ExportTypeEvents},
				{"facts", types.ExportTypeFacts},
				{"bursts", types.ExportTypeBursts},
				{"cv", types.ExportTypeCV},
				{"profile", types.ExportTypeProfile},
			}

			for _, tc := range testCases {
				modal.SetArtifactType(tc.value)
				Expect(modal.GetArtifactType()).To(Equal(tc.expected))
			}
		})

		It("handles all formats", func() {
			testCases := []struct {
				value    string
				expected types.ExportFormat
			}{
				{"json", types.ExportFormatJSON},
				{"yaml", types.ExportFormatYAML},
				{"csv", types.ExportFormatCSV},
				{"txt", types.ExportFormatTXT},
				{"markdown", types.ExportFormatMD},
			}

			for _, tc := range testCases {
				modal.SetFormat(tc.value)
				Expect(modal.GetFormat()).To(Equal(tc.expected))
			}
		})

		It("handles all destinations", func() {
			testCases := []struct {
				value    string
				expected types.ExportDestination
			}{
				{"file", types.ExportDestinationFile},
				{"clipboard", types.ExportDestinationClipboard},
			}

			for _, tc := range testCases {
				modal.SetDestination(tc.value)
				Expect(modal.GetDestination()).To(Equal(tc.expected))
			}
		})
	})

	Describe("Small terminal dimensions", func() {
		It("handles very small width", func() {
			smallModal := components.NewExportConfigWizardModal(40, 20)
			view := smallModal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("handles very small height", func() {
			smallModal := components.NewExportConfigWizardModal(100, 10)
			view := smallModal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})
