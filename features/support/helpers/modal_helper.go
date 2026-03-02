package helpers

import (
	"context"
	"fmt"
	"strings"

	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/testutil/harness"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// ModalHelper provides modal interaction abstraction.
type ModalHelper struct {
	env *harness.TestEnv
	ctx context.Context
}

// NewModalHelper creates a modal helper for the given context.
//
// Expected:
//   - ctx must contain a valid TestEnv set by the BDD environment.
//
// Returns:
//   - A ModalHelper bound to the test environment, or error if env is nil.
//
// Side effects:
//   - None.
func NewModalHelper(ctx context.Context) (*ModalHelper, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return nil, godog.ErrPending
	}
	return &ModalHelper{env: env, ctx: ctx}, nil
}

// OpenModalWithKey presses a key to open a modal (e.g., 'x' for export, 'n' for new).
//
// Expected:
//   - key must be a valid rune corresponding to a modal trigger key.
//
// Returns:
//   - nil always (error interface for consistency).
//
// Side effects:
//   - Sends a key event to the test environment.
func (m *ModalHelper) OpenModalWithKey(key rune) error {
	m.env.PressKeyRune(key)
	return nil
}

// CloseModal closes the modal with Escape.
//
// Returns:
//   - nil always (error interface for consistency).
//
// Side effects:
//   - Sends Escape key event to the test environment.
func (m *ModalHelper) CloseModal() error {
	m.env.PressKey(tea.KeyEscape)
	return nil
}

// ConfirmModal confirms the modal action (Enter or 'y').
//
// Returns:
//   - nil always (error interface for consistency).
//
// Side effects:
//   - Sends confirm key event to the test environment.
func (m *ModalHelper) ConfirmModal() error {
	m.env.Confirm()
	return nil
}

// IsModalVisible checks if a modal with the given title/text is visible.
//
// Expected:
//   - title must be a non-empty string.
//
// Returns:
//   - true if the title appears in the current view.
//
// Side effects:
//   - None.
func (m *ModalHelper) IsModalVisible(title string) bool {
	view := m.env.GetView()
	return strings.Contains(view, title)
}

// WaitForModalToAppear waits for a modal with the given text to appear.
//
// Expected:
//   - text must be a non-empty string to search for.
//
// Returns:
//   - nil if modal appears within timeout, or error from gomega assertion.
//
// Side effects:
//   - Polls the test environment view repeatedly.
func (m *ModalHelper) WaitForModalToAppear(text string) error {
	gomega.Eventually(func() string {
		return m.env.GetView()
	}, "5s", "100ms").Should(gomega.ContainSubstring(text))
	return nil
}

// WaitForModalToDisappear waits for modal text to disappear.
//
// Expected:
//   - text must be a non-empty string to search for.
//
// Returns:
//   - nil if modal disappears within timeout, or error from gomega assertion.
//
// Side effects:
//   - Polls the test environment view repeatedly.
func (m *ModalHelper) WaitForModalToDisappear(text string) error {
	gomega.Eventually(func() string {
		return m.env.GetView()
	}, "5s", "100ms").ShouldNot(gomega.ContainSubstring(text))
	return nil
}

// NavigateToPreview presses Enter to navigate to preview (in review screen).
//
// Returns:
//   - nil always (error interface for consistency).
//
// Side effects:
//   - Sends confirm key event to the test environment.
func (m *ModalHelper) NavigateToPreview() error {
	m.env.Confirm()
	return nil
}

// PressKeyToExport presses the specified key to open export modal.
//
// Expected:
//   - key must be a non-empty string; only the first character is used.
//
// Returns:
//   - nil always (error interface for consistency).
//
// Side effects:
//   - Sends a key event to the test environment if key is non-empty.
func (m *ModalHelper) PressKeyToExport(key string) error {
	if key != "" {
		m.env.PressKeyRune(rune(key[0]))
	}
	return nil
}

// AssertModalContains verifies the modal contains the specified text.
//
// Expected:
//   - text must be a non-empty string to search for.
//
// Returns:
//   - nil if text found, or error describing the missing text.
//
// Side effects:
//   - None.
func (m *ModalHelper) AssertModalContains(text string) error {
	view := m.env.GetView()
	if !strings.Contains(view, text) {
		return fmt.Errorf("modal does not contain '%s'", text)
	}
	return nil
}
