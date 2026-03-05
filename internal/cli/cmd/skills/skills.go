package skills

import (
	"errors"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/spf13/cobra"
)

// NewSkillsCmd creates the parent "skills" command with subcommands.
func NewSkillsCmd(ctx cliutil.ServiceContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skills",
		Short: "Manage skills",
		Long:  "Commands for managing skills and their categorization",
	}

	cmd.AddCommand(NewRecategorizeCmd(ctx))

	return cmd
}

// NewRecategorizeCmd creates the "recategorize" subcommand.
func NewRecategorizeCmd(ctx cliutil.ServiceContext) *cobra.Command {
	return &cobra.Command{
		Use:   "recategorize",
		Short: "Recategorize all skills",
		Long:  "Re-run skill categorization based on keyword matching",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc := ctx.Service()
			if svc == nil {
				return errors.New("service not initialized")
			}
			code := RecategorizeSkills(svc, cmd.OutOrStdout(), cmd.ErrOrStderr())
			if code != 0 {
				return errors.New("recategorization failed")
			}
			return nil
		},
	}
}
