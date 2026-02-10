package helpers

import (
	"context"
	"fmt"
	"strings"

	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// ModalHelper provides modal interaction abstraction.
type ModalHelper struct {
	env *e2e.TestEnv
	ctx context.Context
}

// NewModalHelper creates a modal helper for the given context.
func NewModalHelper(ctx context.Context) (*ModalHelper, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return nil, godog.ErrPending
	}
	return &ModalHelper{env: env, ctx: ctx}, nil
}

// OpenModalWithKey presses a key to open a modal (e.g., 'x' for export, 'n' for new).
func (m *ModalHelper) OpenModalWithKey(key rune) error {
	m.env.PressKeyRune(key)
	return nil
}

// CloseModal closes the modal with Escape.
func (m *ModalHelper) CloseModal() error {
	m.env.PressKey(tea.KeyEscape)
	return nil
}

// ConfirmModal confirms the modal action (Enter or 'y').
func (m *ModalHelper) ConfirmModal() error {
	m.env.Confirm()
	return nil
}

// IsModalVisible checks if a modal with the given title/text is visible.
func (m *ModalHelper) IsModalVisible(title string) bool {
	view := m.env.GetView()
	return strings.Contains(view, title)
}

// WaitForModalToAppear waits for a modal with the given text to appear.
func (m *ModalHelper) WaitForModalToAppear(text string) error {
	gomega.Eventually(func() string {
		return m.env.GetView()
	}, "5s", "100ms").Should(gomega.ContainSubstring(text))
	return nil
}

// WaitForModalToDisappear waits for modal text to disappear.
func (m *ModalHelper) WaitForModalToDisappear(text string) error {
	gomega.Eventually(func() string {
		return m.env.GetView()
	}, "5s", "100ms").ShouldNot(gomega.ContainSubstring(text))
	return nil
}

// NavigateToPreview presses Enter to navigate to preview (in review screen).
func (m *ModalHelper) NavigateToPreview() error {
	m.env.Confirm()
	return nil
}

// PressKeyToExport presses the specified key to open export modal.
func (m *ModalHelper) PressKeyToExport(key string) error {
	if len(key) > 0 {
		m.env.PressKeyRune(rune(key[0]))
	}
	return nil
}

// AssertModalContains verifies the modal contains the specified text.
func (m *ModalHelper) AssertModalContains(text string) error {
	view := m.env.GetView()
	if !strings.Contains(view, text) {
		return fmt.Errorf("modal does not contain '%s'", text)
	}
	return nil
}
