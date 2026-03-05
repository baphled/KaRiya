package cmd

import (
	"github.com/spf13/cobra"
)

var version = "dev"

// NewRootCmd creates the root command for the KaRiya CLI
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kariya",
		Short: "Career Event Capture for Engineers",
		Long: `KaRiya is a terminal user interface for capturing career achievements 
and generating tailored CVs. Built with Go and Bubble Tea, it transforms your 
career events into professional, role-specific CVs directly from your terminal.`,
		Version:      version,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Will launch TUI in future - for now return help
			return cmd.Help()
		},
	}

	// Override version template to match current output format
	cmd.SetVersionTemplate("kariya version {{.Version}}\n")

	return cmd
}
