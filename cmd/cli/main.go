package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/cli/bootstrap"
	"github.com/baphled/kariya/internal/cli/importer"
	cliservice "github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/logger"
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
		showVersion  = false
		showHelp     = false
		dbPath       = ""
		mode         = ""
		listEvents   = false
		inMemory     = false
		importPath   = ""
		importSkip   = false
		reviewFacts  = false
		detectBursts = false
		extractFacts = false
		showBursts   = false
		showFacts    = false
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
				i++
			}
		case "--mode":
			if i+1 < len(args) {
				mode = args[i+1]
				if mode != "timeline" && mode != "backfill" && mode != "manual" {
					fmt.Fprintf(errOut, "Error: Invalid mode '%s'. Valid modes are: timeline, backfill, manual\n", mode)
					return 1
				}
				i++
			}
		case "--list":
			listEvents = true
		case "--in-memory":
			inMemory = true
		case "--import":
			if i+1 >= len(args) {
				fmt.Fprintf(errOut, "Error: --import flag requires a file path\n")
				return 1
			}
			importPath = args[i+1]
			i++
		case "--skip-import-review":
			importSkip = true
		case "--review-facts":
			reviewFacts = true
		case "--detect-bursts":
			detectBursts = true
		case "--extract-facts":
			extractFacts = true
		case "--show-bursts":
			showBursts = true
		case "--show-facts":
			showFacts = true
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

	// Set up repository
	var repo career.Repository

	if inMemory {
		repo = career.NewMemoryRepository()
	} else {
		if dbPath == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Fprintf(errOut, "Error getting home directory: %v\n", err)
				return 1
			}
			kariyaDir := filepath.Join(homeDir, ".kariya")
			dbPath = filepath.Join(kariyaDir, "events.db")

			if err := os.MkdirAll(kariyaDir, 0750); err != nil {
				fmt.Fprintf(errOut, "Error creating kariya directory: %v\n", err)
				return 1
			}
		}

		// Open database connection
		db, err := sql.Open("sqlite", dbPath)
		if err != nil {
			fmt.Fprintf(errOut, "Error opening database at '%s': %v\n", dbPath, err)
			return 1
		}

		// Run migrations
		if err := career.RunMigrations(db); err != nil {
			fmt.Fprintf(errOut, "Error running migrations: %v\n", err)
			return 1
		}

		// Create repository with existing connection
		repo = career.NewSQLiteRepositoryWithDB(db)
	}

	svc := careerservice.NewService(repo)

	// Initialize fact and burst repositories
	if inMemory {
		// Use in-memory repositories for facts and bursts
		factRepo := career.NewMemoryFactRepository()
		svc.SetFactRepository(factRepo)

		burstRepo := career.NewMemoryBurstRepository()
		svc.SetBurstRepository(burstRepo)
	} else {
		// Use SQLite repositories for facts and bursts (migrations already run)
		sqliteRepo, ok := repo.(*career.SQLiteRepository)
		if ok && sqliteRepo != nil {
			db := sqliteRepo.GetDB()

			// Use the *WithDB constructors since migrations are already applied
			factRepo := career.NewSQLiteFactRepositoryWithDB(db)
			svc.SetFactRepository(factRepo)

			burstRepo := career.NewSQLiteBurstRepositoryWithDB(db)
			svc.SetBurstRepository(burstRepo)

			skillRepo := career.NewSQLiteSkillRepositoryWithDB(db)
			svc.SetSkillRepository(skillRepo)
		}
	}

	cliSvc := cliservice.NewCLIEventService(svc)

	// Handle burst/fact operations
	if detectBursts {
		return handleDetectBursts(svc, out, errOut)
	}

	if extractFacts {
		return handleExtractFacts(svc, out, errOut)
	}

	if showBursts {
		return handleShowBursts(svc, out, errOut)
	}

	if showFacts {
		return handleShowFacts(svc, out, errOut)
	}

	// Handle non-interactive import
	if importPath != "" {
		return handleNonInteractiveImport(importPath, importSkip, reviewFacts, svc, out, errOut)
	}

	// Run bootstrap (handles onboarding and service initialization)
	log := logger.DefaultLogger()
	bootstrapResult, err := bootstrap.Run(svc, log)
	if err != nil {
		fmt.Fprintf(errOut, "Error during bootstrap: %v\n", err)
		return 1
	}

	// User aborted onboarding (Ctrl+C)
	if bootstrapResult == nil {
		return 0
	}

	// Initialize application model with bootstrap results
	model := app.NewModel(cliSvc, svc, bootstrapResult)

	if mode != "" {
		model.SetInitialCaptureMode(mode)
	}

	if listEvents {
		model.SetInitialScreen(app.ListScreen)
	}

	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(errOut, "Error running program: %v\n", err)
		return 1
	}
	return 0
}

// handleDetectBursts re-runs burst detection on all events
func handleDetectBursts(svc *careerservice.Service, out io.Writer, errOut io.Writer) int {
	ctx := context.Background()

	events, err := svc.ListEvents(ctx, career.ListFilters{Limit: 10000})
	if err != nil {
		fmt.Fprintf(errOut, "Error retrieving events: %v\n", err)
		return 1
	}

	if len(events) == 0 {
		fmt.Fprintf(out, "No events found in database.\n")
		return 0
	}

	eventIDs := make([]string, len(events))
	for i, event := range events {
		eventIDs[i] = event.ID
	}

	fmt.Fprintf(out, "Detecting bursts from %d events...\n", len(events))

	suggestions, err := svc.SuggestBursts(ctx, eventIDs)
	if err != nil {
		fmt.Fprintf(errOut, "Error detecting bursts: %v\n", err)
		return 1
	}

	if len(suggestions) == 0 {
		fmt.Fprintf(out, "No bursts detected.\n")
		return 0
	}

	fmt.Fprintf(out, "\n=== Burst Detection Results ===\n")
	fmt.Fprintf(out, "Detected %d bursts from %d events:\n\n", len(suggestions), len(events))

	// Display suggestions first
	for i, burst := range suggestions {
		burstName := burst.Name
		if burstName == "" {
			burstName = fmt.Sprintf("Burst %d", i+1)
		}
		fmt.Fprintf(out, "%d. %s\n", i+1, burstName)
		fmt.Fprintf(out, "   Events: %d\n", len(burst.EventIDs))
		if burst.Description != "" {
			fmt.Fprintf(out, "   Description: %s\n", burst.Description)
		}
		fmt.Fprintf(out, "   Confidence: %.1f%%\n", burst.ConfidenceScore*100)
		fmt.Fprintf(out, "\n")
	}

	// Save all suggestions as bursts
	savedBursts, err := svc.SaveBurstSuggestions(ctx, suggestions)
	if err != nil {
		fmt.Fprintf(errOut, "Error saving burst suggestions: %v\n", err)
		return 1
	}

	savedCount := len(savedBursts)
	if savedCount > 0 {
		fmt.Fprintf(out, "✓ Burst detection complete! Saved %d of %d bursts to database.\n", savedCount, len(suggestions))
	} else {
		fmt.Fprintf(out, "⚠ Burst detection complete, but no bursts were saved (repository may not be configured).\n")
	}

	return 0
}

// handleExtractFacts re-runs fact extraction on all events
func handleExtractFacts(svc *careerservice.Service, out io.Writer, errOut io.Writer) int {
	ctx := context.Background()

	events, err := svc.ListEvents(ctx, career.ListFilters{Limit: 10000})
	if err != nil {
		fmt.Fprintf(errOut, "Error retrieving events: %v\n", err)
		return 1
	}

	if len(events) == 0 {
		fmt.Fprintf(out, "No events found in database.\n")
		return 0
	}

	fmt.Fprintf(out, "Extracting facts from %d events...\n", len(events))

	factCount := 0
	competencyCount := make(map[string]int)

	for _, event := range events {
		facts, err := svc.ExtractFactsFromEvent(ctx, event)
		if err != nil {
			fmt.Fprintf(errOut, "Warning: Failed to extract facts from event %s: %v\n", event.ID, err)
			continue
		}

		for i := range facts {
			if err := svc.SaveFact(ctx, &facts[i]); err != nil {
				// Silently skip if fact repository is not configured (expected in some scenarios)
				if !errors.Is(err, careerservice.ErrFactRepositoryNotConfigured) {
					fmt.Fprintf(errOut, "Warning: Failed to save fact: %v\n", err)
				}
			} else {
				factCount++
				if len(facts[i].CompetencyCategories) > 0 {
					for _, comp := range facts[i].CompetencyCategories {
						competencyCount[comp]++
					}
				}
			}
		}
	}

	if factCount == 0 {
		fmt.Fprintf(out, "No facts extracted.\n")
		return 0
	}

	fmt.Fprintf(out, "\n=== Fact Extraction Results ===\n")
	fmt.Fprintf(out, "Extracted %d facts from %d events\n\n", factCount, len(events))

	if len(competencyCount) > 0 {
		fmt.Fprintf(out, "Competency breakdown:\n")
		competencies := make([]string, 0, len(competencyCount))
		for c := range competencyCount {
			competencies = append(competencies, c)
		}
		sort.Strings(competencies)

		for _, c := range competencies {
			fmt.Fprintf(out, "  - %s: %d facts\n", c, competencyCount[c])
		}
		fmt.Fprintf(out, "\n")
	}

	fmt.Fprintf(out, "✓ Fact extraction complete!\n")
	return 0
}

// handleShowBursts displays all existing bursts
func handleShowBursts(svc *careerservice.Service, out io.Writer, errOut io.Writer) int {
	ctx := context.Background()

	burstRepo := svc.GetBurstRepository()
	if burstRepo == nil {
		fmt.Fprintf(errOut, "Error: Burst repository not configured\n")
		return 1
	}

	bursts, err := burstRepo.List(ctx, career.BurstListFilters{Limit: 10000})
	if err != nil {
		fmt.Fprintf(errOut, "Error retrieving bursts: %v\n", err)
		return 1
	}

	if len(bursts) == 0 {
		fmt.Fprintf(out, "No bursts found.\n")
		return 0
	}

	fmt.Fprintf(out, "\n=== Bursts ===\n")

	fmt.Fprintf(out, "Total bursts: %d\n\n", len(bursts))

	for i, burst := range bursts {
		burstName := burst.Name
		if burstName == "" {
			burstName = fmt.Sprintf("Burst %d", i+1)
		}
		fmt.Fprintf(out, "%d. %s\n", i+1, burstName)
		fmt.Fprintf(out, "   ID: %s\n", burst.ID)
		fmt.Fprintf(out, "   Events: %d\n", len(burst.EventIDs))
		if burst.Description != "" {
			fmt.Fprintf(out, "   Description: %s\n", burst.Description)
		}
		if burst.Description != "" {
			fmt.Fprintf(out, "   Competency Focus: %s\n", burst.Description)
		}
		fmt.Fprintf(out, "   Created: %s\n", burst.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Fprintf(out, "\n")
	}

	return 0
}

// handleShowFacts displays all existing facts
func handleShowFacts(svc *careerservice.Service, out io.Writer, errOut io.Writer) int {
	ctx := context.Background()

	factRepo := svc.GetFactRepository()
	if factRepo == nil {
		fmt.Fprintf(errOut, "Error: Fact repository not configured\n")
		return 1
	}

	facts, err := factRepo.List(ctx, career.FactListFilters{Limit: 10000})
	if err != nil {
		fmt.Fprintf(errOut, "Error retrieving facts: %v\n", err)
		return 1
	}

	if len(facts) == 0 {
		fmt.Fprintf(out, "No facts found in database.\n")
		return 0
	}

	fmt.Fprintf(out, "\n=== Existing Facts ===\n")
	fmt.Fprintf(out, "Total facts: %d\n\n", len(facts))

	for i, fact := range facts {
		fmt.Fprintf(out, "%d. %s\n", i+1, fact.Text)
		fmt.Fprintf(out, "   ID: %s\n", fact.ID)
		fmt.Fprintf(out, "   Source Event: %s\n", fact.SourceEventID)
		if len(fact.CompetencyCategories) > 0 {
			fmt.Fprintf(out, "   Competencies: %s\n", strings.Join(fact.CompetencyCategories, ", "))
		}
		if fact.RoleFit != "" {
			fmt.Fprintf(out, "   Role Fit: %s\n", fact.RoleFit)
		}
		if len(fact.AudienceRelevance) > 0 {
			fmt.Fprintf(out, "   Target Audiences: %s\n", strings.Join(fact.AudienceRelevance, ", "))
		}
		fmt.Fprintf(out, "   Created: %s\n", fact.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Fprintf(out, "\n")
	}

	return 0
}

// handleNonInteractiveImport performs import without showing the interactive UI.
func handleNonInteractiveImport(filePath string, _ bool, _ bool, svc *careerservice.Service, out io.Writer, errOut io.Writer) int {
	if _, err := os.Stat(filePath); err != nil {
		fmt.Fprintf(errOut, "Error: Cannot access import file '%s': %v\n", filePath, err)
		return 1
	}

	ctx := context.Background()

	importService := importer.NewImportService(svc)

	file, err := os.Open(filePath) // #nosec G304 -- user-provided import file path (intentional)
	if err != nil {
		fmt.Fprintf(errOut, "Error opening import file: %v\n", err)
		return 1
	}
	defer file.Close()

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

	// Select all rows
	selectedRows := make([]int, len(parsedRows))
	for i := range parsedRows {
		selectedRows[i] = i
	}

	fmt.Fprintf(out, "Importing %d rows...\n", len(parsedRows))
	result, err := importService.ImportRows(ctx, parsedRows, selectedRows)
	if err != nil {
		fmt.Fprintf(errOut, "Error during import: %v\n", err)
		return 1
	}

	fmt.Fprintf(out, "\nImport Complete ===\n")
	fmt.Fprintf(out, "Total rows processed: %d\n", result.TotalRows)
	fmt.Fprintf(out, "Successfully imported: %d\n", result.SuccessCount)
	fmt.Fprintf(out, "Skipped: %d\n", result.SkippedCount)
	if result.FailedCount > 0 {
		fmt.Fprintf(out, "Failed: %d\n", result.FailedCount)
		if len(result.FailedRows) > 0 {
			fmt.Fprintf(out, "\nFailed rows (first 5):\n")
			for i, failedRow := range result.FailedRows {
				if i >= 5 {
					break
				}
				if failedRow.Event != nil {
					fmt.Fprintf(out, "  - %s\n", failedRow.Event.Text)
				}
			}
		}
	}

	if len(result.BurstSuggestions) > 0 {
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

	if result.ExtractedFactsCount > 0 {
		fmt.Fprintf(out, "\n=== Fact Extraction ===\n")
		fmt.Fprintf(out, "Extracted %d facts from %d events\n", result.ExtractedFactsCount, len(result.CreatedEvents))

		if len(result.FactsByCompetency) > 0 {
			fmt.Fprintf(out, "\nCompetency breakdown:\n")
			competencies := make([]string, 0, len(result.FactsByCompetency))
			for c := range result.FactsByCompetency {
				competencies = append(competencies, c)
			}
			sort.Strings(competencies)
			for _, c := range competencies {
				fmt.Fprintf(out, "  - %s: %d facts\n", c, result.FactsByCompetency[c])
			}
		}
		fmt.Fprintf(out, "\nThese facts highlight key competencies and achievements.\n")
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
	fmt.Fprintln(out, "  --detect-bursts            Re-run burst detection on all events")
	fmt.Fprintln(out, "  --extract-facts            Re-run fact extraction on all events")
	fmt.Fprintln(out, "  --show-bursts              Display all existing bursts")
	fmt.Fprintln(out, "  --show-facts               Display all existing facts")
	fmt.Fprintln(out, "\nExamples:")
	fmt.Fprintln(out, "  kariya                                    # Start with default database")
	fmt.Fprintln(out, "  kariya --db ./events.db                  # Use custom database path")
	fmt.Fprintln(out, "  kariya --mode timeline                   # Start in timeline journaling mode")
	fmt.Fprintln(out, "  kariya --db ./events.db --list           # Open with events list")
	fmt.Fprintln(out, "  kariya --in-memory                       # Start with in-memory storage")
	fmt.Fprintln(out, "  kariya --import events.csv               # Import events from CSV file")
	fmt.Fprintln(out, "  kariya --import events.csv --skip-import-review")
	fmt.Fprintln(out, "                                           # Import without confirmation")
	fmt.Fprintln(out, "  kariya --detect-bursts                   # Re-detect all bursts")
	fmt.Fprintln(out, "  kariya --extract-facts                   # Re-extract all facts")
	fmt.Fprintln(out, "  kariya --show-bursts                     # List all bursts")
	fmt.Fprintln(out, "  kariya --show-facts                      # List all facts")
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
