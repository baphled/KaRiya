// Package bursts provides Cobra CLI commands for burst detection and listing.
package bursts

import (
	"errors"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/spf13/cobra"
)

// NewBurstsCmd creates the parent "bursts" command with subcommands.
func NewBurstsCmd(ctx cliutil.ServiceContext) *cobra.Command {
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
func NewDetectCmd(ctx cliutil.ServiceContext) *cobra.Command {
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
func NewListCmd(ctx cliutil.ServiceContext) *cobra.Command {
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
