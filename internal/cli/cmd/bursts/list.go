package bursts

import (
	"context"
	"fmt"
	"io"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
)

// ListBursts displays all existing bursts.
//
// Expected:
//   - svc: Career service instance providing burst repository access
//   - out: Writer for formatted output
//   - _: Error writer (unused)
//   - _: Bubble Tea program options (unused)
//
// Returns:
//   - Exit code: 0 on success, 1 on error or repository not configured
//
// Side effects:
//   - Queries burst repository via context
//   - Writes error/info messages to stderr via cliutil
//   - Writes formatted burst list to out writer
func ListBursts(svc *careerservice.Service, out io.Writer, _ io.Writer, _ ...tea.ProgramOption) int {
	ctx := context.Background()

	burstRepo := svc.GetBurstRepository()
	if burstRepo == nil {
		cliutil.PrintError("Burst repository not configured")
		return 1
	}

	bursts, err := burstRepo.List(ctx, career.BurstListFilters{Limit: 10000})
	if err != nil {
		cliutil.PrintError(fmt.Sprintf("Error retrieving bursts: %v", err))
		return 1
	}

	if len(bursts) == 0 {
		cliutil.PrintInfo("No bursts found.")
		return 0
	}

	output := FormatBurstList(bursts)
	fmt.Fprint(out, output)

	return 0
}
