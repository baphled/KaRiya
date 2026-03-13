package modals

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	themes2 "github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// BurstDetailModal wraps feedback.DetailModal for viewing burst details.
// This uses the generic UIKit DetailModal with burst-specific content rendering.
//
// Usage:
//
//	modal := modals.NewBurstDetailModal(burst, theme)
//	modal.SetDimensions(width, height)
//	modal.Show()
//
//	// In Update:
//	model, cmd := modal.Update(msg)
//
//	// In View:
//	if modal.IsVisible() {
//	    return behaviors.RenderModalOverlay(modal, background)
//	}
type BurstDetailModal struct {
	modal *feedback.DetailModal
	burst *career.Burst
	theme themes.Theme
}

// NewBurstDetailModal creates a new burst detail modal.
//
// Expected:
//   - burst must be a non-nil *career.Burst pointer.
//   - theme must be a valid Theme instance (can be nil).
//
// Returns:
//   - A fully initialized BurstDetailModal ready for use.
//
// Side effects:
//   - Initializes footer badges.
func NewBurstDetailModal(burst *career.Burst, theme themes.Theme) *BurstDetailModal {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	content := renderBurstDetailContent(burst, theme)

	modal := feedback.NewDetailModal("Burst Details", content)
	if theme != nil {
		modal = modal.WithTheme(theme)
	}

	m := &BurstDetailModal{
		modal: modal,
		burst: burst,
		theme: theme,
	}

	// Initialize footer badges immediately.
	m.updateFooterBadges()

	return m
}

// updateFooterBadges updates the footer badges.
func (m *BurstDetailModal) updateFooterBadges() {
	theme := m.theme
	if theme == nil {
		theme = themes2.Default()
	}

	badges := []*primitives.Badge{
		primitives.ViewEventsBadge(theme),
		primitives.ViewFactsBadge(theme),
		primitives.ViewSkillsBadge(theme),
		primitives.EditBadge(theme),
		primitives.DeleteBadge(theme),
		primitives.ConfirmActionBadge(theme),
		primitives.InferSkillsBadge(theme),
		primitives.CloseBadge(theme),
	}

	m.modal = m.modal.WithFooterBadges(badges...)
}

// Init initializes the modal.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *BurstDetailModal) Init() tea.Cmd {
	m.updateFooterBadges()
	return m.modal.Init()
}

// Update handles keyboard input and window sizing.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Model: the updated model.
//   - tea.Cmd: command to execute.
//
// Side effects:
//   - Delegates to underlying modal's Update method.
func (m *BurstDetailModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.modal.Update(msg)
}

// View renders the modal content.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *BurstDetailModal) View() string {
	return m.modal.View()
}

// IsVisible returns whether the modal is currently visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *BurstDetailModal) IsVisible() bool {
	return m.modal.IsVisible()
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *BurstDetailModal) Show() {
	m.modal.Show()
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *BurstDetailModal) Hide() {
	m.modal.Hide()
}

// SetDimensions sets the terminal dimensions.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (m *BurstDetailModal) SetDimensions(width, height int) {
	m.modal.SetDimensions(width, height)
}

// SetBurst updates the burst being displayed.
//
// Expected:
//   - burst must be valid.
//
// Side effects:
//   - None.
func (m *BurstDetailModal) SetBurst(burst *career.Burst) {
	m.burst = burst
	content := renderBurstDetailContent(burst, m.theme)
	m.modal.SetContent(content)
}

// GetBurst returns the burst being displayed.
//
// Returns:
//   - A fully initialized career.Burst ready for use.
//
// Side effects:
//   - None.
func (m *BurstDetailModal) GetBurst() *career.Burst {
	return m.burst
}

// renderBurstDetailContent renders the burst details as formatted text.
func renderBurstDetailContent(burst *career.Burst, theme themes.Theme) string {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	// Guard against nil burst.
	if burst == nil {
		return primitives.NewText("No burst selected", theme).Render()
	}

	var b strings.Builder

	// Title with confirmation status indicator.
	if burst.Confirmed {
		confirmedBadge := primitives.SuccessText("✓ Confirmed", theme).Bold().Render()
		b.WriteString(confirmedBadge + "\n\n")
	}

	// Name.
	b.WriteString(primitives.NewText("Name:", theme).Bold().Render())
	b.WriteString(" " + burst.Name + "\n\n")

	// Description.
	if burst.Description != "" {
		b.WriteString(primitives.NewText("Description:", theme).Bold().Render())
		b.WriteString("\n" + burst.Description + "\n\n")
	}

	// Event count.
	b.WriteString(primitives.NewText("Events:", theme).Bold().Render())
	fmt.Fprintf(&b, " %d\n\n", len(burst.EventIDs))

	// Confirmation details.
	if burst.Confirmed && burst.ConfirmedAt != nil {
		b.WriteString(primitives.NewText("Confirmed:", theme).Bold().Render())
		b.WriteString(" " + burst.ConfirmedAt.Format("2006-01-02 15:04") + "\n\n")
	}

	// Created/Updated timestamps.
	b.WriteString(primitives.NewText("Created:", theme).Bold().Render())
	b.WriteString(" " + burst.CreatedAt.Format("2006-01-02") + "\n")

	b.WriteString(primitives.NewText("Updated:", theme).Bold().Render())
	b.WriteString(" " + burst.UpdatedAt.Format("2006-01-02") + "\n")

	return b.String()
}
