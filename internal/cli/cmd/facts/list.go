package facts

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// ListFacts displays all existing facts.
func ListFacts(svc *careerservice.Service, out io.Writer, _ io.Writer) int {
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

	th := theme.Default()
	fmt.Fprintln(out, primitives.Title("Existing Facts", th).Render())
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
