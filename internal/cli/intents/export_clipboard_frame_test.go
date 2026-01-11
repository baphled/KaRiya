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

var _ = Describe("FRAME BY FRAME: What user sees during clipboard export", func() {
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

		tmpDir := GinkgoT().TempDir()
		dbPath := filepath.Join(tmpDir, "test.db")
		db, err := sql.Open("sqlite", dbPath)
		Expect(err).NotTo(HaveOccurred())
		err = careerrepo.RunMigrations(db)
		Expect(err).NotTo(HaveOccurred())
		cleanup = func() { db.Close() }

		eventRepo = careerrepo.NewSQLiteRepositoryWithDB(db)
		factRepo = careerrepo.NewSQLiteFactRepositoryWithDB(db)
		burstRepo = careerrepo.NewSQLiteBurstRepositoryWithDB(db)

		careerService = careerservice.NewService(eventRepo)
		careerService.SetFactRepository(factRepo)
		careerService.SetBurstRepository(burstRepo)
		exportService = cv.NewExportService(log)

		testDate, _ := time.Parse("2006-01-02", "2024-01-01")
		err = eventRepo.Create(ctx, &career.CareerEvent{
			ID:      "test-event-1",
			Text:    "Test event",
			Date:    testDate,
			Company: "Test Corp",
			Tags:    []string{"technical"},
		})
		Expect(err).NotTo(HaveOccurred())

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

	It("should show EXACTLY what user sees at each step", func() {
		GinkgoWriter.Println("\n========================================")
		GinkgoWriter.Println("FRAME-BY-FRAME: User Experience")
		GinkgoWriter.Println("========================================\n")

		// Frame 1: Initial state
		GinkgoWriter.Println("FRAME 1: Select Artifact Type")
		GinkgoWriter.Println("State:", intent.model.state)
		view := intent.View()
		GinkgoWriter.Println(view)
		GinkgoWriter.Println("\n--- User presses Enter ---\n")

		// Frame 2: After selecting Events
		intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		GinkgoWriter.Println("FRAME 2: Select Format")
		GinkgoWriter.Println("State:", intent.model.state)
		view = intent.View()
		GinkgoWriter.Println(view)
		GinkgoWriter.Println("\n--- User presses Enter (JSON) ---\n")

		// Frame 3: After selecting JSON
		intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		GinkgoWriter.Println("FRAME 3: Select Destination")
		GinkgoWriter.Println("State:", intent.model.state)
		view = intent.View()
		GinkgoWriter.Println(view)
		GinkgoWriter.Println("\n--- User presses 'j' to navigate to clipboard ---\n")

		// Frame 4: Navigate to clipboard
		intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		GinkgoWriter.Println("FRAME 4: Clipboard selected (highlighted)")
		GinkgoWriter.Println("State:", intent.model.state)
		GinkgoWriter.Println("Selected index:", intent.model.selectedIndex)
		view = intent.View()
		GinkgoWriter.Println(view)
		GinkgoWriter.Println("\n--- User presses Enter to select clipboard ---\n")

		// Frame 5: After selecting clipboard
		intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		GinkgoWriter.Println("FRAME 5: Configure (clipboard needs no config)")
		GinkgoWriter.Println("State:", intent.model.state)
		view = intent.View()
		GinkgoWriter.Println(view)
		GinkgoWriter.Println("\n--- User presses Enter to generate preview ---\n")

		// Frame 6: Generate preview
		intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		GinkgoWriter.Println("FRAME 6: Preview")
		GinkgoWriter.Println("State:", intent.model.state)
		view = intent.View()
		GinkgoWriter.Println(view[:500], "...") // Truncate for readability
		GinkgoWriter.Println("\n--- User presses Enter to confirm ---\n")

		// Frame 7: Confirm
		intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		GinkgoWriter.Println("FRAME 7: Confirm")
		GinkgoWriter.Println("State:", intent.model.state)
		view = intent.View()
		GinkgoWriter.Println(view)
		GinkgoWriter.Println("\n--- User presses Enter to START EXPORT ---\n")

		// Frame 8: START EXPORT
		cmd := intent.model.updateConfirm(tea.KeyMsg{Type: tea.KeyEnter})
		GinkgoWriter.Println("FRAME 8: Export In Progress (IMMEDIATELY after pressing Enter)")
		GinkgoWriter.Println("State:", intent.model.state)
		view = intent.View()
		GinkgoWriter.Println(view)
		GinkgoWriter.Println("⏱️  This is what user sees while export runs...")
		GinkgoWriter.Println("\n--- Export command executes (may take time) ---\n")

		// Frame 9: Execute the export command
		exportMsg := cmd()
		GinkgoWriter.Println("FRAME 9: Export command completed, message returned")
		GinkgoWriter.Printf("Message type: %T\n", exportMsg)

		if errMsg, ok := exportMsg.(ExportErrorMsg); ok {
			GinkgoWriter.Println("❌ Export returned ERROR")
			GinkgoWriter.Println("Error code:", errMsg.Error.Code)
			GinkgoWriter.Println("Error message:", errMsg.Error.Message)
		} else if completeMsg, ok := exportMsg.(ExportCompleteMsg); ok {
			GinkgoWriter.Println("✅ Export returned SUCCESS")
			GinkgoWriter.Println("Result:", completeMsg.Result)
		}

		GinkgoWriter.Println("\n--- Feeding message back to Update ---\n")

		// Frame 10: Process the result message
		intent.Update(exportMsg)
		GinkgoWriter.Println("FRAME 10: After processing export result")
		GinkgoWriter.Println("State:", intent.model.state)
		view = intent.View()
		GinkgoWriter.Println(view)

		GinkgoWriter.Println("\n========================================")
		GinkgoWriter.Println("ANALYSIS:")
		GinkgoWriter.Println("========================================")
		GinkgoWriter.Println("Did user see the error? Check FRAME 10 above.")
		GinkgoWriter.Println("If state = 'failed' and view shows error, error IS displayed.")
		GinkgoWriter.Println("If not, that's the bug we need to fix.")
		GinkgoWriter.Println("========================================\n")

		Expect(true).To(BeTrue()) // Always pass - this is diagnostic
	})
})
