package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/cli/bootstrap"
	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	cliservice "github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/logger"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command with all subcommands.
func NewRootCmd(version string) *cobra.Command {
	ctx := &CLIContext{}

	cmd := &cobra.Command{
		Use:          "kariya",
		Short:        "KaRiya - Career Event Capture for Engineers",
		Long:         "A terminal user interface for capturing career achievements and generating tailored CVs",
		Version:      version,
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			dbPath, err := cmd.Flags().GetString("db")
			if err != nil {
				return err
			}
			inMemory, err := cmd.Flags().GetBool("in-memory")
			if err != nil {
				return err
			}
			ctx.DBPath = dbPath
			ctx.InMemory = inMemory
			return ctx.InitService()
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			mode, err := cmd.Flags().GetString("mode")
			if err != nil {
				return err
			}
			return runInteractive(ctx.Service, mode, version, cmd.ErrOrStderr())
		},
	}

	cmd.PersistentFlags().String("db", "", "Path to SQLite database (default: ~/.kariya/events.db)")
	cmd.PersistentFlags().String("mode", "", "Start in specific capture mode (timeline, backfill, manual)")
	cmd.PersistentFlags().Bool("in-memory", false, "Use in-memory storage (data not persisted)")

	// Add subcommands
	cmd.AddCommand(NewImportCmd(ctx))
	cmd.AddCommand(NewBurstsCmd(ctx))
	cmd.AddCommand(NewFactsCmd(ctx))
	cmd.AddCommand(newSkillsCmd(ctx))

	return cmd
}

// runInteractive runs the interactive TUI application.
func runInteractive(svc *careerservice.Service, mode, version string, errOut io.Writer) error {
	cliSvc := cliservice.NewCLIEventService(svc)

	log := logger.DefaultLogger()
	bootstrapResult, err := bootstrap.Run(svc, log)
	if err != nil {
		if errors.Is(err, bootstrap.ErrUserAborted) {
			return nil
		}
		fmt.Fprintf(errOut, "Error during bootstrap: %v\n", err)
		return err
	}

	displayVersion := version
	if displayVersion == "" {
		displayVersion = "dev"
	}

	model := app.NewModel(cliSvc, svc, bootstrapResult, app.WithVersion(displayVersion))

	if mode != "" {
		model.SetInitialCaptureMode(mode)
	}

	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(errOut, "Error running program: %v\n", err)
		return err
	}
	return nil
}

// NewImportCmd creates the import command.
func NewImportCmd(ctx *CLIContext) *cobra.Command {
	return &cobra.Command{
		Use:   "import <file>",
		Short: "Import career events from CSV file",
		Long:  "Import career events from a CSV file with automatic burst and fact detection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opener := cliutil.OSFileOpener{}
			runner := cliutil.DefaultProgressRunner{}

			var opts []tea.ProgramOption
			if cmd.OutOrStdout() == io.Discard || cmd.OutOrStdout() != nil {
				// In many test environments, we want to avoid TTY issues
				opts = append(opts, tea.WithInput(nil))
				if cmd.OutOrStdout() == io.Discard {
					opts = append(opts, tea.WithOutput(io.Discard))
				}
			}

			if code := HandleImport(
				args[0], false, false, ctx.Service, cmd.OutOrStdout(), cmd.ErrOrStderr(), opener, runner, opts...,
			); code != 0 {
				return errors.New("import failed")
			}
			return nil
		},
	}
}

// NewBurstsCmd creates the bursts command with subcommands.
func NewBurstsCmd(ctx *CLIContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bursts",
		Short: "Manage bursts",
		Long:  "Commands for detecting and managing burst groupings",
	}

	cmd.AddCommand(NewTUICommand(ctx, CommandConfig{
		Name:    "detect",
		Short:   "Detect bursts from events",
		Long:    "Re-run burst detection on all events",
		Action:  DetectBursts,
		FailMsg: "burst detection failed",
	}))

	cmd.AddCommand(NewListCommand(ctx, ListConfig{
		Name:    "list",
		Short:   "List all bursts",
		Long:    "Display all existing bursts",
		Action:  ListBursts,
		FailMsg: "listing bursts failed",
	}))

	return cmd
}

// NewFactsCmd creates the facts command with subcommands.
func NewFactsCmd(ctx *CLIContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "facts",
		Short: "Manage facts",
		Long:  "Commands for extracting and managing facts",
	}

	cmd.AddCommand(NewTUICommand(ctx, CommandConfig{
		Name:    "extract",
		Short:   "Extract facts from events",
		Long:    "Re-run fact extraction on all events",
		Action:  ExtractFacts,
		FailMsg: "fact extraction failed",
	}))

	cmd.AddCommand(NewListCommand(ctx, ListConfig{
		Name:    "list",
		Short:   "List all facts",
		Long:    "Display all existing facts",
		Action:  ListFacts,
		FailMsg: "listing facts failed",
	}))

	return cmd
}

// newSkillsCmd creates the skills command with subcommands.
func newSkillsCmd(ctx *CLIContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skills",
		Short: "Manage skills",
		Long:  "Commands for managing skills and their categorization",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "recategorize",
		Short: "Recategorize all skills",
		Long:  "Re-run skill categorization based on keyword matching",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if code := handleRecategorizeSkills(ctx.Service, cmd.OutOrStdout(), cmd.ErrOrStderr()); code != 0 {
				return errors.New("recategorization failed")
			}
			return nil
		},
	})

	return cmd
}
