package harness

import (
	"time"

	"github.com/baphled/kariya/internal/cli/app"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// SelectIntent navigates to and selects a menu item by its index (0-based).
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SelectIntent(index int) *TestEnv {
	e.T.Helper()

	// Ensure we're in menu state before attempting navigation
	if !e.IsInMenuState() {
		e.T.Fatalf("Cannot select intent: not in menu state. Current view:\n%s", e.GetView())
	}

	// Navigate to the menu item from position 0
	// Press 'g' to ensure we're at the top of the menu first
	e.PressKeyRune('g')

	// Navigate down to the desired index
	for range index {
		e.PressKeyRune('j')
	}

	// Select the intent
	e.PressKey(tea.KeyEnter)

	return e
}

// SelectIntentByName navigates to and selects a menu item by its intent name.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SelectIntentByName(name string) *TestEnv {
	e.T.Helper()

	intentOrder := map[string]int{
		"capture_event":    0,
		"browse_timeline":  1,
		"manage_skills":    2,
		"generate_cv":      3,
		"burst_management": 4,
		"fact_management":  5,
	}

	index, ok := intentOrder[name]
	if !ok {
		e.T.Fatalf("unknown intent name: %s", name)
	}

	return e.SelectIntent(index)
}

// PressKey sends a key message to the model.
//
// Expected:
//   - keytype must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) PressKey(key tea.KeyType) *TestEnv {
	e.T.Helper()

	modelInterface, cmd := e.Model.Update(tea.KeyMsg{Type: key})
	model, ok := modelInterface.(*app.Model)
	if !ok {
		e.T.Fatal("model type assertion failed: expected *app.Model")
	}
	e.Model = model

	// Execute any returned command
	e.executeCmd(cmd)

	return e
}

// PressKeyRune sends a rune key message to the model.
//
// Expected:
//   - rune must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) PressKeyRune(r rune) *TestEnv {
	e.T.Helper()

	modelInterface, cmd := e.Model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	model, ok := modelInterface.(*app.Model)
	if !ok {
		e.T.Fatal("model type assertion failed: expected *app.Model")
	}
	e.Model = model

	// Execute any returned command
	e.executeCmd(cmd)

	return e
}

// PressKeys sends multiple keys in sequence.
//
// Expected:
//   - interface{} must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) PressKeys(keys ...interface{}) *TestEnv {
	e.T.Helper()

	for _, key := range keys {
		switch k := key.(type) {
		case tea.KeyType:
			e.PressKey(k)
		case rune:
			e.PressKeyRune(k)
		case string:
			for _, r := range k {
				e.PressKeyRune(r)
			}
		default:
			e.T.Fatalf("unsupported key type: %T", key)
		}
	}

	return e
}

// TypeText types a string character by character.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) TypeText(text string) *TestEnv {
	e.T.Helper()

	for _, r := range text {
		e.PressKeyRune(r)
	}

	return e
}

// NavigateDown moves down in a list (j or down arrow).
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) NavigateDown() *TestEnv {
	return e.PressKeyRune('j')
}

// NavigateUp moves up in a list (k or up arrow).
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) NavigateUp() *TestEnv {
	return e.PressKeyRune('k')
}

// Confirm presses Enter to confirm an action.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) Confirm() *TestEnv {
	return e.PressKey(tea.KeyEnter)
}

// Cancel presses Escape to cancel/go back.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) Cancel() *TestEnv {
	return e.PressKey(tea.KeyEscape)
}

// GoBack presses Escape to go back.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) GoBack() *TestEnv {
	return e.Cancel()
}

// Quit presses 'q' to quit.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) Quit() *TestEnv {
	return e.PressKeyRune('q')
}

// Tab presses Tab to move to next field.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) Tab() *TestEnv {
	return e.PressKey(tea.KeyTab)
}

// SubmitHuhForm submits a huh form by pressing Enter.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SubmitHuhForm() *TestEnv {
	return e.Confirm()
}

// ClearTextField clears a text field by moving to end and pressing backspace.
// This is useful for clearing pre-populated huh form fields.
//
// Expected:
//   - maxChars should be a reasonable upper bound for the text length.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) ClearTextField(maxChars int) *TestEnv {
	e.T.Helper()

	e.PressKey(tea.KeyCtrlE)

	for range maxChars {
		e.PressKey(tea.KeyBackspace)
	}

	return e
}

// NextFormField moves to the next field in a huh form.
// This sends the huh.NextField message directly to properly navigate forms.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) NextFormField() *TestEnv {
	e.T.Helper()

	msg := huh.NextField()
	e.updateModelAndExecute(msg)

	return e
}

// executeCmd executes commands returned by Update, but only for specific message types
// that are essential for state transitions (like form submission).
//
// Most Bubble Tea commands (cursor blink, window resize) are ignored because they
// cause infinite loops or stuck goroutines in tests. We only care about messages
// that actually change application state.
//
// Commands that take longer than 5s to execute are skipped to avoid hanging tests.
// The 5s timeout accommodates database operations (including inference) while still
// filtering out stuck goroutines. Cursor blink ticks (530ms) are fast enough to pass.
func (e *TestEnv) executeCmd(cmd tea.Cmd) {
	if cmd == nil {
		return
	}

	// Execute command with timeout to skip stuck commands
	type result struct {
		msg tea.Msg
	}
	done := make(chan result, 1)
	go func() {
		done <- result{msg: cmd()}
	}()

	select {
	case r := <-done:
		if r.msg == nil {
			return
		}
		e.processCmdResult(r.msg)
	case <-time.After(5 * time.Second):
		// Command took too long - skip it to avoid hanging
		return
	}
}

// updateModelAndExecute updates the model with a message and executes any returned command.
// This is a helper to reduce cognitive complexity in processCmdResult.
//
// Expected:
//   - msg must be a valid tea.Msg.
//
// Returns:
//   - None.
//
// Side effects:
//   - Updates the Model field of the TestEnv.
//   - Recursively executes any returned command.
func (e *TestEnv) updateModelAndExecute(msg tea.Msg) {
	modelInterface, nextCmd := e.Model.Update(msg)
	model, ok := modelInterface.(*app.Model)
	if !ok {
		e.T.Fatal("model type assertion failed: expected *app.Model")
	}
	e.Model = model
	e.executeCmd(nextCmd)
}

// SendMessage sends a message directly to the model.
//
// Expected:
//   - msg must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SendMessage(msg tea.Msg) *TestEnv {
	e.T.Helper()

	modelInterface, cmd := e.Model.Update(msg)
	model, ok := modelInterface.(*app.Model)
	if !ok {
		e.T.Fatal("model type assertion failed: expected *app.Model")
	}
	e.Model = model
	e.executeCmd(cmd)

	return e
}

// InitModel initializes the model by calling Init() and sending a WindowSizeMsg.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) InitModel() *TestEnv {
	e.T.Helper()

	_ = e.Model.Init()

	// Send a WindowSizeMsg to trigger form layout
	e.SendMessage(tea.WindowSizeMsg{Width: TerminalWidth, Height: TerminalHeightShared})

	// Type and delete a character to force the form to render its fields
	// This workaround activates huh's internal rendering state
	e.PressKeyRune('x')
	e.PressKey(tea.KeyBackspace)

	return e
}

// PressEnterWithFormProcessing presses Enter and processes any internal form messages.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) PressEnterWithFormProcessing() *TestEnv {
	e.T.Helper()

	// Send Enter key
	modelInterface, cmd := e.Model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model, ok := modelInterface.(*app.Model)
	if !ok {
		e.T.Fatal("model type assertion failed: expected *app.Model")
	}
	e.Model = model

	// Process all resulting messages (including internal form messages)
	// This allows huh's group transitions to complete
	e.processFormCmds(cmd, 10)

	return e
}

// SendMessageWithFormProcessing sends a message and processes all resulting
// internal form messages (including huh init, focus, and group transitions).
//
// Expected:
//   - msg must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SendMessageWithFormProcessing(msg tea.Msg) *TestEnv {
	e.T.Helper()

	modelInterface, cmd := e.Model.Update(msg)
	model, ok := modelInterface.(*app.Model)
	if ok {
		e.Model = model
	}

	if cmd != nil {
		e.processFormCmds(cmd, 10)
	}

	return e
}

// PressKeyWithFormProcessing sends a key and processes all resulting internal form messages.
func (e *TestEnv) PressKeyWithFormProcessing(key tea.KeyType) *TestEnv {
	e.T.Helper()
	return e.SendMessageWithFormProcessing(tea.KeyMsg{Type: key})
}

// PressKeyRuneWithFormProcessing sends a rune key and processes all resulting
// internal form messages. Use this when a key press opens or initializes a
// huh form (e.g., pressing 'e' to open the metadata editor).
//
// Expected:
//   - r must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) PressKeyRuneWithFormProcessing(r rune) *TestEnv {
	e.T.Helper()

	return e.SendMessageWithFormProcessing(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
}

// TabWithFormProcessing sends a Tab key and processes all resulting internal
// form messages. Use this when navigating between fields in a huh form.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) TabWithFormProcessing() *TestEnv {
	e.T.Helper()

	return e.SendMessageWithFormProcessing(tea.KeyMsg{Type: tea.KeyTab})
}

// processFormCmds processes commands from form interactions.
// Unlike executeCmd, this processes ALL messages (including huh internals)
// but has a depth limit to prevent infinite loops. Commands that take longer
// than 600ms (cursor blink ticks) are skipped to avoid blocking.
func (e *TestEnv) processFormCmds(cmd tea.Cmd, maxDepth int) {
	if cmd == nil || maxDepth <= 0 {
		return
	}

	type result struct {
		msg tea.Msg
	}
	done := make(chan result, 1)
	go func() {
		done <- result{msg: cmd()}
	}()

	var msg tea.Msg
	select {
	case r := <-done:
		msg = r.msg
	case <-time.After(600 * time.Millisecond):
		return
	}

	if msg == nil {
		return
	}

	switch typedMsg := msg.(type) {
	case tea.BatchMsg:
		for _, batchCmd := range typedMsg {
			e.processFormCmds(batchCmd, maxDepth-1)
		}
	case nil:
		return
	default:
		modelInterface, nextCmd := e.Model.Update(msg)
		model, ok := modelInterface.(*app.Model)
		if !ok {
			e.T.Fatal("model type assertion failed: expected *app.Model")
		}
		e.Model = model
		e.processFormCmds(nextCmd, maxDepth-1)
	}
}

// CompleteOnboarding simulates completing the onboarding wizard.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) CompleteOnboarding(name, email string) *TestEnv {
	e.T.Helper()

	if !e.IsInOnboardingState() {
		return e
	}

	// Initialize the model first (required for huh forms to work)
	e.InitModel()

	// Step 1: Welcome + Name
	// Type the name
	e.TypeText(name)
	// Press Enter to advance to next step (uses form processing to handle internal huh messages)
	e.PressEnterWithFormProcessing()

	// Step 2: Contact - Email + Location
	// Type the email
	e.TypeText(email)
	// Press Enter to accept email and move to Location field
	e.PressEnterWithFormProcessing()
	// Press Enter again to accept empty Location and advance to Step 3
	e.PressEnterWithFormProcessing()

	// Step 3: Professional Details - 3 optional fields (Title, GitHub, Portfolio)
	// Press Enter 3 times to accept all empty fields and complete
	e.PressEnterWithFormProcessing()
	e.PressEnterWithFormProcessing()
	e.PressEnterWithFormProcessing()

	return e
}
