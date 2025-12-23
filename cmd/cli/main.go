package main

import (
	"fmt"
	"log"
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
	// Parse CLI flags
	showVersion := false
	showHelp := false

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--version", "-v":
			showVersion = true
		case "--help", "-h":
			showHelp = true
		}
	}

	// Handle version flag
	if showVersion {
		fmt.Printf("KaRiya CLI v%s\n", version)
		os.Exit(0)
	}

	// Handle help flag
	if showHelp {
		printHelp()
		os.Exit(0)
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
		log.Fatalf("Error running program: %v", err)
	}
}

func printHelp() {
	fmt.Println("KaRiya CLI - Career Journaling Tool")
	fmt.Println("\nUsage: kariya [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -v, --version   Show version information")
	fmt.Println("  -h, --help      Show this help message")
	fmt.Println("\nCommands:")
	fmt.Println("  capture         Start a new career event capture")
	fmt.Println("  list            List recent career events")
}
