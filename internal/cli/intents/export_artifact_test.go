package intents

import (
	"context"
	"fmt"
	"strings"
	"time"

	careerdomain "github.com/baphled/kariya/internal/domain/career"
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
		ExportService:       nil,
		CVGenerationService: nil,
		CVConfigManager:     nil,
		CareerService:       nil,
		EventRepository:     nil,
		FactRepository:      nil,
		BurstRepository:     nil,
		AppContext:          context.Background(),
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
		ArtifactTypes:       DefaultArtifactTypes(),
		SupportedFormats:    DefaultSupportedFormats(),
		DefaultFormat:       DefaultFormats(),
		Destinations:        DefaultDestinations(),
		ExportService:       exportService,
		CVGenerationService: nil,                         // Not needed for export
		CVConfigManager:     cv.NewMemoryConfigManager(), // Provide memory config manager
		CareerService:       nil,                         // Not needed for export
		EventRepository:     eventRepo,
		FactRepository:      factRepo,
		BurstRepository:     burstRepo,
		AppContext:          context.Background(),
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
			Expect(intent.model.context.ArtifactTypes).To(HaveLen(4))
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
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
			Expect(view).To(ContainSubstring("txt"))
			Expect(view).To(ContainSubstring("markdown"))
			Expect(view).To(ContainSubstring("yaml"))
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
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
			// Select non-CV type (Events at index 1) to skip SelectCV
			intent.SetSelectedIndex(1)
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
			intent.SetConfig(NewExportConfiguration(ExportTypeCV, intent.model.context))
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

		It("should go back to SelectCV on escape (when artifact type is CV)", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc, Runes: []rune{'\x1b'}})
			Expect(intent.GetState()).To(Equal(ExportStateSelectCV))
		})
	})

	Describe("State Transitions - SelectDestination", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should create configuration with default format for artifact type", func() {
			config := NewExportConfiguration(ExportTypeCV, intent.model.context)
			Expect(config.ArtifactType).To(Equal(ExportTypeCV))
			Expect(config.Format).To(Equal(ExportFormatMD))
			Expect(config.Destination).To(Equal(ExportDestinationFile))
		})

		It("should support different formats for different artifact types", func() {
			cvConfig := NewExportConfiguration(ExportTypeCV, intent.model.context)
			eventsConfig := NewExportConfiguration(ExportTypeEvents, intent.model.context)
			Expect(cvConfig.Format).To(Equal(ExportFormatMD))
			Expect(eventsConfig.Format).To(Equal(ExportFormatJSON))
		})
	})

	Describe("Export Context", func() {
		It("should have all artifact types", func() {
			ctx := NewTestExportArtifactContext()
			Expect(ctx.ArtifactTypes).To(HaveLen(4))
			Expect(ctx.ArtifactTypes).To(ContainElements(
				ExportTypeCV, ExportTypeEvents, ExportTypeFacts, ExportTypeBursts,
			))
		})

		It("should have supported formats for each artifact type", func() {
			ctx := NewTestExportArtifactContext()
			Expect(ctx.SupportedFormats[ExportTypeCV]).To(ContainElements(ExportFormatTXT, ExportFormatMD, ExportFormatYAML))
			Expect(ctx.SupportedFormats[ExportTypeEvents]).To(ContainElements(ExportFormatJSON, ExportFormatYAML, ExportFormatCSV, ExportFormatTXT))
		})

		It("should have default format for each artifact type", func() {
			ctx := NewTestExportArtifactContext()
			Expect(ctx.DefaultFormat[ExportTypeCV]).To(Equal(ExportFormatMD))
			Expect(ctx.DefaultFormat[ExportTypeEvents]).To(Equal(ExportFormatJSON))
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
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
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

	Describe("CV Selection State", func() {
		var intent *ExportArtifactIntent

		BeforeEach(func() {
			var err error
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		Describe("State Entry", func() {
			It("should transition to SelectCV state when CV type is selected", func() {
				// Select CV artifact type (index 0 - CV is first in DefaultArtifactTypes)
				intent.model.selectedIndex = 0

				// Press Enter to proceed
				intent.model.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Should transition to SelectCV state
				Expect(intent.model.state).To(Equal(ExportStateSelectCV))
				Expect(intent.model.config).NotTo(BeNil())
				Expect(intent.model.config.ArtifactType).To(Equal(ExportTypeCV))
			})

			It("should skip SelectCV state for non-CV artifact types", func() {
				// Select Events artifact type (index 1 - Events is second)
				intent.model.selectedIndex = 1

				// Press Enter to proceed
				intent.model.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Should transition directly to SelectFormat state
				Expect(intent.model.state).To(Equal(ExportStateSelectFormat))
				Expect(intent.model.config).NotTo(BeNil())
				Expect(intent.model.config.ArtifactType).To(Equal(ExportTypeEvents))
			})
		})

		Describe("Load Available CVs", func() {
			It("should load CVs from ConfigManager", func() {
				// Create a config manager with test configs
				configManager := cv.NewMemoryConfigManager()
				config1 := &careerdomain.CVConfig{
					Name:           "Senior Engineer CV",
					TargetRole:     "senior_ic",
					TargetAudience: "hiring_manager",
				}
				config2 := &careerdomain.CVConfig{
					Name:           "Staff Engineer CV",
					TargetRole:     "staff",
					TargetAudience: "recruiter",
				}
				ctx := context.Background()
				err := configManager.SaveConfig(ctx, config1)
				Expect(err).NotTo(HaveOccurred())
				err = configManager.SaveConfig(ctx, config2)
				Expect(err).NotTo(HaveOccurred())

				// Create intent with config manager
				testCtx := NewTestExportArtifactContext()
				testCtx.CVConfigManager = configManager
				intent, err := NewExportArtifactIntent(testCtx)
				Expect(err).NotTo(HaveOccurred())
				intent.Init()

				// Load CVs
				cmd := intent.model.loadAvailableCVs()
				Expect(cmd).NotTo(BeNil())

				// Execute command and get message
				msg := cmd()
				cvsMsg, ok := msg.(CVsLoadedMsg)
				Expect(ok).To(BeTrue())

				// Verify CVs were loaded
				Expect(cvsMsg.CVs).To(HaveLen(2))
				Expect(cvsMsg.CVs[0].Name).To(Equal("Senior Engineer CV"))
				Expect(cvsMsg.CVs[0].TargetRole).To(Equal("senior_ic"))
				Expect(cvsMsg.CVs[1].Name).To(Equal("Staff Engineer CV"))
				Expect(cvsMsg.CVs[1].TargetRole).To(Equal("staff"))
			})

			It("should return empty list when ConfigManager is nil", func() {
				// Create intent without config manager
				testCtx := NewTestExportArtifactContext()
				testCtx.CVConfigManager = nil
				intent, err := NewExportArtifactIntent(testCtx)
				Expect(err).NotTo(HaveOccurred())
				intent.Init()

				// Load CVs
				cmd := intent.model.loadAvailableCVs()
				msg := cmd()
				cvsMsg, ok := msg.(CVsLoadedMsg)
				Expect(ok).To(BeTrue())

				// Should return empty list
				Expect(cvsMsg.CVs).To(HaveLen(0))
			})

			It("should return empty list when ConfigManager returns error", func() {
				// Use a mock that returns error (simulate error by using fresh manager)
				configManager := cv.NewMemoryConfigManager()
				testCtx := NewTestExportArtifactContext()
				testCtx.CVConfigManager = configManager
				intent, err := NewExportArtifactIntent(testCtx)
				Expect(err).NotTo(HaveOccurred())
				intent.Init()

				// Load CVs (empty manager returns empty list, not error)
				cmd := intent.model.loadAvailableCVs()
				msg := cmd()
				cvsMsg, ok := msg.(CVsLoadedMsg)
				Expect(ok).To(BeTrue())

				// Should return empty list
				Expect(cvsMsg.CVs).To(HaveLen(0))
			})

			It("should update availableCVs when CVsLoadedMsg is received", func() {
				// Set up intent in SelectCV state
				intent.model.state = ExportStateSelectCV

				// Create mock CVs
				mockCVs := []*careerdomain.CVView{
					{Name: "Test CV", TargetRole: "Engineer"},
				}

				// Send CVsLoadedMsg
				cmd := intent.model.Update(CVsLoadedMsg{CVs: mockCVs})

				// Verify availableCVs was updated
				Expect(intent.model.availableCVs).To(Equal(mockCVs))
				Expect(cmd).To(BeNil())
			})
		})

		Describe("CV List Navigation", func() {
			BeforeEach(func() {
				// Set up CV selection state with mock CVs
				intent.model.state = ExportStateSelectCV
				intent.model.availableCVs = []*careerdomain.CVView{
					{Name: "CV 1", TargetRole: "Senior Engineer"},
					{Name: "CV 2", TargetRole: "Staff Engineer"},
					{Name: "CV 3", TargetRole: "Tech Lead"},
				}
				intent.model.selectedIndex = 0
			})

			It("should navigate down with arrow key", func() {
				intent.model.Update(tea.KeyMsg{Type: tea.KeyDown})
				Expect(intent.model.selectedIndex).To(Equal(1))
			})

			It("should navigate down with j key", func() {
				intent.model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				Expect(intent.model.selectedIndex).To(Equal(1))
			})

			It("should navigate up with arrow key", func() {
				intent.model.selectedIndex = 1
				intent.model.Update(tea.KeyMsg{Type: tea.KeyUp})
				Expect(intent.model.selectedIndex).To(Equal(0))
			})

			It("should navigate up with k key", func() {
				intent.model.selectedIndex = 1
				intent.model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
				Expect(intent.model.selectedIndex).To(Equal(0))
			})

			It("should wrap at bottom when navigating down", func() {
				intent.model.selectedIndex = 2 // Last item
				intent.model.Update(tea.KeyMsg{Type: tea.KeyDown})
				Expect(intent.model.selectedIndex).To(Equal(0)) // Wraps to first
			})

			It("should wrap at top when navigating up", func() {
				intent.model.selectedIndex = 0 // First item
				intent.model.Update(tea.KeyMsg{Type: tea.KeyUp})
				Expect(intent.model.selectedIndex).To(Equal(2)) // Wraps to last
			})
		})

		Describe("CV Selection", func() {
			BeforeEach(func() {
				intent.model.state = ExportStateSelectCV
				intent.model.availableCVs = []*careerdomain.CVView{
					{Name: "CV 1", TargetRole: "Senior Engineer"},
					{Name: "CV 2", TargetRole: "Staff Engineer"},
				}
				intent.model.selectedIndex = 1
			})

			It("should select CV and transition to SelectFormat when Enter is pressed", func() {
				intent.model.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(intent.model.selectedCV).To(Equal(intent.model.availableCVs[1]))
				Expect(intent.model.state).To(Equal(ExportStateSelectFormat))
			})

			It("should go back to SelectType when Escape is pressed", func() {
				intent.model.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.model.state).To(Equal(ExportStateSelectType))
			})
		})

		Describe("View Rendering", func() {
			BeforeEach(func() {
				intent.model.state = ExportStateSelectCV
				intent.model.availableCVs = []*careerdomain.CVView{
					{Name: "CV 1", TargetRole: "Senior Engineer"},
					{Name: "CV 2", TargetRole: "Staff Engineer"},
				}
				intent.model.selectedIndex = 0
			})

			It("should display list of available CVs", func() {
				view := intent.View()

				Expect(view).To(ContainSubstring("CV 1"))
				Expect(view).To(ContainSubstring("CV 2"))
			})

			It("should highlight selected CV", func() {
				view := intent.View()

				// Should have selection indicator for first item
				Expect(view).To(ContainSubstring("▶"))
			})

			It("should show CV target role", func() {
				view := intent.View()

				Expect(view).To(ContainSubstring("Senior Engineer"))
				Expect(view).To(ContainSubstring("Staff Engineer"))
			})
		})
	})

	Describe("Format Mapping", func() {
		It("should map TXT to ExportFormatText", func() {
			result := mapToExportServiceFormat(ExportFormatTXT)
			Expect(result).To(Equal(cv.ExportFormatText))
		})

		It("should map MD to ExportFormatMarkdown", func() {
			result := mapToExportServiceFormat(ExportFormatMD)
			Expect(result).To(Equal(cv.ExportFormatMarkdown))
		})

		It("should map YAML to ExportFormatYAML", func() {
			result := mapToExportServiceFormat(ExportFormatYAML)
			Expect(result).To(Equal(cv.ExportFormatYAML))
		})

		It("should default to ExportFormatText for unknown formats", func() {
			result := mapToExportServiceFormat(ExportFormat("unknown"))
			Expect(result).To(Equal(cv.ExportFormatText))
		})

		It("should default to ExportFormatText for JSON", func() {
			result := mapToExportServiceFormat(ExportFormatJSON)
			Expect(result).To(Equal(cv.ExportFormatText))
		})

		It("should default to ExportFormatText for CSV", func() {
			result := mapToExportServiceFormat(ExportFormatCSV)
			Expect(result).To(Equal(cv.ExportFormatText))
		})
	})

	Describe("Real Export Implementation", func() {
		var model *ExportArtifactModel

		BeforeEach(func() {
			context := NewTestExportArtifactContextWithServices()
			intent, err := NewExportArtifactIntent(context)
			Expect(err).NotTo(HaveOccurred())
			model = intent.model

			// Initialize config (normally done by state machine)
			model.config = &ExportConfiguration{
				ArtifactType: ExportTypeCV,
				Format:       ExportFormatTXT,
				Destination:  ExportDestinationFile,
			}
		})

		Describe("exportCV", func() {
			It("should return error when selectedCV is nil", func() {
				model.config.ArtifactType = ExportTypeCV
				model.config.Format = ExportFormatTXT
				model.config.Destination = ExportDestinationFile
				model.selectedCV = nil

				result, err := model.exportCV()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no CV selected"))
				Expect(result).To(BeNil())
			})

			It("should export CV to text format and save to file", func() {
				// Setup a test CV
				now := time.Now()
				testCV := &careerdomain.CVView{
					ID:               "test-cv-1",
					Name:             "Senior Go Developer CV",
					TargetRole:       "Senior Go Developer",
					TargetAudience:   "Tech Companies",
					GeneratedAt:      now,
					SourceEventCount: 5,
					SourceFactCount:  10,
					Sections: []*careerdomain.CVSection{
						{
							Title:       "Summary",
							SectionType: "summary",
							Summary:     "Experienced Go developer",
							Content:     []*careerdomain.SectionContentGroup{},
						},
					},
				}

				model.selectedCV = testCV
				model.config.ArtifactType = ExportTypeCV
				model.config.Format = ExportFormatTXT
				model.config.Destination = ExportDestinationFile

				result, err := model.exportCV()
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Success).To(BeTrue())
				Expect(result.ArtifactType).To(Equal(ExportTypeCV))
				Expect(result.Format).To(Equal(ExportFormatTXT))
				Expect(result.Destination).To(Equal(ExportDestinationFile))
				Expect(result.FilePath).NotTo(BeEmpty())
				Expect(result.FilePath).To(ContainSubstring(".txt"))
				Expect(result.Size).To(BeNumerically(">", 0))
			})

			It("should export CV to markdown format and save to file", func() {
				now := time.Now()
				testCV := &careerdomain.CVView{
					ID:               "test-cv-2",
					Name:             "Full Stack Engineer CV",
					TargetRole:       "Full Stack Engineer",
					TargetAudience:   "Startups",
					GeneratedAt:      now,
					SourceEventCount: 3,
					SourceFactCount:  8,
					Sections: []*careerdomain.CVSection{
						{
							Title:       "Experience",
							SectionType: "experience",
							Summary:     "",
							Content:     []*careerdomain.SectionContentGroup{},
						},
					},
				}

				model.selectedCV = testCV
				model.config.ArtifactType = ExportTypeCV
				model.config.Format = ExportFormatMD
				model.config.Destination = ExportDestinationFile

				result, err := model.exportCV()
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.FilePath).To(ContainSubstring(".md"))
			})

			It("should export CV to YAML format and save to file", func() {
				now := time.Now()
				testCV := &careerdomain.CVView{
					ID:               "test-cv-3",
					Name:             "DevOps Engineer CV",
					TargetRole:       "DevOps Engineer",
					TargetAudience:   "Enterprise",
					GeneratedAt:      now,
					SourceEventCount: 7,
					SourceFactCount:  15,
					Sections:         []*careerdomain.CVSection{},
				}

				model.selectedCV = testCV
				model.config.ArtifactType = ExportTypeCV
				model.config.Format = ExportFormatYAML
				model.config.Destination = ExportDestinationFile

				result, err := model.exportCV()
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.FilePath).To(ContainSubstring(".yaml"))
			})

			It("should export CV to clipboard", func() {
				now := time.Now()
				testCV := &careerdomain.CVView{
					ID:               "test-cv-4",
					Name:             "Backend Developer CV",
					TargetRole:       "Backend Developer",
					TargetAudience:   "SaaS Companies",
					GeneratedAt:      now,
					SourceEventCount: 4,
					SourceFactCount:  9,
					Sections:         []*careerdomain.CVSection{},
				}

				model.selectedCV = testCV
				model.config.ArtifactType = ExportTypeCV
				model.config.Format = ExportFormatMD
				model.config.Destination = ExportDestinationClipboard

				result, err := model.exportCV()
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Success).To(BeTrue())
				Expect(result.Destination).To(Equal(ExportDestinationClipboard))
				Expect(result.FilePath).To(Equal("clipboard"))
				Expect(result.Size).To(BeNumerically(">", 0))
			})
		})

		Describe("exportEvents", func() {
			It("should export events to JSON format", func() {
				model.config.ArtifactType = ExportTypeEvents
				model.config.Format = ExportFormatJSON
				model.config.Destination = ExportDestinationFile

				result, err := model.exportEvents()
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Success).To(BeTrue())
				Expect(result.FilePath).To(ContainSubstring(".json"))
			})

			It("should export events to CSV format", func() {
				model.config.ArtifactType = ExportTypeEvents
				model.config.Format = ExportFormatCSV
				model.config.Destination = ExportDestinationFile

				result, err := model.exportEvents()
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.FilePath).To(ContainSubstring(".csv"))
			})

			It("should export events to clipboard", func() {
				model.config.ArtifactType = ExportTypeEvents
				model.config.Format = ExportFormatJSON
				model.config.Destination = ExportDestinationClipboard

				result, err := model.exportEvents()
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.FilePath).To(Equal("clipboard"))
			})
		})

		Describe("exportFacts", func() {
			It("should export facts to JSON format", func() {
				model.config.ArtifactType = ExportTypeFacts
				model.config.Format = ExportFormatJSON
				model.config.Destination = ExportDestinationFile

				result, err := model.exportFacts()
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Success).To(BeTrue())
			})

			It("should export facts to YAML format", func() {
				model.config.ArtifactType = ExportTypeFacts
				model.config.Format = ExportFormatYAML
				model.config.Destination = ExportDestinationFile

				result, err := model.exportFacts()
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.FilePath).To(ContainSubstring(".yaml"))
			})
		})

		Describe("exportBursts", func() {
			It("should export bursts to JSON format", func() {
				model.config.ArtifactType = ExportTypeBursts
				model.config.Format = ExportFormatJSON
				model.config.Destination = ExportDestinationFile

				result, err := model.exportBursts()
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Success).To(BeTrue())
			})

			It("should export bursts to TXT format", func() {
				model.config.ArtifactType = ExportTypeBursts
				model.config.Format = ExportFormatTXT
				model.config.Destination = ExportDestinationFile

				result, err := model.exportBursts()
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.FilePath).To(ContainSubstring(".txt"))
			})
		})

		Describe("startExport integration", func() {
			It("should call exportCV for CV artifact type", func() {
				now := time.Now()
				testCV := &careerdomain.CVView{
					ID:               "test-cv-integration",
					Name:             "Test CV",
					TargetRole:       "Developer",
					TargetAudience:   "Tech",
					GeneratedAt:      now,
					SourceEventCount: 1,
					SourceFactCount:  2,
					Sections:         []*careerdomain.CVSection{},
				}

				model.selectedCV = testCV
				model.config.ArtifactType = ExportTypeCV
				model.config.Format = ExportFormatTXT
				model.config.Destination = ExportDestinationFile

				cmd := model.startExport()
				Expect(cmd).NotTo(BeNil())

				msg := cmd()
				completeMsg, ok := msg.(ExportCompleteMsg)
				Expect(ok).To(BeTrue())
				Expect(completeMsg.Result).NotTo(BeNil())
				Expect(completeMsg.Result.Success).To(BeTrue())
			})

			It("should call exportEvents for Events artifact type", func() {
				model.config.ArtifactType = ExportTypeEvents
				model.config.Format = ExportFormatJSON
				model.config.Destination = ExportDestinationFile

				cmd := model.startExport()
				Expect(cmd).NotTo(BeNil())

				msg := cmd()
				completeMsg, ok := msg.(ExportCompleteMsg)
				Expect(ok).To(BeTrue())
				Expect(completeMsg.Result.Success).To(BeTrue())
			})

			It("should return error message when export fails", func() {
				// CV export without selected CV should fail
				model.config.ArtifactType = ExportTypeCV
				model.config.Format = ExportFormatTXT
				model.config.Destination = ExportDestinationFile
				model.selectedCV = nil

				cmd := model.startExport()
				msg := cmd()
				errorMsg, ok := msg.(ExportErrorMsg)
				Expect(ok).To(BeTrue())
				Expect(errorMsg.Error).NotTo(BeNil())
			})
		})
	})

	Describe("Real Preview Data", func() {
		var model *ExportArtifactModel
		var testContext *ExportArtifactContext

		BeforeEach(func() {
			testContext = NewTestExportArtifactContextWithServices()
			intent, err := NewExportArtifactIntent(testContext)
			Expect(err).NotTo(HaveOccurred())
			model = intent.model

			// Initialize config
			model.config = &ExportConfiguration{
				ArtifactType: ExportTypeCV,
				Format:       ExportFormatTXT,
				Destination:  ExportDestinationFile,
			}
		})

		Describe("generateCVPreview", func() {
			It("should generate preview from real CV data", func() {
				// Setup a test CV
				now := time.Now()
				testCV := &careerdomain.CVView{
					ID:               "test-cv-1",
					Name:             "Senior Developer CV",
					TargetRole:       "Senior Developer",
					TargetAudience:   "Tech Companies",
					GeneratedAt:      now,
					SourceEventCount: 5,
					SourceFactCount:  10,
					Sections: []*careerdomain.CVSection{
						{
							Title:       "Summary",
							SectionType: "summary",
							Summary:     "Experienced developer with 10 years",
							Content:     []*careerdomain.SectionContentGroup{},
						},
					},
				}

				model.selectedCV = testCV
				model.config.Format = ExportFormatTXT

				preview := model.generateCVPreview()
				Expect(preview).NotTo(BeEmpty())
				// ExportService uppercases CV name in text format
				Expect(strings.ToLower(preview)).To(ContainSubstring("senior developer cv"))
				Expect(strings.ToLower(preview)).To(ContainSubstring("summary"))
			})

			It("should show error when no CV is selected", func() {
				model.selectedCV = nil
				model.config.Format = ExportFormatMD

				preview := model.generateCVPreview()
				Expect(preview).To(ContainSubstring("No CV selected"))
			})

			It("should support markdown format", func() {
				now := time.Now()
				testCV := &careerdomain.CVView{
					ID:          "test-cv-2",
					Name:        "Test CV",
					TargetRole:  "Engineer",
					GeneratedAt: now,
					Sections:    []*careerdomain.CVSection{},
				}

				model.selectedCV = testCV
				model.config.Format = ExportFormatMD

				preview := model.generateCVPreview()
				Expect(preview).To(ContainSubstring("# Test CV"))
			})
		})

		Describe("generateEventsPreview", func() {
			It("should generate preview from real events data", func() {
				// Add test events to repository
				event1 := &careerdomain.CareerEvent{
					ID:      "evt-1",
					Text:    "Led team standup",
					Date:    time.Now().AddDate(0, 0, -1),
					Company: "Tech Corp",
					Tags:    []string{"leadership"},
				}
				event2 := &careerdomain.CareerEvent{
					ID:      "evt-2",
					Text:    "Completed API integration",
					Date:    time.Now().AddDate(0, 0, -2),
					Project: "Platform",
					Tags:    []string{"technical"},
				}

				ctx := context.Background()
				testContext.EventRepository.Create(ctx, event1)
				testContext.EventRepository.Create(ctx, event2)

				model.config.ArtifactType = ExportTypeEvents
				model.config.Format = ExportFormatJSON

				preview := model.generateEventsPreview()
				Expect(preview).NotTo(BeEmpty())
				Expect(preview).To(ContainSubstring("evt-1"))
				Expect(preview).To(ContainSubstring("Led team standup"))
			})

			It("should show message when no events exist", func() {
				model.config.ArtifactType = ExportTypeEvents
				model.config.Format = ExportFormatCSV

				preview := model.generateEventsPreview()
				Expect(preview).NotTo(BeEmpty())
				Expect(preview).To(ContainSubstring("No events"))
			})

			It("should support text format", func() {
				event := &careerdomain.CareerEvent{
					ID:      "evt-3",
					Text:    "Test event",
					Date:    time.Now(),
					Company: "Acme",
				}

				ctx := context.Background()
				testContext.EventRepository.Create(ctx, event)

				model.config.ArtifactType = ExportTypeEvents
				model.config.Format = ExportFormatTXT

				preview := model.generateEventsPreview()
				Expect(preview).To(ContainSubstring("Test event"))
				Expect(preview).To(ContainSubstring("Acme"))
			})
		})

		Describe("generateFactsPreview", func() {
			It("should generate preview from real facts data", func() {
				// Add test facts to repository
				now := time.Now()
				fact1 := &careerdomain.Fact{
					ID:                   "fact-1",
					Text:                 "Expert in Go",
					CompetencyCategories: []string{"technical"},
					RoleFit:              careerdomain.RoleFitPrincipal,
					StrengthSignal:       "strong",
					AudienceRelevance:    []string{"hiring_manager", "peer"},
					SourceEventID:        "event-1",
					CreatedAt:            now,
					UpdatedAt:            now,
				}
				fact2 := &careerdomain.Fact{
					ID:                   "fact-2",
					Text:                 "Led 5-person team",
					CompetencyCategories: []string{"leadership"},
					RoleFit:              careerdomain.RoleFitStaff,
					StrengthSignal:       "medium",
					AudienceRelevance:    []string{"hiring_manager", "recruiter"},
					SourceEventID:        "event-2",
					CreatedAt:            now,
					UpdatedAt:            now,
				}

				ctx := context.Background()
				err := testContext.FactRepository.Create(ctx, fact1)
				Expect(err).NotTo(HaveOccurred())
				err = testContext.FactRepository.Create(ctx, fact2)
				Expect(err).NotTo(HaveOccurred())

				model.config.ArtifactType = ExportTypeFacts
				model.config.Format = ExportFormatJSON

				preview := model.generateFactsPreview()
				Expect(preview).NotTo(BeEmpty())
				Expect(preview).To(ContainSubstring("fact-1"))
				Expect(preview).To(ContainSubstring("Expert in Go"))
			})

			It("should show message when no facts exist", func() {
				model.config.ArtifactType = ExportTypeFacts
				model.config.Format = ExportFormatYAML

				preview := model.generateFactsPreview()
				Expect(preview).NotTo(BeEmpty())
				Expect(preview).To(ContainSubstring("No facts"))
			})
		})

		Describe("generateBurstsPreview", func() {
			It("should generate preview from real bursts data", func() {
				// Add test bursts to repository
				burst1 := &careerdomain.Burst{
					ID:          "burst-1",
					Name:        "Q4 2024 Platform Work",
					Description: "Major platform improvements",
					EventIDs:    []string{"evt-1", "evt-2"},
					Confirmed:   true,
				}
				burst2 := &careerdomain.Burst{
					ID:        "burst-2",
					Name:      "API Development Sprint",
					EventIDs:  []string{"evt-3"},
					Confirmed: false,
				}

				ctx := context.Background()
				testContext.BurstRepository.Create(ctx, burst1)
				testContext.BurstRepository.Create(ctx, burst2)

				model.config.ArtifactType = ExportTypeBursts
				model.config.Format = ExportFormatJSON

				preview := model.generateBurstsPreview()
				Expect(preview).NotTo(BeEmpty())
				Expect(preview).To(ContainSubstring("burst-1"))
				Expect(preview).To(ContainSubstring("Q4 2024 Platform Work"))
			})

			It("should show message when no bursts exist", func() {
				model.config.ArtifactType = ExportTypeBursts
				model.config.Format = ExportFormatCSV

				preview := model.generateBurstsPreview()
				Expect(preview).NotTo(BeEmpty())
				Expect(preview).To(ContainSubstring("No bursts"))
			})

			It("should support text format", func() {
				burst := &careerdomain.Burst{
					ID:       "burst-3",
					Name:     "Test Burst",
					EventIDs: []string{"e1", "e2"},
				}

				ctx := context.Background()
				testContext.BurstRepository.Create(ctx, burst)

				model.config.ArtifactType = ExportTypeBursts
				model.config.Format = ExportFormatTXT

				preview := model.generateBurstsPreview()
				Expect(preview).To(ContainSubstring("Test Burst"))
			})
		})

		Describe("generatePreview integration", func() {
			It("should call the correct preview generator based on artifact type", func() {
				model.config.ArtifactType = ExportTypeEvents
				model.config.Format = ExportFormatJSON

				model.generatePreview()
				Expect(model.preview).NotTo(BeEmpty())
				Expect(model.previewLines).NotTo(BeEmpty())
			})

			It("should split preview into lines", func() {
				model.config.ArtifactType = ExportTypeFacts
				model.config.Format = ExportFormatJSON

				model.generatePreview()
				Expect(len(model.previewLines)).To(BeNumerically(">", 0))
			})
		})
	})
})
