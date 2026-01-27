package app

import (
	"context"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/display"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/logger"
	careerservice "github.com/baphled/kariya/internal/service/career"
	cv "github.com/baphled/kariya/internal/service/career/cv"
)

// AppState represents the current state of the application.
type AppState string

const (
	StateMenu   AppState = "menu"
	StateIntent AppState = "intent"
)

// MenuItem represents a menu option.
type MenuItem struct {
	Name   string
	Intent string
	Help   string
}

// Model is the root Bubble Tea model for the KaRiya application.
type Model struct {
	// Core services.
	cliService      *service.CLIEventService
	careerService   *careerservice.Service
	logger          *logger.Logger
	intentRouter    *intents.DefaultIntentRouter
	configManager   cv.ConfigManager
	cvGenService    cv.CVGenerationService
	cvExportService *cv.ExportService

	// Theme for consistent styling.
	theme themes.Theme

	// UI state.
	state       AppState
	width       int
	height      int
	showingHelp bool

	// Menu state.
	selectedMenuIndex int
	menuItems         []MenuItem
	logo              *display.Logo

	// Terminal info for responsive rendering.
	terminalInfo *terminal.Info

	// Context for intent creation.
	ctx context.Context

	// Info modal for blocking user feedback (empty state warnings).
	infoModal *feedback.InfoModal

	// Initial navigation settings (set before Init).
	initialScreen      Screen
	initialCaptureMode string
}
