package cmd

import (
	"fmt"
	"io"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// CommandConfig defines the configuration for a command that uses TUI progress.
type CommandConfig struct {
	Name    string
	Short   string
	Long    string
	Action  func(svc *careerservice.Service, out, errOut io.Writer, runner cliutil.ProgressRunner, opts ...tea.ProgramOption) int
	FailMsg string
}

// NewTUICommand creates a new cobra command with standardized TUI options.
func NewTUICommand(ctx *CLIContext, cfg CommandConfig) *cobra.Command {
	return &cobra.Command{
		Use:   cfg.Name,
		Short: cfg.Short,
		Long:  cfg.Long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			runner := cliutil.DefaultProgressRunner{}
			var opts []tea.ProgramOption
			opts = append(opts, tea.WithInput(nil))
			if cmd.OutOrStdout() == io.Discard {
				opts = append(opts, tea.WithOutput(io.Discard))
			}

			if code := cfg.Action(ctx.Service(), cmd.OutOrStdout(), cmd.ErrOrStderr(), runner, opts...); code != 0 {
				return fmt.Errorf("%s", cfg.FailMsg)
			}
			return nil
		},
	}
}

// ListConfig defines the configuration for a simple list command.
type ListConfig struct {
	Name    string
	Short   string
	Long    string
	Action  func(svc *careerservice.Service, out, errOut io.Writer) int
	FailMsg string
}

// NewListCommand creates a new cobra command for listing resources.
func NewListCommand(ctx *CLIContext, cfg ListConfig) *cobra.Command {
	return &cobra.Command{
		Use:   cfg.Name,
		Short: cfg.Short,
		Long:  cfg.Long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if code := cfg.Action(ctx.Service(), cmd.OutOrStdout(), cmd.ErrOrStderr()); code != 0 {
				return fmt.Errorf("%s", cfg.FailMsg)
			}
			return nil
		},
	}
}
