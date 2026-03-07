package skills

import (
	"errors"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// NewSkillsCmd creates the parent "skills" command with subcommands.
//
// Expected:
//   - ctx: ServiceProvider with initialized Service
//
// Returns:
//   - Configured cobra.Command with recategorize subcommand
//
// Side effects:
//   - Adds recategorize subcommand to the parent command
func NewSkillsCmd(ctx cliutil.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skills",
		Short: "Manage skills",
		Long:  "Commands for managing skills and their categorization",
	}

	cmd.AddCommand(NewRecategorizeCmd(ctx))

	return cmd
}

// NewRecategorizeCmd creates the "recategorize" subcommand.
//
// Expected:
//   - ctx: ServiceProvider with initialized Service
//   - opts: optional tea.ProgramOption for controlling terminal behavior
//
// Returns:
//   - Configured cobra.Command for skill recategorization
//
// Side effects:
//   - None (command configuration only)
func NewRecategorizeCmd(ctx cliutil.ServiceProvider, opts ...tea.ProgramOption) *cobra.Command {
	return &cobra.Command{
		Use:   "recategorize",
		Short: "Recategorize all skills",
		Long:  "Re-run skill categorization based on keyword matching",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc := ctx.Service()
			if svc == nil {
				return errors.New("service not initialized")
			}
			code := RecategorizeSkills(svc, cmd.OutOrStdout(), cmd.ErrOrStderr(), opts...)
			if code != 0 {
				return errors.New("recategorization failed")
			}
			return nil
		},
	}
}
