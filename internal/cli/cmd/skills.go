package cmd

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/service/career/technology"
)

// handleRecategorizeSkills recategorizes all skills based on keyword matching.
func handleRecategorizeSkills(svc *careerservice.Service, out io.Writer, errOut io.Writer) int {
	ctx := context.Background()

	skillRepo := svc.GetSkillRepository()
	if skillRepo == nil {
		fmt.Fprintf(errOut, "Error: Skill repository not configured\n")
		return 1
	}

	result, err := technology.RecategorizeSkills(ctx, skillRepo)
	if err != nil {
		fmt.Fprintf(errOut, "Error recategorizing skills: %v\n", err)
		return 1
	}

	fmt.Fprintf(out, "Recategorized %d skills", result.Updated)
	if len(result.ByCategory) > 0 {
		categories := make([]string, 0, len(result.ByCategory))
		for cat := range result.ByCategory {
			categories = append(categories, cat)
		}
		sort.Strings(categories)

		parts := make([]string, 0, len(categories))
		for _, cat := range categories {
			parts = append(parts, fmt.Sprintf("%d %s", result.ByCategory[cat], cat))
		}
		fmt.Fprintf(out, " (%s)", strings.Join(parts, ", "))
	}
	fmt.Fprintf(out, "\n")

	if result.NoMatch > 0 {
		fmt.Fprintf(out, "Skipped %d skills with no keyword match\n", result.NoMatch)
	}
	if result.Skipped > 0 {
		fmt.Fprintf(out, "Skipped %d skills already correctly categorized\n", result.Skipped)
	}

	return 0
}
