package cmd

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"

	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/cli/bootstrap"
	"github.com/baphled/kariya/internal/cli/cmd/bursts"
	"github.com/baphled/kariya/internal/cli/cmd/facts"
	importcmd "github.com/baphled/kariya/internal/cli/cmd/import"
	"github.com/baphled/kariya/internal/cli/cmd/skills"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/logger"
)

// NewRootCmd creates the root command for the KaRiya CLI.
func NewRootCmd(version string) *cobra.Command {
	ctx := &CLIContext{}

	cmd := &cobra.Command{
		Use:   "kariya",
		Short: "Career Event Capture for Engineers",
		Long: `KaRiya is a terminal user interface for capturing career achievements 
and generating tailored CVs. Built with Go and Bubble Tea, it transforms your 
career events into professional, role-specific CVs directly from your terminal.`,
		Version:      version,
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// Initialize service for subcommands
			dbPath, err := cmd.Flags().GetString("db")
			if err != nil {
				return err
			}
			inMemory, err := cmd.Flags().GetBool("in-memory")
			if err != nil {
				return err
			}
			ctx.dbPath = dbPath
			ctx.inMemory = inMemory
			return ctx.InitService(cmd.ErrOrStderr())
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Check if both stdin and stdout are terminals; if not, show help
			// This prevents launching TUI when output is piped or in non-interactive shells
			stdinTTY := isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd())
			stdoutTTY := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

			if !stdinTTY || !stdoutTTY {
				return cmd.Help()
			}

			// Launch TUI when no subcommand specified
			careerService := ctx.Service()

			bootstrapResult, err := bootstrap.Run(careerService, logger.DefaultLogger())
			if err != nil {
				return err
			}

			cliService := service.NewCLIEventService(careerService)
			model := app.NewModel(cliService, careerService, bootstrapResult)

			p := tea.NewProgram(model, tea.WithAltScreen())
			_, err = p.Run()
			return err
		},
	}

	// Override version template to match current output format
	cmd.SetVersionTemplate("kariya version {{.Version}}\n")

	// Persistent flags (available to all subcommands)
	cmd.PersistentFlags().String("db", "", "Database path (default: ~/.kariya/kariya.db)")
	cmd.PersistentFlags().Bool("in-memory", false, "Use in-memory database for testing")

	// Add subcommands
	cmd.AddCommand(bursts.NewBurstsCmd(ctx))
	cmd.AddCommand(facts.NewFactsCmd(ctx))
	cmd.AddCommand(skills.NewSkillsCmd(ctx))
	cmd.AddCommand(importcmd.NewImportCmd(ctx))

	return cmd
}
