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
	var (
		showVersion = false
		showHelp    = false
		dbPath      = ""
		mode        = ""
		listEvents  = false
	)

	// Parse command-line arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--version", "-v":
			showVersion = true
		case "--help", "-h":
			showHelp = true
		case "--db", "--database":
			if i+1 < len(args) {
				dbPath = args[i+1]
				i++ // Skip next argument as it's the value
			}
		case "--mode":
			if i+1 < len(args) {
				mode = args[i+1]
				// Validate mode
				if mode != "timeline" && mode != "backfill" && mode != "manual" {
					fmt.Fprintf(out, "Error: Invalid mode '%s'. Valid modes are: timeline, backfill, manual\n", mode)
					return 1
				}
				i++ // Skip next argument as it's the value
			}
		case "--list":
			listEvents = true
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
	// TODO: Replace with proper dependency injection and SQLite when --db is provided
	_ = dbPath // Currently ignored, would use for SQLite initialization
	_ = mode   // Currently ignored, would set initial capture mode
	_ = listEvents // Currently ignored, would show list screen on startup

	repo := career.NewMemoryRepository()
	svc := careerservice.NewService(repo)
	cliSvc := cliservice.NewCLIEventService(svc)

	// Initialize application model
	model := app.NewModel(cliSvc, svc)

	// Initialize BubbleTea program
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(out, "Error running program: %v\n", err)
		return 1
	}
	return 0
}

func printHelpTo(out io.Writer) {
	fmt.Fprintln(out, "KaRiya CLI - Career Journaling Tool")
	fmt.Fprintln(out, "\nUsage: kariya [options]")
	fmt.Fprintln(out, "\nOptions:")
	fmt.Fprintln(out, "  -v, --version              Show version information")
	fmt.Fprintln(out, "  -h, --help                 Show this help message")
	fmt.Fprintln(out, "  --db, --database PATH      Use custom database path (SQLite)")
	fmt.Fprintln(out, "  --mode MODE                Start in specific capture mode")
	fmt.Fprintln(out, "                             Valid modes: timeline, backfill, manual")
	fmt.Fprintln(out, "  --list                     Show recent events on startup")
	fmt.Fprintln(out, "\nCommands:")
	fmt.Fprintln(out, "  capture                    Start a new career event capture")
	fmt.Fprintln(out, "  list                       List recent career events")
}
