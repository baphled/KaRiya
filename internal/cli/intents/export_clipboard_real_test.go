package intents

import (
	"context"
	"database/sql"
	"io"
	"path/filepath"
	"time"

	"github.com/atotto/clipboard"
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

var _ = Describe("REAL WORLD: Clipboard Export", func() {
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

		// Setup test database (exactly like app.go does)
		tmpDir := GinkgoT().TempDir()
		dbPath := filepath.Join(tmpDir, "test.db")

		db, err := sql.Open("sqlite", dbPath)
		Expect(err).NotTo(HaveOccurred())

		err = careerrepo.RunMigrations(db)
		Expect(err).NotTo(HaveOccurred())

		cleanup = func() {
			db.Close()
		}

		// Create repositories (exactly like app.go does)
		eventRepo = careerrepo.NewSQLiteRepositoryWithDB(db)
		factRepo = careerrepo.NewSQLiteFactRepositoryWithDB(db)
		burstRepo = careerrepo.NewSQLiteBurstRepositoryWithDB(db)

		// Create services (EXACTLY like app.go does at line 72)
		careerService = careerservice.NewService(eventRepo)
		careerService.SetFactRepository(factRepo)
		careerService.SetBurstRepository(burstRepo)

		// THIS IS THE KEY - exactly like app.go line 72
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

		// Create intent context (EXACTLY like app.go does at lines 580-591)
		exportCtx := &ExportArtifactContext{
			ArtifactTypes:    DefaultArtifactTypes(),
			SupportedFormats: DefaultSupportedFormats(),
			DefaultFormat:    DefaultFormats(),
			Destinations:     DefaultDestinations(),
			ExportService:    exportService,
			CareerService:    careerService,
			EventRepository:  eventRepo,
			FactRepository:   factRepo,
			BurstRepository:  burstRepo,
			AppContext:       ctx,
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

	Describe("REAL clipboard behavior (like production app)", func() {
		It("should show ACTUAL clipboard state and error", func() {
			GinkgoWriter.Println("=== REAL WORLD TEST ===")
			GinkgoWriter.Printf("clipboard.Unsupported = %v\n", clipboard.Unsupported)

			// Try to actually write to clipboard directly
			testErr := clipboard.WriteAll("test content")
			if testErr != nil {
				GinkgoWriter.Printf("Direct clipboard.WriteAll() failed: %v\n", testErr)
			} else {
				GinkgoWriter.Println("Direct clipboard.WriteAll() succeeded!")
			}

			// Now check what our ExportService reports
			sysClip := &cv.SystemClipboard{}
			GinkgoWriter.Printf("SystemClipboard.IsUnsupported() = %v\n", sysClip.IsUnsupported())

			// Try our ExportService method
			exportErr := exportService.CopyToClipboard(ctx, "test via export service")
			if exportErr != nil {
				GinkgoWriter.Printf("ExportService.CopyToClipboard() failed: %v\n", exportErr)
			} else {
				GinkgoWriter.Println("ExportService.CopyToClipboard() succeeded!")
			}

			// Now run through the FULL user flow
			GinkgoWriter.Println("\n=== FULL USER FLOW ===")

			// User navigates to export
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})                     // Select events
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})                     // Select JSON
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}) // Navigate to clipboard
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})                     // Select clipboard
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})                     // Generate preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})                     // Go to confirm

			// Start export
			cmd := intent.model.updateConfirm(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())

			// Execute the export
			exportMsg := cmd()
			intent.Update(exportMsg)

			// What happened?
			GinkgoWriter.Printf("Final state: %s\n", intent.model.state)

			if intent.model.state == ExportStateFailed {
				GinkgoWriter.Printf("❌ Export FAILED\n")
				GinkgoWriter.Printf("Error code: %s\n", intent.model.error.Code)
				GinkgoWriter.Printf("Error message: %s\n", intent.model.error.Message)

				// Show what user sees
				view := intent.View()
				GinkgoWriter.Println("\n=== WHAT USER SEES ===")
				GinkgoWriter.Println(view)
				GinkgoWriter.Println("======================")

			} else if intent.model.state == ExportStateComplete {
				GinkgoWriter.Printf("✅ Export SUCCEEDED\n")
				GinkgoWriter.Printf("Result: %+v\n", intent.model.result)
			} else {
				GinkgoWriter.Printf("⚠️  Unexpected state: %s\n", intent.model.state)
			}

			// This test always passes - it just shows what happens
			Expect(true).To(BeTrue())
		})
	})
})
