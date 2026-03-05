package facts

import (
	"errors"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/spf13/cobra"
)

// NewFactsCmd creates the parent "facts" command with subcommands.
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
