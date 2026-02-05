package generatecv

import (
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	tea "github.com/charmbracelet/bubbletea"
)

// NewIntent constructs a GenerateCV intent initialised with the given
// context, selecting a default profile when one is available.
//
// Expected:
//   - ctx must pass Validate (non-nil, at least one profile and one event).
//
// Returns:
//   - A ready-to-activate intent and nil error on success.
//   - Nil intent and a validation error when context is invalid.
//
// Side effects:
//   - None.
func NewIntent(ctx *IntentContext) (*Intent, error) {
	if err := ctx.Validate(); err != nil {
		return nil, err
	}

	selectedProfile := ctx.DefaultProfile
	if selectedProfile == nil && len(ctx.AvailableProfiles) > 0 {
		selectedProfile = ctx.AvailableProfiles[0]
	}

	baseIntent := intents.NewBaseIntent()

	return &Intent{
		BaseIntent: baseIntent,
		context:    ctx,
		state: &model{
			context:         ctx,
			currentState:    StateConfiguring,
			selectedProfile: selectedProfile,
		},
		active: true,
	}, nil
}

// Init prepares the intent for its first render cycle by creating the
// configuration wizard modal and showing it.
//
// Returns:
//   - A tea.Cmd that initialises the wizard modal.
//
// Side effects:
//   - Creates and shows the wizard modal.
func (i *Intent) Init() tea.Cmd {
	if i.context.DefaultProfile != nil {
		i.state.selectedProfile = i.context.DefaultProfile
	}

	return i.initWizardFlow()
}

// Update advances the intent state machine by processing a single Bubble Tea
// message, delegating to the wizard flow or active screen handlers.
//
// Expected:
//   - msg must be a valid tea.Msg (key press, window resize, or async result).
//
// Returns:
//   - A tea.Cmd for follow-up work (async generation, quit, etc.), or nil.
//
// Side effects:
//   - Mutates internal state, active screens, and modal visibility.
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	return i.updateWizardFlow(msg)
}

// View produces the terminal UI string for the intent's current state.
//
// Returns:
//   - A rendered string for the terminal.
//
// Side effects:
//   - None.
func (i *Intent) View() string {
	if !i.active {
		return "GenerateCV intent is not active"
	}

	return i.wizardView()
}

// Result retrieves the outcome of the intent after it becomes inactive.
//
// Returns:
//   - A fully initialized IntentResult, or nil if the intent is still active.
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

// setCompleted marks the intent as completed with success.
func (i *Intent) setCompleted() {
	i.result = &intents.IntentResult[*Result]{
		Status: intents.Completed,
		Data: &Result{
			GeneratedCV:     i.state.generatedCV,
			SelectedProfile: i.state.selectedProfile,
			AcceptedFields:  make(map[string]bool),
		},
		Metadata: map[string]interface{}{
			"profile":     i.state.selectedProfile.ID,
			"audience":    i.state.selectedAudience,
			"timestamp":   time.Now(),
			"event_count": len(i.context.Events),
			"fact_count":  len(i.context.Facts),
		},
	}
	i.active = false
}

// setCancelled marks the intent as cancelled by the user.
func (i *Intent) setCancelled() {
	i.result = &intents.IntentResult[*Result]{
		Status: intents.Cancelled,
	}
	i.active = false
}
