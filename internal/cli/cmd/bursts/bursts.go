// Package bursts provides Cobra CLI commands for burst detection and listing.
package bursts

import (
	"errors"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/spf13/cobra"
)

// NewBurstsCmd creates the parent "bursts" command with subcommands.
//
// Expected:
//   - ctx: ServiceProvider with initialized Service
//
// Returns:
//   - Configured cobra.Command with detect and list subcommands
//
// Side effects:
//   - Adds detect and list subcommands to the parent command
func NewBurstsCmd(ctx cliutil.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bursts",
		Short: "Manage bursts",
		Long:  "Commands for detecting and managing burst patterns in events",
	}

	cmd.AddCommand(NewDetectCmd(ctx))
	cmd.AddCommand(NewListCmd(ctx))

	return cmd
}

// NewDetectCmd creates the "detect" subcommand.
//
// Expected:
//   - ctx: ServiceProvider with initialized Service
//
// Returns:
//   - Configured cobra.Command for burst detection
//
// Side effects:
//   - None (command configuration only)
func NewDetectCmd(ctx cliutil.ServiceProvider) *cobra.Command {
	return &cobra.Command{
		Use:   "detect",
		Short: "Detect burst patterns in events",
		Long:  "Analyze events and detect burst patterns, then save them to the repository",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc := ctx.Service()
			if svc == nil {
				return errors.New("service not initialized")
			}
			code := DetectBursts(svc, cmd.OutOrStdout(), cmd.ErrOrStderr())
			if code != 0 {
				return errors.New("burst detection failed")
			}
			return nil
		},
	}
}

// NewListCmd creates the "list" subcommand.
//
// Expected:
//   - ctx: ServiceProvider with initialized Service
//
// Returns:
//   - Configured cobra.Command for listing bursts
//
// Side effects:
//   - None (command configuration only)
func NewListCmd(ctx cliutil.ServiceProvider) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all bursts",
		Long:  "Display all existing bursts from the repository",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc := ctx.Service()
			if svc == nil {
				return errors.New("service not initialized")
			}
			code := ListBursts(svc, cmd.OutOrStdout(), cmd.ErrOrStderr())
			if code != 0 {
				return errors.New("listing bursts failed")
			}
			return nil
		},
	}
}
