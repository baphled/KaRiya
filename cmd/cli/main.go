package main

import (
	"fmt"
	"io"
	"os"

	"github.com/baphled/kariya/internal/cli/app"
	cliservice "github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	version = "0.1.0"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout))
}

func run(args []string, out io.Writer) int {
	// Parse CLI flags
	showVersion := false
	showHelp := false

	for _, arg := range args {
		switch arg {
		case "--version", "-v":
			showVersion = true
		case "--help", "-h":
			showHelp = true
		}
	}

	// Handle version flag
	if showVersion {
		fmt.Fprintf(out, "KaRiya CLI v%s\n", version)
		return 0
	}

	// Handle help flag
	if showHelp {
		printHelpTo(out)
		return 0
	}

	// Set up dependencies
	// TODO: Replace with proper dependency injection
	repo := career.NewMemoryRepository()
	svc := careerservice.NewService(repo)
	cliSvc := cliservice.NewCLIEventService(svc)

	// Initialize application model
	model := app.NewModel(cliSvc, svc)

	// Initialize BubbleTea program
	p := tea.NewProgram(model)
	if err := p.Start(); err != nil {
		fmt.Fprintf(out, "Error running program: %v\n", err)
		return 1
	}
	return 0
}

func printHelpTo(out io.Writer) {
	fmt.Fprintln(out, "KaRiya CLI - Career Journaling Tool")
	fmt.Fprintln(out, "\nUsage: kariya [options]")
	fmt.Fprintln(out, "\nOptions:")
	fmt.Fprintln(out, "  -v, --version   Show version information")
	fmt.Fprintln(out, "  -h, --help      Show this help message")
	fmt.Fprintln(out, "\nCommands:")
	fmt.Fprintln(out, "  capture         Start a new career event capture")
	fmt.Fprintln(out, "  list            List recent career events")
}
