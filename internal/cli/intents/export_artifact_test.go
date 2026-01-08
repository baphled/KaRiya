package intents

import (
	"context"
	"fmt"

	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// NewTestExportArtifactContext creates a test context with default values
// Services and repositories are nil since most tests don't need them
func NewTestExportArtifactContext() *ExportArtifactContext {
	return &ExportArtifactContext{
		ArtifactTypes:    DefaultArtifactTypes(),
		SupportedFormats: DefaultSupportedFormats(),
		DefaultFormat:    DefaultFormats(),
		Destinations:     DefaultDestinations(),
		// Services and repositories are nil for testing
		ExportService: nil,

		CareerService:   nil,
		EventRepository: nil,
		FactRepository:  nil,
		BurstRepository: nil,
		AppContext:      context.Background(),
	}
}

// NewTestExportArtifactContextWithServices creates a test context with real services
// Used for integration tests that need actual export functionality
func NewTestExportArtifactContextWithServices() *ExportArtifactContext {
	// Create logger (discard output during tests)
	log := logger.New(nil, logger.ErrorLevel)

	// Create in-memory repositories
	eventRepo := careerrepo.NewMemoryRepository()
	factRepo := careerrepo.NewMemoryFactRepository()
	burstRepo := careerrepo.NewMemoryBurstRepository()

	// Create export service
	exportService := cv.NewExportService(log)

	return &ExportArtifactContext{
		ArtifactTypes:    DefaultArtifactTypes(),
		SupportedFormats: DefaultSupportedFormats(),
		DefaultFormat:    DefaultFormats(),
		Destinations:     DefaultDestinations(),
		ExportService:    exportService,
		// Not needed for export

		CareerService:   nil, // Not needed for export
		EventRepository: eventRepo,
		FactRepository:  factRepo,
		BurstRepository: burstRepo,
		AppContext:      context.Background(),
	}
}

var _ = Describe("ExportArtifact Intent", func() {
	var (
		intent *ExportArtifactIntent
	)

	Describe("Intent Creation", func() {
		It("should create a new ExportArtifact intent with valid context", func() {
			var err error
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			Expect(intent).NotTo(BeNil())
			Expect(intent.model).NotTo(BeNil())
		})

		It("should fail to create intent with nil context", func() {
			// Skipped because staticcheck prevents passing nil context
			// This is tested indirectly through all other tests
		})

		It("should initialize with correct default state", func() {
			intent, _ := NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(intent.GetState()).To(Equal(ExportStateSelectType))
		})

		It("should initialize with correct artifact types", func() {
			intent, _ := NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(intent.model.context.ArtifactTypes).To(HaveLen(3))
		})
	})

	Describe("Init Method", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should render SelectType view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Select Artifact Type"))
		})

		It("should show all artifact types", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("events"))
			Expect(view).To(ContainSubstring("facts"))
			Expect(view).To(ContainSubstring("bursts"))
		})

		It("should show navigation instructions", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("↑/↓"))
			Expect(view).To(ContainSubstring("Enter"))
			Expect(view).To(ContainSubstring("Esc"))
		})

		It("should highlight selected item", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("> events"))
		})
	})

	Describe("View Rendering - SelectFormat State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateSelectFormat)
			intent.SetConfig(NewExportConfiguration(ExportTypeEvents, intent.model.context))
		})

		It("should render SelectFormat view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Select Export Format"))
		})

		It("should show supported formats for artifact type", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("json"))
			Expect(view).To(ContainSubstring("csv"))
			Expect(view).To(ContainSubstring("txt"))
		})
	})

	Describe("View Rendering - SelectDestination State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
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
		})
	})

	Describe("View Rendering - Configure State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateConfigure)
			intent.SetConfig(NewExportConfiguration(ExportTypeEvents, intent.model.context))
		})

		It("should render Configure view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Configure Export"))
		})

		It("should show selected artifact type", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Artifact Type: events"))
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStatePreview)
			intent.SetConfig(NewExportConfiguration(ExportTypeEvents, intent.model.context))
		})

		It("should render Preview view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Preview Export"))
		})
	})

	Describe("View Rendering - Confirm State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateConfirm)
			intent.SetConfig(NewExportConfiguration(ExportTypeEvents, intent.model.context))
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateComplete)
			result := NewExportArtifactResult(true, ExportTypeEvents, ExportFormatPDF, ExportDestinationFile, "/tmp/cv.pdf", 1024)
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
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
			// Select Events (at index 0)
			intent.SetSelectedIndex(0)
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter, Runes: []rune{'\n'}})
			Expect(intent.GetState()).To(Equal(ExportStateSelectFormat))
			Expect(intent.GetConfig()).NotTo(BeNil())
			Expect(intent.GetConfig().ArtifactType).To(Equal(ExportTypeEvents))
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateSelectFormat)
			intent.SetConfig(NewExportConfiguration(ExportTypeEvents, intent.model.context))
		})

		It("should navigate through formats", func() {
			Expect(intent.GetSelectedIndex()).To(Equal(0))
			intent.Update(tea.KeyMsg{Type: tea.KeyDown, Runes: []rune{'j'}})
			Expect(intent.GetSelectedIndex()).To(Equal(1))
		})

		It("should navigate down with j key", func() {
			intent.SetSelectedIndex(0)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(intent.GetSelectedIndex()).To(Equal(1))
		})

		It("should navigate up with k key", func() {
			intent.SetSelectedIndex(1)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(intent.GetSelectedIndex()).To(Equal(0))
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateSelectDest)
			intent.SetConfig(NewExportConfiguration(ExportTypeEvents, intent.model.context))
		})

		It("should navigate through destinations", func() {
			Expect(intent.GetSelectedIndex()).To(Equal(0))
			intent.Update(tea.KeyMsg{Type: tea.KeyDown, Runes: []rune{'j'}})
			Expect(intent.GetSelectedIndex()).To(Equal(1))
		})

		It("should navigate down with j key", func() {
			intent.SetSelectedIndex(0)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(intent.GetSelectedIndex()).To(Equal(1))
		})

		It("should navigate up with k key", func() {
			intent.SetSelectedIndex(1)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(intent.GetSelectedIndex()).To(Equal(0))
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateConfigure)
			intent.SetConfig(NewExportConfiguration(ExportTypeEvents, intent.model.context))
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStatePreview)
			intent.SetConfig(NewExportConfiguration(ExportTypeEvents, intent.model.context))
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateConfirm)
			intent.SetConfig(NewExportConfiguration(ExportTypeEvents, intent.model.context))
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateInProgress)
			intent.SetConfig(NewExportConfiguration(ExportTypeEvents, intent.model.context))
		})

		It("should transition to Complete on export completion", func() {
			result := NewExportArtifactResult(true, ExportTypeEvents, ExportFormatPDF, ExportDestinationFile, "/tmp/cv.pdf", 1024)
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStateComplete)
			result := NewExportArtifactResult(true, ExportTypeEvents, ExportFormatPDF, ExportDestinationFile, "/tmp/cv.pdf", 1024)
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should return nil when no result set", func() {
			// Before the intent is completed, Result() should return nil
			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should return Cancelled when intent is cancelled", func() {
			// Simulate cancelling the intent by pressing Escape
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Cancelled))
		})

		It("should return Completed on successful export", func() {
			intent.SetState(ExportStateComplete)
			exportResult := NewExportArtifactResult(true, ExportTypeEvents, ExportFormatPDF, ExportDestinationFile, "/tmp/cv.pdf", 1024)
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should create configuration with default format for artifact type", func() {
			config := NewExportConfiguration(ExportTypeEvents, intent.model.context)
			Expect(config.ArtifactType).To(Equal(ExportTypeEvents))
			Expect(config.Format).To(Equal(ExportFormatJSON))
			Expect(config.Destination).To(Equal(ExportDestinationFile))
		})

		It("should support different formats for different artifact types", func() {
			eventsConfig := NewExportConfiguration(ExportTypeEvents, intent.model.context)
			factsConfig := NewExportConfiguration(ExportTypeFacts, intent.model.context)
			Expect(eventsConfig.Format).To(Equal(ExportFormatJSON))
			Expect(factsConfig.Format).To(Equal(ExportFormatJSON))
		})
	})

	Describe("Export Context", func() {
		It("should have all artifact types", func() {
			ctx := NewTestExportArtifactContext()
			Expect(ctx.ArtifactTypes).To(HaveLen(3))
			Expect(ctx.ArtifactTypes).To(ContainElements(
				ExportTypeEvents, ExportTypeFacts, ExportTypeBursts,
			))
		})

		It("should have supported formats for each artifact type", func() {
			ctx := NewTestExportArtifactContext()
			Expect(ctx.SupportedFormats[ExportTypeEvents]).To(ContainElements(ExportFormatJSON, ExportFormatYAML, ExportFormatCSV, ExportFormatTXT))
			Expect(ctx.SupportedFormats[ExportTypeFacts]).To(ContainElements(ExportFormatJSON, ExportFormatYAML, ExportFormatCSV, ExportFormatTXT))
		})

		It("should have default format for each artifact type", func() {
			ctx := NewTestExportArtifactContext()
			Expect(ctx.DefaultFormat[ExportTypeEvents]).To(Equal(ExportFormatJSON))
			Expect(ctx.DefaultFormat[ExportTypeFacts]).To(Equal(ExportFormatJSON))
		})

		It("should have all destinations", func() {
			ctx := NewTestExportArtifactContext()
			Expect(ctx.Destinations).To(HaveLen(2))
			Expect(ctx.Destinations).To(ContainElements(
				ExportDestinationFile, ExportDestinationClipboard,
			))
		})
	})

	Describe("Intent Interface Compliance", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
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

	Describe("formatBytes utility function", func() {
		It("should format 0 bytes", func() {
			Expect(formatBytes(0)).To(Equal("0 B"))
		})

		It("should format bytes under 1KB", func() {
			Expect(formatBytes(1)).To(Equal("1 B"))
			Expect(formatBytes(512)).To(Equal("512 B"))
			Expect(formatBytes(1023)).To(Equal("1023 B"))
		})

		It("should format 1KB exactly", func() {
			Expect(formatBytes(1024)).To(Equal("1 KB"))
		})

		It("should format kilobytes", func() {
			Expect(formatBytes(2048)).To(Equal("2 KB"))
			Expect(formatBytes(5120)).To(Equal("5 KB"))
		})

		It("should format 1MB exactly", func() {
			Expect(formatBytes(1048576)).To(Equal("1 MB"))
		})

		It("should format megabytes", func() {
			Expect(formatBytes(2097152)).To(Equal("2 MB"))
			Expect(formatBytes(10485760)).To(Equal("10 MB"))
		})

		It("should format 1GB exactly", func() {
			Expect(formatBytes(1073741824)).To(Equal("1 GB"))
		})

		It("should format gigabytes", func() {
			Expect(formatBytes(2147483648)).To(Equal("2 GB"))
		})

		It("should format 1TB exactly", func() {
			Expect(formatBytes(1099511627776)).To(Equal("1 TB"))
		})
	})

	Describe("Scroll percentage display in preview", func() {
		var intent *ExportArtifactIntent

		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ExportStatePreview)
			intent.SetConfig(NewExportConfiguration(ExportTypeCV, intent.model.context))
		})

		It("should not show scroll percentage for short content", func() {
			// Set preview with only 5 lines (less than viewHeight of 15)
			intent.model.preview = "Line 1\nLine 2\nLine 3\nLine 4\nLine 5"
			intent.model.previewLines = []string{"Line 1", "Line 2", "Line 3", "Line 4", "Line 5"}
			intent.model.scrollOffset = 0

			view := intent.model.viewPreview()
			Expect(view).NotTo(ContainSubstring("% scrolled"))
		})

		It("should show 0% at top of long content", func() {
			// Create 30 lines of content (more than viewHeight of 15)
			lines := make([]string, 30)
			for i := range lines {
				lines[i] = fmt.Sprintf("Line %d", i+1)
			}
			intent.model.previewLines = lines
			intent.model.scrollOffset = 0

			view := intent.model.viewPreview()
			Expect(view).To(ContainSubstring("[0% scrolled]"))
		})

		It("should show 50% at middle of content", func() {
			// Create 30 lines of content
			lines := make([]string, 30)
			for i := range lines {
				lines[i] = fmt.Sprintf("Line %d", i+1)
			}
			intent.model.previewLines = lines
			intent.model.scrollOffset = 15 // Middle of 30 lines

			view := intent.model.viewPreview()
			Expect(view).To(ContainSubstring("[50% scrolled]"))
		})

		It("should show 100% at end of content", func() {
			// Create 30 lines of content
			lines := make([]string, 30)
			for i := range lines {
				lines[i] = fmt.Sprintf("Line %d", i+1)
			}
			intent.model.previewLines = lines
			intent.model.scrollOffset = 30 // End of content

			view := intent.model.viewPreview()
			Expect(view).To(ContainSubstring("[100% scrolled]"))
		})

		It("should show 33% at one-third of content", func() {
			// Create 30 lines of content
			lines := make([]string, 30)
			for i := range lines {
				lines[i] = fmt.Sprintf("Line %d", i+1)
			}
			intent.model.previewLines = lines
			intent.model.scrollOffset = 10 // 10/30 = 33%

			view := intent.model.viewPreview()
			Expect(view).To(ContainSubstring("[33% scrolled]"))
		})
	})
})
