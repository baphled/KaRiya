package bursts

import (
	"context"
	"fmt"
	"io"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// ListBursts displays all existing bursts.
func ListBursts(svc *careerservice.Service, out io.Writer, _ io.Writer) int {
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

	th := theme.Default()
	fmt.Fprintln(out, primitives.Title("Bursts", th).Render())

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
