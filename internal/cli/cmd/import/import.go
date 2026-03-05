// Package importcmd provides Cobra CLI command for importing events from CSV files.
package importcmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/cli/importer"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// ImportParams holds parameters for the import operation.
type ImportParams struct {
	Reader  io.Reader
	Service *careerservice.Service
	Out     io.Writer
	ErrOut  io.Writer
}

// NewImportCmd creates the "import" command.
func NewImportCmd(ctx cliutil.ServiceContext) *cobra.Command {
	var filePath string

	importCmd := &cobra.Command{
		Use:   "import",
		Short: "Import events from CSV",
		Long:  "Import events from a CSV file and process them for facts and bursts",
		RunE: func(cobraCmd *cobra.Command, _ []string) error {
			svc := ctx.Service()
			if svc == nil {
				return errors.New("service not initialized")
			}

			// Validate and open file
			if err := validateFilePath(filePath); err != nil {
				cliutil.PrintError(fmt.Sprintf("Cannot access import file '%s': %v", filePath, err))
				return errors.New("file access failed")
			}

			file, err := os.Open(filePath) // #nosec G304 -- User-provided import file path is intentional
			if err != nil {
				cliutil.PrintError(fmt.Sprintf("Error opening import file: %v", err))
				return errors.New("file open failed")
			}
			defer file.Close()

			code := HandleImport(ImportParams{
				Reader:  file,
				Service: svc,
				Out:     cobraCmd.OutOrStdout(),
				ErrOut:  cobraCmd.ErrOrStderr(),
			})
			if code != 0 {
				return errors.New("import failed")
			}
			return nil
		},
	}

	importCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to CSV file to import")
	if err := importCmd.MarkFlagRequired("file"); err != nil {
		return nil
	}

	return importCmd
}

// HandleImport executes the import logic from an io.Reader.
// It parses the CSV, imports all rows, and displays results including
// burst suggestions and fact extraction summaries.
//
// Expected:
//   - params.Reader must be a valid io.Reader with CSV data.
//   - params.Service must be a valid career Service instance.
//   - params.Out and params.ErrOut must be valid io.Writer instances.
//   - opts are optional Bubble Tea program options for testing.
//
// Returns:
//   - 0 on successful import (even if some rows were skipped).
//   - 1 on error (parsing or import failure).
//
// Side effects:
//   - Reads from params.Reader.
//   - Writes to params.Out and params.ErrOut.
//   - Modifies the database via params.Service.
func HandleImport(params ImportParams, opts ...tea.ProgramOption) int {
	if params.Reader == nil {
		cliutil.PrintError("Reader cannot be nil")
		return 1
	}

	ctx := context.Background()
	importService := importer.NewImportService(params.Service)

	parsedRows, err := parseCSV(ctx, params.Reader, importService, opts...)
	if err != nil {
		return 1
	}

	if len(parsedRows) == 0 {
		cliutil.PrintError("No valid rows found in CSV file")
		return 1
	}

	cliutil.PrintInfo(fmt.Sprintf("Found %d rows to import", len(parsedRows)))

	result, err := importAllRows(ctx, importService, parsedRows, opts...)
	if err != nil {
		cliutil.PrintError(fmt.Sprintf("Error during import: %v", err))
		return 1
	}

	displayImportResults(params.Out, result)
	return determineExitCode(result)
}

func validateFilePath(filePath string) error {
	_, err := os.Stat(filePath)
	return err
}

func parseCSV(
	ctx context.Context,
	reader io.Reader,
	importService *importer.ImportService,
	opts ...tea.ProgramOption,
) ([]*importer.ParsedRow, error) {
	var parsedRows []*importer.ParsedRow
	err := cliutil.RunWithSpinner("Parsing CSV...", func() error {
		var err error
		parsedRows, err = importService.PrepareImport(ctx, reader)
		return err
	}, opts...)
	if err != nil {
		cliutil.PrintError(fmt.Sprintf("Error parsing CSV: %v", err))
		return nil, err
	}

	return parsedRows, nil
}

func importAllRows(
	ctx context.Context,
	importService *importer.ImportService,
	parsedRows []*importer.ParsedRow,
	opts ...tea.ProgramOption,
) (*importer.ImportResult, error) {
	selectedRows := make([]int, len(parsedRows))
	for i := range parsedRows {
		selectedRows[i] = parsedRows[i].RowNumber
	}

	var result *importer.ImportResult
	err := cliutil.RunWithSpinner(fmt.Sprintf("Importing %d rows...", len(parsedRows)), func() error {
		var err error
		result, err = importService.ImportRows(ctx, parsedRows, selectedRows)
		return err
	}, opts...)
	return result, err
}

func displayImportResults(out io.Writer, result *importer.ImportResult) {
	fmt.Fprintf(out, "\nImport Complete ===\n")
	fmt.Fprintf(out, "Total rows processed: %d\n", result.TotalRows)
	fmt.Fprintf(out, "Successfully imported: %d\n", result.SuccessCount)
	fmt.Fprintf(out, "Skipped: %d\n", result.SkippedCount)

	displayFailedRows(out, result)
	displayBurstSuggestions(out, result)
	displayFactExtraction(out, result)
}

func displayFailedRows(out io.Writer, result *importer.ImportResult) {
	if result.FailedCount == 0 {
		return
	}

	fmt.Fprintf(out, "Failed: %d\n", result.FailedCount)
	if len(result.FailedRows) == 0 {
		return
	}

	fmt.Fprintf(out, "\nFailed rows (first 5):\n")
	maxDisplay := 5
	if len(result.FailedRows) < maxDisplay {
		maxDisplay = len(result.FailedRows)
	}

	for i := range maxDisplay {
		if result.FailedRows[i].Event != nil {
			fmt.Fprintf(out, "  - %s\n", result.FailedRows[i].Event.Text)
		}
	}
}

func displayBurstSuggestions(out io.Writer, result *importer.ImportResult) {
	if len(result.BurstSuggestions) == 0 {
		return
	}

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

func displayFactExtraction(out io.Writer, result *importer.ImportResult) {
	if result.ExtractedFactsCount == 0 {
		return
	}

	fmt.Fprintf(out, "\n=== Fact Extraction ===\n")
	fmt.Fprintf(out, "Extracted %d facts from %d events\n", result.ExtractedFactsCount, len(result.CreatedEvents))

	if len(result.FactsByCompetency) == 0 {
		return
	}

	fmt.Fprintf(out, "\nCompetency breakdown:\n")
	competencies := make([]string, 0, len(result.FactsByCompetency))
	for c := range result.FactsByCompetency {
		competencies = append(competencies, c)
	}
	sort.Strings(competencies)

	for _, c := range competencies {
		fmt.Fprintf(out, "  - %s: %d facts\n", c, result.FactsByCompetency[c])
	}
	fmt.Fprintf(out, "\nThese facts highlight key competencies and achievements.\n")
}

func determineExitCode(result *importer.ImportResult) int {
	if result.SuccessCount > 0 {
		cliutil.PrintSuccess("Import successful!")
		return 0
	}

	if result.SkippedCount > 0 {
		cliutil.PrintInfo("All rows were skipped (possibly duplicates).")
		return 0
	}

	cliutil.PrintError("Import failed - no rows were imported.")
	return 1
}
