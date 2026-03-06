package facts

import (
	"errors"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/spf13/cobra"
)

// NewFactsCmd creates the parent "facts" command with subcommands.
//
// Expected:
//   - ctx: ServiceContext with initialized Service
//
// Returns:
//   - Configured cobra.Command with extract and list subcommands
//
// Side effects:
//   - Adds extract and list subcommands to the parent command
func NewFactsCmd(ctx cliutil.ServiceContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "facts",
		Short: "Manage facts",
		Long:  "Commands for extracting and managing facts from events",
	}

	cmd.AddCommand(NewExtractCmd(ctx))
	cmd.AddCommand(NewListCmd(ctx))

	return cmd
}

// NewExtractCmd creates the "extract" subcommand.
//
// Expected:
//   - ctx: ServiceContext with initialized Service
//
// Returns:
//   - Configured cobra.Command for fact extraction
//
// Side effects:
//   - None (command configuration only)
func NewExtractCmd(ctx cliutil.ServiceContext) *cobra.Command {
	return &cobra.Command{
		Use:   "extract",
		Short: "Extract facts from events",
		Long:  "Analyze events and extract facts, then save them to the repository",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc := ctx.Service()
			if svc == nil {
				return errors.New("service not initialized")
			}
			code := ExtractFacts(svc, cmd.OutOrStdout(), cmd.ErrOrStderr())
			if code != 0 {
				return errors.New("fact extraction failed")
			}
			return nil
		},
	}
}

// NewListCmd creates the "list" subcommand.
//
// Expected:
//   - ctx: ServiceContext with initialized Service
//
// Returns:
//   - Configured cobra.Command for listing facts
//
// Side effects:
//   - None (command configuration only)
func NewListCmd(ctx cliutil.ServiceContext) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all facts",
		Long:  "Display all existing facts from the repository",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc := ctx.Service()
			if svc == nil {
				return errors.New("service not initialized")
			}
			code := ListFacts(svc, cmd.OutOrStdout(), cmd.ErrOrStderr())
			if code != 0 {
				return errors.New("listing facts failed")
			}
			return nil
		},
	}
}
