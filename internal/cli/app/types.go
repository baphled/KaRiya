package app

import (
	"context"

	"github.com/baphled/kariya/internal/cli/intents"
	configscreens "github.com/baphled/kariya/internal/cli/screens/configure"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/display"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/logger"
	careerservice "github.com/baphled/kariya/internal/service/career"
	cv "github.com/baphled/kariya/internal/service/career/cv"
)

// State identifies which top-level phase the application is in, controlling
// whether the root Model renders the main menu or delegates to an active intent.
type State string

const (
	// StateMenu is the landing state where the user sees the application logo and a
	// selectable list of available intents such as Browse Timeline, Capture Event, and
	// Manage Skills. Keyboard input is handled by the menu navigation logic.
	StateMenu State = "menu"
	// StateIntent is the active state entered after the user selects a menu item. All
	// input and rendering are forwarded to the intent returned by the IntentRouter. The
	// application returns to StateMenu when the intent completes or the user presses Escape.
	StateIntent State = "intent"
)

// MenuItem holds the display label, intent identifier, and help text for a
// single entry in the main menu. The IntentRouter uses the Intent field to
// resolve the corresponding workflow when the user confirms their selection.
type MenuItem struct {
	Name   string
	Intent string
	Help   string
}

// Model is the root Bubble Tea model for the KaRiya application. It owns the
// top-level state machine, service dependencies, and visual chrome such as
// the logo and help modal. All keyboard and resize events flow through its
// Update method before being dispatched to either the menu or the active intent.
type Model struct {
	cliService      *service.CLIEventService
	careerService   *careerservice.Service
	logger          *logger.Logger
	intentRouter    *intents.DefaultIntentRouter
	configManager   cv.ConfigManager
	cvGenService    cv.CVGenerationService
	cvExportService *cv.ExportService

	theme themes.Theme

	state       State
	width       int
	height      int
	showingHelp bool

	selectedMenuIndex int
	menuItems         []MenuItem
	logo              *display.Logo

	terminalInfo *terminal.Info

	ctx context.Context

	infoModal   *feedback.InfoModal
	configModal *configscreens.SettingsModal

	initialScreen      Screen
	initialCaptureMode string
}
