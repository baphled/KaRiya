package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/cli/importer"
	cliservice "github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	version = "0.1.0"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, out io.Writer, errOut io.Writer) int {
	// Parse CLI flags
	var (
		showVersion = false
		showHelp    = false
		dbPath      = ""
		mode        = ""
		listEvents  = false
		inMemory    = false
		importPath  = ""
		importSkip  = false
	)

	// Parse command-line arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--version", "-v":
			showVersion = true
		case "--help", "-h":
			showHelp = true
		case "--db", "--database":
			if i+1 < len(args) {
				dbPath = args[i+1]
				i++ // Skip next argument as it's the value
			}
		case "--mode":
			if i+1 < len(args) {
				mode = args[i+1]
				// Validate mode
				if mode != "timeline" && mode != "backfill" && mode != "manual" {
					fmt.Fprintf(errOut, "Error: Invalid mode '%s'. Valid modes are: timeline, backfill, manual\n", mode)
					return 1
				}
				i++ // Skip next argument as it's the value
			}
		case "--list":
			listEvents = true
		case "--in-memory":
			inMemory = true
		case "--import":
			if i+1 < len(args) {
				importPath = args[i+1]
				i++ // Skip next argument as it's the value
			} else {
				fmt.Fprintf(errOut, "Error: --import flag requires a file path\n")
				return 1
			}
		case "--skip-import-review":
			importSkip = true
		}
	}

	// Handle version flag
	if showVersion {
		fmt.Fprintf(out, "KaRiya CLI v%s\n", version)
		return 0
	}

	// Handle help flag
	if showHelp {
		printHelpTo(out)
		return 0
	}

	// Set up repository based on flags
	var repo career.Repository
	var err error

	if inMemory {
		// Use in-memory repository if explicitly requested
		repo = career.NewMemoryRepository()
	} else {
		// Determine database path
		if dbPath == "" {
			// Use default path: ~/.kariya/events.db
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Fprintf(errOut, "Error getting home directory: %v\n", err)
				return 1
			}
			kariyaDir := filepath.Join(homeDir, ".kariya")
			dbPath = filepath.Join(kariyaDir, "events.db")

			// Create directory if it doesn't exist
			if err := os.MkdirAll(kariyaDir, 0755); err != nil {
				fmt.Fprintf(errOut, "Error creating kariya directory: %v\n", err)
				return 1
			}
		}

		// Use SQLite repository with specified or default path
		repo, err = career.NewSQLiteRepository(dbPath)
		if err != nil {
			fmt.Fprintf(errOut, "Error initializing database at '%s': %v\n", dbPath, err)
			return 1
		}
	}

	svc := careerservice.NewService(repo)

	// Initialize fact and burst repositories if using SQLite (not in-memory)
	if !inMemory {
		// Get the DB connection from the event repository
		sqliteRepo, ok := repo.(*career.SQLiteRepository)
		if ok && sqliteRepo != nil {
			db := sqliteRepo.GetDB()

			// Initialize fact repository
			factRepo, err := career.NewSQLiteFactRepository(db)
			if err != nil {
				// Log warning but continue - facts are optional enhancement
				fmt.Fprintf(errOut, "Warning: Failed to initialize fact repository: %v\n", err)
			} else {
				svc.SetFactRepository(factRepo)
			}

			// Initialize burst repository
			burstRepo, err := career.NewSQLiteBurstRepository(db)
			if err != nil {
				// Log warning but continue - bursts are optional enhancement
				fmt.Fprintf(errOut, "Warning: Failed to initialize burst repository: %v\n", err)
			} else {
				svc.SetBurstRepository(burstRepo)
			}
		}
	}

	cliSvc := cliservice.NewCLIEventService(svc)

	// Handle non-interactive import if --import flag is provided
	if importPath != "" {
		return handleNonInteractiveImport(importPath, importSkip, svc, out, errOut)
	}

	// Initialize application model
	model := app.NewModel(cliSvc, svc)

	// Set initial capture mode if specified
	if mode != "" {
		model.SetInitialCaptureMode(mode)
	}

	// Set initial screen if --list flag is provided
	if listEvents {
		model.SetInitialScreen(app.ListScreen)
	}

	// Initialize BubbleTea program with mouse support
	p := tea.NewProgram(model, tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(errOut, "Error running program: %v\n", err)
		return 1
	}
	return 0
}

// handleNonInteractiveImport performs import without showing the interactive UI
func handleNonInteractiveImport(filePath string, skipReview bool, svc *careerservice.Service, out io.Writer, errOut io.Writer) int {
	// Validate file exists
	if _, err := os.Stat(filePath); err != nil {
		fmt.Fprintf(errOut, "Error: Cannot access import file '%s': %v\n", filePath, err)
		return 1
	}

	ctx := context.Background()

	// Create import service
	importService := importer.NewImportService(svc)

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(errOut, "Error opening import file: %v\n", err)
		return 1
	}
	defer file.Close()

	// Parse CSV
	fmt.Fprintf(out, "Parsing CSV file: %s\n", filePath)
	parsedRows, err := importService.PrepareImport(ctx, file)
	if err != nil {
		fmt.Fprintf(errOut, "Error parsing CSV: %v\n", err)
		return 1
	}

	if len(parsedRows) == 0 {
		fmt.Fprintf(errOut, "Error: No valid rows found in CSV file\n")
		return 1
	}

	fmt.Fprintf(out, "Found %d rows to import\n", len(parsedRows))

	// Show preview if not skipping review
	if !skipReview {
		fmt.Fprintf(out, "\nPreview of rows to import:\n")
		fmt.Fprintf(out, "%-50s | %-15s | %-20s\n", "Event Text", "Date", "Company")
		fmt.Fprintf(out, "%s+%s+%s\n", "---------------------------------------------------", "----------------", "---------------------")
		for i, row := range parsedRows {
			if i >= 5 { // Show first 5 rows
				fmt.Fprintf(out, "... and %d more rows\n", len(parsedRows)-5)
				break
			}

			// Skip invalid or duplicate rows in preview
			if !row.IsValid || row.IsDuplicate {
				continue
			}

			dateStr := ""
			if row.Event != nil && !row.Event.Date.IsZero() {
				dateStr = row.Event.Date.Format("2006-01-02")
			}
			eventText := ""
			company := ""
			if row.Event != nil {
				eventText = row.Event.Text
				company = row.Event.Company
			}
			if len(eventText) > 50 {
				eventText = eventText[:47] + "..."
			}
			fmt.Fprintf(out, "%-50s | %-15s | %-20s\n", eventText, dateStr, company)
		}

		// Ask for confirmation
		fmt.Fprintf(out, "\nProceed with import? (y/n): ")
		var response string
		_, err := fmt.Scanln(&response)
		if err != nil || (response != "y" && response != "Y") {
			fmt.Fprintf(out, "Import cancelled.\n")
			return 0
		}
	}

	// Perform import
	fmt.Fprintf(out, "Importing %d rows...\n", len(parsedRows))
	selectedRows := make([]int, len(parsedRows))
	for i := range parsedRows {
		selectedRows[i] = i
	}

	result, err := importService.ImportRows(ctx, parsedRows, selectedRows)
	if err != nil {
		fmt.Fprintf(errOut, "Error during import: %v\n", err)
		return 1
	}

	// Print results
	fmt.Fprintf(out, "\n=== Import Complete ===\n")
	fmt.Fprintf(out, "Total rows processed: %d\n", result.TotalRows)
	fmt.Fprintf(out, "Successfully imported: %d\n", result.SuccessCount)
	fmt.Fprintf(out, "Skipped: %d\n", result.SkippedCount)
	if result.FailedCount > 0 {
		fmt.Fprintf(out, "Failed: %d\n", result.FailedCount)
		if len(result.FailedRows) > 0 {
			fmt.Fprintf(out, "\nFailed rows:\n")
			for _, failedRow := range result.FailedRows {
				fmt.Fprintf(out, "  - Row %d", failedRow.RowNumber)
				if failedRow.Event != nil {
					fmt.Fprintf(out, ": %s", failedRow.Event.Text)
				}
				fmt.Fprintf(out, "\n")
				if len(failedRow.ValidationErrors) > 0 {
					for _, errMsg := range failedRow.ValidationErrors {
						fmt.Fprintf(out, "    Error: %s\n", errMsg)
					}
				}
			}
		}
	}

	// Display burst detection results (Task 2.1)
	if result.BurstSuggestions != nil && len(result.BurstSuggestions) > 0 {
		fmt.Fprintf(out, "\n=== Burst Suggestions ===\n")
		fmt.Fprintf(out, "Detected %d potential bursts from imported events:\n\n", len(result.BurstSuggestions))
		for i, burst := range result.BurstSuggestions {
			burstName := burst.Name
			if burstName == "" {
				burstName = fmt.Sprintf("Burst %d", i+1)
			}
			fmt.Fprintf(out, "%d. %s\n", i+1, burstName)
			fmt.Fprintf(out, "   Events: %d | Confidence: %.1f%%\n", len(burst.EventIDs), burst.ConfidenceScore*100)
		}
		fmt.Fprintf(out, "\nThese bursts represent potential project groupings or themes.\n")
	}

	if result.SuccessCount > 0 {
		fmt.Fprintf(out, "\n✓ Import successful!\n")
		return 0
	} else if result.SkippedCount > 0 {
		fmt.Fprintf(out, "\n⚠ All rows were skipped (possibly duplicates).\n")
		return 0
	} else {
		fmt.Fprintf(errOut, "\n✗ Import failed - no rows were imported.\n")
		return 1
	}
}

func printHelpTo(out io.Writer) {
	fmt.Fprintln(out, "KaRiya CLI - Career Journaling Tool")
	fmt.Fprintln(out, "\nUsage: kariya [options]")
	fmt.Fprintln(out, "\nOptions:")
	fmt.Fprintln(out, "  -v, --version              Show version information")
	fmt.Fprintln(out, "  -h, --help                 Show this help message")
	fmt.Fprintln(out, "  --db, --database PATH      Use custom database path (SQLite)")
	fmt.Fprintln(out, "                             Default: ~/.kariya/events.db")
	fmt.Fprintln(out, "  --mode MODE                Start in specific capture mode")
	fmt.Fprintln(out, "                             Valid modes: timeline, backfill, manual")
	fmt.Fprintln(out, "  --list                     Show recent events on startup")
	fmt.Fprintln(out, "  --in-memory                Use in-memory storage (data not persisted)")
	fmt.Fprintln(out, "  --import PATH              Import career events from CSV file")
	fmt.Fprintln(out, "  --skip-import-review       Skip confirmation prompt during import")
	fmt.Fprintln(out, "\nExamples:")
	fmt.Fprintln(out, "  kariya                                    # Start with default database")
	fmt.Fprintln(out, "  kariya --db ./events.db                  # Use custom database path")
	fmt.Fprintln(out, "  kariya --mode timeline                   # Start in timeline journaling mode")
	fmt.Fprintln(out, "  kariya --db ./events.db --list           # Open with events list")
	fmt.Fprintln(out, "  kariya --in-memory                       # Start with in-memory storage")
	fmt.Fprintln(out, "  kariya --import events.csv               # Import events from CSV file")
	fmt.Fprintln(out, "  kariya --import events.csv --skip-import-review")
	fmt.Fprintln(out, "                                           # Import without confirmation")
	fmt.Fprintln(out, "  kariya --version                         # Show version")
	fmt.Fprintln(out, "\nCSV File Format:")
	fmt.Fprintln(out, "  The CSV file should have the following columns (in any order):")
	fmt.Fprintln(out, "  - Text (required): Event description")
	fmt.Fprintln(out, "  - Date (required): Event date (YYYY-MM-DD format)")
	fmt.Fprintln(out, "  - Company (optional): Company name")
	fmt.Fprintln(out, "  - Project (optional): Project name")
	fmt.Fprintln(out, "  - Tags (optional): Semicolon-separated tags (e.g., technical;leadership)")
	fmt.Fprintln(out, "\nFor more information, visit: https://github.com/baphled/kariya")
}
