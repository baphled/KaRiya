// Package browsetimeline implements the BrowseTimeline intent for browsing career events.
package browsetimeline

import (
	"context"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/tui/intents"
	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/display"

	eventviews "github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
)

const errorLoadingSkillsTitle = "Error Loading Skills"

// NewIntent creates a new BrowseTimeline intent with the given context,
// initializing all internal state including filters, screens, and modal registry.
//
// Expected:
//   - ctx must be non-nil and pass validation.
//
// Returns:
//   - A fully initialized Intent ready for Init, or an error if validation fails.
//
// Side effects:
//   - Allocates internal collections and a modal registry.
func NewIntent(ctx *IntentValidator) (*Intent, error) {
	if err := ctx.Validate(); err != nil {
		return nil, err
	}

	base := intents.NewBaseIntent()

	intent := &Intent{
		BaseIntent:     base,
		context:        ctx,
		state:          StateTimeline,
		filteredEvents: ctx.Events,
		selectedIndex:  0,
		filters:        ctx.InitialFilters,
		filterStack:    behaviors.NewFilterStack(),
		selectedFacts:  make([]*career.Fact, 0),
		viewedEvents:   make([]*career.Event, 0),
		active:         true,
		modalRegistry:  intents.NewModalRegistry(),
		skillService:   ctx.CLISkillCreator,
	}

	return intent, nil
}

// Init activates the intent by applying initial filters and transitioning
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (i *Intent) Init() tea.Cmd {
	i.applyFilters()

	if len(i.filteredEvents) > 0 {
		i.selectedEvent = i.filteredEvents[0]
	}

	i.state = StateTimeline
	i.transitionToView(eventviews.NewListView(display.EventsFromDomain(i.filteredEvents)))
	return nil
}

// Update is the main message loop that delegates to views, keyboard shortcuts,
// or the active screen in priority order.
//
// Expected:
//   - msg is a Bubble Tea message from the runtime.
//
// Returns:
//   - A command to execute, or nil if the message was fully handled.
//
// Side effects:
//   - May open or close modals, transition screens, or mark the intent inactive.
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	if cmd, handled := i.handleSkillMessages(msg); handled {
		return cmd
	}

	if cmd, handled := i.routeModalAndKeys(msg); handled {
		return cmd
	}

	cmd, result := i.activeView.Update(msg)
	if result == nil {
		return cmd
	}

	if result.Type() == widgets.ResultCancel {
		cmd = i.handleCancelViewResult(cmd)
		return cmd
	}

	return tea.Batch(cmd, i.handleViewResult(result))
}

// View renders the current state of the intent, including any visible modal overlays.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (i *Intent) View() string {
	if i.activeView == nil {
		return "No active view"
	}
	view := i.CreateViewWithBreadcrumbs("Main Menu", "Browse Timeline", i.getStateName())
	view.WithContent(i.activeView.RenderContent())
	view.WithHelp(i.activeView.HelpText())
	baseView := view.Render()
	i.rebuildModalRegistry()
	return i.modalRegistry.RenderOverlay(baseView)
}

// renderTimelineView renders the timeline list view with modal overlays.

// Result provides the outcome of the intent after it becomes inactive,
//
// Returns:
//   - A fully initialized intents.IntentResult[interface{}] ready for use.
//
// Side effects:
//   - None.
func (i *Intent) Result() *intents.IntentResult[interface{}] {
	if i.result == nil {
		return nil
	}

	return &intents.IntentResult[interface{}]{
		Status:   i.result.Status,
		Data:     i.result.Data,
		Error:    i.result.Error,
		Metadata: i.result.Metadata,
	}
}

// setCancelled marks the intent as cancelled.
func (i *Intent) setCancelled() {
	i.result = &intents.IntentResult[*Result]{
		Status: intents.Cancelled,
	}
	i.active = false
}

// getContext returns a context for service calls.
//
// Currently returns context.Background() which does not support
// cancellation or timeouts. When BaseIntent gains a Context() method
// for lifecycle management, this should be updated to use that instead.
// For now, service calls using this context may not respond to intent
// teardown immediately.
func (i *Intent) getContext() context.Context {
	return context.Background()
}

// transitionToScreen sets the active screen and updates state.
type viewConfigurer interface {
	SetTerminalInfo(width, height int)
	SetTheme(theme interface{})
	SetLogo(logo interface{}, spacing int)
}

func (i *Intent) transitionToView(v widgets.View) {
	i.activeView = v
	if cfg, ok := v.(viewConfigurer); ok {
		termInfo := i.GetTerminalInfo()
		if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
			cfg.SetTerminalInfo(termInfo.Width, termInfo.Height)
		}
		if theme := i.Theme(); theme != nil {
			cfg.SetTheme(theme)
		}
		if logo := i.GetLogo(); logo != nil {
			cfg.SetLogo(logo, i.GetLogoSpacing())
		}
	}
}
