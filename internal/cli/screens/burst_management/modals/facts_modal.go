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

// BurstFactsModal wraps feedback.DetailModal for viewing facts extracted from a burst.
type BurstFactsModal struct {
	modal   *feedback.DetailModal
	burstID string
	facts   []*career.Fact
	theme   themes.Theme
}

// NewBurstFactsModal creates a new burst facts modal.
//
// Expected:
//   - burstid must be a valid string.
//   - burstname must be a valid string.
//   - facts must be a valid slice of *career.Fact.
//   - theme must be a valid Theme instance (can be nil).
//
// Returns:
//   - A fully initialized BurstFactsModal ready for use.
//
// Side effects:
//   - None.
func NewBurstFactsModal(burstID string, burstName string, facts []*career.Fact, theme themes.Theme) *BurstFactsModal {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	content := renderFactsContent(facts, theme)
	title := fmt.Sprintf("Facts from Burst: %s (%d)", burstName, len(facts))

	modal := feedback.NewDetailModal(title, content)
	if theme != nil {
		modal = modal.WithTheme(theme)
	}

	return &BurstFactsModal{
		modal:   modal,
		burstID: burstID,
		facts:   facts,
		theme:   theme,
	}
}

// Init initializes the modal.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *BurstFactsModal) Init() tea.Cmd {
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
//   - Delegates to underlying modal.
func (m *BurstFactsModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.modal.Update(msg)
}

// View renders the modal content.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *BurstFactsModal) View() string {
	return m.modal.View()
}

// IsVisible returns whether the modal is currently visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *BurstFactsModal) IsVisible() bool {
	return m.modal.IsVisible()
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *BurstFactsModal) Show() {
	m.modal.Show()
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *BurstFactsModal) Hide() {
	m.modal.Hide()
}

// SetDimensions sets the terminal dimensions.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (m *BurstFactsModal) SetDimensions(width, height int) {
	m.modal.SetDimensions(width, height)
}

// SetFacts updates the facts being displayed.
//
// Expected:
//   - fact must be valid.
//
// Side effects:
//   - None.
func (m *BurstFactsModal) SetFacts(facts []*career.Fact) {
	m.facts = facts
	content := renderFactsContent(facts, m.theme)
	m.modal.SetContent(content)
	m.modal.SetTitle(fmt.Sprintf("Facts (%d)", len(facts)))
}

// GetBurstID returns the burst ID this modal is showing facts for.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *BurstFactsModal) GetBurstID() string {
	return m.burstID
}

// renderFactsContent renders the facts as formatted text.
func renderFactsContent(facts []*career.Fact, theme themes.Theme) string {
	if theme == nil {
		theme = themes2.Default()
	}

	if len(facts) == 0 {
		return primitives.WarningText("No facts extracted yet. Confirm the burst in detail view to extract facts.", theme).Render()
	}

	var b strings.Builder

	for idx, fact := range facts {
		b.WriteString(fmt.Sprintf("%d. %s\n", idx+1, fact.Text))

		if len(fact.CompetencyCategories) > 0 {
			catText := "   Categories: " + strings.Join(fact.CompetencyCategories, ", ")
			b.WriteString(primitives.NewText(catText, theme).
				Foreground(theme.SecondaryColor()).Render())
			b.WriteString("\n")
		}

		if fact.StrengthSignal != "" {
			strengthText := "   Strength: " + fact.StrengthSignal
			b.WriteString(primitives.NewText(strengthText, theme).
				Foreground(theme.SecondaryColor()).Render())
			b.WriteString("\n")
		}

		if idx < len(facts)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}
