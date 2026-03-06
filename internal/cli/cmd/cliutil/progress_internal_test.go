package cliutil

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestSpinnerModelCtrlC verifies that spinnerModel.Update() correctly handles ctrl+c input.
// This is a direct unit test of the model, independent of Bubble Tea's event loop races.
func TestSpinnerModelCtrlC(t *testing.T) {
	resultChan := make(chan error, 1)
	m := newSpinnerModel("test", resultChan)

	// Send ctrl+c key message
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	sm := updated.(spinnerModel)
	if !sm.quitting {
		t.Error("expected quitting to be true after ctrl+c")
	}
	if cmd == nil {
		t.Error("expected quit command")
	}
}

// TestSpinnerModelIgnoresOtherKeys verifies that spinnerModel.Update() ignores non-ctrl+c keys.
func TestSpinnerModelIgnoresOtherKeys(t *testing.T) {
	resultChan := make(chan error, 1)
	m := newSpinnerModel("test", resultChan)

	// Send regular key message
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	sm := updated.(spinnerModel)
	if sm.quitting {
		t.Error("should not quit on regular key")
	}
	if cmd != nil {
		t.Error("should return nil cmd for unhandled key")
	}
}

// TestProgressModelCtrlC verifies that progressModel.Update() correctly handles ctrl+c input.
// This is a direct unit test of the model, independent of Bubble Tea's event loop races.
func TestProgressModelCtrlC(t *testing.T) {
	progressChan := make(chan int)
	resultChan := make(chan error, 1)
	m := newProgressModel("test", 10, progressChan, resultChan)

	// Send ctrl+c key message
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	pm := updated.(progressModel)
	if !pm.quitting {
		t.Error("expected quitting to be true after ctrl+c")
	}
	if cmd == nil {
		t.Error("expected quit command")
	}
}

// TestProgressModelIgnoresOtherKeys verifies that progressModel.Update() ignores non-ctrl+c keys.
func TestProgressModelIgnoresOtherKeys(t *testing.T) {
	progressChan := make(chan int)
	resultChan := make(chan error, 1)
	m := newProgressModel("test", 10, progressChan, resultChan)

	// Send regular key message
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

	pm := updated.(progressModel)
	if pm.quitting {
		t.Error("should not quit on regular key")
	}
	if cmd != nil {
		t.Error("should return nil cmd for unhandled key")
	}
}
