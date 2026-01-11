package intents

import (
	"context"
	"database/sql"
	"io"
	"path/filepath"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	_ "modernc.org/sqlite"
)

var _ = Describe("ExportArtifact Clipboard E2E", func() {
	var (
		intent        *ExportArtifactIntent
		ctx           context.Context
		eventRepo     careerrepo.Repository
		factRepo      careerrepo.FactRepository
		burstRepo     careerrepo.BurstRepository
		careerService *careerservice.Service
		exportService *cv.ExportService
		log           *logger.Logger
		cleanup       func()
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		ctx = context.Background()

		// Setup test database
		tmpDir := GinkgoT().TempDir()
		dbPath := filepath.Join(tmpDir, "test.db")

		db, err := sql.Open("sqlite", dbPath)
		Expect(err).NotTo(HaveOccurred())

		err = careerrepo.RunMigrations(db)
		Expect(err).NotTo(HaveOccurred())

		cleanup = func() {
			db.Close()
		}

		// Create repositories
		eventRepo = careerrepo.NewSQLiteRepositoryWithDB(db)
		factRepo = careerrepo.NewSQLiteFactRepositoryWithDB(db)
		burstRepo = careerrepo.NewSQLiteBurstRepositoryWithDB(db)

		// Create services
		careerService = careerservice.NewService(eventRepo)
		careerService.SetFactRepository(factRepo)
		careerService.SetBurstRepository(burstRepo)
		exportService = cv.NewExportService(log)

		// Create test data
		testDate, _ := time.Parse("2006-01-02", "2024-01-01")
		err = eventRepo.Create(ctx, &career.CareerEvent{
			ID:      "test-event-1",
			Text:    "Developed new feature for testing clipboard export",
			Date:    testDate,
			Company: "Test Corp",
			Project: "Clipboard Testing",
			Tags:    []string{"technical"},
		})
		Expect(err).NotTo(HaveOccurred())

		// Create intent context
		exportCtx := &ExportArtifactContext{
			ArtifactTypes: []ExportArtifactType{ExportTypeEvents},
			SupportedFormats: map[ExportArtifactType][]ExportFormat{
				ExportTypeEvents: {ExportFormatJSON, ExportFormatCSV},
			},
			DefaultFormat: map[ExportArtifactType]ExportFormat{
				ExportTypeEvents: ExportFormatJSON,
			},
			Destinations:    []ExportDestination{ExportDestinationFile, ExportDestinationClipboard},
			ExportService:   exportService,
			CareerService:   careerService,
			EventRepository: eventRepo,
			FactRepository:  factRepo,
			BurstRepository: burstRepo,
			AppContext:      ctx,
		}

		var initErr error
		intent, initErr = NewExportArtifactIntent(exportCtx)
		Expect(initErr).NotTo(HaveOccurred())
		intent.Init()
	})

	AfterEach(func() {
		if cleanup != nil {
			cleanup()
		}
	})

	Describe("Clipboard Export User Flow", func() {
		It("should show error when clipboard export fails", func() {
			// Navigate through the flow:
			// 1. Select artifact type (events)
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.model.state).To(Equal(ExportStateSelectFormat))

			// 2. Select format (JSON - default)
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.model.state).To(Equal(ExportStateSelectDest))

			// 3. Navigate to clipboard destination
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}) // Move down to clipboard
			Expect(intent.model.selectedIndex).To(Equal(1))                   // clipboard is index 1

			// 4. Select clipboard
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.model.state).To(Equal(ExportStateConfigure))

			// 5. Press Enter again to generate preview (clipboard needs no path configuration)
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.model.state).To(Equal(ExportStatePreview))

			// 6. Confirm export (Enter to go to confirm state)
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.model.state).To(Equal(ExportStateConfirm))

			// 7. Press Enter to actually start export
			// This calls startExport() which returns a tea.Cmd
			cmd := intent.model.updateConfirm(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.model.state).To(Equal(ExportStateInProgress))
			Expect(cmd).NotTo(BeNil(), "startExport() should return a command")

			// 8. Execute the export command to get the result message
			// In a real TUI, Bubble Tea runtime would execute this automatically
			exportMsg := cmd()
			Expect(exportMsg).NotTo(BeNil(), "Export command should return a message")

			// 9. Feed the message back to process the result
			intent.Update(exportMsg)

			// The state should transition based on clipboard availability
			// In headless environment: ExportStateFailed
			// With clipboard: ExportStateComplete

			if intent.model.state == ExportStateFailed {
				// Headless environment - clipboard failed
				GinkgoWriter.Println("✅ Clipboard export failed as expected in headless environment")
				Expect(intent.model.error).NotTo(BeNil())
				Expect(intent.model.error.Message).To(ContainSubstring("Export failed"))

				// Verify error is shown in view
				view := intent.View()
				Expect(view).To(ContainSubstring("Export Failed"))
				Expect(view).To(ContainSubstring("Error:"))
				GinkgoWriter.Printf("Error message shown to user: %s\n", intent.model.error.Message)

			} else if intent.model.state == ExportStateComplete {
				// Clipboard succeeded (macOS/Windows or Linux with display)
				GinkgoWriter.Println("✅ Clipboard export succeeded on this platform")
				Expect(intent.model.result).NotTo(BeNil())
				Expect(intent.model.result.Success).To(BeTrue())
				Expect(intent.model.result.FilePath).To(Equal("clipboard"))

			} else {
				Fail("Unexpected state: " + string(intent.model.state))
			}
		})

		It("should allow retry after clipboard failure", func() {
			// Simulate the same flow as above but test retry
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})                     // Select events
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})                     // Select JSON
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}) // Navigate to clipboard
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})                     // Select clipboard
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})                     // Generate preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})                     // Go to confirm

			// Start export and execute the command
			cmd := intent.model.updateConfirm(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())
			exportMsg := cmd()
			intent.Update(exportMsg)

			if intent.model.state == ExportStateFailed {
				// Test retry option
				view := intent.View()
				Expect(view).To(ContainSubstring("Retry"))

				// Press 'R' to retry
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
				// Should go back to confirm (retry = go back to confirmation step)
				Expect(intent.model.state).To(Equal(ExportStateConfirm))
				GinkgoWriter.Println("✅ Retry option works after clipboard failure")
			}
		})
	})
})
