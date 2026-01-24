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
func (m *BurstFactsModal) Init() tea.Cmd {
	return m.modal.Init()
}

// Update handles keyboard input and window sizing.
func (m *BurstFactsModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.modal.Update(msg)
}

// View renders the modal content.
func (m *BurstFactsModal) View() string {
	return m.modal.View()
}

// IsVisible returns whether the modal is currently visible.
func (m *BurstFactsModal) IsVisible() bool {
	return m.modal.IsVisible()
}

// Show makes the modal visible.
func (m *BurstFactsModal) Show() {
	m.modal.Show()
}

// Hide hides the modal.
func (m *BurstFactsModal) Hide() {
	m.modal.Hide()
}

// SetDimensions sets the terminal dimensions.
func (m *BurstFactsModal) SetDimensions(width, height int) {
	m.modal.SetDimensions(width, height)
}

// SetFacts updates the facts being displayed.
func (m *BurstFactsModal) SetFacts(facts []*career.Fact) {
	m.facts = facts
	content := renderFactsContent(facts, m.theme)
	m.modal.SetContent(content)
	m.modal.SetTitle(fmt.Sprintf("Facts (%d)", len(facts)))
}

// GetBurstID returns the burst ID this modal is showing facts for.
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
			b.WriteString(primitives.NewText(fmt.Sprintf("   Categories: %s", strings.Join(fact.CompetencyCategories, ", ")), theme).Foreground(theme.SecondaryColor()).Render())
			b.WriteString("\n")
		}

		if fact.StrengthSignal != "" {
			b.WriteString(primitives.NewText(fmt.Sprintf("   Strength: %s", fact.StrengthSignal), theme).Foreground(theme.SecondaryColor()).Render())
			b.WriteString("\n")
		}

		if idx < len(facts)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}
