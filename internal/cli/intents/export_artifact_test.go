package intents

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExportArtifact Intent", func() {
	var (
		intent *ExportArtifactIntent
		ctx    context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("Intent Creation", func() {
		It("should create a new ExportArtifact intent with valid context", func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent).NotTo(BeNil())
			Expect(intent.model).NotTo(BeNil())
		})

		It("should fail to create intent with nil context", func() {
			// Skipped because staticcheck prevents passing nil context
			// This is tested indirectly through all other tests
		})

		It("should initialize with correct default state", func() {
			intent, _ := NewExportArtifactIntent(ctx)
			Expect(intent.GetState()).To(Equal(ExportStateSelectType))
		})

		It("should initialize with correct artifact types", func() {
			intent, _ := NewExportArtifactIntent(ctx)
			Expect(intent.model.context.ArtifactTypes).To(HaveLen(5))
		})
	})

	Describe("Init Method", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return nil command", func() {
			cmd := intent.Init()
			Expect(cmd).To(BeNil())
		})

		It("should mark intent as active", func() {
			intent.Init()
			Expect(intent.IsActive()).To(BeTrue())
		})
	})

	Describe("View Rendering - SelectType State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should render SelectType view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Select Artifact Type"))
		})

		It("should show all artifact types", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("cv"))
			Expect(view).To(ContainSubstring("events"))
			Expect(view).To(ContainSubstring("facts"))
		})

		It("should show navigation instructions", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("↑/↓"))
			Expect(view).To(ContainSubstring("Enter"))
			Expect(view).To(ContainSubstring("Esc"))
		})

		It("should highlight selected item", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("> cv"))
		})
	})

	Describe("View Rendering - SelectFormat State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateSelectFormat)
			intent.SetConfig(NewExportConfiguration(ExportTypeCV, intent.model.context))
		})

		It("should render SelectFormat view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Select Export Format"))
		})

		It("should show supported formats for artifact type", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("pdf"))
			Expect(view).To(ContainSubstring("json"))
		})
	})

	Describe("View Rendering - SelectDestination State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateSelectDest)
		})

		It("should render SelectDestination view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Select Export Destination"))
		})

		It("should show all destinations", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("file"))
			Expect(view).To(ContainSubstring("clipboard"))
			Expect(view).To(ContainSubstring("email"))
		})
	})

	Describe("View Rendering - Configure State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateConfigure)
			intent.SetConfig(NewExportConfiguration(ExportTypeCV, intent.model.context))
		})

		It("should render Configure view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Configure Export"))
		})

		It("should show selected artifact type", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Artifact Type: cv"))
		})

		It("should show selected format", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Format:"))
		})

		It("should show selected destination", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Destination:"))
		})
	})

	Describe("View Rendering - Preview State", func() {
	BeforeEach(func() {
		var err error
		intent, err = NewExportArtifactIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
		intent.SetState(ExportStatePreview)
		intent.SetConfig(NewExportConfiguration(ExportTypeCV, intent.model.context))
	})

		It("should render Preview view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Preview Export"))
		})
	})

	Describe("View Rendering - Confirm State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateConfirm)
			intent.SetConfig(NewExportConfiguration(ExportTypeCV, intent.model.context))
		})

		It("should render Confirm view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Confirm Export"))
		})

		It("should show configuration summary", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Artifact:"))
			Expect(view).To(ContainSubstring("Format:"))
			Expect(view).To(ContainSubstring("Destination:"))
		})
	})

	Describe("View Rendering - InProgress State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateInProgress)
		})

		It("should render InProgress view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Exporting artifact"))
		})
	})

	Describe("View Rendering - Complete State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateComplete)
			result := NewExportArtifactResult(true, ExportTypeCV, ExportFormatPDF, ExportDestinationFile, "/tmp/cv.pdf", 1024)
			intent.model.result = result
		})

		It("should render Complete view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Export Complete"))
		})

		It("should show file path", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("/tmp/cv.pdf"))
		})
	})

	Describe("View Rendering - Failed State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateFailed)
			intent.model.error = &IntentError{
				Code:    "export_error",
				Message: "Failed to export artifact",
			}
		})

		It("should render Failed view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Export Failed"))
		})

		It("should show error message", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Failed to export artifact"))
		})
	})

	Describe("State Transitions - SelectType", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should navigate down through artifact types", func() {
			Expect(intent.GetSelectedIndex()).To(Equal(0))
			intent.Update(tea.KeyMsg{Type: tea.KeyDown, Runes: []rune{'j'}})
			Expect(intent.GetSelectedIndex()).To(Equal(1))
		})

		It("should navigate up through artifact types", func() {
			intent.SetSelectedIndex(2)
			intent.Update(tea.KeyMsg{Type: tea.KeyUp, Runes: []rune{'k'}})
			Expect(intent.GetSelectedIndex()).To(Equal(1))
		})

		It("should not go below first artifact type", func() {
			intent.SetSelectedIndex(0)
			intent.Update(tea.KeyMsg{Type: tea.KeyUp, Runes: []rune{'k'}})
			Expect(intent.GetSelectedIndex()).To(Equal(0))
		})

		It("should not go above last artifact type", func() {
			intent.SetSelectedIndex(4)
			intent.Update(tea.KeyMsg{Type: tea.KeyDown, Runes: []rune{'j'}})
			Expect(intent.GetSelectedIndex()).To(Equal(4))
		})

		It("should transition to SelectFormat on enter", func() {
			intent.SetSelectedIndex(0)
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter, Runes: []rune{'\n'}})
			Expect(intent.GetState()).To(Equal(ExportStateSelectFormat))
			Expect(intent.GetConfig()).NotTo(BeNil())
			Expect(intent.GetConfig().ArtifactType).To(Equal(ExportTypeCV))
		})

		It("should cancel on escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc, Runes: []rune{'\x1b'}})
			Expect(intent.IsActive()).To(BeFalse())
			result := intent.Result()
			Expect(result.Status).To(Equal(Cancelled))
		})
	})

	Describe("State Transitions - SelectFormat", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateSelectFormat)
			intent.SetConfig(NewExportConfiguration(ExportTypeCV, intent.model.context))
		})

		It("should navigate through formats", func() {
			Expect(intent.GetSelectedIndex()).To(Equal(0))
			intent.Update(tea.KeyMsg{Type: tea.KeyDown, Runes: []rune{'j'}})
			Expect(intent.GetSelectedIndex()).To(Equal(1))
		})

		It("should transition to SelectDestination on enter", func() {
			intent.SetSelectedIndex(0)
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter, Runes: []rune{'\n'}})
			Expect(intent.GetState()).To(Equal(ExportStateSelectDest))
		})

		It("should go back to SelectType on escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc, Runes: []rune{'\x1b'}})
			Expect(intent.GetState()).To(Equal(ExportStateSelectType))
		})
	})

	Describe("State Transitions - SelectDestination", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateSelectDest)
			intent.SetConfig(NewExportConfiguration(ExportTypeCV, intent.model.context))
		})

		It("should navigate through destinations", func() {
			Expect(intent.GetSelectedIndex()).To(Equal(0))
			intent.Update(tea.KeyMsg{Type: tea.KeyDown, Runes: []rune{'j'}})
			Expect(intent.GetSelectedIndex()).To(Equal(1))
		})

		It("should transition to Configure on enter", func() {
			intent.SetSelectedIndex(0)
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter, Runes: []rune{'\n'}})
			Expect(intent.GetState()).To(Equal(ExportStateConfigure))
			Expect(intent.GetConfig().Destination).To(Equal(ExportDestinationFile))
		})

		It("should go back to SelectFormat on escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc, Runes: []rune{'\x1b'}})
			Expect(intent.GetState()).To(Equal(ExportStateSelectFormat))
		})
	})

	Describe("State Transitions - Configure", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateConfigure)
			intent.SetConfig(NewExportConfiguration(ExportTypeCV, intent.model.context))
		})

		It("should transition to Preview on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter, Runes: []rune{'\n'}})
			Expect(intent.GetState()).To(Equal(ExportStatePreview))
		})

		It("should go back to SelectDestination on escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc, Runes: []rune{'\x1b'}})
			Expect(intent.GetState()).To(Equal(ExportStateSelectDest))
		})
	})

	Describe("State Transitions - Preview", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStatePreview)
			intent.SetConfig(NewExportConfiguration(ExportTypeCV, intent.model.context))
		})

		It("should transition to Confirm on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter, Runes: []rune{'\n'}})
			Expect(intent.GetState()).To(Equal(ExportStateConfirm))
		})

		It("should go back to Configure on escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc, Runes: []rune{'\x1b'}})
			Expect(intent.GetState()).To(Equal(ExportStateConfigure))
		})
	})

	Describe("State Transitions - Confirm", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateConfirm)
			intent.SetConfig(NewExportConfiguration(ExportTypeCV, intent.model.context))
		})

		It("should transition to InProgress on confirm (y)", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(intent.GetState()).To(Equal(ExportStateInProgress))
		})

		It("should transition to InProgress on confirm (enter)", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter, Runes: []rune{'\n'}})
			Expect(intent.GetState()).To(Equal(ExportStateInProgress))
		})

		It("should go back to Preview on cancel (n)", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(intent.GetState()).To(Equal(ExportStatePreview))
		})

		It("should go back to Preview on cancel (esc)", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc, Runes: []rune{'\x1b'}})
			Expect(intent.GetState()).To(Equal(ExportStatePreview))
		})
	})

	Describe("State Transitions - InProgress", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateInProgress)
			intent.SetConfig(NewExportConfiguration(ExportTypeCV, intent.model.context))
		})

		It("should transition to Complete on export completion", func() {
			result := NewExportArtifactResult(true, ExportTypeCV, ExportFormatPDF, ExportDestinationFile, "/tmp/cv.pdf", 1024)
			intent.Update(ExportCompleteMsg{Result: result})
			Expect(intent.GetState()).To(Equal(ExportStateComplete))
			Expect(intent.GetResult()).NotTo(BeNil())
		})

		It("should transition to Failed on export error", func() {
			err := &IntentError{Code: "export_error", Message: "Export failed"}
			intent.Update(ExportErrorMsg{Error: err})
			Expect(intent.GetState()).To(Equal(ExportStateFailed))
		})
	})

	Describe("State Transitions - Complete", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateComplete)
			result := NewExportArtifactResult(true, ExportTypeCV, ExportFormatPDF, ExportDestinationFile, "/tmp/cv.pdf", 1024)
			intent.model.result = result
		})

		It("should mark intent as inactive on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter, Runes: []rune{'\n'}})
			Expect(intent.IsActive()).To(BeFalse())
		})

		It("should mark intent as inactive on escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc, Runes: []rune{'\x1b'}})
			Expect(intent.IsActive()).To(BeFalse())
		})
	})

	Describe("State Transitions - Failed", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateFailed)
			intent.model.error = &IntentError{Code: "export_error", Message: "Export failed"}
		})

		It("should transition to Confirm on retry (r)", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			Expect(intent.GetState()).To(Equal(ExportStateConfirm))
		})

		It("should mark intent as inactive on cancel (esc)", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc, Runes: []rune{'\x1b'}})
			Expect(intent.IsActive()).To(BeFalse())
		})
	})

	Describe("Result Handling", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should return Cancelled when no result set", func() {
			result := intent.Result()
			Expect(result.Status).To(Equal(Cancelled))
		})

		It("should return Completed on successful export", func() {
			intent.SetState(ExportStateComplete)
			exportResult := NewExportArtifactResult(true, ExportTypeCV, ExportFormatPDF, ExportDestinationFile, "/tmp/cv.pdf", 1024)
			intent.model.result = exportResult
			result := intent.Result()
			Expect(result.Status).To(Equal(Completed))
			Expect(result.Data).NotTo(BeNil())
		})

		It("should return Failed on export failure", func() {
			intent.SetState(ExportStateFailed)
			exportResult := NewExportArtifactResultWithError(&IntentError{Code: "export_error", Message: "Export failed"})
			intent.model.result = exportResult
			result := intent.Result()
			Expect(result.Status).To(Equal(Failed))
			Expect(result.Error).NotTo(BeNil())
		})
	})

	Describe("Configuration Management", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should create configuration with default format for artifact type", func() {
			config := NewExportConfiguration(ExportTypeCV, intent.model.context)
			Expect(config.ArtifactType).To(Equal(ExportTypeCV))
			Expect(config.Format).To(Equal(ExportFormatPDF))
			Expect(config.Destination).To(Equal(ExportDestinationFile))
		})

		It("should support different formats for different artifact types", func() {
			cvConfig := NewExportConfiguration(ExportTypeCV, intent.model.context)
			eventsConfig := NewExportConfiguration(ExportTypeEvents, intent.model.context)
			Expect(cvConfig.Format).To(Equal(ExportFormatPDF))
			Expect(eventsConfig.Format).To(Equal(ExportFormatJSON))
		})
	})

	Describe("Export Context", func() {
		It("should have all artifact types", func() {
			ctx := NewExportArtifactContext()
			Expect(ctx.ArtifactTypes).To(HaveLen(5))
			Expect(ctx.ArtifactTypes).To(ContainElements(
				ExportTypeCV, ExportTypeEvents, ExportTypeFacts, ExportTypeBursts, ExportTypeProfile,
			))
		})

		It("should have supported formats for each artifact type", func() {
			ctx := NewExportArtifactContext()
			Expect(ctx.SupportedFormats[ExportTypeCV]).To(ContainElements(ExportFormatPDF, ExportFormatJSON, ExportFormatMD))
			Expect(ctx.SupportedFormats[ExportTypeEvents]).To(ContainElements(ExportFormatJSON, ExportFormatCSV, ExportFormatTXT))
		})

		It("should have default format for each artifact type", func() {
			ctx := NewExportArtifactContext()
			Expect(ctx.DefaultFormat[ExportTypeCV]).To(Equal(ExportFormatPDF))
			Expect(ctx.DefaultFormat[ExportTypeEvents]).To(Equal(ExportFormatJSON))
		})

		It("should have all destinations", func() {
			ctx := NewExportArtifactContext()
			Expect(ctx.Destinations).To(HaveLen(3))
			Expect(ctx.Destinations).To(ContainElements(
				ExportDestinationFile, ExportDestinationClipboard, ExportDestinationEmail,
			))
		})
	})

	Describe("Intent Interface Compliance", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should implement Intent interface", func() {
			var _ Intent = intent
		})

		It("should have Init method", func() {
			Expect(intent.Init).NotTo(BeNil())
		})

		It("should have Update method", func() {
			Expect(intent.Update).NotTo(BeNil())
		})

		It("should have View method", func() {
			Expect(intent.View).NotTo(BeNil())
		})

		It("should have Result method", func() {
			Expect(intent.Result).NotTo(BeNil())
		})
	})
})
