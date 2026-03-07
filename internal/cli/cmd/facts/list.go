package facts

import (
	"context"
	"fmt"
	"io"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
)

// ListFacts displays all existing facts.
//
// Expected:
//   - svc: Career service instance providing fact repository access
//   - out: Writer for formatted output
//   - _: Error writer (unused)
//   - _: Bubble Tea program options (unused)
//
// Returns:
//   - Exit code: 0 on success, 1 on error or repository not configured
//
// Side effects:
//   - Queries fact repository via context
//   - Writes error/info messages to stderr via cliutil
//   - Writes formatted fact list to out writer
func ListFacts(svc *careerservice.Service, out io.Writer, _ io.Writer, _ ...tea.ProgramOption) int {
	ctx := context.Background()

	factRepo := svc.GetFactRepository()
	if factRepo == nil {
		cliutil.PrintError("Fact repository not configured")
		return 1
	}

	facts, err := factRepo.List(ctx, career.FactListFilters{Limit: 10000})
	if err != nil {
		cliutil.PrintError(fmt.Sprintf("Error retrieving facts: %v", err))
		return 1
	}

	if len(facts) == 0 {
		cliutil.PrintInfo("No facts found in database.")
		return 0
	}

	DisplayFactList(out, facts)
	return 0
}

// DisplayFactList writes formatted fact list to the output.
//
// Expected:
//   - out: Writer to output formatted facts
//   - facts: Slice of domain.Fact pointers to display
//
// Returns:
//   - None
//
// Side effects:
//   - Writes formatted output to provided writer
func DisplayFactList(out io.Writer, facts []*domain.Fact) {
	output := FormatFactList(facts)
	fmt.Fprint(out, output)
}
