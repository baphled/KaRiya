package facts

import (
	"errors"
	"io"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
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
			return ExecuteExtractFacts(ctx.Service(), cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}
}

// ExecuteExtractFacts extracts facts from events in the database.
//
// Expected:
//   - svc: initialized career Service
//   - out/errOut: output writers
//
// Returns:
//   - nil on success
//   - error if service is nil or extraction fails
//
// Side effects:
//   - Reads from database
//   - Writes facts to database
//   - Writes output to writers
func ExecuteExtractFacts(svc *careerservice.Service, out, errOut io.Writer, opts ...tea.ProgramOption) error {
	if svc == nil {
		return errors.New("service not initialized")
	}
	code := ExtractFacts(svc, out, errOut, opts...)
	if code != 0 {
		return errors.New("fact extraction failed")
	}
	return nil
}

// ExecuteListFacts lists all facts from the database.
//
// Expected:
//   - svc: initialized career Service
//   - out/errOut: output writers
//
// Returns:
//   - nil on success
//   - error if service is nil or list fails
//
// Side effects:
//   - Reads from database
//   - Writes output to writers
func ExecuteListFacts(svc *careerservice.Service, out, errOut io.Writer, opts ...tea.ProgramOption) error {
	if svc == nil {
		return errors.New("service not initialized")
	}
	code := ListFacts(svc, out, errOut, opts...)
	if code != 0 {
		return errors.New("list facts failed")
	}
	return nil
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
			return ExecuteListFacts(ctx.Service(), cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}
}
