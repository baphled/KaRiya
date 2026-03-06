// Package importcmd provides Cobra CLI command for importing events from CSV files.
package importcmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

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
//
// Expected:
//   - ctx: ServiceContext with initialized Service
//
// Returns:
//   - Configured cobra.Command for importing CSV files
//
// Side effects:
//   - Registers "file" flag on the command
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
			if err := ValidateFilePath(filePath); err != nil {
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

// ValidateFilePath checks if a file exists and is accessible.
//
// Expected:
//   - filePath: Path to file to validate
//
// Returns:
//   - nil if file exists and is accessible
//   - error from os.Stat if file cannot be accessed
//
// Side effects:
//   - Calls os.Stat to check file existence
func ValidateFilePath(filePath string) error {
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

	DisplayFailedRows(out, result)
	DisplayBurstSuggestions(out, result)
	DisplayFactExtraction(out, result)
}

// DisplayFailedRows is a thin wrapper around FormatFailedRows.
//
// Expected:
//   - out: io.Writer to write output to
//   - result: ImportResult with FailedCount and FailedRows
//
// Side effects:
//   - Writes formatted failed rows output to out
func DisplayFailedRows(out io.Writer, result *importer.ImportResult) {
	output := FormatFailedRows(result)
	if output != "" {
		fmt.Fprint(out, output)
	}
}

// DisplayBurstSuggestions is a thin wrapper around FormatBurstSuggestions.
//
// Expected:
//   - out: io.Writer to write output to
//   - result: ImportResult with BurstSuggestions
//
// Side effects:
//   - Writes formatted burst suggestions to out
func DisplayBurstSuggestions(out io.Writer, result *importer.ImportResult) {
	output := FormatBurstSuggestions(result)
	if output != "" {
		fmt.Fprint(out, output)
	}
}

// DisplayFactExtraction is a thin wrapper around FormatFactExtraction.
//
// Expected:
//   - out: io.Writer to write output to
//   - result: ImportResult with ExtractedFactsCount and FactsByCompetency
//
// Side effects:
//   - Writes formatted fact extraction output to out
func DisplayFactExtraction(out io.Writer, result *importer.ImportResult) {
	output := FormatFactExtraction(result)
	if output != "" {
		fmt.Fprint(out, output)
	}
}

func determineExitCode(result *importer.ImportResult) int {
	return CalculateExitCode(result)
}
