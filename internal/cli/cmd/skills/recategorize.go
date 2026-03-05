// Package skills provides CLI commands for managing skills and their categorization.
package skills

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/service/career/technology"
	tea "github.com/charmbracelet/bubbletea"
)

// RecategorizeSkills recategorizes all skills based on keyword matching.
func RecategorizeSkills(svc *careerservice.Service, _ io.Writer, _ io.Writer, opts ...tea.ProgramOption) int {
	ctx := context.Background()

	skillRepo := svc.GetSkillRepository()
	if skillRepo == nil {
		cliutil.PrintError("Error: Skill repository not configured")
		return 1
	}

	var result *technology.RecategorizeResult
	err := cliutil.RunWithSpinner("Recategorizing skills...", func() error {
		var err error
		result, err = technology.RecategorizeSkills(ctx, skillRepo)
		return err
	}, opts...)

	if err != nil {
		cliutil.PrintError(fmt.Sprintf("Error recategorizing skills: %v", err))
		return 1
	}

	cliutil.PrintSuccess(fmt.Sprintf("Recategorized %d skills", result.Updated))
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
		cliutil.PrintInfo("Breakdown: " + strings.Join(parts, ", "))
	}

	if result.NoMatch > 0 {
		cliutil.PrintInfo(fmt.Sprintf("Skipped %d skills with no keyword match", result.NoMatch))
	}
	if result.Skipped > 0 {
		cliutil.PrintInfo(fmt.Sprintf("Skipped %d skills already correctly categorized", result.Skipped))
	}

	return 0
}
