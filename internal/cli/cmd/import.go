package cmd

import (
	"context"
	"fmt"
	"io"
	"sort"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/cli/importer"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
)

// ImportOptions defines configuration for the import command.
type ImportOptions struct {
	FilePath string
	Service  *careerservice.Service
	Out      io.Writer
	ErrOut   io.Writer
	Opener   cliutil.FileOpener
	Runner   cliutil.ProgressRunner
	TeaOpts  []tea.ProgramOption
}

// HandleImport executes the import logic for a CSV file.
//
// Deprecated: Use runImport with ImportOptions for new code.
// Maintained for backward compatibility with existing command signature.
//
//nolint:revive // Maintain backward compatibility with existing command signature
func HandleImport(
	filePath string,
	_ bool,
	_ bool,
	svc *careerservice.Service,
	out io.Writer,
	errOut io.Writer,
	opener cliutil.FileOpener,
	runner cliutil.ProgressRunner,
	opts ...tea.ProgramOption,
) int {
	return runImport(ImportOptions{
		FilePath: filePath,
		Service:  svc,
		Out:      out,
		ErrOut:   errOut,
		Opener:   opener,
		Runner:   runner,
		TeaOpts:  opts,
	})
}

// runImport executes the import logic.
//
// This function handles parsing, execution, and reporting for the import command.
//
//nolint:funlen // Comprehensive import logic including parsing, execution, and reporting
func runImport(opts ImportOptions) int {
	if _, err := opts.Opener.Stat(opts.FilePath); err != nil {
		fmt.Fprintf(opts.ErrOut, "Error: Cannot access import file '%s': %v\n", opts.FilePath, err)
		return 1
	}

	ctx := context.Background()
	importService := importer.NewImportService(opts.Service)

	file, err := opts.Opener.Open(opts.FilePath)
	if err != nil {
		fmt.Fprintf(opts.ErrOut, "Error opening import file: %v\n", err)
		return 1
	}
	defer file.Close()

	var parsedRows []*importer.ParsedRow
	msg := "Parsing CSV file: " + opts.FilePath
	err = opts.Runner.RunWithSpinner(msg, func() error {
		var err error
		parsedRows, err = importService.PrepareImport(ctx, file)
		if err != nil {
			return fmt.Errorf("import operation failed: %w", err)
		}
		return nil
	}, opts.TeaOpts...)

	if err != nil {
		fmt.Fprintf(opts.ErrOut, "Error parsing CSV: %v\n", err)
		return 1
	}

	if len(parsedRows) == 0 {
		fmt.Fprintf(opts.ErrOut, "Error: No valid rows found in CSV file\n")
		return 1
	}

	selectedRows := make([]int, len(parsedRows))
	for i, row := range parsedRows {
		selectedRows[i] = row.RowNumber
	}

	var result *importer.ImportResult
	msg = fmt.Sprintf("Importing %d rows...", len(parsedRows))
	err = opts.Runner.RunWithSpinner(msg, func() error {
		var err error
		result, err = importService.ImportRows(ctx, parsedRows, selectedRows)
		if err != nil {
			return fmt.Errorf("import operation failed: %w", err)
		}
		return nil
	}, opts.TeaOpts...)

	if err != nil {
		fmt.Fprintf(opts.ErrOut, "Error during import: %v\n", err)
		return 1
	}

	DisplayImportResults(result, opts.Out)

	if result.SuccessCount > 0 {
		fmt.Fprintf(opts.Out, "\n✓ Import successful!\n")
		return 0
	} else if result.SkippedCount > 0 {
		fmt.Fprintf(opts.Out, "\n⚠ All rows were skipped (possibly duplicates).\n")
		return 0
	} else {
		fmt.Fprintf(opts.ErrOut, "\n✗ Import failed - no rows were imported.\n")
		return 1
	}
}

// DisplayImportResults displays the results of an import operation.
// DisplayImportResults prints a summary of the import, including success,
// skipped, and failed counts, as well as burst suggestions and fact extraction
// results.
func DisplayImportResults(result *importer.ImportResult, out io.Writer) {
	fmt.Fprintf(out, "\nImport Complete ===\n")
	fmt.Fprintf(out, "Total rows processed: %d\n", result.TotalRows)
	fmt.Fprintf(out, "Successfully imported: %d\n", result.SuccessCount)
	fmt.Fprintf(out, "Skipped: %d\n", result.SkippedCount)
	if result.FailedCount > 0 {
		fmt.Fprintf(out, "Failed: %d\n", result.FailedCount)
	}

	if len(result.BurstSuggestions) > 0 {
		fmt.Fprintf(out, "\n=== Burst Suggestions ===\n")
		fmt.Fprintf(out, "Detected %d potential bursts from imported events:\n\n", len(result.BurstSuggestions))
		for i, burst := range result.BurstSuggestions {
			fmt.Fprintf(out, "%d. %s\n", i+1, burst.Name)
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
}
