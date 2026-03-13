package burst

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	themes2 "github.com/baphled/kariya/internal/ui/uikit/theme"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
)

// Detail wraps feedback.DetailModal for viewing burst details.
// This uses the generic UIKit DetailModal with burst-specific content rendering.
//
// Usage:
//
//	modal := modals.NewDetail(burst, theme)
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
type Detail struct {
	*widgets.ModalDelegate
	burst display.Burst
	theme themes.Theme
}

// NewDetail creates a new burst detail modal.
//
// Expected:
//   - burst must be a non-nil *career.Burst pointer.
//   - theme must be a valid Theme instance (can be nil).
//
// Returns:
//   - A fully initialized Detail ready for use.
//
// Side effects:
//   - Initializes footer badges.
func NewDetail(burst display.Burst, theme themes.Theme) *Detail {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	content := renderDetailContent(burst, theme)

	modal := feedback.NewDetailModal("Burst Details", content)
	if theme != nil {
		modal = modal.WithTheme(theme)
	}

	m := &Detail{
		ModalDelegate: widgets.NewModalDelegate(modal),
		burst:         burst,
		theme:         theme,
	}

	// Initialize footer badges immediately.
	m.updateFooterBadges()

	return m
}

// updateFooterBadges updates the footer badges.
func (m *Detail) updateFooterBadges() {
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

	m.SetModal(m.Modal().WithFooterBadges(badges...))
}

// Init initializes the modal.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *Detail) Init() tea.Cmd {
	m.updateFooterBadges()
	return m.Modal().Init()
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
func (m *Detail) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.ModalUpdate(msg)
}

// View renders the modal content.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *Detail) View() string {
	return m.ModalView()
}

// SetBurst updates the burst being displayed.
//
// Expected:
//   - burst must be valid.
//
// Side effects:
//   - None.
func (m *Detail) SetBurst(burst display.Burst) {
	m.burst = burst
	content := renderDetailContent(burst, m.theme)
	m.SetContent(content)
}

// GetBurst returns the burst being displayed.
//
// Returns:
//   - A display.Burst value.
//
// Side effects:
//   - None.
func (m *Detail) GetBurst() display.Burst {
	return m.burst
}

// renderDetailContent renders the burst details as formatted text.
func renderDetailContent(burst display.Burst, theme themes.Theme) string {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	if isEmptyDisplayBurst(burst) {
		return primitives.NewText("No burst selected", theme).Render()
	}

	var b strings.Builder

	if burst.Confirmed {
		confirmedBadge := primitives.SuccessText("✓ Confirmed", theme).Bold().Render()
		b.WriteString(confirmedBadge + "\n\n")
	}

	b.WriteString(primitives.NewText("Name:", theme).Bold().Render())
	b.WriteString(" " + burst.Name + "\n\n")

	if burst.Description != "" {
		b.WriteString(primitives.NewText("Description:", theme).Bold().Render())
		b.WriteString("\n" + burst.Description + "\n\n")
	}

	b.WriteString(primitives.NewText("Events:", theme).Bold().Render())
	fmt.Fprintf(&b, " %d\n\n", len(burst.EventIDs))

	if burst.Confirmed && burst.ConfirmedAt != nil {
		b.WriteString(primitives.NewText("Confirmed:", theme).Bold().Render())
		b.WriteString(" " + burst.ConfirmedAt.Format("2006-01-02 15:04") + "\n\n")
	}

	b.WriteString(primitives.NewText("Created:", theme).Bold().Render())
	b.WriteString(" " + burst.CreatedAt.Format("2006-01-02") + "\n")

	b.WriteString(primitives.NewText("Updated:", theme).Bold().Render())
	b.WriteString(" " + burst.UpdatedAt.Format("2006-01-02") + "\n")

	return b.String()
}

func isEmptyDisplayBurst(burst display.Burst) bool {
	return burst.ID == "" &&
		burst.Name == "" &&
		burst.Description == "" &&
		len(burst.EventIDs) == 0 &&
		!burst.Confirmed &&
		burst.ConfirmedAt == nil &&
		burst.CreatedAt.IsZero() &&
		burst.UpdatedAt.IsZero()
}
